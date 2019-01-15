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
	"github.com/google/cel-go/interpreter/functions"
)

func TestDefaultDispatcher_Dispatch(t *testing.T) {
	state := NewEvalState(3)
	state.SetValue(1, types.Int(2))
	state.SetValue(2, types.Int(2))
	call := NewCall(0, operators.Equals, []int64{1, 2})
	disp(state).Dispatch(call)
	res, _ := state.Value(0)
	if res != types.True {
		t.Errorf("Got '%v', wanted 'true'", res)
	}
}

func TestDefaultDispatcher_DispatchOverload(t *testing.T) {
	state := NewEvalState(3)
	state.SetValue(1, types.Int(100))
	state.SetValue(2, types.Int(200))
	call := NewCallOverload(0,
		operators.Equals,
		[]int64{1, 2},
		overloads.Equals)
	disp(state).Dispatch(call)
	res, _ := state.Value(0)
	if res != types.False {
		t.Errorf("Got '%v', wanted 'false'", res)
	}
}

func BenchmarkDefaultDispatcher_Dispatch(b *testing.B) {
	call := NewCall(0, operators.NotEquals, []int64{1, 2})
	state := NewEvalState(3)
	state.SetValue(1, types.Int(1))
	state.SetValue(2, types.Int(2))
	d := disp(state)
	for i := 0; i < b.N; i++ {
		d.Dispatch(call)
	}
}

func BenchmarkDefaultDispatcher_DispatchOverload(b *testing.B) {
	call := NewCallOverload(0,
		operators.NotEquals,
		[]int64{1, 2},
		overloads.NotEquals)
	state := NewEvalState(3)
	state.SetValue(1, types.Int(2))
	state.SetValue(2, types.Int(2))
	d := disp(state)
	for i := 0; i < b.N; i++ {
		d.Dispatch(call)
	}
}

func disp(state MutableEvalState) Dispatcher {
	dispatcher := NewDispatcher()
	dispatcher.Add(functions.StandardOverloads()...)
	return dispatcher.Init(state)
}
