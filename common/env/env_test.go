// Copyright 2025 Google LLC
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

package env

import (
	"errors"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "context_env",
			want: &Config{
				Name:      "context-env",
				Container: "google.expr",
				Imports:   []*Import{NewImport("google.expr.proto3.test.TestAllTypes")},
				StdLib: &LibrarySubset{
					IncludeMacros: []string{"has"},
					IncludeFunctions: []*Function{
						{Name: operators.Equals},
						{Name: operators.NotEquals},
						{Name: operators.LogicalNot},
						{Name: operators.Less},
						{Name: operators.LessEquals},
						{Name: operators.Greater},
						{Name: operators.GreaterEquals},
					},
				},
				Extensions: []*Extension{
					NewExtension("optional", math.MaxUint32),
					NewExtension("strings", 1),
				},
				ContextVariable: NewContextVariable("google.expr.proto3.test.TestAllTypes"),
				Functions: []*Function{
					NewFunction("coalesce", []*Overload{
						NewOverload("coalesce_wrapped_int",
							[]*TypeDesc{NewTypeDesc("google.protobuf.Int64Value"), NewTypeDesc("int")},
							NewTypeDesc("int")),
						NewOverload("coalesce_wrapped_double",
							[]*TypeDesc{NewTypeDesc("google.protobuf.DoubleValue"), NewTypeDesc("double")},
							NewTypeDesc("double")),
						NewOverload("coalesce_wrapped_uint",
							[]*TypeDesc{NewTypeDesc("google.protobuf.UInt64Value"), NewTypeDesc("uint")},
							NewTypeDesc("uint")),
					}),
				},
			},
		},
		{
			name: "extended_env",
			want: &Config{
				Name:      "extended-env",
				Container: "google.expr",
				Extensions: []*Extension{
					NewExtension("optional", 2),
					NewExtension("math", math.MaxUint32),
				},
				Variables: []*Variable{
					NewVariable("msg", NewTypeDesc("google.expr.proto3.test.TestAllTypes")),
				},
				Functions: []*Function{
					NewFunction("isEmpty", []*Overload{
						NewMemberOverload("wrapper_string_isEmpty",
							NewTypeDesc("google.protobuf.StringValue"), nil,
							NewTypeDesc("bool")),
						NewMemberOverload("list_isEmpty",
							NewTypeDesc("list", NewTypeParam("T")), nil,
							NewTypeDesc("bool")),
					}),
				},
			},
		},
		{
			name: "subset_env",
			want: &Config{
				Name: "subset-env",
				StdLib: &LibrarySubset{
					ExcludeMacros: []string{"map", "filter"},
					ExcludeFunctions: []*Function{
						{Name: operators.Add, Overloads: []*Overload{
							{ID: overloads.AddBytes},
							{ID: overloads.AddList},
							{ID: overloads.AddString},
						}},
						{Name: overloads.Matches},
						{Name: overloads.TypeConvertTimestamp, Overloads: []*Overload{
							{ID: overloads.StringToTimestamp},
						}},
						{Name: overloads.TypeConvertDuration, Overloads: []*Overload{
							{ID: overloads.StringToDuration},
						}},
					},
				},
				Variables: []*Variable{
					NewVariable("x", NewTypeDesc("int")),
					NewVariable("y", NewTypeDesc("double")),
					NewVariable("z", NewTypeDesc("uint")),
				},
			},
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			fileName := fmt.Sprintf("testdata/%s.yaml", tc.name)
			data, err := os.ReadFile(fileName)
			if err != nil {
				t.Fatalf("os.ReadFile(%q) failed: %v", fileName, err)
			}
			got := unmarshalYAML(t, data)
			if got.Container != tc.want.Container {
				t.Errorf("Container got %s, wanted %s", got.Container, tc.want.Container)
			}
			if !reflect.DeepEqual(got.Imports, tc.want.Imports) {
				t.Errorf("Imports got %v, wanted %v", got.Imports, tc.want.Imports)
			}
			if !reflect.DeepEqual(got.StdLib, tc.want.StdLib) {
				t.Errorf("StdLib got %v, wanted %v", got.StdLib, tc.want.StdLib)
			}
			if !reflect.DeepEqual(got.ContextVariable, tc.want.ContextVariable) {
				t.Errorf("ContextVariable got %v, wanted %v", got.ContextVariable, tc.want.ContextVariable)
			}
			if len(got.Variables) != len(tc.want.Variables) {
				t.Errorf("Variables count got %d, wanted %d", len(got.Variables), len(tc.want.Variables))
			} else {
				for i, v := range got.Variables {
					wv := tc.want.Variables[i]
					if !reflect.DeepEqual(v, wv) {
						t.Errorf("Variables[%d] not equal, got %v, wanted %v", i, v, wv)
					}
				}
			}
			if len(got.Functions) != len(tc.want.Functions) {
				t.Errorf("Functions count got %d, wanted %d", len(got.Functions), len(tc.want.Functions))
			} else {
				for i, f := range got.Functions {
					wf := tc.want.Functions[i]
					if f.Name != wf.Name {
						t.Errorf("Functions[%d] not equal, got %v, wanted %v", i, f.Name, wf.Name)
					}
					if len(f.Overloads) != len(wf.Overloads) {
						t.Errorf("Function %s got overload count: %d, wanted %d", f.Name, len(f.Overloads), len(wf.Overloads))
					}
					for j, o := range f.Overloads {
						wo := wf.Overloads[j]
						if !reflect.DeepEqual(o, wo) {
							t.Errorf("Overload[%d] got %v, wanted %v", j, o, wo)
						}
					}
				}
			}
		})
	}
}

func TestNewImport(t *testing.T) {
	imp := NewImport("qualified.type.name")
	if imp.Name != "qualified.type.name" {
		t.Errorf("NewImport() got name: %s, wanted %s", imp.Name, "qualified.type.name")
	}
}

func TestNewContextVariable(t *testing.T) {
	ctx := NewContextVariable("qualified.type.name")
	if ctx.TypeName != "qualified.type.name" {
		t.Errorf("NewContextVariable() got name: %s, wanted %s", ctx.TypeName, "qualified.type.name")
	}
}

func TestConfigAddVariableDecls(t *testing.T) {
	tests := []struct {
		name string
		in   *decls.VariableDecl
		out  *Variable
	}{
		{
			name: "nil var decl",
		},
		{
			name: "simple var decl",
			in:   decls.NewVariable("var", types.StringType),
			out:  NewVariable("var", NewTypeDesc("string")),
		},
		{
			name: "parameterized var decl",
			in:   decls.NewVariable("var", types.NewListType(types.NewTypeParamType("T"))),
			out:  NewVariable("var", NewTypeDesc("list", NewTypeParam("T"))),
		},
		{
			name: "opaque var decl",
			in:   decls.NewVariable("var", types.NewOpaqueType("bitvector")),
			out:  NewVariable("var", NewTypeDesc("bitvector")),
		},
		{
			name: "proto var decl",
			in:   decls.NewVariable("var", types.NewObjectType("google.type.Expr")),
			out:  NewVariable("var", NewTypeDesc("google.type.Expr")),
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			conf := NewConfig().AddVariableDecls(tc.in)
			if len(conf.Variables) != 1 {
				t.Fatalf("AddVariableDecls() did not add declaration to conf: %v", conf)
			}
			if !reflect.DeepEqual(conf.Variables[0], tc.out) {
				t.Errorf("AddVariableDecls() added %v, wanted %v", conf.Variables, tc.out)
			}
		})
	}
}

func TestConfigAddVariableDeclsEmpty(t *testing.T) {
	if len(NewConfig().AddVariables().Variables) != 0 {
		t.Error("AddVariables() with no args failed")
	}
}

func TestConfigAddFunctionDecls(t *testing.T) {
	tests := []struct {
		name string
		in   *decls.FunctionDecl
		out  *Function
	}{
		{
			name: "nil function decl",
		},
		{
			name: "global function decl",
			in: mustNewFunction(t, "size",
				decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType),
			),
			out: NewFunction("size", []*Overload{
				NewOverload("size_string", []*TypeDesc{NewTypeDesc("string")}, NewTypeDesc("int")),
			}),
		},
		{
			name: "global function decl - nullable arg",
			in: mustNewFunction(t, "size",
				decls.Overload("size_wrapper_string", []*types.Type{types.NewNullableType(types.StringType)}, types.IntType),
			),
			out: NewFunction("size", []*Overload{
				NewOverload("size_wrapper_string", []*TypeDesc{NewTypeDesc("google.protobuf.StringValue")}, NewTypeDesc("int")),
			}),
		},
		{
			name: "member function decl - nullable arg",
			in: mustNewFunction(t, "size",
				decls.MemberOverload("list_size", []*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
				decls.MemberOverload("string_size", []*types.Type{types.StringType}, types.IntType),
			),
			out: NewFunction("size", []*Overload{
				NewMemberOverload("list_size", NewTypeDesc("list", NewTypeParam("T")), []*TypeDesc{}, NewTypeDesc("int")),
				NewMemberOverload("string_size", NewTypeDesc("string"), []*TypeDesc{}, NewTypeDesc("int")),
			}),
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			conf := NewConfig().AddFunctionDecls(tc.in)
			if len(conf.Functions) != 1 {
				t.Fatalf("AddFunctionDecls() did not add declaration to conf: %v", conf)
			}
			if !reflect.DeepEqual(conf.Functions[0], tc.out) {
				t.Errorf("AddFunctionDecls() added %v, wanted %v", conf.Functions, tc.out)
			}
		})
	}
}

func TestVariableGetType(t *testing.T) {
	tests := []struct {
		name string
		v    *Variable
		t    *TypeDesc
	}{
		{
			name: "nil-safety check",
			v:    nil,
			t:    nil,
		},
		{
			name: "nil type access",
			v:    &Variable{},
			t:    nil,
		},
		{
			name: "nested type desc",
			v:    &Variable{TypeDesc: &TypeDesc{}},
			t:    &TypeDesc{},
		},
		{
			name: "field type desc",
			v:    &Variable{Type: &TypeDesc{}},
			t:    &TypeDesc{},
		},
		{
			name: "nested type desc precedence",
			v: &Variable{
				TypeDesc: &TypeDesc{TypeName: "type.name.EmbeddedType"},
				Type:     &TypeDesc{TypeName: "type.name.FieldType"},
			},
			t: &TypeDesc{TypeName: "type.name.EmbeddedType"},
		},
	}

	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.v.GetType(), tc.t) {
				t.Errorf("GetType() got %v, wanted %v", tc.v.GetType(), tc.t)
			}
		})
	}
}

func TestVariableAsCELVariable(t *testing.T) {
	tests := []struct {
		name string
		v    *Variable
		want any
	}{
		{
			name: "nil-safety check",
			v:    nil,
			want: errors.New("nil Variable"),
		},
		{
			name: "no variable name",
			v:    &Variable{},
			want: errors.New("invalid variable"),
		},
		{
			name: "no type",
			v: &Variable{
				Name: "hello",
			},
			want: errors.New("no type specified"),
		},
		{
			name: "bad type",
			v: &Variable{
				Name:     "hello",
				TypeDesc: &TypeDesc{},
			},
			want: errors.New("declare a type name"),
		},
		{
			name: "int type",
			v:    NewVariable("int_var", NewTypeDesc("int")),
			want: decls.NewVariable("int_var", types.IntType),
		},
		{
			name: "uint type",
			v: &Variable{
				Name:     "uint_var",
				TypeDesc: NewTypeDesc("uint"),
			},
			want: decls.NewVariable("uint_var", types.UintType),
		},
		{
			name: "dyn type",
			v: &Variable{
				Name:     "dyn_var",
				TypeDesc: NewTypeDesc("dyn"),
			},
			want: decls.NewVariable("dyn_var", types.DynType),
		},
		{
			name: "list type",
			v: &Variable{
				Name:     "list_var",
				TypeDesc: NewTypeDesc("list", NewTypeParam("T")),
			},
			want: decls.NewVariable("list_var", types.NewListType(types.NewTypeParamType("T"))),
		},
		{
			name: "map type",
			v: &Variable{
				Name: "map_var",
				TypeDesc: &TypeDesc{
					TypeName: "map",
					Params: []*TypeDesc{
						NewTypeDesc("string"),
						NewTypeDesc("optional_type", NewTypeParam("T")),
					},
				},
			},
			want: decls.NewVariable("map_var",
				types.NewMapType(types.StringType, types.NewOptionalType(types.NewTypeParamType("T")))),
		},
		{
			name: "set type",
			v: &Variable{
				Name:     "set_var",
				TypeDesc: NewTypeDesc("set", NewTypeDesc("string")),
			},
			want: decls.NewVariable("set_var", types.NewOpaqueType("set", types.StringType)),
		},
		{
			name: "string type - nested type precedence",
			v: &Variable{
				Name:     "hello",
				TypeDesc: NewTypeDesc("string"),
				Type:     NewTypeDesc("int"),
			},
			want: decls.NewVariable("hello", types.StringType),
		},
		{
			name: "wrapper type variable",
			v: &Variable{
				Name:     "msg",
				TypeDesc: NewTypeDesc("google.protobuf.StringValue"),
			},
			want: decls.NewVariable("msg", types.NewNullableType(types.StringType)),
		},
	}

	tp, err := types.NewRegistry()
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	tp.RegisterType(types.NewOpaqueType("set", types.NewTypeParamType("T")))
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			gotVar, err := tc.v.AsCELVariable(tp)
			if err != nil {
				wantErr, ok := tc.want.(error)
				if !ok {
					t.Fatalf("AsCELVariable() got error %v, wanted %v", err, tc.want)
				}
				if !strings.Contains(err.Error(), wantErr.Error()) {
					t.Fatalf("AsCELVariable() got error %v, wanted error contining %v", err, wantErr)
				}
				return
			}
			if !gotVar.DeclarationIsEquivalent(tc.want.(*decls.VariableDecl)) {
				t.Errorf("AsCELVariable() got %v, wanted %v", gotVar, tc.want)
			}
		})
	}
}

func TestTypeDescString(t *testing.T) {
	tests := []struct {
		desc *TypeDesc
		want string
	}{
		{desc: NewTypeDesc("string"), want: "string"},
		{desc: NewTypeDesc("list", NewTypeParam("T")), want: "list(T)"},
		{desc: NewTypeDesc("map", NewTypeDesc("string"), NewTypeParam("T")), want: "map(string,T)"},
	}
	for _, tc := range tests {
		if tc.desc.String() != tc.want {
			t.Errorf("String() got %s, wanted %s", tc.desc.String(), tc.want)
		}
	}
}

func TestFunctionAsCELFunction(t *testing.T) {
	tests := []struct {
		name string
		f    *Function
		want any
	}{
		{
			name: "nil function",
			f:    nil,
			want: errors.New("nil Function"),
		},
		{
			name: "unnamed function",
			f:    &Function{},
			want: errors.New("must declare a name"),
		},
		{
			name: "no overloads",
			f:    NewFunction("no_overloads", []*Overload{}),
			want: errors.New("must declare an overload"),
		},
		{
			name: "nil overload",
			f:    NewFunction("no_overloads", []*Overload{nil}),
			want: errors.New("nil Overload"),
		},
		{
			name: "no return type",
			f: NewFunction("size", []*Overload{
				NewOverload("size_string", []*TypeDesc{NewTypeDesc("string")}, nil),
			}),
			want: errors.New("missing return type"),
		},
		{
			name: "bad return type",
			f: NewFunction("size", []*Overload{
				NewOverload("size_string", []*TypeDesc{NewTypeDesc("string")}, NewTypeDesc("")),
			}),
			want: errors.New("invalid type"),
		},
		{
			name: "bad arg type",
			f: NewFunction("size", []*Overload{
				NewOverload("size_string", []*TypeDesc{NewTypeDesc("")}, NewTypeDesc("")),
			}),
			want: errors.New("invalid type"),
		},
		{
			name: "bad target type",
			f: &Function{Name: "size",
				Overloads: []*Overload{
					{ID: "string_size",
						Target: &TypeDesc{},
						Args:   []*TypeDesc{},
						Return: &TypeDesc{TypeName: "int"},
					},
				},
			},
			want: errors.New("invalid type"),
		},
		{
			name: "global function",
			f: &Function{Name: "size",
				Overloads: []*Overload{
					{ID: "size_string",
						Args:   []*TypeDesc{{TypeName: "string"}},
						Return: &TypeDesc{TypeName: "int"}},
				},
			},
			want: mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
		},
		{
			name: "member function",
			f: &Function{Name: "size",
				Overloads: []*Overload{
					{ID: "string_size",
						Target: &TypeDesc{TypeName: "string"},
						Return: &TypeDesc{TypeName: "int"}},
				},
			},
			want: mustNewFunction(t, "size", decls.MemberOverload("string_size", []*types.Type{types.StringType}, types.IntType)),
		},
	}
	tp, err := types.NewRegistry()
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	tp.RegisterType(types.NewOpaqueType("set", types.NewTypeParamType("T")))
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			gotFn, err := tc.f.AsCELFunction(tp)
			if err != nil {
				wantErr, ok := tc.want.(error)
				if !ok {
					t.Fatalf("AsCELFunction() got error %v, wanted %v", err, tc.want)
				}
				if !strings.Contains(err.Error(), wantErr.Error()) {
					t.Fatalf("AsCELFunction() got error %v, wanted error contining %v", err, wantErr)
				}
				return
			}
			assertFuncEquals(t, gotFn, tc.want.(*decls.FunctionDecl))
		})
	}
}

func TestTypeDescAsCELTypeErrors(t *testing.T) {
	tests := []struct {
		name string
		t    *TypeDesc
		want any
	}{
		{
			name: "nil-safety check",
			t:    nil,
			want: errors.New("nil TypeDesc"),
		},
		{
			name: "no type name",
			t:    &TypeDesc{},
			want: errors.New("invalid type"),
		},
		{
			name: "invalid optional_type",
			t:    &TypeDesc{TypeName: "optional_type"},
			want: errors.New("unexpected param count"),
		},
		{
			name: "invalid optional param type",
			t:    &TypeDesc{TypeName: "optional_type", Params: []*TypeDesc{{}}},
			want: errors.New("invalid type"),
		},
		{
			name: "invalid list",
			t:    &TypeDesc{TypeName: "list"},
			want: errors.New("unexpected param count"),
		},
		{
			name: "invalid list param type",
			t:    &TypeDesc{TypeName: "list", Params: []*TypeDesc{{}}},
			want: errors.New("invalid type"),
		},
		{
			name: "invalid map",
			t:    &TypeDesc{TypeName: "map"},
			want: errors.New("unexpected param count"),
		},
		{
			name: "invalid map key type",
			t:    &TypeDesc{TypeName: "map", Params: []*TypeDesc{{}, {}}},
			want: errors.New("invalid type"),
		},
		{
			name: "invalid map value type",
			t:    &TypeDesc{TypeName: "map", Params: []*TypeDesc{{TypeName: "string"}, {}}},
			want: errors.New("invalid type"),
		},
		{
			name: "invalid set",
			t:    &TypeDesc{TypeName: "set", Params: []*TypeDesc{{}}},
			want: errors.New("invalid type"),
		},
		{
			name: "undefined type identifier",
			t:    &TypeDesc{TypeName: "vector"},
			want: errors.New("undefined type"),
		},
	}
	tp, err := types.NewRegistry()
	if err != nil {
		t.Fatalf("types.NewRegistry() failed: %v", err)
	}
	tp.RegisterType(types.NewOpaqueType("set", types.NewTypeParamType("T")))
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			gotVar, err := tc.t.AsCELType(tp)
			if err != nil {
				wantErr, ok := tc.want.(error)
				if !ok {
					t.Fatalf("AsCELType() got error %v, wanted %v", err, tc.want)
				}
				if !strings.Contains(err.Error(), wantErr.Error()) {
					t.Fatalf("AsCELType() got error %v, wanted error contining %v", err, wantErr)
				}
				return
			}
			if !reflect.DeepEqual(gotVar, tc.want.(*decls.VariableDecl)) {
				t.Errorf("AsCELType() got %v, wanted %v", gotVar, tc.want)
			}
		})
	}
}

func TestSubsetFunction(t *testing.T) {
	tests := []struct {
		name     string
		lib      *LibrarySubset
		orig     *decls.FunctionDecl
		subset   *decls.FunctionDecl
		included bool
	}{
		{
			name:     "nil lib, included",
			orig:     mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			subset:   mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: true,
		},
		{
			name:     "empty, included",
			lib:      &LibrarySubset{},
			orig:     mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			subset:   mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: true,
		},
		{
			name: "lib, not included allow-list",
			lib: &LibrarySubset{
				IncludeFunctions: []*Function{
					{Name: "int"},
				},
			},
			orig:     mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: false,
		},
		{
			name: "lib, included whole function",
			lib: &LibrarySubset{
				IncludeFunctions: []*Function{
					{Name: "size"},
				},
			},
			orig:     mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			subset:   mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: true,
		},
		{
			name: "lib, included overload subset",
			lib: &LibrarySubset{
				IncludeFunctions: []*Function{
					{Name: "size", Overloads: []*Overload{{ID: "size_string"}}},
				},
			},
			orig: mustNewFunction(t, "size",
				decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType),
				decls.Overload("size_list", []*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
			),
			subset:   mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: true,
		},
		{
			name: "lib, included deny-list",
			lib: &LibrarySubset{
				ExcludeFunctions: []*Function{
					{Name: "int"},
				},
			},
			orig:     mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			subset:   mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: true,
		},
		{
			name: "lib, excluded whole function",
			lib: &LibrarySubset{
				ExcludeFunctions: []*Function{
					{Name: "size"},
				},
			},
			orig:     mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: false,
		},
		{
			name: "lib, excluded partial function",
			lib: &LibrarySubset{
				ExcludeFunctions: []*Function{
					{Name: "size", Overloads: []*Overload{{ID: "size_list"}}},
				},
			},
			orig: mustNewFunction(t, "size",
				decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType),
				decls.Overload("size_list", []*types.Type{types.NewListType(types.NewTypeParamType("T"))}, types.IntType),
			),
			subset:   mustNewFunction(t, "size", decls.Overload("size_string", []*types.Type{types.StringType}, types.IntType)),
			included: true,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			got, included := tc.lib.SubsetFunction(tc.orig)
			if included != tc.included {
				t.Fatalf("SubsetFunction() got included %t, wanted %t", included, tc.included)
			}
			if !tc.included {
				return
			}
			assertFuncEquals(t, got, tc.subset)
		})
	}
}

func TestSubsetMacro(t *testing.T) {
	tests := []struct {
		name      string
		lib       *LibrarySubset
		macroName string
		included  bool
	}{
		{
			name:      "nil lib, included",
			macroName: "has",
			included:  true,
		},
		{
			name:      "empty, included",
			lib:       NewLibrarySubset(),
			macroName: "has",
			included:  true,
		},
		{
			name:      "empty, included",
			lib:       &LibrarySubset{DisableMacros: true},
			macroName: "has",
			included:  false,
		},
		{
			name: "lib, not included allow-list",
			lib: &LibrarySubset{
				IncludeMacros: []string{"exists"},
			},
			macroName: "has",
			included:  false,
		},
		{
			name: "lib, included allow-list",
			lib: &LibrarySubset{
				IncludeMacros: []string{"exists"},
			},
			macroName: "exists",
			included:  true,
		},
		{
			name: "lib, not included deny-list",
			lib: &LibrarySubset{
				ExcludeMacros: []string{"exists"},
			},
			macroName: "exists",
			included:  false,
		},
		{
			name: "lib, included deny-list",
			lib: &LibrarySubset{
				ExcludeMacros: []string{"exists"},
			},
			macroName: "has",
			included:  true,
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			included := tc.lib.SubsetMacro(tc.macroName)
			if included != tc.included {
				t.Fatalf("SubsetMacro() got included %t, wanted %t", included, tc.included)
			}
		})
	}
}

func TestNewExtension(t *testing.T) {
	tests := []struct {
		name    string
		version uint32
		want    *Extension
	}{
		{
			name:    "strings",
			version: math.MaxUint32,
			want:    &Extension{Name: "strings", Version: "latest"},
		},
		{
			name:    "bindings",
			version: 1,
			want:    &Extension{Name: "bindings", Version: "1"},
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			got := NewExtension(tc.name, tc.version)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("NewExtension() got %v, wanted %v", got, tc.want)
			}
		})
	}
}

func TestExtensionGetVersion(t *testing.T) {
	tests := []struct {
		name string
		ext  *Extension
		want any
	}{
		{
			name: "nil extension",
			want: errors.New("nil Extension"),
		},
		{
			name: "unset version",
			ext:  &Extension{},
			want: uint32(0),
		},
		{
			name: "numeric version",
			ext:  &Extension{Version: "1"},
			want: uint32(1),
		},
		{
			name: "latest version",
			ext:  &Extension{Version: "latest"},
			want: uint32(math.MaxUint32),
		},
		{
			name: "bad version",
			ext:  &Extension{Version: "1.0"},
			want: errors.New("invalid syntax"),
		},
	}
	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			ver, err := tc.ext.GetVersion()
			if err != nil {
				wantErr, ok := tc.want.(error)
				if !ok {
					t.Fatalf("GetVersion() got error %v, wanted %v", err, tc.want)
				}
				if !strings.Contains(err.Error(), wantErr.Error()) {
					t.Fatalf("GetVersion() got error %v, wanted error contining %v", err, wantErr)
				}
				return
			}
			if tc.want.(uint32) != ver {
				t.Fatalf("GetVersion() got %d, wanted %v", ver, tc.want)
			}
		})
	}
}

func mustNewFunction(t *testing.T, name string, opts ...decls.FunctionOpt) *decls.FunctionDecl {
	t.Helper()
	fn, err := decls.NewFunction(name, opts...)
	if err != nil {
		t.Fatalf("decls.NewFunction() failed: %v", err)
	}
	return fn
}

func assertFuncEquals(t *testing.T, got, want *decls.FunctionDecl) {
	t.Helper()
	if got.Name() != want.Name() {
		t.Fatalf("got function name %s, wanted %s", got.Name(), want.Name())
	}
	if len(got.OverloadDecls()) != len(want.OverloadDecls()) {
		t.Fatalf("got overload count %d, wanted %d", len(got.OverloadDecls()), len(want.OverloadDecls()))
	}
	for i, gotOverload := range got.OverloadDecls() {
		wantOverload := want.OverloadDecls()[i]
		if gotOverload.ID() != wantOverload.ID() {
			t.Errorf("got overload id: %s, wanted: %s", gotOverload.ID(), wantOverload.ID())
		}
		if gotOverload.IsMemberFunction() != wantOverload.IsMemberFunction() {
			t.Errorf("got is member function %t, wanted %t", gotOverload.IsMemberFunction(), wantOverload.IsMemberFunction())
		}
		if len(gotOverload.ArgTypes()) != len(wantOverload.ArgTypes()) {
			t.Fatalf("got arg count %d, wanted %d", len(gotOverload.ArgTypes()), len(wantOverload.ArgTypes()))
		}
		for i, p := range gotOverload.ArgTypes() {
			wp := wantOverload.ArgTypes()[i]
			if !p.IsExactType(wp) {
				t.Errorf("got arg[%d] type %v, wanted %v", i, p, wp)
			}
		}
		if len(gotOverload.TypeParams()) != len(wantOverload.TypeParams()) {
			t.Fatalf("got type param count %d, wanted %d", len(gotOverload.TypeParams()), len(wantOverload.TypeParams()))
		}
		for i, p := range gotOverload.TypeParams() {
			wp := wantOverload.TypeParams()[i]
			if p != wp {
				t.Errorf("got type param[%d] %s, wanted %s", i, p, wp)
			}
		}
		if !gotOverload.ResultType().IsExactType(wantOverload.ResultType()) {
			t.Errorf("got result type %v, wanted %v", gotOverload.ResultType(), wantOverload.ResultType())
		}
	}
}

func unmarshalYAML(t *testing.T, data []byte) *Config {
	t.Helper()
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		t.Fatalf("yaml.Unmarshal(%q) failed: %v", string(data), err)
	}
	return config
}

func marshalYAML(t *testing.T, config *Config) []byte {
	t.Helper()
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal(%q) failed: %v", string(data), err)
	}
	return data
}
