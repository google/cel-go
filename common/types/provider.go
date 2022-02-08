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
	"time"

	"github.com/google/cel-go/common/types/pb"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	dpb "google.golang.org/protobuf/types/known/durationpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	rpcpb "google.golang.org/genproto/googleapis/rpc/status"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

type protoTypeRegistry struct {
	revTypeMap map[string]ref.Type
	pbdb       *pb.Db
}

// NewRegistry accepts a list of proto message instances and returns a type
// provider which can create new instances of the provided message or any
// message that proto depends upon in its FileDescriptor.
func NewRegistry(types ...proto.Message) (ref.TypeRegistry, error) {
	p := &protoTypeRegistry{
		revTypeMap: make(map[string]ref.Type),
		pbdb:       pb.NewDb(),
	}
	err := p.RegisterType(
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
	if err != nil {
		return nil, err
	}
	// This block ensures that the well-known protobuf types are registered by default.
	for _, fd := range p.pbdb.FileDescriptions() {
		err = p.registerAllTypes(fd)
		if err != nil {
			return nil, err
		}
	}
	for _, msgType := range types {
		err = p.RegisterMessage(msgType)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
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
	enumVal, found := p.pbdb.DescribeEnum(enumName)
	if !found {
		return NewErr("unknown enum name '%s'", enumName)
	}
	return Int(enumVal.Value())
}

func (p *protoTypeRegistry) FindFieldType(messageType string,
	fieldName string) (*ref.FieldType, bool) {
	msgType, found := p.pbdb.DescribeType(messageType)
	if !found {
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
	if enumVal, found := p.pbdb.DescribeEnum(identName); found {
		return Int(enumVal.Value()), true
	}
	return nil, false
}

func (p *protoTypeRegistry) FindType(typeName string) (*exprpb.Type, bool) {
	if _, found := p.pbdb.DescribeType(typeName); !found {
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
	td, found := p.pbdb.DescribeType(typeName)
	if !found {
		return NewErr("unknown type '%s'", typeName)
	}
	msg := td.New()
	fieldMap := td.FieldMap()
	for name, value := range fields {
		field, found := fieldMap[name]
		if !found {
			return NewErr("no such field: %s", name)
		}
		err := msgSetField(msg, field, value)
		if err != nil {
			return &Err{err}
		}
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

// NativeToValue converts various "native" types to ref.Val with this specific implementation
// providing support for custom proto-based types.
//
// This method should be the inverse of ref.Val.ConvertToNative.
func (p *protoTypeRegistry) NativeToValue(value interface{}) ref.Val {
	if val, found := nativeToValue(p, value); found {
		return val
	}
	switch v := value.(type) {
	case proto.Message:
		typeName := string(v.ProtoReflect().Descriptor().FullName())
		td, found := p.pbdb.DescribeType(typeName)
		if !found {
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
	case protoreflect.List:
		return NewProtoList(p, v)
	case protoreflect.Message:
		return p.NativeToValue(v.Interface())
	case protoreflect.Value:
		return p.NativeToValue(v.Interface())
	}
	return UnsupportedRefValConversionErr(value)
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
	return UnsupportedRefValConversionErr(value)
}

// nativeToValue returns the converted (ref.Val, true) of a conversion is found,
// otherwise (nil, false)
func nativeToValue(a ref.TypeAdapter, value interface{}) (ref.Val, bool) {
	switch v := value.(type) {
	case nil:
		return NullValue, true
	case *Bool:
		if v != nil {
			return *v, true
		}
	case *Bytes:
		if v != nil {
			return *v, true
		}
	case *Double:
		if v != nil {
			return *v, true
		}
	case *Int:
		if v != nil {
			return *v, true
		}
	case *String:
		if v != nil {
			return *v, true
		}
	case *Uint:
		if v != nil {
			return *v, true
		}
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
	case *dpb.Duration:
		return Duration{Duration: v.AsDuration()}, true
	case time.Duration:
		return Duration{Duration: v}, true
	case *tpb.Timestamp:
		return Timestamp{Time: v.AsTime()}, true
	case time.Time:
		return Timestamp{Time: v}, true
	case *bool:
		if v != nil {
			return Bool(*v), true
		}
	case *float32:
		if v != nil {
			return Double(*v), true
		}
	case *float64:
		if v != nil {
			return Double(*v), true
		}
	case *int:
		if v != nil {
			return Int(*v), true
		}
	case *int32:
		if v != nil {
			return Int(*v), true
		}
	case *int64:
		if v != nil {
			return Int(*v), true
		}
	case *string:
		if v != nil {
			return String(*v), true
		}
	case *uint:
		if v != nil {
			return Uint(*v), true
		}
	case *uint32:
		if v != nil {
			return Uint(*v), true
		}
	case *uint64:
		if v != nil {
			return Uint(*v), true
		}
	case []byte:
		return Bytes(v), true
	// specializations for common lists types.
	case []string:
		return NewStringList(a, v), true
	case []ref.Val:
		return NewRefValList(a, v), true
	// specializations for common map types.
	case map[string]string:
		return NewStringStringMap(a, v), true
	case map[string]interface{}:
		return NewStringInterfaceMap(a, v), true
	case map[ref.Val]ref.Val:
		return NewRefValMap(a, v), true
	// additional specializations may be added upon request / need.
	case *anypb.Any:
		if v == nil {
			return UnsupportedRefValConversionErr(v), true
		}
		unpackedAny, err := v.UnmarshalNew()
		if err != nil {
			return NewErr("anypb.UnmarshalNew() failed for type %q: %v", v.GetTypeUrl(), err), true
		}
		return a.NativeToValue(unpackedAny), true
	case *structpb.NullValue, structpb.NullValue:
		return NullValue, true
	case *structpb.ListValue:
		return NewJSONList(a, v), true
	case *structpb.Struct:
		return NewJSONStruct(a, v), true
	case ref.Val:
		return v, true
	case protoreflect.EnumNumber:
		return Int(v), true
	case proto.Message:
		if v == nil {
			return UnsupportedRefValConversionErr(v), true
		}
		typeName := string(v.ProtoReflect().Descriptor().FullName())
		td, found := pb.DefaultDb.DescribeType(typeName)
		if !found {
			return nil, false
		}
		val, unwrapped := td.MaybeUnwrap(v)
		if !unwrapped {
			return nil, false
		}
		return a.NativeToValue(val), true
	// Note: dynamicpb.Message implements the proto.Message _and_ protoreflect.Message interfaces
	// which means that this case must appear after handling a proto.Message type.
	case protoreflect.Message:
		return a.NativeToValue(v.Interface()), true
	default:
		refValue := reflect.ValueOf(v)
		if refValue.Kind() == reflect.Ptr {
			if refValue.IsNil() {
				return UnsupportedRefValConversionErr(v), true
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

// RefValueToExprValue converts between ref.Val and exprpb.ExprValue.
func RefValueToExprValue(res ref.Val, err error) (*exprpb.ExprValue, error) {
	if err != nil {
		s := status.Convert(err).Proto()
		return &exprpb.ExprValue{
			Kind: &exprpb.ExprValue_Error{
				Error: &exprpb.ErrorSet{
					Errors: []*rpcpb.Status{s},
				},
			},
		}, nil
	}
	if IsUnknown(res) {
		return &exprpb.ExprValue{
			Kind: &exprpb.ExprValue_Unknown{
				Unknown: &exprpb.UnknownSet{
					Exprs: res.Value().([]int64),
				},
			}}, nil
	}
	v, err := RefValueToValue(res)
	if err != nil {
		return nil, err
	}
	return &exprpb.ExprValue{
		Kind: &exprpb.ExprValue_Value{Value: v}}, nil
}

// RefValueToValue converts between ref.Val and Value.
// The ref.Val must not be error or unknown.
func RefValueToValue(res ref.Val) (*exprpb.Value, error) {
	switch res.Type() {
	case BoolType:
		return &exprpb.Value{
			Kind: &exprpb.Value_BoolValue{BoolValue: res.Value().(bool)}}, nil
	case BytesType:
		return &exprpb.Value{
			Kind: &exprpb.Value_BytesValue{BytesValue: res.Value().([]byte)}}, nil
	case DoubleType:
		return &exprpb.Value{
			Kind: &exprpb.Value_DoubleValue{DoubleValue: res.Value().(float64)}}, nil
	case IntType:
		return &exprpb.Value{
			Kind: &exprpb.Value_Int64Value{Int64Value: res.Value().(int64)}}, nil
	case ListType:
		l := res.(traits.Lister)
		sz := l.Size().(Int)
		elts := make([]*exprpb.Value, 0, int64(sz))
		for i := Int(0); i < sz; i++ {
			v, err := RefValueToValue(l.Get(i))
			if err != nil {
				return nil, err
			}
			elts = append(elts, v)
		}
		return &exprpb.Value{
			Kind: &exprpb.Value_ListValue{
				ListValue: &exprpb.ListValue{Values: elts}}}, nil
	case MapType:
		mapper := res.(traits.Mapper)
		sz := mapper.Size().(Int)
		entries := make([]*exprpb.MapValue_Entry, 0, int64(sz))
		for it := mapper.Iterator(); it.HasNext().(Bool); {
			k := it.Next()
			v := mapper.Get(k)
			kv, err := RefValueToValue(k)
			if err != nil {
				return nil, err
			}
			vv, err := RefValueToValue(v)
			if err != nil {
				return nil, err
			}
			entries = append(entries, &exprpb.MapValue_Entry{Key: kv, Value: vv})
		}
		return &exprpb.Value{
			Kind: &exprpb.Value_MapValue{
				MapValue: &exprpb.MapValue{Entries: entries}}}, nil
	case NullType:
		return &exprpb.Value{
			Kind: &exprpb.Value_NullValue{}}, nil
	case StringType:
		return &exprpb.Value{
			Kind: &exprpb.Value_StringValue{StringValue: res.Value().(string)}}, nil
	case TypeType:
		typeName := res.(ref.Type).TypeName()
		return &exprpb.Value{Kind: &exprpb.Value_TypeValue{TypeValue: typeName}}, nil
	case UintType:
		return &exprpb.Value{
			Kind: &exprpb.Value_Uint64Value{Uint64Value: res.Value().(uint64)}}, nil
	case DurationType:
		d, ok := res.Value().(time.Duration)
		if !ok {
			return nil, status.New(codes.InvalidArgument, "Expected time.Duration").Err()
		}
		any, err := anypb.New(dpb.New(d))
		if err != nil {
			return nil, err
		}
		return &exprpb.Value{
			Kind: &exprpb.Value_ObjectValue{ObjectValue: any}}, nil
	case TimestampType:
		t, ok := res.Value().(time.Time)
		if !ok {
			return nil, status.New(codes.InvalidArgument, "Expected time.Time").Err()
		}
		any, err := anypb.New(tpb.New(t))
		if err != nil {
			return nil, err
		}
		return &exprpb.Value{
			Kind: &exprpb.Value_ObjectValue{ObjectValue: any}}, nil
	default:
		// Object type
		pb, ok := res.Value().(proto.Message)
		if !ok {
			return nil, status.New(codes.InvalidArgument, "Expected proto message").Err()
		}
		any, err := anypb.New(pb)
		if err != nil {
			return nil, err
		}
		return &exprpb.Value{
			Kind: &exprpb.Value_ObjectValue{ObjectValue: any}}, nil
	}
}

func msgSetField(target protoreflect.Message, field *pb.FieldDescription, val ref.Val) error {
	if field.IsList() {
		lv := target.NewField(field.Descriptor())
		list, ok := val.(traits.Lister)
		if !ok {
			return unsupportedTypeConversionError(field, val)
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
			return unsupportedTypeConversionError(field, val)
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
		return fieldTypeConversionError(field, err)
	}
	switch v.(type) {
	case proto.Message:
		v = v.(proto.Message).ProtoReflect()
	}
	target.Set(field.Descriptor(), protoreflect.ValueOf(v))
	return nil
}

func msgSetListField(target protoreflect.List, listField *pb.FieldDescription, listVal traits.Lister) error {
	elemReflectType := listField.ReflectType().Elem()
	for i := Int(0); i < listVal.Size().(Int); i++ {
		elem := listVal.Get(i)
		elemVal, err := elem.ConvertToNative(elemReflectType)
		if err != nil {
			return fieldTypeConversionError(listField, err)
		}
		switch ev := elemVal.(type) {
		case proto.Message:
			elemVal = ev.ProtoReflect()
		}
		target.Append(protoreflect.ValueOf(elemVal))
	}
	return nil
}

func msgSetMapField(target protoreflect.Map, mapField *pb.FieldDescription, mapVal traits.Mapper) error {
	targetKeyType := mapField.KeyType.ReflectType()
	targetValType := mapField.ValueType.ReflectType()
	it := mapVal.Iterator()
	for it.HasNext() == True {
		key := it.Next()
		val := mapVal.Get(key)
		k, err := key.ConvertToNative(targetKeyType)
		if err != nil {
			return fieldTypeConversionError(mapField, err)
		}
		v, err := val.ConvertToNative(targetValType)
		if err != nil {
			return fieldTypeConversionError(mapField, err)
		}
		switch v.(type) {
		case proto.Message:
			v = v.(proto.Message).ProtoReflect()
		}
		target.Set(protoreflect.ValueOf(k).MapKey(), protoreflect.ValueOf(v))
	}
	return nil
}

func unsupportedTypeConversionError(field *pb.FieldDescription, val ref.Val) error {
	msgName := field.Descriptor().ContainingMessage().FullName()
	return fmt.Errorf("unsupported field type for %v.%v: %v", msgName, field.Name(), val.Type())
}

func fieldTypeConversionError(field *pb.FieldDescription, err error) error {
	msgName := field.Descriptor().ContainingMessage().FullName()
	return fmt.Errorf("field type conversion error for %v.%v value type: %v", msgName, field.Name(), err)
}
