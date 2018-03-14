package adapters

import (
	"github.com/google/cel-go/interpreter/types/objects"
	"fmt"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"reflect"
)

type TypeConverter func(*reflect.Value, interface{}) interface{}

func ProtoToExpr(value interface{}) (interface{}, error) {
	switch value.(type) {
	// Compatible types
	case bool, float64, int64, uint64, string, []byte,
		*dpb.Duration, *tspb.Timestamp, structpb.NullValue:
		return value, nil
	// Upconverted numeric types
	case int32:
		return int64(value.(int32)), nil
	case uint32:
		return uint64(value.(uint32)), nil
	case float32:
		return float64(value.(float32)), nil
	// Views on complex types which may have a mix of compatible types as well
	// as ones that need upconversion.
	case proto.Message:
		return NewMsgAdapter(value), nil
	default:
		// lists and maps should be all that is left
		refValue := reflect.ValueOf(value)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		refKind := refValue.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			return NewListAdapter(value), nil
		case reflect.Map:
			return NewMapAdapter(value), nil
		}
	}
	return nil, fmt.Errorf("unimplemented proto to expr conversion for:"+
		"%T %v", value, value)
}

func ExprToProto(refType reflect.Type, value interface{}) (interface{}, error) {
	refKind := refType.Kind()
	switch refKind {
	case reflect.Bool, reflect.Float64, reflect.Int64, reflect.Uint64, reflect.String:
		return value, nil
	case reflect.Int32:
		switch value.(type) {
		case structpb.NullValue, int32:
			return value, nil
		default:
			return int32(value.(int64)), nil
		}
	case reflect.Uint32:
		switch value.(type) {
		case uint32:
			return value, nil
		default:
			return uint32(value.(uint64)), nil
		}
	case reflect.Float32:
		switch value.(type) {
		case float32:
			return value, nil
		default:
			return float32(value.(float64)), nil
		}
	case reflect.Map:
		switch value.(type) {
		case objects.Protoer:
			return value.(objects.Protoer).ToProto(refType)
		default:
			return NewMapAdapter(value).ToProto(refType)
		}
	case reflect.Slice, reflect.Array:
		switch value.(type) {
		case objects.Protoer:
			return value.(objects.Protoer).ToProto(refType)
		case []byte:
			return value.([]byte), nil
		default:
			return NewListAdapter(value).ToProto(refType)
		}
	case reflect.Struct, reflect.Ptr:
		switch value.(type) {
		case proto.Message, structpb.NullValue:
			return value, nil
		case MsgAdapter:
			return value.(MsgAdapter).Value(), nil
		}
	}
	return nil, fmt.Errorf("unimplemented expr to proto conversion for:"+
		" %v to type %v", value, refKind)
}
