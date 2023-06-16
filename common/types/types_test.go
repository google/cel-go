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

package types

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestTypeString(t *testing.T) {
	tests := []struct {
		in  *Type
		out string
	}{
		{
			in:  NewListType(IntType),
			out: "list(int)",
		},
		{
			in:  NewMapType(UintType, DoubleType),
			out: "map(uint, double)",
		},
		{
			in:  BoolType,
			out: "bool",
		},
		{
			in:  DynType,
			out: "dyn",
		},
		{
			in:  NullType,
			out: "null_type",
		},
		{
			in:  NewNullableType(BoolType),
			out: "wrapper(bool)",
		},
		{
			in:  NewOptionalType(NewListType(StringType)),
			out: "optional(list(string))",
		},
		{
			in:  NewObjectType("my.type.Message"),
			out: "my.type.Message",
		},
		{
			in:  NewObjectType("google.protobuf.Int32Value"),
			out: "wrapper(int)",
		},
		{
			in:  NewObjectType("google.protobuf.UInt32Value"),
			out: "wrapper(uint)",
		},
		{
			in:  NewObjectType("google.protobuf.Value"),
			out: "dyn",
		},
		{
			in:  NewTypeTypeWithParam(StringType),
			out: "type(string)",
		},
		{
			in:  NewTypeParamType("T"),
			out: "T",
		},
	}
	for _, tst := range tests {
		if tst.in.String() != tst.out {
			t.Errorf("String() got %v, wanted %v", tst.in, tst.out)
		}
	}
}

func TestTypeIsType(t *testing.T) {
	tests := []struct {
		t1     *Type
		t2     *Type
		isType bool
	}{
		{
			t1:     StringType,
			t2:     StringType,
			isType: true,
		},
		{
			t1:     StringType,
			t2:     IntType,
			isType: false,
		},
		{
			t1:     NewOptionalType(StringType),
			t2:     NewOptionalType(IntType),
			isType: false,
		},
		{
			t1:     NewOptionalType(UintType),
			t2:     NewOptionalType(UintType),
			isType: true,
		},
		{
			t1:     NewMapType(BoolType, IntType),
			t2:     NewMapType(BoolType, IntType),
			isType: true,
		},
		{
			t1:     NewMapType(NewTypeParamType("K1"), IntType),
			t2:     NewMapType(NewTypeParamType("K2"), IntType),
			isType: true,
		},
		{
			t1:     NewMapType(NewTypeParamType("K1"), NewObjectType("my.msg.First")),
			t2:     NewMapType(NewTypeParamType("K2"), NewObjectType("my.msg.Last")),
			isType: false,
		},
	}
	for _, tst := range tests {
		if tst.t1.IsType(tst.t2) != tst.isType {
			t.Errorf("%v.IsType(%v) got %v, wanted %v", tst.t1, tst.t2, !tst.isType, tst.isType)
		}
	}
}

func TestTypeIsAssignableType(t *testing.T) {
	tests := []struct {
		t1           *Type
		t2           *Type
		isAssignable bool
	}{
		{
			t1:           NewNullableType(DoubleType),
			t2:           NullType,
			isAssignable: true,
		},
		{
			t1:           NewNullableType(DoubleType),
			t2:           DoubleType,
			isAssignable: true,
		},
		{
			t1:           NewOpaqueType("vector", NewNullableType(DoubleType)),
			t2:           NewOpaqueType("vector", NullType),
			isAssignable: true,
		},
		{
			t1:           NewOpaqueType("vector", NewNullableType(DoubleType)),
			t2:           NewOpaqueType("vector", DoubleType),
			isAssignable: true,
		},
		{
			t1:           NewOpaqueType("vector", DynType),
			t2:           NewOpaqueType("vector", NewNullableType(IntType)),
			isAssignable: true,
		},
		{
			t1:           NewObjectType("my.msg.MsgName"),
			t2:           NewObjectType("my.msg.MsgName"),
			isAssignable: true,
		},
		{
			t1:           NewMapType(NewTypeParamType("K"), IntType),
			t2:           NewMapType(StringType, IntType),
			isAssignable: true,
		},
		{
			t1:           NewMapType(StringType, IntType),
			t2:           NewMapType(NewTypeParamType("K"), IntType),
			isAssignable: false,
		},
		{
			t1:           NewOpaqueType("vector", DoubleType),
			t2:           NewOpaqueType("vector", NewNullableType(IntType)),
			isAssignable: false,
		},
		{
			t1:           NewOpaqueType("vector", NewNullableType(DoubleType)),
			t2:           NewOpaqueType("vector", DynType),
			isAssignable: false,
		},
		{
			t1:           NewObjectType("my.msg.MsgName"),
			t2:           NewObjectType("my.msg.MsgName2"),
			isAssignable: false,
		},
	}
	for _, tst := range tests {
		if tst.t1.IsAssignableType(tst.t2) != tst.isAssignable {
			t.Errorf("%v.IsAssignableType(%v) got %v, wanted %v", tst.t1, tst.t2, !tst.isAssignable, tst.isAssignable)
		}
	}
}

func TestTypeIsAssignableRuntimeType(t *testing.T) {
	if !NewNullableType(DoubleType).IsAssignableRuntimeType(NullValue) {
		t.Error("nullable double cannot be assigned from null")
	}
	if !NewNullableType(DoubleType).IsAssignableRuntimeType(Double(0.0)) {
		t.Error("nullable double cannot be assigned from double")
	}
	if !NewMapType(StringType, DurationType).IsAssignableRuntimeType(
		DefaultTypeAdapter.NativeToValue(map[string]time.Duration{})) {
		t.Error("map(string, duration) not assignable to map at runtime")
	}
	if !NewMapType(StringType, DurationType).IsAssignableRuntimeType(
		DefaultTypeAdapter.NativeToValue(map[string]time.Duration{"one": time.Duration(1)})) {
		t.Error("map(string, duration) not assignable to map at runtime")
	}
	if !NewMapType(StringType, DynType).IsAssignableRuntimeType(
		DefaultTypeAdapter.NativeToValue(map[string]time.Duration{"one": time.Duration(1)})) {
		t.Error("map(string, dyn) not assignable to map at runtime")
	}
	if NewMapType(StringType, DynType).IsAssignableRuntimeType(
		DefaultTypeAdapter.NativeToValue(map[int64]time.Duration{1: time.Duration(1)})) {
		t.Error("map(string, dyn) must not be assignable to map(int, duration) at runtime")
	}
}

func TestTypeToExprType(t *testing.T) {
	tests := []struct {
		in             *Type
		out            *exprpb.Type
		unidirectional bool
	}{
		{
			in:  NewOpaqueType("vector", DoubleType, DoubleType),
			out: chkdecls.NewAbstractType("vector", chkdecls.Double, chkdecls.Double),
		},
		{
			in:  AnyType,
			out: chkdecls.Any,
		},
		{
			in:  BoolType,
			out: chkdecls.Bool,
		},
		{
			in:  BytesType,
			out: chkdecls.Bytes,
		},
		{
			in:  DoubleType,
			out: chkdecls.Double,
		},
		{
			in:  DurationType,
			out: chkdecls.Duration,
		},
		{
			in:  DynType,
			out: chkdecls.Dyn,
		},
		{
			in:  IntType,
			out: chkdecls.Int,
		},
		{
			in:  NewListType(NewTypeParamType("T")),
			out: chkdecls.NewListType(chkdecls.NewTypeParamType("T")),
		},
		{
			in:  NewMapType(NewTypeParamType("K"), NewTypeParamType("V")),
			out: chkdecls.NewMapType(chkdecls.NewTypeParamType("K"), chkdecls.NewTypeParamType("V")),
		},
		{
			in:  NullType,
			out: chkdecls.Null,
		},
		{
			in:  NewObjectType("google.type.Expr"),
			out: chkdecls.NewObjectType("google.type.Expr"),
		},
		{
			in:  StringType,
			out: chkdecls.String,
		},
		{
			in:  TimestampType,
			out: chkdecls.Timestamp,
		},
		{
			in:  TypeType,
			out: chkdecls.NewTypeType(nil),
		},
		{
			in:  UintType,
			out: chkdecls.Uint,
		},
		{
			in:  NewNullableType(BoolType),
			out: chkdecls.NewWrapperType(chkdecls.Bool),
		},
		{
			in:  NewNullableType(BytesType),
			out: chkdecls.NewWrapperType(chkdecls.Bytes),
		},
		{
			in:  NewNullableType(DoubleType),
			out: chkdecls.NewWrapperType(chkdecls.Double),
		},
		{
			in:  NewNullableType(IntType),
			out: chkdecls.NewWrapperType(chkdecls.Int),
		},
		{
			in:  NewNullableType(StringType),
			out: chkdecls.NewWrapperType(chkdecls.String),
		},
		{
			in:  NewNullableType(UintType),
			out: chkdecls.NewWrapperType(chkdecls.Uint),
		},
		{
			in:  NewTypeTypeWithParam(NewTypeTypeWithParam(DynType)),
			out: chkdecls.NewTypeType(chkdecls.NewTypeType(chkdecls.Dyn)),
		},
		{
			in:             NewObjectType("google.protobuf.Any"),
			out:            chkdecls.Any,
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.Duration"),
			out:            chkdecls.Duration,
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.Timestamp"),
			out:            chkdecls.Timestamp,
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.Value"),
			out:            chkdecls.Dyn,
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.ListValue"),
			out:            chkdecls.NewListType(chkdecls.Dyn),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.Struct"),
			out:            chkdecls.NewMapType(chkdecls.String, chkdecls.Dyn),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.BoolValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Bool),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.BytesValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Bytes),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.DoubleValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Double),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.FloatValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Double),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.Int32Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Int),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.Int64Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Int),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.StringValue"),
			out:            chkdecls.NewWrapperType(chkdecls.String),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.UInt32Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Uint),
			unidirectional: true,
		},
		{
			in:             NewObjectType("google.protobuf.UInt64Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Uint),
			unidirectional: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			got, err := TypeToExprType(tc.in)
			if err != nil {
				t.Fatalf("TypeToExprType(%v) failed: %v", tc.in, err)
			}
			if !proto.Equal(got, tc.out) {
				t.Errorf("TypeToExprType(%v) returned %v, wanted %v", tc.in, got, tc.out)
			}
			if tc.unidirectional {
				return
			}
			roundTrip, err := ExprTypeToType(got)
			if err != nil {
				t.Fatalf("ExprTypeToType(%v) failed: %v", got, err)
			}
			if !tc.in.IsType(roundTrip) {
				t.Errorf("ExprTypeToType(%v) returned %v, wanted %v", got, roundTrip, tc.in)
			}
		})
	}
}

func TestTypeToExprTypeInvalid(t *testing.T) {
	tests := []struct {
		in  *Type
		out string
	}{
		{
			in:  &Type{Kind: ListKind, runtimeTypeName: "list"},
			out: "invalid list",
		},
		{
			in: &Type{
				Kind: ListKind,
				Parameters: []*Type{
					{Kind: MapKind, runtimeTypeName: "map"},
				},
				runtimeTypeName: "list",
			},
			out: "invalid map",
		},
		{
			in:  &Type{Kind: MapKind, runtimeTypeName: "map"},
			out: "invalid map",
		},
		{
			in: &Type{
				Kind: MapKind,
				Parameters: []*Type{
					StringType,
					{Kind: MapKind, runtimeTypeName: "map"},
				},
				runtimeTypeName: "map",
			},
			out: "invalid map",
		},
		{
			in: &Type{
				Kind: MapKind,
				Parameters: []*Type{
					{Kind: MapKind, runtimeTypeName: "map"},
					StringType,
				},
				runtimeTypeName: "map",
			},
			out: "invalid map",
		},
		{
			in: &Type{
				Kind:            TypeKind,
				Parameters:      []*Type{{Kind: ListKind, runtimeTypeName: "list"}},
				runtimeTypeName: "type",
			},
			out: "invalid list",
		},
		{
			in: NewOpaqueType("bad_list", &Type{
				Kind:            ListKind,
				runtimeTypeName: "list",
			}),
			out: "invalid list",
		},
		{
			in:  &Type{},
			out: "missing type conversion",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			_, err := TypeToExprType(tc.in)
			if err == nil || !strings.Contains(err.Error(), tc.out) {
				t.Fatalf("TypeToExprType(%v) got %v, wanted error %v", tc.in, err, tc.out)
			}
		})
	}
}

func TestExprTypeToType(t *testing.T) {
	tests := []struct {
		in  *exprpb.Type
		out *Type
	}{
		{
			in:  chkdecls.NewObjectType("google.protobuf.Any"),
			out: AnyType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Duration"),
			out: DurationType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Timestamp"),
			out: TimestampType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Value"),
			out: DynType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.ListValue"),
			out: NewListType(DynType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Struct"),
			out: NewMapType(StringType, DynType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.BoolValue"),
			out: NewNullableType(BoolType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.BytesValue"),
			out: NewNullableType(BytesType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.DoubleValue"),
			out: NewNullableType(DoubleType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.FloatValue"),
			out: NewNullableType(DoubleType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Int32Value"),
			out: NewNullableType(IntType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Int64Value"),
			out: NewNullableType(IntType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.StringValue"),
			out: NewNullableType(StringType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.UInt32Value"),
			out: NewNullableType(UintType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.UInt64Value"),
			out: NewNullableType(UintType),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			got, err := ExprTypeToType(tc.in)
			if err != nil {
				t.Fatalf("ExprTypeToType(%v) failed: %v", tc.in, err)
			}
			if !got.IsType(tc.out) {
				t.Errorf("ExprTypeToType(%v) returned %v, wanted %v", tc.in, got, tc.out)
			}
		})
	}
}

func TestExprTypeToTypeInvalid(t *testing.T) {
	tests := []struct {
		in  *exprpb.Type
		out string
	}{
		{
			in:  &exprpb.Type{},
			out: "unsupported type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_Primitive{}},
			out: "unsupported primitive type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_WellKnown{}},
			out: "unsupported well-known type",
		},
		{
			in:  chkdecls.NewListType(&exprpb.Type{}),
			out: "unsupported type",
		},
		{
			in:  chkdecls.NewMapType(&exprpb.Type{}, chkdecls.Dyn),
			out: "unsupported type",
		},
		{
			in:  chkdecls.NewMapType(chkdecls.Dyn, &exprpb.Type{}),
			out: "unsupported type",
		},
		{
			in:  chkdecls.NewAbstractType("bad", &exprpb.Type{}),
			out: "unsupported type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{}},
			out: "unsupported primitive type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_Type{Type: &exprpb.Type{TypeKind: &exprpb.Type_Function{}}}},
			out: "unsupported type",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			_, err := ExprTypeToType(tc.in)
			if err == nil || !strings.Contains(err.Error(), tc.out) {
				t.Fatalf("ExprTypeToType(%v) got %v, wanted error %v", tc.in, err, tc.out)
			}
		})
	}
}

func TestTypeHasTrait(t *testing.T) {
	if !BoolType.HasTrait(traits.ComparerType) {
		t.Error("BoolType.HasTrait(ComparerType) returned false")
	}
}

func TestTypeConvertToType(t *testing.T) {
	_, err := BoolType.ConvertToNative(reflect.TypeOf(true))
	if err == nil {
		t.Error("ConvertToNative() did not error")
	}
}
