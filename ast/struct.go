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

type CreateStructExpression struct {
	BaseExpression

	Entries []*StructEntry
}

type StructEntry struct {
	BaseExpression

	Key   Expression
	Value Expression
}

func (e *CreateStructExpression) String() string {
	return ToDebugString(e)
}

func (e *CreateStructExpression) writeDebugString(w *debugWriter) {
	w.append("{")
	if len(e.Entries) > 0 {
		w.appendLine()
		w.addIndent()
		for i, e := range e.Entries {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.appendExpression(e)
		}

		w.removeIndent()
		w.appendLine()
	}
	w.append("}")
	w.adorn(e)
}

func (e *StructEntry) String() string {
	return ToDebugString(e)
}

func (e *StructEntry) writeDebugString(w *debugWriter) {
	w.appendExpression(e.Key)
	w.append(":")
	w.appendExpression(e.Value)
	w.adorn(e)
}

func NewCreateStruct(id int64, l common.Location, entries ...*StructEntry) *CreateStructExpression {
	return &CreateStructExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Entries:        entries,
	}
}

func NewStructEntry(id int64, l common.Location, key Expression, value Expression) *StructEntry {
	return &StructEntry{
		BaseExpression: BaseExpression{id: id, location: l},
		Key:            key,
		Value:          value,
	}
}
