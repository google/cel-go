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
	"celgo/ast"
	"celgo/common"
	"celgo/semantics"
	"celgo/semantics/types"
)

type Env struct {
	errors       *TypeErrors
	typeProvider TypeProvider

	declarations *Scopes
}

func NewEnv(errors *common.Errors, typeProvider TypeProvider) *Env {
	declarations := NewScopes()
	declarations.Push()

	return &Env{
		errors:       &TypeErrors{errors},
		typeProvider: typeProvider,
		declarations: declarations,
	}
}

func (env *Env) AddFunction(function *semantics.Function) {
	current := env.declarations.FindFunction(function.Name())

	if current != nil {
		function = function.Merge(current)
	}
	env.declarations.AddFunction(function)
}

func (env *Env) AddIdent(ident *semantics.Ident) {
	current := env.declarations.FindIdentInScope(ident.Name())
	if current != nil {
		panic("ident already exists")
	}

	env.declarations.AddIdent(ident)
}

func (env *Env) LookupIdent(container string, typeName string) *semantics.Ident {
	for _, candidate := range qualifiedTypeNameCandidates(container, typeName) {
		if ident := env.declarations.FindIdent(candidate); ident != nil {
			return ident
		}

		// Next try to import the name as a reference to a message type. If found,
		// the declaration is added to the outest (global) scope of the
		// environment, so next time we can access it faster.
		if t := env.typeProvider.LookupType(candidate); t != nil {
			decl := semantics.NewIdent(candidate, t, nil)
			env.declarations.AddIdent(decl)
			return decl
		}

		// Next try to import this as an enum value by splitting the name in a type prefix and
		// the enum inside.
		if enumValue, found := env.typeProvider.LookupEnumValue(candidate); found {
			decl := semantics.NewIdent(candidate, types.Int64, ast.NewInt64Constant(0, common.NoLocation, enumValue))
			env.declarations.AddIdent(decl)
		}
	}

	return nil
}

func (env *Env) LookupFunction(container string, typeName string) *semantics.Function {
	for _, candidate := range qualifiedTypeNameCandidates(container, typeName) {
		if fn := env.declarations.FindFunction(candidate); fn != nil {
			return fn
		}
	}
	return nil
}

func (env *Env) enterScope() {
	env.declarations.Push()
}

func (env *Env) exitScope() {
	env.declarations.Pop()
}
