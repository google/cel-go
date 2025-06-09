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

// Package ext contains CEL extension libraries where each library defines a related set of
// constants, functions, macros, or other configuration settings which may not be covered by
// the core CEL spec.

package ext

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestRegex(t *testing.T) {
	regexTests := []struct {
		expr string
	}{
		// Tests for replace Function
		{expr: "regex.replace('foo bar', '(fo)o (ba)r', '$2 $1') == 'ba fo'"},
		{expr: "regex.replace('banana', 'ana', 'x') == 'bxna'"},
		{expr: "regex.replace('abc', 'b(.)', 'x$1') == 'axc'"},
		{expr: "regex.replace('hello world hello', 'hello', 'hi') == 'hi world hi'"},
		{expr: "regex.replace('apple pie', 'p', 'X') == 'aXXle Xie'"},
		{expr: "regex.replace('remove all spaces', '\\\\s', '') == 'removeallspaces'"},
		{expr: "regex.replace('digit:99919291992', '\\\\d+', '3') == 'digit:3'"},
		{expr: "regex.replace('foo bar baz', '\\\\w+', '($0)') == '(foo) (bar) (baz)'"},
		{expr: "regex.replace('', 'a', 'b') == ''"},
		{expr: "regex.replace('', 'a', 'b') == ''"},
		{expr: "regex.replace('banana', 'a', 'x') == 'bxnxnx'"},
		{expr: "regex.replace('banana', 'a', 'x', 0) == 'banana'"},
		{expr: "regex.replace('banana', 'a', 'x', 1) == 'bxnana'"},
		{expr: "regex.replace('banana', 'a', 'x', 2) == 'bxnxna'"},
		{expr: "regex.replace('banana', 'a', 'x', 100) == 'bxnxnx'"},
		{expr: "regex.replace('banana', 'a', 'x', -1) == 'bxnxnx'"},
		// {expr: "regex.replace('banana', 'a', 'x', -100) == 'banana'"},
		{expr: "regex.replace('cat-dog dog-cat cat-dog dog-cat', '(cat)-(dog)', '$2-$1', 1) == 'dog-cat dog-cat cat-dog dog-cat'"},
		{expr: "regex.replace('cat-dog dog-cat cat-dog dog-cat', '(cat)-(dog)', '$2-$1', 2) == 'dog-cat dog-cat dog-cat dog-cat'"},
		{expr: "regex.replace('a.b.c', '\\\\.', '-', 1) == 'a-b.c'"},
		{expr: "regex.replace('a.b.c', '\\\\.', '-', -1) == 'a-b-c'"},

		// Tests for capture Function
		{expr: "regex.capture('hello world', 'hello(.*)') == optional.of(' world')"},
		{expr: "regex.capture('item-A, item-B', 'item-(\\\\w+)') == optional.of('A')"},
		{expr: "regex.capture('The color is red', 'The color is (\\\\w+)') == optional.of('red')"},
		{expr: "regex.capture('The color is red', 'The color is \\\\w+') == optional.of('The color is red')"},
		{expr: "regex.capture('phone: 415-5551212', 'phone: ((\\\\d{3})-)?') == optional.of('415-')"},
		{expr: "regex.capture('brand', 'brand') == optional.of('brand')"},

		// Tests for captureAll Function
		{expr: "regex.captureAll('id:123, id:456', 'assa') == []"},
		{expr: "regex.captureAll('phone: 5551212', 'phone: ((\\\\d{3})-)?') == []"},
		{expr: "regex.captureAll('id:123, id:456', 'id:\\\\d+') == ['id:123', 'id:456']"},
		{expr: "regex.captureAll('testuser@', '(?P<username>.*)@') == ['testuser']"},
		{expr: "regex.captureAll('banananana', '(ana)') == ['ana', 'ana']"},
		{expr: "regex.captureAll('Name: John Doe, Age:321', 'Name: (?P<Name>.*), Age:(?P<Age>\\\\d+)')== ['John Doe', '321']"},
		{expr: "regex.captureAll('testuser@testdomain', '(.*)@([^.]*)') == ['testuser', 'testdomain']"},
		{expr: "regex.captureAll('The user testuser belongs to testdomain', 'The (user|domain) (?P<Username>.*) belongs (to) (?P<Domain>.*)') == ['user', 'testuser', 'to', 'testdomain']"},

		// Tests for captureAllNamed Function
		{expr: "regex.captureAllNamed('testuser@', '(?P<username>.*)@') == {'username': 'testuser'}"},
		{expr: "regex.captureAllNamed('Name: John Doe, Age:321', 'Name: (?P<Name>.*), Age:(?P<Age>\\\\d+)') == {'Name': 'John Doe', 'Age': '321'}"},
		{expr: "regex.captureAllNamed('id:123, id:456', 'assa') == {}"},
		{expr: "regex.captureAllNamed('id:123, id:456', 'id:\\\\d+') == {}"},
		{expr: "regex.captureAllNamed('testuser@testdomain', '(.*)@([^.]*)') == {}"},
		{expr: "regex.captureAllNamed('The user testuser belongs to testdomain', 'The (user|domain) (?P<Username>.*) belongs to (?P<Domain>.*)') == {'Username': 'testuser', 'Domain': 'testdomain'}"},
		{expr: "regex.captureAllNamed('', '(?P<name>\\\\w+)') == {}"},
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

func testRegexEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	baseOpts := []cel.EnvOption{
		Regex(),
		cel.OptionalTypes(),
	}
	env, err := cel.NewEnv(append(baseOpts, opts...)...)
	if err != nil {
		t.Fatalf("cel.NewEnv(Regex()) failed: %v", err)
	}
	return env
}
