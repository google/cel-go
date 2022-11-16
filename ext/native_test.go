// Copyright 2022 Google LLC
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

package ext

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/pb"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/test"

	proto3pb "github.com/google/cel-go/test/proto3pb"
)

func TestNativeTypes(t *testing.T) {
	var nativeTests = []struct {
		expr string
		out  any
		in   any
	}{
		{
			expr: `ext.TestAllTypes{
				NestedVal: ext.TestNestedType{NestedMapVal: {1: false}},
				BoolVal: true,
				BytesVal: b'hello',
				DurationVal: duration('5s'),
				DoubleVal: 1.5,
				FloatVal: 2.5,
				Int32Val: 10,
				Int64Val: 20,
				StringVal: 'hello world',
				TimestampVal: timestamp('2011-08-06T01:23:45Z'),
				Uint32Val: 100u,
				Uint64Val: 200u,
				ListVal: [
					ext.TestNestedType{
						NestedListVal:['goodbye', 'cruel', 'world'],
						NestedMapVal: {42: true},
					},
				],
				MapVal: {'map-key': ext.TestAllTypes{BoolVal: true}},
			}`,
			out: &TestAllTypes{
				NestedVal:    &TestNestedType{NestedMapVal: map[int64]bool{1: false}},
				BoolVal:      true,
				BytesVal:     []byte("hello"),
				DurationVal:  time.Second * 5,
				DoubleVal:    1.5,
				FloatVal:     2.5,
				Int32Val:     10,
				Int64Val:     20,
				StringVal:    "hello world",
				TimestampVal: mustParseTime(t, "2011-08-06T01:23:45Z"),
				Uint32Val:    uint32(100),
				Uint64Val:    uint64(200),
				ListVal: []*TestNestedType{
					{
						NestedListVal: []string{"goodbye", "cruel", "world"},
						NestedMapVal:  map[int64]bool{42: true},
					},
				},
				MapVal: map[string]TestAllTypes{"map-key": {BoolVal: true}},
			},
		},
		{
			expr: `ext.TestAllTypes{					
					PbVal: test.TestAllTypes{single_int32: 123}
				}.PbVal`,
			out: &proto3pb.TestAllTypes{SingleInt32: 123},
		},
		{
			expr: `ext.TestAllTypes{PbVal: test.TestAllTypes{}} == 
			ext.TestAllTypes{PbVal: test.TestAllTypes{single_bool: false}}`,
		},
		{expr: `ext.TestNestedType{} == TestNestedType{}`},
		{expr: `ext.TestAllTypes{}.BoolVal != true`},
		{expr: `!has(ext.TestAllTypes{}.BoolVal) && !has(ext.TestAllTypes{}.NestedVal)`},
		{expr: `type(ext.TestAllTypes) == type`},
		{expr: `type(ext.TestAllTypes{}) == ext.TestAllTypes`},
		{expr: `type(ext.TestAllTypes{}) == ext.TestAllTypes`},
		{expr: `ext.TestAllTypes != test.TestAllTypes`},
		{expr: `ext.TestAllTypes{BoolVal: true} != dyn(test.TestAllTypes{single_bool: true})`},
		{expr: `ext.TestAllTypes{}.NestedVal == ext.TestNestedType{}`},
		{expr: `ext.TestNestedType{} == ext.TestAllTypes{}.NestedStructVal`},
		{expr: `ext.TestAllTypes{}.NestedStructVal == ext.TestNestedType{}`},
		{expr: `ext.TestAllTypes{}.ListVal.size() == 0`},
		{expr: `ext.TestAllTypes{}.MapVal.size() == 0`},
		{expr: `ext.TestAllTypes{}.TimestampVal == timestamp(0)`},
		{expr: `test.TestAllTypes{}.single_timestamp == timestamp(0)`},
		{expr: `[TestAllTypes{BoolVal: true}, TestAllTypes{BoolVal: false}].exists(t, t.BoolVal == true)`},
		{
			expr: `tests.all(t, t.Int32Val > 17)`,
			in: map[string]any{
				"tests": []*TestAllTypes{{Int32Val: 18}, {Int32Val: 19}, {Int32Val: 20}},
			},
		},
	}
	env := testNativeEnv(t)
	for i, tst := range nativeTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			cAst, iss := env.Check(pAst)
			if iss.Err() != nil {
				t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, cAst)
			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatal(err)
				}
				in := tc.in
				if in == nil {
					in = cel.NoVars()
				}
				out, _, err := prg.Eval(in)
				if err != nil {
					t.Fatal(err)
				}
				want := tc.out
				if want == nil {
					want = true
				}
				wantPB, isPB := want.(proto.Message)
				if isPB && !pb.Equal(wantPB, out.Value().(proto.Message)) {
					t.Errorf("got %v, wanted %v for expr: %s", out.Value(), want, tc.expr)
				}
				if !isPB && !reflect.DeepEqual(out.Value(), want) {
					t.Errorf("got %v, wanted %v for expr: %s", out.Value(), want, tc.expr)
				}
			}
		})
	}
}

func TestNativeTypesStaticErrors(t *testing.T) {
	var nativeTests = []struct {
		expr string
		err  string
	}{
		{
			expr: `TestAllTypos{}`,
			err: `ERROR: <input>:1:13: undeclared reference to 'TestAllTypos' (in container 'ext')
			 | TestAllTypos{}
			 | ............^`,
		},
		{
			expr: `ext.TestAllTypes{bool_val: false}`,
			err: `ERROR: <input>:1:26: undefined field 'bool_val'
			| ext.TestAllTypes{bool_val: false}
			| .........................^`,
		},
		{
			expr: `ext.TestAllTypes{UnsupportedVal: null}`,
			err: `ERROR: <input>:1:32: undefined field 'UnsupportedVal'
			| ext.TestAllTypes{UnsupportedVal: null}
			| ...............................^`,
		},
		{
			expr: `ext.TestAllTypes{UnsupportedListVal: null}`,
			err: `ERROR: <input>:1:36: undefined field 'UnsupportedListVal'
			| ext.TestAllTypes{UnsupportedListVal: null}
			| ...................................^`,
		},
		{
			expr: `ext.TestAllTypes{UnsupportedMapVal: null}`,
			err: `ERROR: <input>:1:35: undefined field 'UnsupportedMapVal'
			| ext.TestAllTypes{UnsupportedMapVal: null}
			| ..................................^`,
		},
	}
	env := testNativeEnv(t)
	for i, tst := range nativeTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			_, iss := env.Compile(tc.expr)
			if iss.Err() == nil {
				t.Fatalf("env.Compile(%v) succeeded, wanted error", tc.expr)
			}
			if !test.Compare(iss.Err().Error(), tc.err) {
				t.Errorf("env.Compile(%v) got %v, wanted error %s", tc.expr, iss.Err(), tc.err)
			}
		})
	}
}

func TestNativeTypesRuntimeErrors(t *testing.T) {
	var nativeTests = []struct {
		expr string
		err  string
	}{
		{
			expr: `TestAllTypos{}`,
			err:  `unknown type: TestAllTypos`,
		},
		{
			expr: `ext.TestAllTypes{bool_val: false}`,
			err:  `no such field: bool_val`,
		},
		{
			expr: `ext.TestAllTypes{UnsupportedVal: null}`,
			err:  `no such field: UnsupportedVal`,
		},
		{
			expr: `ext.TestAllTypes{UnsupportedListVal: null}`,
			err:  `no such field: UnsupportedListVal`,
		},
		{
			expr: `ext.TestAllTypes{UnsupportedMapVal: null}`,
			err:  `no such field: UnsupportedMapVal`,
		},
		{
			expr: `ext.TestAllTypes{privateVal: null}`,
			err:  `no such field: privateVal`,
		},
		{
			expr: `ext.TestAllTypes{}.UnsupportedMapVal`,
			err:  `no such field: UnsupportedMapVal`,
		},
		{
			expr: `ext.TestAllTypes{}.privateVal`,
			err:  `no such field: privateVal`,
		},
		{
			expr: `ext.TestAllTypes{BoolVal: 'false'}`,
			err:  `unsupported native conversion from string to 'bool'`,
		},
		{
			expr: `has(ext.TestAllTypes{}.BadFieldName)`,
			err:  `no such field: BadFieldName`,
		},
		{
			expr: `ext.TestAllTypes{}[42]`,
			err:  `no such overload`,
		},
		{
			expr: `ext.TestAllTypes{Int32Val: 9223372036854775807}`,
			err:  `integer overflow`,
		},
		{
			expr: `ext.TestAllTypes{Uint32Val: 9223372036854775807u}`,
			err:  `unsigned integer overflow`,
		},
	}
	env := testNativeEnv(t)
	for i, tst := range nativeTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			ast, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			prg, err := env.Program(ast)
			if err != nil {
				if !strings.Contains(err.Error(), tc.err) {
					t.Fatal(err)
				}
				return
			}
			out, _, err := prg.Eval(cel.NoVars())
			if err == nil || !strings.Contains(err.Error(), tc.err) {
				var got any = err
				if err == nil {
					got = out
				}
				t.Fatalf("prg.Eval() got %v, wanted error %v", got, tc.err)
			}
		})
	}
}

func TestNativeTypesErrors(t *testing.T) {
	envTests := []struct {
		nativeType any
		err        string
	}{
		{
			nativeType: reflect.TypeOf(1),
			err:        "unsupported reflect.Type",
		},
		{
			nativeType: reflect.ValueOf(1),
			err:        "unsupported reflect.Type",
		},
		{
			nativeType: 1,
			err:        "must be reflect.Type or reflect.Value",
		},
	}
	for i, tst := range envTests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			_, err := cel.NewEnv(NativeTypes(tc.nativeType))
			if err == nil || !strings.Contains(err.Error(), tc.err) {
				t.Errorf("cel.NewEnv(NativeTypes(%v)) got error %v, wanted %v", tc.nativeType, err, tc.err)
			}
		})
	}
}

func TestNativeTypesConvertToNative(t *testing.T) {
	env := testNativeEnv(t, NativeTypes(reflect.TypeOf(TestNestedType{})))
	adapter := env.TypeAdapter()
	conversions := []struct {
		in  any
		out any
		err string
	}{
		{
			in:  &TestAllTypes{BoolVal: true},
			out: &TestAllTypes{BoolVal: true},
		},
		{
			in:  TestAllTypes{BoolVal: true},
			out: &TestAllTypes{BoolVal: true},
		},
		{
			in:  &TestAllTypes{BoolVal: true},
			out: TestAllTypes{BoolVal: true},
		},
		{
			in:  nil,
			out: types.NullValue,
		},
		{
			in:  &TestAllTypes{BoolVal: true},
			out: &proto3pb.TestAllTypes{},
			err: "type conversion error",
		},
	}
	for _, c := range conversions {
		inVal := adapter.NativeToValue(c.in)
		if types.IsError(inVal) {
			t.Fatalf("adapter.NativeToValue(%v) failed: %v", c.in, inVal)
		}
		out, err := inVal.ConvertToNative(reflect.TypeOf(c.out))
		if err != nil {
			if c.err != "" {
				if !strings.Contains(err.Error(), c.err) {
					t.Fatalf("%v.ConvertToNative(%T) got %v, wanted error %v", c.in, c.out, err, c.err)
				}
				return
			}
			t.Fatalf("%v.ConvertToNative(%T) failed: %v", c.in, c.out, err)
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("%v.ConvertToNative(%T) got %v, wanted %v", c.in, c.out, out, c.out)
		}
	}
}

func TestNativeTypesConvertToExprTypeErrors(t *testing.T) {
	unsupportedTypes := []reflect.Type{
		reflect.TypeOf(make(map[string]chan string)),
		reflect.TypeOf(make([]chan int, 0)),
		reflect.TypeOf(make(map[chan int]bool, 0)),
	}
	for _, ut := range unsupportedTypes {
		if _, converted := convertToExprType(ut); converted {
			t.Errorf("convertToExprType(%v) succeeded when it should have failed", ut)
		}
	}
}

func TestConvertToTypeErrors(t *testing.T) {
	env := testNativeEnv(t, NativeTypes(reflect.TypeOf(TestNestedType{})))
	adapter := env.TypeAdapter()
	conversions := []struct {
		in  any
		out any
		err string
	}{
		{
			in:  &TestAllTypes{BoolVal: true},
			out: &TestAllTypes{BoolVal: true},
		},
		{
			in:  TestAllTypes{BoolVal: true},
			out: &TestAllTypes{BoolVal: true},
		},
		{
			in:  &TestAllTypes{BoolVal: true},
			out: TestAllTypes{BoolVal: true},
		},
		{
			in:  &TestAllTypes{BoolVal: true},
			out: &proto3pb.TestAllTypes{},
			err: "type conversion error",
		},
	}
	for _, c := range conversions {
		inVal := adapter.NativeToValue(c.in)
		outVal := adapter.NativeToValue(c.out)
		if types.IsError(inVal) {
			t.Fatalf("adapter.NativeToValue(%v) failed: %v", c.in, inVal)
		}
		if types.IsError(outVal) {
			t.Fatalf("adapter.NativeToValue(%v) failed: %v", c.out, outVal)
		}
		conv := inVal.ConvertToType(outVal.Type())
		if c.err != "" {
			if !types.IsError(conv) {
				t.Fatalf("%v.ConvertToType(%v) got %v, wanted error %v", c.in, outVal.Type(), conv, c.err)
			}
			convErr := conv.(*types.Err)
			if !strings.Contains(convErr.Error(), c.err) {
				t.Fatalf("%v.ConvertToType(%v) got %v, wanted error %v", c.in, outVal.Type(), conv, c.err)
			}
			return
		}
		if conv != inVal {
			t.Errorf("%v.ConvertToType(%v) got %v, wanted %v", c.in, outVal.Type(), conv, c.err)
		}
		conv = inVal.ConvertToType(types.TypeType)
		if conv.Type() != types.TypeType || conv.(ref.Type) != inVal.Type() {
			t.Errorf("%v.ConvertToType(Type) got %v, wanted %v", inVal, conv, inVal.Type())
		}
	}
}

func TestNativeTypesWithOptional(t *testing.T) {
	var nativeTests = []struct {
		expr string
	}{
		{expr: `!optional.ofNonZeroValue(ext.TestAllTypes{}).hasValue()`},
		{expr: `!ext.TestAllTypes{}.?BoolVal.orValue(false)`},
		{expr: `!ext.TestAllTypes{}.?BoolVal.hasValue()`},
		{expr: `!ext.TestAllTypes{BoolVal: false}.?BoolVal.hasValue()`},
		{expr: `ext.TestAllTypes{BoolVal: true}.?BoolVal.hasValue()`},
		{expr: `ext.TestAllTypes{}.NestedVal.?NestedMapVal.orValue({}).size() == 0`},
	}
	env := testNativeEnv(t, cel.OptionalTypes())
	for i, tst := range nativeTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			cAst, iss := env.Check(pAst)
			if iss.Err() != nil {
				t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, cAst)
			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatal(err)
				}
				out, _, err := prg.Eval(cel.NoVars())
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(out.Value(), true) {
					t.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}

func TestNativeTypeConvertToType(t *testing.T) {
	nt, err := newNativeType(reflect.TypeOf(&TestAllTypes{}))
	if err != nil {
		t.Fatalf("newNativeType() failed: %v", err)
	}
	if nt.ConvertToType(types.TypeType) != types.TypeType {
		t.Error("ConvertToType(Type) failed")
	}
	if !types.IsError(nt.ConvertToType(types.StringType)) {
		t.Errorf("ConvertToType(String) got %v, wanted error", nt.ConvertToType(types.StringType))
	}
}

func TestNativeTypeConvertToNative(t *testing.T) {
	nt, err := newNativeType(reflect.TypeOf(&TestAllTypes{}))
	if err != nil {
		t.Fatalf("newNativeType() failed: %v", err)
	}
	out, err := nt.ConvertToNative(reflect.TypeOf(1))
	if err == nil {
		t.Errorf("nt.ConvertToNative(1) produced %v, wanted error", out)
	}
}

func TestNativeTypeHasTrait(t *testing.T) {
	nt, err := newNativeType(reflect.TypeOf(&TestAllTypes{}))
	if err != nil {
		t.Fatalf("newNativeType() failed: %v", err)
	}
	if !nt.HasTrait(traits.IndexerType) || !nt.HasTrait(traits.FieldTesterType) {
		t.Error("nt.HasTrait() failed indicate support for presence test and field access.")
	}
}

func TestNativeTypeValue(t *testing.T) {
	nt, err := newNativeType(reflect.TypeOf(&TestAllTypes{}))
	if err != nil {
		t.Fatalf("newNativeType() failed: %v", err)
	}
	if nt.Value() != nt.String() {
		t.Errorf("nt.Value() got %v, wanted %v", nt.Value(), nt.String())
	}
}

// testEnv initializes the test environment common to all tests.
func testNativeEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	envOpts := []cel.EnvOption{
		cel.Container("ext"),
		cel.Abbrevs("google.expr.proto3.test"),
		cel.Types(&proto3pb.TestAllTypes{}),
		cel.Variable("tests", cel.ListType(cel.ObjectType("ext.TestAllTypes"))),
	}
	envOpts = append(envOpts, opts...)
	envOpts = append(envOpts,
		NativeTypes(
			reflect.TypeOf(&TestNestedType{}),
			reflect.ValueOf(&TestAllTypes{}),
		),
	)
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		t.Fatalf("cel.NewEnv(NativeTypes()) failed: %v", err)
	}
	return env
}

func mustParseTime(t *testing.T, timestamp string) time.Time {
	t.Helper()
	out, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("time.Parse(%q) failed: %v", timestamp, err)
	}
	return out
}

type TestNestedType struct {
	NestedListVal []string
	NestedMapVal  map[int64]bool
}

type TestAllTypes struct {
	NestedVal       *TestNestedType
	NestedStructVal TestNestedType
	BoolVal         bool
	BytesVal        []byte
	DurationVal     time.Duration
	DoubleVal       float64
	FloatVal        float32
	Int32Val        int32
	Int64Val        int32
	StringVal       string
	TimestampVal    time.Time
	Uint32Val       uint32
	Uint64Val       uint64
	ListVal         []*TestNestedType
	MapVal          map[string]TestAllTypes
	PbVal           *proto3pb.TestAllTypes

	// channel types are not supported
	UnsupportedVal     chan string
	UnsupportedListVal []chan string
	UnsupportedMapVal  map[int]chan string

	// unexported types can be found but not set or accessed
	privateVal map[string]string
}
