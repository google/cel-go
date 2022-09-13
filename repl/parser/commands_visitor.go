// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Commands
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// A complete Visitor for a parse tree produced by CommandsParser.
type CommandsVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by CommandsParser#startCommand.
	VisitStartCommand(ctx *StartCommandContext) interface{}

	// Visit a parse tree produced by CommandsParser#command.
	VisitCommand(ctx *CommandContext) interface{}

	// Visit a parse tree produced by CommandsParser#let.
	VisitLet(ctx *LetContext) interface{}

	// Visit a parse tree produced by CommandsParser#declare.
	VisitDeclare(ctx *DeclareContext) interface{}

	// Visit a parse tree produced by CommandsParser#varDecl.
	VisitVarDecl(ctx *VarDeclContext) interface{}

	// Visit a parse tree produced by CommandsParser#fnDecl.
	VisitFnDecl(ctx *FnDeclContext) interface{}

	// Visit a parse tree produced by CommandsParser#param.
	VisitParam(ctx *ParamContext) interface{}

	// Visit a parse tree produced by CommandsParser#delete.
	VisitDelete(ctx *DeleteContext) interface{}

	// Visit a parse tree produced by CommandsParser#simple.
	VisitSimple(ctx *SimpleContext) interface{}

	// Visit a parse tree produced by CommandsParser#empty.
	VisitEmpty(ctx *EmptyContext) interface{}

	// Visit a parse tree produced by CommandsParser#exprCmd.
	VisitExprCmd(ctx *ExprCmdContext) interface{}

	// Visit a parse tree produced by CommandsParser#qualId.
	VisitQualId(ctx *QualIdContext) interface{}

	// Visit a parse tree produced by CommandsParser#startType.
	VisitStartType(ctx *StartTypeContext) interface{}

	// Visit a parse tree produced by CommandsParser#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by CommandsParser#typeId.
	VisitTypeId(ctx *TypeIdContext) interface{}

	// Visit a parse tree produced by CommandsParser#typeParamList.
	VisitTypeParamList(ctx *TypeParamListContext) interface{}

	// Visit a parse tree produced by CommandsParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by CommandsParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by CommandsParser#conditionalOr.
	VisitConditionalOr(ctx *ConditionalOrContext) interface{}

	// Visit a parse tree produced by CommandsParser#conditionalAnd.
	VisitConditionalAnd(ctx *ConditionalAndContext) interface{}

	// Visit a parse tree produced by CommandsParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by CommandsParser#calc.
	VisitCalc(ctx *CalcContext) interface{}

	// Visit a parse tree produced by CommandsParser#MemberExpr.
	VisitMemberExpr(ctx *MemberExprContext) interface{}

	// Visit a parse tree produced by CommandsParser#LogicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by CommandsParser#Negate.
	VisitNegate(ctx *NegateContext) interface{}

	// Visit a parse tree produced by CommandsParser#SelectOrCall.
	VisitSelectOrCall(ctx *SelectOrCallContext) interface{}

	// Visit a parse tree produced by CommandsParser#PrimaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by CommandsParser#Index.
	VisitIndex(ctx *IndexContext) interface{}

	// Visit a parse tree produced by CommandsParser#CreateMessage.
	VisitCreateMessage(ctx *CreateMessageContext) interface{}

	// Visit a parse tree produced by CommandsParser#IdentOrGlobalCall.
	VisitIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) interface{}

	// Visit a parse tree produced by CommandsParser#Nested.
	VisitNested(ctx *NestedContext) interface{}

	// Visit a parse tree produced by CommandsParser#CreateList.
	VisitCreateList(ctx *CreateListContext) interface{}

	// Visit a parse tree produced by CommandsParser#CreateStruct.
	VisitCreateStruct(ctx *CreateStructContext) interface{}

	// Visit a parse tree produced by CommandsParser#ConstantLiteral.
	VisitConstantLiteral(ctx *ConstantLiteralContext) interface{}

	// Visit a parse tree produced by CommandsParser#exprList.
	VisitExprList(ctx *ExprListContext) interface{}

	// Visit a parse tree produced by CommandsParser#fieldInitializerList.
	VisitFieldInitializerList(ctx *FieldInitializerListContext) interface{}

	// Visit a parse tree produced by CommandsParser#mapInitializerList.
	VisitMapInitializerList(ctx *MapInitializerListContext) interface{}

	// Visit a parse tree produced by CommandsParser#Int.
	VisitInt(ctx *IntContext) interface{}

	// Visit a parse tree produced by CommandsParser#Uint.
	VisitUint(ctx *UintContext) interface{}

	// Visit a parse tree produced by CommandsParser#Double.
	VisitDouble(ctx *DoubleContext) interface{}

	// Visit a parse tree produced by CommandsParser#String.
	VisitString(ctx *StringContext) interface{}

	// Visit a parse tree produced by CommandsParser#Bytes.
	VisitBytes(ctx *BytesContext) interface{}

	// Visit a parse tree produced by CommandsParser#BoolTrue.
	VisitBoolTrue(ctx *BoolTrueContext) interface{}

	// Visit a parse tree produced by CommandsParser#BoolFalse.
	VisitBoolFalse(ctx *BoolFalseContext) interface{}

	// Visit a parse tree produced by CommandsParser#Null.
	VisitNull(ctx *NullContext) interface{}
}
