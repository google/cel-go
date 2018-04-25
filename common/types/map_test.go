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
	"github.com/google/cel-go/common/types/traits"
	"testing"
)

func TestMapValue_Equal(t *testing.T) {
	mapValue := NewDynamicMap(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}}).(traits.Mapper)
	if mapValue.Equal(mapValue) != True {
		t.Error("Map value was not equal to itself")
	}
	if nestedVal := mapValue.Get(String("nested")); IsError(nestedVal) {
		t.Error(nestedVal)
	} else if mapValue.Equal(nestedVal) == True ||
		nestedVal.Equal(mapValue) == True {
		t.Error("Same length, but different key names")
	}
}

func TestMapValue_Get(t *testing.T) {
	mapValue := NewDynamicMap(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}}).(traits.Mapper)
	if nestedVal := mapValue.Get(String("nested")); IsError(nestedVal) {
		t.Error(nestedVal)
	} else if floatVal := nestedVal.(traits.Indexer).Get(Int(1)); IsError(floatVal) {
		t.Error(floatVal)
	} else if floatVal.Equal(Double(-1.0)) != True {
		t.Error("Nested map access of float property not float64")
	}
}

func TestMapValue_Iterator(t *testing.T) {
	mapValue := NewDynamicMap(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}}).(traits.Mapper)
	it := mapValue.Iterator()
	var i = 0
	var fieldNames []interface{}
	for ; it.HasNext() == True; i++ {
		if value := mapValue.Get(it.Next()); IsError(value) {
			t.Error(value)
		} else {
			fieldNames = append(fieldNames, value)
		}
	}
	if len(fieldNames) != 2 {
		t.Errorf("Did not find the correct number of fields: %v", fieldNames)
	}
	if it.Next() != nil {
		t.Error("Iterator ran off the end of the field names")
	}
}
