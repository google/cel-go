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
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker"
)

func TestSets(t *testing.T) {
	setsTests := []struct {
		expr          string
		vars          []cel.EnvOption
		in            map[string]any
		hints         map[string]uint64
		estimatedCost checker.CostEstimate
		actualCost    uint64
	}{
		// set containment
		{
			expr:  `sets.contains(x, [1, 2, 3])`,
			vars:  []cel.EnvOption{cel.Variable("x", cel.ListType(cel.IntType))},
			in:    map[string]any{"x": []int64{5, 4, 3, 2, 1}},
			hints: map[string]uint64{"x": 10},
			// min cost is input 'x' length 0, 10 for list creation, 2 for arg costs
			// max cost is input 'x' lenght 10, 10 for list creation, 2 for arg costs
			estimatedCost: checker.CostEstimate{Min: 12, Max: 42},
			// actual cost is 'x' length 5 * list literal length 3, 10 for list creation, 2 for arg cost
			actualCost: 27,
		},
		{
			expr: `sets.contains(x, [1, 1, 1, 1, 1])`,
			vars: []cel.EnvOption{cel.Variable("x", cel.ListType(cel.IntType))},
			in:   map[string]any{"x": []int64{5, 4, 3, 2, 1}},
			// min cost is input 'x' length 0, 10 for list creation, 2 for arg costs
			// max cost is effectively infinite due to missing size hint for 'x'
			estimatedCost: checker.CostEstimate{Min: 12, Max: math.MaxUint64},
			// actual cost is 'x' length 5 * list literal length 5, 10 for list creation, 2 for arg cost
			actualCost: 37,
		},
		{
			expr:          `sets.contains([], [])`,
			estimatedCost: checker.CostEstimate{Min: 21, Max: 21},
			actualCost:    21,
		},
		{
			expr:          `sets.contains([1], [])`,
			estimatedCost: checker.CostEstimate{Min: 21, Max: 21},
			actualCost:    21,
		},
		{
			expr:          `sets.contains([1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 22, Max: 22},
			actualCost:    22,
		},
		{
			expr:          `sets.contains([1], [1, 1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.contains([1, 1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.contains([2, 1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.contains([1, 2, 3, 4], [2, 3])`,
			estimatedCost: checker.CostEstimate{Min: 29, Max: 29},
			actualCost:    29,
		},
		{
			expr:          `sets.contains([1], [1.0, 1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.contains([1, 2], [2u, 2.0])`,
			estimatedCost: checker.CostEstimate{Min: 25, Max: 25},
			actualCost:    25,
		},
		{
			expr:          `sets.contains([1, 2u], [2, 2.0])`,
			estimatedCost: checker.CostEstimate{Min: 25, Max: 25},
			actualCost:    25,
		},
		{
			expr:          `sets.contains([1, 2.0, 3u], [1.0, 2u, 3])`,
			estimatedCost: checker.CostEstimate{Min: 30, Max: 30},
			actualCost:    30,
		},
		{
			expr: `sets.contains([[1], [2, 3]], [[2, 3.0]])`,
			// 10 for each list creation, top-level list sizes are 2, 1
			estimatedCost: checker.CostEstimate{Min: 53, Max: 53},
			actualCost:    53,
		},
		{
			expr:          `!sets.contains([1], [2])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `!sets.contains([1], [1, 2])`,
			estimatedCost: checker.CostEstimate{Min: 24, Max: 24},
			actualCost:    24,
		},
		{
			expr:          `!sets.contains([1], ["1", 1])`,
			estimatedCost: checker.CostEstimate{Min: 24, Max: 24},
			actualCost:    24,
		},
		{
			expr:          `!sets.contains([1], [1.1, 1u])`,
			estimatedCost: checker.CostEstimate{Min: 24, Max: 24},
			actualCost:    24,
		},

		// set equivalence (note the cost factor is higher as it's basically two contains checks)
		{
			expr:          `sets.equivalent([], [])`,
			estimatedCost: checker.CostEstimate{Min: 21, Max: 21},
			actualCost:    21,
		},
		{
			expr:          `sets.equivalent([1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.equivalent([1], [1, 1])`,
			estimatedCost: checker.CostEstimate{Min: 25, Max: 25},
			actualCost:    25,
		},
		{
			expr:          `sets.equivalent([1, 1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 25, Max: 25},
			actualCost:    25,
		},
		{
			expr:          `sets.equivalent([1], [1u, 1.0])`,
			estimatedCost: checker.CostEstimate{Min: 25, Max: 25},
			actualCost:    25,
		},
		{
			expr:          `sets.equivalent([1], [1u, 1.0])`,
			estimatedCost: checker.CostEstimate{Min: 25, Max: 25},
			actualCost:    25,
		},
		{
			expr:          `sets.equivalent([1, 2, 3], [3u, 2.0, 1])`,
			estimatedCost: checker.CostEstimate{Min: 39, Max: 39},
			actualCost:    39,
		},
		{
			expr:          `sets.equivalent([[1.0], [2, 3]], [[1], [2, 3.0]])`,
			estimatedCost: checker.CostEstimate{Min: 69, Max: 69},
			actualCost:    69,
		},
		{
			expr:          `!sets.equivalent([2, 1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 26, Max: 26},
			actualCost:    26,
		},
		{
			expr:          `!sets.equivalent([1], [1, 2])`,
			estimatedCost: checker.CostEstimate{Min: 26, Max: 26},
			actualCost:    26,
		},
		{
			expr:          `!sets.equivalent([1, 2], [2u, 2, 2.0])`,
			estimatedCost: checker.CostEstimate{Min: 34, Max: 34},
			actualCost:    34,
		},
		{
			expr:          `!sets.equivalent([1, 2], [1u, 2, 2.3])`,
			estimatedCost: checker.CostEstimate{Min: 34, Max: 34},
			actualCost:    34,
		},

		// set intersection
		{
			expr:          `sets.intersects([1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 22, Max: 22},
			actualCost:    22,
		},
		{
			expr:          `sets.intersects([1], [1, 1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.intersects([1, 1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.intersects([2, 1], [1])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.intersects([1], [1, 2])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.intersects([1], [1.0, 2])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `sets.intersects([1, 2], [2u, 2, 2.0])`,
			estimatedCost: checker.CostEstimate{Min: 27, Max: 27},
			actualCost:    27,
		},
		{
			expr:          `sets.intersects([1, 2], [1u, 2, 2.3])`,
			estimatedCost: checker.CostEstimate{Min: 27, Max: 27},
			actualCost:    27,
		},
		{
			expr:          `sets.intersects([[1], [2, 3]], [[1, 2], [2, 3.0]])`,
			estimatedCost: checker.CostEstimate{Min: 65, Max: 65},
			actualCost:    65,
		},
		{
			expr:          `!sets.intersects([], [])`,
			estimatedCost: checker.CostEstimate{Min: 22, Max: 22},
			actualCost:    22,
		},
		{
			expr:          `!sets.intersects([1], [])`,
			estimatedCost: checker.CostEstimate{Min: 22, Max: 22},
			actualCost:    22,
		},
		{
			expr:          `!sets.intersects([1], [2])`,
			estimatedCost: checker.CostEstimate{Min: 23, Max: 23},
			actualCost:    23,
		},
		{
			expr:          `!sets.intersects([1], ["1", 2])`,
			estimatedCost: checker.CostEstimate{Min: 24, Max: 24},
			actualCost:    24,
		},
		{
			expr:          `!sets.intersects([1], [1.1, 2u])`,
			estimatedCost: checker.CostEstimate{Min: 24, Max: 24},
			actualCost:    24,
		},
	}

	for _, tst := range setsTests {
		tc := tst
		t.Run(tc.expr, func(t *testing.T) {
			env := testSetsEnv(t, tc.vars...)
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

			hints := map[string]uint64{}
			if len(tc.hints) != 0 {
				hints = tc.hints
			}
			est, err := env.EstimateCost(cAst, testSetsCostEstimator{hints: hints})
			if err != nil {
				t.Fatalf("env.EstimateCost() failed: %v", err)
			}
			if !reflect.DeepEqual(est, tc.estimatedCost) {
				t.Errorf("env.EstimateCost() got %v, wanted %v", est, tc.estimatedCost)
			}
			asts = append(asts, cAst)

			for _, ast := range asts {
				prgOpts := []cel.ProgramOption{}
				if ast.IsChecked() {
					prgOpts = append(prgOpts, cel.CostTracking(nil))
				}
				prg, err := env.Program(ast, prgOpts...)
				if err != nil {
					t.Fatalf("env.Program() failed: %v", err)
				}
				in := tc.in
				if in == nil {
					in = map[string]any{}
				}
				out, det, err := prg.Eval(in)
				if err != nil {
					t.Fatalf("prg.Eval() failed: %v", err)
				}
				if out.Value() != true {
					t.Errorf("prg.Eval() got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
				if det.ActualCost() != nil && *det.ActualCost() != tc.actualCost {
					t.Errorf("prg.Eval() had cost %v, wanted %v", *det.ActualCost(), tc.actualCost)
				}
			}
		})
	}
}

func testSetsEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	baseOpts := []cel.EnvOption{Sets()}
	env, err := cel.NewEnv(append(baseOpts, opts...)...)
	if err != nil {
		t.Fatalf("cel.NewEnv(Sets()) failed: %v", err)
	}
	return env
}

type testSetsCostEstimator struct {
	hints map[string]uint64
}

func (tc testSetsCostEstimator) EstimateSize(element checker.AstNode) *checker.SizeEstimate {
	if l, ok := tc.hints[strings.Join(element.Path(), ".")]; ok {
		return &checker.SizeEstimate{Min: 0, Max: l}
	}
	return nil
}

func (testSetsCostEstimator) EstimateCallCost(function, overloadID string, target *checker.AstNode, args []checker.AstNode) *checker.CallEstimate {
	return nil
}
