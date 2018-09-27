// Package decls provides helpers for creating variable and function declarations.
package decls

import (
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/struct"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	// Error type used to communicate issues during type-checking.
	Error = &expr.Type{
		TypeKind: &expr.Type_Error{
			Error: &empty.Empty{}}}

	// Dyn is a top-type used to represent any value.
	Dyn = &expr.Type{
		TypeKind: &expr.Type_Dyn{
			Dyn: &empty.Empty{}}}

	// Commonly used types.
	Bool   = NewPrimitiveType(expr.Type_BOOL)
	Bytes  = NewPrimitiveType(expr.Type_BYTES)
	Double = NewPrimitiveType(expr.Type_DOUBLE)
	Int    = NewPrimitiveType(expr.Type_INT64)
	Null   = &expr.Type{
		TypeKind: &expr.Type_Null{
			Null: structpb.NullValue_NULL_VALUE}}
	String = NewPrimitiveType(expr.Type_STRING)
	Uint   = NewPrimitiveType(expr.Type_UINT64)

	// Well-known types.
	// TODO: Replace with an abstract type registry.
	Any       = NewWellKnownType(expr.Type_ANY)
	Duration  = NewWellKnownType(expr.Type_DURATION)
	Timestamp = NewWellKnownType(expr.Type_TIMESTAMP)
)

// NewFunctionType creates a function invocation contract, typically only used
// by type-checking steps after overload resolution.
func NewFunctionType(resultType *expr.Type,
	argTypes ...*expr.Type) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_Function{
			Function: &expr.Type_FunctionType{
				ResultType: resultType,
				ArgTypes:   argTypes}}}
}

// NewFunction creates a named function declaration with one or more overloads.
func NewFunction(name string,
	overloads ...*expr.Decl_FunctionDecl_Overload) *expr.Decl {
	return &expr.Decl{
		Name: name,
		DeclKind: &expr.Decl_Function{
			Function: &expr.Decl_FunctionDecl{
				Overloads: overloads}}}
}

// NewIdent creates a named identifier declaration with an optional literal
// value.
//
// Literal values are typically only associated with enum identifiers.
func NewIdent(name string, t *expr.Type, v *expr.Constant) *expr.Decl {
	return &expr.Decl{
		Name: name,
		DeclKind: &expr.Decl_Ident{
			Ident: &expr.Decl_IdentDecl{
				Type:  t,
				Value: v}}}
}

// NewInstanceOverload creates a instance function overload contract.
func NewInstanceOverload(id string, argTypes []*expr.Type,
	resultType *expr.Type) *expr.Decl_FunctionDecl_Overload {
	return &expr.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: true}
}

// NewListType generates a new list with elements of a certain type.
func NewListType(elem *expr.Type) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_ListType_{
			ListType: &expr.Type_ListType{
				ElemType: elem}}}
}

// NewMapType generates a new map with typed keys and values.
func NewMapType(key *expr.Type, value *expr.Type) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_MapType_{
			MapType: &expr.Type_MapType{
				KeyType:   key,
				ValueType: value}}}
}

// NewObjectType creates an object type for a qualified type name.
func NewObjectType(typeName string) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_MessageType{
			MessageType: typeName}}
}

// NewOverload creates a function overload declaration which contains a unique
// overload id as well as the expected argument and result types. Overloads
// must be aggregated within a Function declaration.
func NewOverload(id string, argTypes []*expr.Type,
	resultType *expr.Type) *expr.Decl_FunctionDecl_Overload {
	return &expr.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: false}
}

// NewParameterizedOverload creates a parametric function instance overload
// type.
func NewParameterizedInstanceOverload(id string,
	argTypes []*expr.Type,
	resultType *expr.Type,
	typeParams []string) *expr.Decl_FunctionDecl_Overload {
	return &expr.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: true}
}

// NewParameterizedOverload creates a parametric function overload type.
func NewParameterizedOverload(id string,
	argTypes []*expr.Type,
	resultType *expr.Type,
	typeParams []string) *expr.Decl_FunctionDecl_Overload {
	return &expr.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: false}
}

// NewPrimitiveType creates a type for a primitive value. See the var declarations
// for Int, Uint, etc.
func NewPrimitiveType(primitive expr.Type_PrimitiveType) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_Primitive{
			Primitive: primitive}}
}

// NewTypeType creates a new type designating a type.
func NewTypeType(nested *expr.Type) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_Type{
			Type: nested}}
}

// NewTypeParamType creates a type corresponding to a named, contextual parameter.
func NewTypeParamType(name string) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_TypeParam{
			TypeParam: name}}
}

// NewWellKnownType creates a type corresponding to a protobuf well-known type
// value.
func NewWellKnownType(wellKnown expr.Type_WellKnownType) *expr.Type {
	return &expr.Type{
		TypeKind: &expr.Type_WellKnown{
			WellKnown: wellKnown}}
}

// NewWrapperType creates a wrapped primitive type instance. Wrapped types
// are roughly equivalent to a nullable, or optionally valued type.
func NewWrapperType(wrapped *expr.Type) *expr.Type {
	primitive := wrapped.GetPrimitive()
	if primitive == expr.Type_PRIMITIVE_TYPE_UNSPECIFIED {
		// TODO: return an error
		panic("Wrapped type must be a primitive")
	}
	return &expr.Type{
		TypeKind: &expr.Type_Wrapper{
			Wrapper: primitive}}
}
