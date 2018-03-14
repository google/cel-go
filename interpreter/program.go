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
	"github.com/google/cel-go/common"
	"fmt"
	expr "github.com/google/cel-spec/proto/v1"
	"strings"
)

const (
	UnknownError = "undef"
)

type Program interface {
	Container() string
	Instructions() InstructionStepper
	Metadata() Metadata
}

type InstructionStepper interface {
	Next() (Instruction, bool)
	JumpCount(count int) bool
}

type exprProgram struct {
	ast          *expr.Expr
	container    string
	instructions []Instruction
	metadata     *exprMetadata
}

var _ Program = &exprProgram{}

func NewProgram(container string, ast *expr.Expr,
	info *expr.SourceInfo) *exprProgram {
	metadata := newExprMetadata(info)
	walker := NewAstWalker(metadata)
	return &exprProgram{
		ast,
		container,
		walker.Walk(ast),
		metadata}
}

func (p *exprProgram) String() string {
	instStrs := make([]string, len(p.instructions), len(p.instructions))
	for i, inst := range p.instructions {
		instStrs[i] = fmt.Sprintf("%d: %v", i, inst)
	}
	return strings.Join(instStrs, "\n")
}

func (p *exprProgram) Container() string {
	return p.container
}

func (p *exprProgram) Instructions() InstructionStepper {
	return &exprStepper{p, 0}
}

func (p *exprProgram) Metadata() Metadata {
	return p.metadata
}

// The exprStepper keeps a cursor pointed at the next instruction to execute
// in the program.
type exprStepper struct {
	program     *exprProgram
	instruction int
}

var _ InstructionStepper = &exprStepper{}

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

var _ Metadata = &exprMetadata{}

func newExprMetadata(info *expr.SourceInfo) *exprMetadata {
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
		return common.NewLocation(m.info.Location, line, int(column)), true
	}
	return nil, false
}

func (m *exprMetadata) CharacterOffset(exprId int64) (int32, bool) {
	position, found := m.info.Positions[exprId]
	return position, found
}

func (m *exprMetadata) Expressions() []int64 {
	expressions := make([]int64, len(m.info.Positions))
	i := 0
	for key := range m.info.Positions {
		expressions[i] = key
		i++
	}
	return expressions
}
