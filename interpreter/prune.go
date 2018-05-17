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

type astPruner struct {
	expr  *expr.Expr
	state EvalState
}

// TODO Consider having a separate walk of the AST that finds common
// subexpressions. This can be called before or after constant folding to find
// common subexpressions.

// Prunes the given AST based on the given EvalState and generates a new AST.
// Given AST is copied on write and a new AST is returned.
// Couple of typical use cases this interface would be:
//
// A)
// 1) Evaluate expr with some unknowns,
// 2) If result is unknown:
//   a) PruneAst
//   b) Goto 1
// Functional call results which are known would be effectively cached across
// iterations.
//
// B)
// 1) Compile the expression (maybe via a service and maybe after checking a
//    compiled expression does not exists in local cache)
// 2) Prepare the environment and the interpreter. Activation might be empty.
// 3) Eval the expression. This might return unknown or error or a concrete
//    value.
// 4) PruneAst
// 4) Maybe cache the expression
// This is effectively constant folding the expression. How the environment is
// prepared in step 2 is flexible. For example, If the caller caches the
// compiled and constant folded expressions, but is not willing to constant
// fold(and thus cache results of) some external calls, then they can prepare
// the overloads accordingly.
func PruneAst(expr *expr.Expr, state EvalState) *expr.Expr {
	pruner := &astPruner{
		expr:  expr,
		state: state}
	_, newExpr := pruner.prune(expr)
	return newExpr
}

func (p *astPruner) createLiteral(node *expr.Expr, val *expr.Literal) *expr.Expr {
	newExpr := *node
	newExpr.ExprKind = &expr.Expr_LiteralExpr{LiteralExpr: val}
	return &newExpr
}

func (p *astPruner) maybePruneAndOr(node *expr.Expr) (bool, *expr.Expr) {
	if !p.existsWithUnknownValue(node.GetId()) {
		return false, nil
	}

	call := node.GetCallExpr()

	// We know result is unknown, so we have at least one unknown arg
	// and if one side is a known value, we know we can ignore it.

	if p.existsWithKnownValue(call.Args[0].GetId()) {
		return true, call.Args[1]
	}

	if p.existsWithKnownValue(call.Args[1].GetId()) {
		return true, call.Args[0]
	}
	return false, nil
}

func (p *astPruner) maybePruneConditional(node *expr.Expr) (bool, *expr.Expr) {
	if !p.existsWithUnknownValue(node.GetId()) {
		return false, nil
	}

	call := node.GetCallExpr()
	condVal, condValueExists := p.value(call.Args[0].GetId())
	if !condValueExists || types.IsUnknownOrError(condVal) {
		return false, nil
	}

	if condVal.Value().(bool) {
		return true, call.Args[1]
	} else {
		return true, call.Args[2]
	}
}

func (p *astPruner) maybePruneFunction(node *expr.Expr) (bool, *expr.Expr) {
	call := node.GetCallExpr()
	if call.Function == operators.LogicalOr || call.Function == operators.LogicalAnd {
		return p.maybePruneAndOr(node)
	}
	if call.Function == operators.Conditional {
		return p.maybePruneConditional(node)
	}

	return false, nil
}

func (p *astPruner) prune(node *expr.Expr) (bool, *expr.Expr) {
	if node == nil {
		return false, node
	}
	if val, valueExists := p.value(node.GetId()); valueExists && !types.IsUnknownOrError(val) {

		// TODO if we have a list or struct, create a list/struct
		// expression. This is useful especially if these expressions
		// are result of a function call.

		switch val.Type() {
		case types.BoolType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_BoolValue{val.Value().(bool)}})
		case types.IntType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_Int64Value{val.Value().(int64)}})
		case types.UintType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_Uint64Value{val.Value().(uint64)}})
		case types.StringType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_StringValue{val.Value().(string)}})
		case types.DoubleType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_DoubleValue{val.Value().(float64)}})
		case types.BytesType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_BytesValue{val.Value().([]byte)}})
		case types.NullType:
			return true, p.createLiteral(node, &expr.Literal{&expr.Literal_NullValue{val.Value().(structpb.NullValue)}})
		}
	}

	// We have either an unknown/error value, or something we dont want to
	// transform, or expression was not evaluated. If possible, drill down
	// more.

	switch node.ExprKind.(type) {
	case *expr.Expr_SelectExpr:
		if pruned, operand := p.prune(node.GetSelectExpr().Operand); pruned {
			newExpr := *node
			newSelect := *newExpr.GetSelectExpr()
			newSelect.Operand = operand
			newExpr.GetExprKind().(*expr.Expr_SelectExpr).SelectExpr = &newSelect
			return true, &newExpr
		}
	case *expr.Expr_CallExpr:
		if pruned, newExpr := p.maybePruneFunction(node); pruned {
			_, newExpr = p.prune(newExpr)
			return true, newExpr
		}
		newCall := *node.GetCallExpr()
		var prunedCall bool
		var prunedArg bool
		for i, arg := range node.GetCallExpr().Args {
			if prunedArg, newCall.Args[i] = p.prune(arg); prunedArg {
				prunedCall = true
			}
		}
		if prunedTarget, newTarget := p.prune(node.GetCallExpr().Target); prunedTarget {
			prunedCall = true
			newCall.Target = newTarget
		}
		if prunedCall {
			newExpr := *node
			newExpr.GetExprKind().(*expr.Expr_CallExpr).CallExpr = &newCall
			return true, &newExpr
		}
	case *expr.Expr_ListExpr:
		newList := *node.GetListExpr()
		var prunedList bool
		var prunedElem bool
		for i, elem := range node.GetListExpr().Elements {
			if prunedElem, newList.Elements[i] = p.prune(elem); prunedElem {
				prunedList = true
			}
		}
		if prunedList {
			newExpr := *node
			newExpr.GetExprKind().(*expr.Expr_ListExpr).ListExpr = &newList
			return true, &newExpr
		}
	case *expr.Expr_StructExpr:
		newStruct := *node.GetStructExpr()
		var prunedStruct bool
		var prunedEntry bool
		for i, entry := range node.GetStructExpr().Entries {
			newEntry := *entry
			if x, ok := newEntry.GetKeyKind().(*expr.Expr_CreateStruct_Entry_MapKey); ok {
				if pruned, newKey := p.prune(entry.GetMapKey()); pruned {
					prunedEntry = true
					x.MapKey = newKey
				}
			}
			if pruned, newValue := p.prune(entry.Value); pruned {
				prunedEntry = true
				newEntry.Value = newValue
			}
			if prunedEntry {
				prunedStruct = true
				newStruct.Entries[i] = &newEntry
			}
		}
		if prunedStruct {
			newExpr := *node
			newExpr.GetExprKind().(*expr.Expr_StructExpr).StructExpr = &newStruct
			return true, &newExpr
		}
	case *expr.Expr_ComprehensionExpr:
		if pruned, newIterRange := p.prune(node.GetComprehensionExpr().IterRange); pruned {
			newExpr := *node
			newCompre := *newExpr.GetComprehensionExpr()
			newCompre.IterRange = newIterRange
			newExpr.GetExprKind().(*expr.Expr_ComprehensionExpr).ComprehensionExpr = &newCompre
			return true, &newExpr
		}
	}
	return false, node
}

func (p *astPruner) value(id int64) (ref.Value, bool) {
	val, found := p.state.Value(p.state.GetRuntimeExpressionId(id))
	return val, (found && val != nil)
}

func (p *astPruner) existsWithUnknownValue(id int64) bool {
	val, valueExists := p.value(id)
	return valueExists && types.IsUnknown(val)
}

func (p *astPruner) existsWithKnownValue(id int64) bool {
	val, valueExists := p.value(id)
	return valueExists && !types.IsUnknown(val)
}
