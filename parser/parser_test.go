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

package parser

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/ast"
	"github.com/google/cel-go/test"
	"reflect"
)

var testCases = []testInfo{
	{
		I: `"A"`,
		P: `"A"^#1:*ast.StringConstant#`,
	},
	{
		I: `true`,
		P: `true^#1:*ast.BoolConstant#`,
	},
	{
		I: `false`,
		P: `false^#1:*ast.BoolConstant#`,
	},
	{
		I: `0`,
		P: `0^#1:*ast.Int64Constant#`,
	},
	{
		I: `42`,
		P: `42^#1:*ast.Int64Constant#`,
	},
	{
		I: `0u`,
		P: `0u^#1:*ast.Uint64Constant#`,
	},
	{
		I: `23u`,
		P: `23u^#1:*ast.Uint64Constant#`,
	},
	{
		I: `24u`,
		P: `24u^#1:*ast.Uint64Constant#`,
	},
	{
		I: `-1`,
		P: `-_(
    		  1^#1:*ast.Int64Constant#)^#2:*ast.CallExpression#`,
	},
	{
		I: `b"abc"`,
		P: `b"abc"^#1:*ast.BytesConstant#`,
	},
	{
		I: `23.39`,
		P: `23.39^#1:*ast.DoubleConstant#`,
	},
	{
		I: `!a`,
		P: `!_(
    		  a^#1:*ast.IdentExpression#)^#2:*ast.CallExpression#`,
	},
	{
		I: `null`,
		P: `null^#1:*ast.NullConstant#`,
	},
	{
		I: `a`,
		P: `a^#1:*ast.IdentExpression#`,
	},
	{
		I: `a?b:c`,
		P: `
			_?_:_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#,
    		  c^#3:*ast.IdentExpression#)^#4:*ast.CallExpression#`,
	},
	{
		I: `a || b`,
		P: `_||_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a + b`,
		P: `_+_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a - b`,
		P: `_-_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a * b`,
		P: `_*_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a / b`,
		P: `_/_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a % b`,
		P: `_%_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a in b`,
		P: `_in_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a == b`,
		P: `_==_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a != b`,
		P: `_!=_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a > b`,
		P: `_>_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a >= b`,
		P: `_>=_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a < b`,
		P: `_<_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a <= b`,
		P: `_<=_(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a.b`,
		P: `a^#1:*ast.IdentExpression#.b^#2:*ast.SelectExpression#`,
	},
	{
		I: `a.b.c`,
		P: `a^#1:*ast.IdentExpression#.b^#2:*ast.SelectExpression#.c^#3:*ast.SelectExpression#`,
	},
	{
		I: `a[b]`,
		P: `_[_](
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},

	// TODO: This is an error.
	//{
	//  I: `foo{ a=b, c=d }`,
	//  P: `
	//
	//  `,
	//},

	{
		I: `foo{ }`,
		P: `foo{}^#2:*ast.CreateMessageExpression#`,
	},
	{
		I: `foo{ a:b }`,
		P: `foo{
    		  a:b^#2:*ast.IdentExpression#^#3:*ast.FieldEntry#}^#4:*ast.CreateMessageExpression#`,
	},
	{
		I: `foo{ a:b, c:d }`,
		P: `foo{
    		  a:b^#2:*ast.IdentExpression#^#3:*ast.FieldEntry#,
    		  c:d^#4:*ast.IdentExpression#^#5:*ast.FieldEntry#}^#6:*ast.CreateMessageExpression#`,
	},
	{
		I: `{}`,
		P: `{}^#1:*ast.CreateStructExpression#`,
	},

	{
		I: `{a:b, c:d}`,
		P: `{
    		  a^#1:*ast.IdentExpression#:b^#2:*ast.IdentExpression#^#3:*ast.StructEntry#,
    		  c^#4:*ast.IdentExpression#:d^#5:*ast.IdentExpression#^#6:*ast.StructEntry#}^#7:*ast.CreateStructExpression#`,
	},
	{
		I: `[]`,
		P: `[]^#1:*ast.CreateListExpression#`,
	},
	{
		I: `[a]`,
		P: `[
    		  a^#2:*ast.IdentExpression#]^#1:*ast.CreateListExpression#`,
	},
	{
		I: `[a, b, c]`,
		P: `[
    		  a^#2:*ast.IdentExpression#,
    		  b^#3:*ast.IdentExpression#,
    		  c^#4:*ast.IdentExpression#]^#1:*ast.CreateListExpression#`,
	},
	{
		I: `(a)`,
		P: `a^#1:*ast.IdentExpression#`,
	},
	{
		I: `((a))`,
		P: `a^#1:*ast.IdentExpression#`,
	},
	{
		// Deprecated "in"
		I: `in(a,b)`,
		P: `in(
    		  a^#1:*ast.IdentExpression#,
    		  b^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a()`,
		P: `a()^#1:*ast.CallExpression#`,
	},

	{
		I: `a(b)`,
		P: `a(
    		  b^#1:*ast.IdentExpression#)^#2:*ast.CallExpression#`,
	},

	{
		I: `a(b, c)`,
		P: `a(
    		  b^#1:*ast.IdentExpression#,
    		  c^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
	},
	{
		I: `a.b()`,
		P: `a^#1:*ast.IdentExpression#.b()^#2:*ast.CallExpression#`,
	},

	{
		I: `a.b(c)`,
		P: `a^#1:*ast.IdentExpression#.b(
    		  c^#2:*ast.IdentExpression#)^#3:*ast.CallExpression#`,
		L: `a^#1[1,0]#.b(
    		  c^#2[1,4]#
    		)^#3[1,1]#`,
	},

	// Parse error tests
	{
		I: `*@a | b`,
		E: `
ERROR: <input>:1:2: Syntax error: token recognition error at: '@'
 | *@a | b
 | .^
ERROR: <input>:1:1: Syntax error: extraneous input '*' expecting {'in', '[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}
 | *@a | b
 | ^
ERROR: <input>:1:5: Syntax error: token recognition error at: '| '
 | *@a | b
 | ....^
ERROR: <input>:1:7: Syntax error: extraneous input 'b' expecting <EOF>
 | *@a | b
 | ......^`,
	},
	{
		I: `a | b`,
		E: `
ERROR: <input>:1:3: Syntax error: token recognition error at: '| '
 | a | b
 | ..^
ERROR: <input>:1:5: Syntax error: extraneous input 'b' expecting <EOF>
 | a | b
 | ....^`,
	},

	// Macro tests
	{
		I: `has(m.f)`,
		P: `m^#1:*ast.IdentExpression#.f~test-only~^#3:*ast.SelectExpression#`,
		L: `m^#1[1,4]#.f~test-only~^#3[1,3]#`,
	},
	{
		I: `m.exists_one(v, f)`,
		P: `__comprehension__(
    		  // Variable
    		  v,
    		  // Target
    		  m^#1:*ast.IdentExpression#,
    		  // Accumulator
    		  __result__,
    		  // Init
    		  0^#4:*ast.Int64Constant#,
    		  // LoopCondition
    		  _<=_(
    		    __result__^#7:*ast.IdentExpression#,
    		    1^#5:*ast.Int64Constant#)^#6:*ast.CallExpression#,
    		  // LoopStep
    		  _?_:_(
    		    f^#3:*ast.IdentExpression#,
    		    _+_(
    		      __result__^#10:*ast.IdentExpression#,
    		      1^#5:*ast.Int64Constant#)^#9:*ast.CallExpression#,
    		    __result__^#11:*ast.IdentExpression#)^#8:*ast.CallExpression#,
    		  // Result
    		  _==_(
    		    __result__^#13:*ast.IdentExpression#,
    		    1^#5:*ast.Int64Constant#)^#12:*ast.CallExpression#)^#14:*ast.ComprehensionExpression#`,
	},
	{
		I: `m.map(v, f)`,
		P: `__comprehension__(
    		  // Variable
    		  v,
    		  // Target
    		  m^#1:*ast.IdentExpression#,
    		  // Accumulator
    		  __result__,
    		  // Init
    		  []^#5:*ast.CreateListExpression#,
    		  // LoopCondition
    		  true^#6:*ast.BoolConstant#,
    		  // LoopStep
    		  _+_(
    		    __result__^#4:*ast.IdentExpression#,
    		    [
    		      f^#3:*ast.IdentExpression#]^#8:*ast.CreateListExpression#)^#7:*ast.CallExpression#,
    		    // Result
    		    __result__^#4:*ast.IdentExpression#)^#9:*ast.ComprehensionExpression#`,
	},

	{
		I: `m.map(v, p, f)`,
		P: `__comprehension__(
    		  // Variable
    		  v,
    		  // Target
    		  m^#1:*ast.IdentExpression#,
    		  // Accumulator
    		  __result__,
    		  // Init
    		  []^#6:*ast.CreateListExpression#,
    		  // LoopCondition
    		  true^#7:*ast.BoolConstant#,
    		  // LoopStep
    		  _?_:_(
    		    p^#3:*ast.IdentExpression#,
    		    _+_(
    		      __result__^#5:*ast.IdentExpression#,
    		      [
    		        f^#4:*ast.IdentExpression#]^#9:*ast.CreateListExpression#)^#8:*ast.CallExpression#,
    		      __result__^#5:*ast.IdentExpression#)^#10:*ast.CallExpression#,
    		    // Result
    		    __result__^#5:*ast.IdentExpression#)^#11:*ast.ComprehensionExpression#`,
	},

	{
		I: `m.filter(v, p)`,
		P: `__comprehension__(
    		  // Variable
    		  v,
    		  // Target
    		  m^#1:*ast.IdentExpression#,
    		  // Accumulator
    		  __result__,
    		  // Init
    		  []^#5:*ast.CreateListExpression#,
    		  // LoopCondition
    		  true^#6:*ast.BoolConstant#,
    		  // LoopStep
    		  _?_:_(
    		    p^#3:*ast.IdentExpression#,
    		    _+_(
    		      __result__^#4:*ast.IdentExpression#,
    		      [
    		        v^#2:*ast.IdentExpression#]^#8:*ast.CreateListExpression#)^#7:*ast.CallExpression#,
    		      __result__^#4:*ast.IdentExpression#)^#9:*ast.CallExpression#,
    		    // Result
    		    __result__^#4:*ast.IdentExpression#)^#10:*ast.ComprehensionExpression#`,
	},

	// Tests from Java parser
	{
		I: `[] + [1,2,3,] + [4]`,
		P: `_+_(
    		  _+_(
    		    []^#1:*ast.CreateListExpression#,
    		    [
    		      1^#3:*ast.Int64Constant#,
    		      2^#4:*ast.Int64Constant#,
    		      3^#5:*ast.Int64Constant#
    		    ]^#2:*ast.CreateListExpression#
    		  )^#6:*ast.CallExpression#,
    		  [
    		    4^#8:*ast.Int64Constant#
    		  ]^#7:*ast.CreateListExpression#
    		)^#9:*ast.CallExpression#`,
	},
	{
		I: `{1:2u, 2:3u}`,
		P: `{
    		  1^#1:*ast.Int64Constant#:2u^#2:*ast.Uint64Constant#^#3:*ast.StructEntry#,
    		  2^#4:*ast.Int64Constant#:3u^#5:*ast.Uint64Constant#^#6:*ast.StructEntry#
    		}^#7:*ast.CreateStructExpression#`,
	},
	{
		I: `TestAllTypes{single_int32: 1, single_int64: 2}`,
		P: `TestAllTypes{
    		  single_int32:1^#2:*ast.Int64Constant#^#3:*ast.FieldEntry#,
    		  single_int64:2^#4:*ast.Int64Constant#^#5:*ast.FieldEntry#
    		}^#6:*ast.CreateMessageExpression#`,
	},
	{
		I: `TestAllTypes(){single_int32: 1, single_int64: 2}`,
		E: `
ERROR: <input>:1:13: expected a qualified name
 | TestAllTypes(){single_int32: 1, single_int64: 2}
 | ............^
`,
	},
	{
		I: `size(x) == x.size()`,
		P: `_==_(
    		  size(
    		    x^#1:*ast.IdentExpression#
    		  )^#2:*ast.CallExpression#,
    		  x^#3:*ast.IdentExpression#.size()^#4:*ast.CallExpression#
    		)^#5:*ast.CallExpression#`,
	},
	{
		I: `1 + $`,
		E: `
ERROR: <input>:1:5: Syntax error: token recognition error at: '$'
 | 1 + $
 | ....^
ERROR: <input>:1:6: Syntax error: mismatched input '<EOF>' expecting {'in', '[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}
 | 1 + $
 | .....^
`,
	},
	{
		I: `1 + 2
3 +`,
		E: `
ERROR: <input>:2:1: Syntax error: mismatched input '3' expecting <EOF>
 | 3 +
 | ^
		`,
	},
	{
		I: `"\""`,
		P: `"\""^#1:*ast.StringConstant#`,
	},
	{
		I: `[1,3,4][0]`,
		P: `_[_](
    		  [
    		    1^#2:*ast.Int64Constant#,
    		    3^#3:*ast.Int64Constant#,
    		    4^#4:*ast.Int64Constant#
    		  ]^#1:*ast.CreateListExpression#,
    		  0^#5:*ast.Int64Constant#
    		)^#6:*ast.CallExpression#`,
	},
	{
		I: `1.all(2, 3)`,
		E: `
ERROR: <input>:1:7: The argument must be a simple name
 | 1.all(2, 3)
 | ......^
		`,
	},
	{
		I: `x["a"].single_int32 == 23`,
		P: `_==_(
    		  _[_](
    		    x^#1:*ast.IdentExpression#,
    		    "a"^#2:*ast.StringConstant#
    		  )^#3:*ast.CallExpression#.single_int32^#4:*ast.SelectExpression#,
    		  23^#5:*ast.Int64Constant#
    		)^#6:*ast.CallExpression#`,
	},
	{
		I: `x.single_nested_message != null`,
		P: `_!=_(
    		  x^#1:*ast.IdentExpression#.single_nested_message^#2:*ast.SelectExpression#,
    		  null^#3:*ast.NullConstant#
    		)^#4:*ast.CallExpression#`,
	},
	{
		I: `false && !true || false ? 2 : 3`,
		P: `_?_:_(
    		  _||_(
    		    _&&_(
    		      false^#1:*ast.BoolConstant#,
    		      !_(
    		        true^#2:*ast.BoolConstant#
    		      )^#3:*ast.CallExpression#
    		    )^#4:*ast.CallExpression#,
    		    false^#5:*ast.BoolConstant#
    		  )^#6:*ast.CallExpression#,
    		  2^#7:*ast.Int64Constant#,
    		  3^#8:*ast.Int64Constant#
    		)^#9:*ast.CallExpression#`,
	},
	{
		I: `b"abc" + B"def"`,
		P: `_+_(
    		  b"abc"^#1:*ast.BytesConstant#,
    		  b"def"^#2:*ast.BytesConstant#
    		)^#3:*ast.CallExpression#`,
	},
	{
		I: `1 + 2 * 3 - 1 / 2 == 6 % 1`,
		P: `_==_(
    		  _-_(
    		    _+_(
    		      1^#1:*ast.Int64Constant#,
    		      _*_(
    		        2^#2:*ast.Int64Constant#,
    		        3^#3:*ast.Int64Constant#
    		      )^#4:*ast.CallExpression#
    		    )^#5:*ast.CallExpression#,
    		    _/_(
    		      1^#6:*ast.Int64Constant#,
    		      2^#7:*ast.Int64Constant#
    		    )^#8:*ast.CallExpression#
    		  )^#9:*ast.CallExpression#,
    		  _%_(
    		    6^#10:*ast.Int64Constant#,
    		    1^#11:*ast.Int64Constant#
    		  )^#12:*ast.CallExpression#
    		)^#13:*ast.CallExpression#`,
	},
	{
		I: `1 + +`,
		E: `
ERROR: <input>:1:5: Syntax error: mismatched input '+' expecting {'in', '[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}
 | 1 + +
 | ....^
ERROR: <input>:1:6: Syntax error: mismatched input '<EOF>' expecting {'in', '[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}
 | 1 + +
 | .....^
		`,
	},
	{
		I: `"abc" + "def"`,
		P: `_+_(
    		  "abc"^#1:*ast.StringConstant#,
    		  "def"^#2:*ast.StringConstant#
    		)^#3:*ast.CallExpression#`,
	},
}

type testInfo struct {
	// I contains the input expression to be parsed.
	I string

	// P contains the type/id adorned debug output of the expression tree.
	P string

	// E contains the expected error output for a failed parse, or "" if the parse is expected to be successful.
	E string

	// L contains the expected source adorned debug output of the expression tree.
	L string
}

type kindAndIdAdorner struct {
}

func (k *kindAndIdAdorner) GetMetadata(e ast.Expression) string {
	return fmt.Sprintf("^#%d:%s#", e.Id(), reflect.TypeOf(e))
}

type locationAdorner struct {
}

func (k *locationAdorner) GetMetadata(e ast.Expression) string {
	return fmt.Sprintf("^#%d[%d,%d]#", e.Id(), e.Location().Line(), e.Location().Column())
}

func Test(t *testing.T) {
	for i, tst := range testCases {
		name := fmt.Sprintf("%d %s", i, tst.I)
		t.Run(name, func(tt *testing.T) {

			expression, errors := ParseText(tst.I)
			if len(errors.GetErrors()) > 0 {
				actualErr := errors.ToDisplayString()
				if tst.E == "" {
					tt.Fatalf("Unexpected errors: %v", actualErr)
				} else if !test.Compare(actualErr, tst.E) {
					tt.Fatalf(test.DiffMessage("Error mismatch", actualErr, tst.E))
				}
				return
			} else if tst.E != "" {
				tt.Fatalf("Expected error not thrown: '%s'", tst.E)
			}

			actualWithKind := ast.ToAdornedDebugString(expression, &kindAndIdAdorner{})
			if !test.Compare(actualWithKind, tst.P) {
				tt.Fatal(test.DiffMessage("structure", actualWithKind, tst.P))
			}

			if tst.L != "" {
				actualWithLocation := ast.ToAdornedDebugString(expression, &locationAdorner{})
				if !test.Compare(actualWithLocation, tst.L) {
					tt.Fatal(test.DiffMessage("location", actualWithLocation, tst.L))
				}
			}
		})
	}
}
