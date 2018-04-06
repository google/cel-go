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
	"github.com/golang/protobuf/ptypes"
	testExpr "github.com/google/cel-go/interpreter/testing"
	"reflect"
	"testing"
)

func TestObjectValue_Get(t *testing.T) {
	existsMsg := NewProtoValue(testExpr.Exists.Expr)
	compreExpr, err := existsMsg.Get("ComprehensionExpr")
	if err != nil {
		t.Error(err)
	}
	iterVar, err := compreExpr.(ObjectValue).Get("IterVar")
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
	field, err := selectExpr.(ObjectValue).Get("Field")
	if err != nil {
		t.Error(err)
	}
	if field != "" {
		t.Error("Selected field on unset message, but field was non-empty")
	}
	field2, err := selectExpr.(ObjectValue).Get("Field")
	if field2 != field {
		t.Error("Selected cached field not equal to original field")
	}
}

func TestObjectValue_Any(t *testing.T) {
	anyVal, err := ptypes.MarshalAny(testExpr.Exists.Expr)
	if err != nil {
		t.Error(err)
	}
	existsMsg := NewProtoValue(anyVal)
	compre, err := existsMsg.Get("ComprehensionExpr")
	if err != nil {
		t.Error(err)
	}
	if compre == nil {
		t.Error("Comprehension was null")
	}
}

func TestObjectValue_Iterator(t *testing.T) {
	existsMsg := NewProtoValue(testExpr.Exists.Expr)
	it := existsMsg.Iterator()
	var fields []interface{}
	for it.HasNext() {
		fields = append(fields, it.Next())
	}
	if !reflect.DeepEqual(fields, []interface{}{"ComprehensionExpr", "Id"}) {
		t.Errorf("Got %v, wanted %v", fields, []interface{}{"ComprehensionExpr", "Id"})
	}
}
