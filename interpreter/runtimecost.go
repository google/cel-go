// Copyright 2022 Google LLC
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

// Package interpreter provides functions to evaluate parsed expressions with
// the option to augment the evaluation with inputs and functions supplied at
// evaluation time.

package interpreter

import (
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types/ref"
	"math"
	"reflect"
)

// ActualCostEstimator provides function call cost estimations at runtime
// CallCost returns an estimated cost for the function overload invocation with the given args, or nil if it has not estimate to provide.
// CEL attempts to provide reasonable estimates for it's standard function library, so CallCost should typically not need to provide an estimate for CELs standard function.
type ActualCostEstimator interface {
	CallCost(overloadId string, args []ref.Val) *uint64
}

// CostTracker represents the information needed for tacking runtime cost
type CostTracker struct {
	estimator         *checker.CostEstimator
	CallCostEstimator ActualCostEstimator
	state             EvalState
	cost              uint64
}

// actualSize returns the size of value
func (c CostTracker) actualSize(i Interpretable) uint64 {
	if v, ok := c.state.Value(i.ID()); ok {
		reflectV := reflect.ValueOf(v.Value())
		switch reflectV.Kind() {
		// Note that the CEL bytes type is implemented with Go byte slices, therefore also supported by the following
		// code.
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			return uint64(reflectV.Len())
		}
	}
	return 1
}

// ActualCost returns the runtime cost
func (c CostTracker) ActualCost() uint64 {
	return c.cost
}

type evalCostTrackerConst struct {
	InterpretableConst
	tracker *CostTracker
}

// Eval implements the Interpretable interface method.
func (e *evalCostTrackerConst) Eval(ctx Activation) ref.Val {
	val := e.InterpretableConst.Eval(ctx)
	return val
}

type evalCostTrackerAttribute struct {
	InterpretableAttribute
	tracker *CostTracker
}

// Eval implements the Interpretable interface method.
func (e *evalCostTrackerAttribute) Eval(ctx Activation) ref.Val {
	val := e.InterpretableAttribute.Eval(ctx)
	e.tracker.cost += e.InterpretableAttribute.RuntimeCost(ctx, e.tracker.state)
	return val
}

type evalCostTrackerCall struct {
	InterpretableCall
	tracker *CostTracker
}

// Eval implements the Interpretable interface method.
func (e *evalCostTrackerCall) Eval(ctx Activation) ref.Val {
	val := e.InterpretableCall.Eval(ctx)
	args := e.InterpretableCall.Args()
	var argValues []ref.Val
	for _, arg := range args {
		if v, ok := e.tracker.state.Value(arg.ID()); ok {
			argValues = append(argValues, v)
		} else {
			// FIXME: abort.
		}
	}
	if e.tracker.CallCostEstimator != nil {
		cost := e.tracker.CallCostEstimator.CallCost(e.OverloadID(), argValues)
		if cost != nil {
			e.tracker.cost += *cost
			return val
		}
	}
	// if user didn't specify, the default way of calculating runtime cost would be used.
	// if user has their own implementation of ActualCostEstimator, make sure to cover the mapping between overloadId and cost calculation
	switch e.OverloadID() {
	// O(n) functions
	case overloads.StartsWithString, overloads.EndsWithString, overloads.StringToBytes, overloads.BytesToString:
		e.tracker.cost += uint64(math.Ceil(float64(e.tracker.actualSize(args[0])) * 0.1))
	case overloads.InList:
		// If a list is composed entirely of constant values this is O(1), but we don't account for that here.
		// We just assume all list containment checks are O(n).
		e.tracker.cost += e.tracker.actualSize(args[1])
	// O(min(m, n)) functions
	case overloads.LessString, overloads.GreaterString, overloads.LessEqualsString, overloads.GreaterEqualsString,
		overloads.LessBytes, overloads.GreaterBytes, overloads.LessEqualsBytes, overloads.GreaterEqualsBytes,
		overloads.Equals, overloads.NotEquals:
		// When we check the equality of 2 scalar values (e.g. 2 integers, 2 floating-point numbers, 2 booleans etc.),
		// the CostTracker.actualSize() function by definition returns 1 for each operand, resulting in an overall cost
		// of 1.
		lhsCost := e.tracker.actualSize(args[0])
		rhsCost := e.tracker.actualSize(args[1])
		if lhsCost > rhsCost {
			e.tracker.cost += rhsCost
		} else {
			e.tracker.cost += lhsCost
		}
	// O(m+n) functions
	case overloads.AddString, overloads.AddBytes, overloads.AddList:
		// In the worst case scenario, we would need to reallocate a new backing store and copy both operands over.
		e.tracker.cost += uint64(math.Ceil(float64(e.tracker.actualSize(args[0])+e.tracker.actualSize(args[1])) * 0.1))
	// O(nm) functions
	case overloads.MatchesString:
		// https://swtch.com/~rsc/regexp/regexp1.html applies to RE2 implementation supported by CEL
		strCost := uint64(math.Ceil(float64(e.tracker.actualSize(args[0])) * 0.1))
		// We don't know how many expressions are in the regex, just the string length (a huge
		// improvement here would be to somehow get a count the number of expressions in the regex or
		// how many states are in the regex state machine and use that to measure regex cost).
		// For now, we're making a guess that each expression in a regex is typically at least 4 chars
		// in length.
		regexCost := uint64(math.Ceil(float64(e.tracker.actualSize(args[1])) * 0.25))
		e.tracker.cost += strCost * regexCost
	case overloads.ContainsString:
		strCost := uint64(math.Ceil(float64(e.tracker.actualSize(args[0])) * 0.1))
		substrCost := uint64(math.Ceil(float64(e.tracker.actualSize(args[1])) * 0.1))
		e.tracker.cost += strCost * substrCost
	default:
		// The following operations are assumed to have O(1) complexity.
		// 1. Concatenation of 2 lists: see the implementation of the concatList type.
		// 2. Computing the size of strings, byte sequences, lists and maps: presumably, the length of each of these
		//    data structures are cached and can be retrieved in constant time.
		e.tracker.cost += 1

	}

	return val
}

type evalCostTrackerOp struct {
	InterpretableOp
	tracker *CostTracker
}

// Eval implements the Interpretable interface method.
func (e *evalCostTrackerOp) Eval(ctx Activation) ref.Val {
	val := e.InterpretableOp.Eval(ctx)
	e.tracker.cost += e.InterpretableOp.RuntimeCost()

	return val
}

func TrackCost(evalState EvalState, tracker *CostTracker) InterpretableDecorator {
	return func(i Interpretable) (Interpretable, error) {
		tracker.state = evalState
		switch t := i.(type) {
		// this is to remove extra interpretable nodes
		case *evalCostTrackerAttribute, *evalCostTrackerCall, *evalCostTrackerConst, *evalCostTrackerOp:
			return i, nil
		case InterpretableConst:
			return &evalCostTrackerConst{InterpretableConst: t, tracker: tracker}, nil
		case InterpretableAttribute:
			return &evalCostTrackerAttribute{InterpretableAttribute: t, tracker: tracker}, nil
		case InterpretableCall:
			return &evalCostTrackerCall{InterpretableCall: t, tracker: tracker}, nil
		case InterpretableOp:
			return &evalCostTrackerOp{InterpretableOp: t, tracker: tracker}, nil
		}
		return i, nil
	}
}

func calRuntimeCost(i interface{}, ctx Activation, evalState EvalState) uint64 {
	c, ok := i.(RuntimeCoster)
	if !ok {
		return math.MaxInt64
	}
	return c.RuntimeCost(ctx, evalState)
}
