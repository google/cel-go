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

package checker

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/checker/types"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/parser"
	"github.com/google/cel-spec/proto/checked/v1/checked"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type Env struct {
	errors       *typeErrors
	typeProvider types.TypeProvider

	declarations *decls.Scopes
}

func NewEnv(errors *common.Errors, typeProvider types.TypeProvider) *Env {
	declarations := decls.NewScopes()
	declarations.Push()

	return &Env{
		errors:       &typeErrors{errors},
		typeProvider: typeProvider,
		declarations: declarations,
	}
}

func NewStandardEnv(errors *common.Errors, typeProvider types.TypeProvider) *Env {
	e := NewEnv(errors, typeProvider)
	e.Add(StandardDeclarations()...)
	return e
}

func (e *Env) Add(decls ...*checked.Decl) {
	for _, decl := range decls {
		switch decl.DeclKind.(type) {
		case *checked.Decl_Ident:
			e.addIdent(decl)
		case *checked.Decl_Function:
			e.addFunction(decl)
		}
	}
}

func (e *Env) addOverload(f *checked.Decl, overload *checked.Decl_FunctionDecl_Overload) {
	function := f.GetFunction()
	emptyMappings := types.NewMapping()
	overloadFunction := types.NewFunction(overload.GetResultType(),
		overload.GetParams()...)
	overloadErased := types.Substitute(emptyMappings, overloadFunction, true)
	for _, existing := range function.GetOverloads() {
		existingFunction := types.NewFunction(existing.GetResultType(),
			existing.GetParams()...)
		existingErased := types.Substitute(emptyMappings, existingFunction, true)
		overlap := (types.IsAssignable(emptyMappings, overloadErased, existingErased) != nil ||
			types.IsAssignable(emptyMappings, existingErased, overloadErased) != nil)
		if overlap &&
			overload.GetIsInstanceFunction() == existing.GetIsInstanceFunction() {
			e.errors.overlappingOverload(common.NoLocation, f.Name, overload.GetOverloadId(), overloadFunction,
				existing.GetOverloadId(), existingFunction)
			return
		}
	}

	for _, macro := range parser.AllMacros {
		if macro.GetName() == f.Name && macro.GetIsInstanceStyle() == overload.GetIsInstanceFunction() &&
			macro.GetArgCount() == len(overload.GetParams()) {
			e.errors.overlappingMacro(common.NoLocation, f.Name, macro.GetArgCount())
			return
		}
	}
	function.Overloads = append(function.GetOverloads(), overload)
}

func (e *Env) addFunction(decl *checked.Decl) {
	current := e.declarations.FindFunction(decl.Name)
	if current == nil {
		//Add the function declaration without overloads and check the overloads below.
		current = decls.NewFunction(decl.Name)
		e.declarations.AddFunction(current)
	}

	for _, overload := range decl.GetFunction().GetOverloads() {
		e.addOverload(current, overload)
	}
}

func (e *Env) addIdent(decl *checked.Decl) {
	current := e.declarations.FindIdentInScope(decl.Name)
	if current != nil {
		panic("ident already exists")
	}
	e.declarations.AddIdent(decl)
}

func (e *Env) LookupIdent(container string, typeName string) *checked.Decl {
	for _, candidate := range qualifiedTypeNameCandidates(container, typeName) {
		if ident := e.declarations.FindIdent(candidate); ident != nil {
			return ident
		}

		// Next try to import the name as a reference to a message type. If found,
		// the declaration is added to the outest (global) scope of the
		// environment, so next time we can access it faster.
		if t := e.typeProvider.LookupType(candidate); t != nil {
			decl := decls.NewIdent(candidate, t, nil)
			e.declarations.AddIdent(decl)
			return decl
		}

		// Next try to import this as an enum value by splitting the name in a type prefix and
		// the enum inside.
		if enumValue, found := e.typeProvider.LookupEnumValue(candidate); found {
			decl := decls.NewIdent(candidate,
				types.Int64,
				&expr.Constant{
					ConstantKind: &expr.Constant_Int64Value{
						Int64Value: enumValue}})
			e.declarations.AddIdent(decl)
		}
	}

	return nil
}

func (e *Env) LookupFunction(container string, typeName string) *checked.Decl {
	for _, candidate := range qualifiedTypeNameCandidates(container, typeName) {
		if fn := e.declarations.FindFunction(candidate); fn != nil {
			return fn
		}
	}
	return nil
}

func (e *Env) enterScope() {
	e.declarations.Push()
}

func (e *Env) exitScope() {
	e.declarations.Pop()
}
