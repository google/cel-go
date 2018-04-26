package types

import(
	"testing"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/duration"
)

func TestString_Add(t *testing.T) {
	if String("hello").Add(String(" world")) != String("hello world") {
		t.Error("Adding two strings did not produce a concatenated value.")
	}
	if !IsError(String("goodbye").Add(Int(1))) {
		t.Error("Adding a string an int did not cause an error")
	}
}

func TestString_Compare(t *testing.T) {
	a := String("a")
	b := String("bbbb")
	c := String("c")
	if a.Compare(b) != IntNegOne {
		t.Error("'a' was not less than 'bbbb'")
	}
	if a.Compare(a) != IntZero {
		t.Error("'a' was not equal to 'a'")
	}
	if c.Compare(b) != IntOne {
		t.Error("'c' was not greater than 'bbbb'")
	}
	if !IsError(a.Compare(True)) {
		t.Error("Comparison to a non-string type did not generate an error.")
	}
}

func TestString_ConvertToType(t *testing.T) {
	if !String("-1").ConvertToType(IntType).Equal(IntNegOne).(Bool) {
		t.Error("String could not be converted to int")
	}
	if !String("false").ConvertToType(BoolType).Equal(False).(Bool) {
		t.Error("String could not be coverted to bool")
	}
	if !String("1").ConvertToType(UintType).Equal(Uint(1)).(Bool) {
		t.Error("String could not be converted to uint")
	}
	if !String("2017-01-01T00:00:00Z").ConvertToType(TimestampType).
		Equal(Timestamp{&timestamp.Timestamp{Seconds: 1483228800}}).(Bool) {
		t.Error("String could not be converted to timestamp")
	}
	if !String("1h5s").ConvertToType(DurationType).
		Equal(Duration{&duration.Duration{Seconds: 3605}}).(Bool) {
		t.Error("String could not be converted to duration")
	}
	if !String("2.5").ConvertToType(DoubleType).Equal(Double(2.5)).(Bool) {
		t.Error("String could not be converted to double")
	}
	if !String("hello").ConvertToType(BytesType).Equal(Bytes("hello")).(Bool) {
		t.Error("String could not be converted to bytes")
	}
	if !String("goodbye").ConvertToType(TypeType).Equal(StringType).(Bool) {
		t.Error("String could not be converted to type")
	}
	if !String("goodbye").ConvertToType(StringType).Equal(String("goodbye")).(Bool) {
		t.Error("String could not be converted to itself")
	}
}

func TestString_Equal(t *testing.T) {
	if !String("hello").Equal(String("hello")).(Bool) {
		t.Error("Two equivalent strings were not equal")
	}
	if String("hello").Equal(String("hell")).(Bool) {
		t.Error("Two inqueal strings were found equal")
	}
	if String("c").Equal(Int(99)).(Bool) {
		t.Error("String 'c' was found equal to int 99")
	}
}

func TestString_Match(t *testing.T) {
	str := String("hello 1 world")
	sw := String("^hello")
	ew := String("\\d world$")
	if !str.Match(sw).(Bool) {
		t.Error("Did not match starts with pattern.")
	}
	if !str.Match(ew).(Bool) {
		t.Error("Did not match ends with pattern.")
	}
	if str.Match(String("ello 1 worlds")).(Bool) {
		t.Error("Matched an incorrect pattern.")
	}
	if !IsError(str.Match(Int(1))) {
		t.Error("Matched a non-string pattern.")
	}
}

func TestString_Size(t *testing.T) {
	if String("").Size().(Int) != 0 {
		t.Error("Empty string had a non-zero size")
	}
	if String("hello world").Size().(Int) != 11 {
		t.Error("String with eleven characters had incorrect size")
	}
}
