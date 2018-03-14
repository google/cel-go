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

package functions

const (
	// Boolean logic overloads
	LogicalNot             = "logical_not"
	EqualsBytes            = "equals_bytes"
	NotEqualsBytes         = "not_equals_bytes"
	LessBool               = "less_bool"
	LessInt64              = "less_int64"
	LessUint64             = "less_uint64"
	LessDouble             = "less_double"
	LessString             = "less_string"
	LessTimestamp          = "less_timestamp"
	LessDuration           = "less_duration"
	LessEqualsBool         = "less_equals_bool"
	LessEqualsInt64        = "less_equals_int64"
	LessEqualsUint64       = "less_equals_uint64"
	LessEqualsDouble       = "less_equals_double"
	LessEqualsString       = "less_equals_string"
	LessEqualsTimestamp    = "less_equals_timestamp"
	LessEqualsDuration     = "less_equals_duration"
	GreaterBool            = "greater_bool"
	GreaterInt64           = "greater_int64"
	GreaterUint64          = "greater_uint64"
	GreaterDouble          = "greater_double"
	GreaterString          = "greater_string"
	GreaterTimestamp       = "greater_timestamp"
	GreaterDuration        = "greater_duration"
	GreaterEqualsBool      = "greater_equals_bool"
	GreaterEqualsInt64     = "greater_equals_int64"
	GreaterEqualsUint64    = "greater_equals_uint64"
	GreaterEqualsDouble    = "greater_equals_double"
	GreaterEqualsString    = "greater_equals_string"
	GreaterEqualsTimestamp = "greater_equals_timestamp"
	GreaterEqualsDuration  = "greater_equals_duration"

	// Math overloads
	AddInt64                   = "add_int64"
	AddUint64                  = "add_uint64"
	AddDouble                  = "add_double"
	AddString                  = "add_string"
	AddBytes                   = "add_bytes"
	AddList                    = "add_list"
	AddTimestampDuration       = "add_timestamp_duration"
	AddDurationTimestamp       = "add_duration_timestamp"
	AddDurationDuration        = "add_duration_duration"
	SubtractInt64              = "subtract_int64"
	SubtractUint64             = "subtract_uint64"
	SubtractDouble             = "subtract_double"
	SubtractTimestampTimestamp = "subtract_timestamp_timestamp"
	SubtractTimestampDuration  = "subtract_timestamp_duration"
	SubtractDurationDuration   = "subtract_duration_duration"
	MultiplyInt64              = "multiply_int64"
	MultiplyUint64             = "multiply_uint64"
	MultiplyDouble             = "multiply_double"
	DivideInt64                = "divide_int64"
	DivideUint64               = "divide_uint64"
	DivideDouble               = "divide_double"
	ModuloInt64                = "modulo_int64"
	ModuloUint64               = "modulo_uint64"
	NegateInt64                = "negate_int64"
	NegateDouble               = "negate_double"

	// In operators
	InList = "in_list"
	InMap  = "in_map"

	// Size overloads
	Size       = "size"
	SizeString = "size_string"
	SizeBytes  = "size_bytes"
	SizeList   = "size_list"
	SizeMap    = "size_map"

	// Matches function
	MatchString = "matches"

	// Time-based functions
	TimeGetFullYear     = "getFullYear"
	TimeGetMonth        = "getMonth"
	TimeGetDayOfYear    = "getDayOfYear"
	TimeGetDate         = "getDate"
	TimeGetDayOfMonth   = "getDayOfMonth"
	TimeGetDayOfWeek    = "getDayOfWeek"
	TimeGetHours        = "getHours"
	TimeGetMinutes      = "getMinutes"
	TimeGetSeconds      = "getSeconds"
	TimeGetMilliseconds = "getMilliseconds"

	// Timestamp overloads for time functions without timezones.
	TimestampToYear                = "timestamp_to_year"
	TimestampToMonth               = "timestamp_to_month"
	TimestampToDayOfYear           = "timestamp_to_day_of_year"
	TimestampToDayOfMonthZeroBased = "timestamp_to_day_of_month_zero_based"
	TimestampToDayOfMonthOneBased  = "timestamp_to_day_of_month_one_based"
	TimestampToDayOfWeek           = "timestamp_to_day_of_week"
	TimestampToHours               = "timestamp_to_hours"
	TimestampToMinutes             = "timestamp_to_minutes"
	TimestampToSeconds             = "timestamp_to_seconds"
	TimestampToMilliseconds        = "timestamp_to_milliseconds"

	// Timestamp overloads for time functions with timezones.
	TimestampToYearWithTz                = "timestamp_to_year_with_tz"
	TimestampToMonthWithTz               = "timestamp_to_month_with_tz"
	TimestampToDayOfYearWithTz           = "timestamp_to_day_of_year_with_tz"
	TimestampToDayOfMonthZeroBasedWithTz = "timestamp_to_day_of_month_zero_based_with_tz"
	TimestampToDayOfMonthOneBasedWithTz  = "timestamp_to_day_of_month_one_based_with_tz"
	TimestampToDayOfWeekWithTz           = "timestamp_to_day_of_week_with_tz"
	TimestampToHoursWithTz               = "timestamp_to_hours_with_tz"
	TimestampToMinutesWithTz             = "timestamp_to_minutes_with_tz"
	TimestampToSecondsWithTz             = "timestamp_to_seconds_tz"
	TimestampToMillisecondsWithTz        = "timestamp_to_milliseconds_with_tz"

	// Duration overloads for time functions.
	DurationToHours        = "duration_to_hours"
	DurationToMinutes      = "duration_to_minutes"
	DurationToSeconds      = "duration_to_seconds"
	DurationToMilliseconds = "duration_to_milliseconds"

	// Type conversion methods and overloads
	TypeConvertInt       = "int"
	TypeConvertUint      = "uint"
	TypeConvertDouble    = "double"
	TypeConvertBool      = "bool"
	TypeConvertString    = "string"
	TypeConvertBytes     = "bytes"
	TypeConvertTimestamp = "timestamp"
	TypeConvertDuration  = "duration"

	// Int conversion functions.
	IntFromInt       = "int_from_int"
	IntFromUint      = "int_from_uint"
	IntFromDouble    = "int_from_double"
	IntFromString    = "int_from_string"
	IntFromTimestamp = "int_from_timestamp"
	IntFromDuration  = "int_from_duration"

	// Uint conversion functions.
	UintFromUint   = "uint_from_uint"
	UintFromInt    = "uint_from_int"
	UintFromDouble = "uint_from_double"
	UintFromString = "uint_from_string"

	// Double conversion functions.
	DoubleFromDouble = "double_from_double"
	DoubleFromInt    = "double_from_int"
	DoubleFromUint   = "double_from_uint"
	DoubleFromString = "double_from_string"

	// Bool conversion functions.
	BoolFromBool   = "bool_from_bool"
	BoolFromString = "bool_from_string"

	// Bytes conversion functions.
	BytesFromBytes  = "bytes_from_bytes"
	BytesFromString = "bytes_from_string"

	// String conversion functions.
	StringFromString    = "string_from_string"
	StringFromBool      = "string_from_bool"
	StringFromInt       = "string_from_int"
	StringFromUint      = "string_from_uint"
	StringFromDouble    = "string_from_double"
	StringFromBytes     = "string_from_bytes"
	StringFromTimestamp = "string_from_timestamp"
	StringFromDuration  = "string_from_duration"

	// Timestamp conversion functions
	TimestampFromTimestamp = "timestamp_from_timestamp"
	TimestampFromString    = "timestamp_from_string"
	TimestampFromInt       = "timestamp_from_int"

	// Convert duration from string
	DurationFromDuration = "duration_from_duration"
	DurationFromString   = "duration_from_string"
	DurationFromInt      = "duration_from_int"

	// Type inspection methods.
	TypeOfValue = "type"

	// Comprehensions helper methods.
	Iterator = "@iterator"
	HasNext  = "@hasNext"
	Next     = "@next"
)
