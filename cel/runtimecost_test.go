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

package cel

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/test/proto3pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

type testInfo struct {
	env     *Env
	in      interface{}
	lhsExpr string
	rhsExpr string
}

func computeCosts(t *testing.T, info *testInfo) (lhsCost, rhsCost uint64) {
	t.Helper()

	env := info.env
	if env == nil {
		emptyEnv, err := NewEnv()
		if err != nil {
			t.Fatalf("Failed to create empty environment, error: %v", err)
		}
		env = emptyEnv
	}
	ctx := constructActivation(t, info.in)
	lhsCost = computeCost(t, env, info.lhsExpr, &ctx)
	rhsCost = computeCost(t, env, info.rhsExpr, &ctx)

	return lhsCost, rhsCost
}

func TestTrackCostAdvanced(t *testing.T) {
	var equalCases = []testInfo{
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
	for i, testCase := range equalCases {
		lhsCost, rhsCost := computeCosts(t, &testCase)
		if lhsCost != rhsCost {
			t.Errorf(`Expected equal cost case #%d, expressions "%s" vs. "%s", respective cost %d vs. %d`, i,
				testCase.lhsExpr, testCase.rhsExpr, lhsCost, rhsCost)
		}
	}
	var smallerCases = []testInfo{
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
	for i, testCase := range smallerCases {
		lhsCost, rhsCost := computeCosts(t, &testCase)
		if lhsCost >= rhsCost {
			t.Errorf(`Expected smaller cost case #%d, expect the cost of expression "%s" to be strictly smaller than "%s", but got %d vs. %d respectively`,
				i, testCase.lhsExpr, testCase.rhsExpr, lhsCost, rhsCost)
		}
	}
}

func computeCost(t *testing.T, env *Env, expr string, ctx *interpreter.Activation) uint64 {
	t.Helper()

	// Compile and check expression.
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		t.Fatalf(`Failed to compile expression "%s", error: %v`, expr, iss.Err())
	}
	checked_ast, iss := env.Check(ast)
	if iss.Err() != nil {
		t.Fatalf(`Failed to check expression "%s", error: %v`, expr, iss.Err())
	}

	// Evaluate expression.
	program, err := env.Program(checked_ast, EvalOptions(OptTrackCost))
	if err != nil {
		t.Fatalf(`Failed to construct Program from expression "%s", error: %v`, expr, err)
	}
	_, details, err := program.Eval(*ctx)
	if err != nil {
		t.Fatalf(`Failed to evaluate expression "%s", error: %v`, expr, err)
	}
	costPtr := details.ActualCost()
	if costPtr == nil {
		t.Fatalf(`Null pointer returned for the cost of expression "%s"`, expr)
	}
	return *costPtr
}

func constructActivation(t *testing.T, in interface{}) interpreter.Activation {
	t.Helper()
	if in == nil {
		return interpreter.EmptyActivation()
	}
	a, err := interpreter.NewActivation(in)
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

func (e testRuntimeCostEstimator) CallCost(overloadId string, args []ref.Val) *uint64 {
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

	switch overloadId {
	case overloads.TimestampToYear:
		return &timeToYearCost
	default:
		return nil
	}
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
				"input": []proto3pb.TestAllTypes{},
			},
		},
		{
			name:  "nested all comprehension",
			decls: []*exprpb.Decl{decls.NewVar("input", nestedList)},
			expr:  `input.all(x, x.all(y, true))`,
			want:  2,
			in: map[string]interface{}{
				"input": []proto3pb.TestAllTypes{},
			},
		},
		{
			name: "all comprehension on literal",
			expr: `[1, 2, 3].all(x, true)`,
			want: 23,
		},
		{
			name:  "variable cost function",
			decls: []*exprpb.Decl{decls.NewVar("input", decls.String)},
			expr:  `input.matches('[0-9]')`,
			want:  101,
			in:    map[string]interface{}{"input": string(randSeq(500))},
		},
		{
			name: "variable cost function with constant",
			expr: `'123'.matches('[0-9]')`,
			want: 2,
		},
		{
			name: "or",
			expr: `true || false`,
			want: 1,
		},
		{
			name: "and",
			expr: `true && false`,
			want: 1,
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
			want: 1,
		},
		{
			name: "string size",
			expr: `size("123")`,
			want: 1,
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
			want: 101,
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
			want: 9,
			in:   map[string]interface{}{"input1": []proto3pb.TestAllTypes{proto3pb.TestAllTypes{}}, "input2": []proto3pb.TestAllTypes{proto3pb.TestAllTypes{}}, "x": 1},
		},
		{
			name: "comprehension over map",
			expr: `input.all(k, input[k].single_int32 > 3)`,
			decls: []*exprpb.Decl{
				decls.NewVar("input", allMap),
			},
			want: 11,
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
			want: 3,
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e, err := NewEnv(
				Declarations(tc.decls...),
				Types(&proto3pb.TestAllTypes{}),
				CustomTypeAdapter(types.DefaultTypeAdapter))
			if err != nil {
				t.Fatalf("environment creation error: %s\n", err)
			}
			ast, iss := e.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatal(iss.Err())
			}

			ctx := constructActivation(t, tc.in)
			checked_ast, iss := e.Check(ast)
			if iss.Err() != nil {
				t.Fatalf(`Failed to check expression with error: %v`, iss.Err())
			}
			// Evaluate expression.
			var program Program
			if tc.testFuncCost {
				program, err = e.Program(checked_ast, EvalOptions(OptTrackCost), CallCostEstimator(testRuntimeCostEstimator{}))
			} else {
				program, err = e.Program(checked_ast, EvalOptions(OptTrackCost))
			}

			if err != nil {
				t.Fatalf(`Failed to construct Program with error: %v`, err)
			}
			_, details, err := program.Eval(ctx)
			if err != nil {
				t.Fatalf(`Failed to evaluate expression with error: %v`, err)
			}
			actualCost := details.ActualCost()
			if actualCost == nil {
				t.Fatalf(`Null pointer returned for the cost of expression "%s"`, tc.expr)
			}
			if *actualCost != tc.want {
				t.Fatalf("runtime cost %d does not match expected runtime cost %d", *actualCost, tc.want)
			}

		})
	}
}
