// Copyright 2019 Google LLC
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

package parser

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/google/cel-go/common"
)

func TestUnparse_Identical(t *testing.T) {
	tests := map[string]string{
		"call_add":            `a + b - c`,
		"call_and":            `a && b && c && d && e`,
		"call_and_or":         `a || b && (c || d) && e`,
		"call_cond":           `a ? b : c`,
		"call_index":          `a[1]["b"]`,
		"call_mul":            `a * (b / c) % 0`,
		"call_mul_add":        `a + b * c`,
		"call_mul_add_nested": `(a + b) * c / (d - e)`,
		"call_mul_nested":     `a * b / c % 0`,
		"call_not":            `!true`,
		"call_neg":            `-num`,
		"call_or":             `a || b || c || d || e`,
		"func_global":         `size(a ? (b ? c : d) : e)`,
		"func_member":         `a.hello("world")`,
		"func_no_arg":         `zero()`,
		"func_one_arg":        `one("a")`,
		"func_two_args":       `and(d, 32u)`,
		"func_var_args":       `max(a, b, 100)`,
		"list_empty":          `[]`,
		"list_one":            `[1]`,
		"list_many":           `["hello, world", "goodbye, world", "sure, why not?"]`,
		"lit_bytes":           `b"Ã¿"`,
		"lit_double":          `-42.101`,
		"lit_false":           `false`,
		"lit_int":             `-405069`,
		"lit_null":            `null`,
		"lit_string":          `"hello:\t'world'"`,
		"lit_true":            `true`,
		"lit_uint":            `42u`,
		"ident":               `my_ident`,
		"macro_has":           `has(hello.world)`,
		"map_empty":           `{}`,
		"map_lit_key":         `{"a": a.b.c, b"b": bytes(a.b.c)}`,
		"map_expr_key": `
            {a: a,
             b: a.b,
             c: a.b.c,
             a ? b : c: false,
             a || b: true}`,
		"msg_empty":  `v1alpha1.Expr{}`,
		"msg_fields": `v1alpha1.Expr{id: 1, call_expr: v1alpha1.Call_Expr{function: "name"}}`,
		"select":     `a.b.c`,
	}

	for name, in := range tests {
		t.Run(name, func(tt *testing.T) {
			p, iss := Parse(common.NewTextSource(in))
			if len(iss.GetErrors()) > 0 {
				tt.Fatal(iss.ToDisplayString())
			}
			out, err := Unparse(p.GetExpr(), p.GetSourceInfo())
			if err != nil {
				tt.Error(err)
			}
			if out != in {
				tt.Errorf("Got '%s', wanted '%s'", out, in)
			}
			p2, iss := Parse(common.NewTextSource(out))
			before := p.GetExpr()
			after := p2.GetExpr()
			if !proto.Equal(before, after) {
				tt.Errorf("Second parse differs from the first. Got '%v', wanted '%v'",
					before, after)
			}
		})
	}
}

func TestUnparse_Equivalent(t *testing.T) {
	tests := map[string][]string{
		"call_add":   {`a+b-c`, `a + b - c`},
		"call_cond":  {`a ? b          : c`, `a ? b :          c`},
		"call_index": {`a[  1  ]["b"]`, `a[  1]  ["b"]`},
		"select":     {`a . b . c`, `a .b  .c`},
	}

	for name, in := range tests {
		t.Run(name, func(tt *testing.T) {
			p, iss := Parse(common.NewTextSource(in[0]))
			if len(iss.GetErrors()) > 0 {
				tt.Fatal(iss.ToDisplayString())
			}
			out, err := Unparse(p.GetExpr(), p.GetSourceInfo())
			if err != nil {
				tt.Error(err)
			}
			if out != in[1] {
				tt.Errorf("Got '%s', wanted '%s'", out, in[1])
			}
			p2, iss := Parse(common.NewTextSource(out))
			before := p.GetExpr()
			after := p2.GetExpr()
			if !proto.Equal(before, after) {
				tt.Errorf("Second parse differs from the first. Got '%v', wanted '%v'",
					before, after)
			}
		})
	}
}
