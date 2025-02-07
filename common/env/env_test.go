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
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/types"
)

func TestConfig(t *testing.T) {
	conf := NewConfig()
	if conf == nil {
		t.Fatal("got nil config, wanted non-nil value")
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
			v: &Variable{
				Name:     "int_var",
				TypeDesc: &TypeDesc{TypeName: "int"},
			},
			want: decls.NewVariable("int_var", types.IntType),
		},
		{
			name: "uint type",
			v: &Variable{
				Name:     "uint_var",
				TypeDesc: &TypeDesc{TypeName: "uint"},
			},
			want: decls.NewVariable("uint_var", types.UintType),
		},
		{
			name: "dyn type",
			v: &Variable{
				Name:     "dyn_var",
				TypeDesc: &TypeDesc{TypeName: "dyn"},
			},
			want: decls.NewVariable("dyn_var", types.DynType),
		},
		{
			name: "list type",
			v: &Variable{
				Name:     "list_var",
				TypeDesc: &TypeDesc{TypeName: "list", Params: []*TypeDesc{{TypeName: "T", IsTypeParam: true}}},
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
						{TypeName: "string"},
						{TypeName: "optional_type",
							Params: []*TypeDesc{{TypeName: "T", IsTypeParam: true}}},
					},
				},
			},
			want: decls.NewVariable("map_var",
				types.NewMapType(types.StringType, types.NewOptionalType(types.NewTypeParamType("T")))),
		},
		{
			name: "set type",
			v: &Variable{
				Name: "set_var",
				TypeDesc: &TypeDesc{
					TypeName: "set",
					Params: []*TypeDesc{
						{TypeName: "string"},
					},
				},
			},
			want: decls.NewVariable("set_var", types.NewOpaqueType("set", types.StringType)),
		},
		{
			name: "string type - nested type precedence",
			v: &Variable{
				Name:     "hello",
				TypeDesc: &TypeDesc{TypeName: "string"},
				Type:     &TypeDesc{TypeName: "int"},
			},
			want: decls.NewVariable("hello", types.StringType),
		},
		{
			name: "wrapper type variable",
			v: &Variable{
				Name:     "msg",
				TypeDesc: &TypeDesc{TypeName: "google.protobuf.StringValue"},
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
			f:    &Function{Name: "no_overloads"},
			want: errors.New("must declare an overload"),
		},
		{
			name: "nil overload",
			f:    &Function{Name: "no_overloads", Overloads: []*Overload{nil}},
			want: errors.New("nil Overload"),
		},
		{
			name: "no return type",
			f: &Function{Name: "size",
				Overloads: []*Overload{
					{ID: "size_string",
						Args: []*TypeDesc{{TypeName: "string"}},
					},
				},
			},
			want: errors.New("missing return type"),
		},
		{
			name: "bad return type",
			f: &Function{Name: "size",
				Overloads: []*Overload{
					{ID: "size_string",
						Args:   []*TypeDesc{{TypeName: "string"}},
						Return: &TypeDesc{},
					},
				},
			},
			want: errors.New("invalid type"),
		},
		{
			name: "bad arg type",
			f: &Function{Name: "size",
				Overloads: []*Overload{
					{ID: "size_string",
						Args:   []*TypeDesc{{}},
						Return: &TypeDesc{},
					},
				},
			},
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
			name: "invalid optional",
			t:    &TypeDesc{TypeName: "optional"},
			want: errors.New("unexpected param count"),
		},
		{
			name: "invalid optional param type",
			t:    &TypeDesc{TypeName: "optional", Params: []*TypeDesc{{}}},
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
			lib:       &LibrarySubset{},
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
