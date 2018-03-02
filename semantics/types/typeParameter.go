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

type TypeParameter struct {
	name string
}

var _ Type = &TypeParameter{}

func (p *TypeParameter) Kind() TypeKind {
	return KindTypeParameter
}

func (p *TypeParameter) Equals(t Type) bool {
	other, ok := t.(*TypeParameter)
	if !ok {
		return false
	}
	return other.name == p.name
}

func (p *TypeParameter) String() string {
	return p.name
}

func NewTypeParam(name string) *TypeParameter {
	return &TypeParameter{
		name: name,
	}
}
