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
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestamp_Add(t *testing.T) {
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

func TestTimestamp_ConvertToNative_Any(t *testing.T) {
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

func TestTimestamp_ConvertToNative_Json(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	val, err := ts.ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	}
	want := &structpb.Value{
		Kind: &structpb.Value_StringValue{
			StringValue: "1970-01-01T02:05:06Z",
		},
	}
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
}

func TestTimestamp_ConvertToNative_Timestamp(t *testing.T) {
	// 1970-01-01T02:05:06Z
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	val, err := ts.ConvertToNative(timestampValueType)
	if err != nil {
		t.Error(err)
	}
	want := tpb.New(ts.Time)
	if !proto.Equal(val.(proto.Message), want) {
		t.Errorf("Got '%v', expected '%v'", val, want)
	}
}

func TestTimestamp_Subtract(t *testing.T) {
	ts := Timestamp{Time: time.Unix(7506, 0).UTC()}
	val := ts.Subtract(Duration{Duration: time.Duration(3600)*time.Second + time.Duration(1000)})
	if val.ConvertToType(TypeType) != TimestampType {
		t.Error("Could not add duration and timestamp")
	}
	expected := Timestamp{Time: time.Unix(3905, 999999000).UTC()}
	if !expected.Compare(val).Equal(IntZero).(Bool) {
		t.Errorf("Got '%v', expected '%v'", val, expected)
	}
}

func TestTimestamp_RecieveGetDayOfYear(t *testing.T) {
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
}

func TestTimestamp_ReceiveGetMonth(t *testing.T) {
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

func TestTimestamp_ReceiveGetHours(t *testing.T) {
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

func TestTimestamp_ReceiveGetMinutes(t *testing.T) {
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

func TestTimestamp_ReceiveGetSeconds(t *testing.T) {
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
