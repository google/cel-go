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

type ComprehensionExpression struct {
	BaseExpression

	Variable      string
	Target        Expression
	Accumulator   string
	Init          Expression
	LoopCondition Expression
	LoopStep      Expression
	Result        Expression
}

func (e *ComprehensionExpression) String() string {
	return ToDebugString(e)
}

func (e *ComprehensionExpression) writeDebugString(w *debugWriter) {
	w.append("__comprehension__(")
	w.addIndent()
	w.appendLine()
	w.append("// Variable")
	w.appendLine()
	w.append(e.Variable)
	w.append(",")
	w.appendLine()
	w.append("// Target")
	w.appendLine()
	w.appendExpression(e.Target)
	w.append(",")
	w.appendLine()
	w.append("// Accumulator")
	w.appendLine()
	w.append(e.Accumulator)
	w.append(",")
	w.appendLine()
	w.append("// Init")
	w.appendLine()
	w.appendExpression(e.Init)
	w.append(",")
	w.appendLine()
	w.append("// LoopCondition")
	w.appendLine()
	w.appendExpression(e.LoopCondition)
	w.append(",")
	w.appendLine()
	w.append("// LoopStep")
	w.appendLine()
	w.appendExpression(e.LoopStep)
	w.append(",")
	w.appendLine()
	w.append("// Result")
	w.appendLine()
	w.appendExpression(e.Result)
	w.append(")")
	w.removeIndent()
	w.adorn(e)
}

func NewComprehension(
	id int64,
	l common.Location,
	variable string,
	target Expression,
	acc string,
	init Expression,
	condition Expression,
	step Expression,
	result Expression) *ComprehensionExpression {

	return &ComprehensionExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Variable:       variable,
		Target:         target,
		Accumulator:    acc,
		Init:           init,
		LoopCondition:  condition,
		LoopStep:       step,
		Result:         result,
	}
}
