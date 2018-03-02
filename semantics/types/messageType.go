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

type MessageType struct {
	name string
}

var _ Type = &MessageType{}

func (m *MessageType) Kind() TypeKind {
	return KindMessage
}

func (m *MessageType) Equals(t Type) bool {
	if mt, ok := t.(*MessageType); ok {
		return t.(*MessageType).name == mt.name
	}

	return false
}

func (m *MessageType) String() string {
	return m.name
}

func (m *MessageType) Name() string {
	return m.name
}

func NewMessage(name string) *MessageType {
	return &MessageType{
		name: name,
	}
}
