package types

import (
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"testing"
)

func TestDuration_Add(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	if !d.Add(d).Equal(Duration{&duration.Duration{Seconds: 15012}}).(Bool) {
		t.Error("Adding duration and itself did not double it.")
	}
}

func TestDuration_Compare(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	lt := Duration{&duration.Duration{Seconds: -10}}
	if d.Compare(lt).(Int) != IntOne {
		t.Error("Larger duration was not considered greater than smaller one.")
	}
	if lt.Compare(d).(Int) != IntNegOne {
		t.Error("Smaller duration was not less than larger one.")
	}
	if d.Compare(d).(Int) != IntZero {
		t.Error("Durations were not considered equal.")
	}
	if !IsError(d.Compare(False)) {
		t.Error("Unexpected comparison supported.")
	}
}

func TestDuration_ReceiveGetHours(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	hr := d.Receive(overloads.TimeGetHours, overloads.DurationToHours, []ref.Value{})
	if !hr.Equal(Int(2)).(Bool) {
		t.Error("Expected 2 hours, got", hr)
	}
}

func TestDuration_ReceiveGetMinutes(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	min := d.Receive(overloads.TimeGetMinutes, overloads.DurationToMinutes, []ref.Value{})
	if !min.Equal(Int(125)).(Bool) {
		t.Error("Expected 5 minutes, got", min)
	}
}

func TestDuration_ReceiveGetSeconds(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	sec := d.Receive(overloads.TimeGetSeconds, overloads.DurationToSeconds, []ref.Value{})
	if !sec.Equal(Int(7506)).(Bool) {
		t.Error("Expected 6 seconds, got", sec)
	}
}

func TestDuration_Subtract(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	if !d.Subtract(d).ConvertToType(IntType).Equal(IntZero).(Bool) {
		t.Error("Subtracting a duration from itself did not equal zero.")
	}
}
