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
	"math"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Math returns a cel.EnvOption to configure namespaced math helper macros and
// functions.
//
// Note, all macros use the 'math' namespace; however, at the time of macro
// expansion the namespace looks just like any other identifier. If you are
// currently using a variable named 'math', the macro will likely work just as
// intended; however, there is some chance for collision.
//
// # Math.Greatest
//
// Returns the greatest valued number present in the arguments to the macro.
//
// Greatest is a variable argument count macro which must take at least one
// argument. Simple numeric and list literals are supported as valid argument
// types; however, other literals will be flagged as errors during macro
// expansion. If the argument expression does not resolve to a numeric or
// list(numeric) type during type-checking, or during runtime then an error
// will be produced. If a list argument is empty, this too will produce an
// error.
//
//	math.greatest(<arg>, ...) -> <double|int|uint>
//
// Examples:
//
//	math.greatest(1)      // 1
//	math.greatest(1u, 2u) // 2u
//	math.greatest(-42.0, -21.5, -100.0)   // -21.5
//	math.greatest([-42.0, -21.5, -100.0]) // -21.5
//	math.greatest(numbers) // numbers must be list(numeric)
//
//	math.greatest()         // parse error
//	math.greatest('string') // parse error
//	math.greatest(a, b)     // check-time error if a or b is non-numeric
//	math.greatest(dyn('string')) // runtime error
//
// # Math.Least
//
// Returns the least valued number present in the arguments to the macro.
//
// Least is a variable argument count macro which must take at least one
// argument. Simple numeric and list literals are supported as valid argument
// types; however, other literals will be flagged as errors during macro
// expansion. If the argument expression does not resolve to a numeric or
// list(numeric) type during type-checking, or during runtime then an error
// will be produced. If a list argument is empty, this too will produce an
// error.
//
//	math.least(<arg>, ...) -> <double|int|uint>
//
// Examples:
//
//	math.least(1)      // 1
//	math.least(1u, 2u) // 1u
//	math.least(-42.0, -21.5, -100.0)   // -100.0
//	math.least([-42.0, -21.5, -100.0]) // -100.0
//	math.least(numbers) // numbers must be list(numeric)
//
//	math.least()         // parse error
//	math.least('string') // parse error
//	math.least(a, b)     // check-time error if a or b is non-numeric
//	math.least(dyn('string')) // runtime error
func Math() cel.EnvOption {
	return cel.Lib(&mathLib{version: math.MaxUint32})
}

const (
	mathNamespace = "math"
	leastMacro    = "least"
	greatestMacro = "greatest"
	bitAndMacro   = "bitAnd"
	bitOrMacro    = "bitOr"
	bitXorMacro   = "bitXor"

	// Min-max functions
	minFunc = "math.@min"
	maxFunc = "math.@max"

	// Rounding functions
	ceilFunc  = "math.ceil"
	floorFunc = "math.floor"
	roundFunc = "math.round"
	truncFunc = "math.trunc"

	// Floating point helper functions
	isInfFunc    = "math.isInf"
	isNanFunc    = "math.isNaN"
	isFiniteFunc = "math.isFinite"

	// Signedness functions
	absFunc  = "math.abs"
	signFunc = "math.sign"

	// Bitwise functions
	bitAndFunc        = "math.@bitAnd"
	bitOrFunc         = "math.@bitOr"
	bitXorFunc        = "math.@bitXor"
	bitNotFunc        = "math.bitNot"
	bitShiftLeftFunc  = "math.bitShiftLeft"
	bitShiftRightFunc = "math.bitShiftRight"
)

var (
	errIntOverflow = types.NewErr("integer overflow")
)

type MathOption func(*mathLib) *mathLib

func MathVersion(version uint32) MathOption {
	return func(lib *mathLib) *mathLib {
		lib.version = version
		return lib
	}
}

type mathLib struct {
	version uint32
}

// LibraryName implements the SingletonLibrary interface method.
func (*mathLib) LibraryName() string {
	return "cel.lib.ext.math"
}

// CompileOptions implements the Library interface method.
func (lib *mathLib) CompileOptions() []cel.EnvOption {
	opts := []cel.EnvOption{
		cel.Macros(
			// math.least(num, ...)
			cel.ReceiverVarArgMacro(leastMacro, mathLeast),
			// math.greatest(num, ...)
			cel.ReceiverVarArgMacro(greatestMacro, mathGreatest),
		),
		cel.Function(minFunc,
			cel.Overload("math_@min_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
				cel.UnaryBinding(identity)),
			cel.Overload("math_@min_int", []*cel.Type{cel.IntType}, cel.IntType,
				cel.UnaryBinding(identity)),
			cel.Overload("math_@min_uint", []*cel.Type{cel.UintType}, cel.UintType,
				cel.UnaryBinding(identity)),
			cel.Overload("math_@min_double_double", []*cel.Type{cel.DoubleType, cel.DoubleType}, cel.DoubleType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_uint_uint", []*cel.Type{cel.UintType, cel.UintType}, cel.UintType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_int_uint", []*cel.Type{cel.IntType, cel.UintType}, cel.DynType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_int_double", []*cel.Type{cel.IntType, cel.DoubleType}, cel.DynType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_double_int", []*cel.Type{cel.DoubleType, cel.IntType}, cel.DynType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_double_uint", []*cel.Type{cel.DoubleType, cel.UintType}, cel.DynType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_uint_int", []*cel.Type{cel.UintType, cel.IntType}, cel.DynType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_uint_double", []*cel.Type{cel.UintType, cel.DoubleType}, cel.DynType,
				cel.BinaryBinding(minPair)),
			cel.Overload("math_@min_list_double", []*cel.Type{cel.ListType(cel.DoubleType)}, cel.DoubleType,
				cel.UnaryBinding(minList)),
			cel.Overload("math_@min_list_int", []*cel.Type{cel.ListType(cel.IntType)}, cel.IntType,
				cel.UnaryBinding(minList)),
			cel.Overload("math_@min_list_uint", []*cel.Type{cel.ListType(cel.UintType)}, cel.UintType,
				cel.UnaryBinding(minList)),
		),
		cel.Function(maxFunc,
			cel.Overload("math_@max_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
				cel.UnaryBinding(identity)),
			cel.Overload("math_@max_int", []*cel.Type{cel.IntType}, cel.IntType,
				cel.UnaryBinding(identity)),
			cel.Overload("math_@max_uint", []*cel.Type{cel.UintType}, cel.UintType,
				cel.UnaryBinding(identity)),
			cel.Overload("math_@max_double_double", []*cel.Type{cel.DoubleType, cel.DoubleType}, cel.DoubleType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_uint_uint", []*cel.Type{cel.UintType, cel.UintType}, cel.UintType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_int_uint", []*cel.Type{cel.IntType, cel.UintType}, cel.DynType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_int_double", []*cel.Type{cel.IntType, cel.DoubleType}, cel.DynType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_double_int", []*cel.Type{cel.DoubleType, cel.IntType}, cel.DynType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_double_uint", []*cel.Type{cel.DoubleType, cel.UintType}, cel.DynType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_uint_int", []*cel.Type{cel.UintType, cel.IntType}, cel.DynType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_uint_double", []*cel.Type{cel.UintType, cel.DoubleType}, cel.DynType,
				cel.BinaryBinding(maxPair)),
			cel.Overload("math_@max_list_double", []*cel.Type{cel.ListType(cel.DoubleType)}, cel.DoubleType,
				cel.UnaryBinding(maxList)),
			cel.Overload("math_@max_list_int", []*cel.Type{cel.ListType(cel.IntType)}, cel.IntType,
				cel.UnaryBinding(maxList)),
			cel.Overload("math_@max_list_uint", []*cel.Type{cel.ListType(cel.UintType)}, cel.UintType,
				cel.UnaryBinding(maxList)),
		),
	}
	if lib.version >= 1 {
		opts = append(opts,
			cel.Macros(
				// math.bitOr(num, ...)
				cel.ReceiverVarArgMacro(bitOrMacro, mathBitLogic(bitOrFunc)),
				// math.bitAnd(num, ...)
				cel.ReceiverVarArgMacro(bitAndMacro, mathBitLogic(bitAndFunc)),
				// math.bitXor(num, ...)
				cel.ReceiverVarArgMacro(bitXorMacro, mathBitLogic(bitXorFunc)),
			),
			// Rounding function declarations
			cel.Function(ceilFunc,
				cel.Overload("math_ceil_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
					cel.UnaryBinding(ceil))),
			cel.Function(floorFunc,
				cel.Overload("math_floor_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
					cel.UnaryBinding(floor))),
			cel.Function(roundFunc,
				cel.Overload("math_round_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
					cel.UnaryBinding(round))),
			cel.Function(truncFunc,
				cel.Overload("math_trunc_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
					cel.UnaryBinding(trunc))),

			// Floating point helpers
			cel.Function(isInfFunc,
				cel.Overload("math_isInf_double", []*cel.Type{cel.DoubleType}, cel.BoolType,
					cel.UnaryBinding(isInf))),
			cel.Function(isNanFunc,
				cel.Overload("math_isNaN_double", []*cel.Type{cel.DoubleType}, cel.BoolType,
					cel.UnaryBinding(isNaN))),
			cel.Function(isFiniteFunc,
				cel.Overload("math_isFinite_double", []*cel.Type{cel.DoubleType}, cel.BoolType,
					cel.UnaryBinding(isFinite))),

			// Signedness functions
			cel.Function(absFunc,
				cel.Overload("math_abs_double", []*cel.Type{cel.DoubleType}, cel.DoubleType,
					cel.UnaryBinding(absDouble)),
				cel.Overload("math_abs_int", []*cel.Type{cel.IntType}, cel.IntType,
					cel.UnaryBinding(absInt)),
				cel.Overload("math_abs_uint", []*cel.Type{cel.UintType}, cel.UintType,
					cel.UnaryBinding(identity)),
			),
			cel.Function(signFunc,
				cel.Overload("math_sign_double", []*cel.Type{cel.DoubleType}, cel.DoubleType),
				cel.Overload("math_sign_int", []*cel.Type{cel.IntType}, cel.IntType),
				cel.Overload("math_sign_uint", []*cel.Type{cel.UintType}, cel.UintType),
				cel.SingletonUnaryBinding(sign),
			),

			// Bitwise operator declarations
			cel.Function(bitAndFunc,
				cel.Overload("math_@bitAnd_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
					cel.BinaryBinding(bitAndPairInt)),
				cel.Overload("math_@bitAnd_uint_uint", []*cel.Type{cel.UintType, cel.UintType}, cel.UintType,
					cel.BinaryBinding(bitAndPairUint)),
				cel.Overload("math_@bitAnd_list_int", []*cel.Type{cel.ListType(cel.IntType)}, cel.IntType,
					cel.UnaryBinding(bitAndListInt)),
				cel.Overload("math_@bitAnd_list_uint", []*cel.Type{cel.ListType(cel.UintType)}, cel.UintType,
					cel.UnaryBinding(bitAndListUint)),
			),
			cel.Function(bitOrFunc,
				cel.Overload("math_@bitOr_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
					cel.BinaryBinding(bitOrPairInt)),
				cel.Overload("math_@bitOr_uint_uint", []*cel.Type{cel.UintType, cel.UintType}, cel.UintType,
					cel.BinaryBinding(bitOrPairUint)),
				cel.Overload("math_@bitOr_list_int", []*cel.Type{cel.ListType(cel.IntType)}, cel.IntType,
					cel.UnaryBinding(bitOrListInt)),
				cel.Overload("math_@bitOr_list_uint", []*cel.Type{cel.ListType(cel.UintType)}, cel.UintType,
					cel.UnaryBinding(bitOrListUint)),
			),
			cel.Function(bitXorFunc,
				cel.Overload("math_@bitXor_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
					cel.BinaryBinding(bitXorPairInt)),
				cel.Overload("math_@bitXor_uint_uint", []*cel.Type{cel.UintType, cel.UintType}, cel.UintType,
					cel.BinaryBinding(bitXorPairUint)),
				cel.Overload("math_@bitXor_list_int", []*cel.Type{cel.ListType(cel.IntType)}, cel.IntType,
					cel.UnaryBinding(bitXorListInt)),
				cel.Overload("math_@bitXor_list_uint", []*cel.Type{cel.ListType(cel.UintType)}, cel.UintType,
					cel.UnaryBinding(bitXorListUint)),
			),
			cel.Function(bitNotFunc,
				cel.Overload("math_bitNot_uint_int", []*cel.Type{cel.UintType}, cel.IntType,
					cel.UnaryBinding(bitNotInt)),
				cel.Overload("math_bitNot_int_uint", []*cel.Type{cel.IntType}, cel.UintType,
					cel.UnaryBinding(bitNotUint)),
			),
			cel.Function(bitShiftLeftFunc,
				cel.Overload("math_bitShiftLeft_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
					cel.BinaryBinding(bitShiftLeftIntInt)),
				cel.Overload("math_bitShiftLeft_uint_int", []*cel.Type{cel.UintType, cel.IntType}, cel.UintType,
					cel.BinaryBinding(bitShiftLeftUintInt)),
			),
			cel.Function(bitShiftRightFunc,
				cel.Overload("math_bitShiftRight_int_int", []*cel.Type{cel.IntType, cel.IntType}, cel.IntType,
					cel.BinaryBinding(bitShiftRightIntInt)),
				cel.Overload("math_bitShiftRight_uint_int", []*cel.Type{cel.UintType, cel.IntType}, cel.UintType,
					cel.BinaryBinding(bitShiftRightUintInt)),
			),
		)
	}
	return opts
}

// ProgramOptions implements the Library interface method.
func (*mathLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func mathLeast(meh cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	if !macroTargetMatchesNamespace(mathNamespace, target) {
		return nil, nil
	}
	switch len(args) {
	case 0:
		return nil, meh.NewError(target.ID(), "math.least() requires at least one argument")
	case 1:
		if isListLiteralWithNumericArgs(args[0]) || isNumericArgType(args[0]) {
			return meh.NewCall(minFunc, args[0]), nil
		}
		return nil, meh.NewError(args[0].ID(), "math.least() invalid single argument value")
	case 2:
		err := checkInvalidArgs(meh, "math.least()", args)
		if err != nil {
			return nil, err
		}
		return meh.NewCall(minFunc, args...), nil
	default:
		err := checkInvalidArgs(meh, "math.least()", args)
		if err != nil {
			return nil, err
		}
		return meh.NewCall(minFunc, meh.NewList(args...)), nil
	}
}

func mathGreatest(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	if !macroTargetMatchesNamespace(mathNamespace, target) {
		return nil, nil
	}
	switch len(args) {
	case 0:
		return nil, mef.NewError(target.ID(), "math.greatest() requires at least one argument")
	case 1:
		if isListLiteralWithNumericArgs(args[0]) || isNumericArgType(args[0]) {
			return mef.NewCall(maxFunc, args[0]), nil
		}
		return nil, mef.NewError(args[0].ID(), "math.greatest() invalid single argument value")
	case 2:
		err := checkInvalidArgs(mef, "math.greatest()", args)
		if err != nil {
			return nil, err
		}
		return mef.NewCall(maxFunc, args...), nil
	default:
		err := checkInvalidArgs(mef, "math.greatest()", args)
		if err != nil {
			return nil, err
		}
		return mef.NewCall(maxFunc, mef.NewList(args...)), nil
	}
}

func mathBitLogic(bitFuncName string) cel.MacroFactory {
	return func(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
		bitFuncDisplayName := bitFuncName + "()"
		if !macroTargetMatchesNamespace(mathNamespace, target) {
			return nil, nil
		}
		switch len(args) {
		case 0:
			return nil, mef.NewError(target.ID(), bitFuncDisplayName+" requires at least one argument")
		case 1:
			if isListLiteralWithNumericArgs(args[0]) || isNumericArgType(args[0]) {
				return mef.NewCall(bitFuncName, mef.NewList(args[0])), nil
			}
			return nil, mef.NewError(args[0].ID(), bitFuncDisplayName+" invalid single argument value")
		case 2:
			err := checkInvalidArgs(mef, bitFuncDisplayName, args)
			if err != nil {
				return nil, err
			}
			return mef.NewCall(bitFuncName, args...), nil
		default:
			err := checkInvalidArgs(mef, bitFuncDisplayName, args)
			if err != nil {
				return nil, err
			}
			return mef.NewCall(bitFuncName, mef.NewList(args...)), nil
		}
	}
}

func identity(val ref.Val) ref.Val {
	return val
}

func ceil(val ref.Val) ref.Val {
	v := val.(types.Double)
	return types.Double(math.Ceil(float64(v)))
}

func floor(val ref.Val) ref.Val {
	v := val.(types.Double)
	return types.Double(math.Floor(float64(v)))
}

func round(val ref.Val) ref.Val {
	v := val.(types.Double)
	return types.Double(math.Round(float64(v)))
}

func trunc(val ref.Val) ref.Val {
	v := val.(types.Double)
	return types.Double(math.Trunc(float64(v)))
}

func isInf(val ref.Val) ref.Val {
	v := val.(types.Double)
	return types.Bool(math.IsInf(float64(v), 0))
}

func isFinite(val ref.Val) ref.Val {
	v := float64(val.(types.Double))
	return types.Bool(!math.IsInf(v, 0) && !math.IsNaN(v))
}

func isNaN(val ref.Val) ref.Val {
	v := val.(types.Double)
	return types.Bool(math.IsNaN(float64(v)))
}

func absDouble(val ref.Val) ref.Val {
	v := float64(val.(types.Double))
	return types.Double(math.Abs(v))
}

func absInt(val ref.Val) ref.Val {
	v := int64(val.(types.Int))
	if v == math.MinInt64 {
		return errIntOverflow
	}
	if v >= 0 {
		return val
	}
	return -types.Int(v)
}

func sign(val ref.Val) ref.Val {
	switch v := val.(type) {
	case types.Double:
		if isNaN(v) == types.True {
			return v
		}
		zero := types.Double(0)
		if v > zero {
			return types.Double(1)
		}
		if v < zero {
			return types.Double(-1)
		}
		return zero
	case types.Int:
		return v.Compare(types.IntZero)
	case types.Uint:
		if v == types.Uint(0) {
			return types.Uint(0)
		}
		return types.Uint(1)
	default:
		return maybeSuffixError(val, "math.sign")
	}
}

func bitAndListInt(values ref.Val) ref.Val {
	return bitOpList(values, "math.bitAnd(list)", bitAndPairInt)
}

func bitAndListUint(values ref.Val) ref.Val {
	return bitOpList(values, "math.bitAnd(list)", bitAndPairUint)
}

func bitOrListInt(values ref.Val) ref.Val {
	return bitOpList(values, "math.bitOr(list)", bitOrPairInt)
}

func bitOrListUint(values ref.Val) ref.Val {
	return bitOpList(values, "math.bitOr(list)", bitOrPairUint)
}

func bitXorListInt(values ref.Val) ref.Val {
	return bitOpList(values, "math.bitXor(list)", bitXorPairInt)
}

func bitXorListUint(values ref.Val) ref.Val {
	return bitOpList(values, "math.bitXor(list)", bitXorPairUint)
}

func bitOpList(values ref.Val, bitOpName string, bitOp func(value, bits ref.Val) ref.Val) ref.Val {
	l := values.(traits.Lister)
	size := l.Size().(types.Int)
	if size == types.IntZero {
		return types.NewErr("%s argument must not be empty", bitOpName)
	}
	result := l.Get(types.IntZero)
	for i := types.IntOne; i < size; i++ {
		result = bitOp(result, l.Get(i))
	}
	return result
}

func bitAndPairInt(first, second ref.Val) ref.Val {
	l := first.(types.Int)
	r := second.(types.Int)
	return l & r
}

func bitAndPairUint(first, second ref.Val) ref.Val {
	l := first.(types.Uint)
	r := second.(types.Uint)
	return l & r
}

func bitOrPairInt(first, second ref.Val) ref.Val {
	l := first.(types.Int)
	r := second.(types.Int)
	return l | r
}

func bitOrPairUint(first, second ref.Val) ref.Val {
	l := first.(types.Uint)
	r := second.(types.Uint)
	return l | r
}

func bitXorPairInt(first, second ref.Val) ref.Val {
	l := first.(types.Int)
	r := second.(types.Int)
	return l ^ r
}

func bitXorPairUint(first, second ref.Val) ref.Val {
	l := first.(types.Uint)
	r := second.(types.Uint)
	return l ^ r
}

func bitNotInt(value ref.Val) ref.Val {
	v := value.(types.Int)
	return ^v
}

func bitNotUint(value ref.Val) ref.Val {
	v := value.(types.Uint)
	return ^v
}

func bitShiftLeftIntInt(value, bits ref.Val) ref.Val {
	v := value.(types.Int)
	bs := bits.(types.Int)
	if bs < types.IntZero {
		return types.NewErr("math.bitShiftLeft() invalid shift count: %d", bs)
	}
	return v << bs
}

func bitShiftLeftUintInt(value, bits ref.Val) ref.Val {
	v := value.(types.Uint)
	bs := bits.(types.Int)
	if bs < types.IntZero {
		return types.NewErr("math.bitShiftLeft() invalid shift count: %d", bs)
	}
	return v << bs
}

func bitShiftRightIntInt(value, bits ref.Val) ref.Val {
	v := value.(types.Int)
	bs := bits.(types.Int)
	if bs < types.IntZero {
		return types.NewErr("math.bitShiftRight() invalid shift count: %d", bs)
	}
	return v >> bs
}

func bitShiftRightUintInt(value, bits ref.Val) ref.Val {
	v := value.(types.Uint)
	bs := bits.(types.Int)
	if bs < types.IntZero {
		return types.NewErr("math.bitShiftRight() invalid shift count: %d", bs)
	}
	return v >> bs
}

func minPair(first, second ref.Val) ref.Val {
	cmp, ok := first.(traits.Comparer)
	if !ok {
		return types.MaybeNoSuchOverloadErr(first)
	}
	out := cmp.Compare(second)
	if types.IsUnknownOrError(out) {
		return maybeSuffixError(out, "math.@min")
	}
	if out == types.IntOne {
		return second
	}
	return first
}

func minList(numList ref.Val) ref.Val {
	l := numList.(traits.Lister)
	size := l.Size().(types.Int)
	if size == types.IntZero {
		return types.NewErr("math.@min(list) argument must not be empty")
	}
	min := l.Get(types.IntZero)
	for i := types.IntOne; i < size; i++ {
		min = minPair(min, l.Get(i))
	}
	switch min.Type() {
	case types.IntType, types.DoubleType, types.UintType, types.UnknownType:
		return min
	default:
		return types.NewErr("no such overload: math.@min")
	}
}

func maxPair(first, second ref.Val) ref.Val {
	cmp, ok := first.(traits.Comparer)
	if !ok {
		return types.MaybeNoSuchOverloadErr(first)
	}
	out := cmp.Compare(second)
	if types.IsUnknownOrError(out) {
		return maybeSuffixError(out, "math.@max")
	}
	if out == types.IntNegOne {
		return second
	}
	return first
}

func maxList(numList ref.Val) ref.Val {
	l := numList.(traits.Lister)
	size := l.Size().(types.Int)
	if size == types.IntZero {
		return types.NewErr("math.@max(list) argument must not be empty")
	}
	max := l.Get(types.IntZero)
	for i := types.IntOne; i < size; i++ {
		max = maxPair(max, l.Get(i))
	}
	switch max.Type() {
	case types.IntType, types.DoubleType, types.UintType, types.UnknownType:
		return max
	default:
		return types.NewErr("no such overload: math.@max")
	}
}

func checkInvalidArgs(meh cel.MacroExprFactory, funcName string, args []ast.Expr) *cel.Error {
	for _, arg := range args {
		err := checkInvalidArgLiteral(funcName, arg)
		if err != nil {
			return meh.NewError(arg.ID(), err.Error())
		}
	}
	return nil
}

func checkInvalidArgLiteral(funcName string, arg ast.Expr) error {
	if !isNumericArgType(arg) {
		return fmt.Errorf("%s simple literal arguments must be numeric", funcName)
	}
	return nil
}

func isNumericArgType(arg ast.Expr) bool {
	switch arg.Kind() {
	case ast.LiteralKind:
		c := ref.Val(arg.AsLiteral())
		switch c.(type) {
		case types.Double, types.Int, types.Uint:
			return true
		default:
			return false
		}
	case ast.ListKind, ast.MapKind, ast.StructKind:
		return false
	default:
		return true
	}
}

func isListLiteralWithNumericArgs(arg ast.Expr) bool {
	switch arg.Kind() {
	case ast.ListKind:
		list := arg.AsList()
		if list.Size() == 0 {
			return false
		}
		for _, e := range list.Elements() {
			if !isNumericArgType(e) {
				return false
			}
		}
		return true
	}
	return false
}

func maybeSuffixError(val ref.Val, suffix string) ref.Val {
	if types.IsError(val) {
		msg := val.(*types.Err).String()
		if !strings.Contains(msg, suffix) {
			return types.NewErr("%s: %s", msg, suffix)
		}
	}
	return val
}
