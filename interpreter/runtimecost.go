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
)

type CostTracker struct {
	estimator *checker.CostEstimator
	state     EvalState
	cost      uint64 // TODO: does this need to be an atomic?
}

func (c CostTracker) actualSize(i Interpretable, ctx Activation) uint64 {
	if v, ok := c.state.Value(i.ID()); ok {
		switch o := v.Value().(type) {
		case string:
			return uint64(len(o))
		}
	}
	return 0
}

func (c CostTracker) ActualCost() uint64 {
	return c.cost
}

type evalCostTrackerConst struct {
	InterpretableConst
	tracker *CostTracker
}

func (e *evalCostTrackerConst) Eval(ctx Activation) ref.Val {
	e.tracker.cost += 1
	val := e.InterpretableConst.Eval(ctx)
	return val
}

type evalCostTrackerAttribute struct {
	InterpretableAttribute
	tracker *CostTracker
}

func (e *evalCostTrackerAttribute) Eval(ctx Activation) ref.Val {
	// TODO: adjust the cost based on e.InterpretableAttribute.Attr()
	e.tracker.cost += 1
	val := e.InterpretableAttribute.Eval(ctx)
	return val
}

type evalCostTrackerCall struct {
	InterpretableCall
	tracker *CostTracker
}

func (e *evalCostTrackerCall) Eval(ctx Activation) ref.Val {
	val := e.InterpretableCall.Eval(ctx)
	args := e.InterpretableCall.Args()
	switch e.OverloadID() {
	// O(n) functions
	case overloads.StartsWithString, overloads.EndsWithString, overloads.StringToBytes, overloads.BytesToString:
		e.tracker.cost += e.tracker.actualSize(args[0], ctx) // TODO: multiply by 0.1
	case overloads.InList:
		// If a list is composed entirely of constant values this is O(1), but we don't account for that here.
		// We just assume all list containment checks are O(n).
		e.tracker.cost += e.tracker.actualSize(args[1], ctx)
	// O(nm) functions
	case overloads.MatchesString:
		// https://swtch.com/~rsc/regexp/regexp1.html applies to RE2 implementation supported by CEL
		strCost := e.tracker.actualSize(args[0], ctx) // TODO: multiply by 0.1
		// We don't know how many expressions are in the regex, just the string length (a huge
		// improvement here would be to somehow get a count the number of expressions in the regex or
		// how many states are in the regex state machine and use that to measure regex cost).
		// For now, we're making a guess that each expression in a regex is typically at least 4 chars
		// in length.
		regexCost := e.tracker.actualSize(args[1], ctx) // TODO: multiply by 0.25
		e.tracker.cost += strCost * regexCost
	case overloads.ContainsString:
		strCost := e.tracker.actualSize(args[0], ctx)    // TODO: multiply by 0.1
		substrCost := e.tracker.actualSize(args[1], ctx) // TODO: multiply by 0.1
		e.tracker.cost += strCost * substrCost
	}

	// TODO: add cost based on overload id of function
	return val
}

func TrackCost(evalState EvalState, tracker *CostTracker) InterpretableDecorator {
	return func(i Interpretable) (Interpretable, error) {
		tracker.state = evalState
		switch t := i.(type) {
		case InterpretableConst:
			return &evalCostTrackerConst{InterpretableConst: t, tracker: tracker}, nil
		case InterpretableAttribute:
			return &evalCostTrackerAttribute{InterpretableAttribute: t, tracker: tracker}, nil
		case InterpretableCall:
			return &evalCostTrackerCall{InterpretableCall: t, tracker: tracker}, nil
		}
		return i, nil
	}
}
