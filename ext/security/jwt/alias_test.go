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

package jwt_test

import (
	"testing"

	"github.com/google/cel-go/ext/security/jwt"
	canonicalcel "cel.dev/cel-go/cel"
)

func TestJwtAlias(t *testing.T) {
	env, err := canonicalcel.NewEnv(jwt.Library())
	if err != nil {
		t.Fatalf("NewEnv with jwt.Library failed: %v", err)
	}
	if env == nil {
		t.Fatal("expected non-nil env")
	}
}
