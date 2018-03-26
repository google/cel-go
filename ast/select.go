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
