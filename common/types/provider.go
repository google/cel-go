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
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/google/cel-go/common/types/pb"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	anypb "google.golang.org/protobuf/types/known/anypb"
	dpb "google.golang.org/protobuf/types/known/durationpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type protoTypeRegistry struct {
	revTypeMap map[string]ref.Type
	pbdb       *pb.Db
}

// NewRegistry accepts a list of proto message instances and returns a type
// provider which can create new instances of the provided message or any
// message that proto depends upon in its FileDescriptor.
func NewRegistry(types ...proto.Message) ref.TypeRegistry {
	p := &protoTypeRegistry{
		revTypeMap: make(map[string]ref.Type),
		pbdb:       pb.NewDb(),
	}
	p.RegisterType(
		BoolType,
		BytesType,
		DoubleType,
		DurationType,
		IntType,
		ListType,
		MapType,
		NullType,
		StringType,
		TimestampType,
		TypeType,
		UintType)

	for _, msgType := range types {
		err := p.RegisterMessage(msgType)
		if err != nil {
			panic(err)
		}
	}
	return p
}

// NewEmptyRegistry returns a registry which is completely unconfigured.
func NewEmptyRegistry() ref.TypeRegistry {
	return &protoTypeRegistry{
		revTypeMap: make(map[string]ref.Type),
		pbdb:       pb.NewDb(),
	}
}

// Copy implements the ref.TypeRegistry interface method which copies the current state of the
// registry into its own memory space.
func (p *protoTypeRegistry) Copy() ref.TypeRegistry {
	copy := &protoTypeRegistry{
		revTypeMap: make(map[string]ref.Type),
		pbdb:       p.pbdb.Copy(),
	}
	for k, v := range p.revTypeMap {
		copy.revTypeMap[k] = v
	}
	return copy
}

func (p *protoTypeRegistry) EnumValue(enumName string) ref.Val {
	enumVal, err := p.pbdb.DescribeEnum(enumName)
	if err != nil {
		return NewErr("unknown enum name '%s'", enumName)
	}
	return Int(enumVal.Value())
}

func (p *protoTypeRegistry) FindFieldType(messageType string,
	fieldName string) (*ref.FieldType, bool) {
	msgType, err := p.pbdb.DescribeType(messageType)
	if err != nil {
		return nil, false
	}
	field, found := msgType.FieldByName(fieldName)
	if !found {
		return nil, false
	}
	return &ref.FieldType{
			Type:    field.CheckedType(),
			IsSet:   field.IsSet,
			GetFrom: field.GetFrom},
		true
}

func (p *protoTypeRegistry) FindIdent(identName string) (ref.Val, bool) {
	if t, found := p.revTypeMap[identName]; found {
		return t.(ref.Val), true
	}
	if enumVal, err := p.pbdb.DescribeEnum(identName); err == nil {
		return Int(enumVal.Value()), true
	}
	return nil, false
}

func (p *protoTypeRegistry) FindType(typeName string) (*exprpb.Type, bool) {
	if _, err := p.pbdb.DescribeType(typeName); err != nil {
		return nil, false
	}
	if typeName != "" && typeName[0] == '.' {
		typeName = typeName[1:]
	}
	return &exprpb.Type{
		TypeKind: &exprpb.Type_Type{
			Type: &exprpb.Type{
				TypeKind: &exprpb.Type_MessageType{
					MessageType: typeName}}}}, true
}

func (p *protoTypeRegistry) NewValue(typeName string, fields map[string]ref.Val) ref.Val {
	td, err := p.pbdb.DescribeType(typeName)
	if err != nil {
		return NewErr("unknown type '%s'", typeName)
	}
	msg := td.New()
	fieldMap := td.FieldMap()
	for name, value := range fields {
		field, found := fieldMap[name]
		if !found {
			return NewErr("no such field: %s", name)
		}
		msgSetField(msg, field, value)
	}
	return p.NativeToValue(msg.Interface())
}

func (p *protoTypeRegistry) RegisterDescriptor(fileDesc protoreflect.FileDescriptor) error {
	fd, err := p.pbdb.RegisterDescriptor(fileDesc)
	if err != nil {
		return err
	}
	return p.registerAllTypes(fd)
}

func (p *protoTypeRegistry) RegisterMessage(message proto.Message) error {
	fd, err := p.pbdb.RegisterMessage(message)
	if err != nil {
		return err
	}
	return p.registerAllTypes(fd)
}

func (p *protoTypeRegistry) RegisterType(types ...ref.Type) error {
	for _, t := range types {
		p.revTypeMap[t.TypeName()] = t
	}
	// TODO: generate an error when the type name is registered more than once.
	return nil
}

func (p *protoTypeRegistry) registerAllTypes(fd *pb.FileDescription) error {
	for _, typeName := range fd.GetTypeNames() {
		err := p.RegisterType(NewObjectTypeValue(typeName))
		if err != nil {
			return err
		}
	}
	return nil
}

// NativeToValue converts various "native" types to ref.Val with this specific implementation
// providing support for custom proto-based types.
//
// This method should be the inverse of ref.Val.ConvertToNative.
func (p *protoTypeRegistry) NativeToValue(value interface{}) ref.Val {
	switch v := value.(type) {
	case ref.Val:
		return v
	// Adapt common types and aggregate specializations using the DefaultTypeAdapter.
	case bool, *bool,
		float32, *float32, float64, *float64,
		int, *int, int32, *int32, int64, *int64,
		string, *string,
		uint, *uint, uint32, *uint32, uint64, *uint64,
		[]byte,
		[]string,
		map[string]string:
		return DefaultTypeAdapter.NativeToValue(value)
	// Adapt well-known proto-types using the DefaultTypeAdapter.
	case *dpb.Duration,
		*tpb.Timestamp,
		*structpb.ListValue,
		structpb.NullValue,
		*structpb.Struct,
		*structpb.Value,
		*wrapperspb.BoolValue,
		*wrapperspb.BytesValue,
		*wrapperspb.DoubleValue,
		*wrapperspb.FloatValue,
		*wrapperspb.Int32Value,
		*wrapperspb.Int64Value,
		*wrapperspb.StringValue,
		*wrapperspb.UInt32Value,
		*wrapperspb.UInt64Value:
		return DefaultTypeAdapter.NativeToValue(value)
	// Override the Any type by ensuring that custom proto-types are considered on recursive calls.
	case *anypb.Any:
		if v == nil {
			return NewErr("unsupported type conversion: '%T'", value)
		}
		unpackedAny, err := v.UnmarshalNew()
		if err != nil {
			return NewErr("anypb.UnmarshalNew() failed for type %q: %v", v.GetTypeUrl(), err)
		}
		return p.NativeToValue(unpackedAny)
	// Convert custom proto types to CEL values based on type's presence within the pb.Db.
	case proto.Message:
		typeName := string(v.ProtoReflect().Descriptor().FullName())
		td, err := p.pbdb.DescribeType(typeName)
		if err != nil {
			return NewErr("unknown type: '%s'", typeName)
		}
		typeVal, found := p.FindIdent(typeName)
		if !found {
			return NewErr("unknown type: '%s'", typeName)
		}
		return NewObject(p, td, typeVal.(*TypeValue), v)
	// Override default handling for list and maps to ensure that blends of Go + proto types
	// are appropriately adapted on recursive calls or subsequent inspection of the aggregate
	// value.
	default:
		refValue := reflect.ValueOf(value)
		if refValue.Kind() == reflect.Ptr {
			if refValue.IsNil() {
				return NewErr("unsupported type conversion: '%T'", value)
			}
			refValue = refValue.Elem()
		}
		refKind := refValue.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			return NewDynamicList(p, value)
		case reflect.Map:
			return NewDynamicMap(p, value)
		}
	}
	// By default return the default type adapter's conversion to CEL.
	return DefaultTypeAdapter.NativeToValue(value)
}

// defaultTypeAdapter converts go native types to CEL values.
type defaultTypeAdapter struct{}

var (
	// DefaultTypeAdapter adapts canonical CEL types from their equivalent Go values.
	DefaultTypeAdapter = &defaultTypeAdapter{}
)

// NativeToValue implements the ref.TypeAdapter interface.
func (a *defaultTypeAdapter) NativeToValue(value interface{}) ref.Val {
	switch value.(type) {
	case nil:
		return NullValue
	case *Bool:
		if ptr := value.(*Bool); ptr != nil {
			return ptr
		}
	case *Bytes:
		if ptr := value.(*Bytes); ptr != nil {
			return ptr
		}
	case *Double:
		if ptr := value.(*Double); ptr != nil {
			return ptr
		}
	case *Int:
		if ptr := value.(*Int); ptr != nil {
			return ptr
		}
	case *String:
		if ptr := value.(*String); ptr != nil {
			return ptr
		}
	case *Uint:
		if ptr := value.(*Uint); ptr != nil {
			return ptr
		}
	case ref.Val:
		return value.(ref.Val)
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
	case *bool:
		if ptr := value.(*bool); ptr != nil {
			return Bool(*ptr)
		}
	case *float32:
		if ptr := value.(*float32); ptr != nil {
			return Double(*ptr)
		}
	case *float64:
		if ptr := value.(*float64); ptr != nil {
			return Double(*ptr)
		}
	case *int:
		if ptr := value.(*int); ptr != nil {
			return Int(*ptr)
		}
	case *int32:
		if ptr := value.(*int32); ptr != nil {
			return Int(*ptr)
		}
	case *int64:
		if ptr := value.(*int64); ptr != nil {
			return Int(*ptr)
		}
	case *string:
		if ptr := value.(*string); ptr != nil {
			return String(*ptr)
		}
	case *uint:
		if ptr := value.(*uint); ptr != nil {
			return Uint(*ptr)
		}
	case *uint32:
		if ptr := value.(*uint32); ptr != nil {
			return Uint(*ptr)
		}
	case *uint64:
		if ptr := value.(*uint64); ptr != nil {
			return Uint(*ptr)
		}
	case []byte:
		return Bytes(value.([]byte))
	case []string:
		return NewStringList(a, value.([]string))
	case map[string]string:
		return NewStringStringMap(a, value.(map[string]string))
	case *dpb.Duration:
		if ptr := value.(*dpb.Duration); ptr != nil {
			return Duration{ptr}
		}
	case *structpb.ListValue:
		if ptr := value.(*structpb.ListValue); ptr != nil {
			return NewJSONList(a, ptr)
		}
	case structpb.NullValue, *structpb.NullValue:
		return NullValue
	case *structpb.Struct:
		if ptr := value.(*structpb.Struct); ptr != nil {
			return NewJSONStruct(a, ptr)
		}
	case *structpb.Value:
		v := value.(*structpb.Value)
		if v == nil {
			return NullValue
		}
		switch v.Kind.(type) {
		case *structpb.Value_BoolValue:
			return a.NativeToValue(v.GetBoolValue())
		case *structpb.Value_ListValue:
			return a.NativeToValue(v.GetListValue())
		case *structpb.Value_NullValue:
			return NullValue
		case *structpb.Value_NumberValue:
			return a.NativeToValue(v.GetNumberValue())
		case *structpb.Value_StringValue:
			return a.NativeToValue(v.GetStringValue())
		case *structpb.Value_StructValue:
			return a.NativeToValue(v.GetStructValue())
		}
	case *tpb.Timestamp:
		if ptr := value.(*tpb.Timestamp); ptr != nil {
			return Timestamp{ptr}
		}
	case *anypb.Any:
		val := value.(*anypb.Any)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		unpackedAny, err := val.UnmarshalNew()
		if err != nil {
			return NewErr("anypb.UnmarshalNew() failed for type %q: %v", val.GetTypeUrl(), err)
		}
		return a.NativeToValue(unpackedAny)
	case *wrapperspb.BoolValue:
		val := value.(*wrapperspb.BoolValue)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Bool(val.GetValue())
	case *wrapperspb.BytesValue:
		val := value.(*wrapperspb.BytesValue)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Bytes(val.GetValue())
	case *wrapperspb.DoubleValue:
		val := value.(*wrapperspb.DoubleValue)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Double(val.GetValue())
	case *wrapperspb.FloatValue:
		val := value.(*wrapperspb.FloatValue)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Double(val.GetValue())
	case *wrapperspb.Int32Value:
		val := value.(*wrapperspb.Int32Value)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Int(val.GetValue())
	case *wrapperspb.Int64Value:
		val := value.(*wrapperspb.Int64Value)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Int(val.GetValue())
	case *wrapperspb.StringValue:
		val := value.(*wrapperspb.StringValue)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return String(val.GetValue())
	case *wrapperspb.UInt32Value:
		val := value.(*wrapperspb.UInt32Value)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Uint(val.GetValue())
	case *wrapperspb.UInt64Value:
		val := value.(*wrapperspb.UInt64Value)
		if val == nil {
			return NewErr("unsupported type conversion")
		}
		return Uint(val.GetValue())
	default:
		refValue := reflect.ValueOf(value)
		if refValue.Kind() == reflect.Ptr {
			if refValue.IsNil() {
				return NewErr("unsupported type conversion: '%T'", value)
			}
			refValue = refValue.Elem()
		}
		refKind := refValue.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			return NewDynamicList(a, value)
		case reflect.Map:
			return NewDynamicMap(a, value)
		// type aliases of primitive types cannot be asserted as that type, but rather need
		// to be downcast to int32 before being converted to a CEL representation.
		case reflect.Int32:
			intType := reflect.TypeOf(int32(0))
			return Int(refValue.Convert(intType).Interface().(int32))
		case reflect.Int64:
			intType := reflect.TypeOf(int64(0))
			return Int(refValue.Convert(intType).Interface().(int64))
		case reflect.Uint32:
			uintType := reflect.TypeOf(uint32(0))
			return Uint(refValue.Convert(uintType).Interface().(uint32))
		case reflect.Uint64:
			uintType := reflect.TypeOf(uint64(0))
			return Uint(refValue.Convert(uintType).Interface().(uint64))
		case reflect.Float32:
			doubleType := reflect.TypeOf(float32(0))
			return Double(refValue.Convert(doubleType).Interface().(float32))
		case reflect.Float64:
			doubleType := reflect.TypeOf(float64(0))
			return Double(refValue.Convert(doubleType).Interface().(float64))
		}
	}
	return NewErr("unsupported type conversion: '%T'", value)
}

func msgSetField(target protoreflect.Message, field *pb.FieldDescription, val ref.Val) error {
	if field.IsList() {
		lv := target.NewField(field.Descriptor())
		list, ok := val.(traits.Lister)
		if !ok {
			msgName := field.Descriptor().ContainingMessage().FullName()
			return fmt.Errorf("unsupported field type for %v.%v: %v", msgName, field.Name(), val.Type())
		}
		err := msgSetListField(lv.List(), field, list)
		if err != nil {
			return err
		}
		target.Set(field.Descriptor(), lv)
		return nil
	}
	if field.IsMap() {
		mv := target.NewField(field.Descriptor())
		mp, ok := val.(traits.Mapper)
		if !ok {
			msgName := field.Descriptor().ContainingMessage().FullName()
			return fmt.Errorf("unsupported field type for %v.%v: %v", msgName, field.Name(), val.Type())
		}
		err := msgSetMapField(mv.Map(), field, mp)
		if err != nil {
			return err
		}
		target.Set(field.Descriptor(), mv)
		return nil
	}
	v, err := val.ConvertToNative(field.ReflectType())
	if err != nil {
		msgName := field.Descriptor().ContainingMessage().FullName()
		return fmt.Errorf("field type conversion error for %v.%v: %v", msgName, field.Name(), err)
	}
	switch v.(type) {
	case proto.Message:
		v = v.(proto.Message).ProtoReflect()
	}
	target.Set(field.Descriptor(), protoreflect.ValueOf(v))
	return nil
}

func msgSetListField(target protoreflect.List, elemType *pb.FieldDescription, listVal traits.Lister) error {
	elemReflectType := elemType.ReflectType().Elem()
	for i := Int(0); i < listVal.Size().(Int); i++ {
		elem := listVal.Get(i)
		elemVal, err := elem.ConvertToNative(elemReflectType)
		if err != nil {
			msgName := elemType.Descriptor().ContainingMessage().FullName()
			return fmt.Errorf("field type conversion error for %v.%v: %v", msgName, elemType.Name(), err)
		}
		switch ev := elemVal.(type) {
		case proto.Message:
			elemVal = ev.ProtoReflect()
		}
		target.Append(protoreflect.ValueOf(elemVal))
	}
	return nil
}

func msgSetMapField(target protoreflect.Map, entryType *pb.FieldDescription, mapVal traits.Mapper) error {
	targetKeyType := entryType.KeyType.ReflectType()
	targetValType := entryType.ValueType.ReflectType()
	it := mapVal.Iterator()
	for it.HasNext() == True {
		key := it.Next()
		val := mapVal.Get(key)
		k, err := key.ConvertToNative(targetKeyType)
		if err != nil {
			msgName := entryType.Descriptor().ContainingMessage().FullName()
			return fmt.Errorf("field type conversion error for %v.%v key type: %v", msgName, entryType.Name(), err)
		}
		v, err := val.ConvertToNative(targetValType)
		if err != nil {
			msgName := entryType.Descriptor().ContainingMessage().FullName()
			return fmt.Errorf("field type conversion error for %v.%v value type: %v", msgName, entryType.Name(), err)
		}
		switch v.(type) {
		case proto.Message:
			v = v.(proto.Message).ProtoReflect()
		}
		target.Set(protoreflect.ValueOf(k).MapKey(), protoreflect.ValueOf(v))
	}
	return nil
}
