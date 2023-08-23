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

package cel_test

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/test/proto3pb"
)

func TestConstantFoldingOptimizer(t *testing.T) {
	tests := []struct {
		expr   string
		folded string
	}{
		{
			expr:   `[1, 1 + 2, 1 + (2 + 3)]`,
			folded: `[1, 3, 6]`,
		},
		{
			expr:   `6 in [1, 1 + 2, 1 + (2 + 3)]`,
			folded: `true`,
		},
		{
			expr:   `5 in [1, 1 + 2, 1 + (2 + 3)]`,
			folded: `false`,
		},
		{
			expr:   `x in [1, 1 + 2, 1 + (2 + 3)]`,
			folded: `x in [1, 3, 6]`,
		},
		{
			expr:   `1 in [1, x + 2, 1 + (2 + 3)]`,
			folded: `1 in [1, x + 2, 6]`,
		},
		{
			expr:   `{'hello': 'world'}.hello == x`,
			folded: `"world" == x`,
		},
		{
			expr:   `{'hello': 'world'}.?hello.orValue('default') == x`,
			folded: `"world" == x`,
		},
		{
			expr:   `{'hello': 'world'}['hello'] == x`,
			folded: `"world" == x`,
		},
		{
			expr:   `optional.of("hello")`,
			folded: `optional.of("hello")`,
		},
		{
			expr:   `optional.ofNonZeroValue("")`,
			folded: `optional.none()`,
		},
		{
			expr:   `{?'hello': optional.of('world')}['hello'] == x`,
			folded: `"world" == x`,
		},
		{
			expr:   `duration(string(7 * 24) + 'h')`,
			folded: `duration("604800s")`,
		},
		{
			expr:   `timestamp("1970-01-01T00:00:00Z")`,
			folded: `timestamp("1970-01-01T00:00:00Z")`,
		},
		{
			expr:   `[1, 1 + 1, 1 + 2, 2 + 3].exists(i, i < 10)`,
			folded: `true`,
		},
		{
			expr:   `[1, 1 + 1, 1 + 2, 2 + 3].exists(i, i < 1 % 2)`,
			folded: `false`,
		},
		{
			expr:   `[1, 2, 3].map(i, [1, 2, 3].map(j, i * j))`,
			folded: `[[1, 2, 3], [2, 4, 6], [3, 6, 9]]`,
		},
		{
			expr:   `[{}, {"a": 1}, {"b": 2}].filter(m, has(m.a))`,
			folded: `[{"a": 1}]`,
		},
		{
			expr:   `[{}, {"a": 1}, {"b": 2}].filter(m, has({'a': true}.a))`,
			folded: `[{}, {"a": 1}, {"b": 2}]`,
		},
		{
			expr:   `type(1)`,
			folded: `int`,
		},
		{
			expr:   `[google.expr.proto3.test.TestAllTypes{single_int32: 2 + 3}].map(i, i)[0]`,
			folded: `google.expr.proto3.test.TestAllTypes{single_int32: 5}`,
		},
		{
			expr:   `[1, ?optional.ofNonZeroValue(0)]`,
			folded: `[1]`,
		},
		{
			expr:   `[1, x, ?optional.ofNonZeroValue(3), ?x.?y]`,
			folded: `[1, x, 3, ?x.?y]`,
		},
		{
			expr:   `[1, x, ?optional.ofNonZeroValue(3), ?x.?y].size() > 3`,
			folded: `[1, x, 3, ?x.?y].size() > 3`,
		},
		{
			expr:   `{?'a': optional.of('hello'), ?x : optional.of(1), ?'b': optional.none()}`,
			folded: `{"a": "hello", ?x: optional.of(1)}`,
		},
		{
			expr:   `true ? x + 1 : x + 2`,
			folded: `x + 1`,
		},
		{
			expr:   `false ? x + 1 : x + 2`,
			folded: `x + 2`,
		},
		{
			expr:   `false ? x + 'world' : 'hello' + 'world'`,
			folded: `"helloworld"`,
		},
		{
			expr:   `null`,
			folded: `null`,
		},
		{
			expr:   `google.expr.proto3.test.TestAllTypes{?single_int32: optional.ofNonZeroValue(1)}`,
			folded: `google.expr.proto3.test.TestAllTypes{single_int32: 1}`,
		},
		{
			expr:   `google.expr.proto3.test.TestAllTypes{?single_int32: optional.ofNonZeroValue(0)}`,
			folded: `google.expr.proto3.test.TestAllTypes{}`,
		},
		{
			expr:   `google.expr.proto3.test.TestAllTypes{single_int32: x, repeated_int32: [1, 2, 3]}`,
			folded: `google.expr.proto3.test.TestAllTypes{single_int32: x, repeated_int32: [1, 2, 3]}`,
		},
		{
			expr:   `x + dyn([1, 2] + [3, 4])`,
			folded: `x + [1, 2, 3, 4]`,
		},
		{
			expr:   `dyn([1, 2]) + [3.0, 4.0]`,
			folded: `[1, 2, 3.0, 4.0]`,
		},
		{
			expr:   `{'a': dyn([1, 2]), 'b': x}`,
			folded: `{"a": [1, 2], "b": x}`,
		},
	}
	e, err := cel.NewEnv(
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.Types(&proto3pb.TestAllTypes{}),
		cel.Variable("x", cel.DynType))
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.expr, func(t *testing.T) {
			checked, iss := e.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("Compile() failed: %v", iss.Err())
			}
			opt := cel.NewStaticOptimizer(cel.NewConstantFoldingOptimizer())
			optimized, iss := opt.Optimize(e, checked)
			if iss.Err() != nil {
				t.Fatalf("Optimize() generated an invalid AST: %v", iss.Err())
			}
			folded, err := cel.AstToString(optimized)
			if err != nil {
				t.Fatalf("cel.AstToString() failed: %v", err)
			}
			if folded != tc.folded {
				t.Errorf("got %q, wanted %q", folded, tc.folded)
			}
		})
	}
}
