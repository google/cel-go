package types

import (
	"testing"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"github.com/golang/protobuf/ptypes/duration"
)

func TestDuration_ReceiveGetHours(t *testing.T) {
	d := Duration{&duration.Duration{Seconds:7506}}
	hr := d.Receive(overloads.TimeGetHours, overloads.DurationToHours, []ref.Value{})
	if !hr.Equal(Int(2)).(Bool) {
		t.Error("Expected 2 hours, got", hr)
	}
}

func TestDuration_ReceiveGetMinutes(t *testing.T) {
	d := Duration{&duration.Duration{Seconds:7506}}
	min := d.Receive(overloads.TimeGetMinutes, overloads.DurationToMinutes, []ref.Value{})
	if !min.Equal(Int(125)).(Bool) {
		t.Error("Expected 5 minutes, got", min)
	}
}

func TestDuration_ReceiveGetSeconds(t *testing.T) {
	d := Duration{&duration.Duration{Seconds:7506}}
	sec := d.Receive(overloads.TimeGetSeconds, overloads.DurationToSeconds, []ref.Value{})
	if !sec.Equal(Int(7506)).(Bool) {
		t.Error("Expected 6 seconds, got", sec)
	}
}