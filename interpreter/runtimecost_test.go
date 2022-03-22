// Copyright 2022 Google LLC
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
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/parser"
	"github.com/google/cel-go/test/proto3pb"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestTrackCostAdvanced(t *testing.T) {
	var equalCases = []struct {
		in      interface{}
		lhsExpr string
		rhsExpr string
	}{
		{
			lhsExpr: `1`,
			rhsExpr: `2`,
		},
		{
			lhsExpr: `"abc".contains("d")`,
			rhsExpr: `"def".contains("d")`,
		},
		{
			lhsExpr: `1 in [4, 5, 6]`,
			rhsExpr: `2 in [15, 17, 16]`,
		},
	}
	for _, tc := range equalCases {
		t.Run(tc.lhsExpr+" vs "+tc.rhsExpr, func(t *testing.T) {
			ctx := constructActivation(t, tc.in)
			lhsCost, _, err := computeCost(t, tc.lhsExpr, nil, ctx, nil)
			if err != nil {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to eval expression due: %v", err)
			}
			rhsCost, _, err := computeCost(t, tc.rhsExpr, nil, ctx, nil)
			if err != nil {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to eval expression due: %v", err)
			}
			if lhsCost != rhsCost {
				t.Errorf(`Interpreter.Eval(activation Activation) failed return a cost for %s of %d equal to a cost for %s of %d`,
					tc.lhsExpr, lhsCost, tc.rhsExpr, rhsCost)
			}
		})

	}
	var smallerCases = []struct {
		in      interface{}
		lhsExpr string
		rhsExpr string
	}{
		{
			lhsExpr: `1`,
			rhsExpr: `1 + 2`,
		},
		{
			lhsExpr: `"abc".contains("d")`,
			rhsExpr: `"abcdhdflsfiehfieubdkwjbdwgxvuyagwsdwdnw qdbgquyidvbwqi".contains("e")`,
		},
		{
			lhsExpr: `1 in [4, 5, 6]`,
			rhsExpr: `1 in [4, 5, 6, 7, 8, 9]`,
		},
	}
	for _, tc := range smallerCases {
		t.Run(tc.lhsExpr+" vs "+tc.rhsExpr, func(t *testing.T) {
			ctx := constructActivation(t, tc.in)
			lhsCost, _, err := computeCost(t, tc.lhsExpr, nil, ctx, nil)
			if err != nil {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to eval expression due: %v", err)
			}
			rhsCost, _, err := computeCost(t, tc.rhsExpr, nil, ctx, nil)
			if err != nil {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to eval expression due: %v", err)
			}
			if lhsCost >= rhsCost {
				t.Errorf(`Interpreter.Eval(activation Activation) failed return a cost for %s of %d less than the cost for %s of %d`,
					tc.lhsExpr, lhsCost, tc.rhsExpr, rhsCost)
			}
		})
	}
}

func computeCost(t *testing.T, expr string, decls []*exprpb.Decl, ctx Activation, limit *uint64) (cost uint64, est checker.CostEstimate, err error) {
	t.Helper()

	s := common.NewTextSource(expr)
	p, err := parser.NewParser(parser.Macros(parser.AllMacros...))
	if err != nil {
		t.Fatalf("Failed to initialize parser: %v", err)
	}
	parsed, errs := p.Parse(s)
	if len(errs.GetErrors()) != 0 {
		t.Fatalf(`Failed to Parse expression "%s", error: %v`, expr, errs.GetErrors())
	}

	cont := containers.DefaultContainer
	reg := newTestRegistry(t, &proto3pb.TestAllTypes{})
	attrs := NewAttributeFactory(cont, reg, reg)
	env := newTestEnv(t, cont, reg)
	err = env.Add(decls...)
	if err != nil {
		t.Fatalf("Failed to initialize env: %v", err)
	}

	checked, errs := checker.Check(parsed, s, env)
	if len(errs.GetErrors()) != 0 {
		t.Fatalf(`Failed to check expression "%s", error: %v`, expr, errs.GetErrors())
	}
	est = checker.Cost(checked, testCostEstimator{})
	interp := NewStandardInterpreter(cont, reg, reg, attrs)
	costTracker := &CostTracker{Estimator: &testRuntimeCostEstimator{}, Limit: limit}
	prg, err := interp.NewInterpretable(checked, Observe(CostObserver(costTracker)))
	if err != nil {
		t.Fatalf(`Failed to check expression "%s", error: %v`, expr, errs.GetErrors())
	}

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case EvalCancelledError:
				err = t
			default:
				err = fmt.Errorf("internal error: %v", r)
			}
		}
	}()
	prg.Eval(ctx)
	// TODO: enable this once all attributes are properly pushed and popped from stack.
	//if len(costTracker.stack) != 1 {
	//	t.Fatalf(`Expected resulting stack size to be 1 but got %d: %#+v`, len(costTracker.stack), costTracker.stack)
	//}
	return costTracker.cost, est, err
}

func constructActivation(t *testing.T, in interface{}) Activation {
	t.Helper()
	if in == nil {
		return EmptyActivation()
	}
	a, err := NewActivation(in)
	if err != nil {
		t.Fatalf("NewActivation(%v) failed: %v", in, err)
	}
	return a
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randSeq(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

type testRuntimeCostEstimator struct {
}

var timeToYearCost uint64 = 7

func (e testRuntimeCostEstimator) CallCost(function, overloadID string, args []ref.Val, result ref.Val) *uint64 {
	argsSize := make([]uint64, len(args))
	for i, arg := range args {
		reflectV := reflect.ValueOf(arg.Value())
		switch reflectV.Kind() {
		// Note that the CEL bytes type is implemented with Go byte slices, therefore also supported by the following
		// code.
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			argsSize[i] = uint64(reflectV.Len())
		default:
			argsSize[i] = 1
		}
	}

	switch overloadID {
	case overloads.TimestampToYear:
		return &timeToYearCost
	default:
		return nil
	}
}

type testCostEstimator struct {
	hints map[string]int64
}

func (tc testCostEstimator) EstimateSize(element checker.AstNode) *checker.SizeEstimate {
	if l, ok := tc.hints[strings.Join(element.Path(), ".")]; ok {
		return &checker.SizeEstimate{Min: 0, Max: uint64(l)}
	}
	return nil
}

func (tc testCostEstimator) EstimateCallCost(function, overloadID string, target *checker.AstNode, args []checker.AstNode) *checker.CallEstimate {
	switch overloadID {
	case overloads.TimestampToYear:
		return &checker.CallEstimate{CostEstimate: checker.CostEstimate{Min: 7, Max: 7}}
	}
	return nil
}

func TestRuntimeCost(t *testing.T) {
	allTypes := decls.NewObjectType("google.expr.proto3.test.TestAllTypes")
	allList := decls.NewListType(allTypes)
	intList := decls.NewListType(decls.Int)
	nestedList := decls.NewListType(allList)

	allMap := decls.NewMapType(decls.String, allTypes)
	nestedMap := decls.NewMapType(decls.String, allMap)
	cases := []struct {
		name         string
		expr         string
		decls        []*exprpb.Decl
		want         uint64
		in           interface{}
		testFuncCost bool
		limit        uint64

		expectExceedsLimit bool
	}{
		{
			name: "const",
			expr: `"Hello World!"`,
			want: 0,
		},
		{
			name:  "identity",
			expr:  `input`,
			decls: []*exprpb.Decl{decls.NewVar("input", intList)},
			want:  1,
			in:    map[string]interface{}{"input": []int{1, 2}},
		},
		{
			name:  "select: map",
			expr:  `input['key']`,
			decls: []*exprpb.Decl{decls.NewVar("input", decls.NewMapType(decls.String, decls.String))},
			want:  2,
			in:    map[string]interface{}{"input": map[string]string{"key": "v"}},
		},
		{
			name:  "select: array index",
			expr:  `input[0]`,
			decls: []*exprpb.Decl{decls.NewVar("input", decls.NewListType(decls.String))},
			want:  2,
			in:    map[string]interface{}{"input": []string{"v"}},
		},
		{
			name:  "select: field",
			expr:  `input.single_int32`,
			decls: []*exprpb.Decl{decls.NewVar("input", allTypes)},
			want:  2,
			in: map[string]interface{}{
				"input": &proto3pb.TestAllTypes{
					RepeatedBool: []bool{false},
					MapInt64NestedType: map[int64]*proto3pb.NestedTestAllTypes{
						1: {},
					},
					MapStringString: map[string]string{},
				},
			},
		},
		{
			name:  "expr select: map",
			expr:  `input['ke' + 'y']`,
			decls: []*exprpb.Decl{decls.NewVar("input", decls.NewMapType(decls.String, decls.String))},
			want:  3,
			in:    map[string]interface{}{"input": map[string]string{"key": "v"}},
		},
		{
			name:  "expr select: array index",
			expr:  `input[3-3]`,
			decls: []*exprpb.Decl{decls.NewVar("input", decls.NewListType(decls.String))},
			want:  3,
			in:    map[string]interface{}{"input": []string{"v"}},
		},
		{
			name:  "select: field test only",
			expr:  `has(input.single_int32)`,
			decls: []*exprpb.Decl{decls.NewVar("input", decls.NewObjectType("google.expr.proto3.test.TestAllTypes"))},
			want:  0,
			in: map[string]interface{}{
				"input": &proto3pb.TestAllTypes{
					RepeatedBool: []bool{false},
					MapInt64NestedType: map[int64]*proto3pb.NestedTestAllTypes{
						1: {},
					},
					MapStringString: map[string]string{},
				},
			},
		},
		{
			name:         "estimated function call",
			expr:         `input.getFullYear()`,
			decls:        []*exprpb.Decl{decls.NewVar("input", decls.Timestamp)},
			want:         8,
			in:           map[string]interface{}{"input": time.Now()},
			testFuncCost: true,
		},
		{
			name: "create list",
			expr: `[1, 2, 3]`,
			want: 10,
		},
		{
			name: "create struct",
			expr: `google.expr.proto3.test.TestAllTypes{single_int32: 1, single_float: 3.14, single_string: 'str'}`,
			want: 40,
		},
		{
			name: "create map",
			expr: `{"a": 1, "b": 2, "c": 3}`,
			want: 30,
		},
		{
			name:  "all comprehension",
			decls: []*exprpb.Decl{decls.NewVar("input", allList)},
			expr:  `input.all(x, true)`,
			want:  2,
			in: map[string]interface{}{
				"input": []*proto3pb.TestAllTypes{},
			},
		},
		{
			name:  "nested all comprehension",
			decls: []*exprpb.Decl{decls.NewVar("input", nestedList)},
			expr:  `input.all(x, x.all(y, true))`,
			want:  2,
			in: map[string]interface{}{
				"input": []*proto3pb.TestAllTypes{},
			},
		},
		{
			name: "all comprehension on literal",
			expr: `[1, 2, 3].all(x, true)`,
			want: 20,
		},
		{
			name:  "variable cost function",
			decls: []*exprpb.Decl{decls.NewVar("input", decls.String)},
			expr:  `input.matches('[0-9]')`,
			want:  103,
			in:    map[string]interface{}{"input": string(randSeq(500))},
		},
		{
			name: "variable cost function with constant",
			expr: `'123'.matches('[0-9]')`,
			want: 2,
		},
		{
			name: "or",
			expr: `false || false`,
			want: 0,
		},
		{
			name: "or short-circuit",
			expr: `true || false`,
			want: 0,
		},
		{
			name: "and",
			expr: `true && false`,
			want: 0,
		},
		{
			name: "and short-circuit",
			expr: `false && true`,
			want: 0,
		},
		{
			name: "lt",
			expr: `1 < 2`,
			want: 1,
		},
		{
			name: "lte",
			expr: `1 <= 2`,
			want: 1,
		},
		{
			name: "eq",
			expr: `1 == 2`,
			want: 1,
		},
		{
			name: "gt",
			expr: `2 > 1`,
			want: 1,
		},
		{
			name: "gte",
			expr: `2 >= 1`,
			want: 1,
		},
		{
			name: "in",
			expr: `2 in [1, 2, 3]`,
			want: 13,
		},
		{
			name: "plus",
			expr: `1 + 1`,
			want: 1,
		},
		{
			name: "minus",
			expr: `1 - 1`,
			want: 1,
		},
		{
			name: "/",
			expr: `1 / 1`,
			want: 1,
		},
		{
			name: "/",
			expr: `1 * 1`,
			want: 1,
		},
		{
			name: "%",
			expr: `1 % 1`,
			want: 1,
		},
		{
			name: "ternary",
			expr: `true ? 1 : 2`,
			want: 0,
		},
		{
			name: "string size",
			expr: `size("123")`,
			want: 1,
		},
		{
			name: "str eq str",
			expr: `'12345678901234567890' == '123456789012345678901234567890'`,
			want: 2,
		},
		{
			name:  "bytes to string conversion",
			decls: []*exprpb.Decl{decls.NewVar("input", decls.Bytes)},
			expr:  `string(input)`,
			want:  51,
			in:    map[string]interface{}{"input": randSeq(500)},
		},
		{
			name:  "string to bytes conversion",
			decls: []*exprpb.Decl{decls.NewVar("input", decls.String)},
			expr:  `bytes(input)`,
			want:  51,
			in:    map[string]interface{}{"input": string(randSeq(500))},
		},
		{
			name: "int to string conversion",
			expr: `string(1)`,
			want: 1,
		},
		{
			name: "contains",
			expr: `input.contains(arg1)`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
				decls.NewVar("arg1", decls.String),
			},
			want: 2502,
			in:   map[string]interface{}{"input": string(randSeq(500)), "arg1": string(randSeq(500))},
		},
		{
			name: "matches",
			expr: `input.matches('\\d+a\\d+b')`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
			},
			want: 103,
			in:   map[string]interface{}{"input": string(randSeq(500)), "arg1": string(randSeq(500))},
		},
		{
			name: "startsWith",
			expr: `input.startsWith(arg1)`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
				decls.NewVar("arg1", decls.String),
			},
			want: 3,
			in:   map[string]interface{}{"input": "idc", "arg1": string(randSeq(500))},
		},
		{
			name: "endsWith",
			expr: `input.endsWith(arg1)`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
				decls.NewVar("arg1", decls.String),
			},
			want: 3,
			in:   map[string]interface{}{"input": "idc", "arg1": string(randSeq(500))},
		},
		{
			name: "size receiver",
			expr: `input.size()`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
			},
			want: 2,
			in:   map[string]interface{}{"input": "500", "arg1": "500"},
		},
		{
			name: "size",
			expr: `size(input)`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
			},
			want: 2,
			in:   map[string]interface{}{"input": "500", "arg1": "500"},
		},
		{
			name: "ternary eval",
			expr: `(x > 2 ? input1 : input2).all(y, true)`,
			decls: []*exprpb.Decl{
				decls.NewVar("x", decls.Int),
				decls.NewVar("input1", allList),
				decls.NewVar("input2", allList),
			},
			want: 6,
			in:   map[string]interface{}{"input1": []*proto3pb.TestAllTypes{{}}, "input2": []*proto3pb.TestAllTypes{{}}, "x": 1},
		},
		{
			name: "ternary eval trivial, true",
			expr: `true ? false : 1 > 3`,
			want: 0,
			in:   map[string]interface{}{},
		},
		{
			name: "ternary eval trivial, false",
			expr: `false ? false : 1 > 3`,
			want: 1,
			in:   map[string]interface{}{},
		},
		{
			name: "comprehension over map",
			expr: `input.all(k, input[k].single_int32 > 3)`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", allMap),
			},
			want: 9,
			in:   map[string]interface{}{"input": map[string]interface{}{"val": &proto3pb.TestAllTypes{}}},
		},
		{
			name: "comprehension over nested map of maps",
			expr: `input.all(k, input[k].all(x, true))`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", nestedMap),
			},
			want: 2,
			in:   map[string]interface{}{"input": map[string]interface{}{}},
		},
		{
			name: "string size of map keys",
			expr: `input.all(k, k.contains(k))`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", nestedMap),
			},
			want: 2,
			in:   map[string]interface{}{"input": map[string]interface{}{}},
		},
		{
			name: "comprehension variable shadowing",
			expr: `input.all(k, input[k].all(k, true) && k.contains(k))`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", nestedMap),
			},
			want: 2,
			in:   map[string]interface{}{"input": map[string]interface{}{}},
		},
		{
			name: "comprehension variable shadowing",
			expr: `input.all(k, input[k].all(k, true) && k.contains(k))`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", nestedMap),
			},
			want: 2,
			in:   map[string]interface{}{"input": map[string]interface{}{}},
		},
		{
			name: "list concat",
			expr: `(list1 + list2).all(x, true)`,
			decls: []*exprpb.Decl{
				decls.NewVar("list1", decls.NewListType(decls.Int)),
				decls.NewVar("list2", decls.NewListType(decls.Int)),
			},
			want: 4,
			in:   map[string]interface{}{"list1": []int{}, "list2": []int{}},
		},
		{
			name: "str concat",
			expr: `"abcdefg".contains(str1 + str2)`,
			decls: []*exprpb.Decl{
				decls.NewVar("str1", decls.String),
				decls.NewVar("str2", decls.String),
			},
			want: 6,
			in:   map[string]interface{}{"str1": "val1", "str2": "val2222222"},
		},
		{
			name: "at limit",
			expr: `"abcdefg".contains(str1 + str2)`,
			decls: []*exprpb.Decl{
				decls.NewVar("str1", decls.String),
				decls.NewVar("str2", decls.String),
			},
			in:    map[string]interface{}{"str1": "val1", "str2": "val2222222"},
			limit: 6,
			want:  6,
		},
		{
			name: "above limit",
			expr: `"abcdefg".contains(str1 + str2)`,
			decls: []*exprpb.Decl{
				decls.NewVar("str1", decls.String),
				decls.NewVar("str2", decls.String),
			},
			in:                 map[string]interface{}{"str1": "val1", "str2": "val2222222"},
			limit:              5,
			expectExceedsLimit: true,
		},
		{
			name:  "ternary as operand",
			expr:  `(1 > 2 ? 5 : 3) > 1`,
			decls: []*exprpb.Decl{},
			in:    map[string]interface{}{},
			want:  2,
		},
		{
			name:  "ternary as operand",
			expr:  `(1 > 2 || 2 > 1) == true`,
			decls: []*exprpb.Decl{},
			in:    map[string]interface{}{},
			want:  3,
		},
		{
			name:  "list map literal",
			expr:  `[{'k1': 1}, {'k2': 2}].all(x, true)`,
			decls: []*exprpb.Decl{},
			in:    map[string]interface{}{},
			want:  77,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := constructActivation(t, tc.in)

			var costLimit *uint64
			if tc.limit > 0 {
				costLimit = &tc.limit
			}
			actualCost, est, err := computeCost(t, tc.expr, tc.decls, ctx, costLimit)
			if err != nil {
				if tc.expectExceedsLimit {
					return
				}
				t.Fatalf("Interpreter.Eval(activation Activation) failed due to: %v", err)
			}
			if tc.expectExceedsLimit {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to return a cost exceeded error for limit %d, got cost %d", tc.limit, actualCost)
			}
			if actualCost != tc.want {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to return expected runtime cost %d, got %d", tc.want, actualCost)
			}
			if est.Min > actualCost || est.Max < actualCost {
				t.Fatalf("Interpreter.Eval(activation Activation) failed to return cost in range of estimate cost [%d, %d], got %d",
					est.Min, est.Max, actualCost)
			}
		})
	}
}
