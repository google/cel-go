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
	traitspb "github.com/google/cel-go/common/types/traits"
	refpb "github.com/google/cel-go/common/types/ref"
)

// Uint type implementation which supports comparison and math operators.
type Uint uint64

var (
	// UintType singleton.
	UintType = NewTypeValue("uint",
		traitspb.AdderType,
		traitspb.ComparerType,
		traitspb.DividerType,
		traitspb.ModderType,
		traitspb.MultiplierType,
		traitspb.SubtractorType)
)

const (
	uintZero = Uint(0)
)

func (i Uint) Add(other refpb.Value) refpb.Value {
	if UintType != other.Type() {
		return NewErr("unsupported overload")
	}
	return i + other.(Uint)
}

func (i Uint) Compare(other refpb.Value) refpb.Value {
	if UintType != other.Type() {
		return NewErr("unsupported overload")
	}
	if i < other.(Uint) {
		return IntNegOne
	}
	if i > other.(Uint) {
		return IntOne
	}
	return IntZero
}

func (i Uint) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	value := i.Value()
	switch typeDesc.Kind() {
	case reflect.Uint32:
		return uint32(value.(uint64)), nil
	case reflect.Uint64:
		return value, nil
	case reflect.Ptr:
		if typeDesc == jsonValueType {
			return &structpb.Value{
				Kind: &structpb.Value_NumberValue{
					NumberValue: float64(i)}}, nil
		}
		switch typeDesc.Elem().Kind() {
		case reflect.Uint32:
			p := uint32(i)
			return &p, nil
		case reflect.Uint64:
			p := uint64(i)
			return &p, nil
		}
	case reflect.Interface:
		if reflect.TypeOf(i).Implements(typeDesc) {
			return i, nil
		}
	}
	return nil, fmt.Errorf("unsupported type conversion from 'uint' to %v", typeDesc)
}

func (i Uint) ConvertToType(typeVal refpb.Type) refpb.Value {
	switch typeVal {
	case IntType:
		return Int(i)
	case UintType:
		return i
	case DoubleType:
		return Double(i)
	case StringType:
		return String(fmt.Sprintf("%d", uint64(i)))
	case TypeType:
		return UintType
	}
	return NewErr("type conversion error from '%s' to '%s'", UintType, typeVal)
}

func (i Uint) Divide(other refpb.Value) refpb.Value {
	if UintType != other.Type() {
		return NewErr("unsupported overload")
	}
	otherUint := other.(Uint)
	if otherUint == uintZero {
		return NewErr("divide by zero")
	}
	return i / otherUint
}

func (i Uint) Equal(other refpb.Value) refpb.Value {
	return Bool(UintType == other.Type() &&
		i.Value() == other.Value())
}

func (i Uint) Modulo(other refpb.Value) refpb.Value {
	if UintType != other.Type() {
		return NewErr("unsupported overload")
	}
	otherUint := other.(Uint)
	if otherUint == uintZero {
		return NewErr("modulus by zero")
	}
	return i % otherUint
}

func (i Uint) Multiply(other refpb.Value) refpb.Value {
	if UintType != other.Type() {
		return NewErr("unsupported overload")
	}
	return i * other.(Uint)
}

func (i Uint) Subtract(subtrahend refpb.Value) refpb.Value {
	if UintType != subtrahend.Type() {
		return NewErr("unsupported overload")
	}
	return i - subtrahend.(Uint)
}

func (i Uint) Type() refpb.Type {
	return UintType
}

func (i Uint) Value() interface{} {
	return uint64(i)
}
