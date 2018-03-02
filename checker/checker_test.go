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

package checker

import (
	"fmt"
	"testing"

	"celgo/common"
	"celgo/parser"
	"celgo/semantics"
	"celgo/semantics/types"
	"celgo/test"
)

var testCases = []testInfo{
	// Const types
	{
		I:    `"A"`,
		R:    `"A"~string`,
		Type: types.String,
	},
	{
		I:    `12`,
		R:    `12~int`,
		Type: types.Int64,
	},
	{
		I:    `12u`,
		R:    `12u~uint`,
		Type: types.Uint64,
	},
	{
		I:    `true`,
		R:    `true~bool`,
		Type: types.Bool,
	},
	{
		I:    `false`,
		R:    `false~bool`,
		Type: types.Bool,
	},
	{
		I:    `12.23`,
		R:    `12.23~double`,
		Type: types.Double,
	},
	{
		I:    `null`,
		R:    `null~null`,
		Type: types.Null,
	},
	{
		I:    `b"ABC"`,
		R:    `b"ABC"~bytes`,
		Type: types.Bytes,
	},

	// Ident types
	{
		I:    `is`,
		R:    `is~string^is`,
		Type: types.String,
		Env:  testEnvs["default"],
	},
	{
		I:    `ii`,
		R:    `ii~int^ii`,
		Type: types.Int64,
		Env:  testEnvs["default"],
	},
	{
		I:    `iu`,
		R:    `iu~uint^iu`,
		Type: types.Uint64,
		Env:  testEnvs["default"],
	},
	{
		I:    `iz`,
		R:    `iz~bool^iz`,
		Type: types.Bool,
		Env:  testEnvs["default"],
	},
	{
		I:    `id`,
		R:    `id~double^id`,
		Type: types.Double,
		Env:  testEnvs["default"],
	},
	{
		I:    `ix`,
		R:    `ix~null^ix`,
		Type: types.Null,
		Env:  testEnvs["default"],
	},
	{
		I:    `ib`,
		R:    `ib~bytes^ib`,
		Type: types.Bytes,
		Env:  testEnvs["default"],
	},
	{
		I:    `id`,
		R:    `id~double^id`,
		Type: types.Double,
		Env:  testEnvs["default"],
	},

	{
		I:    `[]`,
		R:    `[]~list(dyn)`,
		Type: types.NewList(types.Dynamic),
	},
	{
		I:    `[1]`,
		R:    `[1~int]~list(int)`,
		Type: types.NewList(types.Int64),
	},

	{
		I:    `[1, "A"]`,
		R:    `[1~int, "A"~string]~list(int)`,
		Type: types.NewList(types.Int64),
		Error: `
ERROR: <input>:1:5: type 'string' does not match previous type 'int' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.
 | [1, "A"]
 | ....^
		`,
	},
	{
		I:    `foo`,
		R:    `foo~!error!`,
		Type: types.Error,
		Error: `
ERROR: <input>:1:1: undeclared reference to 'foo' (in container '')
| foo
| ^`,
	},

	// Call resolution
	{
		I:    `fg_s()`,
		R:    `fg_s()~string^fg_s_0`,
		Type: types.String,
		Env:  testEnvs["default"],
	},
	{
		I:    `is.fi_s_s()`,
		R:    `is~string^is.fi_s_s()~string^fi_s_s_0`,
		Type: types.String,
		Env:  testEnvs["default"],
	},
	{
		I:    `1 + 2`,
		R:    `_+_(1~int, 2~int)~int^add_int64`,
		Type: types.Int64,
		Env:  testEnvs["default"],
	},
	{
		I:    `1 + ii`,
		R:    `_+_(1~int, ii~int^ii)~int^add_int64`,
		Type: types.Int64,
		Env:  testEnvs["default"],
	},
	{
		I:    `[1] + [2]`,
		R:    `_+_([1~int]~list(int), [2~int]~list(int))~list(int)^add_list`,
		Type: types.NewList(types.Int64),
		Env:  testEnvs["default"],
	},

	// Tests from Java implementation
	{
		I:    `[] + [1,2,3,] + [4]`,
		Type: types.NewList(types.Int64),
		R: `
	_+_(
		_+_(
			[]~list(int),
			[1~int, 2~int, 3~int]~list(int))~list(int)^add_list,
			[4~int]~list(int))
	~list(int)^add_list
	`,
	},

	{
		I:    `[1, 2u] + []`,
		Type: types.NewList(types.Int64),
		Error: `
ERROR: <input>:1:5: type 'uint' does not match previous type 'int' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.
  | [1, 2u] + []
  | ....^
		`,
	},

	{
		I:    `{1:2u, 2:3u}`,
		Type: types.NewMap(types.Int64, types.Uint64),
		R:    `{1~int : 2u~uint, 2~int : 3u~uint}~map(int, uint)`,
	},

	{
		I:    `{"a":1, "b":2}.a`,
		Type: types.Int64,
		R:    `{"a"~string : 1~int, "b"~string : 2~int}~map(string, int).a~int`,
	},
	{
		I: `{1:2u, 2u:3}`,
		Error: `
ERROR: <input>:1:8: type 'uint' does not match previous type 'int' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.
 | {1:2u, 2u:3}
 | .......^
ERROR: <input>:1:11: type 'int' does not match previous type 'uint' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.
| {1:2u, 2u:3}
| ..........^
	`,
	},

	{
		I:         `TestAllTypes{single_int32: 1, single_int64: 2}`,
		Container: "google.api.tools.expr.test",
		R: `
	TestAllTypes{single_int32 : 1~int, single_int64 : 2~int}
	  ~google.api.tools.expr.test.TestAllTypes
	    ^google.api.tools.expr.test.TestAllTypes`,
		Type: types.NewMessage("google.api.tools.expr.test.TestAllTypes"),
	},

	{
		I:         `TestAllTypes{single_int32: 1u}`,
		Container: "google.api.tools.expr.test.",
		Error: `
	ERROR: <input>:1:26: expected type of field 'single_int32' is 'int' but provided type is 'uint'
	  | TestAllTypes{single_int32: 1u}
	  | .........................^`,
	},

	{
		I:         `TestAllTypes{single_int32: 1, undefined: 2}`,
		Container: "google.api.tools.expr.test.",
		Error: `
	ERROR: <input>:1:40: undefined field 'undefined'
	  | TestAllTypes{single_int32: 1, undefined: 2}
	  | .......................................^`,
	},

	{
		I: `size(x) == x.size()`,
		R: `
_==_(size(x~list(int)^x)~int^size_list, x~list(int)^x.size()~int^my_size)
  ~bool^equals`,
		Env: env{
			functions: []*semantics.Function{
				semantics.NewFunction("size",
					semantics.NewOverload("my_size", true, types.Int64, types.NewList(types.NewTypeParam("A")))),
			},
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewList(types.Int64), nil),
			},
		},
		Type: types.Bool,
	},
	{
		I: `int(1u) + int(uint("1"))`,
		R: `
_+_(int(1u~uint)~int^uint64_to_int64,
      int(uint("1"~string)~uint^string_to_uint64)~int^uint64_to_int64)
  ~int^add_int64`,
		Type: types.Int64,
	},

	{
		I: `false && !true || false ? 2 : 3`,
		R: `
_?_:_(_||_(_&&_(false~bool, !_(true~bool)~bool^logical_not)~bool^logical_and,
            false~bool)
        ~bool^logical_or,
      2~int,
      3~int)
  ~int^conditional
`,
		Type: types.Int64,
	},

	{
		I:    `b"abc" + b"def"`,
		R:    `_+_(b"abc"~bytes, b"def"~bytes)~bytes^add_bytes`,
		Type: types.Bytes,
	},

	{
		I: `1.0 + 2.0 * 3.0 - 1.0 / 2.20202 != 66.6`,
		R: `
_!=_(_-_(_+_(1~double, _*_(2~double, 3~double)~double^multiply_double)
           ~double^add_double,
           _/_(1~double, 2.20202~double)~double^divide_double)
       ~double^subtract_double,
      66.6~double)
  ~bool^not_equals`,
		Type: types.Bool,
	},

	{
		I:    `1 + 2 * 3 - 1 / 2 == 6 % 1`,
		R:    ` _==_(_-_(_+_(1~int, _*_(2~int, 3~int)~int^multiply_int64)~int^add_int64, _/_(1~int, 2~int)~int^divide_int64)~int^subtract_int64, _%_(6~int, 1~int)~int^modulo_int64)~bool^equals`,
		Type: types.Bool,
	},

	{
		I:    `"abc" + "def"`,
		R:    `_+_("abc"~string, "def"~string)~string^add_string`,
		Type: types.String,
	},

	{
		I: `1u + 2u * 3u - 1u / 2u == 6u % 1u`,
		R: `_==_(_-_(_+_(1u~uint, _*_(2u~uint, 3u~uint)~uint^multiply_uint64)
	         ~uint^add_uint64,
	         _/_(1u~uint, 2u~uint)~uint^divide_uint64)
	     ~uint^subtract_uint64,
	    _%_(6u~uint, 1u~uint)~uint^modulo_uint64)
	~bool^equals`,
		Type: types.Bool,
	},

	{
		I: `x.single_int32 != null`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.Proto2Message"), nil),
			},
		},
		Error: `
	ERROR: <input>:1:2: [internal] unexpected failed resolution of 'google.api.tools.expr.test.Proto2Message'
	  | x.single_int32 != null
	  | .^
	`,
	},

	{
		I: `x.single_value + 1 / x.single_struct == 23`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `_==_(_+_(x~google.api.tools.expr.test.TestAllTypes^x.single_value~dyn,
	         _/_(1~int,
	               x~google.api.tools.expr.test.TestAllTypes^x.single_struct~dyn)
	           ~int^divide_int64)
	     ~int^add_int64,
	    23~int)
	~bool^equals`,
		Type: types.Bool,
	},

	{
		I: `x.single_value[23] + x.single_struct`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `_+_(_[_](x~google.api.tools.expr.test.TestAllTypes^x.single_value~dyn, 23~int)
	    ~dyn^index_list|index_map,
	    x~google.api.tools.expr.test.TestAllTypes^x.single_struct~dyn)
	~dyn^add_int64|add_uint64|add_double|add_string|add_bytes|add_list`,
		Type: types.Dynamic,
	},

	// TODO: We have minimal proto support. This would get fixed with better proto support.
	//{
	//	I: `TestAllTypes.NestedEnum.BAR != 99`,
	//	Container: "google.api.tools.expr.test",
	//	R: `_!=_(TestAllTypes.NestedEnum.BAR
	//     ~int^google.api.tools.expr.test.TestAllTypes.NestedEnum.BAR,
	//    99~int)
	//~bool^not_equals`,
	//},

	{
		I:    `size([] + [1])`,
		R:    `size(_+_([]~list(int), [1~int]~list(int))~list(int)^add_list)~int^size_list`,
		Type: types.Int64,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
	},

	{
		I: `x + y`,
		R: ``,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewList(types.NewMessage("google.api.tools.expr.test.TestAllTypes")), nil),
				semantics.NewIdent("y", types.NewList(types.Int64), nil),
			},
		},
		Error: `
ERROR: <input>:1:3: found no matching overload for '_+_' applied to '(list(google.api.tools.expr.test.TestAllTypes), list(int))'
  | x + y
  | ..^
		`,
	},

	{
		I: `x[1u]`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewList(types.NewMessage("google.api.tools.expr.test.TestAllTypes")), nil),
			},
		},
		Error: `
ERROR: <input>:1:2: found no matching overload for '_[_]' applied to '(list(google.api.tools.expr.test.TestAllTypes), uint)'
  | x[1u]
  | .^
`,
	},

	{
		I: `(x + x)[1].single_int32 == size(x)`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewList(types.NewMessage("google.api.tools.expr.test.TestAllTypes")), nil),
			},
		},
		R: `
_==_(_[_](_+_(x~list(google.api.tools.expr.test.TestAllTypes)^x,
                x~list(google.api.tools.expr.test.TestAllTypes)^x)
            ~list(google.api.tools.expr.test.TestAllTypes)^add_list,
           1~int)
       ~google.api.tools.expr.test.TestAllTypes^index_list
       .
       single_int32
       ~int,
      size(x~list(google.api.tools.expr.test.TestAllTypes)^x)~int^size_list)
  ~bool^equals
	`,
		Type: types.Bool,
	},

	{
		I: `x.repeated_int64[x.single_int32] == 23`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
_==_(_[_](x~google.api.tools.expr.test.TestAllTypes^x.repeated_int64~list(int),
           x~google.api.tools.expr.test.TestAllTypes^x.single_int32~int)
       ~int^index_list,
      23~int)
  ~bool^equals`,
		Type: types.Bool,
	},

	{
		I: `size(x.map_int64_nested_type) == 0`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
_==_(size(x~google.api.tools.expr.test.TestAllTypes^x.map_int64_nested_type
            ~map(int, google.api.tools.expr.test.NestedTestAllTypes))
       ~int^size_map,
      0~int)
  ~bool^equals
		`,
		Type: types.Bool,
	},

	{
		I: `x.repeated_int64.map(x, double(x))`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
		__comprehension__(
    		  // Variable
    		  x,
    		  // Target
    		  x~google.api.tools.expr.test.TestAllTypes^x.repeated_int64~list(int),
    		  // Accumulator
    		  __result__,
    		  // Init
    		  []~list(double),
    		  // LoopCondition
    		  true~bool,
    		  // LoopStep
    		  _+_(
    		    __result__~list(double)^__result__,
    		    [
    		      double(
    		        x~int^x
    		      )~double^int64_to_double
    		    ]~list(double)
    		  )~list(double)^add_list,
    		  // Result
    		  __result__~list(double)^__result__)~list(double)
		`,
		Type: types.NewList(types.Double),
	},

	{
		I: `x.repeated_int64.map(x, x > 0, double(x))`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
	__comprehension__(
    		  // Variable
    		  x,
    		  // Target
    		  x~google.api.tools.expr.test.TestAllTypes^x.repeated_int64~list(int),
    		  // Accumulator
    		  __result__,
    		  // Init
    		  []~list(double),
    		  // LoopCondition
    		  true~bool,
    		  // LoopStep
    		  _?_:_(
    		    _>_(
    		      x~int^x,
    		      0~int
    		    )~bool^greater_int64,
    		    _+_(
    		      __result__~list(double)^__result__,
    		      [
    		        double(
    		          x~int^x
    		        )~double^int64_to_double
    		      ]~list(double)
    		    )~list(double)^add_list,
    		    __result__~list(double)^__result__
    		  )~list(double)^conditional,
    		  // Result
    		  __result__~list(double)^__result__)~list(double)
		`,
		Type: types.NewList(types.Double),
	},

	{
		I: `x[2].single_int32 == 23`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMap(types.String, types.NewMessage("google.api.tools.expr.test.TestAllTypes")), nil),
			},
		},
		Error: `
ERROR: <input>:1:2: found no matching overload for '_[_]' applied to '(map(string, google.api.tools.expr.test.TestAllTypes), int)'
  | x[2].single_int32 == 23
  | .^
		`,
	},

	{
		I: `x["a"].single_int32 == 23`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMap(types.String, types.NewMessage("google.api.tools.expr.test.TestAllTypes")), nil),
			},
		},
		R: `
		_==_(_[_](x~map(string, google.api.tools.expr.test.TestAllTypes)^x, "a"~string)
		~google.api.tools.expr.test.TestAllTypes^index_map
		.
		single_int32
		~int,
		23~int)
		~bool^equals`,
		Type: types.Bool,
	},

	{
		I: `x.single_nested_message.bb == 43 && has(x.single_nested_message)`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},

		// Our implementation code is expanding the macro
		R: `_&&_(
    		  _==_(
    		    x~google.api.tools.expr.test.TestAllTypes^x.single_nested_message~google.api.tools.expr.test.TestAllTypes.NestedMessage.bb~int,
    		    43~int
    		  )~bool^equals,
    		  x~google.api.tools.expr.test.TestAllTypes^x.single_nested_message~test-only~~bool
    		)~bool^logical_and`,
		//	R: `
		//	_&&_(_==_(x~google.api.tools.expr.test.TestAllTypes^x.single_nested_message
		//	~google.api.tools.expr.test.TestAllTypes.NestedMessage
		//	.
		//	bb
		//	~int,
		//	43~int)
		//	~bool^equals,
		//	has(x~google.api.tools.expr.test.TestAllTypes^x.single_nested_message)
		//	~bool)
		//	~bool^logical_and
		//	`,
		Type: types.Bool,
	},

	{
		I: `x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		Error: `
ERROR: <input>:1:24: undefined field 'undefined'
 | x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)
 | .......................^
ERROR: <input>:1:39: undefined field 'undefined'
 | x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)
 | ......................................^
ERROR: <input>:1:56: field 'single_int32' does not support presence check
 | x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)
 | .......................................................^
ERROR: <input>:1:79: field 'repeated_int32' does not support presence check
 | x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)
 | ..............................................................................^
		`,
	},

	{
		I: `x.single_nested_message != null`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
		_!=_(x~google.api.tools.expr.test.TestAllTypes^x.single_nested_message
		~google.api.tools.expr.test.TestAllTypes.NestedMessage,
		null~null)
		~bool^not_equals
		`,
		Type: types.Bool,
	},

	{
		I: `x.single_int64 != null`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		Error: `
ERROR: <input>:1:16: found no matching overload for '_!=_' applied to '(int, null)'
 | x.single_int64 != null
 | ...............^
		`,
	},

	{
		I: `x.single_int64_wrapper == null`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
		_==_(x~google.api.tools.expr.test.TestAllTypes^x.single_int64_wrapper
		~wrapper(int),
		null~null)
		~bool^equals
		`,
		Type: types.Bool,
	},
	//
	//{
	//	I: `x.repeated_int64.all(e, e > 0) && x.repeated_int64.exists(e, e < 0) && x.repeated_int64.exists_one(e, e == 0)`,
	//	Env: env{
	//		idents: []*semantics.Ident{
	//			semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
	//		},
	//	},
	//	R: `
	//	_&&_(_&&_(__comprehension__(e,
	//	x~google.api.tools.expr.test.TestAllTypes^x.repeated_int64
	//	~list(int),
	//	__result__,
	//	true~bool,
	//	__result__~bool^__result__,
	//	_&&_(__result__~bool^__result__,
	//	_>_(e~int^e, 0~int)~bool^greater_int64)
	//	~bool^logical_and,
	//	__result__~bool^__result__)
	//	~bool,
	//	__comprehension__(e,
	//	x~google.api.tools.expr.test.TestAllTypes^x.repeated_int64
	//	~list(int),
	//	__result__,
	//	false~bool,
	//	!_(__result__~bool^__result__)~bool^logical_not,
	//	_||_(__result__~bool^__result__,
	//	_<_(e~int^e, 0~int)~bool^less_int64)
	//	~bool^logical_or,
	//	__result__~bool^__result__)
	//	~bool)
	//	~bool^logical_and,
	//	__comprehension__(e,
	//	x~google.api.tools.expr.test.TestAllTypes^x.repeated_int64
	//	~list(int),
	//	__result__,
	//	0~int,
	//	_<=_(__result__~int^__result__, 1~int)~bool^less_equals_int64,
	//	_?_:_(_==_(e~int^e, 0~int)~bool^equals,
	//	_+_(__result__~int^__result__, 1~int)~int^add_int64,
	//	__result__~int^__result__)
	//	~int^conditional,
	//	_==_(__result__~int^__result__, 1~int)~bool^equals)
	//	~bool)
	//	~bool^logical_and
	//	`,
	//	Type: types.Bool,
	//},

	{
		I: `x.all(e, 0)`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		Error: `
ERROR: <input>:1:1: expression of type 'google.api.tools.expr.test.TestAllTypes' cannot be range of a comprehension (must be list, map, or dynamic)
 | x.all(e, 0)
 | ^
ERROR: <input>:1:2: found no matching overload for '_&&_' applied to '(bool, int)'
 | x.all(e, 0)
 | .^
		`,
	},

	{
		I: `.google.api.tools.expr.test.TestAllTypes`,
		R: `	.google.api.tools.expr.test.TestAllTypes
	~type(google.api.tools.expr.test.TestAllTypes)
	^google.api.tools.expr.test.TestAllTypes`,
		Type: types.NewTypeType(types.NewMessage("google.api.tools.expr.test.TestAllTypes")),
	},

	{
		I:         `test.TestAllTypes`,
		Container: "google.api.tools.expr",
		R: `
	test.TestAllTypes
	~type(google.api.tools.expr.test.TestAllTypes)
	^google.api.tools.expr.test.TestAllTypes
		`,
		Type: types.NewTypeType(types.NewMessage("google.api.tools.expr.test.TestAllTypes")),
	},

	{
		I: `1 + x`,
		Error: `
ERROR: <input>:1:5: undeclared reference to 'x' (in container '')
 | 1 + x
 | ....^`,
	},

	{
		I:         `x`,
		Container: "container",
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("container.x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R:    `x~google.api.tools.expr.test.TestAllTypes^container.x`,
		Type: types.NewMessage("google.api.tools.expr.test.TestAllTypes"),
	},

	{
		I: `list(int) == type([1]) && map(int, uint) == type({1:2u})`,
		R: `
_&&_(_==_(list(int~type(int)^int)~type(list(int))^list_type,
           type([1~int]~list(int))~type(list(int))^type)
       ~bool^equals,
      _==_(map(int~type(int)^int, uint~type(uint)^uint)
             ~type(map(int, uint))^map_type,
            type({1~int : 2u~uint}~map(int, uint))~type(map(int, uint))^type)
        ~bool^equals)
  ~bool^logical_and
	`,
		Type: types.Bool,
	},

	{
		I: `myfun(1, true, 3u) + 1.myfun(false, 3u).myfun(true, 42u)`,
		Env: env{
			functions: []*semantics.Function{
				semantics.NewFunction("myfun",
					semantics.NewOverload("myfun_instance", true, types.Int64, types.Int64, types.Bool, types.Uint64)),
				semantics.NewFunction("myfun",
					semantics.NewOverload("myfun_static", false, types.Int64, types.Int64, types.Bool, types.Uint64)),
			},
		},
		// Original "expected" is pretty broken. I think our code is doing-the-right-thing(TM).
		R: `_+_(
    		  myfun(
    		    1~int,
    		    true~bool,
    		    3u~uint
    		  )~int^myfun_static,
    		  1~int.myfun(
    		    false~bool,
    		    3u~uint
    		  )~int^myfun_instance.myfun(
    		    true~bool,
    		    42u~uint
    		  )~int^myfun_instance
    		)~int^add_int64`,
		//R: `
		//_+_(myfun(1~int, true~bool, 3u~uint)~int^myfun_static,
		//42u~uint.
		//myfun(3u~uint.myfun(1~int, false~bool)~int^myfun_instance, true~bool)
		//~int^myfun_instance)
		//~int^add_int64`,
		Type: types.Int64,
	},

	// TODO: This doesn't seem like the right path for this kind of error. Env should catch the case and throw
	// during env setup
	//	{
	//		I: `false`,
	//		Env: env{
	//			functions: []*semantics.Function{
	//				semantics.NewFunction("has",
	//					semantics.NewOverload("has_id", false, types.Dynamic, types.Dynamic)),
	//			},
	//		},
	//		Error: `
	//ERROR: <input>:1:1: overload for name 'has' with 1 argument(s) overlaps with predefined macro
	// | false
	// | ^
	//		`,
	//	},

	{
		I: `size(x) > 4`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
			functions: []*semantics.Function{
				semantics.NewFunction("size",
					semantics.NewOverload("size_message", false, types.Int64, types.NewMessage("google.api.tools.expr.test.TestAllTypes"))),
			},
		},
		Type: types.Bool,
	},

	{
		I: `x.single_int64_wrapper + 1 != 23`,
		Env: env{
			idents: []*semantics.Ident{
				semantics.NewIdent("x", types.NewMessage("google.api.tools.expr.test.TestAllTypes"), nil),
			},
		},
		R: `
		_!=_(_+_(x~google.api.tools.expr.test.TestAllTypes^x.single_int64_wrapper
		~wrapper(int),
		1~int)
		~int^add_int64,
		23~int)
		~bool^not_equals
		`,
		Type: types.Bool,
	},
}

var typeProvider = initTypeProvider()

func initTypeProvider() TypeProvider {
	provider := NewInMemoryTypeProvider()

	nestedAllTypes := types.NewMessage("google.api.tools.expr.test.NestedTestAllTypes")

	nestedMessageType := types.NewMessage("google.api.tools.expr.test.TestAllTypes.NestedMessage")
	provider.AddType(nestedMessageType.Name(), nestedMessageType)
	provider.AddFieldType(nestedMessageType, "bb", types.Int64, false)

	allTypes := types.NewMessage("google.api.tools.expr.test.TestAllTypes")
	provider.AddType(allTypes.Name(), allTypes)

	provider.AddFieldType(allTypes, "single_struct", types.Dynamic, false)
	provider.AddFieldType(allTypes, "single_value", types.Dynamic, false)
	provider.AddFieldType(allTypes, "single_int32", types.Int64, false)
	provider.AddFieldType(allTypes, "single_int64", types.Int64, false)
	provider.AddFieldType(allTypes, "repeated_int32", types.NewList(types.Int64), false)
	provider.AddFieldType(allTypes, "repeated_int64", types.NewList(types.Int64), false)
	provider.AddFieldType(allTypes, "single_int64_wrapper", types.NewWrapper(types.Int64), false)
	provider.AddFieldType(allTypes, "map_int64_nested_type",
		types.NewMap(
			types.Int64,
			nestedAllTypes), true)

	provider.AddFieldType(allTypes, "single_nested_message", nestedMessageType, true)

	provider.AddEnum("google.api.tools.expr.test.TestAllTypes.NestedEnum.FOO", 0)
	provider.AddEnum("google.api.tools.expr.test.TestAllTypes.NestedEnum.BAR", 1)
	provider.AddEnum("google.api.tools.expr.test.TestAllTypes.NestedEnum.BAZ", 2)

	//enumType := types.NewMessage("google.api.tools.expr.test.TestAllTypes.NestedEnum")
	//provider.AddType(enumType.Name(), enumType)

	return provider
}

var testEnvs = map[string]env{
	"default": env{
		functions: []*semantics.Function{
			semantics.NewFunction("fg_s",
				semantics.NewOverload("fg_s_0", false, types.String)),
			semantics.NewFunction("fi_s_s",
				semantics.NewOverload("fi_s_s_0", true, types.String, types.String)),
		},
		idents: []*semantics.Ident{
			semantics.NewIdent("is", types.String, nil),
			semantics.NewIdent("ii", types.Int64, nil),
			semantics.NewIdent("iu", types.Uint64, nil),
			semantics.NewIdent("iz", types.Bool, nil),
			semantics.NewIdent("ib", types.Bytes, nil),
			semantics.NewIdent("id", types.Double, nil),
			semantics.NewIdent("ix", types.Null, nil),
		},
	},
}

type testInfo struct {
	// I contains the input expression to be parsed.
	I string

	// R contains the result output.
	R string

	// Type is the expected type of the expression
	Type types.Type

	// Container is the container name to use for test.
	Container string

	// Env is the environment to use for testing.
	Env env

	// Error is the expected error for negative test cases.
	Error string
}

type env struct {
	idents    []*semantics.Ident
	functions []*semantics.Function
}

func Test(t *testing.T) {
	for i, tst := range testCases {
		name := fmt.Sprintf("%d %s", i, tst.I)
		t.Run(name, func(tt *testing.T) {

			expression, errors := parser.ParseText(tst.I)
			if len(errors.GetErrors()) > 0 {
				tt.Fatalf("Unexpected parse errors: %v", errors.GetErrors()[0].ToDisplayString())
				return
			}

			errors = common.NewErrors()
			env := NewEnv(errors, typeProvider)
			AddStandard(env)

			if tst.Env.idents != nil {
				for _, ident := range tst.Env.idents {
					env.AddIdent(ident)
				}
			}
			if tst.Env.functions != nil {
				for _, fn := range tst.Env.functions {
					env.AddFunction(fn)
				}
			}

			semantics := Check(env, tst.Container, expression)
			if len(errors.GetErrors()) > 0 {
				errorString := errors.String()
				if tst.Error != "" {
					if !test.Compare(errorString, tst.Error) {
						tt.Error(test.DiffMessage("Error mismatch", errorString, tst.Error))
					}
				} else {
					tt.Errorf("Unexpected type-check errors: %v", errorString)
				}
			} else if tst.Error != "" {
				tt.Errorf("Expected error not thrown: %s", tst.Error)
			}

			actual := semantics.GetType(expression)

			if tst.Error == "" {
				if actual == nil || !actual.Equals(tst.Type) {
					tt.Error(test.DiffMessage("Type Error", actual, tst.Type))
				}
			}

			if tst.R != "" {
				actualStr := print(expression, semantics)
				if !test.Compare(actualStr, tst.R) {
					tt.Error(test.DiffMessage("Structure error", actualStr, tst.R))
				}
			}
		})
	}
}
