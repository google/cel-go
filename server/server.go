package server

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/parser"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	rpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CelServer struct{}

func (s *CelServer) Parse(ctx context.Context, in *expr.ParseRequest) (*expr.ParseResponse, error) {
	if in.CelSource == "" {
		st := status.New(codes.InvalidArgument, "No source code.")
		return nil, st.Err()
	}
	// NOTE: syntax_version isn't currently used
	src := common.NewStringSource(in.CelSource, in.SourceLocation)
	var macs parser.Macros
	if in.DisableMacros {
		macs = parser.NoMacros
	} else {
		macs = parser.AllMacros
	}
	pexpr, errs := parser.Parse(src, macs)
	resp := expr.ParseResponse{}
	if len(errs.GetErrors()) == 0 {
		// Success
		resp.ParsedExpr = pexpr
	} else {
		// Failure
		appendErrors(errs, &resp.Issues)
	}
	return &resp, nil
}

func (s *CelServer) Check(ctx context.Context, in *expr.CheckRequest) (*expr.CheckResponse, error) {
	if in.ParsedExpr == nil {
		st := status.New(codes.InvalidArgument, "No parsed expression.")
		return nil, st.Err()
	}
	if in.ParsedExpr.SourceInfo == nil {
		st := status.New(codes.InvalidArgument, "No source info.")
		return nil, st.Err()
	}
	pkg := packages.NewPackage(in.Container)
	typeProvider := types.NewProvider()
	errs := common.NewErrors(common.NewInfoSource(in.ParsedExpr.SourceInfo))
	var env *checker.Env
	if in.NoStdEnv {
		env = checker.NewEnv(pkg, typeProvider, errs)
	} else {
		env = checker.NewStandardEnv(pkg, typeProvider, errs)
	}
	env.Add(in.TypeEnv...)
	c := checker.Check(in.ParsedExpr, env)
	resp := expr.CheckResponse{}
	if len(errs.GetErrors()) == 0 {
		// Success
		resp.CheckedExpr = c
	} else {
		// Failure
		appendErrors(errs, &resp.Issues)
	}
	return &resp, nil
}

func (s *CelServer) Eval(ctx context.Context, in *expr.EvalRequest) (*expr.EvalResponse, error) {
	pkg := packages.NewPackage(in.Container)
	typeProvider := types.NewProvider()
	i := interpreter.NewStandardIntepreter(pkg, typeProvider)
	var prog interpreter.Program
	switch in.ExprKind.(type) {
	case *expr.EvalRequest_ParsedExpr:
		parsed := in.GetParsedExpr()
		prog = interpreter.NewProgram(parsed.Expr, parsed.SourceInfo)
	case *expr.EvalRequest_CheckedExpr:
		prog = interpreter.NewCheckedProgram(in.GetCheckedExpr())
	default:
		st := status.New(codes.InvalidArgument, "No expression.")
		return nil, st.Err()
	}
	ev := i.NewInterpretable(prog)
	args := make(map[string]interface{})
	for name, exprValue := range in.Bindings {
		refVal, err := ExprValueToRefValue(exprValue)
		if err != nil {
			return nil, fmt.Errorf("can't convert binding %s: %s", name, err)
		}
		args[name] = refVal
	}
	// NOTE: the EvalState is currently discarded
	result, _ := ev.Eval(interpreter.NewActivation(args))
	resultExprVal, err := RefValueToExprValue(result)
	if err != nil {
		return nil, fmt.Errorf("con't convert result: %s", err)
	}
	return &expr.EvalResponse{Result: resultExprVal}, nil
}

// appendErrors converts the errors from errs to Status messages
// and appends them to the list of issues.
func appendErrors(errs *common.Errors, issues *[]*rpc.Status) {
	for _, e := range errs.GetErrors() {
		status := ErrToStatus(e, expr.IssueDetails_ERROR)
		*issues = append(*issues, status)
	}
}

// ErrToStatus converts an Error to a Status message with the given severity.
func ErrToStatus(e common.Error, severity expr.IssueDetails_Severity) *rpc.Status {
	detail := expr.IssueDetails{
		Severity: severity,
		Position: &expr.SourcePosition{
			Line:   int32(e.Location.Line()),
			Column: int32(e.Location.Column()),
		},
	}
	s := status.New(codes.InvalidArgument, e.Message)
	sd, err := s.WithDetails(&detail)
	if err == nil {
		return sd.Proto()
	} else {
		return s.Proto()
	}
}

// TODO(jimlarson): The following conversion code should be moved to
// common/types/provider.go and consolidated/refactored as appropriate.
// In particular, make judicious use of types.NativeToValue().

func RefValueToExprValue(res ref.Value) (*expr.ExprValue, error) {
	if types.IsError(res) {
		return &expr.ExprValue{
			Kind: &expr.ExprValue_Error{}}, nil
	}
	if types.IsUnknown(res) {
		return &expr.ExprValue{
			Kind: &expr.ExprValue_Unknown{}}, nil
	}
	v, err := RefValueToValue(res)
	if err != nil {
		return nil, err
	}
	return &expr.ExprValue{
		Kind: &expr.ExprValue_Value{Value: v}}, nil
}

var (
	typeNameToTypeValue = map[string]*types.TypeValue{
		"bool":      types.BoolType,
		"bytes":     types.BytesType,
		"double":    types.DoubleType,
		"null_type": types.NullType,
		"int":       types.IntType,
		"list":      types.ListType,
		"map":       types.MapType,
		"string":    types.StringType,
		"type":      types.TypeType,
		"uint":      types.UintType,
	}
)

// Convert res, which must not be error or unknown, to a Value proto.
func RefValueToValue(res ref.Value) (*expr.Value, error) {
	switch res.Type() {
	case types.BoolType:
		return &expr.Value{
			Kind: &expr.Value_BoolValue{res.Value().(bool)}}, nil
	case types.BytesType:
		return &expr.Value{
			Kind: &expr.Value_BytesValue{res.Value().([]byte)}}, nil
	case types.DoubleType:
		return &expr.Value{
			Kind: &expr.Value_DoubleValue{res.Value().(float64)}}, nil
	case types.IntType:
		return &expr.Value{
			Kind: &expr.Value_Int64Value{res.Value().(int64)}}, nil
	case types.ListType:
		l := res.(traits.Lister)
		sz := l.Size().(types.Int)
		elts := make([]*expr.Value, int64(sz))
		for i := types.Int(0); i < sz; i++ {
			v, err := RefValueToValue(l.Get(i))
			if err != nil {
				return nil, err
			}
			elts = append(elts, v)
		}
		return &expr.Value{
			Kind: &expr.Value_ListValue{
				&expr.ListValue{Values: elts}}}, nil
	case types.MapType:
		mapper := res.(traits.Mapper)
		sz := mapper.Size().(types.Int)
		entries := make([]*expr.MapValue_Entry, int64(sz))
		for it := mapper.Iterator(); it.HasNext().(types.Bool); {
			k := it.Next()
			v := mapper.Get(k)
			kv, err := RefValueToValue(k)
			if err != nil {
				return nil, err
			}
			vv, err := RefValueToValue(v)
			if err != nil {
				return nil, err
			}
			entries = append(entries, &expr.MapValue_Entry{Key: kv, Value: vv})
		}
		return &expr.Value{
			Kind: &expr.Value_MapValue{
				&expr.MapValue{Entries: entries}}}, nil
	case types.NullType:
		return &expr.Value{
			Kind: &expr.Value_NullValue{}}, nil
	case types.StringType:
		return &expr.Value{
			Kind: &expr.Value_StringValue{res.Value().(string)}}, nil
	case types.TypeType:
		typeName := res.(ref.Type).TypeName()
		return &expr.Value{Kind: &expr.Value_TypeValue{typeName}}, nil
	case types.UintType:
		return &expr.Value{
			Kind: &expr.Value_Uint64Value{res.Value().(uint64)}}, nil
	default:
		// Object type
		pb, ok := res.Value().(proto.Message)
		if !ok {
			return nil, status.New(codes.InvalidArgument, "Expected proto message").Err()
		}
		any, err := ptypes.MarshalAny(pb)
		if err != nil {
			return nil, err
		}
		return &expr.Value{
			Kind: &expr.Value_ObjectValue{any}}, nil
	}
}

func ExprValueToRefValue(ev *expr.ExprValue) (ref.Value, error) {
	switch ev.Kind.(type) {
	case *expr.ExprValue_Value:
		return ValueToRefValue(ev.GetValue())
	case *expr.ExprValue_Error:
		// An error ExprValue is a repeated set of rpc.Status
		// messages, with no convention for the status details.
		// To convert this to a types.Err, we need to convert
		// these Status messages to a single string, and be
		// able to decompose that string on output so we can
		// round-trip arbitrary ExprValue messages.
		// TODO(jimlarson) make a convention for this.
		return types.NewErr("XXX add details later"), nil
	case *expr.ExprValue_Unknown:
		return types.Unknown(ev.GetUnknown().Exprs), nil
	}
	return nil, status.New(codes.InvalidArgument, "unknown ExprValue kind").Err()
}

func ValueToRefValue(v *expr.Value) (ref.Value, error) {
	switch v.Kind.(type) {
	case *expr.Value_NullValue:
		return types.NullValue, nil
	case *expr.Value_BoolValue:
		return types.Bool(v.GetBoolValue()), nil
	case *expr.Value_Int64Value:
		return types.Int(v.GetInt64Value()), nil
	case *expr.Value_Uint64Value:
		return types.Uint(v.GetUint64Value()), nil
	case *expr.Value_DoubleValue:
		return types.Double(v.GetDoubleValue()), nil
	case *expr.Value_StringValue:
		return types.String(v.GetStringValue()), nil
	case *expr.Value_BytesValue:
		return types.Bytes(v.GetBytesValue()), nil
	case *expr.Value_ObjectValue:
		any := v.GetObjectValue()
		var msg ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(any, &msg); err != nil {
			return nil, err
		}
		return types.NewObject(msg.Message), nil
	case *expr.Value_MapValue:
		m := v.GetMapValue()
		entries := make(map[ref.Value]ref.Value)
		for _, entry := range m.Entries {
			key, err := ValueToRefValue(entry.Key)
			if err != nil {
				return nil, err
			}
			pb, err := ValueToRefValue(entry.Value)
			if err != nil {
				return nil, err
			}
			entries[key] = pb
		}
		return types.NewDynamicMap(entries), nil
	case *expr.Value_ListValue:
		l := v.GetListValue()
		elts := make([]ref.Value, len(l.Values))
		for i, e := range l.Values {
			rv, err := ValueToRefValue(e)
			if err != nil {
				return nil, err
			}
			elts[i] = rv
		}
		return types.NewValueList(elts), nil
	case *expr.Value_TypeValue:
		typeName := v.GetTypeValue()
		tv, ok := typeNameToTypeValue[typeName]
		if ok {
			return tv, nil
		} else {
			return types.NewObjectTypeValue(typeName), nil
		}
	}
	return nil, status.New(codes.InvalidArgument, "unknown value").Err()
}
