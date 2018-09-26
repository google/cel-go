package parser

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/google/cel-go/common"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type parserHelper struct {
	source    common.Source
	errors    *parseErrors
	macros    map[string]Macro
	nextId    int64
	positions map[int64]int32
}

func newParserHelper(source common.Source, macros Macros) *parserHelper {
	macroMap := make(map[string]Macro)
	for _, m := range macros {
		macroMap[makeMacroKey(m.name, m.args, m.instanceStyle)] = m
	}

	return &parserHelper{
		errors:    &parseErrors{common.NewErrors(source)},
		source:    source,
		macros:    macroMap,
		nextId:    1,
		positions: make(map[int64]int32),
	}
}

func (p *parserHelper) getSourceInfo() *expr.SourceInfo {
	return &expr.SourceInfo{
		Location:    p.source.Description(),
		Positions:   p.positions,
		LineOffsets: p.source.LineOffsets()}
}

func (p *parserHelper) reportError(ctx interface{}, format string, args ...interface{}) *expr.Expr {
	var location common.Location
	switch ctx.(type) {
	case common.Location:
		location = ctx.(common.Location)
	case antlr.Token, antlr.ParserRuleContext:
		err := p.newExpr(ctx)
		location = p.getLocation(err.Id)
	}
	err := p.newExpr(ctx)
	// Provide arguments to the report error.
	p.errors.ReportError(location, format, args...)
	return err
}

func (p *parserHelper) newLiteral(ctx interface{}, value *expr.Constant) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_ConstExpr{ConstExpr: value}
	return exprNode
}

func (p *parserHelper) newLiteralBool(ctx interface{}, value bool) *expr.Expr {
	return p.newLiteral(ctx,
		&expr.Constant{ConstantKind: &expr.Constant_BoolValue{value}})
}

func (p *parserHelper) newLiteralString(ctx interface{}, value string) *expr.Expr {
	return p.newLiteral(ctx,
		&expr.Constant{ConstantKind: &expr.Constant_StringValue{value}})
}

func (p *parserHelper) newLiteralBytes(ctx interface{}, value []byte) *expr.Expr {
	return p.newLiteral(ctx,
		&expr.Constant{ConstantKind: &expr.Constant_BytesValue{value}})
}

func (p *parserHelper) newLiteralInt(ctx interface{}, value int64) *expr.Expr {
	return p.newLiteral(ctx,
		&expr.Constant{ConstantKind: &expr.Constant_Int64Value{value}})
}

func (p *parserHelper) newLiteralUint(ctx interface{}, value uint64) *expr.Expr {
	return p.newLiteral(ctx, &expr.Constant{ConstantKind: &expr.Constant_Uint64Value{value}})
}

func (p *parserHelper) newLiteralDouble(ctx interface{}, value float64) *expr.Expr {
	return p.newLiteral(ctx,
		&expr.Constant{ConstantKind: &expr.Constant_DoubleValue{value}})
}

func (p *parserHelper) newIdent(ctx interface{}, name string) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_IdentExpr{IdentExpr: &expr.Expr_Ident{Name: name}}
	return exprNode
}

func (p *parserHelper) newSelect(ctx interface{}, operand *expr.Expr, field string) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_SelectExpr{
		SelectExpr: &expr.Expr_Select{Operand: operand, Field: field}}
	return exprNode
}

func (p *parserHelper) newPresenceTest(ctx interface{}, operand *expr.Expr, field string) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_SelectExpr{
		SelectExpr: &expr.Expr_Select{Operand: operand, Field: field, TestOnly: true}}
	return exprNode
}

func (p *parserHelper) newGlobalCall(ctx interface{}, function string, args ...*expr.Expr) *expr.Expr {
	if macro, found := p.macros[makeMacroKey(function, len(args), false)]; found {
		return macro.expander(p, ctx, nil, args)
	}
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_CallExpr{
		CallExpr: &expr.Expr_Call{Function: function, Args: args}}
	return exprNode
}

func (p *parserHelper) newMemberCall(ctx interface{}, function string, target *expr.Expr, args ...*expr.Expr) *expr.Expr {
	if macro, found := p.macros[makeMacroKey(function, len(args), true)]; found {
		return macro.expander(p, ctx, target, args)
	}
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_CallExpr{
		CallExpr: &expr.Expr_Call{Function: function, Target: target, Args: args}}
	return exprNode
}

func (p *parserHelper) newList(ctx interface{}, elements ...*expr.Expr) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_ListExpr{
		ListExpr: &expr.Expr_CreateList{Elements: elements}}
	return exprNode
}

func (p *parserHelper) newMap(ctx interface{}, entries ...*expr.Expr_CreateStruct_Entry) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_StructExpr{
		StructExpr: &expr.Expr_CreateStruct{Entries: entries}}
	return exprNode
}

func (p *parserHelper) newMapEntry(ctx interface{}, key *expr.Expr, value *expr.Expr) *expr.Expr_CreateStruct_Entry {
	return &expr.Expr_CreateStruct_Entry{
		Id:      p.id(ctx),
		KeyKind: &expr.Expr_CreateStruct_Entry_MapKey{MapKey: key},
		Value:   value}
}

func (p *parserHelper) newObject(ctx interface{},
	typeName string,
	entries ...*expr.Expr_CreateStruct_Entry) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_StructExpr{
		StructExpr: &expr.Expr_CreateStruct{
			MessageName: typeName,
			Entries:     entries}}
	return exprNode
}

func (p *parserHelper) newObjectField(ctx interface{}, field string, value *expr.Expr) *expr.Expr_CreateStruct_Entry {
	return &expr.Expr_CreateStruct_Entry{
		Id:      p.id(ctx),
		KeyKind: &expr.Expr_CreateStruct_Entry_FieldKey{FieldKey: field},
		Value:   value}
}

func (p *parserHelper) newComprehension(ctx interface{}, iterVar string,
	iterRange *expr.Expr,
	accuVar string,
	accuInit *expr.Expr,
	condition *expr.Expr,
	step *expr.Expr,
	result *expr.Expr) *expr.Expr {
	exprNode := p.newExpr(ctx)
	exprNode.ExprKind = &expr.Expr_ComprehensionExpr{
		ComprehensionExpr: &expr.Expr_Comprehension{
			AccuVar:       accuVar,
			AccuInit:      accuInit,
			IterVar:       iterVar,
			IterRange:     iterRange,
			LoopCondition: condition,
			LoopStep:      step,
			Result:        result}}
	return exprNode
}

func (p *parserHelper) newExpr(ctx interface{}) *expr.Expr {
	return &expr.Expr{Id: p.id(ctx)}
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
	location := common.NewLocation(token.GetLine(), token.GetColumn())
	id := p.nextId
	p.positions[id], _ = p.source.LocationOffset(location)
	p.nextId++
	return id
}

func (p *parserHelper) getLocation(id int64) common.Location {
	characterOffset := p.positions[id]
	location, _ := p.source.OffsetLocation(characterOffset)
	return location
}

// ANTLR Parse listener implementations
func (p *parserHelper) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	// TODO: Snippet
	l := common.NewLocation(line, column)
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
