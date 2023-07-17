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
	"reflect"
	"strings"
	"testing"
	"time"

	"google.golang.org/protobuf/proto"

	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestFunctionBindings(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
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
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType,
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
		MemberOverload("string_split",
			[]*types.Type{types.StringType}, types.NewListType(types.StringType),
			UnaryBinding(func(str ref.Val) ref.Val {
				s := str.(types.String)
				return splitImpl(string(s), "", -1)
			})),
		MemberOverload("string_split_string",
			[]*types.Type{types.StringType, types.StringType}, types.NewListType(types.StringType),
			BinaryBinding(func(str, sep ref.Val) ref.Val {
				s := str.(types.String)
				delim := sep.(types.String)
				return splitImpl(string(s), string(delim), -1)
			})),
		MemberOverload("string_split_string_int",
			[]*types.Type{types.StringType, types.StringType, types.IntType}, types.NewListType(types.StringType),
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
			celUnk := binding.Binary(types.Bytes("hi"), types.NewUnknown(1, types.NewAttributeTrail("x")))
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
		Overload("now", []*types.Type{}, types.TimestampType,
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
		Overload("size_map",
			[]*types.Type{types.NewMapType(types.NewTypeParamType("K"), types.NewTypeParamType("V"))}, types.IntType),
		Overload("size_list",
			[]*types.Type{types.NewListType(types.NewTypeParamType("V"))}, types.IntType),
		Overload("size_string",
			[]*types.Type{types.StringType}, types.IntType),
		MemberOverload("map_size",
			[]*types.Type{types.NewMapType(types.NewTypeParamType("K"), types.NewTypeParamType("V"))}, types.IntType),
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("V"))}, types.IntType),
		MemberOverload("string_size",
			[]*types.Type{types.StringType}, types.IntType),
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
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
		MemberOverload("map_size",
			[]*types.Type{types.NewMapType(types.NewTypeParamType("K"), types.NewTypeParamType("V"))}, types.IntType),
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
		MemberOverload("vector_size",
			[]*types.Type{types.NewOpaqueType("vector", types.NewTypeParamType("T"))}, types.IntType),
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
	if sizeMerged.Name() != "size" {
		t.Errorf("Merge() produced a function with name %v, wanted 'size'", sizeMerged.Name())
	}
	if len(sizeMerged.overloads) != 3 {
		t.Errorf("Merge() produced %d overloads, wanted 3", len(sizeFunc.overloads))
	}
	overloads := map[string]bool{
		"list_size":   true,
		"map_size":    true,
		"vector_size": true,
	}
	for _, o := range sizeMerged.overloads {
		delete(overloads, o.ID())
	}
	if len(overloads) != 0 {
		t.Errorf("Merge() did not include overloads: %v", overloads)
	}
}

func TestFunctionMergeWrongName(t *testing.T) {
	sizeFunc, err := NewFunction("size",
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("sizeN",
		MemberOverload("vector_size",
			[]*types.Type{types.NewOpaqueType("vector", types.NewTypeParamType("T"))}, types.IntType),
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
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("list_size2",
			[]*types.Type{types.NewListType(types.NewTypeParamType("K"))}, types.IntType),
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
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T")), types.IntType}, types.IntType),
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
		MemberOverload("arg_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("arg_size",
			[]*types.Type{types.NewMapType(types.IntType, types.StringType)}, types.IntType),
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
		MemberOverload("list_size",
			[]*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
		SingletonUnaryBinding(func(ref.Val) ref.Val {
			return types.IntZero
		}),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	sizeVecFunc, err := NewFunction("size",
		MemberOverload("string_size",
			[]*types.Type{types.StringType}, types.IntType),
		SingletonUnaryBinding(func(ref.Val) ref.Val {
			return types.IntZero
		}),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	_, err = sizeFunc.Merge(sizeVecFunc)
	if err == nil || !strings.Contains(err.Error(), "already has a singleton") {
		t.Fatalf("Merge() expected to fail, got: %v", err)
	}
}

func TestFunctionAddDuplicateOverloads(t *testing.T) {
	_, err := NewFunction("max",
		Overload("max_int", []*types.Type{types.IntType}, types.IntType),
		Overload("max_int", []*types.Type{types.IntType}, types.IntType),
	)
	if err != nil {
		t.Fatalf("NewFunction() with duplicate overload signature failed: %v", err)
	}
}

func TestFunctionAddCollidingOverloads(t *testing.T) {
	_, err := NewFunction("max",
		Overload("max_int", []*types.Type{types.IntType}, types.IntType),
		Overload("max_int2", []*types.Type{types.IntType}, types.IntType),
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
		Overload("id_any", []*types.Type{types.AnyType}, types.AnyType,
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
		Overload("id_any", []*types.Type{types.AnyType}, types.AnyType),
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
		Overload("right_double_double", []*types.Type{types.DoubleType, types.DoubleType}, types.DoubleType),
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
		Overload("id_any", []*types.Type{types.AnyType}, types.AnyType),
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
		Overload("id_any", []*types.Type{types.AnyType}, types.AnyType,
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
		Overload("id_any", []*types.Type{}, types.AnyType,
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
		Overload("id_any", []*types.Type{}, types.AnyType,
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
		Overload("right_double_double", []*types.Type{types.DoubleType, types.DoubleType}, types.DoubleType,
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
		Overload("id_any", []*types.Type{types.AnyType}, types.AnyType,
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
			[]*types.Type{types.NewMapType(
				types.NewTypeParamType("K"), types.NewTypeParamType("V")),
				types.NewTypeParamType("K"),
				types.NewTypeParamType("V"),
			},
			types.NewTypeParamType("V"),
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
			[]*types.Type{types.NewMapType(
				types.NewTypeParamType("K"), types.NewTypeParamType("V")),
				types.NewTypeParamType("K"),
				types.NewTypeParamType("V"),
			},
			types.NewTypeParamType("V"),
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

func TestFunctionGetTypeParams(t *testing.T) {
	fn, err := NewFunction("deep_type_params",
		Overload("no_type_params", []*types.Type{}, types.DynType),
		Overload("one_type_param", []*types.Type{types.BoolType}, types.NewTypeParamType("K")),
		Overload("deep_type_params",
			[]*types.Type{types.NewTypeParamType("E1"),
				types.NewMapType(types.NewTypeParamType("K"), types.NewTypeParamType("V"))},
			types.NewTypeParamType("V"),
		),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	if len(fn.OverloadDecls()) != 3 {
		t.Fatal("fn.OverloadDecls() not equal to 3")
	}
	o1 := fn.OverloadDecls()[0]
	o2 := fn.OverloadDecls()[1]
	o3 := fn.OverloadDecls()[2]
	if len(o1.TypeParams()) != 0 {
		t.Errorf("overload %v did not have zero type-params", o1)
	}
	if len(o2.TypeParams()) != 1 && !reflect.DeepEqual(o2.TypeParams(), []string{"K"}) {
		t.Errorf("overload %v did not have a single type param", o2)
	}
	if len(o3.TypeParams()) != 3 {
		t.Errorf("overload %v did not have three type params", o3)
	}
}

func TestFunctionDisableDeclaration(t *testing.T) {
	fn, err := NewFunction("in",
		DisableDeclaration(true),
		Overload("in_list",
			[]*types.Type{types.NewListType(types.NewTypeParamType("K")), types.NewTypeParamType("K")},
			types.BoolType,
		),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	if !fn.IsDeclarationDisabled() {
		t.Error("got declaration enabled, wanted disabled")
	}
}

func TestFunctionEnableDeclaration(t *testing.T) {
	fn, err := NewFunction("in",
		DisableDeclaration(false),
		Overload("in_list",
			[]*types.Type{types.NewListType(types.NewTypeParamType("K")), types.NewTypeParamType("K")},
			types.BoolType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	if fn.IsDeclarationDisabled() {
		t.Error("got declaration disabled, wanted enabled")
	}
	fn2, err := NewFunction("in",
		DisableDeclaration(true),
		Overload("in_list",
			[]*types.Type{types.NewListType(types.NewTypeParamType("K")), types.NewTypeParamType("K")},
			types.BoolType),
	)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	if !fn2.IsDeclarationDisabled() {
		t.Error("got declaration enabled, wanted disabled")
	}

	// enabled -> disabled
	merged, err := fn.Merge(fn2)
	if err != nil {
		t.Fatalf("fn.Merge(fn2) failed: %v", err)
	}
	if !merged.IsDeclarationDisabled() {
		t.Error("got declaration enabled, wanted disabled")
	}
	// disabled -> enabled
	merged2, err := fn2.Merge(fn)
	if err != nil {
		t.Fatalf("fn.Merge(fn2) failed: %v", err)
	}
	if merged2.IsDeclarationDisabled() {
		t.Error("got declaration disabled, wanted enabled")
	}
}

func TestFunctionDeclToExprDecl(t *testing.T) {
	tests := []struct {
		fn     *FunctionDecl
		exDecl *exprpb.Decl
	}{
		{
			fn: testFunction(t, "equals",
				Overload("equals_value_value",
					[]*types.Type{types.NewTypeParamType("T"), types.NewTypeParamType("T")}, types.BoolType)),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId: "equals_value_value",
								Params: []*exprpb.Type{
									chkdecls.NewTypeParamType("T"),
									chkdecls.NewTypeParamType("T"),
								},
								TypeParams: []string{"T"},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
		{
			fn: testFunction(t, "equals",
				MemberOverload("value_equals_value",
					[]*types.Type{types.NewTypeParamType("T"), types.NewTypeParamType("T")}, types.BoolType)),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId:         "value_equals_value",
								IsInstanceFunction: true,
								Params: []*exprpb.Type{
									chkdecls.NewTypeParamType("T"),
									chkdecls.NewTypeParamType("T"),
								},
								TypeParams: []string{"T"},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
		{
			fn: testFunction(t, "equals",
				Overload("equals_int_uint",
					[]*types.Type{types.IntType, types.UintType}, types.BoolType)),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId: "equals_int_uint",
								Params: []*exprpb.Type{
									chkdecls.Int,
									chkdecls.Uint,
								},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
		{
			fn: testFunction(t, "equals",
				MemberOverload("int_equals_uint",
					[]*types.Type{types.IntType, types.UintType}, types.BoolType)),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId:         "int_equals_uint",
								IsInstanceFunction: true,
								Params: []*exprpb.Type{
									chkdecls.Int,
									chkdecls.Uint,
								},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
		{
			fn: testFunction(t, "equals",
				MemberOverload("list_optional_value_equals_list_optional_value", []*types.Type{
					types.NewListType(types.NewOptionalType(types.NewTypeParamType("T"))),
					types.NewListType(types.NewOptionalType(types.NewTypeParamType("T"))),
				}, types.BoolType)),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId:         "list_optional_value_equals_list_optional_value",
								IsInstanceFunction: true,
								Params: []*exprpb.Type{
									chkdecls.NewListType(chkdecls.NewOptionalType(chkdecls.NewTypeParamType("T"))),
									chkdecls.NewListType(chkdecls.NewOptionalType(chkdecls.NewTypeParamType("T"))),
								},
								TypeParams: []string{"T"},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
		{
			fn: testFunction(t, "equals",
				MemberOverload("int_equals_uint",
					[]*types.Type{types.IntType, types.UintType}, types.BoolType),
				MemberOverload("uint_equals_int",
					[]*types.Type{types.UintType, types.IntType}, types.BoolType)),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId:         "int_equals_uint",
								IsInstanceFunction: true,
								Params: []*exprpb.Type{
									chkdecls.Int,
									chkdecls.Uint,
								},
								ResultType: chkdecls.Bool,
							},
							{
								OverloadId:         "uint_equals_int",
								IsInstanceFunction: true,
								Params: []*exprpb.Type{
									chkdecls.Uint,
									chkdecls.Int,
								},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
		// Test to make sure the overloads are in stable order: int_uint | uint_int | int_int | uint_uint
		{
			fn: testMerge(t,
				testFunction(t, "equals",
					Overload("int_equals_uint", []*types.Type{types.IntType, types.UintType}, types.BoolType),
					Overload("uint_equals_int", []*types.Type{types.UintType, types.IntType}, types.BoolType)),
				testFunction(t, "equals",
					Overload("int_equals_int", []*types.Type{types.IntType, types.IntType}, types.BoolType),
					Overload("int_equals_uint", []*types.Type{types.IntType, types.UintType}, types.BoolType),
					Overload("uint_equals_uint", []*types.Type{types.UintType, types.UintType}, types.BoolType))),
			exDecl: &exprpb.Decl{
				Name: "equals",
				DeclKind: &exprpb.Decl_Function{
					Function: &exprpb.Decl_FunctionDecl{
						Overloads: []*exprpb.Decl_FunctionDecl_Overload{
							{
								OverloadId: "int_equals_uint",
								Params: []*exprpb.Type{
									chkdecls.Int,
									chkdecls.Uint,
								},
								ResultType: chkdecls.Bool,
							},
							{
								OverloadId: "uint_equals_int",
								Params: []*exprpb.Type{
									chkdecls.Uint,
									chkdecls.Int,
								},
								ResultType: chkdecls.Bool,
							},
							{
								OverloadId: "int_equals_int",
								Params: []*exprpb.Type{
									chkdecls.Int,
									chkdecls.Int,
								},
								ResultType: chkdecls.Bool,
							},
							{
								OverloadId: "uint_equals_uint",
								Params: []*exprpb.Type{
									chkdecls.Uint,
									chkdecls.Uint,
								},
								ResultType: chkdecls.Bool,
							},
						},
					},
				},
			},
		},
	}
	for _, tst := range tests {
		exDecl, err := FunctionDeclToExprDecl(tst.fn)
		if err != nil {
			t.Fatalf("FunctionDeclToExprDecl(%v) failed: %v", tst.fn, err)
		}
		if !proto.Equal(exDecl, tst.exDecl) {
			t.Errorf("got not equal, wanted %v == %v", exDecl, tst.exDecl)
		}
	}
}

func TestFunctionDeclToExprDeclInvalid(t *testing.T) {
	fn1 := testFunction(t, "bad_equals",
		MemberOverload("bad_equals_param", []*types.Type{{}, types.UintType}, types.BoolType))
	ex1, err := FunctionDeclToExprDecl(fn1)
	if err == nil {
		t.Errorf("FunctionDeclToExprDecl(bad_equals) succeeded: %v, wanted error", ex1)
	}
	fn2 := testFunction(t, "bad_equals",
		Overload("bad_equals_out", []*types.Type{types.IntType, types.UintType}, &types.Type{}))
	ex2, err := FunctionDeclToExprDecl(fn2)
	if err == nil {
		t.Errorf("FunctionDeclToExprDecl(bad_equals) succeeded: %v, wanted error", ex2)
	}
}

func TestNewVariable(t *testing.T) {
	a := NewVariable("a", types.BoolType)
	if !a.DeclarationIsEquivalent(a) {
		t.Error("NewVariable(a, bool) does not equal itself")
	}
	if !a.DeclarationIsEquivalent(NewVariable("a", types.BoolType)) {
		t.Error("NewVariable(a, bool) does not equal itself")
	}
	a1 := NewVariable("a", types.IntType)
	if a.DeclarationIsEquivalent(a1) {
		t.Error("NewVariable(a, int).DeclarationEquals(NewVariable(a, bool))")
	}
}

func TestNewConstant(t *testing.T) {
	a := NewConstant("a", types.IntType, types.Int(42))
	if !a.DeclarationIsEquivalent(a) {
		t.Error("NewConstant(a, int) does not equal itself")
	}
	if !a.DeclarationIsEquivalent(NewVariable("a", types.IntType)) {
		t.Error("NewConstant(a, int) is not declaration equivalent to int variable")
	}
}

func TestTypeVariable(t *testing.T) {
	tests := []struct {
		t *types.Type
		v *VariableDecl
	}{
		{
			t: types.AnyType,
			v: NewVariable("google.protobuf.Any", types.NewTypeTypeWithParam(types.AnyType)),
		},
		{
			t: types.DynType,
			v: NewVariable("dyn", types.NewTypeTypeWithParam(types.DynType)),
		},
		{
			t: types.NewObjectType("google.protobuf.Int32Value"),
			v: NewVariable("int", types.NewTypeTypeWithParam(types.NewNullableType(types.IntType))),
		},
		{
			t: types.NewObjectType("google.protobuf.Int32Value"),
			v: NewVariable("int", types.NewTypeTypeWithParam(types.NewNullableType(types.IntType))),
		},
	}
	for _, tst := range tests {
		if !TypeVariable(tst.t).DeclarationIsEquivalent(tst.v) {
			t.Errorf("got not equal %v.Equals(%v)", TypeVariable(tst.t), tst.v)
		}
	}
}

func TestVariableDeclToExprDecl(t *testing.T) {
	a, err := VariableDeclToExprDecl(NewVariable("a", types.BoolType))
	if err != nil {
		t.Fatalf("VariableDeclToExprDecl() failed: %v", err)
	}
	if !proto.Equal(a, chkdecls.NewVar("a", chkdecls.Bool)) {
		t.Error("proto.Equal() returned false, wanted true")
	}

}

func TestVariableDeclToExprDeclInvalid(t *testing.T) {
	out, err := VariableDeclToExprDecl(NewVariable("bad", &types.Type{}))
	if err == nil {
		t.Fatalf("VariableDeclToExprDecl() succeeded: %v, wanted error", out)
	}
}

func testMerge(t *testing.T, funcs ...*FunctionDecl) *FunctionDecl {
	t.Helper()
	fn := funcs[0]
	var err error
	for i := 1; i < len(funcs); i++ {
		fn, err = fn.Merge(funcs[i])
		if err != nil {
			t.Fatalf("%v.Merge(%v) failed: %v", fn, funcs[i], err)
		}
	}
	return fn
}

func testFunction(t *testing.T, name string, opts ...FunctionOpt) *FunctionDecl {
	t.Helper()
	fn, err := NewFunction(name, opts...)
	if err != nil {
		t.Fatalf("NewFunction() failed: %v", err)
	}
	return fn
}
