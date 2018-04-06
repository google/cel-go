// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// package common defines types common to parsing and other diagnostics.
package common

import (
	"testing"
)

// TestErrors reporting and recording.
func TestErrors(t *testing.T) {
	source := NewStringSource("a.b\n&&arg(missing, paren", "errors-test")
	errors := NewErrors(source)
	errors.ReportError(NewLocation(1, 1), "No such field")
	if len(errors.GetErrors()) != 1 {
		t.Errorf("%s first error not recorded", t.Name())
	}
	errors.ReportError(NewLocation(2, 20),
		"Syntax error, missing paren")
	if len(errors.GetErrors()) != 2 {
		t.Errorf("%s second error not recorded", t.Name())
	}
	expected :=
		"ERROR: errors-test:1:2: No such field\n" +
			" | a.b\n" +
			" | .^\n" +
			"ERROR: errors-test:2:21: Syntax error, missing paren\n" +
			" | &&arg(missing, paren\n" +
			" | ....................^"
	actual := errors.ToDisplayString()
	if actual != expected {
		t.Errorf("%s got %s, wanted %s", t.Name(), actual, expected)
	}
}
