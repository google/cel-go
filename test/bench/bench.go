// Copyright 2023 Google LLC
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

// Package bench defines a structure for benchmarked tests against custom CEL environments.
package bench

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
)

// Case represents a human-readable expression and an expected output given an input
type Case struct {
	// Expr is a human-readable expression which is expected to compile.
	Expr string

	// Options indicate additional pieces of configuration such as CEL libraries, variables, and functions.
	Options []cel.EnvOption

	// In is expected to be a map[string]any or interpreter.Activation instance representing the input to the expression.
	In any

	// Out is the expected CEL valued output.
	Out ref.Val
}

var (
	// ReferenceCases represent canonical CEL expressions for common use cases.
	ReferenceCases = []*Case{
		{
			Expr: `string_value == 'value'`,
			Options: []cel.EnvOption{
				cel.Variable("string_value", cel.StringType),
			},
			In: map[string]any{
				"string_value": "value",
			},
			Out: types.True,
		},
		{
			Expr: `string_value != 'value'`,
			Options: []cel.EnvOption{
				cel.Variable("string_value", cel.StringType),
			},
			In: map[string]any{
				"string_value": "value",
			},
			Out: types.False,
		},
		{
			Expr: `'value' in list_value`,
			Options: []cel.EnvOption{
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"list_value": []string{"a", "b", "c", "value"},
			},
			Out: types.True,
		},
		{
			Expr: `!('value' in list_value)`,
			Options: []cel.EnvOption{
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"list_value": []string{"a", "b", "c", "d"},
			},
			Out: types.True,
		},
		{
			Expr: `x in ['a', 'b', 'c', 'd']`,
			Options: []cel.EnvOption{
				cel.Variable("x", cel.StringType),
			},
			In: map[string]any{
				"x": "c",
			},
			Out: types.True,
		},
		{
			Expr: `!(x in ['a', 'b', 'c', 'd'])`,
			Options: []cel.EnvOption{
				cel.Variable("x", cel.StringType),
			},
			In: map[string]any{
				"x": "e",
			},
			Out: types.True,
		},
		{
			Expr: `x in list_value`,
			Options: []cel.EnvOption{
				cel.Variable("x", cel.StringType),
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"x":          "c",
				"list_value": []string{"a", "b", "c", "d"},
			},
			Out: types.True,
		},
		{
			Expr: `!(x in list_value)`,
			Options: []cel.EnvOption{
				cel.Variable("x", cel.StringType),
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"x":          "e",
				"list_value": []string{"a", "b", "c", "d"},
			},
			Out: types.True,
		},
		{
			Expr: `list_value.exists(e, e.contains('cd'))`,
			Options: []cel.EnvOption{
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"list_value": []string{"abc", "bcd", "cde", "def"},
			},
			Out: types.True,
		},
		{
			Expr: `list_value.exists(e, e.startsWith('cd'))`,
			Options: []cel.EnvOption{
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"list_value": []string{"abc", "bcd", "cde", "def"},
			},
			Out: types.True,
		},
		{
			Expr: `list_value.exists(e, e.matches('cd*'))`,
			Options: []cel.EnvOption{
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"list_value": []string{"abc", "bcd", "cde", "def"},
			},
			Out: types.True,
		},
		{
			Expr: `list_value.filter(e, e.matches('^cd+')) == ['cde']`,
			Options: []cel.EnvOption{
				cel.Variable("list_value", cel.ListType(cel.StringType)),
			},
			In: map[string]any{
				"list_value": []string{"abc", "bcd", "cde", "def"},
			},
			Out: types.True,
		},
		{
			Expr: `'formatted list: %s, size: %d'.format([['abc', 'cde'], 2])`,
			Options: []cel.EnvOption{
				ext.Strings(),
			},
			In:  map[string]any{},
			Out: types.String(`formatted list: ["abc", "cde"], size: 2`),
		},
	}
)

// RunReferenceCases evaluates the set of ReferenceCases against a custom CEL environment.
//
// See: bench_test.go for an example.
func RunReferenceCases(b *testing.B, env *cel.Env) {
	b.Helper()
	for _, rc := range ReferenceCases {
		RunCase(b, env, rc)
	}
}

// RunCase evaluates a single test case against a custom environment, running three different
// variants of the expression: optimized, unoptimized, and trace.
//
// * `optimized` - applies the cel.EvalOptions(cel.OptOptimize) flag.
// * `unoptimized` - no optimization flags applied.
// * `trace` - observes the evaluation state of an expression.
//
// In many cases the evaluation times may be similar, but when running comparisons against the
// baseline CEL environment, it may be useful to characterize the performance of the custom
// environment against the baseline.
func RunCase(b *testing.B, env *cel.Env, bc *Case) {
	b.Helper()
	var err error
	if len(bc.Options) > 0 {
		env, err = env.Extend(bc.Options...)
		if err != nil {
			b.Fatalf("env.Extend() failed: %v", err)
		}
	}
	ast, iss := env.Compile(bc.Expr)
	if iss.Err() != nil {
		b.Fatalf("env.Compile(%v) failed: %v", bc.Expr, iss.Err())
	}
	opts := map[string][]cel.ProgramOption{
		"optimized":   {cel.EvalOptions(cel.OptOptimize)},
		"unoptimized": {},
		"trace":       {cel.EvalOptions(cel.OptTrackState)},
	}
	optOrder := []string{"optimized", "unoptimized", "trace"}
	for _, name := range optOrder {
		opt := opts[name]
		b.Run(fmt.Sprintf("%s/%s", bc.Expr, name), func(b *testing.B) {
			prg, err := env.Program(ast, opt...)
			if err != nil {
				b.Fatalf("env.Program(%v) failed: %v", bc.Expr, err)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				var input any = cel.NoVars()
				if bc.In != nil {
					input = bc.In
				}
				out, _, err := prg.Eval(input)
				if err != nil {
					b.Fatalf("prg.Eval(%v) failed: %v", input, err)
				}
				if out.Equal(bc.Out) != types.True {
					b.Fatalf("prg.Eval(%v) got %v, wanted %v", input, out, bc.Out)
				}
			}
		})
	}
}
