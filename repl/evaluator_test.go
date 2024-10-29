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

package repl

import (
	"strings"
	"testing"

	"github.com/google/cel-go/cel"

	proto2pb "github.com/google/cel-go/test/proto2pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var testTextDescriptorFile string = "testdata/attribute_context_fds.textproto"

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
	} else if val.Value().(int64) != 42 {
		t.Errorf("eval.Evaluate('y') got %v, wanted 42", val.Value())
	}

}

func TestAddLetFn(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	eval.AddLetFn("fn", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "int")},
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"),
		"x * x - y * y")

	eval.AddLetVar("testcases", "[[1, 2], [2, 3], [3, 4], [10, 20]]", mustParseType(t, "list(list(int))"))

	result, _, err := eval.Evaluate("testcases.all(e, fn(e[0], e[1]) == (e[0] - e[1]) * (e[0] + e[1]))")

	if err != nil {
		t.Errorf("eval.Evaluate() got error %v wanted nil", err)
	} else if !result.Value().(bool) {
		t.Errorf("eval.Evaluate() got %v wanted true", result.Value())
	}
}

func TestAddLetFnComposed(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("square", []letFunctionParam{{
		identifier: "x",
		typeHint:   mustParseType(t, "int"),
	}}, mustParseType(t, "int"), "x * x")

	if err != nil {
		t.Errorf("eval.AddLetFn(square x -> int) got error %v expected nil", err)
	}

	err = eval.AddLetFn("squareDiff", []letFunctionParam{
		{
			identifier: "x",
			typeHint:   mustParseType(t, "int"),
		},
		{
			identifier: "y",
			typeHint:   mustParseType(t, "int"),
		}}, mustParseType(t, "int"), "square(x) - square(y)")

	if err != nil {
		t.Errorf("eval.AddLetFn(squareDiff x, y -> int) got error %v expected nil", err)
	}

	result, _, err := eval.Evaluate("squareDiff(4, 3)")

	if err != nil {
		t.Errorf("eval.Evaluate() got error %v wanted nil", err)
	} else if result.Value().(int64) != 7 {
		t.Errorf("eval.Evaluate() got %v wanted true", result.Value())
	}
}

func TestAddLetFnErrorOnTypeChange(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("square", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"),
		"x * x")

	if err != nil {
		t.Errorf("eval.AddLetFn('square x -> int') got error %v expected nil", err)
	}

	eval.AddLetVar("y", "square(1) + 1", nil)

	// Overloads not yet supported
	err = eval.AddLetFn("square", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "double")}},
		mustParseType(t, "double"),
		"x * x")

	if err == nil {
		t.Error("eval.AddLetFn('square x -> double') got nil, expected error")
	}

	result, _, err := eval.Evaluate("y")

	if err != nil {
		t.Errorf("eval.Evaluate() got error %v wanted nil", err)
	} else if result.Value().(int64) != 2 {
		t.Errorf("eval.Evaluate() got %v wanted true", result.Value())
	}
}

func TestAddLetFnErrorTypeMismatch(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("fn", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "int")},
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "double"),
		"x * x - y * y")

	if err == nil {
		t.Error("eval.AddLetFn('fn x : int, y : int -> double') got nil expected error")
	}
}

func TestAddLetFnErrorBadExpr(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("fn", []letFunctionParam{
		{identifier: "y", typeHint: mustParseType(t, "string")}},
		mustParseType(t, "int"),
		"2 - y")

	if err == nil {
		t.Error("eval.AddLetFn('fn y : string -> int = 2 - y') got nil wanted error")
	}
}

func TestAddDeclFn(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddDeclFn("fn", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "int")},
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"))

	if err != nil {
		t.Errorf("eval.AddDeclFn('fn(x:int, y:int): int') got error %v wanted nil", err)
	}

	err = eval.AddLetVar("z", "fn(1, 2)", nil)

	if err != nil {
		t.Errorf("eval.AddLetVar('z = fn(1, 2)') got error %v wanted nil", err)
	}
}

func TestAddDeclFnError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddDeclFn("fn", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "int")},
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"))

	if err != nil {
		t.Errorf("eval.AddDeclFn('fn(x:int, y:int): int') got error %v wanted nil", err)
	}

	err = eval.AddLetVar("z", "fn(1, 2)", nil)

	if err != nil {
		t.Errorf("eval.AddLetVar('z = fn(1, 2)') got error %v wanted nil", err)
	}

	err = eval.AddDeclFn("fn", []letFunctionParam{
		{identifier: "x", typeHint: mustParseType(t, "string")},
		{identifier: "y", typeHint: mustParseType(t, "string")}},
		mustParseType(t, "int"))

	if err == nil {
		t.Errorf("eval.AddDeclFn('fn(x:string, y:string): int') got nil wanted error")
	}

}

func TestDelLetFn(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("fn", []letFunctionParam{
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"),
		"2 - y")

	if err != nil {
		t.Fatalf("eval.AddLetFn('fn (y : string): int -> 2 - y') failed: %s", err)
	}

	err = eval.DelLetFn("fn")
	if err != nil {
		t.Errorf("eval.DelLetFn('fn') got error %s, wanted nil", err)
	}

	err = eval.AddLetVar("x", "fn(2)", nil)
	if err == nil {
		t.Error("eval.AddLetVar('fn(2)') got nil, wanted error")
	}
}

func TestDelLetFnError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("fn", []letFunctionParam{
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"),
		"2 - y")

	if err != nil {
		t.Fatalf("eval.AddLetFn('fn (y : string): int -> 2 - y') failed: %s", err)
	}

	err = eval.AddLetVar("x", "fn(2)", nil)
	if err != nil {
		t.Errorf("eval.AddLetVar('fn(2)') got error %s, wanted nil", err)
	}

	err = eval.DelLetFn("fn")
	if err == nil {
		t.Error("eval.DelLetFn('fn') got nil, wanted error")
	}
}

func TestInstanceFunction(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddLetFn("int.plus", []letFunctionParam{
		{identifier: "y", typeHint: mustParseType(t, "int")}},
		mustParseType(t, "int"),
		"this + y")

	if err != nil {
		t.Fatalf("eval.AddLetFn('int.plus(y : int): int -> this + y') failed: %s", err)
	}

	val, _, err := eval.Evaluate("40.plus(2)")
	if err != nil {
		t.Errorf("eval.Eval('40.plus(2)') got error %v, wanted nil", err)
	}
	if val.Value().(int64) != 42 {
		t.Errorf("eval.Eval('40.plus(2)') got %s, wanted 42", val)
	}
}

type testContainerOption struct {
	container string
}

func (o *testContainerOption) String() string {
	return o.container
}

func (o *testContainerOption) Option() cel.EnvOption {
	return cel.Container(o.container)
}

func TestSetOption(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddOption(&testContainerOption{container: "google.protobuf"})
	if err != nil {
		t.Errorf("eval.AddOption() got error: %v, wanted nil", err)
	}

	val, _, err := eval.Evaluate("Int64Value{value: 42} == 42")
	if err != nil {
		t.Errorf("eval.Evaluate() got error: %v, expected nil", err)
	}

	if !val.Value().(bool) {
		t.Error("eval.Evaluate() got false expected true")
	}
}

func TestSetOptionError(t *testing.T) {
	eval, err := NewEvaluator()
	if err != nil {
		t.Fatalf("NewEvaluator() failed with: %v", err)
	}

	err = eval.AddOption(&testContainerOption{container: "google.protobuf"})
	if err != nil {
		t.Errorf("eval.AddOption() got error: %v, wanted nil", err)
	}

	err = eval.AddLetVar("x", "Int64Value{value: 42} == 42", nil)
	if err != nil {
		t.Errorf("eval.Evaluate() got error: %v, expected nil", err)
	}

	err = eval.AddOption(&testContainerOption{container: ""})
	if err == nil {
		t.Error("eval.AddOption() got nil expected error")
	}
}

func TestProcess(t *testing.T) {
	var testCases = []struct {
		name             string
		commands         []Cmder
		wantText         string
		wantExit         bool
		wantError        bool
		ignoreWhitespace bool
	}{
		{
			name: "CompileResult",
			commands: []Cmder{
				&compileCmd{
					expr: "3u",
				},
			},
			wantText: `type_map:  {
				key:  1
				value:  {
				  primitive:  UINT64
				}
			  }
			  source_info:  {
				location:  "<input>"
				line_offsets:  3
				positions:  {
				  key:  1
				  value:  0
				}
			  }
			  expr:  {
				id:  1
				const_expr:  {
				  uint64_value:  3
				}
			  }`,
			wantExit:         false,
			wantError:        false,
			ignoreWhitespace: true,
		},
		{
			name: "FormatNumberResult",
			commands: []Cmder{
				&evalCmd{
					expr: "1u + 2u",
				},
			},
			wantText:  "3 : uint",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "FormatStringResult",
			commands: []Cmder{
				&evalCmd{
					expr: `'a' + r'b\1'`,
				},
			},
			wantText:  `ab\1 : string`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "FormatListResult",
			commands: []Cmder{
				&evalCmd{
					expr: `['abc', 123, 3.14, duration('2m')]`,
				},
			},
			wantText:  `["abc", 123, 3.140000, duration("120s")] : list(dyn)`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "FormatMapResult",
			commands: []Cmder{
				&evalCmd{
					expr: `{1: 123, 2: 3.14, 3: duration('2m'), 4: b'123'}`,
				},
			},
			wantText:  `{1:123, 2:3.140000, 3:duration("120s"), 4:b"123"} : map(int, dyn)`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionBasic",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"google.protobuf",
					},
				},
			},
			wantText:  "",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionContainer",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"google.protobuf",
					},
				},
				&evalCmd{
					expr: "Int64Value{value: 20}",
				},
			},
			wantText:  "20 : wrapper(int)",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionExtensionOptional",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--extension",
						"optional",
					},
				},
				&evalCmd{
					expr: "optional.none().orValue('default')",
				},
			},
			wantText:  "default : string",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionExtensionStrings",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--extension",
						"strings",
					},
				},
				&evalCmd{
					expr: "'test'.substring(2)",
				},
			},
			wantText:  "st : string",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionExtensionProtos",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--extension",
						"protos",
					},
				},
				&evalCmd{
					expr: "proto.getExt(google.expr.proto2.test.ExampleType{}, google.expr.proto2.test.int32_ext) == 0",
				},
			},
			wantText:  "true : bool",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionExtensionMath",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--extension",
						"math",
					},
				},
				&evalCmd{
					expr: "math.greatest(1,2)",
				},
			},
			wantText:  "2 : int",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionExtensionEncoders",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--extension",
						"encoders",
					},
				},
				&evalCmd{
					expr: "base64.encode(b'hello')",
				},
			},
			wantText:  "aGVsbG8= : string",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionExtensionAll",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--extension",
						"all",
					},
				},
				&evalCmd{
					expr: "'test'.substring(2) == 'st' && " +
						"proto.getExt(google.expr.proto2.test.ExampleType{}, google.expr.proto2.test.int32_ext) == 0 && " +
						"math.greatest(1,2) == 2 && " +
						"base64.encode(b'hello') == 'aGVsbG8='",
				},
			},
			wantText:  "true : bool",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "OptionEnablePartialEval",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--enable_partial_eval",
					},
				},
				&letVarCmd{
					identifier: "x",
					typeHint:   mustParseType(t, "int"),
					src:        "",
				},
				&letVarCmd{
					identifier: "y",
					typeHint:   mustParseType(t, "int"),
					src:        "10",
				},
				&evalCmd{
					expr: "x + y > 10 || y > 10",
				},
			},
			wantText:  "Unknown x (1)",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "PartialEvalDisabled",
			commands: []Cmder{
				&letVarCmd{
					identifier: "x",
					typeHint:   mustParseType(t, "int"),
					src:        "",
				},
				&letVarCmd{
					identifier: "y",
					typeHint:   mustParseType(t, "int"),
					src:        "10",
				},
				&evalCmd{
					expr: "x + y > 10 || y > 10",
				},
			},
			wantText:  "",
			wantExit:  false,
			wantError: true,
		},
		{
			name: "PartialEvalFiltered",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--enable_partial_eval",
					},
				},
				&letVarCmd{
					identifier: "x",
					typeHint:   mustParseType(t, "int"),
					src:        "",
				},
				&letVarCmd{
					identifier: "y",
					typeHint:   mustParseType(t, "int"),
					src:        "11",
				},
				&evalCmd{
					expr: "x + y > 10 || y > 10",
				},
			},
			wantText:  "true : bool",
			wantExit:  false,
			wantError: false,
		},
		{
			name: "LoadDescriptorsError",
			commands: []Cmder{
				&simpleCmd{
					cmd: "load_descriptors",
					args: []string{
						"<not a file>",
					},
				},
			},
			wantText:  "",
			wantExit:  false,
			wantError: true,
		},
		{
			name: "LoadDescriptorsAddsTypes",
			commands: []Cmder{
				&simpleCmd{
					cmd: "load_descriptors",
					args: []string{
						testTextDescriptorFile,
					},
				},
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"google.rpc.context",
					},
				},
				&evalCmd{
					expr: "AttributeContext.Request{host: 'www.example.com'}",
				},
			},
			wantText:  `host:"www.example.com" : google.rpc.context.AttributeContext.Request`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "LoadDescriptorsPackageSpec",
			commands: []Cmder{
				&simpleCmd{
					cmd: "load_descriptors",
					args: []string{
						"--pkg",
						"cel-spec-test-types",
					},
				},
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"cel.expr.conformance",
					},
				},
				&evalCmd{
					expr: "proto2.TestAllTypes{single_int64: 42}.single_int64 == proto3.TestAllTypes{single_int64: 42}.single_int64",
				},
			},
			wantText:  `true : bool`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "LoadDescriptorsPackageRpc",
			commands: []Cmder{
				&simpleCmd{
					cmd: "load_descriptors",
					args: []string{
						"--pkg",
						"google-rpc",
					},
				},
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"google.rpc.context",
					},
				},
				&evalCmd{
					expr: "AttributeContext{source: AttributeContext.Peer{ip: '192.168.0.1'}}.source.ip == '192.168.0.1'",
				},
			},
			wantText:  `true : bool`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "LoadDescriptorsPackageUnknown",
			commands: []Cmder{
				&simpleCmd{
					cmd: "load_descriptors",
					args: []string{
						"--pkg",
						"'unknown'",
					},
				},
			},
			wantText:  "",
			wantExit:  false,
			wantError: true,
		},
		{
			name: "Status",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"google",
					},
				},
				&letVarCmd{
					identifier: "x",
					typeHint:   mustParseType(t, "int"),
					src:        "1",
				},
				&letFnCmd{
					identifier: "fn",
					src:        "2",
					resultType: mustParseType(t, "int"),
					params:     nil,
				},
				&simpleCmd{
					cmd:  "status",
					args: []string{},
				},
			},
			wantText: `// Options
%option --container 'google'

// Functions
%let fn() : int -> 2

// Variables
%let x = 1
`,
			wantExit:  false,
			wantError: false,
		},
		{
			name: "Reset",
			commands: []Cmder{
				&simpleCmd{
					cmd: "option",
					args: []string{
						"--container",
						"google",
					},
				},
				&letVarCmd{
					identifier: "x",
					typeHint:   mustParseType(t, "int"),
					src:        "1",
				},
				&letFnCmd{
					identifier: "fn",
					src:        "2",
					resultType: mustParseType(t, "int"),
					params:     nil,
				},
				&simpleCmd{
					cmd:  "reset",
					args: nil,
				},
				&simpleCmd{
					cmd:  "status",
					args: []string{},
				},
			},
			wantText: `// Options

// Functions

// Variables
`,
			wantExit:  false,
			wantError: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			eval, err := NewEvaluator()
			if err != nil {
				t.Fatalf("NewEvaluator returned error: %v, wanted nil", err)
			}
			eval.env, _ = eval.env.Extend(cel.Types(&proto2pb.ExampleType{}, &proto2pb.ExternalMessageType{}))
			n := len(tc.commands)
			for _, cmd := range tc.commands[:n-1] {
				// only need output of last command
				eval.Process(cmd)
			}
			text, exit, err := eval.Process(tc.commands[n-1])

			gotErr := false
			if err != nil {
				gotErr = true
			}
			wantText := tc.wantText
			if tc.ignoreWhitespace {
				text = stripWhitespace(text)
				wantText = stripWhitespace(wantText)
			}
			if text != wantText || exit != tc.wantExit || (gotErr != tc.wantError) {
				t.Errorf("For command %s got (output: '%s' exit: %v err: %v (%v)) wanted (output: '%s' exit: %v err: %v)",
					tc.commands[n-1], text, exit, gotErr, err, tc.wantText, tc.wantExit, tc.wantError)
			}
		})
	}
}

func TestProcessOptionError(t *testing.T) {
	var testCases = []struct {
		name     string
		command  Cmder
		errorMsg string
	}{
		{
			name: "OptionContainerNotEnoughArgs",
			command: &simpleCmd{
				cmd: "option",
				args: []string{
					"--container",
				},
			},
			errorMsg: "container: not enough arguments",
		},
		{
			name: "OptionExtensionNotEnoughArgs",
			command: &simpleCmd{
				cmd: "option",
				args: []string{
					"--extension",
				},
			},
			errorMsg: "extension: not enough arguments",
		},
		{
			name: "OptionExtensionInvalid",
			command: &simpleCmd{
				cmd: "option",
				args: []string{
					"--extension",
					"'bogus'",
				},
			},
			errorMsg: "extension: Unknown option: 'bogus'. Available options are: ['all', 'bindings', 'encoders', 'lists', 'math', 'optional', 'protos', 'sets', 'strings']",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			eval, err := NewEvaluator()
			if err != nil {
				t.Fatalf("NewEvaluator returned error: %v, wanted nil", err)
			}
			_, _, err = eval.Process(tc.command)

			if err == nil {
				t.Fatalf("Expected an error processing command: %s", tc.command)
			}

			if err.Error() != tc.errorMsg {
				t.Errorf("For command %s got (error: '%s') wanted (error: '%s')",
					tc.command, err.Error(), tc.errorMsg)
			}
		})
	}
}

func stripWhitespace(a string) string {
	a = strings.Replace(a, " ", "", -1)
	a = strings.Replace(a, "\n", "", -1)
	a = strings.Replace(a, "\t", "", -1)
	return strings.Replace(a, "\r", "", -1)
}
