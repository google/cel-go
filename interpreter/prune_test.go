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

package interpreter

import (
	"github.com/google/cel-go/common/debug"
	operatorspb "github.com/google/cel-go/common/operators"
	testpb "github.com/google/cel-go/test"
	exprpb "github.com/google/cel-spec/proto/v1/syntax"
	"testing"
)

type testInfo struct {
	E *exprpb.Expr
	P string
}

var testCases = []testInfo{
	{
		E: testpb.ExprCall(2, operatorspb.LogicalAnd,
			testpb.ExprLiteral(1, true),
			testpb.ExprLiteral(3, false)),
		P: `false`,
	},
	{
		E: testpb.ExprCall(4, operatorspb.LogicalAnd,
			testpb.ExprCall(2, operatorspb.LogicalOr,
				testpb.ExprLiteral(1, true),
				testpb.ExprLiteral(3, false)),
			testpb.ExprIdent(5, "x")),
		P: `x`,
	},
	{
		E: testpb.ExprCall(4, operatorspb.LogicalAnd,
			testpb.ExprCall(2, operatorspb.LogicalOr,
				testpb.ExprLiteral(1, false),
				testpb.ExprLiteral(3, false)),
			testpb.ExprIdent(5, "x")),
		P: `false`,
	},
	{
		E: testpb.ExprCall(2, operatorspb.LogicalAnd,
			testpb.ExprIdent(1, "a"),
			testpb.ExprComprehension(3,
				"x",
				testpb.ExprList(7,
					testpb.ExprLiteral(4, int64(1)),
					testpb.ExprLiteral(5, uint64(1)),
					testpb.ExprLiteral(6, float64(1.0))),
				"_accu_",
				testpb.ExprLiteral(8, false),
				testpb.ExprCall(10,
					operatorspb.LogicalNot,
					testpb.ExprIdent(9, "_accu_")),
				testpb.ExprCall(13,
					operatorspb.Equals,
					testpb.ExprCall(12, "type", testpb.ExprIdent(11, "x")),
					testpb.ExprIdent(14, "uint")),
				testpb.ExprIdent(15, "_accu_"))),
		P: `a`,
	},
	{
		E: testpb.ExprMap(8,
			testpb.ExprEntry(2,
				testpb.ExprLiteral(1, "hello"),
				testpb.ExprMemberCall(3,
					"size",
					testpb.ExprLiteral(4, "world"))),
			testpb.ExprEntry(6,
				testpb.ExprLiteral(5, "bytes"),
				testpb.ExprLiteral(7, []byte("bytes-string")))),
		P: `{"hello":5, "bytes":b"bytes-string"}`,
	},
	{
		E: testpb.ExprCall(1, operatorspb.Less,
			testpb.ExprLiteral(2, int64(2)),
			testpb.ExprLiteral(3, int64(3))),
		P: `true`,
	},
	{
		E: testpb.ExprCall(8, operatorspb.Conditional,
			testpb.ExprLiteral(1, true),
			testpb.ExprCall(3,
				operatorspb.Less,
				testpb.ExprIdent(2, "b"),
				testpb.ExprLiteral(4, 1.2)),
			testpb.ExprCall(6,
				operatorspb.Equals,
				testpb.ExprIdent(5, "c"),
				testpb.ExprList(7, testpb.ExprLiteral(7, "hello")))),
		P: `_<_(b,1.2)`,
	},
}

func Test(t *testing.T) {
	for _, tst := range testCases {
		pExpr := &exprpb.ParsedExpr{Expr: tst.E}
		program := NewProgram(pExpr.Expr, pExpr.SourceInfo)
		interpretable := interpreter.NewInterpretable(program)
		_, state := interpretable.Eval(
			NewActivation(map[string]interface{}{}))
		newExpr := PruneAst(pExpr.Expr, state)
		actual := debug.ToDebugString(newExpr)
		if !testpb.Compare(actual, tst.P) {
			t.Fatal(testpb.DiffMessage("structure", actual, tst.P))
		}
	}
}
