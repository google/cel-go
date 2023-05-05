// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interpreter

import (
	"testing"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/parser"
	"github.com/google/cel-go/test"
)

type testInfo struct {
	in   any
	expr string
	out  string
}

var testCases = []testInfo{
	{
		in: map[string]any{
			"msg": map[string]string{"foo": "bar"},
		},
		expr: `msg`,
		out:  `{"foo": "bar"}`,
	},
	{
		expr: `true && false`,
		out:  `false`,
	},
	{
		in:   unknownActivation("x"),
		expr: `(true || false) && x`,
		out:  `x`,
	},
	{
		in:   unknownActivation("x"),
		expr: `(false || false) && x`,
		out:  `false`,
	},
	{
		in:   unknownActivation("a"),
		expr: `a && [1, 1u, 1.0].exists(x, type(x) == uint)`,
		out:  `a`,
	},
	{
		in:   unknownActivation("this"),
		expr: `this in []`,
		out:  `false`,
	},
	{
		in:   unknownActivation("this"),
		expr: `this in {}`,
		out:  `false`,
	},
	{
		in:   partialActivation(map[string]any{"rules": []string{}}, "this"),
		expr: `this in rules`,
		out:  `false`,
	},
	{
		in:   partialActivation(map[string]any{"rules": map[string]any{"not_in": []string{}}}, "this"),
		expr: `this.size() > 0 ? this in rules.not_in : !(this in rules.not_in)`,
		out:  `(this.size() > 0) ? false : true`,
	},
	{
		in: partialActivation(map[string]any{"rules": map[string]any{"not_in": []string{}}}, "this"),
		expr: `this.size() > 0 ? this in rules.not_in : 
			  !(this in rules.not_in) ? true : false`,
		out: `(this.size() > 0) ? false : true`,
	},
	{
		expr: `{'hello': 'world'.size()}`,
		out:  `{"hello": 5}`,
	},
	{
		expr: `[b'bytes-string']`,
		out:  `[b"\142\171\164\145\163\055\163\164\162\151\156\147"]`,
	},
	{
		expr: `[b'bytes'] + [b'-' + b'string']`,
		out:  `[b"\142\171\164\145\163", b"\055\163\164\162\151\156\147"]`,
	},
	{
		expr: `1u + 3u`,
		out:  `4u`,
	},
	{
		expr: `2 < 3`,
		out:  `true`,
	},
	{
		in:   unknownActivation(),
		expr: `test == null`,
		out:  `test == null`,
	},
	{
		in:   unknownActivation(),
		expr: `test == null || true`,
		out:  `true`,
	},
	{
		in:   unknownActivation(),
		expr: `test == null && false`,
		out:  `false`,
	},
	{
		in:   unknownActivation("b", "c"),
		expr: `true ? b < 1.2 : c == ['hello']`,
		out:  `b < 1.2`,
	},
	{
		in:   unknownActivation("b", "c"),
		expr: `false ? b < 1.2 : c == ['hello']`,
		out:  `c == ["hello"]`,
	},
	{
		in:   unknownActivation(),
		expr: `[1+3, 2+2, 3+1, four]`,
		out:  `[4, 4, 4, four]`,
	},
	{
		in:   unknownActivation(),
		expr: `test == {'a': 1, 'field': 2}.field`,
		out:  `test == 2`,
	},
	{
		in:   unknownActivation(),
		expr: `test in {'a': 1, 'field': [2, 3]}.field`,
		out:  `test in [2, 3]`,
	},
	{
		in:   unknownActivation(),
		expr: `test == {'field': [1 + 2, 2 + 3]}`,
		out:  `test == {"field": [3, 5]}`,
	},
	{
		in:   unknownActivation(),
		expr: `test in {'a': 1, 'field': [test, 3]}.field`,
		out:  `test in {"a": 1, "field": [test, 3]}.field`,
	},
	{
		in:   partialActivation(map[string]any{"foo": "bar"}, "r.attr"),
		expr: `foo == "bar" && r.attr.loc in ["GB", "US"]`,
		out:  `r.attr.loc in ["GB", "US"]`,
	},
	// TODO: the output of an expression like this relies on either
	// a) doing replacements on the original macro call, or
	// b) mutating the macro call tracking data rather than the core
	//    expression in order to render the partial correctly.
	// {
	// 	in:   unknownActivation(),
	// 	expr: `[1+3, 2+2, 3+1, four].exists(x, x == four)`,
	// 	out:  `[4, 4, 4, four].exists(x, x == four)`,
	// },
}

func TestPrune(t *testing.T) {
	p, err := parser.NewParser(
		parser.PopulateMacroCalls(true),
		parser.Macros(parser.AllMacros...),
	)
	if err != nil {
		t.Fatalf("parser.NewParser() failed: %v", err)
	}
	for i, tst := range testCases {
		ast, iss := p.Parse(common.NewStringSource(tst.expr, "<input>"))
		if len(iss.GetErrors()) > 0 {
			t.Fatalf(iss.ToDisplayString())
		}
		state := NewEvalState()
		reg := newTestRegistry(t)
		attrs := NewPartialAttributeFactory(containers.DefaultContainer, reg, reg)
		interp := NewStandardInterpreter(containers.DefaultContainer, reg, reg, attrs)

		interpretable, _ := interp.NewUncheckedInterpretable(
			ast.GetExpr(),
			ExhaustiveEval(), Observe(EvalStateObserver(state)))
		interpretable.Eval(testActivation(t, tst.in))
		newExpr := PruneAst(ast.GetExpr(), ast.GetSourceInfo().GetMacroCalls(), state)
		actual, err := parser.Unparse(newExpr.GetExpr(), newExpr.GetSourceInfo())
		if err != nil {
			t.Error(err)
		}
		if !test.Compare(actual, tst.out) {
			t.Errorf("prune[%d], diff: %s", i, test.DiffMessage("structure", actual, tst.out))
		}
	}
}

func unknownActivation(vars ...string) PartialActivation {
	pats := make([]*AttributePattern, len(vars), len(vars))
	for i, v := range vars {
		pats[i] = NewAttributePattern(v)
	}
	a, _ := NewPartialActivation(map[string]any{}, pats...)
	return a
}

func partialActivation(in map[string]any, vars ...string) PartialActivation {
	pats := make([]*AttributePattern, len(vars), len(vars))
	for i, v := range vars {
		pats[i] = NewAttributePattern(v)
	}
	a, _ := NewPartialActivation(in, pats...)
	return a
}

func testActivation(t *testing.T, in any) Activation {
	t.Helper()
	if in == nil {
		return EmptyActivation()
	}
	a, err := NewActivation(in)
	if err != nil {
		t.Fatalf("NewActivation(%v) failed: %v", in, err)
	}
	return a
}
