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
	"math"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// ActualCostEstimator provides function call cost estimations at runtime
// CallCost returns an estimated cost for the function overload invocation with the given args, or nil if it has no
// estimate to provide. CEL attempts to provide reasonable estimates for its standard function library, so CallCost
// should typically not need to provide an estimate for CELs standard function.
type ActualCostEstimator interface {
	CallCost(overloadId string, args []ref.Val) *uint64
}

// CostObserver provides an observer that tracks runtime cost.
func CostObserver(tracker *CostTracker) EvalObserver {
	observer := func(id int64, programStep interface{}, val ref.Val) {
		switch t := programStep.(type) {
		case InterpretableConst:
			// zero cost
		case ConstantQualifier:
			tracker.cost += 1
		case InterpretableAttribute:
			tracker.cost += 1
		case Qualifier:
			tracker.stack.pop()
			tracker.cost += 1
		case InterpretableCall:
			argVals := make([]ref.Val, len(t.Args()))
			argsFound := true
			for i := len(t.Args()) - 1; i >= 0; i-- {
				if v, ok := tracker.stack.pop(); ok {
					argVals[i] = v
				} else {
					// should never happen
					argsFound = false
				}
			}
			if argsFound {
				tracker.cost += tracker.costCall(t, argVals)
			}
		case InterpretableConstructor:
			switch t.Type() {
			case types.ListType:
				tracker.cost += 10
			case types.MapType:
				tracker.cost += 30
			default:
				tracker.cost += 40
			}
		case InterpretableBooleanBinaryOp:
			tracker.cost += 1
		}
		tracker.stack.push(val)
		if tracker.Limit != nil && tracker.cost > *tracker.Limit {
			panic(EvalCanceledError{Cause: ActualCostLimitExceeded})
		}
	}
	return observer
}

// CostTracker represents the information needed for tacking runtime cost
type CostTracker struct {
	Estimator ActualCostEstimator
	Limit     *uint64

	cost  uint64
	stack refValStack
}

// ActualCost returns the runtime cost
func (c CostTracker) ActualCost() uint64 {
	return c.cost
}

func (c CostTracker) costCall(call InterpretableCall, argValues []ref.Val) uint64 {
	var cost uint64
	if c.Estimator != nil {
		callCost := c.Estimator.CallCost(call.OverloadID(), argValues)
		if callCost != nil {
			cost += *callCost
			return cost
		}
	}
	// if user didn't specify, the default way of calculating runtime cost would be used.
	// if user has their own implementation of ActualCostEstimator, make sure to cover the mapping between overloadId and cost calculation
	switch call.OverloadID() {
	// O(n) functions
	case overloads.StartsWithString, overloads.EndsWithString, overloads.StringToBytes, overloads.BytesToString:
		cost += uint64(math.Ceil(float64(c.actualSize(argValues[0])) * 0.1))
	case overloads.InList:
		// If a list is composed entirely of constant values this is O(1), but we don't account for that here.
		// We just assume all list containment checks are O(n).
		cost += c.actualSize(argValues[1])
	// O(min(m, n)) functions
	case overloads.LessString, overloads.GreaterString, overloads.LessEqualsString, overloads.GreaterEqualsString,
		overloads.LessBytes, overloads.GreaterBytes, overloads.LessEqualsBytes, overloads.GreaterEqualsBytes,
		overloads.Equals, overloads.NotEquals:
		// When we check the equality of 2 scalar values (e.g. 2 integers, 2 floating-point numbers, 2 booleans etc.),
		// the CostTracker.actualSize() function by definition returns 1 for each operand, resulting in an overall cost
		// of 1.
		lhsCost := c.actualSize(argValues[0])
		rhsCost := c.actualSize(argValues[1])
		if lhsCost > rhsCost {
			cost += rhsCost
		} else {
			cost += lhsCost
		}
	// O(m+n) functions
	case overloads.AddString, overloads.AddBytes, overloads.AddList:
		// In the worst case scenario, we would need to reallocate a new backing store and copy both operands over.
		cost += uint64(math.Ceil(float64(c.actualSize(argValues[0])+c.actualSize(argValues[1])) * 0.1))
	// O(nm) functions
	case overloads.MatchesString:
		// https://swtch.com/~rsc/regexp/regexp1.html applies to RE2 implementation supported by CEL
		strCost := uint64(math.Ceil(float64(c.actualSize(argValues[0])) * 0.1))
		// We don't know how many expressions are in the regex, just the string length (a huge
		// improvement here would be to somehow get a count the number of expressions in the regex or
		// how many states are in the regex state machine and use that to measure regex cost).
		// For now, we're making a guess that each expression in a regex is typically at least 4 chars
		// in length.
		regexCost := uint64(math.Ceil(float64(c.actualSize(argValues[1])) * 0.25))
		cost += strCost * regexCost
	case overloads.ContainsString:
		strCost := uint64(math.Ceil(float64(c.actualSize(argValues[0])) * 0.1))
		substrCost := uint64(math.Ceil(float64(c.actualSize(argValues[1])) * 0.1))
		cost += strCost * substrCost
	default:
		// The following operations are assumed to have O(1) complexity.
		// 1. Concatenation of 2 lists: see the implementation of the concatList type.
		// 2. Computing the size of strings, byte sequences, lists and maps: presumably, the length of each of these
		//    data structures are cached and can be retrieved in constant time.
		cost += 1

	}
	return cost
}

// actualSize returns the size of value
func (c CostTracker) actualSize(value ref.Val) uint64 {
	if sz, ok := value.(traits.Sizer); ok {
		return uint64(sz.Size().Value().(int64))
	}
	return 1
}

// refValStack keeps track of values of the stack for cost calculation purposes
type refValStack []ref.Val

func (s *refValStack) push(value ref.Val) {
	*s = append(*s, value)
}

func (s *refValStack) pop() (ref.Val, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	idx := len(*s) - 1
	el := (*s)[idx]
	*s = (*s)[:idx]
	return el, true
}
