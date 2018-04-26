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
	paramA := newTypeParam("A")
	typeParamAList := []string{"A"}
	listOfA := newList(paramA)
	paramB := newTypeParam("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := newMap(paramA, paramB)

	var idents []*checked.Decl
	for _, t := range []*checked.Type{Int, Uint, Bool, Double, Bytes, String} {
		idents = append(idents, decls.NewIdent(FormatCheckedType(t), newType(t), nil))
	}
	idents = append(idents,
		decls.NewIdent("list", newType(listOfA), nil),
		decls.NewIdent("map", newType(mapOfAB), nil))

	// Booleans
	// TODO: allow the conditional to return a heterogenous type.
	return append(idents, []*checked.Decl{
		decls.NewFunction(operators.Conditional,
			decls.NewParameterizedOverload(overloads.Conditional,
				[]*checked.Type{Bool, paramA, paramA}, paramA,
				typeParamAList)),

		decls.NewFunction(operators.LogicalAnd,
			decls.NewOverload(overloads.LogicalAnd,
				[]*checked.Type{Bool, Bool}, Bool)),

		decls.NewFunction(operators.LogicalOr,
			decls.NewOverload(overloads.LogicalOr,
				[]*checked.Type{Bool, Bool}, Bool)),

		decls.NewFunction(operators.LogicalNot,
			decls.NewOverload(overloads.LogicalNot,
				[]*checked.Type{Bool}, Bool)),

		decls.NewFunction(overloads.Matches,
			decls.NewInstanceOverload(overloads.MatchString,
				[]*checked.Type{String, String}, Bool)),

		// Relations

		decls.NewFunction(operators.Less,
			decls.NewOverload(overloads.LessBool,
				[]*checked.Type{Bool, Bool}, Bool),
			decls.NewOverload(overloads.LessInt64,
				[]*checked.Type{Int, Int}, Bool),
			decls.NewOverload(overloads.LessUint64,
				[]*checked.Type{Uint, Uint}, Bool),
			decls.NewOverload(overloads.LessDouble,
				[]*checked.Type{Double, Double}, Bool),
			decls.NewOverload(overloads.LessString,
				[]*checked.Type{String, String}, Bool),
			decls.NewOverload(overloads.LessBytes,
				[]*checked.Type{Bytes, Bytes}, Bool),
			decls.NewOverload(overloads.LessTimestamp,
				[]*checked.Type{Timestamp, Timestamp}, Bool),
			decls.NewOverload(overloads.LessDuration,
				[]*checked.Type{Duration, Duration}, Bool)),

		decls.NewFunction(operators.LessEquals,
			decls.NewOverload(overloads.LessEqualsBool,
				[]*checked.Type{Bool, Bool}, Bool),
			decls.NewOverload(overloads.LessEqualsInt64,
				[]*checked.Type{Int, Int}, Bool),
			decls.NewOverload(overloads.LessEqualsUint64,
				[]*checked.Type{Uint, Uint}, Bool),
			decls.NewOverload(overloads.LessEqualsDouble,
				[]*checked.Type{Double, Double}, Bool),
			decls.NewOverload(overloads.LessEqualsString,
				[]*checked.Type{String, String}, Bool),
			decls.NewOverload(overloads.LessEqualsBytes,
				[]*checked.Type{Bytes, Bytes}, Bool),
			decls.NewOverload(overloads.LessEqualsTimestamp,
				[]*checked.Type{Timestamp, Timestamp}, Bool),
			decls.NewOverload(overloads.LessEqualsDuration,
				[]*checked.Type{Duration, Duration}, Bool)),

		decls.NewFunction(operators.Greater,
			decls.NewOverload(overloads.GreaterBool,
				[]*checked.Type{Bool, Bool}, Bool),
			decls.NewOverload(overloads.GreaterInt64,
				[]*checked.Type{Int, Int}, Bool),
			decls.NewOverload(overloads.GreaterUint64,
				[]*checked.Type{Uint, Uint}, Bool),
			decls.NewOverload(overloads.GreaterDouble,
				[]*checked.Type{Double, Double}, Bool),
			decls.NewOverload(overloads.GreaterString,
				[]*checked.Type{String, String}, Bool),
			decls.NewOverload(overloads.GreaterBytes,
				[]*checked.Type{Bytes, Bytes}, Bool),
			decls.NewOverload(overloads.GreaterTimestamp,
				[]*checked.Type{Timestamp, Timestamp}, Bool),
			decls.NewOverload(overloads.GreaterDuration,
				[]*checked.Type{Duration, Duration}, Bool)),

		decls.NewFunction(operators.GreaterEquals,
			decls.NewOverload(overloads.GreaterEqualsBool,
				[]*checked.Type{Bool, Bool}, Bool),
			decls.NewOverload(overloads.GreaterEqualsInt64,
				[]*checked.Type{Int, Int}, Bool),
			decls.NewOverload(overloads.GreaterEqualsUint64,
				[]*checked.Type{Uint, Uint}, Bool),
			decls.NewOverload(overloads.GreaterEqualsDouble,
				[]*checked.Type{Double, Double}, Bool),
			decls.NewOverload(overloads.GreaterEqualsString,
				[]*checked.Type{String, String}, Bool),
			decls.NewOverload(overloads.GreaterEqualsBytes,
				[]*checked.Type{Bytes, Bytes}, Bool),
			decls.NewOverload(overloads.GreaterEqualsTimestamp,
				[]*checked.Type{Timestamp, Timestamp}, Bool),
			decls.NewOverload(overloads.GreaterEqualsDuration,
				[]*checked.Type{Duration, Duration}, Bool)),

		decls.NewFunction(operators.Equals,
			decls.NewParameterizedOverload(overloads.Equals,
				[]*checked.Type{paramA, paramA}, Bool,
				typeParamAList)),

		decls.NewFunction(operators.NotEquals,
			decls.NewParameterizedOverload(overloads.NotEquals,
				[]*checked.Type{paramA, paramA}, Bool,
				typeParamAList)),

		// Algebra

		decls.NewFunction(operators.Subtract,
			decls.NewOverload(overloads.SubtractInt64,
				[]*checked.Type{Int, Int}, Int),
			decls.NewOverload(overloads.SubtractUint64,
				[]*checked.Type{Uint, Uint}, Uint),
			decls.NewOverload(overloads.SubtractDouble,
				[]*checked.Type{Double, Double}, Double),
			decls.NewOverload(overloads.SubtractTimestampTimestamp,
				[]*checked.Type{Timestamp, Timestamp}, Duration),
			decls.NewOverload(overloads.SubtractTimestampDuration,
				[]*checked.Type{Timestamp, Duration}, Timestamp),
			decls.NewOverload(overloads.SubtractDurationDuration,
				[]*checked.Type{Duration, Duration}, Duration)),

		decls.NewFunction(operators.Multiply,
			decls.NewOverload(overloads.MultiplyInt64,
				[]*checked.Type{Int, Int}, Int),
			decls.NewOverload(overloads.MultiplyUint64,
				[]*checked.Type{Uint, Uint}, Uint),
			decls.NewOverload(overloads.MultiplyDouble,
				[]*checked.Type{Double, Double}, Double)),

		decls.NewFunction(operators.Divide,
			decls.NewOverload(overloads.DivideInt64,
				[]*checked.Type{Int, Int}, Int),
			decls.NewOverload(overloads.DivideUint64,
				[]*checked.Type{Uint, Uint}, Uint),
			decls.NewOverload(overloads.DivideDouble,
				[]*checked.Type{Double, Double}, Double)),

		decls.NewFunction(operators.Modulo,
			decls.NewOverload(overloads.ModuloInt64,
				[]*checked.Type{Int, Int}, Int),
			decls.NewOverload(overloads.ModuloUint64,
				[]*checked.Type{Uint, Uint}, Uint)),

		decls.NewFunction(operators.Add,
			decls.NewOverload(overloads.AddInt64,
				[]*checked.Type{Int, Int}, Int),
			decls.NewOverload(overloads.AddUint64,
				[]*checked.Type{Uint, Uint}, Uint),
			decls.NewOverload(overloads.AddDouble,
				[]*checked.Type{Double, Double}, Double),
			decls.NewOverload(overloads.AddString,
				[]*checked.Type{String, String}, String),
			decls.NewOverload(overloads.AddBytes,
				[]*checked.Type{Bytes, Bytes}, Bytes),
			decls.NewParameterizedOverload(overloads.AddList,
				[]*checked.Type{listOfA, listOfA}, listOfA,
				typeParamAList),
			decls.NewOverload(overloads.AddTimestampDuration,
				[]*checked.Type{Timestamp, Duration}, Timestamp),
			decls.NewOverload(overloads.AddDurationTimestamp,
				[]*checked.Type{Duration, Timestamp}, Timestamp),
			decls.NewOverload(overloads.AddDurationDuration,
				[]*checked.Type{Duration, Duration}, Duration)),

		decls.NewFunction(operators.Negate,
			decls.NewOverload(overloads.NegateInt64,
				[]*checked.Type{Int}, Int),
			decls.NewOverload(overloads.NegateDouble,
				[]*checked.Type{Double}, Double)),

		// Index

		decls.NewFunction(operators.Index,
			decls.NewParameterizedOverload(overloads.IndexList,
				[]*checked.Type{listOfA, Int}, paramA,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.IndexMap,
				[]*checked.Type{mapOfAB, paramA}, paramB,
				typeParamABList)),
		//decls.NewOverload(overloads.IndexMessage,
		//	[]*checked.Type{Dyn, String}, Dyn)),

		// Collections

		decls.NewFunction(overloads.Size,
			decls.NewInstanceOverload(overloads.SizeStringInst,
				[]*checked.Type{String}, Int),
			decls.NewInstanceOverload(overloads.SizeBytesInst,
				[]*checked.Type{Bytes}, Int),
			decls.NewParameterizedInstanceOverload(overloads.SizeListInst,
				[]*checked.Type{listOfA}, Int, typeParamAList),
			decls.NewParameterizedInstanceOverload(overloads.SizeMapInst,
				[]*checked.Type{mapOfAB}, Int, typeParamABList),
			decls.NewOverload(overloads.SizeString,
				[]*checked.Type{String}, Int),
			decls.NewOverload(overloads.SizeBytes,
				[]*checked.Type{Bytes}, Int),
			decls.NewParameterizedOverload(overloads.SizeList,
				[]*checked.Type{listOfA}, Int, typeParamAList),
			decls.NewParameterizedOverload(overloads.SizeMap,
				[]*checked.Type{mapOfAB}, Int, typeParamABList)),

		decls.NewFunction(operators.In,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checked.Type{paramA, listOfA}, Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checked.Type{paramA, mapOfAB}, Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*checked.Type{Dyn, String},Bool)),

		// Deprecated 'in()' function

		decls.NewFunction(overloads.DeprecatedIn,
			decls.NewParameterizedOverload(overloads.InList,
				[]*checked.Type{paramA, listOfA}, Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*checked.Type{paramA, mapOfAB}, Bool,
				typeParamABList)),
		//decls.NewOverload(overloads.InMessage,
		//	[]*checked.Type{Dyn, String},Bool)),

		// Conversions to type

		decls.NewFunction(overloads.TypeConvertType,
			decls.NewParameterizedOverload(overloads.TypeConvertType,
				[]*checked.Type{paramA}, newType(paramA), typeParamAList)),

		// Conversions to int

		decls.NewFunction(overloads.TypeConvertInt,
			decls.NewOverload(overloads.IntToInt, []*checked.Type{Int}, Int),
			decls.NewOverload(overloads.UintToInt, []*checked.Type{Uint}, Int),
			decls.NewOverload(overloads.DoubleToInt, []*checked.Type{Double}, Int),
			decls.NewOverload(overloads.StringToInt, []*checked.Type{String}, Int),
			decls.NewOverload(overloads.TimestampToInt, []*checked.Type{Timestamp}, Int),
			decls.NewOverload(overloads.DurationToInt, []*checked.Type{Duration}, Int)),

		// Conversions to uint

		decls.NewFunction(overloads.TypeConvertUint,
			decls.NewOverload(overloads.UintToUint, []*checked.Type{Uint}, Uint),
			decls.NewOverload(overloads.IntToUint, []*checked.Type{Int}, Uint),
			decls.NewOverload(overloads.DoubleToUint, []*checked.Type{Double}, Uint),
			decls.NewOverload(overloads.StringToUint, []*checked.Type{String}, Uint)),

		// Conversions to double

		decls.NewFunction(overloads.TypeConvertDouble,
			decls.NewOverload(overloads.DoubleToDouble, []*checked.Type{Double}, Double),
			decls.NewOverload(overloads.IntToDouble, []*checked.Type{Int}, Double),
			decls.NewOverload(overloads.UintToDouble, []*checked.Type{Uint}, Double),
			decls.NewOverload(overloads.StringToDouble, []*checked.Type{String}, Double)),

		// Conversions to bool

		decls.NewFunction(overloads.TypeConvertBool,
			decls.NewOverload(overloads.BoolToBool, []*checked.Type{Bool}, Bool),
			decls.NewOverload(overloads.StringToBool, []*checked.Type{String}, Bool)),

		// Conversions to string

		decls.NewFunction(overloads.TypeConvertString,
			decls.NewOverload(overloads.StringToString, []*checked.Type{String}, String),
			decls.NewOverload(overloads.BoolToString, []*checked.Type{Bool}, String),
			decls.NewOverload(overloads.IntToString, []*checked.Type{Int}, String),
			decls.NewOverload(overloads.UintToString, []*checked.Type{Uint}, String),
			decls.NewOverload(overloads.DoubleToString, []*checked.Type{Double}, String),
			decls.NewOverload(overloads.BytesToString, []*checked.Type{Bytes}, String),
			decls.NewOverload(overloads.TimestampToString, []*checked.Type{Timestamp}, String),
			decls.NewOverload(overloads.DurationToString, []*checked.Type{Duration}, String)),

		// Conversions to bytes

		decls.NewFunction(overloads.TypeConvertBytes,
			decls.NewOverload(overloads.BytesToBytes, []*checked.Type{Bytes}, Bytes),
			decls.NewOverload(overloads.StringToBytes, []*checked.Type{String}, Bytes)),

		// Conversions to timestamps

		decls.NewFunction(overloads.TypeConvertTimestamp,
			decls.NewOverload(overloads.TimestampToTimestamp,
				[]*checked.Type{Timestamp}, Timestamp),
			decls.NewOverload(overloads.StringToTimestamp,
				[]*checked.Type{String}, Timestamp),
			decls.NewOverload(overloads.IntToTimestamp,
				[]*checked.Type{Int}, Timestamp)),

		// Conversions to durations

		decls.NewFunction(overloads.TypeConvertDuration,
			decls.NewOverload(overloads.DurationToDuration,
				[]*checked.Type{Duration}, Duration),
			decls.NewOverload(overloads.StringToDuration,
				[]*checked.Type{String}, Duration),
			decls.NewOverload(overloads.IntToDuration,
				[]*checked.Type{Int}, Duration)),

		// Conversions to Dyn

		decls.NewFunction(overloads.TypeConvertDyn,
			decls.NewParameterizedOverload(overloads.ToDyn, []*checked.Type{paramA}, Dyn,
				typeParamAList)),

		// Date/time functions

		decls.NewFunction(overloads.TimeGetFullYear,
			decls.NewInstanceOverload(overloads.TimestampToYear,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToYearWithTz,
				[]*checked.Type{Timestamp, String}, Int)),

		decls.NewFunction(overloads.TimeGetMonth,
			decls.NewInstanceOverload(overloads.TimestampToMonth,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToMonthWithTz,
				[]*checked.Type{Timestamp, String}, Int)),

		decls.NewFunction(overloads.TimeGetDayOfYear,
			decls.NewInstanceOverload(overloads.TimestampToDayOfYear,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfYearWithTz,
				[]*checked.Type{Timestamp, String}, Int)),

		decls.NewFunction(overloads.TimeGetDayOfMonth,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBased,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				[]*checked.Type{Timestamp, String}, Int)),

		decls.NewFunction(overloads.TimeGetDate,
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBased,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				[]*checked.Type{Timestamp, String}, Int)),

		decls.NewFunction(overloads.TimeGetDayOfWeek,
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeek,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToDayOfWeekWithTz,
				[]*checked.Type{Timestamp, String}, Int)),

		decls.NewFunction(overloads.TimeGetHours,
			decls.NewInstanceOverload(overloads.TimestampToHours,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToHoursWithTz,
				[]*checked.Type{Timestamp, String}, Int),
			decls.NewInstanceOverload(overloads.DurationToHours,
				[]*checked.Type{Duration}, Int)),

		decls.NewFunction(overloads.TimeGetMinutes,
			decls.NewInstanceOverload(overloads.TimestampToMinutes,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToMinutesWithTz,
				[]*checked.Type{Timestamp, String}, Int),
			decls.NewInstanceOverload(overloads.DurationToMinutes,
				[]*checked.Type{Duration}, Int)),

		decls.NewFunction(overloads.TimeGetSeconds,
			decls.NewInstanceOverload(overloads.TimestampToSeconds,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToSecondsWithTz,
				[]*checked.Type{Timestamp, String}, Int),
			decls.NewInstanceOverload(overloads.DurationToSeconds,
				[]*checked.Type{Duration}, Int)),

		decls.NewFunction(overloads.TimeGetMilliseconds,
			decls.NewInstanceOverload(overloads.TimestampToMilliseconds,
				[]*checked.Type{Timestamp}, Int),
			decls.NewInstanceOverload(overloads.TimestampToMillisecondsWithTz,
				[]*checked.Type{Timestamp, String}, Int),
			decls.NewInstanceOverload(overloads.DurationToMilliseconds,
				[]*checked.Type{Duration}, Int))}...)
}
