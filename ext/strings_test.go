// Copyright 2020 Google LLC
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

package ext

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

var stringTests = []struct {
	expr      string
	err       string
	parseOnly bool
}{
	// Success test cases.
	{expr: `'hello'.after('l') == 'lo'`},
	{expr: `'hello'.after('l', 3) == 'o'`},
	{expr: `'hello'.after('none') == ''`},
	{expr: `'tacocat'.before('c') == 'ta'`},
	{expr: `'tacocat'.before('c', 3) == 'taco'`},
	{expr: `'tacocat'.before('none') == 'tacocat'`},
	{expr: `'tacocat'.charAt(3) == 'o'`},
	{expr: `'tacocat'.indexOf('a') == 1`},
	{expr: `'tacocat'.indexOf('a', 3) == 5`},
	{expr: `'tacocat'.indexOf('none') == -1`},
	{expr: `'MixedCase'.lower().after('mixed') == 'case'`},
	{expr: `'MixedCase'.upper().after('MIXED') == 'CASE'`},
	{expr: `"12 days 12 hours".replace("{0}", "2") == "12 days 12 hours"`},
	{expr: `"{0} days {0} hours".replace("{0}", "2") == "2 days 2 hours"`},
	{expr: `"{0} days {0} hours".replace("{0}", "2", 1).replace("{0}", "23") == "2 days 23 hours"`},
	{expr: `"hello world".split(" ") == ["hello", "world"]`},
	{expr: `"hello world events!".split(" ", 1) == ["hello world events!"]`},
	{expr: `"hello world events!".split(" ", 2) == ["hello", "world events!"]`},
	{expr: `"tacocat".substring(4) == "cat"`},
	{expr: `"tacocat".substring(0, 4) == "taco"`},
	{expr: `"tacocat".substring(4, 4) == ""`},
	{expr: `"   trim   ".trim() == "trim"`},
	// Error test cases.
	{
		expr: `'hello'.after('l', 30) == 'o'`,
		err:  "index out of range: 30",
	},
	{
		expr: `'tacocat'.before('c', 30) == 'taco'`,
		err:  "index out of range: 30",
	},
	{
		expr: `'tacocat'.charAt(30) == ''`,
		err:  "index out of range: 30",
	},
	{
		expr: `'tacocat'.indexOf('a', 30) == -1`,
		err:  "index out of range: 30",
	},
	{
		expr: `"tacocat".substring(40) == "cat"`,
		err:  "index out of range: 40",
	},
	{
		expr: `"tacocat".substring(-1) == "cat"`,
		err:  "index out of range: -1",
	},
	{
		expr: `"tacocat".substring(1, 50) == "cat"`,
		err:  "index out of range: 50",
	},
	{
		expr: `"tacocat".substring(49, 50) == "cat"`,
		err:  "index out of range: 49",
	},
	{
		expr: `"tacocat".substring(4, 3) == ""`,
		err:  "invalid substring range. start: 4, end: 3",
	},
	{
		expr:      `42.lower() == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `42.charAt(2) == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `'hello'.charAt(true) == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `'hello'.substring(1, 2, 3) == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `30.substring(true, 3) == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"tacocat".substring(true, 3) == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
	{
		expr:      `"tacocat".substring(0, false) == ""`,
		err:       "no such overload",
		parseOnly: true,
	},
}

func TestStrings(t *testing.T) {
	env, err := cel.NewEnv(Strings())
	if err != nil {
		t.Fatal(err)
	}
	for i, tst := range stringTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(tt *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				tt.Fatal(iss.Err())
			}
			asts = append(asts, pAst)
			if !tc.parseOnly {
				cAst, iss := env.Check(pAst)
				if iss.Err() != nil {
					tt.Fatal(iss.Err())
				}
				asts = append(asts, cAst)
			}
			for _, ast := range asts {
				exe, err := env.Program(ast)
				if err != nil {
					tt.Fatal(err)
				}
				out, _, err := exe.Eval(cel.NoVars())
				if tc.err != "" {
					if err == nil {
						tt.Fatalf("got value %v, wanted error %s for expr: %s", out.Value(), tc.err, tc.expr)
					}
					if tc.err != err.Error() {
						tt.Errorf("got error %v, wanted error %s for expr: %s", err, tc.err, tc.expr)
					}
				} else if err != nil {
					tt.Fatal(err)
				} else if out.Value() != true {
					tt.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}
