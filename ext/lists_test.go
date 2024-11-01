// Copyright 2023 Google LLC
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
	proto2pb "github.com/google/cel-go/test/proto2pb"
)

func TestLists(t *testing.T) {
	listsTests := []struct {
		expr string
		err  string
	}{
		{expr: `lists.range(4) == [0,1,2,3]`},
		{expr: `lists.range(0) == []`},
		{expr: `[5,1,2,3].reverse() == [3,2,1,5]`},
		{expr: `[].reverse() == []`},
		{expr: `[1].reverse() == [1]`},
		{expr: `['are', 'you', 'as', 'bored', 'as', 'I', 'am'].reverse() == ['am', 'I', 'as', 'bored', 'as', 'you', 'are']`},
		{expr: `[false, true, true].reverse().reverse() == [false, true, true]`},
		{expr: `[1,2,3,4].slice(0, 4) == [1,2,3,4]`},
		{expr: `[1,2,3,4].slice(0, 0) == []`},
		{expr: `[1,2,3,4].slice(1, 1) == []`},
		{expr: `[1,2,3,4].slice(4, 4) == []`},
		{expr: `[1,2,3,4].slice(1, 3) == [2, 3]`},
		{expr: `[1,2,3,4].slice(3, 0)`, err: "cannot slice(3, 0), start index must be less than or equal to end index"},
		{expr: `[1,2,3,4].slice(0, 10)`, err: "cannot slice(0, 10), list is length 4"},
		{expr: `[1,2,3,4].slice(-5, 10)`, err: "cannot slice(-5, 10), negative indexes not supported"},
		{expr: `[1,2,3,4].slice(-5, -3)`, err: "cannot slice(-5, -3), negative indexes not supported"},

		{expr: `dyn([]).flatten() == []`},
		{expr: `dyn([1,2,3,4]).flatten() == [1,2,3,4]`},
		{expr: `[1,[2,[3,4]]].flatten() == [1,2,[3,4]]`},
		{expr: `[1,2,[],[],[3,4]].flatten() == [1,2,3,4]`},
		{expr: `[1,[2,[3,4]]].flatten(2) == [1,2,3,4]`},
		{expr: `[1,[2,[3,[4]]]].flatten(-1) == [1,2,3,4]`, err: "level must be non-negative"},
		{expr: `[].sort() == []`},
		{expr: `[1].sort() == [1]`},
		{expr: `[4, 3, 2, 1].sort() == [1, 2, 3, 4]`},
		{expr: `["d", "a", "b", "c"].sort() == ["a", "b", "c", "d"]`},
		{expr: `["d", 3, 2, "c"].sort() == ["a", "b", "c", "d"]`, err: "list elements must have the same type"},
		{expr: `[].sortBy(e, e) == []`},
		{expr: `["a"].sortBy(e, e) == ["a"]`},
		{expr: `[-3, 1, -5, -2, 4].sortBy(e, -(e * e)) == [-5, 4, -3, -2, 1]`},
		{expr: `[-3, 1, -5, -2, 4].map(e, e * 2).sortBy(e, -(e * e)) == [-10, 8, -6, -4, 2]`},
		{expr: `lists.range(3).sortBy(e, -e) == [2, 1, 0]`},
		{expr: `["a", "c", "b", "first"].sortBy(e, e == "first" ? "" : e) == ["first", "a", "b", "c"]`},
		{expr: `[ExampleType{name: 'foo'}, ExampleType{name: 'bar'}, ExampleType{name: 'baz'}].sortBy(e, e.name) == [ExampleType{name: 'bar'}, ExampleType{name: 'baz'}, ExampleType{name: 'foo'}]`},
		{expr: `[].distinct() == []`},
		{expr: `[1].distinct() == [1]`},
		{expr: `[-2, 5, -2, 1, 1, 5, -2, 1].distinct() == [-2, 5, 1]`},
		{expr: `['c', 'a', 'a', 'b', 'a', 'b', 'c', 'c'].distinct() == ['c', 'a', 'b']`},
		{expr: `[1, 2.0, "c", 3, "c", 1].distinct() == [1, 2.0, "c", 3]`},
		{expr: `[1, 1.0, 2].distinct() == [1, 2]`},
		{expr: `[[1], [1], [2]].distinct() == [[1], [2]]`},
		{expr: `[ExampleType{name: 'a'}, ExampleType{name: 'b'}, ExampleType{name: 'a'}].distinct() == [ExampleType{name: 'a'}, ExampleType{name: 'b'}]`},
		{expr: `![].first().hasValue()`},
		{expr: `[1, 2, 3].first().value() == 1`},
		{expr: `![].last().hasValue()`},
		{expr: `[1, 2, 3].last().value() == 3`},
		{expr: `'/path/to'.split('/').filter(t, t.size() > 0).first().value() == 'path'`},
		{expr: `'/path/to'.split('/').filter(t, t.size() > 0).last().value() == 'to'`},
	}

	env := testListsEnv(t, cel.OptionalTypes(), Strings())
	for i, tst := range listsTests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			cAst, iss := env.Check(pAst)
			if iss.Err() != nil {
				t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, cAst)

			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatalf("env.Program() failed: %v", err)
				}
				out, _, err := prg.Eval(cel.NoVars())
				if tc.err != "" {
					if err == nil {
						t.Fatalf("got value %v, wanted error %s for expr: %s",
							out.Value(), tc.err, tc.expr)
					}
					if !strings.Contains(err.Error(), tc.err) {
						t.Errorf("got error %v, wanted error %s for expr: %s", err, tc.err, tc.expr)
					}
				} else if err != nil {
					t.Fatal(err)
				} else if out.Value() != true {
					t.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}

func testListsEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	baseOpts := []cel.EnvOption{
		Lists(),
		cel.Container("google.expr.proto2.test"),
		cel.Types(&proto2pb.ExampleType{},
			&proto2pb.ExternalMessageType{},
		)}
	env, err := cel.NewEnv(append(baseOpts, opts...)...)
	if err != nil {
		t.Fatalf("cel.NewEnv(Lists()) failed: %v", err)
	}
	return env
}
