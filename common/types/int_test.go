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

package types

import (
	"math"
	"reflect"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func TestIntAdd(t *testing.T) {
	if !Int(4).Add(Int(-3)).Equal(Int(1)).(Bool) {
		t.Error("Adding two ints did not match expected value.")
	}
	if !IsError(Int(-1).Add(String("-1"))) {
		t.Error("Adding non-int to int was not an error.")
	}
	if lhs, rhs := math.MaxInt64, 1; !IsError(Int(lhs).Add(Int(rhs))) {
		t.Errorf("Expected adding %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := math.MinInt64, -1; !IsError(Int(lhs).Add(Int(rhs))) {
		t.Errorf("Expected adding %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := math.MaxInt64-1, 1; !Int(lhs).Add(Int(rhs)).Equal(Int(math.MaxInt64)).(Bool) {
		t.Errorf("Expected adding %d and %d to yield %d", lhs, rhs, math.MaxInt64)
	}
	if lhs, rhs := math.MinInt64+1, -1; !Int(lhs).Add(Int(rhs)).Equal(Int(math.MinInt64)).(Bool) {
		t.Errorf("Expected adding %d and %d to yield %d", lhs, rhs, math.MaxInt64)
	}
}

func TestIntCompare(t *testing.T) {
	lt := Int(-1300)
	gt := Int(204)
	if !lt.Compare(gt).Equal(IntNegOne).(Bool) {
		t.Error("Comparison did not yield - 1")
	}
	if !gt.Compare(lt).Equal(IntOne).(Bool) {
		t.Error("Comparison did not yield 1")
	}
	if !gt.Compare(gt).Equal(IntZero).(Bool) {
		t.Error(("Comparison did not yield 0"))
	}
	if !IsError(gt.Compare(TypeType)) {
		t.Error("Got comparison value, expected error.")
	}
}

func TestIntConvertToNative_Any(t *testing.T) {
	val, err := Int(math.MaxInt64).ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	want, err := anypb.New(wrapperspb.Int64(math.MaxInt64))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', wanted %v", val, want)
	}
}

func TestIntConvertToNative_Error(t *testing.T) {
	val, err := Int(1).ConvertToNative(jsonStructType)
	if err == nil {
		t.Errorf("Got '%v', expected error", val)
	}
}

func TestIntConvertToNative_Int32(t *testing.T) {
	val, err := Int(20050).ConvertToNative(reflect.TypeOf(int32(0)))
	if err != nil {
		t.Error(err)
	} else if val.(int32) != 20050 {
		t.Errorf("Got '%v', expected 20050", val)
	}
}

func TestIntConvertToNative_Int64(t *testing.T) {
	// Value greater than max int32.
	val, err := Int(4147483648).ConvertToNative(reflect.TypeOf(int64(0)))
	if err != nil {
		t.Error(err)
	} else if val.(int64) != 4147483648 {
		t.Errorf("Got '%v', expected 4147483648", val)
	}
}

func TestIntConvertToNative_Json(t *testing.T) {
	// Value can be represented accurately as a JSON number.
	val, err := Int(maxIntJSON).ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	} else if !proto.Equal(val.(proto.Message),
		structpb.NewNumberValue(9007199254740991.0)) {
		t.Errorf("Got '%v', expected a json number for a 32-bit int", val)
	}

	// Value converts to a JSON decimal string.
	val, err = Int(maxIntJSON + 1).ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	} else if !proto.Equal(val.(proto.Message), structpb.NewStringValue("9007199254740992")) {
		t.Errorf("Got '%v', expected a json string for a 64-bit int", val)
	}
}

func TestIntConvertToNative_Ptr_Int32(t *testing.T) {
	ptrType := int32(0)
	val, err := Int(20050).ConvertToNative(reflect.TypeOf(&ptrType))
	if err != nil {
		t.Error(err)
	} else if *val.(*int32) != 20050 {
		t.Errorf("Got '%v', expected 20050", val)
	}
}

func TestIntConvertToNative_Ptr_Int64(t *testing.T) {
	// Value greater than max int32.
	ptrType := int64(0)
	val, err := Int(math.MaxInt32 + 1).ConvertToNative(reflect.TypeOf(&ptrType))
	if err != nil {
		t.Error(err)
	} else if *val.(*int64) != math.MaxInt32+1 {
		t.Errorf("Got '%v', expected MaxInt32 + 1", val)
	}
}

func TestIntConvertToNative_Wrapper(t *testing.T) {
	val, err := Int(math.MinInt32).ConvertToNative(int32WrapperType)
	if err != nil {
		t.Error(err)
	}
	want := wrapperspb.Int32(math.MinInt32)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', wanted %v", val, want)
	}

	val, err = Int(math.MinInt64).ConvertToNative(int64WrapperType)
	if err != nil {
		t.Error(err)
	}
	want2 := wrapperspb.Int64(math.MinInt64)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want2) {
		t.Errorf("Got '%v', wanted %v", val, want2)
	}
}

func TestIntConvertToType(t *testing.T) {
	if !Int(-4).ConvertToType(IntType).Equal(Int(-4)).(Bool) {
		t.Error("Unsuccessful type conversion to int")
	}
	if !IsError(Int(-4).ConvertToType(UintType)) {
		t.Error("Got uint, expected error.")
	}
	if !Int(-4).ConvertToType(DoubleType).Equal(Double(-4)).(Bool) {
		t.Error("Unsuccessful type conversion to double")
	}
	if !Int(-4).ConvertToType(StringType).Equal(String("-4")).(Bool) {
		t.Error("Unsuccessful type conversion to string")
	}
	if !Int(-4).ConvertToType(TypeType).Equal(IntType).(Bool) {
		t.Error("Unsuccessful type conversion to type")
	}
	if !IsError(Int(-4).ConvertToType(DurationType)) {
		t.Error("Got duration, expected error.")
	}
	tm := time.Unix(946684800, 0).UTC()
	celts := Timestamp{Time: tm}
	if !Int(946684800).ConvertToType(TimestampType).Equal(celts).(Bool) {
		t.Error("unsuccessful type conversion to timestamp")
	}
}

func TestIntDivide(t *testing.T) {
	if !Int(3).Divide(Int(2)).Equal(Int(1)).(Bool) {
		t.Error("Dividing two ints did not match expectations.")
	}
	if !IsError(IntZero.Divide(IntZero)) {
		t.Error("Divide by zero did not cause error.")
	}
	if !IsError(Int(1).Divide(Double(-1))) {
		t.Error("Division permitted without express type-conversion.")
	}
	if lhs, rhs := math.MinInt64, -1; !IsError(Int(lhs).Divide(Int(rhs))) {
		t.Errorf("Expected dividing %d and %d result in overflow.", lhs, rhs)
	}
}

func TestIntEqual(t *testing.T) {
	if !IsError(Int(0).Equal(False)) {
		t.Error("Int equal to non-int type resulted in non-error.")
	}
}

func TestIntModulo(t *testing.T) {
	if !Int(21).Modulo(Int(2)).Equal(Int(1)).(Bool) {
		t.Error("Unexpected result from modulus operator.")
	}
	if !IsError(Int(21).Modulo(IntZero)) {
		t.Error("Modulus by zero did not cause error.")
	}
	if !IsError(Int(21).Modulo(uintZero)) {
		t.Error("Modulus permitted between different types without type conversion.")
	}
	if lhs, rhs := math.MinInt64, -1; !IsError(Int(lhs).Modulo(Int(rhs))) {
		t.Errorf("Expected modulo %d and %d result in overflow.", lhs, rhs)
	}
}

func TestIntMultiply(t *testing.T) {
	if !Int(2).Multiply(Int(-2)).Equal(Int(-4)).(Bool) {
		t.Error("Multiplying two values did not match expectations.")
	}
	if !IsError(Int(1).Multiply(Double(-4.0))) {
		t.Error("Multiplication permitted without express type-conversion.")
	}
	if lhs, rhs := math.MaxInt64/2, 3; !IsError(Int(lhs).Multiply(Int(rhs))) {
		t.Errorf("Expected multiplying %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := math.MinInt64/2, 3; !IsError(Int(lhs).Multiply(Int(rhs))) {
		t.Errorf("Expected multiplying %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := math.MaxInt64/2, 2; !Int(lhs).Multiply(Int(rhs)).Equal(Int(math.MaxInt64 - 1)).(Bool) {
		t.Errorf("Expected multiplying %d and %d to yield %d", lhs, rhs, math.MaxInt64-1)
	}
	if lhs, rhs := math.MinInt64/2, 2; !Int(lhs).Multiply(Int(rhs)).Equal(Int(math.MinInt64)).(Bool) {
		t.Errorf("Expected multiplying %d and %d to yield %d", lhs, rhs, math.MinInt64)
	}
	if lhs, rhs := math.MaxInt64/2, -2; !Int(lhs).Multiply(Int(rhs)).Equal(Int(math.MinInt64 + 2)).(Bool) {
		t.Errorf("Expected multiplying %d and %d to yield %d", lhs, rhs, math.MinInt64+2)
	}
	if lhs, rhs := (math.MinInt64+2)/2, -2; !Int(lhs).Multiply(Int(rhs)).Equal(Int(math.MaxInt64 - 1)).(Bool) {
		t.Errorf("Expected multiplying %d and %d to yield %d", lhs, rhs, math.MaxInt64-1)
	}
	if lhs, rhs := math.MinInt64, -1; !IsError(Int(lhs).Multiply(Int(rhs))) {
		t.Errorf("Expected multiplying %d and %d result in overflow.", lhs, rhs)
	}
}

func TestIntNegate(t *testing.T) {
	if !Int(1).Negate().Equal(Int(-1)).(Bool) {
		t.Error("Negating int value did not succeed")
	}
	if v := math.MinInt64; !IsError(Int(v).Negate()) {
		t.Errorf("Expected negating %d to result in overflow.", v)
	}
	if v := math.MaxInt64; !Int(v).Negate().Equal(Int(math.MinInt64 + 1)).(Bool) {
		t.Errorf("Expected negating %d to yield %d", v, math.MinInt64+1)
	}
}

func TestIntSubtract(t *testing.T) {
	if !Int(4).Subtract(Int(-3)).Equal(Int(7)).(Bool) {
		t.Error("Subtracting two ints did not match expected value.")
	}
	if !IsError(Int(1).Subtract(Uint(1))) {
		t.Error("Subtraction permitted without express type-conversion.")
	}
	if lhs, rhs := math.MaxInt64, -1; !IsError(Int(lhs).Subtract(Int(rhs))) {
		t.Errorf("Expected subtracting %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := math.MinInt64, 1; !IsError(Int(lhs).Subtract(Int(rhs))) {
		t.Errorf("Expected subtracting %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := math.MaxInt64-1, -1; !Int(lhs).Subtract(Int(rhs)).Equal(Int(math.MaxInt64)).(Bool) {
		t.Errorf("Expected subtracting %d and %d to yield %d", lhs, rhs, math.MaxInt64)
	}
	if lhs, rhs := math.MinInt64+1, 1; !Int(lhs).Subtract(Int(rhs)).Equal(Int(math.MinInt64)).(Bool) {
		t.Errorf("Expected subtracting %d and %d to yield %d", lhs, rhs, math.MinInt64)
	}
}
