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

type ListType struct {
	ElementType Type
}

var _ Type = &ListType{}

func (l *ListType) Kind() TypeKind {
	return KindList
}

func (l *ListType) Equals(t Type) bool {
	if other, ok := t.(*ListType); ok {
		return l.ElementType.Equals(other.ElementType)
	}

	return false
}

func (l *ListType) String() string {
	return fmt.Sprintf("list(%s)", l.ElementType.String())
}

func NewList(elementType Type) *ListType {
	return &ListType{
		ElementType: elementType,
	}
}
