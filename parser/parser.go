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
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/parser/gen"

	structpb "github.com/golang/protobuf/ptypes/struct"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// reservedIds are not legal to use as variables.  We exclude them post-parse, as they *are* valid
// field names for protos, and it would complicate the grammar to distinguish the cases.
var reservedIds = []string{
	"as", "break", "const", "continue", "else", "false", "for", "function", "if",
	"import", "in", "let", "loop", "package", "namespace", "null", "return",
	"true", "var", "void", "while",
}

// Parse converts a source input a parsed expression.
// This function calls ParseWithMacros with AllMacros.
func Parse(source common.Source) (*exprpb.ParsedExpr, *common.Errors) {
	return ParseWithMacros(source, AllMacros)
}

// ParseWithMacros converts a source input and macros set to a parsed expression.
func ParseWithMacros(source common.Source, macros []Macro) (*exprpb.ParsedExpr, *common.Errors) {
	macroMap := make(map[string]Macro)
	for _, m := range macros {
		macroMap[m.MacroKey()] = m
	}
	p := parser{
		errors: &parseErrors{common.NewErrors(source)},
		helper: newParserHelper(source),
		macros: macroMap,
	}
	e := p.parse(source.Content())
	return &exprpb.ParsedExpr{
		Expr:       e,
		SourceInfo: p.helper.getSourceInfo(),
	}, p.errors.Errors
}

type parser struct {
	gen.BaseCELVisitor
	errors *parseErrors
	helper *parserHelper
	macros map[string]Macro
}

var _ gen.CELVisitor = (*parser)(nil)

func (p *parser) parse(expression string) *exprpb.Expr {
	stream := antlr.NewInputStream(expression)
	lexer := gen.NewCELLexer(stream)
	prsr := gen.NewCELParser(antlr.NewCommonTokenStream(lexer, 0))

	lexer.RemoveErrorListeners()
	prsr.RemoveErrorListeners()
	lexer.AddErrorListener(p)
	prsr.AddErrorListener(p)

	return p.Visit(prsr.Start()).(*exprpb.Expr)
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
	case *gen.MemberExprContext:
		return p.VisitMemberExpr(tree.(*gen.MemberExprContext))
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
	case *gen.CreateListContext:
		return p.VisitCreateList(tree.(*gen.CreateListContext))
	case *gen.CreateMessageContext:
		return p.VisitCreateMessage(tree.(*gen.CreateMessageContext))
	case *gen.CreateStructContext:
		return p.VisitCreateStruct(tree.(*gen.CreateStructContext))
	}

	// Report at least one error if the parser reaches an unknown parse element.
	// Typically, this happens if the parser has already encountered a syntax error elsewhere.
	if len(p.errors.GetErrors()) == 0 {
		txt := "<<nil>>"
		if tree != nil {
			txt = fmt.Sprintf("<<%T>>", tree)
		}
		return p.reportError(common.NoLocation, "unknown parse element encountered: %s", txt)
	}
	return p.helper.newExpr(common.NoLocation)

}

// Visit a parse tree produced by CELParser#start.
func (p *parser) VisitStart(ctx *gen.StartContext) interface{} {
	return p.Visit(ctx.Expr())
}

// Visit a parse tree produced by CELParser#expr.
func (p *parser) VisitExpr(ctx *gen.ExprContext) interface{} {
	result := p.Visit(ctx.GetE()).(*exprpb.Expr)
	if ctx.GetOp() == nil {
		return result
	}
	opID := p.helper.id(ctx.GetOp())
	ifTrue := p.Visit(ctx.GetE1()).(*exprpb.Expr)
	ifFalse := p.Visit(ctx.GetE2()).(*exprpb.Expr)
	return p.globalCallOrMacro(opID, operators.Conditional, result, ifTrue, ifFalse)
}

// Visit a parse tree produced by CELParser#conditionalOr.
func (p *parser) VisitConditionalOr(ctx *gen.ConditionalOrContext) interface{} {
	result := p.Visit(ctx.GetE()).(*exprpb.Expr)
	if ctx.GetOps() == nil {
		return result
	}
	b := newBalancer(p.helper, operators.LogicalOr, result)
	rest := ctx.GetE1()
	for i, op := range ctx.GetOps() {
		if i >= len(rest) {
			return p.reportError(ctx, "unexpected character, wanted '||'")
		}
		next := p.Visit(rest[i]).(*exprpb.Expr)
		opID := p.helper.id(op)
		b.addTerm(opID, next)
	}
	return b.balance()
}

// Visit a parse tree produced by CELParser#conditionalAnd.
func (p *parser) VisitConditionalAnd(ctx *gen.ConditionalAndContext) interface{} {
	result := p.Visit(ctx.GetE()).(*exprpb.Expr)
	if ctx.GetOps() == nil {
		return result
	}
	b := newBalancer(p.helper, operators.LogicalAnd, result)
	rest := ctx.GetE1()
	for i, op := range ctx.GetOps() {
		if i >= len(rest) {
			return p.reportError(ctx, "unexpected character, wanted '&&'")
		}
		next := p.Visit(rest[i]).(*exprpb.Expr)
		opID := p.helper.id(op)
		b.addTerm(opID, next)
	}
	return b.balance()
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
		lhs := p.Visit(ctx.Relation(0)).(*exprpb.Expr)
		opID := p.helper.id(ctx.GetOp())
		rhs := p.Visit(ctx.Relation(1)).(*exprpb.Expr)
		return p.globalCallOrMacro(opID, op, lhs, rhs)
	}
	return p.reportError(ctx, "operator not found")
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
		lhs := p.Visit(ctx.Calc(0)).(*exprpb.Expr)
		opID := p.helper.id(ctx.GetOp())
		rhs := p.Visit(ctx.Calc(1)).(*exprpb.Expr)
		return p.globalCallOrMacro(opID, op, lhs, rhs)
	}
	return p.reportError(ctx, "operator not found")
}

func (p *parser) VisitUnary(ctx *gen.UnaryContext) interface{} {
	return p.helper.newLiteralString(ctx, "<<error>>")
}

// Visit a parse tree produced by CELParser#MemberExpr.
func (p *parser) VisitMemberExpr(ctx *gen.MemberExprContext) interface{} {
	switch ctx.Member().(type) {
	case *gen.PrimaryExprContext:
		return p.VisitPrimaryExpr(ctx.Member().(*gen.PrimaryExprContext))
	case *gen.SelectOrCallContext:
		return p.VisitSelectOrCall(ctx.Member().(*gen.SelectOrCallContext))
	case *gen.IndexContext:
		return p.VisitIndex(ctx.Member().(*gen.IndexContext))
	case *gen.CreateMessageContext:
		return p.VisitCreateMessage(ctx.Member().(*gen.CreateMessageContext))
	}
	return p.reportError(ctx, "unsupported simple expression")
}

// Visit a parse tree produced by CELParser#LogicalNot.
func (p *parser) VisitLogicalNot(ctx *gen.LogicalNotContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Member())
	}
	opID := p.helper.id(ctx.GetOps()[0])
	target := p.Visit(ctx.Member()).(*exprpb.Expr)
	return p.globalCallOrMacro(opID, operators.LogicalNot, target)
}

func (p *parser) VisitNegate(ctx *gen.NegateContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Member())
	}
	opID := p.helper.id(ctx.GetOps()[0])
	target := p.Visit(ctx.Member()).(*exprpb.Expr)
	return p.globalCallOrMacro(opID, operators.Negate, target)
}

// Visit a parse tree produced by CELParser#SelectOrCall.
func (p *parser) VisitSelectOrCall(ctx *gen.SelectOrCallContext) interface{} {
	operand := p.Visit(ctx.Member()).(*exprpb.Expr)
	// Handle the error case where no valid identifier is specified.
	if ctx.GetId() == nil {
		return p.helper.newExpr(ctx)
	}
	id := ctx.GetId().GetText()
	if ctx.GetOpen() != nil {
		opID := p.helper.id(ctx.GetOpen())
		return p.receiverCallOrMacro(opID, id, operand, p.visitList(ctx.GetArgs())...)
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

	return p.reportError(ctx, "invalid primary expression")
}

// Visit a parse tree produced by CELParser#Index.
func (p *parser) VisitIndex(ctx *gen.IndexContext) interface{} {
	target := p.Visit(ctx.Member()).(*exprpb.Expr)
	opID := p.helper.id(ctx.GetOp())
	index := p.Visit(ctx.GetIndex()).(*exprpb.Expr)
	return p.globalCallOrMacro(opID, operators.Index, target, index)
}

// Visit a parse tree produced by CELParser#CreateMessage.
func (p *parser) VisitCreateMessage(ctx *gen.CreateMessageContext) interface{} {
	target := p.Visit(ctx.Member()).(*exprpb.Expr)
	objID := p.helper.id(ctx.GetOp())
	if messageName, found := p.extractQualifiedName(target); found {
		entries := p.VisitIFieldInitializerList(ctx.GetEntries()).([]*exprpb.Expr_CreateStruct_Entry)
		return p.helper.newObject(objID, messageName, entries...)
	}
	return p.helper.newExpr(objID)
}

// Visit a parse tree of field initializers.
func (p *parser) VisitIFieldInitializerList(ctx gen.IFieldInitializerListContext) interface{} {
	if ctx == nil || ctx.GetFields() == nil {
		// This is the result of a syntax error handled elswhere, return empty.
		return []*exprpb.Expr_CreateStruct_Entry{}
	}

	result := make([]*exprpb.Expr_CreateStruct_Entry, len(ctx.GetFields()))
	cols := ctx.GetCols()
	vals := ctx.GetValues()
	for i, f := range ctx.GetFields() {
		if i >= len(cols) || i >= len(vals) {
			// This is the result of a syntax error detected elsewhere.
			return []*exprpb.Expr_CreateStruct_Entry{}
		}
		initID := p.helper.id(cols[i])
		value := p.Visit(vals[i]).(*exprpb.Expr)
		field := p.helper.newObjectField(initID, f.GetText(), value)
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
	// Handle reserved identifiers.
	id := ctx.GetId().GetText()
	for _, r := range reservedIds {
		if r == id {
			return p.reportError(ctx, "reserved identifier: %s", r)
		}
	}
	identName += id
	if ctx.GetOp() != nil {
		opID := p.helper.id(ctx.GetOp())
		return p.globalCallOrMacro(opID, identName, p.visitList(ctx.GetArgs())...)
	}
	return p.helper.newIdent(ctx.GetId(), identName)
}

// Visit a parse tree produced by CELParser#Nested.
func (p *parser) VisitNested(ctx *gen.NestedContext) interface{} {
	return p.Visit(ctx.GetE())
}

// Visit a parse tree produced by CELParser#CreateList.
func (p *parser) VisitCreateList(ctx *gen.CreateListContext) interface{} {
	listID := p.helper.id(ctx.GetOp())
	return p.helper.newList(listID, p.visitList(ctx.GetElems())...)
}

// Visit a parse tree produced by CELParser#CreateStruct.
func (p *parser) VisitCreateStruct(ctx *gen.CreateStructContext) interface{} {
	structID := p.helper.id(ctx.GetOp())
	entries := []*exprpb.Expr_CreateStruct_Entry{}
	if ctx.GetEntries() != nil {
		entries = p.Visit(ctx.GetEntries()).([]*exprpb.Expr_CreateStruct_Entry)
	}
	return p.helper.newMap(structID, entries...)
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
	return p.reportError(ctx, "invalid literal")
}

// Visit a parse tree produced by CELParser#mapInitializerList.
func (p *parser) VisitMapInitializerList(ctx *gen.MapInitializerListContext) interface{} {
	if ctx == nil || ctx.GetKeys() == nil {
		// This is the result of a syntax error handled elswhere, return empty.
		return []*exprpb.Expr_CreateStruct_Entry{}
	}

	result := make([]*exprpb.Expr_CreateStruct_Entry, len(ctx.GetCols()))
	keys := ctx.GetKeys()
	vals := ctx.GetValues()
	for i, col := range ctx.GetCols() {
		colID := p.helper.id(col)
		if i >= len(keys) || i >= len(vals) {
			// This is the result of a syntax error detected elsewhere.
			return []*exprpb.Expr_CreateStruct_Entry{}
		}
		key := p.Visit(keys[i]).(*exprpb.Expr)
		value := p.Visit(vals[i]).(*exprpb.Expr)
		entry := p.helper.newMapEntry(colID, key, value)
		result[i] = entry
	}
	return result
}

// Visit a parse tree produced by CELParser#Int.
func (p *parser) VisitInt(ctx *gen.IntContext) interface{} {
	text := ctx.GetTok().GetText()
	if ctx.GetSign() != nil {
		text = ctx.GetSign().GetText() + text
	}
	i, err := strconv.ParseInt(text, 10, 64)
	if err != nil {
		return p.reportError(ctx, "invalid int literal")
	}
	return p.helper.newLiteralInt(ctx, i)
}

// Visit a parse tree produced by CELParser#Uint.
func (p *parser) VisitUint(ctx *gen.UintContext) interface{} {
	text := ctx.GetTok().GetText()
	// trim the 'u' designator included in the uint literal.
	text = text[:len(text)-1]
	i, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return p.reportError(ctx, "invalid uint literal")
	}
	return p.helper.newLiteralUint(ctx, i)
}

// Visit a parse tree produced by CELParser#Double.
func (p *parser) VisitDouble(ctx *gen.DoubleContext) interface{} {
	txt := ctx.GetTok().GetText()
	if ctx.GetSign() != nil {
		txt = ctx.GetSign().GetText() + txt
	}
	f, err := strconv.ParseFloat(txt, 64)
	if err != nil {
		return p.reportError(ctx, "invalid double literal")
	}
	return p.helper.newLiteralDouble(ctx, f)

}

// Visit a parse tree produced by CELParser#String.
func (p *parser) VisitString(ctx *gen.StringContext) interface{} {
	s := p.unquote(ctx, ctx.GetText(), false)
	return p.helper.newLiteralString(ctx, s)
}

// Visit a parse tree produced by CELParser#Bytes.
func (p *parser) VisitBytes(ctx *gen.BytesContext) interface{} {
	b := []byte(p.unquote(ctx, ctx.GetTok().GetText()[1:], true))
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
		&exprpb.Constant{
			ConstantKind: &exprpb.Constant_NullValue{
				NullValue: structpb.NullValue_NULL_VALUE}})
}

func (p *parser) visitList(ctx gen.IExprListContext) []*exprpb.Expr {
	if ctx == nil {
		return []*exprpb.Expr{}
	}
	return p.visitSlice(ctx.GetE())
}

func (p *parser) visitSlice(expressions []gen.IExprContext) []*exprpb.Expr {
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
	p.reportError(location, "expected a qualified name")
	return "", false
}

func (p *parser) unquote(ctx interface{}, value string, isBytes bool) string {
	text, err := unescape(value, isBytes)
	if err != nil {
		p.reportError(ctx, "%s", err.Error())
		return value
	}
	return text
}

func (p *parser) reportError(ctx interface{}, format string, args ...interface{}) *exprpb.Expr {
	var location common.Location
	switch ctx.(type) {
	case common.Location:
		location = ctx.(common.Location)
	case antlr.Token, antlr.ParserRuleContext:
		err := p.helper.newExpr(ctx)
		location = p.helper.getLocation(err.Id)
	}
	err := p.helper.newExpr(ctx)
	// Provide arguments to the report error.
	p.errors.ReportError(location, format, args...)
	return err
}

// ANTLR Parse listener implementations
func (p *parser) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// TODO: Snippet
	l := p.helper.source.NewLocation(line, column)
	p.errors.syntaxError(l, msg)
}

func (p *parser) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	// Intentional
}

func (p *parser) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	// Intentional
}

func (p *parser) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	// Intentional
}

func (p *parser) globalCallOrMacro(exprID int64, function string, args ...*exprpb.Expr) *exprpb.Expr {
	if expr, found := p.expandMacro(exprID, function, nil, args...); found {
		return expr
	}
	return p.helper.newGlobalCall(exprID, function, args...)
}

func (p *parser) receiverCallOrMacro(exprID int64, function string, target *exprpb.Expr, args ...*exprpb.Expr) *exprpb.Expr {
	if expr, found := p.expandMacro(exprID, function, target, args...); found {
		return expr
	}
	return p.helper.newReceiverCall(exprID, function, target, args...)
}

func (p *parser) expandMacro(exprID int64, function string, target *exprpb.Expr, args ...*exprpb.Expr) (*exprpb.Expr, bool) {
	macro, found := p.macros[makeMacroKey(function, len(args), target != nil)]
	if !found {
		macro, found = p.macros[makeVarArgMacroKey(function, target != nil)]
		if !found {
			return nil, false
		}
	}
	eh := exprHelperPool.Get().(*exprHelper)
	defer exprHelperPool.Put(eh)
	eh.parserHelper = p.helper
	eh.id = exprID
	expr, err := macro.Expander()(eh, target, args)
	if err != nil {
		if err.Location != nil {
			return p.reportError(err.Location, err.Message), true
		}
		return p.reportError(p.helper.getLocation(exprID), err.Message), true
	}
	return expr, true
}
