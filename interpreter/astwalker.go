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
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/operators"
	"fmt"
	expr "github.com/google/cel-spec/proto/v1"
)

const (
	GenSymFormat = "_sym_@%d"
)

type AstWalker struct {
	metadata  Metadata
	genExprId int64
	genSymId  string
	scope     *blockScope
}

func NewAstWalker(metadata Metadata) *AstWalker {
	var maxId int64 = 0
	for _, id := range metadata.Expressions() {
		if maxId < id {
			maxId = id + 1
		}
	}
	return &AstWalker{
		metadata,
		maxId,
		fmt.Sprintf(GenSymFormat, maxId),
		newScope()}
}

func (w *AstWalker) Walk(node *expr.Expr) []Instruction {
	switch node.ExprKind.(type) {
	case *expr.Expr_CallExpr:
		return w.walkCall(node)
	case *expr.Expr_IdentExpr:
		return w.walkIdent(node)
	case *expr.Expr_SelectExpr:
		return w.walkSelect(node)
	case *expr.Expr_ConstExpr:
		return []Instruction{w.walkConst(node)}
	case *expr.Expr_ListExpr:
		return w.walkList(node)
	case *expr.Expr_StructExpr:
		return w.walkStruct(node)
	case *expr.Expr_ComprehensionExpr:
		return w.walkComprehension(node)
	}
	return []Instruction{}
}

func (w *AstWalker) walkConst(node *expr.Expr) Instruction {
	constExpr := node.GetConstExpr()
	var value interface{} = nil
	switch constExpr.ConstantKind.(type) {
	case *expr.Constant_BoolValue:
		value = constExpr.GetBoolValue()
	case *expr.Constant_BytesValue:
		value = constExpr.GetBytesValue()
	case *expr.Constant_DoubleValue:
		value = constExpr.GetDoubleValue()
	case *expr.Constant_DurationValue:
		value = constExpr.GetDurationValue()
	case *expr.Constant_Int64Value:
		value = constExpr.GetInt64Value()
	case *expr.Constant_NullValue:
		value = constExpr.GetNullValue()
	case *expr.Constant_StringValue:
		value = constExpr.GetStringValue()
	case *expr.Constant_TimestampValue:
		value = constExpr.GetTimestampValue()
	case *expr.Constant_Uint64Value:
		value = constExpr.GetUint64Value()
	}
	return NewConst(node.Id, value)
}

func (w *AstWalker) walkIdent(node *expr.Expr) []Instruction {
	identName := node.GetIdentExpr().Name
	if _, found := w.scope.ref(identName); !found {
		ident := NewIdent(node.Id, identName)
		w.scope.setRef(identName, node.Id)
		return []Instruction{ident}
	}
	return []Instruction{}
}

func (w *AstWalker) walkSelect(node *expr.Expr) []Instruction {
	sel := node.GetSelectExpr()
	operandId := w.getId(sel.Operand)
	return append(
		w.Walk(sel.Operand),
		NewSelect(node.Id, operandId, sel.Field))
}

func (w *AstWalker) walkCall(node *expr.Expr) []Instruction {
	call := node.GetCallExpr()
	function := call.Function
	argGroups, argGroupLens, argIds := w.walkCallArgs(call)
	argCount := len(argIds)

	// Compute the instruction set.
	var instructions []Instruction
	switch function {
	case operators.LogicalAnd, operators.LogicalOr:
		var instructionCount = argCount - 1
		for _, argGroupLen := range argGroupLens {
			instructionCount += argGroupLen
		}
		var evalCount = 0
		for i, argGroup := range argGroups {
			evalCount += argGroupLens[i]
			instructions = append(instructions, argGroup...)
			if i != argCount-1 {
				instructions = append(instructions,
					NewJump(
						argIds[i],
						instructionCount-evalCount,
						function == operators.LogicalOr))
			}
			evalCount += 1
		}
		return append(instructions, NewCall(node.Id, call.Function, argIds))

	case operators.Conditional:
		// Compute the conditional jump, with two jumps, one for false,
		// and one for true
		// 0: condition
		// 1: jump to <END> on undefined/error
		// 2: jump to <ELSE> on false
		// 3: <IF> expr
		// 4: jump to <END>
		// 5: <ELSE> expr
		// 6: <END> ternary
		conditionId, condition := argIds[0], argGroups[0]
		trueId, trueVal := argIds[1], argGroups[1]
		falseVal := argGroups[2]
		instructions = append(instructions, condition...)
		// Jump on undefined instruction over both returns to end call.
		instructions = append(instructions,
			NewJump(conditionId, len(trueVal)+len(falseVal)+3, UnknownError))
		// Jump to the else.
		instructions = append(instructions,
			NewJump(conditionId, len(trueVal)+2, false))
		// Jump over the else.
		instructions = append(instructions, trueVal...)
		// Jump to the end.
		instructions = append(instructions,
			NewJump(trueId, len(falseVal), nil))
		instructions = append(instructions, falseVal...)
		return append(instructions, NewCall(node.Id, call.Function, argIds))

	default:
		for _, argGroup := range argGroups {
			instructions = append(instructions, argGroup...)
		}
		return append(instructions, NewCall(node.Id, call.Function, argIds))
	}
}

func (w *AstWalker) walkList(node *expr.Expr) []Instruction {
	listExpr := node.GetListExpr()
	var elementIds []int64
	var elementSteps []Instruction
	for _, elem := range listExpr.GetElements() {
		elementIds = append(elementIds, w.getId(elem))
		elementSteps = append(elementSteps, w.Walk(elem)...)
	}
	return append(elementSteps, NewList(node.Id, elementIds))
}

func (w *AstWalker) walkStruct(node *expr.Expr) []Instruction {
	structExpr := node.GetStructExpr()
	keyValues := make(map[int64]int64)
	fieldValues := make(map[string]int64)
	var entrySteps []Instruction
	for _, entry := range structExpr.GetEntries() {
		valueId := w.getId(entry.GetValue())
		switch entry.KeyKind.(type) {
		case *expr.Expr_CreateStruct_Entry_FieldKey:
			fieldValues[entry.GetFieldKey()] = valueId
		case *expr.Expr_CreateStruct_Entry_MapKey:
			keyValues[w.getId(entry.GetMapKey())] = valueId
			entrySteps = append(entrySteps, w.Walk(entry.GetMapKey())...)
		}
		entrySteps = append(entrySteps, w.Walk(entry.GetValue())...)
	}
	if len(structExpr.MessageName) == 0 {
		return append(entrySteps, NewMap(node.Id, keyValues))
	} else {
		return append(entrySteps,
			NewType(node.Id, structExpr.MessageName, fieldValues))
	}
}

func (w *AstWalker) walkComprehension(node *expr.Expr) []Instruction {
	// expr: list.all(x, x < 10)
	//
	// 0: list                            # iter-range
	// 1: push-scope accu, iterVar, it
	// 2: accu = true                     # init
	// 3: it = list.iterator()            # iterator()
	// 4: iterVar = it.next()             # it.next()
	// 5: <LOOP> accu                     # loopCondition
	// 6: jump <END> if !accu
	// 7: accu = accu && x < 10           # loopStep
	// 8: jump <LOOP>
	// 9: result = accu                   # result
	// 10: comp = result
	// 11: pop-scope
	comprehensionExpr := node.GetComprehensionExpr()
	comprehensionRange := comprehensionExpr.GetIterRange()
	comprehensionAccu := comprehensionExpr.GetAccuInit()
	comprehensionLoop := comprehensionExpr.GetLoopCondition()
	comprehensionStep := comprehensionExpr.GetLoopStep()
	result := comprehensionExpr.GetResult()

	// iter-range
	rangeSteps := w.Walk(comprehensionRange)

	// Push Scope with the accumulator, iter var, and iterator
	iteratorId := w.nextExprId()
	iterNextId := w.nextExprId()
	iterSymId := w.nextSymId()
	accuId := w.getId(comprehensionAccu)
	loopId := w.getId(comprehensionLoop)
	stepId := w.getId(comprehensionStep)
	pushScopeStep := NewPushScope(
		node.GetId(),
		map[string]int64{
			comprehensionExpr.AccuVar: accuId,
			comprehensionExpr.IterVar: iterNextId,
			iterSymId:                 iteratorId})
	currScope := newScope()
	currScope.setRef(comprehensionExpr.AccuVar, accuId)
	currScope.setRef(comprehensionExpr.IterVar, iterNextId)
	currScope.setRef(iterSymId, iteratorId)
	w.pushScope(currScope)
	// accu-init
	accuInitSteps := w.Walk(comprehensionAccu)

	// iter-init
	iterInitStep :=
		NewCall(iteratorId, functions.Iterator, []int64{w.getId(comprehensionRange)})

	// <LOOP>
	// Loop instruction breakdown
	// 1:                       <LOOP> it.hasNext()
	// 2:                       jmpif false, <END>
	// 3+len(cond):             x = it.next()
	// 3+len(cond):             <cond>
	// 4+len(cond):             jmpif false, <END>
	// 4+len(cond)+len(step):   <step>
	// 5+len(cond)+len(step):   mov step, accu
	// 6+len(cond)+len(step):   jmp LOOP
	// 7+len(cond)+len(step)    <result>
	// <END>

	// loop-condition and step, +len(condSteps), +len(stepSteps)
	loopConditionSteps := w.Walk(comprehensionLoop)
	loopStepSteps := w.Walk(comprehensionStep)
	loopInstructionCount := 7 + len(loopConditionSteps) + len(loopStepSteps)

	// iter-hasNext
	iterHasNextId := w.nextExprId()
	iterHasNextStep :=
		NewCall(iterHasNextId, functions.HasNext, []int64{iteratorId})
	// jump <END> if !it.hasNext()
	jumpIterEndStep :=
		NewJump(iterHasNextId, loopInstructionCount-2, false)
	// eval x = it.next()
	// eval <cond>
	// jump <END> if condition false
	jumpConditionFalseStep :=
		NewJump(loopId, loopInstructionCount-4, false)

	// iter-next
	nextIterVarStep := NewCall(iterNextId, functions.Next, []int64{iteratorId})
	// assign the loop-step to the accu var
	accuUpdateStep := NewMov(stepId, accuId)
	// jump <LOOP>
	jumpCondStep := NewJump(stepId, -(loopInstructionCount - 2), nil)

	// <END> result
	resultSteps := w.Walk(result)
	compResultUpdateStep := NewMov(w.getId(result), w.getId(node))
	popScopeStep := NewPopScope(w.getId(node))
	w.popScope()

	var instructions []Instruction
	instructions = append(instructions, rangeSteps...)
	instructions = append(instructions, pushScopeStep)
	instructions = append(instructions, accuInitSteps...)
	instructions = append(instructions, iterInitStep)
	instructions = append(instructions, iterHasNextStep,
		jumpIterEndStep)
	instructions = append(instructions, loopConditionSteps...)
	instructions = append(instructions, jumpConditionFalseStep, nextIterVarStep)
	instructions = append(instructions, loopStepSteps...)
	instructions = append(instructions, accuUpdateStep, jumpCondStep)
	instructions = append(instructions, resultSteps...)
	instructions = append(instructions, compResultUpdateStep, popScopeStep)
	return instructions
}

func (w *AstWalker) walkCallArgs(call *expr.Expr_Call) (
	argGroups [][]Instruction, argGroupLens []int, argIds []int64) {
	args := getArgs(call)
	argCount := len(args)
	argGroups = make([][]Instruction, argCount)
	argGroupLens = make([]int, argCount)
	argIds = make([]int64, argCount)
	for i, arg := range getArgs(call) {
		argIds[i] = w.getId(arg)
		argGroups[i] = w.Walk(arg)
		argGroupLens[i] = len(argGroups[i])
	}
	return
}

func getArgs(call *expr.Expr_Call) []*expr.Expr {
	var argSet []*expr.Expr
	if call.Target != nil {
		argSet = append(argSet, call.Target)
	}
	if call.GetArgs() != nil {
		argSet = append(argSet, call.GetArgs()...)
	}
	return argSet
}

func (w *AstWalker) nextSymId() string {
	defer w.incExprId()
	return fmt.Sprintf(GenSymFormat, w.genExprId)
}

func (w *AstWalker) nextExprId() int64 {
	defer w.incExprId()
	return w.genExprId
}

func (w *AstWalker) incExprId() {
	w.genExprId++
}

func (w *AstWalker) pushScope(scope *blockScope) {
	scope.parent = w.scope
	w.scope = scope
}

func (w *AstWalker) popScope() {
	w.scope = w.scope.parent
}

func (w *AstWalker) getId(expr *expr.Expr) int64 {
	id := expr.GetId()
	if ident := expr.GetIdentExpr(); ident != nil {
		if altId, found := w.scope.ref(ident.Name); found {
			return altId
		}
	}
	return id
}

type blockScope struct {
	parent     *blockScope
	references map[string]int64
}

func newScope() *blockScope {
	return &blockScope{references: make(map[string]int64)}
}

func (b *blockScope) ref(ident string) (int64, bool) {
	if inst, found := b.references[ident]; found {
		return inst, found
	} else if b.parent != nil {
		return b.parent.ref(ident)
	} else {
		return 0, false
	}
}

func (b *blockScope) setRef(ident string, exprId int64) {
	b.references[ident] = exprId
}
