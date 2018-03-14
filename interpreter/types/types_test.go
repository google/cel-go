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
	"github.com/google/cel-go/interpreter/types/adapters"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/cel-spec/proto/v1"
	"testing"
)

func TestTypeOf_Type(t *testing.T) {
	if val, found := TypeOf(IntType); !found || val != TypeType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Null(t *testing.T) {
	null := structpb.NullValue_NULL_VALUE
	if val, found := TypeOf(null); !found || val != NullType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Timestamp(t *testing.T) {
	ts := timestamp.Timestamp{}
	if val, found := TypeOf(&ts); !found || val != TimestampType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Duration(t *testing.T) {
	dur := duration.Duration{}
	if val, found := TypeOf(&dur); !found || val != DurationType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Primitive(t *testing.T) {
	if val, found := TypeOf(true); !found || val != BoolType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
	if val, found := TypeOf(int64(1)); !found || val != IntType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
	if val, found := TypeOf(uint64(10000000)); !found || val != UintType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
	if val, found := TypeOf(float64(-1.2)); !found || val != DoubleType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
	if val, found := TypeOf(""); !found || val != StringType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Bytes(t *testing.T) {
	if val, found := TypeOf([]byte("hello")); !found || val != BytesType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Map(t *testing.T) {
	if val, found := TypeOf(adapters.NewMapAdapter(map[string]int32{"0": 1, "1": 2, "2": 3})); !found || val != MapType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_List(t *testing.T) {
	if val, found := TypeOf(adapters.NewListAdapter([]int32{1, 2, 3})); !found || val != ListType {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_Proto(t *testing.T) {
	expr := adapters.NewMsgAdapter(&syntax_proto.Expr{})
	if val, found := TypeOf(expr); !found ||
		!val.Equal(MessageType("google.api.expr.v1.Expr")) {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestTypeOf_NotFound(t *testing.T) {
	if val, found := TypeOf(1); !found || !val.Equal(DynType) {
		t.Errorf("Unexpected type result. found: %t, type: %v", found, val)
	}
}

func TestExprType_Equal(t *testing.T) {
	dur := &duration.Duration{}
	durType, _ := TypeOf(dur)
	durProtoType := MessageType("google.protobuf.Duration")
	if !durType.Equal(durProtoType) {
		t.Errorf("Type of duration proto by instance and by name not equal")
	}

	ts := &timestamp.Timestamp{}
	tsType, _ := TypeOf(ts)
	tsProtoType := MessageType("google.protobuf.Timestamp")
	if !tsType.Equal(tsProtoType) {
		t.Errorf("Type of timestamp proto by instance and by name not equal")
	}
	if durType.Equal(tsType) {
		t.Errorf("Duration and proto types should not be equal.")
	}
}
