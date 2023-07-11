// Copyright 2023 Google LLC
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
	"fmt"
	"testing"

	"github.com/google/cel-go/common/types/ref"
)

func TestIsUnknown(t *testing.T) {
	if IsUnknown(&Unknown{}) != true {
		t.Error("IsUnknown(Unknown{}) returned false, wanted true")
	}
	if IsUnknown(Bool(true)) != false {
		t.Error("IsUnknown(true) returned true, wanted false")
	}
}

func TestNewAttribute(t *testing.T) {
	if NewAttributeTrail("") != unspecifiedAttribute {
		t.Error("An empty attribute must be the unspecified attribute")
	}
	if NewAttributeTrail("v").Equal(unspecifiedAttribute) {
		t.Error("A non-empty attribute must not be equal to the unspecified attribute")
	}
}

func TestAttributeEquals(t *testing.T) {
	tests := []struct {
		a     *AttributeTrail
		b     *AttributeTrail
		equal bool
	}{
		{
			a:     unspecifiedAttribute,
			b:     NewAttributeTrail(""),
			equal: true,
		},
		{
			a:     NewAttributeTrail("a"),
			b:     NewAttributeTrail(""),
			equal: false,
		},
		{
			a:     NewAttributeTrail("a"),
			b:     NewAttributeTrail("a"),
			equal: true,
		},
		{
			a:     QualifyAttribute[string](NewAttributeTrail("a"), "b"),
			b:     NewAttributeTrail("a"),
			equal: false,
		},
		{
			a:     QualifyAttribute[string](NewAttributeTrail("a"), "b"),
			b:     QualifyAttribute[int64](NewAttributeTrail("a"), 1),
			equal: false,
		},
		{
			a:     QualifyAttribute[string](NewAttributeTrail("a"), "b"),
			b:     QualifyAttribute[string](NewAttributeTrail("a"), "b"),
			equal: true,
		},
	}
	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := tc.a.Equal(tc.b)
			if out != tc.equal {
				t.Errorf("%v.Equal(%v) got %v, wanted %v", tc.a, tc.b, out, tc.equal)
			}
		})
	}
}

func TestAttributeString(t *testing.T) {
	tests := []struct {
		attr *AttributeTrail
		out  string
	}{
		{
			attr: unspecifiedAttribute,
			out:  "<unspecified>",
		},
		{
			attr: NewAttributeTrail("a"),
			out:  "a",
		},
		{
			attr: QualifyAttribute[bool](NewAttributeTrail("a"), false),
			out:  "a[false]",
		},
		{
			attr: QualifyAttribute[string](NewAttributeTrail("a"), "b"),
			out:  "a.b",
		},
		{
			attr: QualifyAttribute(QualifyAttribute[string](NewAttributeTrail("a"), "b"), "$this"),
			out:  `a.b["$this"]`,
		},
		{
			attr: QualifyAttribute[int64](NewAttributeTrail("a"), 12),
			out:  "a[12]",
		},
		{
			attr: QualifyAttribute[uint64](NewAttributeTrail("a"), 24),
			out:  "a[24u]",
		},
	}
	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := tc.attr.String()
			if out != tc.out {
				t.Errorf("%v.String() got %v, wanted %v", tc.attr, out, tc.out)
			}
		})
	}
}

func TestUnknownContains(t *testing.T) {
	tests := []struct {
		unk   *Unknown
		other *Unknown
		out   bool
	}{
		{
			unk:   NewUnknown(1, nil),
			other: NewUnknown(1, unspecifiedAttribute),
			out:   true,
		},
		{
			unk:   NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), true)),
			other: NewUnknown(4, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			out:   false,
		},
		{
			unk:   NewUnknown(3, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			other: NewUnknown(4, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			out:   false,
		},
		{
			unk:   NewUnknown(3, QualifyAttribute[string](NewAttributeTrail("a"), "c")),
			other: NewUnknown(3, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			out:   false,
		},
		{
			unk: MergeUnknowns(
				NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), true)),
				NewUnknown(4, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			),
			other: NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), true)),
			out:   true,
		},
		{
			unk: NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), true)),
			other: MergeUnknowns(
				NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), true)),
				NewUnknown(4, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			),
			out: false,
		},
	}
	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := tc.unk.Contains(tc.other)
			if out != tc.out {
				t.Errorf("%v.Contains(%v) got %v, wanted %v", tc.unk, tc.other, out, tc.out)
			}
		})
	}
}

func TestUnknownString(t *testing.T) {
	tests := []struct {
		unk *Unknown
		out string
	}{
		{
			unk: NewUnknown(1, nil),
			out: "<unspecified> (1)",
		},
		{
			unk: NewUnknown(1, unspecifiedAttribute),
			out: "<unspecified> (1)",
		},
		{
			unk: NewUnknown(2, NewAttributeTrail("a")),
			out: "a (2)",
		},
		{
			unk: NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), false)),
			out: "a[false] (3)",
		},
		{
			unk: MergeUnknowns(
				NewUnknown(3, QualifyAttribute[bool](NewAttributeTrail("a"), true)),
				NewUnknown(4, QualifyAttribute[string](NewAttributeTrail("a"), "b")),
			),
			out: "a[true] (3), a.b (4)",
		},
	}
	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := tc.unk.String()
			if out != tc.out {
				t.Errorf("%v.String() got %v, wanted %v", tc.unk, out, tc.out)
			}
		})
	}
}

func TestMaybeMergeUnknowns(t *testing.T) {
	tests := []struct {
		in    ref.Val
		unk   *Unknown
		want  *Unknown
		isUnk bool
	}{
		// one of the unknowns is empty
		{
			in:    String(""),
			unk:   nil,
			isUnk: false,
		},
		// both unknowns are empty
		{
			in:    String(""),
			unk:   &Unknown{},
			want:  &Unknown{},
			isUnk: true,
		},
		{
			in:    newUnk(t, 2, "x"),
			unk:   newUnk(t, 1, "y"),
			want:  MergeUnknowns(newUnk(t, 2, "x"), newUnk(t, 1, "y")),
			isUnk: true,
		},
		{
			in:    newUnk(t, 2, "x"),
			want:  newUnk(t, 2, "x"),
			isUnk: true,
		},
	}

	for i, tst := range tests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			got, isUnk := MaybeMergeUnknowns(tc.in, tc.unk)
			if isUnk != tc.isUnk {
				t.Errorf("MaybeMergeUnknowns(%v, %v) got %v, wanted %v", tc.in, tc.unk, got, tc.want)
			}
		})
	}
}

func newUnk(t *testing.T, id int64, varName string) *Unknown {
	t.Helper()
	attr := NewAttributeTrail(varName)
	unk := NewUnknown(id, attr)
	return unk
}
