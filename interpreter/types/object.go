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
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/google/cel-go/interpreter/types/traits"
	"reflect"
)

// ObjectValue wraps accesses to protobuf messages in order to ensure
// appropriate type conversions between Go and CEL.
//
// ObjectValue values are comparable, iterable, and support indexed access by
// field name.
//
// Note: iterators will iterate over all default and non-default fields, but
// for oneof field sets will only include the field name of the non-default
// oneof.
type ObjectValue interface {
	traits.Equaler
	traits.Indexer
	traits.Iterable

	// Value of the underlying message.
	Value() interface{}
}

type protoValue struct {
	value       interface{}
	refValue    *reflect.Value
	isAny       bool
	fieldValues map[string]interface{}
}

func NewProtoValue(value interface{}) ObjectValue {
	// Unwrap any.Any values on construction.
	_, isAny := value.(*any.Any)
	if isAny {
		// TODO: check whether the any being null is valid.
		dynAny := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(value.(*any.Any), dynAny); err != nil {
			// TODO: find a better way to handle unknown types.
			panic(err)
		}
		value = dynAny.Message
	}
	refValue := reflect.ValueOf(value)
	return &protoValue{
		value,
		&refValue,
		isAny,
		make(map[string]interface{})}
}

func (m *protoValue) Value() interface{} {
	return m.value
}

func (m *protoValue) Get(field interface{}) (interface{}, error) {
	if fieldStr, ok := field.(string); !ok {
		return nil, fmt.Errorf("unexpected field type")
	} else if fieldValue, found := m.fieldValues[fieldStr]; found {
		return fieldValue, nil
	} else {
		if refGetter := m.refValue.MethodByName("Get" + fieldStr); !refGetter.IsValid() {
			return nil, fmt.Errorf("no such field")
		} else if refFieldValue := refGetter.Call([]reflect.Value{})[0]; !refFieldValue.IsValid() {
			return nil, fmt.Errorf("no such field")
		} else if fieldValue := refFieldValue.Interface(); fieldValue == nil {
			// TODO: test getting an empty Any field.
			refField := m.refValue.FieldByName(fieldStr)
			return NewProtoValue(reflect.New(refField.Type()).Interface()), nil
		} else if convertedValue, err := ProtoToExpr(fieldValue); err != nil {
			return nil, err
		} else {
			m.fieldValues[fieldStr] = convertedValue
			return convertedValue, nil
		}
	}
}

func (m *protoValue) Equal(other interface{}) bool {
	if adapter, ok := other.(ObjectValue); ok {
		first := m.value.(proto.Message)
		second := adapter.Value().(proto.Message)
		return proto.Equal(first, second)
	}
	return false
}

func (m *protoValue) Iterator() traits.Iterator {
	refType := m.refValue.Type()
	return &msgIterator{
		msgValue:    m,
		msgRefValue: m.refValue,
		msgType:     refType,
		cursor:      int64(0),
		len:         int64(refType.NumMethod())}
}

type msgIterator struct {
	msgValue    ObjectValue
	msgRefValue *reflect.Value
	msgType     reflect.Type
	cursor      int64
	cursorField string
	len         int64
}

func (it *msgIterator) HasNext() bool {
	if it.cursorField != "" {
		return true
	}
	it.cursorField = ""
	for it.cursor < it.len {
		index := it.cursor
		it.cursor += 1
		refMethod := it.msgType.Method(int(index))
		if refMethod.Name[:3] == "Get" {
			refValue := *it.msgRefValue
			fieldRefValue := refMethod.Func.Call([]reflect.Value{refValue})[0]
			if fieldRefValue.IsValid() {
				fieldRefKind := fieldRefValue.Kind()
				if fieldRefKind != reflect.Ptr && fieldRefKind != reflect.Interface {
					it.cursorField = refMethod.Name[3:]
					return true
				}
				if !fieldRefValue.IsNil() {
					fieldValue := fieldRefValue.Interface()
					if _, isProto := fieldValue.(proto.Message); isProto {
						it.cursorField = refMethod.Name[3:]
						return true
					}
				}
			}
		}
	}
	return false
}

func (it *msgIterator) Next() interface{} {
	if it.cursorField != "" {
		next := it.cursorField
		it.cursorField = ""
		return next
	}
	it.HasNext()
	return nil
}
