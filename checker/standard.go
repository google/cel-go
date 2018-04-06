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

package checker

import (
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/checker/types"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-spec/proto/checked/v1/checked"
)

func StandardDeclarations() []*checked.Decl {
	// Some shortcuts we use when building declarations.
	paramA := types.NewTypeParam("A")
	typeParamAList := []string{"A"}
	listOfA := types.NewList(paramA)
	paramB := types.NewTypeParam("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := types.NewMap(paramA, paramB)

	var idents []*checked.Decl
	for _, t := range []*checked.Type{types.Int64, types.Uint64, types.Bool, types.Double, types.Bytes, types.String} {
		idents = append(idents, decls.NewIdent(types.FormatType(t), types.NewType(t), nil))
	}
	idents = append(idents,
		decls.NewIdent("list", types.NewType(listOfA), nil),
		decls.NewIdent("map", types.NewType(mapOfAB), nil))

	// Booleans
	// TODO: allow the conditional to return a heterogenous type.
	return append(idents, []*checked.Decl{
		decls.NewFunction(operators.Conditional,
			decls.NewParameterizedOverload(overloads.Conditional,
				[]*checked.Type{types.Bool, paramA, paramA}, paramA,
				typeParamAList)),

		decls.NewFunction(operators.LogicalAnd,
			decls.NewOverload(overloads.LogicalAnd,
				[]*checked.Type{types.Bool, types.Bool}, types.Bool)),

		decls.NewFunction(operators.LogicalOr,
			decls.NewOverload(overloads.LogicalOr,
				[]*checked.Type{types.Bool, types.Bool}, types.Bool)),

		decls.NewFunction(operators.LogicalNot,
			decls.NewOverload(overloads.LogicalNot,
				[]*checked.Type{types.Bool}, types.Bool)),

		decls.NewFunction(overloads.Matches,
			decls.NewInstanceOverload(overloads.MatchString,
				[]*checked.Type{types.String, types.String}, types.Bool)),

		// Relations

		decls.NewFunction(operators.Less,
			decls.NewOverload(overloads.LessBool,
				[]*checked.Type{types.Bool, types.Bool}, types.Bool),
			decls.NewOverload(overloads.LessInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Bool),
			decls.NewOverload(overloads.LessUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Bool),
			decls.NewOverload(overloads.LessDouble,
				[]*checked.Type{types.Double, types.Double}, types.Bool),
			decls.NewOverload(overloads.LessString,
				[]*checked.Type{types.String, types.String}, types.Bool),
			decls.NewOverload(overloads.LessBytes,
				[]*checked.Type{types.Bytes, types.Bytes}, types.Bool),
			decls.NewOverload(overloads.LessTimestamp,
				[]*checked.Type{types.Timestamp, types.Timestamp}, types.Bool),
			decls.NewOverload(overloads.LessDuration,
				[]*checked.Type{types.Duration, types.Duration}, types.Bool)),

		decls.NewFunction(operators.LessEquals,
			decls.NewOverload(overloads.LessEqualsBool,
				[]*checked.Type{types.Bool, types.Bool}, types.Bool),
			decls.NewOverload(overloads.LessEqualsInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Bool),
			decls.NewOverload(overloads.LessEqualsUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Bool),
			decls.NewOverload(overloads.LessEqualsDouble,
				[]*checked.Type{types.Double, types.Double}, types.Bool),
			decls.NewOverload(overloads.LessEqualsString,
				[]*checked.Type{types.String, types.String}, types.Bool),
			decls.NewOverload(overloads.LessEqualsBytes,
				[]*checked.Type{types.Bytes, types.Bytes}, types.Bool),
			decls.NewOverload(overloads.LessEqualsTimestamp,
				[]*checked.Type{types.Timestamp, types.Timestamp}, types.Bool),
			decls.NewOverload(overloads.LessEqualsDuration,
				[]*checked.Type{types.Duration, types.Duration}, types.Bool)),

		decls.NewFunction(operators.Greater,
			decls.NewOverload(overloads.GreaterBool,
				[]*checked.Type{types.Bool, types.Bool}, types.Bool),
			decls.NewOverload(overloads.GreaterInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Bool),
			decls.NewOverload(overloads.GreaterUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Bool),
			decls.NewOverload(overloads.GreaterDouble,
				[]*checked.Type{types.Double, types.Double}, types.Bool),
			decls.NewOverload(overloads.GreaterString,
				[]*checked.Type{types.String, types.String}, types.Bool),
			decls.NewOverload(overloads.GreaterBytes,
				[]*checked.Type{types.Bytes, types.Bytes}, types.Bool),
			decls.NewOverload(overloads.GreaterTimestamp,
				[]*checked.Type{types.Timestamp, types.Timestamp}, types.Bool),
			decls.NewOverload(overloads.GreaterDuration,
				[]*checked.Type{types.Duration, types.Duration}, types.Bool)),

		decls.NewFunction(operators.GreaterEquals,
			decls.NewOverload(overloads.GreaterEqualsBool,
				[]*checked.Type{types.Bool, types.Bool}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsDouble,
				[]*checked.Type{types.Double, types.Double}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsString,
				[]*checked.Type{types.String, types.String}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsBytes,
				[]*checked.Type{types.Bytes, types.Bytes}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsTimestamp,
				[]*checked.Type{types.Timestamp, types.Timestamp}, types.Bool),
			decls.NewOverload(overloads.GreaterEqualsDuration,
				[]*checked.Type{types.Duration, types.Duration}, types.Bool)),

		decls.NewFunction(operators.Equals,
			decls.NewParameterizedOverload(overloads.Equals,
				[]*checked.Type{paramA, paramA}, types.Bool,
				typeParamAList)),

		decls.NewFunction(operators.NotEquals,
			decls.NewParameterizedOverload(overloads.NotEquals,
				[]*checked.Type{paramA, paramA}, types.Bool,
				typeParamAList)),

		// Algebra

		decls.NewFunction(operators.Subtract,
			decls.NewOverload(overloads.SubtractInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Int64),
			decls.NewOverload(overloads.SubtractUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Uint64),
			decls.NewOverload(overloads.SubtractDouble,
				[]*checked.Type{types.Double, types.Double}, types.Double),
			decls.NewOverload(overloads.SubtractTimestampTimestamp,
				[]*checked.Type{types.Timestamp, types.Timestamp}, types.Duration),
			decls.NewOverload(overloads.SubtractTimestampDuration,
				[]*checked.Type{types.Timestamp, types.Duration}, types.Timestamp),
			decls.NewOverload(overloads.SubtractDurationDuration,
				[]*checked.Type{types.Duration, types.Duration}, types.Duration)),

		decls.NewFunction(operators.Multiply,
			decls.NewOverload(overloads.MultiplyInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Int64),
			decls.NewOverload(overloads.MultiplyUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Uint64),
			decls.NewOverload(overloads.MultiplyDouble,
				[]*checked.Type{types.Double, types.Double}, types.Double)),

		decls.NewFunction(operators.Divide,
			decls.NewOverload(overloads.DivideInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Int64),
			decls.NewOverload(overloads.DivideUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Uint64),
			decls.NewOverload(overloads.DivideDouble,
				[]*checked.Type{types.Double, types.Double}, types.Double)),

		decls.NewFunction(operators.Modulo,
			decls.NewOverload(overloads.ModuloInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Int64),
			decls.NewOverload(overloads.ModuloUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Uint64)),

		decls.NewFunction(operators.Add,
			decls.NewOverload(overloads.AddInt64,
				[]*checked.Type{types.Int64, types.Int64}, types.Int64),
			decls.NewOverload(overloads.AddUint64,
				[]*checked.Type{types.Uint64, types.Uint64}, types.Uint64),
			decls.NewOverload(overloads.AddDouble,
				[]*checked.Type{types.Double, types.Double}, types.Double),
			decls.NewOverload(overloads.AddString,
				[]*checked.Type{types.String, types.String}, types.String),
			decls.NewOverload(overloads.AddBytes,
				[]*checked.Type{types.Bytes, types.Bytes}, types.Bytes),
			decls.NewParameterizedOverload(overloads.AddList,
				[]*checked.Type{listOfA, listOfA}, listOfA,
				typeParamAList),
			decls.NewOverload(overloads.AddTimestampDuration,
				[]*checked.Type{types.Timestamp, types.Duration}, types.Timestamp),
			decls.NewOverload(overloads.AddDurationTimestamp,
				[]*checked.Type{types.Duration, types.Timestamp}, types.Timestamp),
			decls.NewOverload(overloads.AddDurationDuration,
				[]*checked.Type{types.Duration, types.Duration}, types.Duration)),

		decls.NewFunction(operators.Negate,
			decls.NewOverload(overloads.NegateInt64,
				[]*checked.Type{types.Int64}, types.Int64),
			decls.NewOverload(overloads.NegateDouble,
				[]*checked.Type{types.Double}, types.Double)),

		// Index

		decls.NewFunction(operators.Index,
			decls.NewParameterizedOverload(overloads.IndexList,
				[]*checked.Type{listOfA, types.Int64}, paramA,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.IndexMap,
				[]*checked.Type{mapOfAB, paramA}, paramB,
				typeParamABList)),
		//decls.NewOverload(overloads.IndexMessage,
		//	[]*checked.Type{types.Dyn, types.String}, types.Dyn)),

		// Collections

		decls.NewFunction(overloads.Size,
			decls.NewInstanceOverload(overloads.SizeStringInst,
				[]*checked.Type{types.String}, types.Int64),
			decls.NewInstanceOverload(overloads.SizeBytesInst,
				[]*checked.Type{types.Bytes}, types.Int64),
			decls.NewParameterizedInstanceOverload(overloads.SizeListInst,
				[]*checked.Type{listOfA}, types.Int64, typeParamAList),
			decls.NewParameterizedInstanceOverload(overloads.SizeMapInst,
				[]*checked.Type{mapOfAB}, types.Int64, typeParamABList),
			decls.NewOverload(overloads.SizeString,
				[]*checked.Type{types.String}, types.Int64),
			decls.NewOverload(overloads.SizeBytes,
				[]*checked.Type{types.Bytes}, types.Int64),
			decls.NewParameterizedOverload(overloads.SizeList,
				[]*checked.Type{listOfA}, types.Int64, typeParamAList),
			decls.NewParameterizedOverload(overloads.SizeMap,
				[]*checked.Type{mapOfAB}, types.Int64, typeParamABList)),

		decls.NewFunction(operators.In,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checked.Type{paramA, listOfA}, types.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checked.Type{paramA, mapOfAB}, types.Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*checked.Type{types.Dyn, types.String},types.Bool)),

		// Deprecated 'in()' function

		decls.NewFunction(overloads.DeprecatedIn,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checked.Type{paramA, listOfA}, types.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checked.Type{paramA, mapOfAB}, types.Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*checked.Type{types.Dyn, types.String},types.Bool)),

		// Conversions to type

		decls.NewFunction(overloads.TypeConvertType,
			decls.NewParameterizedOverload(overloads.TypeConvertType,
				[]*checked.Type{paramA}, types.NewType(paramA), typeParamAList)),

		// Conversions to int

		decls.NewFunction(overloads.TypeConvertInt,
			decls.NewOverload(overloads.IntToInt, []*checked.Type{types.Int64}, types.Int64),
			decls.NewOverload(overloads.UintToInt, []*checked.Type{types.Uint64}, types.Int64),
			decls.NewOverload(overloads.DoubleToInt, []*checked.Type{types.Double}, types.Int64),
			decls.NewOverload(overloads.StringToInt, []*checked.Type{types.String}, types.Int64),
			decls.NewOverload(overloads.TimestampToInt, []*checked.Type{types.Timestamp}, types.Int64),
			decls.NewOverload(overloads.DurationToInt, []*checked.Type{types.Duration}, types.Int64)),

		// Conversions to uint

		decls.NewFunction(overloads.TypeConvertUint,
			decls.NewOverload(overloads.UintToUint, []*checked.Type{types.Uint64}, types.Uint64),
			decls.NewOverload(overloads.IntToUint, []*checked.Type{types.Int64}, types.Uint64),
			decls.NewOverload(overloads.DoubleToUint, []*checked.Type{types.Double}, types.Uint64),
			decls.NewOverload(overloads.StringToUint, []*checked.Type{types.String}, types.Uint64)),

		// Conversions to double

		decls.NewFunction(overloads.TypeConvertDouble,
			decls.NewOverload(overloads.DoubleToDouble, []*checked.Type{types.Double}, types.Double),
			decls.NewOverload(overloads.IntToDouble, []*checked.Type{types.Int64}, types.Double),
			decls.NewOverload(overloads.UintToDouble, []*checked.Type{types.Uint64}, types.Double),
			decls.NewOverload(overloads.StringToDouble, []*checked.Type{types.String}, types.Double)),

		// Conversions to bool

		decls.NewFunction(overloads.TypeConvertBool,
			decls.NewOverload(overloads.BoolToBool, []*checked.Type{types.Bool}, types.Bool),
			decls.NewOverload(overloads.StringToBool, []*checked.Type{types.String}, types.Bool)),

		// Conversions to string

		decls.NewFunction(overloads.TypeConvertString,
			decls.NewOverload(overloads.StringToString, []*checked.Type{types.String}, types.String),
			decls.NewOverload(overloads.BoolToString, []*checked.Type{types.Bool}, types.String),
			decls.NewOverload(overloads.IntToString, []*checked.Type{types.Int64}, types.String),
			decls.NewOverload(overloads.UintToString, []*checked.Type{types.Uint64}, types.String),
			decls.NewOverload(overloads.DoubleToString, []*checked.Type{types.Double}, types.String),
			decls.NewOverload(overloads.BytesToString, []*checked.Type{types.Bytes}, types.String),
			decls.NewOverload(overloads.TimestampToString, []*checked.Type{types.Timestamp}, types.String),
			decls.NewOverload(overloads.DurationToString, []*checked.Type{types.Duration}, types.String)),

		// Conversions to bytes

		decls.NewFunction(overloads.TypeConvertBytes,
			decls.NewOverload(overloads.BytesToBytes, []*checked.Type{types.Bytes}, types.Bytes),
			decls.NewOverload(overloads.StringToBytes, []*checked.Type{types.String}, types.Bytes)),

		// Conversions to timestamps

		decls.NewFunction(overloads.TypeConvertTimestamp,
			decls.NewOverload(overloads.TimestampToTimestamp,
				[]*checked.Type{types.Timestamp}, types.Timestamp),
			decls.NewOverload(overloads.StringToTimestamp,
				[]*checked.Type{types.String}, types.Timestamp),
			decls.NewOverload(overloads.IntToTimestamp,
				[]*checked.Type{types.Int64}, types.Timestamp)),

		// Conversions to durations

		decls.NewFunction(overloads.TypeConvertDuration,
			decls.NewOverload(overloads.DurationToDuration,
				[]*checked.Type{types.Duration}, types.Duration),
			decls.NewOverload(overloads.StringToDuration,
				[]*checked.Type{types.String}, types.Duration),
			decls.NewOverload(overloads.IntToDuration,
				[]*checked.Type{types.Int64}, types.Duration)),

		// Conversions to Dyn

		decls.NewFunction(overloads.TypeConvertDyn,
			decls.NewParameterizedOverload(overloads.ToDyn, []*checked.Type{paramA}, types.Dyn,
				typeParamAList)),

		// Date/time functions

		decls.NewFunction(overloads.TimeGetFullYear,
			decls.NewInstanceOverload(overloads.TimestampToYear,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToYearWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64)),

		decls.NewFunction(overloads.TimeGetMonth,
			decls.NewInstanceOverload(overloads.TimestampToMonth,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToMonthWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64)),

		decls.NewFunction(overloads.TimeGetDayOfYear,
			decls.NewInstanceOverload(overloads.TimestampToDayOfYear,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToDayOfYearWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64)),

		decls.NewFunction(overloads.TimeGetDayOfMonth,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBased,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64)),

		decls.NewFunction(overloads.TimeGetDate,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBased,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64)),

		decls.NewFunction(overloads.TimeGetDayOfWeek,
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeek,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeekWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64)),

		decls.NewFunction(overloads.TimeGetHours,
			decls.NewInstanceOverload(overloads.TimestampToHours,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToHoursWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64),
			decls.NewInstanceOverload(overloads.DurationToHours,
				[]*checked.Type{types.Duration}, types.Int64)),

		decls.NewFunction(overloads.TimeGetMinutes,
			decls.NewInstanceOverload(overloads.TimestampToMinutes,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToMinutesWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64),
			decls.NewInstanceOverload(overloads.DurationToMinutes,
				[]*checked.Type{types.Duration}, types.Int64)),

		decls.NewFunction(overloads.TimeGetSeconds,
			decls.NewInstanceOverload(overloads.TimestampToSeconds,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToSecondsWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64),
			decls.NewInstanceOverload(overloads.DurationToSeconds,
				[]*checked.Type{types.Duration}, types.Int64)),

		decls.NewFunction(overloads.TimeGetMilliseconds,
			decls.NewInstanceOverload(overloads.TimestampToMilliseconds,
				[]*checked.Type{types.Timestamp}, types.Int64),
			decls.NewInstanceOverload(overloads.TimestampToMillisecondsWithTz,
				[]*checked.Type{types.Timestamp, types.String}, types.Int64),
			decls.NewInstanceOverload(overloads.DurationToMilliseconds,
				[]*checked.Type{types.Duration}, types.Int64))}...)
}
