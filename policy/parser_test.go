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
	for _, tst := range policyTests {
		srcFile := readPolicy(t, fmt.Sprintf("testdata/%s/policy.yaml", tst.name))
		parser, err := NewParser(tst.parseOpts...)
		if err != nil {
			t.Fatalf("NewParser() failed: %v", err)
		}
		p, iss := parser.Parse(srcFile)
		if iss.Err() != nil {
			t.Fatalf("parse() failed: %v", iss.Err())
		}
		if p.Name().Value != tst.name {
			t.Errorf("policy name is %v, wanted 'required_labels'", p.name)
		}
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		txt string
		err string
	}{
		{
			txt: `
name:
  illegal: yaml-type`,
			err: `ERROR: <input>:3:3: got yaml node type tag:yaml.org,2002:map, wanted type(s) [tag:yaml.org,2002:str !txt]
 |   illegal: yaml-type
 | ..^`,
		},
		{
			txt: `
rule:
  custom: yaml-type`,
			err: `ERROR: <input>:3:3: unsupported rule tag: custom
 |   custom: yaml-type
 | ..^`,
		},
		{
			txt: `
inputs:
  - name: a
  - name: b`,
			err: `ERROR: <input>:1:2: unsupported policy tag: inputs
 | 
 | ^`,
		},
		{
			txt: `
rule:
  variables:
    - name: "true"
      alt_name: "bool_true"`,
			err: `ERROR: <input>:5:7: unsupported variable tag: alt_name
 |       alt_name: "bool_true"
 | ......^`,
		},
		{
			txt: `
rule:
  match:
    - name: "true"
      alt_name: "bool_true"`,
			err: `ERROR: <input>:4:7: unsupported match tag: name
 |     - name: "true"
 | ......^
ERROR: <input>:5:7: unsupported match tag: alt_name
 |       alt_name: "bool_true"
 | ......^`,
		},
		{
			txt: `
- rule:
    id: a`,
			err: `ERROR: <input>:1:2: got yaml node type tag:yaml.org,2002:seq, wanted type(s) [tag:yaml.org,2002:map]
 | 
 | ^`,
		},
		{
			txt: `
rule:
  match:
    - condition: "true"
      output: "world"
      rule:
        match:
          - output: "hello"`,
			err: `ERROR: <input>:6:7: only the rule or the output may be set
 |       rule:
 | ......^`,
		},
		{
			txt: `
rule:
  match:
    - condition: "true"
      rule:
        match:
          - output: "hello"
      output: "world"`,
			err: `ERROR: <input>:8:7: only the rule or the output may be set
 |       output: "world"
 | ......^`,
		},
	}

	for _, tst := range tests {
		parser, err := NewParser()
		if err != nil {
			t.Fatalf("NewParser() failed: %v", err)
		}
		_, iss := parser.Parse(StringSource(tst.txt, "<input>"))
		if iss.Err() == nil {
			t.Fatalf("parser.Parse(%q) did not error, wanted %s", tst.txt, tst.err)
		}

		if iss.Err().Error() != tst.err {
			t.Errorf("parser.Parse(%q) got error %v, wanted error %s", tst.txt, iss.Err(), tst.err)
		}
	}
}
