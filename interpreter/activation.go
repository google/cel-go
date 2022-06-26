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
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/cel-go/common/types/ref"
)

// Activation used to resolve identifiers by name and references by id.
//
// An Activation is the primary mechanism by which a caller supplies input into a CEL program.
type Activation interface {
	context.Context
	// ResolveName returns a value from the activation by qualified name, or false if the name
	// could not be found.
	ResolveName(name string) (interface{}, bool)

	// Parent returns the parent of the current activation, may be nil.
	// If non-nil, the parent will be searched during resolve calls.
	Parent() Activation
}

// EmptyActivation returns a variable-free activation.
func EmptyActivation() Activation {
	return &emptyActivation{ctx: context.Background()}
}

// emptyActivation is a variable-free activation.
type emptyActivation struct {
	ctx context.Context
}

func (a *emptyActivation) Deadline() (deadline time.Time, ok bool) { return a.ctx.Deadline() }
func (a *emptyActivation) Done() <-chan struct{}                   { return a.ctx.Done() }
func (a *emptyActivation) Err() error                              { return a.ctx.Err() }
func (a *emptyActivation) Value(key interface{}) interface{}       { return a.ctx.Value }
func (a *emptyActivation) ResolveName(string) (interface{}, bool)  { return nil, false }
func (a *emptyActivation) Parent() Activation                      { return nil }

// NewActivation returns an activation based on a map-based binding where the map keys are
// expected to be qualified names used with ResolveName calls.
//
// The input `bindings` may either be of type `Activation` or `map[string]interface{}`.
//
// Lazy bindings may be supplied within the map-based input in either of the following forms:
// - func() interface{}
// - func() ref.Val
//
// The output of the lazy binding will overwrite the variable reference in the internal map.
//
// Values which are not represented as ref.Val types on input may be adapted to a ref.Val using
// the ref.TypeAdapter configured in the environment.
func NewActivation(bindings interface{}) (Activation, error) {
	if bindings == nil {
		return nil, errors.New("bindings must be non-nil")
	}
	a, isActivation := bindings.(Activation)
	if isActivation {
		return a, nil
	}
	m, isMap := bindings.(map[string]interface{})
	if !isMap {
		return nil, fmt.Errorf(
			"activation input must be an activation or map[string]interface: got %T",
			bindings)
	}
	return &mapActivation{ctx: context.Background(), bindings: m}, nil
}

// mapActivation which implements Activation and maps of named values.
//
// Named bindings may lazily supply values by providing a function which accepts no arguments and
// produces an interface value.
type mapActivation struct {
	ctx      context.Context
	bindings map[string]interface{}
}

func (a *mapActivation) Deadline() (deadline time.Time, ok bool) { return a.ctx.Deadline() }
func (a *mapActivation) Done() <-chan struct{}                   { return a.ctx.Done() }
func (a *mapActivation) Err() error                              { return a.ctx.Err() }
func (a *mapActivation) Value(key interface{}) interface{}       { return a.ctx.Value }

// Parent implements the Activation interface method.
func (a *mapActivation) Parent() Activation {
	return nil
}

// ResolveName implements the Activation interface method.
func (a *mapActivation) ResolveName(name string) (interface{}, bool) {
	obj, found := a.bindings[name]
	if !found {
		return nil, false
	}
	fn, isLazy := obj.(func() ref.Val)
	if isLazy {
		obj = fn()
		a.bindings[name] = obj
	}
	fnRaw, isLazy := obj.(func() interface{})
	if isLazy {
		obj = fnRaw()
		a.bindings[name] = obj
	}
	return obj, found
}

// hierarchicalActivation which implements Activation and contains a parent and
// child activation.
type hierarchicalActivation struct {
	doneOnce *sync.Once
	doneChan <-chan struct{}
	parent   Activation
	child    Activation
}

func (a *hierarchicalActivation) Deadline() (deadline time.Time, ok bool) {
	if d1, ok := a.child.Deadline(); ok {
		if d2, ok := a.parent.Deadline(); ok {
			if d1.Before(d2) {
				return d1, true
			} else {
				return d2, true
			}
		}
		return d1, ok
	}
	return a.parent.Deadline()
}

func (a *hierarchicalActivation) Done() <-chan struct{} {
	a.doneOnce.Do(func() {
		if d1 := a.child.Done(); d1 != nil {
			if d2 := a.parent.Done(); d2 != nil {
				c := make(chan struct{})
				a.doneChan = c
				go func() {
					select {
					case s := <-d1:
						c <- s
					case s := <-d2:
						c <- s
					}
				}()
			}
			a.doneChan = d1
		}
		a.doneChan = a.parent.Done()
	})
	return a.doneChan
}

func (a *hierarchicalActivation) Err() error {
	if err := a.child.Err(); err != nil {
		return err
	} else if err = a.parent.Err(); err != nil {
		return err
	}
	return nil
}

func (a *hierarchicalActivation) Value(key interface{}) interface{} {
	if v := a.child.Value(key); v != nil {
		return v
	}
	return a.parent.Value(key)
}

// Parent implements the Activation interface method.
func (a *hierarchicalActivation) Parent() Activation {
	return a.parent
}

// ResolveName implements the Activation interface method.
func (a *hierarchicalActivation) ResolveName(name string) (interface{}, bool) {
	if object, found := a.child.ResolveName(name); found {
		return object, found
	}
	return a.parent.ResolveName(name)
}

// NewHierarchicalActivation takes two activations and produces a new one which prioritizes
// resolution in the child first and parent(s) second.
func NewHierarchicalActivation(parent Activation, child Activation) Activation {
	return &hierarchicalActivation{&sync.Once{}, nil, parent, child}
}

// NewPartialActivation returns an Activation which contains a list of AttributePattern values
// representing field and index operations that should result in a 'types.Unknown' result.
//
// The `bindings` value may be any value type supported by the interpreter.NewActivation call,
// but is typically either an existing Activation or map[string]interface{}.
func NewPartialActivation(bindings interface{},
	unknowns ...*AttributePattern) (PartialActivation, error) {
	a, err := NewActivation(bindings)
	if err != nil {
		return nil, err
	}
	return &partActivation{Activation: a, unknowns: unknowns}, nil
}

// PartialActivation extends the Activation interface with a set of UnknownAttributePatterns.
type PartialActivation interface {
	Activation

	// UnknownAttributePaths returns a set of AttributePattern values which match Attribute
	// expressions for data accesses whose values are not yet known.
	UnknownAttributePatterns() []*AttributePattern
}

// partActivation is the default implementations of the PartialActivation interface.
type partActivation struct {
	Activation
	unknowns []*AttributePattern
}

// UnknownAttributePatterns implements the PartialActivation interface method.
func (a *partActivation) UnknownAttributePatterns() []*AttributePattern {
	return a.unknowns
}

// varActivation represents a single mutable variable binding.
//
// This activation type should only be used within folds as the fold loop controls the object
// life-cycle.
type varActivation struct {
	parent Activation
	name   string
	val    ref.Val
}

func (a *varActivation) Deadline() (deadline time.Time, ok bool) {
	if a.parent != nil {
		return a.parent.Deadline()
	}
	return time.Time{}, ok
}

func (a *varActivation) Done() <-chan struct{} {
	if a.parent != nil {
		return a.parent.Done()
	}
	return nil
}

func (a *varActivation) Err() error {
	if a.parent != nil {
		return a.parent.Err()
	}
	return nil
}

func (a *varActivation) Value(key interface{}) interface{} {
	if a.parent != nil {
		return a.parent.Value(key)
	}
	return nil
}

// Parent implements the Activation interface method.
func (v *varActivation) Parent() Activation {
	return v.parent
}

// ResolveName implements the Activation interface method.
func (v *varActivation) ResolveName(name string) (interface{}, bool) {
	if name == v.name {
		return v.val, true
	}
	return v.parent.ResolveName(name)
}

var (
	// pool of var activations to reduce allocations during folds.
	varActivationPool = &sync.Pool{
		New: func() interface{} {
			return &varActivation{}
		},
	}
)
