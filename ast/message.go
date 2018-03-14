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

package ast

import "github.com/google/cel-go/common"

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
