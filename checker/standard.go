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

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	standardDeclarations []*exprpb.Decl
)

func addOverloadDoc(ol *exprpb.Decl_FunctionDecl_Overload, doc string) *exprpb.Decl_FunctionDecl_Overload {
	ol.Doc = doc
	return ol
}

func init() {
	// Some shortcuts we use when building declarations.
	paramA := decls.NewTypeParamType("A")
	typeParamAList := []string{"A"}
	listOfA := decls.NewListType(paramA)
	paramB := decls.NewTypeParamType("B")
	typeParamABList := []string{"A", "B"}
	mapOfAB := decls.NewMapType(paramA, paramB)

	var idents []*exprpb.Decl
	for _, t := range []*exprpb.Type{
		decls.Int, decls.Uint, decls.Bool,
		decls.Double, decls.Bytes, decls.String} {
		idents = append(idents,
			decls.NewVar(FormatCheckedType(t), decls.NewTypeType(t)))
	}
	idents = append(idents,
		decls.NewVar("list", decls.NewTypeType(listOfA)),
		decls.NewVar("map", decls.NewTypeType(mapOfAB)),
		decls.NewVar("null_type", decls.NewTypeType(decls.Null)),
		decls.NewVar("type", decls.NewTypeType(decls.NewTypeType(nil))))

	standardDeclarations = append(standardDeclarations, idents...)
	standardDeclarations = append(standardDeclarations, []*exprpb.Decl{
		// Booleans
		decls.NewFunction(operators.Conditional,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.Conditional,
				[]*exprpb.Type{decls.Bool, paramA, paramA}, paramA,
				typeParamAList), "The conditional operator. See above for evaluation semantics. Will evaluate the test and only one of the remaining sub-expressions.")),

		decls.NewFunction(operators.LogicalAnd,
			addOverloadDoc(decls.NewOverload(overloads.LogicalAnd,
				[]*exprpb.Type{decls.Bool, decls.Bool}, decls.Bool), "logical and")),

		decls.NewFunction(operators.LogicalOr,
			addOverloadDoc(decls.NewOverload(overloads.LogicalOr,
				[]*exprpb.Type{decls.Bool, decls.Bool}, decls.Bool), "logical or")),

		decls.NewFunction(operators.LogicalNot,
			addOverloadDoc(decls.NewOverload(overloads.LogicalNot,
				[]*exprpb.Type{decls.Bool}, decls.Bool), "logical not")),

		decls.NewFunction(operators.NotStrictlyFalse,
			decls.NewOverload(overloads.NotStrictlyFalse,
				[]*exprpb.Type{decls.Bool}, decls.Bool)),

		decls.NewFunction(operators.Equals,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.Equals,
				[]*exprpb.Type{paramA, paramA}, decls.Bool,
				typeParamAList), "equality")),

		decls.NewFunction(operators.NotEquals,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.NotEquals,
				[]*exprpb.Type{paramA, paramA}, decls.Bool,
				typeParamAList), "inequality")),

		// Algebra.

		decls.NewFunction(operators.Subtract,
			addOverloadDoc(decls.NewOverload(overloads.SubtractInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Int), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.SubtractUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Uint), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.SubtractDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Double), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.SubtractTimestampTimestamp,
				[]*exprpb.Type{decls.Timestamp, decls.Timestamp}, decls.Duration), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.SubtractTimestampDuration,
				[]*exprpb.Type{decls.Timestamp, decls.Duration}, decls.Timestamp), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.SubtractDurationDuration,
				[]*exprpb.Type{decls.Duration, decls.Duration}, decls.Duration), "arithmetic")),

		decls.NewFunction(operators.Multiply,
			addOverloadDoc(decls.NewOverload(overloads.MultiplyInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Int), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.MultiplyUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Uint), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.MultiplyDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Double), "arithmetic")),

		decls.NewFunction(operators.Divide,
			addOverloadDoc(decls.NewOverload(overloads.DivideInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Int), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.DivideUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Uint), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.DivideDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Double), "arithmetic")),

		decls.NewFunction(operators.Modulo,
			addOverloadDoc(decls.NewOverload(overloads.ModuloInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Int), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.ModuloUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Uint), "arithmetic")),

		decls.NewFunction(operators.Add,
			addOverloadDoc(decls.NewOverload(overloads.AddInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Int), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.AddUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Uint), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.AddDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Double), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.AddString,
				[]*exprpb.Type{decls.String, decls.String}, decls.String), "string concatenation. Space and time cost proportional to the sum of the input sizes."),
			addOverloadDoc(decls.NewOverload(overloads.AddBytes,
				[]*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bytes), "bytes concatenation"),
			addOverloadDoc(decls.NewParameterizedOverload(overloads.AddList,
				[]*exprpb.Type{listOfA, listOfA}, listOfA,
				typeParamAList), "list concatenation. Space and time cost proportional to the sum of the input sizes."),
			addOverloadDoc(decls.NewOverload(overloads.AddTimestampDuration,
				[]*exprpb.Type{decls.Timestamp, decls.Duration}, decls.Timestamp), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.AddDurationTimestamp,
				[]*exprpb.Type{decls.Duration, decls.Timestamp}, decls.Timestamp), "arithmetic"),
			addOverloadDoc(decls.NewOverload(overloads.AddDurationDuration,
				[]*exprpb.Type{decls.Duration, decls.Duration}, decls.Duration), "arithmetic")),

		decls.NewFunction(operators.Negate,
			addOverloadDoc(decls.NewOverload(overloads.NegateInt64,
				[]*exprpb.Type{decls.Int}, decls.Int), "negation"),
			addOverloadDoc(decls.NewOverload(overloads.NegateDouble,
				[]*exprpb.Type{decls.Double}, decls.Double), "negation")),

		// Index.

		decls.NewFunction(operators.Index,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.IndexList,
				[]*exprpb.Type{listOfA, decls.Int}, paramA,
				typeParamAList), "list indexing. Constant time cost. "),
			addOverloadDoc(decls.NewParameterizedOverload(overloads.IndexMap,
				[]*exprpb.Type{mapOfAB, paramA}, paramB,
				typeParamABList), "map indexing. For string keys, cost is proportional to the size of the map keys times the size of the index string. ")),

		// Collections.

		decls.NewFunction(overloads.Size,
			addOverloadDoc(decls.NewInstanceOverload(overloads.SizeStringInst,
				[]*exprpb.Type{decls.String}, decls.Int), "string size"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.SizeBytesInst,
				[]*exprpb.Type{decls.Bytes}, decls.Int), "length size"),
			addOverloadDoc(decls.NewParameterizedInstanceOverload(overloads.SizeListInst,
				[]*exprpb.Type{listOfA}, decls.Int, typeParamAList), "list size. Time cost proportional to the length of the list."),
			addOverloadDoc(decls.NewParameterizedInstanceOverload(overloads.SizeMapInst,
				[]*exprpb.Type{mapOfAB}, decls.Int, typeParamABList), "map size. Time cost proportional to the number of entries."),
			addOverloadDoc(decls.NewOverload(overloads.SizeString,
				[]*exprpb.Type{decls.String}, decls.Int), "string length"),
			addOverloadDoc(decls.NewOverload(overloads.SizeBytes,
				[]*exprpb.Type{decls.Bytes}, decls.Int), "bytes length"),
			addOverloadDoc(decls.NewParameterizedOverload(overloads.SizeList,
				[]*exprpb.Type{listOfA}, decls.Int, typeParamAList), "list size. Time cost proportional to the length of the list."),
			addOverloadDoc(decls.NewParameterizedOverload(overloads.SizeMap,
				[]*exprpb.Type{mapOfAB}, decls.Int, typeParamABList), "map size. Time cost proportional to the number of entries.")),

		decls.NewFunction(operators.In,
			decls.NewParameterizedOverload(overloads.InList,
				[]*exprpb.Type{paramA, listOfA}, decls.Bool,
				typeParamAList),
			decls.NewParameterizedOverload(overloads.InMap,
				[]*exprpb.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList)),

		// Deprecated 'in()' function.

		decls.NewFunction(overloads.DeprecatedIn,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.InList,
				[]*exprpb.Type{paramA, listOfA}, decls.Bool,
				typeParamAList), "list membership. Time cost proportional to the product of the size of both arguments."),
			addOverloadDoc(decls.NewParameterizedOverload(overloads.InMap,
				[]*exprpb.Type{paramA, mapOfAB}, decls.Bool,
				typeParamABList), "map key membership. Time cost proportional to the product of the size of both arguments.")),

		// Conversions to type.

		decls.NewFunction(overloads.TypeConvertType,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.TypeConvertType,
				[]*exprpb.Type{paramA}, decls.NewTypeType(paramA), typeParamAList), "returns type of value")),

		// Conversions to int.

		decls.NewFunction(overloads.TypeConvertInt,
			addOverloadDoc(decls.NewOverload(overloads.IntToInt, []*exprpb.Type{decls.Int}, decls.Int), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.UintToInt, []*exprpb.Type{decls.Uint}, decls.Int), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.DoubleToInt, []*exprpb.Type{decls.Double}, decls.Int), "Type conversion. Rounds toward zero, then errors if result is out of range."),
			addOverloadDoc(decls.NewOverload(overloads.StringToInt, []*exprpb.Type{decls.String}, decls.Int), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.TimestampToInt, []*exprpb.Type{decls.Timestamp}, decls.Int), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.DurationToInt, []*exprpb.Type{decls.Duration}, decls.Int), "Convert timestamp to int64 in seconds since Unix epoch.")),

		// Conversions to uint.

		decls.NewFunction(overloads.TypeConvertUint,
			addOverloadDoc(decls.NewOverload(overloads.UintToUint, []*exprpb.Type{decls.Uint}, decls.Uint), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.IntToUint, []*exprpb.Type{decls.Int}, decls.Uint), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.DoubleToUint, []*exprpb.Type{decls.Double}, decls.Uint), "Type conversion. Rounds toward zero, then errors if result is out of range."),
			addOverloadDoc(decls.NewOverload(overloads.StringToUint, []*exprpb.Type{decls.String}, decls.Uint), "type conversion")),

		// Conversions to double.

		decls.NewFunction(overloads.TypeConvertDouble,
			addOverloadDoc(decls.NewOverload(overloads.DoubleToDouble, []*exprpb.Type{decls.Double}, decls.Double), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.IntToDouble, []*exprpb.Type{decls.Int}, decls.Double), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.UintToDouble, []*exprpb.Type{decls.Uint}, decls.Double), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.StringToDouble, []*exprpb.Type{decls.String}, decls.Double), "type conversion")),

		// Conversions to bool.

		decls.NewFunction(overloads.TypeConvertBool,
			addOverloadDoc(decls.NewOverload(overloads.BoolToBool, []*exprpb.Type{decls.Bool}, decls.Bool), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.StringToBool, []*exprpb.Type{decls.String}, decls.Bool), "type conversion")),

		// Conversions to string.

		decls.NewFunction(overloads.TypeConvertString,
			addOverloadDoc(decls.NewOverload(overloads.StringToString, []*exprpb.Type{decls.String}, decls.String), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.BoolToString, []*exprpb.Type{decls.Bool}, decls.String), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.IntToString, []*exprpb.Type{decls.Int}, decls.String), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.UintToString, []*exprpb.Type{decls.Uint}, decls.String), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.DoubleToString, []*exprpb.Type{decls.Double}, decls.String), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.BytesToString, []*exprpb.Type{decls.Bytes}, decls.String), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.TimestampToString, []*exprpb.Type{decls.Timestamp}, decls.String), "type conversion, using the same format as timestamp string parsing"),
			addOverloadDoc(decls.NewOverload(overloads.DurationToString, []*exprpb.Type{decls.Duration}, decls.String), "type conversion, using the same format as duration string parsing")),

		// Conversions to bytes.

		decls.NewFunction(overloads.TypeConvertBytes,
			addOverloadDoc(decls.NewOverload(overloads.BytesToBytes, []*exprpb.Type{decls.Bytes}, decls.Bytes), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.StringToBytes, []*exprpb.Type{decls.String}, decls.Bytes), "type conversion")),

		// Conversions to timestamps.

		decls.NewFunction(overloads.TypeConvertTimestamp,
			addOverloadDoc(decls.NewOverload(overloads.TimestampToTimestamp,
				[]*exprpb.Type{decls.Timestamp}, decls.Timestamp), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.StringToTimestamp,
				[]*exprpb.Type{decls.String}, decls.Timestamp), `Type conversion of strings to timestamps according to RFC3339. Example: "1972-01-01T10:00:20.021-05:00"`),
			addOverloadDoc(decls.NewOverload(overloads.IntToTimestamp,
				[]*exprpb.Type{decls.Int}, decls.Timestamp), "type conversion")),

		// Conversions to durations.

		decls.NewFunction(overloads.TypeConvertDuration,
			addOverloadDoc(decls.NewOverload(overloads.DurationToDuration,
				[]*exprpb.Type{decls.Duration}, decls.Duration), "type conversion"),
			addOverloadDoc(decls.NewOverload(overloads.StringToDuration,
				[]*exprpb.Type{decls.String}, decls.Duration), `Type conversion. Duration strings should support the following suffixes: "h" (hour), "m" (minute), "s" (second), "ms" (millisecond), "us" (microsecond), and "ns" (nanosecond). Duration strings may be zero, negative, fractional, and/or compound. Examples: "0", "-1.5h", "1m6s"`),
			addOverloadDoc(decls.NewOverload(overloads.IntToDuration,
				[]*exprpb.Type{decls.Int}, decls.Duration), "type conversion")),

		// Conversions to Dyn.

		decls.NewFunction(overloads.TypeConvertDyn,
			addOverloadDoc(decls.NewParameterizedOverload(overloads.ToDyn,
				[]*exprpb.Type{paramA}, decls.Dyn,
				typeParamAList), "type conversion")),

		// String functions.

		decls.NewFunction(overloads.Contains,
			addOverloadDoc(decls.NewInstanceOverload(overloads.ContainsString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "Tests whether the string operand contains the substring. Time cost proportional to the product of sizes of the arguments.")),
		decls.NewFunction(overloads.EndsWith,
			addOverloadDoc(decls.NewInstanceOverload(overloads.EndsWithString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "Tests whether the string operand ends with the suffix argument. Time cost proportional to the product of the sizes of the arguments.")),
		decls.NewFunction(overloads.Matches,
			addOverloadDoc(decls.NewOverload(overloads.Matches,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "Matches first argument against regular expression in second argument. Time cost proportional to the product of the sizes of the arguments."),
			addOverloadDoc(decls.NewInstanceOverload(overloads.MatchesString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "Matches the self argument against regular expression in first argument. Time cost proportional to the product of the sizes of the arguments.")),
		decls.NewFunction(overloads.StartsWith,
			addOverloadDoc(decls.NewInstanceOverload(overloads.StartsWithString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "Tests whether the string operand starts with the prefix argument. Time cost proportional to the product of the sizes of the arguments.")),

		// Date/time functions.

		decls.NewFunction(overloads.TimeGetFullYear,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToYear,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get year from the date in UTC"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToYearWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get year from the date with timezone")),

		decls.NewFunction(overloads.TimeGetMonth,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToMonth,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get month from the date in UTC, 0-11"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToMonthWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get month from the date with timezone, 0-11")),

		decls.NewFunction(overloads.TimeGetDayOfYear,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfYear,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get day of year from the date in UTC, zero-based indexing"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfYearWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get day of year from the date with timezone, zero-based indexing")),

		decls.NewFunction(overloads.TimeGetDayOfMonth,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBased,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get day of month from the date in UTC, zero-based indexing"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfMonthZeroBasedWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get day of month from the date with timezone, zero-based indexing")),

		decls.NewFunction(overloads.TimeGetDate,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBased,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get day of month from the date in UTC, one-based indexing"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfMonthOneBasedWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get day of month from the date with timezone, one-based indexing")),

		decls.NewFunction(overloads.TimeGetDayOfWeek,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfWeek,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get day of week from the date in UTC, zero-based, zero for Sunday"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToDayOfWeekWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get day of week from the date with timezone, zero-based, zero for Sunday")),

		decls.NewFunction(overloads.TimeGetHours,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToHours,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get hours from the date in UTC, 0-23"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToHoursWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get hours from the date with timezone, 0-23"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.DurationToHours,
				[]*exprpb.Type{decls.Duration}, decls.Int), "get hours from duration")),

		decls.NewFunction(overloads.TimeGetMinutes,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToMinutes,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get minutes from the date in UTC, 0-59"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToMinutesWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get minutes from the date with timezone, 0-59"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.DurationToMinutes,
				[]*exprpb.Type{decls.Duration}, decls.Int), "get minutes from duration")),

		decls.NewFunction(overloads.TimeGetSeconds,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToSeconds,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get seconds from the date in UTC, 0-59"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToSecondsWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get seconds from the date with timezone, 0-59"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.DurationToSeconds,
				[]*exprpb.Type{decls.Duration}, decls.Int), "get seconds from duration")),

		decls.NewFunction(overloads.TimeGetMilliseconds,
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToMilliseconds,
				[]*exprpb.Type{decls.Timestamp}, decls.Int), "get milliseconds from the date in UTC, 0-999"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.TimestampToMillisecondsWithTz,
				[]*exprpb.Type{decls.Timestamp, decls.String}, decls.Int), "get milliseconds from the date with timezone, 0-999"),
			addOverloadDoc(decls.NewInstanceOverload(overloads.DurationToMilliseconds,
				[]*exprpb.Type{decls.Duration}, decls.Int), "milliseconds from duration, 0-999")),

		// Relations.
		decls.NewFunction(operators.Less,
			addOverloadDoc(decls.NewOverload(overloads.LessBool,
				[]*exprpb.Type{decls.Bool, decls.Bool}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessInt64Double,
				[]*exprpb.Type{decls.Int, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessInt64Uint64,
				[]*exprpb.Type{decls.Int, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessUint64Double,
				[]*exprpb.Type{decls.Uint, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessUint64Int64,
				[]*exprpb.Type{decls.Uint, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessDoubleInt64,
				[]*exprpb.Type{decls.Double, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessDoubleUint64,
				[]*exprpb.Type{decls.Double, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessBytes,
				[]*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessTimestamp,
				[]*exprpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessDuration,
				[]*exprpb.Type{decls.Duration, decls.Duration}, decls.Bool), "ordering")),

		decls.NewFunction(operators.LessEquals,
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsBool,
				[]*exprpb.Type{decls.Bool, decls.Bool}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsInt64Double,
				[]*exprpb.Type{decls.Int, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsInt64Uint64,
				[]*exprpb.Type{decls.Int, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsUint64Double,
				[]*exprpb.Type{decls.Uint, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsUint64Int64,
				[]*exprpb.Type{decls.Uint, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsDoubleInt64,
				[]*exprpb.Type{decls.Double, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsDoubleUint64,
				[]*exprpb.Type{decls.Double, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsBytes,
				[]*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsTimestamp,
				[]*exprpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.LessEqualsDuration,
				[]*exprpb.Type{decls.Duration, decls.Duration}, decls.Bool), "ordering")),

		decls.NewFunction(operators.Greater,
			addOverloadDoc(decls.NewOverload(overloads.GreaterBool,
				[]*exprpb.Type{decls.Bool, decls.Bool}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterInt64Double,
				[]*exprpb.Type{decls.Int, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterInt64Uint64,
				[]*exprpb.Type{decls.Int, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterUint64Double,
				[]*exprpb.Type{decls.Uint, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterUint64Int64,
				[]*exprpb.Type{decls.Uint, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterDoubleInt64,
				[]*exprpb.Type{decls.Double, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterDoubleUint64,
				[]*exprpb.Type{decls.Double, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterBytes,
				[]*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterTimestamp,
				[]*exprpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterDuration,
				[]*exprpb.Type{decls.Duration, decls.Duration}, decls.Bool), "ordering")),

		decls.NewFunction(operators.GreaterEquals,
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsBool,
				[]*exprpb.Type{decls.Bool, decls.Bool}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsInt64,
				[]*exprpb.Type{decls.Int, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsInt64Double,
				[]*exprpb.Type{decls.Int, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsInt64Uint64,
				[]*exprpb.Type{decls.Int, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsUint64,
				[]*exprpb.Type{decls.Uint, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsUint64Double,
				[]*exprpb.Type{decls.Uint, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsUint64Int64,
				[]*exprpb.Type{decls.Uint, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsDouble,
				[]*exprpb.Type{decls.Double, decls.Double}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsDoubleInt64,
				[]*exprpb.Type{decls.Double, decls.Int}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsDoubleUint64,
				[]*exprpb.Type{decls.Double, decls.Uint}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsString,
				[]*exprpb.Type{decls.String, decls.String}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsBytes,
				[]*exprpb.Type{decls.Bytes, decls.Bytes}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsTimestamp,
				[]*exprpb.Type{decls.Timestamp, decls.Timestamp}, decls.Bool), "ordering"),
			addOverloadDoc(decls.NewOverload(overloads.GreaterEqualsDuration,
				[]*exprpb.Type{decls.Duration, decls.Duration}, decls.Bool), "ordering")),
	}...)
}

// StandardDeclarations returns the Decls for all functions and constants in the evaluator.
func StandardDeclarations() []*exprpb.Decl {
	return standardDeclarations
}
