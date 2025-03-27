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

// Package stdlib contains all of the standard library function declarations and definitions for CEL.
package stdlib

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

var (
	stdFunctions []*decls.FunctionDecl
	stdTypes     []*decls.VariableDecl
	utcTZ        = types.String("UTC")
)

func init() {
	paramA := types.NewTypeParamType("A")
	paramB := types.NewTypeParamType("B")
	listOfA := types.NewListType(paramA)
	mapOfAB := types.NewMapType(paramA, paramB)

	stdTypes = []*decls.VariableDecl{
		decls.TypeVariable(types.BoolType),
		decls.TypeVariable(types.BytesType),
		decls.TypeVariable(types.DoubleType),
		decls.TypeVariable(types.DurationType),
		decls.TypeVariable(types.IntType),
		decls.TypeVariable(listOfA),
		decls.TypeVariable(mapOfAB),
		decls.TypeVariable(types.NullType),
		decls.TypeVariable(types.StringType),
		decls.TypeVariable(types.TimestampType),
		decls.TypeVariable(types.TypeType),
		decls.TypeVariable(types.UintType),
	}

	stdFunctions = []*decls.FunctionDecl{
		// Logical operators. Special-cased within the interpreter.
		// Note, the singleton binding prevents extensions from overriding the operator behavior.
		function(operators.Conditional,
			decls.FunctionDoc(
				`The ternary operator tests a boolean predicate and returns the left-hand side `+
					`(truthy) expression if true, or the right-hand side (falsy) expression if false`),
			decls.Overload(overloads.Conditional, argTypes(types.BoolType, paramA, paramA), paramA,
				decls.OverloadIsNonStrict(),
				decls.OverloadDoc(
					`'hello'.contains('lo') ? 'hi' : 'bye' // 'hi'`,
					`32 % 3 == 0 ? 'divisible' : 'not divisible' // 'not divisible'`)),
			decls.SingletonFunctionBinding(noFunctionOverrides)),

		function(operators.LogicalAnd,
			decls.FunctionDoc(common.MultilineDescription{
				`logically AND two boolean values. Errors and unknown values`,
				`are valid inputs and will not halt evaluation.`}),
			decls.Overload(overloads.LogicalAnd, argTypes(types.BoolType, types.BoolType), types.BoolType,
				decls.OverloadIsNonStrict(),
				decls.OverloadDoc(
					`true && true   // true`,
					`true && false  // false`,
					`error && true  // error`,
					`error && false // false`)),
			decls.SingletonBinaryBinding(noBinaryOverrides)),

		function(operators.LogicalOr,
			decls.FunctionDoc(common.MultilineDescription{
				`logically OR two boolean values. Errors and unknown values`,
				`are valid inputs and will not halt evaluation.`}),
			decls.Overload(overloads.LogicalOr, argTypes(types.BoolType, types.BoolType), types.BoolType,
				decls.OverloadIsNonStrict(),
				decls.OverloadDoc(
					`true || false // true`,
					`false || false // false`,
					`error || true // true`,
					`error || error // true`)),
			decls.SingletonBinaryBinding(noBinaryOverrides)),

		function(operators.LogicalNot,
			decls.FunctionDoc(`logically negate a boolean value.`),
			decls.Overload(overloads.LogicalNot, argTypes(types.BoolType), types.BoolType,
				decls.OverloadDoc(
					`!true // false`,
					`!false // true`,
					`!error // error`)),
			decls.SingletonUnaryBinding(func(val ref.Val) ref.Val {
				b, ok := val.(types.Bool)
				if !ok {
					return types.MaybeNoSuchOverloadErr(val)
				}
				return b.Negate()
			})),

		// Comprehension short-circuiting related function
		function(operators.NotStrictlyFalse,
			decls.Overload(overloads.NotStrictlyFalse, argTypes(types.BoolType), types.BoolType,
				decls.OverloadIsNonStrict(),
				decls.UnaryBinding(notStrictlyFalse))),
		// Deprecated: __not_strictly_false__
		function(operators.OldNotStrictlyFalse,
			decls.DisableDeclaration(true), // safe deprecation
			decls.Overload(operators.OldNotStrictlyFalse, argTypes(types.BoolType), types.BoolType,
				decls.OverloadIsNonStrict(),
				decls.UnaryBinding(notStrictlyFalse))),

		// Equality / inequality. Special-cased in the interpreter
		function(operators.Equals,
			decls.FunctionDoc(`compare two values of the same type for equality`),
			decls.Overload(overloads.Equals, argTypes(paramA, paramA), types.BoolType,
				decls.OverloadDoc(
					`1 == 1 // true`,
					`'hello' == 'world' // false`,
					`bytes('hello') == b'hello' // true`,
					`duration('1h') == duration('60m') // true`,
					`dyn(3.0) == 3 // true`)),
			decls.SingletonBinaryBinding(noBinaryOverrides)),
		function(operators.NotEquals,
			decls.FunctionDoc(`compare two values of the same type for inequality`),
			decls.Overload(overloads.NotEquals, argTypes(paramA, paramA), types.BoolType,
				decls.OverloadDoc(
					`1 != 2     // true`,
					`"a" != "a" // false`,
					`3.0 != 3.1 // true`)),
			decls.SingletonBinaryBinding(noBinaryOverrides)),

		// Mathematical operators
		function(operators.Add,
			decls.FunctionDoc(common.MultilineDescription{
				`adds two numeric values or concatenates two strings, bytes,`,
				`or lists.`}),
			decls.Overload(overloads.AddBytes,
				argTypes(types.BytesType, types.BytesType), types.BytesType,
				decls.OverloadDoc(`b'hi' + bytes('ya') // b'hiya'`)),
			decls.Overload(overloads.AddDouble,
				argTypes(types.DoubleType, types.DoubleType), types.DoubleType,
				decls.OverloadDoc(`3.14 + 1.59 // 4.73`)),
			decls.Overload(overloads.AddDurationDuration,
				argTypes(types.DurationType, types.DurationType), types.DurationType,
				decls.OverloadDoc(`duration('1m') + duration('1s') // duration('1m1s')`)),
			decls.Overload(overloads.AddDurationTimestamp,
				argTypes(types.DurationType, types.TimestampType), types.TimestampType,
				decls.OverloadDoc(`duration('24h') + timestamp('2023-01-01T00:00:00Z') // timestamp('2023-01-02T00:00:00Z')`)),
			decls.Overload(overloads.AddTimestampDuration,
				argTypes(types.TimestampType, types.DurationType), types.TimestampType,
				decls.OverloadDoc(`timestamp('2023-01-01T00:00:00Z') + duration('24h1m2s') // timestamp('2023-01-02T00:01:02Z')`)),
			decls.Overload(overloads.AddInt64,
				argTypes(types.IntType, types.IntType), types.IntType,
				decls.OverloadDoc(`1 + 2 // 3`)),
			decls.Overload(overloads.AddList,
				argTypes(listOfA, listOfA), listOfA,
				decls.OverloadDoc(`[1] + [2, 3] // [1, 2, 3]`)),
			decls.Overload(overloads.AddString,
				argTypes(types.StringType, types.StringType), types.StringType,
				decls.OverloadDoc(`"Hello, " + "world!" // "Hello, world!"`)),
			decls.Overload(overloads.AddUint64,
				argTypes(types.UintType, types.UintType), types.UintType,
				decls.OverloadDoc(`22u + 33u // 55u`)),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Adder).Add(rhs)
			}, traits.AdderType)),
		function(operators.Divide,
			decls.FunctionDoc(`divide two numbers`),
			decls.Overload(overloads.DivideDouble,
				argTypes(types.DoubleType, types.DoubleType), types.DoubleType,
				decls.OverloadDoc(`7.0 / 2.0 // 3.5`)),
			decls.Overload(overloads.DivideInt64,
				argTypes(types.IntType, types.IntType), types.IntType,
				decls.OverloadDoc(`10 / 2 // 5`)),
			decls.Overload(overloads.DivideUint64,
				argTypes(types.UintType, types.UintType), types.UintType,
				decls.OverloadDoc(`42u / 2u // 21u`)),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Divider).Divide(rhs)
			}, traits.DividerType)),
		function(operators.Modulo,
			decls.FunctionDoc(`compute the modulus of one value into another`),
			decls.Overload(overloads.ModuloInt64,
				argTypes(types.IntType, types.IntType), types.IntType),
			decls.Overload(overloads.ModuloUint64,
				argTypes(types.UintType, types.UintType), types.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Modder).Modulo(rhs)
			}, traits.ModderType)),
		function(operators.Multiply,
			decls.FunctionDoc(`multiply two numbers`),
			decls.Overload(overloads.MultiplyDouble,
				argTypes(types.DoubleType, types.DoubleType), types.DoubleType),
			decls.Overload(overloads.MultiplyInt64,
				argTypes(types.IntType, types.IntType), types.IntType),
			decls.Overload(overloads.MultiplyUint64,
				argTypes(types.UintType, types.UintType), types.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Multiplier).Multiply(rhs)
			}, traits.MultiplierType)),
		function(operators.Negate,
			decls.FunctionDoc(`negate a numeric value`),
			decls.Overload(overloads.NegateDouble, argTypes(types.DoubleType), types.DoubleType),
			decls.Overload(overloads.NegateInt64, argTypes(types.IntType), types.IntType),
			decls.SingletonUnaryBinding(func(val ref.Val) ref.Val {
				if types.IsBool(val) {
					return types.MaybeNoSuchOverloadErr(val)
				}
				return val.(traits.Negater).Negate()
			}, traits.NegatorType)),
		function(operators.Subtract,
			decls.FunctionDoc(`subtract two numbers, or two time-related values`),
			decls.Overload(overloads.SubtractDouble,
				argTypes(types.DoubleType, types.DoubleType), types.DoubleType,
				decls.OverloadDoc(`10.5 - 2.0 // 8.5`)),
			decls.Overload(overloads.SubtractDurationDuration,
				argTypes(types.DurationType, types.DurationType), types.DurationType,
				decls.OverloadDoc(`duration('1m') - duration('1s') // duration('59s')`)),
			decls.Overload(overloads.SubtractInt64,
				argTypes(types.IntType, types.IntType), types.IntType,
				decls.OverloadDoc(`5 - 3 // 2`)),
			decls.Overload(overloads.SubtractTimestampDuration,
				argTypes(types.TimestampType, types.DurationType), types.TimestampType,
				decls.OverloadDoc(common.MultilineDescription{
					`timestamp('2023-01-10T12:00:00Z')`,
					`  - duration('12h') // timestamp('2023-01-10T00:00:00Z')`})),
			decls.Overload(overloads.SubtractTimestampTimestamp,
				argTypes(types.TimestampType, types.TimestampType), types.DurationType,
				decls.OverloadDoc(common.MultilineDescription{
					`timestamp('2023-01-10T12:00:00Z')`,
					`  - timestamp('2023-01-10T00:00:00Z') // duration('12h')`})),
			decls.Overload(overloads.SubtractUint64,
				argTypes(types.UintType, types.UintType), types.UintType,
				decls.OverloadDoc(common.MultilineDescription{
					`// the subtraction result must be positive, otherwise an overflow`,
					`// error is generated.`,
					`42u - 3u // 39u`})),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Subtractor).Subtract(rhs)
			}, traits.SubtractorType)),

		// Relations operators

		function(operators.Less,
			decls.FunctionDoc(`compare whether the first value is less than the second`),
			decls.Overload(overloads.LessBool,
				argTypes(types.BoolType, types.BoolType), types.BoolType),
			decls.Overload(overloads.LessInt64,
				argTypes(types.IntType, types.IntType), types.BoolType),
			decls.Overload(overloads.LessInt64Double,
				argTypes(types.IntType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.LessInt64Uint64,
				argTypes(types.IntType, types.UintType), types.BoolType),
			decls.Overload(overloads.LessUint64,
				argTypes(types.UintType, types.UintType), types.BoolType),
			decls.Overload(overloads.LessUint64Double,
				argTypes(types.UintType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.LessUint64Int64,
				argTypes(types.UintType, types.IntType), types.BoolType),
			decls.Overload(overloads.LessDouble,
				argTypes(types.DoubleType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.LessDoubleInt64,
				argTypes(types.DoubleType, types.IntType), types.BoolType),
			decls.Overload(overloads.LessDoubleUint64,
				argTypes(types.DoubleType, types.UintType), types.BoolType),
			decls.Overload(overloads.LessString,
				argTypes(types.StringType, types.StringType), types.BoolType),
			decls.Overload(overloads.LessBytes,
				argTypes(types.BytesType, types.BytesType), types.BoolType),
			decls.Overload(overloads.LessTimestamp,
				argTypes(types.TimestampType, types.TimestampType), types.BoolType),
			decls.Overload(overloads.LessDuration,
				argTypes(types.DurationType, types.DurationType), types.BoolType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntNegOne {
					return types.True
				}
				if cmp == types.IntOne || cmp == types.IntZero {
					return types.False
				}
				return cmp
			}, traits.ComparerType)),

		function(operators.LessEquals,
			decls.FunctionDoc(`compare whether the first value is less than or equal to the second`),
			decls.Overload(overloads.LessEqualsBool,
				argTypes(types.BoolType, types.BoolType), types.BoolType),
			decls.Overload(overloads.LessEqualsInt64,
				argTypes(types.IntType, types.IntType), types.BoolType),
			decls.Overload(overloads.LessEqualsInt64Double,
				argTypes(types.IntType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.LessEqualsInt64Uint64,
				argTypes(types.IntType, types.UintType), types.BoolType),
			decls.Overload(overloads.LessEqualsUint64,
				argTypes(types.UintType, types.UintType), types.BoolType),
			decls.Overload(overloads.LessEqualsUint64Double,
				argTypes(types.UintType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.LessEqualsUint64Int64,
				argTypes(types.UintType, types.IntType), types.BoolType),
			decls.Overload(overloads.LessEqualsDouble,
				argTypes(types.DoubleType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.LessEqualsDoubleInt64,
				argTypes(types.DoubleType, types.IntType), types.BoolType),
			decls.Overload(overloads.LessEqualsDoubleUint64,
				argTypes(types.DoubleType, types.UintType), types.BoolType),
			decls.Overload(overloads.LessEqualsString,
				argTypes(types.StringType, types.StringType), types.BoolType),
			decls.Overload(overloads.LessEqualsBytes,
				argTypes(types.BytesType, types.BytesType), types.BoolType),
			decls.Overload(overloads.LessEqualsTimestamp,
				argTypes(types.TimestampType, types.TimestampType), types.BoolType),
			decls.Overload(overloads.LessEqualsDuration,
				argTypes(types.DurationType, types.DurationType), types.BoolType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntNegOne || cmp == types.IntZero {
					return types.True
				}
				if cmp == types.IntOne {
					return types.False
				}
				return cmp
			}, traits.ComparerType)),

		function(operators.Greater,
			decls.FunctionDoc(`compare whether the first value is greater than to the second`),
			decls.Overload(overloads.GreaterBool,
				argTypes(types.BoolType, types.BoolType), types.BoolType),
			decls.Overload(overloads.GreaterInt64,
				argTypes(types.IntType, types.IntType), types.BoolType),
			decls.Overload(overloads.GreaterInt64Double,
				argTypes(types.IntType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.GreaterInt64Uint64,
				argTypes(types.IntType, types.UintType), types.BoolType),
			decls.Overload(overloads.GreaterUint64,
				argTypes(types.UintType, types.UintType), types.BoolType),
			decls.Overload(overloads.GreaterUint64Double,
				argTypes(types.UintType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.GreaterUint64Int64,
				argTypes(types.UintType, types.IntType), types.BoolType),
			decls.Overload(overloads.GreaterDouble,
				argTypes(types.DoubleType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.GreaterDoubleInt64,
				argTypes(types.DoubleType, types.IntType), types.BoolType),
			decls.Overload(overloads.GreaterDoubleUint64,
				argTypes(types.DoubleType, types.UintType), types.BoolType),
			decls.Overload(overloads.GreaterString,
				argTypes(types.StringType, types.StringType), types.BoolType),
			decls.Overload(overloads.GreaterBytes,
				argTypes(types.BytesType, types.BytesType), types.BoolType),
			decls.Overload(overloads.GreaterTimestamp,
				argTypes(types.TimestampType, types.TimestampType), types.BoolType),
			decls.Overload(overloads.GreaterDuration,
				argTypes(types.DurationType, types.DurationType), types.BoolType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntOne {
					return types.True
				}
				if cmp == types.IntNegOne || cmp == types.IntZero {
					return types.False
				}
				return cmp
			}, traits.ComparerType)),

		function(operators.GreaterEquals,
			decls.FunctionDoc(`compare whether the first value is greater than or equal to the second`),
			decls.Overload(overloads.GreaterEqualsBool,
				argTypes(types.BoolType, types.BoolType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsInt64,
				argTypes(types.IntType, types.IntType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsInt64Double,
				argTypes(types.IntType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsInt64Uint64,
				argTypes(types.IntType, types.UintType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsUint64,
				argTypes(types.UintType, types.UintType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsUint64Double,
				argTypes(types.UintType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsUint64Int64,
				argTypes(types.UintType, types.IntType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsDouble,
				argTypes(types.DoubleType, types.DoubleType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsDoubleInt64,
				argTypes(types.DoubleType, types.IntType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsDoubleUint64,
				argTypes(types.DoubleType, types.UintType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsString,
				argTypes(types.StringType, types.StringType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsBytes,
				argTypes(types.BytesType, types.BytesType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsTimestamp,
				argTypes(types.TimestampType, types.TimestampType), types.BoolType),
			decls.Overload(overloads.GreaterEqualsDuration,
				argTypes(types.DurationType, types.DurationType), types.BoolType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntOne || cmp == types.IntZero {
					return types.True
				}
				if cmp == types.IntNegOne {
					return types.False
				}
				return cmp
			}, traits.ComparerType)),

		// Indexing
		function(operators.Index,
			decls.FunctionDoc(`select a value from a list by index, or value from a map by key`),
			decls.Overload(overloads.IndexList, argTypes(listOfA, types.IntType), paramA),
			decls.Overload(overloads.IndexMap, argTypes(mapOfAB, paramA), paramB),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Indexer).Get(rhs)
			}, traits.IndexerType)),

		// Collections operators
		function(operators.In,
			decls.FunctionDoc(`test whether a value exists in a list, or a key exists in a map`),
			decls.Overload(overloads.InList, argTypes(paramA, listOfA), types.BoolType),
			decls.Overload(overloads.InMap, argTypes(paramA, mapOfAB), types.BoolType),
			decls.SingletonBinaryBinding(inAggregate)),
		function(operators.OldIn,
			decls.DisableDeclaration(true), // safe deprecation
			decls.Overload(overloads.InList, argTypes(paramA, listOfA), types.BoolType),
			decls.Overload(overloads.InMap, argTypes(paramA, mapOfAB), types.BoolType),
			decls.SingletonBinaryBinding(inAggregate)),
		function(overloads.DeprecatedIn,
			decls.DisableDeclaration(true), // safe deprecation
			decls.Overload(overloads.InList, argTypes(paramA, listOfA), types.BoolType),
			decls.Overload(overloads.InMap, argTypes(paramA, mapOfAB), types.BoolType),
			decls.SingletonBinaryBinding(inAggregate)),
		function(overloads.Size,
			decls.FunctionDoc(`compute the size of a list or map, the number of characters in a string, or the number of bytes in a sequence`),
			decls.Overload(overloads.SizeBytes, argTypes(types.BytesType), types.IntType),
			decls.MemberOverload(overloads.SizeBytesInst, argTypes(types.BytesType), types.IntType),
			decls.Overload(overloads.SizeList, argTypes(listOfA), types.IntType),
			decls.MemberOverload(overloads.SizeListInst, argTypes(listOfA), types.IntType),
			decls.Overload(overloads.SizeMap, argTypes(mapOfAB), types.IntType),
			decls.MemberOverload(overloads.SizeMapInst, argTypes(mapOfAB), types.IntType),
			decls.Overload(overloads.SizeString, argTypes(types.StringType), types.IntType),
			decls.MemberOverload(overloads.SizeStringInst, argTypes(types.StringType), types.IntType),
			decls.SingletonUnaryBinding(func(val ref.Val) ref.Val {
				return val.(traits.Sizer).Size()
			}, traits.SizerType)),

		// Type conversions
		function(overloads.TypeConvertType,
			decls.FunctionDoc(`convert a value to its type identifier`),
			decls.Overload(overloads.TypeConvertType, argTypes(paramA), types.NewTypeTypeWithParam(paramA)),
			decls.SingletonUnaryBinding(convertToType(types.TypeType))),

		// Bool conversions
		function(overloads.TypeConvertBool,
			decls.FunctionDoc(`convert a value to a boolean`),
			decls.Overload(overloads.BoolToBool, argTypes(types.BoolType), types.BoolType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.StringToBool, argTypes(types.StringType), types.BoolType,
				decls.UnaryBinding(convertToType(types.BoolType)))),

		// Bytes conversions
		function(overloads.TypeConvertBytes,
			decls.FunctionDoc(`convert a value to bytes`),
			decls.Overload(overloads.BytesToBytes, argTypes(types.BytesType), types.BytesType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.StringToBytes, argTypes(types.StringType), types.BytesType,
				decls.UnaryBinding(convertToType(types.BytesType)))),

		// Double conversions
		function(overloads.TypeConvertDouble,
			decls.FunctionDoc(`convert a value to a double`),
			decls.Overload(overloads.DoubleToDouble, argTypes(types.DoubleType), types.DoubleType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.IntToDouble, argTypes(types.IntType), types.DoubleType,
				decls.UnaryBinding(convertToType(types.DoubleType))),
			decls.Overload(overloads.StringToDouble, argTypes(types.StringType), types.DoubleType,
				decls.UnaryBinding(convertToType(types.DoubleType))),
			decls.Overload(overloads.UintToDouble, argTypes(types.UintType), types.DoubleType,
				decls.UnaryBinding(convertToType(types.DoubleType)))),

		// Duration conversions
		function(overloads.TypeConvertDuration,
			decls.FunctionDoc(`convert a value to a google.protobuf.Duration`),
			decls.Overload(overloads.DurationToDuration, argTypes(types.DurationType), types.DurationType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.IntToDuration, argTypes(types.IntType), types.DurationType,
				decls.UnaryBinding(convertToType(types.DurationType))),
			decls.Overload(overloads.StringToDuration, argTypes(types.StringType), types.DurationType,
				decls.UnaryBinding(convertToType(types.DurationType)))),

		// Dyn conversions
		function(overloads.TypeConvertDyn,
			decls.FunctionDoc(`indicate that the type is dynamic for type-checking purposes`),
			decls.Overload(overloads.ToDyn, argTypes(paramA), types.DynType),
			decls.SingletonUnaryBinding(identity)),

		// Int conversions
		function(overloads.TypeConvertInt,
			decls.FunctionDoc(`convert a value to an int`),
			decls.Overload(overloads.IntToInt, argTypes(types.IntType), types.IntType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.DoubleToInt, argTypes(types.DoubleType), types.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.DurationToInt, argTypes(types.DurationType), types.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.StringToInt, argTypes(types.StringType), types.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.TimestampToInt, argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.UintToInt, argTypes(types.UintType), types.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
		),

		// String conversions
		function(overloads.TypeConvertString,
			decls.FunctionDoc(`convert a value to a string`),
			decls.Overload(overloads.StringToString, argTypes(types.StringType), types.StringType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.BoolToString, argTypes(types.BoolType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.BytesToString, argTypes(types.BytesType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.DoubleToString, argTypes(types.DoubleType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.DurationToString, argTypes(types.DurationType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.IntToString, argTypes(types.IntType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.TimestampToString, argTypes(types.TimestampType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.UintToString, argTypes(types.UintType), types.StringType,
				decls.UnaryBinding(convertToType(types.StringType)))),

		// Timestamp conversions
		function(overloads.TypeConvertTimestamp,
			decls.FunctionDoc(`convert a value to a google.protobuf.Timestamp`),
			decls.Overload(overloads.TimestampToTimestamp, argTypes(types.TimestampType), types.TimestampType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.IntToTimestamp, argTypes(types.IntType), types.TimestampType,
				decls.UnaryBinding(convertToType(types.TimestampType))),
			decls.Overload(overloads.StringToTimestamp, argTypes(types.StringType), types.TimestampType,
				decls.UnaryBinding(convertToType(types.TimestampType)))),

		// Uint conversions
		function(overloads.TypeConvertUint,
			decls.FunctionDoc(`convert a value to a uint`),
			decls.Overload(overloads.UintToUint, argTypes(types.UintType), types.UintType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.DoubleToUint, argTypes(types.DoubleType), types.UintType,
				decls.UnaryBinding(convertToType(types.UintType))),
			decls.Overload(overloads.IntToUint, argTypes(types.IntType), types.UintType,
				decls.UnaryBinding(convertToType(types.UintType))),
			decls.Overload(overloads.StringToUint, argTypes(types.StringType), types.UintType,
				decls.UnaryBinding(convertToType(types.UintType)))),

		// String functions
		function(overloads.Contains,
			decls.FunctionDoc(`test whether a string contains a substring`),
			decls.MemberOverload(overloads.ContainsString,
				argTypes(types.StringType, types.StringType), types.BoolType,
				decls.BinaryBinding(types.StringContains)),
			decls.DisableTypeGuards(true)),
		function(overloads.EndsWith,
			decls.FunctionDoc(`test whether a string ends with a substring suffix`),
			decls.MemberOverload(overloads.EndsWithString,
				argTypes(types.StringType, types.StringType), types.BoolType,
				decls.BinaryBinding(types.StringEndsWith)),
			decls.DisableTypeGuards(true)),
		function(overloads.StartsWith,
			decls.FunctionDoc(`test whether a string starts with a substring prefix`),
			decls.MemberOverload(overloads.StartsWithString,
				argTypes(types.StringType, types.StringType), types.BoolType,
				decls.BinaryBinding(types.StringStartsWith)),
			decls.DisableTypeGuards(true)),
		function(overloads.Matches,
			decls.FunctionDoc(`test whether a string matches an RE2 regular expression`),
			decls.Overload(overloads.Matches, argTypes(types.StringType, types.StringType), types.BoolType),
			decls.MemberOverload(overloads.MatchesString,
				argTypes(types.StringType, types.StringType), types.BoolType),
			decls.SingletonBinaryBinding(func(str, pat ref.Val) ref.Val {
				return str.(traits.Matcher).Match(pat)
			}, traits.MatcherType)),

		// Timestamp / duration functions
		function(overloads.TimeGetFullYear,
			decls.FunctionDoc(`get the 0-based full year from a timestamp, UTC unless an IANA timezone is specified.`),
			decls.MemberOverload(overloads.TimestampToYear,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetFullYear(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToYearWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetFullYear))),

		function(overloads.TimeGetMonth,
			decls.FunctionDoc(`get the 0-based month from a timestamp, UTC unless an IANA timezone is specified.`),
			decls.MemberOverload(overloads.TimestampToMonth,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetMonth(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToMonthWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetMonth))),

		function(overloads.TimeGetDayOfYear,
			decls.FunctionDoc(`get the 0-based day of the year from a timestamp, UTC unless an IANA timezone is specified.`),
			decls.MemberOverload(overloads.TimestampToDayOfYear,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetDayOfYear(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToDayOfYearWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetDayOfYear))),

		function(overloads.TimeGetDayOfMonth,
			decls.FunctionDoc(`get the 0-based day of the month from a timestamp, UTC unless an IANA timezone is specified.`),
			decls.MemberOverload(overloads.TimestampToDayOfMonthZeroBased,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetDayOfMonthZeroBased(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetDayOfMonthZeroBased))),

		function(overloads.TimeGetDate,
			decls.FunctionDoc(`get the 1-based day of the month from a timestamp, UTC unless an IANA timezone is specified.`),
			decls.MemberOverload(overloads.TimestampToDayOfMonthOneBased,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetDayOfMonthOneBased(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetDayOfMonthOneBased))),

		function(overloads.TimeGetDayOfWeek,
			decls.FunctionDoc(`get the 0-based day of the week from a timestamp, UTC unless an IANA timezone is specified.`),
			decls.MemberOverload(overloads.TimestampToDayOfWeek,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetDayOfWeek(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToDayOfWeekWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetDayOfWeek))),

		function(overloads.TimeGetHours,
			decls.FunctionDoc(`get the hours portion from a timestamp, or convert a duration to hours`),
			decls.MemberOverload(overloads.TimestampToHours,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetHours(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToHoursWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetHours)),
			decls.MemberOverload(overloads.DurationToHours,
				argTypes(types.DurationType), types.IntType,
				decls.UnaryBinding(types.DurationGetHours))),

		function(overloads.TimeGetMinutes,
			decls.FunctionDoc(`get the minutes portion from a timestamp, or convert a duration to minutes`),
			decls.MemberOverload(overloads.TimestampToMinutes,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetMinutes(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToMinutesWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetMinutes)),
			decls.MemberOverload(overloads.DurationToMinutes,
				argTypes(types.DurationType), types.IntType,
				decls.UnaryBinding(types.DurationGetMinutes))),

		function(overloads.TimeGetSeconds,
			decls.FunctionDoc(`get the seconds portion from a timestamp, or convert a duration to seconds`),
			decls.MemberOverload(overloads.TimestampToSeconds,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetSeconds(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToSecondsWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetSeconds)),
			decls.MemberOverload(overloads.DurationToSeconds,
				argTypes(types.DurationType), types.IntType,
				decls.UnaryBinding(types.DurationGetSeconds))),

		function(overloads.TimeGetMilliseconds,
			decls.FunctionDoc(`get the milliseconds portion from a timestamp`),
			decls.MemberOverload(overloads.TimestampToMilliseconds,
				argTypes(types.TimestampType), types.IntType,
				decls.UnaryBinding(func(ts ref.Val) ref.Val {
					return timestampGetMilliseconds(ts, utcTZ)
				})),
			decls.MemberOverload(overloads.TimestampToMillisecondsWithTz,
				argTypes(types.TimestampType, types.StringType), types.IntType,
				decls.BinaryBinding(timestampGetMilliseconds)),
			decls.MemberOverload(overloads.DurationToMilliseconds,
				argTypes(types.DurationType), types.IntType,
				decls.UnaryBinding(types.DurationGetMilliseconds))),
	}
}

// Functions returns the set of standard library function declarations and definitions for CEL.
func Functions() []*decls.FunctionDecl {
	return stdFunctions
}

// Types returns the set of standard library types for CEL.
func Types() []*decls.VariableDecl {
	return stdTypes
}

func notStrictlyFalse(value ref.Val) ref.Val {
	if types.IsBool(value) {
		return value
	}
	return types.True
}

func inAggregate(lhs ref.Val, rhs ref.Val) ref.Val {
	if rhs.Type().HasTrait(traits.ContainerType) {
		return rhs.(traits.Container).Contains(lhs)
	}
	return types.ValOrErr(rhs, "no such overload")
}

func function(name string, opts ...decls.FunctionOpt) *decls.FunctionDecl {
	fn, err := decls.NewFunction(name, opts...)
	if err != nil {
		panic(err)
	}
	return fn
}

func argTypes(args ...*types.Type) []*types.Type {
	return args
}

func noBinaryOverrides(rhs, lhs ref.Val) ref.Val {
	return types.NoSuchOverloadErr()
}

func noFunctionOverrides(args ...ref.Val) ref.Val {
	return types.NoSuchOverloadErr()
}

func identity(val ref.Val) ref.Val {
	return val
}

func convertToType(t ref.Type) functions.UnaryOp {
	return func(val ref.Val) ref.Val {
		return val.ConvertToType(t)
	}
}

func timestampGetFullYear(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Year())
}

func timestampGetMonth(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	// CEL spec indicates that the month should be 0-based, but the Time value
	// for Month() is 1-based.
	return types.Int(t.Month() - 1)
}

func timestampGetDayOfYear(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.YearDay() - 1)
}

func timestampGetDayOfMonthZeroBased(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Day() - 1)
}

func timestampGetDayOfMonthOneBased(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Day())
}

func timestampGetDayOfWeek(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Weekday())
}

func timestampGetHours(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Hour())
}

func timestampGetMinutes(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Minute())
}

func timestampGetSeconds(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Second())
}

func timestampGetMilliseconds(ts, tz ref.Val) ref.Val {
	t, err := inTimeZone(ts, tz)
	if err != nil {
		return types.NewErrFromString(err.Error())
	}
	return types.Int(t.Nanosecond() / 1000000)
}

func inTimeZone(ts, tz ref.Val) (time.Time, error) {
	t := ts.(types.Timestamp)
	val := string(tz.(types.String))
	ind := strings.Index(val, ":")
	if ind == -1 {
		loc, err := time.LoadLocation(val)
		if err != nil {
			return time.Time{}, err
		}
		return t.In(loc), nil
	}

	// If the input is not the name of a timezone (for example, 'US/Central'), it should be a numerical offset from UTC
	// in the format ^(+|-)(0[0-9]|1[0-4]):[0-5][0-9]$. The numerical input is parsed in terms of hours and minutes.
	hr, err := strconv.Atoi(string(val[0:ind]))
	if err != nil {
		return time.Time{}, err
	}
	min, err := strconv.Atoi(string(val[ind+1:]))
	if err != nil {
		return time.Time{}, err
	}
	var offset int
	if string(val[0]) == "-" {
		offset = hr*60 - min
	} else {
		offset = hr*60 + min
	}
	secondsEastOfUTC := int((time.Duration(offset) * time.Minute).Seconds())
	timezone := time.FixedZone("", secondsEastOfUTC)
	return t.In(timezone), nil
}
