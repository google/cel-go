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
	"encoding/json"
	"reflect"
	"testing"

	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

func TestNewProtoObject(t *testing.T) {
	reg := newTestRegistry(t)
	parsedExpr := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3}}}
	reg.RegisterMessage(parsedExpr)
	obj := reg.NativeToValue(parsedExpr).(traits.Indexer)
	si := obj.Get(String("source_info")).(traits.Indexer)
	lo := si.Get(String("line_offsets")).(traits.Indexer)
	if lo.Get(Int(2)).Equal(Int(3)) != True {
		t.Errorf("Could not select fields by their proto type names")
	}
	expr := obj.Get(String("expr")).(traits.Indexer)
	call := expr.Get(String("call_expr")).(traits.Indexer)
	if call.Get(String("function")).Equal(String("")) != True {
		t.Errorf("Could not traverse through default values for unset fields")
	}
}

func TestProtoObjectConvertToNative(t *testing.T) {
	reg := newTestRegistry(t, &exprpb.Expr{})
	msg := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3}}}
	objVal := reg.NativeToValue(msg)

	// Proto Message
	val, err := objVal.ConvertToNative(reflect.TypeOf(&exprpb.ParsedExpr{}))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), msg) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), msg)
	}

	// Dynamic protobuf
	dynPB := reg.NewValue(
		string(msg.ProtoReflect().Descriptor().FullName()),
		map[string]ref.Val{
			"source_info": reg.NativeToValue(msg.GetSourceInfo()),
		})
	if IsError(dynPB) {
		t.Fatalf("reg.NewValue() failed: %v", dynPB)
	}
	dynVal := reg.NativeToValue(dynPB)
	val, err = dynVal.ConvertToNative(reflect.TypeOf(msg))
	if err != nil {
		t.Fatalf("dynVal.ConvertToNative() failed: %v", err)
	}
	if !proto.Equal(val.(proto.Message), msg) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), msg)
	}

	// google.protobuf.Any
	anyVal, err := objVal.ConvertToNative(anyValueType)
	if err != nil {
		t.Fatalf("objVal.ConvertToNative() failed: %v", err)
	}
	anyMsg := anyVal.(*anypb.Any)
	unpackedAny, err := anyMsg.UnmarshalNew()
	if err != nil {
		t.Fatalf("UnmarshalNew() failed: %v", err)
	}
	if !proto.Equal(unpackedAny, objVal.Value().(proto.Message)) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), unpackedAny)
	}

	// JSON
	jsonVal, err := objVal.ConvertToNative(jsonValueType)
	if err != nil {
		t.Fatalf("objVal.ConvertToNative(%v) failed: %v", jsonValueType, err)
	}
	jsonBytes, err := protojson.Marshal(jsonVal.(proto.Message))
	jsonTxt := string(jsonBytes)
	if err != nil {
		t.Fatalf("protojson.Marshal(%v) failed: %v", jsonVal, err)
	}
	outMap := map[string]interface{}{}
	err = json.Unmarshal(jsonBytes, &outMap)
	if err != nil {
		t.Fatalf("json.Unmarshal(%q) failed: %v", jsonTxt, err)
	}
	want := map[string]interface{}{
		"sourceInfo": map[string]interface{}{
			"lineOffsets": []interface{}{1.0, 2.0, 3.0},
		},
	}
	if !reflect.DeepEqual(outMap, want) {
		t.Errorf("got json '%v', expected %v", outMap, want)
	}
}

func TestProtoObjectIsSet(t *testing.T) {
	msg := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3},
		},
	}
	reg := newTestRegistry(t, msg)
	objVal := reg.NativeToValue(msg).(*protoObj)
	if objVal.IsSet(String("source_info")) != True {
		t.Error("got 'source_info' not set, wanted set")
	}
	if objVal.IsSet(String("expr")) != False {
		t.Error("got 'expr' set, wanted not set")
	}
	if !IsError(objVal.IsSet(String("bad_field"))) {
		t.Error("got 'bad_field' wanted error")
	}
	if !IsError(objVal.IsSet(IntZero)) {
		t.Error("got field '0' wanted error")
	}
}

func TestProtoObjectGet(t *testing.T) {
	msg := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3},
		},
	}
	reg := newTestRegistry(t, msg)
	objVal := reg.NativeToValue(msg).(*protoObj)
	if objVal.Get(String("source_info")).Equal(reg.NativeToValue(msg.GetSourceInfo())) != True {
		t.Error("could not get 'source_info'")
	}
	if objVal.Get(String("expr")).Equal(reg.NativeToValue(&exprpb.Expr{})) != True {
		t.Errorf("did not get 'expr' default value: %v", objVal.Get(String("expr")))
	}
	if !IsError(objVal.Get(String("bad_field"))) {
		t.Error("got 'bad_field' wanted error")
	}
	if !IsError(objVal.Get(IntZero)) {
		t.Error("got field '0' wanted error")
	}
}

func TestProtoObjectConvertToType(t *testing.T) {
	msg := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3},
		},
	}
	reg := newTestRegistry(t, msg)
	objVal := reg.NativeToValue(msg)
	tv := objVal.Type().(*TypeValue)
	if objVal.ConvertToType(TypeType).Equal(tv) != True {
		t.Errorf("got non-type value: %v, wanted objet type", objVal.ConvertToType(TypeType))
	}
	if objVal.ConvertToType(objVal.Type()) != objVal {
		t.Error("identity type conversion failed")
	}
}
