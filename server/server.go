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
	cspb "github.com/google/cel-spec/proto/v1/cel_service"
	"github.com/google/cel-spec/proto/v1/eval"
	"github.com/google/cel-spec/proto/v1/value"
	"github.com/googleapis/googleapis/google/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CelServer struct{}

func (s *CelServer) Parse(ctx context.Context, in *cspb.ParseRequest) (*cspb.ParseResponse, error) {
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
	expr, errs := parser.Parse(src, macs)
	resp := cspb.ParseResponse{}
	if len(errs.GetErrors()) == 0 {
		// Success
		resp.ParsedExpr = expr
	} else {
		// Failure
		appendErrors(errs, &resp.Issues)
	}
	return &resp, nil
}

func (s *CelServer) Check(ctx context.Context, in *cspb.CheckRequest) (*cspb.CheckResponse, error) {
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
	resp := cspb.CheckResponse{}
	if len(errs.GetErrors()) == 0 {
		// Success
		resp.CheckedExpr = c
	} else {
		// Failure
		appendErrors(errs, &resp.Issues)
	}
	return &resp, nil
}

func (s *CelServer) Eval(ctx context.Context, in *cspb.EvalRequest) (*cspb.EvalResponse, error) {
	pkg := packages.NewPackage(in.Container)
	typeProvider := types.NewProvider()
	i := interpreter.NewStandardInterpreter(pkg, typeProvider)
	var prog interpreter.Program
	switch in.ExprKind.(type) {
	case *cspb.EvalRequest_ParsedExpr:
		parsed := in.GetParsedExpr()
		prog = interpreter.NewProgram(parsed.Expr, parsed.SourceInfo)
	case *cspb.EvalRequest_CheckedExpr:
		prog = interpreter.NewCheckedProgram(in.GetCheckedExpr())
	default:
		st := status.New(codes.InvalidArgument, "No expression.")
		return nil, st.Err()
	}
	eval := i.NewInterpretable(prog)
	args := make(map[string]interface{})
	for name, exprValue := range in.Bindings {
		refVal, err := ExprValueToRefValue(exprValue)
		if err != nil {
			return nil, fmt.Errorf("can't convert binding %s: %s", name, err)
		}
		args[name] = refVal
	}
	// NOTE: the EvalState is currently discarded
	result, _ := eval.Eval(interpreter.NewActivation(args))
	resultExprVal, err := RefValueToExprValue(result)
	if err != nil {
		return nil, fmt.Errorf("con't convert result: %s", err)
	}
	return &cspb.EvalResponse{Result: resultExprVal}, nil
}

// appendErrors converts the errors from errs to Status messages
// and appends them to the list of issues.
func appendErrors(errs *common.Errors, issues *[]*rpc.Status) {
	for _, e := range errs.GetErrors() {
		status := ErrToStatus(e, cspb.StatusDetails_ERROR)
		*issues = append(*issues, status)
	}
}

// ErrToStatus converts an Error to a Status message with the given severity.
func ErrToStatus(e common.Error, severity cspb.StatusDetails_Severity) *rpc.Status {
	detail := cspb.StatusDetails{
		Severity: severity,
		Line:     int32(e.Location.Line()),
		Column:   int32(e.Location.Column()),
	}
	// TODO: simply use the following when we unify app-level
	// and gRPC-level Status messages.
	// return status.New(codes.InvalidArgument, e.message).WithDetails(detail).Proto()
	s := rpc.Status{
		Code:    int32(codes.InvalidArgument),
		Message: e.Message,
	}
	any, err := ptypes.MarshalAny(&detail)
	if err == nil {
		s.Details = append(s.Details, any)
	}
	return &s
}

// TODO(jimlarson): The following conversion code should be moved to
// common/types/provider.go and consolidated/refactored as appropriate.
// In particular, make judicious use of types.NativeToValue().

func RefValueToExprValue(res ref.Value) (*eval.ExprValue, error) {
	if types.IsError(res) {
		return &eval.ExprValue{
			Kind: &eval.ExprValue_Error{}}, nil
	}
	if types.IsUnknown(res) {
		return &eval.ExprValue{
			Kind: &eval.ExprValue_Unknown{}}, nil
	}
	v, err := RefValueToValue(res)
	if err != nil {
		return nil, err
	}
	return &eval.ExprValue{
		Kind: &eval.ExprValue_Value{Value: v}}, nil
}

var (
	typeNameToBasicType = map[string]value.TypeValue_BasicType{
		"bool":      value.TypeValue_BOOL_TYPE,
		"bytes":     value.TypeValue_BYTES_TYPE,
		"double":    value.TypeValue_DOUBLE_TYPE,
		"null_type": value.TypeValue_NULL_TYPE,
		"int":       value.TypeValue_INT_TYPE,
		"list":      value.TypeValue_LIST_TYPE,
		"map":       value.TypeValue_MAP_TYPE,
		"string":    value.TypeValue_STRING_TYPE,
		"type":      value.TypeValue_TYPE_TYPE,
		"uint":      value.TypeValue_UINT_TYPE,
	}
	basicTypeToTypeValue = map[value.TypeValue_BasicType]*types.TypeValue{
		value.TypeValue_NULL_TYPE:   types.NullType,
		value.TypeValue_BOOL_TYPE:   types.BoolType,
		value.TypeValue_INT_TYPE:    types.IntType,
		value.TypeValue_UINT_TYPE:   types.UintType,
		value.TypeValue_DOUBLE_TYPE: types.DoubleType,
		value.TypeValue_STRING_TYPE: types.StringType,
		value.TypeValue_BYTES_TYPE:  types.BytesType,
		value.TypeValue_TYPE_TYPE:   types.TypeType,
		value.TypeValue_MAP_TYPE:    types.MapType,
		value.TypeValue_LIST_TYPE:   types.ListType,
	}
)

// Convert res, which must not be error or unknown, to a Value proto.
func RefValueToValue(res ref.Value) (*value.Value, error) {
	switch res.Type() {
	case types.BoolType:
		return &value.Value{
			Kind: &value.Value_BoolValue{res.Value().(bool)}}, nil
	case types.BytesType:
		return &value.Value{
			Kind: &value.Value_BytesValue{res.Value().([]byte)}}, nil
	case types.DoubleType:
		return &value.Value{
			Kind: &value.Value_DoubleValue{res.Value().(float64)}}, nil
	case types.IntType:
		return &value.Value{
			Kind: &value.Value_Int64Value{res.Value().(int64)}}, nil
	case types.ListType:
		l := res.(traits.Lister)
		sz := l.Size().(types.Int)
		elts := make([]*value.Value, int64(sz))
		for i := types.Int(0); i < sz; i++ {
			v, err := RefValueToValue(l.Get(i))
			if err != nil {
				return nil, err
			}
			elts = append(elts, v)
		}
		return &value.Value{
			Kind: &value.Value_ListValue{
				&value.ListValue{Values: elts}}}, nil
	case types.MapType:
		mapper := res.(traits.Mapper)
		sz := mapper.Size().(types.Int)
		entries := make([]*value.MapValue_Entry, int64(sz))
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
			entries = append(entries, &value.MapValue_Entry{Key: kv, Value: vv})
		}
		return &value.Value{
			Kind: &value.Value_MapValue{
				&value.MapValue{Entries: entries}}}, nil
	case types.NullType:
		return &value.Value{
			Kind: &value.Value_NullValue{}}, nil
	case types.StringType:
		return &value.Value{
			Kind: &value.Value_StringValue{res.Value().(string)}}, nil
	case types.TypeType:
		typeName := res.(ref.Type).TypeName()
		var tv *value.TypeValue
		if basicType, found := typeNameToBasicType[typeName]; found {
			// Names a basic type.
			tv = &value.TypeValue{
				DesignatorKind: &value.TypeValue_BasicType_{basicType}}
		} else {
			// Otherwise names a proto.
			tv = &value.TypeValue{
				DesignatorKind: &value.TypeValue_ObjectType{typeName}}
		}
		return &value.Value{Kind: &value.Value_TypeValue{tv}}, nil
	case types.UintType:
		return &value.Value{
			Kind: &value.Value_Uint64Value{res.Value().(uint64)}}, nil
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
		return &value.Value{
			Kind: &value.Value_ObjectValue{any}}, nil
	}
}

func ExprValueToRefValue(ev *eval.ExprValue) (ref.Value, error) {
	switch ev.Kind.(type) {
	case *eval.ExprValue_Value:
		return ValueToRefValue(ev.GetValue())
	case *eval.ExprValue_Error:
		// An error ExprValue is a repeated set of rpc.Status
		// messages, with no convention for the status details.
		// To convert this to a types.Err, we need to convert
		// these Status messages to a single string, and be
		// able to decompose that string on output so we can
		// round-trip arbitrary ExprValue messages.
		// TODO(jimlarson) make a convention for this.
		return types.NewErr("XXX add details later"), nil
	case *eval.ExprValue_Unknown:
		return types.Unknown(ev.GetUnknown().Exprs), nil
	}
	return nil, status.New(codes.InvalidArgument, "unknown ExprValue kind").Err()
}

func ValueToRefValue(v *value.Value) (ref.Value, error) {
	switch v.Kind.(type) {
	case *value.Value_NullValue:
		return types.NullValue, nil
	case *value.Value_BoolValue:
		return types.Bool(v.GetBoolValue()), nil
	case *value.Value_Int64Value:
		return types.Int(v.GetInt64Value()), nil
	case *value.Value_Uint64Value:
		return types.Uint(v.GetUint64Value()), nil
	case *value.Value_DoubleValue:
		return types.Double(v.GetDoubleValue()), nil
	case *value.Value_StringValue:
		return types.String(v.GetStringValue()), nil
	case *value.Value_BytesValue:
		return types.Bytes(v.GetBytesValue()), nil
	case *value.Value_ObjectValue:
		any := v.GetObjectValue()
		var msg ptypes.DynamicAny
		if err := ptypes.UnmarshalAny(any, &msg); err != nil {
			return nil, err
		}
		return types.NewObject(msg.Message), nil
	case *value.Value_MapValue:
		m := v.GetMapValue()
		entries := make(map[ref.Value]ref.Value)
		for _, entry := range m.Entries {
			key, err := ValueToRefValue(entry.Key)
			if err != nil {
				return nil, err
			}
			value, err := ValueToRefValue(entry.Value)
			if err != nil {
				return nil, err
			}
			entries[key] = value
		}
		return types.NewDynamicMap(entries), nil
	case *value.Value_ListValue:
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
	case *value.Value_TypeValue:
		var t *value.TypeValue
		t = v.GetTypeValue()
		switch t.DesignatorKind.(type) {
		case *value.TypeValue_BasicType_:
			bt := t.GetBasicType()
			tv, ok := basicTypeToTypeValue[bt]
			if ok {
				return tv, nil
			}
			return nil, status.New(codes.InvalidArgument, "unknown basic type").Err()
		case *value.TypeValue_ObjectType:
			o := t.GetObjectType()
			return types.NewObjectTypeValue(o), nil
		}
		return nil, status.New(codes.InvalidArgument, "unknown type designator kind").Err()
	}
	return nil, status.New(codes.InvalidArgument, "unknown value").Err()
}
