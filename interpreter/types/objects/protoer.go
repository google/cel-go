package objects

import "reflect"

type Protoer interface {
	ToProto(typeDesc reflect.Type) (interface{}, error)
}
