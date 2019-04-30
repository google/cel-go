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

package types

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	tpb "github.com/golang/protobuf/ptypes/timestamp"
)

// Timestamp type implementation which supports add, compare, and subtract
// operations. Timestamps are also capable of participating in dynamic
// function dispatch to instance methods.
type Timestamp struct {
	*tpb.Timestamp
}

var (
	// TimestampType singleton.
	TimestampType = NewTypeValue("google.protobuf.Timestamp",
		traits.AdderType,
		traits.ComparerType,
		traits.ReceiverType,
		traits.SubtractorType)
)

// Add implements traits.Adder.Add.
func (t Timestamp) Add(other ref.Val) ref.Val {
	switch other.Type() {
	case DurationType:
		return other.(Duration).Add(t)
	}
	return ValOrErr(other, "no such overload")
}

// Compare implements traits.Comparer.Compare.
func (t Timestamp) Compare(other ref.Val) ref.Val {
	if TimestampType != other.Type() {
		return ValOrErr(other, "no such overload")
	}
	ts1, err := ptypes.Timestamp(t.Timestamp)
	if err != nil {
		return &Err{err}
	}
	ts2, err := ptypes.Timestamp(other.(Timestamp).Timestamp)
	if err != nil {
		return &Err{err}
	}
	ts := ts1.Sub(ts2)
	if ts < 0 {
		return IntNegOne
	}
	if ts > 0 {
		return IntOne
	}
	return IntZero
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (t Timestamp) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if typeDesc == timestampValueType {
		return t.Value(), nil
	}
	// If the timestamp is already assignable to the desired type return it.
	if reflect.TypeOf(t).AssignableTo(typeDesc) {
		return t, nil
	}
	return nil, fmt.Errorf("type conversion error from "+
		"'google.protobuf.Duration' to '%v'", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (t Timestamp) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case StringType:
		return String(ptypes.TimestampString(t.Timestamp))
	case IntType:
		if ts, err := ptypes.Timestamp(t.Timestamp); err == nil {
			// Return the Unix time in seconds since 1970
			return Int(ts.Unix())
		}
	case TimestampType:
		return t
	case TypeType:
		return TimestampType
	}
	return NewErr("type conversion error from '%s' to '%s'", TimestampType, typeVal)
}

// Equal implements ref.Val.Equal.
func (t Timestamp) Equal(other ref.Val) ref.Val {
	if TimestampType != other.Type() {
		return ValOrErr(other, "no such overload")
	}
	return Bool(proto.Equal(t.Timestamp, other.Value().(proto.Message)))
}

// Receive implements traits.Reciever.Receive.
func (t Timestamp) Receive(function string, overload string, args []ref.Val) ref.Val {
	ts := t.Timestamp
	tstamp, err := ptypes.Timestamp(ts)
	if err != nil {
		return &Err{err}
	}
	switch len(args) {
	case 0:
		if f, found := timestampZeroArgOverloads[function]; found {
			return f(tstamp)
		}
	case 1:
		if f, found := timestampOneArgOverloads[function]; found {
			return f(tstamp, args[0])
		}
	}
	return NewErr("no such overload")
}

// Subtract implements traits.Subtractor.Subtract.
func (t Timestamp) Subtract(subtrahend ref.Val) ref.Val {
	switch subtrahend.Type() {
	case DurationType:
		ts, err := ptypes.Timestamp(t.Timestamp)
		if err != nil {
			return &Err{err}
		}
		dur, err := ptypes.Duration(subtrahend.(Duration).Duration)
		if err != nil {
			return &Err{err}
		}
		tstamp, err := ptypes.TimestampProto(ts.Add(-dur))
		if err != nil {
			return &Err{err}
		}
		return Timestamp{tstamp}
	case TimestampType:
		ts1, err := ptypes.Timestamp(t.Timestamp)
		if err != nil {
			return &Err{err}
		}
		ts2, err := ptypes.Timestamp(subtrahend.(Timestamp).Timestamp)
		if err != nil {
			return &Err{err}
		}
		return Duration{ptypes.DurationProto(ts1.Sub(ts2))}
	}
	return ValOrErr(subtrahend, "no such overload")
}

// Type implements ref.Val.Type.
func (t Timestamp) Type() ref.Type {
	return TimestampType
}

// Value implements ref.Val.Value.
func (t Timestamp) Value() interface{} {
	return t.Timestamp
}

var (
	timestampValueType = reflect.TypeOf(&tpb.Timestamp{})

	timestampZeroArgOverloads = map[string]func(time.Time) ref.Val{
		overloads.TimeGetFullYear:     timestampGetFullYear,
		overloads.TimeGetMonth:        timestampGetMonth,
		overloads.TimeGetDayOfYear:    timestampGetDayOfYear,
		overloads.TimeGetDate:         timestampGetDayOfMonthOneBased,
		overloads.TimeGetDayOfMonth:   timestampGetDayOfMonthZeroBased,
		overloads.TimeGetDayOfWeek:    timestampGetDayOfWeek,
		overloads.TimeGetHours:        timestampGetHours,
		overloads.TimeGetMinutes:      timestampGetMinutes,
		overloads.TimeGetSeconds:      timestampGetSeconds,
		overloads.TimeGetMilliseconds: timestampGetMilliseconds}

	timestampOneArgOverloads = map[string]func(time.Time, ref.Val) ref.Val{
		overloads.TimeGetFullYear:     timestampGetFullYearWithTz,
		overloads.TimeGetMonth:        timestampGetMonthWithTz,
		overloads.TimeGetDayOfYear:    timestampGetDayOfYearWithTz,
		overloads.TimeGetDate:         timestampGetDayOfMonthOneBasedWithTz,
		overloads.TimeGetDayOfMonth:   timestampGetDayOfMonthZeroBasedWithTz,
		overloads.TimeGetDayOfWeek:    timestampGetDayOfWeekWithTz,
		overloads.TimeGetHours:        timestampGetHoursWithTz,
		overloads.TimeGetMinutes:      timestampGetMinutesWithTz,
		overloads.TimeGetSeconds:      timestampGetSecondsWithTz,
		overloads.TimeGetMilliseconds: timestampGetMillisecondsWithTz}
)

type timestampVisitor func(time.Time) ref.Val

func timestampGetFullYear(t time.Time) ref.Val {
	return Int(t.Year())
}
func timestampGetMonth(t time.Time) ref.Val {
	// CEL spec indicates that the month should be 0-based, but the Time value
	// for Month() is 1-based.
	return Int(t.Month() - 1)
}
func timestampGetDayOfYear(t time.Time) ref.Val {
	return Int(t.YearDay())
}
func timestampGetDayOfMonthZeroBased(t time.Time) ref.Val {
	return Int(t.Day() - 1)
}
func timestampGetDayOfMonthOneBased(t time.Time) ref.Val {
	return Int(t.Day())
}
func timestampGetDayOfWeek(t time.Time) ref.Val {
	return Int(t.Weekday())
}
func timestampGetHours(t time.Time) ref.Val {
	return Int(t.Hour())
}
func timestampGetMinutes(t time.Time) ref.Val {
	return Int(t.Minute())
}
func timestampGetSeconds(t time.Time) ref.Val {
	return Int(t.Second())
}
func timestampGetMilliseconds(t time.Time) ref.Val {
	return Int(t.Nanosecond() / 1000000)
}

func timestampGetFullYearWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetFullYear)(t)
}
func timestampGetMonthWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetMonth)(t)
}
func timestampGetDayOfYearWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetDayOfYear)(t)
}
func timestampGetDayOfMonthZeroBasedWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetDayOfMonthZeroBased)(t)
}
func timestampGetDayOfMonthOneBasedWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetDayOfMonthOneBased)(t)
}
func timestampGetDayOfWeekWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetDayOfWeek)(t)
}
func timestampGetHoursWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetHours)(t)
}
func timestampGetMinutesWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetMinutes)(t)
}
func timestampGetSecondsWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetSeconds)(t)
}
func timestampGetMillisecondsWithTz(t time.Time, tz ref.Val) ref.Val {
	return timeZone(tz, timestampGetMilliseconds)(t)
}

func timeZone(tz ref.Val, visitor timestampVisitor) timestampVisitor {
	return func(t time.Time) ref.Val {
		if StringType != tz.Type() {
			return ValOrErr(tz, "no such overload")
		}
		loc, err := time.LoadLocation(string(tz.(String)))
		if err == nil {
			return visitor(t.In(loc))
		}
		return &Err{err}
	}
}
