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

package common

import (
	"fmt"
	"strings"
)

// Error represents an error message.
type Error struct {
	// Location within the source.
	Location Location

	// Message is the actual error message.
	Message string
}

// ToDisplayString returns the error in a user-friendly format.
func (e *Error) ToDisplayString() string {

	result := fmt.Sprintf("ERROR: %s:%d:%d: %s", e.Location.Source().Name(), e.Location.Line(), e.Location.Column(), e.Message)

	if snippet, found := e.Location.Source().Snippet(e.Location.Line()); found {
		result += "\n | "
		result += snippet
		result += "\n | "
		result += strings.Repeat(".", int(e.Location.Column()-1))
		result += "^"

	}
	return result
}
