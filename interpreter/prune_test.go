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
	"fmt"
	"github.com/google/cel-go/common/debug"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/test"
	expr "github.com/google/cel-spec/proto/v1/syntax"
	"testing"
)

type testInfo struct {
	E *expr.Expr
	P string
}

var testCases = []testInfo{
	{
		E: test.ExprCall(2, operators.LogicalAnd,
			test.ExprLiteral(1, true),
			test.ExprLiteral(3, false)),
		P: `false`,
	},
	{
		E: test.ExprCall(4, operators.LogicalAnd,
			test.ExprCall(2, operators.LogicalOr,
				test.ExprLiteral(1, true),
				test.ExprLiteral(3, false)),
			test.ExprIdent(5, "x")),
		P: `_&&_(true, x)`,
	},
	{
		E: test.ExprCall(2, operators.LogicalAnd,
			test.ExprIdent(1, "a"),
			test.ExprComprehension(3,
				"x",
				test.ExprList(7,
					test.ExprLiteral(4, int64(1)),
					test.ExprLiteral(5, uint64(1)),
					test.ExprLiteral(6, float64(1.0))),
				"_accu_",
				test.ExprLiteral(8, false),
				test.ExprCall(10,
					operators.LogicalNot,
					test.ExprIdent(9, "_accu_")),
				test.ExprCall(13,
					operators.Equals,
					test.ExprCall(12, "type", test.ExprIdent(11, "x")),
					test.ExprIdent(14, "uint")),
				test.ExprIdent(15, "_accu_"))),
		P: `_&&_(a, true)`,
	},
}

func Test(t *testing.T) {
	for i, tst := range testCases {
		fmt.Printf("TestCase %d\n", i)
		ConstantFoldExpr(&expr.ParsedExpr{Expr: tst.E}, interpreter, "")
		actual := debug.ToDebugString(tst.E)
		if !test.Compare(actual, tst.P) {
			t.Fatal(test.DiffMessage("structure", actual, tst.P))
		}
	}
}
