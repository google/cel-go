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
	"github.com/google/cel-go/interpreter/types/aspects"
	"testing"
)

func TestMapValue_Equal(t *testing.T) {
	mapValue := NewMapValue(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}})
	if !mapValue.Equal(mapValue) {
		t.Error("Map value was not equal to itself")
	}
	if nestedVal, err := mapValue.Get("nested"); err != nil {
		t.Error(err)
	} else if mapValue.Equal(nestedVal) || nestedVal.(MapValue).Equal(mapValue) {
		t.Error("Same length, but different key names")
	}
}

func TestMapValue_Get(t *testing.T) {
	mapValue := NewMapValue(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}})
	if nestedVal, err := mapValue.Get("nested"); err != nil {
		t.Error(err)
	} else if floatVal, err := nestedVal.(aspects.Indexer).Get(int64(1)); err != nil {
		t.Error(err)
	} else if floatVal != float64(-1.0) {
		t.Error("Nested map access of float property not float64")
	}
}

func TestMapValue_Iterator(t *testing.T) {
	mapValue := NewMapValue(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}})
	it := mapValue.Iterator()
	var i = 0
	var fieldNames []interface{}
	for ; it.HasNext(); i++ {
		if value, err := mapValue.Get(it.Next()); err != nil {
			t.Error(err)
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
