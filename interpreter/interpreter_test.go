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
	"reflect"
	"testing"

	"github.com/google/cel-go/common"

	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"
	"github.com/google/cel-go/test"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestExhaustiveInterpreter_ConditionalExpr(t *testing.T) {
	// a ? b < 1.0 : c == ["hello"]
	// Operator "_==_" is at Expr 6, should be evaluated in exhaustive mode
	// even though "a" is true
	program := NewExhaustiveProgram(
		test.Conditional.Expr,
		test.Conditional.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a": true,
			"b": 0.999,
			"c": types.NewStringList([]string{"hello"})}))
	ev, _ := state.Value(6)
	// "==" should be evaluated in exhaustive mode though unnecessary
	if ev != types.True {
		t.Errorf("Else expression expected to be true, got: %v", ev)
	}
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestExhaustiveInterpreter_ConditionalExprErr(t *testing.T) {
	// a ? b < 1.0 : c == ["hello"]
	// Operator "<" is at Expr 3, "_==_" is at Expr 6.
	// Both should be evaluated in exhaustive mode though a is not provided
	program := NewExhaustiveProgram(
		test.Conditional.Expr,
		test.Conditional.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"b": 1.001,
			"c": types.NewStringList([]string{"hello"})}))
	iv, _ := state.Value(3)
	// "<" should be evaluated in exhaustive mode though unnecessary
	if iv != types.False {
		t.Errorf("If expression expected to be false, got: %v", iv)
	}
	ev, _ := state.Value(6)
	// "==" should be evaluated in exhaustive mode though unnecessary
	if ev != types.True {
		t.Errorf("Else expression expected to be true, got: %v", ev)
	}
	if result.Type() != types.UnknownType {
		t.Errorf("Expected unknown result, got: %v", result)
	}
}

func TestExhaustiveInterpreter_LogicalOrEquals(t *testing.T) {
	// a || b == "b"
	// Operator "==" is at Expr 4, should be evaluated though "a" is true
	program := NewExhaustiveProgram(
		test.LogicalOrEquals.Expr,
		test.LogicalOrEquals.Info(t.Name()))

	// TODO: make the type identifiers part of the standard declaration set.
	provider := types.NewProvider(&exprpb.Expr{})
	i := NewStandardInterpreter(packages.NewPackage("test"), provider)
	interpretable := i.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a": true,
			"b": "b",
		}))
	rhv, _ := state.Value(4)
	// "==" should be evaluated in exhaustive mode though unnecessary
	if rhv != types.True {
		t.Errorf("Right hand side expression expected to be true, got: %v", rhv)
	}
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_CallExpr(t *testing.T) {
	program := NewProgram(
		test.Equality.Expr,
		test.Equality.Info(t.Name()))
	intr := NewStandardInterpreter(
		packages.NewPackage("google.api.expr"),
		types.NewProvider(&exprpb.ParsedExpr{}))
	interpretable := intr.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{"a": int64(41)}))
	if result != types.False {
		t.Errorf("Expected false, got: %v", result)
	}
	if ident, found := state.Value(1); !found || ident != types.Int(41) {
		t.Errorf("State of ident 'a' != 41, got: %v", ident)
	}
}

func TestInterpreter_SelectExpr(t *testing.T) {
	program := NewProgram(
		test.Select.Expr,
		test.Select.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a.b": types.NewDynamicMap(map[string]bool{"c": true}),
		}))
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_ConditionalExpr(t *testing.T) {
	// a ? b < 1.0 : c == ["hello"]
	// Operator "<" is at Expr 3, "_==_" is at Expr 6.
	program := NewProgram(
		test.Conditional.Expr,
		test.Conditional.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a": true,
			"b": 0.999,
			"c": types.NewStringList([]string{"hello"})}))
	ev, _ := state.Value(6)
	// "_==_" should not be evaluated in normal mode since a is true
	if ev != nil {
		t.Errorf("Else expression expected to be nil, got: %v", ev)
	}
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_ComprehensionExpr(t *testing.T) {
	result, _ := evalExpr(t, "[1, 1u, 1.0].exists(x, type(x) == uint)")
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_NonStrictExistsComprehension(t *testing.T) {
	result, _ := evalExpr(t, "[0, 2, 4].exists(x, 4/x == 2 && 4/(4-x) == 2)")
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_NonStrictAllComprehension(t *testing.T) {
	result, _ := evalExpr(t, "![0, 2, 4].all(x, 4/x != 2 && 4/(4-x) != 2)")
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_NonStrictAllWithInput(t *testing.T) {
	parsed := parseExpr(t,
		`code == "111" && ["a", "b"].all(x, x in tags)
		|| code == "222" && ["a", "b"].all(x, x in tags)`)
	pgrm := NewProgram(parsed.Expr, parsed.SourceInfo)
	i := interpreter.NewInterpretable(pgrm)
	result, _ := i.Eval(NewActivation(map[string]interface{}{
		"code": "222",
		"tags": []string{"a", "b"},
	}))
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_ExistsOne(t *testing.T) {
	result, _ := evalExpr(t, "[1, 2, 3].exists_one(x, x % 2 == 0)")
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_Map(t *testing.T) {
	result, _ := evalExpr(t, "[1, 2, 3].map(x, x * 2) == [2, 4, 6]")
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_Filter(t *testing.T) {
	result, _ := evalExpr(t, "[1, 2, 3].filter(x, x > 2) == [3]")
	if result != types.True {
		t.Errorf("Got %v, wanted true", result)
	}
}

func TestInterpreter_LogicalAnd(t *testing.T) {
	// a && {c: true}.c
	program := NewProgram(
		test.LogicalAnd.Expr,
		test.LogicalAnd.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	// TODO: make the type identifiers part of the standard declaration set.
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{"a": true}))
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_LogicalAndMissingType(t *testing.T) {
	// a && {c: true}.c
	program := NewProgram(
		test.LogicalAndMissingType.Expr,
		test.LogicalAndMissingType.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{"a": false}))
	if result != types.False {
		t.Errorf("Got: %v, wanted true", result)
	}
	result, _ = interpretable.Eval(
		NewActivation(map[string]interface{}{"a": true}))
	if !types.IsError(result) {
		t.Errorf("Got: %v, wanted error", result)
	}
}

func TestInterpreter_LogicalOr(t *testing.T) {
	// {c: false}.c || a
	program := NewProgram(
		test.LogicalOr.Expr,
		test.LogicalOr.Info(t.Name()))

	// TODO: make the type identifiers part of the standard declaration set.
	provider := types.NewProvider(&exprpb.Expr{})
	i := NewStandardInterpreter(packages.NewPackage("test"), provider)
	interpretable := i.NewInterpretable(program)
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{"a": true}))
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_LogicalOrEquals(t *testing.T) {
	// a || b == "b"
	// Operator "==" is at Expr 4, should not be evaluated since "a" is true
	program := NewProgram(
		test.LogicalOrEquals.Expr,
		test.LogicalOrEquals.Info(t.Name()))

	// TODO: make the type identifiers part of the standard declaration set.
	provider := types.NewProvider(&exprpb.Expr{})
	i := NewStandardInterpreter(packages.NewPackage("test"), provider)
	interpretable := i.NewInterpretable(program)
	result, state := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a": true,
			"b": "b",
		}))
	rhv, _ := state.Value(4)
	// "==" should not be evaluated in normal mode since it is unnecessary
	if rhv != nil {
		t.Errorf("Right hand side expression expected to be nil, got: %v", rhv)
	}
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_BuildObject(t *testing.T) {
	src := common.NewTextSource(
		"v1alpha1.Expr{id: 1, " +
			"const_expr: v1alpha1.Constant{ " +
			"string_value: \"oneof_test\"}}")
	parsed, errors := parser.Parse(src)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	pkgr := packages.NewPackage("google.api.expr")
	provider := types.NewProvider(&exprpb.Expr{})
	env := checker.NewStandardEnv(pkgr, provider)
	checked, errors := checker.Check(parsed, src, env)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	i := NewStandardInterpreter(pkgr, provider)
	eval := i.NewInterpretable(NewCheckedProgram(checked))
	result, _ := eval.Eval(emptyActivation)
	expected := &exprpb.Expr{Id: 1,
		ExprKind: &exprpb.Expr_ConstExpr{
			ConstExpr: &exprpb.Constant{
				ConstantKind: &exprpb.Constant_StringValue{
					StringValue: "oneof_test"}}}}
	if !proto.Equal(result.(ref.Value).Value().(proto.Message), expected) {
		t.Errorf("Could not build object properly. Got '%v', wanted '%v'",
			result.(ref.Value).Value(),
			expected)
	}
}

func TestInterpreter_ConstantReturnValue(t *testing.T) {
	parsed, err := parser.Parse(common.NewTextSource("1"))
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(emptyActivation)
	if int64(res.(types.Int)) != int64(1) {
		t.Errorf("Got '%v', wanted 1", res)
	}
}

func TestInterpreter_InList(t *testing.T) {
	parsed, err := parser.Parse(common.NewTextSource("1 in [1, 2, 3]"))
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(emptyActivation)
	if res != types.True {
		t.Errorf("Got '%v', wanted 'true'", res)
	}
}

func TestInterpreter_BuildMap(t *testing.T) {
	parsed, err := parser.Parse(common.NewTextSource("{'b': '''hi''', 'c': name}"))
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(NewActivation(map[string]interface{}{"name": "tristan"}))
	value, _ := res.(ref.Value).ConvertToNative(
		reflect.TypeOf(map[string]string{}))
	mapVal := value.(map[string]string)
	if mapVal["b"] != "hi" || mapVal["c"] != "tristan" {
		t.Errorf("Got '%v', expected map[b:hi c:tristan]", value)
	}
}

func TestInterpreter_MapIndex(t *testing.T) {
	parsed, err := parser.Parse(common.NewTextSource("{'a':1}['a']"))
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(emptyActivation)
	if res != types.Int(1) {
		t.Errorf("Got '%v', wanted 1", res)
	}
}

func BenchmarkInterpreter_ConditionalExpr(b *testing.B) {
	// a ? b < 1.0 : c == ["hello"]
	program := NewProgram(
		test.Conditional.Expr,
		test.Conditional.Info(b.Name()))
	interpretable := interpreter.NewInterpretable(program)
	activation := NewActivation(map[string]interface{}{
		"a": types.False,
		"b": types.Double(0.999),
		"c": types.NativeToValue([]string{"hello"})})
	for i := 0; i < b.N; i++ {
		interpretable.Eval(activation)
	}
}

func BenchmarkInterpreter_EqualsCall(b *testing.B) {
	// type(x) == uint
	activation := NewActivation(map[string]interface{}{
		"x": types.Uint(20)})
	d := NewDispatcher()
	d.Add(functions.StandardOverloads()...)
	evalState := NewEvalState(4)
	for i := 0; i < b.N; i++ {
		xRef, _ := activation.ResolveName("x")
		evalState.SetValue(1, xRef)
		xRef, _ = evalState.Value(1)
		typeOfXRef := xRef.ConvertToType(types.TypeType)
		evalState.SetValue(2, typeOfXRef)
		typeOfXRef, _ = evalState.Value(2)
		evalState.SetValue(3, typeOfXRef.Equal(types.UintType))
	}
}

func BenchmarkInterpreter_EqualsDispatch(b *testing.B) {
	// type(x) == uint
	activation := NewActivation(map[string]interface{}{
		"x": types.Uint(20)})
	d := NewDispatcher()
	d.Add(functions.StandardOverloads()...)
	p := types.NewProvider()
	callTypeOf := NewCall(2, "type", []int64{1})
	callEq := NewCall(3, "_==_", []int64{1, 2})
	evalState := NewEvalState(4)
	for i := 0; i < b.N; i++ {
		xRef, _ := activation.ResolveName("x")
		evalState.SetValue(1, xRef)
		xRef, _ = evalState.Value(1)
		ctxType := &CallContext{
			call:       callTypeOf,
			args:       []ref.Value{xRef},
			activation: activation,
		}
		evalState.SetValue(callTypeOf.GetID(), d.Dispatch(ctxType))
		typeOfXRef, _ := evalState.Value(callTypeOf.GetID())
		// not-found here.
		activation.ResolveName("uint")
		uintType, _ := p.FindIdent("uint")
		ctxEq := &CallContext{
			call:       callEq,
			args:       []ref.Value{typeOfXRef, uintType},
			activation: activation,
		}
		evalState.SetValue(callEq.GetID(), d.Dispatch(ctxEq))
	}
}

func BenchmarkInterpreter_EqualInstructions(b *testing.B) {
	// type(x) == uint
	program := NewProgram(
		test.TypeEquality.Expr,
		test.TypeEquality.Info(b.Name()))
	interpretable := interpreter.NewInterpretable(program)
	activation := NewActivation(map[string]interface{}{
		"x": types.Uint(20)})
	for i := 0; i < b.N; i++ {
		interpretable.Eval(activation)
	}
}

func BenchmarkInterpreter_ComprehensionExpr(b *testing.B) {
	// [1, 1u, 1.0].exists(x, type(x) == uint)
	program := NewProgram(
		test.Exists.Expr,
		test.Exists.Info(b.Name()))
	interpretable := interpreter.NewInterpretable(program)
	for i := 0; i < b.N; i++ {
		interpretable.Eval(emptyActivation)
	}
}

func BenchmarkInterpreter_ComprehensionExprWithInput(b *testing.B) {
	// elems.exists(x, type(x) == uint)
	program := NewProgram(
		test.ExistsWithInput.Expr,
		test.ExistsWithInput.Info(b.Name()))
	interpretable := interpreter.NewInterpretable(program)
	activation := NewActivation(map[string]interface{}{
		"elems": types.NativeToValue([]interface{}{0, 1, 2, 3, 4, uint(5), 6})})
	for i := 0; i < b.N; i++ {
		interpretable.Eval(activation)
	}
}

var (
	interpreter = NewStandardInterpreter(
		packages.DefaultPackage,
		types.NewProvider(&exprpb.ParsedExpr{}))
	emptyActivation = NewActivation(map[string]interface{}{})
)

func parseExpr(t *testing.T, src string) *exprpb.ParsedExpr {
	t.Helper()
	s := common.NewTextSource(src)
	parsed, errors := parser.Parse(s)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}
	return parsed
}

func evalExpr(t *testing.T, src string) (ref.Value, EvalState) {
	t.Helper()
	parsed := parseExpr(t, src)
	pgrm := NewProgram(parsed.Expr, parsed.SourceInfo)
	eval := interpreter.NewInterpretable(pgrm)
	return eval.Eval(emptyActivation)
}
