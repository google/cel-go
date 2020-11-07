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

	for _, fd := range p.pbdb.FileDescriptions() {
		p.registerAllTypes(fd)
	}
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
	case proto.Message:
		if val, found := nativeToValue(p, value); found {
			return val
		}
		typeName := string(v.ProtoReflect().Descriptor().FullName())
		td, err := p.pbdb.DescribeType(typeName)
		if err != nil {
			return NewErr("unknown type: '%s'", typeName)
		}
		unwrapped, isUnwrapped := td.MaybeUnwrap(v)
		if isUnwrapped {
			return p.NativeToValue(unwrapped)
		}
		typeVal, found := p.FindIdent(typeName)
		if !found {
			return NewErr("unknown type: '%s'", typeName)
		}
		return NewObject(p, td, typeVal.(*TypeValue), v)
	case *pb.Map:
		return NewProtoMap(p, v)
	case *pb.List:
		return NewProtoList(p, v)
	case protoreflect.Message:
		return p.NativeToValue(v.Interface())
	case protoreflect.Value:
		return p.NativeToValue(v.Interface())
	}
	if val, found := nativeToValue(p, value); found {
		return val
	}
	return NoSuchTypeConversionForValue(value)
}

// defaultTypeAdapter converts go native types to CEL values.
type defaultTypeAdapter struct{}

var (
	// DefaultTypeAdapter adapts canonical CEL types from their equivalent Go values.
	DefaultTypeAdapter = &defaultTypeAdapter{}
)

// NativeToValue implements the ref.TypeAdapter interface.
func (a *defaultTypeAdapter) NativeToValue(value interface{}) ref.Val {
	if val, found := nativeToValue(a, value); found {
		return val
	}
	return NewErr("unsupported type conversion for %T to ref.Val", value)
}

func nativeToValue(a ref.TypeAdapter, value interface{}) (ref.Val, bool) {
	switch v := value.(type) {
	case ref.Val:
		return v, true
	case nil:
		return NullValue, true
	case *Bool:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return *v, true
	case *Bytes:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return *v, true
	case *Double:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return *v, true
	case *Int:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return *v, true
	case *String:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return *v, true
	case *Uint:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return *v, true
	case bool:
		return Bool(v), true
	case int:
		return Int(v), true
	case int32:
		return Int(v), true
	case int64:
		return Int(v), true
	case uint:
		return Uint(v), true
	case uint32:
		return Uint(v), true
	case uint64:
		return Uint(v), true
	case float32:
		return Double(v), true
	case float64:
		return Double(v), true
	case string:
		return String(v), true
	case *bool:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Bool(*v), true
	case *float32:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Double(*v), true
	case *float64:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Double(*v), true
	case *int:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Int(*v), true
	case *int32:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Int(*v), true
	case *int64:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Int(*v), true
	case *string:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return String(*v), true
	case *uint:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Uint(*v), true
	case *uint32:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Uint(*v), true
	case *uint64:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Uint(*v), true
	case []byte:
		return Bytes(value.([]byte)), true
	case []string:
		return NewStringList(a, value.([]string)), true
	case map[string]string:
		return NewStringStringMap(a, value.(map[string]string)), true
	case *dpb.Duration:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Duration{Duration: v}, true
	case *structpb.ListValue:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return NewJSONList(a, v), true
	case structpb.NullValue, *structpb.NullValue:
		return NullValue, true
	case *structpb.Struct:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return NewJSONStruct(a, v), true
	case *structpb.Value:
		if v == nil {
			return NullValue, true
		}
		switch v.Kind.(type) {
		case *structpb.Value_BoolValue:
			return nativeToValue(a, v.GetBoolValue())
		case *structpb.Value_ListValue:
			return nativeToValue(a, v.GetListValue())
		case *structpb.Value_NullValue:
			return NullValue, true
		case *structpb.Value_NumberValue:
			return nativeToValue(a, v.GetNumberValue())
		case *structpb.Value_StringValue:
			return nativeToValue(a, v.GetStringValue())
		case *structpb.Value_StructValue:
			return nativeToValue(a, v.GetStructValue())
		}
	case *tpb.Timestamp:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Timestamp{Timestamp: v}, true
	case *anypb.Any:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		unpackedAny, err := v.UnmarshalNew()
		if err != nil {
			return NewErr("anypb.UnmarshalNew() failed for type %q: %v", v.GetTypeUrl(), err), true
		}
		return a.NativeToValue(unpackedAny), true
	case *wrapperspb.BoolValue:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Bool(v.GetValue()), true
	case *wrapperspb.BytesValue:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Bytes(v.GetValue()), true
	case *wrapperspb.DoubleValue:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Double(v.GetValue()), true
	case *wrapperspb.FloatValue:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Double(v.GetValue()), true
	case *wrapperspb.Int32Value:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Int(v.GetValue()), true
	case *wrapperspb.Int64Value:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Int(v.GetValue()), true
	case *wrapperspb.StringValue:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return String(v.GetValue()), true
	case *wrapperspb.UInt32Value:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Uint(v.GetValue()), true
	case *wrapperspb.UInt64Value:
		if v == nil {
			return NoSuchTypeConversionForValue(v), true
		}
		return Uint(v.GetValue()), true
	default:
		refValue := reflect.ValueOf(v)
		if refValue.Kind() == reflect.Ptr {
			if refValue.IsNil() {
				return NoSuchTypeConversionForValue(v), true
			}
			refValue = refValue.Elem()
		}
		refKind := refValue.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			return NewDynamicList(a, v), true
		case reflect.Map:
			return NewDynamicMap(a, v), true
		// type aliases of primitive types cannot be asserted as that type, but rather need
		// to be downcast to int32 before being converted to a CEL representation.
		case reflect.Int32:
			intType := reflect.TypeOf(int32(0))
			return Int(refValue.Convert(intType).Interface().(int32)), true
		case reflect.Int64:
			intType := reflect.TypeOf(int64(0))
			return Int(refValue.Convert(intType).Interface().(int64)), true
		case reflect.Uint32:
			uintType := reflect.TypeOf(uint32(0))
			return Uint(refValue.Convert(uintType).Interface().(uint32)), true
		case reflect.Uint64:
			uintType := reflect.TypeOf(uint64(0))
			return Uint(refValue.Convert(uintType).Interface().(uint64)), true
		case reflect.Float32:
			doubleType := reflect.TypeOf(float32(0))
			return Double(refValue.Convert(doubleType).Interface().(float32)), true
		case reflect.Float64:
			doubleType := reflect.TypeOf(float64(0))
			return Double(refValue.Convert(doubleType).Interface().(float64)), true
		}
	}
	return nil, false
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
