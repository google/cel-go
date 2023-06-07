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
	if IsUnknown(Unknown{}) != true {
		t.Error("IsUnknown(Unknown{}) returned false, wanted true")
	}
	if IsUnknown(Bool(true)) != false {
		t.Error("IsUnknown(true) returned true, wanted false")
	}
}

func TestMaybeMergeUnknowns(t *testing.T) {
	tests := []struct {
		in    ref.Val
		unk   Unknown
		want  Unknown
		isUnk bool
	}{
		{
			in:    String(""),
			unk:   nil,
			want:  nil,
			isUnk: false,
		},
		{
			in:    String(""),
			unk:   Unknown{},
			want:  Unknown{},
			isUnk: true,
		},
		{
			in:    Unknown{2},
			unk:   Unknown{1},
			want:  Unknown{1, 2},
			isUnk: true,
		},
		{
			in:    Unknown{2},
			unk:   nil,
			want:  Unknown{2},
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
