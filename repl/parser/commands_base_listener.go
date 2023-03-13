// Code generated from ./Commands.g4 by ANTLR 4.12.0. DO NOT EDIT.

package parser // Commands
import "github.com/antlr/antlr4/runtime/Go/antlr/v4"

// BaseCommandsListener is a complete listener for a parse tree produced by CommandsParser.
type BaseCommandsListener struct{}

var _ CommandsListener = &BaseCommandsListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseCommandsListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseCommandsListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseCommandsListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseCommandsListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStartCommand is called when production startCommand is entered.
func (s *BaseCommandsListener) EnterStartCommand(ctx *StartCommandContext) {}

// ExitStartCommand is called when production startCommand is exited.
func (s *BaseCommandsListener) ExitStartCommand(ctx *StartCommandContext) {}

// EnterCommand is called when production command is entered.
func (s *BaseCommandsListener) EnterCommand(ctx *CommandContext) {}

// ExitCommand is called when production command is exited.
func (s *BaseCommandsListener) ExitCommand(ctx *CommandContext) {}

// EnterLet is called when production let is entered.
func (s *BaseCommandsListener) EnterLet(ctx *LetContext) {}

// ExitLet is called when production let is exited.
func (s *BaseCommandsListener) ExitLet(ctx *LetContext) {}

// EnterDeclare is called when production declare is entered.
func (s *BaseCommandsListener) EnterDeclare(ctx *DeclareContext) {}

// ExitDeclare is called when production declare is exited.
func (s *BaseCommandsListener) ExitDeclare(ctx *DeclareContext) {}

// EnterVarDecl is called when production varDecl is entered.
func (s *BaseCommandsListener) EnterVarDecl(ctx *VarDeclContext) {}

// ExitVarDecl is called when production varDecl is exited.
func (s *BaseCommandsListener) ExitVarDecl(ctx *VarDeclContext) {}

// EnterFnDecl is called when production fnDecl is entered.
func (s *BaseCommandsListener) EnterFnDecl(ctx *FnDeclContext) {}

// ExitFnDecl is called when production fnDecl is exited.
func (s *BaseCommandsListener) ExitFnDecl(ctx *FnDeclContext) {}

// EnterParam is called when production param is entered.
func (s *BaseCommandsListener) EnterParam(ctx *ParamContext) {}

// ExitParam is called when production param is exited.
func (s *BaseCommandsListener) ExitParam(ctx *ParamContext) {}

// EnterDelete is called when production delete is entered.
func (s *BaseCommandsListener) EnterDelete(ctx *DeleteContext) {}

// ExitDelete is called when production delete is exited.
func (s *BaseCommandsListener) ExitDelete(ctx *DeleteContext) {}

// EnterSimple is called when production simple is entered.
func (s *BaseCommandsListener) EnterSimple(ctx *SimpleContext) {}

// ExitSimple is called when production simple is exited.
func (s *BaseCommandsListener) ExitSimple(ctx *SimpleContext) {}

// EnterEmpty is called when production empty is entered.
func (s *BaseCommandsListener) EnterEmpty(ctx *EmptyContext) {}

// ExitEmpty is called when production empty is exited.
func (s *BaseCommandsListener) ExitEmpty(ctx *EmptyContext) {}

// EnterExprCmd is called when production exprCmd is entered.
func (s *BaseCommandsListener) EnterExprCmd(ctx *ExprCmdContext) {}

// ExitExprCmd is called when production exprCmd is exited.
func (s *BaseCommandsListener) ExitExprCmd(ctx *ExprCmdContext) {}

// EnterQualId is called when production qualId is entered.
func (s *BaseCommandsListener) EnterQualId(ctx *QualIdContext) {}

// ExitQualId is called when production qualId is exited.
func (s *BaseCommandsListener) ExitQualId(ctx *QualIdContext) {}

// EnterStartType is called when production startType is entered.
func (s *BaseCommandsListener) EnterStartType(ctx *StartTypeContext) {}

// ExitStartType is called when production startType is exited.
func (s *BaseCommandsListener) ExitStartType(ctx *StartTypeContext) {}

// EnterType is called when production type is entered.
func (s *BaseCommandsListener) EnterType(ctx *TypeContext) {}

// ExitType is called when production type is exited.
func (s *BaseCommandsListener) ExitType(ctx *TypeContext) {}

// EnterTypeId is called when production typeId is entered.
func (s *BaseCommandsListener) EnterTypeId(ctx *TypeIdContext) {}

// ExitTypeId is called when production typeId is exited.
func (s *BaseCommandsListener) ExitTypeId(ctx *TypeIdContext) {}

// EnterTypeParamList is called when production typeParamList is entered.
func (s *BaseCommandsListener) EnterTypeParamList(ctx *TypeParamListContext) {}

// ExitTypeParamList is called when production typeParamList is exited.
func (s *BaseCommandsListener) ExitTypeParamList(ctx *TypeParamListContext) {}

// EnterStart is called when production start is entered.
func (s *BaseCommandsListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseCommandsListener) ExitStart(ctx *StartContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseCommandsListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseCommandsListener) ExitExpr(ctx *ExprContext) {}

// EnterConditionalOr is called when production conditionalOr is entered.
func (s *BaseCommandsListener) EnterConditionalOr(ctx *ConditionalOrContext) {}

// ExitConditionalOr is called when production conditionalOr is exited.
func (s *BaseCommandsListener) ExitConditionalOr(ctx *ConditionalOrContext) {}

// EnterConditionalAnd is called when production conditionalAnd is entered.
func (s *BaseCommandsListener) EnterConditionalAnd(ctx *ConditionalAndContext) {}

// ExitConditionalAnd is called when production conditionalAnd is exited.
func (s *BaseCommandsListener) ExitConditionalAnd(ctx *ConditionalAndContext) {}

// EnterRelation is called when production relation is entered.
func (s *BaseCommandsListener) EnterRelation(ctx *RelationContext) {}

// ExitRelation is called when production relation is exited.
func (s *BaseCommandsListener) ExitRelation(ctx *RelationContext) {}

// EnterCalc is called when production calc is entered.
func (s *BaseCommandsListener) EnterCalc(ctx *CalcContext) {}

// ExitCalc is called when production calc is exited.
func (s *BaseCommandsListener) ExitCalc(ctx *CalcContext) {}

// EnterMemberExpr is called when production MemberExpr is entered.
func (s *BaseCommandsListener) EnterMemberExpr(ctx *MemberExprContext) {}

// ExitMemberExpr is called when production MemberExpr is exited.
func (s *BaseCommandsListener) ExitMemberExpr(ctx *MemberExprContext) {}

// EnterLogicalNot is called when production LogicalNot is entered.
func (s *BaseCommandsListener) EnterLogicalNot(ctx *LogicalNotContext) {}

// ExitLogicalNot is called when production LogicalNot is exited.
func (s *BaseCommandsListener) ExitLogicalNot(ctx *LogicalNotContext) {}

// EnterNegate is called when production Negate is entered.
func (s *BaseCommandsListener) EnterNegate(ctx *NegateContext) {}

// ExitNegate is called when production Negate is exited.
func (s *BaseCommandsListener) ExitNegate(ctx *NegateContext) {}

// EnterMemberCall is called when production MemberCall is entered.
func (s *BaseCommandsListener) EnterMemberCall(ctx *MemberCallContext) {}

// ExitMemberCall is called when production MemberCall is exited.
func (s *BaseCommandsListener) ExitMemberCall(ctx *MemberCallContext) {}

// EnterSelect is called when production Select is entered.
func (s *BaseCommandsListener) EnterSelect(ctx *SelectContext) {}

// ExitSelect is called when production Select is exited.
func (s *BaseCommandsListener) ExitSelect(ctx *SelectContext) {}

// EnterPrimaryExpr is called when production PrimaryExpr is entered.
func (s *BaseCommandsListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production PrimaryExpr is exited.
func (s *BaseCommandsListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterIndex is called when production Index is entered.
func (s *BaseCommandsListener) EnterIndex(ctx *IndexContext) {}

// ExitIndex is called when production Index is exited.
func (s *BaseCommandsListener) ExitIndex(ctx *IndexContext) {}

// EnterIdentOrGlobalCall is called when production IdentOrGlobalCall is entered.
func (s *BaseCommandsListener) EnterIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) {}

// ExitIdentOrGlobalCall is called when production IdentOrGlobalCall is exited.
func (s *BaseCommandsListener) ExitIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) {}

// EnterNested is called when production Nested is entered.
func (s *BaseCommandsListener) EnterNested(ctx *NestedContext) {}

// ExitNested is called when production Nested is exited.
func (s *BaseCommandsListener) ExitNested(ctx *NestedContext) {}

// EnterCreateList is called when production CreateList is entered.
func (s *BaseCommandsListener) EnterCreateList(ctx *CreateListContext) {}

// ExitCreateList is called when production CreateList is exited.
func (s *BaseCommandsListener) ExitCreateList(ctx *CreateListContext) {}

// EnterCreateStruct is called when production CreateStruct is entered.
func (s *BaseCommandsListener) EnterCreateStruct(ctx *CreateStructContext) {}

// ExitCreateStruct is called when production CreateStruct is exited.
func (s *BaseCommandsListener) ExitCreateStruct(ctx *CreateStructContext) {}

// EnterCreateMessage is called when production CreateMessage is entered.
func (s *BaseCommandsListener) EnterCreateMessage(ctx *CreateMessageContext) {}

// ExitCreateMessage is called when production CreateMessage is exited.
func (s *BaseCommandsListener) ExitCreateMessage(ctx *CreateMessageContext) {}

// EnterConstantLiteral is called when production ConstantLiteral is entered.
func (s *BaseCommandsListener) EnterConstantLiteral(ctx *ConstantLiteralContext) {}

// ExitConstantLiteral is called when production ConstantLiteral is exited.
func (s *BaseCommandsListener) ExitConstantLiteral(ctx *ConstantLiteralContext) {}

// EnterExprList is called when production exprList is entered.
func (s *BaseCommandsListener) EnterExprList(ctx *ExprListContext) {}

// ExitExprList is called when production exprList is exited.
func (s *BaseCommandsListener) ExitExprList(ctx *ExprListContext) {}

// EnterListInit is called when production listInit is entered.
func (s *BaseCommandsListener) EnterListInit(ctx *ListInitContext) {}

// ExitListInit is called when production listInit is exited.
func (s *BaseCommandsListener) ExitListInit(ctx *ListInitContext) {}

// EnterFieldInitializerList is called when production fieldInitializerList is entered.
func (s *BaseCommandsListener) EnterFieldInitializerList(ctx *FieldInitializerListContext) {}

// ExitFieldInitializerList is called when production fieldInitializerList is exited.
func (s *BaseCommandsListener) ExitFieldInitializerList(ctx *FieldInitializerListContext) {}

// EnterOptField is called when production optField is entered.
func (s *BaseCommandsListener) EnterOptField(ctx *OptFieldContext) {}

// ExitOptField is called when production optField is exited.
func (s *BaseCommandsListener) ExitOptField(ctx *OptFieldContext) {}

// EnterMapInitializerList is called when production mapInitializerList is entered.
func (s *BaseCommandsListener) EnterMapInitializerList(ctx *MapInitializerListContext) {}

// ExitMapInitializerList is called when production mapInitializerList is exited.
func (s *BaseCommandsListener) ExitMapInitializerList(ctx *MapInitializerListContext) {}

// EnterOptExpr is called when production optExpr is entered.
func (s *BaseCommandsListener) EnterOptExpr(ctx *OptExprContext) {}

// ExitOptExpr is called when production optExpr is exited.
func (s *BaseCommandsListener) ExitOptExpr(ctx *OptExprContext) {}

// EnterInt is called when production Int is entered.
func (s *BaseCommandsListener) EnterInt(ctx *IntContext) {}

// ExitInt is called when production Int is exited.
func (s *BaseCommandsListener) ExitInt(ctx *IntContext) {}

// EnterUint is called when production Uint is entered.
func (s *BaseCommandsListener) EnterUint(ctx *UintContext) {}

// ExitUint is called when production Uint is exited.
func (s *BaseCommandsListener) ExitUint(ctx *UintContext) {}

// EnterDouble is called when production Double is entered.
func (s *BaseCommandsListener) EnterDouble(ctx *DoubleContext) {}

// ExitDouble is called when production Double is exited.
func (s *BaseCommandsListener) ExitDouble(ctx *DoubleContext) {}

// EnterString is called when production String is entered.
func (s *BaseCommandsListener) EnterString(ctx *StringContext) {}

// ExitString is called when production String is exited.
func (s *BaseCommandsListener) ExitString(ctx *StringContext) {}

// EnterBytes is called when production Bytes is entered.
func (s *BaseCommandsListener) EnterBytes(ctx *BytesContext) {}

// ExitBytes is called when production Bytes is exited.
func (s *BaseCommandsListener) ExitBytes(ctx *BytesContext) {}

// EnterBoolTrue is called when production BoolTrue is entered.
func (s *BaseCommandsListener) EnterBoolTrue(ctx *BoolTrueContext) {}

// ExitBoolTrue is called when production BoolTrue is exited.
func (s *BaseCommandsListener) ExitBoolTrue(ctx *BoolTrueContext) {}

// EnterBoolFalse is called when production BoolFalse is entered.
func (s *BaseCommandsListener) EnterBoolFalse(ctx *BoolFalseContext) {}

// ExitBoolFalse is called when production BoolFalse is exited.
func (s *BaseCommandsListener) ExitBoolFalse(ctx *BoolFalseContext) {}

// EnterNull is called when production Null is entered.
func (s *BaseCommandsListener) EnterNull(ctx *NullContext) {}

// ExitNull is called when production Null is exited.
func (s *BaseCommandsListener) ExitNull(ctx *NullContext) {}
