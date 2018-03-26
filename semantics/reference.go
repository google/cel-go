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
)

type Reference interface {
	Equals(Reference) bool
	String() string
}

type FunctionReference struct {
	overloads []string
}

var _ Reference = &FunctionReference{}

func (f *FunctionReference) Equals(r Reference) bool {
	if other, ok := r.(*FunctionReference); ok {
		if len(other.overloads) == len(f.overloads) {
			for i, o := range f.overloads {
				if o != other.overloads[i] {
					return false
				}
			}
			return true
		}
	}

	return false
}

func (f *FunctionReference) String() string {
	result := ""
	for i, o := range f.overloads {
		if i > 0 {
			result += "|"
		}
		result += o
	}
	return result
}

type IdentReference struct {
	name  string
	value ast.Constant
}

var _ Reference = &IdentReference{}

func (f *IdentReference) Equals(r Reference) bool {
	if other, ok := r.(*IdentReference); ok {
		return other.name == f.name
	}

	return false
}

func (i *IdentReference) String() string {
	return i.name
}

func NewIdentReference(name string, constant ast.Constant) *IdentReference {
	return &IdentReference{
		name:  name,
		value: constant,
	}
}

func NewFunctionReference(overload string) *FunctionReference {
	return &FunctionReference{
		overloads: []string{overload},
	}
}

func (f *FunctionReference) AddOverloadReference(overload string) *FunctionReference {
	return &FunctionReference{
		overloads: append(f.overloads, overload),
	}
}

func (f *FunctionReference) Overloads() []string {
	return f.overloads[:]
}

func (r *IdentReference) Name() string {
	return r.name
}

func (r *IdentReference) Value() ast.Constant {
	return r.value
}
