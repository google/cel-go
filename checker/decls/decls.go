// Package decls provides helpers for creating variable and function declarations.
package decls

import (
	"github.com/google/cel-spec/proto/checked/v1/checked"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

func NewIdent(name string, t *checked.Type, v *expr.Constant) *checked.Decl {
	return &checked.Decl{
		Name: name,
		DeclKind: &checked.Decl_Ident{
			Ident: &checked.Decl_IdentDecl{
				Type:  t,
				Value: v}}}
}

func NewFunction(name string, overloads ...*checked.Decl_FunctionDecl_Overload) *checked.Decl {
	return &checked.Decl{
		Name: name,
		DeclKind: &checked.Decl_Function{
			Function: &checked.Decl_FunctionDecl{
				Overloads: overloads}}}
}

func NewOverload(id string, argTypes []*checked.Type, resultType *checked.Type) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: false}
}

func NewInstanceOverload(id string, argTypes []*checked.Type, resultType *checked.Type) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		IsInstanceFunction: true}
}

func NewParameterizedOverload(id string, argTypes []*checked.Type, resultType *checked.Type, typeParams []string) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: false}
}

func NewParameterizedInstanceOverload(id string, argTypes []*checked.Type, resultType *checked.Type, typeParams []string) *checked.Decl_FunctionDecl_Overload {
	return &checked.Decl_FunctionDecl_Overload{
		OverloadId:         id,
		ResultType:         resultType,
		Params:             argTypes,
		TypeParams:         typeParams,
		IsInstanceFunction: true}
}
