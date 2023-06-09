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
	"testing"
	"time"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
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
			in:  NullableType(BoolType),
			out: "bool",
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
			out: "int",
		},
		{
			in:  ObjectType("google.protobuf.UInt32Value"),
			out: "uint",
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
