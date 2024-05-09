// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	parser := NewParser()
	for _, tst := range policyTests {
		srcFile := readPolicy(t, fmt.Sprintf("testdata/%s/policy.yaml", tst.name))
		p, iss := parser.Parse(srcFile)
		if iss.Err() != nil {
			t.Fatalf("parse() failed: %v", iss.Err())
		}
		if p.name.Value != tst.name {
			t.Errorf("policy name is %v, wanted 'required_labels'", p.name)
		}
	}
}
