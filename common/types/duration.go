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
	dpb "github.com/golang/protobuf/ptypes/duration"
	overloadspb "github.com/google/cel-go/common/overloads"
	refpb "github.com/google/cel-go/common/types/ref"
	traitspb "github.com/google/cel-go/common/types/traits"
)

// Duration type that implements refpb.Value and supports add, compare, negate,
// and subtract operators. This type is also a receiver which means it can
// participate in dispatch to receiver functions.
type Duration struct {
	*dpb.Duration
}

var (
	// DurationType singleton.
	DurationType = NewTypeValue("google.protobuf.Duration",
		traitspb.AdderType,
		traitspb.ComparerType,
		traitspb.NegatorType,
		traitspb.ReceiverType,
		traitspb.SubtractorType)
)

func (d Duration) Add(other refpb.Value) refpb.Value {
	switch other.Type() {
	case DurationType:
		dur1, err := ptypespb.Duration(d.Duration)
		if err != nil {
			return &Err{err}
		}
		dur2, err := ptypespb.Duration(other.(Duration).Duration)
		if err != nil {
			return &Err{err}
		}
		return Duration{ptypespb.DurationProto(dur1 + dur2)}
	case TimestampType:
		dur, err := ptypespb.Duration(d.Duration)
		if err != nil {
			return &Err{err}
		}
		ts, err := ptypespb.Timestamp(other.(Timestamp).Timestamp)
		if err != nil {
			return &Err{err}
		}
		tstamp, err := ptypespb.TimestampProto(ts.Add(dur))
		if err != nil {
			return &Err{err}
		}
		return Timestamp{tstamp}
	}
	return NewErr("unsupported overload")
}

func (d Duration) Compare(other refpb.Value) refpb.Value {
	if DurationType != other.Type() {
		return NewErr("unsupported overload")
	}
	dur1, err := ptypespb.Duration(d.Duration)
	if err != nil {
		return &Err{err}
	}
	dur2, err := ptypespb.Duration(other.(Duration).Duration)
	if err != nil {
		return &Err{err}
	}
	dur := dur1 - dur2
	if dur < 0 {
		return IntNegOne
	}
	if dur > 0 {
		return IntOne
	}
	return IntZero
}

func (d Duration) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if typeDesc == durationValueType {
		return d.Value(), nil
	}
	// If the duration is already assignable to the desired type return it.
	if reflect.TypeOf(d).AssignableTo(typeDesc) {
		return d, nil
	}
	return nil, fmt.Errorf("type conversion error from "+
		"'google.protobuf.Duration' to '%v'", typeDesc)
}

func (d Duration) ConvertToType(typeVal refpb.Type) refpb.Value {
	switch typeVal {
	case StringType:
		if dur, err := ptypespb.Duration(d.Duration); err == nil {
			return String(dur.String())
		}
	case IntType:
		if dur, err := ptypespb.Duration(d.Duration); err == nil {
			return Int(dur)
		}
	case DurationType:
		return d
	case TypeType:
		return DurationType
	}
	return NewErr("type conversion error from '%s' to '%s'", DurationType, typeVal)
}

func (d Duration) Equal(other refpb.Value) refpb.Value {
	return Bool(DurationType == other.Type() &&
		protopb.Equal(d.Duration, other.Value().(protopb.Message)))
}

func (d Duration) Negate() refpb.Value {
	dur, err := ptypespb.Duration(d.Duration)
	if err != nil {
		return &Err{err}
	}
	return Duration{ptypespb.DurationProto(-dur)}
}

func (d Duration) Receive(function string, overload string, args []refpb.Value) refpb.Value {
	dur, err := ptypespb.Duration(d.Duration)
	if err != nil {
		return &Err{err}
	}
	if len(args) == 0 {
		if f, found := durationZeroArgOverloads[function]; found {
			return f(dur)
		}
	}
	return NewErr("unsupported overload")
}

func (d Duration) Subtract(subtrahend refpb.Value) refpb.Value {
	if DurationType != subtrahend.Type() {
		return NewErr("unsupported overload")
	}
	return d.Add(subtrahend.(Duration).Negate())
}

func (d Duration) Type() refpb.Type {
	return DurationType
}

func (d Duration) Value() interface{} {
	return d.Duration
}

var (
	durationValueType = reflect.TypeOf(&dpb.Duration{})

	durationZeroArgOverloads = map[string]func(time.Duration) refpb.Value{
		overloadspb.TimeGetHours: func(dur time.Duration) refpb.Value {
			return Int(dur.Hours())
		},
		overloadspb.TimeGetMinutes: func(dur time.Duration) refpb.Value {
			return Int(dur.Minutes())
		},
		overloadspb.TimeGetSeconds: func(dur time.Duration) refpb.Value {
			return Int(dur.Seconds())
		},
		overloadspb.TimeGetMilliseconds: func(dur time.Duration) refpb.Value {
			return Int(dur.Nanoseconds() / 1000000)
		}}
)
