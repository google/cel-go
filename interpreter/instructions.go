package interpreter

import (
	"fmt"
	"strings"
)

type Instruction interface {
	fmt.Stringer
	GetId() int64
}

var _ Instruction = &ConstExpr{}
var _ Instruction = &IdentExpr{}
var _ Instruction = &SelectExpr{}
var _ Instruction = &CallExpr{}
var _ Instruction = &CreateListExpr{}
var _ Instruction = &CreateMapExpr{}
var _ Instruction = &CreateTypeExpr{}
var _ Instruction = &MovInst{}
var _ Instruction = &JumpInst{}
var _ Instruction = &PushScopeInst{}
var _ Instruction = &PopScopeInst{}

type baseInstruction struct {
	Id int64
}

func (e *baseInstruction) GetId() int64 {
	return e.Id
}

type ConstExpr struct {
	*baseInstruction
	Value interface{}
}

func (e *ConstExpr) String() string {
	return fmt.Sprintf("const %v, r%d", e.Value, e.GetId())
}

func NewConst(exprId int64, value interface{}) *ConstExpr {
	return &ConstExpr{&baseInstruction{exprId}, value}
}

type IdentExpr struct {
	*baseInstruction
	Name string
}

func (e *IdentExpr) String() string {
	return fmt.Sprintf("local '%s', r%d", e.Name, e.GetId())
}

func NewIdent(exprId int64, name string) *IdentExpr {
	return &IdentExpr{&baseInstruction{exprId}, name}
}

type CallExpr struct {
	*baseInstruction
	Function string
	Args     []int64
	Overload string
}

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

func NewCall(exprId int64, function string, argIds []int64) *CallExpr {
	return &CallExpr{&baseInstruction{exprId}, function, argIds, ""}
}

func NewCallOverload(exprId int64, function string, argIds []int64, overload string) *CallExpr {
	return &CallExpr{&baseInstruction{exprId}, function, argIds, overload}
}

type SelectExpr struct {
	*baseInstruction
	Operand int64
	Field   string
}

func (e *SelectExpr) String() string {
	return fmt.Sprintf("call  select(%d, '%s'), r%d",
		e.Operand, e.Field, e.GetId())
}

func NewSelect(exprId int64, operandId int64, field string) *SelectExpr {
	return &SelectExpr{&baseInstruction{exprId}, operandId, field}
}

type CreateListExpr struct {
	*baseInstruction
	Elements []int64
}

func (e *CreateListExpr) String() string {
	return fmt.Sprintf("mov   list(%v), r%d", e.Elements, e.GetId())
}

func NewList(exprId int64, elements []int64) *CreateListExpr {
	return &CreateListExpr{&baseInstruction{exprId}, elements}
}

type CreateMapExpr struct {
	*baseInstruction
	KeyValues map[int64]int64
}

func (e *CreateMapExpr) String() string {
	return fmt.Sprintf("mov   map(%v), r%d", e.KeyValues, e.GetId())
}

func NewMap(exprId int64, keyValues map[int64]int64) *CreateMapExpr {
	return &CreateMapExpr{&baseInstruction{exprId}, keyValues}
}

type CreateTypeExpr struct {
	*baseInstruction
	Name        string
	FieldValues map[string]int64
}

func (e *CreateTypeExpr) String() string {
	return fmt.Sprintf("mov   type(%s%v), r%d", e.Name, e.FieldValues, e.GetId())
}

func NewType(exprId int64, name string,
	fieldValues map[string]int64) *CreateTypeExpr {
	return &CreateTypeExpr{&baseInstruction{exprId}, name, fieldValues}
}

// TODO: These expressions get a bit iffy.
// TODO: treat this like any other call?
type JumpInst struct {
	*baseInstruction
	Count   int
	OnValue interface{} // may be nil
}

func (e *JumpInst) String() string {
	if e.OnValue != nil {
		return fmt.Sprintf("jump  %d if r%d == %v", e.Count, e.GetId(), e.OnValue)
	} else {
		return fmt.Sprintf("jump  %d", e.Count)
	}
}

func NewJump(exprId int64, instructionCount int, jumpOnValue interface{}) *JumpInst {
	return &JumpInst{&baseInstruction{exprId}, instructionCount, jumpOnValue}
}

type MovInst struct {
	*baseInstruction
	ToExprId int64
}

func (e *MovInst) String() string {
	return fmt.Sprintf("mov   r%d, r%d", e.GetId(), e.ToExprId)
}

func NewMov(exprId int64, toExprId int64) *MovInst {
	return &MovInst{&baseInstruction{exprId}, toExprId}
}

type PushScopeInst struct {
	*baseInstruction
	Declarations map[string]int64
}

func (e *PushScopeInst) String() string {
	return fmt.Sprintf("block  %v", e.Declarations)
}

func NewPushScope(exprId int64, declarations map[string]int64) *PushScopeInst {
	return &PushScopeInst{&baseInstruction{exprId}, declarations}
}

type PopScopeInst struct {
	*baseInstruction
}

func (e *PopScopeInst) String() string {
	return fmt.Sprintf("end")
}

func NewPopScope(exprId int64) *PopScopeInst {
	return &PopScopeInst{&baseInstruction{exprId}}
}
