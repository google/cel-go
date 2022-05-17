// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"
)

func argsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, x := range a {
		if x != b[i] {
			return false
		}
	}
	return true
}

func TestParse(t *testing.T) {
	var testCases = []struct {
		commandLine string
		wantCmd     string
		wantArgs    []string
		wantExpr    string
	}{
		{
			commandLine: "%let x = 1",
			wantCmd:     "let",
			wantArgs:    []string{"x"},
			wantExpr:    "1",
		},
		{
			commandLine: "%let x: int = 1",
			wantCmd:     "let",
			wantArgs:    []string{"x", "int"},
			wantExpr:    "1",
		},
		{
			commandLine: `%eval x + 2`,
			wantCmd:     "eval",
			wantArgs:    []string{},
			wantExpr:    "x + 2",
		},
		{
			commandLine: "x + 2",
			wantCmd:     "eval",
			wantArgs:    []string{},
			wantExpr:    "x + 2",
		},
		{
			commandLine: `%exit`,
			wantCmd:     "exit",
			wantArgs:    []string{},
			wantExpr:    "",
		},
		{
			commandLine: "   ",
			wantCmd:     "null",
			wantArgs:    []string{},
			wantExpr:    "",
		},
		{
			commandLine: `%delete x`,
			wantCmd:     "delete",
			wantArgs:    []string{"x"},
			wantExpr:    "",
		},
		{
			commandLine: `%declare x: int`,
			wantCmd:     "declare",
			wantArgs:    []string{"x", "int"},
			wantExpr:    "",
		},
	}
	for _, tc := range testCases {
		cmd, args, expr, err := Parse(tc.commandLine)
		if err != nil {
			t.Errorf("Parse(\"%s\") failed: %s", tc.commandLine, err)
		}
		if tc.wantCmd != cmd || !argsEqual(tc.wantArgs, args) || tc.wantExpr != expr {
			t.Errorf("Parse('%s') got (%s, %v, %s) wanted (%s, %v, %s)", tc.commandLine,
				cmd, args, expr, tc.wantCmd, tc.wantArgs, tc.wantExpr)
		}
	}
}

func TestParseErrors(t *testing.T) {
	var testCases = []struct {
		commandLine string
	}{
		{
			// not an identifier
			commandLine: "%let 123 = 1",
		},
		{
			// no assignment
			commandLine: "%let x: int",
		},
		{
			// type required for declare
			commandLine: "%declare x",
		},
		{
			// not an identifier
			commandLine: "%declare 123",
		},
		{
			// not an identifier
			commandLine: "%delete 123",
		},
	}
	for _, tc := range testCases {
		_, _, _, err := Parse(tc.commandLine)
		if err == nil {
			t.Errorf("Parse(\"%s\") ok wanted error", tc.commandLine)
		}
	}
}
