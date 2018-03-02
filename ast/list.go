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
