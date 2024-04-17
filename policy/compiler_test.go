package policy

import (
	"testing"

	"github.com/google/cel-go/cel"
)

func TestCompile(t *testing.T) {
	srcFile := readPolicy(t, "testdata/required_labels.yaml")
	p, iss := parse(srcFile)
	if iss.Err() != nil {
		t.Fatalf("parse() failed: %v", iss.Err())
	}
	if p.name.value != "required_labels" {
		t.Errorf("policy name is %v, wanted 'required_labels'", p.name)
	}
	env, err := cel.NewEnv(
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking(),
		cel.ExtendedValidations(),
		cel.Variable("rule.labels", cel.MapType(cel.StringType, cel.StringType)),
		cel.Variable("resource.labels", cel.MapType(cel.StringType, cel.StringType)),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	_, iss = compile(env, p)
	if iss.Err() != nil {
		t.Errorf("compile() failed: %v", iss.Err())
	}
}
