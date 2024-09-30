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
	mapInsert         = "cel.@mapInsert"
	mapInsertOverload = "@mapInsert_map_key_value"
)

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
	vPrimeType := cel.TypeParamType("V1")
	mapKVType := cel.MapType(kType, vType)
	mapKVPrimeType := cel.MapType(kType, vPrimeType)
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
		),
		cel.Function(mapInsert,
			cel.Overload(mapInsertOverload, []*cel.Type{mapKVType, kType, vType}, mapKVPrimeType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					m := args[0].(traits.Mapper)
					k := args[1]
					v := args[2]
					return types.InsertMapKeyValue(m, k, v)
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

func extractIterVar(meh cel.MacroExprFactory, target ast.Expr) (string, *cel.Error) {
	iterVar, found := extractIdent(target)
	if !found {
		return "", meh.NewError(target.ID(), "iteration variable must be a simple name")
	}
	return iterVar, nil
}
