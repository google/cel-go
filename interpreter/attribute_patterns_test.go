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
	pattern *AttributePattern
	pkg     string
	matches []attr
	misses  []attr
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
	"ns.var": {
		pattern: NewAttributePattern("ns.app.var"),
		pkg:     "ns.app",
		matches: []attr{
			{name: "ns.app.var"},
			{name: "ns.app.var", quals: []interface{}{int64(0)}},
			{name: "ns", quals: []interface{}{"app", "var"}, maybe: true},
		},
		misses: []attr{
			{name: "ns.var"},
			{name: "ns", quals: []interface{}{"var"}, maybe: true},
		},
	},
}

func TestAttributePattern(t *testing.T) {
	reg := types.NewRegistry()
	for nm, tc := range patternTests {
		tst := tc
		t.Run(nm, func(tt *testing.T) {
			pkg := packages.DefaultPackage
			if tst.pkg != "" {
				pkg = packages.NewPackage(tst.pkg)
			}
			fac := NewPartialAttributeFactory(pkg, reg, reg)
			for i, match := range tst.matches {
				m := match
				tt.Run(fmt.Sprintf("[%d]", i), func(ttt *testing.T) {
					id := int64(1)
					var attr Attribute
					if m.maybe {
						attr = fac.MaybeAttribute(1, m.name)
					} else {
						attr = fac.AbsoluteAttribute(1, m.name)
					}
					for _, q := range m.quals {
						qual, _ := fac.NewQualifier(nil, id, q)
						attr.AddQualifier(qual)
						id++
					}
					partVars, _ := NewPartialActivation(EmptyActivation(), tst.pattern)
					val, err := attr.Resolve(partVars)
					if err != nil {
						t.Fatalf("Got error: %s, wanted unknown", err)
					}
					_, isUnk := val.(types.Unknown)
					if !isUnk {
						t.Fatalf("Got value %v, wanted unknown", val)
					}
				})
			}
		})
	}
}
