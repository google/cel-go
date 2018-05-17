package types

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
)

var (
	jsonStructType = reflect.TypeOf(&structpb.Struct{})
)

type jsonStruct struct {
	*structpb.Struct
}

func NewJsonStruct(st *structpb.Struct) traits.Mapper {
	return &jsonStruct{st}
}

func (m *jsonStruct) Contains(index ref.Value) ref.Value {
	return !Bool(IsError(m.Get(index)))
}

func (m *jsonStruct) ConvertToNative(refType reflect.Type) (interface{}, error) {
	switch refType.Kind() {
	case reflect.Map:
		otherKey := refType.Key()
		otherElem := refType.Elem()
		if refType.Key().Kind() == reflect.String {
			nativeMap := reflect.MakeMapWithSize(refType, int(m.Size().(Int)))
			it := m.Iterator()
			for it.HasNext() == True {
				key := it.Next()
				refKeyValue, err := key.ConvertToNative(otherKey)
				if err != nil {
					return nil, err
				}
				refElemValue, err := m.Get(key).ConvertToNative(otherElem)
				if err != nil {
					return nil, err
				}
				nativeMap.SetMapIndex(
					reflect.ValueOf(refKeyValue),
					reflect.ValueOf(refElemValue))
			}
			return nativeMap.Interface(), nil
		}
	case reflect.Ptr:
		if refType == jsonValueType {
			return &structpb.Value{
				Kind: &structpb.Value_StructValue{
					StructValue: m.Struct}}, nil
		}
		if refType == jsonStructType {
			return m.Struct, nil
		}
	}
	return nil, fmt.Errorf(
		"no conversion found from map type to native type."+
			" map type: google.protobuf.Struct, native type: %v", refType)
}

func (m *jsonStruct) ConvertToType(typeVal ref.Type) ref.Value {
	switch typeVal {
	case MapType:
		return m
	case TypeType:
		return MapType
	}
	return NewErr("type conversion error from '%s' to '%s'", MapType, typeVal)
}

func (m *jsonStruct) Equal(other ref.Value) ref.Value {
	if MapType != other.Type() {
		return False
	}
	otherMap := other.(traits.Mapper)
	if m.Size() != otherMap.Size() {
		return False
	}
	it := m.Iterator()
	for it.HasNext() == True {
		key := it.Next()
		if otherVal := otherMap.Get(key); IsError(otherVal.Type()) {
			return False
		} else if thisVal := m.Get(key); IsError(thisVal.Type()) {
			return False
		} else if thisVal.Equal(otherVal) != True {
			return False
		}
	}
	return True
}

func (m *jsonStruct) Get(key ref.Value) ref.Value {
	if StringType != key.Type() {
		return NewErr("unsupported key type: '%v", key.Type())
	}
	fields := m.Struct.GetFields()
	value, found := fields[string(key.(String))]
	if !found {
		return NewErr("no such key: '%v'", key)
	}
	return NativeToValue(value)
}

func (m *jsonStruct) Iterator() traits.Iterator {
	f := m.GetFields()
	keys := make([]string, len(m.GetFields()))
	i := 0
	for k, _ := range f {
		keys[i] = k
		i++
	}
	return &jsonValueMapIterator{
		baseIterator: &baseIterator{},
		len:          len(keys),
		mapKeys:      keys}
}

func (m *jsonStruct) Size() ref.Value {
	return Int(len(m.GetFields()))
}

func (m *jsonStruct) Type() ref.Type {
	return MapType
}

func (m *jsonStruct) Value() interface{} {
	return m.Struct
}

type jsonValueMapIterator struct {
	*baseIterator
	cursor  int
	len     int
	mapKeys []string
}

func (it *jsonValueMapIterator) HasNext() ref.Value {
	return Bool(it.cursor < it.len)
}

func (it *jsonValueMapIterator) Next() ref.Value {
	if it.HasNext() == True {
		index := it.cursor
		it.cursor += 1
		return String(it.mapKeys[index])
	}
	return nil
}
