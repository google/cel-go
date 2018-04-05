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

// MapValue mediates type conversions between Go and CEL types for
// map values.
//
// MapValues are comparable, iterable, support indexed access by key,
// and may be converted to a protobuf representation.
type MapValue interface {
	traits.Equaler
	traits.Indexer
	traits.Protoer
	traits.Iterable

	// Value of the underlying map.
	Value() interface{}

	// Number of entries in the map.
	Len() int64
}

func NewMapValue(value interface{}) MapValue {
	refValue := reflect.ValueOf(value)
	mapType := refValue.Type()
	return &mapValue{
		value,
		&refValue,
		mapType,
		mapType.Key().Kind(),
		mapType.Elem().Kind()}
}

type mapValue struct {
	value    interface{}
	refValue *reflect.Value
	mapType  reflect.Type
	keyKind  reflect.Kind
	elemKind reflect.Kind
}

func (m *mapValue) Value() interface{} {
	return m.value
}

func (m *mapValue) Equal(other interface{}) bool {
	adapter, ok := other.(*mapValue)
	if !ok || m.Len() != adapter.Len() {
		return false
	}
	adapterVal := reflect.ValueOf(adapter.Value())
	otherKeyType := adapterVal.Type().Key()
	if !otherKeyType.ConvertibleTo(m.mapType.Key()) ||
		!m.mapType.Key().ConvertibleTo(otherKeyType) {
		return false
	}
	for _, key := range m.refValue.MapKeys() {
		if otherRefVal := adapterVal.MapIndex(key); !otherRefVal.IsValid() {
			return false
		} else if thisRefVal := m.refValue.MapIndex(key); !thisRefVal.IsValid() {
			return false
		} else {
			thisVal := thisRefVal.Interface()
			otherVal := otherRefVal.Interface()
			if thisExprVal, err := ProtoToExpr(thisVal); err != nil {
				fmt.Print(err)
				return false
			} else if otherExprVal, err := ProtoToExpr(otherVal); err != nil {
				fmt.Print(err)
				return false
			} else if thisEqualerVal, ok := thisExprVal.(traits.Equaler); ok {
				if !thisEqualerVal.Equal(otherExprVal) {
					return false
				}
			} else if !reflect.DeepEqual(thisExprVal, otherExprVal) {
				return false
			}
		}
	}
	return true
}

func (m *mapValue) Get(key interface{}) (interface{}, error) {
	if protoKey, err := ExprToProto(m.mapType.Key(), key); err != nil {
		return nil, err
	} else if value := m.refValue.MapIndex(reflect.ValueOf(protoKey)); !value.IsValid() {
		return nil, fmt.Errorf("no such key")
	} else {
		return ProtoToExpr(value.Interface())
	}
}

func (m *mapValue) ToProto(refType reflect.Type) (interface{}, error) {
	protoKey := refType.Key()
	protoKeyKind := protoKey.Kind()
	protoElem := refType.Elem()
	protoElemKind := protoElem.Kind()
	if protoKeyKind == m.keyKind && protoElemKind == m.elemKind {
		return m.value, nil
	} else if protoKey.ConvertibleTo(m.mapType.Key()) &&
		protoElem.ConvertibleTo(m.mapType.Elem()) {
		elemCount := int(m.Len())
		protoMap := reflect.MakeMapWithSize(refType, elemCount)
		for _, refKey := range m.refValue.MapKeys() {
			if refKeyValue, err := ExprToProto(protoKey,
				refKey.Interface()); err != nil {
				return nil, err
			} else {
				if refElemValue, err := ExprToProto(protoElem,
					m.refValue.MapIndex(refKey).Interface()); err != nil {
					return nil, err
				} else {
					protoMap.SetMapIndex(reflect.ValueOf(refKeyValue),
						reflect.ValueOf(refElemValue))
				}
			}
		}
		return protoMap.Interface(), nil
	}
	return nil, fmt.Errorf(
		"no conversion found from map type to proto."+
			" map: %v, proto: %v", m.mapType, refType)
}

func (m *mapValue) Len() int64 {
	return int64(m.refValue.Len())
}

func (m *mapValue) Iterator() traits.Iterator {
	return &mapIterator{
		mapValue: m,
		mapKeys:  m.refValue.MapKeys(),
		cursor:   0,
		len:      m.Len()}
}

type mapIterator struct {
	mapValue *mapValue
	mapKeys  []reflect.Value
	cursor   int64
	len      int64
}

func (it *mapIterator) HasNext() bool {
	return it.cursor < it.len
}

func (it *mapIterator) Next() interface{} {
	if it.HasNext() {
		index := it.cursor
		it.cursor += 1
		refKey := it.mapKeys[index]
		key, err := ProtoToExpr(refKey.Interface())
		if err != nil {
			return err
		}
		return key
	}
	return nil
}
