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
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/test"
	"github.com/google/cel-go/tools/compiler"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/testing/protocmp"

	celpb "cel.dev/expr"
	conformancepb "cel.dev/expr/conformance/test"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	descpb "google.golang.org/protobuf/types/descriptorpb"
	dynamicpb "google.golang.org/protobuf/types/dynamicpb"
)

var (
	celExpression         string
	celExpressionType     int
	expressionType        compiler.ExpressionType
	testSuitePath         string
	fileDescriptorSetPath string
	configPath            string
	baseConfigPath        string
)

func init() {
	flag.StringVar(&testSuitePath, "test_suite_path", "", "path to a test suite")
	flag.StringVar(&fileDescriptorSetPath, "file_descriptor_set", "", "path to a file descriptor set")
	flag.StringVar(&configPath, "config_path", "", "path to a config file")
	flag.StringVar(&baseConfigPath, "base_config_path", "", "path to a base config file")
	flag.StringVar(&celExpression, "cel_expr", "", "CEL expression to test")
	flag.IntVar(&celExpressionType, "cel_expression_type", 0, "type of the CEL expression")
	flag.Parse()
	expressionType = compiler.ExpressionType(celExpressionType)
}

func loadInput(path string, format compiler.FileFormat, out protoreflect.ProtoMessage) error {
	unmarshaller := proto.Unmarshal
	if format == compiler.TextProto {
		unmarshaller = prototext.Unmarshal
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %v", path, err)
	}
	return unmarshaller(data, out)
}

// TestRunnerOption is used to configure the following attributes of the Test Runner:
// - set the Compiler
// - add Input Expressions
// - set the test suite file path
// - set the test suite parser based on the file format: YAML or Textproto
type TestRunnerOption func(*TestRunner) (*TestRunner, error)

// TestSuiteParserTextproto returns a cel.expr.conformance.test.TestSuite
// message which is used to set up Tests.
// - In case the message is serialized in a Textproto file, TestSuiteParserTextproto
// is invoked with the path of the file passed as an argument.
// - Alternatively, TestSuiteParserTextproto can also be invoked to generate the
// cel.expr.conformance.test.TestSuite message at runtime.
type TestSuiteParserTextproto func(*testing.T, any) (*conformancepb.TestSuite, error)

// TestSuiteParserYAML returns a test.TestSuite object which is used to set up Tests.
// - In case the object is serialized in a YAML file, TestSuiteParserYAML is invoked
// with the path of the file passed as an argument.
// - Alternatively, TestSuiteParserYAML can also be invoked to generate the test.TestSuite
// object at runtime.
//
// Note: If TestSuiteParserTextproto is already configured then TestSuiteParserYAML
// will not be invoked.
type TestSuiteParserYAML func(*testing.T, any) (*test.Suite, error)

// TriggerTests triggers tests for a CEL policy, expression or checked expression
// with the provided set of options. The options can be used to:
// - configure the Compiler used for parsing and compiling the expression
// - configure the Test Runner used for parsing and executing the tests
func TriggerTests(t *testing.T, opts ...TestRunnerOption) {
	compilerOpt := testRunnerCompilerFromFlags()
	testSuiteParserOpt := TestSuiteParser(testSuitePath)
	fileDescriptorSetOpt := AddFileDescriptorSet(fileDescriptorSetPath)
	testRunnerExprOpt := testRunnerExpressionsFromFlags()
	opts = append([]TestRunnerOption{compilerOpt, testSuiteParserOpt, fileDescriptorSetOpt, testRunnerExprOpt}, opts...)
	tr, err := NewTestRunner(opts...)
	if err != nil {
		t.Fatalf("error creating test runner: %v", err)
	}
	programs, err := tr.Programs(t)
	if err != nil {
		t.Fatalf("error creating programs: %v", err)
	}
	tests, err := tr.Tests(t)
	if err != nil {
		t.Fatalf("error creating tests: %v", err)
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := tr.ExecuteTest(t, programs, test)
			if err != nil {
				t.Fatalf("error executing test: %v", err)
			}
		})
	}
}

func testRunnerCompilerFromFlags() TestRunnerOption {
	var opts []any
	if fileDescriptorSetPath != "" {
		opts = append(opts, compiler.TypeDescriptorSetFile(fileDescriptorSetPath))
	}
	if baseConfigPath != "" {
		opts = append(opts, compiler.EnvironmentFile(baseConfigPath))
	}
	if configPath != "" {
		opts = append(opts, compiler.EnvironmentFile(configPath))
	}
	return func(tr *TestRunner) (*TestRunner, error) {
		c, err := compiler.NewCompiler(opts...)
		if err != nil {
			return nil, err
		}
		tr.Compiler = c
		return tr, nil
	}
}

func testRunnerExpressionsFromFlags() TestRunnerOption {
	return func(tr *TestRunner) (*TestRunner, error) {
		if celExpression != "" {
			switch expressionType {
			case compiler.CompiledExpressionFile:
				tr.Expressions = append(tr.Expressions, &compiler.CompiledExpression{Path: celExpression})
			case compiler.PolicyFile, compiler.ExpressionFile:
				tr.Expressions = append(tr.Expressions, &compiler.FileExpression{Path: celExpression})
			case compiler.RawExpressionString:
				tr.Expressions = append(tr.Expressions, &compiler.RawExpression{Value: celExpression})
			default:
				return nil, fmt.Errorf("unsupported expression type: %v", expressionType)
			}
		}
		return tr, nil
	}
}

// TestSuiteParser provides a TestRunnerOption that sets the test suite file path and the test
// suite parser based on the file format: YAML or Textproto.
func TestSuiteParser(path string) TestRunnerOption {
	return func(tr *TestRunner) (*TestRunner, error) {
		if testSuitePath == "" {
			return tr, nil
		}
		tr.TestSuiteFilePath = path
		testSuiteFormat := compiler.InferFileFormat(testSuitePath)
		switch testSuiteFormat {
		case compiler.TextProto:
			tr.TestSuiteParserTextproto = defaultTestSuiteParserTextproto
		case compiler.TextYAML:
			tr.TestSuiteParserYAML = defaultTestSuiteParserYAML
		default:
			return nil, fmt.Errorf("unsupported test suite file format: %v", testSuiteFormat)
		}
		return tr, nil
	}
}

// TestRunner provides a structure to hold the different components required to execute tests for
// a list of Input Expressions. The TestRunner can be configured with the following options:
// - Compiler: The compiler used for parsing and compiling the input expressions.
// - Input Expressions: The list of input expressions to be tested.
// - Test Suite File Path: The path to the test suite file.
// - File Descriptor Set Path: The path to the file descriptor set file.
// - Test Suite Parser Textproto: A parser for a custom test suite file serialized in Textproto
// format. This option is required if the provided test suite is not a cel.spec.expr.conformance.test.TestSuite message.
// - Test Suite Parser YAML: A parser for a custom test suite file serialized in YAML format.
// This option is required if the provided test suite is not a test.TestSuite object.
//
// The TestRunner provides the following methods:
// - Programs: Creates a list of CEL programs from the input expressions.
// - Tests: Creates a list of tests from the test suite file.
// - ExecuteTest: Executes a single
type TestRunner struct {
	compiler.Compiler
	Expressions              []compiler.InputExpression
	TestSuiteFilePath        string
	FileDescriptorSetPath    string
	TestSuiteParserTextproto TestSuiteParserTextproto
	TestSuiteParserYAML      TestSuiteParserYAML
}

// Test represents a single test case to be executed. It encompasses the following:
// - name: The name of the test case.
// - input: The input to be used for evaluating the CEL expression.
// - resultMatcher: A function that takes in the result of evaluating the CEL expression and
// returns a TestResult.
type Test struct {
	name          string
	input         interpreter.Activation
	resultMatcher func(ref.Val, error) TestResult
}

// NewTest creates a new Test with the provided name, input and result matcher.
func NewTest(name string, input interpreter.Activation, resultMatcher func(ref.Val, error) TestResult) *Test {
	return &Test{
		name:          name,
		input:         input,
		resultMatcher: resultMatcher,
	}
}

// TestResult represents the result of a test case execution. It contains the validation result
// along with the expected result and any errors encountered during the execution.
// - Success: Whether the result matcher condition validating the test case was satisfied.
// - Wanted: The expected result of the test case.
// - Error: Any error encountered during the execution.
type TestResult struct {
	Success bool
	Wanted  string
	Error   error
}

// NewTestRunner creates a Test Runner with the provided options.
// The options can be used to:
// - configure the Compiler used for parsing and compiling the input expressions
// - configure the Test Runner used for parsing and executing the tests
func NewTestRunner(opts ...TestRunnerOption) (*TestRunner, error) {
	tr := &TestRunner{}
	var err error
	for _, opt := range opts {
		tr, err = opt(tr)
		if err != nil {
			return nil, err
		}
	}
	return tr, nil
}

// AddFileDescriptorSet creates a Test Runner Option which adds a file descriptor set to the test
// runner. The file descriptor set is used to register proto messages in the global proto registry.
func AddFileDescriptorSet(path string) TestRunnerOption {
	return func(tr *TestRunner) (*TestRunner, error) {
		if path != "" {
			tr.FileDescriptorSetPath = path
		}
		return tr, nil
	}
}

func registerMessages(path string) error {
	if path == "" {
		return nil
	}
	fds, err := fileDescriptorSet(path)
	if err != nil {
		return err
	}
	for _, file := range fds.GetFile() {
		reflectFD, err := protodesc.NewFile(file, protoregistry.GlobalFiles)
		if err != nil {
			return fmt.Errorf("protodesc.NewFile(%q) failed: %v", file.GetName(), err)
		}
		if _, err := protoregistry.GlobalFiles.FindFileByPath(reflectFD.Path()); err == nil {
			continue
		}
		err = protoregistry.GlobalFiles.RegisterFile(reflectFD)
		if err != nil {
			return fmt.Errorf("protoregistry.GlobalFiles.RegisterFile() failed: %v", err)
		}
		for i := 0; i < reflectFD.Messages().Len(); i++ {
			msg := reflectFD.Messages().Get(i)
			msgType := dynamicpb.NewMessageType(msg)
			err = protoregistry.GlobalTypes.RegisterMessage(msgType)
			if err != nil && !strings.Contains(err.Error(), "already registered") {
				return fmt.Errorf("protoregistry.GlobalTypes.RegisterMessage(%q) failed: %v", msgType, err)
			}
		}
	}
	return nil
}

func fileDescriptorSet(path string) (*descpb.FileDescriptorSet, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file descriptor set file %q: %v", fileDescriptorSetPath, err)
	}
	fds := &descpb.FileDescriptorSet{}
	if err := proto.Unmarshal(bytes, fds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal file descriptor set file %q: %v", fileDescriptorSetPath, err)
	}
	return fds, nil
}

func defaultTestSuiteParserTextproto(t *testing.T, path any) (*conformancepb.TestSuite, error) {
	t.Helper()
	testSuite := &conformancepb.TestSuite{}
	err := loadInput(path.(string), compiler.TextProto, testSuite)
	return testSuite, err
}

func defaultTestSuiteParserYAML(t *testing.T, path any) (*test.Suite, error) {
	t.Helper()
	testSuiteBytes, err := os.ReadFile(path.(string))
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(%q) failed: %v", path.(string), err)
	}
	testSuite := &test.Suite{}
	err = yaml.Unmarshal(testSuiteBytes, testSuite)
	return testSuite, err
}

// Programs creates a list of CEL programs from the input expressions configured in the test runner
// using the provided program options.
func (tr *TestRunner) Programs(t *testing.T, opts ...cel.ProgramOption) ([]cel.Program, error) {
	t.Helper()
	if tr.Compiler == nil {
		return nil, fmt.Errorf("compiler is not set")
	}
	e, err := tr.CreateEnv()
	if err != nil {
		return nil, err
	}
	var programs []cel.Program
	for _, expr := range tr.Expressions {
		ast, _, err := expr.CreateAST(tr.Compiler)
		if err != nil {
			return nil, err
		}
		prg, err := e.Program(ast, opts...)
		if err != nil {
			return nil, err
		}
		programs = append(programs, prg)
	}
	return programs, nil
}

// Tests creates a list of tests from the test suite file and test suite parser configured in the
// test runner.
func (tr *TestRunner) Tests(t *testing.T) ([]*Test, error) {
	if tr.Compiler == nil {
		return nil, fmt.Errorf("compiler is not set")
	}
	if tr.TestSuiteParserTextproto != nil {
		err := registerMessages(tr.FileDescriptorSetPath)
		if err != nil {
			return nil, fmt.Errorf("registerMessages(%q) failed: %v", tr.FileDescriptorSetPath, err)
		}
		testSuite, err := tr.TestSuiteParserTextproto(t, tr.TestSuiteFilePath)
		if err != nil {
			return nil, err
		}
		return tr.createTestsFromTextproto(t, testSuite)
	}
	if tr.TestSuiteParserYAML != nil {
		testSuite, err := tr.TestSuiteParserYAML(t, tr.TestSuiteFilePath)
		if err != nil {
			return nil, err
		}
		return tr.createTestsFromYAML(t, testSuite)
	}
	return nil, nil
}

func (tr *TestRunner) createTestsFromTextproto(t *testing.T, testSuite *conformancepb.TestSuite) ([]*Test, error) {
	var tests []*Test
	for _, section := range testSuite.GetSections() {
		sectionName := section.GetName()
		for _, testCase := range section.GetTests() {
			testName := fmt.Sprintf("%s/%s", sectionName, testCase.GetName())
			testInput, err := tr.createTestInputFromPB(t, testCase)
			if err != nil {
				return nil, err
			}
			testResultMatcher, err := tr.createResultMatcherFromPB(t, testCase)
			if err != nil {
				return nil, err
			}
			tests = append(tests, NewTest(testName, testInput, testResultMatcher))
		}
	}
	return tests, nil
}

func (tr *TestRunner) createTestInputFromPB(t *testing.T, testCase *conformancepb.TestCase) (interpreter.Activation, error) {
	t.Helper()
	input := map[string]any{}
	e, err := tr.CreateEnv()
	if err != nil {
		return nil, err
	}
	if testCase.GetInputContext() != nil {
		if len(testCase.GetInput()) != 0 {
			return nil, fmt.Errorf("only one of input and input_context can be provided at a time")
		}
		switch testInput := testCase.GetInputContext().GetInputContextKind().(type) {
		case *conformancepb.InputContext_ContextExpr:
			refVal, err := tr.eval(testInput.ContextExpr)
			if err != nil {
				return nil, fmt.Errorf("eval(%q) failed: %w", testInput.ContextExpr, err)
			}
			ctx, err := refVal.ConvertToNative(
				reflect.TypeOf((*proto.Message)(nil)).Elem())
			if err != nil {
				return nil, fmt.Errorf("context variable is not a valid proto: %w", err)
			}
			return cel.ContextProtoVars(ctx.(proto.Message))
		case *conformancepb.InputContext_ContextMessage:
			refVal := e.CELTypeAdapter().NativeToValue(testInput.ContextMessage)
			ctx, err := refVal.ConvertToNative(reflect.TypeOf((*proto.Message)(nil)).Elem())
			if err != nil {
				return nil, fmt.Errorf("context variable is not a valid proto: %w", err)
			}
			return cel.ContextProtoVars(ctx.(proto.Message))
		}
	}
	for k, v := range testCase.GetInput() {
		switch v.GetKind().(type) {
		case *conformancepb.InputValue_Value:
			input[k], err = cel.ProtoAsValue(e.CELTypeAdapter(), v.GetValue())
			if err != nil {
				return nil, fmt.Errorf("cel.ProtoAsValue(%q) failed: %w", v, err)
			}
		case *conformancepb.InputValue_Expr:
			input[k], err = tr.eval(v.GetExpr())
			if err != nil {
				return nil, fmt.Errorf("eval(%q) failed: %w", v.GetExpr(), err)
			}
		}
	}
	return interpreter.NewActivation(input)
}

func (tr *TestRunner) createResultMatcherFromPB(t *testing.T, testCase *conformancepb.TestCase) (func(ref.Val, error) TestResult, error) {
	t.Helper()
	if testCase.GetOutput() == nil {
		return nil, fmt.Errorf("expected output is nil")
	}
	successResult := TestResult{Success: true}
	e, err := tr.CreateEnv()
	if err != nil {
		return nil, err
	}
	switch testOutput := testCase.GetOutput().GetResultKind().(type) {
	case *conformancepb.TestOutput_ResultValue:
		return func(val ref.Val, err error) TestResult {
			want := e.CELTypeAdapter().NativeToValue(testOutput.ResultValue)
			if err != nil {
				return TestResult{Success: false, Wanted: fmt.Sprintf("simple value %v", want), Error: err}
			}
			outputVal, err := refValueToExprValue(val)
			if err != nil {
				return TestResult{Success: false, Wanted: fmt.Sprintf("simple value %v", want), Error: fmt.Errorf("refValueToExprValue(%q) failed: %v", val, err)}
			}
			testResultVal, err := canonicalValueToV1Alpha1(testOutput.ResultValue)
			if err != nil {
				return TestResult{Success: false, Wanted: fmt.Sprintf("simple value %v", want), Error: fmt.Errorf("canonicalValueToV1Alpha1(%q) failed: %v", testOutput.ResultValue, err)}
			}
			testVal := &exprpb.ExprValue{
				Kind: &exprpb.ExprValue_Value{Value: testResultVal}}

			if diff := cmp.Diff(testVal, outputVal, protocmp.Transform(),
				protocmp.SortRepeatedFields(&exprpb.MapValue{}, "entries")); diff != "" {
				return TestResult{Success: false, Wanted: fmt.Sprintf("simple value %v", want), Error: fmt.Errorf("mismatched test output with diff (-want +got):\n%s", diff)}
			}
			return successResult
		}, nil
	case *conformancepb.TestOutput_ResultExpr:
		return func(val ref.Val, err error) TestResult {
			if err != nil {
				return TestResult{Success: false, Error: err}
			}
			testOut, err := tr.eval(testOutput.ResultExpr)
			if err != nil {
				return TestResult{Success: false, Error: fmt.Errorf("eval(%q) failed: %v", testOutput.ResultExpr, err)}
			}
			if optOut, ok := val.(*types.Optional); ok {
				if optOut.Equal(types.OptionalNone) == types.True {
					if testOut.Equal(types.OptionalNone) != types.True {
						return TestResult{Success: false, Wanted: fmt.Sprintf("optional value %v", testOut), Error: fmt.Errorf("policy eval got %v", val)}
					}
				} else if testOut.Equal(optOut.GetValue()) != types.True {
					return TestResult{Success: false, Wanted: fmt.Sprintf("optional value %v", testOut), Error: fmt.Errorf("policy eval got %v", val)}
				}
			} else if val.Equal(testOut) != types.True {
				return TestResult{Success: false, Wanted: fmt.Sprintf("optional value %v", testOut), Error: fmt.Errorf("policy eval got %v", val)}
			}
			return successResult
		}, nil
	case *conformancepb.TestOutput_EvalError:
		return func(val ref.Val, err error) TestResult {
			failureResult := TestResult{Success: false, Wanted: fmt.Sprintf("error %v", testOutput.EvalError)}
			if err == nil {
				return failureResult
			}
			// Compare the evaluated error with the expected error message only.
			for _, want := range testOutput.EvalError.GetErrors() {
				if strings.Contains(err.Error(), want.GetMessage()) {
					return successResult
				}
			}
			return failureResult
		}, nil
	case *conformancepb.TestOutput_Unknown:
		//  TODO: to implement
	}
	return nil, nil
}

func refValueToExprValue(refVal ref.Val) (*exprpb.ExprValue, error) {
	if types.IsUnknown(refVal) {
		return &exprpb.ExprValue{
			Kind: &exprpb.ExprValue_Unknown{
				Unknown: &exprpb.UnknownSet{
					Exprs: refVal.Value().([]int64),
				},
			}}, nil
	}
	v, err := cel.RefValueToValue(refVal)
	if err != nil {
		return nil, err
	}
	return &exprpb.ExprValue{
		Kind: &exprpb.ExprValue_Value{Value: v}}, nil
}

func canonicalValueToV1Alpha1(val *celpb.Value) (*exprpb.Value, error) {
	var v1val exprpb.Value
	b, err := prototext.Marshal(val)
	if err != nil {
		return nil, err
	}
	if err := prototext.Unmarshal(b, &v1val); err != nil {
		return nil, err
	}
	return &v1val, nil
}

func (tr *TestRunner) eval(expr string) (ref.Val, error) {
	e, err := tr.CreateEnv()
	if err != nil {
		return nil, err
	}
	e, err = e.Extend(cel.OptionalTypes())
	if err != nil {
		return nil, fmt.Errorf("e.Extend() failed: %v", err)
	}
	ast, iss := e.Compile(expr)
	if iss.Err() != nil {
		return nil, fmt.Errorf("e.Compile(%q) failed: %v", expr, iss.Err())
	}
	prg, err := e.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("e.Program(%q) failed: %v", expr, err)
	}
	out, _, err := prg.Eval(cel.NoVars())
	if err != nil {
		return nil, fmt.Errorf("prg.Eval(%q) failed: %v", expr, err)
	}
	return out, nil
}

func (tr *TestRunner) createTestsFromYAML(t *testing.T, testSuite *test.Suite) ([]*Test, error) {
	var tests []*Test
	for _, section := range testSuite.Sections {
		for _, testCase := range section.Tests {
			testName := fmt.Sprintf("%s/%s", section.Name, testCase.Name)
			testInput, err := tr.createTestInput(t, testCase)
			if err != nil {
				return nil, err
			}
			testResultMatcher, err := tr.createResultMatcher(t, testCase.Output)
			if err != nil {
				return nil, err
			}
			tests = append(tests, NewTest(testName, testInput, testResultMatcher))
		}
	}
	return tests, nil
}

func (tr *TestRunner) createTestInput(t *testing.T, testCase *test.Case) (interpreter.Activation, error) {
	t.Helper()
	if testCase.InputContext != nil && testCase.InputContext.ContextExpr != "" {
		if len(testCase.Input) != 0 {
			return nil, fmt.Errorf("only one of input and input_context can be provided at a time")
		}
		contextExpr := testCase.InputContext.ContextExpr
		out, err := tr.eval(contextExpr)
		if err != nil {
			return nil, fmt.Errorf("eval(%q) failed: %w", contextExpr, err)
		}
		ctx, err := out.ConvertToNative(reflect.TypeOf((*proto.Message)(nil)).Elem())
		if err != nil {
			return nil, fmt.Errorf("context variable is not a valid proto: %w", err)
		}
		return cel.ContextProtoVars(ctx.(proto.Message))
	}
	input := map[string]any{}
	for k, v := range testCase.Input {
		if v.Expr != "" {
			val, err := tr.eval(v.Expr)
			if err != nil {
				return nil, fmt.Errorf("eval(%q) failed: %w", v.Expr, err)
			}
			input[k] = val
			continue
		}
		input[k] = v.Value
	}
	return interpreter.NewActivation(input)
}

func (tr *TestRunner) createResultMatcher(t *testing.T, testOutput *test.Output) (func(ref.Val, error) TestResult, error) {
	t.Helper()
	e, err := tr.CreateEnv()
	if err != nil {
		return nil, err
	}
	successResult := TestResult{Success: true}
	if testOutput.Value != nil {
		want := e.CELTypeAdapter().NativeToValue(testOutput.Value)
		return func(out ref.Val, err error) TestResult {
			if err == nil {
				if out.Equal(want) == types.True {
					return successResult
				}
				if optOut, ok := out.(*types.Optional); ok {
					if optOut.HasValue() && optOut.GetValue().Equal(want) == types.True {
						return successResult
					}
				}
			}
			return TestResult{Success: false, Wanted: fmt.Sprintf("simple value %v", want), Error: err}
		}, nil
	}
	if testOutput.Expr != "" {
		want, err := tr.eval(testOutput.Expr)
		if err != nil {
			return nil, fmt.Errorf("eval(%q) failed: %w", testOutput.Expr, err)
		}
		return func(out ref.Val, err error) TestResult {
			if err == nil {
				if out.Equal(want) == types.True {
					return successResult
				}
				if optOut, ok := out.(*types.Optional); ok {
					if optOut.HasValue() && optOut.GetValue().Equal(want) == types.True {
						return successResult
					}
				}
			}
			return TestResult{Success: false, Wanted: fmt.Sprintf("simple value %v", want), Error: err}
		}, nil
	}
	if testOutput.ErrorSet != nil {
		return func(out ref.Val, err error) TestResult {
			failureResult := TestResult{Success: false, Wanted: fmt.Sprintf("error %v", testOutput.ErrorSet)}
			if err == nil {
				return failureResult
			}
			for _, want := range testOutput.ErrorSet {
				if strings.Contains(err.Error(), want) {
					return successResult
				}
			}
			return failureResult
		}, nil
	}
	if testOutput.UnknownSet != nil {
		//  TODO: to implement
	}
	return nil, nil
}

// ExecuteTest executes the test case against the provided list of programs and returns an error if
// the test fails.
func (tr *TestRunner) ExecuteTest(t *testing.T, programs []cel.Program, test *Test) error {
	t.Helper()
	if tr.Compiler == nil {
		return fmt.Errorf("compiler is not set")
	}
	for _, program := range programs {
		out, _, err := program.Eval(test.input)
		if testResult := test.resultMatcher(out, err); !testResult.Success {
			return fmt.Errorf("test: %s \n wanted: %v \n failed: %v", test.name, testResult.Wanted, testResult.Error)
		}
	}
	return nil
}
