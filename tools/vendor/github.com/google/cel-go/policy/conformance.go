// Copyright 2024 Google LLC
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

// TestSuite describes a set of tests divided by section.
type TestSuite struct {
	Description string         `yaml:"description"`
	Sections    []*TestSection `yaml:"section"`
}

// TestSection describes a related set of tests associated with a behavior.
type TestSection struct {
	Name  string      `yaml:"name"`
	Tests []*TestCase `yaml:"tests"`
}

// TestCase describes a named test scenario with a set of inputs and expected outputs.
//
// Note, when a test requires additional functions to be provided to execute, the test harness
// must supply these functions.
type TestCase struct {
	Name   string               `yaml:"name"`
	Input  map[string]TestInput `yaml:"input"`
	Output string               `yaml:"output"`
}

// TestInput represents an input literal value or expression.
type TestInput struct {
	// Value is a simple literal value.
	Value any `yaml:"value"`

	// Expr is a CEL expression based input.
	Expr string `yaml:"expr"`

	// ContextExpr is a CEL expression which is used as cel.ContextProtoVars
	ContextExpr string `yaml:"context_expr"`
}
