package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
	"testing"
)

func TestJsonListValue_Add(t *testing.T) {
	listA := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	listB := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_NumberValue{2}},
		{&structpb.Value_NumberValue{3}}}})
	list := listA.Add(listB)
	nativeVal, err := list.ConvertToNative(jsonListValueType)
	if err != nil {
		t.Error(err)
	}
	expected := &structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}},
		{&structpb.Value_NumberValue{2}},
		{&structpb.Value_NumberValue{3}}}}
	if !proto.Equal(nativeVal.(proto.Message), expected) {
		t.Errorf("Concatenated lists did not combine as expected."+
			" Got '%v', expected '%v'", nativeVal, expected)
	}
}

func TestJsonListValue_Contains(t *testing.T) {
	list := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	if !list.Contains(Double(1)).(Bool) {
		t.Error("Expected value list to contain number '1'", list)
	}
	if list.Contains(Double(2)).(Bool) {
		t.Error("Expected value list to not contain number '2'", list)
	}
}

func TestJsonListValue_ConvertToNative_Json(t *testing.T) {
	list := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	listVal, err := list.ConvertToNative(jsonListValueType)
	if err != nil {
		t.Error(err)
	}
	if listVal != list.Value().(proto.Message) {
		t.Error("List did not convert to its underlying representation.")
	}

	val, err := list.ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message),
		&structpb.Value{Kind: &structpb.Value_ListValue{
			ListValue: listVal.(*structpb.ListValue)}}) {
		t.Errorf("Messages were not equal, got '%v'", val)
	}
}

func TestJsonListValue_ConvertToNative_Slice(t *testing.T) {
	list := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	listVal, err := list.ConvertToNative(reflect.TypeOf([]*structpb.Value{}))
	if err != nil {
		t.Error(err)
	}
	for i, v := range listVal.([]*structpb.Value) {
		if !list.Get(Int(i)).Equal(NativeToValue(v)).(Bool) {
			t.Errorf("elem[%d] Got '%v', expected '%v'",
				i, v, list.Get(Int(i)))
		}
	}
}

func TestJsonListValue_ConvertToType(t *testing.T) {
	list := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	if list.ConvertToType(TypeType) != ListType {
		t.Error("Json list type was not a list.")
	}
	if list.ConvertToType(ListType) != list {
		t.Error("Json list not convertible to itself.")
	}
	if !IsError(list.ConvertToType(MapType)) {
		t.Error("Json list converted to an unsupported type.")
	}
}

func TestJsonListValue_Equal(t *testing.T) {
	listA := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	listB := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_NumberValue{2}},
		{&structpb.Value_NumberValue{3}}}})
	if listA.Equal(listB).(Bool) || listB.Equal(listA).(Bool) {
		t.Error("Lists with different elements considered equal.")
	}
	if !listA.Equal(listA).(Bool) {
		t.Error("List was not equal to itself.")
	}
	if listA.Add(listA).Equal(listB).(Bool) {
		t.Error("Lists of different size were equal.")
	}
	if listA.Equal(True).(Bool) {
		t.Error("Equality of different type returned true.")
	}
}

func TestJsonListValue_Get_OutOfRange(t *testing.T) {
	list := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}}}})
	if !IsError(list.Get(Int(-1))) {
		t.Error("Negative index did not result in error.")
	}
	if !IsError(list.Get(Int(2))) {
		t.Error("Index out of range did not result in error.")
	}
	if !IsError(list.Get(Uint(1))) {
		t.Error("Index of incorrect type did not result in error.")
	}
}

func TestJsonListValue_Iterator(t *testing.T) {
	list := NewJsonList(&structpb.ListValue{[]*structpb.Value{
		{&structpb.Value_StringValue{"hello"}},
		{&structpb.Value_NumberValue{1}},
		{&structpb.Value_NumberValue{2}},
		{&structpb.Value_NumberValue{3}}}})
	it := list.Iterator()
	for i := Int(0); it.HasNext() != False; i++ {
		v := it.Next()
		if v.Equal(list.Get(i)) != True {
			t.Errorf("elem[%d] Got '%v', expected '%v'", v, list.Get(i))
		}
	}

	if it.HasNext() != False {
		t.Error("Iterator indicated more elements were left")
	}
	if it.Next() != nil {
		t.Error("Calling Next() for a complete iterator resulted in a non-nil value.")
	}
}
