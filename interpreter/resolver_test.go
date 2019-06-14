// Copyright 2019 Google LLC
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

package interpreter

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	structpb "github.com/golang/protobuf/ptypes/struct"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestResolver(t *testing.T) {
	vars, _ := NewActivation(map[string]interface{}{
		"a": []interface{}{
			map[string]interface{}{"b": []string{"c", "d"}},
			map[string]bool{"b": true},
		}})
	res := NewResolver(types.DefaultTypeAdapter)
	val, found := res.Resolve(vars,
		newExprVarAttribute(1, "a",
			newExprPathElem(2, types.IntZero),
			newExprPathElem(3, types.String("b"))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if !reflect.DeepEqual(val, []string{"c", "d"}) {
		t.Errorf("Got %v, wanted ['c', 'd']", val)
	}

	// Pretend that the 'val' found earlier was the result of a function call, select element 1.
	elem := res.ResolveRelative(val, vars,
		newExprRelAttribute(newExprPathElem(4, types.Int(1))))
	if elem != "d" {
		t.Errorf("Got %v, wanted 'd'", elem)
	}

	// Attempt to select element 2. This should error.
	elem = res.ResolveRelative(val, vars,
		newExprRelAttribute(newExprPathElem(5, types.Int(2))))
	if !types.IsError(elem.(ref.Val)) {
		t.Errorf("Got %v, wanted error", elem)
	}
}

func TestResolver_AnyCustomProto(t *testing.T) {
	anyExpr, err := ptypes.MarshalAny(&exprpb.ParsedExpr{Expr: &exprpb.Expr{Id: 10}})
	if err != nil {
		t.Fatal(err)
	}
	vars, _ := NewActivation(map[string]interface{}{"any": anyExpr})
	reg := types.NewRegistry(&exprpb.ParsedExpr{})
	res := NewResolver(reg)
	val, found := res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.String("expr")),
			newExprPathElem(3, types.String("id"))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if val != types.Int(10) {
		t.Errorf("Got %v, wanted 10", val)
	}
}

func TestResolver_AnyWellKnownProto(t *testing.T) {
	anyJSON, err := ptypes.MarshalAny(&structpb.ListValue{Values: []*structpb.Value{
		{Kind: &structpb.Value_BoolValue{BoolValue: false}},
		{Kind: &structpb.Value_NullValue{}},
		{Kind: &structpb.Value_NumberValue{NumberValue: 1.5}},
		{Kind: &structpb.Value_StringValue{StringValue: "hello"}},
		{Kind: &structpb.Value_ListValue{ListValue: &structpb.ListValue{}}}}})
	if err != nil {
		t.Fatal(err)
	}
	vars, _ := NewActivation(map[string]interface{}{"any": anyJSON})
	res := NewResolver(types.DefaultTypeAdapter)

	val, found := res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.Int(0))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if val != false {
		t.Errorf("Got %v, wanted false", val)
	}

	val, found = res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.Int(1))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if val != types.NullValue {
		t.Errorf("Got %v, wanted null", val)
	}

	val, found = res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.Int(2))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if val != 1.5 {
		t.Errorf("Got %v, wanted 1.5", val)
	}

	val, found = res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.Int(3))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if val != "hello" {
		t.Errorf("Got %v, wanted hello", val)
	}

	val, found = res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.Int(4))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if proto.Equal(val.(proto.Message), &structpb.ListValue{}) {
		t.Errorf("Got %v, wanted []", val)
	}

	val, found = res.Resolve(vars,
		newExprVarAttribute(1, "any",
			newExprPathElem(2, types.Int(5))))
	if !found {
		t.Error("Got not found, wanted found")
	}
	if !types.IsError(val.(ref.Val)) {
		t.Errorf("Got %v, wanted error", val)
	}
}

func TestListeningResolver(t *testing.T) {
	vars, _ := NewActivation(map[string]interface{}{
		"elems": []interface{}{
			0,
			uint(1),
			map[string]string{"str": "what is the answer to the ultimate question?"},
			map[string]int{"int": 42},
		},
	})
	resolvedAttrs := make(map[string]ResolveStatus)
	statusListener := func(attr Attribute, status ResolveStatus) {
		attrName := attr.Variable().Name()
		for _, elem := range attr.Path() {
			elemName := elem.ToValue(vars)
			attrName += fmt.Sprintf(".%v", elemName)
		}
		if _, found := resolvedAttrs[attrName]; found {
			t.Errorf("Attribute %s previously resolved.", attrName)
		}
		resolvedAttrs[attrName] = status
	}

	res := NewListeningResolver(NewResolver(types.DefaultTypeAdapter), statusListener)
	res.Resolve(vars, newExprVarAttribute(1, "elems"))
	if resolvedAttrs["elems"] != FoundAttribute {
		t.Errorf("Got %v, wanted 'elems' == FoundAttribute", resolvedAttrs)
	}

	res.Resolve(vars, newExprVarAttribute(1, "elems", newExprPathElem(2, types.Int(1))))
	if resolvedAttrs["elems.1"] != FoundAttribute {
		t.Errorf("Got %v, wanted 'elems.1' == FoundAttribute", resolvedAttrs)
	}

	res.Resolve(vars, newExprVarAttribute(1, "elems",
		newExprPathElem(3, types.Int(3)),
		newExprPathElem(4, types.String("int"))))
	if resolvedAttrs["elems.3.int"] != FoundAttribute {
		t.Errorf("Got %v, wanted 'elems.3.int' == FoundAttribute", resolvedAttrs)
	}

	res.Resolve(vars, newExprVarAttribute(404, "not_found"))
	if resolvedAttrs["not_found"] != NoSuchVariable {
		t.Errorf("Got %v, wanted 'not_found' == NoSuchVariable", resolvedAttrs)
	}

	res.Resolve(vars, newExprVarAttribute(1, "elems",
		newExprPathElem(400, types.String("no_such_key"))))
	if resolvedAttrs["elems.no_such_key"] != NoSuchAttribute {
		t.Errorf("Got %v, wanted 'elems.no_such_key' == NoSuchAttribute", resolvedAttrs)
	}
}

func TestUnknownResolver(t *testing.T) {
	unkVarsStop, _ := NewUnknownAttribute("vars", "stop")
	unkVarsAlso, _ := NewUnknownAttribute("vars", "also")
	vars, _ := NewPartialActivation(
		map[string]interface{}{
			"vars": map[string]string{
				"start": "beginning",
			},
		},
		[]Attribute{unkVarsStop, unkVarsAlso},
	)

	resolvedAttrs := make(map[string]ResolveStatus)
	statusListener := func(attr Attribute, status ResolveStatus) {
		attrName := attr.Variable().Name()
		for _, elem := range attr.Path() {
			elemName := elem.ToValue(vars)
			attrName += fmt.Sprintf(".%v", elemName)
		}
		if _, found := resolvedAttrs[attrName]; found {
			t.Errorf("Attribute %s previously resolved.", attrName)
		}
		resolvedAttrs[attrName] = status
	}
	res := NewListeningResolver(
		NewUnknownResolver(NewResolver(types.DefaultTypeAdapter)),
		statusListener)

	// When a top-level variable contains an unknown attribute, then it too becomes unknown.
	// Otherwise, the use of this value in comparisons might yield invalid results.
	res.Resolve(vars, newExprVarAttribute(1, "vars"))
	if resolvedAttrs["vars"] != UnknownAttribute {
		t.Errorf("Got %v, wanted 'vars' == UnknownAttribute", resolvedAttrs)
	}

	// Fully qualified attributes with known concrete values are still known however.
	res.Resolve(vars, newExprVarAttribute(1, "vars",
		newExprPathElem(2, types.String("start"))))
	if resolvedAttrs["vars.start"] != FoundAttribute {
		t.Errorf("Got %v, wanted 'vars.start' == FoundAttribute", resolvedAttrs)
	}

	// Note, a partially known map potentially poses problems with respect to map equality and
	// containment tests. At present only top-level unknown variables are recommended.
	res.Resolve(vars, newExprVarAttribute(1, "vars",
		newExprPathElem(3, types.String("stop"))))
	if resolvedAttrs["vars.stop"] != UnknownAttribute {
		t.Errorf("Got %v, wanted 'vars.stop' == UnknownAttribute", resolvedAttrs)
	}

	// Selecting below where an unknown is declared should also result in an unknown attribute
	// which extends only as far as the path where the first unknown occurs.
	res.Resolve(vars, newExprVarAttribute(1, "vars",
		newExprPathElem(4, types.String("also")),
		newExprPathElem(5, types.Int(12))))
	if resolvedAttrs["vars.also"] != UnknownAttribute {
		t.Errorf("Got %v, wanted 'vars.also' == UnknownAttribute", resolvedAttrs)
	}
}
