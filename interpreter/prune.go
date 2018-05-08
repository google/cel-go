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
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type exprPruner struct {
	expr    *expr.Expr
	program Program
	state   EvalState
}

func PruneExpr(expr *expr.Expr, program Program, state EvalState) {
	pruner := &exprPruner{
		expr:    expr,
		program: program,
		state:   state}
	pruner.prune(expr)
}

func ConstantFoldExpr(expr *expr.ParsedExpr, interpreter Interpreter, container string) interface{} {
	program := NewProgram(expr.Expr, expr.SourceInfo, container)
	interpretable := interpreter.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{}))
	// We could prune even in case of error but that is probably not very useful.
	if !types.IsError(result) {
		PruneExpr(expr.Expr, program, state)
	}
	return result
}

func (p *exprPruner) createLiteral(node *expr.Expr, value *expr.Literal) {
	node.ExprKind = &expr.Expr_LiteralExpr{LiteralExpr: value}
}

func (p *exprPruner) prune(node *expr.Expr) {
	if node == nil {
		return
	}
	if val, notErrorOrUnknown := p.value(p.program.GetRuntimeExpressionId(node.GetId())); notErrorOrUnknown {

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

	// We have either an unknown/error value, or something we dont want to transform. If
	// possible, drill down more.

	switch node.ExprKind.(type) {
	case *expr.Expr_SelectExpr:
		p.prune(node.GetSelectExpr().Operand)
	case *expr.Expr_CallExpr:
		// TODO if we have logical and/or with an unknown here, transform it such that only unknown branches are left.

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
	if object, found := p.state.Value(id); found {
		fmt.Printf("index %d value %v \n", id, object)
		return object, !(types.IsUnknown(object) || types.IsError(object))
	}
	fmt.Printf("index %d not found \n", id)
	return nil, false
}
