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
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"

	"gopkg.in/yaml.v3"

	proto3pb "github.com/google/cel-go/test/proto3pb"
)

func TestConfig(t *testing.T) {
	tests := []string{
		`name: hello`,
		`description: empty`,
		`container: pb.pkg`,
		`
extensions:
  - name: "bindings"
  - name: "encoders"
  - name: "lists"
  - name: "math"
  - name: "optional"
  - name: "protos"
  - name: "sets"
  - name: "strings"
    version: 1`,
		`
functions:
  - name: "coalesce"
    overloads:
      - id: "null_coalesce_int"
        target:
          type_name: "null_type"
        args:
          - type_name: "int"
        return:
          type_name: "int"
      - id: "coalesce_null_int"
        args:
          - type_name: "null_type"
          - type_name: "int"
        return:
          type_name: "int"          
      - id: "int_coalesce_int"
        target: 
          type_name: "int"
        args:
          - type_name: "int"
        return: 
          type_name: "int"
      - id: "optional_T_coalesce_T"
        target: 
          type_name: "optional_type"
          params:
            - type_name: "T"
              is_type_param: true
        args:
          - type_name: "T"
            is_type_param: true
        return: 
          type_name: "T"
          is_type_param: true  
`,
		`
variables:
- name: "request"
  type:
    type_name: "map"
    params:
      - type_name: "string"
      - type_name: "dyn"
`,
		`
variables:
- name: "request"
  type:
    type_name: "google.expr.proto3.test.TestAllTypes"
`,
	}
	baseEnv, err := cel.NewEnv(
		cel.OptionalTypes(),
		cel.Types(&proto3pb.TestAllTypes{}),
	)
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	for _, tst := range tests {
		c := parseConfigYaml(t, tst)
		_, err := c.AsEnvOptions(baseEnv)
		if err != nil {
			t.Errorf("AsEnvOptions() generated error: %v", err)
		}
	}
}

func TestConfigErrors(t *testing.T) {
	tests := []struct {
		config string
		err    string
	}{
		{
			config: `
extensions:
  - name: "bad_name"`,
			err: "unrecognized extension: bad_name",
		},
		{
			config: `
variables:
  - name: "bad_type"
    type:
      type_name: "strings"`,
			err: "undefined type name: strings",
		},
		{
			config: `
variables:
  - name: "bad_list"
    type:
      type_name: "list"`,
			err: "list type has unexpected param count: 0",
		},
		{
			config: `
variables:
  - name: "bad_map"
    type:
      type_name: "map"
      params:
        - type_name: "string"`,
			err: "map type has unexpected param count: 1",
		},
		{
			config: `
variables:
  - name: "bad_list_type_param"
    type:
      type_name: "list"
      params:
        - type_name: "number"`,
			err: "undefined type name: number",
		},
		{
			config: `
variables:
  - name: "bad_map_type_param"
    type:
      type_name: "map"
      params:
        - type_name: "string"
        - type_name: "optional"`,
			err: "undefined type name: optional",
		},
		{
			config: `
functions:
  - name: "bad_return"
    overloads:
      - id: "zero_arity"
        return:
          type_name: "mystery"`,
			err: "undefined type name: mystery",
		},
		{
			config: `
functions:
  - name: "bad_target"
    overloads:
      - id: "unary_member"
        target:
          type_name: "unknown"
        return:
          type_name: "null_type"`,
			err: "undefined type name: unknown",
		},
		{
			config: `
functions:
  - name: "bad_arg"
    overloads:
      - id: "unary_global"
        args:
          - type_name: "unknown"
        return:
          type_name: "null_type"`,
			err: "undefined type name: unknown",
		},
	}
	baseEnv, err := cel.NewEnv(cel.OptionalTypes())
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	for _, tst := range tests {
		c := parseConfigYaml(t, tst.config)
		_, err := c.AsEnvOptions(baseEnv)
		if err == nil || err.Error() != tst.err {
			t.Errorf("AsEnvOptions() got error: %v, wanted %s", err, tst.err)
		}
	}
}

func TestExtensionResolver(t *testing.T) {
	ext := `
extensions:
  - name: "math"
  - name: "strings_en_US"
    version: 1`

	baseEnv, err := cel.NewEnv()
	if err != nil {
		t.Fatalf("cel.NewEnv() failed: %v", err)
	}
	c := parseConfigYaml(t, ext)
	for _, e := range c.Extensions {
		e.ExtensionResolver = stringLocaleResolver{}
	}
	opts, err := c.AsEnvOptions(baseEnv)
	if err != nil {
		t.Errorf("AsEnvOptions() generated error: %v", err)
	}
	extEnv, err := baseEnv.Extend(opts...)
	if err != nil {
		t.Fatalf("baseEnv.Extend() failed: %v", err)
	}
	if !extEnv.HasLibrary("cel.lib.ext.strings") || !extEnv.HasLibrary("cel.lib.ext.math") {
		t.Error("extended env did not contain standardized or custom extensions")
	}
}

func parseConfigYaml(t *testing.T, doc string) *Config {
	config := &Config{}
	if err := yaml.Unmarshal([]byte(doc), config); err != nil {
		t.Fatalf("yaml.Unmarshal(%q) failed: %v", doc, err)
	}
	return config
}

type stringLocaleResolver struct{}

func (stringLocaleResolver) ResolveExtension(name string) (ExtensionFactory, bool) {
	parts := strings.SplitN(name, "_", 2)
	if len(parts) == 2 && parts[0] == "strings" {
		return func(version uint32) cel.EnvOption {
			return ext.Strings(
				ext.StringsLocale(parts[1]),
				ext.StringsVersion(version),
			)
		}, true
	}
	return nil, false
}
