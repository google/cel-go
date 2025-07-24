// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package cel

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"

	proto3pb "github.com/google/cel-go/test/proto3pb"
)

// testAnnotationFactory is a simple factory for testing.
type testAnnotationFactory struct {
	name       string
	value      any
	IsExpr     bool
	err        error
	applicable func(expr ast.Expr) bool
}

func (f *testAnnotationFactory) GenerateAnnotation(expr ast.Expr, a *ast.AST) (*Annotation, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.applicable == nil || f.applicable(expr) {
		return &Annotation{Name: f.name, Value: f.value, IsExpr: f.IsExpr}, nil
	}
	return nil, nil
}

type annotationTest struct {
	name        string
	expr        string
	factories   []AnnotationFactory
	expectError string
	want        string
	customCheck func(t *testing.T, a *ast.AST)
}

func TestAnnotationOptimizer(t *testing.T) {
	testValue := "test"
	tests := []annotationTest{
		{
			name: "no annotation factories",
			expr: "1 + 2",
			want: `1 + 2`,
		},
		{
			name: "annotation factory error",
			expr: "1",
			factories: []AnnotationFactory{
				&testAnnotationFactory{err: fmt.Errorf("factory failed")},
			},
			expectError: "error generating annotation: factory failed",
		},
		{
			name: "multiple annotations",
			expr: "1",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann1", value: testValue},
				&testAnnotationFactory{name: "ann2", value: testValue},
			},
			want: `cel.@annotation(1, [{"is_expr": false, "name": "ann1", "value": "test"}, {"is_expr": false, "name": "ann2", "value": "test"}])`,
		},
		{
			name: "annotation with nil value",
			expr: "false",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "nil_value_ann", value: nil, IsExpr: false},
			},
			want: `cel.@annotation(false, [{"is_expr": false, "name": "nil_value_ann", "value": null}])`,
		},
		{
			name: "multiple factories, some returning nil annotation",
			expr: "true && false",
			factories: []AnnotationFactory{
				&testAnnotationFactory{
					name:  "true_ann",
					value: testValue,
					applicable: func(e ast.Expr) bool {
						return e.Kind() == ast.LiteralKind && e.AsLiteral().Value() == true
					},
				},
				&testAnnotationFactory{
					name:  "false_ann",
					value: testValue,
					applicable: func(e ast.Expr) bool {
						return e.Kind() == ast.LiteralKind && e.AsLiteral().Value() == false
					},
				},
				&testAnnotationFactory{
					name:       "no_apply_ann",
					value:      testValue,
					applicable: func(e ast.Expr) bool { return false }, // Always returns nil
				},
			},
			want: `cel.@annotation(true, [{"is_expr": false, "name": "true_ann", "value": "test"}]) &&
cel.@annotation(false, [{"is_expr": false, "name": "false_ann", "value": "test"}])`,
		},
		{
			name: "Annotating call expression",
			expr: "size('hello')",
			factories: []AnnotationFactory{
				&testAnnotationFactory{
					name:       "call_ann",
					value:      testValue,
					applicable: func(e ast.Expr) bool { return e.Kind() == ast.CallKind },
				},
				&testAnnotationFactory{
					name:       "literal_ann",
					value:      testValue,
					applicable: func(e ast.Expr) bool { return e.Kind() == ast.LiteralKind },
				},
			},
			want: `cel.@annotation(size(cel.@annotation("hello", [{"is_expr": false, "name": "literal_ann", "value": "test"}])), [{"is_expr": false, "name": "call_ann", "value": "test"}])`,
		},
		{
			name: "Annotating member call",
			expr: "'hello'.size()",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation(cel.@annotation("hello", [{"is_expr": false, "name": "ann", "value": "test"}]).size(), [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating list",
			expr: "[1, 2]",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation([cel.@annotation(1, [{"is_expr": false, "name": "ann", "value": "test"}]), cel.@annotation(2, [{"is_expr": false, "name": "ann", "value": "test"}])], [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating map",
			expr: "{'a': 1, 'b': 2}",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation({cel.@annotation("a", [{"is_expr": false, "name": "ann", "value": "test"}]): cel.@annotation(1, [{"is_expr": false, "name": "ann", "value": "test"}]), cel.@annotation("b", [{"is_expr": false, "name": "ann", "value": "test"}]): cel.@annotation(2, [{"is_expr": false, "name": "ann", "value": "test"}])}, [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating map with optional entry",
			expr: "{'a': 1, ?'b': optional.of(2)}",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation({cel.@annotation("a", [{"is_expr": false, "name": "ann", "value": "test"}]): cel.@annotation(1, [{"is_expr": false, "name": "ann", "value": "test"}]), ?cel.@annotation("b", [{"is_expr": false, "name": "ann", "value": "test"}]): cel.@annotation(optional.of(cel.@annotation(2, [{"is_expr": false, "name": "ann", "value": "test"}])), [{"is_expr": false, "name": "ann", "value": "test"}])}, [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating struct",
			expr: "google.expr.proto3.test.TestAllTypes{single_int32: 1, single_string: 's'}",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation(google.expr.proto3.test.TestAllTypes{single_int32: cel.@annotation(1, [{"is_expr": false, "name": "ann", "value": "test"}]), single_string: cel.@annotation("s", [{"is_expr": false, "name": "ann", "value": "test"}])}, [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating struct with optional field",
			expr: "google.expr.proto3.test.TestAllTypes{single_string: 's', ?single_int32: optional.of(1)}",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation(google.expr.proto3.test.TestAllTypes{single_string: cel.@annotation("s", [{"is_expr": false, "name": "ann", "value": "test"}]), ?single_int32: cel.@annotation(optional.of(cel.@annotation(1, [{"is_expr": false, "name": "ann", "value": "test"}])), [{"is_expr": false, "name": "ann", "value": "test"}])}, [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating select",
			expr: "msg.single_int32",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation(cel.@annotation(msg, [{"is_expr": false, "name": "ann", "value": "test"}]).single_int32, [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating presence test",
			expr: "has(msg.single_int32)",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation(has(cel.@annotation(msg, [{"is_expr": false, "name": "ann", "value": "test"}]).single_int32), [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "Annotating comprehension",
			expr: "[1].map(x, x * 2)",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			customCheck: func(t *testing.T, a *ast.AST) {
				t.Helper()
				call := a.Expr().AsCall()
				if call == nil {
					t.Fatalf("expected call expression, got %v", a.Expr().Kind())
				}
				if call.FunctionName() != "cel.@annotation" {
					t.Errorf("expected cel.@annotation call, got %s", call.FunctionName())
				}
				compExpr := call.Args()[0]
				if compExpr.Kind() != ast.ComprehensionKind {
					t.Fatalf("expected comprehension, got %v", compExpr.Kind())
				}
				comp := compExpr.AsComprehension()
				sourceInfo := a.SourceInfo()
				checkPart := func(name string, expr ast.Expr, want string) {
					t.Helper()
					partAst := ast.NewAST(expr, sourceInfo)
					got, err := AstToString(&Ast{impl: partAst})
					if err != nil {
						t.Fatalf("AstToString(%s) failed: %v", name, err)
					}
					if got != want {
						t.Errorf("%s got\n%s\nwant\n%s", name, got, want)
					}
				}
				checkPart("iterRange", comp.IterRange(), `cel.@annotation([cel.@annotation(1, [{"is_expr": false, "name": "ann", "value": "test"}])], [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("accuInit", comp.AccuInit(), `cel.@annotation([], [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("loopCondition", comp.LoopCondition(), `cel.@annotation(true, [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("loopStep", comp.LoopStep(), `cel.@annotation(cel.@annotation(@result, [{"is_expr": false, "name": "ann", "value": "test"}]) + cel.@annotation([cel.@annotation(cel.@annotation(x, [{"is_expr": false, "name": "ann", "value": "test"}]) * cel.@annotation(2, [{"is_expr": false, "name": "ann", "value": "test"}]), [{"is_expr": false, "name": "ann", "value": "test"}])], [{"is_expr": false, "name": "ann", "value": "test"}]), [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("result", comp.Result(), `cel.@annotation(@result, [{"is_expr": false, "name": "ann", "value": "test"}])`)
			},
		},
		{
			name: "Annotating map comprehension with 2 variables",
			expr: "{'a': 'b'}.exists(k, k == v)",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			customCheck: func(t *testing.T, a *ast.AST) {
				t.Helper()
				call := a.Expr().AsCall()
				if call == nil {
					t.Fatalf("expected call expression, got %v", a.Expr().Kind())
				}
				if call.FunctionName() != "cel.@annotation" {
					t.Errorf("expected cel.@annotation call, got %s", call.FunctionName())
				}
				compExpr := call.Args()[0]
				if compExpr.Kind() != ast.ComprehensionKind {
					t.Fatalf("expected comprehension, got %v", compExpr.Kind())
				}
				comp := compExpr.AsComprehension()
				sourceInfo := a.SourceInfo()
				checkPart := func(name string, expr ast.Expr, want string) {
					t.Helper()
					partAst := ast.NewAST(expr, sourceInfo)
					got, err := AstToString(&Ast{impl: partAst})
					if err != nil {
						t.Fatalf("AstToString(%s) failed: %v", name, err)
					}
					if got != want {
						t.Errorf("%s got\n%s\nwant\n%s", name, got, want)
					}
				}
				checkPart("iterRange", comp.IterRange(), `cel.@annotation({cel.@annotation("a", [{"is_expr": false, "name": "ann", "value": "test"}]): cel.@annotation("b", [{"is_expr": false, "name": "ann", "value": "test"}])}, [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("accuInit", comp.AccuInit(), `cel.@annotation(false, [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("loopCondition", comp.LoopCondition(), `cel.@annotation(@not_strictly_false(cel.@annotation(!(cel.@annotation(@result, [{"is_expr": false, "name": "ann", "value": "test"}])), [{"is_expr": false, "name": "ann", "value": "test"}])), [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("loopStep", comp.LoopStep(), `cel.@annotation(cel.@annotation(@result, [{"is_expr": false, "name": "ann", "value": "test"}]) ||
cel.@annotation(cel.@annotation(k, [{"is_expr": false, "name": "ann", "value": "test"}]) == cel.@annotation(v, [{"is_expr": false, "name": "ann", "value": "test"}]), [{"is_expr": false, "name": "ann", "value": "test"}]), [{"is_expr": false, "name": "ann", "value": "test"}])`)
				checkPart("result", comp.Result(), `cel.@annotation(@result, [{"is_expr": false, "name": "ann", "value": "test"}])`)
			},
		},
		{
			name: "Annotating simple literal",
			expr: "true",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "literal_ann", value: testValue},
			},
			want: `cel.@annotation(true, [{"is_expr": false, "name": "literal_ann", "value": "test"}])`,
		},
		{
			name: "Annotating ident",
			expr: "x",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "ann", value: testValue},
			},
			want: `cel.@annotation(x, [{"is_expr": false, "name": "ann", "value": "test"}])`,
		},
		{
			name: "ternary operator",
			expr: "true ? 'yes' : 'no'",
			factories: []AnnotationFactory{
				&testAnnotationFactory{name: "tern_ann", value: testValue},
			},
			want: `cel.@annotation((cel.@annotation(true, [{"is_expr": false, "name": "tern_ann", "value": "test"}])) ? (cel.@annotation("yes", [{"is_expr": false, "name": "tern_ann", "value": "test"}])) : (cel.@annotation("no", [{"is_expr": false, "name": "tern_ann", "value": "test"}])), [{"is_expr": false, "name": "tern_ann", "value": "test"}])`,
		},
	}

	e, err := NewEnv(
		OptionalTypes(),
		EnableMacroCallTracking(),
		Types(&proto3pb.TestAllTypes{}),
		Variable("x", DynType),
		Variable("v", types.StringType),
		Variable("msg", types.NewObjectType("google.expr.proto3.test.TestAllTypes")),
	)
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			checked, iss := e.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("Compile(%q) failed: %v", tc.expr, iss.Err())
			}

			var factories []AnnotationFactory
			if tc.factories != nil {
				factories = tc.factories
			}

			optimizer, err := NewAnnotationOptimizer(AnnotationFactories(factories), AnnotationEnv(e))
			if err != nil {
				t.Fatalf("NewAnnotationOptimizer() failed: %v", err)
			}

			opt := NewStaticOptimizer(optimizer)
			optimized, iss := opt.Optimize(e, checked)

			if tc.customCheck != nil {
				tc.customCheck(t, optimized.NativeRep())
				return
			}

			if tc.expectError != "" {
				if iss.Err() == nil {
					t.Fatalf("Optimize() succeeded, wanted error: %v", tc.expectError)
				}
				if !strings.Contains(iss.Err().Error(), tc.expectError) {
					t.Errorf("Optimize() got error %v, wanted error containing %v", iss.Err(), tc.expectError)
				}
				return
			}

			if iss.Err() != nil {
				t.Fatalf("Optimize() failed: %v", iss.Err())
			}

			if tc.want == "" {
				t.Fatalf("tc.want must be set for test case %q", tc.name)
			}
			got, err := AstToString(optimized)
			if err != nil {
				t.Fatalf("ast.AstToString() failed: %v", err)
			}
			if got != tc.want {
				t.Errorf("ast.AstToString() got\n%s\nwant\n%s", got, tc.want)
			}
		})
	}
}
