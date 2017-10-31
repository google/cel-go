package semantics

import (
	"fmt"

	"celgo/ast"
	"celgo/semantics/types"
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
