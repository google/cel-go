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
	"reflect"
	"strconv"

	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Bool type that implements ref.Val and supports comparison and negation.
type Bool bool

var (
	// BoolType singleton.
	BoolType = NewTypeValue("bool",
		traits.ComparerType,
		traits.NegatorType)

	// boolWrapperType golang reflected type for protobuf bool wrapper type.
	boolWrapperType = reflect.TypeOf(&wrapperspb.BoolValue{})
)

// Boolean constants
var (
	False = Bool(false)
	True  = Bool(true)
)

// Compare implements the traits.Comparer interface method.
func (b Bool) Compare(other ref.Val) ref.Val {
	otherBool, ok := other.(Bool)
	if !ok {
		return ValOrErr(other, "no such overload")
	}
	if b == otherBool {
		return IntZero
	}
	if !b && otherBool {
		return IntNegOne
	}
	return IntOne
}

// ConvertToNative implements the ref.Val interface method.
func (b Bool) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	switch typeDesc.Kind() {
	case reflect.Bool:
		return bool(b), nil
	case reflect.Ptr:
		switch typeDesc {
		case anyValueType:
			// Primitives must be wrapped before being set on an Any field.
			return anypb.New(wrapperspb.Bool(bool(b)))
		case boolWrapperType:
			// Convert the bool to a protobuf.BoolValue.
			return wrapperspb.Bool(bool(b)), nil
		case jsonValueType:
			return structpb.NewBoolValue(bool(b)), nil
		default:
			if typeDesc.Elem().Kind() == reflect.Bool {
				p := bool(b)
				return &p, nil
			}
		}
	case reflect.Interface:
		bv := b.Value()
		if reflect.TypeOf(bv).Implements(typeDesc) {
			return bv, nil
		}
		if reflect.TypeOf(b).Implements(typeDesc) {
			return b, nil
		}
	}
	return nil, fmt.Errorf("type conversion error from bool to '%v'", typeDesc)
}

// ConvertToType implements the ref.Val interface method.
func (b Bool) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case StringType:
		return String(strconv.FormatBool(bool(b)))
	case BoolType:
		return b
	case TypeType:
		return BoolType
	}
	return NewErr("type conversion error from '%v' to '%v'", BoolType, typeVal)
}

// Equal implements the ref.Val interface method.
func (b Bool) Equal(other ref.Val) ref.Val {
	otherBool, ok := other.(Bool)
	if !ok {
		return ValOrErr(other, "no such overload")
	}
	return Bool(b == otherBool)
}

// Negate implements the traits.Negater interface method.
func (b Bool) Negate() ref.Val {
	return !b
}

// Type implements the ref.Val interface method.
func (b Bool) Type() ref.Type {
	return BoolType
}

// Value implements the ref.Val interface method.
func (b Bool) Value() interface{} {
	return bool(b)
}

// IsBool returns whether the input ref.Val or ref.Type is equal to BoolType.
func IsBool(elem interface{}) bool {
	switch elem := elem.(type) {
	case ref.Type:
		return elem == BoolType
	case ref.Val:
		return IsBool(elem.Type())
	}
	return false
}
