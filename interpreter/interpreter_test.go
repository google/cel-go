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
	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"
	"github.com/google/cel-go/test"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"reflect"
	"testing"
)

func TestInterpreter_CallExpr(t *testing.T) {
	program := NewProgram(
		test.Equality.Expr,
		test.Equality.Info(t.Name()))
	intr := NewStandardInterpreter(
		packages.NewPackage("google.api.expr"),
		types.NewProvider(&expr.ParsedExpr{}))
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
	program := NewProgram(
		test.Conditional.Expr,
		test.Conditional.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{
			"a": true,
			"b": 0.999}))
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_ComprehensionExpr(t *testing.T) {
	// [1, 1u, 1.0].exists(x, type(x) == uint)
	program := NewProgram(
		test.Exists.Expr,
		test.Exists.Info(t.Name()))

	interpretable := interpreter.NewInterpretable(program)
	// TODO: make the type identifiers part of the standard declaration set.
	result, _ := interpretable.Eval(
		NewActivation(map[string]interface{}{}))
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_BuildObject(t *testing.T) {
	parsed, errors := parser.ParseText("v1alpha1.Expr{id: 1, " +
		"const_expr: v1alpha1.Constant{string_value: \"oneof_test\"}}")
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	pkgr := packages.NewPackage("google.api.expr")
	provider := types.NewProvider(&expr.Expr{})
	env := checker.NewStandardEnv(pkgr, provider, errors)
	checked := checker.Check(parsed, env)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	i := NewStandardInterpreter(pkgr, provider)
	eval := i.NewInterpretable(NewCheckedProgram(checked))
	result, _ := eval.Eval(NewActivation(map[string]interface{}{}))
	expected := &expr.Expr{Id: 1,
		ExprKind: &expr.Expr_ConstExpr{
			ConstExpr: &expr.Constant{
				ConstantKind: &expr.Constant_StringValue{
					StringValue: "oneof_test"}}}}
	if !proto.Equal(result.(ref.Value).Value().(proto.Message), expected) {
		t.Errorf("Could not build object properly. Got '%v', wanted '%v'",
			result.(ref.Value).Value(),
			expected)
	}
}

func TestInterpreter_ConstantReturnValue(t *testing.T) {
	parsed, err := parser.ParseText("1")
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(NewActivation(map[string]interface{}{}))
	if int64(res.(types.Int)) != int64(1) {
		t.Errorf("Got '%v', wanted 1", res)
	}
}

func TestInterpreter_InList(t *testing.T) {
	parsed, err := parser.ParseText("1 in [1, 2, 3]")
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(NewActivation(map[string]interface{}{}))
	if res != types.True {
		t.Errorf("Got '%v', wanted 'true'", res)
	}
}

func TestInterpreter_BuildMap(t *testing.T) {
	parsed, err := parser.ParseText("{'b': '''hi''', 'c': name}")
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
	parsed, err := parser.ParseText("{'a':1}['a']")
	if len(err.GetErrors()) != 0 {
		t.Error(err)
	}
	prg := NewProgram(parsed.GetExpr(), parsed.GetSourceInfo())
	i := interpreter.NewInterpretable(prg)
	res, _ := i.Eval(NewActivation(map[string]interface{}{}))
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
		evalState.SetValue(callTypeOf.GetId(), d.Dispatch(ctxType))
		typeOfXRef, _ := evalState.Value(callTypeOf.GetId())
		// not-found here.
		activation.ResolveName("uint")
		uintType, _ := p.FindIdent("uint")
		ctxEq := &CallContext{
			call:       callEq,
			args:       []ref.Value{typeOfXRef, uintType},
			activation: activation,
		}
		evalState.SetValue(callEq.GetId(), d.Dispatch(ctxEq))
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
	activation := NewActivation(map[string]interface{}{})
	for i := 0; i < b.N; i++ {
		interpretable.Eval(activation)
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
		types.NewProvider(&expr.ParsedExpr{}))
)
