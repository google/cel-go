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
	"fmt"
	"strings"
	"testing"

	"google.golang.org/protobuf/proto"
)

func cmdMatches(t testing.TB, got Cmd, expected Cmd) (result bool) {
	t.Helper()

	defer func() {
		// catch type assertion errors
		if v := recover(); v != nil {
			result = false
		}
	}()

	switch expected := expected.(type) {
	case *evalCmd:
		return got.(*evalCmd).expr == expected.expr
	case *delCmd:
		return got.(*delCmd).identifier == expected.identifier
	case *simpleCmd:
		return got.(*simpleCmd).cmd == expected.cmd
	case *letFnCmd:
		got := got.(*letFnCmd)
		if got.identifier != expected.identifier ||
			got.src != expected.src ||
			!proto.Equal(got.resultType, expected.resultType) ||
			len(got.params) != len(expected.params) {
			return false
		}
		for i, p := range expected.params {
			if got.params[i].identifier != p.identifier ||
				!proto.Equal(got.params[i].typeHint, p.typeHint) {
				return false
			}
		}
		return true
	case *letVarCmd:
		got := got.(*letVarCmd)
		return got.identifier == expected.identifier &&
			proto.Equal(got.typeHint, expected.typeHint) &&
			got.src == expected.src
	}

	return false
}

func (c *evalCmd) String() string {
	return fmt.Sprintf("%%eval %s", c.expr)
}

func (c *letVarCmd) String() string {
	return fmt.Sprintf("%%let %s : %s = %s", c.identifier, UnparseType(c.typeHint), c.src)
}

func fmtParam(p letFunctionParam) string {
	return fmt.Sprintf("%s : %s", p.identifier, UnparseType(p.typeHint))
}

func fmtParams(ps []letFunctionParam) string {
	buf := make([]string, len(ps))
	for i, p := range ps {
		buf[i] = fmtParam(p)
	}
	return strings.Join(buf, ", ")
}

func (c *letFnCmd) String() string {
	return fmt.Sprintf("%%let %s (%s) : %s -> %s", c.identifier, fmtParams(c.params), UnparseType(c.resultType), c.src)
}

func (c *delCmd) String() string {
	return fmt.Sprintf("%%delete %s", c.identifier)
}

func (c *simpleCmd) String() string {
	return fmt.Sprintf("%%%s", c.cmd)
}

func TestParse(t *testing.T) {
	var testCases = []struct {
		commandLine string
		wantCmd     Cmd
	}{
		{
			commandLine: "%let x = 1",
			wantCmd: &letVarCmd{
				identifier: "x",
				typeHint:   nil,
				src:        "1",
			},
		},
		{
			commandLine: "%let x: int = 1",
			wantCmd: &letVarCmd{
				identifier: "x",
				typeHint:   mustParseType(t, "int"),
				src:        "1",
			},
		},
		{
			commandLine: `%eval x + 2`,
			wantCmd:     &evalCmd{expr: "x + 2"},
		},
		{
			commandLine: "x + 2",
			wantCmd:     &evalCmd{expr: "x + 2"},
		},
		{
			commandLine: `%exit`,
			wantCmd:     &simpleCmd{cmd: "exit"},
		},
		{
			commandLine: "   ",
			wantCmd:     &simpleCmd{"null"},
		},
		{
			commandLine: `%delete x`,
			wantCmd:     &delCmd{identifier: "x"},
		},
		{
			commandLine: `%declare x: int`,
			wantCmd: &letVarCmd{
				identifier: "x",
				typeHint:   mustParseType(t, "int"),
				src:        "",
			},
		},
		{
			commandLine: `%let fn (x : int) : int -> x + 2`,
			wantCmd: &letFnCmd{
				identifier: "fn",
				params: []letFunctionParam{
					{identifier: "x", typeHint: mustParseType(t, "int")},
				},
				resultType: mustParseType(t, "int"),
				src:        "x + 2",
			},
		},
		{
			commandLine: `%let fn () : int -> 2 + 2`,
			wantCmd: &letFnCmd{
				identifier: "fn",
				params:     []letFunctionParam{},
				resultType: mustParseType(t, "int"),
				src:        "2 + 2",
			},
		},
		{
			commandLine: `%let fn (x:int, y :int) : int -> x + y`,
			wantCmd: &letFnCmd{
				identifier: "fn",
				params: []letFunctionParam{
					{identifier: "x", typeHint: mustParseType(t, "int")},
					{identifier: "y", typeHint: mustParseType(t, "int")},
				},
				resultType: mustParseType(t, "int"),
				src:        "x + y",
			},
		},
		{
			commandLine: `%declare fn (x : int) : int`,
			wantCmd: &letFnCmd{
				identifier: "fn",
				params: []letFunctionParam{
					{identifier: "x", typeHint: mustParseType(t, "int")},
				},
				resultType: mustParseType(t, "int"),
				src:        "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.commandLine, func(t *testing.T) {
			cmd, err := Parse(tc.commandLine)
			if err != nil {
				t.Errorf("Parse(\"%s\") failed: %s", tc.commandLine, err)
			}
			if cmd == nil || !cmdMatches(t, cmd, tc.wantCmd) {
				t.Errorf("Parse('%s') got (%s) wanted (%s)", tc.commandLine,
					cmd, tc.wantCmd)
			}
		})
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
			// no assignment
			commandLine: "%let fn() : int",
		},
		{
			// missing types
			commandLine: "%let fn (x) -> x + 1",
		},
		{
			// missing arg id
			commandLine: "%let fn (: int) : int -> x + 1",
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
			commandLine: "%declare 1() : int",
		},
		{
			// not an identifier
			commandLine: "%delete 123",
		},
	}
	for _, tc := range testCases {
		_, err := Parse(tc.commandLine)
		if err == nil {
			t.Errorf("Parse(\"%s\") ok wanted error", tc.commandLine)
		}
	}
}
