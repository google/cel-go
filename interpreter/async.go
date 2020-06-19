// Copyright 2020 Google LLC
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
	"sync"
	"sync/atomic"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

// AsyncInterpretable supports evaluation with async extension functions.
type AsyncInterpretable interface {
	// ID of the interpretable expression.
	ID() int64

	// AsyncEval drives the evaluation of an Interpretable program which may invoke asynchronous
	// calls.
	AsyncEval(context.Context, *AsyncActivation) ref.Val
}

// asyncEval is the default implementation of the AsyncInterpretable interface.
type asyncEval struct {
	Interpretable
}

// AsyncEval implements the AsyncInterpretable interface method and applies the following
// algorithm:
//
// - Evaluate synchronously. If the result is not of types.Unkonwn, return the result.
// - Otherwise for each unknown value determine if the expression id maps to the expression id
//   of an async call.
// - For each async call and unique argument set found, invoke the async implementation.
// - Async call results are memoized with their corresponding argument sets.
// - Revaluate synchronously. Repeat if unknowns are still present and progress (at least one
//   async call invocation has been made).
//
// The process of incremental evaluation may be strictly more expensive than some alternative
// approaches; however, the async call is expected to dominate the execution time by several
// orders of magnitude.
func (async *asyncEval) AsyncEval(ctx context.Context, vars *AsyncActivation) ref.Val {
	res := async.Interpretable.Eval(vars)
	for types.IsUnknown(res) {
		progress := false
		unk := res.(types.Unknown)
		for _, id := range unk {
			call, found := vars.findCall(id)
			if !found {
				continue
			}
			call.evals.Range(func(idx, rawArgs interface{}) bool {
				asyncResult, hasResult := call.results.Load(idx)
				if hasResult {
					return true
				}
				progress = true
				args := rawArgs.([]ref.Val)
				asyncResult = call.impl(ctx, vars, args)
				call.results.Store(idx, asyncResult)
				return false
			})
		}
		if !progress {
			break
		}
		res = async.Interpretable.Eval(vars)
	}
	return res
}

// NewAsyncActivation returns an AsyncActivation capable of tracking async calls relevant to a
// single evaluation pass through an AsyncInterpretable object.
//
func NewAsyncActivation(vars Activation) *AsyncActivation {
	return &AsyncActivation{Activation: vars}
}

// AsyncActivation tracks calls by expression identifier and overload name for use in subsequent
// async function invocations when the result cannot be computed without context provided by the
// remote call.
//
// Note, this object is concurrency safe, but should not be reused across invocations.
type AsyncActivation struct {
	Activation
	asyncCalls sync.Map // map[int64]*asyncCall
}

// findCall returns an async call matching the given expression id, if any.
func (act *AsyncActivation) findCall(id int64) (*asyncCall, bool) {
	call, found := act.asyncCalls.Load(id)
	if found {
		return call.(*asyncCall), true
	}
	return nil, false
}

// findCallbyOverload returns an async call matching the given overload identifer, if any.
//
// Note, the same overload may be referenced in multiple locations within the expression.
// This call will return the first instance of the call evaluated for the purpose of aggregating
// related arugments sets to the underlying function and deduping them.
func (act *AsyncActivation) findCallByOverload(overload string) (*asyncCall, bool) {
	var call *asyncCall
	act.asyncCalls.Range(func(_, ac interface{}) bool {
		asyncCall := ac.(*asyncCall)
		if asyncCall.overload == overload {
			call = asyncCall
			return false
		}
		return true
	})
	if call != nil {
		return call, true
	}
	return nil, false
}

// putCall stores the asyncCall by its expression id.
func (act *AsyncActivation) putCall(id int64, call *asyncCall) {
	act.asyncCalls.Store(id, call)
}

// newAsyncCall creates a new asyncCall instance capable of tracking argument sets to the call.
func newAsyncCall(id int64, overload string, impl functions.AsyncOp) *asyncCall {
	var idx int32
	return &asyncCall{
		id:          id,
		overload:    overload,
		impl:        impl,
		nextCallIdx: &idx,
	}
}

// asyncCall tracks argument sets and associated results for a given async call.
type asyncCall struct {
	id          int64
	impl        functions.AsyncOp
	overload    string
	nextCallIdx *int32
	evals       sync.Map // [][]ref.Val
	results     sync.Map // []ref.Val
}

// Eval implements the Interpretable interface method, tracks argument sets, and returns memoized
// results, if present.
func (ac *asyncCall) Eval(args []ref.Val) ref.Val {
	argsFound := false
	idx := atomic.LoadInt32(ac.nextCallIdx)
	ac.evals.Range(func(i, rawArgs interface{}) bool {
		asyncArgs := rawArgs.([]ref.Val)
		for j, arg := range args {
			if arg.Equal(asyncArgs[j]) != types.True {
				return true
			}
		}
		idx = i.(int32)
		argsFound = true
		return false
	})
	val, resultFound := ac.results.Load(idx)
	// args found and there's a result, return it.
	if resultFound {
		return val.(ref.Val)
	}
	// args found, but no result yet, return unknown.
	if argsFound {
		return types.Unknown{ac.id}
	}
	// Args not tracked, track them here.
	swapped := false
	for !swapped {
		swapped = atomic.CompareAndSwapInt32(ac.nextCallIdx, idx, idx+1)
		idx++
	}
	ac.evals.Store(idx, args)
	return types.Unknown{ac.id}
}

// evalAsyncCall implements the InterpretableCall interface and is the entry point into async
// function evaluation coordination.
type evalAsyncCall struct {
	id       int64
	function string
	overload string
	args     []Interpretable
	impl     functions.AsyncOp
}

// ID returns the expression identifier where the call is referenced.
func (async *evalAsyncCall) ID() int64 {
	return async.id
}

// Function implements the InterpretableCall interface method.
func (async *evalAsyncCall) Function() string {
	return async.function
}

// OverloadID implements the InterpretableCall interface method.
func (async *evalAsyncCall) OverloadID() string {
	return async.overload
}

// Args returns the argument to the unary function.
func (async *evalAsyncCall) Args() []Interpretable {
	return async.args
}

// Eval tracks the arguments provided to the underlying async invocation, deduping argument sets
// if possible.
func (async *evalAsyncCall) Eval(vars Activation) ref.Val {
	argVals := make([]ref.Val, len(async.args), len(async.args))
	// Early return if any argument to the function is unknown or error.
	for i, arg := range async.args {
		argVals[i] = arg.Eval(vars)
		if types.IsUnknownOrError(argVals[i]) {
			return argVals[i]
		}
	}
	// Attempt to return the result if one exists. Note, this can be pretty expensive,
	// but presumably cheaper than actually invoking the async call.
	asyncVars := vars.(*AsyncActivation)
	calls, found := asyncVars.findCallByOverload(async.overload)
	if !found {
		calls = newAsyncCall(async.id, async.overload, async.impl)
	}
	asyncVars.putCall(async.id, calls)
	return calls.Eval(argVals)
}
