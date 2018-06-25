// Package decls provides helpers for creating variable and function declarations.
package decls

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-spec/proto/checked/v1/checked"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

var (
	// Error type used to communicate issues during type-checking.
	Error = &checked.Type{
		TypeKind: &checked.Type_Error{
			Error: &empty.Empty{}}}

	// Dyn is a top-type used to represent any value.
	Dyn = &checked.Type{
		TypeKind: &checked.Type_Dyn{
			Dyn: &empty.Empty{}}}

	// Commonly used types.
	Bool   = NewPrimitiveType(checked.Type_BOOL)
	Bytes  = NewPrimitiveType(checked.Type_BYTES)
	Double = NewPrimitiveType(checked.Type_DOUBLE)
	Int    = NewPrimitiveType(checked.Type_INT64)
	Null   = &checked.Type{
		TypeKind: &checked.Type_Null{
			Null: structpb.NullValue_NULL_VALUE}}
	String = NewPrimitiveType(checked.Type_STRING)
	Uint   = NewPrimitiveType(checked.Type_UINT64)

	// Well-known types.
	// TODO: Replace with an abstract type registry.
	Any       = NewWellKnownType(checked.Type_ANY)
	Duration  = NewWellKnownType(checked.Type_DURATION)
	Timestamp = NewWellKnownType(checked.Type_TIMESTAMP)
)

// NewFunctionType creates a function invocation contract, typically only used
// by type-checking steps after overload resolution.
func NewFunctionType(resultType *checked.Type,
	argTypes ...*checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Function{
			Function: &checked.Type_FunctionType{
				ResultType: resultType,
				ArgTypes:   argTypes}}}
}

// NewFunction creates a named function declaration with one or more overloads.
func NewFunction(name string,
	overloads ...*checked.Decl_FunctionDecl_Overload) *checked.Decl {
	return &checked.Decl{
		Name: name,
		DeclKind: &checked.Decl_Function{
			Function: &checked.Decl_FunctionDecl{
				Overloads: overloads}}}
}

// NewIdent creates a named identifier declaration with an optional literal
// value.
//
// Literal values are typically only associated with enum identifiers.
func NewIdent(name string, t *checked.Type, v *expr.Literal) *checked.Decl {
	return &checked.Decl{
		Name: name,
		DeclKind: &checked.Decl_Ident{
			Ident: &checked.Decl_IdentDecl{
				Type:  t,
				Value: v}}}
}

// NewInstanceOverload creates a instance function overload contract.
func NewInstanceOverload(id string, argTypes []*checked.Type,
	resultType *checked.Type) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: true}
}

// NewListType generates a new list with elements of a certain type.
func NewListType(elem *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_ListType_{
			ListType: &checked.Type_ListType{
				ElemType: elem}}}
}

// NewMapType generates a new map with typed keys and values.
func NewMapType(key *checked.Type, value *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_MapType_{
			MapType: &checked.Type_MapType{
				KeyType:   key,
				ValueType: value}}}
}

// NewObjectType creates an object type for a qualified type name.
func NewObjectType(typeName string) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_MessageType{
			MessageType: typeName}}
}

// NewOverload creates a function overload declaration which contains a unique
// overload id as well as the expected argument and result types. Overloads
// must be aggregated within a Function declaration.
func NewOverload(id string, argTypes []*checked.Type,
	resultType *checked.Type) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: false}
}

// NewParameterizedOverload creates a parametric function instance overload
// type.
func NewParameterizedInstanceOverload(id string,
	argTypes []*checked.Type,
	resultType *checked.Type,
	typeParams []string) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: true}
}

// NewParameterizedOverload creates a parametric function overload type.
func NewParameterizedOverload(id string,
	argTypes []*checked.Type,
	resultType *checked.Type,
	typeParams []string) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: false}
}

// NewPrimitiveType creates a type for a primitive value. See the var declarations
// for Int, Uint, etc.
func NewPrimitiveType(primitive checked.Type_PrimitiveType) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Primitive{
			Primitive: primitive}}
}

// NewTypeType creates a new type designating a type.
func NewTypeType(nested *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Type{
			Type: nested}}
}

// NewTypeParamType creates a type corresponding to a named, contextual parameter.
func NewTypeParamType(name string) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_TypeParam{
			TypeParam: name}}
}

// NewWellKnownType creates a type corresponding to a protobuf well-known type
// value.
func NewWellKnownType(wellKnown checked.Type_WellKnownType) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_WellKnown{
			WellKnown: wellKnown}}
}

// NewWrapperType creates a wrapped primitive type instance. Wrapped types
// are roughly equivalent to a nullable, or optionally valued type.
func NewWrapperType(wrapped *checked.Type) *checked.Type {
	primitive := wrapped.GetPrimitive()
	if primitive == checked.Type_PRIMITIVE_TYPE_UNSPECIFIED {
		// TODO: return an error
		panic("Wrapped type must be a primitive")
	}
	return &checked.Type{
		TypeKind: &checked.Type_Wrapper{
			Wrapper: primitive}}
}
