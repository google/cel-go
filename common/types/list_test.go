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
	"github.com/google/cel-go/common/types/traits"
	"reflect"
	"testing"
)

func TestListValue_Get(t *testing.T) {
	list := NewDynamicList([]int32{1, 2, 3}).(traits.Lister)
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
}

func TestListValue_Iterator(t *testing.T) {
	list := NewDynamicList([]int32{1, 2, 3}).(traits.Lister)
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

func TestConcatListValue_Get(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0}).(traits.Lister)
	listB := NewDynamicList([]float64{3.0}).(traits.Lister)
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
	listA := NewDynamicList([]float32{1.0, 2.0}).(traits.Lister)
	listB := NewDynamicList([]float64{3.0}).(traits.Lister)
	list := listA.Add(listB).(traits.Lister)
	it := list.Iterator()
	var i = int64(0)
	for ; it.HasNext() == True; i++ {
		elem := it.Next()
		fmt.Printf("it[%d] elem %v\n", i, elem)
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

func TestConcatListValue_Equal(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0}).(traits.Lister)
	listB := NewDynamicList([]float64{3.0})
	listConcat := listA.Add(listB)
	// Note the internal type of list raw and concat list are slightly different.
	listRaw := NewDynamicList([]interface{}{
		float32(1.0), float64(2.0), float64(3.0)})
	if listRaw.Equal(listConcat) != True ||
		listConcat.Equal(listRaw) != True {
		t.Errorf("Concat list and raw list were not equal, got '%v', expected '%v'",
			listConcat.Value(),
			listRaw.Value())
	}
	if listRaw.Equal(listA) == True {
		t.Errorf("Lists of unequal length considered equal")
	}
}

func TestListValue_Contains(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0}).(traits.Lister)
	listB := NewDynamicList([]float64{3.0})
	list := listA.Add(listB).(traits.Lister).Add(listA).(traits.Lister)
	if list.Contains(Int(3)) == True {
		t.Error("List contains succeeded with wrong type")
	}
	if list.Contains(Double(3)) != True {
		t.Error("List contains did not succeed")
	}
}

func TestListValue_ToProto(t *testing.T) {
	listA := NewDynamicList([]float32{1.0, 2.0}).(traits.Lister)
	listB := NewDynamicList([]float64{3.0})
	listConcat := listA.Add(listB)
	if protoList, err := listConcat.ConvertToNative(reflect.TypeOf([]float32{})); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(protoList, []float32{1.0, 2.0, 3.0}) {
		t.Errorf("Could not convert to []float32: %v", protoList)
	}
}

func TestListValue_NestedList(t *testing.T) {
	listUint32 := []uint32{1, 2}
	nestedUint32 := NewDynamicList([]interface{}{listUint32}).(traits.Lister)
	listUint64 := []uint64{1, 2}
	nestedUint64 := NewDynamicList([]interface{}{listUint64}).(traits.Lister)
	if nestedUint32.Equal(nestedUint64) != True {
		t.Error("Could not find nested list")
	}
	if nestedUint32.Contains(NewDynamicList(listUint64)) != True ||
		nestedUint64.Contains(NewDynamicList(listUint32)) != True {
		t.Error("Could not find type compatible nested lists")
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
