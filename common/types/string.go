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
	"regexp"
	"strconv"
	"strings"
	"time"

	ptypespb "github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	refpb "github.com/google/cel-go/common/types/ref"
	pbpb "github.com/google/cel-go/common/types/pb"
)

// String type implementation which supports addition, comparison, matching,
// and size functions.
type String string

var (
	// StringType singleton.
	StringType = NewTypeValue("string",
		pbpb.AdderType,
		pbpb.ComparerType,
		pbpb.MatcherType,
		pbpb.SizerType)
)

func (s String) Add(other refpb.Value) refpb.Value {
	if StringType != other.Type() {
		return NewErr("unsupported overload")
	}
	return s + other.(String)
}

func (s String) Compare(other refpb.Value) refpb.Value {
	if StringType != other.Type() {
		return NewErr("unsupported overload")
	}
	return Int(strings.Compare(string(s), string(other.(String))))
}

func (s String) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	switch typeDesc.Kind() {
	case reflect.String:
		return string(s), nil
	case reflect.Ptr:
		if typeDesc == jsonValueType {
			return &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: s.Value().(string)}}, nil
		}
		if typeDesc.Elem().Kind() == reflect.String {
			p := string(s)
			return &p, nil
		}
	case reflect.Interface:
		if reflect.TypeOf(s).Implements(typeDesc) {
			return s, nil
		}
	}
	return nil, fmt.Errorf(
		"unsupported native conversion from string to '%v'", typeDesc)
}

func (s String) ConvertToType(typeVal refpb.Type) refpb.Value {
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
			return Duration{ptypespb.DurationProto(d)}
		}
	case TimestampType:
		if t, err := time.Parse(time.RFC3339, string(s)); err == nil {
			if ts, err := ptypespb.TimestampProto(t); err == nil {
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

func (s String) Equal(other refpb.Value) refpb.Value {
	return Bool(StringType == other.Type() && s.Value() == other.Value())
}

func (s String) Match(pattern refpb.Value) refpb.Value {
	if pattern.Type() != StringType {
		return NewErr("unsupported overload")
	}
	matched, err := regexp.MatchString(string(pattern.(String)), string(s))
	if err != nil {
		return &Err{err}
	}
	return Bool(matched)
}

func (s String) Size() refpb.Value {
	return Int(len(string(s)))
}

func (s String) Type() refpb.Type {
	return StringType
}

func (s String) Value() interface{} {
	return string(s)
}
