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
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func TestCompile(t *testing.T) {
	for _, tst := range policyTests {
		r := newRunner(t, tst.name, tst.expr, tst.parseOpts, tst.envOpts...)
		r.run(t)
	}
}

func TestCompileError(t *testing.T) {
	for _, tst := range policyErrorTests {
		_, _, iss := compile(t, tst.name, []ParserOption{}, []cel.EnvOption{}, tst.compilerOpts)
		if iss.Err() == nil {
			t.Fatalf("compile(%s) did not error, wanted %s", tst.name, tst.err)
		}
		if iss.Err().Error() != tst.err {
			t.Errorf("compile(%s) got error %s, wanted %s", tst.name, iss.Err().Error(), tst.err)
		}
	}
}

func BenchmarkCompile(b *testing.B) {
	for _, tst := range policyTests {
		r := newRunner(b, tst.name, tst.expr, tst.parseOpts, tst.envOpts...)
		r.bench(b)
	}
}

func newRunner(t testing.TB, name, expr string, parseOpts []ParserOption, opts ...cel.EnvOption) *runner {
	r := &runner{
		name:      name,
		envOpts:   opts,
		parseOpts: parseOpts,
		expr:      expr}
	r.setup(t)
	return r
}

type runner struct {
	name         string
	envOpts      []cel.EnvOption
	parseOpts    []ParserOption
	compilerOpts []CompilerOption
	env          *cel.Env
	expr         string
	prg          cel.Program
}

func compile(t testing.TB, name string, parseOpts []ParserOption, envOpts []cel.EnvOption, compilerOpts []CompilerOption) (*cel.Env, *cel.Ast, *cel.Issues) {
	config := readPolicyConfig(t, fmt.Sprintf("testdata/%s/config.yaml", name))
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
	env, err := cel.NewEnv(
		cel.DefaultUTCTimeZone(true),
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.ExtendedValidations())
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	// Configure any custom environment options.
	env, err = env.Extend(envOpts...)
	if err != nil {
		t.Fatalf("env.Extend() with env options %v, failed: %v", config, err)
	}
	// Configure declarations
	configOpts, err := config.AsEnvOptions(env)
	if err != nil {
		t.Fatalf("config.AsEnvOptions() failed: %v", err)
	}
	env, err = env.Extend(configOpts...)
	if err != nil {
		t.Fatalf("env.Extend() with config options %v, failed: %v", config, err)
	}
	ast, iss := Compile(env, policy, compilerOpts...)
	return env, ast, iss
}

func (r *runner) setup(t testing.TB) {
	env, ast, iss := compile(t, r.name, r.parseOpts, r.envOpts, r.compilerOpts)
	if iss.Err() != nil {
		t.Fatalf("Compile() failed: %v", iss.Err())
	}
	pExpr, err := cel.AstToString(ast)
	if err != nil {
		t.Fatalf("cel.AstToString() failed: %v", err)
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
				for k, v := range tc.Input {
					if v.Expr == "" {
						input[k] = v.Value
						continue
					}
					input[k] = r.eval(t, v.Expr)
				}
				out, _, err := r.prg.Eval(input)
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
				for i := 0; i < b.N; i++ {
					_, _, err := r.prg.Eval(tc.Input)
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

func normalize(s string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(s, " ", ""), "\n", ""),
		"\t", "")
}
