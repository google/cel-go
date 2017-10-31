package semantics

import (
	"celgo/ast"
	"celgo/semantics/types"
)

var Error = &Ident{
	name: "*ident_error*",
	t:    types.Error,
}

type Ident struct {
	name  string
	t     types.Type
	value ast.Constant
}

var _ Declaration = &Ident{}

func NewIdent(name string, t types.Type, value ast.Constant) *Ident {
	return &Ident{
		name:  name,
		t:     t,
		value: value,
	}
}

func (i *Ident) Name() string {
	return i.name
}

func (i *Ident) Type() types.Type {
	return i.t
}

func (i *Ident) Value() ast.Constant {
	return i.value
}
