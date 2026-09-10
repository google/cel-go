// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package celtest provides compatibility aliases forwarding to cel.dev/cel-go/tools/celtest.
package celtest

import (
	"testing"

	"cel.dev/cel-go/cel"
	"cel.dev/cel-go/common/types/ref"
	canonceltest "cel.dev/cel-go/tools/celtest"
	"github.com/google/cel-go/tools/compiler"
	descpb "google.golang.org/protobuf/types/descriptorpb"
)

// Types
type (
	TestRunnerOption  = canonceltest.TestRunnerOption
	ActivationFactory = canonceltest.ActivationFactory
	TestSuiteParser   = canonceltest.TestSuiteParser
	TestRunner        = canonceltest.TestRunner
	Test              = canonceltest.Test
	TestResult        = canonceltest.TestResult
	Program           = canonceltest.Program
)

// Functions

func UpdateTestResourcesPaths(testResourcesDir string, paths []*string) error {
	return canonceltest.UpdateTestResourcesPaths(testResourcesDir, paths)
}

func TriggerTests(t *testing.T, testRunnerOpts ...TestRunnerOption) {
	canonceltest.TriggerTests(t, testRunnerOpts...)
}

func TestRunnerOptionsFromFlags(testResourcesDir string, testRunnerOpts []TestRunnerOption, testCompilerOpts ...any) TestRunnerOption {
	return canonceltest.TestRunnerOptionsFromFlags(testResourcesDir, testRunnerOpts, testCompilerOpts...)
}

func TestInputActivationFactory(f ActivationFactory) TestRunnerOption {
	return canonceltest.TestInputActivationFactory(f)
}

func TestSuite(path string) TestRunnerOption {
	return canonceltest.TestSuite(path)
}

func DefaultTestSuiteParser(path string) TestRunnerOption {
	return canonceltest.DefaultTestSuiteParser(path)
}

func TestSuiteParserOption(p TestSuiteParser) TestRunnerOption {
	return canonceltest.TestSuiteParserOption(p)
}

func DebugAST(ast *cel.Ast) string {
	return canonceltest.DebugAST(ast)
}

func NewTest(name string, input cel.PartialActivation, resultMatcher func(ref.Val, error) TestResult) *Test {
	return canonceltest.NewTest(name, input, resultMatcher)
}

func NewTestRunner(opts ...TestRunnerOption) (*TestRunner, error) {
	return canonceltest.NewTestRunner(opts...)
}

func TestExpression(value string) TestRunnerOption {
	return canonceltest.TestExpression(value)
}

func TestCompiler(compileOpts ...any) TestRunnerOption {
	return canonceltest.TestCompiler(compileOpts...)
}

func CustomTestCompiler(c compiler.Compiler) TestRunnerOption {
	return canonceltest.CustomTestCompiler(c)
}

func FileDescriptorSet(path string) TestRunnerOption {
	return canonceltest.FileDescriptorSet(path)
}

func AddFileDescriptorSetProto(fds *descpb.FileDescriptorSet) TestRunnerOption {
	return canonceltest.AddFileDescriptorSetProto(fds)
}

func PartialEvalProgramOption() TestRunnerOption {
	return canonceltest.PartialEvalProgramOption()
}

func EnableCoverage() TestRunnerOption {
	return canonceltest.EnableCoverage()
}
