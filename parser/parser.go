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
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/runes"
	"github.com/google/cel-go/parser/gen"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// Parser encapsulates the context necessary to perform parsing for different expressions.
type Parser struct {
	options
}

// NewParser builds and returns a new Parser using the provided options.
func NewParser(opts ...Option) (*Parser, error) {
	p := &Parser{}
	for _, opt := range opts {
		if err := opt(&p.options); err != nil {
			return nil, err
		}
	}
	if p.maxRecursionDepth == 0 {
		p.maxRecursionDepth = 250
	}
	if p.maxRecursionDepth == -1 {
		p.maxRecursionDepth = int((^uint(0)) >> 1)
	}
	if p.errorRecoveryTokenLookaheadLimit == 0 {
		p.errorRecoveryTokenLookaheadLimit = 256
	}
	if p.errorRecoveryLimit == 0 {
		p.errorRecoveryLimit = 30
	}
	if p.errorRecoveryLimit == -1 {
		p.errorRecoveryLimit = int((^uint(0)) >> 1)
	}
	if p.expressionSizeCodePointLimit == 0 {
		p.expressionSizeCodePointLimit = 100_000
	}
	if p.expressionSizeCodePointLimit == -1 {
		p.expressionSizeCodePointLimit = int((^uint(0)) >> 1)
	}
	// Bool is false by default, so populateMacroCalls will be false by default
	return p, nil
}

// mustNewParser does the work of NewParser and panics if an error occurs.
//
// This function is only intended for internal use and is for backwards compatibility in Parse and
// ParseWithMacros, where we know the options will result in an error.
func mustNewParser(opts ...Option) *Parser {
	p, err := NewParser(opts...)
	if err != nil {
		panic(err)
	}
	return p
}

// Parse parses the expression represented by source and returns the result.
func (p *Parser) Parse(source common.Source) (*exprpb.ParsedExpr, *common.Errors) {
	impl := parser{
		errors:                           &parseErrors{common.NewErrors(source)},
		helper:                           newParserHelper(source),
		macros:                           p.macros,
		maxRecursionDepth:                p.maxRecursionDepth,
		errorRecoveryLimit:               p.errorRecoveryLimit,
		errorRecoveryLookaheadTokenLimit: p.errorRecoveryTokenLookaheadLimit,
		populateMacroCalls:               p.populateMacroCalls,
	}
	buf, ok := source.(runes.Buffer)
	if !ok {
		buf = runes.NewBuffer(source.Content())
	}
	var e *exprpb.Expr
	if buf.Len() > p.expressionSizeCodePointLimit {
		e = impl.reportError(common.NoLocation,
			"expression code point size exceeds limit: size: %d, limit %d",
			buf.Len(), p.expressionSizeCodePointLimit)
	} else {
		e = impl.parse(buf, source.Description())
	}
	return &exprpb.ParsedExpr{
		Expr:       e,
		SourceInfo: impl.helper.getSourceInfo(),
	}, impl.errors.Errors
}

// reservedIds are not legal to use as variables.  We exclude them post-parse, as they *are* valid
// field names for protos, and it would complicate the grammar to distinguish the cases.
var reservedIds = map[string]struct{}{
	"as":        {},
	"break":     {},
	"const":     {},
	"continue":  {},
	"else":      {},
	"false":     {},
	"for":       {},
	"function":  {},
	"if":        {},
	"import":    {},
	"in":        {},
	"let":       {},
	"loop":      {},
	"package":   {},
	"namespace": {},
	"null":      {},
	"return":    {},
	"true":      {},
	"var":       {},
	"void":      {},
	"while":     {},
}

// Parse converts a source input a parsed expression.
// This function calls ParseWithMacros with AllMacros.
//
// Deprecated: Use NewParser().Parse() instead.
func Parse(source common.Source) (*exprpb.ParsedExpr, *common.Errors) {
	return mustNewParser(Macros(AllMacros...)).Parse(source)
}

type recursionError struct {
	message string
}

// Error implements error.
func (re *recursionError) Error() string {
	return re.message
}

var _ error = &recursionError{}

type recursionListener struct {
	maxDepth      int
	ruleTypeDepth map[int]*int
}

func (rl *recursionListener) VisitTerminal(node antlr.TerminalNode) {}

func (rl *recursionListener) VisitErrorNode(node antlr.ErrorNode) {}

func (rl *recursionListener) EnterEveryRule(ctx antlr.ParserRuleContext) {
	if ctx == nil {
		return
	}
	ruleIndex := ctx.GetRuleIndex()
	depth, found := rl.ruleTypeDepth[ruleIndex]
	if !found {
		var counter = 1
		rl.ruleTypeDepth[ruleIndex] = &counter
		depth = &counter
	} else {
		*depth++
	}
	if *depth > rl.maxDepth {
		panic(&recursionError{
			message: fmt.Sprintf("expression recursion limit exceeded: %d", rl.maxDepth),
		})
	}
}

func (rl *recursionListener) ExitEveryRule(ctx antlr.ParserRuleContext) {
	if ctx == nil {
		return
	}
	ruleIndex := ctx.GetRuleIndex()
	if depth, found := rl.ruleTypeDepth[ruleIndex]; found && *depth > 0 {
		*depth--
	}
}

var _ antlr.ParseTreeListener = &recursionListener{}

type recoveryLimitError struct {
	message string
}

// Error implements error.
func (rl *recoveryLimitError) Error() string {
	return rl.message
}

type lookaheadLimitError struct {
	message string
}

func (ll *lookaheadLimitError) Error() string {
	return ll.message
}

var _ error = &recoveryLimitError{}

type recoveryLimitErrorStrategy struct {
	*antlr.DefaultErrorStrategy
	errorRecoveryLimit               int
	errorRecoveryTokenLookaheadLimit int
	recoveryAttempts                 int
}

type lookaheadConsumer struct {
	antlr.Parser
	errorRecoveryTokenLookaheadLimit int
	lookaheadAttempts                int
}

func (lc *lookaheadConsumer) Consume() antlr.Token {
	if lc.lookaheadAttempts >= lc.errorRecoveryTokenLookaheadLimit {
		panic(&lookaheadLimitError{
			message: fmt.Sprintf("error recovery token lookahead limit exceeded: %d", lc.errorRecoveryTokenLookaheadLimit),
		})
	}
	lc.lookaheadAttempts++
	return lc.Parser.Consume()
}

func (rl *recoveryLimitErrorStrategy) Recover(recognizer antlr.Parser, e antlr.RecognitionException) {
	rl.checkAttempts(recognizer)
	lc := &lookaheadConsumer{Parser: recognizer, errorRecoveryTokenLookaheadLimit: rl.errorRecoveryTokenLookaheadLimit}
	rl.DefaultErrorStrategy.Recover(lc, e)
}

func (rl *recoveryLimitErrorStrategy) RecoverInline(recognizer antlr.Parser) antlr.Token {
	rl.checkAttempts(recognizer)
	lc := &lookaheadConsumer{Parser: recognizer, errorRecoveryTokenLookaheadLimit: rl.errorRecoveryTokenLookaheadLimit}
	return rl.DefaultErrorStrategy.RecoverInline(lc)
}

func (rl *recoveryLimitErrorStrategy) checkAttempts(recognizer antlr.Parser) {
	if rl.recoveryAttempts == rl.errorRecoveryLimit {
		rl.recoveryAttempts++
		msg := fmt.Sprintf("error recovery attempt limit exceeded: %d", rl.errorRecoveryLimit)
		recognizer.NotifyErrorListeners(msg, nil, nil)
		panic(&recoveryLimitError{
			message: msg,
		})
	}
	rl.recoveryAttempts++
}

var _ antlr.ErrorStrategy = &recoveryLimitErrorStrategy{}

type parser struct {
	gen.BaseCELVisitor
	errors                           *parseErrors
	helper                           *parserHelper
	macros                           map[string]Macro
	recursionDepth                   int
	maxRecursionDepth                int
	errorRecoveryLimit               int
	errorRecoveryLookaheadTokenLimit int
	populateMacroCalls               bool
}

var (
	_ gen.CELVisitor = (*parser)(nil)

	lexerPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			l := gen.NewCELLexer(nil)
			l.RemoveErrorListeners()
			return l
		},
	}

	parserPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			p := gen.NewCELParser(nil)
			p.RemoveErrorListeners()
			return p
		},
	}
)

func (p *parser) parse(expr runes.Buffer, desc string) *exprpb.Expr {
	// TODO: get rid of these pools once https://github.com/antlr/antlr4/pull/3571 is in a release
	lexer := lexerPool.Get().(*gen.CELLexer)
	prsr := parserPool.Get().(*gen.CELParser)

	// Unfortunately ANTLR Go runtime is missing (*antlr.BaseParser).RemoveParseListeners, so this is
	// good enough until that is exported.
	prsrListener := &recursionListener{
		maxDepth:      p.maxRecursionDepth,
		ruleTypeDepth: map[int]*int{},
	}

	defer func() {
		// Reset the lexer and parser before putting them back in the pool.
		lexer.RemoveErrorListeners()
		prsr.RemoveParseListener(prsrListener)
		prsr.RemoveErrorListeners()
		lexer.SetInputStream(nil)
		prsr.SetInputStream(nil)
		lexerPool.Put(lexer)
		parserPool.Put(prsr)
	}()

	lexer.SetInputStream(newCharStream(expr, desc))
	prsr.SetInputStream(antlr.NewCommonTokenStream(lexer, 0))

	lexer.AddErrorListener(p)
	prsr.AddErrorListener(p)
	prsr.AddParseListener(prsrListener)

	prsr.SetErrorHandler(&recoveryLimitErrorStrategy{
		DefaultErrorStrategy:             antlr.NewDefaultErrorStrategy(),
		errorRecoveryLimit:               p.errorRecoveryLimit,
		errorRecoveryTokenLookaheadLimit: p.errorRecoveryLookaheadTokenLimit,
	})

	defer func() {
		if val := recover(); val != nil {
			switch err := val.(type) {
			case *lookaheadLimitError:
				p.errors.ReportError(common.NoLocation, err.Error())
			case *recursionError:
				p.errors.ReportError(common.NoLocation, err.Error())
			case *recoveryLimitError:
				// do nothing, listeners already notified and error reported.
			default:
				panic(val)
			}
		}
	}()

	return p.Visit(prsr.Start()).(*exprpb.Expr)
}

// Visitor implementations.
func (p *parser) Visit(tree antlr.ParseTree) interface{} {
	t := unnest(tree)
	switch tree := t.(type) {
	case *gen.StartContext:
		return p.VisitStart(tree)
	case *gen.ExprContext:
		p.checkAndIncrementRecursionDepth()
		out := p.VisitExpr(tree)
		p.decrementRecursionDepth()
		return out
	case *gen.ConditionalAndContext:
		return p.VisitConditionalAnd(tree)
	case *gen.ConditionalOrContext:
		return p.VisitConditionalOr(tree)
	case *gen.RelationContext:
		p.checkAndIncrementRecursionDepth()
		out := p.VisitRelation(tree)
		p.decrementRecursionDepth()
		return out
	case *gen.CalcContext:
		p.checkAndIncrementRecursionDepth()
		out := p.VisitCalc(tree)
		p.decrementRecursionDepth()
		return out
	case *gen.LogicalNotContext:
		return p.VisitLogicalNot(tree)
	case *gen.IdentOrGlobalCallContext:
		return p.VisitIdentOrGlobalCall(tree)
	case *gen.SelectContext:
		p.checkAndIncrementRecursionDepth()
		out := p.VisitSelect(tree)
		p.decrementRecursionDepth()
		return out
	case *gen.MemberCallContext:
		p.checkAndIncrementRecursionDepth()
		out := p.VisitMemberCall(tree)
		p.decrementRecursionDepth()
		return out
	case *gen.MapInitializerListContext:
		return p.VisitMapInitializerList(tree)
	case *gen.NegateContext:
		return p.VisitNegate(tree)
	case *gen.IndexContext:
		p.checkAndIncrementRecursionDepth()
		out := p.VisitIndex(tree)
		p.decrementRecursionDepth()
		return out
	case *gen.UnaryContext:
		return p.VisitUnary(tree)
	case *gen.CreateListContext:
		return p.VisitCreateList(tree)
	case *gen.CreateMessageContext:
		return p.VisitCreateMessage(tree)
	case *gen.CreateStructContext:
		return p.VisitCreateStruct(tree)
	case *gen.IntContext:
		return p.VisitInt(tree)
	case *gen.UintContext:
		return p.VisitUint(tree)
	case *gen.DoubleContext:
		return p.VisitDouble(tree)
	case *gen.StringContext:
		return p.VisitString(tree)
	case *gen.BytesContext:
		return p.VisitBytes(tree)
	case *gen.BoolFalseContext:
		return p.VisitBoolFalse(tree)
	case *gen.BoolTrueContext:
		return p.VisitBoolTrue(tree)
	case *gen.NullContext:
		return p.VisitNull(tree)
	}

	// Report at least one error if the parser reaches an unknown parse element.
	// Typically, this happens if the parser has already encountered a syntax error elsewhere.
	if len(p.errors.GetErrors()) == 0 {
		txt := "<<nil>>"
		if t != nil {
			txt = fmt.Sprintf("<<%T>>", t)
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

// VisitSelect visits a parse tree produced by CELParser#Select.
func (p *parser) VisitSelect(ctx *gen.SelectContext) interface{} {
	operand := p.Visit(ctx.Member()).(*exprpb.Expr)
	// Handle the error case where no valid identifier is specified.
	if ctx.GetId() == nil {
		return p.helper.newExpr(ctx)
	}
	id := ctx.GetId().GetText()
	return p.helper.newSelect(ctx.GetOp(), operand, id)
}

// VisitMemberCall visits a parse tree produced by CELParser#MemberCall.
func (p *parser) VisitMemberCall(ctx *gen.MemberCallContext) interface{} {
	operand := p.Visit(ctx.Member()).(*exprpb.Expr)
	// Handle the error case where no valid identifier is specified.
	if ctx.GetId() == nil {
		return p.helper.newExpr(ctx)
	}
	id := ctx.GetId().GetText()
	opID := p.helper.id(ctx.GetOpen())
	return p.receiverCallOrMacro(opID, id, operand, p.visitList(ctx.GetArgs())...)
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
	messageName := ""
	for _, id := range ctx.GetIds() {
		if len(messageName) != 0 {
			messageName += "."
		}
		messageName += id.GetText()
	}
	if ctx.GetLeadingDot() != nil {
		messageName = "." + messageName
	}
	objID := p.helper.id(ctx.GetOp())
	entries := p.VisitIFieldInitializerList(ctx.GetEntries()).([]*exprpb.Expr_CreateStruct_Entry)
	return p.helper.newObject(objID, messageName, entries...)
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
	if _, ok := reservedIds[id]; ok {
		return p.reportError(ctx, "reserved identifier: %s", id)
	}
	identName += id
	if ctx.GetOp() != nil {
		opID := p.helper.id(ctx.GetOp())
		return p.globalCallOrMacro(opID, identName, p.visitList(ctx.GetArgs())...)
	}
	return p.helper.newIdent(ctx.GetId(), identName)
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
	base := 10
	if strings.HasPrefix(text, "0x") {
		base = 16
		text = text[2:]
	}
	if ctx.GetSign() != nil {
		text = ctx.GetSign().GetText() + text
	}
	i, err := strconv.ParseInt(text, base, 64)
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
	base := 10
	if strings.HasPrefix(text, "0x") {
		base = 16
		text = text[2:]
	}
	i, err := strconv.ParseUint(text, base, 64)
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
		location = p.helper.getLocation(err.GetId())
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
	// Hack to keep existing error messages consistent with previous versions of CEL when a reserved word
	// is used as an identifier. This behavior needs to be overhauled to provide consistent, normalized error
	// messages out of ANTLR to prevent future breaking changes related to error message content.
	if strings.Contains(msg, "no viable alternative") {
		msg = reservedIdentifier.ReplaceAllString(msg, mismatchedReservedIdentifier)
	}
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
	if p.populateMacroCalls {
		p.helper.addMacroCall(expr.GetId(), function, target, args...)
	}
	return expr, true
}

func (p *parser) checkAndIncrementRecursionDepth() {
	p.recursionDepth++
	if p.recursionDepth > p.maxRecursionDepth {
		panic(&recursionError{message: "max recursion depth exceeded"})
	}
}

func (p *parser) decrementRecursionDepth() {
	p.recursionDepth--
}

// unnest traverses down the left-hand side of the parse graph until it encounters the first compound
// parse node or the first leaf in the parse graph.
func unnest(tree antlr.ParseTree) antlr.ParseTree {
	for tree != nil {
		switch t := tree.(type) {
		case *gen.ExprContext:
			// conditionalOr op='?' conditionalOr : expr
			if t.GetOp() != nil {
				return t
			}
			// conditionalOr
			tree = t.GetE()
		case *gen.ConditionalOrContext:
			// conditionalAnd (ops=|| conditionalAnd)*
			if t.GetOps() != nil && len(t.GetOps()) > 0 {
				return t
			}
			// conditionalAnd
			tree = t.GetE()
		case *gen.ConditionalAndContext:
			// relation (ops=&& relation)*
			if t.GetOps() != nil && len(t.GetOps()) > 0 {
				return t
			}
			// relation
			tree = t.GetE()
		case *gen.RelationContext:
			// relation op relation
			if t.GetOp() != nil {
				return t
			}
			// calc
			tree = t.Calc()
		case *gen.CalcContext:
			// calc op calc
			if t.GetOp() != nil {
				return t
			}
			// unary
			tree = t.Unary()
		case *gen.MemberExprContext:
			// member expands to one of: primary, select, index, or create message
			tree = t.Member()
		case *gen.PrimaryExprContext:
			// primary expands to one of identifier, nested, create list, create struct, literal
			tree = t.Primary()
		case *gen.NestedContext:
			// contains a nested 'expr'
			tree = t.GetE()
		case *gen.ConstantLiteralContext:
			// expands to a primitive literal
			tree = t.Literal()
		default:
			return t
		}
	}
	return tree
}

var (
	reservedIdentifier           = regexp.MustCompile("no viable alternative at input '.(true|false|null)'")
	mismatchedReservedIdentifier = "mismatched input '$1' expecting IDENTIFIER"
)
