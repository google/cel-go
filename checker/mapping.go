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

package checker

import (
	"celgo/semantics/types"
	"fmt"
)

type Mapping struct {
	mapping map[string]types.Type
}

func newMapping() *Mapping {
	return &Mapping{
		mapping: make(map[string]types.Type),
	}
}

func (s *Mapping) Add(from types.Type, to types.Type) {
	s.mapping[typeKey(from)] = to
}

func (s *Mapping) Find(from types.Type) (types.Type, bool) {
	if r, found := s.mapping[typeKey(from)]; found {
		return r, found
	}
	return nil, false
}

func (s *Mapping) Copy() *Mapping {
	c := newMapping()

	for k, v := range s.mapping {
		c.mapping[k] = v
	}
	return c
}

func (s *Mapping) String() string {
	result := "{"

	for k, v := range s.mapping {
		result += fmt.Sprintf("%v => %v   ", k, v)
	}

	result += "}"
	return result
}
