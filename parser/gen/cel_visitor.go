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

// Generated from /Users/ozben/go/src/celgo/bin/../parser/gen/CEL.g4 by ANTLR 4.7.

package gen // CEL
import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by CELParser.
type CELVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by CELParser#start.
	VisitStart(ctx *StartContext) interface{}

	// Visit a parse tree produced by CELParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by CELParser#conditionalOr.
	VisitConditionalOr(ctx *ConditionalOrContext) interface{}

	// Visit a parse tree produced by CELParser#conditionalAnd.
	VisitConditionalAnd(ctx *ConditionalAndContext) interface{}

	// Visit a parse tree produced by CELParser#relation.
	VisitRelation(ctx *RelationContext) interface{}

	// Visit a parse tree produced by CELParser#calc.
	VisitCalc(ctx *CalcContext) interface{}

	// Visit a parse tree produced by CELParser#StatementExpr.
	VisitStatementExpr(ctx *StatementExprContext) interface{}

	// Visit a parse tree produced by CELParser#LogicalNot.
	VisitLogicalNot(ctx *LogicalNotContext) interface{}

	// Visit a parse tree produced by CELParser#Negate.
	VisitNegate(ctx *NegateContext) interface{}

	// Visit a parse tree produced by CELParser#SelectOrCall.
	VisitSelectOrCall(ctx *SelectOrCallContext) interface{}

	// Visit a parse tree produced by CELParser#PrimaryExpr.
	VisitPrimaryExpr(ctx *PrimaryExprContext) interface{}

	// Visit a parse tree produced by CELParser#Index.
	VisitIndex(ctx *IndexContext) interface{}

	// Visit a parse tree produced by CELParser#CreateMessage.
	VisitCreateMessage(ctx *CreateMessageContext) interface{}

	// Visit a parse tree produced by CELParser#IdentOrGlobalCall.
	VisitIdentOrGlobalCall(ctx *IdentOrGlobalCallContext) interface{}

	// Visit a parse tree produced by CELParser#DeprecatedIn.
	VisitDeprecatedIn(ctx *DeprecatedInContext) interface{}

	// Visit a parse tree produced by CELParser#Nested.
	VisitNested(ctx *NestedContext) interface{}

	// Visit a parse tree produced by CELParser#CreateList.
	VisitCreateList(ctx *CreateListContext) interface{}

	// Visit a parse tree produced by CELParser#CreateStruct.
	VisitCreateStruct(ctx *CreateStructContext) interface{}

	// Visit a parse tree produced by CELParser#ConstantLiteral.
	VisitConstantLiteral(ctx *ConstantLiteralContext) interface{}

	// Visit a parse tree produced by CELParser#exprList.
	VisitExprList(ctx *ExprListContext) interface{}

	// Visit a parse tree produced by CELParser#fieldInitializerList.
	VisitFieldInitializerList(ctx *FieldInitializerListContext) interface{}

	// Visit a parse tree produced by CELParser#mapInitializerList.
	VisitMapInitializerList(ctx *MapInitializerListContext) interface{}

	// Visit a parse tree produced by CELParser#Int.
	VisitInt(ctx *IntContext) interface{}

	// Visit a parse tree produced by CELParser#Uint.
	VisitUint(ctx *UintContext) interface{}

	// Visit a parse tree produced by CELParser#Double.
	VisitDouble(ctx *DoubleContext) interface{}

	// Visit a parse tree produced by CELParser#String.
	VisitString(ctx *StringContext) interface{}

	// Visit a parse tree produced by CELParser#Bytes.
	VisitBytes(ctx *BytesContext) interface{}

	// Visit a parse tree produced by CELParser#BoolTrue.
	VisitBoolTrue(ctx *BoolTrueContext) interface{}

	// Visit a parse tree produced by CELParser#BoolFalse.
	VisitBoolFalse(ctx *BoolFalseContext) interface{}

	// Visit a parse tree produced by CELParser#Null.
	VisitNull(ctx *NullContext) interface{}
}
