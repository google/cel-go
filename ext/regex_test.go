// Copyright 2025 Google LLC
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
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker"
)

func TestRegex(t *testing.T) {
	regexTests := []struct {
		expr string
	}{
		// Tests for replace Function
		{expr: "regex.replace('abc', '^', 'start_') == 'start_abc'"},
		{expr: "regex.replace('abc', '$', '_end') == 'abc_end'"},
		{expr: `regex.replace('a-b', r'\b', '|') == '|a|-|b|'`},
		{expr: `regex.replace('foo bar', '(fo)o (ba)r', r'\2 \1') == 'ba fo'`},
		{expr: `regex.replace('foo bar', 'foo', r'\\') == '\\ bar'`},
		{expr: "regex.replace('banana', 'ana', 'x') == 'bxna'"},
		{expr: `regex.replace('abc', 'b(.)', r'x\1') == 'axc'`},
		{expr: "regex.replace('hello world hello', 'hello', 'hi') == 'hi world hi'"},
		{expr: `regex.replace('ac', 'a(b)?c', r'[\1]') == '[]'`},
		{expr: "regex.replace('apple pie', 'p', 'X') == 'aXXle Xie'"},
		{expr: `regex.replace('remove all spaces', r'\s', '') == 'removeallspaces'`},
		{expr: `regex.replace('digit:99919291992', r'\d+', '3') == 'digit:3'`},
		{expr: `regex.replace('foo bar baz', r'\w+', r'(\0)') == '(foo) (bar) (baz)'`},
		{expr: "regex.replace('', 'a', 'b') == ''"},
		{expr: `regex.replace('User: Alice, Age: 30', r'User: (?P<name>\w+), Age: (?P<age>\d+)', '${name} is ${age} years old') == '${name} is ${age} years old'`},
		{expr: `regex.replace('User: Alice, Age: 30', r'User: (?P<name>\w+), Age: (?P<age>\d+)', r'\1 is \2 years old') == 'Alice is 30 years old'`},
		{expr: "regex.replace('hello ☃', '☃', '❄') == 'hello ❄'"},
		{expr: `regex.replace('id=123', r'id=(?P<value>\d+)', r'value: \1') == 'value: 123'`},
		{expr: "regex.replace('banana', 'a', 'x') == 'bxnxnx'"},
		{expr: `regex.replace(regex.replace('%(foo) %(bar) %2', r'%\((\w+)\)', r'${\1}'),r'%(\d+)', r'$\1') == '${foo} ${bar} $2'`},

		// Tests for replace Function with count variable
		{expr: "regex.replace('banana', 'a', 'x', 0) == 'banana'"},
		{expr: "regex.replace('banana', 'a', 'x', 1) == 'bxnana'"},
		{expr: "regex.replace('banana', 'a', 'x', 2) == 'bxnxna'"},
		{expr: "regex.replace('banana', 'a', 'x', 100) == 'bxnxnx'"},
		{expr: "regex.replace('banana', 'a', 'x', -1) == 'bxnxnx'"},
		{expr: "regex.replace('banana', 'a', 'x', -100) == 'bxnxnx'"},
		{expr: `regex.replace('cat-dog dog-cat cat-dog dog-cat', '(cat)-(dog)', r'\2-\1', 1) == 'dog-cat dog-cat cat-dog dog-cat'`},
		{expr: `regex.replace('cat-dog dog-cat cat-dog dog-cat', '(cat)-(dog)', r'\2-\1', 2) == 'dog-cat dog-cat dog-cat dog-cat'`},
		{expr: `regex.replace('a.b.c', r'\.', '-', 1) == 'a-b.c'`},
		{expr: `regex.replace('a.b.c', r'\.', '-', -1) == 'a-b-c'`},
		{expr: `regex.replace('abc def', r'(abc)', r'\\1') == r'\1 def'`},
		{expr: `regex.replace('abc def', r'(abc)', r'\\2') == r'\2 def'`},
		{expr: `regex.replace('abc def', r'(abc)', r'\\{word}') == '\\{word} def'`},
		{expr: `regex.replace('abc def', r'(abc)', r'\\word') == '\\word def'`},

		// Tests for extract Function
		{expr: "regex.extract('hello world', 'hello(.*)') == optional.of(' world')"},
		{expr: `regex.extract('item-A, item-B', r'item-(\w+)') == optional.of('A')`},
		{expr: `regex.extract('The color is red', r'The color is (\w+)') == optional.of('red')`},
		{expr: `regex.extract('The color is red', r'The color is \w+') == optional.of('The color is red')`},
		{expr: "regex.extract('brand', 'brand') == optional.of('brand')"},
		{expr: "regex.extract('hello world', 'goodbye (.*)') == optional.none()"},
		{expr: "regex.extract('HELLO', 'hello') == optional.none()"},
		{expr: `regex.extract('', r'\w+') == optional.none()`},
		{expr: "regex.extract('4122345432', '22').or(optional.of('777')) == optional.of('22')"},
		{expr: "regex.extract('4122345432', '22').orValue('777') == '22'"},

		// Tests for extractAll Function
		{expr: "regex.extractAll('id:123, id:456', 'assa') == []"},
		{expr: `regex.extractAll('id:123, id:456', r'id:\d+') == ['id:123', 'id:456']`},
		{expr: `regex.extractAll('Files: f_1.txt, f_2.csv', r'f_(\d+)') == ['1', '2']`},
		{expr: "regex.extractAll('testuser@', '(?P<username>.*)@') == ['testuser']"},
		{expr: "regex.extractAll('testuser@gmail.com, a@y.com, 2312321wsamkldjq2w2@sdad.com', '(?P<username>.*)@') == ['testuser@gmail.com, a@y.com, 2312321wsamkldjq2w2']"},
		{expr: `regex.extractAll('testuser@gmail.com, a@y.com, 2312321wsamkldjq2w2@sdad.com', r'(?P<username>\w+)@') == ['testuser', 'a', '2312321wsamkldjq2w2']`},
		{expr: "regex.extractAll('banananana', '(ana)') == ['ana', 'ana']"},
		{expr: `regex.extractAll('item:a1, topic:b2', r'(?:item:|topic:)([a-z]\d)') == ['a1', 'b2']`},
		{expr: "regex.extractAll('val=a, val=, val=c', 'val=([^,]*)') == ['a', 'c']"},
		{expr: "regex.extractAll('key=, key=, key=', 'key=([^,]*)') == []"},
		{expr: `regex.extractAll('a b c', r'(\S*)\s*') == ['a', 'b', 'c']`},
	}

	env := testRegexEnv(t)
	for i, tst := range regexTests {
		tr := tst // capture range variable
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tr.expr)
			if iss.Err() != nil {
				t.Fatalf("Parse(%s) failed: %v", tr.expr, iss.Err())
			}
			asts = append(asts, pAst)
			cAst, iss := env.Check(pAst)
			if iss.Err() != nil {
				t.Fatalf("Check(%s) failed: %v", tr.expr, iss.Err())
			}
			asts = append(asts, cAst)

			for _, ast := range asts {
				prog, err := env.Program(ast)
				if err != nil {
					t.Fatalf("Program(%s) failed: %v", tr.expr, err)
				}
				out, _, err := prog.Eval(cel.NoVars())
				if err != nil {
					t.Fatalf("Eval(%s) failed: %v", tr.expr, err)
				}
				if out.Value() != true {
					t.Errorf("Eval(%s) returned %v, want true", tr.expr, out.Value())
				}
			}
		})
	}
}

func TestRegexStaticErrors(t *testing.T) {
	tests := []struct {
		expr string
		err  string
	}{
		{
			expr: "regex.replace('abc', '^', 1)",
			err:  "found no matching overload for 'regex.replace' applied to '(string, string, int)'",
		},
		{
			expr: "regex.replace('abc', '^', '1','')",
			err:  "found no matching overload for 'regex.replace' applied to '(string, string, string, string)'",
		},
		{
			expr: "regex.extract('foo bar', 1)",
			err:  "found no matching overload for 'regex.extract' applied to '(string, int)'",
		},
		{
			expr: "regex.extract('foo bar', 1, 'bar')",
			err:  "found no matching overload for 'regex.extract' applied to '(string, int, string)'",
		},
		{
			expr: "regex.extractAll()",
			err:  "found no matching overload for 'regex.extractAll' applied to '()'",
		},
	}
	env := testRegexEnv(t)
	for i, tst := range tests {
		tr := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, iss := env.Compile(tr.expr)
			if iss.Err() == nil || !strings.Contains(iss.Err().Error(), tr.err) {
				t.Errorf("env.Compile(%q) got error %v, wanted %v", tr.expr, iss.Err(), tr.err)
			}
		})
	}
}

func TestRegexRuntimeErrors(t *testing.T) {
	tests := []struct {
		expr string
		err  string
	}{
		{
			expr: "regex.extract('foo bar', '(')",
			err:  "error parsing regexp: missing closing ): `(`",
		},
		{
			expr: "regex.extractAll('foo bar', '[a-z')",
			err:  "error parsing regexp: missing closing ]: `[a-z`",
		},
		{
			expr: `regex.replace('id=123', r'id=(?P<value>\d+)', r'value: \values')`,
			err:  `invalid replacement string: 'value: \values' \ must be followed by a digit`,
		},
		{
			expr: `regex.replace('test', '(.)', r'\2')`,
			err:  "replacement string references group 2 but regex has only 1 group(s)",
		},
		{
			expr: `regex.replace('id=123', r'id=(?P<value>\d+)', r'value: \')`,
			err:  `invalid replacement string: 'value: \' \ not allowed at end`,
		},
		{
			expr: `regex.extract('phone: 415-5551212', r'phone: ((\d{3})-)?')`,
			err:  `regular expression has more than one capturing group: "phone: ((\\d{3})-)?"`,
		},
		{
			expr: `regex.extractAll('Name: John Doe, Age:321', r'Name: (?P<Name>.*), Age:(?P<Age>\d+)')`,
			err:  `regular expression has more than one capturing group: "Name: (?P<Name>.*), Age:(?P<Age>\\d+)"`,
		},
		{
			expr: `regex.extractAll('testuser@testdomain', '(.*)@([^.]*)')`,
			err:  `regular expression has more than one capturing group: "(.*)@([^.]*)"`,
		},
		{
			expr: `regex.extractAll('The user testuser belongs to testdomain', 'The (user|domain) (?P<Username>.*) belongs (to) (?P<Domain>.*)')`,
			err:  `regular expression has more than one capturing group: "The (user|domain) (?P<Username>.*) belongs (to) (?P<Domain>.*)"`,
		},
	}

	env := testRegexEnv(t)
	for i, tst := range tests {
		tr := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			ast, iss := env.Compile(tr.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Compile(%q) failed with error %v", tr.expr, iss.Err())
			}
			prg, err := env.Program(ast)
			if err != nil {
				t.Fatalf("env.Program(ast) failed: %v", err)
			}
			in := cel.NoVars()
			_, _, err = prg.Eval(in)
			if err == nil || !strings.Contains(err.Error(), tr.err) {
				t.Errorf("prg.Eval() got %v, wanted %v", err, tr.err)
			}
		})
	}
}

func TestRegexEnvCreationErrors(t *testing.T) {
	tests := []struct {
		name string
		opts []cel.EnvOption
	}{
		{
			name: "no optional types",
			opts: []cel.EnvOption{Regex()},
		},
		{
			name: "optional types after regex",
			opts: []cel.EnvOption{Regex(), cel.OptionalTypes()},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := cel.NewEnv(tc.opts...)
			if err == nil || !strings.Contains(err.Error(), "regex library requires the optional library") {
				t.Fatalf("prg.Eval() got %v, wanted regex library requires the optional library", err)
			}
		})
	}
}

func TestRegexVersion(t *testing.T) {
	_, err := cel.NewEnv(cel.OptionalTypes(), Regex(RegexVersion(0)))
	if err != nil {
		t.Fatalf("Regex(0) failed: %v", err)
	}
}

func testRegexEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	baseOpts := []cel.EnvOption{
		cel.OptionalTypes(),
		Regex(),
	}
	env, err := cel.NewEnv(append(baseOpts, opts...)...)
	if err != nil {
		t.Fatalf("cel.NewEnv(Regex()) failed: %v", err)
	}
	return env
}

func TestRegexCosts(t *testing.T) {
	tests := []struct {
		expr          string
		vars          []cel.EnvOption
		in            map[string]any
		hints         map[string]uint64
		estimatedCost checker.CostEstimate
		actualCost    uint64
	}{
		{
			expr:          `regex.extract('hello world', 'hello (.*)') == optional.of('world')`,
			estimatedCost: checker.CostEstimate{Min: 19, Max: 20},
			actualCost:    6,
		},
		{
			expr:          "regex.extract('4122345432', '22').orValue('777') == '22'",
			estimatedCost: checker.CostEstimate{Min: 14, Max: 14},
			actualCost:    3,
		},
		{
			expr:          "regex.extract('4122345432', '22').or(optional.of('777')) == optional.of('22')",
			estimatedCost: checker.CostEstimate{Min: 16, Max: 1844674407370955279},
			actualCost:    4,
		},
		{
			expr:          "regex.extract('hello world', 'goodbye (.*)') == optional.none()",
			estimatedCost: checker.CostEstimate{Min: 21, Max: 22},
			actualCost:    7,
		},
		{
			expr:          "regex.extractAll('id:123, id:456', 'assa') == []",
			estimatedCost: checker.CostEstimate{Min: 24, Max: 38},
			actualCost:    12,
		},
		{
			expr: `regex.extractAll('id:123, id:456', r'id:\d+') == ['id:123', 'id:456']`,
			// Estimated Cost = search cost + a list allocation cost.
			// - List Allocation Cost: This is the base cost to create a list plus the worst-case
			//   size of its contents, which is the full target string's length (14).
			// - Search Cost: A multiplicative cost based on target (len 14) and regex (len 5) length.
			estimatedCost: checker.CostEstimate{Min: 25, Max: 39},
			// Actual Cost: The result is a list containing two strings: 'id:123' and 'id:456'.
			// The combined length of the contents is 6 + 6 = 12.
			// The cost is the content size plus overhead for the call, search, and list creation (~4),
			// so 12 + 4 = 16.
			actualCost: 16,
		},
		{
			expr:          `regex.extractAll('a b c', r'(\S*)\s*') == ['a', 'b', 'c']`,
			estimatedCost: checker.CostEstimate{Min: 24, Max: 29},
			actualCost:    16,
		},
		{
			expr:          `regex.extractAll('testuser@gmail.com, a@y.com, 2312321wsamkldjq2w2@sdad.com', r'(?P<username>\w+)@') == ['testuser', 'a', '2312321wsamkldjq2w2']`,
			estimatedCost: checker.CostEstimate{Min: 51, Max: 108},
			actualCost:    40,
		},
		{
			expr:          "regex.replace('hello world hello', 'hello', 'hi') == 'hi world hi'",
			estimatedCost: checker.CostEstimate{Min: 56, Max: 57},
			actualCost:    16,
		},
		{
			expr:          `regex.replace('ac', 'a(b)?c', r'[\1]') == '[]'`,
			estimatedCost: checker.CostEstimate{Min: 13, Max: 13},
			actualCost:    4,
		},
		{
			expr:          "regex.replace('apple pie', 'p', 'X') == 'aXXle Xie'",
			estimatedCost: checker.CostEstimate{Min: 20, Max: 20},
			actualCost:    11,
		},
		{
			expr: "regex.replace('aaaaaa', 'a', '-what-') == '-what--what--what--what--what--what-'",
			// Estimated Cost = search cost + build cost.
			// - Build Cost: pessimistic estimate: len(target) * (len(replacement) + 1) = 6 * (6 + 1) = 42
			// - Search Cost: A smaller multiplicative cost based on target and regex length.
			// The total is ~42 + search_cost, which aligns with the 44-47 range.
			estimatedCost: checker.CostEstimate{Min: 44, Max: 47},
			// Actual Cost: The result string '-what--what--...-' has a length of 36.
			// The actual cost is the length of the result plus a small, fixed overhead for the
			// function call and search (~5), so 36 + 5 = 41.
			actualCost: 41,
		},
		// --- Constant Cost Cases ---
		// The next three tests are interesting because the estimated and actual costs
		// remain constant, even though the number of replacements changes via the 'count' argument.
		//
		// Estimated Cost: Cost estimation is based on the static sizes of the target,
		// regex, and replacement strings, without considering the 'count' argument.
		// Since these inputs are identical for all three tests, the estimate is also identical.
		// It's the sum of a pessimistic Build Cost (~12) and a Search Cost (~2).
		//
		// Actual Cost: This is constant because the replacement 'x' (len 1) has the
		// same length as the matched pattern 'a' (len 1). Therefore, the final string's
		// length is always 6, regardless of how many replacements are made. The cost is
		// len(result) + base_overhead = 6 + 2 = 8.
		{
			expr:          "regex.replace('banana', 'a', 'x', 0) == 'banana'",
			estimatedCost: checker.CostEstimate{Min: 14, Max: 14},
			actualCost:    8,
		},
		{
			expr:          "regex.replace('banana', 'a', 'x', 1) == 'bxnana'",
			estimatedCost: checker.CostEstimate{Min: 14, Max: 14},
			actualCost:    8,
		},
		{
			expr:          "regex.replace('banana', 'a', 'x', 100) == 'bxnxnx'",
			estimatedCost: checker.CostEstimate{Min: 14, Max: 14},
			actualCost:    8,
		},
		{
			expr:          `regex.replace('foo bar', r'(foo bar)', r'\1\1\1\1\1' ) == 'foo barfoo barfoo barfoo barfoo bar'`,
			estimatedCost: checker.CostEstimate{Min: 81, Max: 84},
			actualCost:    41,
		},
		{
			expr:          `regex.replace('foo bar', r'(foo bar)', '') == ''`,
			estimatedCost: checker.CostEstimate{Min: 10, Max: 10},
			actualCost:    2,
		},
	}
	for _, test := range tests {
		tc := test
		t.Run(tc.expr, func(t *testing.T) {
			env := testRegexEnv(t, tc.vars...)
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("Parse(%s) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			cAst, iss := env.Check(pAst)
			if iss.Err() != nil {
				t.Fatalf("Check(%s) failed: %v", tc.expr, iss.Err())
			}

			testCheckCost(t, env, cAst, tc.hints, tc.estimatedCost)
			asts = append(asts, cAst)
			for _, ast := range asts {
				testEvalWithCost(t, env, ast, tc.in, tc.actualCost)
			}
		})
	}
}
