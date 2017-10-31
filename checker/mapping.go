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
