// Copyright 2023 Google LLC
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

package decls

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Kind indicates a CEL type's kind which is used to differentiate quickly between simple and complex types.
type Kind uint

const (
	// DynKind represents a dynamic type. This kind only exists at type-check time.
	DynKind Kind = iota

	// AnyKind represents a google.protobuf.Any type. This kind only exists at type-check time.
	AnyKind

	// BoolKind represents a boolean type.
	BoolKind

	// BytesKind represents a bytes type.
	BytesKind

	// DoubleKind represents a double type.
	DoubleKind

	// DurationKind represents a CEL duration type.
	DurationKind

	// IntKind represents an integer type.
	IntKind

	// ListKind represents a list type.
	ListKind

	// MapKind represents a map type.
	MapKind

	// NullTypeKind represents a null type.
	NullTypeKind

	// OpaqueKind represents an abstract type which has no accessible fields.
	OpaqueKind

	// StringKind represents a string type.
	StringKind

	// StructKind represents a structured object with typed fields.
	StructKind

	// TimestampKind represents a a CEL time type.
	TimestampKind

	// TypeKind represents the CEL type.
	TypeKind

	// TypeParamKind represents a parameterized type whose type name will be resolved at type-check time, if possible.
	TypeParamKind

	// UintKind represents a uint type.
	UintKind
)

var (
	// AnyType represents the google.protobuf.Any type.
	AnyType = &Type{
		Kind:        AnyKind,
		runtimeType: types.NewTypeValue("google.protobuf.Any"),
	}
	// BoolType represents the bool type.
	BoolType = &Type{
		Kind:        BoolKind,
		runtimeType: types.BoolType,
	}
	// BytesType represents the bytes type.
	BytesType = &Type{
		Kind:        BytesKind,
		runtimeType: types.BytesType,
	}
	// DoubleType represents the double type.
	DoubleType = &Type{
		Kind:        DoubleKind,
		runtimeType: types.DoubleType,
	}
	// DurationType represents the CEL duration type.
	DurationType = &Type{
		Kind:        DurationKind,
		runtimeType: types.DurationType,
	}
	// DynType represents a dynamic CEL type whose type will be determined at runtime from context.
	DynType = &Type{
		Kind:        DynKind,
		runtimeType: types.NewTypeValue("dyn"),
	}
	// IntType represents the int type.
	IntType = &Type{
		Kind:        IntKind,
		runtimeType: types.IntType,
	}
	// NullType represents the type of a null value.
	NullType = &Type{
		Kind:        NullTypeKind,
		runtimeType: types.NullType,
	}
	// StringType represents the string type.
	StringType = &Type{
		Kind:        StringKind,
		runtimeType: types.StringType,
	}
	// TimestampType represents the time type.
	TimestampType = &Type{
		Kind:        TimestampKind,
		runtimeType: types.TimestampType,
	}
	// TypeType represents a CEL type
	TypeType = &Type{
		Kind:        TypeKind,
		runtimeType: types.TypeType,
	}
	// UintType represents a uint type.
	UintType = &Type{
		Kind:        UintKind,
		runtimeType: types.UintType,
	}
)

// Type holds a reference to a runtime type with an optional type-checked set of type parameters.
type Type struct {
	// Kind indicates general category of the type.
	Kind Kind

	// Parameters holds the optional type-checked set of type Parameters that are used during static analysis.
	Parameters []*Type

	// runtimeType is the runtime type of the declaration.
	runtimeType ref.Type

	// isAssignableType function determines whether one type is assignable to this type.
	// A nil value for the isAssignableType function falls back to equality of kind, runtimeType, and parameters.
	isAssignableType func(other *Type) bool

	// isAssignableRuntimeType function determines whether the runtime type (with erasure) is assignable to this type.
	// A nil value for the isAssignableRuntimeType function falls back to the equality of the type or type name.
	isAssignableRuntimeType func(other ref.Val) bool
}

// IsType indicates whether two types have the same kind, type name, and parameters.
func (t *Type) IsType(other *Type) bool {
	if t.Kind != other.Kind || len(t.Parameters) != len(other.Parameters) {
		return false
	}
	if t.Kind != TypeParamKind && t.RuntimeTypeName() != other.RuntimeTypeName() {
		return false
	}
	for i, p := range t.Parameters {
		if !p.IsType(other.Parameters[i]) {
			return false
		}
	}
	return true
}

// IsAssignableType determines whether the current type is type-check assignable from the input fromType.
func (t *Type) IsAssignableType(fromType *Type) bool {
	if t.isAssignableType != nil {
		return t.isAssignableType(fromType)
	}
	return t.defaultIsAssignableType(fromType)
}

// IsAssignableRuntimeType determines whether the current type is runtime assignable from the input runtimeType.
//
// At runtime, parameterized types are erased and so a function which type-checks to support a map(string, string)
// will have a runtime assignable type of a map.
func (t *Type) IsAssignableRuntimeType(val ref.Val) bool {
	if t.isAssignableRuntimeType != nil {
		return t.isAssignableRuntimeType(val)
	}
	return t.defaultIsAssignableRuntimeType(val)
}

// RuntimeTypeName indicates the type-erased type name associated with the type.
func (t *Type) RuntimeTypeName() string {
	return t.runtimeType.TypeName()
}

// String returns a human-readable definition of the type name.
func (t *Type) String() string {
	if len(t.Parameters) == 0 {
		return t.runtimeType.TypeName()
	}
	params := make([]string, len(t.Parameters))
	for i, p := range t.Parameters {
		params[i] = p.String()
	}
	return fmt.Sprintf("%s(%s)", t.runtimeType.TypeName(), strings.Join(params, ", "))
}

// isDyn indicates whether the type is dynamic in any way.
func (t *Type) isDyn() bool {
	return t.Kind == DynKind || t.Kind == AnyKind || t.Kind == TypeParamKind
}

// defaultIsAssignableType provides the standard definition of what it means for one type to be assignable to another
// where any of the following may return a true result:
// - The from types are the same instance
// - The target type is dynamic
// - The fromType has the same kind and type name as the target type, and all parameters of the target type
//
//	are IsAssignableType() from the parameters of the fromType.
func (t *Type) defaultIsAssignableType(fromType *Type) bool {
	if t == fromType || t.isDyn() {
		return true
	}
	if t.Kind != fromType.Kind ||
		t.runtimeType.TypeName() != fromType.runtimeType.TypeName() ||
		len(t.Parameters) != len(fromType.Parameters) {
		return false
	}
	for i, tp := range t.Parameters {
		fp := fromType.Parameters[i]
		if !tp.IsAssignableType(fp) {
			return false
		}
	}
	return true
}

// defaultIsAssignableRuntimeType inspects the type and in the case of list and map elements, the key and element types
// to determine whether a ref.Val is assignable to the declared type for a function signature.
func (t *Type) defaultIsAssignableRuntimeType(val ref.Val) bool {
	valType := val.Type()
	if !(t.runtimeType == valType || t.isDyn() || t.runtimeType.TypeName() == valType.TypeName()) {
		return false
	}
	switch t.runtimeType {
	case types.ListType:
		elemType := t.Parameters[0]
		l := val.(traits.Lister)
		if l.Size() == types.IntZero {
			return true
		}
		it := l.Iterator()
		for it.HasNext() == types.True {
			elemVal := it.Next()
			return elemType.IsAssignableRuntimeType(elemVal)
		}
	case types.MapType:
		keyType := t.Parameters[0]
		elemType := t.Parameters[1]
		m := val.(traits.Mapper)
		if m.Size() == types.IntZero {
			return true
		}
		it := m.Iterator()
		for it.HasNext() == types.True {
			keyVal := it.Next()
			elemVal := m.Get(keyVal)
			return keyType.IsAssignableRuntimeType(keyVal) && elemType.IsAssignableRuntimeType(elemVal)
		}
	}
	return true
}

// ListType creates an instances of a list type value with the provided element type.
func ListType(elemType *Type) *Type {
	return &Type{
		Kind:        ListKind,
		runtimeType: types.ListType,
		Parameters:  []*Type{elemType},
	}
}

// MapType creates an instance of a map type value with the provided key and value types.
func MapType(keyType, valueType *Type) *Type {
	return &Type{
		Kind:        MapKind,
		runtimeType: types.MapType,
		Parameters:  []*Type{keyType, valueType},
	}
}

// NullableType creates an instance of a nullable type with the provided wrapped type.
//
// Note: only primitive types are supported as wrapped types.
func NullableType(wrapped *Type) *Type {
	return &Type{
		Kind:        wrapped.Kind,
		runtimeType: wrapped.runtimeType,
		Parameters:  wrapped.Parameters,
		isAssignableType: func(other *Type) bool {
			return NullType.IsAssignableType(other) || wrapped.IsAssignableType(other)
		},
		isAssignableRuntimeType: func(other ref.Val) bool {
			return NullType.IsAssignableRuntimeType(other) || wrapped.IsAssignableRuntimeType(other)
		},
	}
}

// OptionalType creates an abstract parameterized type instance corresponding to CEL's notion of optional.
func OptionalType(param *Type) *Type {
	return OpaqueType("optional", param)
}

// OpaqueType creates an abstract parameterized type with a given name.
func OpaqueType(name string, params ...*Type) *Type {
	return &Type{
		Kind:        OpaqueKind,
		runtimeType: types.NewTypeValue(name),
		Parameters:  params,
	}
}

// ObjectType creates a type references to an externally defined type, e.g. a protobuf message type.
func ObjectType(typeName string) *Type {
	// Function sanitizes object types on the fly
	if wkt, found := checkedWellKnowns[typeName]; found {
		return wkt
	}
	return &Type{
		Kind:        StructKind,
		runtimeType: types.NewObjectTypeValue(typeName),
	}
}

// TypeParamType creates a parameterized type instance.
func TypeParamType(paramName string) *Type {
	return &Type{
		Kind:        TypeParamKind,
		runtimeType: types.NewTypeValue(paramName),
	}
}

// TypeTypeWithParam creates a type with a type parameter.
// Used for type-checking purposes, but equivalent to TypeType otherwise.
func TypeTypeWithParam(param *Type) *Type {
	return &Type{
		Kind:        TypeKind,
		Parameters:  []*Type{param},
		runtimeType: types.TypeType,
	}
}

var (
	checkedWellKnowns = map[string]*Type{
		// Wrapper types.
		"google.protobuf.BoolValue":   NullableType(BoolType),
		"google.protobuf.BytesValue":  NullableType(BytesType),
		"google.protobuf.DoubleValue": NullableType(DoubleType),
		"google.protobuf.FloatValue":  NullableType(DoubleType),
		"google.protobuf.Int64Value":  NullableType(IntType),
		"google.protobuf.Int32Value":  NullableType(IntType),
		"google.protobuf.UInt64Value": NullableType(UintType),
		"google.protobuf.UInt32Value": NullableType(UintType),
		"google.protobuf.StringValue": NullableType(StringType),
		// Well-known types.
		"google.protobuf.Any":       AnyType,
		"google.protobuf.Duration":  DurationType,
		"google.protobuf.Timestamp": TimestampType,
		// Json types.
		"google.protobuf.ListValue": ListType(DynType),
		"google.protobuf.NullValue": NullType,
		"google.protobuf.Struct":    MapType(StringType, DynType),
		"google.protobuf.Value":     DynType,
	}
)
