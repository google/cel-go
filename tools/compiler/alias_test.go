// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiler_test

import (
	"testing"

	"github.com/google/cel-go/tools/compiler"
)

func TestCompilerAlias(t *testing.T) {
	c, err := compiler.NewCompiler()
	if err != nil {
		t.Fatalf("compiler.NewCompiler() failed: %v", err)
	}
	env, err := c.CreateEnv()
	if err != nil {
		t.Fatalf("c.CreateEnv() failed: %v", err)
	}
	if env == nil {
		t.Fatal("expected non-nil env")
	}

	format := compiler.InferFileFormat("test.yaml")
	if format != compiler.TextYAML {
		t.Fatalf("expected TextYAML, got %v", format)
	}
}
