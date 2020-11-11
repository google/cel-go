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
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	dynamicpb "google.golang.org/protobuf/types/dynamicpb"
	anypb "google.golang.org/protobuf/types/known/anypb"
	dpb "google.golang.org/protobuf/types/known/durationpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type description interface {
	WrapperField() protoreflect.FieldDescriptor
	Zero() proto.Message
}

// NewTypeDescription produces a TypeDescription value for the fully-qualified proto type name
// with a given descriptor.
//
// The type description creation method also expects the type to be marked clearly as a proto2 or
// proto3 type, and accepts a typeResolver reference for resolving field TypeDescription during
// lazily initialization of the type which is done atomically.
func NewTypeDescription(typeName string, desc protoreflect.MessageDescriptor) *TypeDescription {
	msgType := dynamicpb.NewMessageType(desc)
	msgZero := dynamicpb.NewMessage(desc).Interface()
	fieldMap := map[string]*FieldDescription{}
	fields := desc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fieldMap[string(f.Name())] = NewFieldDescription(f)
	}
	return &TypeDescription{
		typeName:     typeName,
		desc:         desc,
		msgType:      msgType,
		wrapperField: wrapperMsg(desc),
		fieldMap:     fieldMap,
		reflectType:  reflect.TypeOf(msgZero),
		zeroMsg:      msgZero,
	}
}

// TypeDescription is a collection of type metadata relevant to expression
// checking and evaluation.
type TypeDescription struct {
	typeName     string
	desc         protoreflect.MessageDescriptor
	msgType      protoreflect.MessageType
	fieldMap     map[string]*FieldDescription
	wrapperField protoreflect.FieldDescriptor
	reflectType  reflect.Type
	zeroMsg      proto.Message
}

// FieldMap returns a string field name to FieldDescription map.
func (td *TypeDescription) FieldMap() map[string]*FieldDescription {
	return td.fieldMap
}

// FieldByName returns the FieldDescription associated with a field name.
func (td *TypeDescription) FieldByName(name string) (*FieldDescription, bool) {
	fd, found := td.fieldMap[name]
	if !found {
		return nil, false
	}
	return fd, true
}

// MaybeUnwrap accepts a proto message as input and unwraps it to a primitive CEL type if possible.
//
// This method returns the unwrapped value and 'true', else the original value and 'false'.
func (td *TypeDescription) MaybeUnwrap(msg proto.Message) (interface{}, bool) {
	return unwrap(td, msg)
}

// Name of the type.
func (td *TypeDescription) Name() string {
	return string(td.desc.FullName())
}

// New returns a mutable proto message
func (td *TypeDescription) New() protoreflect.Message {
	return td.msgType.New()
}

// ReflectType returns the Golang reflect.Type for this type.
func (td *TypeDescription) ReflectType() reflect.Type {
	return td.reflectType
}

// WrapperField returns the FieldDescriptor for the 'value' field of a wrapper type proto.
//
// When the type is not a google.protobuf wrapper type, the method returns nil.
func (td *TypeDescription) WrapperField() protoreflect.FieldDescriptor {
	return td.wrapperField
}

// Zero returns the zero proto.Message value for this type.
func (td *TypeDescription) Zero() proto.Message {
	return td.zeroMsg
}

// NewFieldDescription creates a new field description from a protoreflect.FieldDescriptor.
func NewFieldDescription(fieldDesc protoreflect.FieldDescriptor) *FieldDescription {
	var reflectType reflect.Type
	var zeroMsg proto.Message
	switch fieldDesc.Kind() {
	case protoreflect.EnumKind:
		reflectType = reflect.TypeOf(protoreflect.EnumNumber(0))
	case protoreflect.MessageKind:
		zeroMsg = dynamicpb.NewMessage(fieldDesc.Message())
		reflectType = reflect.TypeOf(zeroMsg)
	default:
		reflectType = reflect.TypeOf(fieldDesc.Default().Interface())
		if fieldDesc.IsList() {
			parentMsg := dynamicpb.NewMessage(fieldDesc.ContainingMessage())
			listField := parentMsg.NewField(fieldDesc).List()
			elem := listField.NewElement().Interface()
			switch elemType := elem.(type) {
			case protoreflect.Message:
				elem = elemType.Interface()
			}
			reflectType = reflect.TypeOf(elem)
		}
	}
	if fieldDesc.IsList() {
		reflectType = reflect.SliceOf(reflectType)
	}
	var keyType, valType *FieldDescription
	if fieldDesc.IsMap() {
		keyType = NewFieldDescription(fieldDesc.MapKey())
		valType = NewFieldDescription(fieldDesc.MapValue())
	}
	wrapperDesc := wrapperField(fieldDesc)
	return &FieldDescription{
		desc:         fieldDesc,
		KeyType:      keyType,
		ValueType:    valType,
		wrapperField: wrapperDesc,
		reflectType:  reflectType,
		zeroMsg:      zeroMsg,
	}
}

// FieldDescription holds metadata related to fields declared within a type.
type FieldDescription struct {
	desc         protoreflect.FieldDescriptor
	KeyType      *FieldDescription
	ValueType    *FieldDescription
	wrapperField protoreflect.FieldDescriptor
	reflectType  reflect.Type
	zeroMsg      proto.Message
}

// CheckedType returns the type-definition used at type-check time.
func (fd *FieldDescription) CheckedType() *exprpb.Type {
	if fd.desc.IsMap() {
		return &exprpb.Type{
			TypeKind: &exprpb.Type_MapType_{
				MapType: &exprpb.Type_MapType{
					KeyType:   fd.KeyType.typeDefToType(),
					ValueType: fd.ValueType.typeDefToType(),
				},
			},
		}
	}
	if fd.desc.IsList() {
		return &exprpb.Type{
			TypeKind: &exprpb.Type_ListType_{
				ListType: &exprpb.Type_ListType{
					ElemType: fd.typeDefToType()}}}
	}
	return fd.typeDefToType()
}

// Descriptor returns the protoreflect.FieldDescriptor for this type.
func (fd *FieldDescription) Descriptor() protoreflect.FieldDescriptor {
	return fd.desc
}

// IsSet returns whether the field is set on the target value, per the proto presence conventions
// of proto2 or proto3 accordingly.
//
// The input target may either be a reflect.Value or Go struct type.
func (fd *FieldDescription) IsSet(target interface{}) bool {
	switch v := target.(type) {
	case protoreflect.Message:
		return v.Has(fd.desc)
	case proto.Message:
		return v.ProtoReflect().Has(fd.desc)
	case reflect.Value:
		return fd.IsSet(v.Interface())
	default:
		return false
	}
}

// GetFrom returns the accessor method associated with the field on the proto generated struct.
//
// If the field is not set, the proto default value is returned instead.
//
// The input target may either be a reflect.Value or proto.Message type.
func (fd *FieldDescription) GetFrom(target interface{}) (interface{}, error) {
	switch v := target.(type) {
	case proto.Message:
		fieldVal := v.ProtoReflect().Get(fd.desc).Interface()
		switch fv := fieldVal.(type) {
		case protoreflect.EnumNumber:
			return int64(fv), nil
		case protoreflect.List:
			return &List{List: fv, ElemType: fd}, nil
		case protoreflect.Map:
			return &Map{Map: fv, KeyType: fd.KeyType, ValueType: fd.ValueType}, nil
		case protoreflect.Message:
			unwrapped, _ := fd.MaybeUnwrapDynamic(fv)
			return unwrapped, nil
		default:
			return fv, nil
		}
	case reflect.Value:
		return fd.GetFrom(v.Interface())
	default:
		return nil, fmt.Errorf("unsupported field selection target: (%T)%v", target, target)
	}
}

// IsEnum returns true if the field type refers to an enum value.
func (fd *FieldDescription) IsEnum() bool {
	return fd.desc.Kind() == protoreflect.EnumKind
}

// IsMap returns true if the field is of map type.
func (fd *FieldDescription) IsMap() bool {
	return fd.desc.IsMap()
}

// IsMessage returns true if the field is of message type.
func (fd *FieldDescription) IsMessage() bool {
	return fd.desc.Kind() == protoreflect.MessageKind
}

// IsOneof returns true if the field is declared within a oneof block.
func (fd *FieldDescription) IsOneof() bool {
	return fd.desc.ContainingOneof() != nil
}

// IsList returns true if the field is a repeated value.
//
// This method will also return true for map values, so check whether the
// field is also a map.
func (fd *FieldDescription) IsList() bool {
	return fd.desc.IsList()
}

// MaybeUnwrapDynamic takes the reflected protoreflect.Message and determines whether the
// value can be unwrapped to a more primitive CEL type.
//
// This function returns the unwrapped value and 'true' on success, or the original value
// and 'false' otherwise.
func (fd *FieldDescription) MaybeUnwrapDynamic(msg protoreflect.Message) (interface{}, bool) {
	return unwrapDynamic(fd, msg)
}

// Name returns the CamelCase name of the field within the proto-based struct.
func (fd *FieldDescription) Name() string {
	return string(fd.desc.Name())
}

// ReflectType returns the Golang reflect.Type for this field.
func (fd *FieldDescription) ReflectType() reflect.Type {
	return fd.reflectType
}

// String returns a struct-like field definition string.
func (fd *FieldDescription) String() string {
	return fmt.Sprintf("%v.%s `oneof=%t`", fd.desc.ContainingMessage().FullName(), fd.Name(), fd.IsOneof())
}

// WrapperField returns the field descriptor for the 'value' of this field value when the
// field is a wrapper type.
func (fd *FieldDescription) WrapperField() protoreflect.FieldDescriptor {
	return fd.wrapperField
}

// Zero returns the zero value for the protobuf message represented by this field.
//
// If the field is not a proto.Message type, the zero value is nil.
func (fd *FieldDescription) Zero() proto.Message {
	return fd.zeroMsg
}

func (fd *FieldDescription) typeDefToType() *exprpb.Type {
	if fd.desc.Kind() == protoreflect.MessageKind {
		msgType := string(fd.desc.Message().FullName())
		if wk, found := CheckedWellKnowns[msgType]; found {
			return wk
		}
		return checkedMessageType(msgType)
	}
	if fd.desc.Kind() == protoreflect.EnumKind {
		return checkedInt
	}
	return CheckedPrimitives[fd.desc.Kind()]
}

// List wraps the protoreflect.List object with an element type FieldDescription for use in
// retrieving individual elements within CEL value data types.
type List struct {
	protoreflect.List
	ElemType *FieldDescription
}

// Map wraps the protoreflect.Map object with a key and value FieldDescription for use in
// retrieving individual elements within CEL value data types.
type Map struct {
	protoreflect.Map
	KeyType   *FieldDescription
	ValueType *FieldDescription
}

func checkedMessageType(name string) *exprpb.Type {
	return &exprpb.Type{
		TypeKind: &exprpb.Type_MessageType{MessageType: name}}
}

func checkedPrimitive(primitive exprpb.Type_PrimitiveType) *exprpb.Type {
	return &exprpb.Type{
		TypeKind: &exprpb.Type_Primitive{Primitive: primitive}}
}

func checkedWellKnown(wellKnown exprpb.Type_WellKnownType) *exprpb.Type {
	return &exprpb.Type{
		TypeKind: &exprpb.Type_WellKnown{WellKnown: wellKnown}}
}

func checkedWrap(t *exprpb.Type) *exprpb.Type {
	return &exprpb.Type{
		TypeKind: &exprpb.Type_Wrapper{Wrapper: t.GetPrimitive()}}
}

func wrapperField(desc protoreflect.FieldDescriptor) protoreflect.FieldDescriptor {
	if desc.Kind() != protoreflect.MessageKind {
		return nil
	}
	return wrapperMsg(desc.Message())
}

func wrapperMsg(msg protoreflect.MessageDescriptor) protoreflect.FieldDescriptor {
	typeName := string(msg.FullName())
	switch sanitizeProtoName(typeName) {
	case "google.protobuf.BoolValue",
		"google.protobuf.BytesValue",
		"google.protobuf.DoubleValue",
		"google.protobuf.FloatValue",
		"google.protobuf.Int32Value",
		"google.protobuf.Int64Value",
		"google.protobuf.StringValue",
		"google.protobuf.UInt32Value",
		"google.protobuf.UInt64Value":
		return msg.Fields().ByName("value")
	}
	return nil
}

func unwrap(desc description, msg proto.Message) (interface{}, bool) {
	switch v := msg.(type) {
	case *dynamicpb.Message:
		return unwrapDynamic(desc, v)
	case *dpb.Duration:
		return v.AsDuration(), true
	case *tpb.Timestamp:
		return v.AsTime(), true
	case *structpb.Value:
		switch v.GetKind().(type) {
		case *structpb.Value_BoolValue:
			return v.GetBoolValue(), true
		case *structpb.Value_ListValue:
			return v.GetListValue(), true
		case *structpb.Value_NullValue:
			return structpb.NullValue_NULL_VALUE, true
		case *structpb.Value_NumberValue:
			return v.GetNumberValue(), true
		case *structpb.Value_StringValue:
			return v.GetStringValue(), true
		case *structpb.Value_StructValue:
			return v.GetStructValue(), true
		default:
			return structpb.NullValue_NULL_VALUE, true
		}
	case *wrapperspb.BoolValue:
		return v.GetValue(), true
	case *wrapperspb.BytesValue:
		return v.GetValue(), true
	case *wrapperspb.DoubleValue:
		return v.GetValue(), true
	case *wrapperspb.FloatValue:
		return float64(v.GetValue()), true
	case *wrapperspb.Int32Value:
		return int64(v.GetValue()), true
	case *wrapperspb.Int64Value:
		return v.GetValue(), true
	case *wrapperspb.StringValue:
		return v.GetValue(), true
	case *wrapperspb.UInt32Value:
		return uint64(v.GetValue()), true
	case *wrapperspb.UInt64Value:
		return v.GetValue(), true
	}
	return msg, false
}

func unwrapDynamic(desc description, refMsg protoreflect.Message) (interface{}, bool) {
	if desc.WrapperField() != nil {
		if !refMsg.IsValid() {
			return structpb.NullValue_NULL_VALUE, true
		}
		return refMsg.Get(desc.WrapperField()).Interface(), true
	}
	msg := refMsg.Interface()
	if !refMsg.IsValid() {
		msg = desc.Zero()
	}
	typeName := string(refMsg.Descriptor().FullName())
	switch typeName {
	case "google.protobuf.Any":
		unwrapped := &anypb.Any{}
		proto.Merge(unwrapped, msg)
		return unwrapped, true
	case "google.protobuf.Duration":
		unwrapped := &dpb.Duration{}
		proto.Merge(unwrapped, msg)
		return unwrapped.AsDuration(), true
	case "google.protobuf.ListValue":
		unwrapped := &structpb.ListValue{}
		proto.Merge(unwrapped, msg)
		return unwrapped, true
	case "google.protobuf.NullValue":
		return structpb.NullValue_NULL_VALUE, true
	case "google.protobuf.Struct":
		unwrapped := &structpb.Struct{}
		proto.Merge(unwrapped, msg)
		return unwrapped, true
	case "google.protobuf.Timestamp":
		unwrapped := &tpb.Timestamp{}
		proto.Merge(unwrapped, msg)
		return unwrapped.AsTime(), true
	case "google.protobuf.Value":
		unwrapped := &structpb.Value{}
		proto.Merge(unwrapped, msg)
		return unwrap(desc, unwrapped)
	}
	return msg, false
}
