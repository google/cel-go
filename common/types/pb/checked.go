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

package pb

import (
	structpb "github.com/golang/protobuf/ptypes/struct"

	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	emptypb "github.com/golang/protobuf/ptypes/empty"
)

var (
	// CheckedPrimitives map from proto field descriptor type to checkedpb.Type.
	CheckedPrimitives = map[descpb.FieldDescriptorProto_Type]*checkedpb.Type{
		descpb.FieldDescriptorProto_TYPE_BOOL:    checkedBool,
		descpb.FieldDescriptorProto_TYPE_BYTES:   checkedBytes,
		descpb.FieldDescriptorProto_TYPE_DOUBLE:  checkedDouble,
		descpb.FieldDescriptorProto_TYPE_FLOAT:   checkedDouble,
		descpb.FieldDescriptorProto_TYPE_INT32:   checkedInt,
		descpb.FieldDescriptorProto_TYPE_INT64:   checkedInt,
		descpb.FieldDescriptorProto_TYPE_SINT32:  checkedInt,
		descpb.FieldDescriptorProto_TYPE_SINT64:  checkedInt,
		descpb.FieldDescriptorProto_TYPE_UINT32:  checkedUint,
		descpb.FieldDescriptorProto_TYPE_UINT64:  checkedUint,
		descpb.FieldDescriptorProto_TYPE_FIXED32: checkedUint,
		descpb.FieldDescriptorProto_TYPE_FIXED64: checkedUint,
		descpb.FieldDescriptorProto_TYPE_STRING:  checkedString}

	// CheckedWellKnowns map from qualified proto type name to checkedpb.Type for
	// well-known proto types.
	CheckedWellKnowns = map[string]*checkedpb.Type{
		"google.protobuf.DoubleValue": checkedWrap(checkedDouble),
		"google.protobuf.FloatValue":  checkedWrap(checkedDouble),
		"google.protobuf.Int64Value":  checkedWrap(checkedInt),
		"google.protobuf.Int32Value":  checkedWrap(checkedInt),
		"google.protobuf.UInt64Value": checkedWrap(checkedUint),
		"google.protobuf.UInt32Value": checkedWrap(checkedUint),
		"google.protobuf.BoolValue":   checkedWrap(checkedBool),
		"google.protobuf.StringValue": checkedWrap(checkedString),
		"google.protobuf.BytesValue":  checkedWrap(checkedBytes),
		"google.protobuf.NullValue":   checkedNull,
		"google.protobuf.Timestamp":   checkedTimestamp,
		"google.protobuf.Duration":    checkedDuration,
		"google.protobuf.Struct":      checkedDyn,
		"google.protobuf.Value":       checkedDyn,
		"google.protobuf.ListValue":   checkedDyn,
		"google.protobuf.Any":         checkedAny}

	// common types
	checkedBool      = checkedPrimitive(checkedpb.Type_BOOL)
	checkedBytes     = checkedPrimitive(checkedpb.Type_BYTES)
	checkedDouble    = checkedPrimitive(checkedpb.Type_DOUBLE)
	checkedDyn       = &checkedpb.Type{TypeKind: &checkedpb.Type_Dyn{Dyn: &emptypb.Empty{}}}
	checkedInt       = checkedPrimitive(checkedpb.Type_INT64)
	checkedNull      = &checkedpb.Type{TypeKind: &checkedpb.Type_Null{Null: structpb.NullValue_NULL_VALUE}}
	checkedString    = checkedPrimitive(checkedpb.Type_STRING)
	checkedUint      = checkedPrimitive(checkedpb.Type_UINT64)
	checkedAny       = checkedWellKnown(checkedpb.Type_ANY)
	checkedDuration  = checkedWellKnown(checkedpb.Type_DURATION)
	checkedTimestamp = checkedWellKnown(checkedpb.Type_TIMESTAMP)
)
