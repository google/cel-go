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
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/cel-go/common/types/ref"
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
		t.Fatalf("Int.ConvertToNative(int32) failed: %v", err)
	}
	if val.(int32) != 20050 {
		t.Errorf("Got '%v', expected 20050", val)
	}
	val, err = Int(math.MaxInt32 + 1).ConvertToNative(reflect.TypeOf(int32(0)))
	if err == nil {
		t.Errorf("(MaxInt+1).ConvertToNative(int32) did not error, got: %v", val)
	} else if !strings.Contains(err.Error(), "integer overflow") {
		t.Errorf("ConvertToNative(int32) returned unexpected error: %v, wanted integer overflow", err)
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
	tests := []struct {
		name   string
		in     int64
		toType ref.Type
		out    interface{}
	}{
		{
			name:   "IntToType",
			in:     int64(4),
			toType: TypeType,
			out:    IntType.TypeName(),
		},
		{
			name:   "IntToInt",
			in:     int64(4),
			toType: IntType,
			out:    int64(4),
		},
		{
			name:   "IntToUint",
			in:     int64(4),
			toType: UintType,
			out:    uint64(4),
		},
		{
			name:   "IntToUintOverflow",
			in:     -1,
			toType: UintType,
			out:    errUintOverflow,
		},
		{
			name:   "IntToDouble",
			in:     int64(4),
			toType: DoubleType,
			out:    float64(4),
		},
		{
			name:   "IntToString",
			in:     int64(-4),
			toType: StringType,
			out:    "-4",
		},
		{
			name:   "IntToTimestamp",
			in:     int64(946684800),
			toType: TimestampType,
			out:    time.Unix(946684800, 0).UTC(),
		},
		{
			name:   "IntToTimestampPosOverflow",
			in:     maxUnixTime + 1,
			toType: TimestampType,
			out:    errTimestampOverflow,
		},
		{
			name:   "IntToTimestampMinOverflow",
			in:     minUnixTime - 1,
			toType: TimestampType,
			out:    errTimestampOverflow,
		},
		{
			name:   "IntToUnsupportedType",
			in:     int64(4),
			toType: DurationType,
			out:    errors.New("type conversion error"),
		},
	}
	for _, tst := range tests {
		got := Int(tst.in).ConvertToType(tst.toType).Value()
		var eq bool
		switch gotVal := got.(type) {
		case time.Time:
			eq = gotVal.Equal(tst.out.(time.Time))
		case error:
			eq = strings.Contains(gotVal.Error(), tst.out.(error).Error())
		default:
			eq = reflect.DeepEqual(gotVal, tst.out)
		}
		if !eq {
			t.Errorf("Int(%v).ConvertToType(%v) failed, got: %v, wanted: %v",
				tst.in, tst.toType, got, tst.out)
		}
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
