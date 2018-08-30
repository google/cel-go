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

	operatorspb "github.com/google/cel-go/common/operators"
	overloadspb "github.com/google/cel-go/common/overloads"
	typespb "github.com/google/cel-go/common/types"
	refpb "github.com/google/cel-go/common/types/ref"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
)

const (
	// Constant used to generate symbols during AST walking.
	genSymFormat = "_sym_@%d"
)

// WalkExpr produces a set of Instruction values from a CEL expression.
//
// WalkExpr does a post-order traversal of a CEL syntax AST, which means
// expressions are evaluated in a bottom-up fashion just as they would be in
// a recursive execution pattern.
func WalkExpr(expression *exprpb.Expr,
	metadata Metadata,
	dispatcher Dispatcher,
	state MutableEvalState) []Instruction {
	nextId := maxId(expression)
	walker := &astWalker{
		dispatcher: dispatcher,
		genSymId:   nextId,
		genExprId:  nextId,
		metadata:   metadata,
		scope:      newScope(),
		state:      state}
	return walker.walk(expression)
}

// astWalker implementation of the AST walking logic.
type astWalker struct {
	dispatcher Dispatcher
	genExprId  int64
	genSymId   int64
	metadata   Metadata
	scope      *blockScope
	state      MutableEvalState
}

func (w *astWalker) walk(node *exprpb.Expr) []Instruction {
	switch node.ExprKind.(type) {
	case *exprpb.Expr_CallExpr:
		return w.walkCall(node)
	case *exprpb.Expr_IdentExpr:
		return w.walkIdent(node)
	case *exprpb.Expr_SelectExpr:
		return w.walkSelect(node)
	case *exprpb.Expr_LiteralExpr:
		w.walkLiteral(node)
		return []Instruction{}
	case *exprpb.Expr_ListExpr:
		return w.walkList(node)
	case *exprpb.Expr_StructExpr:
		return w.walkStruct(node)
	case *exprpb.Expr_ComprehensionExpr:
		return w.walkComprehension(node)
	}
	return []Instruction{}
}

func (w *astWalker) walkLiteral(node *exprpb.Expr) {
	literal := node.GetLiteralExpr()
	var value refpb.Value = nil
	switch literal.LiteralKind.(type) {
	case *exprpb.Literal_BoolValue:
		value = typespb.Bool(literal.GetBoolValue())
	case *exprpb.Literal_BytesValue:
		value = typespb.Bytes(literal.GetBytesValue())
	case *exprpb.Literal_DoubleValue:
		value = typespb.Double(literal.GetDoubleValue())
	case *exprpb.Literal_Int64Value:
		value = typespb.Int(literal.GetInt64Value())
	case *exprpb.Literal_NullValue:
		value = typespb.Null(literal.GetNullValue())
	case *exprpb.Literal_StringValue:
		value = typespb.String(literal.GetStringValue())
	case *exprpb.Literal_Uint64Value:
		value = typespb.Uint(literal.GetUint64Value())
	}
	w.state.SetValue(node.Id, value)
}

func (w *astWalker) walkIdent(node *exprpb.Expr) []Instruction {
	identName := node.GetIdentExpr().Name
	if _, found := w.scope.ref(identName); !found {
		ident := NewIdent(node.Id, identName)
		w.scope.setRef(identName, node.Id)
		return []Instruction{ident}
	}
	return []Instruction{}
}

func (w *astWalker) walkSelect(node *exprpb.Expr) []Instruction {
	sel := node.GetSelectExpr()
	operandId := w.getId(sel.Operand)
	return append(
		w.walk(sel.Operand),
		NewSelect(node.Id, operandId, sel.Field))
}

func (w *astWalker) walkCall(node *exprpb.Expr) []Instruction {
	call := node.GetCallExpr()
	function := call.Function
	argGroups, argGroupLens, argIds := w.walkCallArgs(call)
	argCount := len(argIds)

	// Compute the instruction set, making sure to special case the behavior of
	// logical and, logical or, and conditional operators.
	var instructions []Instruction
	switch function {
	case operatorspb.LogicalAnd, operatorspb.LogicalOr:
		// Compute the left-hand side with a jump if the value can be used to
		// short-circuit the expression.
		//
		// Instruction layout:
		// 0: lhs expr
		// 1: jump to <END> on true (||), false (&&)
		// 2: rhs expr
		// 3: <END> logical-op(lhs, rhs)
		var instructionCount = argCount - 1
		for _, argGroupLen := range argGroupLens {
			instructionCount += argGroupLen
		}
		var evalCount = 0
		// Logical operators may have more than two arg groups in the future.
		// e.g, and(a, b, c) === a && b && c.
		// Ensure the groups are appropriately laid-out in memory.
		for i, argGroup := range argGroups {
			evalCount += argGroupLens[i]
			instructions = append(instructions, argGroup...)
			if i != argCount-1 {
				instructions = append(instructions,
					NewJump(
						argIds[i],
						instructionCount-evalCount,
						jumpIfEqual(argIds[i], typespb.Bool(function == operatorspb.LogicalOr))))
			}
			evalCount += 1
		}
		return append(instructions, NewCall(node.Id, call.Function, argIds))

	case operatorspb.Conditional:
		// Compute the conditional jump, with two jumps, one for false,
		// and one for true
		//
		// Instruction layout:
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

		// 0: condition
		instructions = append(instructions, condition...)
		// 1: jump to <END> on undefined/error
		instructions = append(instructions,
			NewJump(conditionId, len(trueVal)+len(falseVal)+3,
				jumpIfUnknownOrError(conditionId)))
		// 2: jump to <ELSE> on false.
		instructions = append(instructions,
			NewJump(conditionId, len(trueVal)+2,
				jumpIfEqual(conditionId, typespb.False)))
		// 3: <IF> expr
		instructions = append(instructions, trueVal...)
		// 4: jump to <END>
		instructions = append(instructions,
			NewJump(trueId, len(falseVal)+1, jumpAlways))
		// 5: <ELSE> expr
		instructions = append(instructions, falseVal...)
		// 6: <END> ternary
		return append(instructions, NewCall(node.Id, call.Function, argIds))

	default:
		for _, argGroup := range argGroups {
			instructions = append(instructions, argGroup...)
		}
		return append(instructions, NewCall(node.Id, call.Function, argIds))
	}
}

func (w *astWalker) walkList(node *exprpb.Expr) []Instruction {
	listExpr := node.GetListExpr()
	var elementIds []int64
	var elementSteps []Instruction
	for _, elem := range listExpr.GetElements() {
		elementIds = append(elementIds, w.getId(elem))
		elementSteps = append(elementSteps, w.walk(elem)...)
	}
	return append(elementSteps, NewList(node.Id, elementIds))
}

func (w *astWalker) walkStruct(node *exprpb.Expr) []Instruction {
	structExpr := node.GetStructExpr()
	keyValues := make(map[int64]int64)
	fieldValues := make(map[string]int64)
	var entrySteps []Instruction
	for _, entry := range structExpr.GetEntries() {
		valueId := w.getId(entry.GetValue())
		switch entry.KeyKind.(type) {
		case *exprpb.Expr_CreateStruct_Entry_FieldKey:
			fieldValues[entry.GetFieldKey()] = valueId
		case *exprpb.Expr_CreateStruct_Entry_MapKey:
			keyValues[w.getId(entry.GetMapKey())] = valueId
			entrySteps = append(entrySteps, w.walk(entry.GetMapKey())...)
		}
		entrySteps = append(entrySteps, w.walk(entry.GetValue())...)
	}
	if len(structExpr.MessageName) == 0 {
		return append(entrySteps, NewMap(node.Id, keyValues))
	}
	return append(entrySteps,
		NewObject(node.Id, structExpr.MessageName, fieldValues))
}

func (w *astWalker) walkComprehension(node *exprpb.Expr) []Instruction {
	// Serializing a comprehension into a linear set of executable steps is one
	// of the more complex tasks in AST walking. The challenge being loop
	// termination when errors or unknown values are encountered outside
	// of the accumulation steps.

	// The following example indicate sthe set of steps for the 'all' macro
	//
	// Expr: list.all(x, x < 10)
	//
	// Instruction layout:
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
	rangeSteps := w.walk(comprehensionRange)

	// Push Module with the accumulator, iter var, and iterator
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
	accuInitSteps := w.walk(comprehensionAccu)

	// iter-init
	iterInitStep :=
		NewCall(iteratorId, overloadspb.Iterator, []int64{w.getId(comprehensionRange)})

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
	loopConditionSteps := w.walk(comprehensionLoop)
	loopStepSteps := w.walk(comprehensionStep)
	loopInstructionCount := 7 + len(loopConditionSteps) + len(loopStepSteps)

	// iter-hasNext
	iterHasNextId := w.nextExprId()
	iterHasNextStep :=
		NewCall(iterHasNextId, overloadspb.HasNext, []int64{iteratorId})
	// jump <END> if !it.hasNext()
	jumpIterEndStep :=
		NewJump(iterHasNextId, loopInstructionCount-2, breakIfEnd(iterHasNextId))
	// eval x = it.next()
	// eval <cond>
	// jump <END> if condition false
	jumpConditionFalseStep :=
		NewJump(loopId, loopInstructionCount-4, jumpIfEqual(loopId, typespb.False))

	// iter-next
	nextIterVarStep := NewCall(iterNextId, overloadspb.Next, []int64{iteratorId})
	// assign the loop-step to the accu var
	accuUpdateStep := NewMov(stepId, accuId)
	// jump <LOOP>
	jumpCondStep := NewJump(stepId, -(loopInstructionCount - 2), jumpAlways)

	// <END> result
	resultSteps := w.walk(result)
	compResultUpdateStep := NewMov(w.getId(result), w.getId(node))
	popScopeStep := NewPopScope(w.getId(node))
	w.popScope()

	var instructions []Instruction
	instructions = append(instructions, rangeSteps...)
	instructions = append(instructions, pushScopeStep)
	instructions = append(instructions, accuInitSteps...)
	instructions = append(instructions, iterInitStep, iterHasNextStep, jumpIterEndStep)
	instructions = append(instructions, loopConditionSteps...)
	instructions = append(instructions, jumpConditionFalseStep, nextIterVarStep)
	instructions = append(instructions, loopStepSteps...)
	instructions = append(instructions, accuUpdateStep, jumpCondStep)
	instructions = append(instructions, resultSteps...)
	instructions = append(instructions, compResultUpdateStep, popScopeStep)
	return instructions
}

func (w *astWalker) walkCallArgs(call *exprpb.Expr_Call) (
	argGroups [][]Instruction, argGroupLens []int, argIds []int64) {
	args := getArgs(call)
	argCount := len(args)
	argGroups = make([][]Instruction, argCount)
	argGroupLens = make([]int, argCount)
	argIds = make([]int64, argCount)
	for i, arg := range getArgs(call) {
		argIds[i] = w.getId(arg)
		argGroups[i] = w.walk(arg)
		argGroupLens[i] = len(argGroups[i])
	}
	return // named outputs.
}

// Helper functions.

// getArgs returns a unified set of call args for both global and receiver
// style calls.
func getArgs(call *exprpb.Expr_Call) []*exprpb.Expr {
	var argSet []*exprpb.Expr
	if call.Target != nil {
		argSet = append(argSet, call.Target)
	}
	if call.GetArgs() != nil {
		argSet = append(argSet, call.GetArgs()...)
	}
	return argSet
}

// nextSymId generates an expression-unique identifier name for identifiers
// that need to be produced programmatically.
func (w *astWalker) nextSymId() string {
	nextId := w.genSymId
	w.genSymId++
	return fmt.Sprintf(genSymFormat, nextId)
}

// nextExprId generates expression ids when they are necessary for tracking
// evaluation state, but not captured as part of the AST.
func (w *astWalker) nextExprId() int64 {
	nextId := w.genExprId
	w.genExprId++
	return nextId
}

// pushScope moves a new scope for expression id resolution onto the stack,
// so that the same identifier name may be used in nested contexts (such as
// nested comprehensions), but that the expression ids are kept unique per
// scope.
func (w *astWalker) pushScope(scope *blockScope) {
	scope.parent = w.scope
	w.scope = scope
}

// popScope restores the identifier to expression id mapping defined in the
// prior scope.
func (w *astWalker) popScope() {
	w.scope = w.scope.parent
}

// getId returns the expression id associated with a given identifier if one
// has been set within the current scope, else the expression id.
func (w *astWalker) getId(expr *exprpb.Expr) int64 {
	id := expr.GetId()
	if ident := expr.GetIdentExpr(); ident != nil {
		if altId, found := w.scope.ref(ident.Name); found {
			w.state.SetRuntimeExpressionId(id, altId)
			return altId
		}
	}
	return id
}

// blockScope tracks identifier references within a scope and ensures that for
// all possible references to the same identifier, the same expression id is
// used within generated Instruction values.
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
	}
	return 0, false
}

func (b *blockScope) setRef(ident string, exprId int64) {
	b.references[ident] = exprId
}

func jumpIfUnknownOrError(exprId int64) func(EvalState) bool {
	return func(s EvalState) bool {
		if val, found := s.Value(exprId); found {
			return typespb.IsUnknown(val) || typespb.IsError(val)
		}
		return false
	}
}

func breakIfEnd(conditionId int64) func(EvalState) bool {
	return func(s EvalState) bool {
		if val, found := s.Value(conditionId); found {
			return val == typespb.False ||
				typespb.IsUnknown(val) ||
				typespb.IsError(val)
		}
		return true
	}
}

func jumpIfEqual(exprId int64, value refpb.Value) func(EvalState) bool {
	return func(s EvalState) bool {
		if val, found := s.Value(exprId); found {
			if typespb.IsBool(val.Type()) {
				return bool(val.Equal(value).(typespb.Bool))
			}
		}
		return false
	}
}

func jumpAlways(_ EvalState) bool {
	return true
}

func comprehensionCount(nodes ...*exprpb.Expr) int64 {
	if nodes == nil || len(nodes) == 0 {
		return 0
	}
	count := int64(0)
	for _, node := range nodes {
		if node == nil {
			continue
		}
		switch node.ExprKind.(type) {
		case *exprpb.Expr_SelectExpr:
			count += comprehensionCount(node.GetSelectExpr().GetOperand())
		case *exprpb.Expr_CallExpr:
			call := node.GetCallExpr()
			count += comprehensionCount(call.GetTarget()) + comprehensionCount(call.GetArgs()...)
		case *exprpb.Expr_ListExpr:
			count += comprehensionCount(node.GetListExpr().GetElements()...)
		case *exprpb.Expr_StructExpr:
			for _, entry := range node.GetStructExpr().GetEntries() {
				count += comprehensionCount(entry.GetMapKey()) +
					comprehensionCount(entry.GetValue())
			}
		case *exprpb.Expr_ComprehensionExpr:
			compre := node.GetComprehensionExpr()
			count += 1
			count += comprehensionCount(compre.IterRange) +
				comprehensionCount(compre.AccuInit) +
				comprehensionCount(compre.LoopCondition) +
				comprehensionCount(compre.LoopStep) +
				comprehensionCount(compre.Result)
		}
	}
	return count
}

func maxId(node *exprpb.Expr) int64 {
	if node == nil {
		return 0
	}
	currId := node.Id
	switch node.ExprKind.(type) {
	case *exprpb.Expr_SelectExpr:
		return maxInt(currId, maxId(node.GetSelectExpr().Operand))
	case *exprpb.Expr_CallExpr:
		call := node.GetCallExpr()
		currId = maxInt(currId, maxId(call.Target))
		for _, arg := range call.Args {
			currId = maxInt(currId, maxId(arg))
		}
		return currId
	case *exprpb.Expr_ListExpr:
		list := node.GetListExpr()
		for _, elem := range list.Elements {
			currId = maxInt(currId, maxId(elem))
		}
		return currId
	case *exprpb.Expr_StructExpr:
		str := node.GetStructExpr()
		for _, entry := range str.Entries {
			currId = maxInt(currId, entry.Id, maxId(entry.GetMapKey()), maxId(entry.Value))
		}
		return currId
	case *exprpb.Expr_ComprehensionExpr:
		compre := node.GetComprehensionExpr()
		return maxInt(currId,
			maxId(compre.IterRange),
			maxId(compre.AccuInit),
			maxId(compre.LoopCondition),
			maxId(compre.LoopStep),
			maxId(compre.Result))
	default:
		return currId
	}
}

func maxInt(vals ...int64) int64 {
	var result int64
	for _, val := range vals {
		if val > result {
			result = val
		}
	}
	return result
}
