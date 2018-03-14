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

import (
	"bytes"
	"github.com/google/cel-go/interpreter/types"
	"github.com/google/cel-go/interpreter/types/adapters"
	"github.com/google/cel-go/interpreter/types/objects"
	"github.com/google/cel-go/operators"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	dpb "github.com/golang/protobuf/ptypes/duration"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"regexp"
	"strconv"
	"time"
)

type OverloadImpl func(values ...interface{}) (interface{}, error)

type Overload struct {
	Function  string
	Name      string
	Signature interface{}
	Impl      OverloadImpl
}

func StandardBuiltins() []*Overload {
	return []*Overload{
		// Logical not
		{
			operators.LogicalNot, operators.LogicalNot,
			func(value bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[0].(bool), nil
			}},
		{
			operators.LogicalAnd, operators.LogicalAnd,
			func(lhs, rhs bool) {},
			func(values ...interface{}) (interface{}, error) {
				lhs, lhsIsBool := values[0].(bool)
				rhs, rhsIsBool := values[1].(bool)
				if lhsIsBool && rhsIsBool {
					return lhs && rhs, nil
				} else if lhsIsBool && !lhs ||
					rhsIsBool && !rhs {
					return false, nil
				} else if !rhsIsBool {
					return rhs, nil
				} else {
					return lhs, nil
				}
			}},
		{
			operators.LogicalOr, operators.LogicalOr,
			func(value1, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				lhs, lhsIsBool := values[0].(bool)
				rhs, rhsIsBool := values[1].(bool)
				if lhsIsBool && rhsIsBool {
					return lhs || rhs, nil
				} else if lhsIsBool && lhs ||
					rhsIsBool && rhs {
					return true, nil
				} else if !rhsIsBool {
					return rhs, nil
				} else {
					return lhs, nil
				}
			}},
		{
			operators.Conditional, operators.Conditional,
			func(cond bool, trueVal, falseVal interface{}) {},
			func(values ...interface{}) (interface{}, error) {
				cond, isBool := values[0].(bool)
				trueVal := values[1]
				falseVal := values[2]
				if isBool {
					if cond {
						return trueVal, nil
					} else {
						return falseVal, nil
					}
				} else {
					return values[0], nil
				}
			}},

		// Equality operators
		{operators.Equals, EqualsBytes,
			func(value1, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				value1 := values[0].([]byte)
				value2 := values[1].([]byte)
				return bytes.Equal(value1, value2), nil
			}},

		{operators.Equals, operators.Equals,
			func(value1, value2 interface{}) {},
			func(values ...interface{}) (interface{}, error) {
				if value1, ok := values[0].(objects.Equaler); ok {
					return value1.Equal(values[1]), nil
				} else if value2, ok := values[1].(objects.Equaler); ok {
					return value2.Equal(values[0]), nil
				} else {
					return values[0] == values[1], nil
				}
			}},
		{operators.NotEquals, NotEqualsBytes,
			func(value1, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				value1 := values[0].([]byte)
				value2 := values[1].([]byte)
				return !bytes.Equal(value1, value2), nil
			}},
		{operators.NotEquals, operators.NotEquals,
			func(value1, value2 interface{}) {},
			func(values ...interface{}) (interface{}, error) {
				if value1, ok := values[0].(objects.Equaler); ok {
					return !value1.Equal(values[1]), nil
				} else if value2, ok := values[1].(objects.Equaler); ok {
					return !value2.Equal(values[0]), nil
				} else {
					return values[0] != values[1], nil
				}
			}},

		// Less than operator
		{operators.Less, LessBool,
			func(value1, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[0].(bool) && values[1].(bool), nil
			}},
		{operators.Less, LessInt64,
			func(value1, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) < values[1].(int64), nil
			}},
		{operators.Less, LessUint64,
			func(value1, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) < values[1].(uint64), nil
			}},
		{operators.Less, LessDouble,
			func(value1, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) < values[1].(float64), nil
			}},
		{operators.Less, LessString,
			func(value1, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) < values[1].(string), nil
			}},
		{operators.Less, LessTimestamp,
			func(value1, value2 *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				if ts1, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if ts2, err := ptypes.Timestamp(values[1].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else {
					return ts1.Before(ts2), nil
				}
			}},
		{operators.Less, LessDuration,
			func(value1, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if d1, err := ptypes.Duration(values[0].(*dpb.Duration)); err != nil {
					return nil, err
				} else if d2, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return d1 < d2, nil
				}
			}},

		// Less than or equal operator
		{operators.LessEquals, LessEqualsBool,
			func(value1, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[0].(bool), nil
			}},
		{operators.LessEquals, LessEqualsInt64,
			func(value1, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) <= values[1].(int64), nil
			}},
		{operators.LessEquals, LessEqualsUint64,
			func(value1, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) <= values[1].(uint64), nil
			}},
		{operators.LessEquals, LessEqualsDouble,
			func(value1, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) <= values[1].(float64), nil
			}},
		{operators.LessEquals, LessEqualsString,
			func(value1, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) <= values[1].(string), nil
			}},
		{operators.LessEquals, LessEqualsTimestamp,
			func(value1, value2 *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				if ts1, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if ts2, err := ptypes.Timestamp(values[1].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else {
					return ts1.Before(ts2) || ts1.Equal(ts2), nil
				}
			}},
		{operators.LessEquals, LessEqualsDuration,
			func(value1, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if d1, err := ptypes.Duration(values[0].(*dpb.Duration)); err != nil {
					return nil, err
				} else if d2, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return d1 <= d2, nil
				}
			}},

		// Greater than operator
		{operators.Greater, GreaterBool,
			func(value1 bool, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[1].(bool) && values[0].(bool), nil
			}},
		{operators.Greater, GreaterInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) > values[1].(int64), nil
			}},
		{operators.Greater, GreaterUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) > values[1].(uint64), nil
			}},
		{operators.Greater, GreaterDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) > values[1].(float64), nil
			}},
		{operators.Greater, GreaterString,
			func(value1 string, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) > values[1].(string), nil
			}},
		{operators.Greater, GreaterTimestamp,
			func(value1 *tspb.Timestamp, value2 *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				if ts1, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if ts2, err := ptypes.Timestamp(values[1].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else {
					return ts1.After(ts2), nil
				}
			}},
		{operators.Greater, GreaterDuration,
			func(value1 *dpb.Duration, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if d1, err := ptypes.Duration(values[0].(*dpb.Duration)); err != nil {
					return nil, err
				} else if d2, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return d1 > d2, nil
				}
			}},

		// Greater than equal operators
		{operators.GreaterEquals, GreaterEqualsBool,
			func(value1 bool, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[1].(bool), nil
			}},
		{operators.GreaterEquals, GreaterEqualsInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) >= values[1].(int64), nil
			}},
		{operators.GreaterEquals, GreaterEqualsUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) >= values[1].(uint64), nil
			}},
		{operators.GreaterEquals, GreaterEqualsDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) >= values[1].(float64), nil
			}},
		{operators.GreaterEquals, GreaterEqualsString,
			func(value1 string, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) >= values[1].(string), nil
			}},
		{operators.GreaterEquals, GreaterEqualsTimestamp,
			func(value1 *tspb.Timestamp, value2 *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				if ts1, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if ts2, err := ptypes.Timestamp(values[1].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else {
					return ts1.After(ts2) || ts1.Equal(ts2), nil
				}
			}},
		{operators.GreaterEquals, GreaterEqualsDuration,
			func(value1 *dpb.Duration, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				value1 := values[0].(*dpb.Duration)
				value2 := values[1].(*dpb.Duration)
				if d1, err := ptypes.Duration(value1); err != nil {
					return nil, err
				} else if d2, err := ptypes.Duration(value2); err != nil {
					return nil, err
				} else {
					return d1 >= d2, nil
				}
			}},

		// TODO: Verify overflow, NaN, underflow cases for numeric values.

		// Add operator
		{operators.Add, AddInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) + values[1].(int64), nil
			}},
		{operators.Add, AddUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) + values[1].(uint64), nil
			}},
		{operators.Add, AddDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) + values[1].(float64), nil
			}},
		{operators.Add, AddString,
			func(value1 string, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) + values[1].(string), nil
			}},
		{operators.Add, AddBytes,
			func(value1 []byte, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return append(values[0].([]byte), values[1].([]byte)...), nil
			}},
		{operators.Add, AddList,
			func(value1 adapters.ListAdapter, value2 adapters.ListAdapter) {},
			func(values ...interface{}) (interface{}, error) {
				list1 := values[0].(adapters.ListAdapter)
				list2 := values[1].(adapters.ListAdapter)
				return list1.Concat(list2), nil
			}},
		{operators.Add, AddTimestampDuration,
			func(value1 *tspb.Timestamp, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if ts, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if dur, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return ts.Add(dur), nil
				}
			}},
		{operators.Add, AddDurationTimestamp,
			func(value1 *dpb.Duration, value2 *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				if dur, err := ptypes.Duration(values[0].(*dpb.Duration)); err != nil {
					return nil, err
				} else if ts, err := ptypes.Timestamp(values[1].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else {
					return ts.Add(dur), nil
				}
			}},
		{operators.Add, AddDurationDuration,
			func(value1 *dpb.Duration, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if dur1, err := ptypes.Duration(values[0].(*dpb.Duration)); err != nil {
					return nil, err
				} else if dur2, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return dur1 + dur2, nil
				}
			}},

		// Subtract operators
		{operators.Subtract, SubtractInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) - values[1].(int64), nil
			}},
		{operators.Subtract, SubtractUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) - values[1].(uint64), nil
			}},
		{operators.Subtract, SubtractDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) - values[1].(float64), nil
			}},
		{operators.Subtract, SubtractTimestampTimestamp,
			func(value1 *tspb.Timestamp, value2 *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				if ts1, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if ts2, err := ptypes.Timestamp(values[1].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else {
					return ts1.Sub(ts2), nil
				}
			}},
		{operators.Subtract, SubtractTimestampDuration,
			func(value1 *tspb.Timestamp, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if ts, err := ptypes.Timestamp(values[0].(*tspb.Timestamp)); err != nil {
					return nil, err
				} else if dur, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return ts.Add(-dur), nil
				}
			}},
		{operators.Subtract, SubtractDurationDuration,
			func(value1 *dpb.Duration, value2 *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				if dur1, err := ptypes.Duration(values[0].(*dpb.Duration)); err != nil {
					return nil, err
				} else if dur2, err := ptypes.Duration(values[1].(*dpb.Duration)); err != nil {
					return nil, err
				} else {
					return dur1 - dur2, nil
				}
			}},

		// Multiply operator
		{operators.Multiply, MultiplyInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) * values[1].(int64), nil
			}},
		{operators.Multiply, MultiplyUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) * values[1].(uint64), nil
			}},
		{operators.Multiply, MultiplyDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) * values[1].(float64), nil
			}},

		// Divide operator
		{operators.Divide, DivideInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) / values[1].(int64), nil
			}},
		{operators.Divide, DivideUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) / values[1].(uint64), nil
			}},
		{operators.Divide, DivideDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) / values[1].(float64), nil
			}},

		// Modulo operator
		{operators.Modulo, ModuloInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) % values[1].(int64), nil
			}},
		{operators.Modulo, ModuloUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) % values[1].(uint64), nil
			}},

		// Negate operator
		{operators.Negate, NegateInt64,
			func(value int64) {},
			func(values ...interface{}) (interface{}, error) {
				return -values[0].(int64), nil
			}},
		{operators.Negate, NegateDouble,
			func(value uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return -values[0].(float64), nil
			}},

		// Index operator
		{operators.Index, operators.Index,
			func(value objects.Indexer, index interface{}) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].(objects.Indexer)
				index := values[1]
				return value.Get(index)
			}},

		// Size function
		{Size, SizeString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].(string)
				return len(value), nil
			}},
		{Size, SizeBytes,
			func(value []byte) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].([]byte)
				return len(value), nil
			}},
		{Size, SizeList,
			func(value adapters.ListAdapter) {},
			func(values ...interface{}) (interface{}, error) {
				listValue := values[0].(adapters.ListAdapter)
				return listValue.Len(), nil
			}},
		{Size, SizeMap,
			func(value adapters.MapAdapter) {},
			func(values ...interface{}) (interface{}, error) {
				mapValue := values[0].(adapters.MapAdapter)
				return mapValue.Len(), nil
			}},

		// In operator
		{operators.In, InList,
			func(value interface{}, listValue adapters.ListAdapter) {},
			func(values ...interface{}) (interface{}, error) {
				element := values[0]
				listValue := values[1].(adapters.ListAdapter)
				return listValue.Contains(element), nil
			}},
		{operators.In, InMap,
			func(index interface{}, value objects.Indexer) {},
			func(values ...interface{}) (interface{}, error) {
				index := values[0]
				value := values[1].(objects.Indexer)
				// FIXME: This needs to be a Has() method
				_, err := value.Get(index)
				return err == nil, nil
			}},

		// Matches function
		{MatchString, MatchString,
			func(text string, pattern string) {},
			func(values ...interface{}) (interface{}, error) {
				text := values[0].(string)
				pattern := values[1].(string)
				return regexp.MatchString(pattern, text)
			}},

		// Timestamp member functions.
		{TimeGetFullYear, TimestampToYear,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetFullYear)
			}},
		{TimeGetMonth, TimestampToMonth,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetMonth)
			}},
		{TimeGetDayOfYear, TimestampToDayOfYear,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfYear)
			}},
		{TimeGetDate, TimestampToDayOfMonthZeroBased,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfMonthZeroBased)
			}},
		{TimeGetDayOfMonth, TimestampToDayOfMonthOneBased,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfMonthOneBased)
			}},
		{TimeGetDayOfWeek, TimestampToDayOfWeek,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfWeek)
			}},
		{TimeGetHours, TimestampToHours,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetHours)
			}},
		{TimeGetMinutes, TimestampToMinutes,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetMinutes)
			}},
		{TimeGetSeconds, TimestampToSeconds,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetSeconds)
			}},
		{TimeGetMilliseconds, TimestampToMilliseconds,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetMilliseconds)
			}},

		// Timestamp with time zone.
		{TimeGetFullYear, TimestampToYearWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetFullYear))
			}},
		{TimeGetMonth, TimestampToMonthWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetMonth))
			}},
		{TimeGetDayOfWeek, TimestampToDayOfYearWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfYear))
			}},
		{TimeGetDate, TimestampToDayOfMonthZeroBasedWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfMonthZeroBased))
			}},
		{TimeGetDayOfMonth, TimestampToDayOfMonthOneBasedWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfMonthOneBased))
			}},
		{TimeGetDayOfWeek, TimestampToDayOfWeekWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfWeek))
			}},
		{TimeGetHours, TimestampToHoursWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetHours))
			}},
		{TimeGetMinutes, TimestampToMinutesWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetMinutes))
			}},
		{TimeGetSeconds, TimestampToSecondsWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetSeconds))
			}},
		{TimeGetMilliseconds, TimestampToMillisecondsWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetMilliseconds))
			}},

		// Duration member functions.
		{TimeGetHours, DurationToHours,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Hours(), nil
				} else {
					return nil, err
				}
			}},
		{TimeGetMinutes, DurationToMinutes,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Minutes(), nil
				} else {
					return nil, err
				}
			}},
		{TimeGetSeconds, DurationToSeconds,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Seconds(), nil
				} else {
					return nil, err
				}
			}},
		{TimeGetMilliseconds, DurationToMilliseconds,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Nanoseconds() / 1000, nil
				} else {
					return nil, err
				}
			}},

		// Type conversion functions
		// TODO: verify type conversion safety of numeric values.

		// Int conversions.
		{TypeConvertInt, IntFromInt,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertInt, IntFromUint,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(uint64)
				return int64(num), nil
			}},
		{TypeConvertInt, IntFromDouble,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(float64)
				return int64(num), nil
			}},
		{TypeConvertInt, IntFromString,
			func(num string) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(string)
				return strconv.ParseInt(num, 10, 64)
			}},
		{TypeConvertInt, IntFromTimestamp,
			func(ts *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				// Return the Unix time in seconds since 1970
				ts := values[0].(*tspb.Timestamp)
				if t, err := ptypes.Timestamp(ts); err == nil {
					return t.Unix(), nil
				} else {
					return nil, err
				}
			}},
		{TypeConvertInt, IntFromDuration,
			func(ts *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				dur := values[0].(*dpb.Duration)
				if d, err := ptypes.Duration(dur); err == nil {
					return int64(d), nil
				} else {
					return nil, err
				}
			}},

		// Uint conversions.
		{TypeConvertUint, UintFromUint,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertUint, UintFromInt,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(int64)
				if num >= 0 {
					return uint64(num), nil
				} else {
					return nil,
						fmt.Errorf("unsafe uint conversion of negative int")
				}
			}},
		{TypeConvertUint, UintFromDouble,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(float64)
				if num >= 0.0 {
					return float64(num), nil
				} else {
					return nil,
						fmt.Errorf("unsafe uint conversion of negative double")
				}
			}},
		{TypeConvertUint, UintFromString,
			func(num string) {},
			func(values ...interface{}) (interface{}, error) {
				return strconv.ParseUint(values[0].(string), 10, 64)
			}},

		// Double conversions.
		{TypeConvertDouble, DoubleFromDouble,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertDouble, DoubleFromInt,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(int64)
				return float64(num), nil
			}},
		{TypeConvertDouble, DoubleFromUint,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(uint64)
				return float64(num), nil
			}},
		{TypeConvertDouble, DoubleFromString,
			func(num string) {},
			func(values ...interface{}) (interface{}, error) {
				return strconv.ParseFloat(values[0].(string), 64)
			}},

		// Bool conversions.
		{TypeConvertBool, BoolFromBool,
			func(value bool) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertBool, BoolFromString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return strconv.ParseBool(values[0].(string))
			}},

		// Bytes conversions.
		{TypeConvertBytes, BytesFromBytes,
			func(value []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertBytes, BytesFromString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return []byte(values[0].(string)), nil
			}},

		// String conversions.
		{TypeConvertString, StringFromString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertString, StringFromBool,
			func(value bool) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%t", values[0].(bool)), nil
			}},
		{TypeConvertString, StringFromInt,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%d", values[0].(int64)), nil
			}},
		{TypeConvertString, StringFromUint,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%d", values[0].(uint64)), nil
			}},
		{TypeConvertString, StringFromDouble,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%g", values[0].(float64)), nil
			}},
		{TypeConvertString, StringFromBytes,
			func(value []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%s", values[0].([]byte)), nil
			}},
		{TypeConvertString, StringFromTimestamp,
			func(value *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return ptypes.TimestampString(values[0].(*tspb.Timestamp)), nil
			}},
		{TypeConvertString, StringFromDuration,
			func(value *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				return ptypes.TimestampString(values[0].(*tspb.Timestamp)), nil
			}},

		// Timestamp conversions.
		{TypeConvertTimestamp, TimestampFromTimestamp,
			func(value *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertTimestamp, TimestampFromString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				if t, err := time.Parse(time.RFC3339, values[0].(string)); err != nil {
					return nil, err
				} else {
					return ptypes.TimestampProto(t)
				}
			}},
		{TypeConvertTimestamp, TimestampFromInt,
			func(value int64) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].(int64)
				// Unix timestamp in seconds.
				return ptypes.TimestampProto(time.Unix(value, 0))
			}},
		// Duration conversions.
		{TypeConvertDuration, DurationFromDuration,
			func(value *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{TypeConvertDuration, DurationFromString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return time.ParseDuration(values[0].(string))
			}},
		{TypeConvertDuration, DurationFromInt,
			func(value int64) {},
			func(values ...interface{}) (interface{}, error) {
				// Duration in seconds
				return time.Duration(values[0].(int64)) * time.Second, nil
			}},

		// Type operations.
		{TypeOfValue, TypeOfValue,
			func(value types.Type) {},
			func(values ...interface{}) (interface{}, error) {
				if t, found := types.TypeOf(values[0]); found {
					return t, nil
				} else {
					return nil, fmt.Errorf("undefined type")
				}
			}},

		// Comprehensions functions.
		{Iterator, Iterator,
			func(value objects.Iterable) {},
			func(values ...interface{}) (interface{}, error) {
				iterable := values[0].(objects.Iterable)
				return iterable.Iterator(), nil
			}},
		{HasNext, HasNext,
			func(value objects.Iterator) {},
			func(values ...interface{}) (interface{}, error) {
				it := values[0].(objects.Iterator)
				return it.HasNext(), nil
			}},
		{Next, Next,
			func(value objects.Iterator) {},
			func(values ...interface{}) (interface{}, error) {
				it := values[0].(objects.Iterator)
				return it.Next(), nil
			}},
	}

}

// Timestamp-related helpers.
type timestampVisitor func(time.Time) (interface{}, error)

func timestampGetFullYear(t time.Time) (interface{}, error) {
	return int64(t.Year()), nil
}
func timestampGetMonth(t time.Time) (interface{}, error) {
	return int64(t.Month()), nil
}
func timestampGetDayOfYear(t time.Time) (interface{}, error) {
	return int64(t.YearDay()), nil
}
func timestampGetDayOfMonthZeroBased(t time.Time) (interface{}, error) {
	return int64(t.Day() - 1), nil
}
func timestampGetDayOfMonthOneBased(t time.Time) (interface{}, error) {
	return int64(t.Day()), nil
}
func timestampGetDayOfWeek(t time.Time) (interface{}, error) {
	return int64(t.Weekday()), nil
}
func timestampGetHours(t time.Time) (interface{}, error) {
	return int64(t.Hour()), nil
}
func timestampGetMinutes(t time.Time) (interface{}, error) {
	return int64(t.Minute()), nil
}
func timestampGetSeconds(t time.Time) (interface{}, error) {
	return int64(t.Second()), nil
}
func timestampGetMilliseconds(t time.Time) (interface{}, error) {
	return int64(t.Nanosecond() / 1000), nil
}

func timeZone(value interface{}, visitor timestampVisitor) timestampVisitor {
	return func(t time.Time) (interface{}, error) {
		if loc, err := time.LoadLocation(value.(string)); err == nil {
			return visitor(t.In(loc))
		} else {
			return nil, err
		}
	}
}

func timestampOp(value interface{}, visitor timestampVisitor) (interface{}, error) {
	ts := value.(*tspb.Timestamp)
	if t, err := ptypes.Timestamp(ts); err == nil {
		return visitor(t)
	} else {
		return nil, err
	}
}
