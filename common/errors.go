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

import "fmt"

// Errors is the main error collector mechanism.
type Errors struct {
	errors []Error
}

// NewErrors returns a new Errors instance.
func NewErrors() *Errors {
	return &Errors{
		errors: []Error{},
	}
}

// ReportError captures an error report from the caller.
func (e *Errors) ReportError(l Location, format string, args ...interface{}) {
	e.reportErrorInstance(Error{
		Location: l,
		Message:  fmt.Sprintf(format, args...),
	})
}

// GetErrors returns all the errors that are accumulated so far.
func (e *Errors) GetErrors() []Error {
	return e.errors[:]
}

func (e *Errors) reportErrorInstance(err Error) {
	e.errors = append(e.errors, err)
}

func (e *Errors) String() string {
	result := ""
	for i, err := range e.errors {
		if i > 0 {
			result += "\n"
		}
		result += err.ToDisplayString()
	}
	return result
}
