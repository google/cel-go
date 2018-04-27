package types

import (
	"testing"
)

func TestBool_Compare(t *testing.T) {
	if False.Compare(True).(Int) != IntNegOne {
		t.Error("False was not less than true")
	}
	if True.Compare(False).(Int) != IntOne {
		t.Error("True was not greater than false")
	}
	if True.Compare(True).(Int) != IntZero {
		t.Error("True was not equal to true")
	}
	if False.Compare(False).(Int) != IntZero {
		t.Error("False was not equal to false")
	}
	if !IsError(True.Compare(Uint(0))) {
		t.Error("Was able to compare uncomparable types.")
	}
}

func TestBool_ConvertToType(t *testing.T) {
	if !True.ConvertToType(StringType).Equal(String("true")).(Bool) {
		t.Error("Boolean could not be converted to string")
	}
	if True.ConvertToType(BoolType) != True {
		t.Error("Boolean could not be converted to a boolean.")
	}
	if True.ConvertToType(TypeType) != BoolType {
		t.Error("Boolean could not be converted to a type.")
	}
	if !IsError(True.ConvertToType(TimestampType)) {
		t.Error("Conversion to unsupported type did not error.")
	}
}

func TestBool_Equal(t *testing.T) {
	if !True.Equal(True).(Bool) {
		t.Error("True was not equal to true")
	}
	if False.Equal(True).(Bool) {
		t.Error("False was equal to true")
	}
	if Double(0.0).Equal(False).(Bool) {
		t.Error("Cross-type equality yielded true.")
	}
}

func TestBool_Negate(t *testing.T) {
	if True.Negate() != False {
		t.Error("True did not negate to false.")
	}
	if False.Negate() != True {
		t.Error("False did not negate to true")
	}
}

func TestIsBool(t *testing.T) {
	if !IsBool(True) || !IsBool(False) {
		t.Error("Boolean values did not test as boolean.")
	}
	if IsBool(String("true")) {
		t.Error("Non-boolean value tested as boolean.")
	}
}
