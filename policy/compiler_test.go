package policy

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	celext "github.com/google/cel-go/ext"
)

func TestCompile(t *testing.T) {
	r := newRunner(t, "required_labels",
		cel.Variable("spec", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("resource", cel.MapType(cel.StringType, cel.DynType)))
	r.run(t)
}

func BenchmarkCompile(b *testing.B) {
	r := newRunner(b, "required_labels",
		cel.Variable("spec", cel.MapType(cel.StringType, cel.DynType)),
		cel.Variable("resource", cel.MapType(cel.StringType, cel.DynType)))
	r.bench(b)
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
	srcFile := readPolicy(t, fmt.Sprintf("testdata/%s/policy.yaml", r.name))
	p, iss := Parse(srcFile)
	if iss.Err() != nil {
		t.Fatalf("parse() failed: %v", iss.Err())
	}
	if p.name.value != r.name {
		t.Errorf("policy name is %v, wanted %s", p.name, r.name)
	}
	envOpts := append([]cel.EnvOption{
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.ExtendedValidations(),
		celext.Strings(),
	}, r.envOptions...)
	env, err := cel.NewEnv(envOpts...)
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	ast, iss := Compile(env, p)
	if iss.Err() != nil {
		t.Errorf("Compile() failed: %v", iss.Err())
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
		for i, tst := range s.Tests {
			tc := tst
			t.Run(fmt.Sprintf("%s:%d", s.Name, i), func(t *testing.T) {
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
		for _, tst := range s.Tests {
			tc := tst
			b.Run(fmt.Sprintf("%s", s.Name), func(b *testing.B) {
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
