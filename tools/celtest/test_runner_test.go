// Copyright 2025 Google LLC
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

// Package celtest provides functions for testing CEL policies and expressions.
package celtest

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/policy"
	"github.com/google/cel-go/tools/compiler"
)

type testCase struct {
	name                  string
	celExpr               string
	testSuitePath         string
	fileDescriptorSetPath string
	configPath            string
	opts                  []any
	exitStatus            int
}

func setupTests(t *testing.T) []*testCase {
	t.Helper()
	testCases := []*testCase{
		{
			name:       "unspecified cel expression type",
			celExpr:    "/invalid_path",
			exitStatus: 1,
		},
		{
			name:       "invalid config file path",
			celExpr:    "../../policy/testdata/restricted_destinations/policy.yaml",
			configPath: "/invalid_path.yaml",
			exitStatus: 1,
		},
		{
			name:       "invalid config file",
			celExpr:    "../../policy/testdata/restricted_destinations/policy.yaml",
			configPath: "testdata/invalid_config.yaml",
			exitStatus: 1,
		},
		{
			name:       "invalid context message config file",
			celExpr:    "../../policy/testdata/restricted_destinations/policy.yaml",
			configPath: "testdata/invalid_context_message_config.yaml",
			exitStatus: 1,
		},
		{
			name:          "invalid test suite path",
			celExpr:       "../../policy/testdata/restricted_destinations/policy.yaml",
			testSuitePath: "/invalid path",
			configPath:    "../../policy/testdata/restricted_destinations/config.yaml",
			exitStatus:    0,
		},
		{
			name:          "invalid test missing function decl",
			celExpr:       "../../policy/testdata/restricted_destinations/policy.yaml",
			testSuitePath: "testdata/tests_missing_function.yaml",
			configPath:    "../../policy/testdata/restricted_destinations/config.yaml",
			exitStatus:    0,
		},
		{
			name:          "invalid test missing variable decl",
			celExpr:       "../../policy/testdata/restricted_destinations/policy.yaml",
			testSuitePath: "testdata/invalid_tests_missing_variable.yaml",
			configPath:    "../../policy/testdata/restricted_destinations/config.yaml",
			opts:          []any{locationCodeEnvOption()},
			exitStatus:    0,
		},
		{
			name:          "invalid test output value mismatch",
			celExpr:       "../../policy/testdata/restricted_destinations/policy.yaml",
			testSuitePath: "testdata/invalid_tests_output_value_mismatch.yaml",
			configPath:    "../../policy/testdata/restricted_destinations/config.yaml",
			opts:          []any{locationCodeEnvOption()},
			exitStatus:    1,
		},
		{
			name:          "invalid test expression output",
			celExpr:       "../../policy/testdata/restricted_destinations/policy.yaml",
			testSuitePath: "testdata/invalid_tests_expression_output.yaml",
			configPath:    "../../policy/testdata/restricted_destinations/config.yaml",
			opts:          []any{locationCodeEnvOption()},
			exitStatus:    1,
		},
		{
			name:          "valid checked expression",
			celExpr:       "../../policy/testdata/restricted_destinations/policy.yaml",
			testSuitePath: "../../policy/testdata/restricted_destinations/tests.yaml",
			configPath:    "../../policy/testdata/restricted_destinations/config.yaml",
			opts:          []any{locationCodeEnvOption()},
			exitStatus:    0,
		},
	}
	return testCases
}

func locationCodeEnvOption() cel.EnvOption {
	return cel.Function("locationCode",
		cel.Overload("locationCode_string", []*cel.Type{cel.StringType}, cel.StringType,
			cel.UnaryBinding(locationCode)))
}

func locationCode(ip ref.Val) ref.Val {
	switch ip.(types.String) {
	case "10.0.0.1":
		return types.String("us")
	case "10.0.0.2":
		return types.String("de")
	default:
		return types.String("ir")
	}
}

func k8sParserOpts() policy.ParserOption {
	return policy.ParserOption(func(p *policy.Parser) (*policy.Parser, error) {
		p.TagVisitor = policy.K8sTestTagHandler()
		return p, nil
	})
}

// TestTriggerTestsCustomPolicy tests the TriggerTestsFromCompiler function for a custom policy
// by providing test runner and compiler options without setting the flag variables.
func TestTriggerTestsCustomPolicy(t *testing.T) {
	t.Run("test trigger tests custom policy", func(t *testing.T) {
		envOpt := compiler.EnvironmentFile("../../policy/testdata/k8s/config.yaml")
		testSuiteParser := DefaultTestSuiteParser("../../policy/testdata/k8s/tests.yaml")
		testCELPolicy := TestRunnerOption(func(tr *TestRunner) (*TestRunner, error) {
			tr.Expressions = append(tr.Expressions, &compiler.FileExpression{
				Path: "../../policy/testdata/k8s/policy.yaml",
			})
			return tr, nil
		})
		c, err := compiler.NewCompiler(envOpt, k8sParserOpts())
		if err != nil {
			t.Fatalf("compiler.NewCompiler() failed: %v", err)
		}
		compilerOpt := TestRunnerOption(func(tr *TestRunner) (*TestRunner, error) {
			tr.Compiler = c
			return tr, nil
		})
		opts := []TestRunnerOption{compilerOpt, testSuiteParser, testCELPolicy}
		TriggerTests(t, opts)
	})
}

// TestTriggerTests tests different scenarios of the TriggerTestsFromCompiler function.
func TestTriggerTests(t *testing.T) {
	celExpression = "../../policy/testdata/restricted_destinations/policy.yaml"
	testSuitePath = "../../policy/testdata/restricted_destinations/tests.yaml"
	configPath = "../../policy/testdata/restricted_destinations/config.yaml"
	TriggerTests(t, nil, locationCodeEnvOption())
}
