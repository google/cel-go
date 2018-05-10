package types

import (
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
)

// jsonValueType constant representing the reflected type of a protobuf Value.
var jsonValueType = reflect.TypeOf(&structpb.Value{})
