// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func mustParseType(t testing.TB, name string) *exprpb.Type {
	t.Helper()
	ty, err := ParseType(name)
	if err != nil {
		t.Fatalf("ParseType(%s) failed", name)
	}
	return ty
}

func TestEvalSimple(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	_, _, err = eval.Evaluate("[1, 2, 3]")

	if err != nil {
		t.Errorf("eval.Evaluate('[1, 2, 3]') got %v, wanted non-error", err)
	}
}

func TestEvalSingleLetVar(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetVar("x", "2 + 2", nil)
	if err != nil {
		t.Errorf("eval.AddLetVar('x', '2 + 2') got %v, wanted non-error", err)
	}

	_, _, err = eval.Evaluate("[1, 2, 3, x]")

	if err != nil {
		t.Errorf("eval.Evaluate('[1, 2, 3, x]') got %v, wanted non-error", err)
	}
}

func TestEvalMultiLet(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	eval.AddLetVar("x", "20/5", nil)
	eval.AddLetVar("y", "x * 3", nil)
	eval.AddLetVar("x", "20", nil)

	_, _, err = eval.Evaluate("[1, 2, 3, x, y]")
	if err != nil {
		t.Errorf("eval.Evaluate('[1, 2, 3, x, y]') got %v, wanted non-error", err)
	}
}

func TestEvalError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	eval.AddLetVar("x", "1", nil)
	eval.AddLetVar("y", "0", nil)

	_, _, err = eval.Evaluate("x / y")
	if err == nil {
		t.Errorf("eval.Evaluate('x / y') got non-error, wanted division by zero")
	}
}

func TestLetError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetVar("y", "z + 1", nil)
	if err == nil {
		t.Errorf("eval.AddLetVar('y', 'z + 1') got %v, wanted error", err)
	}

	result, _, err := eval.Evaluate("y")
	if err == nil {
		t.Errorf("eval.Evaluate('y') got result %v, wanted error", result.Value())
	}
}

func TestLetTypeHintError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetVar("y", "10u", mustParseType(t, "int"))
	if err == nil {
		t.Errorf("eval.AddLetVar('y', '10u') got %v, wanted error", err)
	}

	result, _, err := eval.Evaluate("y")
	if err == nil {
		t.Errorf("eval.Evaluate('y') got result %v, wanted error", result)
	}
}

func TestDeclareError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	eval.AddDeclVar("z", mustParseType(t, "double"))
	eval.AddLetVar("y", "z + 10.0", nil)
	err = eval.AddLetVar("z", "'2.0'", nil)
	if err == nil {
		t.Errorf("eval.AddLetVar('z', '\"2.0\"') got %v, wanted error", err)
	}

	err = eval.AddDeclVar("z", mustParseType(t, "string"))
	if err == nil {
		t.Errorf("eval.AddDeclVar('z', string) got %v, wanted error", err)
	}
}

func TestDelError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	eval.AddLetVar("z", "41", nil)
	eval.AddLetVar("y", "z + 1", nil)

	err = eval.DelLetVar("z")
	if err == nil {
		t.Errorf("eval.DelLetVar('z') got %v, wanted error", err)
	}

	val, _, err := eval.Evaluate("y")
	if err != nil {
		t.Errorf("eval.Evaluate('y') failed %v, wanted non-error", err)
	}

	if val.Value().(int64) != 42 {
		t.Errorf("eval.Evaluate('y') got %v, wanted 42", val.Value())
	}
}
