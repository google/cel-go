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

package types

var Error Type = &errorType{}

type errorType struct {
}

var _ Type = &errorType{}

func (e *errorType) Kind() TypeKind {
	return KindError
}

func (e *errorType) Equals(t Type) bool {
	_, ok := t.(*errorType)
	return ok
}

func (e *errorType) String() string {
	return "!error!"
}
