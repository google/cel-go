package ast

import "celgo/common"

type IdentExpression struct {
	BaseExpression

	Name string
}

func (e *IdentExpression) String() string {
	return ToDebugString(e)
}

func (e *IdentExpression) writeDebugString(w *debugWriter) {
	w.append(e.Name)
	w.adorn(e)
}

func NewIdent(id int64, l common.Location, name string) *IdentExpression {
	return &IdentExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Name:           name,
	}
}
