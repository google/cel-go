// Code generated from ./Types.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Types
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by TypesParser.
type TypesVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by TypesParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by TypesParser#type.
	VisitType(ctx *TypeContext) interface{}

	// Visit a parse tree produced by TypesParser#typeId.
	VisitTypeId(ctx *TypeIdContext) interface{}

	// Visit a parse tree produced by TypesParser#typeParamList.
	VisitTypeParamList(ctx *TypeParamListContext) interface{}

	// Visit a parse tree produced by TypesParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by TypesParser#conditionalOr.
	VisitConditionalOr(ctx *ConditionalOrContext) interface{}

	// Visit a parse tree produced by TypesParser#conditionalAnd.
	VisitConditionalAnd(ctx *ConditionalAndContext) interface{}

	// Visit a parse tree produced by TypesParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by TypesParser#calc.
	VisitCalc(ctx *CalcContext) interface{}

	// Visit a parse tree produced by TypesParser#MemberExpr.
	VisitMemberExpr(ctx *MemberExprContext) interface{}

	// Visit a parse tree produced by TypesParser#LogicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by TypesParser#Negate.
	VisitNegate(ctx *NegateContext) interface{}

	// Visit a parse tree produced by TypesParser#SelectOrCall.
	VisitSelectOrCall(ctx *SelectOrCallContext) interface{}

	// Visit a parse tree produced by TypesParser#PrimaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by TypesParser#Index.
	VisitIndex(ctx *IndexContext) interface{}

	// Visit a parse tree produced by TypesParser#CreateMessage.
	VisitCreateMessage(ctx *CreateMessageContext) interface{}

	// Visit a parse tree produced by TypesParser#IdentOrGlobalCall.
	VisitIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) interface{}

	// Visit a parse tree produced by TypesParser#Nested.
	VisitNested(ctx *NestedContext) interface{}

	// Visit a parse tree produced by TypesParser#CreateList.
	VisitCreateList(ctx *CreateListContext) interface{}

	// Visit a parse tree produced by TypesParser#CreateStruct.
	VisitCreateStruct(ctx *CreateStructContext) interface{}

	// Visit a parse tree produced by TypesParser#ConstantLiteral.
	VisitConstantLiteral(ctx *ConstantLiteralContext) interface{}

	// Visit a parse tree produced by TypesParser#exprList.
	VisitExprList(ctx *ExprListContext) interface{}

	// Visit a parse tree produced by TypesParser#fieldInitializerList.
	VisitFieldInitializerList(ctx *FieldInitializerListContext) interface{}

	// Visit a parse tree produced by TypesParser#mapInitializerList.
	VisitMapInitializerList(ctx *MapInitializerListContext) interface{}

	// Visit a parse tree produced by TypesParser#Int.
	VisitInt(ctx *IntContext) interface{}

	// Visit a parse tree produced by TypesParser#Uint.
	VisitUint(ctx *UintContext) interface{}

	// Visit a parse tree produced by TypesParser#Double.
	VisitDouble(ctx *DoubleContext) interface{}

	// Visit a parse tree produced by TypesParser#String.
	VisitString(ctx *StringContext) interface{}

	// Visit a parse tree produced by TypesParser#Bytes.
	VisitBytes(ctx *BytesContext) interface{}

	// Visit a parse tree produced by TypesParser#BoolTrue.
	VisitBoolTrue(ctx *BoolTrueContext) interface{}

	// Visit a parse tree produced by TypesParser#BoolFalse.
	VisitBoolFalse(ctx *BoolFalseContext) interface{}

	// Visit a parse tree produced by TypesParser#Null.
	VisitNull(ctx *NullContext) interface{}
}
