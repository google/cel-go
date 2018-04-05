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

// Adapter package defines utilities for adapting plain Go structs into
// structs suitable for consumption with CEL.
package types

import (
	"fmt"
	"github.com/google/cel-go/interpreter/types/traits"
	"reflect"
)

// ListValue for moderating between plain Go lists and list values, and
// CEL list values.
//
// ListValues are comparable, iterable, support indexed access,
// and may be converted to a protobuf representation.
type ListValue interface {
	traits.Equaler
	traits.Indexer
	traits.Iterable
	traits.Protoer

	// Value of the underlying list.
	Value() interface{}

	// Contains returns true if the input value exists in the list.
	Contains(value interface{}) bool

	// Concat two lists together into an adapted list.
	Concat(other ListValue) ListValue

	// Len returns the length of the list.
	Len() int64
}

// Implementation of a ListValue which uses reflection to mediate between
// Go and CEL types.
type listValue struct {
	value    interface{}
	refValue *reflect.Value
	listType reflect.Type
	elemKind reflect.Kind
}

// NewListValue wraps a list value into a ListValue
func NewListValue(value interface{}) ListValue {
	refValue := reflect.ValueOf(value)
	listType := refValue.Type()
	return &listValue{
		value,
		&refValue,
		listType,
		listType.Elem().Kind()}
}

func (l *listValue) Value() interface{} {
	return l.value
}

// Equal returns true if two lists are deeply equal.
//
// The other value must be a list of the same element type, length,
// where each element at each ordinal is deeply equal.
func (l *listValue) Equal(other interface{}) bool {
	adapter, ok := other.(ListValue)
	elemCount := l.Len()
	if !ok || elemCount != adapter.Len() {
		return false
	}
	for i := int64(0); i < elemCount; i++ {
		elem, _ := l.Get(i)
		otherElem, _ := adapter.Get(i)
		if elemEqualer, ok := elem.(traits.Equaler); ok {
			if !elemEqualer.Equal(otherElem) {
				return false
			}
		} else if !reflect.DeepEqual(elem, otherElem) {
			return false
		}
	}
	return true
}

// Get a list element at the specified index.
func (l *listValue) Get(index interface{}) (interface{}, error) {
	i, ok := index.(int64)
	if !ok {
		return nil, fmt.Errorf("unexpected index type")
	}
	if i < 0 || i >= l.Len() {
		return nil, fmt.Errorf("index out of range")
	}
	value := l.refValue.Index(int(i)).Interface()
	return ProtoToExpr(value)
}

func (l *listValue) ToProto(refType reflect.Type) (interface{}, error) {
	protoElem := refType.Elem()
	protoElemKind := protoElem.Kind()
	if protoElemKind == l.elemKind {
		return l.value, nil
	}
	if protoElem.ConvertibleTo(l.listType.Elem()) {
		elemCount := int(l.Len())
		protoList := reflect.MakeSlice(refType, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			if elem, err := ExprToProto(protoElem,
				l.refValue.Index(int(i)).Interface()); err != nil {
				return nil, err
			} else {
				protoList.Index(i).Set(reflect.ValueOf(elem))
			}
		}
		return protoList.Interface(), nil
	}
	return nil, fmt.Errorf(
		"no conversion found from list type to proto."+
			" list elem: %v, proto elem: %v", l.elemKind, protoElemKind)
}

func (l *listValue) Contains(value interface{}) bool {
	for i := 0; i < int(l.Len()); i++ {
		elem, _ := ProtoToExpr(l.refValue.Index(i).Interface())
		matches := false
		switch elem.(type) {
		case traits.Equaler:
			matches = elem.(traits.Equaler).Equal(value)
		default:
			matches = reflect.DeepEqual(elem, value)
		}
		if matches {
			return true
		}
	}
	return false
}

func (l *listValue) Concat(other ListValue) ListValue {
	return &concatListValue{prev: l, next: other}
}

func (l *listValue) Len() int64 {
	return int64(l.refValue.Len())
}

func (l *listValue) Iterator() traits.Iterator {
	return &listIterator{listValue: l, cursor: 0, len: l.Len()}
}

type concatListValue struct {
	prev  ListValue
	next  ListValue
	value interface{}
}

func (l *concatListValue) Value() interface{} {
	if l.value == nil {
		prevVal := reflect.ValueOf(l.prev.Value())
		nextVal := reflect.ValueOf(l.next.Value())
		merged := make([]interface{}, l.Len(), l.Len())
		prevLen := int(l.prev.Len())
		for i := 0; i < prevLen; i++ {
			merged[i] = prevVal.Index(i).Interface()
		}
		for j := 0; j < int(l.next.Len()); j++ {
			merged[prevLen+j] = nextVal.Index(j).Interface()
		}
		l.value = merged
	}
	return l.value
}

func (l *concatListValue) Get(index interface{}) (interface{}, error) {
	i, ok := index.(int64)
	if !ok {
		return nil, fmt.Errorf("unexpected index type")
	}
	if i < 0 || i >= l.Len() {
		return nil, fmt.Errorf("index out of range")
	}
	prevLen := l.prev.Len()
	if i < prevLen {
		return l.prev.Get(i)
	}
	return l.next.Get(i - prevLen)
}

func (l *concatListValue) Equal(other interface{}) bool {
	adapter, ok := other.(ListValue)
	elemCount := l.Len()
	if !ok || elemCount != adapter.Len() {
		return false
	}
	for i := int64(0); i < elemCount; i++ {
		elem, _ := l.Get(i)
		otherElem, _ := adapter.Get(i)
		if elemEqualer, ok := elem.(traits.Equaler); ok &&
			!elemEqualer.Equal(otherElem) {
			return false
		}
		if !reflect.DeepEqual(elem, otherElem) {
			return false
		}
	}
	return true
}

func (l *concatListValue) ToProto(refType reflect.Type) (interface{}, error) {
	return NewListValue(l.Value()).ToProto(refType)
}

func (l *concatListValue) Contains(value interface{}) bool {
	return l.prev.Contains(value) || l.next.Contains(value)
}

func (l *concatListValue) Concat(other ListValue) ListValue {
	return &concatListValue{prev: l, next: other}
}

func (l *concatListValue) Len() int64 {
	return l.prev.Len() + l.next.Len()
}

func (l *concatListValue) Iterator() traits.Iterator {
	return &listIterator{listValue: l, cursor: 0, len: l.Len()}
}

type listIterator struct {
	listValue ListValue
	cursor    int64
	len       int64
}

func (it *listIterator) HasNext() bool {
	return it.cursor < it.len
}

func (it *listIterator) Next() interface{} {
	if it.HasNext() {
		index := it.cursor
		it.cursor += 1
		element, err := it.listValue.Get(index)
		if err != nil {
			return err
		}
		return element
	}
	return nil
}
