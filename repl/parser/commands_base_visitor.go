// Code generated from java-escape by ANTLR 4.11.1. DO NOT EDIT.

package parser // Commands
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

type BaseCommandsVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseCommandsVisitor) VisitStartCommand(ctx *StartCommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitCommand(ctx *CommandContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitLet(ctx *LetContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitDeclare(ctx *DeclareContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitVarDecl(ctx *VarDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitFnDecl(ctx *FnDeclContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitParam(ctx *ParamContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitDelete(ctx *DeleteContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitSimple(ctx *SimpleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitEmpty(ctx *EmptyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitExprCmd(ctx *ExprCmdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitQualId(ctx *QualIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitStartType(ctx *StartTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitType(ctx *TypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitTypeId(ctx *TypeIdContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitTypeParamList(ctx *TypeParamListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitStart(ctx *StartContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitExpr(ctx *ExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitConditionalOr(ctx *ConditionalOrContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitConditionalAnd(ctx *ConditionalAndContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitRelation(ctx *RelationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitCalc(ctx *CalcContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitMemberExpr(ctx *MemberExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitLogicalNot(ctx *LogicalNotContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitNegate(ctx *NegateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitSelectOrCall(ctx *SelectOrCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitPrimaryExpr(ctx *PrimaryExprContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitIndex(ctx *IndexContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitCreateMessage(ctx *CreateMessageContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitNested(ctx *NestedContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitCreateList(ctx *CreateListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitCreateStruct(ctx *CreateStructContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitConstantLiteral(ctx *ConstantLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitExprList(ctx *ExprListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitFieldInitializerList(ctx *FieldInitializerListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitMapInitializerList(ctx *MapInitializerListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitInt(ctx *IntContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitUint(ctx *UintContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitDouble(ctx *DoubleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitString(ctx *StringContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitBytes(ctx *BytesContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitBoolTrue(ctx *BoolTrueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitBoolFalse(ctx *BoolFalseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseCommandsVisitor) VisitNull(ctx *NullContext) interface{} {
	return v.VisitChildren(ctx)
}
