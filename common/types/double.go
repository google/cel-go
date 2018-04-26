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
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
)

// Double type that implements ref.Value, comparison, and mathematical
// operations.
type Double float64

var (
	// DoubleType singleton.
	DoubleType = NewTypeValue("double",
		traits.AdderType,
		traits.ComparerType,
		traits.DividerType,
		traits.MultiplierType,
		traits.NegatorType,
		traits.SubtractorType)
)

func (d Double) Add(other ref.Value) ref.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	return d + other.(Double)
}

func (d Double) Compare(other ref.Value) ref.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	if d < other.(Double) {
		return IntNegOne
	}
	if d > other.(Double) {
		return IntOne
	}
	return IntZero
}

func (d Double) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	value := d.Value()
	refKind := typeDesc.Kind()
	switch refKind {
	case reflect.Float32:
		return float32(value.(float64)), nil
	case reflect.Float64:
		return value, nil
	}
	return nil, fmt.Errorf("unsupported type conversion from 'double' to %v", typeDesc)
}

func (d Double) ConvertToType(typeVal ref.Type) ref.Value {
	switch typeVal {
	case IntType:
		return Int(float64(d))
	case UintType:
		return Uint(float64(d))
	case DoubleType:
		return d
	case StringType:
		return String(fmt.Sprintf("%g", float64(d)))
	case TypeType:
		return DoubleType
	}
	return NewErr("type conversion error from '%s' to '%s'", DoubleType, typeVal)
}

func (d Double) Divide(other ref.Value) ref.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	if other.(Double) == Double(0) {
		return NewErr("divide by zero")
	}
	return d / other.(Double)
}

func (d Double) Equal(other ref.Value) ref.Value {
	return Bool(DoubleType == other.Type() && d == other)
}

func (d Double) Multiply(other ref.Value) ref.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	return d * other.(Double)
}

func (d Double) Negate() ref.Value {
	return -d
}

func (d Double) Subtract(subtrahend ref.Value) ref.Value {
	if DoubleType != subtrahend.Type() {
		return NewErr("unsupported overload")
	}
	return d - subtrahend.(Double)
}

func (d Double) Type() ref.Type {
	return DoubleType
}

func (d Double) Value() interface{} {
	return float64(d)
}
