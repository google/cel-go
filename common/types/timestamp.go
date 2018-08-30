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

	protopb "github.com/golang/protobuf/proto"
	ptypespb "github.com/golang/protobuf/ptypes"
	tpb "github.com/golang/protobuf/ptypes/timestamp"
	overloadspb "github.com/google/cel-go/common/overloads"
	traitspb "github.com/google/cel-go/common/types/traits"
	refpb "github.com/google/cel-go/common/types/ref"
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
		traitspb.AdderType,
		traitspb.ComparerType,
		traitspb.ReceiverType,
		traitspb.SubtractorType)
)

func (t Timestamp) Add(other refpb.Value) refpb.Value {
	switch other.Type() {
	case DurationType:
		return other.(Duration).Add(t)
	}
	return NewErr("unsupported overload")
}

func (t Timestamp) Compare(other refpb.Value) refpb.Value {
	if TimestampType != other.Type() {
		return NewErr("unsupported overload")
	}
	ts1, err := ptypespb.Timestamp(t.Timestamp)
	if err != nil {
		return &Err{err}
	}
	ts2, err := ptypespb.Timestamp(other.(Timestamp).Timestamp)
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

func (t Timestamp) ConvertToType(typeVal refpb.Type) refpb.Value {
	switch typeVal {
	case StringType:
		return String(ptypespb.TimestampString(t.Timestamp))
	case IntType:
		if ts, err := ptypespb.Timestamp(t.Timestamp); err == nil {
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

func (t Timestamp) Equal(other refpb.Value) refpb.Value {
	return Bool(TimestampType == other.Type() &&
		protopb.Equal(t.Timestamp, other.Value().(protopb.Message)))
}

func (t Timestamp) Receive(function string, overload string, args []refpb.Value) refpb.Value {
	ts := t.Timestamp
	tstamp, err := ptypespb.Timestamp(ts)
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
	return NewErr("unsupported overload")
}

func (t Timestamp) Subtract(subtrahend refpb.Value) refpb.Value {
	switch subtrahend.Type() {
	case DurationType:
		ts, err := ptypespb.Timestamp(t.Timestamp)
		if err != nil {
			return &Err{err}
		}
		dur, err := ptypespb.Duration(subtrahend.(Duration).Duration)
		if err != nil {
			return &Err{err}
		}
		tstamp, err := ptypespb.TimestampProto(ts.Add(-dur))
		if err != nil {
			return &Err{err}
		}
		return Timestamp{tstamp}
	case TimestampType:
		ts1, err := ptypespb.Timestamp(t.Timestamp)
		if err != nil {
			return &Err{err}
		}
		ts2, err := ptypespb.Timestamp(subtrahend.(Timestamp).Timestamp)
		if err != nil {
			return &Err{err}
		}
		return Duration{ptypespb.DurationProto(ts1.Sub(ts2))}
	}
	return NewErr("unsupported overload")
}

func (t Timestamp) Type() refpb.Type {
	return TimestampType
}

func (t Timestamp) Value() interface{} {
	return t.Timestamp
}

var (
	timestampValueType = reflect.TypeOf(&tpb.Timestamp{})

	timestampZeroArgOverloads = map[string]func(time.Time) refpb.Value{
		overloadspb.TimeGetFullYear:     timestampGetFullYear,
		overloadspb.TimeGetMonth:        timestampGetMonth,
		overloadspb.TimeGetDayOfYear:    timestampGetDayOfYear,
		overloadspb.TimeGetDate:         timestampGetDayOfMonthOneBased,
		overloadspb.TimeGetDayOfMonth:   timestampGetDayOfMonthZeroBased,
		overloadspb.TimeGetDayOfWeek:    timestampGetDayOfWeek,
		overloadspb.TimeGetHours:        timestampGetHours,
		overloadspb.TimeGetMinutes:      timestampGetMinutes,
		overloadspb.TimeGetSeconds:      timestampGetSeconds,
		overloadspb.TimeGetMilliseconds: timestampGetMilliseconds}

	timestampOneArgOverloads = map[string]func(time.Time, refpb.Value) refpb.Value{
		overloadspb.TimeGetFullYear:     timestampGetFullYearWithTz,
		overloadspb.TimeGetMonth:        timestampGetMonthWithTz,
		overloadspb.TimeGetDayOfYear:    timestampGetDayOfYearWithTz,
		overloadspb.TimeGetDate:         timestampGetDayOfMonthOneBasedWithTz,
		overloadspb.TimeGetDayOfMonth:   timestampGetDayOfMonthZeroBasedWithTz,
		overloadspb.TimeGetDayOfWeek:    timestampGetDayOfWeekWithTz,
		overloadspb.TimeGetHours:        timestampGetHoursWithTz,
		overloadspb.TimeGetMinutes:      timestampGetMinutesWithTz,
		overloadspb.TimeGetSeconds:      timestampGetSecondsWithTz,
		overloadspb.TimeGetMilliseconds: timestampGetMillisecondsWithTz}
)

type timestampVisitor func(time.Time) refpb.Value

func timestampGetFullYear(t time.Time) refpb.Value {
	return Int(t.Year())
}
func timestampGetMonth(t time.Time) refpb.Value {
	return Int(t.Month())
}
func timestampGetDayOfYear(t time.Time) refpb.Value {
	return Int(t.YearDay())
}
func timestampGetDayOfMonthZeroBased(t time.Time) refpb.Value {
	return Int(t.Day() - 1)
}
func timestampGetDayOfMonthOneBased(t time.Time) refpb.Value {
	return Int(t.Day())
}
func timestampGetDayOfWeek(t time.Time) refpb.Value {
	return Int(t.Weekday())
}
func timestampGetHours(t time.Time) refpb.Value {
	return Int(t.Hour())
}
func timestampGetMinutes(t time.Time) refpb.Value {
	return Int(t.Minute())
}
func timestampGetSeconds(t time.Time) refpb.Value {
	return Int(t.Second())
}
func timestampGetMilliseconds(t time.Time) refpb.Value {
	return Int(t.Nanosecond() / 1000000)
}

func timestampGetFullYearWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetFullYear)(t)
}
func timestampGetMonthWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetMonth)(t)
}
func timestampGetDayOfYearWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetDayOfYear)(t)
}
func timestampGetDayOfMonthZeroBasedWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetDayOfMonthZeroBased)(t)
}
func timestampGetDayOfMonthOneBasedWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetDayOfMonthOneBased)(t)
}
func timestampGetDayOfWeekWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetDayOfWeek)(t)
}
func timestampGetHoursWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetHours)(t)
}
func timestampGetMinutesWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetMinutes)(t)
}
func timestampGetSecondsWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetSeconds)(t)
}
func timestampGetMillisecondsWithTz(t time.Time, tz refpb.Value) refpb.Value {
	return timeZone(tz, timestampGetMilliseconds)(t)
}

func timeZone(tz refpb.Value, visitor timestampVisitor) timestampVisitor {
	return func(t time.Time) refpb.Value {
		if StringType != tz.Type() {
			return NewErr("unsupported overload")
		}
		if loc, err := time.LoadLocation(string(tz.(String))); err == nil {
			return visitor(t.In(loc))
		} else {
			return &Err{err}
		}
	}
}
