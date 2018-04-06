// Copyright 2018 Google LLC
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
	testExpr "github.com/google/cel-go/interpreter/testing"
	"testing"
)

func TestNewProgram_Empty(t *testing.T) {
	program := NewProgram(
		testExpr.Empty.Expr,
		testExpr.Empty.Info(t.Name()),
		"")
	if loc, found := program.Metadata().Location(0); found {
		t.Errorf("Unexpected location found: %v", loc)
	}
	if program.Container() != "" {
		t.Errorf("Unexpected container name: %s", program.Container())
	}
	if step, hasNext := program.Instructions().Next(); hasNext {
		t.Errorf("Unexpected step in empty program: %v", step)
	}
}

func TestNewProgram_LogicalAnd(t *testing.T) {
	program := NewProgram(
		testExpr.LogicalAnd.Expr,
		testExpr.LogicalAnd.Info(t.Name()),
		"")
	if loc, found := program.Metadata().Location(1); found {
		t.Errorf("Unexpected location found: %v", loc)
	}
	if program.Container() != "" {
		t.Errorf("Unexpected container name: %s", program.Container())
	}
	if _, hasNext := program.Instructions().Next(); !hasNext {
		t.Error("Expected a step in program, but found none")
	}
	fmt.Printf("%s\n%s\n\n", t.Name(), program)
}

func TestNewProgram_Conditional(t *testing.T) {
	program := NewProgram(
		testExpr.Conditional.Expr,
		testExpr.Conditional.Info(t.Name()),
		"")
	if loc, found := program.Metadata().Location(1); found {
		t.Errorf("Unexpected location found: %v", loc)
	}
	if program.Container() != "" {
		t.Errorf("Unexpected container name: %s", program.Container())
	}
	if _, hasNext := program.Instructions().Next(); !hasNext {
		t.Error("Expected a step in program, but found none")
	}
	expected := "TestNewProgram_Conditional\n" +
		"0: local 'a', r1\n" +
		"1: jump  10 if r1 == undef\n" +
		"2: jump  5 if r1 == false\n" +
		"3: local 'b', r2\n" +
		"4: const 1, r4\n" +
		"5: call  _<_(r2, r4), r3\n" +
		"6: jump  4\n" +
		"7: local 'c', r5\n" +
		"8: const hello, r7\n" +
		"9: mov   list([7]), r8\n" +
		"10: call  _==_(r5, r8), r6\n" +
		"11: call  _?_:_(r1, r3, r6), r9\n\n"
	actual := fmt.Sprintf("%s\n%s\n\n", t.Name(), program)
	if actual != expected {
		t.Errorf("program did not compile as expected. actual: %v\nexpected: %v",
			actual, expected)
	}
}

func TestNewProgram_Comprehension(t *testing.T) {
	program := NewProgram(
		testExpr.Exists.Expr,
		testExpr.Exists.Info(t.Name()),
		"")
	if loc, found := program.Metadata().Location(1); !found {
		t.Errorf("Unexpected location found: %v", loc)
	}
	if program.Container() != "" {
		t.Errorf("Unexpected container name: %s", program.Container())
	}
	if _, hasNext := program.Instructions().Next(); !hasNext {
		t.Error("Expected a step in program, but found none")
	}
	fmt.Printf("%s\n%s\n\n", t.Name(), program)
}

func TestNewProgram_DynMap(t *testing.T) {
	program := NewProgram(
		testExpr.DynMap.Expr,
		testExpr.DynMap.Info(t.Name()),
		"")
	if loc, found := program.Metadata().Location(1); found {
		t.Errorf("Unexpected location found: %v", loc)
	}
	if program.Container() != "" {
		t.Errorf("Unexpected container name: %s", program.Container())
	}
	if _, hasNext := program.Instructions().Next(); !hasNext {
		t.Error("Expected a step in program, but found none")
	}
	fmt.Printf("%s\n%s\n\n", t.Name(), program)
}
