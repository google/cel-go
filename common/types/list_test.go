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
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/google/cel-go/common/types/traits"
	"reflect"
	"testing"
)

func TestBaseList_Add_Empty(t *testing.T) {
	list := NewDynamicList([]bool{true})
	if list.Add(NewDynamicList([]bool{})) != list {
		t.Error("Adding an empty list created new list.")
	}
	if NewDynamicList([]string{}).Add(list) != list {
		t.Error("Adding list to empty created a new list.")
	}
}

func TestBaseList_Add_Error(t *testing.T) {
	if !IsError(NewDynamicList([]bool{}).Add(String("error"))) {
		t.Error("Addind a non-list value to a list unexpected succeeds.")
	}
}

func TestBaseList_Contains(t *testing.T) {
	list := NewDynamicList([]float32{1.0, 2.0, 3.0})
	if list.Contains(Int(3)) == True {
		t.Error("List contains succeeded with wrong type")
	}
	if list.Contains(Double(3)) != True {
		t.Error("List contains did not succeed")
	}
}

func TestBaseList_ConvertToNative(t *testing.T) {
	list := NewDynamicList([]float64{1.0, 2.0})
	if protoList, err := list.ConvertToNative(reflect.TypeOf([]float32{})); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(protoList, []float32{1.0, 2.0}) {
		t.Errorf("Could not convert to []float32: %v", protoList)
	}
}

func TestBaseList_ConvertToType(t *testing.T) {
	list := NewDynamicList([]string{"h", "e", "l", "l", "o"})
	if list.ConvertToType(ListType) != list {
		t.Error("List was not convertible to itself.")
	}
	if list.ConvertToType(TypeType) != ListType {
		t.Error("Unable to obtain the proper type from the list.")
	}
	if !IsError(list.ConvertToType(MapType)) {
		t.Error("List was able to convert to unexpected type.")
	}
}

func TestBaseList_Equal(t *testing.T) {
	listA := NewDynamicList([]string{"h", "e", "l", "l", "o"})
	listB := NewDynamicList([]string{"h", "e", "l", "p", "!"})
	if listA.Equal(listB) != False {
		t.Error("Lists with different contents returned equal.")
	}
}

func TestBaseList_Get(t *testing.T) {
	list := NewDynamicList([]int32{1, 2, 3})
	if getElem(t, list, 0) != Int(1) ||
		getElem(t, list, 1) != Int(2) ||
		getElem(t, list, 2) != Int(3) {
		t.Errorf("List values by index did not match expectations")
	}
	if val := list.Get(Int(-1)); !IsError(val) {
		t.Errorf("Should not have been able to read a negative index")
	}
	if val := list.Get(Int(3)); !IsError(val) {
		t.Errorf("Should not have been able to read beyond end of list")
	}
	if !IsError(list.Get(Uint(3))) {
		t.Error("Invalid index type did not result in error")
	}
}

func TestBaseList_Iterator(t *testing.T) {
	list := NewDynamicList([]int32{1, 2, 3})
	it := list.Iterator()
	var i = int64(0)
	for ; it.HasNext() == True; i++ {
		elem := it.Next()
		if getElem(t, list, Int(i)) != elem {
			t.Errorf(
				"List iterator returned incorrect value: list[%d]: %v", i, elem)
		}
	}
	if it.Next() != nil {
		t.Errorf("List iterator attempted to continue beyond list size")
	}
	if i != 3 {
		t.Errorf("Iterator did not iterate until last value")
	}
}

func TestBaseList_NestedList(t *testing.T) {
	listUint32 := []uint32{1, 2}
	nestedUint32 := NewDynamicList([]interface{}{listUint32})
	listUint64 := []uint64{1, 2}
	nestedUint64 := NewDynamicList([]interface{}{listUint64})
	if nestedUint32.Equal(nestedUint64) != True {
		t.Error("Could not find nested list")
	}
	if nestedUint32.Contains(NewDynamicList(listUint64)) != True ||
		nestedUint64.Contains(NewDynamicList(listUint32)) != True {
		t.Error("Could not find type compatible nested lists")
	}
}

func TestBaseList_Size(t *testing.T) {
	listUint32 := []uint32{1, 2}
	nestedUint32 := NewDynamicList([]interface{}{listUint32})
	if nestedUint32.Size() != IntOne {
		t.Error("List indicates the incorrect size.")
	}
	if nestedUint32.Get(IntZero).(traits.Sizer).Size() != Int(2) {
		t.Error("Nested list indicates the incorrect size.")
	}
}

func TestConcatList_Add(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	listB := NewStringList([]string{"3"})
	list := listA.Add(listB).(traits.Lister).Add(listA).
		Value().([]interface{})
	expected := []interface{}{
		float32(1.0),
		float32(2.0),
		string("3"),
		float32(1.0),
		float32(2.0)}
	if len(list) != len(expected) {
		t.Errorf("Got '%v', expected '%v'", list, expected)
	} else {
		for i := 0; i < len(list); i++ {
			if expected[i] != list[i] {
				t.Errorf("elem[%d] Got '%v', expected '%v'",
					i, list[i], expected[i])
			}
		}
	}
}

func TestConcatList_ConvertToNative_Json(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	listB := NewDynamicList([]string{"3"})
	list := listA.Add(listB)
	json, err := list.ConvertToNative(jsonValueType)
	if err != nil || IsError(json) {
		t.Errorf("Got '%v', expected '%v'", err, json)
	}
	jsonTxt, err := (&jsonpb.Marshaler{}).MarshalToString(json.(proto.Message))
	if err != nil {
		t.Error(err)
	}
	if jsonTxt != "[1,2,\"3\"]" {
		t.Errorf("Got '%v', expected [1,2,\"3\"]", jsonTxt)
	}
}

func TestConcatList_ConvertToNative_ElementConversionError(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	// Duration is serializable to a string form of json, but there is no
	// concept of a duration literal within CEL, so the serialization to string
	// is not supported here which should cause the conversion to json to fail.
	listB := NewDynamicList([]*duration.Duration{{Seconds: 100}})
	listConcat := listA.Add(listB)
	json, err := listConcat.ConvertToNative(jsonValueType)
	if err == nil {
		t.Errorf("Got '%v', expected error", json)
	}
}

func TestConcatList_ConvertToType(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	listB := NewDynamicList([]*duration.Duration{{Seconds: 100}})
	list := listA.Add(listB)
	if list.ConvertToType(ListType) != list {
		t.Error("List conversion to list failed.")
	}
	if list.ConvertToType(TypeType) != ListType {
		t.Error("List conversion to type failed.")
	}
	if !IsError(list.ConvertToType(MapType)) {
		t.Error("List conversion to map unexpectedly succeeded.")
	}
}

func TestConcatListValue_Equal(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	listB := NewDynamicList([]float64{3.0})
	list := listA.Add(listB)
	// Note the internal type of list raw and concat list are slightly different.
	listRaw := NewDynamicList([]interface{}{
		float32(1.0), float64(2.0), float64(3.0)})
	if listRaw.Equal(list) != True ||
		list.Equal(listRaw) != True {
		t.Errorf("Concat list and raw list were not equal, got '%v', expected '%v'",
			list.Value(),
			listRaw.Value())
	}
	if list.Equal(listA) == True ||
		listRaw.Equal(listA) == True {
		t.Errorf("Lists of unequal length considered equal")
	}
}

func TestConcatListValue_Get(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	listB := NewDynamicList([]float64{3.0})
	list := listA.Add(listB).(traits.Lister)
	if getElem(t, list, 0) != Double(1.0) ||
		getElem(t, list, 1) != Double(2.0) ||
		getElem(t, list, 2) != Double(3.0) {
		t.Errorf("List values by index did not match expectations")
	}
	if val := list.Get(Int(-1)); !IsError(val) {
		t.Errorf("Should not have been able to read a negative index")
	}
	if val := list.Get(Int(3)); !IsError(val) {
		t.Errorf("Should not have been able to read beyond end of list")
	}
}

func TestConcatListValue_Iterator(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0})
	listB := NewDynamicList([]float64{3.0})
	list := listA.Add(listB).(traits.Lister)
	it := list.Iterator()
	var i = int64(0)
	for ; it.HasNext() == True; i++ {
		elem := it.Next()
		if getElem(t, list, Int(i)) != elem {
			t.Errorf(
				"List iterator returned incorrect value: list[%d]: %v", i, elem)
		}
	}
	if it.Next() != nil {
		t.Errorf("List iterator attempted to continue beyond list size")
	}
	if i != 3 {
		t.Errorf("Iterator did not iterate until last value")
	}
}

func TestStringList_Add_Empty(t *testing.T) {
	list := NewStringList([]string{"hello"})
	if list.Add(NewStringList([]string{})) != list {
		t.Error("Adding empty lists resulted in new list creation.")
	}
	if NewStringList([]string{}).Add(list) != list {
		t.Error("Adding empty lists resulted in new list creation.")
	}
}

func TestStringList_Add_Error(t *testing.T) {

}

func TestStringList_Add_Heterogenous(t *testing.T) {
	listA := NewStringList([]string{"hello"})
	listB := NewDynamicList([]int32{1, 2, 3})
	list := listA.Add(listB).(traits.Lister).Value().([]interface{})
	expected := []interface{}{"hello", int32(1), int32(2), int32(3)}
	if len(list) != len(expected) {
		t.Error("Unexpected list size. Got '%d', expected 4", len(list))
	}
	for i, v := range expected {
		if list[i] != v {
			t.Errorf("elem[%d] Got '%v', expected '%v'", i, list[i], v)
		}
	}
}

func TestStringList_Add_StringLists(t *testing.T) {
	listA := NewStringList([]string{"hello"})
	listB := NewStringList([]string{"world", "!"})
	list := listA.Add(listB).(traits.Lister)
	if list.Size() != Int(3) {
		t.Error("Combined list did not have correct size.")
	}
	expected := []string{"hello", "world", "!"}
	for i, v := range expected {
		if list.Get(Int(i)).Equal(String(v)) != True {
			t.Error("elem[%d] Got '%v', expected '%v'", i, list.Get(Int(i)), v)
		}
	}
}

func TestStringList_ConvertToNative(t *testing.T) {
	list := NewStringList([]string{"h", "e", "l", "p"})
	_, err := list.ConvertToNative(reflect.TypeOf([]string{}))
	if err != nil {
		t.Error("Unable to convert string list to itself.")
	}
}

func TestStringList_ConvertToNative_Error(t *testing.T) {
	list := NewStringList([]string{"h", "e", "l", "p"})
	_, err := list.ConvertToNative(jsonStructType)
	if err == nil {
		t.Error("Conversion of list to unsupported type did not error.")
	}
}

func TestStringList_ConvertToNative_Json(t *testing.T) {
	list := NewStringList([]string{"h", "e", "l", "p"})
	json, err := list.ConvertToNative(jsonValueType)
	if err != nil || IsError(json) {
		t.Errorf("Got '%v', expected '%v'", err, json)
	}
	jsonTxt, err := (&jsonpb.Marshaler{}).MarshalToString(json.(proto.Message))
	if err != nil {
		t.Error(err)
	}
	if jsonTxt != "[\"h\",\"e\",\"l\",\"p\"]" {
		t.Errorf("Got '%v', expected [\"h\",\"e\",\"l\",\"p\"]", jsonTxt)
	}

	jsonList, err := list.ConvertToNative(jsonListValueType)
	if err != nil || IsError(jsonList) {
		t.Errorf("Got '%v', expected '%v'", err, jsonList)
	}
	jsonListTxt, err := (&jsonpb.Marshaler{}).MarshalToString(jsonList.(proto.Message))
	if jsonTxt != jsonListTxt {
		t.Errorf("Json value and list value not equal.")
	}
}

func TestStringList_Get_OutOfRange(t *testing.T) {
	list := NewStringList([]string{"hello", "world"})
	if !IsError(list.Get(Int(-1))) {
		t.Error("Negative index did not return error.")
	}
	if !IsError(list.Get(Int(2))) {
		t.Error("Index out of range did not return error.")
	}
	if !IsError(list.Get(Uint(1))) {
		t.Error("Invalid index type did not return error.")
	}
}

func getElem(t *testing.T, list traits.Indexer, index Int) interface{} {
	t.Helper()
	if val := list.Get(index); IsError(val) {
		t.Errorf("Error reading list index %d, %v", index, val)
		return nil
	} else {
		return val
	}
}
