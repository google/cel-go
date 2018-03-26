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

package semantics

import (
	"github.com/google/cel-go/ast"
	"github.com/google/cel-go/semantics/types"
)

var Error = &Ident{
	name: "*ident_error*",
	t:    types.Error,
}

type Ident struct {
	name  string
	t     types.Type
	value ast.Constant
}

var _ Declaration = &Ident{}

func NewIdent(name string, t types.Type, value ast.Constant) *Ident {
	return &Ident{
		name:  name,
		t:     t,
		value: value,
	}
}

func (i *Ident) Name() string {
	return i.name
}

func (i *Ident) Type() types.Type {
	return i.t
}

func (i *Ident) Value() ast.Constant {
	return i.value
}
