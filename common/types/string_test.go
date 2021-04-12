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
	"reflect"
	"testing"
	"time"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
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

func TestString_ConvertToNative_Any(t *testing.T) {
	val, err := String("hello").ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	want, err := anypb.New(wrapperspb.String("hello"))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
}

func TestString_ConvertToNative_Error(t *testing.T) {
	val, err := String("hello").ConvertToNative(reflect.TypeOf(0))
	if err == nil {
		t.Errorf("Got '%v', expected error", val)
	}
}

func TestString_ConvertToNative_Json(t *testing.T) {
	val, err := String("hello").ConvertToNative(jsonValueType)
	pbVal := structpb.NewStringValue("hello")
	if err != nil {
		t.Error(err)
	} else if !proto.Equal(val.(proto.Message), pbVal) {
		t.Errorf("Got '%v', expected json Value type", val)
	}
}

func TestString_ConvertToNative_Ptr(t *testing.T) {
	ptrType := ""
	val, err := String("hello").ConvertToNative(reflect.TypeOf(&ptrType))
	if err != nil {
		t.Error(err)
	} else if *val.(*string) != "hello" {
		t.Errorf("Got '%v', expected 'hello'", val)
	}
}

func TestString_ConvertToNative_String(t *testing.T) {
	val, err := String("hello").ConvertToNative(reflect.TypeOf(""))
	if err != nil {
		t.Error(err)
	} else if val.(string) != "hello" {
		t.Errorf("Got '%v', expected 'hello'", val)
	}
}

func TestString_ConvertToNative_Wrapper(t *testing.T) {
	val, err := String("hello").ConvertToNative(stringWrapperType)
	if err != nil {
		t.Error(err)
	}
	want := wrapperspb.String("hello")
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
}

func TestString_ConvertToType(t *testing.T) {
	if !String("-1").ConvertToType(IntType).Equal(IntNegOne).(Bool) {
		t.Error("String could not be converted to int")
	}
	if !String("false").ConvertToType(BoolType).Equal(False).(Bool) {
		t.Error("String could not be converted to bool")
	}
	if !String("1").ConvertToType(UintType).Equal(Uint(1)).(Bool) {
		t.Error("String could not be converted to uint")
	}
	if !String("2017-01-01T00:00:00Z").ConvertToType(TimestampType).
		Equal(Timestamp{Time: time.Unix(1483228800, 0).UTC()}).(Bool) {
		t.Error("String could not be converted to timestamp")
	}
	if !String("1h5s").ConvertToType(DurationType).
		Equal(Duration{Duration: time.Duration(3605) * time.Second}).(Bool) {
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
	if !IsError(String("map{}").ConvertToType(MapType)) {
		t.Error("Unsupported string conversion resulted in value.")
	}
}

func TestString_Equal(t *testing.T) {
	if !String("hello").Equal(String("hello")).(Bool) {
		t.Error("Two equivalent strings were not equal")
	}
	if String("hello").Equal(String("hell")).(Bool) {
		t.Error("Two inqueal strings were found equal")
	}
	if !IsError(String("c").Equal(Int(99))) {
		t.Error("String 'c' equal to int 99 resulted in non-error")
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

func TestString_Contains(t *testing.T) {
	y := String("goodbye").Receive(
		overloads.Contains,
		overloads.ContainsString,
		[]ref.Val{String("db")})
	if y != True {
		t.Errorf("Got '%v', expected 'true'", y)
	}
	n := String("goodbye").Receive(
		overloads.Contains,
		overloads.ContainsString,
		[]ref.Val{String("ggood")})
	if n == True {
		t.Errorf("Got '%v', expected 'true'", n)
	}
}

func TestString_EndsWith(t *testing.T) {
	y := String("goodbye").Receive(
		overloads.EndsWith,
		overloads.EndsWithString,
		[]ref.Val{String("bye")})
	if y == False {
		t.Errorf("Got '%v', expected 'true'", y)
	}
	n := String("goodbye").Receive(
		overloads.EndsWith,
		overloads.EndsWithString,
		[]ref.Val{String("good")})
	if n == True {
		t.Errorf("Got '%v', expected 'true'", n)
	}
}

func TestString_StartsWith(t *testing.T) {
	y := String("goodbye").Receive(
		overloads.StartsWith,
		overloads.StartsWithString,
		[]ref.Val{String("good")})
	if y != True {
		t.Errorf("Got '%v', expected 'true'", y)
	}
	n := String("goodbye").Receive(
		overloads.StartsWith,
		overloads.StartsWithString,
		[]ref.Val{String("db")})
	if n == True {
		t.Errorf("Got '%v', expected 'true'", n)
	}
}

func TestString_Size(t *testing.T) {
	if String("").Size().(Int) != 0 {
		t.Error("Empty string had a non-zero size")
	}
	if String("hello world").Size().(Int) != 11 {
		t.Error("String with eleven characters had incorrect size")
	}
	if String("\u65e5\u672c\u8a9e").Size().(Int) != 3 {
		t.Error("String size must be code points, not UTF8 bytes")
	}
}
