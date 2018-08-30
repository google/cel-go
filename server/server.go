package server

import (
	"context"
	"fmt"

	protopb "github.com/golang/protobuf/proto"
	ptypespb "github.com/golang/protobuf/ptypes"
	checkerpb "github.com/google/cel-go/checker"
	commonpb "github.com/google/cel-go/common"
	packagespb "github.com/google/cel-go/common/packages"
	typespb "github.com/google/cel-go/common/types"
	refpb "github.com/google/cel-go/common/types/ref"
	traitspb "github.com/google/cel-go/common/types/traits"
	interpreterpb "github.com/google/cel-go/interpreter"
	parserpb "github.com/google/cel-go/parser"
	cspb "github.com/google/cel-spec/proto/v1/cel_service"
	evalpb "github.com/google/cel-spec/proto/v1/eval"
	valuepb "github.com/google/cel-spec/proto/v1/value"
	rpcpb "github.com/googleapis/googleapis/google/rpc"
	codespb "google.golang.org/grpc/codes"
	statuspb "google.golang.org/grpc/status"
)

type CelServer struct{}

func (s *CelServer) Parse(ctx context.Context, in *cspb.ParseRequest) (*cspb.ParseResponse, error) {
	if in.CelSource == "" {
		st := statuspb.New(codespb.InvalidArgument, "No source code.")
		return nil, st.Err()
	}
	// NOTE: syntax_version isn't currently used
	src := commonpb.NewStringSource(in.CelSource, in.SourceLocation)
	var macs parserpb.Macros
	if in.DisableMacros {
		macs = parserpb.NoMacros
	} else {
		macs = parserpb.AllMacros
	}
	expr, errs := parserpb.Parse(src, macs)
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
		st := statuspb.New(codespb.InvalidArgument, "No parsed expression.")
		return nil, st.Err()
	}
	if in.ParsedExpr.SourceInfo == nil {
		st := statuspb.New(codespb.InvalidArgument, "No source info.")
		return nil, st.Err()
	}
	pkg := packagespb.NewPackage(in.Container)
	typeProvider := typespb.NewProvider()
	errs := commonpb.NewErrors(commonpb.NewInfoSource(in.ParsedExpr.SourceInfo))
	var env *checkerpb.Env
	if in.NoStdEnv {
		env = checkerpb.NewEnv(pkg, typeProvider, errs)
	} else {
		env = checkerpb.NewStandardEnv(pkg, typeProvider, errs)
	}
	env.Add(in.TypeEnv...)
	c := checkerpb.Check(in.ParsedExpr, env)
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
	pkg := packagespb.NewPackage(in.Container)
	typeProvider := typespb.NewProvider()
	i := interpreterpb.NewStandardIntepreter(pkg, typeProvider)
	var prog interpreterpb.Program
	switch in.ExprKind.(type) {
	case *cspb.EvalRequest_ParsedExpr:
		parsed := in.GetParsedExpr()
		prog = interpreterpb.NewProgram(parsed.Expr, parsed.SourceInfo)
	case *cspb.EvalRequest_CheckedExpr:
		prog = interpreterpb.NewCheckedProgram(in.GetCheckedExpr())
	default:
		st := statuspb.New(codespb.InvalidArgument, "No expression.")
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
	result, _ := evalpb.Eval(interpreterpb.NewActivation(args))
	resultExprVal, err := RefValueToExprValue(result)
	if err != nil {
		return nil, fmt.Errorf("con't convert result: %s", err)
	}
	return &cspb.EvalResponse{Result: resultExprVal}, nil
}

// appendErrors converts the errors from errs to Status messages
// and appends them to the list of issues.
func appendErrors(errs *commonpb.Errors, issues *[]*rpcpb.Status) {
	for _, e := range errs.GetErrors() {
		status := ErrToStatus(e, cspb.StatusDetails_ERROR)
		*issues = append(*issues, status)
	}
}

// ErrToStatus converts an Error to a Status message with the given severity.
func ErrToStatus(e commonpb.Error, severity cspb.StatusDetails_Severity) *rpcpb.Status {
	detail := cspb.StatusDetails{
		Severity: severity,
		Line:     int32(e.Location.Line()),
		Column:   int32(e.Location.Column()),
	}
	// TODO: simply use the following when we unify app-level
	// and gRPC-level Status messages.
	// return statuspb.New(codespb.InvalidArgument, e.message).WithDetails(detail).Proto()
	s := rpcpb.Status{
		Code:    int32(codespb.InvalidArgument),
		Message: e.Message,
	}
	any, err := ptypespb.MarshalAny(&detail)
	if err == nil {
		s.Details = append(s.Details, any)
	}
	return &s
}

// TODO(jimlarson): The following conversion code should be moved to
// common/types/provider.go and consolidated/refactored as appropriate.
// In particular, make judicious use of types.NativeToValue().

func RefValueToExprValue(res refpb.Value) (*evalpb.ExprValue, error) {
	if typespb.IsError(res) {
		return &evalpb.ExprValue{
			Kind: &evalpb.ExprValue_Error{}}, nil
	}
	if typespb.IsUnknown(res) {
		return &evalpb.ExprValue{
			Kind: &evalpb.ExprValue_Unknown{}}, nil
	}
	v, err := RefValueToValue(res)
	if err != nil {
		return nil, err
	}
	return &evalpb.ExprValue{
		Kind: &evalpb.ExprValue_Value{Value: v}}, nil
}

var (
	typeNameToBasicType = map[string]valuepb.TypeValue_BasicType{
		"bool":      valuepb.TypeValue_BOOL_TYPE,
		"bytes":     valuepb.TypeValue_BYTES_TYPE,
		"double":    valuepb.TypeValue_DOUBLE_TYPE,
		"null_type": valuepb.TypeValue_NULL_TYPE,
		"int":       valuepb.TypeValue_INT_TYPE,
		"list":      valuepb.TypeValue_LIST_TYPE,
		"map":       valuepb.TypeValue_MAP_TYPE,
		"string":    valuepb.TypeValue_STRING_TYPE,
		"type":      valuepb.TypeValue_TYPE_TYPE,
		"uint":      valuepb.TypeValue_UINT_TYPE,
	}
	basicTypeToTypeValue = map[valuepb.TypeValue_BasicType]*typespb.TypeValue{
		valuepb.TypeValue_NULL_TYPE:   typespb.NullType,
		valuepb.TypeValue_BOOL_TYPE:   typespb.BoolType,
		valuepb.TypeValue_INT_TYPE:    typespb.IntType,
		valuepb.TypeValue_UINT_TYPE:   typespb.UintType,
		valuepb.TypeValue_DOUBLE_TYPE: typespb.DoubleType,
		valuepb.TypeValue_STRING_TYPE: typespb.StringType,
		valuepb.TypeValue_BYTES_TYPE:  typespb.BytesType,
		valuepb.TypeValue_TYPE_TYPE:   typespb.TypeType,
		valuepb.TypeValue_MAP_TYPE:    typespb.MapType,
		valuepb.TypeValue_LIST_TYPE:   typespb.ListType,
	}
)

// Convert res, which must not be error or unknown, to a Value proto.
func RefValueToValue(res refpb.Value) (*valuepb.Value, error) {
	switch res.Type() {
	case typespb.BoolType:
		return &valuepb.Value{
			Kind: &valuepb.Value_BoolValue{res.Value().(bool)}}, nil
	case typespb.BytesType:
		return &valuepb.Value{
			Kind: &valuepb.Value_BytesValue{res.Value().([]byte)}}, nil
	case typespb.DoubleType:
		return &valuepb.Value{
			Kind: &valuepb.Value_DoubleValue{res.Value().(float64)}}, nil
	case typespb.IntType:
		return &valuepb.Value{
			Kind: &valuepb.Value_Int64Value{res.Value().(int64)}}, nil
	case typespb.ListType:
		l := res.(traitspb.Lister)
		sz := l.Size().(typespb.Int)
		elts := make([]*valuepb.Value, int64(sz))
		for i := typespb.Int(0); i < sz; i++ {
			v, err := RefValueToValue(l.Get(i))
			if err != nil {
				return nil, err
			}
			elts = append(elts, v)
		}
		return &valuepb.Value{
			Kind: &valuepb.Value_ListValue{
				&valuepb.ListValue{Values: elts}}}, nil
	case typespb.MapType:
		mapper := res.(traitspb.Mapper)
		sz := mapper.Size().(typespb.Int)
		entries := make([]*valuepb.MapValue_Entry, int64(sz))
		for it := mapper.Iterator(); it.HasNext().(typespb.Bool); {
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
			entries = append(entries, &valuepb.MapValue_Entry{Key: kv, Value: vv})
		}
		return &valuepb.Value{
			Kind: &valuepb.Value_MapValue{
				&valuepb.MapValue{Entries: entries}}}, nil
	case typespb.NullType:
		return &valuepb.Value{
			Kind: &valuepb.Value_NullValue{}}, nil
	case typespb.StringType:
		return &valuepb.Value{
			Kind: &valuepb.Value_StringValue{res.Value().(string)}}, nil
	case typespb.TypeType:
		typeName := res.(refpb.Type).TypeName()
		var tv *valuepb.TypeValue
		if basicType, found := typeNameToBasicType[typeName]; found {
			// Names a basic type.
			tv = &valuepb.TypeValue{
				DesignatorKind: &valuepb.TypeValue_BasicType_{basicType}}
		} else {
			// Otherwise names a proto.
			tv = &valuepb.TypeValue{
				DesignatorKind: &valuepb.TypeValue_ObjectType{typeName}}
		}
		return &valuepb.Value{Kind: &valuepb.Value_TypeValue{tv}}, nil
	case typespb.UintType:
		return &valuepb.Value{
			Kind: &valuepb.Value_Uint64Value{res.Value().(uint64)}}, nil
	default:
		// Object type
		pb, ok := res.Value().(protopb.Message)
		if !ok {
			return nil, statuspb.New(codespb.InvalidArgument, "Expected proto message").Err()
		}
		any, err := ptypespb.MarshalAny(pb)
		if err != nil {
			return nil, err
		}
		return &valuepb.Value{
			Kind: &valuepb.Value_ObjectValue{any}}, nil
	}
}

func ExprValueToRefValue(ev *evalpb.ExprValue) (refpb.Value, error) {
	switch ev.Kind.(type) {
	case *evalpb.ExprValue_Value:
		return ValueToRefValue(ev.GetValue())
	case *evalpb.ExprValue_Error:
		// An error ExprValue is a repeated set of rpcpb.Status
		// messages, with no convention for the status details.
		// To convert this to a types.Err, we need to convert
		// these Status messages to a single string, and be
		// able to decompose that string on output so we can
		// round-trip arbitrary ExprValue messages.
		// TODO(jimlarson) make a convention for this.
		return typespb.NewErr("XXX add details later"), nil
	case *evalpb.ExprValue_Unknown:
		return typespb.Unknown(ev.GetUnknown().Exprs), nil
	}
	return nil, statuspb.New(codespb.InvalidArgument, "unknown ExprValue kind").Err()
}

func ValueToRefValue(v *valuepb.Value) (refpb.Value, error) {
	switch v.Kind.(type) {
	case *valuepb.Value_NullValue:
		return typespb.NullValue, nil
	case *valuepb.Value_BoolValue:
		return typespb.Bool(v.GetBoolValue()), nil
	case *valuepb.Value_Int64Value:
		return typespb.Int(v.GetInt64Value()), nil
	case *valuepb.Value_Uint64Value:
		return typespb.Uint(v.GetUint64Value()), nil
	case *valuepb.Value_DoubleValue:
		return typespb.Double(v.GetDoubleValue()), nil
	case *valuepb.Value_StringValue:
		return typespb.String(v.GetStringValue()), nil
	case *valuepb.Value_BytesValue:
		return typespb.Bytes(v.GetBytesValue()), nil
	case *valuepb.Value_ObjectValue:
		any := v.GetObjectValue()
		var msg ptypespb.DynamicAny
		if err := ptypespb.UnmarshalAny(any, &msg); err != nil {
			return nil, err
		}
		return typespb.NewObject(msg.Message), nil
	case *valuepb.Value_MapValue:
		m := v.GetMapValue()
		entries := make(map[refpb.Value]refpb.Value)
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
		return typespb.NewDynamicMap(entries), nil
	case *valuepb.Value_ListValue:
		l := v.GetListValue()
		elts := make([]refpb.Value, len(l.Values))
		for i, e := range l.Values {
			rv, err := ValueToRefValue(e)
			if err != nil {
				return nil, err
			}
			elts[i] = rv
		}
		return typespb.NewValueList(elts), nil
	case *valuepb.Value_TypeValue:
		var t *valuepb.TypeValue
		t = v.GetTypeValue()
		switch t.DesignatorKind.(type) {
		case *valuepb.TypeValue_BasicType_:
			bt := t.GetBasicType()
			tv, ok := basicTypeToTypeValue[bt]
			if ok {
				return tv, nil
			}
			return nil, statuspb.New(codespb.InvalidArgument, "unknown basic type").Err()
		case *valuepb.TypeValue_ObjectType:
			o := t.GetObjectType()
			return typespb.NewObjectTypeValue(o), nil
		}
		return nil, statuspb.New(codespb.InvalidArgument, "unknown type designator kind").Err()
	}
	return nil, statuspb.New(codespb.InvalidArgument, "unknown value").Err()
}
