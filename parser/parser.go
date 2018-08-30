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
	"fmt"
	"reflect"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	structpb "github.com/golang/protobuf/ptypes/struct"
	commonpb "github.com/google/cel-go/common"
	operatorspb "github.com/google/cel-go/common/operators"
	genpb "github.com/google/cel-go/parser/gen"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
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
func ParseText(text string) (*exprpb.ParsedExpr, *commonpb.Errors) {
	return Parse(commonpb.NewStringSource(text, "<input>"), AllMacros)
}

// Parse converts a source input and macros set to a parsed expression.
func Parse(source commonpb.Source, macros Macros) (*exprpb.ParsedExpr, *commonpb.Errors) {
	p := parser{helper: newParserHelper(source, macros)}
	e := p.parse(source.Content())
	return &exprpb.ParsedExpr{
		Expr:       e,
		SourceInfo: p.helper.getSourceInfo(),
	}, p.helper.errors.Errors
}

type parser struct {
	genpb.BaseCELVisitor
	helper *parserHelper
}

var _ genpb.CELVisitor = (*parser)(nil)

func (p *parser) parse(expression string) *exprpb.Expr {
	stream := antlr.NewInputStream(expression)
	lexer := genpb.NewCELLexer(stream)
	prsr := genpb.NewCELParser(antlr.NewCommonTokenStream(lexer, 0))

	lexer.RemoveErrorListeners()
	prsr.RemoveErrorListeners()
	lexer.AddErrorListener(p.helper)
	prsr.AddErrorListener(p.helper)

	return p.Visit(prsr.Start()).(*exprpb.Expr)
}

// Visitor implementations.
func (p *parser) Visit(tree antlr.ParseTree) interface{} {

	switch tree.(type) {
	case *genpb.StartContext:
		return p.VisitStart(tree.(*genpb.StartContext))
	case *genpb.ExprContext:
		return p.VisitExpr(tree.(*genpb.ExprContext))
	case *genpb.ConditionalAndContext:
		return p.VisitConditionalAnd(tree.(*genpb.ConditionalAndContext))
	case *genpb.ConditionalOrContext:
		return p.VisitConditionalOr(tree.(*genpb.ConditionalOrContext))
	case *genpb.RelationContext:
		return p.VisitRelation(tree.(*genpb.RelationContext))
	case *genpb.CalcContext:
		return p.VisitCalc(tree.(*genpb.CalcContext))
	case *genpb.LogicalNotContext:
		return p.VisitLogicalNot(tree.(*genpb.LogicalNotContext))
	case *genpb.StatementExprContext:
		return p.VisitStatementExpr(tree.(*genpb.StatementExprContext))
	case *genpb.PrimaryExprContext:
		return p.VisitPrimaryExpr(tree.(*genpb.PrimaryExprContext))
	case *genpb.SelectOrCallContext:
		return p.VisitSelectOrCall(tree.(*genpb.SelectOrCallContext))
	case *genpb.MapInitializerListContext:
		return p.VisitMapInitializerList(tree.(*genpb.MapInitializerListContext))
	case *genpb.NegateContext:
		return p.VisitNegate(tree.(*genpb.NegateContext))
	case *genpb.IndexContext:
		return p.VisitIndex(tree.(*genpb.IndexContext))
	case *genpb.UnaryContext:
		return p.VisitUnary(tree.(*genpb.UnaryContext))
	}

	text := "<<nil>>"
	if tree != nil {
		text = tree.GetText()
	}
	panic(fmt.Sprintf("unknown parsetree type: '%+v': %+v [%s]", reflect.TypeOf(tree), tree, text))
}

// Visit a parse tree produced by CELParser#start.
func (p *parser) VisitStart(ctx *genpb.StartContext) interface{} {
	return p.Visit(ctx.Expr())
}

// Visit a parse tree produced by CELParser#expr.
func (p *parser) VisitExpr(ctx *genpb.ExprContext) interface{} {
	result := p.Visit(ctx.GetE()).(*exprpb.Expr)
	if ctx.GetOp() == nil {
		return result
	}

	ifTrue := p.Visit(ctx.GetE1()).(*exprpb.Expr)
	ifFalse := p.Visit(ctx.GetE2()).(*exprpb.Expr)
	return p.helper.newGlobalCall(ctx.GetOp(), operatorspb.Conditional, result, ifTrue, ifFalse)
}

// Visit a parse tree produced by CELParser#conditionalOr.
func (p *parser) VisitConditionalOr(ctx *genpb.ConditionalOrContext) interface{} {
	result := p.Visit(ctx.GetE()).(*exprpb.Expr)
	if ctx.GetOps() == nil {
		return result
	}
	for i, op := range ctx.GetOps() {
		next := p.Visit(ctx.GetE1()[i]).(*exprpb.Expr)
		result = p.helper.newGlobalCall(op, operatorspb.LogicalOr, result, next)
	}
	return result
}

// Visit a parse tree produced by CELParser#conditionalAnd.
func (p *parser) VisitConditionalAnd(ctx *genpb.ConditionalAndContext) interface{} {
	result := p.Visit(ctx.GetE()).(*exprpb.Expr)
	if ctx.GetOps() == nil {
		return result
	}
	for i, op := range ctx.GetOps() {
		next := p.Visit(ctx.GetE1()[i]).(*exprpb.Expr)
		result = p.helper.newGlobalCall(op, operatorspb.LogicalAnd, result, next)
	}
	return result
}

// Visit a parse tree produced by CELParser#relation.
func (p *parser) VisitRelation(ctx *genpb.RelationContext) interface{} {
	if ctx.Calc() != nil {
		return p.Visit(ctx.Calc())
	}
	opText := ""
	if ctx.GetOp() != nil {
		opText = ctx.GetOp().GetText()
	}

	if op, found := operatorspb.Find(opText); found {
		lhs := p.Visit(ctx.Relation(0)).(*exprpb.Expr)
		rhs := p.Visit(ctx.Relation(1)).(*exprpb.Expr)
		return p.helper.newGlobalCall(ctx.GetOp(), op, lhs, rhs)
	}
	return p.helper.reportError(ctx, "operator not found")
}

// Visit a parse tree produced by CELParser#calc.
func (p *parser) VisitCalc(ctx *genpb.CalcContext) interface{} {
	if ctx.Unary() != nil {
		return p.Visit(ctx.Unary())
	}
	opText := ""
	if ctx.GetOp() != nil {
		opText = ctx.GetOp().GetText()
	}
	if op, found := operatorspb.Find(opText); found {
		lhs := p.Visit(ctx.Calc(0)).(*exprpb.Expr)
		rhs := p.Visit(ctx.Calc(1)).(*exprpb.Expr)
		return p.helper.newGlobalCall(ctx.GetOp(), op, lhs, rhs)
	}
	return p.helper.reportError(ctx, "operator not found")
}

func (p *parser) VisitUnary(ctx *genpb.UnaryContext) interface{} {
	return p.helper.newLiteralString(ctx, "<<error>>")
}

// Visit a parse tree produced by CELParser#StatementExpr.
func (p *parser) VisitStatementExpr(ctx *genpb.StatementExprContext) interface{} {
	switch ctx.Statement().(type) {
	case *genpb.PrimaryExprContext:
		return p.VisitPrimaryExpr(ctx.Statement().(*genpb.PrimaryExprContext))
	case *genpb.SelectOrCallContext:
		return p.VisitSelectOrCall(ctx.Statement().(*genpb.SelectOrCallContext))
	case *genpb.IndexContext:
		return p.VisitIndex(ctx.Statement().(*genpb.IndexContext))
	case *genpb.CreateMessageContext:
		return p.VisitCreateMessage(ctx.Statement().(*genpb.CreateMessageContext))
	}
	return p.helper.reportError(ctx, "unsupported simple expression")
}

// Visit a parse tree produced by CELParser#LogicalNot.
func (p *parser) VisitLogicalNot(ctx *genpb.LogicalNotContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Statement())
	}
	target := p.Visit(ctx.Statement()).(*exprpb.Expr)
	return p.helper.newGlobalCall(ctx.GetOps()[0], operatorspb.LogicalNot, target)
}

func (p *parser) VisitNegate(ctx *genpb.NegateContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Statement())
	}
	target := p.Visit(ctx.Statement()).(*exprpb.Expr)
	return p.helper.newGlobalCall(ctx.GetOps()[0], operatorspb.Negate, target)
}

// Visit a parse tree produced by CELParser#SelectOrCall.
func (p *parser) VisitSelectOrCall(ctx *genpb.SelectOrCallContext) interface{} {
	operand := p.Visit(ctx.Statement()).(*exprpb.Expr)
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
func (p *parser) VisitPrimaryExpr(ctx *genpb.PrimaryExprContext) interface{} {
	switch ctx.Primary().(type) {
	case *genpb.NestedContext:
		return p.VisitNested(ctx.Primary().(*genpb.NestedContext))
	case *genpb.IdentOrGlobalCallContext:
		return p.VisitIdentOrGlobalCall(ctx.Primary().(*genpb.IdentOrGlobalCallContext))
	case *genpb.CreateListContext:
		return p.VisitCreateList(ctx.Primary().(*genpb.CreateListContext))
	case *genpb.CreateStructContext:
		return p.VisitCreateStruct(ctx.Primary().(*genpb.CreateStructContext))
	case *genpb.ConstantLiteralContext:
		return p.VisitConstantLiteral(ctx.Primary().(*genpb.ConstantLiteralContext))
	}

	return p.helper.reportError(ctx, "invalid primary expression")
}

// Visit a parse tree produced by CELParser#Index.
func (p *parser) VisitIndex(ctx *genpb.IndexContext) interface{} {
	target := p.Visit(ctx.Statement()).(*exprpb.Expr)
	index := p.Visit(ctx.GetIndex()).(*exprpb.Expr)
	return p.helper.newGlobalCall(ctx.GetOp(), operatorspb.Index, target, index)
}

// Visit a parse tree produced by CELParser#CreateMessage.
func (p *parser) VisitCreateMessage(ctx *genpb.CreateMessageContext) interface{} {
	target := p.Visit(ctx.Statement()).(*exprpb.Expr)
	if messageName, found := p.extractQualifiedName(target); found {
		entries := p.VisitIFieldInitializerList(ctx.GetEntries()).([]*exprpb.Expr_CreateStruct_Entry)
		return p.helper.newObject(ctx, messageName, entries...)
	}
	return p.helper.newExpr(ctx)
}

func (p *parser) VisitIFieldInitializerList(ctx genpb.IFieldInitializerListContext) interface{} {
	if ctx == nil || ctx.GetFields() == nil {
		return []*exprpb.Expr_CreateStruct_Entry{}
	}

	result := make([]*exprpb.Expr_CreateStruct_Entry, len(ctx.GetFields()))
	for i, f := range ctx.GetFields() {
		value := p.Visit(ctx.GetValues()[i]).(*exprpb.Expr)
		field := p.helper.newObjectField(ctx.GetCols()[i], f.GetText(), value)
		result[i] = field
	}
	return result
}

// Visit a parse tree produced by CELParser#IdentOrGlobalCall.
func (p *parser) VisitIdentOrGlobalCall(ctx *genpb.IdentOrGlobalCallContext) interface{} {
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
func (p *parser) VisitNested(ctx *genpb.NestedContext) interface{} {
	return p.Visit(ctx.GetE())
}

// Visit a parse tree produced by CELParser#CreateList.
func (p *parser) VisitCreateList(ctx *genpb.CreateListContext) interface{} {
	return p.helper.newList(ctx, p.visitList(ctx.GetElems())...)
}

// Visit a parse tree produced by CELParser#CreateStruct.
func (p *parser) VisitCreateStruct(ctx *genpb.CreateStructContext) interface{} {
	entries := []*exprpb.Expr_CreateStruct_Entry{}
	if ctx.GetEntries() != nil {
		entries = p.Visit(ctx.GetEntries()).([]*exprpb.Expr_CreateStruct_Entry)
	}
	return p.helper.newMap(ctx.GetStart(), entries...)
}

// Visit a parse tree produced by CELParser#ConstantLiteral.
func (p *parser) VisitConstantLiteral(ctx *genpb.ConstantLiteralContext) interface{} {
	switch ctx.Literal().(type) {
	case *genpb.IntContext:
		return p.VisitInt(ctx.Literal().(*genpb.IntContext))
	case *genpb.UintContext:
		return p.VisitUint(ctx.Literal().(*genpb.UintContext))
	case *genpb.DoubleContext:
		return p.VisitDouble(ctx.Literal().(*genpb.DoubleContext))
	case *genpb.StringContext:
		return p.VisitString(ctx.Literal().(*genpb.StringContext))
	case *genpb.BytesContext:
		return p.VisitBytes(ctx.Literal().(*genpb.BytesContext))
	case *genpb.BoolFalseContext:
		return p.VisitBoolFalse(ctx.Literal().(*genpb.BoolFalseContext))
	case *genpb.BoolTrueContext:
		return p.VisitBoolTrue(ctx.Literal().(*genpb.BoolTrueContext))
	case *genpb.NullContext:
		return p.VisitNull(ctx.Literal().(*genpb.NullContext))
	}
	return p.helper.reportError(ctx, "invalid literal")
}

// Visit a parse tree produced by CELParser#exprList.
func (p *parser) VisitExprList(ctx *genpb.ExprListContext) interface{} {
	if ctx == nil || ctx.GetE() == nil {
		return []*exprpb.Expr{}
	}

	result := make([]*exprpb.Expr, len(ctx.GetE()))
	for i, e := range ctx.GetE() {
		exp := p.Visit(e).(*exprpb.Expr)
		result[i] = exp
	}
	return result
}

// Visit a parse tree produced by CELParser#mapInitializerList.
func (p *parser) VisitMapInitializerList(ctx *genpb.MapInitializerListContext) interface{} {
	if ctx == nil || ctx.GetKeys() == nil {
		return []*exprpb.Expr_CreateStruct_Entry{}
	}

	result := make([]*exprpb.Expr_CreateStruct_Entry, len(ctx.GetCols()))
	for i, col := range ctx.GetCols() {
		key := p.Visit(ctx.GetKeys()[i]).(*exprpb.Expr)
		value := p.Visit(ctx.GetValues()[i]).(*exprpb.Expr)
		entry := p.helper.newMapEntry(col, key, value)
		result[i] = entry
	}
	return result
}

// Visit a parse tree produced by CELParser#Int.
func (p *parser) VisitInt(ctx *genpb.IntContext) interface{} {
	i, err := strconv.ParseInt(ctx.GetTok().GetText(), 10, 64)
	if err != nil {
		return p.helper.reportError(ctx, "invalid int literal")
	}
	return p.helper.newLiteralInt(ctx, i)
}

// Visit a parse tree produced by CELParser#Uint.
func (p *parser) VisitUint(ctx *genpb.UintContext) interface{} {
	text := ctx.GetTok().GetText()
	text = text[:len(text)-1]
	i, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return p.helper.reportError(ctx, "invalid uint literal")
	}
	return p.helper.newLiteralUint(ctx, i)
}

// Visit a parse tree produced by CELParser#Double.
func (p *parser) VisitDouble(ctx *genpb.DoubleContext) interface{} {
	f, err := strconv.ParseFloat(ctx.GetTok().GetText(), 64)
	if err != nil {
		return p.helper.reportError(ctx, "invalid double literal")
	}
	return p.helper.newLiteralDouble(ctx, f)

}

// Visit a parse tree produced by CELParser#String.
func (p *parser) VisitString(ctx *genpb.StringContext) interface{} {
	s := p.unquote(ctx, ctx.GetText())
	return p.helper.newLiteralString(ctx, s)
}

// Visit a parse tree produced by CELParser#Bytes.
func (p *parser) VisitBytes(ctx *genpb.BytesContext) interface{} {
	// TODO(ozben): Not sure if this is the right encoding.
	b := []byte(p.unquote(ctx, ctx.GetTok().GetText()[1:]))
	return p.helper.newLiteralBytes(ctx, b)
}

// Visit a parse tree produced by CELParser#BoolTrue.
func (p *parser) VisitBoolTrue(ctx *genpb.BoolTrueContext) interface{} {
	return p.helper.newLiteralBool(ctx, true)
}

// Visit a parse tree produced by CELParser#BoolFalse.
func (p *parser) VisitBoolFalse(ctx *genpb.BoolFalseContext) interface{} {
	return p.helper.newLiteralBool(ctx, false)
}

// Visit a parse tree produced by CELParser#Null.
func (p *parser) VisitNull(ctx *genpb.NullContext) interface{} {
	return p.helper.newLiteral(ctx,
		&exprpb.Literal{
			LiteralKind: &exprpb.Literal_NullValue{
				NullValue: structpb.NullValue_NULL_VALUE}})
}

func (p *parser) visitList(ctx genpb.IExprListContext) []*exprpb.Expr {
	if ctx == nil {
		return []*exprpb.Expr{}
	}
	return p.visitSlice(ctx.GetE())
}

func (p *parser) visitSlice(expressions []genpb.IExprContext) []*exprpb.Expr {
	if expressions == nil {
		return []*exprpb.Expr{}
	}
	result := make([]*exprpb.Expr, len(expressions))
	for i, e := range expressions {
		ex := p.Visit(e).(*exprpb.Expr)
		result[i] = ex
	}
	return result
}

func (p *parser) extractQualifiedName(e *exprpb.Expr) (string, bool) {
	if e == nil {
		return "", false
	}
	switch e.ExprKind.(type) {
	case *exprpb.Expr_IdentExpr:
		return e.GetIdentExpr().Name, true
	case *exprpb.Expr_SelectExpr:
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
	text, err := unescape(value)
	if err != nil {
		p.helper.reportError(ctx, err.Error())
		return value
	}
	return text
}
