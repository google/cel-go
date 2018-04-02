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

// The interpreter package provides functions to evaluate CEL programs against
// a series of inputs and functions supplied at runtime.
package interpreter

// EvalState tracks the values associated with expression ids during execution.
type EvalState interface {
	// Value of the given expression id, false if not found.
	Value(int64) (interface{}, bool)
}

// MutableEvalState permits the mutation of evaluation state for a given
// expression id.
type MutableEvalState interface {
	EvalState
	// SetValue associates an expression id with a value.
	SetValue(int64, interface{})
}

// NewEvalState returns a MutableEvalState.
func NewEvalState() MutableEvalState {
	return &defaultEvalState{make(map[int64]interface{})}
}

type defaultEvalState struct {
	exprValues map[int64]interface{}
}

func (s *defaultEvalState) Value(exprId int64) (interface{}, bool) {
	object, found := s.exprValues[exprId]
	return object, found
}

func (s *defaultEvalState) SetValue(exprId int64, value interface{}) {
	s.exprValues[exprId] = value
}

// TODO: replace this with the value.proto equivalents.
type UnknownValue struct {
	Args []Instruction
}

// TODO: replace this with the value.proto equivalents.
type ErrorValue struct {
	ErrorSet []error
}
