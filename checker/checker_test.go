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

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/parser"
	"github.com/google/cel-go/test"

	"google.golang.org/protobuf/proto"

	proto2pb "github.com/google/cel-go/test/proto2pb"
	proto3pb "github.com/google/cel-go/test/proto3pb"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var testCases = []testInfo{
	// Const types
	{
		in:      `"A"`,
		out:     `"A"~string`,
		outType: decls.String,
	},
	{
		in:      `12`,
		out:     `12~int`,
		outType: decls.Int,
	},
	{
		in:      `12u`,
		out:     `12u~uint`,
		outType: decls.Uint,
	},
	{
		in:      `true`,
		out:     `true~bool`,
		outType: decls.Bool,
	},
	{
		in:      `false`,
		out:     `false~bool`,
		outType: decls.Bool,
	},
	{
		in:      `12.23`,
		out:     `12.23~double`,
		outType: decls.Double,
	},
	{
		in:      `null`,
		out:     `null~null`,
		outType: decls.Null,
	},
	{
		in:      `b"ABC"`,
		out:     `b"ABC"~bytes`,
		outType: decls.Bytes,
	},
	// Ident types
	{
		in:      `is`,
		out:     `is~string^is`,
		outType: decls.String,
		env:     testEnvs["default"],
	},
	{
		in:      `ii`,
		out:     `ii~int^ii`,
		outType: decls.Int,
		env:     testEnvs["default"],
	},
	{
		in:      `iu`,
		out:     `iu~uint^iu`,
		outType: decls.Uint,
		env:     testEnvs["default"],
	},
	{
		in:      `iz`,
		out:     `iz~bool^iz`,
		outType: decls.Bool,
		env:     testEnvs["default"],
	},
	{
		in:      `id`,
		out:     `id~double^id`,
		outType: decls.Double,
		env:     testEnvs["default"],
	},
	{
		in:      `ix`,
		out:     `ix~null^ix`,
		outType: decls.Null,
		env:     testEnvs["default"],
	},
	{
		in:      `ib`,
		out:     `ib~bytes^ib`,
		outType: decls.Bytes,
		env:     testEnvs["default"],
	},
	{
		in:      `id`,
		out:     `id~double^id`,
		outType: decls.Double,
		env:     testEnvs["default"],
	},
	{
		in:      `[]`,
		out:     `[]~list(dyn)`,
		outType: decls.NewListType(decls.Dyn),
	},
	{
		in:      `[1]`,
		out:     `[1~int]~list(int)`,
		outType: decls.NewListType(decls.Int),
	},
	{
		in:      `[1, "A"]`,
		out:     `[1~int, "A"~string]~list(dyn)`,
		outType: decls.NewListType(decls.Dyn),
	},
	{
		in:      `foo`,
		out:     `foo~!error!`,
		outType: decls.Error,
		err: `
ERROR: <input>:1:1: undeclared reference to 'foo' (in container '')
| foo
| ^`,
	},
	// Call resolution
	{
		in:      `fg_s()`,
		out:     `fg_s()~string^fg_s_0`,
		outType: decls.String,
		env:     testEnvs["default"],
	},
	{
		in:      `is.fi_s_s()`,
		out:     `is~string^is.fi_s_s()~string^fi_s_s_0`,
		outType: decls.String,
		env:     testEnvs["default"],
	},
	{
		in:      `1 + 2`,
		out:     `_+_(1~int, 2~int)~int^add_int64`,
		outType: decls.Int,
		env:     testEnvs["default"],
	},
	{
		in:      `1 + ii`,
		out:     `_+_(1~int, ii~int^ii)~int^add_int64`,
		outType: decls.Int,
		env:     testEnvs["default"],
	},
	{
		in:      `[1] + [2]`,
		out:     `_+_([1~int]~list(int), [2~int]~list(int))~list(int)^add_list`,
		outType: decls.NewListType(decls.Int),
		env:     testEnvs["default"],
	},
	{
		in:      `[] + [1,2,3,] + [4]`,
		outType: decls.NewListType(decls.Int),
		out: `
	_+_(
		_+_(
			[]~list(int),
			[1~int, 2~int, 3~int]~list(int))~list(int)^add_list,
			[4~int]~list(int))
	~list(int)^add_list
	`,
	},
	{
		in: `[1, 2u] + []`,
		out: `_+_(
			[
				1~int,
				2u~uint
			]~list(dyn),
			[]~list(dyn)
		)~list(dyn)^add_list`,
		outType: decls.NewListType(decls.Dyn),
	},
	{
		in:      `{1:2u, 2:3u}`,
		outType: decls.NewMapType(decls.Int, decls.Uint),
		out:     `{1~int : 2u~uint, 2~int : 3u~uint}~map(int, uint)`,
	},
	{
		in:      `{"a":1, "b":2}.a`,
		outType: decls.Int,
		out:     `{"a"~string : 1~int, "b"~string : 2~int}~map(string, int).a~int`,
	},
	{
		in:      `{1:2u, 2u:3}`,
		outType: decls.NewMapType(decls.Dyn, decls.Dyn),
		out:     `{1~int : 2u~uint, 2u~uint : 3~int}~map(dyn, dyn)`,
	},
	{
		in:        `TestAllTypes{single_int32: 1, single_int64: 2}`,
		container: "google.expr.proto3.test",
		out: `
		google.expr.proto3.test.TestAllTypes{
			single_int32 : 1~int,
			single_int64 : 2~int
		}~google.expr.proto3.test.TestAllTypes^google.expr.proto3.test.TestAllTypes`,
		outType: decls.NewObjectType("google.expr.proto3.test.TestAllTypes"),
	},
	{
		in:        `TestAllTypes{single_int32: 1u}`,
		container: "google.expr.proto3.test",
		err: `
	ERROR: <input>:1:26: expected type of field 'single_int32' is 'int' but provided type is 'uint'
	  | TestAllTypes{single_int32: 1u}
	  | .........................^`,
	},
	{
		in:        `TestAllTypes{single_int32: 1, undefined: 2}`,
		container: "google.expr.proto3.test",
		err: `
	ERROR: <input>:1:40: undefined field 'undefined'
	  | TestAllTypes{single_int32: 1, undefined: 2}
	  | .......................................^`,
	},
	{
		in: `size(x) == x.size()`,
		out: `
_==_(size(x~list(int)^x)~int^size_list, x~list(int)^x.size()~int^list_size)
  ~bool^equals`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewListType(decls.Int)),
			},
		},
		outType: decls.Bool,
	},
	{
		in: `int(1u) + int(uint("1"))`,
		out: `
_+_(int(1u~uint)~int^uint64_to_int64,
      int(uint("1"~string)~uint^string_to_uint64)~int^uint64_to_int64)
  ~int^add_int64`,
		outType: decls.Int,
	},
	{
		in: `false && !true || false ? 2 : 3`,
		out: `
_?_:_(_||_(_&&_(false~bool, !_(true~bool)~bool^logical_not)~bool^logical_and,
            false~bool)
        ~bool^logical_or,
      2~int,
      3~int)
  ~int^conditional
`,
		outType: decls.Int,
	},
	{
		in:      `b"abc" + b"def"`,
		out:     `_+_(b"abc"~bytes, b"def"~bytes)~bytes^add_bytes`,
		outType: decls.Bytes,
	},
	{
		in: `1.0 + 2.0 * 3.0 - 1.0 / 2.20202 != 66.6`,
		out: `
_!=_(_-_(_+_(1~double, _*_(2~double, 3~double)~double^multiply_double)
           ~double^add_double,
           _/_(1~double, 2.20202~double)~double^divide_double)
       ~double^subtract_double,
      66.6~double)
  ~bool^not_equals`,
		outType: decls.Bool,
	},
	{
		in: `null == null && null != null`,
		out: `
		_&&_(
			_==_(
				null~null,
				null~null
			)~bool^equals,
			_!=_(
				null~null,
				null~null
			)~bool^not_equals
		)~bool^logical_and`,
		outType: decls.Bool,
	},
	{
		in: `1 == 1 && 2 != 1`,
		out: `
		_&&_(
			_==_(
				1~int,
				1~int
			)~bool^equals,
			_!=_(
				2~int,
				1~int
			)~bool^not_equals
		)~bool^logical_and`,
		outType: decls.Bool,
	},
	{
		in:      `1 + 2 * 3 - 1 / 2 == 6 % 1`,
		out:     ` _==_(_-_(_+_(1~int, _*_(2~int, 3~int)~int^multiply_int64)~int^add_int64, _/_(1~int, 2~int)~int^divide_int64)~int^subtract_int64, _%_(6~int, 1~int)~int^modulo_int64)~bool^equals`,
		outType: decls.Bool,
	},
	{
		in:      `"abc" + "def"`,
		out:     `_+_("abc"~string, "def"~string)~string^add_string`,
		outType: decls.String,
	},
	{
		in: `1u + 2u * 3u - 1u / 2u == 6u % 1u`,
		out: `_==_(_-_(_+_(1u~uint, _*_(2u~uint, 3u~uint)~uint^multiply_uint64)
	         ~uint^add_uint64,
	         _/_(1u~uint, 2u~uint)~uint^divide_uint64)
	     ~uint^subtract_uint64,
	    _%_(6u~uint, 1u~uint)~uint^modulo_uint64)
	~bool^equals`,
		outType: decls.Bool,
	},
	{
		in: `x.single_int32 != null`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.Proto2Message")),
			},
		},
		err: `
	ERROR: <input>:1:2: [internal] unexpected failed resolution of 'google.expr.proto3.test.Proto2Message'
	  | x.single_int32 != null
	  | .^
	`,
	},
	{
		in: `x.single_value + 1 / x.single_struct.y == 23`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `_==_(
			_+_(
			  x~google.expr.proto3.test.TestAllTypes^x.single_value~dyn,
			  _/_(
				1~int,
				x~google.expr.proto3.test.TestAllTypes^x.single_struct~map(string, dyn).y~dyn
			  )~int^divide_int64
			)~int^add_int64,
			23~int
		  )~bool^equals`,
		outType: decls.Bool,
	},
	{
		in: `x.single_value[23] + x.single_struct['y']`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `_+_(
			_[_](
			  x~google.expr.proto3.test.TestAllTypes^x.single_value~dyn,
			  23~int
			)~dyn^index_list|index_map,
			_[_](
			  x~google.expr.proto3.test.TestAllTypes^x.single_struct~map(string, dyn),
			  "y"~string
			)~dyn^index_map
		  )~dyn^add_int64|add_uint64|add_double|add_string|add_bytes|add_list|add_timestamp_duration|add_duration_timestamp|add_duration_duration
		  `,
		outType: decls.Dyn,
	},
	{
		in:        `TestAllTypes.NestedEnum.BAR != 99`,
		container: "google.expr.proto3.test",
		out: `_!=_(google.expr.proto3.test.TestAllTypes.NestedEnum.BAR
	     ~int^google.expr.proto3.test.TestAllTypes.NestedEnum.BAR,
	    99~int)
	~bool^not_equals`,
		outType: decls.Bool,
	},
	{
		in:      `size([] + [1])`,
		out:     `size(_+_([]~list(int), [1~int]~list(int))~list(int)^add_list)~int^size_list`,
		outType: decls.Int,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
	},
	{
		in: `x["claims"]["groups"][0].name == "dummy"
		&& x.claims["exp"] == y[1].time
		&& x.claims.structured == {'key': z}
		&& z == 1.0`,
		out: `_&&_(
			_&&_(
				_==_(
					_[_](
						_[_](
							_[_](
								x~map(string, dyn)^x,
								"claims"~string
							)~dyn^index_map,
							"groups"~string
						)~list(dyn)^index_map,
						0~int
					)~dyn^index_list.name~dyn,
					"dummy"~string
				)~bool^equals,
				_==_(
					_[_](
						x~map(string, dyn)^x.claims~dyn,
						"exp"~string
					)~dyn^index_map,
					_[_](
						y~list(dyn)^y,
						1~int
					)~dyn^index_list.time~dyn
				)~bool^equals
			)~bool^logical_and,
			_&&_(
				_==_(
					x~map(string, dyn)^x.claims~dyn.structured~dyn,
					{
						"key"~string:z~dyn^z
					}~map(string, dyn)
				)~bool^equals,
				_==_(
					z~dyn^z,
					1~double
				)~bool^equals
			)~bool^logical_and
		)~bool^logical_and`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.protobuf.Struct")),
				decls.NewVar("y", decls.NewObjectType("google.protobuf.ListValue")),
				decls.NewVar("z", decls.NewObjectType("google.protobuf.Value")),
			},
		},
		outType: decls.Bool,
	},
	{
		in:  `x + y`,
		out: ``,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewListType(decls.NewObjectType("google.expr.proto3.test.TestAllTypes"))),
				decls.NewVar("y", decls.NewListType(decls.Int)),
			},
		},
		err: `
ERROR: <input>:1:3: found no matching overload for '_+_' applied to '(list(google.expr.proto3.test.TestAllTypes), list(int))'
  | x + y
  | ..^
		`,
	},
	{
		in: `x[1u]`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewListType(decls.NewObjectType("google.expr.proto3.test.TestAllTypes"))),
			},
		},
		err: `
ERROR: <input>:1:2: found no matching overload for '_[_]' applied to '(list(google.expr.proto3.test.TestAllTypes), uint)'
  | x[1u]
  | .^
`,
	},
	{
		in: `(x + x)[1].single_int32 == size(x)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewListType(decls.NewObjectType("google.expr.proto3.test.TestAllTypes"))),
			},
		},
		out: `
_==_(_[_](_+_(x~list(google.expr.proto3.test.TestAllTypes)^x,
                x~list(google.expr.proto3.test.TestAllTypes)^x)
            ~list(google.expr.proto3.test.TestAllTypes)^add_list,
           1~int)
       ~google.expr.proto3.test.TestAllTypes^index_list
       .
       single_int32
       ~int,
      size(x~list(google.expr.proto3.test.TestAllTypes)^x)~int^size_list)
  ~bool^equals
	`,
		outType: decls.Bool,
	},
	{
		in: `x.repeated_int64[x.single_int32] == 23`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
_==_(_[_](x~google.expr.proto3.test.TestAllTypes^x.repeated_int64~list(int),
           x~google.expr.proto3.test.TestAllTypes^x.single_int32~int)
       ~int^index_list,
      23~int)
  ~bool^equals`,
		outType: decls.Bool,
	},
	{
		in: `size(x.map_int64_nested_type) == 0`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
_==_(size(x~google.expr.proto3.test.TestAllTypes^x.map_int64_nested_type
            ~map(int, google.expr.proto3.test.NestedTestAllTypes))
       ~int^size_map,
      0~int)
  ~bool^equals
		`,
		outType: decls.Bool,
	},
	{
		in: `x.all(y, y == true)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.Bool),
			},
		},
		out: `
		__comprehension__(
		// Variable
		y,
		// Target
		x~bool^x,
		// Accumulator
		__result__,
		// Init
		true~bool,
		// LoopCondition
		@not_strictly_false(
			__result__~bool^__result__
		)~bool^not_strictly_false,
		// LoopStep
		_&&_(
			__result__~bool^__result__,
			_==_(
			y~!error!^y,
			true~bool
			)~bool^equals
		)~bool^logical_and,
		// Result
		__result__~bool^__result__)~bool
		`,
		err: `ERROR: <input>:1:1: expression of type 'bool' cannot be range of a comprehension (must be list, map, or dynamic)
		| x.all(y, y == true)
		| ^`,
	},
	{
		in: `x.repeated_int64.map(x, double(x))`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
		__comprehension__(
    		  // Variable
    		  x,
    		  // Target
    		  x~google.expr.proto3.test.TestAllTypes^x.repeated_int64~list(int),
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
		outType: decls.NewListType(decls.Double),
	},
	{
		in: `x.repeated_int64.map(x, x > 0, double(x))`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
	__comprehension__(
    		  // Variable
    		  x,
    		  // Target
    		  x~google.expr.proto3.test.TestAllTypes^x.repeated_int64~list(int),
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
		outType: decls.NewListType(decls.Double),
	},
	{
		in: `x[2].single_int32 == 23`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x",
					decls.NewMapType(decls.String,
						decls.NewObjectType("google.expr.proto3.test.TestAllTypes"))),
			},
		},
		err: `
ERROR: <input>:1:2: found no matching overload for '_[_]' applied to '(map(string, google.expr.proto3.test.TestAllTypes), int)'
  | x[2].single_int32 == 23
  | .^
		`,
	},
	{
		in: `x["a"].single_int32 == 23`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x",
					decls.NewMapType(decls.String,
						decls.NewObjectType("google.expr.proto3.test.TestAllTypes"))),
			},
		},
		out: `
		_==_(_[_](x~map(string, google.expr.proto3.test.TestAllTypes)^x, "a"~string)
		~google.expr.proto3.test.TestAllTypes^index_map
		.
		single_int32
		~int,
		23~int)
		~bool^equals`,
		outType: decls.Bool,
	},
	{
		in: `x.single_nested_message.bb == 43 && has(x.single_nested_message)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},

		// Our implementation code is expanding the macro
		out: `_&&_(
    		  _==_(
    		    x~google.expr.proto3.test.TestAllTypes^x.single_nested_message~google.expr.proto3.test.TestAllTypes.NestedMessage.bb~int,
    		    43~int
    		  )~bool^equals,
    		  x~google.expr.proto3.test.TestAllTypes^x.single_nested_message~test-only~~bool
    		)~bool^logical_and`,
		outType: decls.Bool,
	},
	{
		in: `x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		err: `
ERROR: <input>:1:24: undefined field 'undefined'
| x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)
| .......................^
ERROR: <input>:1:39: undefined field 'undefined'
| x.single_nested_message.undefined == x.undefined && has(x.single_int32) && has(x.repeated_int32)
| ......................................^`,
	},
	{
		in: `x.single_nested_message != null`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
		_!=_(x~google.expr.proto3.test.TestAllTypes^x.single_nested_message
		~google.expr.proto3.test.TestAllTypes.NestedMessage,
		null~null)
		~bool^not_equals
		`,
		outType: decls.Bool,
	},
	{
		in: `x.single_int64 != null`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		err: `
ERROR: <input>:1:16: found no matching overload for '_!=_' applied to '(int, null)'
 | x.single_int64 != null
 | ...............^
		`,
	},
	{
		in: `x.single_int64_wrapper == null`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
		_==_(x~google.expr.proto3.test.TestAllTypes^x.single_int64_wrapper
		~wrapper(int),
		null~null)
		~bool^equals
		`,
		outType: decls.Bool,
	},
	{
		in: `x.single_bool_wrapper
		&& x.single_bytes_wrapper == b'hi'
		&& x.single_double_wrapper != 2.0
		&& x.single_float_wrapper == 1.0
		&& x.single_int32_wrapper != 2
		&& x.single_int64_wrapper == 1
		&& x.single_string_wrapper == 'hi'
		&& x.single_uint32_wrapper == 1u
		&& x.single_uint64_wrapper != 42u`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
		_&&_(
			_&&_(
				_&&_(
				_&&_(
					x~google.expr.proto3.test.TestAllTypes^x.single_bool_wrapper~wrapper(bool),
					_==_(
					x~google.expr.proto3.test.TestAllTypes^x.single_bytes_wrapper~wrapper(bytes),
					b"hi"~bytes
					)~bool^equals
				)~bool^logical_and,
				_!=_(
					x~google.expr.proto3.test.TestAllTypes^x.single_double_wrapper~wrapper(double),
					2~double
				)~bool^not_equals
				)~bool^logical_and,
				_&&_(
				_==_(
					x~google.expr.proto3.test.TestAllTypes^x.single_float_wrapper~wrapper(double),
					1~double
				)~bool^equals,
				_!=_(
					x~google.expr.proto3.test.TestAllTypes^x.single_int32_wrapper~wrapper(int),
					2~int
				)~bool^not_equals
				)~bool^logical_and
			)~bool^logical_and,
			_&&_(
				_&&_(
				_==_(
					x~google.expr.proto3.test.TestAllTypes^x.single_int64_wrapper~wrapper(int),
					1~int
				)~bool^equals,
				_==_(
					x~google.expr.proto3.test.TestAllTypes^x.single_string_wrapper~wrapper(string),
					"hi"~string
				)~bool^equals
				)~bool^logical_and,
				_&&_(
				_==_(
					x~google.expr.proto3.test.TestAllTypes^x.single_uint32_wrapper~wrapper(uint),
					1u~uint
				)~bool^equals,
				_!=_(
					x~google.expr.proto3.test.TestAllTypes^x.single_uint64_wrapper~wrapper(uint),
					42u~uint
				)~bool^not_equals
				)~bool^logical_and
			)~bool^logical_and
		)~bool^logical_and`,
		outType: decls.Bool,
	},
	{
		in: `x.single_bool_wrapper == google.protobuf.BoolValue{value: true}
			&& x.single_bytes_wrapper == google.protobuf.BytesValue{value: b'hi'}
			&& x.single_double_wrapper != google.protobuf.DoubleValue{value: 2.0}
			&& x.single_float_wrapper == google.protobuf.FloatValue{value: 1.0}
			&& x.single_int32_wrapper != google.protobuf.Int32Value{value: -2}
			&& x.single_int64_wrapper == google.protobuf.Int64Value{value: 1}
			&& x.single_string_wrapper == google.protobuf.StringValue{value: 'hi'}
			&& x.single_string_wrapper == google.protobuf.Value{string_value: 'hi'}
			&& x.single_uint32_wrapper == google.protobuf.UInt32Value{value: 1u}
			&& x.single_uint64_wrapper != google.protobuf.UInt64Value{value: 42u}`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		outType: decls.Bool,
	},
	{
		in: `x.repeated_int64.exists(y, y > 10) && y < 5`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		err: `ERROR: <input>:1:39: undeclared reference to 'y' (in container '')
		| x.repeated_int64.exists(y, y > 10) && y < 5
		| ......................................^`,
	},
	{
		in: `x.repeated_int64.all(e, e > 0) && x.repeated_int64.exists(e, e < 0) && x.repeated_int64.exists_one(e, e == 0)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `_&&_(
			_&&_(
			  __comprehension__(
				// Variable
				e,
				// Target
				x~google.expr.proto3.test.TestAllTypes^x.repeated_int64~list(int),
				// Accumulator
				__result__,
				// Init
				true~bool,
				// LoopCondition
				@not_strictly_false(
				  __result__~bool^__result__
				)~bool^not_strictly_false,
				// LoopStep
				_&&_(
				  __result__~bool^__result__,
				  _>_(
					e~int^e,
					0~int
				  )~bool^greater_int64
				)~bool^logical_and,
				// Result
				__result__~bool^__result__)~bool,
			  __comprehension__(
				// Variable
				e,
				// Target
				x~google.expr.proto3.test.TestAllTypes^x.repeated_int64~list(int),
				// Accumulator
				__result__,
				// Init
				false~bool,
				// LoopCondition
				@not_strictly_false(
				  !_(
					__result__~bool^__result__
				  )~bool^logical_not
				)~bool^not_strictly_false,
				// LoopStep
				_||_(
				  __result__~bool^__result__,
				  _<_(
					e~int^e,
					0~int
				  )~bool^less_int64
				)~bool^logical_or,
				// Result
				__result__~bool^__result__)~bool
			)~bool^logical_and,
			__comprehension__(
			  // Variable
			  e,
			  // Target
			  x~google.expr.proto3.test.TestAllTypes^x.repeated_int64~list(int),
			  // Accumulator
			  __result__,
			  // Init
			  0~int,
			  // LoopCondition
			  true~bool,
			  // LoopStep
			  _?_:_(
				_==_(
				  e~int^e,
				  0~int
				)~bool^equals,
				_+_(
				  __result__~int^__result__,
				  1~int
				)~int^add_int64,
				__result__~int^__result__
			  )~int^conditional,
			  // Result
			  _==_(
				__result__~int^__result__,
				1~int
			  )~bool^equals)~bool
		  )~bool^logical_and`,
		outType: decls.Bool,
	},

	{
		in: `x.all(e, 0)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		err: `
ERROR: <input>:1:1: expression of type 'google.expr.proto3.test.TestAllTypes' cannot be range of a comprehension (must be list, map, or dynamic)
 | x.all(e, 0)
 | ^
ERROR: <input>:1:6: found no matching overload for '_&&_' applied to '(bool, int)'
 | x.all(e, 0)
 | .....^
		`,
	},
	{
		in: `lists.filter(x, x > 1.5)`,
		out: `__comprehension__(
			// Variable
			x,
			// Target
			lists~dyn^lists,
			// Accumulator
			__result__,
			// Init
			[]~list(dyn),
			// LoopCondition
			true~bool,
			// LoopStep
			_?_:_(
			  _>_(
				x~dyn^x,
				1.5~double
			  )~bool^greater_int64_double|greater_uint64_double|greater_double,
			  _+_(
				__result__~list(dyn)^__result__,
				[
				  x~dyn^x
				]~list(dyn)
			  )~list(dyn)^add_list,
			  __result__~list(dyn)^__result__
			)~list(dyn)^conditional,
			// Result
			__result__~list(dyn)^__result__)~list(dyn)`,
		outType: decls.NewListType(decls.Dyn),
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("lists", decls.Dyn),
			},
		},
	},

	{
		in: `.google.expr.proto3.test.TestAllTypes`,
		out: `google.expr.proto3.test.TestAllTypes
	~type(google.expr.proto3.test.TestAllTypes)
	^google.expr.proto3.test.TestAllTypes`,
		outType: decls.NewTypeType(
			decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
	},

	{
		in:        `test.TestAllTypes`,
		container: "google.expr.proto3",
		out: `
	google.expr.proto3.test.TestAllTypes
	~type(google.expr.proto3.test.TestAllTypes)
	^google.expr.proto3.test.TestAllTypes
		`,
		outType: decls.NewTypeType(
			decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
	},

	{
		in: `1 + x`,
		err: `
ERROR: <input>:1:5: undeclared reference to 'x' (in container '')
 | 1 + x
 | ....^`,
	},

	{
		in: `x == google.protobuf.Any{
				type_url:'types.googleapis.com/google.expr.proto3.test.TestAllTypes'
			} && x.single_nested_message.bb == 43
			|| x == google.expr.proto3.test.TestAllTypes{}
			|| y < x
			|| x >= x`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.Any),
				decls.NewVar("y", decls.NewWrapperType(decls.Int)),
			},
		},
		out: `
		_||_(
			_||_(
				_&&_(
					_==_(
						x~any^x,
						google.protobuf.Any{
							type_url:"types.googleapis.com/google.expr.proto3.test.TestAllTypes"~string
						}~any^google.protobuf.Any
					)~bool^equals,
					_==_(
						x~any^x.single_nested_message~dyn.bb~dyn,
						43~int
					)~bool^equals
				)~bool^logical_and,
				_==_(
					x~any^x,
					google.expr.proto3.test.TestAllTypes{}~google.expr.proto3.test.TestAllTypes^google.expr.proto3.test.TestAllTypes
				)~bool^equals
			)~bool^logical_or,
			_||_(
				_<_(
					y~wrapper(int)^y,
					x~any^x
				)~bool^less_int64|less_int64_double|less_int64_uint64,
				_>=_(
					x~any^x,
					x~any^x
				)~bool^greater_equals_bool|greater_equals_int64|greater_equals_int64_double|greater_equals_int64_uint64|greater_equals_uint64|greater_equals_uint64_double|greater_equals_uint64_int64|greater_equals_double|greater_equals_double_int64|greater_equals_double_uint64|greater_equals_string|greater_equals_bytes|greater_equals_timestamp|greater_equals_duration
			)~bool^logical_or
		)~bool^logical_or
		`,
		outType: decls.Bool,
	},

	{
		in:        `x`,
		container: "container",
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("container.x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out:     `container.x~google.expr.proto3.test.TestAllTypes^container.x`,
		outType: decls.NewObjectType("google.expr.proto3.test.TestAllTypes"),
	},

	{
		in: `list == type([1]) && map == type({1:2u})`,
		out: `
_&&_(_==_(list~type(list(dyn))^list,
           type([1~int]~list(int))~type(list(int))^type)
       ~bool^equals,
      _==_(map~type(map(dyn, dyn))^map,
            type({1~int : 2u~uint}~map(int, uint))~type(map(int, uint))^type)
        ~bool^equals)
  ~bool^logical_and
	`,
		outType: decls.Bool,
	},

	{
		in: `myfun(1, true, 3u) + 1.myfun(false, 3u).myfun(true, 42u)`,
		env: testEnv{
			functions: []*exprpb.Decl{
				decls.NewFunction("myfun",
					decls.NewInstanceOverload("myfun_instance",
						[]*exprpb.Type{decls.Int, decls.Bool, decls.Uint}, decls.Int),
					decls.NewOverload("myfun_static",
						[]*exprpb.Type{decls.Int, decls.Bool, decls.Uint}, decls.Int)),
			},
		},
		out: `_+_(
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
		outType: decls.Int,
	},

	{
		in: `size(x) > 4`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
			functions: []*exprpb.Decl{
				decls.NewFunction("size",
					decls.NewOverload("size_message",
						[]*exprpb.Type{decls.NewObjectType("google.expr.proto3.test.TestAllTypes")},
						decls.Int)),
			},
		},
		outType: decls.Bool,
	},

	{
		in: `x.single_int64_wrapper + 1 != 23`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
		_!=_(_+_(x~google.expr.proto3.test.TestAllTypes^x.single_int64_wrapper
		~wrapper(int),
		1~int)
		~int^add_int64,
		23~int)
		~bool^not_equals
		`,
		outType: decls.Bool,
	},

	{
		in: `x.single_int64_wrapper + y != 23`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
				decls.NewVar("y", decls.NewObjectType("google.protobuf.Int32Value")),
			},
		},
		out: `
		_!=_(
			_+_(
			  x~google.expr.proto3.test.TestAllTypes^x.single_int64_wrapper~wrapper(int),
			  y~wrapper(int)^y
			)~int^add_int64,
			23~int
		  )~bool^not_equals
		`,
		outType: decls.Bool,
	},

	{
		in: `1 in [1, 2, 3]`,
		out: `@in(
    		  1~int,
    		  [
    		    1~int,
    		    2~int,
    		    3~int
    		  ]~list(int)
    		)~bool^in_list`,
		outType: decls.Bool,
	},

	{
		in: `1 in dyn([1, 2, 3])`,
		out: `@in(
			1~int,
			dyn(
			  [
				1~int,
				2~int,
				3~int
			  ]~list(int)
			)~dyn^to_dyn
		  )~bool^in_list|in_map`,
		outType: decls.Bool,
	},

	{
		in: `type(null) == null_type`,
		out: `_==_(
    		  type(
    		    null~null
    		  )~type(null)^type,
    		  null_type~type(null)^null_type
    		)~bool^equals`,
		outType: decls.Bool,
	},

	{
		in: `type(type) == type`,
		out: `_==_(
		  type(
		    type~type(type())^type
		  )~type(type(type()))^type,
		  type~type(type())^type
		)~bool^equals`,
		outType: decls.Bool,
	},
	// Homogeneous aggregate type restriction tests.
	{
		in: `name in [1, 2u, 'string']`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("name", decls.String),
			},
			functions: []*exprpb.Decl{
				decls.NewFunction(operators.In,
					decls.NewOverload(overloads.InList,
						[]*exprpb.Type{
							decls.String,
							decls.NewListType(decls.String),
						}, decls.Bool)),
			},
		},
		opts:          []Option{HomogeneousAggregateLiterals(true)},
		disableStdEnv: true,
		out: `@in(
			name~string^name,
			[
				1~int,
				2u~uint,
				"string"~string
			]~list(string)
		)~bool^in_list`,
		err: `ERROR: <input>:1:13: expected type 'int' but found 'uint'
		| name in [1, 2u, 'string']
		| ............^`,
	},
	{
		in: `name in [1, 2, 3]`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("name", decls.String),
			},
			functions: []*exprpb.Decl{
				decls.NewFunction(operators.In,
					decls.NewOverload(overloads.InList,
						[]*exprpb.Type{
							decls.String,
							decls.NewListType(decls.String),
						}, decls.Bool)),
			},
		},
		opts:          []Option{HomogeneousAggregateLiterals(true)},
		disableStdEnv: true,
		out: `@in(
			name~string^name,
			[
				1~int,
				2~int,
				3~int
			]~list(int)
		)~!error!`,
		err: `ERROR: <input>:1:6: found no matching overload for '@in' applied to '(string, list(int))'
		| name in [1, 2, 3]
		| .....^`,
	},
	{
		in: `name in ["1", "2", "3"]`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("name", decls.String),
			},
			functions: []*exprpb.Decl{
				decls.NewFunction(operators.In,
					decls.NewOverload(overloads.InList,
						[]*exprpb.Type{
							decls.String,
							decls.NewListType(decls.String),
						}, decls.Bool)),
			},
		},
		opts:          []Option{HomogeneousAggregateLiterals(true)},
		disableStdEnv: true,
		out: `@in(
			name~string^name,
			[
				"1"~string,
				"2"~string,
				"3"~string
			]~list(string)
		)~bool^in_list`,
		outType: decls.Bool,
	},
	{
		in: `([[[1]], [[2]], [[3]]][0][0] + [2, 3, {'four': {'five': 'six'}}])[3]`,
		out: `_[_](
			_+_(
				_[_](
					_[_](
						[
							[
								[
									1~int
								]~list(int)
							]~list(list(int)),
							[
								[
									2~int
								]~list(int)
							]~list(list(int)),
							[
								[
									3~int
								]~list(int)
							]~list(list(int))
						]~list(list(list(int))),
						0~int
					)~list(list(int))^index_list,
					0~int
				)~list(int)^index_list,
				[
					2~int,
					3~int,
					{
						"four"~string:{
							"five"~string:"six"~string
						}~map(string, string)
					}~map(string, map(string, string))
				]~list(dyn)
			)~list(dyn)^add_list,
			3~int
		)~dyn^index_list`,
		outType: decls.Dyn,
	},
	{
		in: `[1] + [dyn('string')]`,
		out: `_+_(
			[
				1~int
			]~list(int),
			[
				dyn(
					"string"~string
				)~dyn^to_dyn
			]~list(dyn)
		)~list(dyn)^add_list`,
		outType: decls.NewListType(decls.Dyn),
	},
	{
		in: `[dyn('string')] + [1]`,
		out: `_+_(
			[
				dyn(
					"string"~string
				)~dyn^to_dyn
			]~list(dyn),
			[
				1~int
			]~list(int)
		)~list(dyn)^add_list`,
		outType: decls.NewListType(decls.Dyn),
	},
	{
		in: `[].map(x, [].map(y, x in y && y in x))`,
		err: `
		ERROR: <input>:1:33: found no matching overload for '@in' applied to '(type_param:"_var2" , type_param:"_var0" )'
		| [].map(x, [].map(y, x in y && y in x))
		| ................................^`,
	},
	{
		in: `args.user["myextension"].customAttributes.filter(x, x.name == "hobbies")`,
		out: `__comprehension__(
			// Variable
			x,
			// Target
			_[_](
			args~map(string, dyn)^args.user~dyn,
			"myextension"~string
			)~dyn^index_map.customAttributes~dyn,
			// Accumulator
			__result__,
			// Init
			[]~list(dyn),
			// LoopCondition
			true~bool,
			// LoopStep
			_?_:_(
			_==_(
				x~dyn^x.name~dyn,
				"hobbies"~string
			)~bool^equals,
			_+_(
				__result__~list(dyn)^__result__,
				[
				x~dyn^x
				]~list(dyn)
			)~list(dyn)^add_list,
			__result__~list(dyn)^__result__
			)~list(dyn)^conditional,
			// Result
			__result__~list(dyn)^__result__)~list(dyn)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("args", decls.NewMapType(decls.String, decls.Dyn)),
			},
		},
		outType: decls.NewListType(decls.Dyn),
	},
	{
		in: `a.b + 1 == a[0]`,
		out: `_==_(
			_+_(
			  a~dyn^a.b~dyn,
			  1~int
			)~int^add_int64,
			_[_](
			  a~dyn^a,
			  0~int
			)~dyn^index_list|index_map
		  )~bool^equals`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("a", decls.NewTypeParamType("T")),
			},
		},
		outType: decls.Bool,
	},
	{
		in: `!has(pb2.single_int64)
		&& !has(pb2.repeated_int32)
		&& !has(pb2.map_string_string)
		&& !has(pb3.single_int64)
		&& !has(pb3.repeated_int32)
		&& !has(pb3.map_string_string)`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("pb2", decls.NewObjectType("google.expr.proto2.test.TestAllTypes")),
				decls.NewVar("pb3", decls.NewObjectType("google.expr.proto3.test.TestAllTypes")),
			},
		},
		out: `
		_&&_(
			_&&_(
			  _&&_(
				!_(
				  pb2~google.expr.proto2.test.TestAllTypes^pb2.single_int64~test-only~~bool
				)~bool^logical_not,
				!_(
				  pb2~google.expr.proto2.test.TestAllTypes^pb2.repeated_int32~test-only~~bool
				)~bool^logical_not
			  )~bool^logical_and,
			  !_(
				pb2~google.expr.proto2.test.TestAllTypes^pb2.map_string_string~test-only~~bool
			  )~bool^logical_not
			)~bool^logical_and,
			_&&_(
			  _&&_(
				!_(
				  pb3~google.expr.proto3.test.TestAllTypes^pb3.single_int64~test-only~~bool
				)~bool^logical_not,
				!_(
				  pb3~google.expr.proto3.test.TestAllTypes^pb3.repeated_int32~test-only~~bool
				)~bool^logical_not
			  )~bool^logical_and,
			  !_(
				pb3~google.expr.proto3.test.TestAllTypes^pb3.map_string_string~test-only~~bool
			  )~bool^logical_not
			)~bool^logical_and
		  )~bool^logical_and`,
		outType: decls.Bool,
	},
	{
		in:        `TestAllTypes{}.repeated_nested_message`,
		container: "google.expr.proto2.test",
		out: `
		google.expr.proto2.test.TestAllTypes{}~google.expr.proto2.test.TestAllTypes^
		google.expr.proto2.test.TestAllTypes.repeated_nested_message
		~list(google.expr.proto2.test.TestAllTypes.NestedMessage)`,
		outType: decls.NewListType(
			decls.NewObjectType(
				"google.expr.proto2.test.TestAllTypes.NestedMessage",
			),
		),
	},
	{
		in:        `TestAllTypes{}.repeated_nested_message`,
		container: "google.expr.proto3.test",
		out: `
		google.expr.proto3.test.TestAllTypes{}~google.expr.proto3.test.TestAllTypes^
		google.expr.proto3.test.TestAllTypes.repeated_nested_message
		~list(google.expr.proto3.test.TestAllTypes.NestedMessage)`,
		outType: decls.NewListType(
			decls.NewObjectType(
				"google.expr.proto3.test.TestAllTypes.NestedMessage",
			),
		),
	},
	{
		in: `base64.encode('hello')`,
		env: testEnv{
			functions: []*exprpb.Decl{
				decls.NewFunction("base64.encode",
					decls.NewOverload(
						"base64_encode_string",
						[]*exprpb.Type{decls.String},
						decls.String)),
			},
		},
		out: `
		base64.encode(
			"hello"~string
		)~string^base64_encode_string`,
		outType: decls.String,
	},
	{
		in:        `encode('hello')`,
		container: `base64`,
		env: testEnv{
			functions: []*exprpb.Decl{
				decls.NewFunction("base64.encode",
					decls.NewOverload(
						"base64_encode_string",
						[]*exprpb.Type{decls.String},
						decls.String)),
			},
		},
		out: `
		base64.encode(
			"hello"~string
		)~string^base64_encode_string`,
		outType: decls.String,
	},
	{
		in:      `{}`,
		out:     `{}~map(dyn, dyn)`,
		outType: decls.NewMapType(decls.Dyn, decls.Dyn),
	},
	{
		in: `set([1, 2, 3])`,
		out: `
		set(
		  [
		    1~int,
		    2~int,
		    3~int
		  ]~list(int)
		)~abstract_type:{name:"set" parameter_types:{primitive:INT64}}^set_list`,
		env: testEnv{
			functions: []*exprpb.Decl{
				decls.NewFunction("set",
					decls.NewParameterizedOverload(
						"set_list", []*exprpb.Type{
							decls.NewListType(decls.NewTypeParamType("T")),
						}, decls.NewAbstractType("set", decls.NewTypeParamType("T")),
						[]string{"T"})),
			},
		},
		outType: decls.NewAbstractType("set", decls.Int),
	},
	{
		in: `set([1, 2]) == set([2, 1])`,
		out: `
		_==_(
		  set([1~int, 2~int]~list(int))~abstract_type:{name:"set" parameter_types:{primitive:INT64}}^set_list,
		  set([2~int, 1~int]~list(int))~abstract_type:{name:"set" parameter_types:{primitive:INT64}}^set_list
		)~bool^equals`,
		env: testEnv{
			functions: []*exprpb.Decl{
				decls.NewFunction("set",
					decls.NewParameterizedOverload(
						"set_list", []*exprpb.Type{
							decls.NewListType(decls.NewTypeParamType("T")),
						}, decls.NewAbstractType("set", decls.NewTypeParamType("T")),
						[]string{"T"})),
			},
		},
		outType: decls.Bool,
	},
	{
		in: `set([1, 2]) == x`,
		out: `
		_==_(
		  set([1~int, 2~int]~list(int))~abstract_type:{name:"set" parameter_types:{primitive:INT64}}^set_list,
		  x~abstract_type:{name:"set" parameter_types:{primitive:INT64}}^x
		)~bool^equals`,
		env: testEnv{
			idents: []*exprpb.Decl{
				decls.NewVar("x", decls.NewAbstractType("set", decls.NewTypeParamType("T"))),
			},
			functions: []*exprpb.Decl{
				decls.NewFunction("set",
					decls.NewParameterizedOverload(
						"set_list", []*exprpb.Type{
							decls.NewListType(decls.NewTypeParamType("T")),
						}, decls.NewAbstractType("set", decls.NewTypeParamType("T")),
						[]string{"T"})),
			},
		},
		outType: decls.Bool,
	},
	{
		in: `int{}`,
		err: `
		ERROR: <input>:1:4: 'int' is not a message type
		 | int{}
		 | ...^
		`,
	},
	{
		in: `Msg{}`,
		err: `
		ERROR: <input>:1:4: undeclared reference to 'Msg' (in container '')
		 | Msg{}
		 | ...^
		`,
	},
	{
		in: `fun()`,
		err: `
		ERROR: <input>:1:4: undeclared reference to 'fun' (in container '')
		 | fun()
		 | ...^
		`,
	},
	{
		in: `'string'.fun()`,
		err: `
		ERROR: <input>:1:13: undeclared reference to 'fun' (in container '')
		 | 'string'.fun()
		 | ............^
		`,
	},
	{
		in: `[].length`,
		err: `
		ERROR: <input>:1:3: type 'list_type:{elem_type:{type_param:"_var0"}}' does not support field selection
		 | [].length
		 | ..^
		`,
	},
	{
		in:   `1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1`,
		opts: []Option{CrossTypeNumericComparisons(false)},
		err: `
		ERROR: <input>:1:3: found no matching overload for '_<=_' applied to '(int, double)'
		 | 1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1
		 | ..^
		ERROR: <input>:1:16: found no matching overload for '_<=_' applied to '(uint, double)'
		 | 1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1
		 | ...............^
		ERROR: <input>:1:30: found no matching overload for '_<=_' applied to '(double, int)'
		 | 1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1
		 | .............................^
		ERROR: <input>:1:42: found no matching overload for '_<=_' applied to '(double, uint)'
		 | 1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1
		 | .........................................^
		ERROR: <input>:1:53: found no matching overload for '_<=_' applied to '(int, uint)'
		 | 1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1
		 | ....................................................^
		ERROR: <input>:1:65: found no matching overload for '_<=_' applied to '(uint, int)'
		 | 1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1
		 | ................................................................^
		`,
	},
	{
		in:      `1 <= 1.0 && 1u <= 1.0 && 1.0 <= 1 && 1.0 <= 1u && 1 <= 1u && 1u <= 1`,
		opts:    []Option{CrossTypeNumericComparisons(true)},
		outType: decls.Bool,
		out: `
		_&&_(
			_&&_(
			  _&&_(
				_<=_(
				  1~int,
				  1~double
				)~bool^less_equals_int64_double,
				_<=_(
				  1u~uint,
				  1~double
				)~bool^less_equals_uint64_double
			  )~bool^logical_and,
			  _<=_(
				1~double,
				1~int
			  )~bool^less_equals_double_int64
			)~bool^logical_and,
			_&&_(
			  _&&_(
				_<=_(
				  1~double,
				  1u~uint
				)~bool^less_equals_double_uint64,
				_<=_(
				  1~int,
				  1u~uint
				)~bool^less_equals_int64_uint64
			  )~bool^logical_and,
			  _<=_(
				1u~uint,
				1~int
			  )~bool^less_equals_uint64_int64
			)~bool^logical_and
		  )~bool^logical_and`,
	},
	{
		in:      `1 < 1.0 && 1u < 1.0 && 1.0 < 1 && 1.0 < 1u && 1 < 1u && 1u < 1`,
		opts:    []Option{CrossTypeNumericComparisons(true)},
		outType: decls.Bool,
		out: `
		_&&_(
			_&&_(
			  _&&_(
				_<_(
				  1~int,
				  1~double
				)~bool^less_int64_double,
				_<_(
				  1u~uint,
				  1~double
				)~bool^less_uint64_double
			  )~bool^logical_and,
			  _<_(
				1~double,
				1~int
			  )~bool^less_double_int64
			)~bool^logical_and,
			_&&_(
			  _&&_(
				_<_(
				  1~double,
				  1u~uint
				)~bool^less_double_uint64,
				_<_(
				  1~int,
				  1u~uint
				)~bool^less_int64_uint64
			  )~bool^logical_and,
			  _<_(
				1u~uint,
				1~int
			  )~bool^less_uint64_int64
			)~bool^logical_and
		  )~bool^logical_and`,
	},
	{
		in:      `1 > 1.0 && 1u > 1.0 && 1.0 > 1 && 1.0 > 1u && 1 > 1u && 1u > 1`,
		opts:    []Option{CrossTypeNumericComparisons(true)},
		outType: decls.Bool,
		out: `
		_&&_(
			_&&_(
			  _&&_(
				_>_(
				  1~int,
				  1~double
				)~bool^greater_int64_double,
				_>_(
				  1u~uint,
				  1~double
				)~bool^greater_uint64_double
			  )~bool^logical_and,
			  _>_(
				1~double,
				1~int
			  )~bool^greater_double_int64
			)~bool^logical_and,
			_&&_(
			  _&&_(
				_>_(
				  1~double,
				  1u~uint
				)~bool^greater_double_uint64,
				_>_(
				  1~int,
				  1u~uint
				)~bool^greater_int64_uint64
			  )~bool^logical_and,
			  _>_(
				1u~uint,
				1~int
			  )~bool^greater_uint64_int64
			)~bool^logical_and
		  )~bool^logical_and`,
	},
	{
		in:      `1 >= 1.0 && 1u >= 1.0 && 1.0 >= 1 && 1.0 >= 1u && 1 >= 1u && 1u >= 1`,
		opts:    []Option{CrossTypeNumericComparisons(true)},
		outType: decls.Bool,
		out: `
		_&&_(
			_&&_(
			  _&&_(
				_>=_(
				  1~int,
				  1~double
				)~bool^greater_equals_int64_double,
				_>=_(
				  1u~uint,
				  1~double
				)~bool^greater_equals_uint64_double
			  )~bool^logical_and,
			  _>=_(
				1~double,
				1~int
			  )~bool^greater_equals_double_int64
			)~bool^logical_and,
			_&&_(
			  _&&_(
				_>=_(
				  1~double,
				  1u~uint
				)~bool^greater_equals_double_uint64,
				_>=_(
				  1~int,
				  1u~uint
				)~bool^greater_equals_int64_uint64
			  )~bool^logical_and,
			  _>=_(
				1u~uint,
				1~int
			  )~bool^greater_equals_uint64_int64
			)~bool^logical_and
		  )~bool^logical_and`,
	},
}

var testEnvs = map[string]testEnv{
	"default": {
		functions: []*exprpb.Decl{
			decls.NewFunction("fg_s",
				decls.NewOverload("fg_s_0", []*exprpb.Type{}, decls.String)),
			decls.NewFunction("fi_s_s",
				decls.NewInstanceOverload("fi_s_s_0",
					[]*exprpb.Type{decls.String}, decls.String)),
		},
		idents: []*exprpb.Decl{
			decls.NewVar("is", decls.String),
			decls.NewVar("ii", decls.Int),
			decls.NewVar("iu", decls.Uint),
			decls.NewVar("iz", decls.Bool),
			decls.NewVar("ib", decls.Bytes),
			decls.NewVar("id", decls.Double),
			decls.NewVar("ix", decls.Null),
		},
	},
}

type testInfo struct {
	// in contains the expression to be parsed.
	in string

	// out contains the output.
	out string

	// outType is the expected type of the expression
	outType *exprpb.Type

	// container is the container name to use for test.
	container string

	// env is the environment to use for testing.
	env testEnv

	// err is the expected error for negative test cases.
	err string

	// disableStdEnv indicates whether the standard functions should be disabled.
	disableStdEnv bool

	// opts is the set of checker Option flags to use when type-checking.
	opts []Option
}

type testEnv struct {
	idents    []*exprpb.Decl
	functions []*exprpb.Decl
}

func TestCheck(t *testing.T) {
	for i, tst := range testCases {
		name := fmt.Sprintf("%d %s", i, tst.in)
		tc := tst
		t.Run(name, func(t *testing.T) {
			// Runs the tests in parallel to ensure that there are no data races
			// due to shared mutable state across tests.
			t.Parallel()

			src := common.NewTextSource(tc.in)
			expression, errors := parser.Parse(src)
			if len(errors.GetErrors()) > 0 {
				t.Fatalf("Unexpected parse errors: %v", errors.ToDisplayString())
			}

			reg, err := types.NewRegistry(&proto2pb.TestAllTypes{}, &proto3pb.TestAllTypes{})
			if err != nil {
				t.Fatalf("types.NewRegistry() failed: %v", err)
			}
			cont, err := containers.NewContainer(containers.Name(tc.container))
			if err != nil {
				t.Fatalf("containers.NewContainer() failed: %v", err)
			}
			opts := []Option{CrossTypeNumericComparisons(true)}
			if len(tc.opts) != 0 {
				opts = tc.opts
			}
			env, err := NewEnv(cont, reg, opts...)
			if err != nil {
				t.Fatalf("NewEnv(cont, reg) failed: %v", err)
			}
			if !tc.disableStdEnv {
				env.Add(StandardDeclarations()...)
			}
			if tc.env.idents != nil {
				for _, ident := range tc.env.idents {
					env.Add(ident)
				}
			}
			if tc.env.functions != nil {
				for _, fn := range tc.env.functions {
					env.Add(fn)
				}
			}

			semantics, errors := Check(expression, src, env)
			if len(errors.GetErrors()) > 0 {
				errorString := errors.ToDisplayString()
				if tc.err != "" {
					if !test.Compare(errorString, tc.err) {
						t.Error(test.DiffMessage("Error mismatch", errorString, tc.err))
					}
				} else {
					t.Errorf("Unexpected type-check errors: %v", errorString)
				}
			} else if tc.err != "" {
				t.Errorf("Expected error not thrown: %s", tc.err)
			}

			actual := semantics.TypeMap[expression.Expr.Id]
			if tc.err == "" {
				if actual == nil || !proto.Equal(actual, tc.outType) {
					t.Error(test.DiffMessage("Type Error", actual, tc.outType))
				}
			}

			if tc.out != "" {
				actualStr := Print(expression.Expr, semantics)
				if !test.Compare(actualStr, tc.out) {
					t.Error(test.DiffMessage("Structure error", actualStr, tc.out))
				}
			}
		})
	}
}
