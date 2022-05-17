// Code generated from ./Types.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Types
import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseTypesListener is a complete listener for a parse tree produced by TypesParser.
type BaseTypesListener struct{}

var _ TypesListener = &BaseTypesListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseTypesListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseTypesListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseTypesListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseTypesListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterStart is called when production start is entered.
func (s *BaseTypesListener) EnterStart(ctx *StartContext) {}

// ExitStart is called when production start is exited.
func (s *BaseTypesListener) ExitStart(ctx *StartContext) {}

// EnterType is called when production type is entered.
func (s *BaseTypesListener) EnterType(ctx *TypeContext) {}

// ExitType is called when production type is exited.
func (s *BaseTypesListener) ExitType(ctx *TypeContext) {}

// EnterTypeId is called when production typeId is entered.
func (s *BaseTypesListener) EnterTypeId(ctx *TypeIdContext) {}

// ExitTypeId is called when production typeId is exited.
func (s *BaseTypesListener) ExitTypeId(ctx *TypeIdContext) {}

// EnterTypeParamList is called when production typeParamList is entered.
func (s *BaseTypesListener) EnterTypeParamList(ctx *TypeParamListContext) {}

// ExitTypeParamList is called when production typeParamList is exited.
func (s *BaseTypesListener) ExitTypeParamList(ctx *TypeParamListContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseTypesListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseTypesListener) ExitExpr(ctx *ExprContext) {}

// EnterConditionalOr is called when production conditionalOr is entered.
func (s *BaseTypesListener) EnterConditionalOr(ctx *ConditionalOrContext) {}

// ExitConditionalOr is called when production conditionalOr is exited.
func (s *BaseTypesListener) ExitConditionalOr(ctx *ConditionalOrContext) {}

// EnterConditionalAnd is called when production conditionalAnd is entered.
func (s *BaseTypesListener) EnterConditionalAnd(ctx *ConditionalAndContext) {}

// ExitConditionalAnd is called when production conditionalAnd is exited.
func (s *BaseTypesListener) ExitConditionalAnd(ctx *ConditionalAndContext) {}

// EnterRelation is called when production relation is entered.
func (s *BaseTypesListener) EnterRelation(ctx *RelationContext) {}

// ExitRelation is called when production relation is exited.
func (s *BaseTypesListener) ExitRelation(ctx *RelationContext) {}

// EnterCalc is called when production calc is entered.
func (s *BaseTypesListener) EnterCalc(ctx *CalcContext) {}

// ExitCalc is called when production calc is exited.
func (s *BaseTypesListener) ExitCalc(ctx *CalcContext) {}

// EnterMemberExpr is called when production MemberExpr is entered.
func (s *BaseTypesListener) EnterMemberExpr(ctx *MemberExprContext) {}

// ExitMemberExpr is called when production MemberExpr is exited.
func (s *BaseTypesListener) ExitMemberExpr(ctx *MemberExprContext) {}

// EnterLogicalNot is called when production LogicalNot is entered.
func (s *BaseTypesListener) EnterLogicalNot(ctx *LogicalNotContext) {}

// ExitLogicalNot is called when production LogicalNot is exited.
func (s *BaseTypesListener) ExitLogicalNot(ctx *LogicalNotContext) {}

// EnterNegate is called when production Negate is entered.
func (s *BaseTypesListener) EnterNegate(ctx *NegateContext) {}

// ExitNegate is called when production Negate is exited.
func (s *BaseTypesListener) ExitNegate(ctx *NegateContext) {}

// EnterSelectOrCall is called when production SelectOrCall is entered.
func (s *BaseTypesListener) EnterSelectOrCall(ctx *SelectOrCallContext) {}

// ExitSelectOrCall is called when production SelectOrCall is exited.
func (s *BaseTypesListener) ExitSelectOrCall(ctx *SelectOrCallContext) {}

// EnterPrimaryExpr is called when production PrimaryExpr is entered.
func (s *BaseTypesListener) EnterPrimaryExpr(ctx *PrimaryExprContext) {}

// ExitPrimaryExpr is called when production PrimaryExpr is exited.
func (s *BaseTypesListener) ExitPrimaryExpr(ctx *PrimaryExprContext) {}

// EnterIndex is called when production Index is entered.
func (s *BaseTypesListener) EnterIndex(ctx *IndexContext) {}

// ExitIndex is called when production Index is exited.
func (s *BaseTypesListener) ExitIndex(ctx *IndexContext) {}

// EnterCreateMessage is called when production CreateMessage is entered.
func (s *BaseTypesListener) EnterCreateMessage(ctx *CreateMessageContext) {}

// ExitCreateMessage is called when production CreateMessage is exited.
func (s *BaseTypesListener) ExitCreateMessage(ctx *CreateMessageContext) {}

// EnterIdentOrGlobalCall is called when production IdentOrGlobalCall is entered.
func (s *BaseTypesListener) EnterIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) {}

// ExitIdentOrGlobalCall is called when production IdentOrGlobalCall is exited.
func (s *BaseTypesListener) ExitIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) {}

// EnterNested is called when production Nested is entered.
func (s *BaseTypesListener) EnterNested(ctx *NestedContext) {}

// ExitNested is called when production Nested is exited.
func (s *BaseTypesListener) ExitNested(ctx *NestedContext) {}

// EnterCreateList is called when production CreateList is entered.
func (s *BaseTypesListener) EnterCreateList(ctx *CreateListContext) {}

// ExitCreateList is called when production CreateList is exited.
func (s *BaseTypesListener) ExitCreateList(ctx *CreateListContext) {}

// EnterCreateStruct is called when production CreateStruct is entered.
func (s *BaseTypesListener) EnterCreateStruct(ctx *CreateStructContext) {}

// ExitCreateStruct is called when production CreateStruct is exited.
func (s *BaseTypesListener) ExitCreateStruct(ctx *CreateStructContext) {}

// EnterConstantLiteral is called when production ConstantLiteral is entered.
func (s *BaseTypesListener) EnterConstantLiteral(ctx *ConstantLiteralContext) {}

// ExitConstantLiteral is called when production ConstantLiteral is exited.
func (s *BaseTypesListener) ExitConstantLiteral(ctx *ConstantLiteralContext) {}

// EnterExprList is called when production exprList is entered.
func (s *BaseTypesListener) EnterExprList(ctx *ExprListContext) {}

// ExitExprList is called when production exprList is exited.
func (s *BaseTypesListener) ExitExprList(ctx *ExprListContext) {}

// EnterFieldInitializerList is called when production fieldInitializerList is entered.
func (s *BaseTypesListener) EnterFieldInitializerList(ctx *FieldInitializerListContext) {}

// ExitFieldInitializerList is called when production fieldInitializerList is exited.
func (s *BaseTypesListener) ExitFieldInitializerList(ctx *FieldInitializerListContext) {}

// EnterMapInitializerList is called when production mapInitializerList is entered.
func (s *BaseTypesListener) EnterMapInitializerList(ctx *MapInitializerListContext) {}

// ExitMapInitializerList is called when production mapInitializerList is exited.
func (s *BaseTypesListener) ExitMapInitializerList(ctx *MapInitializerListContext) {}

// EnterInt is called when production Int is entered.
func (s *BaseTypesListener) EnterInt(ctx *IntContext) {}

// ExitInt is called when production Int is exited.
func (s *BaseTypesListener) ExitInt(ctx *IntContext) {}

// EnterUint is called when production Uint is entered.
func (s *BaseTypesListener) EnterUint(ctx *UintContext) {}

// ExitUint is called when production Uint is exited.
func (s *BaseTypesListener) ExitUint(ctx *UintContext) {}

// EnterDouble is called when production Double is entered.
func (s *BaseTypesListener) EnterDouble(ctx *DoubleContext) {}

// ExitDouble is called when production Double is exited.
func (s *BaseTypesListener) ExitDouble(ctx *DoubleContext) {}

// EnterString is called when production String is entered.
func (s *BaseTypesListener) EnterString(ctx *StringContext) {}

// ExitString is called when production String is exited.
func (s *BaseTypesListener) ExitString(ctx *StringContext) {}

// EnterBytes is called when production Bytes is entered.
func (s *BaseTypesListener) EnterBytes(ctx *BytesContext) {}

// ExitBytes is called when production Bytes is exited.
func (s *BaseTypesListener) ExitBytes(ctx *BytesContext) {}

// EnterBoolTrue is called when production BoolTrue is entered.
func (s *BaseTypesListener) EnterBoolTrue(ctx *BoolTrueContext) {}

// ExitBoolTrue is called when production BoolTrue is exited.
func (s *BaseTypesListener) ExitBoolTrue(ctx *BoolTrueContext) {}

// EnterBoolFalse is called when production BoolFalse is entered.
func (s *BaseTypesListener) EnterBoolFalse(ctx *BoolFalseContext) {}

// ExitBoolFalse is called when production BoolFalse is exited.
func (s *BaseTypesListener) ExitBoolFalse(ctx *BoolFalseContext) {}

// EnterNull is called when production Null is entered.
func (s *BaseTypesListener) EnterNull(ctx *NullContext) {}

// ExitNull is called when production Null is exited.
func (s *BaseTypesListener) ExitNull(ctx *NullContext) {}
