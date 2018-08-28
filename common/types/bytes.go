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
	"bytes"
	"fmt"
	refpb "github.com/google/cel-go/common/types/ref"
	traitspb "github.com/google/cel-go/common/types/traits"
	"reflect"
)

// Bytes type that implements refpb.Value and supports add, compare, and size
// operations.
type Bytes []byte

var (
	// BytesType singleton.
	BytesType = NewTypeValue("bytes",
		traitspb.AdderType,
		traitspb.ComparerType,
		traitspb.SizerType)
)

func (b Bytes) Add(other refpb.Value) refpb.Value {
	if BytesType != other.Type() {
		return NewErr("unsupported overload")
	}
	return append(b, other.(Bytes)...)
}

func (b Bytes) Compare(other refpb.Value) refpb.Value {
	if BytesType != other.Type() {
		return NewErr("unsupported overload")
	}
	return Int(bytes.Compare(b, other.(Bytes)))
}

func (b Bytes) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	switch typeDesc.Kind() {
	case reflect.Array, reflect.Slice:
		if typeDesc.Elem().Kind() == reflect.Uint8 {
			return b.Value(), nil
		}
	case reflect.Interface:
		if reflect.TypeOf(b).Implements(typeDesc) {
			return b, nil
		}
	}
	return nil, fmt.Errorf("type conversion error from Bytes to '%v'", typeDesc)
}

func (b Bytes) ConvertToType(typeVal refpb.Type) refpb.Value {
	switch typeVal {
	case StringType:
		return String(b)
	case BytesType:
		return b
	case TypeType:
		return BytesType
	}
	return NewErr("type conversion error from '%s' to '%s'", BytesType, typeVal)
}

func (b Bytes) Equal(other refpb.Value) refpb.Value {
	return Bool(BytesType == other.Type() &&
		bytes.Equal([]byte(b), other.(Bytes)))
}

func (b Bytes) Size() refpb.Value {
	return Int(len(b))
}

func (b Bytes) Type() refpb.Type {
	return BytesType
}

func (b Bytes) Value() interface{} {
	return []byte(b)
}
