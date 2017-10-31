package ast

import "celgo/common"

type CreateListExpression struct {
	BaseExpression

	Entries []Expression
}

func (e *CreateListExpression) String() string {
	return ToDebugString(e)
}

func (e *CreateListExpression) writeDebugString(w *debugWriter) {
	w.append("[")
	if len(e.Entries) > 0 {
		w.appendLine()
		w.addIndent()
		for i, f := range e.Entries {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.appendExpression(f)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append("]")
	w.adorn(e)
}

func NewCreateList(id int64, l common.Location, entries ...Expression) *CreateListExpression {
	return &CreateListExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Entries:        entries,
	}
}
