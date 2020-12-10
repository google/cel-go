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
	"reflect"
	"testing"
	"time"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"

	"google.golang.org/protobuf/proto"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestampAdd(t *testing.T) {
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	val := ts.Add(Duration{time.Duration(3600)*time.Second + time.Duration(1000)})
	if val.ConvertToType(TypeType) != TimestampType {
		t.Error("Could not add duration and timestamp")
	}
	expected := Timestamp{Time: time.Unix(11106, 1000).UTC()}
	if !expected.Compare(val).Equal(IntZero).(Bool) {
		t.Errorf("Got '%v', expected '%v'", val, expected)
	}
	if !IsError(ts.Add(expected)) {
		t.Error("Cannot add two timestamps together")
	}
}

func TestTimestampConvertToNative_Any(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0)}
	val, err := ts.ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	want, err := anypb.New(tpb.New(ts.Time))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
}

func TestTimestampConvertToNative(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	val, err := ts.ConvertToNative(timestampValueType)
	if err != nil {
		t.Error(err)
	}
	var want interface{}
	want = tpb.New(ts.Time)
	if !proto.Equal(val.(proto.Message), want.(proto.Message)) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
	val, err = ts.ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	}
	want = structpb.NewStringValue("1970-01-01T02:05:06Z")
	if !proto.Equal(val.(proto.Message), want.(proto.Message)) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
	val, err = ts.ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	want, err = anypb.New(tpb.New(ts.Time))
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message), want.(proto.Message)) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
	val, err = ts.ConvertToNative(reflect.TypeOf(Timestamp{}))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(val, ts) {
		t.Errorf("got %v wanted %v", val, ts)
	}
	val, err = ts.ConvertToNative(reflect.TypeOf(time.Now()))
	if err != nil {
		t.Error(err)
	}
	want = time.Unix(7506, 0).UTC()
	if !reflect.DeepEqual(val, want) {
		t.Errorf("got %v wanted %v", val, want)
	}
}

func TestTimestampSubtract(t *testing.T) {
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	val := ts.Subtract(Duration{Duration: time.Duration(3600)*time.Second + time.Duration(1000)})
	if val.ConvertToType(TypeType) != TimestampType {
		t.Error("Could not add duration and timestamp")
	}
	expected := Timestamp{Time: time.Unix(3905, 999999000).UTC()}
	if !expected.Compare(val).Equal(IntZero).(Bool) {
		t.Errorf("Got '%v', expected '%v'", val, expected)
	}
	ts2 := Timestamp{Time: time.Unix(6506, 0).UTC()}
	val = ts.Subtract(ts2)
	if val.ConvertToType(TypeType) != DurationType {
		t.Error("Could not subtract timestamps")
	}
	expectedDur := Duration{Duration: time.Duration(1000000000000)}
	if !expectedDur.Compare(val).Equal(IntZero).(Bool) {
		t.Errorf("Got '%v', expected '%v'", val, expectedDur)
	}
}

func TestTimestampGetDayOfYear(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	hr := ts.Receive(overloads.TimeGetDayOfYear, overloads.TimestampToDayOfYear, []ref.Val{})
	if !hr.Equal(Int(0)).(Bool) {
		t.Error("Expected 0, got", hr)
	}
	// 1969-12-31T19:05:06Z
	hrTz := ts.Receive(overloads.TimeGetDayOfYear, overloads.TimestampToDayOfYearWithTz,
		[]ref.Val{String("America/Phoenix")})
	if !hrTz.Equal(Int(364)).(Bool) {
		t.Error("Expected 364, got", hrTz)
	}
	hrTz = ts.Receive(overloads.TimeGetDayOfYear, overloads.TimestampToDayOfYearWithTz,
		[]ref.Val{String("-07:00")})
	if !hrTz.Equal(Int(364)).(Bool) {
		t.Error("Expected 364, got", hrTz)
	}
}

func TestTimestampGetMonth(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	hr := ts.Receive(overloads.TimeGetMonth, overloads.TimestampToMonth, []ref.Val{})
	if !hr.Equal(Int(0)).(Bool) {
		t.Error("Expected 0, got", hr)
	}
	// 1969-12-31T19:05:06Z
	hrTz := ts.Receive(overloads.TimeGetMonth, overloads.TimestampToMonthWithTz,
		[]ref.Val{String("America/Phoenix")})
	if !hrTz.Equal(Int(11)).(Bool) {
		t.Error("Expected 11, got", hrTz)
	}
}

func TestTimestampGetHours(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	hr := ts.Receive(overloads.TimeGetHours, overloads.TimestampToHours, []ref.Val{})
	if !hr.Equal(Int(2)).(Bool) {
		t.Error("Expected 2 hours, got", hr)
	}
	// 1969-12-31T19:05:06Z
	hrTz := ts.Receive(overloads.TimeGetHours, overloads.TimestampToHoursWithTz,
		[]ref.Val{String("America/Phoenix")})
	if !hrTz.Equal(Int(19)).(Bool) {
		t.Error("Expected 19 hours, got", hrTz)
	}
}

func TestTimestampGetMinutes(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	min := ts.Receive(overloads.TimeGetMinutes, overloads.TimestampToMinutes, []ref.Val{})
	if !min.Equal(Int(5)).(Bool) {
		t.Error("Expected 5 minutes, got", min)
	}
	// 1969-12-31T19:05:06Z
	minTz := ts.Receive(overloads.TimeGetMinutes, overloads.TimestampToMinutesWithTz,
		[]ref.Val{String("America/Phoenix")})
	if !minTz.Equal(Int(5)).(Bool) {
		t.Error("Expected 5 minutes, got", minTz)
	}
}

func TestTimestampGetSeconds(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	sec := ts.Receive(overloads.TimeGetSeconds, overloads.TimestampToSeconds, []ref.Val{})
	if !sec.Equal(Int(6)).(Bool) {
		t.Error("Expected 6 seconds, got", sec)
	}
	// 1969-12-31T19:05:06Z
	secTz := ts.Receive(overloads.TimeGetSeconds, overloads.TimestampToSecondsWithTz,
		[]ref.Val{String("America/Phoenix")})
	if !secTz.Equal(Int(6)).(Bool) {
		t.Error("Expected 6 seconds, got", secTz)
	}
}
