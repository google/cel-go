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

package interpreter

import (
	"fmt"
	"github.com/google/cel-go/common"
	"github.com/google/cel-spec/proto/checked/v1/checked"
	expr "github.com/google/cel-spec/proto/v1/syntax"
	"strings"
)

// Program contains instructions with some metadata and optionally run within
// a container, e.g. module name or package name.
type Program interface {
	// Container is the module or package name where the program is run. The
	// container is used to resolve type names and identifiers.
	Container() string

	GetInstruction(exprId int64) Instruction

	// Init returns an InstructionStepper which can be used to iterate
	// through the program. Each call to Init generates a new stepper.
	Init(dispatcher Dispatcher, state MutableEvalState) InstructionStepper

	// MaxInstructionId returns the identifier of the last expression in the
	// program.
	MaxInstructionId() int64

	// Metadata used to determine source locations of sub-expressions.
	Metadata() Metadata
}

// IntructionStepper steps through program instructions and provides an option
// to jump a certain number of instructions forward or back.
type InstructionStepper interface {
	// Next returns the next instruction, or false if the end of the program
	// has been reached.
	Next() (Instruction, bool)

	// JumpCount moves a relative count of instructions forward or back in the
	// program and returns whether the jump was successful.
	//
	// A jump may be unsuccessful if the number of instructions to jump exceeds
	// the beginning or end of the program.
	JumpCount(count int) bool
}

type exprProgram struct {
	container       string
	expression      *expr.Expr
	instructions    []Instruction
	metadata        Metadata
	revInstructions map[int64]int
}

// NewCheckedProgram creates a Program from a checked CEL expression.
func NewCheckedProgram(c *checked.CheckedExpr, container string) Program {
	return NewProgram(c.Expr, c.SourceInfo, container)
}

// NewProgram creates a Program from a CEL expression and source information
// within the specified container.
func NewProgram(expression *expr.Expr,
	info *expr.SourceInfo,
	container string) Program {
	revInstructions := make(map[int64]int)
	return &exprProgram{
		container:       container,
		expression:      expression,
		revInstructions: revInstructions,
		metadata:        newExprMetadata(info)}
}

func (p *exprProgram) Container() string {
	return p.container
}

func (p *exprProgram) GetInstruction(exprId int64) Instruction {
	return p.instructions[p.revInstructions[exprId]]
}

func (p *exprProgram) Init(dispatcher Dispatcher, state MutableEvalState) InstructionStepper {
	if p.instructions == nil {
		p.instructions = WalkExpr(p.expression, p.metadata, dispatcher, state)
		for i, inst := range p.instructions {
			p.revInstructions[inst.GetId()] = i
		}
	}
	return &exprStepper{p, 0}
}

func (p *exprProgram) MaxInstructionId() int64 {
	// The max instruction id is computed as the highest expression id + 1
	// combined with the number of comprehensions times two. Each comprehension
	// introduces two generated ids (one for an iterator and one for current
	// iterator value) once the program is initialized.
	return maxId(p.expression) + comprehensionCount(p.expression)*2
}

func (p *exprProgram) Metadata() Metadata {
	return p.metadata
}

func (p *exprProgram) String() string {
	instStrs := make([]string, len(p.instructions), len(p.instructions))
	for i, inst := range p.instructions {
		instStrs[i] = fmt.Sprintf("%d: %v", i, inst)
	}
	return strings.Join(instStrs, "\n")
}

// exprStepper keeps a cursor pointed at the next instruction to execute
// in the program.
type exprStepper struct {
	program     *exprProgram
	instruction int
}

func (s *exprStepper) Next() (Instruction, bool) {
	if s.instruction < len(s.program.instructions) {
		inst := s.instruction
		s.instruction += 1
		return s.program.instructions[inst], true
	}
	return nil, false
}

func (s *exprStepper) JumpCount(count int) bool {
	// Adjust for the cursor already having been moved.
	offset := count - 1
	candidate := s.instruction + offset
	if candidate >= 0 && candidate < len(s.program.instructions) {
		s.instruction = candidate
		return true
	}
	return false
}

// The exprMetadata type provides helper functions for retrieving source
// locations in a human readable manner based on the data contained within
// the expr.SourceInfo message.
type exprMetadata struct {
	info *expr.SourceInfo
}

func newExprMetadata(info *expr.SourceInfo) Metadata {
	return &exprMetadata{info: info}
}

func (m *exprMetadata) Location(exprId int64) (common.Location, bool) {
	if exprOffset, found := m.CharacterOffset(exprId); found {
		var index = 0
		var lineIndex = 0
		var lineOffset int32 = 0
		for index, lineOffset = range m.info.LineOffsets {
			if lineOffset > exprOffset {
				break
			}
			lineIndex = index
		}
		line := lineIndex + 1
		column := exprOffset - lineOffset
		return common.NewLocation(line, int(column)), true
	}
	return nil, false
}

func (m *exprMetadata) CharacterOffset(exprId int64) (int32, bool) {
	position, found := m.info.Positions[exprId]
	return position, found
}
