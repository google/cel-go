package server

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"testing"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/test"
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	"github.com/google/cel-spec/proto/v1/cel_service"
	"github.com/google/cel-spec/proto/v1/eval"
	"github.com/google/cel-spec/proto/v1/syntax"
	"github.com/google/cel-spec/proto/v1/value"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type serverTest struct {
	cmd    *exec.Cmd
	conn   *grpc.ClientConn
	client cel_service.CelServiceClient
}

var (
	globals = serverTest{}
)

func TestMain(m *testing.M) {
	// Use a helper function to ensure we run shutdown()
	// before calling os.Exit()
	os.Exit(mainHelper(m))
}

func mainHelper(m *testing.M) int {
	err := setup()
	defer shutdown()
	if err != nil {
		// testing.M doesn't have a logging method.  hmm...
		log.Fatal(err)
		return 1
	}
	return m.Run()
}

func setup() error {
	globals.cmd = exec.Command("cel_server")

	out, err := globals.cmd.StdoutPipe()
	if err != nil {
		return err
	}
	globals.cmd.Stderr = os.Stderr // share our error stream

	log.Println("Starting server")
	err = globals.cmd.Start()
	if err != nil {
		return err
	}

	log.Println("Getting server's address")
	var addr string
	_, err = fmt.Fscanf(out, "Listening on %s\n", &addr)
	out.Close()
	if err != nil {
		return err
	}

	log.Println("Connecting to ", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	globals.conn = conn

	log.Println("Creating service client")
	globals.client = cel_service.NewCelServiceClient(conn)
	return nil
}

func shutdown() {
	if globals.conn != nil {
		globals.conn.Close()
		globals.conn = nil
	}
	if globals.cmd != nil {
		globals.cmd.Process.Kill()
		globals.cmd.Wait()
		globals.cmd = nil
	}
}

var (
	parsed = &syntax.ParsedExpr{
		Expr: test.ExprCall(1, operators.Add,
			test.ExprLiteral(2, int64(1)),
			test.ExprLiteral(3, int64(1))),
		SourceInfo: &syntax.SourceInfo{
			Location: "the location",
			Positions: map[int64]int32{
				1: 0,
				2: 0,
				3: 4,
			},
		},
	}
)

func TestParse(t *testing.T) {
	req := cel_service.ParseRequest{
		CelSource: "1 + 1",
	}
	res, err := globals.client.Parse(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Empty result")
	}
	if res.ParsedExpr == nil {
		t.Fatal("Empty parsed expression in result")
	}
	// Could check against 'parsed' above,
	// but the expression ids are arbitrary,
	// and explicit comparison logic is about as
	// much work as normalization would be.
	if res.ParsedExpr.Expr == nil {
		t.Fatal("Empty expression in result")
	}
	switch res.ParsedExpr.Expr.ExprKind.(type) {
	case *syntax.Expr_CallExpr:
		c := res.ParsedExpr.Expr.GetCallExpr()
		if c.Target != nil {
			t.Error("Call has target", c)
		}
		if c.Function != "_+_" {
			t.Error("Wrong function", c)
		}
		if len(c.Args) != 2 {
			t.Error("Too many or few args", c)
		}
		for i, a := range c.Args {
			switch a.ExprKind.(type) {
			case *syntax.Expr_LiteralExpr:
				l := a.GetLiteralExpr()
				switch l.LiteralKind.(type) {
				case *syntax.Literal_Int64Value:
					if l.GetInt64Value() != int64(1) {
						t.Errorf("Arg %d wrong value: %v", i, a)
					}
				default:
					t.Errorf("Arg %d not int: %v", i, a)
				}
			default:
				t.Errorf("Arg %d not literal: %v", i, a)
			}
		}
	default:
		t.Error("Wrong expression type", res.ParsedExpr.Expr)
	}
}

func TestCheck(t *testing.T) {
	// If TestParse() passes, it validates a good chunk
	// of the server mechanisms for data conversion, so we
	// won't be as fussy here..
	req := cel_service.CheckRequest{
		ParsedExpr: parsed,
	}
	res, err := globals.client.Check(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil {
		t.Fatal("Empty result")
	}
	if res.CheckedExpr == nil {
		t.Fatal("No checked expression")
	}
	tp, present := res.CheckedExpr.TypeMap[int64(1)]
	if !present {
		t.Fatal("No type for top level expression", res)
	}
	switch tp.TypeKind.(type) {
	case *checkedpb.Type_Primitive:
		if tp.GetPrimitive() != checkedpb.Type_INT64 {
			t.Error("Bad top-level type", tp)
		}
	default:
		t.Error("Bad top-level type", tp)
	}
}

func TestEval(t *testing.T) {
	req := cel_service.EvalRequest{
		ExprKind: &cel_service.EvalRequest_ParsedExpr{parsed},
	}
	res, err := globals.client.Eval(context.Background(), &req)
	if err != nil {
		t.Fatal(err)
	}
	if res == nil || res.Result == nil {
		t.Fatal("Nil result")
	}
	switch res.Result.Kind.(type) {
	case *eval.ExprValue_Value:
		v := res.Result.GetValue()
		switch v.Kind.(type) {
		case *value.Value_Int64Value:
			if v.GetInt64Value() != int64(2) {
				t.Error("Wrong result for 1 + 1", v)
			}
		default:
			t.Error("Wrong result value type", v)
		}
	default:
		t.Fatal("Result not a value", res.Result)
	}
}

func TestFullUp(t *testing.T) {
	preq := cel_service.ParseRequest{
		CelSource: "x + y",
	}
	pres, err := globals.client.Parse(context.Background(), &preq)
	if err != nil {
		t.Fatal(err)
	}
	parsedExpr := pres.ParsedExpr
	if parsedExpr == nil {
		t.Fatal("Empty parsed expression")
	}

	creq := cel_service.CheckRequest{
		ParsedExpr: parsedExpr,
		TypeEnv: []*checkedpb.Decl{
			decls.NewIdent("x", decls.Int, nil),
			decls.NewIdent("y", decls.Int, nil),
		},
	}
	cres, err := globals.client.Check(context.Background(), &creq)
	if err != nil {
		t.Fatal(err)
	}
	if cres == nil {
		t.Fatal("Empty check result")
	}
	checkedExpr := cres.CheckedExpr
	if checkedExpr == nil {
		t.Fatal("No checked expression")
	}
	tp, present := checkedExpr.TypeMap[int64(1)]
	if !present {
		t.Fatal("No type for top level expression", cres)
	}
	switch tp.TypeKind.(type) {
	case *checkedpb.Type_Primitive:
		if tp.GetPrimitive() != checkedpb.Type_INT64 {
			t.Error("Bad top-level type", tp)
		}
	default:
		t.Error("Bad top-level type", tp)
	}

	ereq := cel_service.EvalRequest{
		ExprKind: &cel_service.EvalRequest_CheckedExpr{checkedExpr},
		Bindings: map[string]*eval.ExprValue{
			"x": exprValueInt64(1),
			"y": exprValueInt64(2),
		},
	}
	eres, err := globals.client.Eval(context.Background(), &ereq)
	if err != nil {
		t.Fatal(err)
	}
	if eres == nil || eres.Result == nil {
		t.Fatal("Nil result")
	}
	switch eres.Result.Kind.(type) {
	case *eval.ExprValue_Value:
		v := eres.Result.GetValue()
		switch v.Kind.(type) {
		case *value.Value_Int64Value:
			if v.GetInt64Value() != int64(3) {
				t.Error("Wrong result for 1 + 2", v)
			}
		default:
			t.Error("Wrong result value type", v)
		}
	default:
		t.Fatal("Result not a value", eres.Result)
	}
}

func exprValueInt64(x int64) *eval.ExprValue {
	return &eval.ExprValue{
		Kind: &eval.ExprValue_Value{
			&value.Value{
				Kind: &value.Value_Int64Value{x},
			},
		},
	}
}

// expectEvalTrue parses, checks, and evaluates the CEL expression in source
// and checks that the result is the boolean value 'true'.
func expectEvalTrue(t *testing.T, source string) {
	// Parse
	preq := cel_service.ParseRequest{
		CelSource: source,
	}
	pres, err := globals.client.Parse(context.Background(), &preq)
	if err != nil {
		t.Fatal(err)
	}
	if pres == nil {
		t.Fatal("Empty parse result")
	}
	parsedExpr := pres.ParsedExpr
	if parsedExpr == nil {
		t.Fatal("Empty parsed expression")
	}
	if parsedExpr.Expr == nil {
		t.Fatal("Empty root expression")
	}
	rootId := parsedExpr.Expr.Id

	// Check
	creq := cel_service.CheckRequest{
		ParsedExpr: parsedExpr,
	}
	cres, err := globals.client.Check(context.Background(), &creq)
	if err != nil {
		t.Fatal(err)
	}
	if cres == nil {
		t.Fatal("Empty check result")
	}
	checkedExpr := cres.CheckedExpr
	if checkedExpr == nil {
		t.Fatal("No checked expression")
	}
	topType, present := checkedExpr.TypeMap[rootId]
	if !present {
		t.Fatal("No type for top level expression", cres)
	}
	switch topType.TypeKind.(type) {
	case *checkedpb.Type_Primitive:
		if topType.GetPrimitive() != checkedpb.Type_BOOL {
			t.Error("Bad top-level type", topType)
		}
	default:
		t.Error("Bad top-level type", topType)
	}

	// Eval
	ereq := cel_service.EvalRequest{
		ExprKind: &cel_service.EvalRequest_CheckedExpr{checkedExpr},
	}
	eres, err := globals.client.Eval(context.Background(), &ereq)
	if err != nil {
		t.Fatal(err)
	}
	if eres == nil || eres.Result == nil {
		t.Fatal("Nil result")
	}
	switch eres.Result.Kind.(type) {
	case *eval.ExprValue_Value:
		v := eres.Result.GetValue()
		switch v.Kind.(type) {
		case *value.Value_BoolValue:
			if !v.GetBoolValue() {
				t.Error("Wrong result", v)
			}
		default:
			t.Error("Wrong result value type", v)
		}
	default:
		t.Fatal("Result not a value", eres.Result)
	}
}

func TestCondTrue(t *testing.T) {
	expectEvalTrue(t, "(true ? 'a' : 'b') == 'a'")
}

func TestCondFalse(t *testing.T) {
	expectEvalTrue(t, "(false ? 'a' : 'b') == 'b'")
}

func TestMapOrderInsignificant(t *testing.T) {
	expectEvalTrue(t, "{1: 'a', 2: 'b'} == {2: 'b', 1: 'a'}")
}

func FailsTestOneMetaType(t *testing.T) {
	expectEvalTrue(t, "type(type(1)) == type(type('foo'))")
}

func FailsTestTypeType(t *testing.T) {
	expectEvalTrue(t, "type(type) == type")
}

func FailsTestNullTypeName(t *testing.T) {
	expectEvalTrue(t, "type(null) == null_type")
}
