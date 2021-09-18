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

package interpreter

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"
	"google.golang.org/protobuf/proto"

	proto2pb "github.com/google/cel-go/test/proto2pb"
	proto3pb "github.com/google/cel-go/test/proto3pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	structpb "google.golang.org/protobuf/types/known/structpb"
	tpb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

type testCase struct {
	name           string
	expr           string
	container      string
	cost           []int64
	exhaustiveCost []int64
	optimizedCost  []int64
	abbrevs        []string
	env            []*exprpb.Decl
	types          []proto.Message
	funcs          []*functions.Overload
	attrs          AttributeFactory
	unchecked      bool

	in  map[string]interface{}
	out interface{}
	err string
}

var (
	testData = []testCase{
		{
			name:           "and_false_1st",
			expr:           `false && true`,
			cost:           []int64{0, 1},
			exhaustiveCost: []int64{1, 1},
			out:            types.False,
		},
		{
			name:           "and_false_2nd",
			expr:           `true && false`,
			cost:           []int64{0, 1},
			exhaustiveCost: []int64{1, 1},
			out:            types.False,
		},
		{
			name:           "and_error_1st_false",
			expr:           `1/0 != 0 && false`,
			cost:           []int64{2, 3},
			exhaustiveCost: []int64{3, 3},
			out:            types.False,
		},
		{
			name:           "and_error_2nd_false",
			expr:           `false && 1/0 != 0`,
			cost:           []int64{0, 3},
			exhaustiveCost: []int64{3, 3},
			out:            types.False,
		},
		{
			name:           "and_error_1st_error",
			expr:           `1/0 != 0 && true`,
			cost:           []int64{2, 3},
			exhaustiveCost: []int64{3, 3},
			err:            "division by zero",
		},
		{
			name:           "and_error_2nd_error",
			expr:           `true && 1/0 != 0`,
			cost:           []int64{0, 3},
			exhaustiveCost: []int64{3, 3},
			err:            "division by zero",
		},
		{
			name:      "call_no_args",
			expr:      `zero()`,
			cost:      []int64{1, 1},
			unchecked: true,
			funcs: []*functions.Overload{
				{
					Operator: "zero",
					Function: func(args ...ref.Val) ref.Val {
						return types.IntZero
					},
				},
			},
			out: types.IntZero,
		},
		{
			name:      "call_one_arg",
			expr:      `neg(1)`,
			cost:      []int64{1, 1},
			unchecked: true,
			funcs: []*functions.Overload{
				{
					Operator:     "neg",
					OperandTrait: traits.NegatorType,
					Unary: func(arg ref.Val) ref.Val {
						return arg.(traits.Negater).Negate()
					},
				},
			},
			out: types.IntNegOne,
		},
		{
			name:      "call_two_arg",
			expr:      `b'abc'.concat(b'def')`,
			cost:      []int64{1, 1},
			unchecked: true,
			funcs: []*functions.Overload{
				{
					Operator:     "concat",
					OperandTrait: traits.AdderType,
					Binary: func(lhs, rhs ref.Val) ref.Val {
						return lhs.(traits.Adder).Add(rhs)
					},
				},
			},
			out: []byte{'a', 'b', 'c', 'd', 'e', 'f'},
		},
		{
			name:      "call_varargs",
			expr:      `addall(a, b, c, d) == 10`,
			cost:      []int64{6, 6},
			unchecked: true,
			funcs: []*functions.Overload{
				{
					Operator:     "addall",
					OperandTrait: traits.AdderType,
					Function: func(args ...ref.Val) ref.Val {
						val := types.Int(0)
						for _, arg := range args {
							val += arg.(types.Int)
						}
						return val
					},
				},
			},
			in: map[string]interface{}{
				"a": 1, "b": 2, "c": 3, "d": 4,
			},
		},
		{
			name: `call_ns_func`,
			expr: `base64.encode('hello')`,
			cost: []int64{1, 1},
			env: []*exprpb.Decl{
				decls.NewFunction("base64.encode",
					decls.NewOverload("base64_encode_string",
						[]*exprpb.Type{decls.String},
						decls.String),
				),
			},
			funcs: []*functions.Overload{
				{
					Operator: "base64.encode",
					Unary:    base64Encode,
				},
				{
					Operator: "base64_encode_string",
					Unary:    base64Encode,
				},
			},
			out: "aGVsbG8=",
		},
		{
			name:      `call_ns_func_unchecked`,
			expr:      `base64.encode('hello')`,
			cost:      []int64{1, 1},
			unchecked: true,
			funcs: []*functions.Overload{
				{
					Operator: "base64.encode",
					Unary:    base64Encode,
				},
			},
			out: "aGVsbG8=",
		},
		{
			name:      `call_ns_func_in_pkg`,
			container: `base64`,
			expr:      `encode('hello')`,
			cost:      []int64{1, 1},
			env: []*exprpb.Decl{
				decls.NewFunction("base64.encode",
					decls.NewOverload("base64_encode_string",
						[]*exprpb.Type{decls.String},
						decls.String),
				),
			},
			funcs: []*functions.Overload{
				{
					Operator: "base64.encode",
					Unary:    base64Encode,
				},
				{
					Operator: "base64_encode_string",
					Unary:    base64Encode,
				},
			},
			out: "aGVsbG8=",
		},
		{
			name:      `call_ns_func_unchecked_in_pkg`,
			expr:      `encode('hello')`,
			cost:      []int64{1, 1},
			container: `base64`,
			unchecked: true,
			funcs: []*functions.Overload{
				{
					Operator: "base64.encode",
					Unary:    base64Encode,
				},
			},
			out: "aGVsbG8=",
		},
		{
			name: "complex",
			expr: `
			!(headers.ip in ["10.0.1.4", "10.0.1.5"]) &&
				((headers.path.startsWith("v1") && headers.token in ["v1", "v2", "admin"]) ||
				(headers.path.startsWith("v2") && headers.token in ["v2", "admin"]) ||
				(headers.path.startsWith("/admin") && headers.token == "admin" && headers.ip in ["10.0.1.2", "10.0.1.2", "10.0.1.2"]))
			`,
			cost:           []int64{3, 24},
			exhaustiveCost: []int64{24, 24},
			optimizedCost:  []int64{2, 20},
			env: []*exprpb.Decl{
				decls.NewVar("headers", decls.NewMapType(decls.String, decls.String)),
			},
			in: map[string]interface{}{
				"headers": map[string]interface{}{
					"ip":    "10.0.1.2",
					"path":  "/admin/edit",
					"token": "admin",
				},
			},
		},
		{
			name: "complex_qual_vars",
			expr: `
			!(headers.ip in ["10.0.1.4", "10.0.1.5"]) &&
				((headers.path.startsWith("v1") && headers.token in ["v1", "v2", "admin"]) ||
				(headers.path.startsWith("v2") && headers.token in ["v2", "admin"]) ||
				(headers.path.startsWith("/admin") && headers.token == "admin" && headers.ip in ["10.0.1.2", "10.0.1.2", "10.0.1.2"]))
			`,
			cost:           []int64{3, 24},
			exhaustiveCost: []int64{24, 24},
			optimizedCost:  []int64{2, 20},
			env: []*exprpb.Decl{
				decls.NewVar("headers.ip", decls.String),
				decls.NewVar("headers.path", decls.String),
				decls.NewVar("headers.token", decls.String),
			},
			in: map[string]interface{}{
				"headers.ip":    "10.0.1.2",
				"headers.path":  "/admin/edit",
				"headers.token": "admin",
			},
		},
		{
			name: "cond",
			expr: `a ? b < 1.2 : c == ['hello']`,
			cost: []int64{3, 3},
			env: []*exprpb.Decl{
				decls.NewVar("a", decls.Bool),
				decls.NewVar("b", decls.Double),
				decls.NewVar("c", decls.NewListType(decls.String)),
			},
			in: map[string]interface{}{
				"a": true,
				"b": 2.0,
				"c": []string{"hello"},
			},
			out: types.False,
		},
		{
			name:          "in_list",
			expr:          `6 in [2, 12, 6]`,
			cost:          []int64{1, 1},
			optimizedCost: []int64{0, 0},
		},
		{
			name: "in_map",
			expr: `'other-key' in {'key': null, 'other-key': 42}`,
			cost: []int64{1, 1},
			out:  types.True,
		},
		{
			name: "in_heterogeneous_map",
			expr: `'hello' in {1: 'one', false: true, 'hello': 'world'}`,
			out:  types.True,
		},
		{
			name: "not_in_heterogeneous_map",
			expr: `!('hello' in {1: 'one', false: true})`,
			err:  "unsupported key type: string",
		},
		{
			name: "not_in_heterogeneous_map_with_same_key_type",
			expr: `!('hello' in {1: 'one', 'world': true})`,
			out:  types.True,
		},
		{
			name:           "index",
			expr:           `m['key'][1] == 42u && m['null'] == null && m[string(0)] == 10`,
			cost:           []int64{2, 9},
			exhaustiveCost: []int64{9, 9},
			optimizedCost:  []int64{2, 8},
			env: []*exprpb.Decl{
				decls.NewVar("m", decls.NewMapType(decls.String, decls.Dyn)),
			},
			in: map[string]interface{}{
				"m": map[string]interface{}{
					"key":  []uint{21, 42},
					"null": nil,
					"0":    10,
				},
			},
		},
		{
			name: "index_relative",
			expr: `([[[1]], [[2]], [[3]]][0][0] + [2, 3, {'four': {'five': 'six'}}])[3].four.five == 'six'`,
			cost: []int64{2, 2},
		},
		{
			name: "list_eq_false_with_error",
			expr: `['string', 1] == [2, 3]`,
			out:  types.False,
		},
		{
			name: "list_eq_error",
			expr: `['string', true] == [2, 3]`,
			err:  "no such overload",
		},
		{
			name: "literal_bool_false",
			expr: `false`,
			cost: []int64{0, 0},
			out:  types.False,
		},
		{
			name: "literal_bool_true",
			expr: `true`,
			cost: []int64{0, 0},
		},
		{
			name: "literal_null",
			expr: `null`,
			cost: []int64{0, 0},
			out:  types.NullValue,
		},
		{
			name: "literal_list",
			expr: `[1, 2, 3]`,
			cost: []int64{0, 0},
			out:  []int64{1, 2, 3},
		},
		{
			name: "literal_map",
			expr: `{'hi': 21, 'world': 42u}`,
			cost: []int64{0, 0},
			out: map[string]interface{}{
				"hi":    21,
				"world": uint(42),
			},
		},
		{
			name:          "literal_equiv_string_bytes",
			expr:          `string(bytes("\303\277")) == '''\303\277'''`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "literal_not_equiv_string_bytes",
			expr:          `string(b"\303\277") != '''\303\277'''`,
			cost:          []int64{2, 2},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "literal_equiv_bytes_string",
			expr:          `string(b"\303\277") == 'ÿ'`,
			cost:          []int64{2, 2},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "literal_bytes_string",
			expr:          `string(b'aaa"bbb')`,
			cost:          []int64{1, 1},
			optimizedCost: []int64{0, 0},
			out:           `aaa"bbb`,
		},
		{
			name:          "literal_bytes_string2",
			expr:          `string(b"""Kim\t""")`,
			cost:          []int64{1, 1},
			optimizedCost: []int64{0, 0},
			out: `Kim	`,
		},
		{
			name:      "literal_pb3_msg",
			container: "google.api.expr",
			types:     []proto.Message{&exprpb.Expr{}},
			expr: `v1alpha1.Expr{
				id: 1,
				const_expr: v1alpha1.Constant{
					string_value: "oneof_test"
				}
			}`,
			cost: []int64{0, 0},
			out: &exprpb.Expr{Id: 1,
				ExprKind: &exprpb.Expr_ConstExpr{
					ConstExpr: &exprpb.Constant{
						ConstantKind: &exprpb.Constant_StringValue{
							StringValue: "oneof_test"}}}},
		},
		{
			name:      "literal_pb_enum",
			container: "google.expr.proto3.test",
			types:     []proto.Message{&proto3pb.TestAllTypes{}},
			expr: `TestAllTypes{
				repeated_nested_enum: [
					0,
					TestAllTypes.NestedEnum.BAZ,
					TestAllTypes.NestedEnum.BAR],
				repeated_int32: [
					TestAllTypes.NestedEnum.FOO,
					TestAllTypes.NestedEnum.BAZ]}`,
			cost: []int64{0, 0},
			out: &proto3pb.TestAllTypes{
				RepeatedNestedEnum: []proto3pb.TestAllTypes_NestedEnum{
					proto3pb.TestAllTypes_FOO,
					proto3pb.TestAllTypes_BAZ,
					proto3pb.TestAllTypes_BAR,
				},
				RepeatedInt32: []int32{0, 2},
			},
		},
		{
			name:          "timestamp_eq_timestamp",
			expr:          `timestamp(0) == timestamp(0)`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "timestamp_ne_timestamp",
			expr:          `timestamp(1) != timestamp(2)`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "timestamp_lt_timestamp",
			expr:          `timestamp(0) < timestamp(1)`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "timestamp_le_timestamp",
			expr:          `timestamp(2) <= timestamp(2)`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "timestamp_gt_timestamp",
			expr:          `timestamp(1) > timestamp(0)`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "timestamp_ge_timestamp",
			expr:          `timestamp(2) >= timestamp(2)`,
			cost:          []int64{3, 3},
			optimizedCost: []int64{1, 1},
		},
		{
			name:          "string_to_timestamp",
			expr:          `timestamp('1986-04-26T01:23:40Z')`,
			cost:          []int64{1, 1},
			optimizedCost: []int64{0, 0},
			out:           &tpb.Timestamp{Seconds: 514862620},
		},
		{
			name:           "macro_all_non_strict",
			expr:           `![0, 2, 4].all(x, 4/x != 2 && 4/(4-x) != 2)`,
			cost:           []int64{5, 38},
			exhaustiveCost: []int64{38, 38},
		},
		{
			name: "macro_all_non_strict_var",
			expr: `code == "111" && ["a", "b"].all(x, x in tags)
				|| code == "222" && ["a", "b"].all(x, x in tags)`,
			env: []*exprpb.Decl{
				decls.NewVar("code", decls.String),
				decls.NewVar("tags", decls.NewListType(decls.String)),
			},
			in: map[string]interface{}{
				"code": "222",
				"tags": []string{"a", "b"},
			},
		},
		{
			name: "macro_exists_lit",
			expr: `[1, 2, 3, 4, 5u, 1.0].exists(e, type(e) == uint)`,
		},
		{
			name: "macro_exists_nonstrict",
			expr: `[0, 2, 4].exists(x, 4/x == 2 && 4/(4-x) == 2)`,
		},
		{
			name:           "macro_exists_var",
			expr:           `elems.exists(e, type(e) == uint)`,
			cost:           []int64{0, 9223372036854775807},
			exhaustiveCost: []int64{0, 9223372036854775807},
			env: []*exprpb.Decl{
				decls.NewVar("elems", decls.NewListType(decls.Dyn)),
			},
			in: map[string]interface{}{
				"elems": []interface{}{0, 1, 2, 3, 4, uint(5), 6},
			},
		},
		{
			name: "macro_exists_one",
			expr: `[1, 2, 3].exists_one(x, (x % 2) == 0)`,
		},
		{
			name: "macro_filter",
			expr: `[1, 2, 3].filter(x, x > 2) == [3]`,
		},
		{
			name:           "macro_has_map_key",
			expr:           `has({'a':1}.a) && !has({}.a)`,
			cost:           []int64{1, 4},
			exhaustiveCost: []int64{4, 4},
		},
		{
			name:      "macro_has_pb2_field",
			container: "google.expr.proto2.test",
			types:     []proto.Message{&proto2pb.TestAllTypes{}},
			env: []*exprpb.Decl{
				decls.NewVar("pb2", decls.NewObjectType("google.expr.proto2.test.TestAllTypes")),
			},
			in: map[string]interface{}{
				"pb2": &proto2pb.TestAllTypes{
					RepeatedBool: []bool{false},
					MapInt64NestedType: map[int64]*proto2pb.NestedTestAllTypes{
						1: {},
					},
					MapStringString: map[string]string{},
				},
			},
			expr: `has(TestAllTypes{standalone_enum: TestAllTypes.NestedEnum.BAR}.standalone_enum)
			&& has(TestAllTypes{standalone_enum: TestAllTypes.NestedEnum.FOO}.standalone_enum)
			&& !has(TestAllTypes{single_nested_enum: TestAllTypes.NestedEnum.FOO}.single_nested_message)
			&& has(TestAllTypes{single_nested_enum: TestAllTypes.NestedEnum.FOO}.single_nested_enum)
			&& !has(TestAllTypes{}.standalone_enum)
			&& !has(pb2.single_int64)
			&& has(pb2.repeated_bool)
			&& !has(pb2.repeated_int32)
			&& has(pb2.map_int64_nested_type)
			&& !has(pb2.map_string_string)`,
			cost:           []int64{1, 29},
			exhaustiveCost: []int64{29, 29},
		},
		{
			name:  "macro_has_pb3_field",
			types: []proto.Message{&proto3pb.TestAllTypes{}},
			env: []*exprpb.Decl{
				decls.NewVar("pb3", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
			container: "google.expr.proto3.test",
			in: map[string]interface{}{
				"pb3": &proto3pb.TestAllTypes{
					RepeatedBool: []bool{false},
					MapInt64NestedType: map[int64]*proto3pb.NestedTestAllTypes{
						1: {},
					},
					MapStringString: map[string]string{},
				},
			},
			expr: `has(TestAllTypes{standalone_enum: TestAllTypes.NestedEnum.BAR}.standalone_enum)
			&& !has(TestAllTypes{standalone_enum: TestAllTypes.NestedEnum.FOO}.standalone_enum)
			&& !has(TestAllTypes{single_nested_enum: TestAllTypes.NestedEnum.FOO}.single_nested_message)
			&& has(TestAllTypes{single_nested_enum: TestAllTypes.NestedEnum.FOO}.single_nested_enum)
			&& !has(TestAllTypes{}.single_nested_message)
			&& has(TestAllTypes{single_nested_message: TestAllTypes.NestedMessage{}}.single_nested_message)
			&& !has(TestAllTypes{}.standalone_enum)
			&& !has(pb3.single_int64)
			&& has(pb3.repeated_bool)
			&& !has(pb3.repeated_int32)
			&& has(pb3.map_int64_nested_type)
			&& !has(pb3.map_string_string)`,
			cost:           []int64{1, 35},
			exhaustiveCost: []int64{35, 35},
		},
		{
			name:           "macro_map",
			expr:           `[1, 2, 3].map(x, x * 2) == [2, 4, 6]`,
			cost:           []int64{6, 14},
			exhaustiveCost: []int64{14, 14},
		},
		{
			name: "matches",
			expr: `input.matches('k.*')
				&& !'foo'.matches('k.*')
				&& !'bar'.matches('k.*')
				&& 'kilimanjaro'.matches('.*ro')`,
			cost:           []int64{2, 10},
			exhaustiveCost: []int64{10, 10},
			env: []*exprpb.Decl{
				decls.NewVar("input", decls.String),
			},
			in: map[string]interface{}{
				"input": "kathmandu",
			},
		},
		{
			name:  "nested_proto_field",
			expr:  `pb3.single_nested_message.bb`,
			cost:  []int64{1, 1},
			types: []proto.Message{&proto3pb.TestAllTypes{}},
			env: []*exprpb.Decl{
				decls.NewVar("pb3",
					decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
			in: map[string]interface{}{
				"pb3": &proto3pb.TestAllTypes{
					NestedType: &proto3pb.TestAllTypes_SingleNestedMessage{
						SingleNestedMessage: &proto3pb.TestAllTypes_NestedMessage{
							Bb: 1234,
						},
					},
				},
			},
			out: types.Int(1234),
		},
		{
			name:  "nested_proto_field_with_index",
			expr:  `pb3.map_int64_nested_type[0].child.payload.single_int32 == 1`,
			cost:  []int64{2, 2},
			types: []proto.Message{&proto3pb.TestAllTypes{}},
			env: []*exprpb.Decl{
				decls.NewVar("pb3",
					decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
			in: map[string]interface{}{
				"pb3": &proto3pb.TestAllTypes{
					MapInt64NestedType: map[int64]*proto3pb.NestedTestAllTypes{
						0: {
							Child: &proto3pb.NestedTestAllTypes{
								Payload: &proto3pb.TestAllTypes{
									SingleInt32: 1,
								},
							},
						},
					},
				},
			},
		},
		{
			name:           "or_true_1st",
			expr:           `ai == 20 || ar["foo"] == "bar"`,
			cost:           []int64{2, 5},
			exhaustiveCost: []int64{5, 5},
			env: []*exprpb.Decl{
				decls.NewVar("ai", decls.Int),
				decls.NewVar("ar", decls.NewMapType(decls.String, decls.String)),
			},
			in: map[string]interface{}{
				"ai": 20,
				"ar": map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name:           "or_true_2nd",
			expr:           `ai == 20 || ar["foo"] == "bar"`,
			cost:           []int64{2, 5},
			exhaustiveCost: []int64{5, 5},
			env: []*exprpb.Decl{
				decls.NewVar("ai", decls.Int),
				decls.NewVar("ar", decls.NewMapType(decls.String, decls.String)),
			},
			in: map[string]interface{}{
				"ai": 2,
				"ar": map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name:           "or_false",
			expr:           `ai == 20 || ar["foo"] == "bar"`,
			cost:           []int64{2, 5},
			exhaustiveCost: []int64{5, 5},
			env: []*exprpb.Decl{
				decls.NewVar("ai", decls.Int),
				decls.NewVar("ar", decls.NewMapType(decls.String, decls.String)),
			},
			in: map[string]interface{}{
				"ai": 2,
				"ar": map[string]string{
					"foo": "baz",
				},
			},
			out: types.False,
		},
		{
			name:           "or_error_1st_error",
			expr:           `1/0 != 0 || false`,
			cost:           []int64{2, 3},
			exhaustiveCost: []int64{3, 3},
			err:            "division by zero",
		},
		{
			name:           "or_error_2nd_error",
			expr:           `false || 1/0 != 0`,
			cost:           []int64{0, 3},
			exhaustiveCost: []int64{3, 3},
			err:            "division by zero",
		},
		{
			name:           "or_error_1st_true",
			expr:           `1/0 != 0 || true`,
			cost:           []int64{2, 3},
			exhaustiveCost: []int64{3, 3},
			out:            types.True,
		},
		{
			name:           "or_error_2nd_true",
			expr:           `true || 1/0 != 0`,
			cost:           []int64{0, 3},
			exhaustiveCost: []int64{3, 3},
			out:            types.True,
		},
		{
			name:      "pkg_qualified_id",
			expr:      `b.c.d != 10`,
			cost:      []int64{2, 2},
			container: "a.b",
			env: []*exprpb.Decl{
				decls.NewVar("a.b.c.d", decls.Int),
			},
			in: map[string]interface{}{
				"a.b.c.d": 9,
			},
		},
		{
			name:      "pkg_qualified_id_unchecked",
			expr:      `c.d != 10`,
			cost:      []int64{2, 2},
			unchecked: true,
			container: "a.b",
			in: map[string]interface{}{
				"a.c.d": 9,
			},
		},
		{
			name:      "pkg_qualified_index_unchecked",
			expr:      `b.c['d'] == 10`,
			cost:      []int64{2, 2},
			unchecked: true,
			container: "a.b",
			in: map[string]interface{}{
				"a.b.c": map[string]int{
					"d": 10,
				},
			},
		},
		{
			name: "select_key",
			expr: `m.strMap['val'] == 'string'
				&& m.floatMap['val'] == 1.5
				&& m.doubleMap['val'] == -2.0
				&& m.intMap['val'] == -3
				&& m.int32Map['val'] == 4
				&& m.int64Map['val'] == -5
				&& m.uintMap['val'] == 6u
				&& m.uint32Map['val'] == 7u
				&& m.uint64Map['val'] == 8u
				&& m.boolMap['val'] == true
				&& m.boolMap['val'] != false`,
			cost:           []int64{2, 32},
			exhaustiveCost: []int64{32, 32},
			env: []*exprpb.Decl{
				decls.NewVar("m", decls.NewMapType(decls.String, decls.Dyn)),
			},
			in: map[string]interface{}{
				"m": map[string]interface{}{
					"strMap":    map[string]string{"val": "string"},
					"floatMap":  map[string]float32{"val": 1.5},
					"doubleMap": map[string]float64{"val": -2.0},
					"intMap":    map[string]int{"val": -3},
					"int32Map":  map[string]int32{"val": 4},
					"int64Map":  map[string]int64{"val": -5},
					"uintMap":   map[string]uint{"val": 6},
					"uint32Map": map[string]uint32{"val": 7},
					"uint64Map": map[string]uint64{"val": 8},
					"boolMap":   map[string]bool{"val": true},
				},
			},
		},
		{
			name: "select_bool_key",
			expr: `m.boolStr[true] == 'string'
				&& m.boolFloat32[true] == 1.5
				&& m.boolFloat64[false] == -2.1
				&& m.boolInt[false] == -3
				&& m.boolInt32[false] == 0
				&& m.boolInt64[true] == 4
				&& m.boolUint[true] == 5u
				&& m.boolUint32[true] == 6u
				&& m.boolUint64[false] == 7u
				&& m.boolBool[true]
				&& m.boolIface[false] == true`,
			cost:           []int64{2, 31},
			exhaustiveCost: []int64{31, 31},
			env: []*exprpb.Decl{
				decls.NewVar("m", decls.NewMapType(decls.String, decls.Dyn)),
			},
			in: map[string]interface{}{
				"m": map[string]interface{}{
					"boolStr":     map[bool]string{true: "string"},
					"boolFloat32": map[bool]float32{true: 1.5},
					"boolFloat64": map[bool]float64{false: -2.1},
					"boolInt":     map[bool]int{false: -3},
					"boolInt32":   map[bool]int32{false: 0},
					"boolInt64":   map[bool]int64{true: 4},
					"boolUint":    map[bool]uint{true: 5},
					"boolUint32":  map[bool]uint32{true: 6},
					"boolUint64":  map[bool]uint64{false: 7},
					"boolBool":    map[bool]bool{true: true},
					"boolIface":   map[bool]interface{}{false: true},
				},
			},
		},
		{
			name: "select_uint_key",
			expr: `m.uintIface[1u] == 'string'
				&& m.uint32Iface[2u] == 1.5
				&& m.uint64Iface[3u] == -2.1
				&& m.uint64String[4u] == 'three'`,
			cost:           []int64{2, 11},
			exhaustiveCost: []int64{11, 11},
			env: []*exprpb.Decl{
				decls.NewVar("m", decls.NewMapType(decls.String, decls.Dyn)),
			},
			in: map[string]interface{}{
				"m": map[string]interface{}{
					"uintIface":    map[uint]interface{}{1: "string"},
					"uint32Iface":  map[uint32]interface{}{2: 1.5},
					"uint64Iface":  map[uint64]interface{}{3: -2.1},
					"uint64String": map[uint64]string{4: "three"},
				},
			},
		},
		{
			name: "select_index",
			expr: `m.strList[0] == 'string'
				&& m.floatList[0] == 1.5
				&& m.doubleList[0] == -2.0
				&& m.intList[0] == -3
				&& m.int32List[0] == 4
				&& m.int64List[0] == -5
				&& m.uintList[0] == 6u
				&& m.uint32List[0] == 7u
				&& m.uint64List[0] == 8u
				&& m.boolList[0] == true
				&& m.boolList[1] != true
				&& m.ifaceList[0] == {}`,
			cost:           []int64{2, 35},
			exhaustiveCost: []int64{35, 35},
			env: []*exprpb.Decl{
				decls.NewVar("m", decls.NewMapType(decls.String, decls.Dyn)),
			},
			in: map[string]interface{}{
				"m": map[string]interface{}{
					"strList":    []string{"string"},
					"floatList":  []float32{1.5},
					"doubleList": []float64{-2.0},
					"intList":    []int{-3},
					"int32List":  []int32{4},
					"int64List":  []int64{-5},
					"uintList":   []uint{6},
					"uint32List": []uint32{7},
					"uint64List": []uint64{8},
					"boolList":   []bool{true, false},
					"ifaceList":  []interface{}{map[string]string{}},
				},
			},
		},
		{
			name: "select_field",
			expr: `a.b.c
				&& pb3.repeated_nested_enum[0] == test.TestAllTypes.NestedEnum.BAR
				&& json.list[0] == 'world'`,
			cost:           []int64{1, 7},
			exhaustiveCost: []int64{7, 7},
			container:      "google.expr.proto3",
			types:          []proto.Message{&proto3pb.TestAllTypes{}},
			env: []*exprpb.Decl{
				decls.NewVar("a.b", decls.NewMapType(decls.String, decls.Bool)),
				decls.NewVar("pb3", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
				decls.NewVar("json", decls.NewMapType(decls.String, decls.Dyn)),
			},
			in: map[string]interface{}{
				"a.b": map[string]bool{
					"c": true,
				},
				"pb3": &proto3pb.TestAllTypes{
					RepeatedNestedEnum: []proto3pb.TestAllTypes_NestedEnum{proto3pb.TestAllTypes_BAR},
				},
				"json": &structpb.Value{
					Kind: &structpb.Value_StructValue{
						StructValue: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"list": {Kind: &structpb.Value_ListValue{
									ListValue: &structpb.ListValue{
										Values: []*structpb.Value{
											structpb.NewStringValue("world"),
										},
									},
								}},
							},
						},
					},
				},
			},
		},
		// pb2 primitive fields may have default values set.
		{
			name: "select_pb2_primitive_fields",
			expr: `!has(a.single_int32)
			&& a.single_int32 == -32
			&& a.single_int64 == -64
			&& a.single_uint32 == 32u
			&& a.single_uint64 == 64u
			&& a.single_float == 3.0
			&& a.single_double == 6.4
			&& a.single_bool
			&& "empty" == a.single_string`,
			cost:           []int64{3, 26},
			exhaustiveCost: []int64{26, 26},
			types:          []proto.Message{&proto2pb.TestAllTypes{}},
			in: map[string]interface{}{
				"a": &proto2pb.TestAllTypes{},
			},
			env: []*exprpb.Decl{
				decls.NewVar("a", decls.NewObjectType("google.expr.proto2.test.TestAllTypes")),
			},
		},
		// Wrapper type nil or value test.
		{
			name: "select_pb3_wrapper_fields",
			expr: `!has(a.single_int32_wrapper) && a.single_int32_wrapper == null
				&& has(a.single_int64_wrapper) && a.single_int64_wrapper == 0
				&& has(a.single_string_wrapper) && a.single_string_wrapper == "hello"
				&& a.single_int64_wrapper == Int32Value{value: 0}`,
			cost:           []int64{3, 21},
			exhaustiveCost: []int64{21, 21},
			types:          []proto.Message{&proto3pb.TestAllTypes{}},
			abbrevs:        []string{"google.protobuf.Int32Value"},
			env: []*exprpb.Decl{
				decls.NewVar("a", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
			in: map[string]interface{}{
				"a": &proto3pb.TestAllTypes{
					SingleInt64Wrapper:  &wrapperspb.Int64Value{},
					SingleStringWrapper: wrapperspb.String("hello"),
				},
			},
		},
		{
			name:      "select_pb3_compare",
			expr:      `a.single_uint64 > 3u`,
			cost:      []int64{2, 2},
			container: "google.expr.proto3.test",
			types:     []proto.Message{&proto3pb.TestAllTypes{}},
			env: []*exprpb.Decl{
				decls.NewVar("a", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
			in: map[string]interface{}{
				"a": &proto3pb.TestAllTypes{
					SingleUint64: 10,
				},
			},
			out: types.True,
		},
		{
			name:      "select_custom_pb3_compare",
			expr:      `a.bb > 100`,
			cost:      []int64{2, 2},
			container: "google.expr.proto3.test",
			types:     []proto.Message{&proto3pb.TestAllTypes_NestedMessage{}},
			env: []*exprpb.Decl{
				decls.NewVar("a",
					decls.NewObjectType("google.expr.proto3.test.TestAllTypes.NestedMessage")),
			},
			attrs: &custAttrFactory{
				AttributeFactory: NewAttributeFactory(
					testContainer("google.expr.proto3.test"),
					types.DefaultTypeAdapter,
					types.NewEmptyRegistry(),
				),
			},
			in: map[string]interface{}{
				"a": &proto3pb.TestAllTypes_NestedMessage{
					Bb: 101,
				},
			},
			out: types.True,
		},
		{
			name: "select_relative",
			expr: `json('{"hi":"world"}').hi == 'world'`,
			cost: []int64{2, 2},
			env: []*exprpb.Decl{
				decls.NewFunction("json",
					decls.NewOverload("string_to_json",
						[]*exprpb.Type{decls.String}, decls.Dyn)),
			},
			funcs: []*functions.Overload{
				{
					Operator: "json",
					Unary: func(val ref.Val) ref.Val {
						str, ok := val.(types.String)
						if !ok {
							return types.MaybeNoSuchOverloadErr(val)
						}
						m := make(map[string]interface{})
						err := json.Unmarshal([]byte(str), &m)
						if err != nil {
							return types.NewErr("invalid json: %v", err)
						}
						return types.DefaultTypeAdapter.NativeToValue(m)
					},
				},
			},
		},
		{
			name: "select_subsumed_field",
			expr: `a.b.c`,
			cost: []int64{1, 1},
			env: []*exprpb.Decl{
				decls.NewVar("a.b.c", decls.Int),
				decls.NewVar("a.b", decls.NewMapType(decls.String, decls.String)),
			},
			in: map[string]interface{}{
				"a.b.c": 10,
				"a.b": map[string]string{
					"c": "ten",
				},
			},
			out: types.Int(10),
		},
		{
			name:      "select_empty_repeated_nested",
			expr:      `TestAllTypes{}.repeated_nested_message.size() == 0`,
			cost:      []int64{2, 2},
			types:     []proto.Message{&proto3pb.TestAllTypes{}},
			container: "google.expr.proto3.test",
			out:       types.True,
		},
	}
)

func BenchmarkInterpreter(b *testing.B) {
	for _, tst := range testData {
		prg, vars, err := program(b, &tst, Optimize())
		if err != nil {
			b.Fatal(err)
		}
		// Benchmark the eval.
		b.Run(tst.name, func(bb *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < bb.N; i++ {
				prg.Eval(vars)
			}
		})
	}
}

func BenchmarkInterpreter_Parallel(b *testing.B) {
	for _, tst := range testData {
		prg, vars, err := program(b, &tst, Optimize())
		if err != nil {
			b.Fatal(err)
		}
		b.ResetTimer()
		b.Run(tst.name,
			func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						prg.Eval(vars)
					}
				})
			})
	}
}

func TestInterpreter(t *testing.T) {
	for _, tst := range testData {
		tc := tst
		prg, vars, err := program(t, &tc)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var want ref.Val = types.True
			if tc.out != nil {
				want = tc.out.(ref.Val)
			}
			got := prg.Eval(vars)
			_, expectUnk := want.(types.Unknown)
			if expectUnk {
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("Got %v, wanted %v", got, want)
				}
			} else if tc.err != "" {
				if !types.IsError(got) || got.(*types.Err).String() != tc.err {
					t.Fatalf("Got %v (%T), wanted error: %s", got, got, tc.err)
				}
			} else if got.Equal(want) != types.True {
				t.Fatalf("Got %v, wanted %v", got, want)
			}

			if tc.cost != nil {
				minCost, maxCost := estimateCost(prg)
				if minCost != tc.cost[0] || maxCost != tc.cost[1] {
					t.Errorf("Got cost interval [%v, %v], wanted %v", minCost, maxCost, tc.cost)
				}
			}
			state := NewEvalState()
			opts := map[string]InterpretableDecorator{
				"optimize":   Optimize(),
				"exhaustive": ExhaustiveEval(state),
				"track":      TrackState(state),
			}
			for mode, opt := range opts {
				prg, vars, err = program(t, &tc, opt)
				if err != nil {
					t.Fatal(err)
				}
				t.Run(mode, func(t *testing.T) {
					got := prg.Eval(vars)
					_, expectUnk := want.(types.Unknown)
					if expectUnk {
						if !reflect.DeepEqual(got, want) {
							t.Errorf("Got %v, wanted %v", got, want)
						}
					} else if tc.err != "" {
						if !types.IsError(got) || got.(*types.Err).String() != tc.err {
							t.Errorf("Got %v (%T), wanted error: %s", got, got, tc.err)
						}
					} else if got.Equal(want) != types.True {
						t.Errorf("Got %v, wanted %v", got, want)
					}
					if mode == "exhaustive" && tc.cost != nil {
						wantedCost := tc.cost
						if tc.exhaustiveCost != nil {
							wantedCost = tc.exhaustiveCost
						}
						minCost, maxCost := estimateCost(prg)
						if minCost != wantedCost[0] || maxCost != wantedCost[1] {
							t.Errorf("Got exhaustive cost interval [%v, %v], wanted %v",
								minCost, maxCost, wantedCost)
						}
					}
					if mode == "optimize" && tc.cost != nil {
						wantedCost := tc.cost
						if tc.optimizedCost != nil {
							wantedCost = tc.optimizedCost
						}
						minCost, maxCost := estimateCost(prg)
						if minCost != wantedCost[0] || maxCost != wantedCost[1] {
							t.Errorf("Got optimize cost interval [%v, %v], wanted %v", minCost, maxCost, tc.cost)
						}
					}
					state.Reset()
				})
			}
		})
	}
}

func TestInterpreter_ProtoAttributeOpt(t *testing.T) {
	inst, _, err := program(t, &testCase{
		name:  "nested_proto_field_with_index",
		expr:  `pb3.map_int64_nested_type[0].child.payload.single_int32`,
		types: []proto.Message{&proto3pb.TestAllTypes{}},
		env: []*exprpb.Decl{
			decls.NewVar("pb3",
				decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
		},
		in: map[string]interface{}{
			"pb3": &proto3pb.TestAllTypes{
				MapInt64NestedType: map[int64]*proto3pb.NestedTestAllTypes{
					0: {
						Child: &proto3pb.NestedTestAllTypes{
							Payload: &proto3pb.TestAllTypes{
								SingleInt32: 1,
							},
						},
					},
				},
			},
		},
	}, Optimize())
	if err != nil {
		t.Fatal(err)
	}
	attr, ok := inst.(InterpretableAttribute)
	if !ok {
		t.Fatalf("got %v, expected attribute value", inst)
	}
	absAttr, ok := attr.Attr().(NamespacedAttribute)
	if !ok {
		t.Fatalf("got %v, expected global attribute", attr.Attr())
	}
	quals := absAttr.Qualifiers()
	if len(quals) != 5 {
		t.Errorf("got %d qualifiers, wanted 5", len(quals))
	}
	if !isFieldQual(quals[0], "map_int64_nested_type") ||
		!isConstQual(quals[1], types.IntZero) ||
		!isFieldQual(quals[2], "child") ||
		!isFieldQual(quals[3], "payload") ||
		!isFieldQual(quals[4], "single_int32") {
		t.Error("unoptimized qualifier types present in optimized attribute")
	}
}

func TestInterpreter_LogicalAndMissingType(t *testing.T) {
	src := common.NewTextSource(`a && TestProto{c: true}.c`)
	parsed, errors := parser.Parse(src)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	reg := newTestRegistry(t)
	cont := containers.DefaultContainer
	attrs := NewAttributeFactory(cont, reg, reg)
	intr := NewStandardInterpreter(cont, reg, reg, attrs)
	i, err := intr.NewUncheckedInterpretable(parsed.GetExpr())
	if err == nil {
		t.Errorf("Got '%v', wanted error", i)
	}
}

func TestInterpreter_ExhaustiveConditionalExpr(t *testing.T) {
	src := common.NewTextSource(`a ? b < 1.0 : c == ['hello']`)
	parsed, errors := parser.Parse(src)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	state := NewEvalState()
	cont := containers.DefaultContainer
	reg := newTestRegistry(t, &exprpb.ParsedExpr{})
	attrs := NewAttributeFactory(cont, reg, reg)
	intr := NewStandardInterpreter(cont, reg, reg, attrs)
	interpretable, _ := intr.NewUncheckedInterpretable(
		parsed.GetExpr(),
		ExhaustiveEval(state))
	vars, _ := NewActivation(map[string]interface{}{
		"a": types.True,
		"b": types.Double(0.999),
		"c": types.NewStringList(reg, []string{"hello"})})
	result := interpretable.Eval(vars)
	// Operator "_==_" is at Expr 7, should be evaluated in exhaustive mode
	// even though "a" is true
	ev, _ := state.Value(7)
	// "==" should be evaluated in exhaustive mode though unnecessary
	if ev != types.True {
		t.Errorf("Else expression expected to be true, got: %v", ev)
	}
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_ExhaustiveLogicalOrEquals(t *testing.T) {
	// a || b == "b"
	// Operator "==" is at Expr 4, should be evaluated though "a" is true
	src := common.NewTextSource(`a || b == "b"`)
	parsed, errors := parser.Parse(src)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	state := NewEvalState()
	reg := newTestRegistry(t, &exprpb.Expr{})
	cont := testContainer("test")
	attrs := NewAttributeFactory(cont, reg, reg)
	interp := NewStandardInterpreter(cont, reg, reg, attrs)
	i, _ := interp.NewUncheckedInterpretable(
		parsed.GetExpr(),
		ExhaustiveEval(state))
	vars, _ := NewActivation(map[string]interface{}{
		"a": true,
		"b": "b",
	})
	result := i.Eval(vars)
	rhv, _ := state.Value(3)
	// "==" should be evaluated in exhaustive mode though unnecessary
	if rhv != types.True {
		t.Errorf("Right hand side expression expected to be true, got: %v", rhv)
	}
	if result != types.True {
		t.Errorf("Expected true, got: %v", result)
	}
}

func TestInterpreter_SetProto2PrimitiveFields(t *testing.T) {
	// Test the use of proto2 primitives within object construction.
	src := common.NewTextSource(
		`input == TestAllTypes{
			single_int32: 1,
			single_int64: 2,
			single_uint32: 3u,
			single_uint64: 4u,
			single_float: -3.3,
			single_double: -2.2,
			single_string: "hello world",
			single_bool: true
		}`)
	parsed, errors := parser.Parse(src)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	cont := testContainer("google.expr.proto2.test")
	reg := newTestRegistry(t, &proto2pb.TestAllTypes{})
	env := checker.NewStandardEnv(cont, reg)
	env.Add(decls.NewVar("input", decls.NewObjectType("google.expr.proto2.test.TestAllTypes")))
	checked, errors := checker.Check(parsed, src, env)
	if len(errors.GetErrors()) != 0 {
		t.Errorf(errors.ToDisplayString())
	}

	attrs := NewAttributeFactory(cont, reg, reg)
	i := NewStandardInterpreter(cont, reg, reg, attrs)
	eval, _ := i.NewInterpretable(checked)
	one := int32(1)
	two := int64(2)
	three := uint32(3)
	four := uint64(4)
	five := float32(-3.3)
	six := float64(-2.2)
	str := "hello world"
	truth := true
	input := &proto2pb.TestAllTypes{
		SingleInt32:  &one,
		SingleInt64:  &two,
		SingleUint32: &three,
		SingleUint64: &four,
		SingleFloat:  &five,
		SingleDouble: &six,
		SingleString: &str,
		SingleBool:   &truth,
	}
	vars, _ := NewActivation(map[string]interface{}{
		"input": reg.NativeToValue(input),
	})
	result := eval.Eval(vars)
	got, ok := result.(ref.Val).Value().(bool)
	if !ok {
		t.Fatalf("Got '%v', wanted 'true'.", result)
	}
	expected := true
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Could not build object properly. Got '%v', wanted '%v'",
			result.(ref.Val).Value(),
			expected)
	}
}

func TestInterpreter_MissingIdentInSelect(t *testing.T) {
	src := common.NewTextSource(`a.b.c`)
	parsed, errors := parser.Parse(src)
	if len(errors.GetErrors()) != 0 {
		t.Fatalf(errors.ToDisplayString())
	}

	cont := testContainer("test")
	reg := newTestRegistry(t)
	env := checker.NewStandardEnv(cont, reg)
	env.Add(decls.NewVar("a.b", decls.Dyn))
	checked, errors := checker.Check(parsed, src, env)
	if len(errors.GetErrors()) != 0 {
		t.Fatalf(errors.ToDisplayString())
	}

	attrs := NewPartialAttributeFactory(cont, reg, reg)
	interp := NewStandardInterpreter(cont, reg, reg, attrs)
	i, _ := interp.NewInterpretable(checked)
	vars, _ := NewPartialActivation(
		map[string]interface{}{
			"a.b": map[string]interface{}{
				"d": "hello",
			},
		},
		NewAttributePattern("a.b").QualString("c"))
	result := i.Eval(vars)
	if !types.IsUnknown(result) {
		t.Errorf("Got %v, wanted unknown", result)
	}

	result = i.Eval(EmptyActivation())
	if !types.IsError(result) {
		t.Errorf("Got %v, wanted error", result)
	}
}

func TestInterpreter_TypeConversionOpt(t *testing.T) {
	tests := []struct {
		in  string
		out ref.Val
		err bool
	}{
		{in: `bool('tru')`, err: true},
		{in: `bool("true")`, out: types.True},
		{in: `bytes("hello")`, out: types.Bytes("hello")},
		{in: `double("_123")`, err: true},
		{in: `double("123.0")`, out: types.Double(123.0)},
		{in: `duration('12hh3')`, err: true},
		{in: `duration('12s')`, out: types.Duration{Duration: time.Duration(12) * time.Second}},
		{in: `dyn(1u)`, out: types.Uint(1)},
		{in: `int('11l')`, err: true},
		{in: `int('11')`, out: types.Int(11)},
		{in: `string('11')`, out: types.String("11")},
		{in: `timestamp('123')`, err: true},
		{in: `timestamp(123)`, out: types.Timestamp{Time: time.Unix(123, 0).UTC()}},
		{in: `type(null)`, out: types.NullType},
		{in: `type(timestamp(int('123')))`, out: types.TimestampType},
		{in: `uint(-1)`, err: true},
		{in: `uint(1)`, out: types.Uint(1)},
	}
	for _, tc := range tests {
		src := common.NewTextSource(tc.in)
		parsed, errors := parser.Parse(src)
		if len(errors.GetErrors()) != 0 {
			t.Fatalf(errors.ToDisplayString())
		}
		cont := containers.DefaultContainer
		reg := newTestRegistry(t)
		env := checker.NewStandardEnv(cont, reg)
		checked, errors := checker.Check(parsed, src, env)
		if len(errors.GetErrors()) != 0 {
			t.Fatalf(errors.ToDisplayString())
		}
		attrs := NewAttributeFactory(cont, reg, reg)
		interp := NewStandardInterpreter(cont, reg, reg, attrs)
		// Show that program planning will now produce an error.
		i, err := interp.NewInterpretable(checked, Optimize())
		if tc.err && err == nil {
			t.Fatalf("got %v, expected error", i)
		}
		if tc.out != nil {
			if err != nil {
				t.Fatal(err)
			}
			ic, isConst := i.(InterpretableConst)
			if !isConst {
				t.Fatalf("got %v, expected constant", ic)
			}
			if tc.out.Equal(ic.Value()) != types.True {
				t.Errorf("got %v, wanted %v", ic.Value(), tc.out)
			}
		}
		// Show how the error returned during program planning is the same as the runtime
		// error which would be produced normally.
		if tc.err {
			i2, err2 := interp.NewInterpretable(checked)
			if err2 != nil {
				t.Fatalf("got error, wanted interpretable: %v", i2)
			}
			errVal := i2.Eval(EmptyActivation())
			errValStr := errVal.(*types.Err).Error()
			if errValStr != err.Error() {
				t.Errorf("got error %s, wanted error %s", errValStr, err.Error())
			}
		}
	}
}

func testContainer(name string) *containers.Container {
	cont, _ := containers.NewContainer(containers.Name(name))
	return cont
}

func program(ctx interface{}, tst *testCase, opts ...InterpretableDecorator) (Interpretable, Activation, error) {
	// Configure the package.
	cont := containers.DefaultContainer
	if tst.container != "" {
		cont = testContainer(tst.container)
	}
	var err error
	if tst.abbrevs != nil {
		cont, err = containers.NewContainer(
			containers.Name(cont.Name()),
			containers.Abbrevs(tst.abbrevs...))
		if err != nil {
			return nil, nil, err
		}
	}
	var reg ref.TypeRegistry
	switch t := ctx.(type) {
	case *testing.T:
		reg = newTestRegistry(t)
		if tst.types != nil {
			reg = newTestRegistry(t, tst.types...)
		}
	case *testing.B:
		reg = newBenchRegistry(t)
		if tst.types != nil {
			reg = newBenchRegistry(t, tst.types...)
		}
	}
	attrs := NewAttributeFactory(cont, reg, reg)
	if tst.attrs != nil {
		attrs = tst.attrs
	}

	// Configure the environment.
	env := checker.NewStandardEnv(cont, reg)
	if tst.env != nil {
		env.Add(tst.env...)
	}
	// Configure the program input.
	vars := EmptyActivation()
	if tst.in != nil {
		vars, _ = NewActivation(tst.in)
	}
	// Adapt the test output, if needed.
	if tst.out != nil {
		tst.out = reg.NativeToValue(tst.out)
	}

	disp := NewDispatcher()
	disp.Add(functions.StandardOverloads()...)
	if tst.funcs != nil {
		disp.Add(tst.funcs...)
	}
	interp := NewInterpreter(disp, cont, reg, reg, attrs)

	// Parse the expression.
	s := common.NewTextSource(tst.expr)
	parsed, errs := parser.Parse(s)
	if len(errs.GetErrors()) != 0 {
		return nil, nil, errors.New(errs.ToDisplayString())
	}
	if tst.unchecked {
		// Build the program plan.
		prg, err := interp.NewUncheckedInterpretable(
			parsed.GetExpr(), opts...)
		if err != nil {
			return nil, nil, err
		}
		return prg, vars, nil
	}
	// Check the expression.
	checked, errs := checker.Check(parsed, s, env)
	if len(errs.GetErrors()) != 0 {
		return nil, nil, errors.New(errs.ToDisplayString())
	}
	// Build the program plan.
	prg, err := interp.NewInterpretable(checked, opts...)
	if err != nil {
		return nil, nil, err
	}
	return prg, vars, nil
}

func base64Encode(val ref.Val) ref.Val {
	str, ok := val.(types.String)
	if !ok {
		return types.MaybeNoSuchOverloadErr(val)
	}
	return types.String(base64.StdEncoding.EncodeToString([]byte(str)))
}

func isConstQual(q Qualifier, val ref.Val) bool {
	c, ok := q.(ConstantQualifier)
	if !ok {
		return false
	}
	return c.Value().Equal(val) == types.True
}

func isFieldQual(q Qualifier, fieldName string) bool {
	f, ok := q.(*fieldQualifier)
	if !ok {
		return false
	}
	return f.Name == fieldName
}

func newBenchRegistry(b *testing.B, msgs ...proto.Message) ref.TypeRegistry {
	b.Helper()
	reg, err := types.NewRegistry(msgs...)
	if err != nil {
		b.Fatalf("types.NewRegistry(%v) failed: %v", msgs, err)
	}
	return reg
}

func newTestRegistry(t *testing.T, msgs ...proto.Message) ref.TypeRegistry {
	t.Helper()
	reg, err := types.NewRegistry(msgs...)
	if err != nil {
		t.Fatalf("types.NewRegistry(%v) failed: %v", msgs, err)
	}
	return reg
}
