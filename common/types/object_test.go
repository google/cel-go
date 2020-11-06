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

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

func TestNewProtoObject(t *testing.T) {
	reg := NewRegistry()
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

func TestProtoObj_ConvertToNative(t *testing.T) {
	reg := NewRegistry(&exprpb.Expr{})
	pbMessage := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3}}}
	objVal := reg.NativeToValue(pbMessage)

	// Proto Message
	val, err := objVal.ConvertToNative(reflect.TypeOf(&exprpb.ParsedExpr{}))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), pbMessage) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), pbMessage)
	}

	// google.protobuf.Any
	anyVal, err := objVal.ConvertToNative(anyValueType)
	if err != nil {
		t.Fatalf("objVal.ConvertToNative() failed: %v", err)
	}
	anyMsg := anyVal.(*anypb.Any)
	unpackedAny, err := anyMsg.UnmarshalNew()
	if err != nil {
		t.Fatalf("")
	}
	if !proto.Equal(unpackedAny, objVal.Value().(proto.Message)) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), unpackedAny)
	}

	// JSON
	json, err := objVal.ConvertToNative(jsonValueType)
	if err != nil {
		t.Fatalf("objVal.ConvertToNative(%v) failed: %v", jsonValueType, err)
	}
	jsonBytes, err := protojson.Marshal(json.(proto.Message))
	jsonTxt := string(jsonBytes)
	if err != nil {
		t.Fatalf("protojson.Marshal(%v) failed: %v", json, err)
	}
	wantTxt := `{"sourceInfo":{"lineOffsets":[1,2,3]}}`
	if jsonTxt != wantTxt {
		t.Errorf("Got %v, wanted %v", jsonTxt, wantTxt)
	}
}
