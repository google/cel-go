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

func readPolicyConfig(t testing.TB, fileName string) *Config {
	t.Helper()
	testCaseBytes, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("os.ReadFile(%s) failed: %v", fileName, err)
	}
	config := &Config{}
	err = yaml.Unmarshal(testCaseBytes, config)
	if err != nil {
		log.Fatalf("error: %v", err)
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
	Tests []*TestCase `yaml:"tests"`
}

type TestCase struct {
	Name   string                 `yaml:"name"`
	Input  map[string]interface{} `yaml:"input"`
	Output string                 `yaml:"output"`
}
