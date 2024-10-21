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

package ext

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

var bindingTests = []struct {
	expr      string
	parseOnly bool
}{
	{expr: `cel.bind(a, 'hell' + 'o' + '!', [a, a, a].join(', ')) ==
	        ['hell' + 'o' + '!', 'hell' + 'o' + '!', 'hell' + 'o' + '!'].join(', ')`},
	// Variable shadowing
	{expr: `cel.bind(a,
		        cel.bind(a, 'world', a + '!'),
		   		'hello ' + a) == 'hello ' + 'world' + '!'`},
}

func TestBindings(t *testing.T) {
	env, err := cel.NewEnv(Bindings(BindingsVersion(0)), Strings())
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings(), Strings()) failed: %v", err)
	}
	for i, tst := range bindingTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			if !tc.parseOnly {
				cAst, iss := env.Check(pAst)
				if iss.Err() != nil {
					t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
				}
				asts = append(asts, cAst)
			}
			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatal(err)
				}
				out, _, err := prg.Eval(cel.NoVars())
				if err != nil {
					t.Fatal(err)
				} else if out.Value() != true {
					t.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}

func TestBindingsNonMatch(t *testing.T) {
	env, err := cel.NewEnv(Bindings(), Strings())
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings(), Strings()) failed: %v", err)
	}
	nonMatchExpr := `ceel.bind(a, 1, a)`
	ast, iss := env.Parse(nonMatchExpr)
	if iss.Err() != nil {
		t.Fatalf("env.Parse(%v) failed: %v", nonMatchExpr, iss.Err())
	}
	if len(ast.SourceInfo().GetMacroCalls()) != 0 {
		t.Fatalf("env.Parse(%v) performed a macro replacement when none was expected: %v",
			nonMatchExpr, ast.SourceInfo().GetMacroCalls())
	}
}

func TestBindingsInvalidIdent(t *testing.T) {
	env, err := cel.NewEnv(Bindings(), Strings())
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings(), Strings()) failed: %v", err)
	}
	invalidIdentExpr := `cel.bind(a.b, 1, a.b)`
	wantErr := "ERROR: <input>:1:11: cel.bind() variable names must be simple identifiers"
	_, iss := env.Parse(invalidIdentExpr)
	if !strings.Contains(iss.Err().Error(), wantErr) {
		t.Fatalf("env.Parse(%v) failed: %v", invalidIdentExpr, iss.Err())
	}
}

func BenchmarkBindings(b *testing.B) {
	env, err := cel.NewEnv(Bindings(), Strings())
	if err != nil {
		b.Fatalf("cel.NewEnv(Bindings(), Strings()) failed: %v", err)
	}
	for i, tst := range bindingTests {
		tc := tst
		ast, iss := env.Compile(tc.expr)
		if iss.Err() != nil {
			b.Fatalf("env.Compile(%q) failed: %v", tc.expr, iss.Err())
		}
		prg, err := env.Program(ast, cel.EvalOptions(cel.OptOptimize))
		if err != nil {
			b.Fatalf("env.Program(ast, Optimize) failed: %v", err)
		}
		// Benchmark the eval.
		b.Run(fmt.Sprintf("[%d]", i), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				prg.Eval(cel.NoVars())
			}
		})
	}
}

func TestBlockEval(t *testing.T) {
	fac := ast.NewExprFactory()
	blockExpr := fac.NewCall(
		1, "cel.@block",
		fac.NewList(2, []ast.Expr{
			fac.NewIdent(3, "x"),
			fac.NewIdent(4, "@index0"),
			fac.NewIdent(5, "@index1"),
		}, []int32{}),
		fac.NewCall(9, operators.Add,
			fac.NewCall(6, operators.Add,
				fac.NewIdent(7, "@index2"),
				fac.NewIdent(8, "@index1")),
			fac.NewIdent(10, "@index0"),
		),
	)
	blockAST := ast.NewAST(blockExpr, nil)
	env, err := cel.NewEnv(
		Bindings(),
		cel.Variable("x", cel.StringType),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings()) failed: %v", err)
	}
	prg, err := env.PlanProgram(blockAST)
	if err != nil {
		t.Fatalf("PlanProgram() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": "hello"})
	if err != nil {
		t.Fatalf("prg.Eval() failed: %v", err)
	}
	if out.Equal(types.String("hellohellohello")) != types.True {
		t.Errorf("got %v, wanted 'hellohelloehello'", out)
	}
}

func TestBlockEval_BadPlan(t *testing.T) {
	fac := ast.NewExprFactory()
	blockExpr := fac.NewCall(
		1, "cel.@block",
		fac.NewList(2, []ast.Expr{
			fac.NewIdent(3, "x"),
			fac.NewIdent(4, "@index0"),
		}, []int32{}),
		fac.NewCall(6, operators.Add,
			fac.NewIdent(7, "@index1"),
			fac.NewIdent(8, "@index0")),
		fac.NewIdent(9, "x"),
	)
	blockAST := ast.NewAST(blockExpr, nil)
	env, err := cel.NewEnv(
		Bindings(BindingsVersion(1)),
		cel.Variable("x", cel.StringType),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings()) failed: %v", err)
	}
	_, err = env.PlanProgram(blockAST)
	if err == nil {
		t.Fatal("PlanProgram() succeeded, expected error")
	}
}

func TestBlockEval_BadBlock(t *testing.T) {
	fac := ast.NewExprFactory()
	blockExpr := fac.NewCall(
		1, "cel.@block",
		fac.NewCall(2, operators.Add,
			fac.NewIdent(3, "@index1"),
			fac.NewIdent(4, "@index0")),
		fac.NewIdent(5, "x"),
	)
	blockAST := ast.NewAST(blockExpr, nil)
	env, err := cel.NewEnv(
		Bindings(BindingsVersion(1)),
		cel.Variable("x", cel.StringType),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings()) failed: %v", err)
	}
	_, err = env.PlanProgram(blockAST)
	if err == nil {
		t.Fatal("PlanProgram() succeeded, expected error")
	}
}

func TestBlockEval_BadIndex(t *testing.T) {
	fac := ast.NewExprFactory()
	blockExpr := fac.NewCall(
		1, "cel.@block",
		fac.NewList(2, []ast.Expr{
			fac.NewIdent(3, "x"),
			fac.NewIdent(4, "@indexNext"),
		}, []int32{}),
		fac.NewCall(6, operators.Add,
			fac.NewIdent(7, "@indexNext"),
			fac.NewIdent(8, "@index0")),
	)
	blockAST := ast.NewAST(blockExpr, nil)
	env, err := cel.NewEnv(
		Bindings(BindingsVersion(1)),
		cel.Variable("x", cel.StringType),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv(Bindings()) failed: %v", err)
	}
	prg, err := env.PlanProgram(blockAST)
	if err != nil {
		t.Fatalf("PlanProgram() failed: %v", err)
	}
	_, _, err = prg.Eval(map[string]any{"x": "hello"})
	if !strings.Contains(err.Error(), "no such attribute") {
		t.Fatalf("prg.Eval() got %v, expected no such attribute error", err)
	}
}
