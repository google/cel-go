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

package testing

import (
	"github.com/google/cel-go/operators"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	expr "github.com/google/cel-spec/proto/v1"
)

type TestExpr struct {
	Expr       *expr.Expr
	SourceInfo *expr.SourceInfo
}

func (t *TestExpr) Info(location string) *expr.SourceInfo {
	info := proto.Clone(t.SourceInfo).(*expr.SourceInfo)
	info.Location = location
	return info
}

var (
	// program with no instructions.
	Empty = &TestExpr{
		&expr.Expr{},

		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// [1, 1u, 1.0].exists(x, type(x) == uint)
	Exists = &TestExpr{
		ExprComprehension(1,
			"x",
			ExprList(5,
				ExprConst(2, int64(1)),
				ExprConst(3, uint64(1)),
				ExprConst(4, float64(1.0))),
			"_accu_",
			ExprConst(6, false),
			ExprCall(8,
				operators.LogicalNot,
				ExprIdent(7, "_accu_")),
			ExprCall(11,
				operators.Equals,
				ExprCall(10,
					"type",
					ExprIdent(9, "x")),
				ExprIdent(12, "uint")),
			ExprIdent(13, "_accu_")),

		&expr.SourceInfo{
			LineOffsets: []int32{0},
			Positions: map[int64]int32{
				0:  12,
				1:  0,
				2:  1,
				3:  4,
				4:  8,
				5:  0,
				6:  18,
				7:  18,
				8:  18,
				9:  18,
				10: 18,
				11: 20,
				12: 20,
				13: 28}}}

	// {"hello": "world".size(),
	//  "dur": duration.Duration{10},
	//  "ts": timestamp.Timestamp{1000},
	//  "null": null,
	//  "bytes": b"bytes-string"}
	DynMap = &TestExpr{
		ExprMap(17,
			ExprEntry(2,
				ExprConst(1, "hello"),
				ExprMemberCall(3,
					"size",
					ExprConst(4, "world"))),
			ExprEntry(6,
				ExprConst(5, "dur"),
				ExprConst(7, &duration.Duration{Seconds: 10})),
			ExprEntry(9,
				ExprConst(8, "ts"),
				ExprConst(10, &timestamp.Timestamp{Seconds: 1000})),
			ExprEntry(12,
				ExprConst(11, "null"),
				ExprConst(13, structpb.NullValue_NULL_VALUE)),
			ExprEntry(15,
				ExprConst(14, "bytes"),
				ExprConst(16, []byte("bytes-string")))),

		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a && TestProto{c: true}.c
	LogicalAnd = &TestExpr{
		ExprCall(2, operators.LogicalAnd,
			ExprIdent(1, "a"),
			ExprSelect(6,
				ExprType(5, "TestProto",
					ExprField(4, "c", ExprConst(5, true))),
				"c")),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a ? b < 1.0 : c == ["hello"]
	Conditional = &TestExpr{
		ExprCall(9, operators.Conditional,
			ExprIdent(1, "a"),
			ExprCall(3,
				operators.Less,
				ExprIdent(2, "b"),
				ExprConst(4, 1.0)),
			ExprCall(6,
				operators.Equals,
				ExprIdent(5, "c"),
				ExprList(8, ExprConst(7, "hello")))),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a.b.c
	Select = &TestExpr{
		ExprSelect(3,
			ExprSelect(2,
				ExprIdent(1, "a"),
				"b"),
			"c"),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}

	// a == 42
	Equality = &TestExpr{
		ExprCall(2,
			operators.Equals,
			ExprIdent(1, "a"),
			ExprConst(3, int64(42))),
		&expr.SourceInfo{
			LineOffsets: []int32{},
			Positions:   map[int64]int32{}}}
)

func ExprIdent(id int64, name string) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_IdentExpr{
		&expr.Expr_Ident{name}}}
}

func ExprSelect(id int64, operand *expr.Expr, field string) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_SelectExpr{
			&expr.Expr_Select{operand, field, false}}}
}

func ExprConst(id int64, value interface{}) *expr.Expr {
	var constExpr *expr.Constant
	switch value.(type) {
	case bool:
		constExpr = &expr.Constant{&expr.Constant_BoolValue{value.(bool)}}
	case int64:
		constExpr = &expr.Constant{&expr.Constant_Int64Value{
			value.(int64)}}
	case uint64:
		constExpr = &expr.Constant{&expr.Constant_Uint64Value{
			value.(uint64)}}
	case float64:
		constExpr = &expr.Constant{&expr.Constant_DoubleValue{
			value.(float64)}}
	case string:
		constExpr = &expr.Constant{&expr.Constant_StringValue{
			value.(string)}}
	case structpb.NullValue:
		constExpr = &expr.Constant{&expr.Constant_NullValue{
			value.(structpb.NullValue)}}
	case []byte:
		constExpr = &expr.Constant{&expr.Constant_BytesValue{
			value.([]byte)}}
	case *timestamp.Timestamp:
		constExpr = &expr.Constant{&expr.Constant_TimestampValue{
			value.(*timestamp.Timestamp)}}
	case *duration.Duration:
		constExpr = &expr.Constant{&expr.Constant_DurationValue{
			value.(*duration.Duration)}}
	default:
		panic("constant type not implemented")
	}
	return &expr.Expr{id, &expr.Expr_ConstExpr{ConstExpr: constExpr}}
}

func ExprCall(id int64, function string, args ...*expr.Expr) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_CallExpr{
			&expr.Expr_Call{nil, function, args}}}
}

func ExprMemberCall(id int64, function string, target *expr.Expr, args ...*expr.Expr) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_CallExpr{
			&expr.Expr_Call{target, function, args}}}
}

func ExprList(id int64, elements ...*expr.Expr) *expr.Expr {
	return &expr.Expr{id,
		&expr.Expr_ListExpr{
			&expr.Expr_CreateList{elements}}}
}

func ExprMap(id int64, entries ...*expr.Expr_CreateStruct_Entry) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_StructExpr{
		&expr.Expr_CreateStruct{Entries: entries}}}
}

func ExprType(id int64, messageName string,
	entries ...*expr.Expr_CreateStruct_Entry) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_StructExpr{
		&expr.Expr_CreateStruct{
			messageName, entries}}}
}

func ExprEntry(id int64, key *expr.Expr,
	value *expr.Expr) *expr.Expr_CreateStruct_Entry {
	return &expr.Expr_CreateStruct_Entry{id,
		&expr.Expr_CreateStruct_Entry_MapKey{key},
		value}
}

func ExprField(id int64, field string,
	value *expr.Expr) *expr.Expr_CreateStruct_Entry {
	return &expr.Expr_CreateStruct_Entry{id,
		&expr.Expr_CreateStruct_Entry_FieldKey{field},
		value}
}

func ExprComprehension(id int64,
	iterVar string, iterRange *expr.Expr,
	accuVar string, accuInit *expr.Expr,
	loopCondition *expr.Expr, loopStep *expr.Expr,
	resultExpr *expr.Expr) *expr.Expr {
	return &expr.Expr{id, &expr.Expr_ComprehensionExpr{
		&expr.Expr_Comprehension{
			iterVar,
			iterRange,
			accuVar,
			accuInit,
			loopCondition,
			loopStep,
			resultExpr}}}
}
