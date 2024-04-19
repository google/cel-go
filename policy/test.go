package policy

import (
	"log"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func readPolicy(t testing.TB, fileName string) *Source {
	t.Helper()
	policyBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	return ByteSource(policyBytes, fileName)
}

func readTestSuite(t testing.TB, fileName string) *TestSuite {
	t.Helper()
	testCaseBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	suite := &TestSuite{}
	err = yaml.Unmarshal(testCaseBytes, suite)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return suite
}

type TestSuite struct {
	Description string     `yaml:"description"`
	Sections    []*Section `yaml:"section"`
}

type Section struct {
	Name  string      `yaml:"name"`
	Tests []*TestCase `yaml:"test"`
}

type TestCase struct {
	Input  map[string]interface{} `yaml:"input"`
	Output string                 `yaml:"output"`
}
