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

// Position is line/column position in an input source.
type Location struct {
	// Line is the line number of the position, starting with 1.
	line int

	// Column is the column number of the position, starting with 1.
	column int

	// Source is the optional name of the source input that was parsed. This is typically a file name.
	source Source
}

// NewLocation creates a new Location instance.
func NewLocation(s Source, l int, c int) Location {
	return Location{
		source: s,
		line:   l,
		column: c,
	}
}

func (l Location) Line() int {
	return l.line
}

func (l Location) Column() int {
	return l.column
}

func (l Location) Source() Source {
	return l.source
}

var NoLocation = Location{
	line:   0,
	column: 0,
	source: nil,
}
