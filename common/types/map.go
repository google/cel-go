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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/cel-go/common/types/pb"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	anypb "google.golang.org/protobuf/types/known/anypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
)

// NewDynamicMap returns a traits.Mapper value with dynamic key, value pairs.
func NewDynamicMap(adapter ref.TypeAdapter, value interface{}) traits.Mapper {
	return &baseMap{
		TypeAdapter: adapter,
		value:       value,
		refValue:    reflect.ValueOf(value)}
}

// NewStringStringMap returns a specialized traits.Mapper with string keys and values.
func NewStringStringMap(adapter ref.TypeAdapter, value map[string]string) traits.Mapper {
	return &stringMap{
		baseMap:   &baseMap{TypeAdapter: adapter, value: value},
		mapStrStr: value,
	}
}

// NewProtoMap returns a specialized traits.Mapper for handling protobuf map values.
func NewProtoMap(adapter ref.TypeAdapter, value *pb.Map) traits.Mapper {
	return &protoMap{
		TypeAdapter: adapter,
		value:       value,
	}
}

var (
	// MapType singleton.
	MapType = NewTypeValue("map",
		traits.ContainerType,
		traits.IndexerType,
		traits.IterableType,
		traits.SizerType)
)

// baseMap is a reflection based map implementation designed to handle a variety of map-like types.
type baseMap struct {
	ref.TypeAdapter
	value    interface{}
	refValue reflect.Value
}

// Contains implements the traits.Container interface method.
func (m *baseMap) Contains(index ref.Val) ref.Val {
	val, found := m.Find(index)
	// When the index is not found and val is non-nil, this is an error or unknown value.
	if !found && val != nil && IsUnknownOrError(val) {
		return val
	}
	return Bool(found)
}

// ConvertToNative implements the ref.Val interface method.
func (m *baseMap) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	// If the map is already assignable to the desired type return it, e.g. interfaces and
	// maps with the same key value types.
	if reflect.TypeOf(m).AssignableTo(typeDesc) {
		return m, nil
	}
	switch typeDesc {
	case anyValueType:
		json, err := m.ConvertToNative(jsonStructType)
		if err != nil {
			return nil, err
		}
		return anypb.New(json.(proto.Message))
	case jsonValueType, jsonStructType:
		jsonEntries, err :=
			m.ConvertToNative(reflect.TypeOf(map[string]*structpb.Value{}))
		if err != nil {
			return nil, err
		}
		jsonMap := &structpb.Struct{
			Fields: jsonEntries.(map[string]*structpb.Value)}
		if typeDesc == jsonStructType {
			return jsonMap, nil
		}
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: jsonMap}}, nil
	}

	// Unwrap pointers, but track their use.
	isPtr := false
	if typeDesc.Kind() == reflect.Ptr {
		tk := typeDesc
		typeDesc = typeDesc.Elem()
		if typeDesc.Kind() == reflect.Ptr {
			return nil, fmt.Errorf("unsupported type conversion to '%v'", tk)
		}
		isPtr = true
	}

	// Establish some basic facts about the map key and value types.
	thisType := m.refValue.Type()
	thisKey := thisType.Key()
	thisKeyKind := thisKey.Kind()
	thisElem := thisType.Elem()
	thisElemKind := thisElem.Kind()

	switch typeDesc.Kind() {
	// Map conversion.
	case reflect.Map:
		otherKey := typeDesc.Key()
		otherKeyKind := otherKey.Kind()
		otherElem := typeDesc.Elem()
		otherElemKind := otherElem.Kind()
		if otherKeyKind == thisKeyKind && otherElemKind == thisElemKind {
			return m.value, nil
		}
		elemCount := m.Size().(Int)
		nativeMap := reflect.MakeMapWithSize(typeDesc, int(elemCount))
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
	case reflect.Struct:
		if thisKeyKind != reflect.String && thisKeyKind != reflect.Interface {
			break
		}
		nativeStructPtr := reflect.New(typeDesc)
		nativeStruct := nativeStructPtr.Elem()
		it := m.Iterator()
		for it.HasNext() == True {
			key := it.Next()
			// Ensure the field name being referenced is exported.
			// Only exported (public) field names can be set by reflection, where the name
			// must be at least one character in length and start with an upper-case letter.
			fieldName := string(key.ConvertToType(StringType).(String))
			switch len(fieldName) {
			case 0:
				return nil, errors.New("type conversion error, unsupported empty field")
			case 1:
				fieldName = strings.ToUpper(fieldName)
			default:
				fieldName = strings.ToUpper(fieldName[0:1]) + fieldName[1:]
			}
			fieldRef := nativeStruct.FieldByName(fieldName)
			if !fieldRef.IsValid() {
				return nil, fmt.Errorf(
					"type conversion error, no such field '%s' in type '%v'",
					fieldName, typeDesc)
			}
			fieldValue, err := m.Get(key).ConvertToNative(fieldRef.Type())
			if err != nil {
				return nil, err
			}
			fieldRef.Set(reflect.ValueOf(fieldValue))
		}
		if isPtr {
			return nativeStructPtr.Interface(), nil
		}
		return nativeStruct.Interface(), nil
	}
	return nil, fmt.Errorf("type conversion error from map to '%v'", typeDesc)
}

// ConvertToType implements the ref.Val interface method.
func (m *baseMap) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case MapType:
		return m
	case TypeType:
		return MapType
	}
	return NewErr("type conversion error from '%s' to '%s'", MapType, typeVal)
}

// Equal implements the ref.Val interface method.
func (m *baseMap) Equal(other ref.Val) ref.Val {
	if MapType != other.Type() {
		return ValOrErr(other, "no such overload")
	}
	otherMap := other.(traits.Mapper)
	if m.Size() != otherMap.Size() {
		return False
	}
	it := m.Iterator()
	for it.HasNext() == True {
		key := it.Next()
		thisVal, _ := m.Find(key)
		otherVal, found := otherMap.Find(key)
		if !found {
			if otherVal == nil {
				return False
			}
			return ValOrErr(otherVal, "no such overload")
		}
		valEq := thisVal.Equal(otherVal)
		if valEq != True {
			return valEq
		}
	}
	return True
}

// Find implements the traits.Mapper interface method.
func (m *baseMap) Find(key ref.Val) (ref.Val, bool) {
	// TODO: There are multiple reasons why a Get could fail. Typically, this is because the key
	// does not exist in the map; however, it's possible that the value cannot be converted to
	// the desired type. Refine this strategy to disambiguate these cases.
	if IsUnknownOrError(key) {
		return key, false
	}
	thisKeyType := m.refValue.Type().Key()
	nativeKey, err := key.ConvertToNative(thisKeyType)
	if err != nil {
		return &Err{err}, false
	}
	nativeKeyVal := reflect.ValueOf(nativeKey)
	value := m.refValue.MapIndex(nativeKeyVal)
	if !value.IsValid() {
		return nil, false
	}
	return m.NativeToValue(value.Interface()), true
}

// Get implements the traits.Indexer interface method.
func (m *baseMap) Get(key ref.Val) ref.Val {
	v, found := m.Find(key)
	if !found {
		return ValOrErr(v, "no such key: %v", key)
	}
	return v
}

// Iterator implements the traits.Iterable interface method.
func (m *baseMap) Iterator() traits.Iterator {
	mapKeys := m.refValue.MapKeys()
	return &mapIterator{
		TypeAdapter: m.TypeAdapter,
		mapKeys:     mapKeys,
		len:         int(m.Size().(Int))}
}

// Size implements the traits.Sizer interface method.
func (m *baseMap) Size() ref.Val {
	return Int(m.refValue.Len())
}

// Type implements the ref.Val interface method.
func (m *baseMap) Type() ref.Type {
	return MapType
}

// Value implements the ref.Val interface method.
func (m *baseMap) Value() interface{} {
	return m.value
}

// stringMap is a specialization to improve the performance of simple key, value pair lookups by
// string as this is the most common usage of maps.
type stringMap struct {
	*baseMap
	mapStrStr map[string]string
}

// Contains implements the traits.Container interface method.
func (m *stringMap) Contains(index ref.Val) ref.Val {
	val, found := m.Find(index)
	// When the index is not found and val is non-nil, this is an error or unknown value.
	if !found && val != nil && IsUnknownOrError(val) {
		return val
	}
	return Bool(found)
}

// ConvertToNative implements the ref.Val interface method.
func (m *stringMap) ConvertToNative(refType reflect.Type) (interface{}, error) {
	if !m.baseMap.refValue.IsValid() {
		m.baseMap.refValue = reflect.ValueOf(m.value)
	}
	return m.baseMap.ConvertToNative(refType)
}

// ConvertToType implements the ref.Val interface method.
func (m *stringMap) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case MapType:
		return m
	default:
		return m.baseMap.ConvertToType(typeVal)
	}
}

// Equal implements the ref.Val interface method.
func (m *stringMap) Equal(other ref.Val) ref.Val {
	if !m.baseMap.refValue.IsValid() {
		m.baseMap.refValue = reflect.ValueOf(m.value)
	}
	return m.baseMap.Equal(other)
}

// Find implements the traits.Mapper interface method.
func (m *stringMap) Find(key ref.Val) (ref.Val, bool) {
	strKey, ok := key.(String)
	if !ok {
		return ValOrErr(key, "no such overload"), false
	}
	val, found := m.mapStrStr[string(strKey)]
	if !found {
		return nil, false
	}
	return String(val), true
}

// Get implements the traits.Indexer interface method.
func (m *stringMap) Get(key ref.Val) ref.Val {
	v, found := m.Find(key)
	if !found {
		return ValOrErr(v, "no such key: %v", key)
	}
	return v
}

// Iterator implements the traits.Iterable interface method.
func (m *stringMap) Iterator() traits.Iterator {
	if !m.baseMap.refValue.IsValid() {
		m.baseMap.refValue = reflect.ValueOf(m.value)
	}
	return m.baseMap.Iterator()
}

// Size implements the traits.Sizer interface method.
func (m *stringMap) Size() ref.Val {
	return Int(len(m.mapStrStr))
}

type protoMap struct {
	ref.TypeAdapter
	value *pb.Map
}

func (m *protoMap) Contains(key ref.Val) ref.Val {
	val, found := m.Find(key)
	// When the index is not found and val is non-nil, this is an error or unknown value.
	if !found && val != nil && IsUnknownOrError(val) {
		return val
	}
	return Bool(found)
}

func (m *protoMap) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	// If the map is already assignable to the desired type return it, e.g. interfaces and
	// maps with the same key value types.
	if reflect.TypeOf(m).AssignableTo(typeDesc) {
		return m, nil
	}
	switch typeDesc {
	case anyValueType:
		json, err := m.ConvertToNative(jsonStructType)
		if err != nil {
			return nil, err
		}
		return anypb.New(json.(proto.Message))
	case jsonValueType, jsonStructType:
		jsonEntries, err :=
			m.ConvertToNative(reflect.TypeOf(map[string]*structpb.Value{}))
		if err != nil {
			return nil, err
		}
		jsonMap := &structpb.Struct{
			Fields: jsonEntries.(map[string]*structpb.Value)}
		if typeDesc == jsonStructType {
			return jsonMap, nil
		}
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: jsonMap}}, nil
	}
	if typeDesc.Kind() != reflect.Map {
		return nil, fmt.Errorf("unsupported type conversion: %v to map", typeDesc)
	}
	keyType := m.value.KeyType.ReflectType()
	valType := m.value.ValueType.ReflectType()
	otherKeyType := typeDesc.Key()
	otherValType := typeDesc.Elem()
	mapVal := reflect.MakeMapWithSize(typeDesc, m.value.Len())
	var err error
	m.value.Range(func(key protoreflect.MapKey, val protoreflect.Value) bool {
		ntvKey := key.Interface()
		ntvVal := val.Interface()
		switch ntvVal.(type) {
		case protoreflect.Message:
			ntvVal = ntvVal.(protoreflect.Message).Interface()
		}
		if keyType == otherKeyType && valType == otherValType {
			mapVal.SetMapIndex(reflect.ValueOf(ntvKey), reflect.ValueOf(ntvVal))
			return true
		}
		celKey := m.NativeToValue(ntvKey)
		celVal := m.NativeToValue(ntvVal)
		ntvKey, err = celKey.ConvertToNative(otherKeyType)
		if err != nil {
			return false
		}
		ntvVal, err = celVal.ConvertToNative(otherValType)
		if err != nil {
			return false
		}
		mapVal.SetMapIndex(reflect.ValueOf(ntvKey), reflect.ValueOf(ntvVal))
		return true
	})
	if err != nil {
		return nil, err
	}
	return mapVal.Interface(), nil
}

func (m *protoMap) ConvertToType(typeVal ref.Type) ref.Val {
	switch typeVal {
	case MapType:
		return m
	case TypeType:
		return MapType
	}
	return NewErr("type conversion error from '%s' to '%s'", MapType, typeVal)
}

func (m *protoMap) Equal(other ref.Val) ref.Val {
	if MapType != other.Type() {
		return ValOrErr(other, "no such overload")
	}
	otherMap := other.(traits.Mapper)
	if m.value.Map.Len() != int(otherMap.Size().(Int)) {
		return False
	}
	var retVal ref.Val = True
	m.value.Range(func(key protoreflect.MapKey, val protoreflect.Value) bool {
		keyVal := m.NativeToValue(key.Interface())
		valVal := m.NativeToValue(val)
		otherVal, found := otherMap.Find(keyVal)
		if !found {
			if otherVal == nil {
				retVal = False
				return false
			}
			retVal = ValOrErr(otherVal, "no such overload")
			return false
		}
		valEq := valVal.Equal(otherVal)
		if valEq != True {
			retVal = valEq
			return false
		}
		return true
	})
	return retVal
}

func (m *protoMap) Find(key ref.Val) (ref.Val, bool) {
	if IsUnknownOrError(key) {
		return key, false
	}
	// TODO: ensure that the key.Value() type is what is desired by the map
	// This means that the type information needs to be provided to the proto map
	// And I'm not sure how to do this.
	ntvKey, err := key.ConvertToNative(m.value.KeyType.ReflectType())
	if err != nil {
		return &Err{err}, false
	}
	val := m.value.Get(protoreflect.ValueOf(ntvKey).MapKey()).Interface()
	switch v := val.(type) {
	case protoreflect.EnumNumber:
		return Int(v), true
	case protoreflect.List:
		return m.NativeToValue(&pb.List{
			List:     v,
			ElemType: m.value.ValueType,
		}), true
	case protoreflect.Map:
		return m.NativeToValue(&pb.Map{
			Map:       v,
			KeyType:   m.value.ValueType.KeyType,
			ValueType: m.value.ValueType.ValueType},
		), true
	case protoreflect.Message:
		return m.NativeToValue(v.Interface()), true
	default:
		return m.NativeToValue(v), true
	}
}

// Get implements the traits.Indexer interface method.
func (m *protoMap) Get(key ref.Val) ref.Val {
	v, found := m.Find(key)
	if !found {
		return ValOrErr(v, "no such key: %v", key)
	}
	return v
}

func (m *protoMap) Iterator() traits.Iterator {
	mapKeys := make([]protoreflect.MapKey, 0, m.value.Len())
	m.value.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		mapKeys = append(mapKeys, k)
		return true
	})
	return &protoMapIterator{
		TypeAdapter: m.TypeAdapter,
		mapKeys:     mapKeys,
		len:         m.value.Len(),
	}
}

func (m *protoMap) Size() ref.Val {
	return Int(m.value.Len())
}

func (m *protoMap) Type() ref.Type {
	return MapType
}

func (m *protoMap) Value() interface{} {
	return m.value.Map
}

type mapIterator struct {
	*baseIterator
	ref.TypeAdapter
	mapKeys []reflect.Value
	cursor  int
	len     int
}

// HasNext implements the traits.Iterator interface method.
func (it *mapIterator) HasNext() ref.Val {
	return Bool(it.cursor < it.len)
}

// Next implements the traits.Iterator interface method.
func (it *mapIterator) Next() ref.Val {
	if it.HasNext() == True {
		index := it.cursor
		it.cursor++
		refKey := it.mapKeys[index]
		return it.NativeToValue(refKey.Interface())
	}
	return nil
}

type protoMapIterator struct {
	*baseIterator
	ref.TypeAdapter
	mapKeys []protoreflect.MapKey
	cursor  int
	len     int
}

// HasNext implements the traits.Iterator interface method.
func (it *protoMapIterator) HasNext() ref.Val {
	return Bool(it.cursor < it.len)
}

// Next implements the traits.Iterator interface method.
func (it *protoMapIterator) Next() ref.Val {
	if it.HasNext() == True {
		index := it.cursor
		it.cursor++
		refKey := it.mapKeys[index]
		return it.NativeToValue(refKey.Interface())
	}
	return nil
}
