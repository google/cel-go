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

import "regexp"

var lineRegexp = regexp.MustCompile("(?m)^")

// Source represents a source location.
type Source interface {
	Name() string
	Snippet(line int) (string, bool)
}

// TestSource is a source from an input text.
type TextSource struct {
	name     string
	contents string
}

var _ Source = &TextSource{}

// NewTextSource returns new TextSource instance.
func NewTextSource(name string, contents string) Source {
	return &TextSource{
		name:     name,
		contents: contents,
	}
}

func (s *TextSource) Name() string {
	return s.name
}

func (s *TextSource) Snippet(line int) (string, bool) {
	if s.contents == "" {
		return "", false
	}

	start := -1
	end := -1
	for i, m := range lineRegexp.FindAllStringIndex(s.contents, -1) {
		if i+1 == line {
			start = m[0]
			continue
		}
		if i == line {
			end = m[0]
			break
		}
	}

	if start == -1 {
		// Source line didn't match.
		return "", false
	}

	if end == -1 {
		end = len(s.contents)
	}

	return s.contents[start:end], true
}
