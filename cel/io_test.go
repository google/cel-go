// Copyright 2019 Google LLC
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

package cel

import (
	"testing"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/test/proto3pb"
	"google.golang.org/protobuf/proto"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	structpb "google.golang.org/types/protobuf/types/known/structpb"
)

func TestValueToRefValue(t *testing.T) {
	stdEnv, err := NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}

	expected := structpb.NewNullValue()
	val, err := RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err := ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	expected = types.True
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	expected = types.Int(0)
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	expected = types.Uint(0)
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	expected = types.Double(0.0)
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	expected = types.Bytes(make([]byte, 0, 5))
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	expected = types.String("abc")
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	reg, err = types.NewRegistry()
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	expected = types.NewDynamicList(reg, []bool{true})
	expected.Add(types.String("abc"))
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	reg, err = types.NewRegistry()
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	expected = types.NewDynamicMap(reg, map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}}).(traits.Mapper)
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}

	reg, err = types.NewRegistry(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	msg := &proto3pb.TestAllTypes{SingleString: "abc"}
	expected = reg.NativeToValue(msg)
	val, err = RefValueToValue(expected)
	if err != nil {
		t.Fatalf("RefValueToValue() failed: &v", err)
	}
	actual, err = ValueToRefValue(stdEnv.TypeAdapter(), val)
	if err != nil {
		t.Fatalf("ValueToRefValue() failed: &v", err)
	}
	if !expected.Equal(actual) {
		t.Fatalf("got val %v, want %v", actual, expected)
	}
 }

func TestAstToProto(t *testing.T) {
	stdEnv, _ := NewEnv(Declarations(
		decls.NewVar("a", decls.Dyn),
		decls.NewVar("b", decls.Dyn),
	))
	ast, iss := stdEnv.Parse("a + b")
	if iss.Err() != nil {
		t.Fatalf("Parse('a + b') failed: %v", iss.Err())
	}
	parsed, err := AstToParsedExpr(ast)
	if err != nil {
		t.Fatalf("AstToParsedExpr() failed: %v", err)
	}
	ast2 := ParsedExprToAst(parsed)
	if !proto.Equal(ast2.Expr(), ast.Expr()) {
		t.Errorf("got expr %v, wanted %v", ast2, ast)
	}
	ast3 := ParsedExprToAstWithSource(parsed, ast.Source())
	if !proto.Equal(ast3.Expr(), ast.Expr()) {
		t.Errorf("got expr %v, wanted %v", ast2, ast)
	}
	if ast3.Source() != ast.Source() {
		t.Errorf("got source %v, wanted %v", ast3.Source(), ast.Source())
	}

	_, err = AstToCheckedExpr(ast)
	if err == nil {
		t.Error("expected error converting unchecked ast")
	}
	ast, iss = stdEnv.Check(ast)
	if iss != nil && iss.Err() != nil {
		t.Fatalf("stdEnv.Check(ast) failed: %v", iss.Err())
	}
	checked, err := AstToCheckedExpr(ast)
	if err != nil {
		t.Fatalf("AstToCheckeExpr(ast) failed: %v", err)
	}
	ast4 := CheckedExprToAst(checked)
	if !proto.Equal(ast4.Expr(), ast.Expr()) {
		t.Fatalf("got ast %v, wanted %v", ast4, ast)
	}
	ast5 := CheckedExprToAstWithSource(checked, ast.Source())
	if !proto.Equal(ast5.Expr(), ast.Expr()) {
		t.Errorf("got expr %v, wanted %v", ast5, ast)
	}
	if ast5.Source() != ast.Source() {
		t.Errorf("got source %v, wanted %v", ast5.Source(), ast.Source())
	}
}

func TestAstToString(t *testing.T) {
	stdEnv, err := NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	in := "a + b - (c ? (-d + 4) : e)"
	ast, iss := stdEnv.Parse(in)
	if iss.Err() != nil {
		t.Fatalf("stdEnv.Parse(%q) failed: %v", in, iss.Err())
	}
	expr, err := AstToString(ast)
	if err != nil {
		t.Fatalf("AstToString(ast) failed: %v", err)
	}
	if expr != in {
		t.Errorf("got %v, wanted %v", expr, in)
	}
}

func TestCheckedExprToAstConstantExpr(t *testing.T) {
	stdEnv, err := NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	in := "10"
	ast, iss := stdEnv.Compile(in)
	if iss.Err() != nil {
		t.Fatalf("stdEnv.Compile(%q) failed: %v", in, iss.Err())
	}
	expr, err := AstToCheckedExpr(ast)
	if err != nil {
		t.Fatalf("AstToCheckedExpr(ast) failed: %v", err)
	}
	ast2 := CheckedExprToAst(expr)
	if !proto.Equal(ast2.Expr(), ast.Expr()) {
		t.Fatalf("got ast %v, wanted %v", ast2, ast)
	}
}

func TestCheckedExprToAstMissingInfo(t *testing.T) {
	stdEnv, err := NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	in := "10"
	ast, iss := stdEnv.Parse(in)
	if iss.Err() != nil {
		t.Fatalf("stdEnv.Compile(%q) failed: %v", in, iss.Err())
	}
	if ast.ResultType() != decls.Dyn {
		t.Fatalf("ast.ResultType() got %v, wanted 'dyn'", ast.ResultType())
	}
	expr, err := AstToParsedExpr(ast)
	if err != nil {
		t.Fatalf("AstToParsedExpr(ast) failed: %v", err)
	}
	checkedExpr := &exprpb.CheckedExpr{
		TypeMap: map[int64]*exprpb.Type{expr.GetExpr().GetId(): decls.Int},
		Expr:    expr.GetExpr(),
	}
	ast2 := CheckedExprToAst(checkedExpr)
	if !ast2.IsChecked() {
		t.Fatal("CheckedExprToAst() did not produce a 'checked' ast")
	}
	if ast2.ResultType() != decls.Int {
		t.Fatalf("ast2.ResultType() got %v, wanted 'int'", ast.ResultType())
	}
}
