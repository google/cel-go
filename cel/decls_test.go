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

package cel

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/checker"
	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var dispatchTests = []struct {
	expr string
	out  any
}{
	{
		expr: "max(-1)",
		out:  types.IntNegOne,
	},
	{
		expr: "max(-1, 0)",
		out:  types.IntZero,
	},
	{
		expr: "max(-1, 0, 1)",
		out:  types.IntOne,
	},
	{
		expr: "max(dyn(1.2))",
		out:  fmt.Errorf("no such overload: max(double)"),
	},

	{
		expr: "max(1, dyn(1.2))",
		out:  fmt.Errorf("no such overload: max(int, double)"),
	},
	{
		expr: "max(1, 2, dyn(1.2))",
		out:  fmt.Errorf("no such overload: max(int, int, double)"),
	},
	{
		expr: "max(err, 1)",
		out:  fmt.Errorf("error argument"),
	},
	{
		expr: "max(err, unk)",
		out:  fmt.Errorf("error argument"),
	},
	{
		expr: "max(unk, unk)",
		out:  types.Unknown{42},
	},
	{
		expr: "max(unk, unk, unk)",
		out:  types.Unknown{42},
	},
}

func TestFunctionMerge(t *testing.T) {
	size := Function("size",
		Overload("size_map", []*Type{MapType(TypeParamType("K"), TypeParamType("V"))}, IntType),
		Overload("size_list", []*Type{ListType(TypeParamType("V"))}, IntType),
		Overload("size_string", []*Type{StringType}, IntType),
		Overload("size_bytes", []*Type{BytesType}, IntType),
		MemberOverload("map_size", []*Type{MapType(TypeParamType("K"), TypeParamType("V"))}, IntType),
		MemberOverload("list_size", []*Type{ListType(TypeParamType("V"))}, IntType),
		MemberOverload("string_size", []*Type{StringType}, IntType),
		MemberOverload("bytes_size", []*Type{BytesType}, IntType),
		SingletonUnaryBinding(func(arg ref.Val) ref.Val {
			return arg.(traits.Sizer).Size()
		}, traits.SizerType),
	)
	// Note, the size() implementation is inherited from the singleton implementation for the size function.
	// It is possible to redefine the singleton, but the singleton approach is incompatible with specialized
	// overloads as the singleton is compatible with the parse-only approach and compiled approach; but
	// a mix of singleton and specialized overloads might result in a singleton which does not encompass
	// dynamic dispatch to all possible overloads.
	sizeExt := Function("size",
		Overload("size_vector", []*Type{OpaqueType("vector", TypeParamType("V"))}, IntType),
		MemberOverload("vector_size", []*Type{OpaqueType("vector", TypeParamType("V"))}, IntType))

	vectorExt := Function("vector",
		Overload("vector_list", []*Type{ListType(TypeParamType("V"))}, OpaqueType("vector", TypeParamType("V")),
			UnaryBinding(func(list ref.Val) ref.Val {
				return list
			})))

	eq := Function(operators.Equals,
		Overload(operators.Equals, []*Type{TypeParamType("T"), TypeParamType("T")}, BoolType,
			BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs.Equal(rhs)
			})))

	e, err := NewCustomEnv(size, sizeExt, vectorExt, eq)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := e.Compile(`[[0].size() == 1,
		{'a': true, 'b': false}.size() == 2,
		'hello'.size() == 5,
		b'hello'.size() == 5,
		vector([1.2, 2.3, 3.4]).size() == 3]`)
	if iss.Err() != nil {
		t.Fatalf("Compile() errored: %v", iss.Err())
	}
	if !reflect.DeepEqual(ast.OutputType(), ListType(BoolType)) {
		t.Errorf("Compile() produced a %v, wanted bool", ast.OutputType())
	}
	if ast.OutputType().String() != "list(bool)" {
		t.Errorf("ast.OutputType().String() produced %s, wanted list(bool)", ast.OutputType())
	}
	prg, err := e.Program(ast)
	if err != nil {
		t.Fatalf("e.Program(ast) failed: %v", err)
	}
	out, _, err := prg.Eval(NoVars())
	if err != nil {
		t.Errorf("prg.Eval() errored: %v", err)
	}
	want := types.DefaultTypeAdapter.NativeToValue([]bool{true, true, true, true, true})
	if out.Equal(want) != types.True {
		t.Errorf("prg.Eval() got %v, wanted %v", out, want)
	}

	_, err = NewCustomEnv(size, size)
	if err == nil || !strings.Contains(err.Error(), "already has singleton binding") {
		t.Errorf("NewCustomEnv(size, size) did not produce the expected error: %v", err)
	}
	e, err = NewCustomEnv(size,
		Function("size",
			Overload("size_int", []*Type{IntType}, IntType,
				UnaryBinding(func(arg ref.Val) ref.Val { return types.Int(2) }),
			),
		),
	)
	if err != nil {
		t.Fatalf("NewCustomEnv(size, <custom>) failed: %v", err)
	}
	_, err = e.Program(ast)
	if err == nil || !strings.Contains(err.Error(), "incompatible with specialized overloads") {
		t.Errorf("NewCustomEnv(size, size_specialization) did not produce the expected error: %v", err)
	}
}

func TestFunctionMergeDuplicate(t *testing.T) {
	maxFunc := Function("max",
		Overload("max_int", []*Type{IntType}, IntType),
		Overload("max_int", []*Type{IntType}, IntType),
	)
	_, err := NewCustomEnv(maxFunc, maxFunc)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}
}

func TestFunctionMergeDeclarationAndDefinition(t *testing.T) {
	idFuncDecl := Function("id", Overload("id", []*Type{TypeParamType("T")}, TypeParamType("T"), OverloadIsNonStrict()))
	idFuncDef := Function("id",
		Overload("id", []*Type{TypeParamType("T")}, TypeParamType("T"), OverloadIsNonStrict(),
			UnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			}),
		),
	)
	env, err := NewCustomEnv(Variable("x", AnyType), idFuncDecl, idFuncDef)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := env.Compile("id(x)")
	if iss.Err() != nil {
		t.Fatalf("env.Compile(id(x)) failed: %v", iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		t.Fatalf("env.Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": true})
	if err != nil {
		t.Fatalf("prg.Eval(x: true) errored: %v", err)
	}
	if out != types.True {
		t.Errorf("prg.Eval() yielded %v, wanted true", out)
	}
}

func TestFunctionMergeCollision(t *testing.T) {
	maxFunc := Function("max",
		Overload("max_int", []*Type{IntType}, IntType),
		Overload("max_int2", []*Type{IntType}, IntType),
	)
	_, err := NewCustomEnv(maxFunc, maxFunc)
	if err == nil {
		t.Fatal("NewCustomEnv() succeeded, wanted collision error")
	}
}

func TestFunctionNoOverloads(t *testing.T) {
	_, err := NewCustomEnv(
		Function("right", SingletonBinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
			return arg2
		})))
	if err == nil || !strings.Contains(err.Error(), "must have at least one overload") {
		t.Errorf("got %v, wanted 'must have at least one overload'", err)
	}
}

func TestSingletonUnaryBinding(t *testing.T) {
	// Test the case where the singleton unary impl is merged with the earlier declaration.
	e, err := NewCustomEnv(
		Variable("x", AnyType),
		Function("id", Overload("id_any", []*Type{AnyType}, AnyType)),
		Function("id", Overload("id_any", []*Type{AnyType}, AnyType), SingletonUnaryBinding(func(arg ref.Val) ref.Val { return arg })),
	)
	if err != nil {
		t.Errorf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := e.Parse("id(x)")
	if iss.Err() != nil {
		t.Fatalf("Parse('id(x)') failed: %v", iss.Err())
	}
	prg, err := e.Program(ast)
	if err != nil {
		t.Fatalf("Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": "hello"})
	if err != nil {
		t.Errorf("prg.Eval(x=hello) failed: %v", err)
	}
	if out.Equal(types.String("hello")) != types.True {
		t.Errorf("Eval got %v, wanted 'hello'", out)
	}
}

func TestSingletonUnaryBindingParameterized(t *testing.T) {
	// Test the case where the singleton unary impl is merged with the earlier declaration.
	e, err := NewCustomEnv(
		Variable("x", AnyType),
		Function("isSorted", MemberOverload("list_int_is_sorted", []*Type{ListType(IntType)}, BoolType)),
		Function("isSorted", MemberOverload("list_uint_is_sorted", []*Type{ListType(UintType)}, BoolType),
			SingletonUnaryBinding(func(arg ref.Val) ref.Val { return types.True })),
	)
	if err != nil {
		t.Errorf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := e.Parse("x.isSorted()")
	if iss.Err() != nil {
		t.Fatalf("Parse('x.isSorted()') failed: %v", iss.Err())
	}
	prg, err := e.Program(ast)
	if err != nil {
		t.Fatalf("Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": []int{1, 2, 3}})
	if err != nil {
		t.Errorf("prg.Eval(x=[1,2,3]) failed: %v", err)
	}
	if out != types.True {
		t.Errorf("Eval got %v, wanted true", out)
	}
}

func TestSingletonBinaryBinding(t *testing.T) {
	_, err := NewCustomEnv(
		Function("right",
			Overload("right_int_int", []*Type{IntType, IntType}, IntType),
			Overload("right_double_double", []*Type{DoubleType, DoubleType}, DoubleType),
			Overload("right_string_string", []*Type{StringType, StringType}, StringType),
			SingletonBinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
				return arg2
			}, traits.ComparerType),
		),
	)
	if err != nil {
		t.Errorf("NewCustomEnv() failed: %v", err)
	}
}

func TestSingletonFunctionBinding(t *testing.T) {
	env, err := NewCustomEnv(
		Variable("unk", DynType),
		Variable("err", DynType),
		Function("dyn",
			Overload("dyn", []*Type{DynType}, DynType),
			SingletonUnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			})),
		Function("max",
			Overload("max_int", []*Type{IntType}, IntType),
			Overload("max_int_int", []*Type{IntType, IntType}, IntType),
			Overload("max_int_int_int", []*Type{IntType, IntType, IntType}, IntType),
			SingletonFunctionBinding(func(args ...ref.Val) ref.Val {
				max := types.Int(math.MinInt64)
				for _, arg := range args {
					i, ok := arg.(types.Int)
					if !ok {
						// With a singleton implementation, the error handling must be explicitly
						// performed as a binding detail of the singleton function.
						// With custom overload implementations, a function guard is automatically
						// added to the function to validate that the runtime types are compatible
						// to provide some basic invocation protections.
						return decls.MaybeNoSuchOverload("max", args...)
					}
					if i > max {
						max = i
					}
				}
				return max
			}),
		),
	)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}

	for _, tc := range dispatchTests {
		tc := tc
		t.Run(fmt.Sprintf("Parse(%s)", tc.expr), func(t *testing.T) {
			testParse(t, env, tc.expr, tc.out)
		})
	}
	for _, tc := range dispatchTests {
		tc := tc
		t.Run(fmt.Sprintf("Compile(%s)", tc.expr), func(t *testing.T) {
			testCompile(t, env, tc.expr, tc.out)
		})
	}
}

func TestUnaryBinding(t *testing.T) {
	_, err := NewCustomEnv(
		Function("dyn",
			Overload("dyn", []*Type{}, DynType,
				// wrong arg count
				UnaryBinding(func(arg ref.Val) ref.Val { return arg }),
			),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "function bound to non-unary overload") {
		t.Errorf("bad binding did not produce expected error: %v", err)
	}

	e, err := NewCustomEnv(
		Function("size",
			Overload("size_non_strict", []*Type{ListType(DynType)}, IntType,
				OverloadIsNonStrict(),
				OverloadOperandTrait(traits.SizerType),
				UnaryBinding(func(arg ref.Val) ref.Val {
					if types.IsUnknownOrError(arg) {
						return arg
					}
					return arg.(traits.Sizer).Size()
				}),
			),
		),
		Variable("x", ListType(DynType)),
	)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := e.Compile("size(x)")
	if iss.Err() != nil {
		t.Fatalf("Compile(size(x)) failed: %v", iss.Err())
	}
	prg, err := e.Program(ast)
	if err != nil {
		t.Fatalf("Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": types.Unknown{1}})
	if err != nil {
		t.Fatalf("prg.Eval(x=unk) failed: %v", err)
	}
	if !reflect.DeepEqual(out, types.Unknown{1}) {
		t.Errorf("prg.Eval(x=unk) returned %v, wanted unknown{1}", out)
	}
}

func TestBinaryBinding(t *testing.T) {
	e, err := NewCustomEnv(
		Function("max",
			Overload("max_int_int", []*Type{IntType, IntType}, IntType,
				OverloadIsNonStrict(),
				BinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
					if types.IsUnknownOrError(arg1) {
						return arg2
					}
					if types.IsUnknownOrError(arg2) {
						return arg1
					}
					i1 := arg1.(types.Int)
					i2 := arg2.(types.Int)
					if i1 > i2 {
						return i1
					}
					return i2
				}),
			),
		),
		Variable("x", IntType),
		Variable("y", IntType),
	)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := e.Parse("max(x, y)")
	if iss.Err() != nil {
		t.Fatalf("Parse(max(x, y))) failed: %v", iss.Err())
	}
	prg, err := e.Program(ast)
	if err != nil {
		t.Fatalf("Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": types.Unknown{1}, "y": 1})
	if err != nil {
		t.Fatalf("prg.Eval(x=unk) failed: %v", err)
	}
	if !reflect.DeepEqual(out, types.IntOne) {
		t.Errorf("prg.Eval(x=unk, y=1) returned %v, wanted 1", out)
	}
	out, _, err = prg.Eval(map[string]any{"x": 2, "y": types.Unknown{2}})
	if err != nil {
		t.Fatalf("prg.Eval(x=2, y=unk) failed: %v", err)
	}
	if !reflect.DeepEqual(out, types.Int(2)) {
		t.Errorf("prg.Eval(x=2, y=unk) returned %v, wanted 2", out)
	}
	out, _, err = prg.Eval(map[string]any{"x": 2, "y": 1})
	if err != nil {
		t.Fatalf("prg.Eval(x=2, y=1) failed: %v", err)
	}
	if !reflect.DeepEqual(out, types.Int(2)) {
		t.Errorf("prg.Eval(x=2, y=1) returned %v, wanted 2", out)
	}

	_, err = NewCustomEnv(
		Function("right",
			Overload("right_int_int", []*Type{IntType, IntType, IntType}, IntType,
				// wrong arg count
				BinaryBinding(func(arg1, arg2 ref.Val) ref.Val { return arg2 }),
			),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "function bound to non-binary overload") {
		t.Errorf("bad binding did not produce expected error: %v", err)
	}
}

func TestFunctionBinding(t *testing.T) {
	e, err := NewCustomEnv(
		Variable("unk", DynType),
		Variable("err", DynType),
		Function("dyn",
			Overload("dyn", []*Type{DynType}, DynType),
			SingletonUnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			})),
		Function("max",
			Overload("max_int", []*Type{IntType}, IntType,
				UnaryBinding(func(arg ref.Val) ref.Val { return arg })),
			Overload("max_int_int", []*Type{IntType, IntType}, IntType,
				BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
					if lhs.(types.Int).Compare(rhs) == types.IntNegOne {
						return rhs
					}
					return lhs
				})),
			Overload("max_int_int_int", []*Type{IntType, IntType, IntType}, IntType,
				FunctionBinding(func(args ...ref.Val) ref.Val {
					max := types.Int(math.MinInt64)
					for _, arg := range args {
						i := arg.(types.Int)
						if i > max {
							max = i
						}
					}
					return max
				})),
		),
	)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}

	for _, tc := range dispatchTests {
		tc := tc
		t.Run(fmt.Sprintf("Parse(%s)", tc.expr), func(t *testing.T) {
			testParse(t, e, tc.expr, tc.out)
		})
	}
	for _, tc := range dispatchTests {
		tc := tc
		t.Run(fmt.Sprintf("Compile(%s)", tc.expr), func(t *testing.T) {
			testCompile(t, e, tc.expr, tc.out)
		})
	}
}

func TestTypeToExprType(t *testing.T) {
	tests := []struct {
		in             *Type
		out            *exprpb.Type
		unidirectional bool
	}{
		{
			in:  OpaqueType("vector", DoubleType, DoubleType),
			out: chkdecls.NewAbstractType("vector", chkdecls.Double, chkdecls.Double),
		},
		{
			in:  AnyType,
			out: chkdecls.Any,
		},
		{
			in:  BoolType,
			out: chkdecls.Bool,
		},
		{
			in:  BytesType,
			out: chkdecls.Bytes,
		},
		{
			in:  DoubleType,
			out: chkdecls.Double,
		},
		{
			in:  DurationType,
			out: chkdecls.Duration,
		},
		{
			in:  DynType,
			out: chkdecls.Dyn,
		},
		{
			in:  IntType,
			out: chkdecls.Int,
		},
		{
			in:  ListType(TypeParamType("T")),
			out: chkdecls.NewListType(chkdecls.NewTypeParamType("T")),
		},
		{
			in:  MapType(TypeParamType("K"), TypeParamType("V")),
			out: chkdecls.NewMapType(chkdecls.NewTypeParamType("K"), chkdecls.NewTypeParamType("V")),
		},
		{
			in:  NullType,
			out: chkdecls.Null,
		},
		{
			in:  ObjectType("google.type.Expr"),
			out: chkdecls.NewObjectType("google.type.Expr"),
		},
		{
			in:  StringType,
			out: chkdecls.String,
		},
		{
			in:  TimestampType,
			out: chkdecls.Timestamp,
		},
		{
			in:  TypeType,
			out: chkdecls.NewTypeType(chkdecls.Dyn),
		},
		{
			in:  UintType,
			out: chkdecls.Uint,
		},
		{
			in:  NullableType(BoolType),
			out: chkdecls.NewWrapperType(chkdecls.Bool),
		},
		{
			in:  NullableType(BytesType),
			out: chkdecls.NewWrapperType(chkdecls.Bytes),
		},
		{
			in:  NullableType(DoubleType),
			out: chkdecls.NewWrapperType(chkdecls.Double),
		},
		{
			in:  NullableType(IntType),
			out: chkdecls.NewWrapperType(chkdecls.Int),
		},
		{
			in:  NullableType(StringType),
			out: chkdecls.NewWrapperType(chkdecls.String),
		},
		{
			in:  NullableType(UintType),
			out: chkdecls.NewWrapperType(chkdecls.Uint),
		},
		{
			in:             ObjectType("google.protobuf.Any"),
			out:            chkdecls.Any,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Duration"),
			out:            chkdecls.Duration,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Timestamp"),
			out:            chkdecls.Timestamp,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Value"),
			out:            chkdecls.Dyn,
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.ListValue"),
			out:            chkdecls.NewListType(chkdecls.Dyn),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Struct"),
			out:            chkdecls.NewMapType(chkdecls.String, chkdecls.Dyn),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.BoolValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Bool),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.BytesValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Bytes),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.DoubleValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Double),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.FloatValue"),
			out:            chkdecls.NewWrapperType(chkdecls.Double),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Int32Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Int),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.Int64Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Int),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.StringValue"),
			out:            chkdecls.NewWrapperType(chkdecls.String),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.UInt32Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Uint),
			unidirectional: true,
		},
		{
			in:             ObjectType("google.protobuf.UInt64Value"),
			out:            chkdecls.NewWrapperType(chkdecls.Uint),
			unidirectional: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			got, err := TypeToExprType(tc.in)
			if err != nil {
				t.Fatalf("TypeToExprType(%v) failed: %v", tc.in, err)
			}
			if !proto.Equal(got, tc.out) {
				t.Errorf("TypeToExprType(%v) returned %v, wanted %v", tc.in, got, tc.out)
			}
			if tc.unidirectional {
				return
			}
			roundTrip, err := ExprTypeToType(got)
			if err != nil {
				t.Fatalf("ExprTypeToType(%v) failed: %v", got, err)
			}
			if !tc.in.IsType(roundTrip) {
				t.Errorf("ExprTypeToType(%v) returned %v, wanted %v", got, roundTrip, tc.in)
			}
		})
	}
}

func TestExprTypeToType(t *testing.T) {
	tests := []struct {
		in  *exprpb.Type
		out *Type
	}{
		{
			in:  chkdecls.NewObjectType("google.protobuf.Any"),
			out: AnyType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Duration"),
			out: DurationType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Timestamp"),
			out: TimestampType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Value"),
			out: DynType,
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.ListValue"),
			out: ListType(DynType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Struct"),
			out: MapType(StringType, DynType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.BoolValue"),
			out: NullableType(BoolType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.BytesValue"),
			out: NullableType(BytesType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.DoubleValue"),
			out: NullableType(DoubleType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.FloatValue"),
			out: NullableType(DoubleType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Int32Value"),
			out: NullableType(IntType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.Int64Value"),
			out: NullableType(IntType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.StringValue"),
			out: NullableType(StringType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.UInt32Value"),
			out: NullableType(UintType),
		},
		{
			in:  chkdecls.NewObjectType("google.protobuf.UInt64Value"),
			out: NullableType(UintType),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			got, err := ExprTypeToType(tc.in)
			if err != nil {
				t.Fatalf("ExprTypeToType(%v) failed: %v", tc.in, err)
			}
			if !got.IsType(tc.out) {
				t.Errorf("ExprTypeToType(%v) returned %v, wanted %v", tc.in, got, tc.out)
			}
		})
	}
}

func TestExprTypeToTypeInvalid(t *testing.T) {
	tests := []struct {
		in  *exprpb.Type
		out string
	}{
		{
			in:  &exprpb.Type{},
			out: "unsupported type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_Primitive{}},
			out: "unsupported primitive type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_WellKnown{}},
			out: "unsupported well-known type",
		},
		{
			in:  chkdecls.NewListType(&exprpb.Type{}),
			out: "unsupported type",
		},
		{
			in:  chkdecls.NewMapType(&exprpb.Type{}, chkdecls.Dyn),
			out: "unsupported type",
		},
		{
			in:  chkdecls.NewMapType(chkdecls.Dyn, &exprpb.Type{}),
			out: "unsupported type",
		},
		{
			in:  chkdecls.NewAbstractType("bad", &exprpb.Type{}),
			out: "unsupported type",
		},
		{
			in:  &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{}},
			out: "unsupported primitive type",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			_, err := ExprTypeToType(tc.in)
			if err == nil || !strings.Contains(err.Error(), tc.out) {
				t.Fatalf("ExprTypeToType(%v) got %v, wanted error %v", tc.in, err, tc.out)
			}
		})
	}
}

func TestExprDeclToDeclaration(t *testing.T) {
	size, err := ExprDeclToDeclaration(
		chkdecls.NewFunction("size", chkdecls.NewOverload("size_string", []*exprpb.Type{chkdecls.String}, chkdecls.Int)),
	)
	if err != nil {
		t.Fatalf("ExprDeclToDeclaration(size) failed: %v", err)
	}
	x, err := ExprDeclToDeclaration(chkdecls.NewVar("x", chkdecls.String))
	if err != nil {
		t.Fatalf("ExprDeclToDeclaration(x) failed: %v", err)
	}
	e, err := NewCustomEnv(size, x)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed: %v", err)
	}
	ast, iss := e.Compile("size(x)")
	if iss.Err() != nil {
		t.Fatalf("Compile(size(x)) failed: %v", iss.Err())
	}
	prg, err := e.Program(ast, Functions(&functions.Overload{
		Operator: overloads.SizeString,
		Unary: func(arg ref.Val) ref.Val {
			str, ok := arg.(types.String)
			if !ok {
				return types.MaybeNoSuchOverloadErr(arg)
			}
			return types.Int(len([]rune(string(str))))
		},
	}))
	if err != nil {
		t.Fatalf("Program(ast) failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"x": "hello"})
	if err != nil {
		t.Fatalf("prg.Eval(x=hello) failed: %v", err)
	}
	if out.Equal(types.Int(5)) != types.True {
		t.Errorf("prg.Eval(size(x)) got %v, wanted 5", out)
	}
}

func TestExprDeclToDeclarationInvalid(t *testing.T) {
	tests := []struct {
		in  *exprpb.Decl
		out string
	}{
		{
			in:  &exprpb.Decl{},
			out: "unsupported decl",
		},
		{
			in: &exprpb.Decl{
				Name: "bad_var",
				DeclKind: &exprpb.Decl_Ident{
					Ident: &exprpb.Decl_IdentDecl{
						Type: chkdecls.NewListType(&exprpb.Type{}),
					},
				},
			},
			out: "unsupported type",
		},
		{
			in: &exprpb.Decl{
				Name: "bad_func_return",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId: "bad_overload",
								ResultType: &exprpb.Type{},
							},
						},
					},
				},
			},
			out: "unsupported type",
		},
		{
			in: &exprpb.Decl{
				Name: "bad_func_arg",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId: "bad_overload",
								Params:     []*exprpb.Type{{}},
								ResultType: chkdecls.Dyn,
							},
						},
					},
				},
			},
			out: "unsupported type",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.in.String(), func(t *testing.T) {
			_, err := ExprDeclToDeclaration(tc.in)
			if err == nil || !strings.Contains(err.Error(), tc.out) {
				t.Fatalf("ExprDeclToDeclarations(%v) got %v, wanted error %v", tc.in, err, tc.out)
			}
		})
	}
}

func testParse(t testing.TB, env *Env, expr string, want any) {
	t.Helper()
	ast, iss := env.Parse(expr)
	if iss.Err() != nil {
		t.Fatalf("env.Parse(%s) failed: %v", expr, iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		t.Fatalf("env.Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"err": types.NewErr("error argument"), "unk": types.Unknown{42}})
	switch want := want.(type) {
	case types.Unknown:
		if !reflect.DeepEqual(want, out.(types.Unknown)) {
			t.Errorf("prg.Eval() got %v, wanted %v", out, want)
		}
	case ref.Val:
		if want.Equal(out) != types.True {
			t.Errorf("prg.Eval() got %v, wanted %v", out, want)
		}
	case error:
		if err == nil || want.Error() != err.Error() {
			t.Errorf("prg.Eval() got error '%v', wanted '%v'", err, want)
		}
	}
}

func testCompile(t testing.TB, env *Env, expr string, want any) {
	t.Helper()
	ast, iss := env.Compile(expr)
	if iss.Err() != nil {
		t.Fatalf("env.Compile(%s) failed: %v", expr, iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		t.Fatalf("env.Program() failed: %v", err)
	}
	out, _, err := prg.Eval(map[string]any{"err": types.NewErr("error argument"), "unk": types.Unknown{42}})
	switch want := want.(type) {
	case types.Unknown:
		if !reflect.DeepEqual(want, out.(types.Unknown)) {
			t.Errorf("prg.Eval() got %v, wanted %v", out, want)
		}
	case ref.Val:
		if want.Equal(out) != types.True {
			t.Errorf("prg.Eval() got %v, wanted %v", out, want)
		}
	case error:
		if err == nil || want.Error() != err.Error() {
			t.Errorf("prg.Eval() got error '%v', wanted '%v'", err, want)
		}
	}
}

func TestStdLibDeclarations(t *testing.T) {
	stdlibFuncs, err := stdlibFuncDecls()
	if err != nil {
		t.Fatalf("stdlibFuncDecls() failed: %v", err)
	}
	_, err = NewCustomEnv(
		Declarations(checker.StandardFunctions()...),
		Declarations(stdlibFuncs...),
	)
	if err != nil {
		t.Fatalf("NewCustomEnv() failed when combining stdlib functions and checker mirrors: %v", err)
	}
	stdlibFuncMap := make(map[string]*exprpb.Decl, len(stdlibFuncs))
	for _, fn := range stdlibFuncs {
		stdlibFuncMap[fn.GetName()] = fn
	}
	for _, fn := range checker.StandardFunctions() {
		stdFn, found := stdlibFuncMap[fn.GetName()]
		if !found {
			t.Errorf("function not found in stdlib: %v", fn)
		}
		if len(stdFn.GetFunction().GetOverloads()) != len(fn.GetFunction().GetOverloads()) {
			t.Errorf("stdlib function and checker function have different overload sets. std: %v, checker: %v", stdFn, fn)
		}
		stdOverloadMap := make(map[string]*exprpb.Decl_FunctionDecl_Overload)
		for _, o := range stdFn.GetFunction().GetOverloads() {
			stdOverloadMap[o.GetOverloadId()] = o
		}
		for _, o := range fn.GetFunction().GetOverloads() {
			stdOverload, found := stdOverloadMap[o.GetOverloadId()]
			if !found {
				t.Errorf("overload not found in stdlib: %v", o)
			}
			if len(stdOverload.GetTypeParams()) != len(o.GetTypeParams()) {
				t.Errorf("type params differ. std: %v, checker: %v", stdOverload, o)
			}
			if stdOverload.GetIsInstanceFunction() != o.GetIsInstanceFunction() {
				t.Errorf("instance function settings differ. std: %v, checker: %v", stdOverload, o)
			}
			if len(stdOverload.GetParams()) != len(o.GetParams()) {
				t.Errorf("arg counts differ. std: %v, checker: %v", stdOverload, o)
			}
			for i, p := range stdOverload.GetParams() {
				if !proto.Equal(p, o.GetParams()[i]) {
					t.Errorf("arg types differ. std: %v, checker: %v", stdOverload.GetResultType(), o.GetResultType())
				}
			}
			if !proto.Equal(stdOverload.GetResultType(), o.GetResultType()) {
				t.Errorf("result types differ. std: %v, checker: %v", stdOverload.GetResultType(), o.GetResultType())
			}
			delete(stdOverloadMap, o.GetOverloadId())
		}
		delete(stdlibFuncMap, fn.GetName())
	}
	for _, fn := range stdlibFuncMap {
		t.Errorf("function in stdlib but not checker standard functions: %v", fn)
	}
}
