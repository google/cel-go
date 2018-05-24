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
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-spec/proto/checked/v1/checked"
)

func StandardDeclarations() []*checked.Decl {
	// Some shortcuts we use when building declarations.
	paramA := decls.NewTypeParamType("A")
	typeParamAList := []string{"A"}
	listOfA := decls.NewListType(paramA)
	paramB := decls.NewTypeParamType("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := decls.NewMapType(paramA, paramB)

	var idents []*checked.Decl
	for _, t := range []*checked.Type{
		decls.Int, decls.Uint, decls.Bool,
		decls.Double, decls.Bytes, decls.String} {
		idents = append(idents,
			decls.NewIdent(FormatCheckedType(t), decls.NewTypeType(t), nil))
	}
	idents = append(idents,
		decls.NewIdent("list", decls.NewTypeType(listOfA), nil),
		decls.NewIdent("map", decls.NewTypeType(mapOfAB), nil))

	// Booleans
	// TODO: allow the conditional to return a heterogenous type.
	return append(idents, []*checked.Decl{
		decls.NewFunction(operators.Conditional,
			decls.NewParameterizedOverload(overloads.Conditional,
				[]*checked.Type{decls.Bool, paramA, paramA}, paramA,
				typeParamAList)),

		decls.NewFunction(operators.LogicalAnd,
			decls.NewOverload(overloads.LogicalAnd,
				[]*checked.Type{decls.Bool, decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.LogicalOr,
			decls.NewOverload(overloads.LogicalOr,
				[]*checked.Type{decls.Bool, decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.LogicalNot,
			decls.NewOverload(overloads.LogicalNot,
				[]*checked.Type{decls.Bool}, decls.Bool)),

		decls.NewFunction(overloads.Matches,
			decls.NewInstanceOverload(overloads.MatchString,
				[]*checked.Type{decls.String, decls.String}, decls.Bool)),

		// Relations

		decls.NewFunction(operators.Less,
			decls.NewOverload(overloads.LessBool,
				[]*checked.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.LessInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.LessUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.LessDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.LessString,
				[]*checked.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.LessBytes,
				[]*checked.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.LessTimestamp,
				[]*checked.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.LessDuration,
				[]*checked.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.LessEquals,
			decls.NewOverload(overloads.LessEqualsBool,
				[]*checked.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsString,
				[]*checked.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsBytes,
				[]*checked.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsTimestamp,
				[]*checked.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsDuration,
				[]*checked.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.Greater,
			decls.NewOverload(overloads.GreaterBool,
				[]*checked.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.GreaterInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.GreaterUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.GreaterDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.GreaterString,
				[]*checked.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.GreaterBytes,
				[]*checked.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.GreaterTimestamp,
				[]*checked.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.GreaterDuration,
				[]*checked.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.GreaterEquals,
			decls.NewOverload(overloads.GreaterEqualsBool,
				[]*checked.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsString,
				[]*checked.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsBytes,
				[]*checked.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsTimestamp,
				[]*checked.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsDuration,
				[]*checked.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.Equals,
			decls.NewParameterizedOverload(overloads.Equals,
				[]*checked.Type{paramA, paramA}, decls.Bool,
				typeParamAList)),

		decls.NewFunction(operators.NotEquals,
			decls.NewParameterizedOverload(overloads.NotEquals,
				[]*checked.Type{paramA, paramA}, decls.Bool,
				typeParamAList)),

		// Algebra

		decls.NewFunction(operators.Subtract,
			decls.NewOverload(overloads.SubtractInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.SubtractUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.SubtractDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Double),
			decls.NewOverload(overloads.SubtractTimestampTimestamp,
				[]*checked.Type{decls.Timestamp, decls.Timestamp}, decls.Duration),
			decls.NewOverload(overloads.SubtractTimestampDuration,
				[]*checked.Type{decls.Timestamp, decls.Duration}, decls.Timestamp),
			decls.NewOverload(overloads.SubtractDurationDuration,
				[]*checked.Type{decls.Duration, decls.Duration}, decls.Duration)),

		decls.NewFunction(operators.Multiply,
			decls.NewOverload(overloads.MultiplyInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.MultiplyUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.MultiplyDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Double)),

		decls.NewFunction(operators.Divide,
			decls.NewOverload(overloads.DivideInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.DivideUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.DivideDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Double)),

		decls.NewFunction(operators.Modulo,
			decls.NewOverload(overloads.ModuloInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.ModuloUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Uint)),

		decls.NewFunction(operators.Add,
			decls.NewOverload(overloads.AddInt64,
				[]*checked.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.AddUint64,
				[]*checked.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.AddDouble,
				[]*checked.Type{decls.Double, decls.Double}, decls.Double),
			decls.NewOverload(overloads.AddString,
				[]*checked.Type{decls.String, decls.String}, decls.String),
			decls.NewOverload(overloads.AddBytes,
				[]*checked.Type{decls.Bytes, decls.Bytes}, decls.Bytes),
			decls.NewParameterizedOverload(overloads.AddList,
				[]*checked.Type{listOfA, listOfA}, listOfA,
				typeParamAList),
			decls.NewOverload(overloads.AddTimestampDuration,
				[]*checked.Type{decls.Timestamp, decls.Duration}, decls.Timestamp),
			decls.NewOverload(overloads.AddDurationTimestamp,
				[]*checked.Type{decls.Duration, decls.Timestamp}, decls.Timestamp),
			decls.NewOverload(overloads.AddDurationDuration,
				[]*checked.Type{decls.Duration, decls.Duration}, decls.Duration)),

		decls.NewFunction(operators.Negate,
			decls.NewOverload(overloads.NegateInt64,
				[]*checked.Type{decls.Int}, decls.Int),
			decls.NewOverload(overloads.NegateDouble,
				[]*checked.Type{decls.Double}, decls.Double)),

		// Index

		decls.NewFunction(operators.Index,
			decls.NewParameterizedOverload(overloads.IndexList,
				[]*checked.Type{listOfA, decls.Int}, paramA,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.IndexMap,
				[]*checked.Type{mapOfAB, paramA}, paramB,
				typeParamABList)),
		//decls.NewOverload(overloads.IndexMessage,
		//	[]*checked.Type{decls.Dyn, decls.String}, decls.Dyn)),

		// Collections

		decls.NewFunction(overloads.Size,
			decls.NewInstanceOverload(overloads.SizeStringInst,
				[]*checked.Type{decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.SizeBytesInst,
				[]*checked.Type{decls.Bytes}, decls.Int),
			decls.NewParameterizedInstanceOverload(overloads.SizeListInst,
				[]*checked.Type{listOfA}, decls.Int, typeParamAList),
			decls.NewParameterizedInstanceOverload(overloads.SizeMapInst,
				[]*checked.Type{mapOfAB}, decls.Int, typeParamABList),
			decls.NewOverload(overloads.SizeString,
				[]*checked.Type{decls.String}, decls.Int),
			decls.NewOverload(overloads.SizeBytes,
				[]*checked.Type{decls.Bytes}, decls.Int),
			decls.NewParameterizedOverload(overloads.SizeList,
				[]*checked.Type{listOfA}, decls.Int, typeParamAList),
			decls.NewParameterizedOverload(overloads.SizeMap,
				[]*checked.Type{mapOfAB}, decls.Int, typeParamABList)),

		decls.NewFunction(operators.In,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checked.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checked.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),


		// Deprecated 'in()' function

		decls.NewFunction(overloads.DeprecatedIn,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checked.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checked.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*checked.Type{Dyn, decls.String},decls.Bool)),

		// Conversions to type

		decls.NewFunction(overloads.TypeConvertType,
			decls.NewParameterizedOverload(overloads.TypeConvertType,
				[]*checked.Type{paramA}, decls.NewTypeType(paramA), typeParamAList)),

		// Conversions to int

		decls.NewFunction(overloads.TypeConvertInt,
			decls.NewOverload(overloads.IntToInt, []*checked.Type{decls.Int}, decls.Int),
			decls.NewOverload(overloads.UintToInt, []*checked.Type{decls.Uint}, decls.Int),
			decls.NewOverload(overloads.DoubleToInt, []*checked.Type{decls.Double}, decls.Int),
			decls.NewOverload(overloads.StringToInt, []*checked.Type{decls.String}, decls.Int),
			decls.NewOverload(overloads.TimestampToInt, []*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewOverload(overloads.DurationToInt, []*checked.Type{decls.Duration}, decls.Int)),

		// Conversions to uint

		decls.NewFunction(overloads.TypeConvertUint,
			decls.NewOverload(overloads.UintToUint, []*checked.Type{decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.IntToUint, []*checked.Type{decls.Int}, decls.Uint),
			decls.NewOverload(overloads.DoubleToUint, []*checked.Type{decls.Double}, decls.Uint),
			decls.NewOverload(overloads.StringToUint, []*checked.Type{decls.String}, decls.Uint)),

		// Conversions to double

		decls.NewFunction(overloads.TypeConvertDouble,
			decls.NewOverload(overloads.DoubleToDouble, []*checked.Type{decls.Double}, decls.Double),
			decls.NewOverload(overloads.IntToDouble, []*checked.Type{decls.Int}, decls.Double),
			decls.NewOverload(overloads.UintToDouble, []*checked.Type{decls.Uint}, decls.Double),
			decls.NewOverload(overloads.StringToDouble, []*checked.Type{decls.String}, decls.Double)),

		// Conversions to bool

		decls.NewFunction(overloads.TypeConvertBool,
			decls.NewOverload(overloads.BoolToBool, []*checked.Type{decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.StringToBool, []*checked.Type{decls.String}, decls.Bool)),

		// Conversions to string

		decls.NewFunction(overloads.TypeConvertString,
			decls.NewOverload(overloads.StringToString, []*checked.Type{decls.String}, decls.String),
			decls.NewOverload(overloads.BoolToString, []*checked.Type{decls.Bool}, decls.String),
			decls.NewOverload(overloads.IntToString, []*checked.Type{decls.Int}, decls.String),
			decls.NewOverload(overloads.UintToString, []*checked.Type{decls.Uint}, decls.String),
			decls.NewOverload(overloads.DoubleToString, []*checked.Type{decls.Double}, decls.String),
			decls.NewOverload(overloads.BytesToString, []*checked.Type{decls.Bytes}, decls.String),
			decls.NewOverload(overloads.TimestampToString, []*checked.Type{decls.Timestamp}, decls.String),
			decls.NewOverload(overloads.DurationToString, []*checked.Type{decls.Duration}, decls.String)),

		// Conversions to bytes

		decls.NewFunction(overloads.TypeConvertBytes,
			decls.NewOverload(overloads.BytesToBytes, []*checked.Type{decls.Bytes}, decls.Bytes),
			decls.NewOverload(overloads.StringToBytes, []*checked.Type{decls.String}, decls.Bytes)),

		// Conversions to timestamps

		decls.NewFunction(overloads.TypeConvertTimestamp,
			decls.NewOverload(overloads.TimestampToTimestamp,
				[]*checked.Type{decls.Timestamp}, decls.Timestamp),
			decls.NewOverload(overloads.StringToTimestamp,
				[]*checked.Type{decls.String}, decls.Timestamp),
			decls.NewOverload(overloads.IntToTimestamp,
				[]*checked.Type{decls.Int}, decls.Timestamp)),

		// Conversions to durations

		decls.NewFunction(overloads.TypeConvertDuration,
			decls.NewOverload(overloads.DurationToDuration,
				[]*checked.Type{decls.Duration}, decls.Duration),
			decls.NewOverload(overloads.StringToDuration,
				[]*checked.Type{decls.String}, decls.Duration),
			decls.NewOverload(overloads.IntToDuration,
				[]*checked.Type{decls.Int}, decls.Duration)),

		// Conversions to Dyn

		decls.NewFunction(overloads.TypeConvertDyn,
			decls.NewParameterizedOverload(overloads.ToDyn,
				[]*checked.Type{paramA}, decls.Dyn,
				typeParamAList)),

		// Date/time functions

		decls.NewFunction(overloads.TimeGetFullYear,
			decls.NewInstanceOverload(overloads.TimestampToYear,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToYearWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMonth,
			decls.NewInstanceOverload(overloads.TimestampToMonth,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMonthWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfYear,
			decls.NewInstanceOverload(overloads.TimestampToDayOfYear,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfYearWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfMonth,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBased,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDate,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBased,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfWeek,
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeek,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeekWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetHours,
			decls.NewInstanceOverload(overloads.TimestampToHours,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToHoursWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToHours,
				[]*checked.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMinutes,
			decls.NewInstanceOverload(overloads.TimestampToMinutes,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMinutesWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToMinutes,
				[]*checked.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetSeconds,
			decls.NewInstanceOverload(overloads.TimestampToSeconds,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToSecondsWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToSeconds,
				[]*checked.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMilliseconds,
			decls.NewInstanceOverload(overloads.TimestampToMilliseconds,
				[]*checked.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMillisecondsWithTz,
				[]*checked.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToMilliseconds,
				[]*checked.Type{decls.Duration}, decls.Int))}...)
}
