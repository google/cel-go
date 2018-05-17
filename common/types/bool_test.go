package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
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

func TestBool_ConvertToNative_Bool(t *testing.T) {
	refType := reflect.TypeOf(true)
	val, err := True.ConvertToNative(refType)
	if err != nil || IsError(val) || !val.(bool) {
		t.Error("Error during conversion to bool", err, val)
	}
}

func TestBool_ConvertToNative_Error(t *testing.T) {
	refType := reflect.TypeOf("")
	val, err := True.ConvertToNative(refType)
	if err == nil {
		t.Error("Got '%v', expected error", val)
	}
}

func TestBool_ConvertToNative_Json(t *testing.T) {
	val, err := True.ConvertToNative(jsonValueType)
	pbVal := &structpb.Value{&structpb.Value_BoolValue{true}}
	if err != nil ||
		IsError(val) ||
		!proto.Equal(val.(proto.Message), pbVal) {
		t.Error("Error during conversion to json Value type", err, val)
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
		t.Error("Got value, expected error.")
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
