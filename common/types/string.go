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
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// String type implementation which supports addition, comparison, matching,
// and size functions.
type String string

var (
	// StringType singleton.
	StringType = NewTypeValue("string",
		traits.AdderType,
		traits.ComparerType,
		traits.MatcherType,
		traits.SizerType)
)

func (s String) Add(other ref.Value) ref.Value {
	if StringType != other.Type() {
		return NewErr("unsupported overload")
	}
	return s + other.(String)
}

func (s String) Compare(other ref.Value) ref.Value {
	if StringType != other.Type() {
		return NewErr("unsupported overload")
	}
	return Int(strings.Compare(string(s), string(other.(String))))
}

func (s String) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if typeDesc == jsonValueType {
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: s.Value().(string)}}, nil
	}
	if typeDesc.Kind() != reflect.String {
		return nil, fmt.Errorf(
			"unsupported native conversion from string to '%v'", typeDesc)
	}
	return s.Value(), nil
}

func (s String) ConvertToType(typeVal ref.Type) ref.Value {
	switch typeVal {
	case IntType:
		if n, err := strconv.ParseInt(string(s), 10, 64); err == nil {
			return Int(n)
		}
	case UintType:
		if n, err := strconv.ParseUint(string(s), 10, 64); err == nil {
			return Uint(n)
		}
	case DoubleType:
		if n, err := strconv.ParseFloat(string(s), 64); err == nil {
			return Double(n)
		}
	case BoolType:
		if b, err := strconv.ParseBool(string(s)); err == nil {
			return Bool(b)
		}
	case BytesType:
		return Bytes(s)
	case DurationType:
		if d, err := time.ParseDuration(string(s)); err == nil {
			return Duration{ptypes.DurationProto(d)}
		}
	case TimestampType:
		if t, err := time.Parse(time.RFC3339, string(s)); err == nil {
			if ts, err := ptypes.TimestampProto(t); err == nil {
				return Timestamp{ts}
			}
		}
	case StringType:
		return s
	case TypeType:
		return StringType
	}
	return NewErr("type conversion error from '%s' to '%s'", StringType, typeVal)
}

func (s String) Equal(other ref.Value) ref.Value {
	return Bool(StringType == other.Type() && s.Value() == other.Value())
}

func (s String) Match(pattern ref.Value) ref.Value {
	if pattern.Type() != StringType {
		return NewErr("unsupported overload")
	}
	matched, err := regexp.MatchString(string(pattern.(String)), string(s))
	if err != nil {
		return &Err{err}
	}
	return Bool(matched)
}

func (s String) Size() ref.Value {
	return Int(len(string(s)))
}

func (s String) Type() ref.Type {
	return StringType
}

func (s String) Value() interface{} {
	return string(s)
}
