// Copyright 2023 Google LLC
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

package decls

import (
	"strings"
	"testing"
	"time"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

func TestFunctionBindings(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	bindings, err := sizeFunc.Bindings()
	if err != nil {
		t.Fatalf("sizeFunc.Bindings() produced an err: %v", err)
	}
	if len(bindings) != 0 {
		t.Errorf("sizeFunc.Bindings() produced %d bindings, wanted none", len(bindings))
	}
	sizeFuncDef, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType,
			UnaryBinding(func(list ref.Val) ref.Val {
				sizer := list.(traits.Sizer)
				return sizer.Size()
			}),
		),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeMerged, err := sizeFunc.Merge(sizeFuncDef)
	if err != nil {
		t.Fatalf("Merge() failed: %v", err)
	}
	bindings, err = sizeMerged.Bindings()
	if err != nil {
		t.Fatalf("sizeFunc.Bindings() produced an err: %v", err)
	}
	if len(bindings) != 2 {
		t.Errorf("sizeFunc.Bindings() got %d bindings, wanted 2", len(bindings))
	}
	empty := types.DefaultTypeAdapter.NativeToValue([]string{})
	in := types.DefaultTypeAdapter.NativeToValue([]string{"1", "2"})
	for _, binding := range bindings {
		if binding.Unary == nil {
			t.Errorf("binding missing unary implementation: %v", binding.Operator)
			continue
		}
		if binding.Unary(in) != types.Int(2) {
			t.Errorf("binding invocation got %v, wanted 2", binding.Unary(in))
		}
		if binding.Unary(empty) != types.IntZero {
			t.Errorf("binding invocation got %v, wanted 0", binding.Unary(empty))
		}
	}
}

func TestFunctionVariableArgBindings(t *testing.T) {
	splitImpl := func(str, delim string, count int64) ref.Val {
		return types.DefaultTypeAdapter.NativeToValue(strings.SplitN(str, delim, int(count)))
	}
	splitFunc, err := NewFunction("split",
		MemberOverload("string_split", []*Type{StringType}, ListType(StringType),
			UnaryBinding(func(str ref.Val) ref.Val {
				s := str.(types.String)
				return splitImpl(string(s), "", -1)
			})),
		MemberOverload("string_split_string", []*Type{StringType, StringType}, ListType(StringType),
			BinaryBinding(func(str, sep ref.Val) ref.Val {
				s := str.(types.String)
				delim := sep.(types.String)
				return splitImpl(string(s), string(delim), -1)
			})),
		MemberOverload("string_split_string_int", []*Type{StringType, StringType, IntType}, ListType(StringType),
			FunctionBinding(func(args ...ref.Val) ref.Val {
				s := args[0].(types.String)
				delim := args[1].(types.String)
				count := args[2].(types.Int)
				return splitImpl(string(s), string(delim), int64(count))
			})),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	bindings, err := splitFunc.Bindings()
	if err != nil {
		t.Fatalf("sizeFunc.Bindings() produced an err: %v", err)
	}
	if len(bindings) != 4 {
		t.Errorf("sizeFunc.Bindings() produced %d bindings, wanted 4", len(bindings))
	}
	in := types.String("hi")
	sep := types.String("")
	out := types.DefaultTypeAdapter.NativeToValue([]string{"h", "i"})
	for _, binding := range bindings {
		if binding.Unary != nil {
			if binding.Unary(in).Equal(out) != types.True {
				t.Errorf("binding invocation got %v, wanted %v", binding.Unary(in), out)
			}
			celErr := binding.Unary(types.Bytes("hi"))
			if !types.IsError(celErr) || !strings.Contains(celErr.(*types.Err).String(), "no such overload") {
				t.Errorf("binding.Unary(bytes) got %v, wanted no such overload", celErr)
			}
		}
		if binding.Binary != nil {
			if binding.Binary(in, sep).Equal(out) != types.True {
				t.Errorf("binding invocation got %v, wanted %v", binding.Binary(in, sep), out)
			}
			celErr := binding.Binary(types.Bytes("hi"), sep)
			if !types.IsError(celErr) || !strings.Contains(celErr.(*types.Err).String(), "no such overload") {
				t.Errorf("binding.Binary(bytes, string) got %v, wanted no such overload", celErr)
			}
			celUnk := binding.Binary(types.Bytes("hi"), types.Unknown{1})
			if !types.IsUnknown(celUnk) {
				t.Errorf("binding.Binary(bytes, unk) got %v, wanted unknown{1}", celUnk)
			}
		}
		if binding.Function != nil {
			if binding.Function(in, sep, types.IntNegOne).Equal(out) != types.True {
				t.Errorf("binding invocation got %v, wanted %v", binding.Function(in, sep, types.IntNegOne), out)
			}
			celErr := binding.Function(types.Bytes("hi"), sep, types.IntOne)
			if !types.IsError(celErr) || !strings.Contains(celErr.(*types.Err).String(), "no such overload") {
				t.Errorf("binding.Function(bytes, string, int) got %v, wanted no such overload", celErr)
			}
			if binding.Operator == "split" {
				if binding.Function(in).Equal(out) != types.True {
					t.Errorf("binding invocation got %v, wanted %v", binding.Function(in), out)
				}
				if binding.Function(in, sep).Equal(out) != types.True {
					t.Errorf("binding invocation got %v, wanted %v", binding.Function(in, sep), out)
				}
				out := binding.Function()
				if !types.IsError(out) || out.(*types.Err).String() != "no such overload: split()" {
					t.Fatalf("binding.Function() got %v, wanted error", out)
				}
			}
		}
	}
}

func TestFunctionZeroArityBinding(t *testing.T) {
	now := types.DefaultTypeAdapter.NativeToValue(time.UnixMilli(1000))
	nowFunc, err := NewFunction("now",
		Overload("now", []*Type{}, TimestampType,
			FunctionBinding(func(args ...ref.Val) ref.Val {
				return now
			})),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	bindings, err := nowFunc.Bindings()
	if err != nil {
		t.Fatalf("nowFunc.Bindings() produced an err: %v", err)
	}
	if len(bindings) != 1 {
		t.Errorf("nowFunc.Bindings() produced %d bindings, wanted one", len(bindings))
	}
	out := bindings[0].Function()
	if out != now {
		t.Errorf("now() got %v, wanted %v", out, now)
	}
}

func TestFunctionSingletonBinding(t *testing.T) {
	size, err := NewFunction("size",
		// Since the singleton requires that the operand is a traits.Sizer, the type guards can be
		// disabled for a minor boost in performance, and because the argument type-checking
		// doesn't actually give much additional benefit. The drawback is that invalid signatures
		// at type-check might be valid at runtime.
		DisableTypeGuards(true),
		Overload("size_map", []*Type{MapType(TypeParamType("K"), TypeParamType("V"))}, IntType),
		Overload("size_list", []*Type{ListType(TypeParamType("V"))}, IntType),
		Overload("size_string", []*Type{StringType}, IntType),
		MemberOverload("map_size", []*Type{MapType(TypeParamType("K"), TypeParamType("V"))}, IntType),
		MemberOverload("list_size", []*Type{ListType(TypeParamType("V"))}, IntType),
		MemberOverload("string_size", []*Type{StringType}, IntType),
		SingletonUnaryBinding(func(arg ref.Val) ref.Val {
			return arg.(traits.Sizer).Size()
		}, traits.SizerType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	bindings, err := size.Bindings()
	if err != nil {
		t.Fatalf("size.Bindings() failed: %v", err)
	}
	if len(bindings) != 1 {
		t.Errorf("size.Bindings() got %d bindings, wanted 1", len(bindings))
	}
	if bindings[0].Unary == nil {
		t.Fatalf("size.Bindings() missing singleton unary binding")
	}
	result := bindings[0].Unary(types.String("hello"))
	if result.Equal(types.Int(5)) != types.True {
		t.Errorf("size('hello') got %v, wanted 5", result)
	}
	// Invalid at type-check, but valid since type guard checks have been disabled
	result = bindings[0].Unary(types.Bytes("hello"))
	if result.Equal(types.Int(5)) != types.True {
		t.Errorf("size(b'hello') got %v, wanted 5", result)
	}
}

func TestFunctionMerge(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType),
		MemberOverload("map_size", []*Type{MapType(TypeParamType("K"), TypeParamType("V"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	out, err := sizeFunc.Merge(sizeFunc)
	if err != nil {
		t.Errorf("sizeFunc.Merge(sizeFunc) failed: %v", err)
	}
	if out != sizeFunc {
		t.Errorf("sizeFunc.Merge(sizeFunc) != sizeFunc: %v", out)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("vector_size", []*Type{OpaqueType("vector", TypeParamType("T"))}, IntType),
		SingletonUnaryBinding(func(sizer ref.Val) ref.Val {
			return sizer.(traits.Sizer).Size()
		}, traits.SizerType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeMerged, err := sizeFunc.Merge(sizeVecFunc)
	if err != nil {
		t.Fatalf("Merge() failed: %v", err)
	}
	if (sizeMerged.Name) != "size" {
		t.Errorf("Merge() produced a function with name %v, wanted 'size'", sizeMerged.Name)
	}
	if len(sizeMerged.Overloads) != 3 {
		t.Errorf("Merge() produced %d overloads, wanted 3", len(sizeFunc.Overloads))
	}
	overloads := map[string]bool{
		"list_size":   true,
		"map_size":    true,
		"vector_size": true,
	}
	for _, o := range sizeMerged.Overloads {
		delete(overloads, o.ID)
	}
	if len(overloads) != 0 {
		t.Errorf("Merge() did not include overloads: %v", overloads)
	}
}

func TestFunctionMergeWrongName(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("sizeN",
		MemberOverload("vector_size", []*Type{OpaqueType("vector", TypeParamType("T"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = sizeFunc.Merge(sizeVecFunc)
	if err == nil || !strings.Contains(err.Error(), "unrelated functions") {
		t.Fatalf("Merge() expected to fail, got: %v", err)
	}
}

func TestFunctionMergeOverloadCollision(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("list_size2", []*Type{ListType(TypeParamType("K"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = sizeFunc.Merge(sizeVecFunc)
	if err == nil || !strings.Contains(err.Error(), "declaration merge failed") {
		t.Fatalf("Merge() expected to fail, got: %v", err)
	}
}

func TestFunctionMergeOverloadArgCountRedefinition(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T")), IntType}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = sizeFunc.Merge(sizeVecFunc)
	if err == nil || !strings.Contains(err.Error(), "redefinition") {
		t.Fatalf("Merge() expected to fail, got: %v", err)
	}
}

func TestFunctionMergeOverloadArgTypeRedefinition(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("arg_size", []*Type{ListType(TypeParamType("T"))}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("arg_size", []*Type{MapType(IntType, StringType)}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = sizeFunc.Merge(sizeVecFunc)
	if err == nil || !strings.Contains(err.Error(), "redefinition") {
		t.Fatalf("Merge() expected to fail, got: %v", err)
	}
}

func TestFunctionMergeSingletonRedefinition(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size", []*Type{ListType(TypeParamType("T"))}, IntType),
		SingletonUnaryBinding(func(ref.Val) ref.Val {
			return types.IntZero
		}),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("string_size", []*Type{StringType}, IntType),
		SingletonUnaryBinding(func(ref.Val) ref.Val {
			return types.IntZero
		}),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = sizeFunc.Merge(sizeVecFunc)
	if err == nil || !strings.Contains(err.Error(), "already has singleton") {
		t.Fatalf("Merge() expected to fail, got: %v", err)
	}
}

func TestFunctionAddDuplicateOverloads(t *testing.T) {
	_, err := NewFunction("max",
		Overload("max_int", []*Type{IntType}, IntType),
		Overload("max_int", []*Type{IntType}, IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() with duplicate overload signature failed: %v", err)
	}
}

func TestFunctionAddCollidingOverloads(t *testing.T) {
	_, err := NewFunction("max",
		Overload("max_int", []*Type{IntType}, IntType),
		Overload("max_int2", []*Type{IntType}, IntType),
	)
	if err == nil || !strings.Contains(err.Error(), "max_int collides with max_int2") {
		t.Fatalf("NewFunction() got %v, wanted collision error", err)
	}
}

func TestFunctionNoOverloads(t *testing.T) {
	_, err := NewFunction("right", SingletonBinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
		return arg2
	}))
	if err == nil || !strings.Contains(err.Error(), "must have at least one overload") {
		t.Errorf("NewFunction() got %v, wanted 'must have at least one overload'", err)
	}
}

func TestSingletonOverloadCollision(t *testing.T) {
	fn, err := NewFunction("id",
		Overload("id_any", []*Type{AnyType}, AnyType,
			UnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			}),
		),
		SingletonUnaryBinding(func(arg ref.Val) ref.Val {
			return arg
		}),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = fn.Bindings()
	if err == nil || !strings.Contains(err.Error(), "incompatible with specialized overloads") {
		t.Errorf("NewFunction() got %v, wanted incompatible with specialized overloads", err)
	}
}

func TestSingletonUnaryBindingRedefinition(t *testing.T) {
	_, err := NewFunction("id",
		Overload("id_any", []*Type{AnyType}, AnyType),
		SingletonUnaryBinding(func(arg ref.Val) ref.Val {
			return arg
		}),
		SingletonUnaryBinding(func(arg ref.Val) ref.Val {
			return arg
		}),
	)
	if err == nil || !strings.Contains(err.Error(), "already has a singleton binding") {
		t.Errorf("NewFunction() got %v, wanted already has a singleton binding", err)
	}
}

func TestSingletonBinaryBindingRedefinition(t *testing.T) {
	_, err := NewFunction("right",
		Overload("right_double_double", []*Type{DoubleType, DoubleType}, DoubleType),
		SingletonBinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
			return arg2
		}, traits.ComparerType),
		SingletonBinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
			return arg2
		}),
	)
	if err == nil || !strings.Contains(err.Error(), "already has a singleton binding") {
		t.Errorf("NewCustomEnv() got %v, wanted already has a singleton binding", err)
	}
}

func TestSingletonFunctionBindingRedefinition(t *testing.T) {
	_, err := NewFunction("id",
		Overload("id_any", []*Type{AnyType}, AnyType),
		SingletonFunctionBinding(func(args ...ref.Val) ref.Val {
			return args[0]
		}, traits.ComparerType),
		SingletonFunctionBinding(func(args ...ref.Val) ref.Val {
			return args[0]
		}, traits.ComparerType),
	)
	if err == nil || !strings.Contains(err.Error(), "already has a singleton binding") {
		t.Errorf("NewCustomEnv() got %v, wanted already has a singleton binding", err)
	}
}

func TestOverloadUnaryBindingRedefinition(t *testing.T) {
	_, err := NewFunction("id",
		Overload("id_any", []*Type{AnyType}, AnyType,
			UnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			}),
			UnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			}),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "already has a binding") {
		t.Errorf("NewFunction() got %v, wanted already has a binding", err)
	}
}

func TestOverloadUnaryBindingArgCountMismatch(t *testing.T) {
	_, err := NewFunction("id",
		Overload("id_any", []*Type{}, AnyType,
			UnaryBinding(func(arg ref.Val) ref.Val {
				return arg
			}),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "non-unary overload") {
		t.Errorf("NewFunction() got %v, wanted non-unary overload", err)
	}
}

func TestOverloadBinaryBindingArgCountMismatch(t *testing.T) {
	_, err := NewFunction("id",
		Overload("id_any", []*Type{}, AnyType,
			BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
				return lhs
			}),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "non-binary overload") {
		t.Errorf("NewFunction() got %v, wanted non-binary overload", err)
	}
}

func TestOverloadBinaryBindingRedefinition(t *testing.T) {
	_, err := NewFunction("right",
		Overload("right_double_double", []*Type{DoubleType, DoubleType}, DoubleType,
			BinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
				return arg2
			}),
			BinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
				return arg2
			}),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "already has a binding") {
		t.Errorf("NewCustomEnv() got %v, wanted already has a binding", err)
	}
}

func TestOverloadFunctionBindingRedefinition(t *testing.T) {
	_, err := NewFunction("id",
		Overload("id_any", []*Type{AnyType}, AnyType,
			FunctionBinding(func(args ...ref.Val) ref.Val {
				return args[0]
			}),
			FunctionBinding(func(args ...ref.Val) ref.Val {
				return args[0]
			}),
		),
	)
	if err == nil || !strings.Contains(err.Error(), "already has a binding") {
		t.Errorf("NewCustomEnv() got %v, wanted already has a binding", err)
	}
}

func TestOverloadIsNonStrict(t *testing.T) {
	fn, err := NewFunction("getOrDefault",
		MemberOverload("get",
			[]*Type{MapType(
				TypeParamType("K"), TypeParamType("V")),
				TypeParamType("K"),
				TypeParamType("V"),
			},
			TypeParamType("V"),
			OverloadOperandTrait(traits.ContainerType|traits.IndexerType),
			OverloadIsNonStrict(),
			FunctionBinding(func(args ...ref.Val) ref.Val {
				container := args[0].(traits.Container)
				key := args[1]
				orValue := args[2]
				if types.Bool(types.IsUnknownOrError(key)) {
					return orValue
				}
				if container.Contains(key) == types.True {
					return args[0].(traits.Indexer).Get(key)
				}
				return orValue
			}),
		),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	bindings, err := fn.Bindings()
	if err != nil {
		t.Fatalf("fn.Binding() failed: %v", err)
	}
	m := types.DefaultTypeAdapter.NativeToValue(map[string]string{"hello": "world"})
	out := bindings[0].Function(m, types.String("hello"), types.String("goodbye"))
	if out.Equal(types.String("world")) != types.True {
		t.Errorf("function got %v, wanted 'world'", out)
	}
	out = bindings[0].Function(m, types.String("missing"), types.String("goodbye"))
	if out.Equal(types.String("goodbye")) != types.True {
		t.Errorf("function got %v, wanted 'goodbye'", out)
	}
	out = bindings[0].Function(m, types.NewErr("no such key"), types.String("goodbye"))
	if out.Equal(types.String("goodbye")) != types.True {
		t.Errorf("function got %v, wanted 'goodbye'", out)
	}
}

func TestOverloadOperandTrait(t *testing.T) {
	fn, err := NewFunction("getOrDefault",
		MemberOverload("get",
			[]*Type{MapType(
				TypeParamType("K"), TypeParamType("V")),
				TypeParamType("K"),
				TypeParamType("V"),
			},
			TypeParamType("V"),
			OverloadOperandTrait(traits.ContainerType|traits.IndexerType),
			FunctionBinding(func(args ...ref.Val) ref.Val {
				container := args[0].(traits.Container)
				key := args[1]
				orValue := args[2]
				if container.Contains(key) == types.True {
					return args[0].(traits.Indexer).Get(key)
				}
				return orValue
			}),
		),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	bindings, err := fn.Bindings()
	if err != nil {
		t.Fatalf("fn.Binding() failed: %v", err)
	}
	m := types.DefaultTypeAdapter.NativeToValue(map[string]string{"hello": "world"})
	out := bindings[0].Function(m, types.String("hello"), types.String("goodbye"))
	if out.Equal(types.String("world")) != types.True {
		t.Errorf("function got %v, wanted 'world'", out)
	}
	out = bindings[0].Function(m, types.String("missing"), types.String("goodbye"))
	if out.Equal(types.String("goodbye")) != types.True {
		t.Errorf("function got %v, wanted 'goodbye'", out)
	}
	noSuchKey := types.NewErr("no such key")
	out = bindings[0].Function(m, noSuchKey, types.String("goodbye"))
	if out != noSuchKey {
		t.Errorf("function got %v, wanted 'no such key'", out)
	}
}

func TestTypeString(t *testing.T) {
	tests := []struct {
		in  *Type
		out string
	}{
		{
			in:  ListType(IntType),
			out: "list(int)",
		},
		{
			in:  MapType(UintType, DoubleType),
			out: "map(uint, double)",
		},
		{
			in:  NullableType(BoolType),
			out: "bool",
		},
		{
			in:  OptionalType(ListType(StringType)),
			out: "optional(list(string))",
		},
		{
			in:  ObjectType("my.type.Message"),
			out: "my.type.Message",
		},
		{
			in:  ObjectType("google.protobuf.Int32Value"),
			out: "int",
		},
		{
			in:  ObjectType("google.protobuf.UInt32Value"),
			out: "uint",
		},
		{
			in:  ObjectType("google.protobuf.Value"),
			out: "dyn",
		},
		{
			in:  TypeTypeWithParam(StringType),
			out: "type(string)",
		},
		{
			in:  TypeParamType("T"),
			out: "T",
		},
	}
	for _, tst := range tests {
		if tst.in.String() != tst.out {
			t.Errorf("String() got %v, wanted %v", tst.in, tst.out)
		}
	}
}

func TestTypeIsType(t *testing.T) {
	tests := []struct {
		t1     *Type
		t2     *Type
		isType bool
	}{
		{
			t1:     StringType,
			t2:     StringType,
			isType: true,
		},
		{
			t1:     StringType,
			t2:     IntType,
			isType: false,
		},
		{
			t1:     OptionalType(StringType),
			t2:     OptionalType(IntType),
			isType: false,
		},
		{
			t1:     OptionalType(UintType),
			t2:     OptionalType(UintType),
			isType: true,
		},
		{
			t1:     MapType(BoolType, IntType),
			t2:     MapType(BoolType, IntType),
			isType: true,
		},
		{
			t1:     MapType(TypeParamType("K1"), IntType),
			t2:     MapType(TypeParamType("K2"), IntType),
			isType: true,
		},
		{
			t1:     MapType(TypeParamType("K1"), ObjectType("my.msg.First")),
			t2:     MapType(TypeParamType("K2"), ObjectType("my.msg.Last")),
			isType: false,
		},
	}
	for _, tst := range tests {
		if tst.t1.IsType(tst.t2) != tst.isType {
			t.Errorf("%v.IsType(%v) got %v, wanted %v", tst.t1, tst.t2, !tst.isType, tst.isType)
		}
	}
}

func TestTypeIsAssignableType(t *testing.T) {
	tests := []struct {
		t1           *Type
		t2           *Type
		isAssignable bool
	}{
		{
			t1:           NullableType(DoubleType),
			t2:           NullType,
			isAssignable: true,
		},
		{
			t1:           NullableType(DoubleType),
			t2:           DoubleType,
			isAssignable: true,
		},
		{
			t1:           OpaqueType("vector", NullableType(DoubleType)),
			t2:           OpaqueType("vector", NullType),
			isAssignable: true,
		},
		{
			t1:           OpaqueType("vector", NullableType(DoubleType)),
			t2:           OpaqueType("vector", DoubleType),
			isAssignable: true,
		},
		{
			t1:           OpaqueType("vector", DynType),
			t2:           OpaqueType("vector", NullableType(IntType)),
			isAssignable: true,
		},
		{
			t1:           ObjectType("my.msg.MsgName"),
			t2:           ObjectType("my.msg.MsgName"),
			isAssignable: true,
		},
		{
			t1:           MapType(TypeParamType("K"), IntType),
			t2:           MapType(StringType, IntType),
			isAssignable: true,
		},
		{
			t1:           MapType(StringType, IntType),
			t2:           MapType(TypeParamType("K"), IntType),
			isAssignable: false,
		},
		{
			t1:           OpaqueType("vector", DoubleType),
			t2:           OpaqueType("vector", NullableType(IntType)),
			isAssignable: false,
		},
		{
			t1:           OpaqueType("vector", NullableType(DoubleType)),
			t2:           OpaqueType("vector", DynType),
			isAssignable: false,
		},
		{
			t1:           ObjectType("my.msg.MsgName"),
			t2:           ObjectType("my.msg.MsgName2"),
			isAssignable: false,
		},
	}
	for _, tst := range tests {
		if tst.t1.IsAssignableType(tst.t2) != tst.isAssignable {
			t.Errorf("%v.IsAssignableType(%v) got %v, wanted %v", tst.t1, tst.t2, !tst.isAssignable, tst.isAssignable)
		}
	}
}

func TestTypeSignatureEquals(t *testing.T) {
	paramA := TypeParamType("A")
	paramB := TypeParamType("B")
	mapOfAB := MapType(paramA, paramB)
	fn, err := NewFunction(overloads.Size,
		MemberOverload(overloads.SizeMapInst, []*Type{mapOfAB}, IntType),
		Overload(overloads.SizeMap, []*Type{mapOfAB}, IntType))
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	if !fn.Overloads[overloads.SizeMap].SignatureEquals(fn.Overloads[overloads.SizeMap]) {
		t.Errorf("SignatureEquals() returned false, wanted true")
	}
	if fn.Overloads[overloads.SizeMap].SignatureEquals(fn.Overloads[overloads.SizeMapInst]) {
		t.Errorf("SignatureEquals() returned false, wanted true")
	}
}

func TestTypeIsAssignableRuntimeType(t *testing.T) {
	if !NullableType(DoubleType).IsAssignableRuntimeType(types.NullValue) {
		t.Error("nullable double cannot be assigned from null")
	}
	if !NullableType(DoubleType).IsAssignableRuntimeType(types.Double(0.0)) {
		t.Error("nullable double cannot be assigned from double")
	}
	if !MapType(StringType, DurationType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[string]time.Duration{})) {
		t.Error("map(string, duration) not assignable to map at runtime")
	}
	if !MapType(StringType, DurationType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[string]time.Duration{"one": time.Duration(1)})) {
		t.Error("map(string, duration) not assignable to map at runtime")
	}
	if !MapType(StringType, DynType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[string]time.Duration{"one": time.Duration(1)})) {
		t.Error("map(string, dyn) not assignable to map at runtime")
	}
	if MapType(StringType, DynType).IsAssignableRuntimeType(
		types.DefaultTypeAdapter.NativeToValue(map[int64]time.Duration{1: time.Duration(1)})) {
		t.Error("map(string, dyn) must not be assignable to map(int, duration) at runtime")
	}
}
