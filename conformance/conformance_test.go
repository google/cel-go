package conformance_test

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/runfiles"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/ext"
	"github.com/google/go-cmp/cmp"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/testing/protocmp"

	test2pb "cel.dev/expr/proto/test/v1/proto2/test_all_types"
	test3pb "cel.dev/expr/proto/test/v1/proto3/test_all_types"
	testpb "cel.dev/expr/proto/test/v1/testpb"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type testsFlag []string

func (t *testsFlag) String() string {
	return strings.Join(([]string)(*t), ",")
}

func (t *testsFlag) Set(v string) error {
	*t = strings.Split(v, ",")
	for i, v := range *t {
		(*t)[i] = strings.TrimSpace(v)
	}
	return nil
}

func (t *testsFlag) Get() any {
	return ([]string)(*t)
}

type skipTestsFlag []string

func (t *skipTestsFlag) String() string {
	return strings.Join(([]string)(*t), ",")
}

func (t *skipTestsFlag) Set(v string) error {
	*t = strings.Split(v, ",")
	for i, v := range *t {
		(*t)[i] = strings.TrimSpace(v)
	}
	return nil
}

func (t *skipTestsFlag) Get() any {
	return ([]string)(*t)
}

var (
	dashboard bool
	tests     testsFlag
	skipTests skipTestsFlag

	envWithMacros *cel.Env
	envNoMacros   *cel.Env
)

func init() {
	flag.BoolVar(&dashboard, "dashboard", false, "Dashboard.")
	flag.Var(&tests, "tests", "Paths to run, separate by a comma.")
	flag.Var(&skipTests, "skip_tests", "Tests to skip, separate by a comma.")

	stdOpts := []cel.EnvOption{
		cel.StdLib(),
		cel.ClearMacros(),
		cel.OptionalTypes(),
		cel.EagerlyValidateDeclarations(true),
		cel.EnableErrorOnBadPresenceTest(true),
		cel.Types(&test2pb.TestAllTypes{}, &test2pb.Proto2ExtensionScopedMessage{}, &test3pb.TestAllTypes{}),
		ext.Bindings(),
		ext.Encoders(),
		ext.Math(),
		ext.Protos(),
		ext.Strings(),
	}

	var err error
	envNoMacros, err = cel.NewCustomEnv(stdOpts...)
	if err != nil {
		log.Fatalf("cel.NewCustomEnv() = %v", err)
	}
	envWithMacros, err = envNoMacros.Extend(cel.Macros(cel.StandardMacros...))
	if err != nil {
		log.Fatalf("cel.NewCustomEnv() = %v", err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse()
	code := m.Run()
	if dashboard {
		code = 0
	}
	os.Exit(code)
}

func shouldSkipTest(s string) bool {
	for _, t := range skipTests {
		if strings.HasPrefix(s, t) {
			n := s[len(t):]
			if n == "" || strings.HasPrefix(n, "/") {
				return true
			}
		}
	}
	return false
}

func refValueToExprValue(res ref.Val) (*exprpb.ExprValue, error) {
	if types.IsUnknown(res) {
		return &exprpb.ExprValue{
			Kind: &exprpb.ExprValue_Unknown{
				Unknown: &exprpb.UnknownSet{
					Exprs: res.Value().([]int64),
				},
			}}, nil
	}
	v, err := cel.RefValueToValue(res)
	if err != nil {
		return nil, err
	}
	return &exprpb.ExprValue{
		Kind: &exprpb.ExprValue_Value{Value: v}}, nil
}

func exprValueToRefValue(adapter types.Adapter, ev *exprpb.ExprValue) (ref.Val, error) {
	switch ev.Kind.(type) {
	case *exprpb.ExprValue_Value:
		return cel.ValueToRefValue(adapter, ev.GetValue())
	case *exprpb.ExprValue_Error:
		// An error ExprValue is a repeated set of statuspb.Status
		// messages, with no convention for the status details.
		// To convert this to a types.Err, we need to convert
		// these Status messages to a single string, and be
		// able to decompose that string on output so we can
		// round-trip arbitrary ExprValue messages.
		// TODO(jimlarson) make a convention for this.
		return types.NewErr("XXX add details later"), nil
	case *exprpb.ExprValue_Unknown:
		var unk *types.Unknown
		for _, id := range ev.GetUnknown().GetExprs() {
			if unk == nil {
				unk = types.NewUnknown(id, nil)
			}
			unk = types.MergeUnknowns(types.NewUnknown(id, nil), unk)
		}
		return unk, nil
	}
	return nil, errors.New("unknown ExprValue kind")
}

func conformanceTest(t *testing.T, name string, pb *testpb.SimpleTest) {
	if shouldSkipTest(name) {
		t.SkipNow()
		return
	}
	var env *cel.Env
	if pb.GetDisableMacros() {
		env = envNoMacros
	} else {
		env = envWithMacros
	}
	src := common.NewStringSource(pb.GetExpr(), pb.GetName())
	ast, iss := env.ParseSource(src)
	if err := iss.Err(); err != nil {
		t.Fatal(err)
	}
	var opts []cel.EnvOption
	if pb.GetContainer() != "" {
		opts = append(opts, cel.Container(pb.GetContainer()))
	}
	opts = append(opts, cel.Declarations(pb.GetTypeEnv()...))
	var err error
	env, err = env.Extend(opts...)
	if err != nil {
		t.Fatal(err)
	}
	if !pb.GetDisableCheck() {
		ast, iss = env.Check(ast)
		if err := iss.Err(); err != nil {
			t.Fatal(err)
		}
	}
	program, err := env.Program(ast)
	if err != nil {
		t.Fatal(err)
	}
	act := make(map[string]any, len(pb.GetBindings()))
	for k, v := range pb.GetBindings() {
		act[k], err = exprValueToRefValue(env.CELTypeAdapter(), v)
		if err != nil {
			t.Fatal(err)
		}
	}
	ret, _, err := program.Eval(act)
	switch m := pb.GetResultMatcher().(type) {
	case *testpb.SimpleTest_Value:
		if err != nil {
			t.Fatalf("program.Eval(): got %v, want nil", err)
		}
		val, err := refValueToExprValue(ret)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(&exprpb.ExprValue{Kind: &exprpb.ExprValue_Value{Value: m.Value}}, val, protocmp.Transform(), protocmp.SortRepeatedFields(&exprpb.MapValue{}, "entries")); diff != "" {
			t.Errorf("program.Eval() diff (-want +got):\n%s", diff)
		}
	case *testpb.SimpleTest_EvalError:
		if err == nil && types.IsError(ret) {
			err = ret.(*types.Err).Unwrap()
		}
		if err == nil {
			t.Errorf("program.Eval(): got nil, want %v", m.EvalError)
		}
	default:
		t.Errorf("unexpected matcher kind: %T", pb.GetResultMatcher())
	}
}

func TestConformance(t *testing.T) {
	var files []*testpb.SimpleTestFile
	for _, path := range tests {
		path = strings.TrimPrefix(path, "external/")
		f, err := runfiles.RlocationFrom(path, "com_google_cel_spec")
		if err != nil {
			log.Fatalf("failed to find runfile %q: %v", path, err)
		}
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatalf("failed to read file %q: %v", f, err)
		}
		file := &testpb.SimpleTestFile{}
		err = prototext.Unmarshal(b, file)
		if err != nil {
			log.Fatalf("failed to parse file %q: %v", path, err)
		}
		files = append(files, file)
	}
	for _, file := range files {
		for _, section := range file.GetSection() {
			for _, test := range section.GetTest() {
				name := fmt.Sprintf("%s/%s/%s", file.GetName(), section.GetName(), test.GetName())
				if test.GetResultMatcher() == nil {
					test.ResultMatcher = &testpb.SimpleTest_Value{
						Value: &exprpb.Value{
							Kind: &exprpb.Value_BoolValue{
								BoolValue: true,
							},
						},
					}
				}
				t.Run(name, func(t *testing.T) {
					conformanceTest(t, name, test)
				})
			}
		}
	}
}
