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
	"log"
	"os"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	"gopkg.in/yaml.v3"
)

var (
	policyTests = []struct {
		name    string
		envOpts []cel.EnvOption
	}{
		{name: "nested_rule"},
		{name: "required_labels"},
		{name: "restricted_destinations", envOpts: []cel.EnvOption{
			cel.Function("locationCode",
				cel.Overload("locationCode_string", []*cel.Type{cel.StringType}, cel.StringType,
					cel.UnaryBinding(func(ip ref.Val) ref.Val {
						switch ip.(types.String) {
						case types.String("10.0.0.1"):
							return types.String("us")
						case types.String("10.0.0.2"):
							return types.String("de")
						default:
							return types.String("ir")
						}
					}))),
		}},
	}
)

func readPolicy(t testing.TB, fileName string) *Source {
	t.Helper()
	policyBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	return ByteSource(policyBytes, fileName)
}

func readPolicyConfig(t testing.TB, fileName string) *Config {
	t.Helper()
	testCaseBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	config := &Config{}
	err = yaml.Unmarshal(testCaseBytes, config)
	if err != nil {
		log.Fatalf("yaml.Unmarshal(%s) error: %v", fileName, err)
	}
	return config
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
		log.Fatalf("yaml.Unmarshal(%s) error: %v", fileName, err)
	}
	return suite
}

// TestSuite describes a set of tests divided by section.
type TestSuite struct {
	Description string         `yaml:"description"`
	Sections    []*TestSection `yaml:"section"`
}

// TestSection describes a related set of tests associated with a behavior.
type TestSection struct {
	Name  string      `yaml:"name"`
	Tests []*TestCase `yaml:"tests"`
}

// TestCase describes a named test scenario with a set of inputs and expected outputs.
//
// Note, when a test requires additional functions to be provided to execute, the test harness
// must supply these functions.
type TestCase struct {
	Name   string                 `yaml:"name"`
	Input  map[string]interface{} `yaml:"input"`
	Output string                 `yaml:"output"`
}
