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
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
)

func StandardDeclarations() []*checkedpb.Decl {
	// Some shortcuts we use when building declarations.
	paramA := decls.NewTypeParamType("A")
	typeParamAList := []string{"A"}
	listOfA := decls.NewListType(paramA)
	paramB := decls.NewTypeParamType("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := decls.NewMapType(paramA, paramB)

	var idents []*checkedpb.Decl
	for _, t := range []*checkedpb.Type{
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
	return append(idents, []*checkedpb.Decl{
		decls.NewFunction(operators.Conditional,
			decls.NewParameterizedOverload(overloads.Conditional,
				[]*checkedpb.Type{decls.Bool, paramA, paramA}, paramA,
				typeParamAList)),

		decls.NewFunction(operators.LogicalAnd,
			decls.NewOverload(overloads.LogicalAnd,
				[]*checkedpb.Type{decls.Bool, decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.LogicalOr,
			decls.NewOverload(overloads.LogicalOr,
				[]*checkedpb.Type{decls.Bool, decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.LogicalNot,
			decls.NewOverload(overloads.LogicalNot,
				[]*checkedpb.Type{decls.Bool}, decls.Bool)),

		decls.NewFunction(overloads.Matches,
			decls.NewInstanceOverload(overloads.MatchString,
				[]*checkedpb.Type{decls.String, decls.String}, decls.Bool)),

		// Relations

		decls.NewFunction(operators.Less,
			decls.NewOverload(overloads.LessBool,
				[]*checkedpb.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.LessInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.LessUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.LessDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.LessString,
				[]*checkedpb.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.LessBytes,
				[]*checkedpb.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.LessTimestamp,
				[]*checkedpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.LessDuration,
				[]*checkedpb.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.LessEquals,
			decls.NewOverload(overloads.LessEqualsBool,
				[]*checkedpb.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsString,
				[]*checkedpb.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsBytes,
				[]*checkedpb.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsTimestamp,
				[]*checkedpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.LessEqualsDuration,
				[]*checkedpb.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.Greater,
			decls.NewOverload(overloads.GreaterBool,
				[]*checkedpb.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.GreaterInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.GreaterUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.GreaterDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.GreaterString,
				[]*checkedpb.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.GreaterBytes,
				[]*checkedpb.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.GreaterTimestamp,
				[]*checkedpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.GreaterDuration,
				[]*checkedpb.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.GreaterEquals,
			decls.NewOverload(overloads.GreaterEqualsBool,
				[]*checkedpb.Type{decls.Bool, decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsString,
				[]*checkedpb.Type{decls.String, decls.String}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsBytes,
				[]*checkedpb.Type{decls.Bytes, decls.Bytes}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsTimestamp,
				[]*checkedpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool),
			decls.NewOverload(overloads.GreaterEqualsDuration,
				[]*checkedpb.Type{decls.Duration, decls.Duration}, decls.Bool)),

		decls.NewFunction(operators.Equals,
			decls.NewParameterizedOverload(overloads.Equals,
				[]*checkedpb.Type{paramA, paramA}, decls.Bool,
				typeParamAList)),

		decls.NewFunction(operators.NotEquals,
			decls.NewParameterizedOverload(overloads.NotEquals,
				[]*checkedpb.Type{paramA, paramA}, decls.Bool,
				typeParamAList)),

		// Algebra

		decls.NewFunction(operators.Subtract,
			decls.NewOverload(overloads.SubtractInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.SubtractUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.SubtractDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Double),
			decls.NewOverload(overloads.SubtractTimestampTimestamp,
				[]*checkedpb.Type{decls.Timestamp, decls.Timestamp}, decls.Duration),
			decls.NewOverload(overloads.SubtractTimestampDuration,
				[]*checkedpb.Type{decls.Timestamp, decls.Duration}, decls.Timestamp),
			decls.NewOverload(overloads.SubtractDurationDuration,
				[]*checkedpb.Type{decls.Duration, decls.Duration}, decls.Duration)),

		decls.NewFunction(operators.Multiply,
			decls.NewOverload(overloads.MultiplyInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.MultiplyUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.MultiplyDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Double)),

		decls.NewFunction(operators.Divide,
			decls.NewOverload(overloads.DivideInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.DivideUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.DivideDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Double)),

		decls.NewFunction(operators.Modulo,
			decls.NewOverload(overloads.ModuloInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.ModuloUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Uint)),

		decls.NewFunction(operators.Add,
			decls.NewOverload(overloads.AddInt64,
				[]*checkedpb.Type{decls.Int, decls.Int}, decls.Int),
			decls.NewOverload(overloads.AddUint64,
				[]*checkedpb.Type{decls.Uint, decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.AddDouble,
				[]*checkedpb.Type{decls.Double, decls.Double}, decls.Double),
			decls.NewOverload(overloads.AddString,
				[]*checkedpb.Type{decls.String, decls.String}, decls.String),
			decls.NewOverload(overloads.AddBytes,
				[]*checkedpb.Type{decls.Bytes, decls.Bytes}, decls.Bytes),
			decls.NewParameterizedOverload(overloads.AddList,
				[]*checkedpb.Type{listOfA, listOfA}, listOfA,
				typeParamAList),
			decls.NewOverload(overloads.AddTimestampDuration,
				[]*checkedpb.Type{decls.Timestamp, decls.Duration}, decls.Timestamp),
			decls.NewOverload(overloads.AddDurationTimestamp,
				[]*checkedpb.Type{decls.Duration, decls.Timestamp}, decls.Timestamp),
			decls.NewOverload(overloads.AddDurationDuration,
				[]*checkedpb.Type{decls.Duration, decls.Duration}, decls.Duration)),

		decls.NewFunction(operators.Negate,
			decls.NewOverload(overloads.NegateInt64,
				[]*checkedpb.Type{decls.Int}, decls.Int),
			decls.NewOverload(overloads.NegateDouble,
				[]*checkedpb.Type{decls.Double}, decls.Double)),

		// Index

		decls.NewFunction(operators.Index,
			decls.NewParameterizedOverload(overloads.IndexList,
				[]*checkedpb.Type{listOfA, decls.Int}, paramA,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.IndexMap,
				[]*checkedpb.Type{mapOfAB, paramA}, paramB,
				typeParamABList)),
		//decls.NewOverload(overloads.IndexMessage,
		//	[]*checkedpb.Type{decls.Dyn, decls.String}, decls.Dyn)),

		// Collections

		decls.NewFunction(overloads.Size,
			decls.NewInstanceOverload(overloads.SizeStringInst,
				[]*checkedpb.Type{decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.SizeBytesInst,
				[]*checkedpb.Type{decls.Bytes}, decls.Int),
			decls.NewParameterizedInstanceOverload(overloads.SizeListInst,
				[]*checkedpb.Type{listOfA}, decls.Int, typeParamAList),
			decls.NewParameterizedInstanceOverload(overloads.SizeMapInst,
				[]*checkedpb.Type{mapOfAB}, decls.Int, typeParamABList),
			decls.NewOverload(overloads.SizeString,
				[]*checkedpb.Type{decls.String}, decls.Int),
			decls.NewOverload(overloads.SizeBytes,
				[]*checkedpb.Type{decls.Bytes}, decls.Int),
			decls.NewParameterizedOverload(overloads.SizeList,
				[]*checkedpb.Type{listOfA}, decls.Int, typeParamAList),
			decls.NewParameterizedOverload(overloads.SizeMap,
				[]*checkedpb.Type{mapOfAB}, decls.Int, typeParamABList)),

		decls.NewFunction(operators.In,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checkedpb.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checkedpb.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),

		// Deprecated 'in()' function

		decls.NewFunction(overloads.DeprecatedIn,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checkedpb.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checkedpb.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*checkedpb.Type{Dyn, decls.String},decls.Bool)),

		// Conversions to type

		decls.NewFunction(overloads.TypeConvertType,
			decls.NewParameterizedOverload(overloads.TypeConvertType,
				[]*checkedpb.Type{paramA}, decls.NewTypeType(paramA), typeParamAList)),

		// Conversions to int

		decls.NewFunction(overloads.TypeConvertInt,
			decls.NewOverload(overloads.IntToInt, []*checkedpb.Type{decls.Int}, decls.Int),
			decls.NewOverload(overloads.UintToInt, []*checkedpb.Type{decls.Uint}, decls.Int),
			decls.NewOverload(overloads.DoubleToInt, []*checkedpb.Type{decls.Double}, decls.Int),
			decls.NewOverload(overloads.StringToInt, []*checkedpb.Type{decls.String}, decls.Int),
			decls.NewOverload(overloads.TimestampToInt, []*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewOverload(overloads.DurationToInt, []*checkedpb.Type{decls.Duration}, decls.Int)),

		// Conversions to uint

		decls.NewFunction(overloads.TypeConvertUint,
			decls.NewOverload(overloads.UintToUint, []*checkedpb.Type{decls.Uint}, decls.Uint),
			decls.NewOverload(overloads.IntToUint, []*checkedpb.Type{decls.Int}, decls.Uint),
			decls.NewOverload(overloads.DoubleToUint, []*checkedpb.Type{decls.Double}, decls.Uint),
			decls.NewOverload(overloads.StringToUint, []*checkedpb.Type{decls.String}, decls.Uint)),

		// Conversions to double

		decls.NewFunction(overloads.TypeConvertDouble,
			decls.NewOverload(overloads.DoubleToDouble, []*checkedpb.Type{decls.Double}, decls.Double),
			decls.NewOverload(overloads.IntToDouble, []*checkedpb.Type{decls.Int}, decls.Double),
			decls.NewOverload(overloads.UintToDouble, []*checkedpb.Type{decls.Uint}, decls.Double),
			decls.NewOverload(overloads.StringToDouble, []*checkedpb.Type{decls.String}, decls.Double)),

		// Conversions to bool

		decls.NewFunction(overloads.TypeConvertBool,
			decls.NewOverload(overloads.BoolToBool, []*checkedpb.Type{decls.Bool}, decls.Bool),
			decls.NewOverload(overloads.StringToBool, []*checkedpb.Type{decls.String}, decls.Bool)),

		// Conversions to string

		decls.NewFunction(overloads.TypeConvertString,
			decls.NewOverload(overloads.StringToString, []*checkedpb.Type{decls.String}, decls.String),
			decls.NewOverload(overloads.BoolToString, []*checkedpb.Type{decls.Bool}, decls.String),
			decls.NewOverload(overloads.IntToString, []*checkedpb.Type{decls.Int}, decls.String),
			decls.NewOverload(overloads.UintToString, []*checkedpb.Type{decls.Uint}, decls.String),
			decls.NewOverload(overloads.DoubleToString, []*checkedpb.Type{decls.Double}, decls.String),
			decls.NewOverload(overloads.BytesToString, []*checkedpb.Type{decls.Bytes}, decls.String),
			decls.NewOverload(overloads.TimestampToString, []*checkedpb.Type{decls.Timestamp}, decls.String),
			decls.NewOverload(overloads.DurationToString, []*checkedpb.Type{decls.Duration}, decls.String)),

		// Conversions to bytes

		decls.NewFunction(overloads.TypeConvertBytes,
			decls.NewOverload(overloads.BytesToBytes, []*checkedpb.Type{decls.Bytes}, decls.Bytes),
			decls.NewOverload(overloads.StringToBytes, []*checkedpb.Type{decls.String}, decls.Bytes)),

		// Conversions to timestamps

		decls.NewFunction(overloads.TypeConvertTimestamp,
			decls.NewOverload(overloads.TimestampToTimestamp,
				[]*checkedpb.Type{decls.Timestamp}, decls.Timestamp),
			decls.NewOverload(overloads.StringToTimestamp,
				[]*checkedpb.Type{decls.String}, decls.Timestamp),
			decls.NewOverload(overloads.IntToTimestamp,
				[]*checkedpb.Type{decls.Int}, decls.Timestamp)),

		// Conversions to durations

		decls.NewFunction(overloads.TypeConvertDuration,
			decls.NewOverload(overloads.DurationToDuration,
				[]*checkedpb.Type{decls.Duration}, decls.Duration),
			decls.NewOverload(overloads.StringToDuration,
				[]*checkedpb.Type{decls.String}, decls.Duration),
			decls.NewOverload(overloads.IntToDuration,
				[]*checkedpb.Type{decls.Int}, decls.Duration)),

		// Conversions to Dyn

		decls.NewFunction(overloads.TypeConvertDyn,
			decls.NewParameterizedOverload(overloads.ToDyn,
				[]*checkedpb.Type{paramA}, decls.Dyn,
				typeParamAList)),

		// Date/time functions

		decls.NewFunction(overloads.TimeGetFullYear,
			decls.NewInstanceOverload(overloads.TimestampToYear,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToYearWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMonth,
			decls.NewInstanceOverload(overloads.TimestampToMonth,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMonthWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfYear,
			decls.NewInstanceOverload(overloads.TimestampToDayOfYear,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfYearWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfMonth,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBased,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDate,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBased,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetDayOfWeek,
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeek,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeekWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int)),

		decls.NewFunction(overloads.TimeGetHours,
			decls.NewInstanceOverload(overloads.TimestampToHours,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToHoursWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToHours,
				[]*checkedpb.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMinutes,
			decls.NewInstanceOverload(overloads.TimestampToMinutes,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMinutesWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToMinutes,
				[]*checkedpb.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetSeconds,
			decls.NewInstanceOverload(overloads.TimestampToSeconds,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToSecondsWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToSeconds,
				[]*checkedpb.Type{decls.Duration}, decls.Int)),

		decls.NewFunction(overloads.TimeGetMilliseconds,
			decls.NewInstanceOverload(overloads.TimestampToMilliseconds,
				[]*checkedpb.Type{decls.Timestamp}, decls.Int),
			decls.NewInstanceOverload(overloads.TimestampToMillisecondsWithTz,
				[]*checkedpb.Type{decls.Timestamp, decls.String}, decls.Int),
			decls.NewInstanceOverload(overloads.DurationToMilliseconds,
				[]*checkedpb.Type{decls.Duration}, decls.Int))}...)
}
