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

type MapType struct {
	KeyType   Type
	ValueType Type
}

var _ Type = &MapType{}

func (m *MapType) Kind() TypeKind {
	return KindMap
}

func (m *MapType) Equals(t Type) bool {
	if other, ok := t.(*MapType); ok {
		return m.KeyType.Equals(other.KeyType) && m.ValueType.Equals(other.ValueType)
	}

	return false
}

func (m *MapType) String() string {
	return fmt.Sprintf("map(%s, %s)", m.KeyType.String(), m.ValueType.String())
}

func NewMap(keyType Type, valueType Type) *MapType {
	return &MapType{
		KeyType:   keyType,
		ValueType: valueType,
	}
}
