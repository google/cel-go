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

// Functions package defines the function names and implementations for
// standard CEL overloads.
package functions

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	dpb "github.com/golang/protobuf/ptypes/duration"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/cel-go/interpreter/types"
	"github.com/google/cel-go/interpreter/types/traits"
	"github.com/google/cel-go/operators"
	"github.com/google/cel-go/overloads"
	"regexp"
	"strconv"
	"time"
)

// Overload defines a named overload of a function, providing a void function
// as a signature to be used for overload resolution, as well as an generic
// overload implementation.
type Overload struct {
	// Function name as written in an expression or defined within
	// operators.go.
	Function string

	// Name of the overload, must be unique.
	Name string

	// Signature of the overload as defined by a void function. The argument
	// count and types will be derived by reflection and will be used to ensure
	// the #Impl is called with right number of arguments of the expected type.
	Signature interface{}

	// Impl of the overload to be called by a dispatcher.
	Impl OverloadImpl
}

// OverloadImpl is a function with accepts zero or more arguments and produces
// an value (as interface{}) or error as a result.
type OverloadImpl func(values ...interface{}) (interface{}, error)

// StandardBuiltins returns the definitions of the built-in CEL overloads.
// TODO: Where possible, organize builtins within types and by traits.
func StandardBuiltins() []*Overload {
	return []*Overload{
		// Logical not (!a)
		{
			operators.LogicalNot, overloads.LogicalNot,
			func(value bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[0].(bool), nil
			}},
		// Logical and (a && b)
		{
			operators.LogicalAnd, overloads.LogicalAnd,
			func(lhs, rhs bool) {},
			func(values ...interface{}) (interface{}, error) {
				lhs, lhsIsBool := values[0].(bool)
				rhs, rhsIsBool := values[1].(bool)
				// both are boolean use natural logic.
				if lhsIsBool && rhsIsBool {
					return lhs && rhs, nil
				}
				// one or the other is boolean and false, return false.
				if lhsIsBool && !lhs ||
					rhsIsBool && !rhs {
					return false, nil
				}
				// if the left-hand side is non-boolean return it as the error.
				if !lhsIsBool {
					return newError("Got '%v', expected argument of type 'bool'", values[0]), nil
				}
				return newError("Got '%v', expected argument of type 'bool'", values[1]), nil
			}},
		// Logical or (a || b)
		{
			operators.LogicalOr, overloads.LogicalOr,
			func(value1, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				lhs, lhsIsBool := values[0].(bool)
				rhs, rhsIsBool := values[1].(bool)
				// both are boolean, use natural logic.
				if lhsIsBool && rhsIsBool {
					return lhs || rhs, nil
				}
				// one or the other is boolean and true, return true
				if lhsIsBool && lhs ||
					rhsIsBool && rhs {
					return true, nil
				}
				// if the left-hand side is non-boolean return it as the error.
				if !lhsIsBool {
					return newError("Got '%v', expected argument of type 'bool'", values[0]), nil
				}
				return newError("Got '%v', expected argument of type 'bool'", values[1]), nil
			}},
		// Conditional operator (a ? b : c)
		{
			operators.Conditional, overloads.Conditional,
			func(cond bool, trueVal, falseVal interface{}) {},
			func(values ...interface{}) (interface{}, error) {
				cond, isBool := values[0].(bool)
				trueVal := values[1]
				falseVal := values[2]
				if isBool {
					if cond {
						return trueVal, nil
					}
					return falseVal, nil
				}
				return newError("Got '%v', expected argument of type 'bool'", values[0]), nil
			}},

		// Equality overloads
		{operators.Equals, overloads.Equals,
			func(value1, value2 interface{}) {}, equals},

		{operators.NotEquals, overloads.NotEquals,
			func(value1, value2 interface{}) {}, notEquals},

		// Less than operator
		{operators.Less, overloads.LessBool,
			func(value1, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[0].(bool) && values[1].(bool), nil
			}},
		{operators.Less, overloads.LessInt64,
			func(value1, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) < values[1].(int64), nil
			}},
		{operators.Less, overloads.LessUint64,
			func(value1, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) < values[1].(uint64), nil
			}},
		{operators.Less, overloads.LessDouble,
			func(value1, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) < values[1].(float64), nil
			}},
		{operators.Less, overloads.LessString,
			func(value1, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) < values[1].(string), nil
			}},
		{operators.Less, overloads.LessBytes,
			func(value1, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return bytes.Compare(values[0].([]byte), values[1].([]byte)) < 0, nil
			}},
		{operators.Less, overloads.LessTimestamp,
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
		{operators.Less, overloads.LessDuration,
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
		{operators.LessEquals, overloads.LessEqualsBool,
			func(value1, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[0].(bool), nil
			}},
		{operators.LessEquals, overloads.LessEqualsInt64,
			func(value1, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) <= values[1].(int64), nil
			}},
		{operators.LessEquals, overloads.LessEqualsUint64,
			func(value1, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) <= values[1].(uint64), nil
			}},
		{operators.LessEquals, overloads.LessEqualsDouble,
			func(value1, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) <= values[1].(float64), nil
			}},
		{operators.LessEquals, overloads.LessEqualsString,
			func(value1, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) <= values[1].(string), nil
			}},
		{operators.Less, overloads.LessEqualsBytes,
			func(value1, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return bytes.Compare(values[0].([]byte), values[1].([]byte)) <= 0, nil
			}},
		{operators.LessEquals, overloads.LessEqualsTimestamp,
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
		{operators.LessEquals, overloads.LessEqualsDuration,
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
		{operators.Greater, overloads.GreaterBool,
			func(value1 bool, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[1].(bool) && values[0].(bool), nil
			}},
		{operators.Greater, overloads.GreaterInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) > values[1].(int64), nil
			}},
		{operators.Greater, overloads.GreaterUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) > values[1].(uint64), nil
			}},
		{operators.Greater, overloads.GreaterDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) > values[1].(float64), nil
			}},
		{operators.Greater, overloads.GreaterString,
			func(value1 string, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) > values[1].(string), nil
			}},
		{operators.Less, overloads.GreaterBytes,
			func(value1, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return bytes.Compare(values[0].([]byte), values[1].([]byte)) > 0, nil
			}},
		{operators.Greater, overloads.GreaterTimestamp,
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
		{operators.Greater, overloads.GreaterDuration,
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
		{operators.GreaterEquals, overloads.GreaterEqualsBool,
			func(value1 bool, value2 bool) {},
			func(values ...interface{}) (interface{}, error) {
				return !values[1].(bool), nil
			}},
		{operators.GreaterEquals, overloads.GreaterEqualsInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) >= values[1].(int64), nil
			}},
		{operators.GreaterEquals, overloads.GreaterEqualsUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) >= values[1].(uint64), nil
			}},
		{operators.GreaterEquals, overloads.GreaterEqualsDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) >= values[1].(float64), nil
			}},
		{operators.GreaterEquals, overloads.GreaterEqualsString,
			func(value1 string, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) >= values[1].(string), nil
			}},
		{operators.Less, overloads.GreaterEqualsBytes,
			func(value1, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return bytes.Compare(values[0].([]byte), values[1].([]byte)) >= 0, nil
			}},
		{operators.GreaterEquals, overloads.GreaterEqualsTimestamp,
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
		{operators.GreaterEquals, overloads.GreaterEqualsDuration,
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
		{operators.Add, overloads.AddInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) + values[1].(int64), nil
			}},
		{operators.Add, overloads.AddUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) + values[1].(uint64), nil
			}},
		{operators.Add, overloads.AddDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) + values[1].(float64), nil
			}},
		{operators.Add, overloads.AddString,
			func(value1 string, value2 string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(string) + values[1].(string), nil
			}},
		{operators.Add, overloads.AddBytes,
			func(value1 []byte, value2 []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return append(values[0].([]byte), values[1].([]byte)...), nil
			}},
		{operators.Add, overloads.AddList,
			func(value1 types.ListValue, value2 types.ListValue) {},
			func(values ...interface{}) (interface{}, error) {
				list1 := values[0].(types.ListValue)
				list2 := values[1].(types.ListValue)
				return list1.Concat(list2), nil
			}},
		{operators.Add, overloads.AddTimestampDuration,
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
		{operators.Add, overloads.AddDurationTimestamp,
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
		{operators.Add, overloads.AddDurationDuration,
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
		{operators.Subtract, overloads.SubtractInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) - values[1].(int64), nil
			}},
		{operators.Subtract, overloads.SubtractUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) - values[1].(uint64), nil
			}},
		{operators.Subtract, overloads.SubtractDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) - values[1].(float64), nil
			}},
		{operators.Subtract, overloads.SubtractTimestampTimestamp,
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
		{operators.Subtract, overloads.SubtractTimestampDuration,
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
		{operators.Subtract, overloads.SubtractDurationDuration,
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
		{operators.Multiply, overloads.MultiplyInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) * values[1].(int64), nil
			}},
		{operators.Multiply, overloads.MultiplyUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) * values[1].(uint64), nil
			}},
		{operators.Multiply, overloads.MultiplyDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) * values[1].(float64), nil
			}},

		// Divide operator
		// TODO: handle divide by zero.
		{operators.Divide, overloads.DivideInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) / values[1].(int64), nil
			}},
		{operators.Divide, overloads.DivideUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) / values[1].(uint64), nil
			}},
		{operators.Divide, overloads.DivideDouble,
			func(value1 float64, value2 float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(float64) / values[1].(float64), nil
			}},

		// Modulo operator
		{operators.Modulo, overloads.ModuloInt64,
			func(value1 int64, value2 int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(int64) % values[1].(int64), nil
			}},
		{operators.Modulo, overloads.ModuloUint64,
			func(value1 uint64, value2 uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0].(uint64) % values[1].(uint64), nil
			}},

		// Negate operator
		{operators.Negate, overloads.NegateInt64,
			func(value int64) {},
			func(values ...interface{}) (interface{}, error) {
				return -values[0].(int64), nil
			}},
		{operators.Negate, overloads.NegateDouble,
			func(value uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return -values[0].(float64), nil
			}},

		// Index operator
		{operators.Index, operators.Index,
			func(value traits.Indexer, index interface{}) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].(traits.Indexer)
				index := values[1]
				return value.Get(index)
			}},

		// Size function
		{overloads.Size, overloads.SizeString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].(string)
				return len(value), nil
			}},
		{overloads.Size, overloads.SizeBytes,
			func(value []byte) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].([]byte)
				return len(value), nil
			}},
		{overloads.Size, overloads.SizeList,
			func(value types.ListValue) {},
			func(values ...interface{}) (interface{}, error) {
				listValue := values[0].(types.ListValue)
				return listValue.Len(), nil
			}},
		{overloads.Size, overloads.SizeMap,
			func(value types.MapValue) {},
			func(values ...interface{}) (interface{}, error) {
				mapValue := values[0].(types.MapValue)
				return mapValue.Len(), nil
			}},

		// In operator
		{operators.In, overloads.InList,
			func(value interface{}, listValue types.ListValue) {},
			func(values ...interface{}) (interface{}, error) {
				element := values[0]
				listValue := values[1].(types.ListValue)
				return listValue.Contains(element), nil
			}},
		{operators.In, overloads.InMap,
			func(index interface{}, value traits.Indexer) {},
			func(values ...interface{}) (interface{}, error) {
				index := values[0]
				value := values[1].(traits.Indexer)
				// FIXME: This needs to be a Has() method
				_, err := value.Get(index)
				return err == nil, nil
			}},

		// Matches function
		{overloads.MatchString, overloads.MatchString,
			func(text string, pattern string) {},
			func(values ...interface{}) (interface{}, error) {
				text := values[0].(string)
				pattern := values[1].(string)
				return regexp.MatchString(pattern, text)
			}},

		// Timestamp member functions.
		{overloads.TimeGetFullYear, overloads.TimestampToYear,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetFullYear)
			}},
		{overloads.TimeGetMonth, overloads.TimestampToMonth,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetMonth)
			}},
		{overloads.TimeGetDayOfYear, overloads.TimestampToDayOfYear,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfYear)
			}},
		{overloads.TimeGetDate, overloads.TimestampToDayOfMonthZeroBased,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfMonthZeroBased)
			}},
		{overloads.TimeGetDayOfMonth, overloads.TimestampToDayOfMonthOneBased,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfMonthOneBased)
			}},
		{overloads.TimeGetDayOfWeek, overloads.TimestampToDayOfWeek,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetDayOfWeek)
			}},
		{overloads.TimeGetHours, overloads.TimestampToHours,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetHours)
			}},
		{overloads.TimeGetMinutes, overloads.TimestampToMinutes,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetMinutes)
			}},
		{overloads.TimeGetSeconds, overloads.TimestampToSeconds,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetSeconds)
			}},
		{overloads.TimeGetMilliseconds, overloads.TimestampToMilliseconds,
			func(timestamp *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timestampGetMilliseconds)
			}},

		// Timestamp with time zone.
		{overloads.TimeGetFullYear, overloads.TimestampToYearWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetFullYear))
			}},
		{overloads.TimeGetMonth, overloads.TimestampToMonthWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetMonth))
			}},
		{overloads.TimeGetDayOfWeek, overloads.TimestampToDayOfYearWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfYear))
			}},
		{overloads.TimeGetDate, overloads.TimestampToDayOfMonthZeroBasedWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfMonthZeroBased))
			}},
		{overloads.TimeGetDayOfMonth, overloads.TimestampToDayOfMonthOneBasedWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfMonthOneBased))
			}},
		{overloads.TimeGetDayOfWeek, overloads.TimestampToDayOfWeekWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetDayOfWeek))
			}},
		{overloads.TimeGetHours, overloads.TimestampToHoursWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetHours))
			}},
		{overloads.TimeGetMinutes, overloads.TimestampToMinutesWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetMinutes))
			}},
		{overloads.TimeGetSeconds, overloads.TimestampToSecondsWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetSeconds))
			}},
		{overloads.TimeGetMilliseconds, overloads.TimestampToMillisecondsWithTz,
			func(timestamp *tspb.Timestamp, tz string) {},
			func(values ...interface{}) (interface{}, error) {
				return timestampOp(values[0], timeZone(values[1], timestampGetMilliseconds))
			}},

		// Duration member functions.
		{overloads.TimeGetHours, overloads.DurationToHours,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Hours(), nil
				} else {
					return nil, err
				}
			}},
		{overloads.TimeGetMinutes, overloads.DurationToMinutes,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Minutes(), nil
				} else {
					return nil, err
				}
			}},
		{overloads.TimeGetSeconds, overloads.DurationToSeconds,
			func(duration *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				duration := values[0].(*dpb.Duration)
				if dur, err := ptypes.Duration(duration); err == nil {
					return dur.Seconds(), nil
				} else {
					return nil, err
				}
			}},
		{overloads.TimeGetMilliseconds, overloads.DurationToMilliseconds,
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
		{overloads.TypeConvertInt, overloads.IntToInt,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertInt, overloads.UintToInt,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(uint64)
				return int64(num), nil
			}},
		{overloads.TypeConvertInt, overloads.DoubleToInt,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(float64)
				return int64(num), nil
			}},
		{overloads.TypeConvertInt, overloads.StringToInt,
			func(num string) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(string)
				return strconv.ParseInt(num, 10, 64)
			}},
		{overloads.TypeConvertInt, overloads.TimestampToInt,
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
		{overloads.TypeConvertInt, overloads.DurationToInt,
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
		{overloads.TypeConvertUint, overloads.UintToUint,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertUint, overloads.IntToUint,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(int64)
				if num >= 0 {
					return uint64(num), nil
				}
				return nil,
					fmt.Errorf("unsafe uint conversion of negative int")
			}},
		{overloads.TypeConvertUint, overloads.DoubleToUint,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(float64)
				if num >= 0.0 {
					return float64(num), nil
				}
				return nil,
					fmt.Errorf("unsafe uint conversion of negative double")
			}},
		{overloads.TypeConvertUint, overloads.StringToUint,
			func(num string) {},
			func(values ...interface{}) (interface{}, error) {
				return strconv.ParseUint(values[0].(string), 10, 64)
			}},

		// Double conversions.
		{overloads.TypeConvertDouble, overloads.DoubleToDouble,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertDouble, overloads.IntToDouble,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(int64)
				return float64(num), nil
			}},
		{overloads.TypeConvertDouble, overloads.UintToDouble,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				num := values[0].(uint64)
				return float64(num), nil
			}},
		{overloads.TypeConvertDouble, overloads.StringToDouble,
			func(num string) {},
			func(values ...interface{}) (interface{}, error) {
				return strconv.ParseFloat(values[0].(string), 64)
			}},

		// Bool conversions.
		{overloads.TypeConvertBool, overloads.BoolToBool,
			func(value bool) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertBool, overloads.StringToBool,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return strconv.ParseBool(values[0].(string))
			}},

		// Bytes conversions.
		{overloads.TypeConvertBytes, overloads.BytesToBytes,
			func(value []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertBytes, overloads.StringToBytes,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return []byte(values[0].(string)), nil
			}},

		// String conversions.
		{overloads.TypeConvertString, overloads.StringToString,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertString, overloads.BoolToString,
			func(value bool) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%t", values[0].(bool)), nil
			}},
		{overloads.TypeConvertString, overloads.IntToString,
			func(num int64) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%d", values[0].(int64)), nil
			}},
		{overloads.TypeConvertString, overloads.UintToString,
			func(num uint64) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%d", values[0].(uint64)), nil
			}},
		{overloads.TypeConvertString, overloads.DoubleToString,
			func(num float64) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%g", values[0].(float64)), nil
			}},
		{overloads.TypeConvertString, overloads.BytesToString,
			func(value []byte) {},
			func(values ...interface{}) (interface{}, error) {
				return fmt.Sprintf("%s", values[0].([]byte)), nil
			}},
		{overloads.TypeConvertString, overloads.TimestampToString,
			func(value *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return ptypes.TimestampString(values[0].(*tspb.Timestamp)), nil
			}},
		{overloads.TypeConvertString, overloads.DurationToString,
			func(value *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				return ptypes.TimestampString(values[0].(*tspb.Timestamp)), nil
			}},

		// Timestamp conversions.
		{overloads.TypeConvertTimestamp, overloads.TimestampToTimestamp,
			func(value *tspb.Timestamp) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertTimestamp, overloads.StringToTimestamp,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				if t, err := time.Parse(time.RFC3339, values[0].(string)); err != nil {
					return nil, err
				} else {
					return ptypes.TimestampProto(t)
				}
			}},
		{overloads.TypeConvertTimestamp, overloads.IntToTimestamp,
			func(value int64) {},
			func(values ...interface{}) (interface{}, error) {
				value := values[0].(int64)
				// Unix timestamp in seconds.
				return ptypes.TimestampProto(time.Unix(value, 0))
			}},
		// Duration conversions.
		{overloads.TypeConvertDuration, overloads.DurationToDuration,
			func(value *dpb.Duration) {},
			func(values ...interface{}) (interface{}, error) {
				return values[0], nil
			}},
		{overloads.TypeConvertDuration, overloads.StringToDuration,
			func(value string) {},
			func(values ...interface{}) (interface{}, error) {
				return time.ParseDuration(values[0].(string))
			}},
		{overloads.TypeConvertDuration, overloads.IntToDuration,
			func(value int64) {},
			func(values ...interface{}) (interface{}, error) {
				// Duration in seconds
				return time.Duration(values[0].(int64)) * time.Second, nil
			}},

		// Type operations.
		{overloads.TypeConvertType, overloads.TypeConvertType,
			func(value types.Type) {},
			func(values ...interface{}) (interface{}, error) {
				if t, found := types.TypeOf(values[0]); found {
					return t, nil
				}
				return nil, fmt.Errorf("undefined type")
			}},

		// Comprehensions functions.
		{overloads.Iterator, overloads.Iterator,
			func(value traits.Iterable) {},
			func(values ...interface{}) (interface{}, error) {
				iterable := values[0].(traits.Iterable)
				return iterable.Iterator(), nil
			}},
		{overloads.HasNext, overloads.HasNext,
			func(value traits.Iterator) {},
			func(values ...interface{}) (interface{}, error) {
				it := values[0].(traits.Iterator)
				return it.HasNext(), nil
			}},
		{overloads.Next, overloads.Next,
			func(value traits.Iterator) {},
			func(values ...interface{}) (interface{}, error) {
				it := values[0].(traits.Iterator)
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

func newError(msg string, value interface{}) interface{} {
	switch value.(type) {
	case error:
		return value
	default:
		return fmt.Errorf(msg, value)
	}
}

func equals(values ...interface{}) (interface{}, error) {
	if value1, ok := values[0].(traits.Equaler); ok {
		return value1.Equal(values[1]), nil
	}
	if value2, ok := values[1].(traits.Equaler); ok {
		return value2.Equal(values[0]), nil
	}
	if value1, ok := values[0].([]byte); ok {
		if value2, ok := values[1].([]byte); ok {
			return bytes.Equal(value1, value2), nil
		}
	}
	return values[0] == values[1], nil
}

func notEquals(values ...interface{}) (interface{}, error) {
	if val, err := equals(values); err != nil {
		return nil, err
	} else {
		return !val.(bool), nil
	}
}
