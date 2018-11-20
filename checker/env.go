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
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/parser"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Env is the environment for type checking.
// It consists of a Packager, a Type Provider, declarations,
// and collection of errors encountered during checking.
type Env struct {
	errors       *typeErrors
	packager     packages.Packager
	typeProvider ref.TypeProvider

	declarations *decls.Scopes
}

// NewEnv returns a new *Env with the given parameters.
func NewEnv(packager packages.Packager,
	typeProvider ref.TypeProvider,
	errors *common.Errors) *Env {
	declarations := decls.NewScopes()
	declarations.Push()

	return &Env{
		errors:       &typeErrors{errors},
		packager:     packager,
		typeProvider: typeProvider,
		declarations: declarations,
	}
}

// NewStandardEnv returns a new *Env with the given params plus standard declarations.
func NewStandardEnv(packager packages.Packager,
	typeProvider ref.TypeProvider,
	errors *common.Errors) *Env {
	e := NewEnv(packager, typeProvider, errors)
	e.Add(StandardDeclarations()...)
	return e
}

// Add adds new Decl protos to the Env.
// Panics on identifiers already in the Env.
// Adds to Env errors if there's an overlap with an existing overload.
func (e *Env) Add(decls ...*exprpb.Decl) {
	for _, decl := range decls {
		switch decl.DeclKind.(type) {
		case *exprpb.Decl_Ident:
			e.addIdent(decl)
		case *exprpb.Decl_Function:
			e.addFunction(decl)
		}
	}
}

// addOverload adds overload to function declaration f.
// If overload overlaps with an existing overload, adds to the errors
// in the Env instead.
func (e *Env) addOverload(f *exprpb.Decl, overload *exprpb.Decl_FunctionDecl_Overload) {
	function := f.GetFunction()
	emptyMappings := newMapping()
	overloadFunction := decls.NewFunctionType(overload.GetResultType(),
		overload.GetParams()...)
	overloadErased := substitute(emptyMappings, overloadFunction, true)
	for _, existing := range function.GetOverloads() {
		existingFunction := decls.NewFunctionType(existing.GetResultType(),
			existing.GetParams()...)
		existingErased := substitute(emptyMappings, existingFunction, true)
		overlap := isAssignable(emptyMappings, overloadErased, existingErased) != nil ||
			isAssignable(emptyMappings, existingErased, overloadErased) != nil
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

// addFunction adds the function Decl to the Env.
// Adds a function decl if one doesn't already exist,
// then adds all overloads from the Decl.
// If overload overlaps with an existing overload, adds to the errors
// in the Env instead.
func (e *Env) addFunction(decl *exprpb.Decl) {
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

// addIdent adds the Decl to the declarations in the Env.
// Panics if an identifier with the same name already exists.
func (e *Env) addIdent(decl *exprpb.Decl) {
	current := e.declarations.FindIdentInScope(decl.Name)
	if current != nil {
		panic("ident already exists")
	}
	e.declarations.AddIdent(decl)
}

// LookupIdent returns a Decl proto for typeName as an identifier in the Env.
// Returns nil if no such identifier is found in the Env.
func (e *Env) LookupIdent(typeName string) *exprpb.Decl {
	for _, candidate := range e.packager.ResolveCandidateNames(typeName) {
		if ident := e.declarations.FindIdent(candidate); ident != nil {
			return ident
		}

		// Next try to import the name as a reference to a message type. If found,
		// the declaration is added to the outest (global) scope of the
		// environment, so next time we can access it faster.
		if t, found := e.typeProvider.FindType(candidate); found {
			decl := decls.NewIdent(candidate, t, nil)
			e.declarations.AddIdent(decl)
			return decl
		}

		// Next try to import this as an enum value by splitting the name in a type prefix and
		// the enum inside.
		if enumValue := e.typeProvider.EnumValue(candidate); enumValue.Type() != types.ErrType {
			decl := decls.NewIdent(candidate,
				decls.Int,
				&exprpb.Constant{
					ConstantKind: &exprpb.Constant_Int64Value{
						Int64Value: int64(enumValue.(types.Int))}})
			e.declarations.AddIdent(decl)
			return decl
		}
	}
	return nil
}

// LookupFunction returns a Decl proto for typeName as a function in env.
// Returns nil if no such function is found in env.
func (e *Env) LookupFunction(typeName string) *exprpb.Decl {
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
