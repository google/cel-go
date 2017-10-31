package checker

import "celgo/semantics"

type Scopes struct {
	scopes []*Group
}

func NewScopes() *Scopes {
	return &Scopes{
		scopes: []*Group{},
	}
}

func (s *Scopes) Push() {
	g := newGroup()
	s.scopes = append(s.scopes, g)
}

func (s *Scopes) Pop() {
	s.scopes = s.scopes[:len(s.scopes)-1]
}

func (s *Scopes) AddIdent(ident *semantics.Ident) {
	s.scopes[0].idents[ident.Name()] = ident
}

func (s *Scopes) FindIdent(name string) *semantics.Ident {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		scope := s.scopes[i]
		if ident, found := scope.idents[name]; found {
			return ident
		}
	}
	return nil
}

func (s *Scopes) FindIdentInScope(name string) *semantics.Ident {
	if ident, found := s.scopes[len(s.scopes)-1].idents[name]; found {
		return ident
	}
	return nil
}

func (s *Scopes) AddFunction(fn *semantics.Function) {
	s.scopes[0].functions[fn.Name()] = fn
}

func (s *Scopes) FindFunction(name string) *semantics.Function {
	for i := len(s.scopes) - 1; i >= 0; i-- {
		scope := s.scopes[i]
		if fn, found := scope.functions[name]; found {
			return fn
		}
	}
	return nil
}

type Group struct {
	idents    map[string]*semantics.Ident
	functions map[string]*semantics.Function
}

func newGroup() *Group {
	return &Group{
		idents:    make(map[string]*semantics.Ident),
		functions: make(map[string]*semantics.Function),
	}
}
