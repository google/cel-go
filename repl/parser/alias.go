// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package parser provides compatibility aliases forwarding to cel.dev/cel-go/repl/parser.
package parser

import (
	canonparser "cel.dev/cel-go/repl/parser"
	antlr "github.com/antlr4-go/antlr/v4"
)

// Type aliases
type (
	BaseCommandsListener         = canonparser.BaseCommandsListener
	BaseCommandsVisitor          = canonparser.BaseCommandsVisitor
	BoolFalseContext             = canonparser.BoolFalseContext
	BoolTrueContext              = canonparser.BoolTrueContext
	BytesContext                 = canonparser.BytesContext
	CalcContext                  = canonparser.CalcContext
	CommandContext               = canonparser.CommandContext
	CommandsLexer                = canonparser.CommandsLexer
	CommandsListener             = canonparser.CommandsListener
	CommandsParser               = canonparser.CommandsParser
	CommandsVisitor              = canonparser.CommandsVisitor
	CompileContext               = canonparser.CompileContext
	ConditionalAndContext        = canonparser.ConditionalAndContext
	ConditionalOrContext         = canonparser.ConditionalOrContext
	ConstantLiteralContext       = canonparser.ConstantLiteralContext
	CreateListContext            = canonparser.CreateListContext
	CreateMessageContext         = canonparser.CreateMessageContext
	CreateStructContext          = canonparser.CreateStructContext
	DeclareContext               = canonparser.DeclareContext
	DeleteContext                = canonparser.DeleteContext
	DoubleContext                = canonparser.DoubleContext
	EmptyContext                 = canonparser.EmptyContext
	EscapeIdentContext           = canonparser.EscapeIdentContext
	EscapedIdentifierContext     = canonparser.EscapedIdentifierContext
	ExprCmdContext               = canonparser.ExprCmdContext
	ExprContext                  = canonparser.ExprContext
	ExprListContext              = canonparser.ExprListContext
	FieldInitializerListContext  = canonparser.FieldInitializerListContext
	FnDeclContext                = canonparser.FnDeclContext
	GlobalCallContext            = canonparser.GlobalCallContext
	HelpContext                  = canonparser.HelpContext
	ICalcContext                 = canonparser.ICalcContext
	ICommandContext              = canonparser.ICommandContext
	ICompileContext              = canonparser.ICompileContext
	IConditionalAndContext       = canonparser.IConditionalAndContext
	IConditionalOrContext        = canonparser.IConditionalOrContext
	IDeclareContext              = canonparser.IDeclareContext
	IDeleteContext               = canonparser.IDeleteContext
	IEmptyContext                = canonparser.IEmptyContext
	IEscapeIdentContext          = canonparser.IEscapeIdentContext
	IExprCmdContext              = canonparser.IExprCmdContext
	IExprContext                 = canonparser.IExprContext
	IExprListContext             = canonparser.IExprListContext
	IFieldInitializerListContext = canonparser.IFieldInitializerListContext
	IFnDeclContext               = canonparser.IFnDeclContext
	IHelpContext                 = canonparser.IHelpContext
	ILetContext                  = canonparser.ILetContext
	IListInitContext             = canonparser.IListInitContext
	ILiteralContext              = canonparser.ILiteralContext
	IMapInitializerListContext   = canonparser.IMapInitializerListContext
	IMemberContext               = canonparser.IMemberContext
	IOptExprContext              = canonparser.IOptExprContext
	IOptFieldContext             = canonparser.IOptFieldContext
	IParamContext                = canonparser.IParamContext
	IParamIdContext              = canonparser.IParamIdContext
	IParseContext                = canonparser.IParseContext
	IPrimaryContext              = canonparser.IPrimaryContext
	IQualIdContext               = canonparser.IQualIdContext
	IRelationContext             = canonparser.IRelationContext
	ISimpleContext               = canonparser.ISimpleContext
	IStartCommandContext         = canonparser.IStartCommandContext
	IStartContext                = canonparser.IStartContext
	IStartTypeContext            = canonparser.IStartTypeContext
	ITypeContext                 = canonparser.ITypeContext
	ITypeIdContext               = canonparser.ITypeIdContext
	ITypeParamListContext        = canonparser.ITypeParamListContext
	IUnaryContext                = canonparser.IUnaryContext
	IVarDeclContext              = canonparser.IVarDeclContext
	IdentContext                 = canonparser.IdentContext
	IndexContext                 = canonparser.IndexContext
	IntContext                   = canonparser.IntContext
	LetContext                   = canonparser.LetContext
	ListInitContext              = canonparser.ListInitContext
	LiteralContext               = canonparser.LiteralContext
	LogicalNotContext            = canonparser.LogicalNotContext
	MapInitializerListContext    = canonparser.MapInitializerListContext
	MemberCallContext            = canonparser.MemberCallContext
	MemberContext                = canonparser.MemberContext
	MemberExprContext            = canonparser.MemberExprContext
	NegateContext                = canonparser.NegateContext
	NestedContext                = canonparser.NestedContext
	NullContext                  = canonparser.NullContext
	OptExprContext               = canonparser.OptExprContext
	OptFieldContext              = canonparser.OptFieldContext
	ParamContext                 = canonparser.ParamContext
	ParamIdContext               = canonparser.ParamIdContext
	ParseContext                 = canonparser.ParseContext
	PrimaryContext               = canonparser.PrimaryContext
	PrimaryExprContext           = canonparser.PrimaryExprContext
	QualIdContext                = canonparser.QualIdContext
	RelationContext              = canonparser.RelationContext
	SelectContext                = canonparser.SelectContext
	SimpleContext                = canonparser.SimpleContext
	SimpleIdentifierContext      = canonparser.SimpleIdentifierContext
	StartCommandContext          = canonparser.StartCommandContext
	StartContext                 = canonparser.StartContext
	StartTypeContext             = canonparser.StartTypeContext
	StringContext                = canonparser.StringContext
	TypeContext                  = canonparser.TypeContext
	TypeIdContext                = canonparser.TypeIdContext
	TypeParamListContext         = canonparser.TypeParamListContext
	UintContext                  = canonparser.UintContext
	UnaryContext                 = canonparser.UnaryContext
	VarDeclContext               = canonparser.VarDeclContext
)

// Constants
const (
	CommandsLexerARROW                      = canonparser.CommandsLexerARROW
	CommandsLexerBYTES                      = canonparser.CommandsLexerBYTES
	CommandsLexerCEL_FALSE                  = canonparser.CommandsLexerCEL_FALSE
	CommandsLexerCEL_TRUE                   = canonparser.CommandsLexerCEL_TRUE
	CommandsLexerCOLON                      = canonparser.CommandsLexerCOLON
	CommandsLexerCOMMA                      = canonparser.CommandsLexerCOMMA
	CommandsLexerCOMMAND                    = canonparser.CommandsLexerCOMMAND
	CommandsLexerCOMMENT                    = canonparser.CommandsLexerCOMMENT
	CommandsLexerDOT                        = canonparser.CommandsLexerDOT
	CommandsLexerEQUALS                     = canonparser.CommandsLexerEQUALS
	CommandsLexerEQUAL_ASSIGN               = canonparser.CommandsLexerEQUAL_ASSIGN
	CommandsLexerESC_IDENTIFIER             = canonparser.CommandsLexerESC_IDENTIFIER
	CommandsLexerEXCLAM                     = canonparser.CommandsLexerEXCLAM
	CommandsLexerFLAG                       = canonparser.CommandsLexerFLAG
	CommandsLexerGREATER                    = canonparser.CommandsLexerGREATER
	CommandsLexerGREATER_EQUALS             = canonparser.CommandsLexerGREATER_EQUALS
	CommandsLexerIDENTIFIER                 = canonparser.CommandsLexerIDENTIFIER
	CommandsLexerIN                         = canonparser.CommandsLexerIN
	CommandsLexerLBRACE                     = canonparser.CommandsLexerLBRACE
	CommandsLexerLBRACKET                   = canonparser.CommandsLexerLBRACKET
	CommandsLexerLESS                       = canonparser.CommandsLexerLESS
	CommandsLexerLESS_EQUALS                = canonparser.CommandsLexerLESS_EQUALS
	CommandsLexerLOGICAL_AND                = canonparser.CommandsLexerLOGICAL_AND
	CommandsLexerLOGICAL_OR                 = canonparser.CommandsLexerLOGICAL_OR
	CommandsLexerLPAREN                     = canonparser.CommandsLexerLPAREN
	CommandsLexerMINUS                      = canonparser.CommandsLexerMINUS
	CommandsLexerNOT_EQUALS                 = canonparser.CommandsLexerNOT_EQUALS
	CommandsLexerNUL                        = canonparser.CommandsLexerNUL
	CommandsLexerNUM_FLOAT                  = canonparser.CommandsLexerNUM_FLOAT
	CommandsLexerNUM_INT                    = canonparser.CommandsLexerNUM_INT
	CommandsLexerNUM_UINT                   = canonparser.CommandsLexerNUM_UINT
	CommandsLexerPARAM_SPECIFIER            = canonparser.CommandsLexerPARAM_SPECIFIER
	CommandsLexerPERCENT                    = canonparser.CommandsLexerPERCENT
	CommandsLexerPLUS                       = canonparser.CommandsLexerPLUS
	CommandsLexerQUESTIONMARK               = canonparser.CommandsLexerQUESTIONMARK
	CommandsLexerRBRACE                     = canonparser.CommandsLexerRBRACE
	CommandsLexerRPAREN                     = canonparser.CommandsLexerRPAREN
	CommandsLexerRPRACKET                   = canonparser.CommandsLexerRPRACKET
	CommandsLexerSLASH                      = canonparser.CommandsLexerSLASH
	CommandsLexerSTAR                       = canonparser.CommandsLexerSTAR
	CommandsLexerSTRING                     = canonparser.CommandsLexerSTRING
	CommandsLexerT__0                       = canonparser.CommandsLexerT__0
	CommandsLexerT__1                       = canonparser.CommandsLexerT__1
	CommandsLexerT__2                       = canonparser.CommandsLexerT__2
	CommandsLexerT__3                       = canonparser.CommandsLexerT__3
	CommandsLexerT__4                       = canonparser.CommandsLexerT__4
	CommandsLexerT__5                       = canonparser.CommandsLexerT__5
	CommandsLexerT__6                       = canonparser.CommandsLexerT__6
	CommandsLexerT__7                       = canonparser.CommandsLexerT__7
	CommandsLexerT__8                       = canonparser.CommandsLexerT__8
	CommandsLexerWHITESPACE                 = canonparser.CommandsLexerWHITESPACE
	CommandsParserARROW                     = canonparser.CommandsParserARROW
	CommandsParserBYTES                     = canonparser.CommandsParserBYTES
	CommandsParserCEL_FALSE                 = canonparser.CommandsParserCEL_FALSE
	CommandsParserCEL_TRUE                  = canonparser.CommandsParserCEL_TRUE
	CommandsParserCOLON                     = canonparser.CommandsParserCOLON
	CommandsParserCOMMA                     = canonparser.CommandsParserCOMMA
	CommandsParserCOMMAND                   = canonparser.CommandsParserCOMMAND
	CommandsParserCOMMENT                   = canonparser.CommandsParserCOMMENT
	CommandsParserDOT                       = canonparser.CommandsParserDOT
	CommandsParserEOF                       = canonparser.CommandsParserEOF
	CommandsParserEQUALS                    = canonparser.CommandsParserEQUALS
	CommandsParserEQUAL_ASSIGN              = canonparser.CommandsParserEQUAL_ASSIGN
	CommandsParserESC_IDENTIFIER            = canonparser.CommandsParserESC_IDENTIFIER
	CommandsParserEXCLAM                    = canonparser.CommandsParserEXCLAM
	CommandsParserFLAG                      = canonparser.CommandsParserFLAG
	CommandsParserGREATER                   = canonparser.CommandsParserGREATER
	CommandsParserGREATER_EQUALS            = canonparser.CommandsParserGREATER_EQUALS
	CommandsParserIDENTIFIER                = canonparser.CommandsParserIDENTIFIER
	CommandsParserIN                        = canonparser.CommandsParserIN
	CommandsParserLBRACE                    = canonparser.CommandsParserLBRACE
	CommandsParserLBRACKET                  = canonparser.CommandsParserLBRACKET
	CommandsParserLESS                      = canonparser.CommandsParserLESS
	CommandsParserLESS_EQUALS               = canonparser.CommandsParserLESS_EQUALS
	CommandsParserLOGICAL_AND               = canonparser.CommandsParserLOGICAL_AND
	CommandsParserLOGICAL_OR                = canonparser.CommandsParserLOGICAL_OR
	CommandsParserLPAREN                    = canonparser.CommandsParserLPAREN
	CommandsParserMINUS                     = canonparser.CommandsParserMINUS
	CommandsParserNOT_EQUALS                = canonparser.CommandsParserNOT_EQUALS
	CommandsParserNUL                       = canonparser.CommandsParserNUL
	CommandsParserNUM_FLOAT                 = canonparser.CommandsParserNUM_FLOAT
	CommandsParserNUM_INT                   = canonparser.CommandsParserNUM_INT
	CommandsParserNUM_UINT                  = canonparser.CommandsParserNUM_UINT
	CommandsParserPARAM_SPECIFIER           = canonparser.CommandsParserPARAM_SPECIFIER
	CommandsParserPERCENT                   = canonparser.CommandsParserPERCENT
	CommandsParserPLUS                      = canonparser.CommandsParserPLUS
	CommandsParserQUESTIONMARK              = canonparser.CommandsParserQUESTIONMARK
	CommandsParserRBRACE                    = canonparser.CommandsParserRBRACE
	CommandsParserRPAREN                    = canonparser.CommandsParserRPAREN
	CommandsParserRPRACKET                  = canonparser.CommandsParserRPRACKET
	CommandsParserRULE_calc                 = canonparser.CommandsParserRULE_calc
	CommandsParserRULE_command              = canonparser.CommandsParserRULE_command
	CommandsParserRULE_compile              = canonparser.CommandsParserRULE_compile
	CommandsParserRULE_conditionalAnd       = canonparser.CommandsParserRULE_conditionalAnd
	CommandsParserRULE_conditionalOr        = canonparser.CommandsParserRULE_conditionalOr
	CommandsParserRULE_declare              = canonparser.CommandsParserRULE_declare
	CommandsParserRULE_delete               = canonparser.CommandsParserRULE_delete
	CommandsParserRULE_empty                = canonparser.CommandsParserRULE_empty
	CommandsParserRULE_escapeIdent          = canonparser.CommandsParserRULE_escapeIdent
	CommandsParserRULE_expr                 = canonparser.CommandsParserRULE_expr
	CommandsParserRULE_exprCmd              = canonparser.CommandsParserRULE_exprCmd
	CommandsParserRULE_exprList             = canonparser.CommandsParserRULE_exprList
	CommandsParserRULE_fieldInitializerList = canonparser.CommandsParserRULE_fieldInitializerList
	CommandsParserRULE_fnDecl               = canonparser.CommandsParserRULE_fnDecl
	CommandsParserRULE_help                 = canonparser.CommandsParserRULE_help
	CommandsParserRULE_let                  = canonparser.CommandsParserRULE_let
	CommandsParserRULE_listInit             = canonparser.CommandsParserRULE_listInit
	CommandsParserRULE_literal              = canonparser.CommandsParserRULE_literal
	CommandsParserRULE_mapInitializerList   = canonparser.CommandsParserRULE_mapInitializerList
	CommandsParserRULE_member               = canonparser.CommandsParserRULE_member
	CommandsParserRULE_optExpr              = canonparser.CommandsParserRULE_optExpr
	CommandsParserRULE_optField             = canonparser.CommandsParserRULE_optField
	CommandsParserRULE_param                = canonparser.CommandsParserRULE_param
	CommandsParserRULE_paramId              = canonparser.CommandsParserRULE_paramId
	CommandsParserRULE_parse                = canonparser.CommandsParserRULE_parse
	CommandsParserRULE_primary              = canonparser.CommandsParserRULE_primary
	CommandsParserRULE_qualId               = canonparser.CommandsParserRULE_qualId
	CommandsParserRULE_relation             = canonparser.CommandsParserRULE_relation
	CommandsParserRULE_simple               = canonparser.CommandsParserRULE_simple
	CommandsParserRULE_start                = canonparser.CommandsParserRULE_start
	CommandsParserRULE_startCommand         = canonparser.CommandsParserRULE_startCommand
	CommandsParserRULE_startType            = canonparser.CommandsParserRULE_startType
	CommandsParserRULE_type                 = canonparser.CommandsParserRULE_type
	CommandsParserRULE_typeId               = canonparser.CommandsParserRULE_typeId
	CommandsParserRULE_typeParamList        = canonparser.CommandsParserRULE_typeParamList
	CommandsParserRULE_unary                = canonparser.CommandsParserRULE_unary
	CommandsParserRULE_varDecl              = canonparser.CommandsParserRULE_varDecl
	CommandsParserSLASH                     = canonparser.CommandsParserSLASH
	CommandsParserSTAR                      = canonparser.CommandsParserSTAR
	CommandsParserSTRING                    = canonparser.CommandsParserSTRING
	CommandsParserT__0                      = canonparser.CommandsParserT__0
	CommandsParserT__1                      = canonparser.CommandsParserT__1
	CommandsParserT__2                      = canonparser.CommandsParserT__2
	CommandsParserT__3                      = canonparser.CommandsParserT__3
	CommandsParserT__4                      = canonparser.CommandsParserT__4
	CommandsParserT__5                      = canonparser.CommandsParserT__5
	CommandsParserT__6                      = canonparser.CommandsParserT__6
	CommandsParserT__7                      = canonparser.CommandsParserT__7
	CommandsParserT__8                      = canonparser.CommandsParserT__8
	CommandsParserWHITESPACE                = canonparser.CommandsParserWHITESPACE
)

// Variables
var (
	CommandsLexerLexerStaticData = canonparser.CommandsLexerLexerStaticData
	CommandsParserStaticData     = canonparser.CommandsParserStaticData
)

// Functions

func CommandsLexerInit() {
	canonparser.CommandsLexerInit()
}

func CommandsParserInit() {
	canonparser.CommandsParserInit()
}

func InitEmptyCalcContext(p *CalcContext) {
	canonparser.InitEmptyCalcContext(p)
}

func InitEmptyCommandContext(p *CommandContext) {
	canonparser.InitEmptyCommandContext(p)
}

func InitEmptyCompileContext(p *CompileContext) {
	canonparser.InitEmptyCompileContext(p)
}

func InitEmptyConditionalAndContext(p *ConditionalAndContext) {
	canonparser.InitEmptyConditionalAndContext(p)
}

func InitEmptyConditionalOrContext(p *ConditionalOrContext) {
	canonparser.InitEmptyConditionalOrContext(p)
}

func InitEmptyDeclareContext(p *DeclareContext) {
	canonparser.InitEmptyDeclareContext(p)
}

func InitEmptyDeleteContext(p *DeleteContext) {
	canonparser.InitEmptyDeleteContext(p)
}

func InitEmptyEmptyContext(p *EmptyContext) {
	canonparser.InitEmptyEmptyContext(p)
}

func InitEmptyEscapeIdentContext(p *EscapeIdentContext) {
	canonparser.InitEmptyEscapeIdentContext(p)
}

func InitEmptyExprCmdContext(p *ExprCmdContext) {
	canonparser.InitEmptyExprCmdContext(p)
}

func InitEmptyExprContext(p *ExprContext) {
	canonparser.InitEmptyExprContext(p)
}

func InitEmptyExprListContext(p *ExprListContext) {
	canonparser.InitEmptyExprListContext(p)
}

func InitEmptyFieldInitializerListContext(p *FieldInitializerListContext) {
	canonparser.InitEmptyFieldInitializerListContext(p)
}

func InitEmptyFnDeclContext(p *FnDeclContext) {
	canonparser.InitEmptyFnDeclContext(p)
}

func InitEmptyHelpContext(p *HelpContext) {
	canonparser.InitEmptyHelpContext(p)
}

func InitEmptyLetContext(p *LetContext) {
	canonparser.InitEmptyLetContext(p)
}

func InitEmptyListInitContext(p *ListInitContext) {
	canonparser.InitEmptyListInitContext(p)
}

func InitEmptyLiteralContext(p *LiteralContext) {
	canonparser.InitEmptyLiteralContext(p)
}

func InitEmptyMapInitializerListContext(p *MapInitializerListContext) {
	canonparser.InitEmptyMapInitializerListContext(p)
}

func InitEmptyMemberContext(p *MemberContext) {
	canonparser.InitEmptyMemberContext(p)
}

func InitEmptyOptExprContext(p *OptExprContext) {
	canonparser.InitEmptyOptExprContext(p)
}

func InitEmptyOptFieldContext(p *OptFieldContext) {
	canonparser.InitEmptyOptFieldContext(p)
}

func InitEmptyParamContext(p *ParamContext) {
	canonparser.InitEmptyParamContext(p)
}

func InitEmptyParamIdContext(p *ParamIdContext) {
	canonparser.InitEmptyParamIdContext(p)
}

func InitEmptyParseContext(p *ParseContext) {
	canonparser.InitEmptyParseContext(p)
}

func InitEmptyPrimaryContext(p *PrimaryContext) {
	canonparser.InitEmptyPrimaryContext(p)
}

func InitEmptyQualIdContext(p *QualIdContext) {
	canonparser.InitEmptyQualIdContext(p)
}

func InitEmptyRelationContext(p *RelationContext) {
	canonparser.InitEmptyRelationContext(p)
}

func InitEmptySimpleContext(p *SimpleContext) {
	canonparser.InitEmptySimpleContext(p)
}

func InitEmptyStartCommandContext(p *StartCommandContext) {
	canonparser.InitEmptyStartCommandContext(p)
}

func InitEmptyStartContext(p *StartContext) {
	canonparser.InitEmptyStartContext(p)
}

func InitEmptyStartTypeContext(p *StartTypeContext) {
	canonparser.InitEmptyStartTypeContext(p)
}

func InitEmptyTypeContext(p *TypeContext) {
	canonparser.InitEmptyTypeContext(p)
}

func InitEmptyTypeIdContext(p *TypeIdContext) {
	canonparser.InitEmptyTypeIdContext(p)
}

func InitEmptyTypeParamListContext(p *TypeParamListContext) {
	canonparser.InitEmptyTypeParamListContext(p)
}

func InitEmptyUnaryContext(p *UnaryContext) {
	canonparser.InitEmptyUnaryContext(p)
}

func InitEmptyVarDeclContext(p *VarDeclContext) {
	canonparser.InitEmptyVarDeclContext(p)
}

func NewBoolFalseContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolFalseContext {
	return canonparser.NewBoolFalseContext(parser, ctx)
}

func NewBoolTrueContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BoolTrueContext {
	return canonparser.NewBoolTrueContext(parser, ctx)
}

func NewBytesContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *BytesContext {
	return canonparser.NewBytesContext(parser, ctx)
}

func NewCalcContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CalcContext {
	return canonparser.NewCalcContext(parser, parent, invokingState)
}

func NewCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CommandContext {
	return canonparser.NewCommandContext(parser, parent, invokingState)
}

func NewCommandsLexer(input antlr.CharStream) *CommandsLexer {
	return canonparser.NewCommandsLexer(input)
}

func NewCommandsParser(input antlr.TokenStream) *CommandsParser {
	return canonparser.NewCommandsParser(input)
}

func NewCompileContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompileContext {
	return canonparser.NewCompileContext(parser, parent, invokingState)
}

func NewConditionalAndContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalAndContext {
	return canonparser.NewConditionalAndContext(parser, parent, invokingState)
}

func NewConditionalOrContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ConditionalOrContext {
	return canonparser.NewConditionalOrContext(parser, parent, invokingState)
}

func NewConstantLiteralContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ConstantLiteralContext {
	return canonparser.NewConstantLiteralContext(parser, ctx)
}

func NewCreateListContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateListContext {
	return canonparser.NewCreateListContext(parser, ctx)
}

func NewCreateMessageContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateMessageContext {
	return canonparser.NewCreateMessageContext(parser, ctx)
}

func NewCreateStructContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CreateStructContext {
	return canonparser.NewCreateStructContext(parser, ctx)
}

func NewDeclareContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeclareContext {
	return canonparser.NewDeclareContext(parser, parent, invokingState)
}

func NewDeleteContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *DeleteContext {
	return canonparser.NewDeleteContext(parser, parent, invokingState)
}

func NewDoubleContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *DoubleContext {
	return canonparser.NewDoubleContext(parser, ctx)
}

func NewEmptyCalcContext() *CalcContext {
	return canonparser.NewEmptyCalcContext()
}

func NewEmptyCommandContext() *CommandContext {
	return canonparser.NewEmptyCommandContext()
}

func NewEmptyCompileContext() *CompileContext {
	return canonparser.NewEmptyCompileContext()
}

func NewEmptyConditionalAndContext() *ConditionalAndContext {
	return canonparser.NewEmptyConditionalAndContext()
}

func NewEmptyConditionalOrContext() *ConditionalOrContext {
	return canonparser.NewEmptyConditionalOrContext()
}

func NewEmptyContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EmptyContext {
	return canonparser.NewEmptyContext(parser, parent, invokingState)
}

func NewEmptyDeclareContext() *DeclareContext {
	return canonparser.NewEmptyDeclareContext()
}

func NewEmptyDeleteContext() *DeleteContext {
	return canonparser.NewEmptyDeleteContext()
}

func NewEmptyEmptyContext() *EmptyContext {
	return canonparser.NewEmptyEmptyContext()
}

func NewEmptyEscapeIdentContext() *EscapeIdentContext {
	return canonparser.NewEmptyEscapeIdentContext()
}

func NewEmptyExprCmdContext() *ExprCmdContext {
	return canonparser.NewEmptyExprCmdContext()
}

func NewEmptyExprContext() *ExprContext {
	return canonparser.NewEmptyExprContext()
}

func NewEmptyExprListContext() *ExprListContext {
	return canonparser.NewEmptyExprListContext()
}

func NewEmptyFieldInitializerListContext() *FieldInitializerListContext {
	return canonparser.NewEmptyFieldInitializerListContext()
}

func NewEmptyFnDeclContext() *FnDeclContext {
	return canonparser.NewEmptyFnDeclContext()
}

func NewEmptyHelpContext() *HelpContext {
	return canonparser.NewEmptyHelpContext()
}

func NewEmptyLetContext() *LetContext {
	return canonparser.NewEmptyLetContext()
}

func NewEmptyListInitContext() *ListInitContext {
	return canonparser.NewEmptyListInitContext()
}

func NewEmptyLiteralContext() *LiteralContext {
	return canonparser.NewEmptyLiteralContext()
}

func NewEmptyMapInitializerListContext() *MapInitializerListContext {
	return canonparser.NewEmptyMapInitializerListContext()
}

func NewEmptyMemberContext() *MemberContext {
	return canonparser.NewEmptyMemberContext()
}

func NewEmptyOptExprContext() *OptExprContext {
	return canonparser.NewEmptyOptExprContext()
}

func NewEmptyOptFieldContext() *OptFieldContext {
	return canonparser.NewEmptyOptFieldContext()
}

func NewEmptyParamContext() *ParamContext {
	return canonparser.NewEmptyParamContext()
}

func NewEmptyParamIdContext() *ParamIdContext {
	return canonparser.NewEmptyParamIdContext()
}

func NewEmptyParseContext() *ParseContext {
	return canonparser.NewEmptyParseContext()
}

func NewEmptyPrimaryContext() *PrimaryContext {
	return canonparser.NewEmptyPrimaryContext()
}

func NewEmptyQualIdContext() *QualIdContext {
	return canonparser.NewEmptyQualIdContext()
}

func NewEmptyRelationContext() *RelationContext {
	return canonparser.NewEmptyRelationContext()
}

func NewEmptySimpleContext() *SimpleContext {
	return canonparser.NewEmptySimpleContext()
}

func NewEmptyStartCommandContext() *StartCommandContext {
	return canonparser.NewEmptyStartCommandContext()
}

func NewEmptyStartContext() *StartContext {
	return canonparser.NewEmptyStartContext()
}

func NewEmptyStartTypeContext() *StartTypeContext {
	return canonparser.NewEmptyStartTypeContext()
}

func NewEmptyTypeContext() *TypeContext {
	return canonparser.NewEmptyTypeContext()
}

func NewEmptyTypeIdContext() *TypeIdContext {
	return canonparser.NewEmptyTypeIdContext()
}

func NewEmptyTypeParamListContext() *TypeParamListContext {
	return canonparser.NewEmptyTypeParamListContext()
}

func NewEmptyUnaryContext() *UnaryContext {
	return canonparser.NewEmptyUnaryContext()
}

func NewEmptyVarDeclContext() *VarDeclContext {
	return canonparser.NewEmptyVarDeclContext()
}

func NewEscapeIdentContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *EscapeIdentContext {
	return canonparser.NewEscapeIdentContext(parser, parent, invokingState)
}

func NewEscapedIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *EscapedIdentifierContext {
	return canonparser.NewEscapedIdentifierContext(parser, ctx)
}

func NewExprCmdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprCmdContext {
	return canonparser.NewExprCmdContext(parser, parent, invokingState)
}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	return canonparser.NewExprContext(parser, parent, invokingState)
}

func NewExprListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprListContext {
	return canonparser.NewExprListContext(parser, parent, invokingState)
}

func NewFieldInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FieldInitializerListContext {
	return canonparser.NewFieldInitializerListContext(parser, parent, invokingState)
}

func NewFnDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FnDeclContext {
	return canonparser.NewFnDeclContext(parser, parent, invokingState)
}

func NewGlobalCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *GlobalCallContext {
	return canonparser.NewGlobalCallContext(parser, ctx)
}

func NewHelpContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *HelpContext {
	return canonparser.NewHelpContext(parser, parent, invokingState)
}

func NewIdentContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentContext {
	return canonparser.NewIdentContext(parser, ctx)
}

func NewIndexContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IndexContext {
	return canonparser.NewIndexContext(parser, ctx)
}

func NewIntContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IntContext {
	return canonparser.NewIntContext(parser, ctx)
}

func NewLetContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LetContext {
	return canonparser.NewLetContext(parser, parent, invokingState)
}

func NewListInitContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ListInitContext {
	return canonparser.NewListInitContext(parser, parent, invokingState)
}

func NewLiteralContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralContext {
	return canonparser.NewLiteralContext(parser, parent, invokingState)
}

func NewLogicalNotContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *LogicalNotContext {
	return canonparser.NewLogicalNotContext(parser, ctx)
}

func NewMapInitializerListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MapInitializerListContext {
	return canonparser.NewMapInitializerListContext(parser, parent, invokingState)
}

func NewMemberCallContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberCallContext {
	return canonparser.NewMemberCallContext(parser, ctx)
}

func NewMemberContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *MemberContext {
	return canonparser.NewMemberContext(parser, parent, invokingState)
}

func NewMemberExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *MemberExprContext {
	return canonparser.NewMemberExprContext(parser, ctx)
}

func NewNegateContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NegateContext {
	return canonparser.NewNegateContext(parser, ctx)
}

func NewNestedContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NestedContext {
	return canonparser.NewNestedContext(parser, ctx)
}

func NewNullContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullContext {
	return canonparser.NewNullContext(parser, ctx)
}

func NewOptExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptExprContext {
	return canonparser.NewOptExprContext(parser, parent, invokingState)
}

func NewOptFieldContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *OptFieldContext {
	return canonparser.NewOptFieldContext(parser, parent, invokingState)
}

func NewParamContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamContext {
	return canonparser.NewParamContext(parser, parent, invokingState)
}

func NewParamIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParamIdContext {
	return canonparser.NewParamIdContext(parser, parent, invokingState)
}

func NewParseContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ParseContext {
	return canonparser.NewParseContext(parser, parent, invokingState)
}

func NewPrimaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *PrimaryContext {
	return canonparser.NewPrimaryContext(parser, parent, invokingState)
}

func NewPrimaryExprContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *PrimaryExprContext {
	return canonparser.NewPrimaryExprContext(parser, ctx)
}

func NewQualIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *QualIdContext {
	return canonparser.NewQualIdContext(parser, parent, invokingState)
}

func NewRelationContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *RelationContext {
	return canonparser.NewRelationContext(parser, parent, invokingState)
}

func NewSelectContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectContext {
	return canonparser.NewSelectContext(parser, ctx)
}

func NewSimpleContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleContext {
	return canonparser.NewSimpleContext(parser, parent, invokingState)
}

func NewSimpleIdentifierContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleIdentifierContext {
	return canonparser.NewSimpleIdentifierContext(parser, ctx)
}

func NewStartCommandContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartCommandContext {
	return canonparser.NewStartCommandContext(parser, parent, invokingState)
}

func NewStartContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartContext {
	return canonparser.NewStartContext(parser, parent, invokingState)
}

func NewStartTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *StartTypeContext {
	return canonparser.NewStartTypeContext(parser, parent, invokingState)
}

func NewStringContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *StringContext {
	return canonparser.NewStringContext(parser, ctx)
}

func NewTypeContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeContext {
	return canonparser.NewTypeContext(parser, parent, invokingState)
}

func NewTypeIdContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeIdContext {
	return canonparser.NewTypeIdContext(parser, parent, invokingState)
}

func NewTypeParamListContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TypeParamListContext {
	return canonparser.NewTypeParamListContext(parser, parent, invokingState)
}

func NewUintContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *UintContext {
	return canonparser.NewUintContext(parser, ctx)
}

func NewUnaryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *UnaryContext {
	return canonparser.NewUnaryContext(parser, parent, invokingState)
}

func NewVarDeclContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *VarDeclContext {
	return canonparser.NewVarDeclContext(parser, parent, invokingState)
}
