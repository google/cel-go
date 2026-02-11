// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestYAMLHelper(t *testing.T) {
	helper := YAMLHelper{}
	tests := []struct {
		name string
		yaml string
		want func(*testing.T, *yaml.Node)
	}{
		{
			name: "list",
			yaml: `
- first
- second
- third`,
			want: func(t *testing.T, node *yaml.Node) {
				if !helper.IsList(node) {
					t.Fatalf("IsList(node) returned false, want true")
				}

				helper.RangeList(node, func(i int, val *yaml.Node) bool {
					switch i {
					case 0:
						if !helper.IsString(val) || val.Value != "first" {
							t.Errorf("val[%d] = %q, want %q", i, val.Value, "first")
						}
					case 1:
						if !helper.IsString(val) || val.Value != "second" {
							t.Errorf("val[%d] = %q, want %q", i, val.Value, "second")
						}
					case 2:
						if !helper.IsString(val) || val.Value != "third" {
							t.Errorf("val[%d] = %q, want %q", i, val.Value, "third")
						}
					}
					return true
				})
			},
		},
		{
			name: "map",
			yaml: `
first: 1
second: 2
third: 3`,
			want: func(t *testing.T, node *yaml.Node) {
				if !helper.IsMap(node) {
					t.Fatalf("IsMap(node) returned false, want true")
				}
				helper.RangeMap(node, func(key *yaml.Node, val *yaml.Node) bool {
					if key.Value == "first" && helper.IsInteger(val) && val.Value != "1" {
						t.Errorf("val[%q] = %q, want %q", key.Value, val.Value, "1")
					}
					if key.Value == "second" && helper.IsInteger(val) && val.Value != "2" {
						t.Errorf("val[%q] = %q, want %q", key.Value, val.Value, "2")
					}
					if key.Value == "third" && helper.IsInteger(val) && val.Value != "3" {
						t.Errorf("val[%q] = %q, want %q", key.Value, val.Value, "3")
					}
					return true
				})
			},
		},
		{
			name: "map_with_empty_value",
			yaml: `
first: 1
second:
third: 3`,
			want: func(t *testing.T, node *yaml.Node) {
				if !helper.IsMap(node) {
					t.Errorf("IsMap(node) returned false, want true")
				}
				helper.RangeMap(node, func(key *yaml.Node, val *yaml.Node) bool {
					if key.Value == "second" && !helper.IsNull(val) {
						t.Errorf("val[%q] = %v, want null", key.Value, val)
					}
					return true
				})
			},
		},
		{
			name: "list_with_mixed_values",
			yaml: `
- 1
- 'hello'
- 1.5
- true
- null
- 2006-01-02T15:04:05Z`,
			want: func(t *testing.T, node *yaml.Node) {
				if !helper.IsList(node) {
					t.Errorf("IsList(node) returned false, want true")
				}
				helper.RangeList(node, func(i int, val *yaml.Node) bool {
					switch i {
					case 0:
						if !helper.IsInteger(val) || !helper.IsNumber(val) {
							t.Errorf("val[%d] = %v, want integer", i, val)
						}
					case 1:
						if !helper.IsString(val) {
							t.Errorf("val[%d] = %v, want string", i, val)
						}
					case 2:
						if !helper.IsDouble(val) && !helper.IsNumber(val) {
							t.Errorf("val[%d] = %v, want double", i, val)
						}
					case 3:
						if !helper.IsBool(val) {
							t.Errorf("val[%d] = %v, want bool", i, val)
						}
					case 4:
						if !helper.IsNull(val) {
							t.Errorf("val[%d] = %v, want null", i, val)
						}
					case 5:
						if !helper.IsTimestamp(val) {
							t.Errorf("val[%d] = %v, want timestamp", i, val)
						}
					}
					return true
				})
			},
		},
		{
			name: "list_string_failure",
			yaml: `- 1`,
			want: func(t *testing.T, node *yaml.Node) {
				if helper.IsMap(node) || helper.IsBool(node) || helper.IsDouble(node) || helper.IsInteger(node) || helper.IsNull(node) ||
					helper.IsNumber(node) || helper.IsString(node) || helper.IsTimestamp(node) {
					t.Errorf("got %v, wanted list", node)
				}
				helper.RangeList(node, func(_ int, elem *yaml.Node) bool {
					if !helper.IsString(elem) {
						return false
					}
					t.Fatal("got string node, wanted non-string node")
					return true
				})
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var node yaml.Node
			err := yaml.Unmarshal([]byte(tc.yaml), &node)
			if err != nil {
				t.Fatalf("yaml.Unmarshal() failed: %v", err)
			}
			tc.want(t, node.Content[0])
		})
	}
}
