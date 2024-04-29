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
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
)

func TestCompile(t *testing.T) {
	for _, tst := range policyTests {
		r := newRunner(t, tst.name, tst.envOpts...)
		r.run(t)
	}
}

func BenchmarkCompile(b *testing.B) {
	for _, tst := range policyTests {
		r := newRunner(b, tst.name, tst.envOpts...)
		r.bench(b)
	}
}

func newRunner(t testing.TB, name string, opts ...cel.EnvOption) *runner {
	r := &runner{name: name, envOptions: opts}
	r.setup(t)
	return r
}

type runner struct {
	name       string
	envOptions []cel.EnvOption
	env        *cel.Env
	prg        cel.Program
}

func (r *runner) setup(t testing.TB) {
	config := readPolicyConfig(t, fmt.Sprintf("testdata/%s/config.yaml", r.name))
	srcFile := readPolicy(t, fmt.Sprintf("testdata/%s/policy.yaml", r.name))
	p, iss := Parse(srcFile)
	if iss.Err() != nil {
		t.Fatalf("parse() failed: %v", iss.Err())
	}
	if p.name.value != r.name {
		t.Errorf("policy name is %v, wanted %s", p.name, r.name)
	}
	env, err := cel.NewEnv(
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.ExtendedValidations())
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
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
	// Configure any implementations
	env, err = env.Extend(r.envOptions...)
	if err != nil {
		t.Fatalf("env.Extend() with config options %v, failed: %v", config, err)
	}
	ast, iss := Compile(env, p)
	if iss.Err() != nil {
		t.Fatalf("Compile() failed: %v", iss.Err())
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
				out, _, err := r.prg.Eval(tc.Input)
				if err != nil {
					t.Fatalf("prg.Eval(tc.Input) failed: %v", err)
				}
				wantExpr, iss := r.env.Compile(tc.Output)
				if iss.Err() != nil {
					t.Fatalf("env.Compile(%q) failed :%v", tc.Output, iss.Err())
				}
				testPrg, err := r.env.Program(wantExpr)
				if err != nil {
					t.Fatalf("env.Program(wantExpr) failed: %v", err)
				}
				testOut, _, err := testPrg.Eval(cel.NoVars())
				if err != nil {
					t.Fatalf("testPrg.Eval() failed: %v", err)
				}
				if optOut, ok := out.(*types.Optional); ok {
					if optOut.Equal(types.OptionalNone) == types.True {
						if testOut.Equal(types.OptionalNone) != types.True {
							t.Errorf("policy eval got %v, wanted %v", out, testOut)
						}
					} else if testOut.Equal(optOut.GetValue()) != types.True {
						t.Errorf("policy eval got %v, wanted %v", out, testOut)
					}
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
