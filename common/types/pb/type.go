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
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// NewTypeDescription produces a TypeDescription value for the fully-qualified proto type name
// with a given descriptor.
//
// The type description creation method also expects the type to be marked clearly as a proto2 or
// proto3 type, and accepts a typeResolver reference for resolving field TypeDescription during
// lazily initialization of the type which is done atomically.
func NewTypeDescription(typeName string, desc protoreflect.MessageDescriptor) *TypeDescription {
	msgType := dynamicpb.NewMessageType(desc)
	fieldMap := map[string]*FieldDescription{}
	fields := desc.Fields()
	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		fieldMap[string(f.Name())] = NewFieldDescription(f)
	}
	return &TypeDescription{
		typeName:    typeName,
		desc:        desc,
		msgType:     msgType,
		fieldMap:    fieldMap,
		reflectType: reflect.TypeOf(msgType.Zero().Interface()),
	}
}

// TypeDescription is a collection of type metadata relevant to expression
// checking and evaluation.
type TypeDescription struct {
	typeName    string
	desc        protoreflect.MessageDescriptor
	msgType     protoreflect.MessageType
	fieldMap    map[string]*FieldDescription
	reflectType reflect.Type
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

// NewFieldDescription creates a new field description from a protoreflect.FieldDescriptor.
func NewFieldDescription(fieldDesc protoreflect.FieldDescriptor) *FieldDescription {
	var reflectType reflect.Type
	switch fieldDesc.Kind() {
	case protoreflect.EnumKind:
		reflectType = reflect.TypeOf(protoreflect.EnumNumber(0))
	case protoreflect.MessageKind:

		reflectType = reflect.TypeOf(dynamicpb.NewMessageType(fieldDesc.Message()).Zero().Interface())
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
	return &FieldDescription{
		descriptor:  fieldDesc,
		KeyType:     keyType,
		ValueType:   valType,
		isWrapper:   isWrapperType(fieldDesc),
		reflectType: reflectType,
	}
}

// FieldDescription holds metadata related to fields declared within a type.
type FieldDescription struct {
	descriptor  protoreflect.FieldDescriptor
	KeyType     *FieldDescription
	ValueType   *FieldDescription
	isWrapper   bool
	reflectType reflect.Type
}

// CheckedType returns the type-definition used at type-check time.
func (fd *FieldDescription) CheckedType() *exprpb.Type {
	if fd.descriptor.IsMap() {
		return &exprpb.Type{
			TypeKind: &exprpb.Type_MapType_{
				MapType: &exprpb.Type_MapType{
					KeyType:   fd.KeyType.typeDefToType(),
					ValueType: fd.ValueType.typeDefToType(),
				},
			},
		}
	}
	if fd.descriptor.IsList() {
		return &exprpb.Type{
			TypeKind: &exprpb.Type_ListType_{
				ListType: &exprpb.Type_ListType{
					ElemType: fd.typeDefToType()}}}
	}
	return fd.typeDefToType()
}

// Descriptor returns the protoreflect.FieldDescriptor for this type.
func (fd *FieldDescription) Descriptor() protoreflect.FieldDescriptor {
	return fd.descriptor
}

// ReflectType returns the Golang reflect.Type for this field.
func (fd *FieldDescription) ReflectType() reflect.Type {
	return fd.reflectType
}

// IsSet returns whether the field is set on the target value, per the proto presence conventions
// of proto2 or proto3 accordingly.
//
// The input target may either be a reflect.Value or Go struct type.
func (fd *FieldDescription) IsSet(target interface{}) bool {
	switch v := target.(type) {
	case protoreflect.Message:
		return v.Has(fd.descriptor)
	case proto.Message:
		return v.ProtoReflect().Has(fd.descriptor)
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
		if fd.IsWrapper() {
			if !fd.IsSet(target) {
				return structpb.NullValue_NULL_VALUE, nil
			}
		}
		fieldVal := v.ProtoReflect().Get(fd.descriptor).Interface()
		switch fv := fieldVal.(type) {
		case protoreflect.EnumNumber:
			return int64(fv), nil
		case protoreflect.List:
			return List{List: fv, ElemType: fd}, nil
		case protoreflect.Map:
			return Map{Map: fv, KeyType: fd.KeyType, ValueType: fd.ValueType}, nil
		case protoreflect.Message:
			return fv.Interface(), nil
		default:
			return fv, nil
		}
	case reflect.Value:
		return fd.GetFrom(v.Interface())
	default:
		return nil, fmt.Errorf("unsupported field selection target: (%T)%v", target, target)
	}
}

// Index returns the field index within a reflected value.
func (fd *FieldDescription) Index() int {
	return fd.descriptor.Index()
}

// IsEnum returns true if the field type refers to an enum value.
func (fd *FieldDescription) IsEnum() bool {
	return fd.descriptor.Kind() == protoreflect.EnumKind
}

// IsMap returns true if the field is of map type.
func (fd *FieldDescription) IsMap() bool {
	return fd.descriptor.IsMap()
}

// IsMessage returns true if the field is of message type.
func (fd *FieldDescription) IsMessage() bool {
	return fd.descriptor.Kind() == protoreflect.MessageKind
}

// IsOneof returns true if the field is declared within a oneof block.
func (fd *FieldDescription) IsOneof() bool {
	return fd.descriptor.ContainingOneof() != nil
}

// IsRepeated returns true if the field is a repeated value.
//
// This method will also return true for map values, so check whether the
// field is also a map.
func (fd *FieldDescription) IsList() bool {
	return fd.descriptor.IsList()
}

// IsWrapper returns true if the field type is a primitive wrapper type.
func (fd *FieldDescription) IsWrapper() bool {
	return fd.isWrapper
}

// Name returns the CamelCase name of the field within the proto-based struct.
func (fd *FieldDescription) Name() string {
	return string(fd.descriptor.Name())
}

// String returns a struct-like field definition string.
func (fd *FieldDescription) String() string {
	return fmt.Sprintf("%v.%s `oneof=%t`", fd.descriptor.ContainingMessage().FullName(), fd.Name(), fd.IsOneof())
}

func (fd *FieldDescription) typeDefToType() *exprpb.Type {
	if fd.descriptor.Kind() == protoreflect.MessageKind {
		msgType := string(fd.descriptor.Message().FullName())
		if wk, found := CheckedWellKnowns[msgType]; found {
			return wk
		}
		return checkedMessageType(msgType)
	}
	if fd.descriptor.Kind() == protoreflect.EnumKind {
		return checkedInt
	}
	return CheckedPrimitives[fd.descriptor.Kind()]
}

type List struct {
	protoreflect.List
	ElemType *FieldDescription
}

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

func isWrapperType(desc protoreflect.FieldDescriptor) bool {
	if desc.Kind() != protoreflect.MessageKind {
		return false
	}
	typeName := string(desc.ContainingMessage().FullName())
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
		return true
	}
	return false
}
