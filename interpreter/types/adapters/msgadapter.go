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

package adapters

import (
	"github.com/google/cel-go/interpreter/types/objects"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"reflect"
)

type MsgAdapter interface {
	objects.Equaler
	objects.Indexer
	objects.Iterable
	Value() interface{}
}

type msgAdapter struct {
	value       interface{}
	refValue    *reflect.Value
	isAny       bool
	fieldValues map[string]interface{}
}

func NewMsgAdapter(value interface{}) MsgAdapter {
	// Unwrap any.Any values on construction.
	_, isAny := value.(*any.Any)
	if isAny {
		// TODO: check whether the any being null is valid.
		dynAny := &ptypes.DynamicAny{}
		if err := ptypes.UnmarshalAny(value.(*any.Any), dynAny); err != nil {
			panic(err)
		}
		value = dynAny.Message
	}
	refValue := reflect.ValueOf(value)
	return &msgAdapter{
		value,
		&refValue,
		isAny,
		make(map[string]interface{})}
}

func (m *msgAdapter) Value() interface{} {
	return m.value
}

func (m *msgAdapter) Get(field interface{}) (interface{}, error) {
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
			return NewMsgAdapter(reflect.New(refField.Type()).Interface()), nil
		} else if convertedValue, err := ProtoToExpr(fieldValue); err != nil {
			return nil, err
		} else {
			m.fieldValues[fieldStr] = convertedValue
			return convertedValue, nil
		}
	}
}

func (m *msgAdapter) Equal(other interface{}) bool {
	if adapter, ok := other.(MsgAdapter); ok {
		first := m.value.(proto.Message)
		second := adapter.Value().(proto.Message)
		return proto.Equal(first, second)
	}
	return false
}

func (m *msgAdapter) Iterator() objects.Iterator {
	refType := m.refValue.Type()
	return &msgIterator{
		msgValue:    m,
		msgRefValue: m.refValue,
		msgType:     refType,
		cursor:      int64(0),
		len:         int64(refType.NumMethod())}
}

type msgIterator struct {
	msgValue    MsgAdapter
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
				} else if !fieldRefValue.IsNil() {
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
