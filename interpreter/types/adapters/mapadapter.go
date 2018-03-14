package adapters

import (
	"github.com/google/cel-go/interpreter/types/objects"
	"fmt"
	"reflect"
)

type MapAdapter interface {
	objects.Equaler
	objects.Indexer
	objects.Protoer
	objects.Iterable
	Value() interface{}
	Len() int64
}

func NewMapAdapter(value interface{}) MapAdapter {
	refValue := reflect.ValueOf(value)
	mapType := refValue.Type()
	return &mapAdapter{value, &refValue, mapType, mapType.Key().Kind(), mapType.Elem().Kind()}
}

var _ MapAdapter = &mapAdapter{}

type mapAdapter struct {
	value    interface{}
	refValue *reflect.Value
	mapType  reflect.Type
	keyKind  reflect.Kind
	elemKind reflect.Kind
}

func (m *mapAdapter) Value() interface{} {
	return m.value
}

func (m *mapAdapter) Equal(other interface{}) bool {
	adapter, ok := other.(*mapAdapter)
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
			} else if thisEqualerVal, ok := thisExprVal.(objects.Equaler); ok {
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

func (m *mapAdapter) Get(key interface{}) (interface{}, error) {
	if protoKey, err := ExprToProto(m.mapType.Key(), key); err != nil {
		return nil, err
	} else if value := m.refValue.MapIndex(reflect.ValueOf(protoKey)); !value.IsValid() {
		return nil, fmt.Errorf("no such key")
	} else {
		return ProtoToExpr(value.Interface())
	}
}

func (m *mapAdapter) ToProto(refType reflect.Type) (interface{}, error) {
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

func (m *mapAdapter) Len() int64 {
	return int64(m.refValue.Len())
}

func (m *mapAdapter) Iterator() objects.Iterator {
	return &mapIterator{
		mapValue: m,
		mapKeys:  m.refValue.MapKeys(),
		cursor:   0,
		len:      m.Len()}
}

type mapIterator struct {
	mapValue *mapAdapter
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
			fmt.Println(err)
		}
		return key
	}
	return nil
}
