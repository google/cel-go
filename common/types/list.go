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
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
)

var (
	// ListType singleton.
	ListType = NewTypeValue("list",
		traits.AdderType,
		traits.ContainerType,
		traits.IndexerType,
		traits.IterableType,
		traits.SizerType)
)

// NewDynamicList returns a traits.Lister with heterogenous elements.
func NewDynamicList(value interface{}) traits.Lister {
	return &baseList{value, reflect.ValueOf(value)}
}

// NewStringList returns a traits.Lister containing only strings.
func NewStringList(elems []string) traits.Lister {
	return &stringList{
		baseList: NewDynamicList(elems).(*baseList),
		elems:    elems}
}

// baseList points to a list containing elements of any type.
type baseList struct {
	value    interface{}
	refValue reflect.Value
}

// concatList combines two list implementations together into a view.
type concatList struct {
	value    interface{}
	prevList traits.Lister
	nextList traits.Lister
}

// stringList is a specialization of the traits.Lister interface which is
// present to demonstrate the ability to specialize Lister implementations.
type stringList struct {
	*baseList
	elems []string
}

func (l *baseList) Add(other ref.Value) ref.Value {
	if other.Type() != ListType {
		return NewErr("no such overload")
	}
	return &concatList{
		prevList: l,
		nextList: other.(traits.Lister)}
}

func (l *concatList) Add(other ref.Value) ref.Value {
	if other.Type() != ListType {
		return NewErr("no such overload")
	}
	return &concatList{
		prevList: l,
		nextList: other.(traits.Lister)}
}

func (l *stringList) Add(other ref.Value) ref.Value {
	if other.Type() != ListType {
		if IsError(other) || IsUnknown(other) {
			return other
		}
		return NewErr("no such overload")
	}
	otherList := other.(traits.Lister)
	switch otherList.(type) {
	case *stringList:
		concatElems := append(l.elems, otherList.(*stringList).elems...)
		return NewStringList(concatElems)
	}
	return &concatList{
		prevList: l.baseList,
		nextList: other.(traits.Lister)}
}

func (l *baseList) Contains(elem ref.Value) ref.Value {
	if IsError(elem) || IsUnknown(elem) {
		return elem
	}
	for i := Int(0); i < l.Size().(Int); i++ {
		if l.Get(i).Equal(elem) == True {
			return True
		}
	}
	return False
}

func (l *concatList) Contains(elem ref.Value) ref.Value {
	return Bool(l.prevList.Contains(elem) == True ||
		l.nextList.Contains(elem) == True)
}

func (l *baseList) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	thisType := l.refValue.Type()
	thisElem := thisType.Elem()
	thisElemKind := thisElem.Kind()
	nativeElem := typeDesc.Elem()
	nativeElemKind := nativeElem.Kind()
	if nativeElem.ConvertibleTo(thisElem) {
		elemCount := int(l.Size().(Int))
		nativeList := reflect.MakeSlice(typeDesc, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			elem := l.Get(Int(i))
			nativeElemVal, err := elem.ConvertToNative(nativeElem)
			if err != nil {
				return nil, err
			}
			nativeList.Index(i).Set(reflect.ValueOf(nativeElemVal))
		}
		return nativeList.Interface(), nil
	}
	return nil, fmt.Errorf(
		"no conversion found from list type to native type."+
			" list elem: %v, native elem type: %v", thisElemKind, nativeElemKind)
}

func (l *concatList) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	combined := &baseList{
		value:    l.Value(),
		refValue: reflect.ValueOf(l.Value())}
	return combined.ConvertToNative(typeDesc)
}

func (l *stringList) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	return l.elems, nil
}

func (l *baseList) ConvertToType(typeVal ref.Type) ref.Value {
	switch typeVal {
	case ListType:
		return l
	case TypeType:
		return ListType
	}
	return NewErr("type conversion error from '%s' to '%s'", ListType, typeVal)
}

func (l *concatList) ConvertToType(typeVal ref.Type) ref.Value {
	switch typeVal {
	case ListType:
		return l
	case TypeType:
		return ListType
	}
	return NewErr("type conversion error from '%s' to '%s'", ListType, typeVal)
}

func (l *baseList) Equal(other ref.Value) ref.Value {
	if ListType != other.Type() {
		return False
	}
	otherList := other.(traits.Lister)
	if l.Size() != otherList.Size() {
		return False
	}
	for i := IntZero; i < l.Size().(Int); i++ {
		thisElem := l.Get(i)
		otherElem := otherList.Get(i)
		if thisElem.Equal(otherElem) != True {
			return False
		}
	}
	return True
}

func (l *concatList) Equal(other ref.Value) ref.Value {
	if ListType != other.Type() {
		return False
	}
	otherList := other.(traits.Lister)
	if l.Size() != otherList.Size() {
		return False
	}
	for i := IntZero; i < l.Size().(Int); i++ {
		thisElem := l.Get(i)
		otherElem := otherList.Get(i)
		if thisElem.Equal(otherElem) != True {
			return False
		}
	}
	return True
}

func (l *baseList) Get(index ref.Value) ref.Value {
	if index.Type() != IntType {
		if IsError(index) || IsUnknown(index) {
			return index
		}
		return NewErr("unsupported index type '%s' in list", index.Type())
	}
	i := index.(Int)
	if i < 0 || i >= l.Size().(Int) {
		return NewErr("index '%d' out of range in list size '%d'", i, l.Size())
	}
	elem := l.refValue.Index(int(i)).Interface()
	return NativeToValue(elem)
}

func (l *concatList) Get(index ref.Value) ref.Value {
	if index.Type() != IntType {
		if IsError(index) || IsUnknown(index) {
			return index
		}
		return NewErr("unsupported index type '%s' in list", index.Type())
	}
	i := index.(Int)
	if i < l.prevList.Size().(Int) {
		return l.prevList.Get(i)
	}
	offset := i - l.prevList.Size().(Int)
	return l.nextList.Get(offset)
}

func (l *stringList) Get(index ref.Value) ref.Value {
	if index.Type() != IntType {
		if IsError(index) || IsUnknown(index) {
			return index
		}
		return NewErr("unsupported index type '%s' in list", index.Type())
	}
	i := index.(Int)
	if i < 0 || i >= l.Size().(Int) {
		return NewErr("index '%d' out of range in list size '%d'", i, l.Size())
	}
	return String(l.elems[i])
}

func (l *baseList) Iterator() traits.Iterator {
	return &listIterator{
		baseIterator: &baseIterator{},
		listValue:    l,
		cursor:       0,
		len:          l.Size().(Int)}
}

func (l *concatList) Iterator() traits.Iterator {
	return &listIterator{
		baseIterator: &baseIterator{},
		listValue:    l,
		cursor:       0,
		len:          l.Size().(Int)}
}

func (l *baseList) Size() ref.Value {
	return Int(l.refValue.Len())
}

func (l *concatList) Size() ref.Value {
	return l.prevList.Size().(Int).Add(l.nextList.Size())
}

func (l *stringList) Size() ref.Value {
	return Int(len(l.elems))
}

func (l *baseList) Type() ref.Type {
	return ListType
}

func (l *concatList) Type() ref.Type {
	return ListType
}

func (l *baseList) Value() interface{} {
	return l.value
}

func (l *concatList) Value() interface{} {
	if l.value == nil {
		prevVal := reflect.ValueOf(l.prevList.Value())
		nextVal := reflect.ValueOf(l.nextList.Value())
		merged := make([]interface{}, l.Size().(Int), l.Size().(Int))
		prevLen := int(l.prevList.Size().(Int))
		for i := 0; i < prevLen; i++ {
			merged[i] = prevVal.Index(i).Interface()
		}
		for j := 0; j < int(l.nextList.Size().(Int)); j++ {
			merged[prevLen+j] = nextVal.Index(j).Interface()
		}
		l.value = merged
	}
	return l.value
}

type listIterator struct {
	*baseIterator
	listValue traits.Lister
	cursor    Int
	len       Int
}

func (it *listIterator) HasNext() ref.Value {
	return Bool(it.cursor < it.len)
}

func (it *listIterator) Next() ref.Value {
	if it.HasNext() == True {
		index := it.cursor
		it.cursor += 1
		return it.listValue.Get(index)
	}
	return nil
}
