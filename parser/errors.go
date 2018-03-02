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

package parser

import (
	"fmt"

	"celgo/common"
)

// ParseErrors is a specialization of Errors.
type ParseErrors struct {
	*common.Errors
}

func (e *ParseErrors) syntaxError(l common.Location, message string) {
	e.ReportError(l, fmt.Sprintf("Syntax error: %s", message))
}

func (e *ParseErrors) invalidHasArgument(l common.Location) {
	e.ReportError(l, "The argument to the function 'has' must be a field selection")
}

func (e *ParseErrors) argumentIsNotIdent(l common.Location) {
	e.ReportError(l, "The argument must be a simple name")
}

func (e *ParseErrors) notAQualifiedName(l common.Location) {
	e.ReportError(l, "expected a qualified name")
}
