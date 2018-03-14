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

package adapters

import (
	"reflect"
	"testing"
)

func TestListAdapter_Get(t *testing.T) {
	list := NewListAdapter([]int32{1, 2, 3})
	if getElem(t, list, 0) != int64(1) ||
		getElem(t, list, 1) != int64(2) ||
		getElem(t, list, 2) != int64(3) {
		t.Errorf("List values by index did not match expectations")
	}
	if _, err := list.Get(-1); err == nil {
		t.Errorf("Should not have been able to read a negative index")
	}
	if _, err := list.Get(3); err == nil {
		t.Errorf("Should not have been able to read beyond end of list")
	}
}

func TestListAdapter_Iterator(t *testing.T) {
	list := NewListAdapter([]int32{1, 2, 3})
	it := list.Iterator()
	var i = int64(0)
	for ; it.HasNext(); i++ {
		elem := it.Next()
		if getElem(t, list, i) != elem {
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

func TestConcatListAdapter_Get(t *testing.T) {
	listA := NewListAdapter([]float32{1.0, 2.0})
	listB := NewListAdapter([]float64{3.0})
	list := listA.Concat(listB)
	if getElem(t, list, 0) != float64(1.0) ||
		getElem(t, list, 1) != float64(2.0) ||
		getElem(t, list, 2) != float64(3.0) {
		t.Errorf("List values by index did not match expectations")
	}
	if _, err := list.Get(-1); err == nil {
		t.Errorf("Should not have been able to read a negative index")
	}
	if _, err := list.Get(3); err == nil {
		t.Errorf("Should not have been able to read beyond end of list")
	}
}

func TestConcatListAdapter_Iterator(t *testing.T) {
	listA := NewListAdapter([]float32{1.0, 2.0})
	listB := NewListAdapter([]float64{3.0})
	list := listA.Concat(listB)
	it := list.Iterator()
	var i = int64(0)
	for ; it.HasNext(); i++ {
		elem := it.Next()
		if getElem(t, list, i) != elem {
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

func TestConcatListAdapter_Equal(t *testing.T) {
	listA := NewListAdapter([]float32{1.0, 2.0})
	listB := NewListAdapter([]float64{3.0})
	listConcat := listA.Concat(listB)
	// Note the internal type of list raw and concat list are slightly different.
	listRaw := NewListAdapter([]interface{}{
		float32(1.0), float64(2.0), float64(3.0)})
	if !listRaw.Equal(listConcat) || !listConcat.Equal(listRaw) {
		t.Errorf("Concat list and raw list were not equal")
	}
	if listRaw.Equal(listA) {
		t.Errorf("Lists of unequal length considered equal")
	}
}

func TestListAdapter_Contains(t *testing.T) {
	listA := NewListAdapter([]float32{1.0, 2.0})
	listB := NewListAdapter([]float64{3.0})
	list := listA.Concat(listB).Concat(listA)
	if list.Contains(float32(3.0)) {
		t.Error("List contains succeeded with wrong type")
	}
	if !list.Contains(float64(3)) {
		t.Error("List contains did not succeed")
	}
}

func TestListAdapter_ToProto(t *testing.T) {
	listA := NewListAdapter([]float32{1.0, 2.0})
	listB := NewListAdapter([]float64{3.0})
	listConcat := listA.Concat(listB)
	if protoList, err := listConcat.ToProto(reflect.TypeOf([]float32{})); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(protoList, []float32{1.0, 2.0, 3.0}) {
		t.Errorf("Could not convert to []float32: %v", protoList)
	}
}

func TestListAdapter_NestedList(t *testing.T) {
	listUint32 := []uint32{1, 2}
	nestedUint32 := NewListAdapter([]interface{}{listUint32})
	listUint64 := []uint64{1, 2}
	nestedUint64 := NewListAdapter([]interface{}{listUint64})
	if !nestedUint32.Equal(nestedUint64) {
		t.Error("Could not find nested list")
	}
	if !nestedUint32.Contains(NewListAdapter(listUint64)) ||
		!nestedUint64.Contains(NewListAdapter(listUint32)) {
		t.Error("Could not find type compatible nested lists")
	}
}

func getElem(t *testing.T, list ListAdapter, index int64) interface{} {
	t.Helper()
	if val, err := list.Get(index); err != nil {
		t.Errorf("Error reading list index %d, %v", index, err)
		return nil
	} else {
		return val
	}
}
