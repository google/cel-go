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
	"time"

	"github.com/google/cel-go/checker/decls"
	"google.golang.org/protobuf/proto"

	proto3pb "github.com/google/cel-go/test/proto3pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	dynamicpb "google.golang.org/protobuf/types/dynamicpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
	dpb "google.golang.org/protobuf/types/known/durationpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

func TestTypeDescription(t *testing.T) {
	pbdb := NewDb()
	types := []string{
		".google.protobuf.Any",
		".google.protobuf.BoolValue",
		".google.protobuf.BytesValue",
		".google.protobuf.DoubleValue",
		".google.protobuf.FloatValue",
		".google.protobuf.Int32Value",
		".google.protobuf.Int64Value",
		".google.protobuf.ListValue",
		".google.protobuf.Struct",
		".google.protobuf.Value",
	}
	for _, typeName := range types {
		if _, found := pbdb.DescribeType(typeName); !found {
			t.Errorf("pbdb.DescribeType(%v) not found", typeName)
		}
	}
}

func TestTypeDescriptionFieldMap(t *testing.T) {
	pbdb := NewDb()
	msg := &proto3pb.NestedTestAllTypes{}
	pbdb.RegisterMessage(msg)
	td, found := pbdb.DescribeType(string(msg.ProtoReflect().Descriptor().FullName()))
	if !found {
		t.Fatalf("pbdb.DescribeType(%v) not found", msg)
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
		t.Fatalf("pbdb.RegisterMessage(%v) failed: %v", &msg, err)
	}
	td, found := pbdb.DescribeType(string(msg.ProtoReflect().Descriptor().FullName()))
	if !found {
		t.Fatalf("pbdb.DescribeType(%v) not found", &msg)
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

func TestFieldDescriptionGetFrom(t *testing.T) {
	pbdb := NewDb()
	msg := &proto3pb.TestAllTypes{
		SingleUint64:       12,
		SingleDuration:     dpb.New(time.Duration(1234)),
		SingleTimestamp:    tpb.New(time.Unix(12345, 0).UTC()),
		SingleBoolWrapper:  wrapperspb.Bool(false),
		SingleInt32Wrapper: wrapperspb.Int32(42),
		StandaloneEnum:     proto3pb.TestAllTypes_BAR,
		NestedType: &proto3pb.TestAllTypes_SingleNestedMessage{
			SingleNestedMessage: &proto3pb.TestAllTypes_NestedMessage{
				Bb: 123,
			},
		},
		SingleValue: structpb.NewStringValue("hello world"),
		SingleStruct: jsonStruct(t, map[string]interface{}{
			"null": nil,
		}),
	}
	msgName := string(msg.ProtoReflect().Descriptor().FullName())
	_, err := pbdb.RegisterMessage(msg)
	if err != nil {
		t.Fatalf("pbdb.RegisterMessage(%q) failed: %v", msgName, err)
	}
	td, found := pbdb.DescribeType(msgName)
	if !found {
		t.Fatalf("pbdb.DescribeType(%q) not found", msgName)
	}
	expected := map[string]interface{}{
		"single_uint64":        uint64(12),
		"single_duration":      time.Duration(1234),
		"single_timestamp":     time.Unix(12345, 0).UTC(),
		"single_bool_wrapper":  false,
		"single_int32_wrapper": int32(42),
		"single_int64_wrapper": structpb.NullValue_NULL_VALUE,
		"single_nested_message": &proto3pb.TestAllTypes_NestedMessage{
			Bb: 123,
		},
		"standalone_enum": int64(1),
		"single_value":    "hello world",
		"single_struct": jsonStruct(t, map[string]interface{}{
			"null": nil,
		}),
	}
	for field, want := range expected {
		f, found := td.FieldByName(field)
		if !found {
			t.Fatalf("td.FieldByName(%q) not found", field)
		}
		got, err := f.GetFrom(msg)
		if err != nil {
			t.Fatalf("field.GetFrom() failed: %v", err)
		}
		switch g := got.(type) {
		case proto.Message:
			if !proto.Equal(g, want.(proto.Message)) {
				t.Errorf("got field %s value %v, wanted %v", field, g, want)
			}
		default:
			if !reflect.DeepEqual(g, want) {
				t.Errorf("got field %s value %v, wanted %v", field, g, want)
			}
		}
	}
}

func TestFieldDescriptionIsSet(t *testing.T) {
	pbdb := NewDb()
	msg := &proto3pb.TestAllTypes{}
	msgName := string(msg.ProtoReflect().Descriptor().FullName())
	_, err := pbdb.RegisterMessage(msg)
	if err != nil {
		t.Fatalf("pbdb.RegisterMessage(%q) failed: %v", msgName, err)
	}
	td, found := pbdb.DescribeType(msgName)
	if !found {
		t.Fatalf("pbdb.DescribeType(%q) not found", msgName)
	}

	tests := []struct {
		msg   interface{}
		field string
		isSet bool
	}{
		{
			msg:   &proto3pb.TestAllTypes{SingleBool: true},
			field: "single_bool",
			isSet: true,
		},
		{
			msg:   &proto3pb.TestAllTypes{},
			field: "single_bool",
			isSet: false,
		},
		{
			msg:   (&proto3pb.TestAllTypes{SingleBool: false}),
			field: "single_bool",
			isSet: false,
		},
		{
			msg:   (&proto3pb.TestAllTypes{}),
			field: "single_bool",
			isSet: false,
		},
		{
			msg:   reflect.ValueOf(&proto3pb.TestAllTypes{}),
			field: "single_bool",
			isSet: false,
		},
		{
			msg:   nil,
			field: "single_any",
			isSet: false,
		},
	}
	for _, tc := range tests {
		f, found := td.FieldByName(tc.field)
		if !found {
			t.Fatalf("td.FieldByName(%q) not found", tc.field)
		}
		if f.IsSet(tc.msg) != tc.isSet {
			t.Errorf("got field %s set: %v, wanted %v", tc.field, f.IsSet(tc.msg), tc.isSet)
		}
	}
}

func TestTypeDescriptionMaybeUnwrap(t *testing.T) {
	pbdb := NewDb()
	_, err := pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatalf("pbdb.RegisterMessage() failed: %v", err)
	}
	msgType := "google.protobuf.Value"
	msgDesc, found := pbdb.DescribeType(msgType)
	if !found {
		t.Fatalf("pbdb.DescribeType(%q) not found", msgType)
	}
	tests := []struct {
		in  proto.Message
		out interface{}
	}{
		{
			in:  msgDesc.Zero(),
			out: structpb.NullValue_NULL_VALUE,
		},
		{
			in:  msgDesc.New().Interface(),
			out: structpb.NullValue_NULL_VALUE,
		},
		{
			in:  anyMsg(t, wrapperspb.Bool(true)),
			out: true,
		},
		{
			in:  anyMsg(t, structpb.NewNumberValue(4.5)),
			out: float64(4.5),
		},
		{
			in:  dynMsg(t, anyMsg(t, structpb.NewNumberValue(4.5))),
			out: float64(4.5),
		},
		{
			in:  dynMsg(t, anyMsg(t, &proto3pb.TestAllTypes{SingleFloat: 123.0})),
			out: &proto3pb.TestAllTypes{SingleFloat: 123.0},
		},
		{
			in:  dynMsg(t, &structpb.ListValue{}),
			out: jsonList(t, []interface{}{}),
		},
		{
			in:  structpb.NewBoolValue(true),
			out: true,
		},
		{
			in:  structpb.NewBoolValue(false),
			out: false,
		},
		{
			in:  structpb.NewNullValue(),
			out: structpb.NullValue_NULL_VALUE,
		},
		{
			in:  &structpb.Value{},
			out: structpb.NullValue_NULL_VALUE,
		},
		{
			in:  structpb.NewNumberValue(1.5),
			out: float64(1.5),
		},
		{
			in:  structpb.NewStringValue("hello world"),
			out: "hello world",
		},
		{
			in:  structpb.NewListValue(jsonList(t, []interface{}{true, 1.0})),
			out: jsonList(t, []interface{}{true, 1.0}),
		},
		{
			in:  structpb.NewStructValue(jsonStruct(t, map[string]interface{}{"hello": "world"})),
			out: jsonStruct(t, map[string]interface{}{"hello": "world"}),
		},
		{
			in:  wrapperspb.Bool(false),
			out: false,
		},
		{
			in:  wrapperspb.Bool(true),
			out: true,
		},
		{
			in:  wrapperspb.Bytes([]byte("hello")),
			out: []byte("hello"),
		},
		{
			in:  wrapperspb.Double(-4.2),
			out: -4.2,
		},
		{
			in:  wrapperspb.Float(4.5),
			out: 4.5,
		},
		{
			in:  wrapperspb.Int32(123),
			out: int64(123),
		},
		{
			in:  wrapperspb.Int64(456),
			out: int64(456),
		},
		{
			in:  wrapperspb.String("goodbye"),
			out: "goodbye",
		},
		{
			in:  wrapperspb.UInt32(1234),
			out: uint64(1234),
		},
		{
			in:  wrapperspb.UInt64(5678),
			out: uint64(5678),
		},
		{
			in:  tpb.New(time.Unix(12345, 0).UTC()),
			out: time.Unix(12345, 0).UTC(),
		},
		{
			in:  dpb.New(time.Duration(345)),
			out: time.Duration(345),
		},
	}
	for _, tc := range tests {
		typeName := string(tc.in.ProtoReflect().Descriptor().FullName())
		td, found := pbdb.DescribeType(typeName)
		if !found {
			t.Fatalf("pbdb.DescribeType(%q) not found", typeName)
		}
		msg, unwrapped := td.MaybeUnwrap(tc.in)
		if !unwrapped {
			t.Errorf("value %v not unwrapped", tc.in)
		}
		switch val := msg.(type) {
		case proto.Message:
			if !proto.Equal(val, tc.out.(proto.Message)) {
				t.Errorf("got value %v, wanted %v", val, tc.out)
			}
		default:
			if !reflect.DeepEqual(val, tc.out) {
				t.Errorf("got value %v, wanted %v", val, tc.out)
			}
		}
	}
}

func BenchmarkTypeDescriptionMaybeUnwrap(b *testing.B) {
	pbdb := NewDb()
	pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	msgType := "google.protobuf.Value"
	msgDesc, found := pbdb.DescribeType(msgType)
	if !found {
		b.Fatalf("pbdb.DescribeType(%q) not found", msgType)
	}
	tests := []struct {
		in proto.Message
	}{
		{in: msgDesc.Zero()},
		{in: msgDesc.New().Interface()},
		{in: dynamicpb.NewMessage((&structpb.ListValue{}).ProtoReflect().Descriptor())},
		{in: structpb.NewBoolValue(true)},
		{in: structpb.NewBoolValue(false)},
		{in: structpb.NewNullValue()},
		{in: &structpb.Value{}},
		{in: structpb.NewNumberValue(1.5)},
		{in: structpb.NewStringValue("hello world")},
		{in: wrapperspb.Bool(false)},
		{in: wrapperspb.Bool(true)},
		{in: wrapperspb.Bytes([]byte("hello"))},
		{in: wrapperspb.Double(-4.2)},
		{in: wrapperspb.Float(4.5)},
		{in: wrapperspb.Int32(123)},
		{in: wrapperspb.Int64(456)},
		{in: wrapperspb.String("goodbye")},
		{in: wrapperspb.UInt32(1234)},
		{in: wrapperspb.UInt64(5678)},
		{in: tpb.New(time.Unix(12345, 0).UTC())},
		{in: dpb.New(time.Duration(345))},
		{in: &proto3pb.TestAllTypes{}},
	}
	for _, tc := range tests {
		typeName := string(tc.in.ProtoReflect().Descriptor().FullName())
		td, found := pbdb.DescribeType(typeName)
		if !found {
			b.Fatalf("pbdb.DescribeType(%q) not found", typeName)
		}
		in := tc.in
		b.Run(typeName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				td.MaybeUnwrap(in)
			}
		})
	}
}

func TestTypeDescriptionCheckedType(t *testing.T) {
	pbdb := NewDb()
	msg := &proto3pb.TestAllTypes{}
	msgName := string(msg.ProtoReflect().Descriptor().FullName())
	_, err := pbdb.RegisterMessage(msg)
	if err != nil {
		t.Fatalf("pbdb.RegisterMessage(%q) failed: %v", msgName, err)
	}
	td, found := pbdb.DescribeType(msgName)
	if !found {
		t.Fatalf("pbdb.DescribeType(%q) not found", msgName)
	}
	field, found := td.FieldByName("map_string_string")
	if !found {
		t.Fatal("td.FieldByName('map_string_string') not found")
	}
	mapType := decls.NewMapType(decls.String, decls.String)
	if !proto.Equal(field.CheckedType(), mapType) {
		t.Errorf("got checked type %v, wanted %v", field.CheckedType(), mapType)
	}
	field, found = td.FieldByName("repeated_nested_message")
	if !found {
		t.Fatal("td.FieldByName('repeated_nested_message') not found")
	}
	listType := decls.NewListType(decls.NewObjectType("google.expr.proto3.test.TestAllTypes.NestedMessage"))
	if !proto.Equal(field.CheckedType(), listType) {
		t.Errorf("got checked type %v, wanted %v", field.CheckedType(), listType)
	}
}

func dynMsg(t *testing.T, msg proto.Message) *dynamicpb.Message {
	t.Helper()
	dynPB := dynamicpb.NewMessage(msg.ProtoReflect().Descriptor())
	proto.Merge(dynPB, msg)
	return dynPB
}

func anyMsg(t *testing.T, msg proto.Message) *anypb.Any {
	t.Helper()
	pb, err := anypb.New(msg)
	if err != nil {
		t.Fatalf("anypb.New(%v) failed: %v", msg, err)
	}
	return pb
}

func jsonList(t *testing.T, elems []interface{}) *structpb.ListValue {
	t.Helper()
	l, err := structpb.NewList(elems)
	if err != nil {
		t.Fatalf("structpb.NewList() failed: %v", err)
	}
	return l
}

func jsonStruct(t *testing.T, entries map[string]interface{}) *structpb.Struct {
	t.Helper()
	s, err := structpb.NewStruct(entries)
	if err != nil {
		t.Fatalf("structpb.NewStruct() failed: %v", err)
	}
	return s
}
