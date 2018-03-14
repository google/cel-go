package adapters

import (
	"github.com/google/cel-go/interpreter/types/objects"
	"testing"
)

func TestMapAdapter_Equal(t *testing.T) {
	mapValue := NewMapAdapter(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}})
	if !mapValue.Equal(mapValue) {
		t.Error("Map value was not equal to itself")
	}
	if nestedVal, err := mapValue.Get("nested"); err != nil {
		t.Error(err)
	} else if mapValue.Equal(nestedVal) || nestedVal.(MapAdapter).Equal(mapValue) {
		t.Error("Same length, but different key names")
	}
}

func TestMapAdapter_Get(t *testing.T) {
	mapValue := NewMapAdapter(map[string]map[int32]float32{
		"nested": {1: -1.0, 2: 2.0},
		"empty":  {}})
	if nestedVal, err := mapValue.Get("nested"); err != nil {
		t.Error(err)
	} else if floatVal, err := nestedVal.(objects.Indexer).Get(int64(1)); err != nil {
		t.Error(err)
	} else if floatVal != float64(-1.0) {
		t.Error("Nested map access of float property not float64")
	}
}

func TestMapAdapter_Iterator(t *testing.T) {
	mapValue := NewMapAdapter(map[string]map[int32]float32{
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
