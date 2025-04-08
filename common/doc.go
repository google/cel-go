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

// Package common defines types and utilities common to expression parsing,
// checking, and interpretation
package common

import (
	"fmt"
	"strings"
	"unicode"
)

// DocKind indicates the type of documentation element.
type DocKind int

const (
	// DocEnv represents environment variable documentation.
	DocEnv DocKind = iota + 1
	// DocFunction represents function documentation.
	DocFunction
	// DocOverload represents function overload documentation.
	DocOverload
	// DocVariable represents variable documentation.
	DocVariable
	// DocMacro represents macro documentation.
	DocMacro
	// DocExample represents example documentation.
	DocExample
)

// MultilineDescription represents a description that can span multiple lines,
// stored as a slice of strings.
type MultilineDescription []string

// Doc holds the documentation details for a specific program element like
// a variable, function, macro, or example.
type Doc struct {
	// Kind specifies the type of documentation element (e.g., Function, Variable).
	Kind DocKind

	// Name is the identifier of the documented element (e.g., function name, variable name).
	Name string

	// Type is the data type associated with the element, primarily used for variables.
	Type string

	// Signature represents the function or overload signature.
	Signature string

	// Description holds the textual description of the element, potentially spanning multiple lines.
	Description MultilineDescription

	// Children holds nested documentation elements, such as overloads for a function
	// or examples for a function/macro.
	Children []*Doc
}

// FormatDescription joins multiple description elements (string, MultilineDescription,
// or []MultilineDescription) into a single string, separated by double newlines ("\n\n").
// It returns the formatted string or an error if an unsupported type is encountered.
func FormatDescription(descs ...any) (string, error) {
	return FormatDescriptionSeparator("\n\n", descs...)
}

// FormatDescriptionSeparator joins multiple description elements (string, MultilineDescription,
// or []MultilineDescription) into a single string using the specified separator.
// It returns the formatted string or an error if an unsupported description type is passed.
func FormatDescriptionSeparator(sep string, descs ...any) (string, error) {
	var builder strings.Builder
	hasDoc := false
	for _, d := range descs {
		if hasDoc {
			builder.WriteString(sep)
		}
		switch v := d.(type) {
		case string:
			builder.WriteString(v)
		case MultilineDescription:
			str := strings.Join(v, "\n")
			builder.WriteString(str)
		case []MultilineDescription:
			for _, md := range v {
				if hasDoc {
					builder.WriteString(sep)
				}
				str := strings.Join(md, "\n")
				builder.WriteString(str)
				hasDoc = true
			}
		default:
			return "", fmt.Errorf("unsupported description type: %T", d)
		}
		hasDoc = true
	}
	return builder.String(), nil
}

// ParseDescription takes a single string containing newline characters and splits
// it into a MultilineDescription. All empty lines will be skipped.
//
// Returns an empty MultilineDescription if the input string is empty.
func ParseDescription(doc string) MultilineDescription {
	var lines MultilineDescription
	if len(doc) != 0 {
		// Split the input string by newline characters.
		for _, line := range strings.Split(doc, "\n") {
			l := strings.TrimRightFunc(line, unicode.IsSpace)
			if len(l) == 0 {
				continue
			}
			lines = append(lines, l)
		}
	}
	// Return an empty slice if the input is empty.
	return lines
}

// ParseDescriptions splits a documentation string into multiple MultilineDescription
// sections, using blank lines as delimiters.
func ParseDescriptions(doc string) []MultilineDescription {
	var examples []MultilineDescription
	if len(doc) != 0 {
		lines := strings.Split(doc, "\n")
		lineStart := 0
		for i, l := range lines {
			// Trim trailing whitespace to identify effectively blank lines.
			l = strings.TrimRightFunc(l, unicode.IsSpace)
			// If a line is blank, it marks the end of the current section.
			if len(l) == 0 {
				// Start the next section after the blank line.
				ex := lines[lineStart:i]
				if len(ex) != 0 {
					examples = append(examples, ex)
				}
				lineStart = i + 1
			}
		}
		// Append the last section if it wasn't terminated by a blank line.
		if lineStart < len(lines) {
			examples = append(examples, lines[lineStart:])
		}
	}
	return examples
}

// NewVariableDoc creates a new Doc struct specifically for documenting a variable.
func NewVariableDoc(name, celType, description string) *Doc {
	return &Doc{
		Kind:        DocVariable,
		Name:        name,
		Type:        celType,
		Description: ParseDescription(description),
	}
}

// NewFunctionDoc creates a new Doc struct for documenting a function.
func NewFunctionDoc(name, description string, overloads ...*Doc) *Doc {
	return &Doc{
		Kind:        DocFunction,
		Name:        name,
		Description: ParseDescription(description),
		Children:    overloads,
	}
}

// NewOverloadDoc creates a new Doc struct for a function example.
func NewOverloadDoc(id, signature string, examples ...*Doc) *Doc {
	return &Doc{
		Kind:      DocOverload,
		Name:      id,
		Signature: signature,
		Children:  examples,
	}
}

// NewMacroDoc creates a new Doc struct for documenting a macro.
func NewMacroDoc(name, description string, examples ...*Doc) *Doc {
	return &Doc{
		Kind:        DocMacro,
		Name:        name,
		Description: ParseDescription(description),
		Children:    examples,
	}
}

// NewExampleDoc creates a new Doc struct specifically for holding an example.
func NewExampleDoc(ex MultilineDescription) *Doc {
	return &Doc{
		Kind:        DocExample,
		Description: ex,
	}
}

// Documentor is an interface for types that can provide their own documentation.
type Documentor interface {
	// Documentation returns the documentation coded by the DocKind to assist
	// with text formatting.
	Documentation() *Doc
}
