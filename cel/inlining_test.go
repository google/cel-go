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

package cel_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"

	"github.com/google/go-cmp/cmp"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"

	proto3pb "github.com/google/cel-go/test/proto3pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestInliningOptimizer(t *testing.T) {
	type varExpr struct {
		name  string
		alias string
		t     *cel.Type
		expr  string
	}
	tests := []struct {
		expr    string
		vars    []varExpr
		inlined string
		folded  string
	}{
		{
			expr: `a || b`,
			vars: []varExpr{
				{
					name: "a",
					t:    cel.BoolType,
				},
				{
					name:  "b",
					alias: "bravo",
					t:     cel.BoolType,
					expr:  `'hello'.contains('lo')`,
				},
			},
			inlined: `a || "hello".contains("lo")`,
			folded:  `true`,
		},
		{
			expr: `a + [a]`,
			vars: []varExpr{
				{
					name:  "a",
					alias: "alpha",
					t:     cel.DynType,
					expr:  `dyn([1, 2])`,
				},
			},
			inlined: `cel.bind(alpha, dyn([1, 2]), alpha + [alpha])`,
			folded:  `[1, 2, [1, 2]]`,
		},
		{
			expr: `a && (a || b)`,
			vars: []varExpr{
				{
					name:  "a",
					alias: "alpha",
					t:     cel.BoolType,
					expr:  `'hello'.contains('lo')`,
				},
				{
					name: "b",
					t:    cel.BoolType,
				},
			},
			inlined: `cel.bind(alpha, "hello".contains("lo"), alpha && (alpha || b))`,
			folded:  `true`,
		},
		{
			expr: `a && b && a`,
			vars: []varExpr{
				{
					name:  "a",
					alias: "alpha",
					t:     cel.BoolType,
					expr:  `'hello'.contains('lo')`,
				},
				{
					name: "b",
					t:    cel.BoolType,
				},
			},
			inlined: `cel.bind(alpha, "hello".contains("lo"), alpha && b && alpha)`,
			folded:  `cel.bind(alpha, true, alpha && b && alpha)`,
		},
		{
			expr: `(c || d) || (a && (a || b))`,
			vars: []varExpr{
				{
					name:  "a",
					alias: "alpha",
					t:     cel.BoolType,
					expr:  `'hello'.contains('lo')`,
				},
				{
					name: "b",
					t:    cel.BoolType,
				},
				{
					name: "c",
					t:    cel.BoolType,
				},
				{
					name: "d",
					t:    cel.BoolType,
					expr: "!false",
				},
			},
			inlined: `c || !false || cel.bind(alpha, "hello".contains("lo"), alpha && (alpha || b))`,
			folded:  `true`,
		},
		{
			expr: `a && (a || b)`,
			vars: []varExpr{
				{
					name: "a",
					t:    cel.BoolType,
				},
				{
					name:  "b",
					alias: "bravo",
					t:     cel.BoolType,
					expr:  `'hello'.contains('lo')`,
				},
			},
			inlined: `a && (a || "hello".contains("lo"))`,
			folded:  `a`,
		},
		{
			expr: `a && b`,
			vars: []varExpr{
				{
					name:  "a",
					alias: "alpha",
					t:     cel.BoolType,
					expr:  `!'hello'.contains('lo')`,
				},
				{
					name:  "b",
					alias: "bravo",
					t:     cel.BoolType,
				},
			},
			inlined: `!"hello".contains("lo") && b`,
			folded:  `false`,
		},
		{
			expr: `operation.system.consumers + operation.destination_consumers`,
			vars: []varExpr{
				{
					name: "operation.system",
					t:    cel.DynType,
				},
				{
					name: "operation.destination_consumers",
					t:    cel.ListType(cel.IntType),
					expr: `productsToConsumers(operation.destination_products)`,
				},
				{
					name: "operation.destination_products",
					t:    cel.ListType(cel.IntType),
					expr: `operation.system.products`,
				},
			},
			inlined: `operation.system.consumers + productsToConsumers(operation.system.products)`,
			folded:  `operation.system.consumers + productsToConsumers(operation.system.products)`,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.expr, func(t *testing.T) {
			opts := []cel.EnvOption{cel.OptionalTypes(),
				cel.EnableMacroCallTracking(),
				cel.Function("productsToConsumers",
					cel.Overload("productsToConsumers_list",
						[]*cel.Type{cel.ListType(cel.IntType)},
						cel.ListType(cel.IntType)))}

			varDecls := make([]cel.EnvOption, len(tc.vars))
			for i, v := range tc.vars {
				varDecls[i] = cel.Variable(v.name, v.t)
			}
			e, err := cel.NewEnv(append(varDecls, opts...)...)
			if err != nil {
				t.Fatalf("NewEnv() failed: %v", err)
			}
			inlinedVars := []*cel.InlineVariable{}
			for _, v := range tc.vars {
				if v.expr == "" {
					continue
				}
				checked, iss := e.Compile(v.expr)
				if iss.Err() != nil {
					t.Fatalf("Compile(%q) failed: %v", v.expr, iss.Err())
				}
				if v.alias == "" {
					inlinedVars = append(inlinedVars, cel.NewInlineVariable(v.name, checked))
				} else {
					inlinedVars = append(inlinedVars, cel.NewInlineVariableWithAlias(v.name, v.alias, checked))
				}
			}
			checked, iss := e.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("Compile() failed: %v", iss.Err())
			}

			opt := cel.NewStaticOptimizer(cel.NewInliningOptimizer(inlinedVars...))
			optimized, iss := opt.Optimize(e, checked)
			if iss.Err() != nil {
				t.Fatalf("Optimize() generated an invalid AST: %v", iss.Err())
			}
			inlined, err := cel.AstToString(optimized)
			if err != nil {
				t.Fatalf("cel.AstToString() failed: %v", err)
			}
			if inlined != tc.inlined {
				t.Errorf("got %q, wanted %q", inlined, tc.inlined)
			}
			folder, err := cel.NewConstantFoldingOptimizer()
			if err != nil {
				t.Fatalf("NewConstantFoldingOptimizer() failed: %v", err)
			}
			opt = cel.NewStaticOptimizer(folder)
			optimized, iss = opt.Optimize(e, optimized)
			if iss.Err() != nil {
				t.Fatalf("Optimize() generated an invalid AST: %v", iss.Err())
			}
			folded, err := cel.AstToString(optimized)
			if err != nil {
				t.Fatalf("cel.AstToString() failed: %v", err)
			}
			if folded != tc.folded {
				t.Errorf("got %q, wanted %q", folded, tc.folded)
			}
		})
	}
}

func TestInliningOptimizerMultiStage(t *testing.T) {
	type varDecl struct {
		name string
		t    *cel.Type
	}
	type inlineVarExpr struct {
		name  string
		alias string
		t     *cel.Type
		expr  string
	}
	tests := []struct {
		expr       string
		vars       []varDecl
		inlineVars []inlineVarExpr
		inlined    string
		folded     string
	}{
		{
			expr: `has(a.b)`,
			vars: []varDecl{
				{
					name: "a",
					t:    cel.MapType(cel.StringType, cel.StringType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "a.b",
					alias: "alpha",
					t:     cel.StringType,
					expr:  `a.b_long`,
				},
			},
			inlined: `has(a.b_long)`,
			folded:  `has(a.b_long)`,
		},
		{
			expr: `has(a.b) ? a.b : 'default'`,
			vars: []varDecl{
				{
					name: "a",
					t:    cel.MapType(cel.StringType, cel.StringType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "a.b",
					alias: "alpha",
					t:     cel.StringType,
					expr:  `'hello'`,
				},
			},
			inlined: `cel.bind(alpha, "hello", (alpha.size() != 0) ? alpha : "default")`,
			folded:  `"hello"`,
		},
		{
			expr: `has(a.b) ? a.b : ['default']`,
			vars: []varDecl{
				{
					name: "a",
					t:    cel.MapType(cel.StringType, cel.ListType(cel.StringType)),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "a.b",
					alias: "alpha",
					t:     cel.StringType,
					expr:  `['hello']`,
				},
			},
			inlined: `cel.bind(alpha, ["hello"], (alpha.size() != 0) ? alpha : ["default"])`,
			folded:  `["hello"]`,
		},
		{
			expr: `0 in msg.map_int64_nested_type`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "nested_map",
					t:    cel.MapType(cel.IntType, cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes")),
				},
			},
			inlineVars: []inlineVarExpr{

				{
					name: "msg.map_int64_nested_type",
					t:    cel.MapType(cel.IntType, cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes")),
					expr: `nested_map`,
				},
			},
			inlined: `0 in nested_map`,
			folded:  `0 in nested_map`,
		},
		{
			expr: `has(msg.single_any)`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "unpacked_wrapper",
					t:    cel.NullableType(cel.StringType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name: "msg.single_any",
					t:    cel.NullableType(cel.StringType),
					expr: `unpacked_wrapper`,
				},
			},
			inlined: `unpacked_wrapper != null`,
			folded:  `unpacked_wrapper != null`,
		},
		{
			expr: `has(msg.single_any) ? msg.single_any : '10'`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "unpacked_wrapper",
					t:    cel.NullableType(cel.StringType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_any",
					t:     cel.NullableType(cel.StringType),
					alias: "wrapped",
					expr:  `unpacked_wrapper`,
				},
			},
			inlined: `cel.bind(wrapped, unpacked_wrapper, (wrapped != null) ? wrapped : "10")`,
			folded:  `cel.bind(wrapped, unpacked_wrapper, (wrapped != null) ? wrapped : "10")`,
		},
		{
			expr: `has(msg.child.payload.single_int32_wrapper)`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes"),
				},
				{
					name: "unpacked_child",
					t:    cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes"),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.child.payload",
					t:     cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes"),
					alias: "payload",
					expr:  `unpacked_child.payload`,
				},
			},
			inlined: `has(unpacked_child.payload.single_int32_wrapper)`,
			folded:  `has(unpacked_child.payload.single_int32_wrapper)`,
		},
		{
			expr: `has(msg.child.payload.single_int32_wrapper)`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes"),
				},
				{
					name: "unpacked_payload",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.child.payload.single_int32_wrapper",
					t:     cel.NullableType(cel.IntType),
					alias: "payload",
					expr:  `unpacked_payload.single_int32_wrapper`,
				},
			},
			inlined: `has(unpacked_payload.single_int32_wrapper)`,
			folded:  `has(unpacked_payload.single_int32_wrapper)`,
		},
		{
			expr: `has(msg.child.payload.single_int32_wrapper) ? msg.child.payload.single_int32_wrapper : 1`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.NestedTestAllTypes"),
				},
				{
					name: "unpacked_payload",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.child.payload.single_int32_wrapper",
					t:     cel.NullableType(cel.IntType),
					alias: "nullable_int",
					expr:  `unpacked_payload.single_int32_wrapper`,
				},
			},
			inlined: `cel.bind(nullable_int, unpacked_payload.single_int32_wrapper, (nullable_int != null) ? nullable_int : 1)`,
			folded:  `cel.bind(nullable_int, unpacked_payload.single_int32_wrapper, (nullable_int != null) ? nullable_int : 1)`,
		},
		{
			expr: `has(msg.single_value) ? msg.single_value : null`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_value",
					t:     cel.NullableType(cel.DoubleType),
					alias: "nullable_float",
					expr:  `dyn(1.5)`,
				},
			},
			inlined: `cel.bind(nullable_float, dyn(1.5), (nullable_float != null) ? nullable_float : null)`,
			folded:  `1.5`,
		},
		{
			expr: `has(msg.single_any) ? msg.single_any : 42`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_any",
					t:     cel.IntType,
					alias: "unpacked_nested",
					expr:  `google.expr.proto3.test.NestedTestAllTypes{}.payload.single_int32`,
				},
			},
			inlined: `has(google.expr.proto3.test.NestedTestAllTypes{}.payload.single_int32) ? google.expr.proto3.test.NestedTestAllTypes{}.payload.single_int32 : 42`,
			folded:  `42`,
		},
		{
			expr: `has(msg.single_any.processing_purpose)`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "unpacked_purpose",
					t:    cel.ListType(cel.IntType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_any.processing_purpose",
					t:     cel.ListType(cel.IntType),
					alias: "unpacked_purpose",
					expr:  `[1, 2, 3].map(i, i * 2)`,
				},
			},
			inlined: `[1, 2, 3].map(i, i * 2).size() != 0`,
			folded:  `true`,
		},
		{
			expr: `has(msg.single_any.processing_purpose) ? msg.single_any.processing_purpose[0] : 42`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "unpacked_purpose",
					t:    cel.ListType(cel.IntType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_any.processing_purpose",
					t:     cel.ListType(cel.IntType),
					alias: "unpacked_purpose",
					expr:  `[1, 2, 3].map(i, i * 2)`,
				},
			},
			inlined: `cel.bind(unpacked_purpose, [1, 2, 3].map(i, i * 2), (unpacked_purpose.size() != 0) ? (unpacked_purpose[0]) : 42)`,
			folded:  `2`,
		},
		{
			expr: `has(msg.single_any.processing_purpose) ? msg.single_any.processing_purpose.map(i, i * 2)[0] : 42`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "unpacked_purpose",
					t:    cel.ListType(cel.IntType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_any.processing_purpose",
					t:     cel.ListType(cel.IntType),
					alias: "unpacked_purpose",
					expr:  `[1, 2, 3].map(i, i * 2)`,
				},
			},
			inlined: `cel.bind(unpacked_purpose, [1, 2, 3].map(i, i * 2), (unpacked_purpose.size() != 0) ? (unpacked_purpose.map(i, i * 2)[0]) : 42)`,
			folded:  `4`,
		},
		{
			expr: `msg.single_any.processing_purpose.filter(j, 
					j < msg.single_any.processing_purpose.size()) == [2]`,
			vars: []varDecl{
				{
					name: "msg",
					t:    cel.ObjectType("google.expr.proto3.test.TestAllTypes"),
				},
				{
					name: "unpacked_purpose",
					t:    cel.ListType(cel.IntType),
				},
			},
			inlineVars: []inlineVarExpr{
				{
					name:  "msg.single_any.processing_purpose",
					t:     cel.ListType(cel.IntType),
					alias: "unpacked_purpose",
					expr:  `[1, 2, 3].map(i, i * 2)`,
				},
			},
			inlined: `cel.bind(unpacked_purpose, [1, 2, 3].map(i, i * 2), unpacked_purpose.filter(j, j < unpacked_purpose.size())) == [2]`,
			folded:  `true`,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.expr, func(t *testing.T) {
			opts := []cel.EnvOption{
				cel.Container("google.expr"),
				cel.Types(&proto3pb.TestAllTypes{}),
				cel.OptionalTypes(),
				cel.EnableMacroCallTracking()}

			varDecls := make([]cel.EnvOption, len(tc.vars))
			for i, v := range tc.vars {
				varDecls[i] = cel.Variable(v.name, v.t)
			}
			e, err := cel.NewEnv(append(varDecls, opts...)...)
			if err != nil {
				t.Fatalf("NewEnv() failed: %v", err)
			}
			inlinedVars := []*cel.InlineVariable{}
			for _, v := range tc.inlineVars {
				if v.expr == "" {
					continue
				}
				checked, iss := e.Compile(v.expr)
				if iss.Err() != nil {
					t.Fatalf("Compile(%q) failed: %v", v.expr, iss.Err())
				}
				if v.alias == "" {
					inlinedVars = append(inlinedVars, cel.NewInlineVariable(v.name, checked))
				} else {
					inlinedVars = append(inlinedVars, cel.NewInlineVariableWithAlias(v.name, v.alias, checked))
				}
			}
			checked, iss := e.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("Compile() failed: %v", iss.Err())
			}

			opt := cel.NewStaticOptimizer(cel.NewInliningOptimizer(inlinedVars...))
			optimized, iss := opt.Optimize(e, checked)
			if iss.Err() != nil {
				t.Fatalf("Optimize() generated an invalid AST: %v", iss.Err())
			}
			inlined, err := cel.AstToString(optimized)
			if err != nil {
				t.Fatalf("cel.AstToString() failed: %v", err)
			}
			if inlined != tc.inlined {
				t.Errorf("got %q, wanted %q", inlined, tc.inlined)
			}
			folder, err := cel.NewConstantFoldingOptimizer()
			if err != nil {
				t.Fatalf("NewConstantFoldingOptimizer() failed: %v", err)
			}
			opt = cel.NewStaticOptimizer(folder)
			optimized, iss = opt.Optimize(e, optimized)
			if iss.Err() != nil {
				t.Fatalf("Optimize() generated an invalid AST: %v", iss.Err())
			}
			folded, err := cel.AstToString(optimized)
			if err != nil {
				t.Fatalf("cel.AstToString() failed: %v", err)
			}
			if folded != tc.folded {
				t.Errorf("got %q, wanted %q", folded, tc.folded)
			}
		})
	}
}

func TestSanitizeExpr(t *testing.T) {
	type varDecl struct {
		name string
		t    *cel.Type
	}
	tests := []struct {
		expr      string
		vars      []varDecl
		sanitized string
	}{
		{
			expr: `has(a.b)`,
			vars: []varDecl{
				{
					name: `a`,
					t:    cel.MapType(cel.StringType, cel.DynType),
				},
			},

			sanitized: `
			id: 4
			`,
		},
		{
			expr: `a.size() > 0 && has(a.b)`,
			vars: []varDecl{
				{
					name: `a`,
					t:    cel.MapType(cel.StringType, cel.DynType),
				},
			},

			sanitized: `
			id: 9
			call_expr: {
				function:"_&&_",
				args:[{
					id:3,
					call_expr:{
						function:"_>_",
						args:[{
							id:2,
							call_expr:{
								target:{id:1,
									ident_expr:{name:"a"}
								},
								function:"size"
							}
						},
						{
							id:4,
							const_expr:{int64_value:0}
						}]
					}
				},
				{
					id:8
				}]}`,
		},
		{
			expr: `a.exists(k, has(a[k].b))`,
			vars: []varDecl{
				{
					name: `a`,
					t:    cel.MapType(cel.StringType, cel.DynType),
				},
			},

			sanitized: `
			id: 17
			`,
		},
		{
			expr: `a.exists(k, has(a[k].b)) && has(a.b)`,
			vars: []varDecl{
				{
					name: `a`,
					t:    cel.MapType(cel.StringType, cel.DynType),
				},
			},

			sanitized: `
			id: 22
			call_expr: {
				function:"_&&_"
				args:[{
					id: 17
				}, {
					id: 21
				}]
			}
			`,
		},
	}

	for _, tst := range tests {
		tc := tst
		t.Run(tc.expr, func(t *testing.T) {

			opts := []cel.EnvOption{
				cel.Container("google.expr"),
				cel.Types(&proto3pb.TestAllTypes{}),
				cel.OptionalTypes(),
				cel.EnableMacroCallTracking()}
			varDecls := make([]cel.EnvOption, len(tc.vars))
			for i, v := range tc.vars {
				varDecls[i] = cel.Variable(v.name, v.t)
			}
			e, err := cel.NewEnv(append(varDecls, opts...)...)
			if err != nil {
				t.Fatalf("NewEnv() failed: %v", err)
			}
			checked, iss := e.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("Compile() failed: %v", iss.Err())
			}
			expr := checked.NativeRep().Expr()
			info := checked.NativeRep().SourceInfo()
			fac := ast.NewExprFactory()
			macroExpr := fac.CopyExpr(expr)
			ast.PostOrderVisit(macroExpr, ast.NewExprVisitor(func(e ast.Expr) {
				if _, exists := info.GetMacroCall(e.ID()); exists {
					e.SetKindCase(nil)
				}
			}))

			macroPBExpr, err := ast.ExprToProto(macroExpr)
			if err != nil {
				t.Fatalf("ast.ExprToProto() failed: %v", err)
			}
			expectedExpr := &exprpb.Expr{}
			if err := prototext.Unmarshal([]byte(tc.sanitized), expectedExpr); err != nil {
				t.Fatalf("prototext.Unmarshal() failed: %v", err)
			}
			if diff := cmp.Diff(expectedExpr, macroPBExpr, protocmp.Transform()); diff != "" {
				t.Errorf("sanitize %s returned unexpected diff (-want +got):\n%s\n%v", tc.expr, diff, prototext.Format(macroPBExpr))
			}
		})
	}
}

func TestUpdateExpr(t *testing.T) {
	fac := ast.NewExprFactory()
	tests := []struct {
		name        string
		expr        ast.Expr
		origInfo    *ast.SourceInfo
		inlined     ast.Expr
		want        ast.Expr
		updatedInfo *ast.SourceInfo
	}{
		{
			name: "rewrite presence test",
			expr: fac.NewPresenceTest(1, fac.NewIdent(2, "a"), "b"),
			origInfo: macrosInfo(
				map[int64]ast.Expr{
					1: fac.NewSelect(0, fac.NewIdent(2, "a"), "b"),
				},
			),
			inlined: fac.NewPresenceTest(1, fac.NewIdent(2, "a"), "b_long"),
			want:    fac.NewPresenceTest(1, fac.NewIdent(2, "a"), "b_long"),
			updatedInfo: macrosInfo(
				map[int64]ast.Expr{
					1: fac.NewSelect(0, fac.NewIdent(2, "a"), "b_long"),
				},
			),
		},
		{
			name: "rewrite presence test",
			expr: fac.NewSelect(1, fac.NewIdent(2, "a"), "b"),
			origInfo: macrosInfo(
				map[int64]ast.Expr{
					1: fac.NewMemberCall(0, "exists",
						fac.NewList(2, []ast.Expr{fac.NewIdent(3, "a")}, []int32{}),
						fac.NewIdent(4, "i"),
						fac.NewUnspecifiedExpr(10)),
					10: fac.NewPresenceTest(0, fac.NewIdent(11, "i"), "j"),
				},
			),
			inlined: fac.NewComprehension(1,
				fac.NewList(2, []ast.Expr{fac.NewIdent(3, "a")}, []int32{}),
				"i",
				"__result__",
				fac.NewLiteral(6, types.True),
				fac.NewLiteral(7, types.True),
				fac.NewCall(8,
					"_&&_",
					fac.NewAccuIdent(9),
					fac.NewPresenceTest(10, fac.NewIdent(11, "i"), "j"),
				),
				fac.NewAccuIdent(12),
			),
			want: fac.NewPresenceTest(1, fac.NewIdent(2, "a"), "b_long"),
			updatedInfo: macrosInfo(
				map[int64]ast.Expr{
					1: fac.NewSelect(0, fac.NewIdent(2, "a"), "b_long"),
				},
			),
		},
	}

	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			target := fac.CopyExpr(tc.expr)
			targetInfo := ast.CopySourceInfo(tc.origInfo)
			cel.UpdateExpr(fac, target, tc.inlined, targetInfo)
			targetPBExpr := pbExpr(t, target)
			expectedPBExpr := pbExpr(t, tc.want)
			if diff := cmp.Diff(expectedPBExpr, targetPBExpr, protocmp.Transform()); diff != "" {
				t.Errorf("UpdateExpr() returned unexpected expr diff (-want +got):\n%s\n%v", diff, prototext.Format(targetPBExpr))
			}
			targetPBInfo := pbInfo(t, targetInfo)
			expectedPBInfo := pbInfo(t, tc.updatedInfo)
			if diff := cmp.Diff(expectedPBInfo, targetPBInfo, protocmp.Transform()); diff != "" {
				t.Errorf("UpdateExpr() returned unexpected info diff (-want +got):\n%s\n%v", diff, prototext.Format(targetPBInfo))
			}
		})
	}
}

func pbExpr(t *testing.T, e ast.Expr) *exprpb.Expr {
	t.Helper()
	pb, err := ast.ExprToProto(e)
	if err != nil {
		t.Fatalf("ast.ExprToProto() failed: %v", err)
	}
	return pb
}

func pbInfo(t *testing.T, info *ast.SourceInfo) *exprpb.SourceInfo {
	t.Helper()
	pb, err := ast.SourceInfoToProto(info)
	if err != nil {
		t.Fatalf("ast.SourceInfoToProto() failed: %v", err)
	}
	return pb
}

func macrosInfo(macros map[int64]ast.Expr) *ast.SourceInfo {
	info := ast.NewSourceInfo(nil)
	for id, call := range macros {
		info.SetMacroCall(id, call)
	}
	return info
}
