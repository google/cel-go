package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	commonpb "github.com/google/cel-go/common"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
)

type parserHelper struct {
	source    commonpb.Source
	errors    *parseErrors
	macros    map[string]Macro
	nextId    int64
	positions map[int64]int32
}

func newParserHelper(source commonpb.Source, macros Macros) *parserHelper {
	macroMap := make(map[string]Macro)
	for _, m := range macros {
		macroMap[makeMacroKey(m.name, m.args, m.instanceStyle)] = m
	}

	return &parserHelper{
		errors:    &parseErrors{commonpb.NewErrors(source)},
		source:    source,
		macros:    macroMap,
		nextId:    1,
		positions: make(map[int64]int32),
	}
}

func (p *parserHelper) getSourceInfo() *exprpb.SourceInfo {
	return &exprpb.SourceInfo{
		Location:    p.source.Description(),
		Positions:   p.positions,
		LineOffsets: p.source.LineOffsets()}
}

func (p *parserHelper) reportError(ctx interface{}, format string, args ...interface{}) *exprpb.Expr {
	var location commonpb.Location
	switch ctx.(type) {
	case commonpb.Location:
		location = ctx.(commonpb.Location)
	case antlr.Token, antlr.ParserRuleContext:
		err := p.newExpr(ctx)
		location = p.getLocation(err.Id)
	}
	err := p.newExpr(ctx)
	// Provide arguments to the report error.
	p.errors.ReportError(location, format, args...)
	return err
}

func (p *parserHelper) newLiteral(ctx interface{}, value *exprpb.Literal) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_LiteralExpr{LiteralExpr: value}
	return exprNode
}

func (p *parserHelper) newLiteralBool(ctx interface{}, value bool) *exprpb.Expr {
	return p.newLiteral(ctx,
		&exprpb.Literal{LiteralKind: &exprpb.Literal_BoolValue{value}})
}

func (p *parserHelper) newLiteralString(ctx interface{}, value string) *exprpb.Expr {
	return p.newLiteral(ctx,
		&exprpb.Literal{LiteralKind: &exprpb.Literal_StringValue{value}})
}

func (p *parserHelper) newLiteralBytes(ctx interface{}, value []byte) *exprpb.Expr {
	return p.newLiteral(ctx,
		&exprpb.Literal{LiteralKind: &exprpb.Literal_BytesValue{value}})
}

func (p *parserHelper) newLiteralInt(ctx interface{}, value int64) *exprpb.Expr {
	return p.newLiteral(ctx,
		&exprpb.Literal{LiteralKind: &exprpb.Literal_Int64Value{value}})
}

func (p *parserHelper) newLiteralUint(ctx interface{}, value uint64) *exprpb.Expr {
	return p.newLiteral(ctx, &exprpb.Literal{LiteralKind: &exprpb.Literal_Uint64Value{value}})
}

func (p *parserHelper) newLiteralDouble(ctx interface{}, value float64) *exprpb.Expr {
	return p.newLiteral(ctx,
		&exprpb.Literal{LiteralKind: &exprpb.Literal_DoubleValue{value}})
}

func (p *parserHelper) newIdent(ctx interface{}, name string) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_IdentExpr{IdentExpr: &exprpb.Expr_Ident{Name: name}}
	return exprNode
}

func (p *parserHelper) newSelect(ctx interface{}, operand *exprpb.Expr, field string) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_SelectExpr{
		SelectExpr: &exprpb.Expr_Select{Operand: operand, Field: field}}
	return exprNode
}

func (p *parserHelper) newPresenceTest(ctx interface{}, operand *exprpb.Expr, field string) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_SelectExpr{
		SelectExpr: &exprpb.Expr_Select{Operand: operand, Field: field, TestOnly: true}}
	return exprNode
}

func (p *parserHelper) newGlobalCall(ctx interface{}, function string, args ...*exprpb.Expr) *exprpb.Expr {
	if macro, found := p.macros[makeMacroKey(function, len(args), false)]; found {
		return macro.expander(p, ctx, nil, args)
	}
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_CallExpr{
		CallExpr: &exprpb.Expr_Call{Function: function, Args: args}}
	return exprNode
}

func (p *parserHelper) newMemberCall(ctx interface{}, function string, target *exprpb.Expr, args ...*exprpb.Expr) *exprpb.Expr {
	if macro, found := p.macros[makeMacroKey(function, len(args), true)]; found {
		return macro.expander(p, ctx, target, args)
	}
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_CallExpr{
		CallExpr: &exprpb.Expr_Call{Function: function, Target: target, Args: args}}
	return exprNode
}

func (p *parserHelper) newList(ctx interface{}, elements ...*exprpb.Expr) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_ListExpr{
		ListExpr: &exprpb.Expr_CreateList{Elements: elements}}
	return exprNode
}

func (p *parserHelper) newMap(ctx interface{}, entries ...*exprpb.Expr_CreateStruct_Entry) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_StructExpr{
		StructExpr: &exprpb.Expr_CreateStruct{Entries: entries}}
	return exprNode
}

func (p *parserHelper) newMapEntry(ctx interface{}, key *exprpb.Expr, value *exprpb.Expr) *exprpb.Expr_CreateStruct_Entry {
	return &exprpb.Expr_CreateStruct_Entry{
		Id:      p.id(ctx),
		KeyKind: &exprpb.Expr_CreateStruct_Entry_MapKey{MapKey: key},
		Value:   value}
}

func (p *parserHelper) newObject(ctx interface{},
	typeName string,
	entries ...*exprpb.Expr_CreateStruct_Entry) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_StructExpr{
		StructExpr: &exprpb.Expr_CreateStruct{
			MessageName: typeName,
			Entries:     entries}}
	return exprNode
}

func (p *parserHelper) newObjectField(ctx interface{}, field string, value *exprpb.Expr) *exprpb.Expr_CreateStruct_Entry {
	return &exprpb.Expr_CreateStruct_Entry{
		Id:      p.id(ctx),
		KeyKind: &exprpb.Expr_CreateStruct_Entry_FieldKey{FieldKey: field},
		Value:   value}
}

func (p *parserHelper) newComprehension(ctx interface{}, iterVar string,
	iterRange *exprpb.Expr,
	accuVar string,
	accuInit *exprpb.Expr,
	condition *exprpb.Expr,
	step *exprpb.Expr,
	result *exprpb.Expr) *exprpb.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &exprpb.Expr_ComprehensionExpr{
		ComprehensionExpr: &exprpb.Expr_Comprehension{
			AccuVar:       accuVar,
			AccuInit:      accuInit,
			IterVar:       iterVar,
			IterRange:     iterRange,
			LoopCondition: condition,
			LoopStep:      step,
			Result:        result}}
	return exprNode
}

func (p *parserHelper) newExpr(ctx interface{}) *exprpb.Expr {
	return &exprpb.Expr{Id: p.id(ctx)}
}

func (p *parserHelper) id(ctx interface{}) int64 {
	var token antlr.Token = nil
	switch ctx.(type) {
	case antlr.ParserRuleContext:
		token = (ctx.(antlr.ParserRuleContext)).GetStart()
	case antlr.Token:
		token = ctx.(antlr.Token)
	default:
		// This should only happen if the ctx is nil
		return -1
	}
	location := commonpb.NewLocation(token.GetLine(), token.GetColumn())
	id := p.nextId
	p.positions[id], _ = p.source.LocationOffset(location)
	p.nextId++
	return id
}

func (p *parserHelper) getLocation(id int64) commonpb.Location {
	characterOffset := p.positions[id]
	location, _ := p.source.OffsetLocation(characterOffset)
	return location
}

// ANTLR Parse listener implementations
func (p *parserHelper) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// TODO: Snippet
	l := commonpb.NewLocation(line, column)
	p.errors.syntaxError(l, msg)
}

func (p *parserHelper) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	// Intentional
}

func (p *parserHelper) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
	// Intentional
}

func (p *parserHelper) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
	// Intentional
}
