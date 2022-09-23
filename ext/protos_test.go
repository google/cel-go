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
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/test"
	"github.com/google/cel-go/test/proto2pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestProtos(t *testing.T) {
	var protosTests = []struct {
		expr      string
		err       string
		parseErr  string
		parseOnly bool
	}{
		// Positive test cases.
		{expr: `proto.getExt(ExampleType{}, google.expr.proto2.test.int32_ext) == 0`},
		{expr: `!proto.hasExt(ExampleType{}, google.expr.proto2.test.int32_wrapper_ext)`},
		{expr: `proto.getExt(ExampleType{}, google.expr.proto2.test.int32_wrapper_ext) == null`},
		{expr: `!proto.hasExt(ExampleType{}, google.expr.proto2.test.nested_example)`},
		{expr: `proto.getExt(ExampleType{}, google.expr.proto2.test.nested_example) == ExampleType{}`},
		{expr: `proto.getExt(ExampleType{}, google.expr.proto2.test.ExtendedExampleType.extended_examples) == []`},
		{expr: `proto.getExt(ExampleType{}, google.expr.proto2.test.ExtendedExampleType.enum_ext) == GlobalEnum.GOO`},
		{expr: `proto.getExt(msg, google.expr.proto2.test.int32_ext) == 42`},
		{expr: `proto.getExt(msg, google.expr.proto2.test.int32_wrapper_ext) == 21`},
		{expr: `proto.hasExt(msg, google.expr.proto2.test.nested_example)`},
		{expr: `proto.getExt(msg, google.expr.proto2.test.nested_example) == ExampleType{name: 'nested'}`},
		{expr: `proto.getExt(msg, google.expr.proto2.test.ExtendedExampleType.extended_examples) == ['example1', 'example2']`},
		{expr: `proto.getExt(msg, google.expr.proto2.test.ExtendedExampleType.enum_ext) == GlobalEnum.GAZ`},

		// Negative test cases.
		{
			expr:      `proto.getExt(ExtendedExampleType{}, google.expr.proto2.test.ExtendedExampleType.enum_ext)`,
			err:       "no such field 'google.expr.proto2.test.ExtendedExampleType.enum_ext'",
			parseOnly: true,
		},
		{
			expr: `proto.getExt(ExtendedExampleType{}, enum_ext)`,
			parseErr: `ERROR: <input>:1:37: invalid extension field
			| proto.getExt(ExtendedExampleType{}, enum_ext)
			| ....................................^`,
		},
		{
			expr: `proto.hasExt(ExtendedExampleType{}, call().enum_ext)`,
			parseErr: `ERROR: <input>:1:43: invalid extension field
			| proto.hasExt(ExtendedExampleType{}, call().enum_ext)
			| ..........................................^`,
		},
		{
			expr: `proto.getExt(ExtendedExampleType{}, has(google.expr.proto2.test.int32_ext))`,
			parseErr: `ERROR: <input>:1:40: invalid extension field
			| proto.getExt(ExtendedExampleType{}, has(google.expr.proto2.test.int32_ext))
			| .......................................^`,
		},

		// Even though 'getExt' is the macro, the call is left unexpanded is an identifier
		// that's not named 'proto'
		{
			expr: `msg.getExt("google.expr.proto2.test.int32_ext", 0) == 42`,
		},
		// Even though 'getExt' is the macro, the call is left unexpanded as the operand is a call
		{
			expr: `dyn(msg).getExt("google.expr.proto2.test.int32_ext", 0) == 42`,
		},
	}

	env, err := cel.NewEnv(
		cel.Container("google.expr.proto2.test"),
		cel.Types(&proto2pb.ExampleType{}),
		cel.Variable("msg", cel.ObjectType("google.expr.proto2.test.ExampleType")),
		cel.Function("getExt",
			cel.MemberOverload("msg_getExt_field_default",
				[]*cel.Type{cel.DynType, cel.StringType, cel.DynType},
				cel.DynType),
			cel.SingletonFunctionBinding(func(args ...ref.Val) ref.Val {
				msg := args[0]
				field := args[1]
				indexer := msg.(traits.Indexer)
				return indexer.Get(field)
			})),
		Protos(),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv(Protos()) failed: %v", err)
	}
	for i, tst := range protosTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				if tc.parseErr == "" {
					t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
				}
				if !test.Compare(iss.Err().Error(), tc.parseErr) {
					t.Errorf("got parse error %v, wanted error %s for expr: %s", iss.Err(), tc.parseErr, tc.expr)
				}
				return
			}
			asts = append(asts, pAst)
			if !tc.parseOnly {
				cAst, iss := env.Check(pAst)
				if iss.Err() != nil {
					t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
				}
				asts = append(asts, cAst)
			}
			msg := &proto2pb.ExampleType{
				Name: proto.String("example0"),
			}
			proto.SetExtension(msg, proto2pb.E_Int32Ext, int32(42))
			proto.SetExtension(msg, proto2pb.E_Int32WrapperExt, wrapperspb.Int32(21))
			proto.SetExtension(msg, proto2pb.E_NestedExample, &proto2pb.ExampleType{Name: proto.String("nested")})
			proto.SetExtension(msg, proto2pb.E_ExtendedExampleType_EnumExt, proto2pb.GlobalEnum_GAZ)
			proto.SetExtension(msg, proto2pb.E_ExtendedExampleType_ExtendedExamples, []string{"example1", "example2"})
			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatal(err)
				}
				out, _, err := prg.Eval(map[string]any{"msg": msg})
				if tc.err != "" {
					if err == nil {
						t.Fatalf("got value %v, wanted error %s for expr: %s",
							out.Value(), tc.err, tc.expr)
					}
					if !test.Compare(err.Error(), tc.err) {
						t.Errorf("got error %v, wanted error %s for expr: %s", err, tc.err, tc.expr)
					}
				} else if err != nil {
					t.Fatal(err)
				} else if out.Value() != true {
					t.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}
