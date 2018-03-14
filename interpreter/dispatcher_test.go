package interpreter

import (
	"github.com/google/cel-go/interpreter/functions"
	"github.com/grafeas/grafeas/server-go/filtering/operators"
	"testing"
)

func TestDefaultDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardBuiltins()...); err != nil {
		t.Error(err)
	}
	call := &CallContext{
		call: NewCall(0,
			operators.Equals,
			[]int64{1, 2}),
		args: []interface{}{int64(1), int64(2)}}
	invokeCall(t, dispatcher, call, false)
}

func TestDefaultDispatcher_DispatchOverload(t *testing.T) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardBuiltins()...); err != nil {
		t.Error(err)
	}
	call := &CallContext{
		call: NewCallOverload(0,
			operators.Equals,
			[]int64{1, 2},
			operators.Equals),
		args: []interface{}{int64(100), int64(200)}}
	invokeCall(t, dispatcher, call, false)
}

func BenchmarkDefaultDispatcher_Dispatch(b *testing.B) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardBuiltins()...); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		call := &CallContext{
			call: NewCall(0,
				operators.Equals,
				[]int64{1, 2}),
			args: []interface{}{int64(1), int64(2)}}
		dispatcher.Dispatch(call)
	}
}

func BenchmarkDefaultDispatcher_DispatchOverload(b *testing.B) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardBuiltins()...); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		call := &CallContext{
			call: NewCallOverload(0,
				operators.Equals,
				[]int64{1, 2},
				operators.Equals),
			args: []interface{}{int64(2), int64(2)}}
		dispatcher.Dispatch(call)
	}
}

func invokeCall(t *testing.T, dispatcher Dispatcher, call *CallContext, expected interface{}) {
	t.Helper()
	if result, err := dispatcher.Dispatch(call); err == nil {
		if result != expected {
			t.Errorf(
				"Unexpected result. expected: %v, got: %v in dispatcher: %v",
				expected, result, dispatcher)
		}
	} else {
		t.Error(err)
	}
}
