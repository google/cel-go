// Copyright 2024 Google LLC
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
	"testing"

	"github.com/google/cel-go/cel"
)

func TestStringCostTracking(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		wantMin uint64
	}{
		{
			name:    "charAt",
			expr:    `"hello world".charAt(0)`,
			wantMin: 2,
		},
		{
			name:    "indexOf",
			expr:    `"hello world".indexOf("world")`,
			wantMin: 2,
		},
		{
			name:    "lastIndexOf",
			expr:    `"hello world".lastIndexOf("o")`,
			wantMin: 2,
		},
		{
			name:    "lowerAscii",
			expr:    `"HELLO".lowerAscii()`,
			wantMin: 2,
		},
		{
			name:    "upperAscii",
			expr:    `"hello".upperAscii()`,
			wantMin: 2,
		},
		{
			name:    "replace",
			expr:    `"hello world".replace("world", "CEL")`,
			wantMin: 2,
		},
		{
			name:    "replace_exponential_growth",
			expr:    `"A".replace("", "AAAAAAAAAA").replace("", "AAAAAAAAAA")`,
			wantMin: 10,
		},
		{
			name:    "split",
			expr:    `"a,b,c,d,e".split(",")`,
			wantMin: 2,
		},
		{
			name:    "substring",
			expr:    `"hello world".substring(0, 5)`,
			wantMin: 2,
		},
		{
			name:    "trim",
			expr:    `"  hello  ".trim()`,
			wantMin: 2,
		},
		{
			name:    "join",
			expr:    `["a", "b", "c", "d", "e"].join("-")`,
			wantMin: 2,
		},
		{
			name:    "reverse",
			expr:    `"hello".reverse()`,
			wantMin: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			env, err := cel.NewEnv(Strings(StringsVersion(3)))
			if err != nil {
				t.Fatalf("cel.NewEnv() failed: %v", err)
			}
			ast, iss := env.Compile(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Compile(%q) failed: %v", tc.expr, iss.Err())
			}
			prg, err := env.Program(ast, cel.CostTracking(nil))
			if err != nil {
				t.Fatalf("env.Program() failed: %v", err)
			}
			_, det, err := prg.Eval(cel.NoVars())
			if err != nil {
				t.Fatalf("prg.Eval() failed: %v", err)
			}
			cost := det.ActualCost()
			if cost == nil {
				t.Fatal("cost tracking returned nil")
			}
			if *cost < tc.wantMin {
				t.Errorf("cost for %q = %d, want at least %d", tc.expr, *cost, tc.wantMin)
			}
		})
	}
}

func TestStringCostLimitEnforced(t *testing.T) {
	env, err := cel.NewEnv(Strings())
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	// Chained replaces that produce exponential output.
	// Without cost tracking, this would be charged ~6. With tracking, the cost
	// scales with output size and should exceed any reasonable limit.
	expr := `"A".replace("", "AAAAAAAAAA").replace("", "AAAAAAAAAA").replace("", "AAAAAAAAAA").replace("", "AAAAAAAAAA").replace("", "AAAAAAAAAA").replace("", "AAAAAAAAAA")`
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		t.Fatalf("env.Compile() failed: %v", iss.Err())
	}
	prg, err := env.Program(ast, cel.CostLimit(1000), cel.CostTracking(nil))
	if err != nil {
		t.Fatalf("env.Program() failed: %v", err)
	}
	_, _, err = prg.Eval(cel.NoVars())
	if err == nil {
		t.Error("expected cost limit exceeded error for exponential string growth, got nil")
	}
}
