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
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type exprPruner struct {
	expr    *expr.Expr
	program Program
	state   EvalState
}

// Prunes the given expression based on the given EvalState. Modifies the expression in place.
// A typical use of this interface would be:
// 1) Evaluate expr with some unknowns,
// 2) If result is unknown:
//   a) Maybe clone the expr(typically would be necessary only for the first iteration) and PruneExpr
//   b) Goto 1
// Functional call results which are known would be effectively cached across iterations.
func PruneExpr(expr *expr.Expr, program Program, state EvalState) {
	pruner := &exprPruner{
		expr:    expr,
		program: program,
		state:   state}
	pruner.prune(expr)
}

// Prunes the given expression by evaluating the expression using the constant values in the expression. Modifies the expression in place.
// A typical use case of this interface would be:
// 1) Compile the expression (maybe via a service and maybe after checking a compiled expression does not exists in local cache)
// 2) Prepare the environment and the interpreter
// 3) ConstantFoldexpr
// 4) Maybe cache the expression
//
// How the environment is prepared in step 2 is flexible. For example, If the caller caches the compiled and constant folded expressions,
// but is not willing to constant fold( and thus cache results of) some external calls, then he can prepare the overloads accordingly.
//
func ConstantFoldExpr(expr *expr.ParsedExpr, interpreter Interpreter, container string) interface{} {
	program := NewProgram(expr.Expr, expr.SourceInfo, container)
	interpretable := interpreter.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{}))
	// Even eval end is error, there might be legitimate cases where we want to constant fold. One such case is
	// when caller deliberately does not provide some overloads.
	PruneExpr(expr.Expr, program, state)
	return result
}

func (p *exprPruner) createLiteral(node *expr.Expr, val *expr.Literal) {
	node.ExprKind = &expr.Expr_LiteralExpr{LiteralExpr: val}
}

func (p *exprPruner) maybePruneAndOr(node *expr.Expr) bool {
	if !p.existsWithUnknownValue(node.GetId()) {
		return false
	}

	call := node.GetCallExpr()

	// We know result is unknown, so we have at least one unknown arg and if one side is a known value, we know we can ignore it.

	if p.existsWithKnownValue(call.Args[0].GetId()) {
		*node = *call.Args[1]
		return true
	}

	if p.existsWithKnownValue(call.Args[1].GetId()) {
		*node = *call.Args[0]
		return true
	}
	return false
}

func (p *exprPruner) maybePruneConditional(node *expr.Expr) bool {
	if !p.existsWithUnknownValue(node.GetId()) {
		return false
	}

	call := node.GetCallExpr()
	condVal, condValueExists := p.value(call.Args[0].GetId())
	if !condValueExists || isUnknownOrError(condVal) {
		return false
	}

	if condVal.Value().(bool) {
		*node = *call.Args[1]
	} else {
		*node = *call.Args[2]
	}
	return true

}

func (p *exprPruner) maybePruneFunction(node *expr.Expr) bool {
	call := node.GetCallExpr()
	if call.Function == operators.LogicalOr || call.Function == operators.LogicalAnd {
		return p.maybePruneAndOr(node)
	}
	if call.Function == operators.Conditional {
		return p.maybePruneConditional(node)
	}

	return false
}

func (p *exprPruner) prune(node *expr.Expr) {
	if node == nil {
		return
	}
	if val, valueExists := p.value(node.GetId()); valueExists && !isUnknownOrError(val) {

		// TODO if we have a list or struct, create a list/struct expression. This is useful especially
		// if these expressions are result of a function call.

		switch val.Type() {
		case types.BoolType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_BoolValue{val.Value().(bool)}})
			return
		case types.IntType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_Int64Value{val.Value().(int64)}})
			return
		case types.UintType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_Uint64Value{val.Value().(uint64)}})
			return
		case types.StringType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_StringValue{val.Value().(string)}})
			return
		case types.DoubleType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_DoubleValue{val.Value().(float64)}})
			return
		case types.BytesType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_BytesValue{val.Value().([]byte)}})
			return
		case types.NullType:
			p.createLiteral(node, &expr.Literal{&expr.Literal_NullValue{val.Value().(structpb.NullValue)}})
			return
		}
	}

	// We have either an unknown/error value, or something we dont want to transform, or expression was not evaluated. If
	// possible, drill down more.

	switch node.ExprKind.(type) {
	case *expr.Expr_SelectExpr:
		p.prune(node.GetSelectExpr().Operand)
	case *expr.Expr_CallExpr:
		if p.maybePruneFunction(node) {
			p.prune(node)
			return
		}

		call := node.GetCallExpr()
		for _, arg := range call.Args {
			p.prune(arg)
		}
	case *expr.Expr_ListExpr:
		list := node.GetListExpr()
		for _, elem := range list.Elements {
			p.prune(elem)
		}
	case *expr.Expr_StructExpr:
		str := node.GetStructExpr()
		for _, entry := range str.Entries {
			p.prune(entry.GetMapKey())
			p.prune(entry.Value)
		}
	case *expr.Expr_ComprehensionExpr:
		compre := node.GetComprehensionExpr()
		p.prune(compre.IterRange)
	}
}

func (p *exprPruner) value(id int64) (ref.Value, bool) {
	val, found := p.state.Value(p.program.GetRuntimeExpressionId(id))
	return val, (found && val != nil)
}

func isUnknown(val ref.Value) bool {
	return types.IsUnknown(val)
}

func isUnknownOrError(val ref.Value) bool {
	return types.IsUnknown(val) || types.IsError(val)
}

func (p *exprPruner) existsWithUnknownValue(id int64) bool {
	val, valueExists := p.value(id)
	return valueExists && isUnknown(val)
}

func (p *exprPruner) existsWithKnownValue(id int64) bool {
	val, valueExists := p.value(id)
	return valueExists && !isUnknown(val)
}
