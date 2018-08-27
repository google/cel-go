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
	"strings"

	exprpb "github.com/google/cel-spec/proto/v1/syntax"
)

// Source interface for filter source contents.
type Source interface {
	// The source content represented as a string, for example a single file,
	// textbox field, or url parameter.
	Content() string

	// Brief description of the source, such as a file name or ui element.
	Description() string

	// The character offsets at which lines occur. The zero-th entry should
	// refer to the break between the first and second line, or EOF if there
	// is only one line of source.
	LineOffsets() []int32

	// The raw character offset at which the a location exists given the
	// location line and column.
	// Returns the line offset and whether the location was found.
	LocationOffset(location Location) (int32, bool)

	// OffsetLocation translates a character offset to a Location, or
	// false if the conversion was not feasible.
	OffsetLocation(offset int32) (Location, bool)

	// Return a line of content and whether the line was found.
	Snippet(line int) (string, bool)

	// IdOffset returns the raw character offset of an expression within
	// the source, or false if the expression cannot be found.
	IdOffset(exprId int64) (int32, bool)

	// IdLocation returns a Location for the given expression id,
	// or false if one cannot be found.  It behaves as the obvious
	// composition of IdOffset() and OffsetLocation().
	IdLocation(exprId int64) (Location, bool)
}

// The sourceImpl type implementation of the Source interface.
type sourceImpl struct {
	contents    string
	description string
	lineOffsets []int32
	idOffsets   map[int64]int32
}

// TODO(jimlarson) "Character offsets" should index the code points
// within the UTF-8 encoded string.  It currently indexes bytes.
// Can be accomplished by using rune[] instead of string for contents.

// Create a new Source given the string contents and description.
func NewStringSource(contents string, description string) Source {
	// Compute line offsets up front as they are referred to frequently.
	lines := strings.Split(contents, "\n")
	offsets := make([]int32, len(lines))
	var offset int32 = 0
	for i, line := range lines {
		offset = offset + int32(len(line)) + 1
		offsets[int32(i)] = offset
	}
	return &sourceImpl{
		contents:    contents,
		description: description,
		lineOffsets: offsets,
		idOffsets:   map[int64]int32{},
	}
}

func NewInfoSource(info *exprpb.SourceInfo) Source {
	return &sourceImpl{
		contents:    "",
		description: info.Location,
		lineOffsets: info.LineOffsets,
		idOffsets:   info.Positions,
	}
}

func (s *sourceImpl) Content() string {
	return s.contents
}

func (s *sourceImpl) Description() string {
	return s.description
}

func (s *sourceImpl) LineOffsets() []int32 {
	return s.lineOffsets
}

func (s *sourceImpl) LocationOffset(location Location) (int32, bool) {
	if lineOffset, found := s.findLineOffset(location.Line()); found {
		return lineOffset + int32(location.Column()), true
	}
	return -1, false
}

func (s *sourceImpl) OffsetLocation(offset int32) (Location, bool) {
	line, lineOffset := s.findLine(offset)
	return NewLocation(int(line), int(offset-lineOffset)), true
}

func (s *sourceImpl) Snippet(line int) (string, bool) {
	charStart, found := s.findLineOffset(line)
	if !found || len(s.contents) == 0 {
		return "", false
	}
	charEnd, found := s.findLineOffset(line + 1)
	if found {
		return s.contents[charStart : charEnd-1], true
	}
	return s.contents[charStart:], true
}

func (s *sourceImpl) IdOffset(exprId int64) (int32, bool) {
	if offset, found := s.idOffsets[exprId]; found {
		return offset, true
	}
	return -1, false
}

func (s *sourceImpl) IdLocation(exprId int64) (Location, bool) {
	if offset, found := s.IdOffset(exprId); found {
		if location, found := s.OffsetLocation(offset); found {
			return location, true
		}
	}
	return NewLocation(1, 0), false
}

// findLineOffset returns the offset where the (1-indexed) line begins,
// or false if line doesn't exist.
func (s *sourceImpl) findLineOffset(line int) (int32, bool) {
	if line == 1 {
		return 0, true
	} else if line > 1 && line <= int(len(s.lineOffsets)) {
		offset := s.lineOffsets[line-2]
		return offset, true
	}
	return -1, false
}

// findLine finds the line that contains the given character offset and
// returns the line number and offset of the beginning of that line.
// Note that the last line is treated as if it contains all offsets
// beyond the end of the actual source.
func (s *sourceImpl) findLine(characterOffset int32) (int32, int32) {
	var line int32 = 1
	for _, lineOffset := range s.lineOffsets {
		if lineOffset > characterOffset {
			break
		} else {
			line += 1
		}
	}
	if line == 1 {
		return line, 0
	}
	return line, s.lineOffsets[line-2]
}
