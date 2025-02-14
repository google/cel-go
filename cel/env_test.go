// Copyright 2021 Google LLC
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
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/env"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	"google.golang.org/protobuf/proto"

	proto3pb "github.com/google/cel-go/test/proto3pb"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestAstNil(t *testing.T) {
	var ast *Ast
	if ast.IsChecked() {
		t.Error("ast.IsChecked() returned true for nil ast")
	}
	if ast.Expr() != nil {
		t.Errorf("ast.Expr() got %v, wanted nil", ast.Expr())
	}
	if ast.SourceInfo() != nil {
		t.Errorf("ast.SourceInfo() got %v, wanted nil", ast.SourceInfo())
	}
	if ast.OutputType() != types.ErrorType {
		t.Errorf("ast.OutputType() got %v, wanted error type", ast.OutputType())
	}
	if ast.Source() != nil {
		t.Errorf("ast.Source() got %v, wanted nil", ast.Source())
	}
}

func TestIssuesNil(t *testing.T) {
	var iss *Issues
	iss = iss.Append(iss)
	if iss.Err() != nil {
		t.Errorf("iss.Err() got %v, wanted nil given nil issue set", iss.Err())
	}
	if len(iss.Errors()) != 0 {
		t.Errorf("iss.Errors() got %v, wanted empty value", iss.Errors())
	}
	if iss.String() != "" {
		t.Errorf("iss.String() returned %v, wanted empty value", iss.String())
	}
}

func TestIssuesEmpty(t *testing.T) {
	iss := NewIssues(common.NewErrors(nil))
	if iss.Err() != nil {
		t.Errorf("iss.Err() got %v, wanted nil given nil issue set", iss.Err())
	}
	if len(iss.Errors()) != 0 {
		t.Errorf("iss.Errors() got %v, wanted empty value", iss.Errors())
	}
	if iss.String() != "" {
		t.Errorf("iss.String() returned %v, wanted empty value", iss.String())
	}
	var iss2 *Issues
	iss3 := iss.Append(iss2)
	iss4 := iss3.Append(nil)
	if !reflect.DeepEqual(iss4, iss) {
		t.Error("Append() with a nil value resulted in the creation of a new issue set")
	}
}

func TestIssuesAppendSelf(t *testing.T) {
	e, err := NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	_, iss := e.Compile("a")
	if len(iss.Errors()) != 1 {
		t.Errorf("iss.Errors() got %v, wanted 1 error", iss.Errors())
	}
	iss = iss.Append(iss)
	if len(iss.Errors()) != 1 {
		t.Errorf("iss.Errors() got %v, wanted 1 error", iss.Errors())
	}
}

func TestIssues(t *testing.T) {
	e, err := NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	_, iss := e.Compile("-")
	_, iss2 := e.Compile("b")
	iss = iss.Append(iss2)
	if len(iss.Errors()) != 3 {
		t.Errorf("iss.Errors() got %v, wanted 3 errors", iss.Errors())
	}

	wantIss := `ERROR: <input>:1:1: undeclared reference to 'b' (in container '')
 | -
 | ^
ERROR: <input>:1:2: Syntax error: no viable alternative at input '-'
 | -
 | .^
ERROR: <input>:1:2: Syntax error: mismatched input '<EOF>' expecting {'[', '{', '(', '.', '-', '!', 'true', 'false', 'null', NUM_FLOAT, NUM_INT, NUM_UINT, STRING, BYTES, IDENTIFIER}
 | -
 | .^`
	if iss.String() != wantIss {
		t.Errorf("iss.String() returned %v, wanted %v", iss.String(), wantIss)
	}
}

func TestFormatCELTypeEquivalence(t *testing.T) {
	values := []*Type{
		AnyType,
		MapType(StringType, DynType),
		types.NewTypeTypeWithParam(ListType(IntType)),
		TypeType,
		NullableType(DoubleType),
	}
	for _, v := range values {
		v := v
		t.Run(v.String(), func(t *testing.T) {
			celStr := FormatCELType(v)
			et, err := TypeToExprType(v)
			if err != nil {
				t.Fatalf("TypeToExprType(%v) failed: %v", v, err)
			}
			exprStr := FormatType(et)
			if celStr != exprStr {
				t.Errorf("FormatCELType(%v) got %s, wanted %s", v, celStr, exprStr)
			}
		})
	}
}

func TestEnvCheckExtendRace(t *testing.T) {
	t.Parallel()
	for i := 0; i < 500; i++ {
		wg := &sync.WaitGroup{}
		wg.Add(2)
		env, err := NewCustomEnv(StdLib())
		if err != nil {
			t.Fatalf("NewCustomEnv() failed: %v", err)
		}
		t.Run(fmt.Sprintf("Compile[%d]", i), func(t *testing.T) {
			go func() {
				defer wg.Done()
				_, _ = env.Compile(`1 + 1 * 20 < 400`)
			}()
		})
		t.Run(fmt.Sprintf("Extend[%d]", i), func(t *testing.T) {
			go func() {
				defer wg.Done()
				_, _ = env.Extend(Variable("bar", BoolType))
			}()
		})
		wg.Wait()
	}
}

func TestEnvPartialVarsError(t *testing.T) {
	env := testEnv(t)
	_, err := env.PartialVars(10)
	if err == nil {
		t.Error("env.PartialVars(10) succeeded, wanted error")
	}
}

func TestTypeProviderInterop(t *testing.T) {
	reg, err := types.NewRegistry(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	tests := []struct {
		name     string
		provider any
	}{
		{
			name:     "custom provider",
			provider: &customCELProvider{provider: reg},
		},
		{
			name:     "custom legacy provider",
			provider: &customLegacyProvider{provider: reg},
		},
		{
			name:     "provider",
			provider: reg,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			env, err := NewEnv(CustomTypeProvider(tc.provider))
			if err != nil {
				t.Fatalf("NewEnv(CustomTypeProvider()) failed: %v", err)
			}
			// Found type
			pbType, found := env.TypeProvider().FindType("google.expr.proto3.test.TestAllTypes")
			if !found {
				t.Fatal("FindType(google.expr.proto3.test.TestAllTypes) failed")
			}
			celType, found := env.CELTypeProvider().FindStructType("google.expr.proto3.test.TestAllTypes")
			if !found {
				t.Fatal("FindStructType(google.expr.proto3.test.TestAllTypes) failed")
			}
			pbConvType, err := types.ExprTypeToType(pbType)
			if err != nil {
				t.Fatalf("types.ExprTypeToType(%v) failed: %v", pbType, err)
			}
			if !celType.IsExactType(pbConvType) {
				t.Errorf("got converted type %v, wanted %v", pbConvType, celType)
			}
			// Found field
			pbFieldType, found := env.TypeProvider().FindFieldType("google.expr.proto3.test.TestAllTypes", "single_int32")
			if !found {
				t.Fatal("FindFieldType(google.expr.proto3.test.TestAllTypes, single_int32) not found")
			}
			celFieldType, found := env.CELTypeProvider().FindStructFieldType("google.expr.proto3.test.TestAllTypes", "single_int32")
			if !found {
				t.Fatal("FindStructFieldType(google.expr.proto3.test.TestAllTypes, single_int32) not found")
			}
			pbConvFieldType, err := types.ExprTypeToType(pbFieldType.Type)
			if err != nil {
				t.Fatalf("types.ExprTypeToType(%v) failed: %v", pbFieldType.Type, err)
			}
			if !celFieldType.Type.IsExactType(pbConvFieldType) {
				t.Errorf("got converted field type %v, wanted %v", pbConvFieldType, celFieldType)
			}
			// Not found type
			_, found = env.TypeProvider().FindType("test.BadTypeName")
			if found {
				t.Fatal("FindType(test.BadTypeName) succeeded")
			}
			_, found = env.CELTypeProvider().FindStructType("test.BadTypeName")
			if found {
				t.Fatal("FindStructType(test.BadTypeName) succeeded")
			}
			// Not found field
			_, found = env.TypeProvider().FindFieldType("google.expr.proto3.test.TestAllTypes", "undefined_field")
			if found {
				t.Fatal("FindFieldType(google.expr.proto3.test.TestAllTypes, undefined_field) was found")
			}
			_, found = env.CELTypeProvider().FindStructFieldType("google.expr.proto3.test.TestAllTypes", "undefined_field")
			if found {
				t.Fatal("FindStructFieldType(google.expr.proto3.test.TestAllTypes, undefined_field) not found")
			}
		})
	}
}

func TestLibraries(t *testing.T) {
	e, err := NewEnv(OptionalTypes())
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	for _, expected := range []string{"cel.lib.std", "cel.lib.optional"} {
		if !e.HasLibrary(expected) {
			t.Errorf("Expected HasLibrary() to return true for '%s'", expected)
		}
		libMap := map[string]struct{}{}
		libraries := e.Libraries()
		for _, lib := range libraries {
			libMap[lib] = struct{}{}
		}
		if len(libraries) != 2 {
			t.Errorf("Expected HasLibrary() to contain exactly 2 libraries but got: %v", libraries)
		}

		if _, ok := libMap[expected]; !ok {
			t.Errorf("Expected Libraries() to include '%s'", expected)
		}
	}
}

func TestFunctions(t *testing.T) {
	e, err := NewEnv(OptionalTypes())
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	for _, expected := range []string{"optional.of", "or"} {
		if !e.HasFunction(expected) {
			t.Errorf("Expected HasFunction() to return true for '%s'", expected)
		}

		if _, ok := e.Functions()[expected]; !ok {
			t.Errorf("Expected Functions() to include '%s'", expected)
		}
	}
}

func TestEnvToConfig(t *testing.T) {
	tests := []struct {
		name string
		opts []EnvOption
		want *env.Config
	}{
		{
			name: "std env",
			want: env.NewConfig("std env"),
		},
		{
			name: "std env - container",
			opts: []EnvOption{
				Container("example.container"),
			},
			want: env.NewConfig("std env - container").SetContainer("example.container"),
		},
		{
			name: "std env - aliases",
			opts: []EnvOption{
				Abbrevs("example.type.name"),
			},
			want: env.NewConfig("std env - aliases").AddImports(env.NewImport("example.type.name")),
		},
		{
			name: "std env disabled",
			opts: []EnvOption{
				func(*Env) (*Env, error) {
					return NewCustomEnv()
				},
			},
			want: env.NewConfig("std env disabled").SetStdLib(
				env.NewLibrarySubset().SetDisabled(true)),
		},
		{
			name: "std env - with variable",
			opts: []EnvOption{
				Variable("var", IntType),
			},
			want: env.NewConfig("std env - with variable").AddVariables(env.NewVariable("var", env.NewTypeDesc("int"))),
		},
		{
			name: "std env - with function",
			opts: []EnvOption{Function("hello", Overload("hello_string", []*Type{StringType}, StringType))},
			want: env.NewConfig("std env - with function").AddFunctions(
				env.NewFunction("hello",
					env.NewOverload("hello_string",
						[]*env.TypeDesc{env.NewTypeDesc("string")}, env.NewTypeDesc("string")),
				)),
		},
		{
			name: "optional lib",
			opts: []EnvOption{
				OptionalTypes(),
			},
			want: env.NewConfig("optional lib").AddExtensions(env.NewExtension("optional", math.MaxUint32)),
		},
		{
			name: "optional lib - versioned",
			opts: []EnvOption{
				OptionalTypes(OptionalTypesVersion(1)),
			},
			want: env.NewConfig("optional lib - versioned").AddExtensions(env.NewExtension("optional", 1)),
		},
		{
			name: "optional lib - alt last()",
			opts: []EnvOption{
				OptionalTypes(),
				Function("last", MemberOverload("string_last", []*Type{StringType}, StringType)),
			},
			want: env.NewConfig("optional lib - alt last()").
				AddExtensions(env.NewExtension("optional", math.MaxUint32)).
				AddFunctions(env.NewFunction("last",
					env.NewMemberOverload("string_last", env.NewTypeDesc("string"), []*env.TypeDesc{}, env.NewTypeDesc("string")),
				)),
		},
		{
			name: "context proto - with extra variable",
			opts: []EnvOption{
				DeclareContextProto((&proto3pb.TestAllTypes{}).ProtoReflect().Descriptor()),
				Variable("extra", StringType),
			},
			want: env.NewConfig("context proto - with extra variable").
				SetContextVariable(env.NewContextVariable("google.expr.proto3.test.TestAllTypes")).
				AddVariables(env.NewVariable("extra", env.NewTypeDesc("string"))),
		},
		{
			name: "context proto",
			opts: []EnvOption{
				DeclareContextProto((&proto3pb.TestAllTypes{}).ProtoReflect().Descriptor()),
			},
			want: env.NewConfig("context proto").SetContextVariable(env.NewContextVariable("google.expr.proto3.test.TestAllTypes")),
		},
	}

	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			e, err := NewEnv(tc.opts...)
			if err != nil {
				t.Fatalf("NewEnv() failed: %v", err)
			}
			gotConfig, err := e.ToConfig(tc.name)
			if err != nil {
				t.Fatalf("ToConfig() failed: %v", err)
			}
			if !reflect.DeepEqual(gotConfig, tc.want) {
				t.Errorf("e.Config() got %v, wanted %v", gotConfig, tc.want)
			}
		})
	}
}

func TestEnvFromConfig(t *testing.T) {
	type exprCase struct {
		name string
		in   any
		expr string
		iss  error
		out  ref.Val
	}
	tests := []struct {
		name       string
		beforeOpts []EnvOption
		afterOpts  []EnvOption
		conf       *env.Config
		exprs      []exprCase
	}{
		{
			name: "std env",
			conf: env.NewConfig("std env"),
			exprs: []exprCase{
				{
					name: "literal",
					expr: "'hello world'",
					out:  types.String("hello world"),
				},
				{
					name: "size",
					expr: "'hello world'.size()",
					out:  types.Int(11),
				},
			},
		},
		{
			name:       "std env - imports",
			beforeOpts: []EnvOption{Types(&proto3pb.TestAllTypes{})},
			conf: env.NewConfig("std env - context proto").
				AddImports(env.NewImport("google.expr.proto3.test.TestAllTypes")),
			exprs: []exprCase{
				{
					name: "literal",
					expr: "TestAllTypes{single_int64: 15}.single_int64",
					out:  types.Int(15),
				},
			},
		},
		{
			name:       "std env - context proto",
			beforeOpts: []EnvOption{Types(&proto3pb.TestAllTypes{})},
			conf: env.NewConfig("std env - context proto").
				SetContainer("google.expr.proto3.test").
				SetContextVariable(env.NewContextVariable("google.expr.proto3.test.TestAllTypes")),
			exprs: []exprCase{
				{
					name: "field select literal",
					in:   mustContextProto(t, &proto3pb.TestAllTypes{SingleInt64: 10}),
					expr: "TestAllTypes{single_int64: single_int64}.single_int64",
					out:  types.Int(10),
				},
			},
		},
		{
			name:       "custom env - variables",
			beforeOpts: []EnvOption{Types(&proto3pb.TestAllTypes{})},
			conf: env.NewConfig("custom env - variables").
				SetStdLib(env.NewLibrarySubset().SetDisabled(true)).
				SetContainer("google.expr.proto3.test").
				AddVariables(env.NewVariable("single_int64", env.NewTypeDesc("int"))),
			exprs: []exprCase{
				{
					name: "field select literal",
					in:   map[string]any{"single_int64": 42},
					expr: "TestAllTypes{single_int64: single_int64}.single_int64",
					out:  types.Int(42),
				},
				{
					name: "invalid operator",
					in:   map[string]any{"single_int64": 42},
					expr: "TestAllTypes{single_int64: single_int64}.single_int64 + 1",
					iss:  errors.New("undeclared reference"),
				},
			},
		},
		{
			name: "custom env - functions",
			afterOpts: []EnvOption{
				Function("plus",
					MemberOverload("int_plus_int", []*Type{IntType, IntType}, IntType,
						BinaryBinding(func(lhs, rhs ref.Val) ref.Val {
							l := lhs.(types.Int)
							r := rhs.(types.Int)
							return l + r
						}),
					),
				)},
			conf: env.NewConfig("custom env - functions").
				SetStdLib(env.NewLibrarySubset().SetDisabled(true)).
				AddVariables(env.NewVariable("x", env.NewTypeDesc("int"))).
				AddFunctions(env.NewFunction("plus",
					env.NewMemberOverload("int_plus_int",
						env.NewTypeDesc("int"),
						[]*env.TypeDesc{env.NewTypeDesc("int")},
						env.NewTypeDesc("int"),
					),
				)),
			exprs: []exprCase{
				{
					name: "plus",
					in:   map[string]any{"x": 42},
					expr: "x.plus(2)",
					out:  types.Int(44),
				},
				{
					name: "plus invalid type",
					in:   map[string]any{"x": 42},
					expr: "x.plus(2.0)",
					iss:  errors.New("no matching overload"),
				},
			},
		},
		{
			name: "pure custom env",
			beforeOpts: []EnvOption{func(*Env) (*Env, error) {
				return NewCustomEnv()
			}},
			conf: env.NewConfig("pure custom env").SetStdLib(
				env.NewLibrarySubset().AddIncludedFunctions([]*env.Function{{Name: "_==_"}}...),
			),
			exprs: []exprCase{
				{
					name: "equals",
					expr: "'hello world' == 'hello'",
					out:  types.False,
				},
				{
					name: "not equals - invalid",
					expr: "'hello world' != 'hello'",
					iss:  errors.New("undeclared reference"),
				},
			},
		},
		{
			name: "std env - allow subset",
			conf: env.NewConfig("std env - allow subset").SetStdLib(
				env.NewLibrarySubset().AddIncludedFunctions([]*env.Function{{Name: "_==_"}}...),
			),
			exprs: []exprCase{
				{
					name: "equals",
					expr: "'hello world' == 'hello'",
					out:  types.False,
				},
				{
					name: "not equals - invalid",
					expr: "'hello world' != 'hello'",
					iss:  errors.New("undeclared reference"),
				},
			},
		},
		{
			name: "std env - deny subset",
			conf: env.NewConfig("std env - deny subset").SetStdLib(
				env.NewLibrarySubset().AddExcludedFunctions([]*env.Function{{Name: "size"}}...),
			),
			exprs: []exprCase{
				{
					name: "size - invalid",
					expr: "'hello world'.size()",
					iss:  errors.New("undeclared reference"),
				},
				{
					name: "equals",
					expr: "'hello world' == 'hello'",
					out:  types.False,
				},
			},
		},
		{
			name: "extensions",
			conf: env.NewConfig("extensions").
				AddVariables(
					env.NewVariable("m",
						env.NewTypeDesc("map", env.NewTypeDesc("string"), env.NewTypeDesc("string")))).
				AddExtensions(env.NewExtension("optional", math.MaxUint32)),
			exprs: []exprCase{
				{
					name: "optional none",
					expr: "optional.none()",
					out:  types.OptionalNone,
				},
				{
					name: "optional key",
					expr: "m.?key.hasValue()",
					in:   map[string]any{"m": map[string]string{"key": "value"}},
					out:  types.True,
				},
			},
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			opts := tc.beforeOpts
			opts = append(opts, FromConfig(tc.conf, func(elem any) (EnvOption, bool) {
				if ext, ok := elem.(*env.Extension); ok && ext.Name == "optional" {
					ver, _ := ext.GetVersion()
					return OptionalTypes(OptionalTypesVersion(ver)), true
				}
				return nil, false
			}))
			opts = append(opts, tc.afterOpts...)
			var e *Env
			var err error
			if tc.conf.StdLib != nil {
				e, err = NewCustomEnv(opts...)
			} else {
				e, err = NewEnv(opts...)
			}
			if err != nil {
				t.Fatalf("NewEnv(FromConfig()) failed: %v", err)
			}
			for _, ex := range tc.exprs {
				t.Run(ex.name, func(t *testing.T) {
					ast, iss := e.Compile(ex.expr)
					if iss.Err() != nil {
						if ex.iss == nil || !strings.Contains(iss.Err().Error(), ex.iss.Error()) {
							t.Errorf("e.Compile() failed with %v, wanted %v", iss.Err(), ex.iss)
						}
						return
					}
					if ex.iss != nil {
						t.Fatalf("e.Compile() succeeded, wanted error %v", ex.iss)
					}
					prg, err := e.Program(ast)
					if err != nil {
						t.Fatalf("e.Program() failed: %v", err)
					}
					var in any = map[string]any{}
					if ex.in != nil {
						in = ex.in
					}
					out, _, err := prg.Eval(in)
					if err != nil {
						t.Fatalf("prg.Eval() failed: %v", err)
					}
					if out.Equal(ex.out) != types.True {
						t.Errorf("prg.Eval() got %v, wanted %v", out, ex.out)
					}
				})
			}
		})
	}
}

func TestEnvFromConfigErrors(t *testing.T) {
	tests := []struct {
		name string
		conf *env.Config
		want error
	}{
		{
			name: "invalid subset",
			conf: env.NewConfig("invalid subset").SetStdLib(env.NewLibrarySubset().SetDisableMacros(true)),
			want: errors.New("invalid subset"),
		},
		{
			name: "invalid import",
			conf: env.NewConfig("invalid import").AddImports(env.NewImport("")),
			want: errors.New("invalid import"),
		},
		{
			name: "invalid context proto",
			conf: env.NewConfig("invalid context proto").SetContextVariable(env.NewContextVariable("invalid")),
			want: errors.New("invalid context proto type"),
		},
		{
			name: "undefined variable type",
			conf: env.NewConfig("undefined variable type").AddVariables(env.NewVariable("undef", env.NewTypeDesc("undefined"))),
			want: errors.New("invalid variable"),
		},
		{
			name: "undefined function type",
			conf: env.NewConfig("undefined function type").AddFunctions(env.NewFunction("invalid", env.NewOverload("invalid", []*env.TypeDesc{}, env.NewTypeDesc("undefined")))),
			want: errors.New("invalid function"),
		},
		{
			name: "unrecognized extension",
			conf: env.NewConfig("unrecognized extension").
				AddExtensions(env.NewExtension("optional", math.MaxUint32)),
			want: errors.New("unrecognized extension"),
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewEnv(FromConfig(tc.conf))
			if err == nil || !strings.Contains(err.Error(), tc.want.Error()) {
				t.Fatalf("NewEnv(FromConfig()) got %v, wanted error containing %v", err, tc.want)
			}
		})
	}
}

func BenchmarkNewCustomEnvLazy(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewCustomEnv(StdLib(), EagerlyValidateDeclarations(false))
		if err != nil {
			b.Fatalf("NewCustomEnv() failed: %v", err)
		}
	}
}

func BenchmarkNewCustomEnvEager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewCustomEnv(StdLib(), EagerlyValidateDeclarations(true))
		if err != nil {
			b.Fatalf("NewCustomEnv() failed: %v", err)
		}
	}
}

func BenchmarkNewEnvLazy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := NewEnv()
		if err != nil {
			b.Fatalf("NewEnv() failed: %v", err)
		}
	}
}

func BenchmarkNewEnvEager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		env, err := NewEnv()
		if err != nil {
			b.Fatalf("NewEnv() failed: %v", err)
		}
		_, iss := env.Compile("123")
		if iss.Err() != nil {
			b.Fatalf("env.Compile(123) failed: %v", iss.Err())
		}
	}
}

func BenchmarkEnvExtendEager(b *testing.B) {
	env, err := NewEnv()
	if err != nil {
		b.Fatalf("NewEnv() failed: %v", err)
	}
	for i := 0; i < b.N; i++ {
		ext, err := env.Extend()
		if err != nil {
			b.Fatalf("env.Extend() failed: %v", err)
		}
		_, iss := ext.Compile("123")
		if iss.Err() != nil {
			b.Fatalf("env.Compile(123) failed: %v", iss.Err())
		}
	}
}

func BenchmarkEnvExtendEagerTypes(b *testing.B) {
	env, err := NewEnv(EagerlyValidateDeclarations(true))
	if err != nil {
		b.Fatalf("NewEnv() failed: %v", err)
	}
	for i := 0; i < b.N; i++ {
		ext, err := env.Extend(CustomTypeProvider(types.NewEmptyRegistry()))
		if err != nil {
			b.Fatalf("env.Extend() failed: %v", err)
		}
		_, iss := ext.Compile("123")
		if iss.Err() != nil {
			b.Fatalf("env.Compile(123) failed: %v", iss.Err())
		}
	}
}

func BenchmarkEnvExtendEagerDecls(b *testing.B) {
	env, err := NewEnv(EagerlyValidateDeclarations(true))
	if err != nil {
		b.Fatalf("NewEnv() failed: %v", err)
	}
	for i := 0; i < b.N; i++ {
		ext, err := env.Extend(
			Variable("test_var", StringType),
			Function(
				operators.In,
				Overload("string_in_string", []*Type{StringType, StringType}, BoolType),
			),
		)
		if err != nil {
			b.Fatalf("env.Extend() failed: %v", err)
		}
		_, iss := ext.Compile("123")
		if iss.Err() != nil {
			b.Fatalf("env.Compile(123) failed: %v", iss.Err())
		}
	}
}

func mustContextProto(t *testing.T, pb proto.Message) Activation {
	t.Helper()
	ctx, err := ContextProtoVars(pb)
	if err != nil {
		t.Fatalf("ContextProtoVars() failed: %v", err)
	}
	return ctx
}

type customLegacyProvider struct {
	provider ref.TypeProvider
}

func (p *customLegacyProvider) EnumValue(enumName string) ref.Val {
	return p.provider.EnumValue(enumName)
}

func (p *customLegacyProvider) FindIdent(identName string) (ref.Val, bool) {
	return p.provider.FindIdent(identName)
}

func (p *customLegacyProvider) FindType(typeName string) (*exprpb.Type, bool) {
	return p.provider.FindType(typeName)
}

func (p *customLegacyProvider) FindFieldType(structType, fieldName string) (*ref.FieldType, bool) {
	return p.provider.FindFieldType(structType, fieldName)
}

func (p *customLegacyProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	return p.provider.NewValue(structType, fields)
}

type customCELProvider struct {
	provider types.Provider
}

func (p *customCELProvider) EnumValue(enumName string) ref.Val {
	return p.provider.EnumValue(enumName)
}

func (p *customCELProvider) FindIdent(identName string) (ref.Val, bool) {
	return p.provider.FindIdent(identName)
}

func (p *customCELProvider) FindStructType(typeName string) (*types.Type, bool) {
	return p.provider.FindStructType(typeName)
}

func (p *customCELProvider) FindStructFieldNames(typeName string) ([]string, bool) {
	return p.provider.FindStructFieldNames(typeName)
}

func (p *customCELProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {
	return p.provider.FindStructFieldType(structType, fieldName)
}

func (p *customCELProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {
	return p.provider.NewValue(structType, fields)
}
