// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/cel-go/interpreter/types/traits"
	"reflect"
)

// ProtoToExpr converts from a proto value into an expression value.
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
		return NewProtoValue(value), nil
	default:
		// lists and maps should be all that is left
		refValue := reflect.ValueOf(value)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		refKind := refValue.Kind()
		switch refKind {
		case reflect.Array, reflect.Slice:
			return NewListValue(value), nil
		case reflect.Map:
			return NewMapValue(value), nil
		}
	}
	return nil, fmt.Errorf("unimplemented proto to expr conversion for:"+
		"%T %v", value, value)
}

// ExprToProto converts a value to a value of the reflected type
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
		case traits.Protoer:
			return value.(traits.Protoer).ToProto(refType)
		default:
			return NewMapValue(value).ToProto(refType)
		}
	case reflect.Slice, reflect.Array:
		switch value.(type) {
		case traits.Protoer:
			return value.(traits.Protoer).ToProto(refType)
		case []byte:
			return value.([]byte), nil
		default:
			return NewListValue(value).ToProto(refType)
		}
	case reflect.Struct, reflect.Ptr:
		switch value.(type) {
		case proto.Message, structpb.NullValue:
			return value, nil
		case ObjectValue:
			return value.(ObjectValue).Value(), nil
		}
	}
	return nil, fmt.Errorf("unimplemented expr to proto conversion for:"+
		" %v to type %v", value, refKind)
}
