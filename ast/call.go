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

type CallExpression struct {
	BaseExpression

	// Target is the target of the call, if this is an instance-style method call.
	Target Expression

	// Function is the name of the function that is being called.
	Function string

	// Args are the arguments to the call.
	Args []Expression
}

func (e *CallExpression) String() string {
	return ToDebugString(e)
}

func (e *CallExpression) writeDebugString(w *debugWriter) {
	if e.Target != nil {
		w.appendExpression(e.Target)
		w.append(".")
	}
	w.append(e.Function)
	w.append("(")
	if len(e.Args) > 0 {
		w.addIndent()
		w.appendLine()
		for i, arg := range e.Args {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.appendExpression(arg)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append(")")
	w.adorn(e)
}

func NewCallFunction(id int64, l common.Location, function string, args ...Expression) *CallExpression {
	return &CallExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Function:       function,
		Args:           args,
	}
}

func NewCallMethod(id int64, l common.Location, function string, target Expression, args ...Expression) *CallExpression {
	return &CallExpression{
		BaseExpression: BaseExpression{id: id, location: l},
		Function:       function,
		Target:         target,
		Args:           args,
	}
}
