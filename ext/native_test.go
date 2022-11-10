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

package ext

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestNativeTypes(t *testing.T) {
	var nativeTests = []struct {
		expr string
		out  any
	}{
		{
			expr: `ext.TestAllTypes{BoolVal: true}`,
			out:  &TestAllTypes{BoolVal: true},
		},
	}
	env := testNativeEnv(t)
	for i, tst := range nativeTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
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
					t.Fatal(err)
				}
				out, _, err := prg.Eval(map[string]any{"msg": msgWithExtensions()})
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(out.Value(), tc.out) {
					t.Errorf("got %v, wanted %v for expr: %s", out.Value(), tc.out, tc.expr)
				}
			}
		})
	}
}

// testEnv initializes the test environment common to all tests.
func testNativeEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	envOpts := []cel.EnvOption{}
	envOpts = append(envOpts, opts...)
	envOpts = append(envOpts,
		NativeTypes(
			reflect.ValueOf(&NestedType{}),
			reflect.ValueOf(&TestAllTypes{}),
		),
	)
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		t.Fatalf("cel.NewEnv(NativeTypes()) failed: %v", err)
	}
	return env
}

type NestedType struct {
	NestedListVal []string
	NestedMapVal  map[int64]bool
}

type TestAllTypes struct {
	NestedVal *NestedType
	BoolVal   bool
	BytesVal  []byte
	DoubleVal float64
	IntVal    int32
	StringVal string
	UintVal   uint64
	ListVal   []*NestedType
	MapVal    map[string]*TestAllTypes
}
