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

type DocKind int

const (
	DocEnv = iota + 1
	DocFunction
	DocOverload
	DocVariable
	DocMacro
	DocExample
)

type MultilineDescription []string

type Doc struct {
	Kind        DocKind
	Name        string
	Type        string
	Signature   string
	Description MultilineDescription
	Children    []*Doc
}

func FormatDescription(descs ...any) (string, error) {
	return FormatDescriptionSeparator("\n\n", descs...)
}

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

func ParseDescription(doc string) MultilineDescription {
	if len(doc) != 0 {
		return strings.Split(doc, "\n")
	}
	return MultilineDescription{}
}

func ParseDescriptions(doc string) []MultilineDescription {
	examples := []MultilineDescription{}
	if len(doc) != 0 {
		lines := strings.Split(doc, "\n")
		lineStart := 0
		for i, l := range lines {
			l = strings.TrimRightFunc(l, unicode.IsSpace)
			if l == "" {
				examples = append(examples, lines[lineStart:i])
				lineStart = i + 1
			}
		}
		if lineStart < len(lines) {
			examples = append(examples, lines[lineStart:])
		}
	}
	return examples
}

func NewVariableDoc(name, celType, description string) *Doc {
	return &Doc{
		Kind:        DocVariable,
		Name:        name,
		Type:        celType,
		Description: strings.Split(description, "\n"),
	}
}

func NewFunctionDoc(name, description string, overloads ...*Doc) *Doc {
	return &Doc{
		Kind:     DocFunction,
		Name:     name,
		Children: overloads,
	}
}

func NewOverloadDoc(id, sig, description string, examples ...*Doc) *Doc {
	return &Doc{
		Kind:      DocOverload,
		Name:      id,
		Signature: sig,
		Children:  examples,
	}
}

func NewMacroDoc(name, description string, examples ...*Doc) *Doc {
	return &Doc{
		Kind:        DocMacro,
		Name:        name,
		Description: strings.Split(description, "\n"),
		Children:    examples,
	}
}

func NewExampleDoc(ex MultilineDescription) *Doc {
	return &Doc{
		Kind:        DocExample,
		Description: ex,
	}
}

type Documentor interface {
	Documentation() *Doc
}
