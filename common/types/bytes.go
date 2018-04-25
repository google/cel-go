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
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
)

// Bytes type that implements ref.Value and supports add, compare, and size
// operations.
type Bytes []byte

var (
	// BytesType singleton.
	BytesType = NewTypeValue("bytes",
		traits.AdderType,
		traits.ComparerType,
		traits.SizerType)
)

func (b Bytes) Add(other ref.Value) ref.Value {
	if BytesType != other.Type() {
		return NewErr("unsupported overload")
	}
	return append(b, other.(Bytes)...)
}

func (b Bytes) Compare(other ref.Value) ref.Value {
	if BytesType != other.Type() {
		return NewErr("unsupported overload")
	}
	return Int(bytes.Compare(b, other.(Bytes)))
}

func (b Bytes) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return b.Value(), nil
}

func (b Bytes) ConvertToType(typeVal ref.Type) ref.Value {
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

func (b Bytes) Equal(other ref.Value) ref.Value {
	return Bool(BytesType == other.Type() &&
		bytes.Equal([]byte(b), other.(Bytes)))
}

func (b Bytes) Size() ref.Value {
	return Int(len(b))
}

func (b Bytes) Type() ref.Type {
	return BytesType
}

func (b Bytes) Value() interface{} {
	return []byte(b)
}
