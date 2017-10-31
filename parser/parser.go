package parser

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"celgo/ast"
	"celgo/common"
	"celgo/operators"
	"celgo/parser/gen"
	"github.com/antlr/antlr4/runtime/Go/antlr"
)

func ParseText(expression string) (ast.Expression, *common.Errors) {
	errors := common.NewErrors()
	expr := Parse(errors, expression, "<input>", AllMacros)
	return expr, errors
}

func Parse(errors *common.Errors, expression string, sourceName string, macros Macros) ast.Expression {
	macroMap := make(map[string]Macro)
	for _, m := range macros {
		macroMap[makeMacroKey(m.name, m.args, m.instanceStyle)] = m
	}

	p := parser{
		errors: &ParseErrors{errors},
		source: common.NewTextSource(sourceName, expression),
		macros: macroMap,
		nextId: 1,
	}

	return p.parse(expression)
}

type parser struct {
	gen.BaseCELVisitor

	errors *ParseErrors
	source common.Source
	macros map[string]Macro
	nextId int64
}

var _ gen.CELVisitor = (*parser)(nil)

func (p *parser) parse(expression string) ast.Expression {
	stream := antlr.NewInputStream(expression)
	lexer := gen.NewCELLexer(stream)
	prsr := gen.NewCELParser(antlr.NewCommonTokenStream(lexer, 0))

	lexer.RemoveErrorListeners()
	prsr.RemoveErrorListeners()
	lexer.AddErrorListener(p)
	prsr.AddErrorListener(p)

	return p.Visit(prsr.Start()).(ast.Expression)
}

// Listener implementations
func (p *parser) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// TODO: Snippet
	l := common.NewLocation(p.source, line, column+1)
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

// Visitor implementations

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
		return &ast.ErrorExpression{}
	}

	text := "<<nil>>"
	if tree != nil {
		text = tree.GetText()
	}
	panic(fmt.Sprintf("Unknown parsetree type: '%+v': %+v [%s]", reflect.TypeOf(tree), tree, text))
}

// Visit a parse tree produced by CELParser#start.
func (p *parser) VisitStart(ctx *gen.StartContext) interface{} {
	return p.Visit(ctx.Expr())
}

// Visit a parse tree produced by CELParser#expr.
func (p *parser) VisitExpr(ctx *gen.ExprContext) interface{} {
	result := p.Visit(ctx.GetE()).(ast.Expression)
	if ctx.GetOp() == nil {
		return result
	}

	ifTrue := p.Visit(ctx.GetE1()).(ast.Expression)
	ifFalse := p.Visit(ctx.GetE2()).(ast.Expression)
	return ast.NewCallFunction(p.id(), p.location(ctx.GetStart()), operators.Conditional, result, ifTrue, ifFalse)
}

// Visit a parse tree produced by CELParser#conditionalOr.
func (p *parser) VisitConditionalOr(ctx *gen.ConditionalOrContext) interface{} {
	result := p.Visit(ctx.GetE()).(ast.Expression)
	if ctx.GetOps() == nil {
		return result
	}

	for i, _ := range ctx.GetOps() {
		next := p.Visit(ctx.GetE1()[i]).(ast.Expression)
		result = ast.NewCallFunction(p.id(), p.location(ctx.GetStart()), operators.LogicalOr, result, next)
	}

	return result
}

// Visit a parse tree produced by CELParser#conditionalAnd.
func (p *parser) VisitConditionalAnd(ctx *gen.ConditionalAndContext) interface{} {
	result := p.Visit(ctx.GetE()).(ast.Expression)
	if ctx.GetOps() == nil {
		return result
	}

	for i, _ := range ctx.GetOps() {
		next := p.Visit(ctx.GetE1()[i]).(ast.Expression)
		result = ast.NewCallFunction(p.id(), p.location(ctx.GetStart()), operators.LogicalAnd, result, next)
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
		lhs := p.Visit(ctx.Relation(0)).(ast.Expression)
		rhs := p.Visit(ctx.Relation(1)).(ast.Expression)
		return ast.NewCallFunction(p.id(), p.location(ctx.GetOp()), op, lhs, rhs)
	}

	return &ast.ErrorExpression{}
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
		lhs := p.Visit(ctx.Calc(0)).(ast.Expression)
		rhs := p.Visit(ctx.Calc(1)).(ast.Expression)
		return ast.NewCallFunction(p.id(), p.location(ctx.GetOp()), op, lhs, rhs)
	}

	return &ast.ErrorExpression{}
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

	return &ast.ErrorExpression{}
}

// Visit a parse tree produced by CELParser#LogicalNot.
func (p *parser) VisitLogicalNot(ctx *gen.LogicalNotContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Statement())
	}

	target := p.Visit(ctx.Statement()).(ast.Expression)

	return ast.NewCallFunction(p.id(), p.location(ctx.GetStart()), operators.LogicalNot, target)

}

func (p *parser) VisitNegate(ctx *gen.NegateContext) interface{} {
	if len(ctx.GetOps())%2 == 0 {
		return p.Visit(ctx.Statement())
	}

	target := p.Visit(ctx.Statement()).(ast.Expression)

	return ast.NewCallFunction(p.id(), p.location(ctx.GetStart()), operators.Negate, target)
}

// Visit a parse tree produced by CELParser#SelectOrCall.
func (p *parser) VisitSelectOrCall(ctx *gen.SelectOrCallContext) interface{} {
	target := p.Visit(ctx.Statement()).(ast.Expression)
	id := ctx.GetId().GetText()

	if ctx.GetOpen() != nil {
		return p.newCallMethodOrMacro(p.location(ctx.GetOp()), id, target, p.visitList(ctx.GetArgs())...)
	}

	return ast.NewSelect(p.id(), p.location(ctx.GetOp()), target, id, false)
}

// Visit a parse tree produced by CELParser#PrimaryExpr.
func (p *parser) VisitPrimaryExpr(ctx *gen.PrimaryExprContext) interface{} {
	switch ctx.Primary().(type) {
	case *gen.NestedContext:
		return p.VisitNested(ctx.Primary().(*gen.NestedContext))
	case *gen.IdentOrGlobalCallContext:
		return p.VisitIdentOrGlobalCall(ctx.Primary().(*gen.IdentOrGlobalCallContext))
		// TODO(b/67832261): Deprecate the 'in' function.
	case *gen.DeprecatedInContext:
		return p.VisitDeprecatedIn(ctx.Primary().(*gen.DeprecatedInContext))
	case *gen.CreateListContext:
		return p.VisitCreateList(ctx.Primary().(*gen.CreateListContext))
	case *gen.CreateStructContext:
		return p.VisitCreateStruct(ctx.Primary().(*gen.CreateStructContext))
	case *gen.ConstantLiteralContext:
		return p.VisitConstantLiteral(ctx.Primary().(*gen.ConstantLiteralContext))
	}

	return &ast.ErrorExpression{}
}

// Visit a parse tree produced by CELParser#Index.
func (p *parser) VisitIndex(ctx *gen.IndexContext) interface{} {
	target := p.Visit(ctx.Statement()).(ast.Expression)
	index := p.Visit(ctx.GetIndex()).(ast.Expression)
	return ast.NewCallFunction(p.id(), p.location(ctx.GetOp()), operators.Index, target, index)
}

// Visit a parse tree produced by CELParser#CreateMessage.
func (p *parser) VisitCreateMessage(ctx *gen.CreateMessageContext) interface{} {
	target := p.Visit(ctx.Statement()).(ast.Expression)

	messageName, found := p.extractQualifiedName(target)
	if !found {
		return &ast.ErrorExpression{}
	}

	entries := p.VisitIFieldInitializerList(ctx.GetEntries()).([]*ast.FieldEntry)

	return ast.NewCreateMessage(p.id(), p.location(ctx.GetStart()), messageName, entries...)
}

func (p *parser) VisitIFieldInitializerList(ctx gen.IFieldInitializerListContext) interface{} {
	if ctx == nil || ctx.GetFields() == nil {
		return []*ast.FieldEntry{}
	}

	result := make([]*ast.FieldEntry, len(ctx.GetFields()))
	for i, f := range ctx.GetFields() {
		value := p.Visit(ctx.GetValues()[i]).(ast.Expression)

		field := ast.NewFieldEntry(p.id(), p.location(ctx.GetCols()[i]), f.GetText(), value)
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
	identName += ctx.GetId().GetText()

	if ctx.GetOp() != nil {
		return p.newCallFunctionOrMacro(p.location(ctx.GetOp()), identName, p.visitList(ctx.GetArgs())...)
	}

	return ast.NewIdent(p.id(), p.location(ctx.GetStart()), identName)
}

// Visit a parse tree produced by CELParser#DeprecatedIn.
func (p *parser) VisitDeprecatedIn(ctx *gen.DeprecatedInContext) interface{} {
	identName := ctx.GetId().GetText()
	return p.newCallFunctionOrMacro(p.location(ctx.GetOp()), identName, p.visitList(ctx.GetArgs())...)
}

// Visit a parse tree produced by CELParser#Nested.
func (p *parser) VisitNested(ctx *gen.NestedContext) interface{} {
	return p.Visit(ctx.GetE())
}

// Visit a parse tree produced by CELParser#CreateList.
func (p *parser) VisitCreateList(ctx *gen.CreateListContext) interface{} {
	return ast.NewCreateList(p.id(), p.location(ctx.GetStart()), p.visitList(ctx.GetElems())...)
}

// Visit a parse tree produced by CELParser#CreateStruct.
func (p *parser) VisitCreateStruct(ctx *gen.CreateStructContext) interface{} {
	entries := []*ast.StructEntry{}
	if ctx.GetEntries() != nil {
		entries = p.Visit(ctx.GetEntries()).([]*ast.StructEntry)
	}

	return ast.NewCreateStruct(p.id(), p.location(ctx.GetStart()), entries...)
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

	return &ast.ErrorExpression{}
}

// Visit a parse tree produced by CELParser#exprList.
func (p *parser) VisitExprList(ctx *gen.ExprListContext) interface{} {
	if ctx == nil || ctx.GetE() == nil {
		return []ast.Expression{}
	}

	result := make([]ast.Expression, len(ctx.GetE()))
	for i, e := range ctx.GetE() {
		exp := p.Visit(e).(ast.Expression)
		result[i] = exp
	}
	return result
}

// Visit a parse tree produced by CELParser#mapInitializerList.
func (p *parser) VisitMapInitializerList(ctx *gen.MapInitializerListContext) interface{} {

	if ctx == nil || ctx.GetKeys() == nil {
		return []*ast.StructEntry{}
	}

	result := make([]*ast.StructEntry, len(ctx.GetCols()))

	for i, _ := range ctx.GetCols() {
		key := p.Visit(ctx.GetKeys()[i]).(ast.Expression)
		value := p.Visit(ctx.GetValues()[i]).(ast.Expression)

		entry := ast.NewStructEntry(p.id(), p.location(ctx.GetStart()), key, value)
		result[i] = entry
	}

	return result
}

// Visit a parse tree produced by CELParser#Int.
func (p *parser) VisitInt(ctx *gen.IntContext) interface{} {
	i, err := strconv.ParseInt(ctx.GetTok().GetText(), 10, 64)
	if err != nil {
		return &ast.ErrorExpression{}
	}

	return ast.NewInt64Constant(p.id(), p.location(ctx.GetStart()), i)
}

// Visit a parse tree produced by CELParser#Uint.
func (p *parser) VisitUint(ctx *gen.UintContext) interface{} {
	text := ctx.GetTok().GetText()
	text = text[:len(text)-1]
	i, err := strconv.ParseUint(text, 10, 64)
	if err != nil {
		return &ast.ErrorExpression{}
	}

	return ast.NewUint64Constant(p.id(), p.location(ctx.GetStart()), i)
}

// Visit a parse tree produced by CELParser#Double.
func (p *parser) VisitDouble(ctx *gen.DoubleContext) interface{} {
	f, err := strconv.ParseFloat(ctx.GetTok().GetText(), 64)
	if err != nil {
		return &ast.ErrorExpression{}
	}
	return ast.NewDoubleConstant(p.id(), p.location(ctx.GetStart()), f)

}

// Visit a parse tree produced by CELParser#String.
func (p *parser) VisitString(ctx *gen.StringContext) interface{} {
	return ast.NewStringConstant(p.id(), p.location(ctx.GetStart()), unquote(ctx.GetText()))
}

// Visit a parse tree produced by CELParser#Bytes.
func (p *parser) VisitBytes(ctx *gen.BytesContext) interface{} {
	// TODO(ozben): Not sure if this is the right encoding.
	b := []byte(unquote(ctx.GetTok().GetText()[1:]))
	return ast.NewBytesConstant(p.id(), p.location(ctx.GetStart()), b)
}

// Visit a parse tree produced by CELParser#BoolTrue.
func (p *parser) VisitBoolTrue(ctx *gen.BoolTrueContext) interface{} {
	return ast.NewBoolConstant(p.id(), p.location(ctx.GetStart()), true)
}

// Visit a parse tree produced by CELParser#BoolFalse.
func (p *parser) VisitBoolFalse(ctx *gen.BoolFalseContext) interface{} {
	return ast.NewBoolConstant(p.id(), p.location(ctx.GetStart()), false)
}

// Visit a parse tree produced by CELParser#Null.
func (p *parser) VisitNull(ctx *gen.NullContext) interface{} {
	return ast.NewNullConstant(p.id(), p.location(ctx.GetStart()))
}

func (p *parser) visitList(ctx gen.IExprListContext) []ast.Expression {
	if ctx == nil {
		return []ast.Expression{}
	}

	return p.visitSlice(ctx.GetE())
}

func (p *parser) visitSlice(expressions []gen.IExprContext) []ast.Expression {
	if expressions == nil {
		return []ast.Expression{}
	}

	result := make([]ast.Expression, len(expressions))
	for i, e := range expressions {
		expr := p.Visit(e).(ast.Expression)
		result[i] = expr
	}

	return result
}

func (p *parser) location(token antlr.Token) common.Location {
	return common.NewLocation(p.source, token.GetLine(), token.GetColumn()+1)
}

func (p *parser) extractQualifiedName(e ast.Expression) (string, bool) {
	if e == nil {
		return "", false
	}

	switch e.(type) {
	case *ast.IdentExpression:
		return e.(*ast.IdentExpression).Name, true
	case *ast.SelectExpression:
		s := e.(*ast.SelectExpression)
		if prefix, found := p.extractQualifiedName(s.Target); found {
			return prefix + "." + s.Field, true
		}
	}

	p.errors.notAQualifiedName(e.Location())
	return "", false
}

func unquote(str string) string {
	// TODO: escape sequences
	if strings.ToLower(str)[0] == 'r' {
		str = str[1:]
	}

	if strings.HasPrefix(str, `""""`) {
		return str[3 : len(str)-3]
	}
	if strings.HasPrefix(str, "''") {
		return str[3 : len(str)-3]
	}
	if strings.HasPrefix(str, `"`) {
		return str[1 : len(str)-1]
	}
	if strings.HasPrefix(str, `'`) {
		return str[1 : len(str)-1]
	}

	panic("Unable to unquote")
}

func (p *parser) newCallMethodOrMacro(loc common.Location, name string, target ast.Expression, args ...ast.Expression) ast.Expression {
	macro, found := p.macros[makeMacroKey(name, len(args), true)]
	if !found {
		return ast.NewCallMethod(p.id(), loc, name, target, args...)
	}

	return macro.expander(p, loc, target, args)
}

func (p *parser) newCallFunctionOrMacro(loc common.Location, name string, args ...ast.Expression) ast.Expression {
	macro, found := p.macros[makeMacroKey(name, len(args), false)]
	if !found {
		return ast.NewCallFunction(p.id(), loc, name, args...)
	}

	return macro.expander(p, loc, nil, args)
}

func (p *parser) id() int64 {
	id := p.nextId
	p.nextId++
	return id
}
