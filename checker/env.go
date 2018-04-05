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
	"github.com/google/cel-spec/proto/checked/v1/checked"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type Env struct {
	errors       *TypeErrors
	typeProvider types.TypeProvider

	declarations *decls.Scopes
}

func NewEnv(errors *common.Errors, typeProvider types.TypeProvider) *Env {
	declarations := decls.NewScopes()
	declarations.Push()

	return &Env{
		errors:       &TypeErrors{errors},
		typeProvider: typeProvider,
		declarations: declarations,
	}
}

func (env *Env) Add(decls ...*checked.Decl) {
	for _, decl := range decls {
		switch decl.DeclKind.(type) {
		case *checked.Decl_Ident:
			env.addIdent(decl)
		case *checked.Decl_Function:
			env.addFunction(decl)
		}
	}
}

func (env *Env) addFunction(decl *checked.Decl) {
	current := env.declarations.FindFunction(decl.Name)
	if current != nil {
		if current.Name != decl.Name {
			return
		}
		// TODO: Check for conflicts.
		function := current.GetFunction()
		function.Overloads = append(function.Overloads,
			decl.GetFunction().Overloads...)
		decl = current
	} else {
		env.declarations.AddFunction(decl)
	}
}

func (env *Env) addIdent(decl *checked.Decl) {
	current := env.declarations.FindIdentInScope(decl.Name)
	if current != nil {
		panic("ident already exists")
	}
	env.declarations.AddIdent(decl)
}

func (env *Env) LookupIdent(container string, typeName string) *checked.Decl {
	for _, candidate := range qualifiedTypeNameCandidates(container, typeName) {
		if ident := env.declarations.FindIdent(candidate); ident != nil {
			return ident
		}

		// Next try to import the name as a reference to a message type. If found,
		// the declaration is added to the outest (global) scope of the
		// environment, so next time we can access it faster.
		if t := env.typeProvider.LookupType(candidate); t != nil {
			decl := decls.NewIdent(candidate, t, nil)
			env.declarations.AddIdent(decl)
			return decl
		}

		// Next try to import this as an enum value by splitting the name in a type prefix and
		// the enum inside.
		if enumValue, found := env.typeProvider.LookupEnumValue(candidate); found {
			decl := decls.NewIdent(candidate,
				types.Int64,
				&expr.Constant{
					ConstantKind: &expr.Constant_Int64Value{
						Int64Value: enumValue}})
			env.declarations.AddIdent(decl)
		}
	}

	return nil
}

func (env *Env) LookupFunction(container string, typeName string) *checked.Decl {
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
