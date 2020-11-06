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
	"strconv"
	"time"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	dpb "google.golang.org/protobuf/types/known/durationpb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

// Duration type that implements ref.Val and supports add, compare, negate,
// and subtract operators. This type is also a receiver which means it can
// participate in dispatch to receiver functions.
type Duration struct {
	*dpb.Duration
}

var (
	// DurationType singleton.
	DurationType = NewTypeValue("google.protobuf.Duration",
		traits.AdderType,
		traits.ComparerType,
		traits.NegatorType,
		traits.ReceiverType,
		traits.SubtractorType)
)

// Add implements traits.Adder.Add.
func (d Duration) Add(other ref.Val) ref.Val {
	switch other.Type() {
	case DurationType:
		dur1 := d.AsDuration()
		dur2 := other.(Duration).AsDuration()
		return Duration{Duration: dpb.New(dur1 + dur2)}
	case TimestampType:
		dur := d.AsDuration()
		ts := other.(Timestamp).AsTime()
		return Timestamp{Timestamp: tpb.New(ts.Add(dur))}
	}
	return ValOrErr(other, "no such overload")
}

// Compare implements traits.Comparer.Compare.
func (d Duration) Compare(other ref.Val) ref.Val {
	otherDur, ok := other.(Duration)
	if !ok {
		return ValOrErr(other, "no such overload")
	}
	dur1 := d.AsDuration()
	dur2 := otherDur.AsDuration()
	dur := dur1 - dur2
	if dur < 0 {
		return IntNegOne
	}
	if dur > 0 {
		return IntOne
	}
	return IntZero
}

// ConvertToNative implements ref.Val.ConvertToNative.
func (d Duration) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	switch typeDesc {
	case anyValueType:
		// Pack the underlying proto value into an Any value.
		return anypb.New(d.Value().(*dpb.Duration))
	case durationValueType:
		// Unwrap the CEL value to its underlying proto value.
		return d.Value(), nil
	case jsonValueType:
		// CEL follows the proto3 to JSON conversion.
		// Note, using jsonpb would wrap the result in extra double quotes.
		v := d.ConvertToType(StringType)
		if IsError(v) {
			return nil, v.(*Err)
		}
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: string(v.(String))},
		}, nil
	}
	// If the duration is already assignable to the desired type return it.
	if reflect.TypeOf(d).AssignableTo(typeDesc) {
		return d, nil
	}
	return nil, fmt.Errorf("type conversion error from "+
		"'google.protobuf.Duration' to '%v'", typeDesc)
}

// ConvertToType implements ref.Val.ConvertToType.
func (d Duration) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case StringType:
		dur := d.AsDuration()
		return String(strconv.FormatFloat(dur.Seconds(), 'f', -1, 64) + "s")
	case IntType:
		dur := d.AsDuration()
		return Int(dur)
	case DurationType:
		return d
	case TypeType:
		return DurationType
	}
	return NewErr("type conversion error from '%s' to '%s'", DurationType, typeVal)
}

// Equal implements ref.Val.Equal.
func (d Duration) Equal(other ref.Val) ref.Val {
	otherDur, ok := other.(Duration)
	if !ok {
		return ValOrErr(other, "no such overload")
	}
	return Bool(proto.Equal(d.Duration, otherDur.Value().(proto.Message)))
}

// Negate implements traits.Negater.Negate.
func (d Duration) Negate() ref.Val {
	dur := d.AsDuration()
	return Duration{Duration: dpb.New(-dur)}
}

// Receive implements traits.Receiver.Receive.
func (d Duration) Receive(function string, overload string, args []ref.Val) ref.Val {
	dur := d.AsDuration()
	if len(args) == 0 {
		if f, found := durationZeroArgOverloads[function]; found {
			return f(dur)
		}
	}
	return NewErr("no such overload")
}

// Subtract implements traits.Subtractor.Subtract.
func (d Duration) Subtract(subtrahend ref.Val) ref.Val {
	subtraDur, ok := subtrahend.(Duration)
	if !ok {
		return ValOrErr(subtrahend, "no such overload")
	}
	return d.Add(subtraDur.Negate())
}

// Type implements ref.Val.Type.
func (d Duration) Type() ref.Type {
	return DurationType
}

// Value implements ref.Val.Value.
func (d Duration) Value() interface{} {
	return d.Duration
}

var (
	durationValueType = reflect.TypeOf(&dpb.Duration{})

	durationZeroArgOverloads = map[string]func(time.Duration) ref.Val{
		overloads.TimeGetHours: func(dur time.Duration) ref.Val {
			return Int(dur.Hours())
		},
		overloads.TimeGetMinutes: func(dur time.Duration) ref.Val {
			return Int(dur.Minutes())
		},
		overloads.TimeGetSeconds: func(dur time.Duration) ref.Val {
			return Int(dur.Seconds())
		},
		overloads.TimeGetMilliseconds: func(dur time.Duration) ref.Val {
			return Int(dur.Nanoseconds() / 1000000)
		}}
)
