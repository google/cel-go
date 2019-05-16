// Copyright 2019 Google LLC
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

package interpreter

import (
	"github.com/google/cel-go/common/types/ref"
)

type Variable interface {
	ID() int64
	Name() string
}

type Attribute interface {
	Variable() Variable
	Path() []*PathElem
	Select(*PathElem) Attribute
}

type PathElem struct {
	ID      int64
	ToValue func(Activation) ref.Val
}

func NewAttribute(id int64, name string) Attribute {
	return &attr{
		variable: &variable{id: id, name: name},
		path:     []*PathElem{},
	}
}

func NewRelativeAttribute(pe *PathElem) Attribute {
	return &attr{
		path:     []*PathElem{pe},
	}
}

func newPathElem(id int64, val ref.Val) *PathElem {
	return &PathElem{
		ID:      id,
		ToValue: func(Activation) ref.Val { return val },
	}
}

type variable struct {
	id   int64
	name string
}

func (v *variable) ID() int64 {
	return v.id
}

func (v *variable) Name() string {
	return v.name
}

type attr struct {
	variable Variable
	path     []*PathElem
}

func (a *attr) Variable() Variable {
	return a.variable
}

func (a *attr) Path() []*PathElem {
	return a.path
}

func (a *attr) Select(pe *PathElem) Attribute {
	return &attr{
		variable: a.variable,
		path:     append(a.path, pe),
	}
}
