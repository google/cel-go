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
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types/ref"
	"strings"
)

// Instruction represents a single step within a CEL program.
type Instruction interface {
	fmt.Stringer
	GetId() int64
}

type baseInstruction struct {
	Id int64
}

// GetId returns the Expr id for the instruction.
func (e *baseInstruction) GetId() int64 {
	return e.Id
}

// ConstExpr is a constant expression.
type ConstExpr struct {
	*baseInstruction
	Value ref.Value
}

// String generates pseudo-assembly for the instruction.
func (e *ConstExpr) String() string {
	return fmt.Sprintf("const %v, r%d", e.Value, e.GetId())
}

// NewLiteral generates a ConstExpr.
func NewLiteral(exprId int64, value ref.Value) *ConstExpr {
	return &ConstExpr{&baseInstruction{exprId}, value}
}

// IdentExpr is an identifier expression.
type IdentExpr struct {
	*baseInstruction
	Name string
}

// String generates pseudo-assembly for the instruction.
func (e *IdentExpr) String() string {
	return fmt.Sprintf("local '%s', r%d", e.Name, e.GetId())
}

// NewIdent generates an IdentExpr.
func NewIdent(exprId int64, name string) *IdentExpr {
	return &IdentExpr{&baseInstruction{exprId}, name}
}

// CallExpr is a call expression where the args are referenced by id.
type CallExpr struct {
	*baseInstruction
	Function string
	Args     []int64
	Overload string
	Strict   bool
}

// String generates pseudo-assembly for the instruction.
func (e *CallExpr) String() string {
	argRegs := make([]string, len(e.Args), len(e.Args))
	for i, arg := range e.Args {
		argRegs[i] = fmt.Sprintf("r%d", arg)
	}
	return fmt.Sprintf("call  %s(%v), r%d",
		e.Function,
		strings.Join(argRegs, ", "),
		e.GetId())
}

// NewCall generates a CallExpr for non-overload calls.
func NewCall(exprId int64, function string, argIds []int64) *CallExpr {
	return &CallExpr{&baseInstruction{exprId}, function, argIds, "", checkIsStrict(function)}
}

// NewCallOverlaod generates a CallExpr for overload calls.
func NewCallOverload(exprId int64, function string, argIds []int64, overload string) *CallExpr {
	return &CallExpr{&baseInstruction{exprId}, function, argIds, overload, checkIsStrict(function)}
}

func checkIsStrict(function string) bool {
	if function != operators.LogicalAnd && function != operators.LogicalOr && function != operators.Conditional {
		return true
	}
	return false
}

// SelectExpr is a select expression where the operand is represented by id.
type SelectExpr struct {
	*baseInstruction
	Operand int64
	Field   string
}

// String generates pseudo-assembly for the instruction.
func (e *SelectExpr) String() string {
	return fmt.Sprintf("call  select(%d, '%s'), r%d",
		e.Operand, e.Field, e.GetId())
}

// NewSelect generates a SelectExpr.
func NewSelect(exprId int64, operandId int64, field string) *SelectExpr {
	return &SelectExpr{&baseInstruction{exprId}, operandId, field}
}

// CrateListExpr will create a new list from the elements referened by their ids.
type CreateListExpr struct {
	*baseInstruction
	Elements []int64
}

// String generates pseudo-assembly for the instruction.
func (e *CreateListExpr) String() string {
	return fmt.Sprintf("mov   list(%v), r%d", e.Elements, e.GetId())
}

// NewList generates a CreateListExpr.
func NewList(exprId int64, elements []int64) *CreateListExpr {
	return &CreateListExpr{&baseInstruction{exprId}, elements}
}

// CreateMapExpr will create a map from the key value pairs where each key and
// value refers to an expression id.
type CreateMapExpr struct {
	*baseInstruction
	KeyValues map[int64]int64
}

// String generates pseudo-assembly for the instruction.
func (e *CreateMapExpr) String() string {
	return fmt.Sprintf("mov   map(%v), r%d", e.KeyValues, e.GetId())
}

// NewMap generates a CreateMapExpr.
func NewMap(exprId int64, keyValues map[int64]int64) *CreateMapExpr {
	return &CreateMapExpr{&baseInstruction{exprId}, keyValues}
}

// CreateObjectExpr generates a new typed object with field values referenced
// by id.
type CreateObjectExpr struct {
	*baseInstruction
	Name        string
	FieldValues map[string]int64
}

// String generates pseudo-assembly for the instruction.
func (e *CreateObjectExpr) String() string {
	return fmt.Sprintf("mov   type(%s%v), r%d", e.Name, e.FieldValues, e.GetId())
}

// NewObject generates a CreateObjectExpr.
func NewObject(exprId int64, name string,
	fieldValues map[string]int64) *CreateObjectExpr {
	return &CreateObjectExpr{&baseInstruction{exprId}, name, fieldValues}
}

// JumpInst represents an conditional jump to an instruction offset.
type JumpInst struct {
	*baseInstruction
	Count       int
	OnCondition func(EvalState) bool
}

// String generates pseudo-assembly for the instruction.
func (e *JumpInst) String() string {
	return fmt.Sprintf("jump  %d if cond<r%d>", e.Count, e.GetId())
}

// NewJump generates a JumpInst.
func NewJump(exprId int64, instructionCount int, cond func(EvalState) bool) *JumpInst {
	return &JumpInst{
		baseInstruction: &baseInstruction{exprId},
		Count:           instructionCount,
		OnCondition:     cond}
}

// MovInst assigns the value of one expression id to another.
type MovInst struct {
	*baseInstruction
	ToExprId int64
}

// String generates pseudo-assembly for the instruction.
func (e *MovInst) String() string {
	return fmt.Sprintf("mov   r%d, r%d", e.GetId(), e.ToExprId)
}

// NewMov generates a MovInst.
func NewMov(exprId int64, toExprId int64) *MovInst {
	return &MovInst{&baseInstruction{exprId}, toExprId}
}

// PushScopeInst results in the generation of a new Activation containing the values
// of the associated declarations.
type PushScopeInst struct {
	*baseInstruction
	Declarations map[string]int64
}

// String generates pseudo-assembly for the instruction.
func (e *PushScopeInst) String() string {
	return fmt.Sprintf("block  %v", e.Declarations)
}

// NewPushScope generates a PushScopeInst.
func NewPushScope(exprId int64, declarations map[string]int64) *PushScopeInst {
	return &PushScopeInst{&baseInstruction{exprId}, declarations}
}

// PopScopeInst resets the current activation to the Activation#Parent() of the
// previous activation.
type PopScopeInst struct {
	*baseInstruction
}

// String generates pseudo-assembly for the instruction.
func (e *PopScopeInst) String() string {
	return fmt.Sprintf("end")
}

// NewPopScope generates a PopScopeInst.
func NewPopScope(exprId int64) *PopScopeInst {
	return &PopScopeInst{&baseInstruction{exprId}}
}
