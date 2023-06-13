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
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	stdFunctions []*decls.FunctionDecl
	stdFnDecls   []*exprpb.Decl
	stdTypes     []*decls.VariableDecl
	stdTypeDecls []*exprpb.Decl
)

func init() {
	paramA := decls.TypeParamType("A")
	paramB := decls.TypeParamType("B")
	listOfA := decls.ListType(paramA)
	mapOfAB := decls.MapType(paramA, paramB)

	stdTypes = []*decls.VariableDecl{
		decls.BoolType.TypeVariable(),
		decls.BytesType.TypeVariable(),
		decls.DoubleType.TypeVariable(),
		decls.DurationType.TypeVariable(),
		decls.IntType.TypeVariable(),
		listOfA.TypeVariable(),
		mapOfAB.TypeVariable(),
		decls.NullType.TypeVariable(),
		decls.StringType.TypeVariable(),
		decls.TimestampType.TypeVariable(),
		decls.TypeType.TypeVariable(),
		decls.UintType.TypeVariable(),
	}

	stdTypeDecls = make([]*exprpb.Decl, 0, len(stdTypes))
	for _, stdType := range stdTypes {
		typeVar, err := decls.VariableDeclToExprDecl(stdType)
		if err != nil {
			panic(err)
		}
		stdTypeDecls = append(stdTypeDecls, typeVar)
	}

	stdFunctions = []*decls.FunctionDecl{
		// Logical operators. Special-cased within the interpreter.
		// Note, the singleton binding prevents extensions from overriding the operator behavior.
		function(operators.Conditional,
			decls.Overload(overloads.Conditional, argTypes(decls.BoolType, paramA, paramA), paramA,
				decls.OverloadIsNonStrict()),
			decls.SingletonFunctionBinding(noFunctionOverrides)),
		function(operators.LogicalAnd,
			decls.Overload(overloads.LogicalAnd, argTypes(decls.BoolType, decls.BoolType), decls.BoolType,
				decls.OverloadIsNonStrict()),
			decls.SingletonBinaryBinding(noBinaryOverrides)),
		function(operators.LogicalOr,
			decls.Overload(overloads.LogicalOr, argTypes(decls.BoolType, decls.BoolType), decls.BoolType,
				decls.OverloadIsNonStrict()),
			decls.SingletonBinaryBinding(noBinaryOverrides)),
		function(operators.LogicalNot,
			decls.Overload(overloads.LogicalNot, argTypes(decls.BoolType), decls.BoolType),
			decls.SingletonUnaryBinding(func(val ref.Val) ref.Val {
				b, ok := val.(types.Bool)
				if !ok {
					return types.MaybeNoSuchOverloadErr(val)
				}
				return b.Negate()
			})),

		// Comprehension short-circuiting related function
		function(operators.NotStrictlyFalse,
			decls.Overload(overloads.NotStrictlyFalse, argTypes(decls.BoolType), decls.BoolType,
				decls.OverloadIsNonStrict(),
				decls.UnaryBinding(notStrictlyFalse))),
		// Deprecated: __not_strictly_false__
		function(operators.OldNotStrictlyFalse,
			decls.DisableDeclaration(true), // safe deprecation
			decls.Overload(operators.OldNotStrictlyFalse, argTypes(decls.BoolType), decls.BoolType,
				decls.OverloadIsNonStrict(),
				decls.UnaryBinding(notStrictlyFalse))),

		// Equality / inequality. Special-cased in the interpreter
		function(operators.Equals,
			decls.Overload(overloads.Equals, argTypes(paramA, paramA), decls.BoolType),
			decls.SingletonBinaryBinding(noBinaryOverrides)),
		function(operators.NotEquals,
			decls.Overload(overloads.NotEquals, argTypes(paramA, paramA), decls.BoolType),
			decls.SingletonBinaryBinding(noBinaryOverrides)),

		// Mathematical operators
		function(operators.Add,
			decls.Overload(overloads.AddBytes,
				argTypes(decls.BytesType, decls.BytesType), decls.BytesType),
			decls.Overload(overloads.AddDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.DoubleType),
			decls.Overload(overloads.AddDurationDuration,
				argTypes(decls.DurationType, decls.DurationType), decls.DurationType),
			decls.Overload(overloads.AddDurationTimestamp,
				argTypes(decls.DurationType, decls.TimestampType), decls.TimestampType),
			decls.Overload(overloads.AddTimestampDuration,
				argTypes(decls.TimestampType, decls.DurationType), decls.TimestampType),
			decls.Overload(overloads.AddInt64,
				argTypes(decls.IntType, decls.IntType), decls.IntType),
			decls.Overload(overloads.AddList,
				argTypes(listOfA, listOfA), listOfA),
			decls.Overload(overloads.AddString,
				argTypes(decls.StringType, decls.StringType), decls.StringType),
			decls.Overload(overloads.AddUint64,
				argTypes(decls.UintType, decls.UintType), decls.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Adder).Add(rhs)
			}, traits.AdderType)),
		function(operators.Divide,
			decls.Overload(overloads.DivideDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.DoubleType),
			decls.Overload(overloads.DivideInt64,
				argTypes(decls.IntType, decls.IntType), decls.IntType),
			decls.Overload(overloads.DivideUint64,
				argTypes(decls.UintType, decls.UintType), decls.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Divider).Divide(rhs)
			}, traits.DividerType)),
		function(operators.Modulo,
			decls.Overload(overloads.ModuloInt64,
				argTypes(decls.IntType, decls.IntType), decls.IntType),
			decls.Overload(overloads.ModuloUint64,
				argTypes(decls.UintType, decls.UintType), decls.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Modder).Modulo(rhs)
			}, traits.ModderType)),
		function(operators.Multiply,
			decls.Overload(overloads.MultiplyDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.DoubleType),
			decls.Overload(overloads.MultiplyInt64,
				argTypes(decls.IntType, decls.IntType), decls.IntType),
			decls.Overload(overloads.MultiplyUint64,
				argTypes(decls.UintType, decls.UintType), decls.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Multiplier).Multiply(rhs)
			}, traits.MultiplierType)),
		function(operators.Negate,
			decls.Overload(overloads.NegateDouble, argTypes(decls.DoubleType), decls.DoubleType),
			decls.Overload(overloads.NegateInt64, argTypes(decls.IntType), decls.IntType),
			decls.SingletonUnaryBinding(func(val ref.Val) ref.Val {
				if types.IsBool(val) {
					return types.MaybeNoSuchOverloadErr(val)
				}
				return val.(traits.Negater).Negate()
			}, traits.NegatorType)),
		function(operators.Subtract,
			decls.Overload(overloads.SubtractDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.DoubleType),
			decls.Overload(overloads.SubtractDurationDuration,
				argTypes(decls.DurationType, decls.DurationType), decls.DurationType),
			decls.Overload(overloads.SubtractInt64,
				argTypes(decls.IntType, decls.IntType), decls.IntType),
			decls.Overload(overloads.SubtractTimestampDuration,
				argTypes(decls.TimestampType, decls.DurationType), decls.TimestampType),
			decls.Overload(overloads.SubtractTimestampTimestamp,
				argTypes(decls.TimestampType, decls.TimestampType), decls.DurationType),
			decls.Overload(overloads.SubtractUint64,
				argTypes(decls.UintType, decls.UintType), decls.UintType),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Subtractor).Subtract(rhs)
			}, traits.SubtractorType)),

		// Relations operators

		function(operators.Less,
			decls.Overload(overloads.LessBool,
				argTypes(decls.BoolType, decls.BoolType), decls.BoolType),
			decls.Overload(overloads.LessInt64,
				argTypes(decls.IntType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.LessInt64Double,
				argTypes(decls.IntType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.LessInt64Uint64,
				argTypes(decls.IntType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.LessUint64,
				argTypes(decls.UintType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.LessUint64Double,
				argTypes(decls.UintType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.LessUint64Int64,
				argTypes(decls.UintType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.LessDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.LessDoubleInt64,
				argTypes(decls.DoubleType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.LessDoubleUint64,
				argTypes(decls.DoubleType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.LessString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType),
			decls.Overload(overloads.LessBytes,
				argTypes(decls.BytesType, decls.BytesType), decls.BoolType),
			decls.Overload(overloads.LessTimestamp,
				argTypes(decls.TimestampType, decls.TimestampType), decls.BoolType),
			decls.Overload(overloads.LessDuration,
				argTypes(decls.DurationType, decls.DurationType), decls.BoolType),
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
			decls.Overload(overloads.LessEqualsBool,
				argTypes(decls.BoolType, decls.BoolType), decls.BoolType),
			decls.Overload(overloads.LessEqualsInt64,
				argTypes(decls.IntType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.LessEqualsInt64Double,
				argTypes(decls.IntType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.LessEqualsInt64Uint64,
				argTypes(decls.IntType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.LessEqualsUint64,
				argTypes(decls.UintType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.LessEqualsUint64Double,
				argTypes(decls.UintType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.LessEqualsUint64Int64,
				argTypes(decls.UintType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.LessEqualsDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.LessEqualsDoubleInt64,
				argTypes(decls.DoubleType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.LessEqualsDoubleUint64,
				argTypes(decls.DoubleType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.LessEqualsString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType),
			decls.Overload(overloads.LessEqualsBytes,
				argTypes(decls.BytesType, decls.BytesType), decls.BoolType),
			decls.Overload(overloads.LessEqualsTimestamp,
				argTypes(decls.TimestampType, decls.TimestampType), decls.BoolType),
			decls.Overload(overloads.LessEqualsDuration,
				argTypes(decls.DurationType, decls.DurationType), decls.BoolType),
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
			decls.Overload(overloads.GreaterBool,
				argTypes(decls.BoolType, decls.BoolType), decls.BoolType),
			decls.Overload(overloads.GreaterInt64,
				argTypes(decls.IntType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.GreaterInt64Double,
				argTypes(decls.IntType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.GreaterInt64Uint64,
				argTypes(decls.IntType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.GreaterUint64,
				argTypes(decls.UintType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.GreaterUint64Double,
				argTypes(decls.UintType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.GreaterUint64Int64,
				argTypes(decls.UintType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.GreaterDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.GreaterDoubleInt64,
				argTypes(decls.DoubleType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.GreaterDoubleUint64,
				argTypes(decls.DoubleType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.GreaterString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType),
			decls.Overload(overloads.GreaterBytes,
				argTypes(decls.BytesType, decls.BytesType), decls.BoolType),
			decls.Overload(overloads.GreaterTimestamp,
				argTypes(decls.TimestampType, decls.TimestampType), decls.BoolType),
			decls.Overload(overloads.GreaterDuration,
				argTypes(decls.DurationType, decls.DurationType), decls.BoolType),
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
			decls.Overload(overloads.GreaterEqualsBool,
				argTypes(decls.BoolType, decls.BoolType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsInt64,
				argTypes(decls.IntType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsInt64Double,
				argTypes(decls.IntType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsInt64Uint64,
				argTypes(decls.IntType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsUint64,
				argTypes(decls.UintType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsUint64Double,
				argTypes(decls.UintType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsUint64Int64,
				argTypes(decls.UintType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsDouble,
				argTypes(decls.DoubleType, decls.DoubleType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsDoubleInt64,
				argTypes(decls.DoubleType, decls.IntType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsDoubleUint64,
				argTypes(decls.DoubleType, decls.UintType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsBytes,
				argTypes(decls.BytesType, decls.BytesType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsTimestamp,
				argTypes(decls.TimestampType, decls.TimestampType), decls.BoolType),
			decls.Overload(overloads.GreaterEqualsDuration,
				argTypes(decls.DurationType, decls.DurationType), decls.BoolType),
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
			decls.Overload(overloads.IndexList, argTypes(listOfA, decls.IntType), paramA),
			decls.Overload(overloads.IndexMap, argTypes(mapOfAB, paramA), paramB),
			decls.SingletonBinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.(traits.Indexer).Get(rhs)
			}, traits.IndexerType)),

		// Collections operators
		function(operators.In,
			decls.Overload(overloads.InList, argTypes(paramA, listOfA), decls.BoolType),
			decls.Overload(overloads.InMap, argTypes(paramA, mapOfAB), decls.BoolType),
			decls.SingletonBinaryBinding(inAggregate)),
		function(operators.OldIn,
			decls.DisableDeclaration(true), // safe deprecation
			decls.Overload(overloads.InList, argTypes(paramA, listOfA), decls.BoolType),
			decls.Overload(overloads.InMap, argTypes(paramA, mapOfAB), decls.BoolType),
			decls.SingletonBinaryBinding(inAggregate)),
		function(overloads.DeprecatedIn,
			decls.DisableDeclaration(true), // safe deprecation
			decls.Overload(overloads.InList, argTypes(paramA, listOfA), decls.BoolType),
			decls.Overload(overloads.InMap, argTypes(paramA, mapOfAB), decls.BoolType),
			decls.SingletonBinaryBinding(inAggregate)),
		function(overloads.Size,
			decls.Overload(overloads.SizeBytes, argTypes(decls.BytesType), decls.IntType),
			decls.MemberOverload(overloads.SizeBytesInst, argTypes(decls.BytesType), decls.IntType),
			decls.Overload(overloads.SizeList, argTypes(listOfA), decls.IntType),
			decls.MemberOverload(overloads.SizeListInst, argTypes(listOfA), decls.IntType),
			decls.Overload(overloads.SizeMap, argTypes(mapOfAB), decls.IntType),
			decls.MemberOverload(overloads.SizeMapInst, argTypes(mapOfAB), decls.IntType),
			decls.Overload(overloads.SizeString, argTypes(decls.StringType), decls.IntType),
			decls.MemberOverload(overloads.SizeStringInst, argTypes(decls.StringType), decls.IntType),
			decls.SingletonUnaryBinding(func(val ref.Val) ref.Val {
				return val.(traits.Sizer).Size()
			}, traits.SizerType)),

		// Type conversions
		function(overloads.TypeConvertType,
			decls.Overload(overloads.TypeConvertType, argTypes(paramA), decls.TypeTypeWithParam(paramA)),
			decls.SingletonUnaryBinding(convertToType(types.TypeType))),

		// Bool conversions
		function(overloads.TypeConvertBool,
			decls.Overload(overloads.BoolToBool, argTypes(decls.BoolType), decls.BoolType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.StringToBool, argTypes(decls.StringType), decls.BoolType,
				decls.UnaryBinding(convertToType(types.BoolType)))),

		// Bytes conversions
		function(overloads.TypeConvertBytes,
			decls.Overload(overloads.BytesToBytes, argTypes(decls.BytesType), decls.BytesType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.StringToBytes, argTypes(decls.StringType), decls.BytesType,
				decls.UnaryBinding(convertToType(types.BytesType)))),

		// Double conversions
		function(overloads.TypeConvertDouble,
			decls.Overload(overloads.DoubleToDouble, argTypes(decls.DoubleType), decls.DoubleType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.IntToDouble, argTypes(decls.IntType), decls.DoubleType,
				decls.UnaryBinding(convertToType(types.DoubleType))),
			decls.Overload(overloads.StringToDouble, argTypes(decls.StringType), decls.DoubleType,
				decls.UnaryBinding(convertToType(types.DoubleType))),
			decls.Overload(overloads.UintToDouble, argTypes(decls.UintType), decls.DoubleType,
				decls.UnaryBinding(convertToType(types.DoubleType)))),

		// Duration conversions
		function(overloads.TypeConvertDuration,
			decls.Overload(overloads.DurationToDuration, argTypes(decls.DurationType), decls.DurationType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.IntToDuration, argTypes(decls.IntType), decls.DurationType,
				decls.UnaryBinding(convertToType(types.DurationType))),
			decls.Overload(overloads.StringToDuration, argTypes(decls.StringType), decls.DurationType,
				decls.UnaryBinding(convertToType(types.DurationType)))),

		// Dyn conversions
		function(overloads.TypeConvertDyn,
			decls.Overload(overloads.ToDyn, argTypes(paramA), decls.DynType),
			decls.SingletonUnaryBinding(identity)),

		// Int conversions
		function(overloads.TypeConvertInt,
			decls.Overload(overloads.IntToInt, argTypes(decls.IntType), decls.IntType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.DoubleToInt, argTypes(decls.DoubleType), decls.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.DurationToInt, argTypes(decls.DurationType), decls.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.StringToInt, argTypes(decls.StringType), decls.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.TimestampToInt, argTypes(decls.TimestampType), decls.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
			decls.Overload(overloads.UintToInt, argTypes(decls.UintType), decls.IntType,
				decls.UnaryBinding(convertToType(types.IntType))),
		),

		// String conversions
		function(overloads.TypeConvertString,
			decls.Overload(overloads.StringToString, argTypes(decls.StringType), decls.StringType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.BoolToString, argTypes(decls.BoolType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.BytesToString, argTypes(decls.BytesType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.DoubleToString, argTypes(decls.DoubleType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.DurationToString, argTypes(decls.DurationType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.IntToString, argTypes(decls.IntType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.TimestampToString, argTypes(decls.TimestampType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType))),
			decls.Overload(overloads.UintToString, argTypes(decls.UintType), decls.StringType,
				decls.UnaryBinding(convertToType(types.StringType)))),

		// Timestamp conversions
		function(overloads.TypeConvertTimestamp,
			decls.Overload(overloads.TimestampToTimestamp, argTypes(decls.TimestampType), decls.TimestampType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.IntToTimestamp, argTypes(decls.IntType), decls.TimestampType,
				decls.UnaryBinding(convertToType(types.TimestampType))),
			decls.Overload(overloads.StringToTimestamp, argTypes(decls.StringType), decls.TimestampType,
				decls.UnaryBinding(convertToType(types.TimestampType)))),

		// Uint conversions
		function(overloads.TypeConvertUint,
			decls.Overload(overloads.UintToUint, argTypes(decls.UintType), decls.UintType,
				decls.UnaryBinding(identity)),
			decls.Overload(overloads.DoubleToUint, argTypes(decls.DoubleType), decls.UintType,
				decls.UnaryBinding(convertToType(types.UintType))),
			decls.Overload(overloads.IntToUint, argTypes(decls.IntType), decls.UintType,
				decls.UnaryBinding(convertToType(types.UintType))),
			decls.Overload(overloads.StringToUint, argTypes(decls.StringType), decls.UintType,
				decls.UnaryBinding(convertToType(types.UintType)))),

		// String functions
		function(overloads.Contains,
			decls.MemberOverload(overloads.ContainsString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType,
				decls.BinaryBinding(types.StringContains)),
			decls.DisableTypeGuards(true)),
		function(overloads.EndsWith,
			decls.MemberOverload(overloads.EndsWithString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType,
				decls.BinaryBinding(types.StringEndsWith)),
			decls.DisableTypeGuards(true)),
		function(overloads.StartsWith,
			decls.MemberOverload(overloads.StartsWithString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType,
				decls.BinaryBinding(types.StringStartsWith)),
			decls.DisableTypeGuards(true)),
		function(overloads.Matches,
			decls.Overload(overloads.Matches, argTypes(decls.StringType, decls.StringType), decls.BoolType),
			decls.MemberOverload(overloads.MatchesString,
				argTypes(decls.StringType, decls.StringType), decls.BoolType),
			decls.SingletonBinaryBinding(func(str, pat ref.Val) ref.Val {
				return str.(traits.Matcher).Match(pat)
			}, traits.MatcherType)),

		// Timestamp / duration functions
		function(overloads.TimeGetFullYear,
			decls.MemberOverload(overloads.TimestampToYear,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToYearWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType)),

		function(overloads.TimeGetMonth,
			decls.MemberOverload(overloads.TimestampToMonth,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToMonthWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType)),

		function(overloads.TimeGetDayOfYear,
			decls.MemberOverload(overloads.TimestampToDayOfYear,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToDayOfYearWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType)),

		function(overloads.TimeGetDayOfMonth,
			decls.MemberOverload(overloads.TimestampToDayOfMonthZeroBased,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType)),

		function(overloads.TimeGetDate,
			decls.MemberOverload(overloads.TimestampToDayOfMonthOneBased,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType)),

		function(overloads.TimeGetDayOfWeek,
			decls.MemberOverload(overloads.TimestampToDayOfWeek,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToDayOfWeekWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType)),

		function(overloads.TimeGetHours,
			decls.MemberOverload(overloads.TimestampToHours,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToHoursWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType),
			decls.MemberOverload(overloads.DurationToHours,
				argTypes(decls.DurationType), decls.IntType)),

		function(overloads.TimeGetMinutes,
			decls.MemberOverload(overloads.TimestampToMinutes,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToMinutesWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType),
			decls.MemberOverload(overloads.DurationToMinutes,
				argTypes(decls.DurationType), decls.IntType)),

		function(overloads.TimeGetSeconds,
			decls.MemberOverload(overloads.TimestampToSeconds,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToSecondsWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType),
			decls.MemberOverload(overloads.DurationToSeconds,
				argTypes(decls.DurationType), decls.IntType)),

		function(overloads.TimeGetMilliseconds,
			decls.MemberOverload(overloads.TimestampToMilliseconds,
				argTypes(decls.TimestampType), decls.IntType),
			decls.MemberOverload(overloads.TimestampToMillisecondsWithTz,
				argTypes(decls.TimestampType, decls.StringType), decls.IntType),
			decls.MemberOverload(overloads.DurationToMilliseconds,
				argTypes(decls.DurationType), decls.IntType)),
	}

	stdFnDecls = make([]*exprpb.Decl, 0, len(stdFunctions))
	for _, fn := range stdFunctions {
		if fn.IsDeclarationDisabled() {
			continue
		}
		ed, err := decls.FunctionDeclToExprDecl(fn)
		if err != nil {
			panic(err)
		}
		stdFnDecls = append(stdFnDecls, ed)
	}
}

// Functions returns the set of standard library function declarations and definitions for CEL.
func Functions() []*decls.FunctionDecl {
	return stdFunctions
}

// FunctionExprDecls returns the legacy style protobuf-typed declarations for all functions and overloads
// in the CEL standard environment.
func FunctionExprDecls() []*exprpb.Decl {
	return stdFnDecls
}

// Types returns the set of standard library types for CEL.
func Types() []*decls.VariableDecl {
	return stdTypes
}

// TypeExprDecls returns the legacy style protobuf-typed declarations for all types in the CEL
// standard environment.
func TypeExprDecls() []*exprpb.Decl {
	return stdTypeDecls
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

func argTypes(args ...*decls.Type) []*decls.Type {
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
