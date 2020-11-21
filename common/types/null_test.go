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

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

func TestNullConvertToNative(t *testing.T) {
	expected := structpb.NewNullValue()
	// Json Value
	val, err := NullValue.ConvertToNative(jsonValueType)
	if err != nil {
		t.Error("Fail to convert Null to jsonValueType")
	}
	if !proto.Equal(expected, val.(proto.Message)) {
		t.Errorf("Messages were not equal, got '%v'", val)
	}

	// google.protobuf.Any
	val, err = NullValue.ConvertToNative(anyValueType)
	if err != nil {
		t.Fatalf("NullValue.ConvertToNative(%v) failed: %v", anyValueType, err)
	}
	data, err := val.(*anypb.Any).UnmarshalNew()
	if err != nil {
		t.Fatalf("val.UnmarshalNew() failed: %v", err)
	}
	if !proto.Equal(expected, data) {
		t.Errorf("Messages were not equal, got '%v'", data)
	}

	// NullValue
	val, err = NullValue.ConvertToNative(reflect.TypeOf(structpb.NullValue_NULL_VALUE))
	if err != nil {
		t.Error("Fail to convert Null to strcutpb.NullValue")
	}
	if val != structpb.NullValue_NULL_VALUE {
		t.Errorf("Messages were not equal, got '%v'", val)
	}
}

func TestNullConvertToType(t *testing.T) {
	if !NullValue.ConvertToType(NullType).Equal(NullValue).(Bool) {
		t.Error("Failed to get NullType of NullValue.")
	}

	if !NullValue.ConvertToType(StringType).Equal(String("null")).(Bool) {
		t.Error("Failed to get StringType of NullValue.")
	}
	if !NullValue.ConvertToType(TypeType).Equal(NullType).(Bool) {
		t.Error("Failed to convert NullValue to type.")
	}
}

func TestNullEqual(t *testing.T) {
	if !NullValue.Equal(NullValue).(Bool) {
		t.Error("NullValue does not equal to itself.")
	}
}

func TestNullType(t *testing.T) {
	if NullValue.Type() != NullType {
		t.Error("NullValue gets incorrect type.")
	}
}

func TestNullValue(t *testing.T) {
	if NullValue.Value() != structpb.NullValue_NULL_VALUE {
		t.Error("NullValue gets incorrect value.")
	}
}
