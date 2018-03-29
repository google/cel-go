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

package providers

import (
	"github.com/google/cel-go/interpreter/types"
	expr "github.com/google/cel-spec/proto/v1/syntax"
	"reflect"
	"testing"
)

func TestTypeProvider_NewValue(t *testing.T) {
	typeProvider := NewTypeProvider(&expr.ParsedExpr{})
	if sourceInfo, err := typeProvider.NewValue(
		"google.api.expr.v1.SourceInfo",
		map[string]interface{}{
			"Location":    "TestTypeProvider_NewValue",
			"LineOffsets": []int64{0, 2},
			"Positions":   map[int64]int64{1: 2, 2: 4},
		}); err != nil {
		t.Error(err)
	} else {
		info := sourceInfo.Value().(*expr.SourceInfo)
		if info.Location != "TestTypeProvider_NewValue" ||
			!reflect.DeepEqual(info.LineOffsets, []int32{0, 2}) ||
			!reflect.DeepEqual(info.Positions, map[int64]int32{1: 2, 2: 4}) {
			t.Errorf("Source info not properly configured: %v", info)
		}
	}
}

func TestTypeProvider_Getters(t *testing.T) {
	typeProvider := NewTypeProvider(&expr.ParsedExpr{})
	if sourceInfo, err := typeProvider.NewValue(
		"google.api.expr.v1.SourceInfo",
		map[string]interface{}{
			"Location":    "TestTypeProvider_GetFieldValue",
			"LineOffsets": []int64{0, 2},
			"Positions":   map[int64]int64{1: 2, 2: 4},
		}); err != nil {
		t.Error(err)
	} else {
		if loc, err := sourceInfo.Get("Location"); err != nil {
			t.Error(err)
		} else if loc != "TestTypeProvider_GetFieldValue" {
			t.Errorf("Expected %s, got %s",
				"TestTypeProvider_GetFieldValue",
				loc)
		}
		if pos, err := sourceInfo.Get("Positions"); err != nil {
			t.Error(err)
		} else if !reflect.DeepEqual(pos.(types.MapValue).Value(),
			map[int64]int32{1: 2, 2: 4}) {
			t.Errorf("Expected map[int64]int32, got %v", pos)
		} else if posKeyVal, err := pos.(types.MapValue).Get(int64(1)); err != nil {
			t.Error(err)
		} else if posKeyVal.(int64) != 2 {
			t.Error("Expected value to be int64, not int32")
		}
		if offsets, err := sourceInfo.Get("LineOffsets"); err != nil {
			t.Error(err)
		} else if offset1, err := offsets.(types.ListValue).Get(int64(1)); err != nil {
			t.Error(err)
		} else if offset1.(int64) != 2 {
			t.Errorf("Expected index 1 to be value 2, was %v", offset1)
		}
	}
}

func BenchmarkTypeProvider_NewValue(b *testing.B) {
	typeProvider := NewTypeProvider(&expr.ParsedExpr{})
	for i := 0; i < b.N; i++ {
		typeProvider.NewValue(
			"google.api.expr.v1.SourceInfo",
			map[string]interface{}{
				"Location":    "BenchmarkTypeProvider_NewValue",
				"LineOffsets": []int64{0, 2},
				"Positions":   map[int64]int64{1: 2, 2: 4},
			})
	}
}
