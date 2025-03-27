// Copyright 2025 Google LLC
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

package cel

import (
	"testing"

	"github.com/google/cel-go/common/env"
)

func TestPromptTemplate(t *testing.T) {
	tests := []struct {
		name    string
		envOpts []EnvOption
		out     string
	}{
		{
			name: "basic",
			out: `You are a software engineer with expertise in networking and application security
authoring boolean Common Expression Language (CEL) expressions to ensure firewall,
networking, authentication, and data access is only permitted when all conditions
are satisified.

Output your response as a CEL expression.

Write the expression with the comment on the first line and the expression on the
subsequent lines. Format the expression using 80-character line limits commonly
found in C++ or Java code.
`,
		},
		{
			name:    "macros",
			envOpts: []EnvOption{Macros(StandardMacros...)},
			out: `You are a software engineer with expertise in networking and application security
authoring boolean Common Expression Language (CEL) expressions to ensure firewall,
networking, authentication, and data access is only permitted when all conditions
are satisified.

Output your response as a CEL expression.

Write the expression with the comment on the first line and the expression on the
subsequent lines. Format the expression using 80-character line limits commonly
found in C++ or Java cod
`,
		},
		{
			name:    "standard_env",
			envOpts: []EnvOption{StdLib(StdLibSubset(env.NewLibrarySubset().SetDisableMacros(true)))},
			out: `You are a software engineer with expertise in networking and application security
authoring boolean Common Expression Language (CEL) expressions to ensure firewall,
networking, authentication, and data access is only permitted when all conditions
are satisified.

Output your response as a CEL expression.

Write the expression with the comment on the first line and the expression on the
subsequent lines. Format the expression using 80-character line limits commonly
found in C++ or Java cod
`,
		},
	}

	for _, tst := range tests {
		tc := tst
		t.Run(tc.name, func(t *testing.T) {
			env, err := NewCustomEnv(tc.envOpts...)
			if err != nil {
				t.Fatalf("cel.NewCustomEnv() failed: %v", err)
			}
			prompt, err := AuthoringPrompt(env)
			if err != nil {
				t.Fatalf("cel.AuthoringPrompt() failed: %v", err)
			}
			out := prompt.Render("<USER_PROMPT>")
			if out != tc.out {
				t.Errorf("got %s, wanted %s", out, tc.out)
			}
		})
	}
}
