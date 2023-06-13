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

package decls

import (
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestTypeString(t *testing.T) {
	tests := []struct {
		in  *Type
		out string
	}{
		{
			in:  ListType(IntType),
			out: "list(int)",
		},
		{
			in:  MapType(UintType, DoubleType),
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
			in:  NullableType(BoolType),
			out: "wrapper(bool)",
		},
		{
			in:  OptionalType(ListType(StringType)),
			out: "optional(list(string))",
		},
		{
			in:  ObjectType("my.type.Message"),
			out: "my.type.Message",
		},
		{
			in:  ObjectType("google.protobuf.Int32Value"),
			out: "wrapper(int)",
		},
		{
			in:  ObjectType("google.protobuf.UInt32Value"),
			out: "wrapper(uint)",
		},
		{
			in:  ObjectType("google.protobuf.Value"),
			out: "dyn",
		},
		{
			in:  TypeTypeWithParam(StringType),
			out: "type(string)",
		},
		{
			in:  TypeParamType("T"),
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
			t1:     OptionalType(StringType),
			t2:     OptionalType(IntType),
			isType: false,
		},
		{
			t1:     OptionalType(UintType),
			t2:     OptionalType(UintType),
			isType: true,
		},
		{
			t1:     MapType(BoolType, IntType),
			t2:     MapType(BoolType, IntType),
			isType: true,
		},
		{
			t1:     MapType(TypeParamType("K1"), IntType),
			t2:     MapType(TypeParamType("K2"), IntType),
			isType: true,
		},
		{
			t1:     MapType(TypeParamType("K1"), ObjectType("my.msg.First")),
			t2:     MapType(TypeParamType("K2"), ObjectType("my.msg.Last")),
			isType: false,
		},
	}
	for _, tst := range tests {
		if tst.t1.IsType(tst.t2) != tst.isType {
			t.Errorf("%v.IsType(%v) got %v, wanted %v", tst.t1, tst.t2, !tst.isType, tst.isType)
		}
	}
}

func TestTypeTypeVariable(t *testing.T) {
	tests := []struct {
		t *Type
		v *VariableDecl
	}{
		{
			t: AnyType,
			v: NewVariable("google.protobuf.Any", TypeTypeWithParam(AnyType)),
		},
		{
			t: DynType,
			v: NewVariable("dyn", TypeTypeWithParam(DynType)),
		},
		{
			t: ObjectType("google.protobuf.Int32Value"),
			v: NewVariable("int", TypeTypeWithParam(NullableType(IntType))),
		},
		{
			t: ObjectType("google.protobuf.Int32Value"),
			v: NewVariable("int", TypeTypeWithParam(NullableType(IntType))),
		},
	}
	for _, tst := range tests {
		if !tst.t.TypeVariable().DeclarationEquals(tst.v) {
			t.Errorf("got not equal %v.Equals(%v)", tst.t.TypeVariable(), tst.v)
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
			t1:           NullableType(DoubleType),
			t2:           NullType,
			isAssignable: true,
		},
		{
			t1:           NullableType(DoubleType),
			t2:           DoubleType,
			isAssignable: true,
		},
		{
			t1:           OpaqueType("vector", NullableType(DoubleType)),
			t2:           OpaqueType("vector", NullType),
			isAssignable: true,
		},
		{
			t1:           OpaqueType("vector", NullableType(DoubleType)),
			t2:           OpaqueType("vector", DoubleType),
			isAssignable: true,
		},
		{
			t1:           OpaqueType("vector", DynType),
			t2:           OpaqueType("vector", NullableType(IntType)),
			isAssignable: true,
		},
		{
			t1:           ObjectType("my.msg.MsgName"),
			t2:           ObjectType("my.msg.MsgName"),
			isAssignable: true,
		},
		{
			t1:           MapType(TypeParamType("K"), IntType),
			t2:           MapType(StringType, IntType),
			isAssignable: true,
		},
		{
			t1:           MapType(StringType, IntType),
			t2:           MapType(TypeParamType("K"), IntType),
			isAssignable: false,
		},
		{
			t1:           OpaqueType("vector", DoubleType),
			t2:           OpaqueType("vector", NullableType(IntType)),
			isAssignable: false,
		},
		{
			t1:           OpaqueType("vector", NullableType(DoubleType)),
			t2:           OpaqueType("vector", DynType),
			isAssignable: false,
		},
		{
			t1:           ObjectType("my.msg.MsgName"),
			t2:           ObjectType("my.msg.MsgName2"),
			isAssignable: false,
		},
	}
	for _, tst := range tests {
		if tst.t1.IsAssignableType(tst.t2) != tst.isAssignable {
			t.Errorf("%v.IsAssignableType(%v) got %v, wanted %v", tst.t1, tst.t2, !tst.isAssignable, tst.isAssignable)
		}
	}
}

func TestTypeSignatureEquals(t *testing.T) {
	paramA := TypeParamType("A")
	paramB := TypeParamType("B")
	mapOfAB := MapType(paramA, paramB)
	fn, err := NewFunction(overloads.Size,
		MemberOverload(overloads.SizeMapInst, []*Type{mapOfAB}, IntType),
		Overload(overloads.SizeMap, []*Type{mapOfAB}, IntType))
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	if !fn.Overloads[overloads.SizeMap].SignatureEquals(fn.Overloads[overloads.SizeMap]) {
		t.Errorf("SignatureEquals() returned false, wanted true")
	}
	if fn.Overloads[overloads.SizeMap].SignatureEquals(fn.Overloads[overloads.SizeMapInst]) {
		t.Errorf("SignatureEquals() returned false, wanted true")
	}
}

func TestTypeIsAssignableRuntimeType(t *testing.T) {
	if !NullableType(DoubleType).IsAssignableRuntimeType(types.NullValue) {
		t.Error("nullable double cannot be assigned from null")
	}
	if !NullableType(DoubleType).IsAssignableRuntimeType(types.Double(0.0)) {
		t.Error("nullable double cannot be assigned from double")
	}
	if !MapType(StringType, DurationType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[string]time.Duration{})) {
		t.Error("map(string, duration) not assignable to map at runtime")
	}
	if !MapType(StringType, DurationType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[string]time.Duration{"one": time.Duration(1)})) {
		t.Error("map(string, duration) not assignable to map at runtime")
	}
	if !MapType(StringType, DynType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[string]time.Duration{"one": time.Duration(1)})) {
		t.Error("map(string, dyn) not assignable to map at runtime")
	}
	if MapType(StringType, DynType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[int64]time.Duration{1: time.Duration(1)})) {
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
			in:  OpaqueType("vector", DoubleType, DoubleType),
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
			in:  ListType(TypeParamType("T")),
			out: chkdecls.NewListType(chkdecls.NewTypeParamType("T")),
		},
		{
			in:  MapType(TypeParamType("K"), TypeParamType("V")),
			out: chkdecls.NewMapType(chkdecls.NewTypeParamType("K"), chkdecls.NewTypeParamType("V")),
		},
		{
			in:  NullType,
			out: chkdecls.Null,
		},
		{
			in:  ObjectType("google.type.Expr"),
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
			in:  NullableType(BoolType),
			out: chkdecls.NewWrapperType(chkdecls.Bool),
		},
		{
			in:  NullableType(BytesType),
			out: chkdecls.NewWrapperType(chkdecls.Bytes),
		},
		{
			in:  NullableType(DoubleType),
			out: chkdecls.NewWrapperType(chkdecls.Double),
		},
		{
			in:  NullableType(IntType),
			out: chkdecls.NewWrapperType(chkdecls.Int),
		},
		{
			in:  NullableType(StringType),
			out: chkdecls.NewWrapperType(chkdecls.String),
		},
		{
			in:  NullableType(UintType),
			out: chkdecls.NewWrapperType(chkdecls.Uint),
		},
		{
			in:  TypeTypeWithParam(TypeTypeWithParam(DynType)),
			out: chkdecls.NewTypeType(chkdecls.NewTypeType(chkdecls.Dyn)),
		},
		{
			in:             ObjectType("google.protobuf.Any"),
			out:            chkdecls.Any,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Duration"),
			out:            chkdecls.Duration,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Timestamp"),
			out:            chkdecls.Timestamp,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Value"),
			out:            chkdecls.Dyn,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.ListValue"),
			out:            chkdecls.NewListType(chkdecls.Dyn),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Struct"),
			out:            chkdecls.NewMapType(chkdecls.String, chkdecls.Dyn),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.BoolValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Bool),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.BytesValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Bytes),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.DoubleValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Double),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.FloatValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Double),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Int32Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Int),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Int64Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Int),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.StringValue"),
			out:            chkdecls.NewWrapperType(chkdecls.String),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.UInt32Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Uint),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.UInt64Value"),
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
			in:  &Type{Kind: ListKind, runtimeType: types.ListType},
			out: "invalid list",
		},
		{
			in: &Type{
				Kind: ListKind,
				Parameters: []*Type{
					{Kind: MapKind, runtimeType: types.MapType},
				},
				runtimeType: types.ListType,
			},
			out: "invalid map",
		},
		{
			in:  &Type{Kind: MapKind, runtimeType: types.MapType},
			out: "invalid map",
		},
		{
			in: &Type{
				Kind: MapKind,
				Parameters: []*Type{
					StringType,
					{Kind: MapKind, runtimeType: types.MapType},
				},
				runtimeType: types.MapType,
			},
			out: "invalid map",
		},
		{
			in: &Type{
				Kind: MapKind,
				Parameters: []*Type{
					{Kind: MapKind, runtimeType: types.MapType},
					StringType,
				},
				runtimeType: types.MapType,
			},
			out: "invalid map",
		},
		{
			in: &Type{
				Kind:        TypeKind,
				Parameters:  []*Type{{Kind: ListKind, runtimeType: types.ListType}},
				runtimeType: types.TypeType,
			},
			out: "invalid list",
		},
		{
			in: OpaqueType("bad_list", &Type{
				Kind:        ListKind,
				runtimeType: types.ListType,
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
			out: ListType(DynType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Struct"),
			out: MapType(StringType, DynType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.BoolValue"),
			out: NullableType(BoolType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.BytesValue"),
			out: NullableType(BytesType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.DoubleValue"),
			out: NullableType(DoubleType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.FloatValue"),
			out: NullableType(DoubleType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Int32Value"),
			out: NullableType(IntType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Int64Value"),
			out: NullableType(IntType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.StringValue"),
			out: NullableType(StringType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.UInt32Value"),
			out: NullableType(UintType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.UInt64Value"),
			out: NullableType(UintType),
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
