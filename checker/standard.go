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
	declspb "github.com/google/cel-go/checker/decls"
	operatorspb "github.com/google/cel-go/common/operators"
	overloadspb "github.com/google/cel-go/common/overloads"
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
)

func StandardDeclarations() []*checkedpb.Decl {
	// Some shortcuts we use when building declarations.
	paramA := declspb.NewTypeParamType("A")
	typeParamAList := []string{"A"}
	listOfA := declspb.NewListType(paramA)
	paramB := declspb.NewTypeParamType("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := declspb.NewMapType(paramA, paramB)

	var idents []*checkedpb.Decl
	for _, t := range []*checkedpb.Type{
		declspb.Int, declspb.Uint, declspb.Bool,
		declspb.Double, declspb.Bytes, declspb.String} {
		idents = append(idents,
			declspb.NewIdent(FormatCheckedType(t), declspb.NewTypeType(t), nil))
	}
	idents = append(idents,
		declspb.NewIdent("list", declspb.NewTypeType(listOfA), nil),
		declspb.NewIdent("map", declspb.NewTypeType(mapOfAB), nil),
		declspb.NewIdent("null_type", declspb.NewTypeType(declspb.Null), nil),
		declspb.NewIdent("type", declspb.NewTypeType(declspb.NewTypeType(nil)), nil))

	// Booleans
	// TODO: allow the conditional to return a heterogenous type.
	return append(idents, []*checkedpb.Decl{
		declspb.NewFunction(operatorspb.Conditional,
			declspb.NewParameterizedOverload(overloadspb.Conditional,
				[]*checkedpb.Type{declspb.Bool, paramA, paramA}, paramA,
				typeParamAList)),

		declspb.NewFunction(operatorspb.LogicalAnd,
			declspb.NewOverload(overloadspb.LogicalAnd,
				[]*checkedpb.Type{declspb.Bool, declspb.Bool}, declspb.Bool)),

		declspb.NewFunction(operatorspb.LogicalOr,
			declspb.NewOverload(overloadspb.LogicalOr,
				[]*checkedpb.Type{declspb.Bool, declspb.Bool}, declspb.Bool)),

		declspb.NewFunction(operatorspb.LogicalNot,
			declspb.NewOverload(overloadspb.LogicalNot,
				[]*checkedpb.Type{declspb.Bool}, declspb.Bool)),

		declspb.NewFunction(overloadspb.Matches,
			declspb.NewInstanceOverload(overloadspb.MatchString,
				[]*checkedpb.Type{declspb.String, declspb.String}, declspb.Bool)),

		// Relations

		declspb.NewFunction(operatorspb.Less,
			declspb.NewOverload(overloadspb.LessBool,
				[]*checkedpb.Type{declspb.Bool, declspb.Bool}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessString,
				[]*checkedpb.Type{declspb.String, declspb.String}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessBytes,
				[]*checkedpb.Type{declspb.Bytes, declspb.Bytes}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessTimestamp,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Timestamp}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessDuration,
				[]*checkedpb.Type{declspb.Duration, declspb.Duration}, declspb.Bool)),

		declspb.NewFunction(operatorspb.LessEquals,
			declspb.NewOverload(overloadspb.LessEqualsBool,
				[]*checkedpb.Type{declspb.Bool, declspb.Bool}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsString,
				[]*checkedpb.Type{declspb.String, declspb.String}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsBytes,
				[]*checkedpb.Type{declspb.Bytes, declspb.Bytes}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsTimestamp,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Timestamp}, declspb.Bool),
			declspb.NewOverload(overloadspb.LessEqualsDuration,
				[]*checkedpb.Type{declspb.Duration, declspb.Duration}, declspb.Bool)),

		declspb.NewFunction(operatorspb.Greater,
			declspb.NewOverload(overloadspb.GreaterBool,
				[]*checkedpb.Type{declspb.Bool, declspb.Bool}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterString,
				[]*checkedpb.Type{declspb.String, declspb.String}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterBytes,
				[]*checkedpb.Type{declspb.Bytes, declspb.Bytes}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterTimestamp,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Timestamp}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterDuration,
				[]*checkedpb.Type{declspb.Duration, declspb.Duration}, declspb.Bool)),

		declspb.NewFunction(operatorspb.GreaterEquals,
			declspb.NewOverload(overloadspb.GreaterEqualsBool,
				[]*checkedpb.Type{declspb.Bool, declspb.Bool}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsString,
				[]*checkedpb.Type{declspb.String, declspb.String}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsBytes,
				[]*checkedpb.Type{declspb.Bytes, declspb.Bytes}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsTimestamp,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Timestamp}, declspb.Bool),
			declspb.NewOverload(overloadspb.GreaterEqualsDuration,
				[]*checkedpb.Type{declspb.Duration, declspb.Duration}, declspb.Bool)),

		declspb.NewFunction(operatorspb.Equals,
			declspb.NewParameterizedOverload(overloadspb.Equals,
				[]*checkedpb.Type{paramA, paramA}, declspb.Bool,
				typeParamAList)),

		declspb.NewFunction(operatorspb.NotEquals,
			declspb.NewParameterizedOverload(overloadspb.NotEquals,
				[]*checkedpb.Type{paramA, paramA}, declspb.Bool,
				typeParamAList)),

		// Algebra

		declspb.NewFunction(operatorspb.Subtract,
			declspb.NewOverload(overloadspb.SubtractInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.SubtractUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Uint),
			declspb.NewOverload(overloadspb.SubtractDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Double),
			declspb.NewOverload(overloadspb.SubtractTimestampTimestamp,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Timestamp}, declspb.Duration),
			declspb.NewOverload(overloadspb.SubtractTimestampDuration,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Duration}, declspb.Timestamp),
			declspb.NewOverload(overloadspb.SubtractDurationDuration,
				[]*checkedpb.Type{declspb.Duration, declspb.Duration}, declspb.Duration)),

		declspb.NewFunction(operatorspb.Multiply,
			declspb.NewOverload(overloadspb.MultiplyInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.MultiplyUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Uint),
			declspb.NewOverload(overloadspb.MultiplyDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Double)),

		declspb.NewFunction(operatorspb.Divide,
			declspb.NewOverload(overloadspb.DivideInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.DivideUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Uint),
			declspb.NewOverload(overloadspb.DivideDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Double)),

		declspb.NewFunction(operatorspb.Modulo,
			declspb.NewOverload(overloadspb.ModuloInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.ModuloUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Uint)),

		declspb.NewFunction(operatorspb.Add,
			declspb.NewOverload(overloadspb.AddInt64,
				[]*checkedpb.Type{declspb.Int, declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.AddUint64,
				[]*checkedpb.Type{declspb.Uint, declspb.Uint}, declspb.Uint),
			declspb.NewOverload(overloadspb.AddDouble,
				[]*checkedpb.Type{declspb.Double, declspb.Double}, declspb.Double),
			declspb.NewOverload(overloadspb.AddString,
				[]*checkedpb.Type{declspb.String, declspb.String}, declspb.String),
			declspb.NewOverload(overloadspb.AddBytes,
				[]*checkedpb.Type{declspb.Bytes, declspb.Bytes}, declspb.Bytes),
			declspb.NewParameterizedOverload(overloadspb.AddList,
				[]*checkedpb.Type{listOfA, listOfA}, listOfA,
				typeParamAList),
			declspb.NewOverload(overloadspb.AddTimestampDuration,
				[]*checkedpb.Type{declspb.Timestamp, declspb.Duration}, declspb.Timestamp),
			declspb.NewOverload(overloadspb.AddDurationTimestamp,
				[]*checkedpb.Type{declspb.Duration, declspb.Timestamp}, declspb.Timestamp),
			declspb.NewOverload(overloadspb.AddDurationDuration,
				[]*checkedpb.Type{declspb.Duration, declspb.Duration}, declspb.Duration)),

		declspb.NewFunction(operatorspb.Negate,
			declspb.NewOverload(overloadspb.NegateInt64,
				[]*checkedpb.Type{declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.NegateDouble,
				[]*checkedpb.Type{declspb.Double}, declspb.Double)),

		// Index

		declspb.NewFunction(operatorspb.Index,
			declspb.NewParameterizedOverload(overloadspb.IndexList,
				[]*checkedpb.Type{listOfA, declspb.Int}, paramA,
				typeParamAList),
			declspb.NewParameterizedOverload(overloadspb.IndexMap,
				[]*checkedpb.Type{mapOfAB, paramA}, paramB,
				typeParamABList)),
		//declspb.NewOverload(overloadspb.IndexMessage,
		//	[]*checkedpb.Type{declspb.Dyn, declspb.String}, declspb.Dyn)),

		// Collections

		declspb.NewFunction(overloadspb.Size,
			declspb.NewInstanceOverload(overloadspb.SizeStringInst,
				[]*checkedpb.Type{declspb.String}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.SizeBytesInst,
				[]*checkedpb.Type{declspb.Bytes}, declspb.Int),
			declspb.NewParameterizedInstanceOverload(overloadspb.SizeListInst,
				[]*checkedpb.Type{listOfA}, declspb.Int, typeParamAList),
			declspb.NewParameterizedInstanceOverload(overloadspb.SizeMapInst,
				[]*checkedpb.Type{mapOfAB}, declspb.Int, typeParamABList),
			declspb.NewOverload(overloadspb.SizeString,
				[]*checkedpb.Type{declspb.String}, declspb.Int),
			declspb.NewOverload(overloadspb.SizeBytes,
				[]*checkedpb.Type{declspb.Bytes}, declspb.Int),
			declspb.NewParameterizedOverload(overloadspb.SizeList,
				[]*checkedpb.Type{listOfA}, declspb.Int, typeParamAList),
			declspb.NewParameterizedOverload(overloadspb.SizeMap,
				[]*checkedpb.Type{mapOfAB}, declspb.Int, typeParamABList)),

		declspb.NewFunction(operatorspb.In,
			declspb.NewParameterizedOverload(overloadspb.InList,
				[]*checkedpb.Type{paramA, listOfA}, declspb.Bool,
				typeParamAList),
			declspb.NewParameterizedOverload(overloadspb.InMap,
				[]*checkedpb.Type{paramA, mapOfAB}, declspb.Bool,
				typeParamABList)),

		// Deprecated 'in()' function

		declspb.NewFunction(overloadspb.DeprecatedIn,
			declspb.NewParameterizedOverload(overloadspb.InList,
				[]*checkedpb.Type{paramA, listOfA}, declspb.Bool,
				typeParamAList),
			declspb.NewParameterizedOverload(overloadspb.InMap,
				[]*checkedpb.Type{paramA, mapOfAB}, declspb.Bool,
				typeParamABList)),
		//declspb.NewOverload(overloadspb.InMessage,
		//	[]*checkedpb.Type{Dyn, declspb.String},declspb.Bool)),

		// Conversions to type

		declspb.NewFunction(overloadspb.TypeConvertType,
			declspb.NewParameterizedOverload(overloadspb.TypeConvertType,
				[]*checkedpb.Type{paramA}, declspb.NewTypeType(paramA), typeParamAList)),

		// Conversions to int

		declspb.NewFunction(overloadspb.TypeConvertInt,
			declspb.NewOverload(overloadspb.IntToInt, []*checkedpb.Type{declspb.Int}, declspb.Int),
			declspb.NewOverload(overloadspb.UintToInt, []*checkedpb.Type{declspb.Uint}, declspb.Int),
			declspb.NewOverload(overloadspb.DoubleToInt, []*checkedpb.Type{declspb.Double}, declspb.Int),
			declspb.NewOverload(overloadspb.StringToInt, []*checkedpb.Type{declspb.String}, declspb.Int),
			declspb.NewOverload(overloadspb.TimestampToInt, []*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewOverload(overloadspb.DurationToInt, []*checkedpb.Type{declspb.Duration}, declspb.Int)),

		// Conversions to uint

		declspb.NewFunction(overloadspb.TypeConvertUint,
			declspb.NewOverload(overloadspb.UintToUint, []*checkedpb.Type{declspb.Uint}, declspb.Uint),
			declspb.NewOverload(overloadspb.IntToUint, []*checkedpb.Type{declspb.Int}, declspb.Uint),
			declspb.NewOverload(overloadspb.DoubleToUint, []*checkedpb.Type{declspb.Double}, declspb.Uint),
			declspb.NewOverload(overloadspb.StringToUint, []*checkedpb.Type{declspb.String}, declspb.Uint)),

		// Conversions to double

		declspb.NewFunction(overloadspb.TypeConvertDouble,
			declspb.NewOverload(overloadspb.DoubleToDouble, []*checkedpb.Type{declspb.Double}, declspb.Double),
			declspb.NewOverload(overloadspb.IntToDouble, []*checkedpb.Type{declspb.Int}, declspb.Double),
			declspb.NewOverload(overloadspb.UintToDouble, []*checkedpb.Type{declspb.Uint}, declspb.Double),
			declspb.NewOverload(overloadspb.StringToDouble, []*checkedpb.Type{declspb.String}, declspb.Double)),

		// Conversions to bool

		declspb.NewFunction(overloadspb.TypeConvertBool,
			declspb.NewOverload(overloadspb.BoolToBool, []*checkedpb.Type{declspb.Bool}, declspb.Bool),
			declspb.NewOverload(overloadspb.StringToBool, []*checkedpb.Type{declspb.String}, declspb.Bool)),

		// Conversions to string

		declspb.NewFunction(overloadspb.TypeConvertString,
			declspb.NewOverload(overloadspb.StringToString, []*checkedpb.Type{declspb.String}, declspb.String),
			declspb.NewOverload(overloadspb.BoolToString, []*checkedpb.Type{declspb.Bool}, declspb.String),
			declspb.NewOverload(overloadspb.IntToString, []*checkedpb.Type{declspb.Int}, declspb.String),
			declspb.NewOverload(overloadspb.UintToString, []*checkedpb.Type{declspb.Uint}, declspb.String),
			declspb.NewOverload(overloadspb.DoubleToString, []*checkedpb.Type{declspb.Double}, declspb.String),
			declspb.NewOverload(overloadspb.BytesToString, []*checkedpb.Type{declspb.Bytes}, declspb.String),
			declspb.NewOverload(overloadspb.TimestampToString, []*checkedpb.Type{declspb.Timestamp}, declspb.String),
			declspb.NewOverload(overloadspb.DurationToString, []*checkedpb.Type{declspb.Duration}, declspb.String)),

		// Conversions to bytes

		declspb.NewFunction(overloadspb.TypeConvertBytes,
			declspb.NewOverload(overloadspb.BytesToBytes, []*checkedpb.Type{declspb.Bytes}, declspb.Bytes),
			declspb.NewOverload(overloadspb.StringToBytes, []*checkedpb.Type{declspb.String}, declspb.Bytes)),

		// Conversions to timestamps

		declspb.NewFunction(overloadspb.TypeConvertTimestamp,
			declspb.NewOverload(overloadspb.TimestampToTimestamp,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Timestamp),
			declspb.NewOverload(overloadspb.StringToTimestamp,
				[]*checkedpb.Type{declspb.String}, declspb.Timestamp),
			declspb.NewOverload(overloadspb.IntToTimestamp,
				[]*checkedpb.Type{declspb.Int}, declspb.Timestamp)),

		// Conversions to durations

		declspb.NewFunction(overloadspb.TypeConvertDuration,
			declspb.NewOverload(overloadspb.DurationToDuration,
				[]*checkedpb.Type{declspb.Duration}, declspb.Duration),
			declspb.NewOverload(overloadspb.StringToDuration,
				[]*checkedpb.Type{declspb.String}, declspb.Duration),
			declspb.NewOverload(overloadspb.IntToDuration,
				[]*checkedpb.Type{declspb.Int}, declspb.Duration)),

		// Conversions to Dyn

		declspb.NewFunction(overloadspb.TypeConvertDyn,
			declspb.NewParameterizedOverload(overloadspb.ToDyn,
				[]*checkedpb.Type{paramA}, declspb.Dyn,
				typeParamAList)),

		// Date/time functions

		declspb.NewFunction(overloadspb.TimeGetFullYear,
			declspb.NewInstanceOverload(overloadspb.TimestampToYear,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToYearWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetMonth,
			declspb.NewInstanceOverload(overloadspb.TimestampToMonth,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToMonthWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetDayOfYear,
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfYear,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfYearWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetDayOfMonth,
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfMonthZeroBased,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfMonthZeroBasedWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetDate,
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfMonthOneBased,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfMonthOneBasedWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetDayOfWeek,
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfWeek,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToDayOfWeekWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetHours,
			declspb.NewInstanceOverload(overloadspb.TimestampToHours,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToHoursWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.DurationToHours,
				[]*checkedpb.Type{declspb.Duration}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetMinutes,
			declspb.NewInstanceOverload(overloadspb.TimestampToMinutes,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToMinutesWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.DurationToMinutes,
				[]*checkedpb.Type{declspb.Duration}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetSeconds,
			declspb.NewInstanceOverload(overloadspb.TimestampToSeconds,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToSecondsWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.DurationToSeconds,
				[]*checkedpb.Type{declspb.Duration}, declspb.Int)),

		declspb.NewFunction(overloadspb.TimeGetMilliseconds,
			declspb.NewInstanceOverload(overloadspb.TimestampToMilliseconds,
				[]*checkedpb.Type{declspb.Timestamp}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.TimestampToMillisecondsWithTz,
				[]*checkedpb.Type{declspb.Timestamp, declspb.String}, declspb.Int),
			declspb.NewInstanceOverload(overloadspb.DurationToMilliseconds,
				[]*checkedpb.Type{declspb.Duration}, declspb.Int))}...)
}
