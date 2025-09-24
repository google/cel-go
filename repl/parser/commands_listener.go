// Code generated from ./Commands.g4 by ANTLR 4.13.1. DO NOT EDIT.

package parser // Commands
import "github.com/antlr4-go/antlr/v4"

// CommandsListener is a complete listener for a parse tree produced by CommandsParser.
type CommandsListener interface {
	antlr.ParseTreeListener

	// EnterStartCommand is called when entering the startCommand production.
	EnterStartCommand(c *StartCommandContext)

	// EnterCommand is called when entering the command production.
	EnterCommand(c *CommandContext)

	// EnterHelp is called when entering the help production.
	EnterHelp(c *HelpContext)

	// EnterLet is called when entering the let production.
	EnterLet(c *LetContext)

	// EnterDeclare is called when entering the declare production.
	EnterDeclare(c *DeclareContext)

	// EnterVarDecl is called when entering the varDecl production.
	EnterVarDecl(c *VarDeclContext)

	// EnterFnDecl is called when entering the fnDecl production.
	EnterFnDecl(c *FnDeclContext)

	// EnterParam is called when entering the param production.
	EnterParam(c *ParamContext)

	// EnterDelete is called when entering the delete production.
	EnterDelete(c *DeleteContext)

	// EnterSimple is called when entering the simple production.
	EnterSimple(c *SimpleContext)

	// EnterEmpty is called when entering the empty production.
	EnterEmpty(c *EmptyContext)

	// EnterCompile is called when entering the compile production.
	EnterCompile(c *CompileContext)

	// EnterParse is called when entering the parse production.
	EnterParse(c *ParseContext)

	// EnterExprCmd is called when entering the exprCmd production.
	EnterExprCmd(c *ExprCmdContext)

	// EnterQualId is called when entering the qualId production.
	EnterQualId(c *QualIdContext)

	// EnterStartType is called when entering the startType production.
	EnterStartType(c *StartTypeContext)

	// EnterType is called when entering the type production.
	EnterType(c *TypeContext)

	// EnterTypeId is called when entering the typeId production.
	EnterTypeId(c *TypeIdContext)

	// EnterTypeParamList is called when entering the typeParamList production.
	EnterTypeParamList(c *TypeParamListContext)

	// EnterStart is called when entering the start production.
	EnterStart(c *StartContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterConditionalOr is called when entering the conditionalOr production.
	EnterConditionalOr(c *ConditionalOrContext)

	// EnterConditionalAnd is called when entering the conditionalAnd production.
	EnterConditionalAnd(c *ConditionalAndContext)

	// EnterRelation is called when entering the relation production.
	EnterRelation(c *RelationContext)

	// EnterCalc is called when entering the calc production.
	EnterCalc(c *CalcContext)

	// EnterMemberExpr is called when entering the MemberExpr production.
	EnterMemberExpr(c *MemberExprContext)

	// EnterLogicalNot is called when entering the LogicalNot production.
	EnterLogicalNot(c *LogicalNotContext)

	// EnterNegate is called when entering the Negate production.
	EnterNegate(c *NegateContext)

	// EnterMemberCall is called when entering the MemberCall production.
	EnterMemberCall(c *MemberCallContext)

	// EnterSelect is called when entering the Select production.
	EnterSelect(c *SelectContext)

	// EnterPrimaryExpr is called when entering the PrimaryExpr production.
	EnterPrimaryExpr(c *PrimaryExprContext)

	// EnterIndex is called when entering the Index production.
	EnterIndex(c *IndexContext)

	// EnterIdent is called when entering the Ident production.
	EnterIdent(c *IdentContext)

	// EnterGlobalCall is called when entering the GlobalCall production.
	EnterGlobalCall(c *GlobalCallContext)

	// EnterNested is called when entering the Nested production.
	EnterNested(c *NestedContext)

	// EnterCreateList is called when entering the CreateList production.
	EnterCreateList(c *CreateListContext)

	// EnterCreateStruct is called when entering the CreateStruct production.
	EnterCreateStruct(c *CreateStructContext)

	// EnterCreateMessage is called when entering the CreateMessage production.
	EnterCreateMessage(c *CreateMessageContext)

	// EnterConstantLiteral is called when entering the ConstantLiteral production.
	EnterConstantLiteral(c *ConstantLiteralContext)

	// EnterExprList is called when entering the exprList production.
	EnterExprList(c *ExprListContext)

	// EnterListInit is called when entering the listInit production.
	EnterListInit(c *ListInitContext)

	// EnterFieldInitializerList is called when entering the fieldInitializerList production.
	EnterFieldInitializerList(c *FieldInitializerListContext)

	// EnterOptField is called when entering the optField production.
	EnterOptField(c *OptFieldContext)

	// EnterMapInitializerList is called when entering the mapInitializerList production.
	EnterMapInitializerList(c *MapInitializerListContext)

	// EnterSimpleIdentifier is called when entering the SimpleIdentifier production.
	EnterSimpleIdentifier(c *SimpleIdentifierContext)

	// EnterEscapedIdentifier is called when entering the EscapedIdentifier production.
	EnterEscapedIdentifier(c *EscapedIdentifierContext)

	// EnterOptExpr is called when entering the optExpr production.
	EnterOptExpr(c *OptExprContext)

	// EnterInt is called when entering the Int production.
	EnterInt(c *IntContext)

	// EnterUint is called when entering the Uint production.
	EnterUint(c *UintContext)

	// EnterDouble is called when entering the Double production.
	EnterDouble(c *DoubleContext)

	// EnterString is called when entering the String production.
	EnterString(c *StringContext)

	// EnterBytes is called when entering the Bytes production.
	EnterBytes(c *BytesContext)

	// EnterBoolTrue is called when entering the BoolTrue production.
	EnterBoolTrue(c *BoolTrueContext)

	// EnterBoolFalse is called when entering the BoolFalse production.
	EnterBoolFalse(c *BoolFalseContext)

	// EnterNull is called when entering the Null production.
	EnterNull(c *NullContext)

	// ExitStartCommand is called when exiting the startCommand production.
	ExitStartCommand(c *StartCommandContext)

	// ExitCommand is called when exiting the command production.
	ExitCommand(c *CommandContext)

	// ExitHelp is called when exiting the help production.
	ExitHelp(c *HelpContext)

	// ExitLet is called when exiting the let production.
	ExitLet(c *LetContext)

	// ExitDeclare is called when exiting the declare production.
	ExitDeclare(c *DeclareContext)

	// ExitVarDecl is called when exiting the varDecl production.
	ExitVarDecl(c *VarDeclContext)

	// ExitFnDecl is called when exiting the fnDecl production.
	ExitFnDecl(c *FnDeclContext)

	// ExitParam is called when exiting the param production.
	ExitParam(c *ParamContext)

	// ExitDelete is called when exiting the delete production.
	ExitDelete(c *DeleteContext)

	// ExitSimple is called when exiting the simple production.
	ExitSimple(c *SimpleContext)

	// ExitEmpty is called when exiting the empty production.
	ExitEmpty(c *EmptyContext)

	// ExitCompile is called when exiting the compile production.
	ExitCompile(c *CompileContext)

	// ExitParse is called when exiting the parse production.
	ExitParse(c *ParseContext)

	// ExitExprCmd is called when exiting the exprCmd production.
	ExitExprCmd(c *ExprCmdContext)

	// ExitQualId is called when exiting the qualId production.
	ExitQualId(c *QualIdContext)

	// ExitStartType is called when exiting the startType production.
	ExitStartType(c *StartTypeContext)

	// ExitType is called when exiting the type production.
	ExitType(c *TypeContext)

	// ExitTypeId is called when exiting the typeId production.
	ExitTypeId(c *TypeIdContext)

	// ExitTypeParamList is called when exiting the typeParamList production.
	ExitTypeParamList(c *TypeParamListContext)

	// ExitStart is called when exiting the start production.
	ExitStart(c *StartContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitConditionalOr is called when exiting the conditionalOr production.
	ExitConditionalOr(c *ConditionalOrContext)

	// ExitConditionalAnd is called when exiting the conditionalAnd production.
	ExitConditionalAnd(c *ConditionalAndContext)

	// ExitRelation is called when exiting the relation production.
	ExitRelation(c *RelationContext)

	// ExitCalc is called when exiting the calc production.
	ExitCalc(c *CalcContext)

	// ExitMemberExpr is called when exiting the MemberExpr production.
	ExitMemberExpr(c *MemberExprContext)

	// ExitLogicalNot is called when exiting the LogicalNot production.
	ExitLogicalNot(c *LogicalNotContext)

	// ExitNegate is called when exiting the Negate production.
	ExitNegate(c *NegateContext)

	// ExitMemberCall is called when exiting the MemberCall production.
	ExitMemberCall(c *MemberCallContext)

	// ExitSelect is called when exiting the Select production.
	ExitSelect(c *SelectContext)

	// ExitPrimaryExpr is called when exiting the PrimaryExpr production.
	ExitPrimaryExpr(c *PrimaryExprContext)

	// ExitIndex is called when exiting the Index production.
	ExitIndex(c *IndexContext)

	// ExitIdent is called when exiting the Ident production.
	ExitIdent(c *IdentContext)

	// ExitGlobalCall is called when exiting the GlobalCall production.
	ExitGlobalCall(c *GlobalCallContext)

	// ExitNested is called when exiting the Nested production.
	ExitNested(c *NestedContext)

	// ExitCreateList is called when exiting the CreateList production.
	ExitCreateList(c *CreateListContext)

	// ExitCreateStruct is called when exiting the CreateStruct production.
	ExitCreateStruct(c *CreateStructContext)

	// ExitCreateMessage is called when exiting the CreateMessage production.
	ExitCreateMessage(c *CreateMessageContext)

	// ExitConstantLiteral is called when exiting the ConstantLiteral production.
	ExitConstantLiteral(c *ConstantLiteralContext)

	// ExitExprList is called when exiting the exprList production.
	ExitExprList(c *ExprListContext)

	// ExitListInit is called when exiting the listInit production.
	ExitListInit(c *ListInitContext)

	// ExitFieldInitializerList is called when exiting the fieldInitializerList production.
	ExitFieldInitializerList(c *FieldInitializerListContext)

	// ExitOptField is called when exiting the optField production.
	ExitOptField(c *OptFieldContext)

	// ExitMapInitializerList is called when exiting the mapInitializerList production.
	ExitMapInitializerList(c *MapInitializerListContext)

	// ExitSimpleIdentifier is called when exiting the SimpleIdentifier production.
	ExitSimpleIdentifier(c *SimpleIdentifierContext)

	// ExitEscapedIdentifier is called when exiting the EscapedIdentifier production.
	ExitEscapedIdentifier(c *EscapedIdentifierContext)

	// ExitOptExpr is called when exiting the optExpr production.
	ExitOptExpr(c *OptExprContext)

	// ExitInt is called when exiting the Int production.
	ExitInt(c *IntContext)

	// ExitUint is called when exiting the Uint production.
	ExitUint(c *UintContext)

	// ExitDouble is called when exiting the Double production.
	ExitDouble(c *DoubleContext)

	// ExitString is called when exiting the String production.
	ExitString(c *StringContext)

	// ExitBytes is called when exiting the Bytes production.
	ExitBytes(c *BytesContext)

	// ExitBoolTrue is called when exiting the BoolTrue production.
	ExitBoolTrue(c *BoolTrueContext)

	// ExitBoolFalse is called when exiting the BoolFalse production.
	ExitBoolFalse(c *BoolFalseContext)

	// ExitNull is called when exiting the Null production.
	ExitNull(c *NullContext)
}
