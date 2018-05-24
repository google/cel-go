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

// Package parser declares an expression parser with support for macro
// expansion.
package parser

import (
	"strconv"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/parser/gen"
	expr "github.com/google/cel-spec/proto/v1/syntax"

	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
)

// ParseText converts a text input into a parsed expression, if valid, as
// well as a list of syntax errors encountered.
//
// By default all macros are enabled. For customization of the input data
// or for customization of the macro set see Parse.
//
// Note: syntax errors may produce parse trees of unusual shape which could
// in segfaults at parse-time. While the code attempts to account for all
// such cases, it is possible a few still remain. These should be fixed by
// adding a repro case to the parser_test.go and appropriate defensive coding
// within the parser.
func ParseText(text string) (*expr.ParsedExpr, *common.Errors) {
	return Parse(common.NewStringSource(text, "<input>"), AllMacros)
}

// Parse converts a source input and macros set to a parsed expression.
func Parse(source common.Source, macros Macros) (*expr.ParsedExpr, *common.Errors) {
	p := parser{helper: newParserHelper(source, macros)}
	e := p.parse(source.Content())
	return &expr.ParsedExpr{
		Expr:       e,
		SourceInfo: p.helper.getSourceInfo(),
	}, p.helper.errors.Errors
}

type parser struct {
	gen.BaseCELVisitor
	helper *parserHelper
}

var _ gen.CELVisitor = (*parser)(nil)

func (p *parser) parse(expression string) *expr.Expr {
	stream := antlr.NewInputStream(expression)
	lexer := gen.NewCELLexer(stream)
	prsr := gen.NewCELParser(antlr.NewCommonTokenStream(lexer, 0))

	lexer.RemoveErrorListeners()
	prsr.RemoveErrorListeners()
	lexer.AddErrorListener(p.helper)
	prsr.AddErrorListener(p.helper)

	return p.Visit(prsr.Start()).(*expr.Expr)
}

// Visitor implementations.
func (p *parser) Visit(tree antlr.ParseTree) interface{} {

	switch tree.(type) {
	case *gen.StartContext:
		return p.VisitStart(tree.(*gen.StartContext))
	case *gen.ExprContext:
		return p.VisitExpr(tree.(*gen.ExprContext))
	case *gen.ConditionalAndContext:
		return p.VisitConditionalAnd(tree.(*gen.ConditionalAndContext))
	case *gen.ConditionalOrContext:
		return p.VisitConditionalOr(tree.(*gen.ConditionalOrContext))
	case *gen.RelationContext:
		return p.VisitRelation(tree.(*gen.RelationContext))
	case *gen.CalcContext:
		return p.VisitCalc(tree.(*gen.CalcContext))
	case *gen.LogicalNotContext:
		return p.VisitLogicalNot(tree.(*gen.LogicalNotContext))
	case *gen.StatementExprContext:
		return p.VisitStatementExpr(tree.(*gen.StatementExprContext))
	case *gen.PrimaryExprContext:
		return p.VisitPrimaryExpr(tree.(*gen.PrimaryExprContext))
	case *gen.SelectOrCallContext:
		return p.VisitSelectOrCall(tree.(*gen.SelectOrCallContext))
	case *gen.MapInitializerListContext:
		return p.VisitMapInitializerList(tree.(*gen.MapInitializerListContext))
	case *gen.NegateContext:
		return p.VisitNegate(tree.(*gen.NegateContext))
	case *gen.IndexContext:
		return p.VisitIndex(tree.(*gen.IndexContext))
	case *gen.UnaryContext:
		return p.VisitUnary(tree.(*gen.UnaryContext))
	}

	text := "<<nil>>"
	if tree != nil {
		text = tree.GetText()
	}
	panic(fmt.Sprintf("unknown parsetree type: '%+v': %+v [%s]", reflect.TypeOf(tree), tree, text))
}

// Visit a parse tree produced by CELParser#start.
func (p *parser) VisitStart(ctx *gen.StartContext) interface{} {
	return p.Visit(ctx.Expr())
}

// Visit a parse tree produced by CELParser#expr.
func (p *parser) VisitExpr(ctx *gen.ExprContext) interface{} {
	result := p.Visit(ctx.GetE()).(*expr.Expr)
	if ctx.GetOp() == nil {
		return result
	}

	ifTrue := p.Visit(ctx.GetE1()).(*expr.Expr)
	ifFalse := p.Visit(ctx.GetE2()).(*expr.Expr)
	return p.helper.newGlobalCall(ctx.GetOp(), operators.Conditional, result, ifTrue, ifFalse)
}

// Visit a parse tree produced by CELParser#conditionalOr.
func (p *parser) VisitConditionalOr(ctx *gen.ConditionalOrContext) interface{} {
	result := p.Visit(ctx.GetE()).(*expr.Expr)
	if ctx.GetOps() == nil {
		return result
	}
	for i, op := range ctx.GetOps() {
		next := p.Visit(ctx.GetE1()[i]).(*expr.Expr)
		result = p.helper.newGlobalCall(op, operators.LogicalOr, result, next)
	}
	return result
}

// Visit a parse tree produced by CELParser#conditionalAnd.
func (p *parser) VisitConditionalAnd(ctx *gen.ConditionalAndContext) interface{} {
	result := p.Visit(ctx.GetE()).(*expr.Expr)
	if ctx.GetOps() == nil {
		return result
	}
	for i, op := range ctx.GetOps() {
		next := p.Visit(ctx.GetE1()[i]).(*expr.Expr)
		result = p.helper.newGlobalCall(op, operators.LogicalAnd, result, next)
	}
	return result
}

// Visit a parse tree produced by CELParser#relation.
func (p *parser) VisitRelation(ctx *gen.RelationContext) interface{} {
	if ctx.Calc() != nil {
		return p.Visit(ctx.Calc())
	}
	opText := ""
	if ctx.GetOp() != nil {
		opText = ctx.GetOp().GetText()
	}

	if op, found := operators.Find(opText); found {
		lhs := p.Visit(ctx.Relation(0)).(*expr.Expr)
		rhs := p.Visit(ctx.Relation(1)).(*expr.Expr)
		return p.helper.newGlobalCall(ctx.GetOp(), op, lhs, rhs)
	}
	return p.helper.reportError(ctx, "operator not found")
}

// Visit a parse tree produced by CELParser#calc.
func (p *parser) VisitCalc(ctx *gen.CalcContext) interface{} {
	if ctx.Unary() != nil {
		return p.Visit(ctx.Unary())
	}
	opText := ""
	if ctx.GetOp() != nil {
		opText = ctx.GetOp().GetText()
	}
	if op, found := operators.Find(opText); found {
		lhs := p.Visit(ctx.Calc(0)).(*expr.Expr)
		rhs := p.Visit(ctx.Calc(1)).(*expr.Expr)
		return p.helper.newGlobalCall(ctx.GetOp(), op, lhs, rhs)
	}
	return p.helper.reportError(ctx, "operator not found")
}

func (p *parser) VisitUnary(ctx *gen.UnaryContext) interface{} {
	return p.helper.newLiteralString(ctx, "<<error>>")
}

// Visit a parse tree produced by CELParser#StatementExpr.
func (p *parser) VisitStatementExpr(ctx *gen.StatementExprContext) interface{} {
	switch ctx.Statement().(type) {
	case *gen.PrimaryExprContext:
		return p.VisitPrimaryExpr(ctx.Statement().(*gen.PrimaryExprContext))
	case *gen.SelectOrCallContext:
		return p.VisitSelectOrCall(ctx.Statement().(*gen.SelectOrCallContext))
	case *gen.IndexContext:
		return p.VisitIndex(ctx.Statement().(*gen.IndexContext))
	case *gen.CreateMessageContext:
		return p.VisitCreateMessage(ctx.Statement().(*gen.CreateMessageContext))
	}
	return p.helper.reportError(ctx, "unsupported simple expression")
}

// Visit a parse tree produced by CELParser#LogicalNot.
func (p *parser) VisitLogicalNot(ctx *gen.LogicalNotContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Statement())
	}
	target := p.Visit(ctx.Statement()).(*expr.Expr)
	return p.helper.newGlobalCall(ctx.GetOps()[0], operators.LogicalNot, target)
}

func (p *parser) VisitNegate(ctx *gen.NegateContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Statement())
	}
	target := p.Visit(ctx.Statement()).(*expr.Expr)
	return p.helper.newGlobalCall(ctx.GetOps()[0], operators.Negate, target)
}

// Visit a parse tree produced by CELParser#SelectOrCall.
func (p *parser) VisitSelectOrCall(ctx *gen.SelectOrCallContext) interface{} {
	operand := p.Visit(ctx.Statement()).(*expr.Expr)
	// Handle the error case where no valid identifier is specified.
	if ctx.GetId() == nil {
		return p.helper.newExpr(ctx)
	}
	id := ctx.GetId().GetText()
	if ctx.GetOpen() != nil {
		return p.helper.newMemberCall(ctx.GetOpen(), id, operand, p.visitList(ctx.GetArgs())...)
	}
	return p.helper.newSelect(ctx.GetOp(), operand, id)
}

// Visit a parse tree produced by CELParser#PrimaryExpr.
func (p *parser) VisitPrimaryExpr(ctx *gen.PrimaryExprContext) interface{} {
	switch ctx.Primary().(type) {
	case *gen.NestedContext:
		return p.VisitNested(ctx.Primary().(*gen.NestedContext))
	case *gen.IdentOrGlobalCallContext:
		return p.VisitIdentOrGlobalCall(ctx.Primary().(*gen.IdentOrGlobalCallContext))
	case *gen.CreateListContext:
		return p.VisitCreateList(ctx.Primary().(*gen.CreateListContext))
	case *gen.CreateStructContext:
		return p.VisitCreateStruct(ctx.Primary().(*gen.CreateStructContext))
	case *gen.ConstantLiteralContext:
		return p.VisitConstantLiteral(ctx.Primary().(*gen.ConstantLiteralContext))
	}

	return p.helper.reportError(ctx, "invalid primary expression")
}

// Visit a parse tree produced by CELParser#Index.
func (p *parser) VisitIndex(ctx *gen.IndexContext) interface{} {
	target := p.Visit(ctx.Statement()).(*expr.Expr)
	index := p.Visit(ctx.GetIndex()).(*expr.Expr)
	return p.helper.newGlobalCall(ctx.GetOp(), operators.Index, target, index)
}

// Visit a parse tree produced by CELParser#CreateMessage.
func (p *parser) VisitCreateMessage(ctx *gen.CreateMessageContext) interface{} {
	target := p.Visit(ctx.Statement()).(*expr.Expr)
	if messageName, found := p.extractQualifiedName(target); found {
		entries := p.VisitIFieldInitializerList(ctx.GetEntries()).([]*expr.Expr_CreateStruct_Entry)
		return p.helper.newObject(ctx, messageName, entries...)
	}
	return p.helper.newExpr(ctx)
}

func (p *parser) VisitIFieldInitializerList(ctx gen.IFieldInitializerListContext) interface{} {
	if ctx == nil || ctx.GetFields() == nil {
		return []*expr.Expr_CreateStruct_Entry{}
	}

	result := make([]*expr.Expr_CreateStruct_Entry, len(ctx.GetFields()))
	for i, f := range ctx.GetFields() {
		value := p.Visit(ctx.GetValues()[i]).(*expr.Expr)
		field := p.helper.newObjectField(ctx.GetCols()[i], f.GetText(), value)
		result[i] = field
	}
	return result
}

// Visit a parse tree produced by CELParser#IdentOrGlobalCall.
func (p *parser) VisitIdentOrGlobalCall(ctx *gen.IdentOrGlobalCallContext) interface{} {
	identName := ""
	if ctx.GetLeadingDot() != nil {
		identName = "."
	}
	// Handle the error case where no valid identifier is specified.
	if ctx.GetId() == nil {
		return p.helper.newExpr(ctx)
	}
	identName += ctx.GetId().GetText()

	if ctx.GetOp() != nil {
		return p.helper.newGlobalCall(ctx.GetOp(), identName, p.visitList(ctx.GetArgs())...)
	}
	return p.helper.newIdent(ctx.GetId(), identName)
}

// Visit a parse tree produced by CELParser#Nested.
func (p *parser) VisitNested(ctx *gen.NestedContext) interface{} {
	return p.Visit(ctx.GetE())
}

// Visit a parse tree produced by CELParser#CreateList.
func (p *parser) VisitCreateList(ctx *gen.CreateListContext) interface{} {
	return p.helper.newList(ctx, p.visitList(ctx.GetElems())...)
}

// Visit a parse tree produced by CELParser#CreateStruct.
func (p *parser) VisitCreateStruct(ctx *gen.CreateStructContext) interface{} {
	entries := []*expr.Expr_CreateStruct_Entry{}
	if ctx.GetEntries() != nil {
		entries = p.Visit(ctx.GetEntries()).([]*expr.Expr_CreateStruct_Entry)
	}
	return p.helper.newMap(ctx.GetStart(), entries...)
}

// Visit a parse tree produced by CELParser#ConstantLiteral.
func (p *parser) VisitConstantLiteral(ctx *gen.ConstantLiteralContext) interface{} {
	switch ctx.Literal().(type) {
	case *gen.IntContext:
		return p.VisitInt(ctx.Literal().(*gen.IntContext))
	case *gen.UintContext:
		return p.VisitUint(ctx.Literal().(*gen.UintContext))
	case *gen.DoubleContext:
		return p.VisitDouble(ctx.Literal().(*gen.DoubleContext))
	case *gen.StringContext:
		return p.VisitString(ctx.Literal().(*gen.StringContext))
	case *gen.BytesContext:
		return p.VisitBytes(ctx.Literal().(*gen.BytesContext))
	case *gen.BoolFalseContext:
		return p.VisitBoolFalse(ctx.Literal().(*gen.BoolFalseContext))
	case *gen.BoolTrueContext:
		return p.VisitBoolTrue(ctx.Literal().(*gen.BoolTrueContext))
	case *gen.NullContext:
		return p.VisitNull(ctx.Literal().(*gen.NullContext))
	}
	return p.helper.reportError(ctx, "invalid literal")
}

// Visit a parse tree produced by CELParser#exprList.
func (p *parser) VisitExprList(ctx *gen.ExprListContext) interface{} {
	if ctx == nil || ctx.GetE() == nil {
		return []*expr.Expr{}
	}

	result := make([]*expr.Expr, len(ctx.GetE()))
	for i, e := range ctx.GetE() {
		exp := p.Visit(e).(*expr.Expr)
		result[i] = exp
	}
	return result
}

// Visit a parse tree produced by CELParser#mapInitializerList.
func (p *parser) VisitMapInitializerList(ctx *gen.MapInitializerListContext) interface{} {
	if ctx == nil || ctx.GetKeys() == nil {
		return []*expr.Expr_CreateStruct_Entry{}
	}

	result := make([]*expr.Expr_CreateStruct_Entry, len(ctx.GetCols()))
	for i, col := range ctx.GetCols() {
		key := p.Visit(ctx.GetKeys()[i]).(*expr.Expr)
		value := p.Visit(ctx.GetValues()[i]).(*expr.Expr)
		entry := p.helper.newMapEntry(col, key, value)
		result[i] = entry
	}
	return result
}

// Visit a parse tree produced by CELParser#Int.
func (p *parser) VisitInt(ctx *gen.IntContext) interface{} {
	i, err := strconv.ParseInt(ctx.GetTok().GetText(), 10, 64)
	if err != nil {
		return p.helper.reportError(ctx, "invalid int literal")
	}
	return p.helper.newLiteralInt(ctx, i)
}

// Visit a parse tree produced by CELParser#Uint.
func (p *parser) VisitUint(ctx *gen.UintContext) interface{} {
	text := ctx.GetTok().GetText()
	text = text[:len(text)-1]
	i, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return p.helper.reportError(ctx, "invalid uint literal")
	}
	return p.helper.newLiteralUint(ctx, i)
}

// Visit a parse tree produced by CELParser#Double.
func (p *parser) VisitDouble(ctx *gen.DoubleContext) interface{} {
	f, err := strconv.ParseFloat(ctx.GetTok().GetText(), 64)
	if err != nil {
		return p.helper.reportError(ctx, "invalid double literal")
	}
	return p.helper.newLiteralDouble(ctx, f)

}

// Visit a parse tree produced by CELParser#String.
func (p *parser) VisitString(ctx *gen.StringContext) interface{} {
	s := p.unquote(ctx, ctx.GetText())
	return p.helper.newLiteralString(ctx, s)
}

// Visit a parse tree produced by CELParser#Bytes.
func (p *parser) VisitBytes(ctx *gen.BytesContext) interface{} {
	// TODO(ozben): Not sure if this is the right encoding.
	b := []byte(p.unquote(ctx, ctx.GetTok().GetText()[1:]))
	return p.helper.newLiteralBytes(ctx, b)
}

// Visit a parse tree produced by CELParser#BoolTrue.
func (p *parser) VisitBoolTrue(ctx *gen.BoolTrueContext) interface{} {
	return p.helper.newLiteralBool(ctx, true)
}

// Visit a parse tree produced by CELParser#BoolFalse.
func (p *parser) VisitBoolFalse(ctx *gen.BoolFalseContext) interface{} {
	return p.helper.newLiteralBool(ctx, false)
}

// Visit a parse tree produced by CELParser#Null.
func (p *parser) VisitNull(ctx *gen.NullContext) interface{} {
	return p.helper.newLiteral(ctx,
		&expr.Literal{
			LiteralKind: &expr.Literal_NullValue{
				NullValue: structpb.NullValue_NULL_VALUE}})
}

func (p *parser) visitList(ctx gen.IExprListContext) []*expr.Expr {
	if ctx == nil {
		return []*expr.Expr{}
	}
	return p.visitSlice(ctx.GetE())
}

func (p *parser) visitSlice(expressions []gen.IExprContext) []*expr.Expr {
	if expressions == nil {
		return []*expr.Expr{}
	}
	result := make([]*expr.Expr, len(expressions))
	for i, e := range expressions {
		ex := p.Visit(e).(*expr.Expr)
		result[i] = ex
	}
	return result
}

func (p *parser) extractQualifiedName(e *expr.Expr) (string, bool) {
	if e == nil {
		return "", false
	}
	switch e.ExprKind.(type) {
	case *expr.Expr_IdentExpr:
		return e.GetIdentExpr().Name, true
	case *expr.Expr_SelectExpr:
		s := e.GetSelectExpr()
		if prefix, found := p.extractQualifiedName(s.Operand); found {
			return prefix + "." + s.Field, true
		}
	}
	// TODO: Add a method to Source to get location from character offset.
	location := p.helper.getLocation(e.Id)
	p.helper.reportError(location, "expected a qualified name")
	return "", false
}

func (p *parser) unquote(ctx interface{}, value string) string {
	if text, err := strconv.Unquote(value); err == nil {
		return text
	}
	p.helper.reportError(ctx, "unable to unquote string")
	return value
}
