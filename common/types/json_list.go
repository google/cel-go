package types

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
)

var (
	jsonListValueType = reflect.TypeOf(&structpb.ListValue{})
)

type jsonListValue struct {
	*structpb.ListValue
}

// NewJsonList creates a traits.Lister implementation backed by a JSON list
// that has been encoded in protocol buffer form.
func NewJsonList(l *structpb.ListValue) traits.Lister {
	return &jsonListValue{l}
}

func (l *jsonListValue) Add(other ref.Value) ref.Value {
	if other.Type() != ListType {
		return NewErr("no such overload")
	}
	switch other.(type) {
	case *jsonListValue:
		otherList := other.(*jsonListValue)
		concatElems := append(l.GetValues(), otherList.GetValues()...)
		return NewJsonList(&structpb.ListValue{Values: concatElems})
	}
	return &concatList{
		prevList: l,
		nextList: other.(traits.Lister)}
}

func (l *jsonListValue) Contains(elem ref.Value) ref.Value {
	for i := Int(0); i < l.Size().(Int); i++ {
		if l.Get(i).Equal(elem) == True {
			return True
		}
	}
	return False
}

func (l *jsonListValue) ConvertToNative(refType reflect.Type) (interface{}, error) {
	switch refType.Kind() {
	case reflect.Array, reflect.Slice:
		elemCount := int(l.Size().(Int))
		nativeList := reflect.MakeSlice(refType, elemCount, elemCount)
		for i := 0; i < elemCount; i++ {
			elem := l.Get(Int(i))
			nativeElemVal, err := elem.ConvertToNative(refType.Elem())
			if err != nil {
				return nil, err
			}
			nativeList.Index(i).Set(reflect.ValueOf(nativeElemVal))
		}
		return nativeList.Interface(), nil
	case reflect.Ptr:
		if refType == jsonValueType {
			return &structpb.Value{
				Kind: &structpb.Value_ListValue{
					ListValue: l.ListValue}}, nil
		}
		if refType == jsonListValueType {
			return l.ListValue, nil
		}
	}
	return nil, fmt.Errorf("no conversion found from list type to native type."+
		" list elem: google.protobuf.Value, native type: %v", refType)
}

func (l *jsonListValue) ConvertToType(typeVal ref.Type) ref.Value {
	switch typeVal {
	case ListType:
		return l
	case TypeType:
		return ListType
	}
	return NewErr("type conversion error from '%s' to '%s'", ListType, typeVal)
}

func (l *jsonListValue) Equal(other ref.Value) ref.Value {
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

func (l *jsonListValue) Get(index ref.Value) ref.Value {
	if IntType != index.Type() {
		return NewErr("unsupported index type: '%v", index.Type())
	}
	i := index.(Int)
	if i < 0 || i >= l.Size().(Int) {
		return NewErr("index '%d' out of range in list size '%d'", i, l.Size())
	}
	elem := l.GetValues()[i]
	return NativeToValue(elem)
}

func (l *jsonListValue) Iterator() traits.Iterator {
	return &jsonValueListIterator{
		baseIterator: &baseIterator{},
		elems:        l.GetValues(),
		len:          len(l.GetValues())}
}

func (l *jsonListValue) Size() ref.Value {
	return Int(len(l.GetValues()))
}

func (l *jsonListValue) Type() ref.Type {
	return ListType
}

func (l *jsonListValue) Value() interface{} {
	return l.ListValue
}

type jsonValueListIterator struct {
	*baseIterator
	cursor int
	elems  []*structpb.Value
	len    int
}

func (it *jsonValueListIterator) HasNext() ref.Value {
	return Bool(it.cursor < it.len)
}

func (it *jsonValueListIterator) Next() ref.Value {
	if it.HasNext() == True {
		index := it.cursor
		it.cursor += 1
		return NativeToValue(it.elems[index])
	}
	return nil
}
