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

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v3"
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
			t.Fatalf("parser.Parse() failed: %v", iss.Err())
		}
		if p.Name().Value != tst.name {
			t.Errorf("policy name is %v, wanted %q", p.name, tst.name)
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
			err: `ERROR: <input>:2:1: unsupported policy tag: inputs
 | inputs:
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
ERROR: <input>:4:7: match does not specify a rule or output
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
			err: `ERROR: <input>:2:1: got yaml node type tag:yaml.org,2002:seq, wanted type(s) [tag:yaml.org,2002:map]
 | - rule:
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
		{
			txt: `
rule:
  match:
    - condition: "true"
      explanation: "hi"
      rule:
        match:
          - output: "hello"`,
			err: `ERROR: <input>:6:7: explanation can only be set on output match cases, not nested rules
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
      explanation: "hi"`,
			err: `ERROR: <input>:8:7: explanation can only be set on output match cases, not nested rules
 |       explanation: "hi"
 | ......^`,
		},
		{
			txt: `
imports:
  - first`,
			err: `ERROR: <input>:3:5: got yaml node type tag:yaml.org,2002:str, wanted type(s) [tag:yaml.org,2002:map]
 |   - first
 | ....^`,
		},
		{
			txt: `
imports:
  first: name`,
			err: `ERROR: <input>:3:3: got yaml node type tag:yaml.org,2002:map, wanted type(s) [tag:yaml.org,2002:seq]
 |   first: name
 | ..^`,
		},
		{
			txt: `
rule:
  - variables: name`,
			err: `ERROR: <input>:3:3: got yaml node type tag:yaml.org,2002:seq, wanted type(s) [tag:yaml.org,2002:map]
 |   - variables: name
 | ..^`,
		},
		{
			txt: `
rule:
  variables: name`,
			err: `ERROR: <input>:3:14: got yaml node type tag:yaml.org,2002:str, wanted type(s) [tag:yaml.org,2002:seq]
 |   variables: name
 | .............^`,
		},
		{
			txt: `
rule:
  variables: 
    - name`,
			err: `ERROR: <input>:4:7: got yaml node type tag:yaml.org,2002:str, wanted type(s) [tag:yaml.org,2002:map]
 |     - name
 | ......^`,
		},
		{
			txt: `
rule:
  match: 
    name: value`,
			err: `ERROR: <input>:4:5: got yaml node type tag:yaml.org,2002:map, wanted type(s) [tag:yaml.org,2002:seq]
 |     name: value
 | ....^`,
		},
		{
			txt: `
rule:
  match: 
    - name`,
			err: `ERROR: <input>:4:7: got yaml node type tag:yaml.org,2002:str, wanted type(s) [tag:yaml.org,2002:map]
 |     - name
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

func TestGetExplanationOutputPolicy(t *testing.T) {
	tst := `
rule:
  match:
    - condition: "false"
      rule:
        match:
          - condition: "1 > 2"
            output: "false"
            explanation: "'bad_inner'"
          - output: "true"
            explanation: "'good_inner'"
    - output: "true"
      explanation: "'good_outer'"
    `

	parser, err := NewParser()
	if err != nil {
		t.Fatalf("NewParser() failed: %v", err)
	}
	policy, iss := parser.Parse(StringSource(tst, "<input>"))
	if iss != nil {
		t.Fatalf("Parse() failed: %v", err)
	}

	explanationPolicy := policy.GetExplanationOutputPolicy()

	want := "'bad_inner'"
	got := explanationPolicy.Rule().Matches()[0].rule.Matches()[0].output.Value
	if got != want {
		t.Errorf("First inner output = %v, wanted %v", got, want)
	}

	want = "1 > 2"
	got = explanationPolicy.Rule().Matches()[0].rule.Matches()[0].condition.Value
	if got != want {
		t.Errorf("First inner condition = %v, wanted %v", got, want)
	}

	want = "'good_inner'"
	got = explanationPolicy.Rule().Matches()[0].rule.Matches()[1].output.Value
	if got != want {
		t.Errorf("Second inner output = %v, wanted %v", got, want)
	}

	want = "'good_outer'"
	got = explanationPolicy.Rule().Matches()[1].output.Value
	if got != want {
		t.Errorf("Second outer output = %v, wanted %v", got, want)
	}
}

type testTagHandler struct {
	defaultTagVisitor

	description string
}

func (t *testTagHandler) PolicyTag(ctx ParserContext, id int64, tagName string, node *yaml.Node, p *Policy) {
	if tagName == "description" {
		t.description = node.Value
	}
}

func TestDescriptionTag(t *testing.T) {
	tst := `name: "test"
description: |-2
   A test description.
rule:
  match:
    - condition: "true"
      output: "true"

`

	parser, err := NewParser()
	if err != nil {
		t.Fatalf("NewParser() failed: %v", err)
	}
	handler := &testTagHandler{}
	parser.TagVisitor = handler
	policy, iss := parser.Parse(StringSource(tst, "<input>"))
	if iss != nil {
		t.Fatalf("Parse() failed: %v", iss.Err())
	}

	if dx := cmp.Diff(" A test description.", handler.description); dx != "" {
		t.Errorf("handler.description (+got, -want): %s", dx)
	}

	if dx := cmp.Diff(" A test description.", policy.Description().Value); dx != "" {
		t.Errorf("policy.Description() (+got, -want): %s", dx)
	}
}
