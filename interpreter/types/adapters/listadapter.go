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
	"reflect"
)

type ListAdapter interface {
	objects.Equaler
	objects.Indexer
	objects.Protoer
	objects.Iterable
	Value() interface{}
	Contains(value interface{}) bool
	Concat(other ListAdapter) ListAdapter
	Len() int64
}

var _ ListAdapter = &listAdapter{}

type listAdapter struct {
	value    interface{}
	refValue *reflect.Value
	listType reflect.Type
	elemKind reflect.Kind
}

func NewListAdapter(value interface{}) ListAdapter {
	refValue := reflect.ValueOf(value)
	listType := refValue.Type()
	return &listAdapter{value, &refValue, listType, listType.Elem().Kind()}
}

func (l *listAdapter) Value() interface{} {
	return l.value
}

func (l *listAdapter) Equal(other interface{}) bool {
	adapter, ok := other.(ListAdapter)
	elemCount := l.Len()
	if !ok || elemCount != adapter.Len() {
		return false
	}
	for i := int64(0); i < elemCount; i++ {
		elem, _ := l.Get(i)
		otherElem, _ := adapter.Get(i)
		if elemEqualer, ok := elem.(objects.Equaler); ok {
			if !elemEqualer.Equal(otherElem) {
				return false
			}
		} else if !reflect.DeepEqual(elem, otherElem) {
			return false
		}
	}
	return true
}

func (l *listAdapter) Get(index interface{}) (interface{}, error) {
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

func (l *listAdapter) ToProto(refType reflect.Type) (interface{}, error) {
	protoElem := refType.Elem()
	protoElemKind := protoElem.Kind()
	if protoElemKind == l.elemKind {
		return l.value, nil
	} else if protoElem.ConvertibleTo(l.listType.Elem()) {
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

func (l *listAdapter) Contains(value interface{}) bool {
	for i := 0; i < int(l.Len()); i++ {
		elem, _ := ProtoToExpr(l.refValue.Index(i).Interface())
		matches := false
		switch elem.(type) {
		case objects.Equaler:
			matches = elem.(objects.Equaler).Equal(value)
		default:
			matches = reflect.DeepEqual(elem, value)
		}
		if matches {
			return true
		}
	}
	return false
}

func (l *listAdapter) Concat(other ListAdapter) ListAdapter {
	return &concatAdapter{prev: l, next: other}
}

func (l *listAdapter) Len() int64 {
	return int64(l.refValue.Len())
}

func (l *listAdapter) Iterator() objects.Iterator {
	return &listIterator{listValue: l, cursor: 0, len: l.Len()}
}

var _ ListAdapter = &concatAdapter{}

type concatAdapter struct {
	prev  ListAdapter
	next  ListAdapter
	value interface{}
}

func (l *concatAdapter) Value() interface{} {
	if l.value == nil {
		prevVal := reflect.ValueOf(l.prev.Value())
		nextVal := reflect.ValueOf(l.next.Value())
		merged := make([]interface{}, l.Len(), l.Len())
		prevLen := int(l.prev.Len())
		// TODO: make this more efficient.
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

func (l *concatAdapter) Get(index interface{}) (interface{}, error) {
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
	} else {
		return l.next.Get(i - prevLen)
	}
}

func (l *concatAdapter) Equal(other interface{}) bool {
	adapter, ok := other.(ListAdapter)
	elemCount := l.Len()
	if !ok || elemCount != adapter.Len() {
		return false
	}
	for i := int64(0); i < elemCount; i++ {
		elem, _ := l.Get(i)
		otherElem, _ := adapter.Get(i)
		if elemEqualer, ok := elem.(objects.Equaler); ok &&
			!elemEqualer.Equal(otherElem) {
			return false
		}
		if !reflect.DeepEqual(elem, otherElem) {
			return false
		}
	}
	return true
}

func (l *concatAdapter) ToProto(refType reflect.Type) (interface{}, error) {
	return NewListAdapter(l.Value()).ToProto(refType)
}

func (l *concatAdapter) Contains(value interface{}) bool {
	return l.prev.Contains(value) || l.next.Contains(value)
}

func (l *concatAdapter) Concat(other ListAdapter) ListAdapter {
	return &concatAdapter{prev: l, next: other}
}

func (l *concatAdapter) Len() int64 {
	return l.prev.Len() + l.next.Len()
}

// Iteration functions
func (l *concatAdapter) Iterator() objects.Iterator {
	return &listIterator{listValue: l, cursor: 0, len: l.Len()}
}

type listIterator struct {
	listValue ListAdapter
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
			fmt.Println(err)
		}
		return element
	}
	return nil
}
