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

package types

import (
	"fmt"
	refpb "github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
)

var (
	// TypeType is the type of a TypeValue.
	TypeType = NewTypeValue("type")
)

// TypeValue is an instance of a Value that describes a value's type.
type TypeValue struct {
	name      string
	traitMask int
}

// NewTypeValue returns *TypeValue which is both a ref.Type and ref.Value.
func NewTypeValue(name string, traits ...int) *TypeValue {
	traitMask := 0
	for _, trait := range traits {
		traitMask |= trait
	}
	return &TypeValue{
		name:      name,
		traitMask: traitMask}
}

// NewObjectTypeValue returns a *TypeValue based on the input name, which is
// annotated with the traits relevant to all objects.
func NewObjectTypeValue(name string) *TypeValue {
	return NewTypeValue(name,
		traits.IndexerType,
		traits.IterableType)
}

func (t *TypeValue) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	// TODO: replace the internal type representation with a proto-value.
	return nil, fmt.Errorf("type conversion not supported for 'type'.")
}

func (t *TypeValue) ConvertToType(typeVal refpb.Type) refpb.Value {
	switch typeVal {
	case TypeType:
		return t
	case StringType:
		return String(t.TypeName())
	}
	return NewErr("type conversion error from '%s' to '%s'", TypeType, typeVal)
}

func (t *TypeValue) Equal(other refpb.Value) refpb.Value {
	return Bool(TypeType == other.Type() && t.Value() == other.Value())
}

func (t *TypeValue) HasTrait(trait int) bool {
	return trait&t.traitMask == trait
}

func (t *TypeValue) Type() refpb.Type {
	return TypeType
}

func (t *TypeValue) TypeName() string {
	return t.name
}

func (t *TypeValue) Value() interface{} {
	return t.name
}

func (t *TypeValue) String() string {
	return t.name
}
