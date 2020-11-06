// Copyright 2019 Google LLC
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

package pb

import (
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	proto3pb "github.com/google/cel-go/test/proto3pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestTypeDescription_Any(t *testing.T) {
	pbdb := NewDb()
	_, err := pbdb.DescribeType(".google.protobuf.Any")
	if err != nil {
		t.Error(err)
	}
}

func TestTypeDescription_Json(t *testing.T) {
	pbdb := NewDb()
	_, err := pbdb.DescribeType(".google.protobuf.Value")
	if err != nil {
		t.Error(err)
	}
}

func TestTypeDescription_JsonNotInTypeInit(t *testing.T) {
	pbdb := NewDb()
	_, err := pbdb.DescribeType(".google.protobuf.ListValue")
	if err != nil {
		t.Error(err)
	}
}

func TestTypeDescription_Wrapper(t *testing.T) {
	pbdb := NewDb()
	_, err := pbdb.DescribeType(".google.protobuf.BoolValue")
	if err != nil {
		t.Error(err)
	}
}

func TestTypeDescription_WrapperNotInTypeInit(t *testing.T) {
	pbdb := NewDb()
	if _, err := pbdb.DescribeType(".google.protobuf.BytesValue"); err != nil {
		t.Error(err)
	}
}

func TestTypeDescriptionFieldMap(t *testing.T) {
	pbdb := NewDb()
	msg := &proto3pb.NestedTestAllTypes{}
	pbdb.RegisterMessage(msg)
	td, err := pbdb.DescribeType(string(msg.ProtoReflect().Descriptor().FullName()))
	if err != nil {
		t.Error(err)
	}
	if len(td.FieldMap()) != 2 {
		t.Errorf("Unexpected field count. got '%d', wanted '%d'", len(td.FieldMap()), 2)
	}
}

func TestFieldDescription(t *testing.T) {
	pbdb := NewDb()
	msg := proto3pb.NestedTestAllTypes{}
	_, err := pbdb.RegisterMessage(&msg)
	if err != nil {
		t.Error(err)
	}
	td, err := pbdb.DescribeType(string(msg.ProtoReflect().Descriptor().FullName()))
	if err != nil {
		t.Error(err)
	}
	fd, found := td.FieldByName("payload")
	if !found {
		t.Error("Field 'payload' not found")
	}
	if fd.Name() != "payload" {
		t.Error("Unexpected struct name for field 'payload'", fd.Name())
	}
	if fd.IsOneof() {
		t.Error("Field payload is listed as a oneof and it is not.")
	}
	if fd.IsMap() {
		t.Error("Field 'payload' is listed as a map and it is not.")
	}
	if !fd.IsMessage() {
		t.Error("Field 'payload' is not marked as a message.")
	}
	if fd.IsEnum() {
		t.Error("Field 'payload' is marked as an enum.")
	}
	if fd.IsList() {
		t.Error("Field 'payload' is marked as repeated.")
	}
	// Access the field by its Go struct name and check to see that it's index
	// matches the one determined by the TypeDescription utils.
	got := fd.CheckedType()
	wanted := &exprpb.Type{
		TypeKind: &exprpb.Type_MessageType{
			MessageType: "google.expr.proto3.test.TestAllTypes",
		},
	}
	if !proto.Equal(got, wanted) {
		t.Error("Field 'payload' had an unexpected checked type.")
	}
}

func TestProtoReflection(t *testing.T) {
	pbdb := NewDb()
	_, err := pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatalf("pbdb.RegisterMessage() failed: %v", err)
	}
	td, err := pbdb.DescribeType("google.expr.proto3.test.TestAllTypes")
	if err != nil {
		t.Fatalf("pbdb.DescribeType() failed: %v", err)
	}
	fieldValues := map[string]interface{}{
		"single_int32":    int32(1),
		"standalone_enum": proto3pb.TestAllTypes_FOO.Number(),
		"repeated_int32":  []int64{1, 2, 3},
		"repeated_nested_message": []*proto3pb.TestAllTypes_NestedMessage{
			{Bb: 8},
			{Bb: 8},
		},
		"map_string_string": map[string]string{
			"hello": "world",
		},
		"map_int64_nested_type": map[int32]interface{}{
			1: &proto3pb.NestedTestAllTypes{
				Child: &proto3pb.NestedTestAllTypes{
					Payload: &proto3pb.TestAllTypes{
						SingleDouble: float64(4.2),
					},
				},
			},
		},
	}
	msg := td.New()
	for f, v := range fieldValues {
		field, found := td.FieldByName(f)
		if !found {
			t.Errorf("td.FieldByName(%q) failed", f)
		}
		t.Logf("%s=%v", f, v)
		setMsgField(msg, field, v)
	}
	got := msg.Interface()
	want := &proto3pb.TestAllTypes{
		SingleInt32:    1,
		StandaloneEnum: proto3pb.TestAllTypes_FOO,
		RepeatedInt32:  []int32{1, 2, 3},
		RepeatedNestedMessage: []*proto3pb.TestAllTypes_NestedMessage{
			{Bb: 8},
			{Bb: 8},
		},
		MapStringString: map[string]string{
			"hello": "world",
		},
		MapInt64NestedType: map[int64]*proto3pb.NestedTestAllTypes{
			1: {
				Child: &proto3pb.NestedTestAllTypes{
					Payload: &proto3pb.TestAllTypes{
						SingleDouble: float64(4.2),
					},
				},
			},
		},
	}
	if !proto.Equal(got, want) {
		t.Errorf("got message %v, wanted %v", got, want)
	}
}

func setMsgField(target protoreflect.Message, field *FieldDescription, val interface{}) {
	if field.IsList() {
		lv := target.NewField(field.Descriptor())
		setListElems(lv.List(), field, val)
		target.Set(field.Descriptor(), lv)
		return
	}
	if field.IsMap() {
		mv := target.NewField(field.Descriptor())
		setMapEntries(mv.Map(), field, val)
		target.Set(field.Descriptor(), mv)
		return
	}
	target.Set(field.Descriptor(), protoreflect.ValueOf(val))
}

func setListElems(target protoreflect.List, elemType *FieldDescription, listVal interface{}) {
	targetElemType := elemType.ReflectType().Elem()
	lv := reflect.ValueOf(listVal)
	for i := 0; i < lv.Len(); i++ {
		elem := lv.Index(i)
		if elem.Type() != targetElemType && targetElemType.Kind() != reflect.Ptr {
			elem = elem.Convert(targetElemType)
		}
		elemVal := elem.Interface()
		switch ev := elemVal.(type) {
		case proto.Message:
			elemVal = ev.ProtoReflect()
		}
		target.Append(protoreflect.ValueOf(elemVal))
	}
}

func setMapEntries(target protoreflect.Map, entryType *FieldDescription, mapVal interface{}) {
	mv := reflect.ValueOf(mapVal)
	targetKeyType := entryType.KeyType.ReflectType()
	targetValType := entryType.ValueType.ReflectType()
	// Ensure the value being set is actually a map value
	for _, k := range mv.MapKeys() {
		val := mv.MapIndex(k)
		if k.Type() != targetKeyType {
			k = k.Convert(targetKeyType)
		}
		if val.Type() != targetValType && targetValType.Kind() != reflect.Ptr {
			val = val.Convert(targetValType)
		}
		entryVal := val.Interface()
		switch ev := entryVal.(type) {
		case proto.Message:
			entryVal = ev.ProtoReflect()
		}
		target.Set(
			protoreflect.ValueOf(k.Interface()).MapKey(),
			protoreflect.ValueOf(entryVal))
	}
}
