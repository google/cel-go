package ast

import "celgo/common"

type CreateMessageExpression struct {
	BaseExpression

	MessageName string
	Fields      []*FieldEntry
}

type FieldEntry struct {
	BaseExpression

	Name        string
	Initializer Expression
}

func (e *CreateMessageExpression) String() string {
	return ToDebugString(e)
}

func (e *CreateMessageExpression) writeDebugString(w *debugWriter) {
	w.append(e.MessageName)
	w.append("{")
	if len(e.Fields) > 0 {
		w.appendLine()
		w.addIndent()
		for i, f := range e.Fields {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.appendExpression(f)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append("}")
	w.adorn(e)
}

func (e *FieldEntry) String() string {
	return ToDebugString(e)
}

func (e *FieldEntry) writeDebugString(w *debugWriter) {
	w.append(e.Name)
	w.append(":")
	w.appendExpression(e.Initializer)
	w.adorn(e)
}

func NewCreateMessage(id int64, l common.Location, messageName string, fields ...*FieldEntry) *CreateMessageExpression {
	return &CreateMessageExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		MessageName:    messageName,
		Fields:         fields,
	}
}

func NewFieldEntry(id int64, l common.Location, name string, initializer Expression) *FieldEntry {
	return &FieldEntry{
		BaseExpression: BaseExpression{id: id, location: l},
		Name:           name,
		Initializer:    initializer,
	}
}
