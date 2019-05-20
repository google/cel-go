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
	"errors"
	"fmt"
	"sync"

	"github.com/google/cel-go/common/types/ref"
)

// Activation describes the variables available to an evaluation scope.
type Activation interface {
	// Find returns a value from the activation by qualified name, or false if the name
	// could not be found.
	Find(name string) (interface{}, bool)

	// Parent returns the parent of the current activation, may be nil.
	// If non-nil, the parent will be searched during resolve calls.
	Parent() Activation
}

// PartialActivation describes a mix of variables and unknown attributes within an evaluation
// scope.
//
// A variable or some portion of its attributes may be treated as known-unknowns by CEL when the
// FindUnknowns function returns a non-empty Attribute set.
//
// Note: PartialActivation must be used in conjunction with an interpreter.UnknownResolver to work.
// Otherwise, the unknowns will be treated as errors.
type PartialActivation interface {
	Activation

	// FindUnknowns returns the collection of unknown Attribute values associated with the
	// variable name provided in the input, if present.
	FindUnknowns(name string) ([]Attribute, bool)
}

// EmptyActivation is an evaluation context with no variables.
func EmptyActivation() Activation {
	return emptyActivation
}

// NewActivation returns an activation based on a map-based binding where the map keys are
// expected to be qualified names used with Find calls.
//
// The input `bindings` may either be of type `Activation` or `map[string]interface{}`.
//
// When the bindings are a `map` form whose values are not of `ref.Val` type, the values will be
// converted to CEL values (if possible) using the `types.DefaultTypeAdapter`.
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
	return &mapActivation{bindings: m}, nil
}

// NewPartialActivation creates a PartialActivation containing a mix of variables and unknown
// attributes.
func NewPartialActivation(bindings interface{}, unknowns []Attribute) (PartialActivation, error) {
	act, err := NewActivation(bindings)
	if err != nil {
		return nil, err
	}
	unknownMap := make(map[string][]Attribute)
	for _, unk := range unknowns {
		varName := unk.Variable().Name()
		attrs, found := unknownMap[varName]
		if !found {
			unknownMap[varName] = []Attribute{unk}
		} else {
			unknownMap[varName] = append(attrs, unk)
		}
	}
	return &partialActivation{Activation: act, unknowns: unknownMap}, nil
}

// UnknownActivation returns a PartialActivation which treats all variables and their attributes
// as unknown.
func UnknownActivation() PartialActivation {
	return unkActivation
}

// NewHierarchicalActivation takes two activations and produces a new one which prioritizes
// resolution in the child first and parent(s) second.
func NewHierarchicalActivation(parent Activation, child Activation) Activation {
	return &hierarchicalActivation{parent, child}
}

// mapActivation which implements Activation and maps of named values.
//
// Named bindings may lazily supply values by providing a function which accepts no arguments and
// produces an interface value.
type mapActivation struct {
	bindings map[string]interface{}
}

// Parent implements the Activation interface method.
func (a *mapActivation) Parent() Activation {
	return nil
}

// Find implements the Activation interface method.
func (a *mapActivation) Find(name string) (interface{}, bool) {
	if object, found := a.bindings[name]; found {
		switch object.(type) {
		// Resolve a lazily bound value.
		case func() ref.Val:
			return object.(func() ref.Val)(), true
		// Otherwise, return the bound value.
		default:
			return object, true
		}
	}
	return nil, false
}

// hierarchicalActivation which implements Activation and contains a parent and
// child activation.
type hierarchicalActivation struct {
	parent Activation
	child  Activation
}

// Parent implements the Activation interface method.
func (a *hierarchicalActivation) Parent() Activation {
	return a.parent
}

// Find implements the Activation interface method.
func (a *hierarchicalActivation) Find(name string) (interface{}, bool) {
	if object, found := a.child.Find(name); found {
		return object, found
	}
	return a.parent.Find(name)
}

// partialActivation implements the PartialActivation interface.
type partialActivation struct {
	Activation
	unknowns map[string][]Attribute
}

// FindUnknowns implements the PartialActivation interface method.
func (a *partialActivation) FindUnknowns(name string) ([]Attribute, bool) {
	unk, found := a.unknowns[name]
	return unk, found
}

// unknownActivation implements the PartialActivation interface.
type unknownActivation struct {
	Activation
}

// FindUnknowns implements the interpreter.PartialActivation interface method.
func (un *unknownActivation) FindUnknowns(name string) ([]Attribute, bool) {
	return emptyAttrs, true
}

// newVarActivation returns a new varActivation instance.
func newVarActivation(parent Activation, name string) *varActivation {
	return &varActivation{
		parent: parent,
		name:   name,
	}
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

// Parent implements the Activation interface method.
func (v *varActivation) Parent() Activation {
	return v.parent
}

// Find implements the Activation interface method.
func (v *varActivation) Find(name string) (interface{}, bool) {
	if name == v.name {
		return v.val, true
	}
	return v.parent.Find(name)
}

var (
	// empty attribute set used to treat all attributes as unknown.
	emptyAttrs = []Attribute{newExprVarAttribute(1, "")}

	// pool of var activations to reduce allocations during folds.
	varActivationPool = &sync.Pool{
		New: func() interface{} {
			return &varActivation{}
		},
	}

	// static empty activation instance.
	emptyActivation = &mapActivation{bindings: map[string]interface{}{}}

	// static unknown activation instance.
	unkActivation = &unknownActivation{Activation: emptyActivation}
)
