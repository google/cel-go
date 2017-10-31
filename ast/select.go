package ast

import "celgo/common"

type SelectExpression struct {
	BaseExpression

	// Target is the target of the selection.
	Target Expression
	// Field is the field that is being selectged.
	Field string

	// TestOnly indicates whether the expression is only for testing existence.
	TestOnly bool
}

func (e *SelectExpression) String() string {
	return ToDebugString(e)
}

func (e *SelectExpression) writeDebugString(w *debugWriter) {
	w.appendExpression(e.Target)
	w.append(".")
	w.append(e.Field)
	if e.TestOnly {
		w.append("~test-only~")
	}
	w.adorn(e)
}

func NewSelect(id int64, l common.Location, target Expression, field string, testonly bool) *SelectExpression {
	return &SelectExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Target:         target,
		Field:          field,
		TestOnly:       testonly,
	}
}
