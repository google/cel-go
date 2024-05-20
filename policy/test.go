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
		expr    string
	}{
		{
			name: "nested_rule",
			expr: `
	cel.bind(variables.permitted_regions, ["us", "uk", "es"], 
	  cel.bind(variables.banned_regions, {"us": false, "ru": false, "ir": false}, 
	  (resource.origin in variables.banned_regions &&
		!(resource.origin in variables.permitted_regions)) 
		? optional.of({"banned": true}) : optional.none()).or(
			optional.of((resource.origin in variables.permitted_regions) 
			? {"banned": false} : {"banned": true})))`,
		},
		{
			name: "required_labels",
			expr: `
	cel.bind(variables.want, spec.labels, 
		cel.bind(variables.missing, variables.want.filter(l, !(l in resource.labels)), 
		cel.bind(variables.invalid, 
			resource.labels.filter(l, l in variables.want &&
				variables.want[l] != resource.labels[l]), 
				(variables.missing.size() > 0) 
				? optional.of("missing one or more required labels: %s".format([variables.missing])) 
				: ((variables.invalid.size() > 0) 
				? optional.of("invalid values provided on one or more labels: %s".format([variables.invalid])) : optional.none()))))`,
		},
		{
			name: "restricted_destinations",
			expr: `
	cel.bind(variables.matches_origin_ip, 
	  locationCode(origin.ip) == spec.origin, 
	  cel.bind(variables.has_nationality, has(request.auth.claims.nationality), 
	    cel.bind(variables.matches_nationality, 
		  variables.has_nationality && request.auth.claims.nationality == spec.origin,
		  cel.bind(variables.matches_dest_ip, 
			locationCode(destination.ip) in spec.restricted_destinations, 
			cel.bind(variables.matches_dest_label, 
			  resource.labels.location in spec.restricted_destinations,
              cel.bind(variables.matches_dest, 
				variables.matches_dest_ip || variables.matches_dest_label, 
				(variables.matches_nationality && variables.matches_dest) 
				? true 
				: ((!variables.has_nationality && variables.matches_origin_ip && variables.matches_dest) 
		        ? true : false)))))))`,
			envOpts: []cel.EnvOption{
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
