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

package interpreter

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/test"
)

func TestAsyncEval_CallTracking(t *testing.T) {
	var tests = []struct {
		lhs         ref.Val
		lhsOverload string
		lhsTimeout  time.Duration
		rhs         ref.Val
		rhsOverload string
		rhsTimeout  time.Duration
		in          map[string]interface{}
		out         ref.Val
	}{
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  50 * time.Millisecond,
			lhsOverload: "same_async",
			rhs:         types.String("y success!"),
			rhsTimeout:  50 * time.Millisecond,
			rhsOverload: "same_async",
			in:          map[string]interface{}{"x": "x", "y": "y"},
			out:         types.True,
		},
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  50 * time.Millisecond,
			lhsOverload: "same_async",
			rhs:         types.String("y success!"),
			rhsTimeout:  1 * time.Millisecond,
			rhsOverload: "alt_async",
			in:          map[string]interface{}{"x": "x", "y": "y"},
			out:         types.True,
		},
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  1 * time.Millisecond,
			lhsOverload: "alt_async",
			rhs:         types.String("y success!"),
			rhsTimeout:  50 * time.Millisecond,
			rhsOverload: "same_async",
			in:          map[string]interface{}{"x": "x", "y": "y"},
			out:         types.True,
		},
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  1 * time.Millisecond,
			lhsOverload: "alt_async",
			rhs:         types.String("y success!"),
			rhsTimeout:  50 * time.Millisecond,
			rhsOverload: "same_async",
			in:          map[string]interface{}{"x": "x", "y": "not y"},
			out:         types.NewErr("context deadline exceeded"),
		},
		{
			lhs:        types.String("x success!"),
			rhs:        types.String("y success!"),
			lhsTimeout: 1 * time.Millisecond,
			rhsTimeout: 1 * time.Millisecond,
			in:         map[string]interface{}{"x": "x", "y": "y"},
			out:        types.NewErr("context deadline exceeded"),
		},
		{
			lhs:        types.Unknown{42},
			lhsTimeout: 50 * time.Millisecond,
			rhs:        types.String("y success!"),
			rhsTimeout: 50 * time.Millisecond,
			in:         map[string]interface{}{"x": "x", "y": "y"},
			out:        types.True,
		},
		{
			lhs:        types.Unknown{42},
			lhsTimeout: 50 * time.Millisecond,
			rhs:        types.String("y success!"),
			rhsTimeout: 50 * time.Millisecond,
			in:         map[string]interface{}{"x": "x", "y": "z"},
			out:        types.Unknown{42},
		},
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  50 * time.Millisecond,
			lhsOverload: "same_async",
			rhs:         types.String("x success!"),
			rhsTimeout:  50 * time.Millisecond,
			rhsOverload: "same_async",
			in:          map[string]interface{}{"x": "x", "y": "x"},
			out:         types.True,
		},
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  50 * time.Millisecond,
			lhsOverload: "same_async",
			rhs:         types.String(" success!"),
			rhsTimeout:  50 * time.Millisecond,
			rhsOverload: "same_async",
			in: map[string]interface{}{
				"x": types.Unknown{42},
				"y": "not y",
			},
			out: types.Unknown{42},
		},
		{
			lhs:         types.String("x success!"),
			lhsTimeout:  50 * time.Millisecond,
			lhsOverload: "same_async",
			rhs:         types.String("y success!"),
			rhsTimeout:  50 * time.Millisecond,
			rhsOverload: "same_async",
			in: map[string]interface{}{
				"x": types.Unknown{42},
				"y": "y",
			},
			out: types.True,
		},
	}
	xVar := &absoluteAttribute{
		id:             2,
		namespaceNames: []string{"x"},
	}
	yVar := &absoluteAttribute{
		id:             4,
		namespaceNames: []string{"y"},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(tt *testing.T) {
			testX := makeTestEq(6,
				tc.lhsOverload, xVar, NewConstValue(8, tc.lhs), tc.lhsTimeout)
			testY := makeTestEq(10,
				tc.rhsOverload, yVar, NewConstValue(12, tc.rhs), tc.rhsTimeout)
			logic := &evalOr{
				id:  14,
				lhs: testX,
				rhs: testY,
			}
			async := &asyncEval{Interpretable: logic}
			in, err := NewActivation(tc.in)
			if err != nil {
				tt.Fatal(err)
			}
			vars := NewAsyncActivation(in)
			ctx := context.TODO()
			out := async.AsyncEval(ctx, vars)
			if !reflect.DeepEqual(out, tc.out) {
				tt.Errorf("got %v, wanted %v", out, tc.out)
			}
			outCached := async.AsyncEval(ctx, vars)
			if !reflect.DeepEqual(out, outCached) {
				tt.Errorf("got %v, wanted %v", outCached, out)
			}
		})
	}
}

func makeTestEq(id int64,
	overload string,
	arg Attribute,
	value Interpretable,
	timeout time.Duration) Interpretable {
	test := &evalAsyncCall{
		id:       id,
		function: "asyncEcho",
		overload: overload,
		args: []Interpretable{
			&evalAttr{
				attr:    arg,
				adapter: types.DefaultTypeAdapter,
			},
		},
		impl: test.FakeRPC(timeout),
	}
	return &evalEq{
		id:  id + 1,
		lhs: test,
		rhs: value,
	}
}
