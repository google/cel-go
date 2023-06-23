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

package ast

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestConvertAST(t *testing.T) {
	ast := &CheckedAST{
		Expr:       &exprpb.Expr{},
		SourceInfo: &exprpb.SourceInfo{},
		TypeMap: map[int64]*types.Type{
			1: types.BoolType,
			2: types.DynType,
		},
		ReferenceMap: map[int64]*ReferenceInfo{
			1: NewFunctionReference(overloads.LogicalNot),
			2: NewIdentReference("TRUE", types.True),
		},
	}

	exprAST := &exprpb.CheckedExpr{
		Expr:       &exprpb.Expr{},
		SourceInfo: &exprpb.SourceInfo{},
		TypeMap: map[int64]*exprpb.Type{
			1: chkdecls.Bool,
			2: chkdecls.Dyn,
		},
		ReferenceMap: map[int64]*exprpb.Reference{
			1: {OverloadId: []string{overloads.LogicalNot}},
			2: {
				Name: "TRUE",
				Value: &exprpb.Constant{
					ConstantKind: &exprpb.Constant_BoolValue{BoolValue: true},
				},
			},
		},
	}

	checkedAST, err := CheckedExprToCheckedAST(exprAST)
	if err != nil {
		t.Fatalf("CheckedExprToCheckedAST() failed: %v", err)
	}
	if !reflect.DeepEqual(checkedAST.ReferenceMap, ast.ReferenceMap) ||
		!reflect.DeepEqual(checkedAST.TypeMap, ast.TypeMap) {
		t.Errorf("conversion to AST did not produce identical results: got %v, wanted %v", checkedAST, ast)
	}
	if !checkedAST.ReferenceMap[1].Equals(ast.ReferenceMap[1]) ||
		!checkedAST.ReferenceMap[2].Equals(ast.ReferenceMap[2]) {
		t.Error("converted reference info values not equal")
	}
	checkedExpr, err := CheckedASTToCheckedExpr(ast)
	if err != nil {
		t.Fatalf("CheckedASTToCheckedExpr() failed: %v", err)
	}
	if !proto.Equal(checkedExpr, exprAST) {
		t.Errorf("conversion to protobuf did not produce identical results: got %v, wanted %v", checkedExpr, exprAST)
	}
}

func TestReferenceInfoEquals(t *testing.T) {
	tests := []struct {
		name  string
		a     *ReferenceInfo
		b     *ReferenceInfo
		equal bool
	}{
		{
			name:  "single overload equal",
			a:     NewFunctionReference(overloads.AddBytes),
			b:     NewFunctionReference(overloads.AddBytes),
			equal: true,
		},
		{
			name:  "single overload not equal",
			a:     NewFunctionReference(overloads.AddBytes),
			b:     NewFunctionReference(overloads.AddDouble),
			equal: false,
		},
		{
			name:  "single and multiple overload not equal",
			a:     NewFunctionReference(overloads.AddBytes),
			b:     NewFunctionReference(overloads.AddBytes, overloads.AddDouble),
			equal: false,
		},
		{
			name:  "multiple overloads equal",
			a:     NewFunctionReference(overloads.AddBytes, overloads.AddDouble),
			b:     NewFunctionReference(overloads.AddDouble, overloads.AddBytes),
			equal: true,
		},
		{
			name:  "identifier reference equal",
			a:     NewIdentReference("BYTES", nil),
			b:     NewIdentReference("BYTES", nil),
			equal: true,
		},
		{
			name:  "identifier reference not equal",
			a:     NewIdentReference("BYTES", nil),
			b:     NewIdentReference("TRUE", nil),
			equal: false,
		},
		{
			name:  "identifier and constant reference not equal",
			a:     NewIdentReference("BYTES", nil),
			b:     NewIdentReference("BYTES", types.Bytes("bytes")),
			equal: false,
		},
		{
			name:  "constant references equal",
			a:     NewIdentReference("BYTES", types.Bytes("bytes")),
			b:     NewIdentReference("BYTES", types.Bytes("bytes")),
			equal: true,
		},
		{
			name:  "constant references not equal",
			a:     NewIdentReference("BYTES", types.Bytes("bytes")),
			b:     NewIdentReference("BYTES", types.Bytes("bytes-other")),
			equal: false,
		},
		{
			name:  "constant and overload reference not equal",
			a:     NewIdentReference("BYTES", types.Bytes("bytes")),
			b:     NewFunctionReference(overloads.AddDouble, overloads.AddBytes),
			equal: false,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			out := tc.a.Equals(tc.b)
			if out != tc.equal {
				t.Errorf("%v.Equals(%v) got %v, wanted %v", tc.a, tc.b, out, tc.equal)
			}
		})
	}
}

func TestReferenceInfoAddOverload(t *testing.T) {
	add := NewFunctionReference(overloads.AddBytes)
	add.AddOverload(overloads.AddDouble)
	if !add.Equals(NewFunctionReference(overloads.AddBytes, overloads.AddDouble)) {
		t.Error("AddOverload() did not produce equal references")
	}
	add.AddOverload(overloads.AddDouble)
	if !add.Equals(NewFunctionReference(overloads.AddBytes, overloads.AddDouble)) {
		t.Error("repeated AddOverload() did not produce equal references")
	}
}

func TestReferenceInfoToReferenceExprError(t *testing.T) {
	out, err := ReferenceInfoToReferenceExpr(NewIdentReference("SECOND", types.Duration{Duration: time.Duration(1) * time.Second}))
	if err == nil {
		t.Errorf("ReferenceInfoToReferenceExpr() got %v, wanted error", out)
	}
}

func TestReferenceExprToReferenceInfoError(t *testing.T) {
	out, err := ReferenceExprToReferenceInfo(&exprpb.Reference{Value: &exprpb.Constant{}})
	if err == nil {
		t.Errorf("ReferenceExprToReferenceInfo() got %v, wanted error", out)
	}
}

func TestConvertVal(t *testing.T) {
	tests := []ref.Val{
		types.True,
		types.Bytes("bytes"),
		types.Double(3.2),
		types.Int(-1),
		types.NullValue,
		types.String("string"),
		types.Uint(27),
	}
	for _, tst := range tests {
		c, err := ValToConstant(tst)
		if err != nil {
			t.Errorf("ValToConstant(%v) failed: %v", tst, err)
		}
		v, err := ConstantToVal(c)
		if err != nil {
			t.Errorf("ValToConstant(%v) failed: %v", c, err)
		}
		if tst.Equal(v) != types.True {
			t.Errorf("roundtrip from %v to %v and back did not produce equal results, got %v, wanted %v", tst, c, v, tst)
		}
	}
}

func TestValToConstantError(t *testing.T) {
	out, err := ValToConstant(types.Duration{Duration: time.Duration(10)})
	if err == nil {
		t.Errorf("ValToConstant() got %v, wanted error", out)
	}
}

func TestConstantToValError(t *testing.T) {
	out, err := ConstantToVal(&exprpb.Constant{})
	if err == nil {
		t.Errorf("ConstantToVal() got %v, wanted error", out)
	}
}
