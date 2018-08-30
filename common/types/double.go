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

	structpb "github.com/golang/protobuf/ptypes/struct"
	refpb "github.com/google/cel-go/common/types/ref"
	traitspb "github.com/google/cel-go/common/types/traits"
)

// Double type that implements refpb.Value, comparison, and mathematical
// operations.
type Double float64

var (
	// DoubleType singleton.
	DoubleType = NewTypeValue("double",
		traitspb.AdderType,
		traitspb.ComparerType,
		traitspb.DividerType,
		traitspb.MultiplierType,
		traitspb.NegatorType,
		traitspb.SubtractorType)
)

func (d Double) Add(other refpb.Value) refpb.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	return d + other.(Double)
}

func (d Double) Compare(other refpb.Value) refpb.Value {
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
	switch typeDesc.Kind() {
	case reflect.Float32:
		return float32(d), nil
	case reflect.Float64:
		return float64(d), nil
	case reflect.Ptr:
		if typeDesc == jsonValueType {
			return &structpb.Value{
				Kind: &structpb.Value_NumberValue{
					NumberValue: float64(d)}}, nil
		}
		switch typeDesc.Elem().Kind() {
		case reflect.Float32:
			p := float32(d)
			return &p, nil
		case reflect.Float64:
			p := float64(d)
			return &p, nil
		}
	case reflect.Interface:
		if reflect.TypeOf(d).Implements(typeDesc) {
			return d, nil
		}
	}
	return nil, fmt.Errorf("type conversion error from Double to '%v'", typeDesc)
}

func (d Double) ConvertToType(typeVal refpb.Type) refpb.Value {
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

func (d Double) Divide(other refpb.Value) refpb.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	if other.(Double) == Double(0) {
		return NewErr("divide by zero")
	}
	return d / other.(Double)
}

func (d Double) Equal(other refpb.Value) refpb.Value {
	return Bool(DoubleType == other.Type() && d == other)
}

func (d Double) Multiply(other refpb.Value) refpb.Value {
	if DoubleType != other.Type() {
		return NewErr("unsupported overload")
	}
	return d * other.(Double)
}

func (d Double) Negate() refpb.Value {
	return -d
}

func (d Double) Subtract(subtrahend refpb.Value) refpb.Value {
	if DoubleType != subtrahend.Type() {
		return NewErr("unsupported overload")
	}
	return d - subtrahend.(Double)
}

func (d Double) Type() refpb.Type {
	return DoubleType
}

func (d Double) Value() interface{} {
	return float64(d)
}
