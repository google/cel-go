package policy

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	srcFile := readPolicy(t, "testdata/required_labels.yaml")
	p, iss := parse(srcFile)
	if iss.Err() != nil {
		t.Fatalf("parse() failed: %v", iss.Err())
	}
	if p.name.value != "required_labels" {
		t.Errorf("policy name is %v, wanted 'required_labels'", p.name)
	}
}

func readPolicy(t *testing.T, fileName string) *Source {
	t.Helper()
	tmplBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	return ByteSource(tmplBytes, fileName)
}
