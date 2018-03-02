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

// CELListener is a complete listener for a parse tree produced by CELParser.
type CELListener interface {
	antlr.ParseTreeListener

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

	// EnterStatementExpr is called when entering the StatementExpr production.
	EnterStatementExpr(c *StatementExprContext)

	// EnterLogicalNot is called when entering the LogicalNot production.
	EnterLogicalNot(c *LogicalNotContext)

	// EnterNegate is called when entering the Negate production.
	EnterNegate(c *NegateContext)

	// EnterSelectOrCall is called when entering the SelectOrCall production.
	EnterSelectOrCall(c *SelectOrCallContext)

	// EnterPrimaryExpr is called when entering the PrimaryExpr production.
	EnterPrimaryExpr(c *PrimaryExprContext)

	// EnterIndex is called when entering the Index production.
	EnterIndex(c *IndexContext)

	// EnterCreateMessage is called when entering the CreateMessage production.
	EnterCreateMessage(c *CreateMessageContext)

	// EnterIdentOrGlobalCall is called when entering the IdentOrGlobalCall production.
	EnterIdentOrGlobalCall(c *IdentOrGlobalCallContext)

	// EnterDeprecatedIn is called when entering the DeprecatedIn production.
	EnterDeprecatedIn(c *DeprecatedInContext)

	// EnterNested is called when entering the Nested production.
	EnterNested(c *NestedContext)

	// EnterCreateList is called when entering the CreateList production.
	EnterCreateList(c *CreateListContext)

	// EnterCreateStruct is called when entering the CreateStruct production.
	EnterCreateStruct(c *CreateStructContext)

	// EnterConstantLiteral is called when entering the ConstantLiteral production.
	EnterConstantLiteral(c *ConstantLiteralContext)

	// EnterExprList is called when entering the exprList production.
	EnterExprList(c *ExprListContext)

	// EnterFieldInitializerList is called when entering the fieldInitializerList production.
	EnterFieldInitializerList(c *FieldInitializerListContext)

	// EnterMapInitializerList is called when entering the mapInitializerList production.
	EnterMapInitializerList(c *MapInitializerListContext)

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

	// ExitStatementExpr is called when exiting the StatementExpr production.
	ExitStatementExpr(c *StatementExprContext)

	// ExitLogicalNot is called when exiting the LogicalNot production.
	ExitLogicalNot(c *LogicalNotContext)

	// ExitNegate is called when exiting the Negate production.
	ExitNegate(c *NegateContext)

	// ExitSelectOrCall is called when exiting the SelectOrCall production.
	ExitSelectOrCall(c *SelectOrCallContext)

	// ExitPrimaryExpr is called when exiting the PrimaryExpr production.
	ExitPrimaryExpr(c *PrimaryExprContext)

	// ExitIndex is called when exiting the Index production.
	ExitIndex(c *IndexContext)

	// ExitCreateMessage is called when exiting the CreateMessage production.
	ExitCreateMessage(c *CreateMessageContext)

	// ExitIdentOrGlobalCall is called when exiting the IdentOrGlobalCall production.
	ExitIdentOrGlobalCall(c *IdentOrGlobalCallContext)

	// ExitDeprecatedIn is called when exiting the DeprecatedIn production.
	ExitDeprecatedIn(c *DeprecatedInContext)

	// ExitNested is called when exiting the Nested production.
	ExitNested(c *NestedContext)

	// ExitCreateList is called when exiting the CreateList production.
	ExitCreateList(c *CreateListContext)

	// ExitCreateStruct is called when exiting the CreateStruct production.
	ExitCreateStruct(c *CreateStructContext)

	// ExitConstantLiteral is called when exiting the ConstantLiteral production.
	ExitConstantLiteral(c *ConstantLiteralContext)

	// ExitExprList is called when exiting the exprList production.
	ExitExprList(c *ExprListContext)

	// ExitFieldInitializerList is called when exiting the fieldInitializerList production.
	ExitFieldInitializerList(c *FieldInitializerListContext)

	// ExitMapInitializerList is called when exiting the mapInitializerList production.
	ExitMapInitializerList(c *MapInitializerListContext)

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
