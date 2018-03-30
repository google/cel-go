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

// The interpreter package provides functions to evaluate CEL programs against
// a series of inputs and functions supplied at runtime.
package interpreter

import (
	"fmt"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/interpreter/types"
	"reflect"
)

// Dispatcher resolves function calls to their appropriate overload.
type Dispatcher interface {

	// Add one or more overloads, returning an error if any Overload has the
	// same Overload#Name.
	Add(overloads ...*functions.Overload) error

	// Dispatch a call to its appropriate Overload and return the result or
	// error.
	Dispatch(ctx *CallContext) (interface{}, error)
}

// CallContext description of a function call invocation.
type CallContext struct {
	call       *CallExpr
	args       []interface{}
	activation Activation
	metadata   Metadata
}

// Function name to be invoked as it is written in an expression.
func (ctx *CallContext) Function() string {
	return ctx.call.Function
}

// Overload name to be invoked, if set.
func (ctx *CallContext) Overload() (string, bool) {
	return ctx.call.Overload, ctx.call.Overload != ""
}

// Args to provide on the overload dispatch.
func (ctx *CallContext) Args() []interface{} {
	return ctx.args
}

// NewDispatcher returns an empty Dispatcher.
//
// Typically this call would be used with functions#StandardBuiltins:
//
//     dispatcher := NewDispatcher()
//     dispatcher.Add(functions.StandardBuiltins())
func NewDispatcher() *defaultDispatcher {
	return &defaultDispatcher{
		overloads: make(map[string]*overload),
		functions: make(map[string]map[int][]*overload)}
}

// Helper types for tracking overloads by various dimensions.
type overloadMap map[string]*overload
type overloadMapByFunctionAndArgCount map[string]map[int][]*overload

type defaultDispatcher struct {
	overloads overloadMap
	functions overloadMapByFunctionAndArgCount
}

var _ Dispatcher = &defaultDispatcher{}

func (d *defaultDispatcher) Add(overloads ...*functions.Overload) error {
	for _, o := range overloads {
		// Determine the arg count and type from the overload signature.
		iface := reflect.TypeOf(o.Signature)
		argCount := iface.NumIn()
		argTypes := make([]types.Type, argCount)
		for i := 0; i < argCount; i++ {
			refType := iface.In(i)
			var argType types.Type = nil
			// If the arg type is an interface{}, this is equivalent to dyn.
			if refType.Kind() == reflect.Interface {
				argType = types.DynType
				// Otherwise, instantiate a zero-value of the argument and inspect
				// the value type.
			} else {
				refVal := reflect.New(refType).Elem()
				if argTypeVal, found := types.TypeOf(refVal.Interface()); found {
					argType = argTypeVal
				}
			}
			// Set the arg type for the argument at the ordinal 'i' to the
			// type that was resolved, otherwise, return an error.
			if argType != nil {
				argTypes[i] = argType
			} else {
				return fmt.Errorf("unrecognized type '%T'"+
					" in function signature for overload '%s'",
					refType, o.Name)
			}
		}
		overloadRef := &overload{o.Function,
			o.Name,
			argCount,
			argTypes,
			o.Impl}
		// Add the overload unless an overload of the same name has already
		// been provided before.
		if _, found := d.overloads[o.Name]; found {
			return fmt.Errorf("overload already exists '%s'", o.Name)
		}
		// Index the overload by function and by arg count.
		// TODO: consider indexing by abstract type.
		d.overloads[o.Name] = overloadRef
		if byFunction, found := d.functions[o.Function]; !found {
			byFunction = make(map[int][]*overload)
			byFunction[argCount] = []*overload{overloadRef}
			d.functions[o.Function] = byFunction
		} else if byArgCount, found := byFunction[argCount]; !found {
			byFunction[argCount] = []*overload{overloadRef}
		} else {
			byFunction[argCount] = append(byArgCount, overloadRef)
		}
	}
	return nil
}

func (d *defaultDispatcher) Dispatch(ctx *CallContext) (interface{}, error) {
	if overload, err := d.findOverload(ctx); err == nil {
		return overload.impl(ctx.args...)
	} else {
		return nil, err
	}
}

// Find the overload that matches the call context. There should be only one,
// otherwise the call is ambiguous and a runtime error produced.
func (d *defaultDispatcher) findOverload(ctx *CallContext) (*overload, error) {
	// TODO: Add location metadata to overloads.
	if overloadId, found := ctx.Overload(); found {
		if match, found := d.overloads[overloadId]; found {
			return match, nil
		}
		return nil, fmt.Errorf(
			"unknown overload id '%s' for function '%s'",
			ctx.call.Overload, ctx.Function())
	}
	if byArgCount, found := d.functions[ctx.Function()]; found {
		args := ctx.Args()
		candidates := byArgCount[len(args)]
		var matches []*overload
		for _, candidate := range candidates {
			if candidate.handlesArgs(args) {
				matches = append(matches, candidate)
			}
		}
		if len(matches) == 1 {
			return matches[0], nil
		}
		return nil, fmt.Errorf("ambiguous overloads for function '%s'. "+
				"candidates: ['%s']",
				ctx.Function(), matches)
	}
	return nil, fmt.Errorf(
		"no matching overload for function '%s'",
		ctx.Function())
}

// overload with the internal representation of a functions#Overload.
type overload struct {
	function   string
	overloadId string
	argCount   int
	argTypes   []types.Type
	impl       functions.OverloadImpl
}

// handleArgs assesses the count and arg type of the input args and indicates
// whether the overload can be invoked with them.
func (o *overload) handlesArgs(args []interface{}) bool {
	for i, arg := range args {
		argType := o.argTypes[i]
		if !argType.IsDyn() {
			if t, found := types.TypeOf(arg); !found {
				return false
			} else if !argType.Equal(t) {
				return false
			}
		}
	}
	return true
}
