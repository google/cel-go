// Package decls provides helpers for creating variable and function declarations.
package decls

import (
	emptypb "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
)

var (
	// Error type used to communicate issues during type-checking.
	Error = &checkedpb.Type{
		TypeKind: &checkedpb.Type_Error{
			Error: &emptypb.Empty{}}}

	// Dyn is a top-type used to represent any value.
	Dyn = &checkedpb.Type{
		TypeKind: &checkedpb.Type_Dyn{
			Dyn: &emptypb.Empty{}}}

	// Commonly used types.
	Bool   = NewPrimitiveType(checkedpb.Type_BOOL)
	Bytes  = NewPrimitiveType(checkedpb.Type_BYTES)
	Double = NewPrimitiveType(checkedpb.Type_DOUBLE)
	Int    = NewPrimitiveType(checkedpb.Type_INT64)
	Null   = &checkedpb.Type{
		TypeKind: &checkedpb.Type_Null{
			Null: structpb.NullValue_NULL_VALUE}}
	String = NewPrimitiveType(checkedpb.Type_STRING)
	Uint   = NewPrimitiveType(checkedpb.Type_UINT64)

	// Well-known types.
	// TODO: Replace with an abstract type registry.
	Any       = NewWellKnownType(checkedpb.Type_ANY)
	Duration  = NewWellKnownType(checkedpb.Type_DURATION)
	Timestamp = NewWellKnownType(checkedpb.Type_TIMESTAMP)
)

// NewFunctionType creates a function invocation contract, typically only used
// by type-checking steps after overload resolution.
func NewFunctionType(resultType *checkedpb.Type,
	argTypes ...*checkedpb.Type) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_Function{
			Function: &checkedpb.Type_FunctionType{
				ResultType: resultType,
				ArgTypes:   argTypes}}}
}

// NewFunction creates a named function declaration with one or more overloads.
func NewFunction(name string,
	overloads ...*checkedpb.Decl_FunctionDecl_Overload) *checkedpb.Decl {
	return &checkedpb.Decl{
		Name: name,
		DeclKind: &checkedpb.Decl_Function{
			Function: &checkedpb.Decl_FunctionDecl{
				Overloads: overloads}}}
}

// NewIdent creates a named identifier declaration with an optional literal
// value.
//
// Literal values are typically only associated with enum identifiers.
func NewIdent(name string, t *checkedpb.Type, v *exprpb.Literal) *checkedpb.Decl {
	return &checkedpb.Decl{
		Name: name,
		DeclKind: &checkedpb.Decl_Ident{
			Ident: &checkedpb.Decl_IdentDecl{
				Type:  t,
				Value: v}}}
}

// NewInstanceOverload creates a instance function overload contract.
func NewInstanceOverload(id string, argTypes []*checkedpb.Type,
	resultType *checkedpb.Type) *checkedpb.Decl_FunctionDecl_Overload {
	return &checkedpb.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: true}
}

// NewListType generates a new list with elements of a certain type.
func NewListType(elem *checkedpb.Type) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_ListType_{
			ListType: &checkedpb.Type_ListType{
				ElemType: elem}}}
}

// NewMapType generates a new map with typed keys and values.
func NewMapType(key *checkedpb.Type, value *checkedpb.Type) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_MapType_{
			MapType: &checkedpb.Type_MapType{
				KeyType:   key,
				ValueType: value}}}
}

// NewObjectType creates an object type for a qualified type name.
func NewObjectType(typeName string) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_MessageType{
			MessageType: typeName}}
}

// NewOverload creates a function overload declaration which contains a unique
// overload id as well as the expected argument and result types. Overloads
// must be aggregated within a Function declaration.
func NewOverload(id string, argTypes []*checkedpb.Type,
	resultType *checkedpb.Type) *checkedpb.Decl_FunctionDecl_Overload {
	return &checkedpb.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: false}
}

// NewParameterizedOverload creates a parametric function instance overload
// type.
func NewParameterizedInstanceOverload(id string,
	argTypes []*checkedpb.Type,
	resultType *checkedpb.Type,
	typeParams []string) *checkedpb.Decl_FunctionDecl_Overload {
	return &checkedpb.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: true}
}

// NewParameterizedOverload creates a parametric function overload type.
func NewParameterizedOverload(id string,
	argTypes []*checkedpb.Type,
	resultType *checkedpb.Type,
	typeParams []string) *checkedpb.Decl_FunctionDecl_Overload {
	return &checkedpb.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: false}
}

// NewPrimitiveType creates a type for a primitive value. See the var declarations
// for Int, Uint, etc.
func NewPrimitiveType(primitive checkedpb.Type_PrimitiveType) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_Primitive{
			Primitive: primitive}}
}

// NewTypeType creates a new type designating a type.
func NewTypeType(nested *checkedpb.Type) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_Type{
			Type: nested}}
}

// NewTypeParamType creates a type corresponding to a named, contextual parameter.
func NewTypeParamType(name string) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_TypeParam{
			TypeParam: name}}
}

// NewWellKnownType creates a type corresponding to a protobuf well-known type
// value.
func NewWellKnownType(wellKnown checkedpb.Type_WellKnownType) *checkedpb.Type {
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_WellKnown{
			WellKnown: wellKnown}}
}

// NewWrapperType creates a wrapped primitive type instance. Wrapped types
// are roughly equivalent to a nullable, or optionally valued type.
func NewWrapperType(wrapped *checkedpb.Type) *checkedpb.Type {
	primitive := wrapped.GetPrimitive()
	if primitive == checkedpb.Type_PRIMITIVE_TYPE_UNSPECIFIED {
		// TODO: return an error
		panic("Wrapped type must be a primitive")
	}
	return &checkedpb.Type{
		TypeKind: &checkedpb.Type_Wrapper{
			Wrapper: primitive}}
}
