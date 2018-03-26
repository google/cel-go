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
	"fmt"

	"github.com/google/cel-go/ast"
	"github.com/google/cel-go/semantics/types"
)

type Semantics struct {
	types      map[int64]types.Type
	references map[int64]Reference
}

func New(types map[int64]types.Type, references map[int64]Reference) *Semantics {
	return &Semantics{
		types:      types,
		references: references,
	}
}

func (s *Semantics) GetType(e ast.Expression) types.Type {
	return s.types[e.Id()]
}

func (s *Semantics) GetReference(e ast.Expression) Reference {
	return s.references[e.Id()]
}

func (s *Semantics) String() string {
	result := "types:\n"
	for k, v := range s.types {
		result += fmt.Sprintf("  e:'%+v'  => t:'%+v'\n", k, v)
	}
	result += "references:\n"
	for k, v := range s.references {
		result += fmt.Sprintf("  e:'%+v'  => r:'%+v'\n", k, v)
	}
	return result
}
