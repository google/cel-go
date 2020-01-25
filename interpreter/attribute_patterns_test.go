// Copyright 2020 Google LLC
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

	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
)

type attr struct {
	maybe bool
	name  string
	quals []interface{}
}

type patternTest struct {
	pattern   *AttributePattern
	container string
	matches   []attr
	misses    []attr
}

var patternTests = map[string]patternTest{
	"var": {
		pattern: NewAttributePattern("var"),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{"field"}},
		},
		misses: []attr{
			{name: "ns.var"},
		},
	},
	"var_namespace": {
		pattern:   NewAttributePattern("ns.app.var"),
		container: "ns.app",
		matches: []attr{
			{name: "ns.app.var"},
			{name: "ns.app.var", quals: []interface{}{int64(0)}},
			{name: "ns", quals: []interface{}{"app", "var", "foo"}, maybe: true},
		},
		misses: []attr{
			{name: "ns.var"},
			{name: "ns", quals: []interface{}{"var"}, maybe: true},
		},
	},
	"var_field": {
		pattern: NewAttributePattern("var").QualString("field"),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{"field"}},
			{name: "var", quals: []interface{}{"field"}, maybe: true},
			{name: "var", quals: []interface{}{"field", uint64(1)}},
		},
		misses: []attr{
			{name: "var", quals: []interface{}{"other"}},
		},
	},
	"var_index": {
		pattern: NewAttributePattern("var").QualInt(0),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{int64(0)}},
			{name: "var", quals: []interface{}{int64(0), false}},
		},
		misses: []attr{
			{name: "var", quals: []interface{}{uint64(0)}},
			{name: "var", quals: []interface{}{int64(1), false}},
		},
	},
	"var_index_uint": {
		pattern: NewAttributePattern("var").QualUint(1),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{uint64(1)}},
			{name: "var", quals: []interface{}{uint64(1), true}},
		},
		misses: []attr{
			{name: "var", quals: []interface{}{uint64(0)}},
			{name: "var", quals: []interface{}{int64(1), false}},
		},
	},
	"var_index_bool": {
		pattern: NewAttributePattern("var").QualBool(true),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{true}},
			{name: "var", quals: []interface{}{true, "name"}},
		},
		misses: []attr{
			{name: "var", quals: []interface{}{false}},
			{name: "none"},
		},
	},
	"var_wildcard": {
		pattern:   NewAttributePattern("ns.var").Wildcard(),
		container: "ns",
		matches: []attr{
			{name: "ns.var"},
			// The maybe attributes consider potential namespacing and field selection
			// when testing variable names.
			{name: "var", quals: []interface{}{true}, maybe: true},
			{name: "var", quals: []interface{}{"name"}, maybe: true},
			{name: "var", quals: []interface{}{"name"}, maybe: true},
		},
		misses: []attr{
			{name: "var", quals: []interface{}{false}},
			{name: "none"},
		},
	},
	"var_wildcard_field": {
		pattern: NewAttributePattern("var").Wildcard().QualString("field"),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{true}},
			{name: "var", quals: []interface{}{int64(10), "field"}},
		},
		misses: []attr{
			{name: "var", quals: []interface{}{int64(10), "other"}},
		},
	},
	"var_wildcard_wildcard": {
		pattern: NewAttributePattern("var").Wildcard().Wildcard(),
		matches: []attr{
			{name: "var"},
			{name: "var", quals: []interface{}{true}},
			{name: "var", quals: []interface{}{int64(10), "field"}},
		},
		misses: []attr{
			{name: "none"},
		},
	},
}

func TestAttributePattern_UnknownResolution(t *testing.T) {
	reg := types.NewRegistry()
	for nm, tc := range patternTests {
		tst := tc
		t.Run(nm, func(tt *testing.T) {
			pkg := packages.DefaultPackage
			if tst.container != "" {
				pkg = packages.NewPackage(tst.container)
			}
			fac := NewPartialAttributeFactory(pkg, reg, reg)
			for i, match := range tst.matches {
				m := match
				tt.Run(fmt.Sprintf("match[%d]", i), func(ttt *testing.T) {
					attr := genAttr(fac, m)
					partVars, _ := NewPartialActivation(EmptyActivation(), tst.pattern)
					val, err := attr.Resolve(partVars)
					if err != nil {
						ttt.Fatalf("Got error: %s, wanted unknown", err)
					}
					_, isUnk := val.(types.Unknown)
					if !isUnk {
						ttt.Fatalf("Got value %v, wanted unknown", val)
					}
				})
			}
			for i, miss := range tst.misses {
				m := miss
				tt.Run(fmt.Sprintf("miss[%d]", i), func(ttt *testing.T) {
					attr := genAttr(fac, m)
					partVars, _ := NewPartialActivation(EmptyActivation(), tst.pattern)
					val, err := attr.Resolve(partVars)
					if err == nil {
						ttt.Fatalf("Got value: %s, wanted error", val)
					}
				})
			}
		})
	}
}

func TestAttributePattern_CrossReference(t *testing.T) {
	reg := types.NewRegistry()
	fac := NewPartialAttributeFactory(packages.DefaultPackage, reg, reg)
	a := fac.AbsoluteAttribute(1, "a")
	b := fac.AbsoluteAttribute(2, "b")
	a.AddQualifier(b)

	// Ensure that var a[b], the dynamic index into var 'a' is the unknown value
	// returned from attribute resolution.
	partVars, _ := NewPartialActivation(
		map[string]interface{}{"a": []int64{1, 2}},
		NewAttributePattern("b"))
	val, err := a.Resolve(partVars)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, types.Unknown{2}) {
		t.Fatalf("Got %v, wanted unknown attribute id for 'b' (2)", val)
	}

	// Ensure that a[b], the dynamic index into var 'a' is the unknown value
	// returned from attribute resolution. Note, both 'a' and 'b' have unknown attribute
	// patterns specified. This changes the evaluation behavior slightly, but the end
	// result is the same.
	partVars, _ = NewPartialActivation(
		map[string]interface{}{"a": []int64{1, 2}},
		NewAttributePattern("a").QualInt(0),
		NewAttributePattern("b"))
	val, err = a.Resolve(partVars)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, types.Unknown{2}) {
		t.Fatalf("Got %v, wanted unknown attribute id for 'b' (2)", val)
	}

	// Note, that only 'a[0].c' will result in an unknown result since both 'a' and 'b'
	// have values. However, since the attribute being pattern matched is just 'a.b',
	// the outcome will indicate that 'a[b]' is unknown.
	partVars, _ = NewPartialActivation(
		map[string]interface{}{"a": []int64{1, 2}, "b": 0},
		NewAttributePattern("a").QualInt(0).QualString("c"))
	val, err = a.Resolve(partVars)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, types.Unknown{2}) {
		t.Fatalf("Got %v, wanted unknown attribute id for 'b' (2)", val)
	}

	// Test a positive case that returns a valid value even though the attribugte factory
	// is the partial attribute factory.
	partVars, _ = NewPartialActivation(
		map[string]interface{}{"a": []int64{1, 2}, "b": 0})
	val, err = a.Resolve(partVars)
	if err != nil {
		t.Fatal(err)
	}
	if val != int64(1) {
		t.Fatalf("Got %v, wanted 1 for a[b]", val)
	}

	// Ensure the unknown attribute id moves when the attribute becomes more specific.
	partVars, _ = NewPartialActivation(
		map[string]interface{}{"a": []int64{1, 2}, "b": 0},
		NewAttributePattern("a").QualInt(0).QualString("c"))
	// Qualify a[b] with 'c', a[b].c
	c, _ := fac.NewQualifier(nil, 3, "c")
	a.AddQualifier(c)
	// The resolve step should return unkonwn
	val, err = a.Resolve(partVars)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(val, types.Unknown{3}) {
		t.Fatalf("Got %v, wanted unknown attribute id for a[b].c (3)", val)
	}
}

func genAttr(fac AttributeFactory, a attr) Attribute {
	id := int64(1)
	var attr Attribute
	if a.maybe {
		attr = fac.MaybeAttribute(1, a.name)
	} else {
		attr = fac.AbsoluteAttribute(1, a.name)
	}
	for _, q := range a.quals {
		qual, _ := fac.NewQualifier(nil, id, q)
		attr.AddQualifier(qual)
		id++
	}
	return attr
}
