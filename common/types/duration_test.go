package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"reflect"
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

func TestDuration_ConvertToNative(t *testing.T) {
	val, err := Duration{&duration.Duration{Seconds: 7506, Nanos: 1000}}.
		ConvertToNative(reflect.TypeOf(&duration.Duration{}))
	if err != nil ||
		!proto.Equal(val.(proto.Message), &duration.Duration{Seconds: 7506, Nanos: 1000}) {
		t.Errorf("Got '%v', expected backing proto message value", err)
	}
}

func TestDuration_ConvertToNative_Error(t *testing.T) {
	val, err := Duration{&duration.Duration{Seconds: 7506, Nanos: 1000}}.
		ConvertToNative(jsonValueType)
	if err == nil {
		t.Errorf("Got '%v', expected error", val)
	}
}

func TestDuration_ConvertToType_Identity(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506, Nanos: 1000}}
	str := d.ConvertToType(StringType).(String)
	if str != "2h5m6.000001s" {
		t.Errorf("Got '%v', wanted 2h5m6.000001s", str)
	}
	i := d.ConvertToType(IntType).(Int)
	if i != Int(7506000001000) {
		t.Errorf("Got '%v', wanted 7506000001000", i)
	}
	if !d.ConvertToType(DurationType).Equal(d).(Bool) {
		t.Errorf("Got '%v', wanted identity", d.ConvertToType(DurationType))
	}
	if d.ConvertToType(TypeType) != DurationType {
		t.Errorf("Got '%v', expected duration type", d.ConvertToType(TypeType))
	}
	if !IsError(d.ConvertToType(UintType)) {
		t.Errorf("Expected error on conversion of duration to uint")
	}
}

func TestDuration_Negate(t *testing.T) {
	neg := Duration{&duration.Duration{Seconds: 1234, Nanos: 1}}.Negate().(Duration)
	if !proto.Equal(neg.Duration, &duration.Duration{Seconds: -1234, Nanos: -1}) {
		t.Errorf("Got '%v', expected seconds: -1234, nanos: -1", neg)
	}
}

func TestDuration_Receive_GetHours(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	hr := d.Receive(overloads.TimeGetHours, overloads.DurationToHours, []ref.Value{})
	if !hr.Equal(Int(2)).(Bool) {
		t.Error("Expected 2 hours, got", hr)
	}
}

func TestDuration_Receive_GetMinutes(t *testing.T) {
	d := Duration{&duration.Duration{Seconds: 7506}}
	min := d.Receive(overloads.TimeGetMinutes, overloads.DurationToMinutes, []ref.Value{})
	if !min.Equal(Int(125)).(Bool) {
		t.Error("Expected 5 minutes, got", min)
	}
}

func TestDuration_Receive_GetSeconds(t *testing.T) {
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
