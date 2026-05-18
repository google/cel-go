// Copyright 2026 Google LLC
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
	"context"
	"testing"
	"time"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func TestFrameInterrupt(t *testing.T) {
	tests := []struct {
		name        string
		concurrency int
	}{
		{
			name:        "without semaphores",
			concurrency: 0,
		},
		{
			name:        "with semaphores",
			concurrency: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			frame := NewExecutionFrame(EmptyActivation())
			defer frame.Close()

			// Test CheckInterrupt with nil ctx
			if frame.CheckInterrupt() {
				t.Error("CheckInterrupt() returned true for nil ctx")
			}

			ctx, cancel := context.WithCancel(context.Background())
			frame.SetContext(ctx, 1)
			frame.SetAsyncMaxConcurrency(tc.concurrency)

			// First check with active context should return false
			if frame.CheckInterrupt() {
				t.Error("frame.CheckInterrupt() returned true, wanted false")
			}

			// Cancel context to trigger interrupt
			cancel()

			// Second check should observe cancellation and return true
			if !frame.CheckInterrupt() {
				t.Error("frame.CheckInterrupt() returned false, wanted true")
			}

			// Third check should hit the cached/loaded interrupted flag and return true
			if !frame.CheckInterrupt() {
				t.Error("frame.CheckInterrupt() (cached) returned false, wanted true")
			}
		})
	}

	t.Run("check frequency not triggered", func(t *testing.T) {
		frame := NewExecutionFrame(EmptyActivation())
		defer frame.Close()

		frame.SetContext(context.Background(), 2)
		// first call adds 1 to count -> 1%2 != 0 -> returns false
		if frame.CheckInterrupt() {
			t.Error("CheckInterrupt() returned true unexpectedly")
		}
	})
}

func TestFrameAsyncErrors(t *testing.T) {
	frame := NewExecutionFrame(EmptyActivation())
	defer frame.Close()

	// ComputeResult without SetContext should return an error
	res := frame.ComputeResult(1, "test", "test", nil, nil)
	if !types.IsError(res) {
		t.Errorf("ComputeResult() returned %v, wanted error", res)
	}
}

func TestFrameCompletions(t *testing.T) {
	frame := NewExecutionFrame(EmptyActivation())
	defer frame.Close()

	ctx := context.Background()
	frame.SetContext(ctx, 1)

	completionCh := make(chan int64, 1)
	frame.SetCompletions(completionCh)

	syncCh := make(chan ref.Val, 1)
	impl := func(ctx context.Context, args ...ref.Val) <-chan ref.Val {
		return syncCh
	}

	// First call returns unknown
	res := frame.ComputeResult(10, "test", "test", impl, nil)
	if !types.IsUnknown(res) {
		t.Fatalf("ComputeResult() returned %v, wanted unknown", res)
	}

	// Send result
	syncCh <- types.String("done")

	// Wait for completion signal
	select {
	case id := <-completionCh:
		if id != res.(*types.Unknown).IDs()[0] {
			t.Errorf("completion signal ID mismatch: got %d, wanted %d", id, res.(*types.Unknown).IDs()[0])
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for completion signal")
	}

	// Second call should return the result
	res2 := frame.ComputeResult(10, "test", "test", impl, nil)
	if res2.Equal(types.String("done")) != types.True {
		t.Errorf("ComputeResult() returned %v, wanted 'done'", res2)
	}
}

func TestEvalAsyncFunc(t *testing.T) {
	frame := NewExecutionFrame(EmptyActivation())
	defer frame.Close()
	frame.SetContext(context.Background(), 1)

	impl := func(ctx context.Context, args ...ref.Val) <-chan ref.Val {
		ch := make(chan ref.Val, 1)
		ch <- types.Int(args[0].(types.Int) + 1)
		return ch
	}

	fn := &evalAsyncFunc{
		id:       1,
		function: "inc",
		overload: "inc_int",
		args:     []InterpretableV2{NewConstValue(1, types.Int(10))},
		impl:     impl,
	}

	// First exec triggers the call
	res := fn.Exec(frame)
	if !types.IsUnknown(res) {
		t.Fatalf("Exec() returned %v, wanted unknown", res)
	}

	// Test Args()
	if len(fn.Args()) != 1 {
		t.Errorf("Args() returned %d, wanted 1", len(fn.Args()))
	}

	// Wait for async result (it's buffered in the impl channel in this test)
	time.Sleep(10 * time.Millisecond)

	// Test Eval() returns result
	res3 := fn.Eval(frame)
	if res3.Equal(types.Int(11)) != types.True {
		t.Errorf("Eval() returned %v, wanted 11", res3)
	}
}

func TestEvalAsyncFuncShortCircuit(t *testing.T) {
	frame := NewExecutionFrame(EmptyActivation())
	defer frame.Close()
	frame.SetContext(context.Background(), 1)

	fn := &evalAsyncFunc{
		id:       1,
		function: "inc",
		args:     []InterpretableV2{NewConstValue(2, types.NewErr("bad arg"))},
	}

	res := fn.Exec(frame)
	if !types.IsError(res) {
		t.Errorf("Exec() returned %v, wanted error", res)
	}

	fn.args = []InterpretableV2{NewConstValue(3, types.NewUnknown(4, nil))}
	res = fn.Exec(frame)
	if !types.IsUnknown(res) {
		t.Errorf("Exec() returned %v, wanted unknown", res)
	}
}

type testAsyncObserver struct {
	startedCalls  int
	finishedCalls int
}

func (o *testAsyncObserver) OnCallStarted(callID int64, function, overload string, args []ref.Val) {
	o.startedCalls++
}

func (o *testAsyncObserver) OnCallFinished(callID int64, function, overload string, res ref.Val) {
	o.finishedCalls++
}

func TestFrameActivationMethods(t *testing.T) {
	baseAct, _ := NewActivation(map[string]any{"x": 1})
	frame := NewExecutionFrame(baseAct)
	defer frame.Close()

	childAct, _ := NewActivation(map[string]any{"y": 2})
	childFrame := frame.push(childAct)

	// Test ResolveName on childFrame
	if v, ok := childFrame.ResolveName("y"); !ok || v != 2 {
		t.Errorf("ResolveName('y') got %v, %v", v, ok)
	}
	if v, ok := childFrame.ResolveName("x"); !ok || v != 1 {
		t.Errorf("ResolveName('x') got %v, %v", v, ok)
	}
	if _, ok := childFrame.ResolveName("missing"); ok {
		t.Errorf("ResolveName('missing') unexpectedly found")
	}

	// Test Parent(), Unwrap(), AsPartialActivation() on childFrame
	if childFrame.Parent() == nil {
		t.Errorf("Parent() returned nil")
	}
	if childFrame.Unwrap() == nil {
		t.Errorf("Unwrap() returned nil")
	}
	if _, ok := childFrame.AsPartialActivation(); ok {
		t.Errorf("AsPartialActivation() unexpectedly true")
	}

	// Pop back to parent
	popped := childFrame.pop()
	if popped != frame {
		t.Errorf("pop() did not return parent frame")
	}

	// Test AsPartialActivation returning true when frame wraps a PartialActivation
	partAct, _ := NewPartialActivation(map[string]any{"z": 3})
	partFrame := NewExecutionFrame(partAct)
	defer partFrame.Close()
	if _, ok := partFrame.AsPartialActivation(); !ok {
		t.Errorf("AsPartialActivation() returned false for PartialActivation")
	}
}

func TestFrameAsyncCallAndObserver(t *testing.T) {
	frame := NewExecutionFrame(EmptyActivation())
	defer frame.Close()

	// Calling PendingAsyncCalls, AsyncCall, NotifyCompletion with nil ctx/asyncCalls
	if count := frame.PendingAsyncCalls(); count != 0 {
		t.Errorf("PendingAsyncCalls() = %d, wanted 0", count)
	}
	if call := frame.AsyncCall(1); call != nil {
		t.Errorf("AsyncCall() = %v, wanted nil", call)
	}
	frame.NotifyCompletion(1) // should be no-op

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	frame.SetContext(ctx, 1)

	obs := &testAsyncObserver{}
	frame.SetAsyncObserver(obs)
	frame.SetAsyncMaxConcurrency(1)

	syncCh := make(chan ref.Val, 1)
	impl := func(ctx context.Context, args ...ref.Val) <-chan ref.Val {
		return syncCh
	}

	argVals := []ref.Val{types.String("arg")}
	res := frame.ComputeResult(100, "fn", "overload", impl, argVals)
	if !types.IsUnknown(res) {
		t.Fatalf("ComputeResult() = %v, wanted unknown", res)
	}

	// Call ComputeResult again with identical parameters to hit cached branch in getOrCreate
	resCached := frame.ComputeResult(100, "fn", "overload", impl, argVals)
	if !types.IsUnknown(resCached) {
		t.Fatalf("ComputeResult() cached = %v, wanted unknown", resCached)
	}

	if count := frame.PendingAsyncCalls(); count != 1 {
		t.Errorf("PendingAsyncCalls() = %d, wanted 1", count)
	}

	callID := res.(*types.Unknown).IDs()[0]
	asyncCall := frame.AsyncCall(callID)
	if asyncCall == nil {
		t.Fatalf("AsyncCall() returned nil")
	}
	if asyncCall.CallID() != callID {
		t.Errorf("CallID() = %d, wanted %d", asyncCall.CallID(), callID)
	}
	if asyncCall.Function() != "fn" {
		t.Errorf("Function() = %s, wanted fn", asyncCall.Function())
	}
	if asyncCall.Overload() != "overload" {
		t.Errorf("Overload() = %s, wanted overload", asyncCall.Overload())
	}

	// Send result to unblock the original async call
	syncCh <- types.String("success")

	// Wait for observer to see finished call
	time.Sleep(20 * time.Millisecond)
	if obs.startedCalls == 0 {
		t.Errorf("observer OnCallStarted not called")
	}
	if obs.finishedCalls == 0 {
		t.Errorf("observer OnCallFinished not called")
	}

	// Test SetAsyncMaxConcurrency(0) resets semaphore
	frame.SetAsyncMaxConcurrency(0)
	if frame.ctx.semaphore != nil {
		t.Errorf("semaphore not set to nil")
	}

	// Test explicit NotifyCompletion
	frame.NotifyCompletion(callID)
}

func TestAsyncCallStateEquals(t *testing.T) {
	baseAcs := &asyncCallState{
		function: "fn",
		overload: "ov",
		argVals:  []ref.Val{types.Int(1)},
	}
	var nilAcs *asyncCallState

	tests := []struct {
		name     string
		receiver *asyncCallState
		other    *asyncCallState
		want     bool
	}{
		{
			name:     "receiver nil",
			receiver: nilAcs,
			other:    baseAcs,
			want:     false,
		},
		{
			name:     "other nil",
			receiver: baseAcs,
			other:    nil,
			want:     false,
		},
		{
			name:     "different function",
			receiver: baseAcs,
			other:    &asyncCallState{function: "other", overload: "ov", argVals: []ref.Val{types.Int(1)}},
			want:     false,
		},
		{
			name:     "different overload",
			receiver: baseAcs,
			other:    &asyncCallState{function: "fn", overload: "other", argVals: []ref.Val{types.Int(1)}},
			want:     false,
		},
		{
			name:     "different arg vals",
			receiver: baseAcs,
			other:    &asyncCallState{function: "fn", overload: "ov", argVals: []ref.Val{types.Int(2)}},
			want:     false,
		},
		{
			name:     "identical",
			receiver: baseAcs,
			other:    &asyncCallState{function: "fn", overload: "ov", argVals: []ref.Val{types.Int(1)}},
			want:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.receiver.equals(tc.other); got != tc.want {
				t.Errorf("%v.equals(%v) = %v, wanted %v", tc.receiver, tc.other, got, tc.want)
			}
		})
	}
}

func TestFrameAsyncCallContextCancelled(t *testing.T) {
	frame := NewExecutionFrame(EmptyActivation())
	defer frame.Close()
	ctx, cancel := context.WithCancel(context.Background())
	frame.SetContext(ctx, 1)
	frame.SetAsyncMaxConcurrency(1)

	implSlow := func(ctx context.Context, args ...ref.Val) <-chan ref.Val {
		return make(chan ref.Val) // blocks forever
	}
	// First call acquires the semaphore
	frame.ComputeResult(200, "fn_slow", "overload", implSlow, nil)

	// Give the first call's goroutine time to acquire the semaphore
	time.Sleep(5 * time.Millisecond)
	if len(frame.ctx.semaphore) != 1 {
		t.Errorf("semaphore length = %d, wanted 1", len(frame.ctx.semaphore))
	}

	// Second call tries to acquire semaphore but it's full, so it blocks in select
	frame.ComputeResult(201, "fn_cancel", "overload", implSlow, nil)

	// Cancel the context so the second call's goroutine hits case <-ctx.Done(): when acquiring semaphore,
	// and the first call's goroutine hits case <-ctx.Done(): when waiting on <-ch.
	cancel()

	// Wait a bit to let goroutines process cancellation
	time.Sleep(10 * time.Millisecond)

	// Verify cancellation successfully released resources and interrupted evaluation
	if len(frame.ctx.semaphore) != 0 {
		t.Errorf("semaphore length after cancel = %d, wanted 0", len(frame.ctx.semaphore))
	}
	if !frame.CheckInterrupt() {
		t.Errorf("CheckInterrupt() returned false after cancellation")
	}

	// Call release with nil directly to cover tracker release nil check
	asyncCallStateTrackerPool.release(nil)
}
