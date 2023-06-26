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

package types

import (
	"fmt"
	"reflect"
	"strings"

	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Kind indicates a CEL type's kind which is used to differentiate quickly between simple
// and complex types.
type Kind uint

const (
	// DynKind represents a dynamic type. This kind only exists at type-check time.
	DynKind Kind = iota + 1

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

	// ErrorKind represents a CEL error type.
	ErrorKind

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

	// UnknownKind represents an unknown value type.
	UnknownKind
)

var (
	// AnyType represents the google.protobuf.Any type.
	AnyType = &Type{
		Kind:            AnyKind,
		runtimeTypeName: "google.protobuf.Any",
		traitMask: traits.FieldTesterType |
			traits.IndexerType,
	}
	// BoolType represents the bool type.
	BoolType = &Type{
		Kind:            BoolKind,
		runtimeTypeName: "bool",
		traitMask: traits.ComparerType |
			traits.NegatorType,
	}
	// BytesType represents the bytes type.
	BytesType = &Type{
		Kind:            BytesKind,
		runtimeTypeName: "bytes",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.SizerType,
	}
	// DoubleType represents the double type.
	DoubleType = &Type{
		Kind:            DoubleKind,
		runtimeTypeName: "double",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.DividerType |
			traits.MultiplierType |
			traits.NegatorType |
			traits.SubtractorType,
	}
	// DurationType represents the CEL duration type.
	DurationType = &Type{
		Kind:            DurationKind,
		runtimeTypeName: "google.protobuf.Duration",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.NegatorType |
			traits.ReceiverType |
			traits.SubtractorType,
	}
	// DynType represents a dynamic CEL type whose type will be determined at runtime from context.
	DynType = &Type{
		Kind:            DynKind,
		runtimeTypeName: "dyn",
	}
	// ErrorType represents a CEL error value.
	ErrorType = &Type{
		Kind:            ErrorKind,
		runtimeTypeName: "error",
	}
	// IntType represents the int type.
	IntType = &Type{
		Kind:            IntKind,
		runtimeTypeName: "int",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.DividerType |
			traits.ModderType |
			traits.MultiplierType |
			traits.NegatorType |
			traits.SubtractorType,
	}
	// ListType represents the runtime list type.
	ListType = NewListType(nil)
	// MapType represents the runtime map type.
	MapType = NewMapType(nil, nil)
	// NullType represents the type of a null value.
	NullType = &Type{
		Kind:            NullTypeKind,
		runtimeTypeName: "null_type",
	}
	// StringType represents the string type.
	StringType = &Type{
		Kind:            StringKind,
		runtimeTypeName: "string",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.MatcherType |
			traits.ReceiverType |
			traits.SizerType,
	}
	// TimestampType represents the time type.
	TimestampType = &Type{
		Kind:            TimestampKind,
		runtimeTypeName: "google.protobuf.Timestamp",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.ReceiverType |
			traits.SubtractorType,
	}
	// TypeType represents a CEL type
	TypeType = &Type{
		Kind:            TypeKind,
		runtimeTypeName: "type",
	}
	// UintType represents a uint type.
	UintType = &Type{
		Kind:            UintKind,
		runtimeTypeName: "uint",
		traitMask: traits.AdderType |
			traits.ComparerType |
			traits.DividerType |
			traits.ModderType |
			traits.MultiplierType |
			traits.SubtractorType,
	}
	// UnknownType represents an unknown value type.
	UnknownType = &Type{
		Kind:            UnknownKind,
		runtimeTypeName: "unknown",
	}
)

var _ ref.Type = &Type{}
var _ ref.Val = &Type{}

// Type holds a reference to a runtime type with an optional type-checked set of type parameters.
type Type struct {
	// Kind indicates general category of the type.
	Kind Kind

	// Parameters holds the optional type-checked set of type Parameters that are used during static analysis.
	Parameters []*Type

	// runtimeTypeName indicates the runtime type name of the type.
	runtimeTypeName string

	// isAssignableType function determines whether one type is assignable to this type.
	// A nil value for the isAssignableType function falls back to equality of kind, runtimeType, and parameters.
	isAssignableType func(other *Type) bool

	// isAssignableRuntimeType function determines whether the runtime type (with erasure) is assignable to this type.
	// A nil value for the isAssignableRuntimeType function falls back to the equality of the type or type name.
	isAssignableRuntimeType func(other ref.Val) bool

	// traitMask is a mask of flags which indicate the capabilities of the type.
	traitMask int
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (t *Type) ConvertToNative(typeDesc reflect.Type) (any, error) {
	return nil, fmt.Errorf("type conversion not supported for 'type'")
}

// ConvertToType implements ref.Val.ConvertToType.
func (t *Type) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case TypeType:
		return TypeType
	case StringType:
		return String(t.TypeName())
	}
	return NewErr("type conversion error from '%s' to '%s'", TypeType, typeVal)
}

// Equal indicates whether two types have the same runtime type name.
//
// The name Equal is a bit of a misnomer, but for historical reasons, this is the
// runtime behavior. For a more accurate definition see IsType().
func (t *Type) Equal(other ref.Val) ref.Val {
	otherType, ok := other.(ref.Type)
	return Bool(ok && t.TypeName() == otherType.TypeName())
}

// HasTrait implements the ref.Type interface method.
func (t *Type) HasTrait(trait int) bool {
	return trait&t.traitMask == trait
}

// IsExactType indicates whether the two types are exactly the same. This check also verifies type parameter type names.
func (t *Type) IsExactType(other *Type) bool {
	return t.isTypeInternal(other, true)
}

// IsEquivalentType indicates whether two types are equivalent. This check ignores type parameter type names.
func (t *Type) IsEquivalentType(other *Type) bool {
	return t.isTypeInternal(other, false)
}

// isTypeInternal checks whether the two types are equivalent or exactly the same based on the checkTypeParamName flag.
func (t *Type) isTypeInternal(other *Type, checkTypeParamName bool) bool {
	if t == other {
		return true
	}
	if t.Kind != other.Kind || len(t.Parameters) != len(other.Parameters) {
		return false
	}
	if (checkTypeParamName || t.Kind != TypeParamKind) && t.TypeName() != other.TypeName() {
		return false
	}
	for i, p := range t.Parameters {
		if !p.isTypeInternal(other.Parameters[i], checkTypeParamName) {
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

// DeclaredTypeName indicates the fully qualified and parameterized type-check type name.
func (t *Type) DeclaredTypeName() string {
	// if the type itself is neither null, nor dyn, but is assignable to null, then it's a wrapper type.
	if t.Kind != NullTypeKind && !t.isDyn() && t.IsAssignableType(NullType) {
		return fmt.Sprintf("wrapper(%s)", t.TypeName())
	}
	return t.TypeName()
}

// Type implements the ref.Val interface method.
func (t *Type) Type() ref.Type {
	return TypeType
}

// Value implements the ref.Val interface method.
func (t *Type) Value() any {
	return t.TypeName()
}

// TypeName returns the type-erased fully qualified runtime type name.
//
// TypeName implements the ref.Type interface method.
func (t *Type) TypeName() string {
	return t.runtimeTypeName
}

// String returns a human-readable definition of the type name.
func (t *Type) String() string {
	if len(t.Parameters) == 0 {
		return t.DeclaredTypeName()
	}
	params := make([]string, len(t.Parameters))
	for i, p := range t.Parameters {
		params[i] = p.String()
	}
	return fmt.Sprintf("%s(%s)", t.DeclaredTypeName(), strings.Join(params, ", "))
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
		t.TypeName() != fromType.TypeName() ||
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
	// If the current type and value type don't agree, then return
	if !(t.isDyn() || t.TypeName() == valType.TypeName()) {
		return false
	}
	switch t.Kind {
	case ListKind:
		elemType := t.Parameters[0]
		l := val.(traits.Lister)
		if l.Size() == IntZero {
			return true
		}
		it := l.Iterator()
		elemVal := it.Next()
		return elemType.IsAssignableRuntimeType(elemVal)
	case MapKind:
		keyType := t.Parameters[0]
		elemType := t.Parameters[1]
		m := val.(traits.Mapper)
		if m.Size() == IntZero {
			return true
		}
		it := m.Iterator()
		keyVal := it.Next()
		elemVal := m.Get(keyVal)
		return keyType.IsAssignableRuntimeType(keyVal) && elemType.IsAssignableRuntimeType(elemVal)
	}
	return true
}

// NewListType creates an instances of a list type value with the provided element type.
func NewListType(elemType *Type) *Type {
	t := &Type{
		Kind:            ListKind,
		Parameters:      []*Type{},
		runtimeTypeName: "list",
		traitMask: traits.AdderType |
			traits.ContainerType |
			traits.IndexerType |
			traits.IterableType |
			traits.SizerType,
	}
	if elemType != nil {
		t.Parameters = append(t.Parameters, elemType)
	}
	return t
}

// NewMapType creates an instance of a map type value with the provided key and value types.
func NewMapType(keyType, valueType *Type) *Type {
	t := &Type{
		Kind:            MapKind,
		Parameters:      []*Type{},
		runtimeTypeName: "map",
		traitMask: traits.ContainerType |
			traits.IndexerType |
			traits.IterableType |
			traits.SizerType,
	}
	if keyType != nil && valueType != nil {
		t.Parameters = append(t.Parameters, keyType, valueType)
	}
	return t
}

// NewNullableType creates an instance of a nullable type with the provided wrapped type.
//
// Note: only primitive types are supported as wrapped types.
func NewNullableType(wrapped *Type) *Type {
	return &Type{
		Kind:            wrapped.Kind,
		Parameters:      wrapped.Parameters,
		runtimeTypeName: wrapped.runtimeTypeName,
		traitMask:       wrapped.traitMask,
		isAssignableType: func(other *Type) bool {
			return NullType.IsAssignableType(other) || wrapped.IsAssignableType(other)
		},
		isAssignableRuntimeType: func(other ref.Val) bool {
			return NullType.IsAssignableRuntimeType(other) || wrapped.IsAssignableRuntimeType(other)
		},
	}
}

// NewOptionalType creates an abstract parameterized type instance corresponding to CEL's notion of optional.
func NewOptionalType(param *Type) *Type {
	return NewOpaqueType("optional", param)
}

// NewOpaqueType creates an abstract parameterized type with a given name.
func NewOpaqueType(name string, params ...*Type) *Type {
	return &Type{
		Kind:            OpaqueKind,
		Parameters:      params,
		runtimeTypeName: name,
	}
}

// NewObjectType creates a type reference to an externally defined type, e.g. a protobuf message type.
func NewObjectType(typeName string) *Type {
	// Function sanitizes object types on the fly
	if wkt, found := checkedWellKnowns[typeName]; found {
		return wkt
	}
	return &Type{
		Kind:            StructKind,
		Parameters:      []*Type{},
		runtimeTypeName: typeName,
		traitMask:       traits.FieldTesterType | traits.IndexerType,
	}
}

// NewObjectTypeValue creates a type reference to an externally defined type.
//
// Deprecated: use cel.ObjectType(typeName)
func NewObjectTypeValue(typeName string) *Type {
	return NewObjectType(typeName)
}

// NewTypeValue creates an opaque type which has a set of optional type traits as defined in
// the common/types/traits package.
//
// Deprecated: use cel.OpaqueType(typeName)
func NewTypeValue(typeName string, traits ...int) *Type {
	traitMask := 0
	for _, trait := range traits {
		traitMask |= trait
	}
	return &Type{
		Kind:            OpaqueKind,
		Parameters:      []*Type{},
		runtimeTypeName: typeName,
		traitMask:       traitMask,
	}
}

// NewTypeParamType creates a parameterized type instance.
func NewTypeParamType(paramName string) *Type {
	return &Type{
		Kind:            TypeParamKind,
		runtimeTypeName: paramName,
	}
}

// NewTypeTypeWithParam creates a type with a type parameter.
// Used for type-checking purposes, but equivalent to TypeType otherwise.
func NewTypeTypeWithParam(param *Type) *Type {
	return &Type{
		Kind:            TypeKind,
		runtimeTypeName: "type",
		Parameters:      []*Type{param},
	}
}

// TypeToExprType converts a CEL-native type representation to a protobuf CEL Type representation.
func TypeToExprType(t *Type) (*exprpb.Type, error) {
	switch t.Kind {
	case AnyKind:
		return chkdecls.Any, nil
	case BoolKind:
		return maybeWrapper(t, chkdecls.Bool), nil
	case BytesKind:
		return maybeWrapper(t, chkdecls.Bytes), nil
	case DoubleKind:
		return maybeWrapper(t, chkdecls.Double), nil
	case DurationKind:
		return chkdecls.Duration, nil
	case DynKind:
		return chkdecls.Dyn, nil
	case ErrorKind:
		return chkdecls.Error, nil
	case IntKind:
		return maybeWrapper(t, chkdecls.Int), nil
	case ListKind:
		if len(t.Parameters) != 1 {
			return nil, fmt.Errorf("invalid list, got %d parameters, wanted one", len(t.Parameters))
		}
		et, err := TypeToExprType(t.Parameters[0])
		if err != nil {
			return nil, err
		}
		return chkdecls.NewListType(et), nil
	case MapKind:
		if len(t.Parameters) != 2 {
			return nil, fmt.Errorf("invalid map, got %d parameters, wanted two", len(t.Parameters))
		}
		kt, err := TypeToExprType(t.Parameters[0])
		if err != nil {
			return nil, err
		}
		vt, err := TypeToExprType(t.Parameters[1])
		if err != nil {
			return nil, err
		}
		return chkdecls.NewMapType(kt, vt), nil
	case NullTypeKind:
		return chkdecls.Null, nil
	case OpaqueKind:
		params := make([]*exprpb.Type, len(t.Parameters))
		for i, p := range t.Parameters {
			pt, err := TypeToExprType(p)
			if err != nil {
				return nil, err
			}
			params[i] = pt
		}
		return chkdecls.NewAbstractType(t.TypeName(), params...), nil
	case StringKind:
		return maybeWrapper(t, chkdecls.String), nil
	case StructKind:
		return chkdecls.NewObjectType(t.TypeName()), nil
	case TimestampKind:
		return chkdecls.Timestamp, nil
	case TypeParamKind:
		return chkdecls.NewTypeParamType(t.TypeName()), nil
	case TypeKind:
		if len(t.Parameters) == 1 {
			p, err := TypeToExprType(t.Parameters[0])
			if err != nil {
				return nil, err
			}
			return chkdecls.NewTypeType(p), nil
		}
		return chkdecls.NewTypeType(nil), nil
	case UintKind:
		return maybeWrapper(t, chkdecls.Uint), nil
	}
	return nil, fmt.Errorf("missing type conversion to proto: %v", t)
}

// ExprTypeToType converts a protobuf CEL type representation to a CEL-native type representation.
func ExprTypeToType(t *exprpb.Type) (*Type, error) {
	switch t.GetTypeKind().(type) {
	case *exprpb.Type_Dyn:
		return DynType, nil
	case *exprpb.Type_AbstractType_:
		paramTypes := make([]*Type, len(t.GetAbstractType().GetParameterTypes()))
		for i, p := range t.GetAbstractType().GetParameterTypes() {
			pt, err := ExprTypeToType(p)
			if err != nil {
				return nil, err
			}
			paramTypes[i] = pt
		}
		return NewOpaqueType(t.GetAbstractType().GetName(), paramTypes...), nil
	case *exprpb.Type_ListType_:
		et, err := ExprTypeToType(t.GetListType().GetElemType())
		if err != nil {
			return nil, err
		}
		return NewListType(et), nil
	case *exprpb.Type_MapType_:
		kt, err := ExprTypeToType(t.GetMapType().GetKeyType())
		if err != nil {
			return nil, err
		}
		vt, err := ExprTypeToType(t.GetMapType().GetValueType())
		if err != nil {
			return nil, err
		}
		return NewMapType(kt, vt), nil
	case *exprpb.Type_MessageType:
		return NewObjectType(t.GetMessageType()), nil
	case *exprpb.Type_Null:
		return NullType, nil
	case *exprpb.Type_Primitive:
		switch t.GetPrimitive() {
		case exprpb.Type_BOOL:
			return BoolType, nil
		case exprpb.Type_BYTES:
			return BytesType, nil
		case exprpb.Type_DOUBLE:
			return DoubleType, nil
		case exprpb.Type_INT64:
			return IntType, nil
		case exprpb.Type_STRING:
			return StringType, nil
		case exprpb.Type_UINT64:
			return UintType, nil
		default:
			return nil, fmt.Errorf("unsupported primitive type: %v", t)
		}
	case *exprpb.Type_TypeParam:
		return NewTypeParamType(t.GetTypeParam()), nil
	case *exprpb.Type_Type:
		if t.GetType().GetTypeKind() != nil {
			p, err := ExprTypeToType(t.GetType())
			if err != nil {
				return nil, err
			}
			return NewTypeTypeWithParam(p), nil
		}
		return TypeType, nil
	case *exprpb.Type_WellKnown:
		switch t.GetWellKnown() {
		case exprpb.Type_ANY:
			return AnyType, nil
		case exprpb.Type_DURATION:
			return DurationType, nil
		case exprpb.Type_TIMESTAMP:
			return TimestampType, nil
		default:
			return nil, fmt.Errorf("unsupported well-known type: %v", t)
		}
	case *exprpb.Type_Wrapper:
		t, err := ExprTypeToType(&exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: t.GetWrapper()}})
		if err != nil {
			return nil, err
		}
		return NewNullableType(t), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", t)
	}
}

func maybeWrapper(t *Type, pbType *exprpb.Type) *exprpb.Type {
	if t.IsAssignableType(NullType) {
		return chkdecls.NewWrapperType(pbType)
	}
	return pbType
}

var (
	checkedWellKnowns = map[string]*Type{
		// Wrapper types.
		"google.protobuf.BoolValue":   NewNullableType(BoolType),
		"google.protobuf.BytesValue":  NewNullableType(BytesType),
		"google.protobuf.DoubleValue": NewNullableType(DoubleType),
		"google.protobuf.FloatValue":  NewNullableType(DoubleType),
		"google.protobuf.Int64Value":  NewNullableType(IntType),
		"google.protobuf.Int32Value":  NewNullableType(IntType),
		"google.protobuf.UInt64Value": NewNullableType(UintType),
		"google.protobuf.UInt32Value": NewNullableType(UintType),
		"google.protobuf.StringValue": NewNullableType(StringType),
		// Well-known types.
		"google.protobuf.Any":       AnyType,
		"google.protobuf.Duration":  DurationType,
		"google.protobuf.Timestamp": TimestampType,
		// Json types.
		"google.protobuf.ListValue": NewListType(DynType),
		"google.protobuf.NullValue": NullType,
		"google.protobuf.Struct":    NewMapType(StringType, DynType),
		"google.protobuf.Value":     DynType,
	}
)
