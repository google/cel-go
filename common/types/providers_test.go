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
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	expr "github.com/google/cel-spec/proto/v1/syntax"
	"reflect"
	"testing"
)

func TestTypeProvider_NewValue(t *testing.T) {
	typeProvider := NewProvider(&expr.ParsedExpr{})
	if sourceInfo := typeProvider.NewValue(
		"google.api.expr.v1.SourceInfo",
		map[string]ref.Value{
			"Location":    String("TestTypeProvider_NewValue"),
			"LineOffsets": NewDynamicList([]int64{0, 2}),
			"Positions":   NewDynamicMap(map[int64]int64{1: 2, 2: 4}),
		}); IsError(sourceInfo) {
		t.Error(sourceInfo)
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
	typeProvider := NewProvider(&expr.ParsedExpr{})
	if sourceInfo := typeProvider.NewValue(
		"google.api.expr.v1.SourceInfo",
		map[string]ref.Value{
			"location":     String("TestTypeProvider_GetFieldValue"),
			"line_offsets": NewDynamicList([]int64{0, 2}),
			"positions":    NewDynamicMap(map[int64]int64{1: 2, 2: 4}),
		}); IsError(sourceInfo) {
		t.Error(sourceInfo)
	} else {
		si := sourceInfo.(traits.Indexer)
		if loc := si.Get(String("location")); IsError(loc) {
			t.Error(loc)
		} else if loc.(String) != "TestTypeProvider_GetFieldValue" {
			t.Errorf("Expected %s, got %s",
				"TestTypeProvider_GetFieldValue",
				loc)
		}
		if pos := si.Get(String("positions")); IsError(pos) {
			t.Error(pos)
		} else if pos.Equal(NewDynamicMap(map[int64]int32{1: 2, 2: 4})) != True {
			t.Errorf("Expected map[int64]int32, got %v", pos)
		} else if posKeyVal := pos.(traits.Indexer).Get(Int(1)); IsError(posKeyVal) {
			t.Error(posKeyVal)
		} else if posKeyVal.(Int) != 2 {
			t.Error("Expected value to be int64, not int32")
		}
		if offsets := si.Get(String("line_offsets")); IsError(offsets) {
			t.Error(offsets)
		} else if offset1 := offsets.(traits.Lister).Get(Int(1)); IsError(offset1) {
			t.Error(offset1)
		} else if offset1.(Int) != 2 {
			t.Errorf("Expected index 1 to be value 2, was %v", offset1)
		}
	}
}

func TestExprToProto(t *testing.T) {
	// Core type conversion tests.
	expectExprToProto(t, True, true)
	expectExprToProto(t, Int(-1), int32(-1))
	expectExprToProto(t, Int(2), int64(2))
	expectExprToProto(t, Uint(3), uint32(3))
	expectExprToProto(t, Uint(4), uint64(4))
	expectExprToProto(t, Double(5.5), float32(5.5))
	expectExprToProto(t, Double(-5.5), float64(-5.5))
	expectExprToProto(t, String("hello"), "hello")
	expectExprToProto(t, Bytes("world"), []byte("world"))
	expectExprToProto(t, NewDynamicList([]int64{1, 2, 3}), []int32{1, 2, 3})
	expectExprToProto(t, NewDynamicMap(map[int64]int64{1: 1, 2: 1, 3: 1}),
		map[int32]int32{1: 1, 2: 1, 3: 1})

	// Null conversion tests.
	expectExprToProto(t, Null(structpb.NullValue_NULL_VALUE), structpb.NullValue_NULL_VALUE)

	// Proto conversion tests.
	parsedExpr := &expr.ParsedExpr{}
	expectExprToProto(t, NewObject(parsedExpr), parsedExpr)
}

func TestProtoToExpr(t *testing.T) {
	// Core type conversions.
	expectProtoToExpr(t, true, True)
	expectProtoToExpr(t, int32(-1), Int(-1))
	expectProtoToExpr(t, int64(2), Int(2))
	expectProtoToExpr(t, uint32(3), Uint(3))
	expectProtoToExpr(t, uint64(4), Uint(4))
	expectProtoToExpr(t, float32(5.5), Double(5.5))
	expectProtoToExpr(t, float64(-5.5), Double(-5.5))
	expectProtoToExpr(t, "hello", String("hello"))
	expectProtoToExpr(t, []byte("world"), Bytes("world"))
	expectProtoToExpr(t, []int32{1, 2, 3}, NewDynamicList([]int32{1, 2, 3}))
	expectProtoToExpr(t, map[int32]int32{1: 1, 2: 1, 3: 1},
		NewDynamicMap(map[int32]int32{1: 1, 2: 1, 3: 1}))
	// Null conversion test.
	expectProtoToExpr(t, structpb.NullValue_NULL_VALUE, Null(structpb.NullValue_NULL_VALUE))
	// Proto conversion test.
	parsedExpr := &expr.ParsedExpr{}
	expectProtoToExpr(t, parsedExpr, NewObject(parsedExpr))
}

func TestUnsupportedConversion(t *testing.T) {
	if val := NativeToValue(nonConvertable{}); !IsError(val) {
		t.Error("Expected error when converting non-proto struct to proto", val)
	}
}

func expectExprToProto(t *testing.T, in ref.Value, out interface{}) {
	t.Helper()
	if val, err := in.ConvertToNative(reflect.TypeOf(out)); err != nil {
		t.Error(err)
	} else {
		var equals bool
		switch val.(type) {
		case []byte:
			equals = bytes.Equal(val.([]byte), out.([]byte))
		case proto.Message:
			equals = proto.Equal(val.(proto.Message), out.(proto.Message))
		case bool, int32, int64, uint32, uint64, float32, float64, string:
			equals = val == out
		default:
			equals = reflect.DeepEqual(val, out)
		}
		if !equals {
			t.Errorf("Unexpected conversion from expr to proto.\n"+
				"expected: %T, actual: %T", val, out)
		}
	}
}

func expectProtoToExpr(t *testing.T, in interface{}, out ref.Value) {
	t.Helper()
	if val := NativeToValue(in); IsError(val) {
		t.Error(val)
	} else {
		if val.Equal(out) != True {
			t.Errorf("Unexpected conversion from expr to proto.\n"+
				"expected: %T, actual: %T", val, out)
		}
	}
}

type nonConvertable struct {
	Field string
}

func BenchmarkTypeProvider_NewValue(b *testing.B) {
	typeProvider := NewProvider(&expr.ParsedExpr{})
	for i := 0; i < b.N; i++ {
		typeProvider.NewValue(
			"google.api.expr.v1.SourceInfo",
			map[string]ref.Value{
				"Location":    String("BenchmarkTypeProvider_NewValue"),
				"LineOffsets": NewDynamicList([]int64{0, 2}),
				"Positions":   NewDynamicMap(map[int64]int64{1: 2, 2: 4}),
			})
	}
}
