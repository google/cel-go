package types

import(
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
}

func TestBytes_Size(t *testing.T) {
	if !Bytes("1234567890").Size().Equal(Int(10)).(Bool) {
		t.Error("Unexpected byte count.")
	}
}