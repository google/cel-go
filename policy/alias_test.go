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

package policy_test

import (
	"testing"

	"github.com/google/cel-go/policy"
	canonicalpolicy "cel.dev/cel-go/policy"
)

func TestTypeAliases(t *testing.T) {
	var p *policy.Policy = nil
	var cp *canonicalpolicy.Policy = p
	_ = cp

	var r *policy.Rule = nil
	var cr *canonicalpolicy.Rule = r
	_ = cr
}
