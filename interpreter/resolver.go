// Copyright 2019 Google LLC
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

package interpreter

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	anypb "github.com/golang/protobuf/ptypes/any"
	structpb "github.com/golang/protobuf/ptypes/struct"
)

type Resolver interface {
	Resolve(Activation, Attribute) (interface{}, bool)

	ResolveRelative(interface{}, Activation, Attribute) interface{}
}

func NewResolver(adapter ref.TypeAdapter, provider ref.TypeProvider) Resolver {
	return &defaultResolver{adapter: adapter, provider: provider}
}

func NewUnknownResolver(adapter ref.TypeAdapter,
	provider ref.TypeProvider,
	unknowns []Attribute) Resolver {
	unknownMap := make(map[string][]Attribute)
	for _, unk := range unknowns {
		varName := unk.Variable().Name()
		unks, found := unknownMap[varName]
		if found {
			unks = append(unks, unk)
		} else {
			unks = []Attribute{unk}
		}
		unknownMap[varName] = unks
	}
	return &unknownResolver{
		defaultResolver: &defaultResolver{adapter: adapter, provider: provider},
		unknowns:        unknownMap,
	}
}

type defaultResolver struct {
	adapter  ref.TypeAdapter
	provider ref.TypeProvider
}

func (res *defaultResolver) Resolve(vars Activation, attr Attribute) (interface{}, bool) {
	varName := attr.Variable().Name()
	attrPath := attr.Path()
	obj, found := vars.Find(varName)
	if found {
		for _, elem := range attrPath {
			obj = res.getElem(obj, elem.ToValue(vars))
		}
		return obj, true
	}
	typ, found := res.provider.FindIdent(varName)
	if found {
		return typ, true
	}
	if !found && len(attrPath) > 0 {
		for _, elem := range attrPath {
			elemVal := elem.ToValue(vars)
			elemStr, ok := elemVal.(types.String)
			if !ok {
				return nil, false
			}
			varName = fmt.Sprintf("%s.%s", varName, elemStr)
			typ, found = res.provider.FindIdent(varName)
			if found {
				return typ, true
			}
		}
		return obj, true
	}
	return nil, false
}

func (res *defaultResolver) ResolveRelative(
	obj interface{},
	vars Activation,
	attr Attribute) interface{} {
	for _, elem := range attr.Path() {
		obj = res.getElem(obj, elem.ToValue(vars))
	}
	return obj
}

func (res *defaultResolver) getElem(obj interface{}, elem ref.Val) interface{} {
	switch obj.(type) {
	case map[string]interface{}:
		m := obj.(map[string]interface{})
		key, ok := elem.(types.String)
		if !ok {
			return types.ValOrErr(elem, "no such overload")
		}
		v, found := m[string(key)]
		if !found {
			return types.NewErr("no such key")
		}
		return v
	case map[string]string:
		return res.getMapStrVal(obj, elem)
	case map[string]int:
		return res.getMapIntVal(obj, elem)
	case map[string]int32:
		return res.getMapInt32Val(obj, elem)
	case map[string]int64:
		return res.getMapInt64Val(obj, elem)
	case []interface{}:
		return res.getListIFaceVal(obj, elem)
	case []string:
		return res.getListStrVal(obj, elem)
	case []int:
		return res.getListIntVal(obj, elem)
	case []int32:
		return res.getListInt32Val(obj, elem)
	case []int64:
		return res.getListInt64Val(obj, elem)
	case proto.Message:
		return res.getProtoField(obj, elem)
	case traits.Indexer:
		indexer := obj.(traits.Indexer)
		return indexer.Get(elem)
	case ref.Val:
		return types.ValOrErr(obj.(ref.Val), "no such overload")
	default:
		objType := reflect.TypeOf(obj)
		objKind := objType.Kind()
		if objKind == reflect.Map ||
			objKind == reflect.Array ||
			objKind == reflect.Slice {
			val := res.adapter.NativeToValue(obj)
			indexer, ok := val.(traits.Indexer)
			if ok {
				return indexer.Get(elem)
			}
			return types.ValOrErr(val, "no such overload")
		}
		return types.NewErr("no such overload")
	}
}

func (res *defaultResolver) getMapStrVal(obj interface{}, k ref.Val) interface{} {
	m := obj.(map[string]string)
	key, ok := k.(types.String)
	if !ok {
		return types.ValOrErr(k, "no such overload")
	}
	v, found := m[string(key)]
	if !found {
		return types.NewErr("no such key")
	}
	return v
}

func (res *defaultResolver) getMapIntVal(obj interface{}, k ref.Val) interface{} {
	m := obj.(map[string]int)
	key, ok := k.(types.String)
	if !ok {
		types.ValOrErr(k, "no such overload")
	}
	v, found := m[string(key)]
	if !found {
		return types.NewErr("no such key")
	}
	return v
}

func (res *defaultResolver) getMapInt32Val(obj interface{}, k ref.Val) interface{} {
	m := obj.(map[string]int32)
	key, ok := k.(types.String)
	if !ok {
		types.ValOrErr(k, "no such overload")
	}
	v, found := m[string(key)]
	if !found {
		return types.NewErr("no such key")
	}
	return v
}

func (res *defaultResolver) getMapInt64Val(obj interface{}, k ref.Val) interface{} {
	m := obj.(map[string]int64)
	key, ok := k.(types.String)
	if !ok {
		types.ValOrErr(k, "no such overload")
	}
	v, found := m[string(key)]
	if !found {
		return types.NewErr("no such key")
	}
	return v
}

func (res *defaultResolver) getListIFaceVal(obj interface{}, i ref.Val) interface{} {
	l := obj.([]interface{})
	idx, ok := i.(types.Int)
	if !ok {
		return types.ValOrErr(i, "no such overload")
	}
	index := int(idx)
	if index < 0 || index >= len(l) {
		return types.ValOrErr(idx, "index out of range")
	}
	return l[index]
}

func (res *defaultResolver) getListStrVal(obj interface{}, i ref.Val) interface{} {
	l := obj.([]string)
	idx, ok := i.(types.Int)
	if !ok {
		return types.ValOrErr(i, "no such overload")
	}
	index := int(idx)
	if index < 0 || index >= len(l) {
		return types.ValOrErr(idx, "index out of range")
	}
	return types.String(l[index])
}

func (res *defaultResolver) getListIntVal(obj interface{}, i ref.Val) interface{} {
	l := obj.([]int)
	idx, ok := i.(types.Int)
	if !ok {
		return types.ValOrErr(i, "no such overload")
	}
	index := int(idx)
	if index < 0 || index >= len(l) {
		return types.ValOrErr(idx, "index out of range")
	}
	return types.Int(l[index])
}

func (res *defaultResolver) getListInt32Val(obj interface{}, i ref.Val) interface{} {
	l := obj.([]int32)
	idx, ok := i.(types.Int)
	if !ok {
		return types.ValOrErr(i, "no such overload")
	}
	index := int(idx)
	if index < 0 || index >= len(l) {
		return types.ValOrErr(idx, "index out of range")
	}
	return types.Int(l[index])
}

func (res *defaultResolver) getListInt64Val(obj interface{}, i ref.Val) interface{} {
	l := obj.([]int64)
	idx, ok := i.(types.Int)
	if !ok {
		return types.ValOrErr(i, "no such overload")
	}
	index := int(idx)
	if index < 0 || index >= len(l) {
		return types.ValOrErr(idx, "index out of range")
	}
	return types.Int(l[index])
}

func (res *defaultResolver) getStructField(m *structpb.Struct, k ref.Val) interface{} {
	fields := m.GetFields()
	key, ok := k.(types.String)
	if !ok {
		return types.ValOrErr(k, "no such overload")
	}
	value, found := fields[string(key)]
	if !found {
		return types.NewErr("no such key")
	}
	return value
}

func (res *defaultResolver) getListValueElem(l *structpb.ListValue, i ref.Val) interface{} {
	idx, ok := i.(types.Int)
	if !ok {
		return types.ValOrErr(i, "no such overload")
	}
	index := int(idx)
	elems := l.GetValues()
	if index < 0 || index >= len(elems) {
		return types.NewErr("index out of range")
	}
	return elems[index]
}

func (res *defaultResolver) getProtoField(obj interface{}, elem ref.Val) interface{} {
	switch obj.(type) {
	case *anypb.Any:
		val := obj.(*anypb.Any)
		if val == nil {
			return types.NewErr("unsupported type conversion: '%T'", obj)
		}
		unpackedAny := ptypes.DynamicAny{}
		if ptypes.UnmarshalAny(val, &unpackedAny) != nil {
			return types.NewErr("unknown type: '%s'", val.GetTypeUrl())
		}
		return res.getProtoField(unpackedAny, elem)
	case *structpb.Value:
		val := obj.(*structpb.Value)
		if val == nil {
			return types.NewErr("no such overload")
		}
		switch val.Kind.(type) {
		case *structpb.Value_StructValue:
			return res.getProtoField(val.GetStructValue(), elem)
		case *structpb.Value_ListValue:
			return res.getProtoField(val.GetListValue(), elem)
		default:
			return types.NewErr("no such overload")
		}
	case *structpb.Struct:
		return res.getStructField(obj.(*structpb.Struct), elem)
	case *structpb.ListValue:
		return res.getListValueElem(obj.(*structpb.ListValue), elem)
	default:
		pb := res.adapter.NativeToValue(obj)
		indexer, ok := pb.(traits.Indexer)
		if !ok {
			return types.ValOrErr(pb, "no such overload")
		}
		return indexer.Get(elem)
	}
}

type unknownResolver struct {
	*defaultResolver
	unknowns map[string][]Attribute
}

func (res *unknownResolver) Resolve(vars Activation, attr Attribute) (interface{}, bool) {
	v, found := res.defaultResolver.Resolve(vars, attr)
	if found {
		val, ok := v.(ref.Val)
		if !ok || !types.IsError(val) {
			return val, true
		}
	}
	varName := attr.Variable().Name()
	varPath := attr.Path()
	candUnknowns, found := res.unknowns[varName]
	for _, cand := range candUnknowns {
		candPath := cand.Path()
		if len(candPath) > len(varPath) {
			continue
		}
		isMatch := true
		for j, candElem := range candPath {
			candElemVal := candElem.ToValue(vars)
			varElem := varPath[j]
			varElemVal := varElem.ToValue(vars)
			if candElemVal.Equal(varElemVal) != types.True {
				isMatch = false
				break
			}
		}
		if isMatch {
			return types.Unknown{attr.Variable().ID()}, true
		}
	}
	return nil, false
}
