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

import (
	"fmt"
	"github.com/google/cel-spec/proto/checked/v1/checked"
)

type Mapping struct {
	mapping map[string]*checked.Type
}

func NewMapping() *Mapping {
	return &Mapping{
		mapping: make(map[string]*checked.Type),
	}
}

func (m *Mapping) Add(from *checked.Type, to *checked.Type) {
	m.mapping[typeKey(from)] = to
}

func (m *Mapping) Find(from *checked.Type) (*checked.Type, bool) {
	if r, found := m.mapping[typeKey(from)]; found {
		return r, found
	}
	return nil, false
}

func (m *Mapping) Copy() *Mapping {
	c := NewMapping()

	for k, v := range m.mapping {
		c.mapping[k] = v
	}
	return c
}

func (m *Mapping) String() string {
	result := "{"

	for k, v := range m.mapping {
		result += fmt.Sprintf("%v => %v   ", k, v)
	}

	result += "}"
	return result
}
