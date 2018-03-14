package types

import (
	"github.com/google/cel-go/interpreter/types/adapters"
	"github.com/google/cel-go/interpreter/types/objects"
	"fmt"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
)

var (
	TypeType      Type = &exprType{name: "type"}
	NullType      Type = &exprType{name: "null"}
	BoolType      Type = &exprType{name: "bool"}
	IntType       Type = &exprType{name: "int"}
	UintType      Type = &exprType{name: "uint"}
	DoubleType    Type = &exprType{name: "double"}
	StringType    Type = &exprType{name: "string"}
	BytesType     Type = &exprType{name: "bytes"}
	MapType       Type = &exprType{name: "map"}
	ListType      Type = &exprType{name: "list"}
	DurationType  Type = &exprType{name: "google.protobuf.Duration"}
	TimestampType Type = &exprType{name: "google.protobuf.Timestamp"}
	DynType       Type = &exprType{name: "dyn", isDyn: true}
	// TODO: handle registration of abstract types, currently hard-coded.
)

type Type interface {
	objects.Equaler
	Name() string
	IsDyn() bool
}

func MessageType(name string) Type {
	if name == TimestampType.Name() {
		return TimestampType
	} else if name == DurationType.Name() {
		return DurationType
	}
	return &exprType{name: name}
}

type exprType struct {
	name  string
	isDyn bool
}

func (t *exprType) Name() string {
	return t.name
}

func (t *exprType) IsDyn() bool {
	return t.isDyn
}

func (t *exprType) Equal(other interface{}) bool {
	otherType, ok := other.(Type)
	return ok && t.name == otherType.Name()
}

func TypeOf(value interface{}) (Type, bool) {
	switch value.(type) {
	case Type:
		return TypeType, true
	case *tspb.Timestamp:
		return TimestampType, true
	case *dpb.Duration:
		return DurationType, true
	case bool:
		return BoolType, true
	case int64:
		return IntType, true
	case uint64:
		return UintType, true
	case float64:
		return DoubleType, true
	case string:
		return StringType, true
	case []byte:
		return BytesType, true
	case adapters.ListAdapter:
		return ListType, true
	case adapters.MapAdapter:
		return MapType, true
	case adapters.MsgAdapter:
		msgAdapter := value.(adapters.MsgAdapter)
		protoValue := msgAdapter.Value().(proto.Message)
		return MessageType(proto.MessageName(protoValue)), true
	case structpb.NullValue:
		return NullType, true
	case interface{}:
		return DynType, true
	}
	return nil, false
}

func (t *exprType) String() string {
	return fmt.Sprint("typeof(%s)", t.name)
}
