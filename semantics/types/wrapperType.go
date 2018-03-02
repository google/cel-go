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

type WrapperType struct {
	primitive *PrimitiveType
}

var _ Type = &WrapperType{}

func (w *WrapperType) Kind() TypeKind {
	return KindWrapper
}

func (w *WrapperType) Equals(t Type) bool {
	other, ok := t.(*WrapperType)
	if !ok {
		return false
	}

	return w.primitive.Equals(other.primitive)
}

func (w *WrapperType) String() string {
	return fmt.Sprintf("wrapper(%s)", w.primitive.String())
}

func (w *WrapperType) Primitive() *PrimitiveType {
	return w.primitive
}

func NewWrapper(p *PrimitiveType) *WrapperType {
	return &WrapperType{
		primitive: p,
	}
}
