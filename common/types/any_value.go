package types

import (
	"github.com/golang/protobuf/ptypes/any"
	"reflect"
)

// anyValueType constant representing the reflected type of google.protobuf.Any.
var anyValueType = reflect.TypeOf(&any.Any{})
