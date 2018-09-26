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
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func StandardDeclarations() []*expr.Decl {
	// Some shortcuts we use when building declarations.
	paramA := decls.NewTypeParamType("A")
	typeParamAList := []string{"A"}
	listOfA := decls.NewListType(paramA)
	paramB := decls.NewTypeParamType("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := decls.NewMapType(paramA, paramB)

	var idents []*expr.Decl
	for _, t := range []*expr.Type{
		decls.Int, decls.Uint, decls.Bool,
		decls.Double, decls.Bytes, decls.String} {
		idents = append(idents,
			decls.NewIdent(FormatCheckedType(t), decls.NewTypeType(t), nil))
	}
	idents = append(idents,
		decls.NewIdent("list", decls.NewTypeType(listOfA), nil),
		decls.NewIdent("map", decls.NewTypeType(mapOfAB), nil),
		decls.NewIdent("null_type", decls.NewTypeType(decls.Null), nil),
		decls.NewIdent("type", decls.NewTypeType(decls.NewTypeType(nil)), nil))

	// Booleans
	// TODO: allow the conditional to return a heterogenous type.
	return append(idents, []*expr.Decl{
		decls.NewFunction(operators.Conditional,
			decls.NewParameterizedOverload(overloads.Conditional,
				[]*expr.Type{decls.Bool, paramA, paramA}, paramA,
				typeParamAList)),

		decls.NewFunction(operators.LogicalAnd,
			decls.NewOverload(overloads.LogicalAnd,
				[]*expr.Type{decls.Bool, decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.LogicalOr,
			decls.NewOverload(overloads.LogicalOr,
				[]*expr.Type{decls.Bool, decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.LogicalNot,
			decls.NewOverload(overloads.LogicalNot,
				[]*expr.Type{decls.Bool}, decls.Bool)),

		decls.NewFunction(overloads.Matches,
			decls.NewInstanceOverload(overloads.MatchString,
				[]*expr.Type{decls.String, decls.String}, decls.Bool)),

		// Relations

		decls.NewFunction(operators.Less,
			decls.NewOverload(overloads.LessBool,
				[]*expr.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.LessInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.LessUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.LessDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.LessString,
				[]*expr.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.LessBytes,
				[]*expr.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.LessTimestamp,
				[]*expr.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.LessDuration,
				[]*expr.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.LessEquals,
			decls.NewOverload(overloads.LessEqualsBool,
				[]*expr.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsString,
				[]*expr.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsBytes,
				[]*expr.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsTimestamp,
				[]*expr.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsDuration,
				[]*expr.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.Greater,
			decls.NewOverload(overloads.GreaterBool,
				[]*expr.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.GreaterInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.GreaterUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.GreaterDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.GreaterString,
				[]*expr.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.GreaterBytes,
				[]*expr.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.GreaterTimestamp,
				[]*expr.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.GreaterDuration,
				[]*expr.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.GreaterEquals,
			decls.NewOverload(overloads.GreaterEqualsBool,
				[]*expr.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsString,
				[]*expr.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsBytes,
				[]*expr.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsTimestamp,
				[]*expr.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsDuration,
				[]*expr.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.Equals,
			decls.NewParameterizedOverload(overloads.Equals,
				[]*expr.Type{paramA, paramA}, decls.Bool,
				typeParamAList)),

		decls.NewFunction(operators.NotEquals,
			decls.NewParameterizedOverload(overloads.NotEquals,
				[]*expr.Type{paramA, paramA}, decls.Bool,
				typeParamAList)),

		// Algebra

		decls.NewFunction(operators.Subtract,
			decls.NewOverload(overloads.SubtractInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.SubtractUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.SubtractDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Double),
			decls.NewOverload(overloads.SubtractTimestampTimestamp,
				[]*expr.Type{decls.Timestamp, decls.Timestamp}, decls.Duration),
			decls.NewOverload(overloads.SubtractTimestampDuration,
				[]*expr.Type{decls.Timestamp, decls.Duration}, decls.Timestamp),
			decls.NewOverload(overloads.SubtractDurationDuration,
				[]*expr.Type{decls.Duration, decls.Duration}, decls.Duration)),

		decls.NewFunction(operators.Multiply,
			decls.NewOverload(overloads.MultiplyInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.MultiplyUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.MultiplyDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Double)),

		decls.NewFunction(operators.Divide,
			decls.NewOverload(overloads.DivideInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.DivideUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.DivideDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Double)),

		decls.NewFunction(operators.Modulo,
			decls.NewOverload(overloads.ModuloInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.ModuloUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Uint)),

		decls.NewFunction(operators.Add,
			decls.NewOverload(overloads.AddInt64,
				[]*expr.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.AddUint64,
				[]*expr.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.AddDouble,
				[]*expr.Type{decls.Double, decls.Double}, decls.Double),
			decls.NewOverload(overloads.AddString,
				[]*expr.Type{decls.String, decls.String}, decls.String),
			decls.NewOverload(overloads.AddBytes,
				[]*expr.Type{decls.Bytes, decls.Bytes}, decls.Bytes),
			decls.NewParameterizedOverload(overloads.AddList,
				[]*expr.Type{listOfA, listOfA}, listOfA,
				typeParamAList),
			decls.NewOverload(overloads.AddTimestampDuration,
				[]*expr.Type{decls.Timestamp, decls.Duration}, decls.Timestamp),
			decls.NewOverload(overloads.AddDurationTimestamp,
				[]*expr.Type{decls.Duration, decls.Timestamp}, decls.Timestamp),
			decls.NewOverload(overloads.AddDurationDuration,
				[]*expr.Type{decls.Duration, decls.Duration}, decls.Duration)),

		decls.NewFunction(operators.Negate,
			decls.NewOverload(overloads.NegateInt64,
				[]*expr.Type{decls.Int}, decls.Int),
			decls.NewOverload(overloads.NegateDouble,
				[]*expr.Type{decls.Double}, decls.Double)),

		// Index

		decls.NewFunction(operators.Index,
			decls.NewParameterizedOverload(overloads.IndexList,
				[]*expr.Type{listOfA, decls.Int}, paramA,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.IndexMap,
				[]*expr.Type{mapOfAB, paramA}, paramB,
				typeParamABList)),
		//decls.NewOverload(overloads.IndexMessage,
		//	[]*expr.Type{decls.Dyn, decls.String}, decls.Dyn)),

		// Collections

		decls.NewFunction(overloads.Size,
			decls.NewInstanceOverload(overloads.SizeStringInst,
				[]*expr.Type{decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.SizeBytesInst,
				[]*expr.Type{decls.Bytes}, decls.Int),
			decls.NewParameterizedInstanceOverload(overloads.SizeListInst,
				[]*expr.Type{listOfA}, decls.Int, typeParamAList),
			decls.NewParameterizedInstanceOverload(overloads.SizeMapInst,
				[]*expr.Type{mapOfAB}, decls.Int, typeParamABList),
			decls.NewOverload(overloads.SizeString,
				[]*expr.Type{decls.String}, decls.Int),
			decls.NewOverload(overloads.SizeBytes,
				[]*expr.Type{decls.Bytes}, decls.Int),
			decls.NewParameterizedOverload(overloads.SizeList,
				[]*expr.Type{listOfA}, decls.Int, typeParamAList),
			decls.NewParameterizedOverload(overloads.SizeMap,
				[]*expr.Type{mapOfAB}, decls.Int, typeParamABList)),

		decls.NewFunction(operators.In,
			decls.NewParameterizedOverload(overloads.InList,
				[]*expr.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*expr.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),

		// Deprecated 'in()' function

		decls.NewFunction(overloads.DeprecatedIn,
			decls.NewParameterizedOverload(overloads.InList,
				[]*expr.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*expr.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*expr.Type{Dyn, decls.String},decls.Bool)),

		// Conversions to type

		decls.NewFunction(overloads.TypeConvertType,
			decls.NewParameterizedOverload(overloads.TypeConvertType,
				[]*expr.Type{paramA}, decls.NewTypeType(paramA), typeParamAList)),

		// Conversions to int

		decls.NewFunction(overloads.TypeConvertInt,
			decls.NewOverload(overloads.IntToInt, []*expr.Type{decls.Int}, decls.Int),
			decls.NewOverload(overloads.UintToInt, []*expr.Type{decls.Uint}, decls.Int),
			decls.NewOverload(overloads.DoubleToInt, []*expr.Type{decls.Double}, decls.Int),
			decls.NewOverload(overloads.StringToInt, []*expr.Type{decls.String}, decls.Int),
			decls.NewOverload(overloads.TimestampToInt, []*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewOverload(overloads.DurationToInt, []*expr.Type{decls.Duration}, decls.Int)),

		// Conversions to uint

		decls.NewFunction(overloads.TypeConvertUint,
			decls.NewOverload(overloads.UintToUint, []*expr.Type{decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.IntToUint, []*expr.Type{decls.Int}, decls.Uint),
			decls.NewOverload(overloads.DoubleToUint, []*expr.Type{decls.Double}, decls.Uint),
			decls.NewOverload(overloads.StringToUint, []*expr.Type{decls.String}, decls.Uint)),

		// Conversions to double

		decls.NewFunction(overloads.TypeConvertDouble,
			decls.NewOverload(overloads.DoubleToDouble, []*expr.Type{decls.Double}, decls.Double),
			decls.NewOverload(overloads.IntToDouble, []*expr.Type{decls.Int}, decls.Double),
			decls.NewOverload(overloads.UintToDouble, []*expr.Type{decls.Uint}, decls.Double),
			decls.NewOverload(overloads.StringToDouble, []*expr.Type{decls.String}, decls.Double)),

		// Conversions to bool

		decls.NewFunction(overloads.TypeConvertBool,
			decls.NewOverload(overloads.BoolToBool, []*expr.Type{decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.StringToBool, []*expr.Type{decls.String}, decls.Bool)),

		// Conversions to string

		decls.NewFunction(overloads.TypeConvertString,
			decls.NewOverload(overloads.StringToString, []*expr.Type{decls.String}, decls.String),
			decls.NewOverload(overloads.BoolToString, []*expr.Type{decls.Bool}, decls.String),
			decls.NewOverload(overloads.IntToString, []*expr.Type{decls.Int}, decls.String),
			decls.NewOverload(overloads.UintToString, []*expr.Type{decls.Uint}, decls.String),
			decls.NewOverload(overloads.DoubleToString, []*expr.Type{decls.Double}, decls.String),
			decls.NewOverload(overloads.BytesToString, []*expr.Type{decls.Bytes}, decls.String),
			decls.NewOverload(overloads.TimestampToString, []*expr.Type{decls.Timestamp}, decls.String),
			decls.NewOverload(overloads.DurationToString, []*expr.Type{decls.Duration}, decls.String)),

		// Conversions to bytes

		decls.NewFunction(overloads.TypeConvertBytes,
			decls.NewOverload(overloads.BytesToBytes, []*expr.Type{decls.Bytes}, decls.Bytes),
			decls.NewOverload(overloads.StringToBytes, []*expr.Type{decls.String}, decls.Bytes)),

		// Conversions to timestamps

		decls.NewFunction(overloads.TypeConvertTimestamp,
			decls.NewOverload(overloads.TimestampToTimestamp,
				[]*expr.Type{decls.Timestamp}, decls.Timestamp),
			decls.NewOverload(overloads.StringToTimestamp,
				[]*expr.Type{decls.String}, decls.Timestamp),
			decls.NewOverload(overloads.IntToTimestamp,
				[]*expr.Type{decls.Int}, decls.Timestamp)),

		// Conversions to durations

		decls.NewFunction(overloads.TypeConvertDuration,
			decls.NewOverload(overloads.DurationToDuration,
				[]*expr.Type{decls.Duration}, decls.Duration),
			decls.NewOverload(overloads.StringToDuration,
				[]*expr.Type{decls.String}, decls.Duration),
			decls.NewOverload(overloads.IntToDuration,
				[]*expr.Type{decls.Int}, decls.Duration)),

		// Conversions to Dyn

		decls.NewFunction(overloads.TypeConvertDyn,
			decls.NewParameterizedOverload(overloads.ToDyn,
				[]*expr.Type{paramA}, decls.Dyn,
				typeParamAList)),

		// Date/time functions

		decls.NewFunction(overloads.TimeGetFullYear,
			decls.NewInstanceOverload(overloads.TimestampToYear,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToYearWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMonth,
			decls.NewInstanceOverload(overloads.TimestampToMonth,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMonthWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfYear,
			decls.NewInstanceOverload(overloads.TimestampToDayOfYear,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfYearWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfMonth,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBased,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDate,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBased,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfWeek,
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeek,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeekWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetHours,
			decls.NewInstanceOverload(overloads.TimestampToHours,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToHoursWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToHours,
				[]*expr.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMinutes,
			decls.NewInstanceOverload(overloads.TimestampToMinutes,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMinutesWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToMinutes,
				[]*expr.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetSeconds,
			decls.NewInstanceOverload(overloads.TimestampToSeconds,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToSecondsWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToSeconds,
				[]*expr.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMilliseconds,
			decls.NewInstanceOverload(overloads.TimestampToMilliseconds,
				[]*expr.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMillisecondsWithTz,
				[]*expr.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToMilliseconds,
				[]*expr.Type{decls.Duration}, decls.Int))}...)
}
