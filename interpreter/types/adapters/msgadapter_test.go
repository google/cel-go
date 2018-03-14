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

package adapters

import (
	testExpr "github.com/google/cel-go/interpreter/testing"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"testing"
)

func TestMsgAdapter_Get(t *testing.T) {
	existsMsg := NewMsgAdapter(testExpr.Exists.Expr)
	compreExpr, err := existsMsg.Get("ComprehensionExpr")
	if err != nil {
		t.Error(err)
	}
	iterVar, err := compreExpr.(MsgAdapter).Get("IterVar")
	if err != nil {
		t.Error(err)
	}
	if iterVar != "x" {
		t.Error("Could not retrieve iter var from comprehension")
	}
	// This field is not set, but should return the default instance.
	selectExpr, err := existsMsg.Get("SelectExpr")
	if err != nil {
		t.Error(err)
	}
	field, err := selectExpr.(MsgAdapter).Get("Field")
	if field != "" {
		t.Error("Selected field on unset message, but field was non-empty")
	}
	field2, err := selectExpr.(MsgAdapter).Get("Field")
	if field2 != field {
		t.Error("Selected cached field not equal to original field")
	}
}

func TestMsgAdapter_Any(t *testing.T) {
	anyVal, err := ptypes.MarshalAny(testExpr.Exists.Expr)
	if err != nil {
		t.Error(err)
	}
	existsMsg := NewMsgAdapter(anyVal)
	compre, err := existsMsg.Get("ComprehensionExpr")
	if err != nil {
		t.Error(err)
	}
	if compre == nil {
		t.Error("Comprehension was null")
	}
}

func TestMsgAdapter_Iterator(t *testing.T) {
	existsMsg := NewMsgAdapter(testExpr.Exists.Expr)
	it := existsMsg.Iterator()
	for it.HasNext() {
		fmt.Println(it.Next())
	}
}
