package policy

import (
	"testing"
)

func TestParse(t *testing.T) {
	srcFile := readPolicy(t, "testdata/required_labels/policy.yaml")
	p, iss := Parse(srcFile)
	if iss.Err() != nil {
		t.Fatalf("parse() failed: %v", iss.Err())
	}
	if p.name.value != "required_labels" {
		t.Errorf("policy name is %v, wanted 'required_labels'", p.name)
	}
}
