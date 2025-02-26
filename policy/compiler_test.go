// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/google/cel-go/interpreter"
)

func TestCompile(t *testing.T) {
	for _, tst := range policyTests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			r := newRunner(tc.name, tc.expr, tc.parseOpts)
			env, ast, iss := r.compile(t, tc.envOpts, []CompilerOption{})
			if iss.Err() != nil {
				t.Fatalf("Compile(%s) failed: %v", r.name, iss.Err())
			}
			r.setup(t, env, ast)
			r.run(t)
		})
	}
}

func TestRuleComposerError(t *testing.T) {
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("NewEnv() failed: %v", err)
	}
	_, err = NewRuleComposer(env, ExpressionUnnestHeight(-1))
	if err == nil || !strings.Contains(err.Error(), "invalid unnest") {
		t.Errorf("NewRuleComposer() got %v, wanted 'invalid unnest'", err)
	}
}

func TestRuleComposerUnnest(t *testing.T) {
	for _, tst := range composerUnnestTests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			r := newRunner(tc.name, tc.expr, []ParserOption{})
			env, rule, iss := r.compileRule(t)
			if iss.Err() != nil {
				t.Fatalf("CompileRule() failed: %v", iss.Err())
			}
			rc, err := NewRuleComposer(env, tc.composerOpts...)
			if err != nil {
				t.Fatalf("NewRuleComposer() failed: %v", err)
			}
			ast, iss := rc.Compose(rule)
			if iss.Err() != nil {
				t.Fatalf("Compose(rule) failed: %v", iss.Err())
			}
			unparsed, err := cel.AstToString(ast)
			if err != nil {
				t.Fatalf("cel.AstToString() failed: %v", err)
			}
			if normalize(unparsed) != normalize(tc.composed) {
				t.Errorf("cel.AstToString() got %s, wanted %s", unparsed, tc.composed)
			}
			if !ast.OutputType().IsEquivalentType(tc.outputType) {
				t.Errorf("ast.OutputType() got %v, wanted %v", ast.OutputType(), tc.outputType)
			}
			r.setup(t, env, ast)
			r.run(t)
		})
	}
}

func TestCompileError(t *testing.T) {
	for _, tst := range policyErrorTests {
		policy := parsePolicy(t, tst.name, []ParserOption{})
		_, _, iss := compile(t, tst.name, policy, []cel.EnvOption{}, tst.compilerOpts)
		if iss.Err() == nil {
			t.Fatalf("compile(%s) did not error, wanted %s", tst.name, tst.err)
		}
		if iss.Err().Error() != tst.err {
			t.Errorf("compile(%s) got error %s, wanted %s", tst.name, iss.Err().Error(), tst.err)
		}
	}
}

func TestCompiledRuleHasOptionalOutput(t *testing.T) {
	env, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	tests := []struct {
		rule     *CompiledRule
		optional bool
	}{
		{rule: &CompiledRule{}, optional: false},
		{
			rule: &CompiledRule{
				matches: []*CompiledMatch{{}},
			},
			optional: true,
		},
		{
			rule: &CompiledRule{
				matches: []*CompiledMatch{{}},
			},
			optional: true,
		},
		{
			rule: &CompiledRule{
				matches: []*CompiledMatch{{cond: mustCompileExpr(t, env, "true")}},
			},
			optional: false,
		},
		{
			rule: &CompiledRule{
				matches: []*CompiledMatch{{cond: mustCompileExpr(t, env, "1 < 0")}},
			},
			optional: true,
		},
	}
	for _, tst := range tests {
		got := tst.rule.HasOptionalOutput()
		if got != tst.optional {
			t.Errorf("rule.HasOptionalOutput() got %v, wanted, %v", got, tst.optional)
		}
	}
}

func TestMaxNestedExpressions_Error(t *testing.T) {
	policyName := "required_labels"
	wantError := `ERROR: testdata/required_labels/policy.yaml:15:8: error configuring compiler option: nested expression limit must be non-negative, non-zero value: -1
 | name: "required_labels"
 | .......^`
	policy := parsePolicy(t, policyName, []ParserOption{})
	_, _, iss := compile(t, policyName, policy, []cel.EnvOption{}, []CompilerOption{MaxNestedExpressions(-1)})
	if iss.Err() == nil {
		t.Fatalf("compile(%s) did not error, wanted %s", policyName, wantError)
	}
	if iss.Err().Error() != wantError {
		t.Errorf("compile(%s) got error %s, wanted %s", policyName, iss.Err().Error(), wantError)
	}
}

func BenchmarkCompile(b *testing.B) {
	for _, tst := range policyTests {
		r := newRunner(tst.name, tst.expr, tst.parseOpts)
		env, ast, iss := r.compile(b, tst.envOpts, []CompilerOption{})
		if iss.Err() != nil {
			b.Fatalf("Compile() failed: %v", iss.Err())
		}
		r.setup(b, env, ast)
		r.bench(b)
	}
}

func newRunner(name, expr string, parseOpts []ParserOption, opts ...cel.EnvOption) *runner {
	return &runner{
		name:      name,
		parseOpts: parseOpts,
		expr:      expr}
}

type runner struct {
	name      string
	parseOpts []ParserOption
	env       *cel.Env
	expr      string
	prg       cel.Program
}

func (r *runner) compile(t testing.TB, envOpts []cel.EnvOption, compilerOpts []CompilerOption) (*cel.Env, *cel.Ast, *cel.Issues) {
	policy := parsePolicy(t, r.name, r.parseOpts)
	return compile(t, r.name, policy, envOpts, compilerOpts)
}

func (r *runner) compileRule(t testing.TB) (*cel.Env, *CompiledRule, *cel.Issues) {
	t.Helper()
	config := readPolicyConfig(t, fmt.Sprintf("testdata/%s/config.yaml", r.name))
	policy := parsePolicy(t, r.name, r.parseOpts)
	env, err := cel.NewCustomEnv(
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.ExtendedValidations(),
		ext.Bindings())
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	// Configure declarations
	env, err = env.Extend(FromConfig(config))
	if err != nil {
		t.Fatalf("env.Extend() with config options %v, failed: %v", config, err)
	}
	rule, iss := CompileRule(env, policy)
	return env, rule, iss
}

func (r *runner) setup(t testing.TB, env *cel.Env, ast *cel.Ast) {
	t.Helper()
	pExpr, err := cel.AstToString(ast)
	if err != nil {
		t.Fatalf("cel.AstToString() failed: %v", err)
	}
	_, err = cel.AstToCheckedExpr(ast)
	if err != nil {
		t.Fatalf("cel.AstToCheckedExpr() failed: %v", err)
	}
	if r.expr != "" && normalize(pExpr) != normalize(r.expr) {
		t.Errorf("cel.AstToString() got %s, wanted %s", pExpr, r.expr)
	}
	prg, err := env.Program(ast, cel.EvalOptions(cel.OptOptimize))
	if err != nil {
		t.Fatalf("env.Program() failed: %v", err)
	}
	r.env = env
	r.prg = prg
}

func (r *runner) run(t *testing.T) {
	tests := readTestSuite(t, fmt.Sprintf("testdata/%s/tests.yaml", r.name))
	for _, s := range tests.Sections {
		section := s.Name
		for _, tst := range s.Tests {
			tc := tst
			t.Run(fmt.Sprintf("%s/%s/%s", r.name, section, tc.Name), func(t *testing.T) {
				input := map[string]any{}
				var err error
				var activation interpreter.Activation
				for k, v := range tc.Input {
					if v.Expr != "" {
						input[k] = r.eval(t, v.Expr)
						continue
					}
					if v.ContextExpr != "" {
						ctx, err := r.eval(t, v.ContextExpr).ConvertToNative(
							reflect.TypeOf(((*proto.Message)(nil))).Elem())
						if err != nil {
							t.Fatalf("context variable is not a valid proto: %v", err)
						}
						activation, err = cel.ContextProtoVars(ctx.(proto.Message))
						if err != nil {
							t.Fatalf("cel.ContextProtoVars() failed: %v", err)
						}
						break
					}
					input[k] = v.Value
				}
				if activation == nil {
					activation, err = interpreter.NewActivation(input)
					if err != nil {
						t.Fatalf("interpreter.NewActivation(input) failed: %v", err)
					}
				}
				out, _, err := r.prg.Eval(activation)
				if err != nil {
					t.Fatalf("prg.Eval(input) failed: %v", err)
				}
				testOut := r.eval(t, tc.Output)
				if optOut, ok := out.(*types.Optional); ok {
					if optOut.Equal(types.OptionalNone) == types.True {
						if testOut.Equal(types.OptionalNone) != types.True {
							t.Errorf("policy eval got %v, wanted %v", out, testOut)
						}
					} else if testOut.Equal(optOut.GetValue()) != types.True {
						t.Errorf("policy eval got %v, wanted %v", out, testOut)
					}
				} else if testOut.Equal(out) != types.True {
					t.Errorf("policy eval got %v, wanted %v", out, testOut)
				}
			})
		}
	}
}

func (r *runner) bench(b *testing.B) {
	tests := readTestSuite(b, fmt.Sprintf("testdata/%s/tests.yaml", r.name))
	for _, s := range tests.Sections {
		section := s.Name
		for _, tst := range s.Tests {
			tc := tst
			b.Run(fmt.Sprintf("%s/%s/%s", r.name, section, tc.Name), func(b *testing.B) {
				input := map[string]any{}
				var err error
				var activation interpreter.Activation
				for k, v := range tc.Input {
					if v.Expr != "" {
						input[k] = r.eval(b, v.Expr)
						continue
					}
					if v.ContextExpr != "" {
						ctx, err := r.eval(b, v.ContextExpr).ConvertToNative(
							reflect.TypeOf(((*proto.Message)(nil))).Elem())
						if err != nil {
							b.Fatalf("context variable is not a valid proto: %v", err)
						}
						activation, err = cel.ContextProtoVars(ctx.(proto.Message))
						if err != nil {
							b.Fatalf("cel.ContextProtoVars() failed: %v", err)
						}
						break
					}
					input[k] = v.Value
				}
				if activation == nil {
					activation, err = interpreter.NewActivation(input)
					if err != nil {
						b.Fatalf("interpreter.NewActivation(input) failed: %v", err)
					}
				}
				for i := 0; i < b.N; i++ {
					_, _, err := r.prg.Eval(activation)
					if err != nil {
						b.Fatalf("policy eval failed: %v", err)
					}
				}
			})
		}
	}
}

func (r *runner) eval(t testing.TB, expr string) ref.Val {
	wantExpr, iss := r.env.Compile(expr)
	if iss.Err() != nil {
		t.Fatalf("env.Compile(%q) failed :%v", expr, iss.Err())
	}
	prg, err := r.env.Program(wantExpr)
	if err != nil {
		t.Fatalf("env.Program(wantExpr) failed: %v", err)
	}
	out, _, err := prg.Eval(cel.NoVars())
	if err != nil {
		t.Fatalf("prg.Eval() failed: %v", err)
	}
	return out
}

func mustCompileExpr(t testing.TB, env *cel.Env, expr string) *cel.Ast {
	t.Helper()
	out, iss := env.Compile(expr)
	if iss.Err() != nil {
		t.Fatalf("env.Compile(%s) failed: %v", expr, iss.Err())
	}
	return out
}

func parsePolicy(t testing.TB, name string, parseOpts []ParserOption) *Policy {
	t.Helper()
	srcFile := readPolicy(t, fmt.Sprintf("testdata/%s/policy.yaml", name))
	parser, err := NewParser(parseOpts...)
	if err != nil {
		t.Fatalf("NewParser() failed: %v", err)
	}
	policy, iss := parser.Parse(srcFile)
	if iss.Err() != nil {
		t.Fatalf("Parse() failed: %v", iss.Err())
	}
	if policy.name.Value != name {
		t.Errorf("policy name is %v, wanted %q", policy.name, name)
	}
	return policy
}

func compile(t testing.TB, name string, policy *Policy, envOpts []cel.EnvOption, compilerOpts []CompilerOption) (*cel.Env, *cel.Ast, *cel.Issues) {
	config := readPolicyConfig(t, fmt.Sprintf("testdata/%s/config.yaml", name))
	env, err := cel.NewCustomEnv(
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.ExtendedValidations(),
		ext.Bindings())
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	// Configure any custom environment options.
	env, err = env.Extend(envOpts...)
	if err != nil {
		t.Fatalf("env.Extend() with env options %v, failed: %v", config, err)
	}
	// Configure declarations
	env, err = env.Extend(FromConfig(config))
	if err != nil {
		t.Fatalf("env.Extend() with config options %v, failed: %v", config, err)
	}
	ast, iss := Compile(env, policy, compilerOpts...)
	return env, ast, iss
}

func normalize(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(s, " ", ""), "\n", ""),
		"\t", "")
}
