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

package types

import "fmt"

type TypeType struct {
	target Type
}

var _ Type = &TypeType{}

func (tt *TypeType) Target() Type {
	return tt.target
}

func (tt *TypeType) Kind() TypeKind {
	return KindType
}

func (tt *TypeType) Equals(t Type) bool {
	other, ok := t.(*TypeType)
	if !ok {
		return false
	}

	return tt.target.Equals(other.target)
}

func (t *TypeType) String() string {
	return fmt.Sprintf("type(%s)", t.target)
}

func NewTypeType(target Type) *TypeType {
	return &TypeType{
		target: target,
	}
}
