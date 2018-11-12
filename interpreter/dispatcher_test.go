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

package interpreter

import (
	"testing"

	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

func TestDefaultDispatcher_Dispatch(t *testing.T) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardOverloads()...); err != nil {
		t.Error(err)
	}
	call := &CallContext{
		call: NewCall(0,
			operators.Equals,
			[]int64{1, 2}),
		args: []ref.Value{types.Int(1), types.Int(2)}}
	invokeCall(t, dispatcher, call, types.False)
}

func TestDefaultDispatcher_DispatchOverload(t *testing.T) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardOverloads()...); err != nil {
		t.Error(err)
	}
	call := &CallContext{
		call: NewCallOverload(0,
			operators.Equals,
			[]int64{1, 2},
			overloads.Equals),
		args: []ref.Value{types.Int(100), types.Int(200)}}
	invokeCall(t, dispatcher, call, types.False)
}

func BenchmarkDefaultDispatcher_Dispatch(b *testing.B) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardOverloads()...); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		call := &CallContext{
			call: NewCall(0,
				operators.Equals,
				[]int64{1, 2}),
			args: []ref.Value{types.Int(1), types.Int(2)}}
		dispatcher.Dispatch(call)
	}
}

func BenchmarkDefaultDispatcher_DispatchOverload(b *testing.B) {
	dispatcher := NewDispatcher()
	if err := dispatcher.Add(functions.StandardOverloads()...); err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		call := &CallContext{
			call: NewCallOverload(0,
				operators.Equals,
				[]int64{1, 2},
				operators.Equals),
			args: []ref.Value{types.Int(2), types.Int(2)}}
		dispatcher.Dispatch(call)
	}
}

func invokeCall(t *testing.T, dispatcher Dispatcher, call *CallContext, expected ref.Value) {
	t.Helper()
	if result := dispatcher.Dispatch(call); types.IsError(result) || types.IsUnknown(result) {
		t.Error(result)
	} else {
		if result != expected {
			t.Errorf(
				"Unexpected result. expected: %v, got: %v in dispatcher: %v",
				expected, result, dispatcher)
		}
	}
}
