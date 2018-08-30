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

	protopb "github.com/golang/protobuf/proto"
	ptypespb "github.com/golang/protobuf/ptypes"
	anypb "github.com/golang/protobuf/ptypes/any"
	refpb "github.com/google/cel-go/common/types/ref"
	traitspb "github.com/google/cel-go/common/types/traits"
	testpb "github.com/google/cel-go/test"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
)

func TestNewProtoObject(t *testing.T) {
	parsedExpr := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3}}}
	obj := NewObject(parsedExpr).(traitspb.Indexer)
	si := obj.Get(String("source_info")).(traitspb.Indexer)
	lo := si.Get(String("line_offsets")).(traitspb.Indexer)
	if lo.Get(Int(2)).Equal(Int(3)) != True {
		t.Errorf("Could not select fields by their proto type names")
	}
	expr := obj.Get(String("expr")).(traitspb.Indexer)
	call := expr.Get(String("call_expr")).(traitspb.Indexer)
	if call.Get(String("function")).Equal(String("")) != True {
		t.Errorf("Could not traverse through default values for unset fields")
	}
}

func TestProtoObject_Iterator(t *testing.T) {
	existsMsg := NewObject(testpb.Exists.Expr).(traitspb.Iterable)
	it := existsMsg.Iterator()
	var fields []refpb.Value
	for it.HasNext() == True {
		fields = append(fields, it.Next())
	}
	if !reflect.DeepEqual(fields, []refpb.Value{String("id"), String("comprehension_expr")}) {
		t.Errorf("Got %v, wanted %v", fields, []interface{}{"id", "comprehension_expr"})
	}
}

func TestProtoObj_ConvertToNative(t *testing.T) {
	pbMessage := &exprpb.ParsedExpr{
		SourceInfo: &exprpb.SourceInfo{
			LineOffsets: []int32{1, 2, 3}}}
	objVal := NewObject(pbMessage)

	// Proto Message
	val, err := objVal.ConvertToNative(reflect.TypeOf(&exprpb.ParsedExpr{}))
	if err != nil {
		t.Error(err)
	}
	if !protopb.Equal(val.(protopb.Message), pbMessage) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), pbMessage)
	}

	// google.protobuf.Any
	anyVal, err := objVal.ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	unpackedAny := ptypespb.DynamicAny{}
	if ptypespb.UnmarshalAny(anyVal.(*anypb.Any), &unpackedAny) != nil {
		NewErr("Failed to unmarshal any")
	}
	if !protopb.Equal(unpackedAny.Message, objVal.Value().(protopb.Message)) {
		t.Errorf("Messages were not equal, expect '%v', got '%v'", objVal.Value(), unpackedAny.Message)
	}
}
