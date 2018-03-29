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

// Activation used to resolve identifiers by name and references by id.
//
// An Activation is the primary mechanism by which a caller supplies input
// into a CEL program.
type Activation interface {

	// ResolveReference returns a value from the activation by expression id,
	// or false if the id-based reference could not be found.
	ResolveReference(exprId int64) (interface{}, bool)

	// ResolveName returns a value from the activation by qualified name, or
	// false if the name could not be found.
	ResolveName(name string) (interface{}, bool)

	// Parent returns the parent of the current activation, may be nil.
	// If non-nil, the parent will be searched during resolve calls.
	Parent() Activation
}

// NewActivation returns an activation based on a map-based binding where the
// map keys are expected to be qualified names used with ResolveName calls.
// TODO: supply references from checked.proto.
func NewActivation(bindings map[string]interface{}) *MapActivation {
	return &MapActivation{bindings: bindings}
}

// MapActivation which implements Activation and maps of named and referenced
// values.
//
// Named bindings may lazily supply values by providing a function which
// accepts no arguments and produces an interface value.
// TODO: consider passing the current activation to the supplier.
type MapActivation struct {
	references map[int64]interface{}
	bindings   map[string]interface{}
}

var _ Activation = &MapActivation{}

func (a *MapActivation) Parent() Activation {
	return nil
}

func (a *MapActivation) ResolveReference(exprId int64) (interface{}, bool) {
	if object, found := a.references[exprId]; found {
		return object, true
	}
	return nil, false
}

func (a *MapActivation) ResolveName(name string) (interface{}, bool) {
	// TODO: Look at how name resolution logic works for enums
	if object, found := a.bindings[name]; found {
		switch object.(type) {
		// Resolve a lazily bound value.
		case func() interface{}:
			return object.(func() interface{})(), true
		// Otherwise, return the bound value.
		default:
			return object, true
		}
	}
	return nil, false
}

// HierarchicalActivation which implements Activation and contains a parent and
// child activation.
type HierarchicalActivation struct {
	parent Activation
	child  Activation
}

var _ Activation = &HierarchicalActivation{}

func (a *HierarchicalActivation) Parent() Activation {
	return a.parent
}

func (a *HierarchicalActivation) ResolveReference(exprId int64) (interface{}, bool) {
	if object, found := a.child.ResolveReference(exprId); found {
		return object, found
	}
	return a.parent.ResolveReference(exprId)
}

func (a *HierarchicalActivation) ResolveName(name string) (interface{}, bool) {
	if object, found := a.child.ResolveName(name); found {
		return object, found
	}
	return a.parent.ResolveName(name)
}

// ExtendActivation takes two activations and produces a new one which prioritizes
// resolution in the child first and parent(s) second.
func ExtendActivation(parent Activation, child Activation) Activation {
	return &HierarchicalActivation{parent, child}
}
