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
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/cel-go/common/types/pb"
	"github.com/google/cel-go/common/types/ref"
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	"reflect"
)

type protoTypeProvider struct {
	revTypeMap map[string]ref.Type
}

// NewProvider accepts a list of proto message instances and returns a type
// provider which can create new instances of the provided message or any
// message that proto depends upon in its FileDescriptor.
func NewProvider(types ...proto.Message) ref.TypeProvider {
	p := &protoTypeProvider{
		revTypeMap: make(map[string]ref.Type)}
	p.RegisterType(
		BoolType,
		BytesType,
		DoubleType,
		DurationType,
		DynType,
		IntType,
		ListType,
		MapType,
		NullType,
		StringType,
		TimestampType,
		TypeType,
		UintType)

	for _, msgType := range types {
		fd, err := pb.DescribeFile(msgType)
		if err != nil {
			panic(err)
		}
		for _, typeName := range fd.GetTypeNames() {
			p.RegisterType(NewObjectTypeValue(typeName))
		}
	}
	return p
}

func (p *protoTypeProvider) EnumValue(enumName string) ref.Value {
	enumVal, err := pb.DescribeEnum(enumName)
	if err != nil {
		return NewErr("unknown enum name '%s'", enumName)
	}
	return Int(enumVal.Value())
}

func (p *protoTypeProvider) FindFieldType(t *checkedpb.Type,
	fieldName string) (*ref.FieldType, bool) {
	switch t.TypeKind.(type) {
	default:
		return nil, false
	case *checkedpb.Type_MessageType:
		msgType, err := pb.DescribeType(t.GetMessageType())
		if err != nil {
			return nil, false
		}
		field, found := msgType.FieldByName(fieldName)
		if !found {
			return nil, false
		}
		return &ref.FieldType{
				Type:             field.CheckedType(),
				SupportsPresence: field.SupportsPresence()},
			true
	}
}

func (p *protoTypeProvider) FindIdent(identName string) (ref.Value, bool) {
	if t, found := p.revTypeMap[identName]; found {
		return t.(ref.Value), true
	}
	if enumVal, err := pb.DescribeEnum(identName); err == nil {
		return Int(enumVal.Value()), true
	}
	return nil, false
}

func (p *protoTypeProvider) FindType(typeName string) (*checkedpb.Type, bool) {
	if _, err := pb.DescribeType(typeName); err != nil {
		return nil, false
	}

	// TODO: verify that well-known message types are handled correctly
	if typeName != "" && typeName[0] == '.' {
		typeName = typeName[1:]
	}
	if checkedType, found := pb.CheckedWellKnowns[typeName]; found {
		return checkedType, true
	}
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_Type{
			Type: &checkedpb.Type{
				TypeKind: &checkedpb.Type_MessageType{
					MessageType: typeName}}}}, true
}

func (p *protoTypeProvider) NewValue(typeName string,
	fields map[string]ref.Value) ref.Value {
	td, err := pb.DescribeType(typeName)
	if err != nil {
		return NewErr("unknown type '%s'", typeName)
	}
	refType := td.ReflectType()
	// create the new type instance.
	value := reflect.New(refType.Elem())
	pbValue := value.Elem()

	// for all of the field names referenced, set the provided value.
	for name, value := range fields {
		fd, found := td.FieldByName(name)
		if !found {
			return NewErr("no such field '%s'", name)
		}
		refField := pbValue.Field(fd.Index())
		if !refField.IsValid() {
			return NewErr("no such field '%s'", name)
		}

		dstType := refField.Type()
		// Oneof fields are defined with wrapper structs that have a single proto.Message
		// field value. The oneof wrapper is not a proto.Message instance.
		if fd.IsOneof() {
			oneofVal := reflect.New(fd.OneofType().Elem())
			refField.Set(oneofVal)
			refField = oneofVal.Elem().Field(0)
			dstType = refField.Type()
		}
		fieldValue, err := value.ConvertToNative(dstType)
		if err != nil {
			return &Err{err}
		}
		refField.Set(reflect.ValueOf(fieldValue))
	}
	return NewObject(value.Interface().(proto.Message))
}

func (p *protoTypeProvider) RegisterType(types ...ref.Type) error {
	for _, t := range types {
		p.revTypeMap[t.TypeName()] = t
	}
	// TODO: generate an error when the type name is registered more than once.
	return nil
}

func NativeToValue(value interface{}) ref.Value {
	switch value.(type) {
	case ref.Value:
		return value.(ref.Value)
	case bool:
		return Bool(value.(bool))
	case int:
		return Int(value.(int))
	case int32:
		return Int(value.(int32))
	case int64:
		return Int(value.(int64))
	case uint:
		return Uint(value.(uint))
	case uint32:
		return Uint(value.(uint32))
	case uint64:
		return Uint(value.(uint64))
	case float32:
		return Double(value.(float32))
	case float64:
		return Double(value.(float64))
	case string:
		return String(value.(string))
	case []byte:
		return Bytes(value.([]byte))
	case []string:
		return NewStringList(value.([]string))
	case *duration.Duration:
		return Duration{value.(*duration.Duration)}
	case *structpb.ListValue:
		return NewJsonList(value.(*structpb.ListValue))
	case structpb.NullValue:
		return NullValue
	case *structpb.Struct:
		return NewJsonStruct(value.(*structpb.Struct))
	case *structpb.Value:
		v := value.(*structpb.Value)
		switch v.Kind.(type) {
		case *structpb.Value_BoolValue:
			return NativeToValue(v.GetBoolValue())
		case *structpb.Value_ListValue:
			return NativeToValue(v.GetListValue())
		case *structpb.Value_NullValue:
			return NullValue
		case *structpb.Value_NumberValue:
			return NativeToValue(v.GetNumberValue())
		case *structpb.Value_StringValue:
			return NativeToValue(v.GetStringValue())
		case *structpb.Value_StructValue:
			return NativeToValue(v.GetStructValue())
		}
	case *timestamp.Timestamp:
		return Timestamp{value.(*timestamp.Timestamp)}
	case *any.Any:
		val := value.(*any.Any)
		unpackedAny := ptypes.DynamicAny{}
		if ptypes.UnmarshalAny(val, &unpackedAny) != nil {
			NewErr("Fail to unmarshal any.")
		}
		return NativeToValue(unpackedAny.Message)
	case proto.Message:
		return NewObject(value.(proto.Message))
	default:
		refValue := reflect.ValueOf(value)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		refKind := refValue.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			return NewDynamicList(value)
		case reflect.Map:
			return NewDynamicMap(value)
		}
	}
	return NewErr("unsupported type conversion for value '%v'", value)
}
