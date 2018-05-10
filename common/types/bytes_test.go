package types

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBytes_Add(t *testing.T) {
	if !Bytes("hello").Add(Bytes("world")).Equal(Bytes("helloworld")).(Bool) {
		t.Error("Byte ranges were not successfully added.")
	}
	if !IsError(Bytes("hello").Add(String("world"))) {
		t.Error("Types combined without conversion.")
	}
}

func TestBytes_Compare(t *testing.T) {
	if !Bytes("1234").Compare(Bytes("2345")).Equal(IntNegOne).(Bool) {
		t.Error("Comparison did not yield -1")
	}
	if !Bytes("2345").Compare(Bytes("1234")).Equal(IntOne).(Bool) {
		t.Error("Comparison did not yield 1")
	}
	if !Bytes("2345").Compare(Bytes("2345")).Equal(IntZero).(Bool) {
		t.Error("Comparison did not yield 0")
	}
	if !IsError(Bytes("1").Compare(String("1"))) {
		t.Error("Comparison permitted without type conversion")
	}
}

func TestBytes_ConvertToNative_ByteSlice(t *testing.T) {
	val, err := Bytes("123").ConvertToNative(reflect.TypeOf([]byte{}))
	if err != nil || IsError(val) || !bytes.Equal(val.([]byte), []byte{49, 50, 51}) {
		t.Errorf("Got '%v', wanted []byte{49, 50, 51}", val)
	}
}

func TestBytes_ConvertToNative_Error(t *testing.T) {
	val, err := Bytes("123").ConvertToNative(reflect.TypeOf(""))
	if err == nil {
		t.Errorf("Got '%v', expected error", val)
	}
}

func TestBytes_ConvertToType(t *testing.T) {
	if !Bytes("hello world").ConvertToType(BytesType).Equal(Bytes("hello world")).(Bool) {
		t.Error("Unsupported type conversion to bytes")
	}
	if !Bytes("hello world").ConvertToType(StringType).Equal(String("hello world")).(Bool) {
		t.Error("Unsupported type conversion to string")
	}
	if !Bytes("hello world").ConvertToType(TypeType).Equal(BytesType).(Bool) {
		t.Error("Unsupported type conversion to type")
	}
	if !IsError(Bytes("hello").ConvertToType(IntType)) {
		t.Errorf("Got value, expected error")
	}
}

func TestBytes_Size(t *testing.T) {
	if !Bytes("1234567890").Size().Equal(Int(10)).(Bool) {
		t.Error("Unexpected byte count.")
	}
}
