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
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/parser"
)

const (
	mapInsert                 = "cel.@mapInsert"
	mapInsertOverloadMap      = "@mapInsert_map_map"
	mapInsertOverloadKeyValue = "@mapInsert_map_key_value"
)

// TwoVarComprehensions introduces support for two-variable comprehensions.
//
// The two-variable form of comprehensions looks similar to the one-variable counterparts.
// Where possible, the same macro names were used and additional macro signatures added.
// The notable distinction for two-variable comprehensions is the introduction of
// `transformList`, `transformMap`, and `transformMapEntry` support for list and map types
// rather than the more traditional `map` and `filter` macros.
func TwoVarComprehensions() cel.EnvOption {
	return cel.Lib(compreV2Lib{})
}

type compreV2Lib struct{}

func (compreV2Lib) LibraryName() string {
	return "cel.lib.ext.comprev2"
}

func (compreV2Lib) CompileOptions() []cel.EnvOption {
	kType := cel.TypeParamType("K")
	vType := cel.TypeParamType("V")
	mapKVType := cel.MapType(kType, vType)
	opts := []cel.EnvOption{
		cel.Macros(
			cel.ReceiverMacro("all", 3, quantifierAll),
			cel.ReceiverMacro("exists", 3, quantifierExists),
			cel.ReceiverMacro("existsOne", 3, quantifierExistsOne),
			cel.ReceiverMacro("exists_one", 3, quantifierExistsOne),
			cel.ReceiverMacro("transformList", 3, transformList),
			cel.ReceiverMacro("transformList", 4, transformList),
			cel.ReceiverMacro("transformMap", 3, transformMap),
			cel.ReceiverMacro("transformMap", 4, transformMap),
			cel.ReceiverMacro("transformMapEntry", 3, transformMapEntry),
			cel.ReceiverMacro("transformMapEntry", 4, transformMapEntry),
		),
		cel.Function(mapInsert,
			cel.Overload(mapInsertOverloadKeyValue, []*cel.Type{mapKVType, kType, vType}, mapKVType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					m := args[0].(traits.Mapper)
					k := args[1]
					v := args[2]
					return types.InsertMapKeyValue(m, k, v)
				})),
			cel.Overload(mapInsertOverloadMap, []*cel.Type{mapKVType, mapKVType}, mapKVType,
				cel.BinaryBinding(func(targetMap, updateMap ref.Val) ref.Val {
					tm := targetMap.(traits.Mapper)
					um := updateMap.(traits.Mapper)
					umIt := um.Iterator()
					for umIt.HasNext() == types.True {
						k := umIt.Next()
						updateOrErr := types.InsertMapKeyValue(tm, k, um.Get(k))
						if types.IsError(updateOrErr) {
							return updateOrErr
						}
						tm = updateOrErr.(traits.Mapper)
					}
					return tm
				})),
		),
	}
	return opts
}

func (compreV2Lib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func quantifierAll(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	iterVar1, err := extractIterVar(mef, args[0])
	if err != nil {
		return nil, err
	}
	iterVar2, err := extractIterVar(mef, args[1])
	if err != nil {
		return nil, err
	}
	return mef.NewComprehensionTwoVar(
		target,
		iterVar1,
		iterVar2,
		parser.AccumulatorName,
		/*accuInit=*/ mef.NewLiteral(types.True),
		/*condition=*/ mef.NewCall(operators.NotStrictlyFalse, mef.NewAccuIdent()),
		/*step=*/ mef.NewCall(operators.LogicalAnd, mef.NewAccuIdent(), args[2]),
		/*result=*/ mef.NewAccuIdent(),
	), nil
}

func quantifierExists(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	iterVar1, err := extractIterVar(mef, args[0])
	if err != nil {
		return nil, err
	}
	iterVar2, err := extractIterVar(mef, args[1])
	if err != nil {
		return nil, err
	}
	return mef.NewComprehensionTwoVar(
		target,
		iterVar1,
		iterVar2,
		parser.AccumulatorName,
		/*accuInit=*/ mef.NewLiteral(types.False),
		/*condition=*/ mef.NewCall(operators.NotStrictlyFalse, mef.NewCall(operators.LogicalNot, mef.NewAccuIdent())),
		/*step=*/ mef.NewCall(operators.LogicalOr, mef.NewAccuIdent(), args[2]),
		/*result=*/ mef.NewAccuIdent(),
	), nil
}

func quantifierExistsOne(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	iterVar1, err := extractIterVar(mef, args[0])
	if err != nil {
		return nil, err
	}
	iterVar2, err := extractIterVar(mef, args[1])
	if err != nil {
		return nil, err
	}
	return mef.NewComprehensionTwoVar(
		target,
		iterVar1,
		iterVar2,
		parser.AccumulatorName,
		/*accuInit=*/ mef.NewLiteral(types.Int(0)),
		/*condition=*/ mef.NewLiteral(types.True),
		/*step=*/ mef.NewCall(operators.Conditional, args[2],
			mef.NewCall(operators.Add, mef.NewAccuIdent(), mef.NewLiteral(types.Int(1))),
			mef.NewAccuIdent()),
		/*result=*/ mef.NewCall(operators.Equals, mef.NewAccuIdent(), mef.NewLiteral(types.Int(1))),
	), nil
}

func transformList(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	iterVar1, err := extractIterVar(mef, args[0])
	if err != nil {
		return nil, err
	}
	iterVar2, err := extractIterVar(mef, args[1])
	if err != nil {
		return nil, err
	}

	var transform ast.Expr
	var filter ast.Expr
	if len(args) == 4 {
		filter = args[2]
		transform = args[3]
	} else {
		filter = nil
		transform = args[2]
	}

	//  __result__ = __result__ + [transform]
	step := mef.NewCall(operators.Add, mef.NewAccuIdent(), mef.NewList(transform))
	if filter != nil {
		//  __result__ = (filter) ? __result__ + [transform] : __result__
		step = mef.NewCall(operators.Conditional, filter, step, mef.NewAccuIdent())
	}

	return mef.NewComprehensionTwoVar(
		target,
		iterVar1,
		iterVar2,
		parser.AccumulatorName,
		/*accuInit=*/ mef.NewList(),
		/*condition=*/ mef.NewLiteral(types.True),
		step,
		/*result=*/ mef.NewAccuIdent(),
	), nil
}

func transformMap(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	iterVar1, err := extractIterVar(mef, args[0])
	if err != nil {
		return nil, err
	}
	iterVar2, err := extractIterVar(mef, args[1])
	if err != nil {
		return nil, err
	}

	var transform ast.Expr
	var filter ast.Expr
	if len(args) == 4 {
		filter = args[2]
		transform = args[3]
	} else {
		filter = nil
		transform = args[2]
	}

	// __result__ = cel.@mapInsert(__result__, iterVar1, transform)
	step := mef.NewCall(mapInsert, mef.NewAccuIdent(), mef.NewIdent(iterVar1), transform)
	if filter != nil {
		// __result__ = (filter) ? cel.@mapInsert(__result__, iterVar1, transform) : __result__
		step = mef.NewCall(operators.Conditional, filter, step, mef.NewAccuIdent())
	}
	return mef.NewComprehensionTwoVar(
		target,
		iterVar1,
		iterVar2,
		parser.AccumulatorName,
		/*accuInit=*/ mef.NewMap(),
		/*condition=*/ mef.NewLiteral(types.True),
		step,
		/*result=*/ mef.NewAccuIdent(),
	), nil
}

func transformMapEntry(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	iterVar1, err := extractIterVar(mef, args[0])
	if err != nil {
		return nil, err
	}
	iterVar2, err := extractIterVar(mef, args[1])
	if err != nil {
		return nil, err
	}

	var transform ast.Expr
	var filter ast.Expr
	if len(args) == 4 {
		filter = args[2]
		transform = args[3]
	} else {
		filter = nil
		transform = args[2]
	}

	// __result__ = cel.@mapInsert(__result__, transform)
	step := mef.NewCall(mapInsert, mef.NewAccuIdent(), transform)
	if filter != nil {
		// __result__ = (filter) ? cel.@mapInsert(__result__, transform) : __result__
		step = mef.NewCall(operators.Conditional, filter, step, mef.NewAccuIdent())
	}
	return mef.NewComprehensionTwoVar(
		target,
		iterVar1,
		iterVar2,
		parser.AccumulatorName,
		/*accuInit=*/ mef.NewMap(),
		/*condition=*/ mef.NewLiteral(types.True),
		step,
		/*result=*/ mef.NewAccuIdent(),
	), nil
}

func extractIterVar(meh cel.MacroExprFactory, target ast.Expr) (string, *cel.Error) {
	iterVar, found := extractIdent(target)
	if !found {
		return "", meh.NewError(target.ID(), "iteration variable must be a simple name")
	}
	return iterVar, nil
}
