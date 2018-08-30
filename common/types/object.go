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
	"reflect"

	protopb "github.com/golang/protobuf/proto"
	ptypespb "github.com/golang/protobuf/ptypes"
	traitspb "github.com/google/cel-go/common/types/traits"
	refpb "github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/pb"
)

type protoObj struct {
	value     protopb.Message
	refValue  reflect.Value
	typeDesc  *pb.TypeDescription
	typeValue *TypeValue
	isAny     bool
}

// NewObject returns an object based on a protopb.Message value which handles
// conversion between protobuf type values and expression type values.
// Objects support indexing and iteration.
func NewObject(value protopb.Message) refpb.Value {
	typeDesc, err := pb.DescribeValue(value)
	if err != nil {
		panic(err)
	}
	return &protoObj{
		value:     value,
		refValue:  reflect.ValueOf(value),
		typeDesc:  typeDesc,
		typeValue: NewObjectTypeValue(typeDesc.Name())}
}

func (o *protoObj) ConvertToNative(typeDesc reflect.Type) (interface{}, error) {
	if typeDesc.AssignableTo(o.refValue.Type()) {
		return o.value, nil
	}
	if typeDesc == anyValueType {
		return ptypespb.MarshalAny(o.Value().(protopb.Message))
	}
	// If the object is already assignable to the desired type return it.
	if reflect.TypeOf(o).AssignableTo(typeDesc) {
		return o, nil
	}
	return nil, fmt.Errorf("type conversion error from '%v' to '%v'",
		o.refValue.Type(), typeDesc)
}

func (o *protoObj) ConvertToType(typeVal refpb.Type) refpb.Value {
	switch typeVal {
	default:
		if o.Type().TypeName() == typeVal.TypeName() {
			return o
		}
	case TypeType:
		return o.typeValue
	}
	return NewErr("type conversion error from '%s' to '%s'",
		o.typeDesc.Name(), typeVal)
}

func (o *protoObj) Equal(other refpb.Value) refpb.Value {
	if o.typeDesc.Name() != other.Type().TypeName() {
		return False
	}
	return Bool(protopb.Equal(o.value, other.Value().(protopb.Message)))
}

func (o *protoObj) Get(index refpb.Value) refpb.Value {
	if index.Type() != StringType {
		return NewErr("illegal object field type '%s'", index.Type())
	}
	protoFieldName := string(index.(String))
	if f, found := o.typeDesc.FieldByName(protoFieldName); found {
		if !f.IsOneof() {
			return getOrDefaultInstance(o.refValue.Elem().Field(f.Index()))
		}

		getter := o.refValue.MethodByName(f.GetterName())
		if getter.IsValid() {
			refField := getter.Call([]reflect.Value{})[0]
			if refField.IsValid() {
				return getOrDefaultInstance(refField)
			}
		}
	}
	return NewErr("no such field '%s'", index)
}

func (o *protoObj) Iterator() traitspb.Iterator {
	return &msgIterator{
		baseIterator: &baseIterator{},
		refValue:     o.refValue,
		typeDesc:     o.typeDesc,
		cursor:       0}
}

func (o *protoObj) Type() refpb.Type {
	return o.typeValue
}

func (o *protoObj) Value() interface{} {
	return o.value
}

type msgIterator struct {
	*baseIterator
	refValue reflect.Value
	typeDesc *pb.TypeDescription
	cursor   int
	len      int
}

func (it *msgIterator) HasNext() refpb.Value {
	return Bool(it.cursor < it.typeDesc.FieldCount())
}

func (it *msgIterator) Next() refpb.Value {
	if it.HasNext() == False {
		return nil
	}
	fieldName, _ := it.typeDesc.FieldNameAtIndex(it.cursor, it.refValue)
	it.cursor += 1
	return String(fieldName)
}

var (
	protoDefaultInstanceMap = make(map[reflect.Type]refpb.Value)
)

func getOrDefaultInstance(refVal reflect.Value) refpb.Value {
	value := refVal.Interface()
	if refVal.Kind() != reflect.Ptr || !refVal.IsNil() {
		return NativeToValue(value)
	}
	return getDefaultInstance(refVal.Type())
}

func getDefaultInstance(refType reflect.Type) refpb.Value {
	if refType.Kind() == reflect.Ptr {
		refType = refType.Elem()
	}
	if defaultValue, found := protoDefaultInstanceMap[refType]; found {
		return defaultValue
	}
	defaultValue := NativeToValue(reflect.New(refType).Interface())
	protoDefaultInstanceMap[refType] = defaultValue
	return defaultValue
}
