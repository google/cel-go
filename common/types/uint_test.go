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

	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUint_Add(t *testing.T) {
	if !Uint(4).Add(Uint(3)).Equal(Uint(7)).(Bool) {
		t.Error("Adding two uints did not match expected value.")
	}
	if !IsError(Uint(1).Add(String("-1"))) {
		t.Error("Adding non-uint to uint was not an error.")
	}
	if lhs, rhs := uint64(math.MaxUint64), 1; !IsError(Uint(lhs).Add(Uint(rhs))) {
		t.Errorf("Expected adding %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := uint64(math.MaxUint64-1), 1; !Uint(lhs).Add(Uint(rhs)).Equal(Uint(math.MaxUint64)).(Bool) {
		t.Errorf("Expected adding %d and %d to yield %d", lhs, rhs, uint64(math.MaxUint64))
	}
}

func TestUint_Compare(t *testing.T) {
	lt := Uint(204)
	gt := Uint(1300)
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
		t.Error("Types not comparable")
	}
}

func TestUint_ConvertToNative_Any(t *testing.T) {
	val, err := Uint(math.MaxUint64).ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	want, err := anypb.New(wrapperspb.UInt64(math.MaxUint64))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got %v, wanted %v", val, want)
	}
}

func TestUint_ConvertToNative_Error(t *testing.T) {
	val, err := Uint(10000).ConvertToNative(reflect.TypeOf(int(0)))
	if err == nil {
		t.Errorf("Got '%v', expected error", val)
	}
}

func TestUint_ConvertToNative_Json(t *testing.T) {
	// Value can be represented accurately as a JSON number.
	val, err := Uint(maxIntJSON).ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	} else if !proto.Equal(val.(proto.Message),
		structpb.NewNumberValue(9007199254740991.0)) {
		t.Errorf("Got '%v', expected a json number for a 32-bit uint", val)
	}

	// Value converts to a JSON decimal string
	val, err = Uint(maxIntJSON + 1).ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	} else if !proto.Equal(val.(proto.Message), structpb.NewStringValue("9007199254740992")) {
		t.Errorf("Got '%v', expected a json string for a 64-bit uint", val)
	}
}

func TestUintConvertToNative_Uint32(t *testing.T) {
	val, err := Uint(20050).ConvertToNative(reflect.TypeOf(uint32(0)))
	if err != nil {
		t.Fatalf("Uint.ConvertToNative(uint32) failed: %v", err)
	}
	if val.(uint32) != 20050 {
		t.Errorf("Got '%v', expected 20050", val)
	}
	val, err = Uint(math.MaxUint32 + 1).ConvertToNative(reflect.TypeOf(uint32(0)))
	if err == nil {
		t.Errorf("(MaxUint+1).ConvertToNative(uint32) did not error, got: %v", val)
	} else if !strings.Contains(err.Error(), "unsigned integer overflow") {
		t.Errorf("ConvertToNative(uint32) returned unexpected error: %v, wanted unsigned integer overflow", err)
	}
}

func TestUint_ConvertToNative_Ptr_Uint32(t *testing.T) {
	ptrType := uint32(0)
	val, err := Uint(10000).ConvertToNative(reflect.TypeOf(&ptrType))
	if err != nil {
		t.Error(err)
	} else if *val.(*uint32) != uint32(10000) {
		t.Errorf("Error converting uint to *uint32. Got '%v', expected 10000.", val)
	}
}

func TestUint_ConvertToNative_Ptr_Uint64(t *testing.T) {
	ptrType := uint64(0)
	val, err := Uint(18446744073709551612).ConvertToNative(reflect.TypeOf(&ptrType))
	if err != nil {
		t.Error(err)
	} else if *val.(*uint64) != uint64(18446744073709551612) {
		t.Errorf("Error converting uint to *uint64. Got '%v', expected 18446744073709551612.", val)
	}
}

func TestUint_ConvertToNative_Wrapper(t *testing.T) {
	val, err := Uint(math.MaxUint32).ConvertToNative(uint32WrapperType)
	if err != nil {
		t.Error(err)
	}
	want := wrapperspb.UInt32(math.MaxUint32)
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got %v, wanted %v", val, want)
	}

	val, err = Uint(math.MaxUint64).ConvertToNative(uint64WrapperType)
	if err != nil {
		t.Error(err)
	}
	want2 := wrapperspb.UInt64(math.MaxUint64)
	if !proto.Equal(val.(proto.Message), want2) {
		t.Errorf("Got %v, wanted %v", val, want2)
	}
}

func TestUint_ConvertToType(t *testing.T) {
	tests := []struct {
		name   string
		in     uint64
		toType ref.Type
		out    interface{}
	}{
		{
			name:   "UintToUint",
			in:     uint64(4),
			toType: UintType,
			out:    uint64(4),
		},
		{
			name:   "UintToType",
			in:     uint64(4),
			toType: TypeType,
			out:    UintType.TypeName(),
		},
		{
			name:   "UintToInt",
			in:     uint64(4),
			toType: IntType,
			out:    int64(4),
		},
		{
			name:   "UintToIntOverflow",
			in:     uint64(math.MaxInt64) + 1,
			toType: IntType,
			out:    errIntOverflow,
		},
		{
			name:   "UintToDouble",
			in:     uint64(4),
			toType: DoubleType,
			out:    float64(4),
		},
		{
			name:   "UintToString",
			in:     uint64(4),
			toType: StringType,
			out:    "4",
		},
		{
			name:   "UintToUnsupportedType",
			in:     uint64(4),
			toType: MapType,
			out:    errors.New("type conversion error"),
		},
	}
	for _, tst := range tests {
		got := Uint(tst.in).ConvertToType(tst.toType).Value()
		var eq bool
		switch gotVal := got.(type) {
		case error:
			eq = strings.Contains(gotVal.Error(), tst.out.(error).Error())
		default:
			eq = reflect.DeepEqual(gotVal, tst.out)
		}
		if !eq {
			t.Errorf("Uint(%v).ConvertToType(%v) failed, got: %v, wanted: %v",
				tst.in, tst.toType, got, tst.out)
		}
	}
}

func TestUint_Divide(t *testing.T) {
	if !Uint(3).Divide(Uint(2)).Equal(Uint(1)).(Bool) {
		t.Error("Dividing two uints did not match expectations.")
	}
	if !IsError(uintZero.Divide(uintZero)) {
		t.Error("Divide by zero did not cause error.")
	}
	if !IsError(Uint(1).Divide(Double(-1))) {
		t.Error("Division permitted without express type-conversion.")
	}
}

func TestUint_Equal(t *testing.T) {
	if !IsError(Uint(0).Equal(False)) {
		t.Error("Uint equal to non-uint type result in non-error")
	}
}

func TestUint_Modulo(t *testing.T) {
	if !Uint(21).Modulo(Uint(2)).Equal(Uint(1)).(Bool) {
		t.Error("Unexpected result from modulus operator.")
	}
	if !IsError(Uint(21).Modulo(uintZero)) {
		t.Error("Modulus by zero did not cause error.")
	}
	if !IsError(Uint(21).Modulo(IntOne)) {
		t.Error("Modulus permitted between different types without type conversion.")
	}
}

func TestUint_Multiply(t *testing.T) {
	if !Uint(2).Multiply(Uint(2)).Equal(Uint(4)).(Bool) {
		t.Error("Multiplying two values did not match expectations.")
	}
	if !IsError(Uint(1).Multiply(Double(-4.0))) {
		t.Error("Multiplication permitted without express type-conversion.")
	}
	if lhs, rhs := uint64(math.MaxUint64/2), 3; !IsError(Uint(lhs).Multiply(Uint(rhs))) {
		t.Errorf("Expected multiplying %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := uint64(math.MaxUint64/2), 2; !Uint(lhs).Multiply(Uint(rhs)).Equal(Uint(uint64(math.MaxUint64 - 1))).(Bool) {
		t.Errorf("Expected multiplying %d and %d to yield %d", lhs, rhs, uint64(math.MaxUint64-1))
	}
}

func TestUint_Subtract(t *testing.T) {
	if !Uint(4).Subtract(Uint(3)).Equal(Uint(1)).(Bool) {
		t.Error("Subtracting two uints did not match expected value.")
	}
	if !IsError(Uint(1).Subtract(Int(1))) {
		t.Error("Subtraction permitted without express type-conversion.")
	}
	if lhs, rhs := uint64(math.MaxUint64-1), uint64(math.MaxUint64); !IsError(Uint(lhs).Subtract(Uint(rhs))) {
		t.Errorf("Expected subtracting %d and %d to result in overflow.", lhs, rhs)
	}
	if lhs, rhs := uint64(math.MaxUint64), uint64(math.MaxUint64); !Uint(lhs).Subtract(Uint(rhs)).Equal(Uint(0)).(Bool) {
		t.Errorf("Expected subtracting %d and %d to yield %d", lhs, rhs, 0)
	}
}
