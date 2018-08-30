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
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	commonpb "github.com/google/cel-go/common"
	declspb "github.com/google/cel-go/checker/decls"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
	packagespb "github.com/google/cel-go/common/packages"
	parserpb "github.com/google/cel-go/parser"
	refpb "github.com/google/cel-go/common/types/ref"
	typespb "github.com/google/cel-go/common/types"
)

type Env struct {
	errors       *typeErrors
	packager     packagespb.Packager
	typeProvider refpb.TypeProvider

	declarations *declspb.Scopes
}

func NewEnv(packager packagespb.Packager,
	typeProvider refpb.TypeProvider,
	errors *commonpb.Errors) *Env {
	declarations := declspb.NewScopes()
	declarations.Push()

	return &Env{
		errors:       &typeErrors{errors},
		packager:     packager,
		typeProvider: typeProvider,
		declarations: declarations,
	}
}

func NewStandardEnv(packager packagespb.Packager,
	typeProvider refpb.TypeProvider,
	errors *commonpb.Errors) *Env {
	e := NewEnv(packager, typeProvider, errors)
	e.Add(StandardDeclarations()...)
	return e
}

func (e *Env) Add(decls ...*checkedpb.Decl) {
	for _, decl := range decls {
		switch decl.DeclKind.(type) {
		case *checkedpb.Decl_Ident:
			e.addIdent(decl)
		case *checkedpb.Decl_Function:
			e.addFunction(decl)
		}
	}
}

func (e *Env) addOverload(f *checkedpb.Decl, overload *checkedpb.Decl_FunctionDecl_Overload) {
	function := f.GetFunction()
	emptyMappings := newMapping()
	overloadFunction := declspb.NewFunctionType(overload.GetResultType(),
		overload.GetParams()...)
	overloadErased := substitute(emptyMappings, overloadFunction, true)
	for _, existing := range function.GetOverloads() {
		existingFunction := declspb.NewFunctionType(existing.GetResultType(),
			existing.GetParams()...)
		existingErased := substitute(emptyMappings, existingFunction, true)
		overlap := isAssignable(emptyMappings, overloadErased, existingErased) != nil ||
			isAssignable(emptyMappings, existingErased, overloadErased) != nil
		if overlap &&
			overload.GetIsInstanceFunction() == existing.GetIsInstanceFunction() {
			e.errors.overlappingOverload(commonpb.NoLocation, f.Name, overload.GetOverloadId(), overloadFunction,
				existing.GetOverloadId(), existingFunction)
			return
		}
	}

	for _, macro := range parserpb.AllMacros {
		if macro.GetName() == f.Name && macro.GetIsInstanceStyle() == overload.GetIsInstanceFunction() &&
			macro.GetArgCount() == len(overload.GetParams()) {
			e.errors.overlappingMacro(commonpb.NoLocation, f.Name, macro.GetArgCount())
			return
		}
	}
	function.Overloads = append(function.GetOverloads(), overload)
}

func (e *Env) addFunction(decl *checkedpb.Decl) {
	current := e.declarations.FindFunction(decl.Name)
	if current == nil {
		//Add the function declaration without overloads and check the overloads below.
		current = declspb.NewFunction(decl.Name)
		e.declarations.AddFunction(current)
	}

	for _, overload := range decl.GetFunction().GetOverloads() {
		e.addOverload(current, overload)
	}
}

func (e *Env) addIdent(decl *checkedpb.Decl) {
	current := e.declarations.FindIdentInScope(decl.Name)
	if current != nil {
		panic("ident already exists")
	}
	e.declarations.AddIdent(decl)
}

func (e *Env) LookupIdent(typeName string) *checkedpb.Decl {
	for _, candidate := range e.packager.ResolveCandidateNames(typeName) {
		if ident := e.declarations.FindIdent(candidate); ident != nil {
			return ident
		}

		// Next try to import the name as a reference to a message type. If found,
		// the declaration is added to the outest (global) scope of the
		// environment, so next time we can access it faster.
		if t, found := e.typeProvider.FindType(candidate); found {
			decl := declspb.NewIdent(candidate, t, nil)
			e.declarations.AddIdent(decl)
			return decl
		}

		// Next try to import this as an enum value by splitting the name in a type prefix and
		// the enum inside.
		if enumValue := e.typeProvider.EnumValue(candidate); enumValue.Type() != typespb.ErrType {
			decl := declspb.NewIdent(candidate,
				declspb.Int,
				&exprpb.Literal{
					LiteralKind: &exprpb.Literal_Int64Value{
						Int64Value: int64(enumValue.(typespb.Int))}})
			e.declarations.AddIdent(decl)
			return decl
		}
	}
	return nil
}

func (e *Env) LookupFunction(typeName string) *checkedpb.Decl {
	for _, candidate := range e.packager.ResolveCandidateNames(typeName) {
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
