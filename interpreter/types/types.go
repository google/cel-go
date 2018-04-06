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

// Package types declares core language type interfaces with wrappers for
// protocol buffer values.
package types

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/cel-go/interpreter/types/traits"
)

// Type instances and their human readable names.
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

// Types may be compared with each other, have a name, and indicate whether
// the underlying type is dynamic.
type Type interface {
	traits.Equaler
	Name() string
	IsDyn() bool
}

// MessageType produces a new Type instance for a qualified message name.
func MessageType(name string) Type {
	if name == TimestampType.Name() {
		return TimestampType
	}
	if name == DurationType.Name() {
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

// TypeOf returns the CEL Type for the input value, or false if the type
// cannot be found.
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
	case ListValue:
		return ListType, true
	case MapValue:
		return MapType, true
	case ObjectValue:
		msgAdapter := value.(ObjectValue)
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
	return fmt.Sprintf("typeof(%s)", t.name)
}
