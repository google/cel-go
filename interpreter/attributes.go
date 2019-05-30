// Copyright 2019 Google LLC
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
	"fmt"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// Attribute refers to a Variable and zero or more field or index traversals.
type Attribute interface {
	// Variable referenced by the attribute.
	//
	// When the Attribute is absolute, the Variable value is non-nil. When the Attribute is
	// relative such as might happen when refering to the field or index of a object-based
	// return value from a function, the Variable value will be nil.
	Variable() Variable

	// Path of field traversals against reference object, either a variable or a contextual
	// intermediate value computed during CEL evaluation.
	//
	// If Variable is nil, the Path must contain at least one PathElem.
	Path() []*PathElem

	// Select a path element from the current Attribute to produce a more fully qualified
	// Attribute.
	Select(*PathElem) Attribute
}

// Variable refers to a type-level expression identifier corresponding to program input.
type Variable interface {
	// ID indicates the expression id where the variable appears.
	//
	// Non-zero when the Variable is identified during program plan time.
	//
	// Zero-ish when the Variable is provided as part of a user-provided Attribute specification
	// for use with partial state evaluation.
	ID() int64

	// Name of the variable, may be simple or namespaced, e.g. 'in' or 'cloud.iam.resource'.
	Name() string
}

// PathElem representing a fragment of the attribute selection path.
type PathElem struct {
	// ID of the expression where the PathElem occurs, non-zero if determined during program plan
	// time.
	ID int64

	// ToValue function which resolve the PathElem to a concrete ref.Val.
	//
	// In most the PathElem will be a constant value; however, there are cases when an Attribute
	// may be cross-referenced, e.g. a[b] where the Attribute of {var: a, path: [b]} depends on
	// the resolution of 'b' before the Attribute reference is fully known. In these cases, the
	// ToValue function reference permits flexible resolution of the cross-reference.
	ToValue func(Activation) ref.Val
}

// NewUnknownAttribute creates an Attribute that refers to a top-level variable or one of its
// descendants fields in a selection path.
//
// A selection path element may be of the following scalar types:
//
// - bool, types.Bool
// - int, int32, int64, types.Int
// - string, types.String
// - uint, uint32, uint64, types.Uint
//
// The string path element '*' is treated as a wildcard match against any path element value.
//
// As absolute Attribute references within an expression are being resolved, the attribute path
// is compared to the set of known-unknowns which refer to the same top-level variable. When the
// unknown Attribute matches any part of the resolved attribute, a types.Unknown value is
// generated. The types.Unknown value refers to the expression id of the last resolved path
// element that matches the unknown.
func NewUnknownAttribute(varName string, pathElems ...interface{}) (Attribute, error) {
	v := &variable{name: varName}
	p := make([]*PathElem, len(pathElems), len(pathElems))
	for i, elem := range pathElems {
		var e ref.Val
		switch elem.(type) {
		case bool:
			e = types.Bool(elem.(bool))
		case int:
			e = types.Int(elem.(int))
		case int32:
			e = types.Int(elem.(int32))
		case int64:
			e = types.Int(elem.(int64))
		case string:
			e = types.String(elem.(string))
		case uint:
			e = types.Uint(elem.(uint))
		case uint32:
			e = types.Uint(elem.(uint32))
		case uint64:
			e = types.Uint(elem.(uint64))
		case types.Bool, types.Int, types.String, types.Uint:
			e = elem.(ref.Val)
		default:
			return nil, fmt.Errorf("invalid scalar element type in attribute path")
		}
		p[i] = &PathElem{
			ToValue: func(Activation) ref.Val { return e },
		}
	}
	return &attribute{variable: v, path: p}, nil
}

// newExprVarAttribute creates a Attribute reference to a Variable and zero or more path elements.
func newExprVarAttribute(id int64, varName string, pathElems ...*PathElem) Attribute {
	return &attribute{
		variable: &variable{id: id, name: varName},
		path:     pathElems,
	}
}

// newExprRelAttribute creates a contextual reference to a field within a dynamically computed
// object.
func newExprRelAttribute(pe *PathElem) Attribute {
	return &attribute{
		path: []*PathElem{pe},
	}
}

// newExprPathElem creates a new PathElem value based on a constant value.
func newExprPathElem(id int64, val ref.Val) *PathElem {
	return &PathElem{
		ID:      id,
		ToValue: func(Activation) ref.Val { return val },
	}
}

// variable is the default implementation of the Variable interface.
type variable struct {
	id   int64
	name string
}

// ID implements the Variable interface method.
func (v *variable) ID() int64 {
	return v.id
}

// Name implements the Variable interface method.
func (v *variable) Name() string {
	return v.name
}

// attribute is the default implementation of the Attribute interface.
type attribute struct {
	variable Variable
	path     []*PathElem
}

// Variable implements the Attribute interface method.
func (a *attribute) Variable() Variable {
	return a.variable
}

// Path implements the Attribute interface method.
func (a *attribute) Path() []*PathElem {
	return a.path
}

// Select implements the Attribute interface method.
func (a *attribute) Select(pe *PathElem) Attribute {
	return &attribute{
		variable: a.variable,
		path:     append(a.path, pe),
	}
}
