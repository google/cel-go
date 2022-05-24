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
//
// Handles parsing command line inputs into canonical commands.
package main

import (
	"fmt"
	"regexp"
	"strings"
)

var exprSubRe = `(?P<expr>.*)`
var exprAssignRe = `(\s+=\s+)` + exprSubRe
var identifier = `[_a-zA-Z][_a-zA-Z0-9]*`
var typeSuffix = `(\s*:\s*(?P<type_hint>[^=]+))`
var ws = regexp.MustCompile(`[ \t\r\n]+`)
var identPrefix = `(?P<ident>` + identifier + `)`
var letRe = regexp.MustCompile(identPrefix + typeSuffix + `?` + exprAssignRe)
var declRe = regexp.MustCompile(identPrefix + typeSuffix)
var delRe = regexp.MustCompile(identPrefix)

var letUsage = `Let introduces a variable whose value is defined by a sub-CEL expression.
%let <identifier> (: <type>)? = <expr>`

var declareUsage = `Declare introduces a variable for type checking, but doesn't define a value for it.
%declare <identifier> : <type>`

var deleteUsage = `Delete removes a variable declaration from the evaluation context.
%delete <identifier>`

func subexpByName(r *regexp.Regexp, m []string, n string) *string {
	idx := r.SubexpIndex(n)
	if idx < 0 {
		return nil
	}
	return &m[idx]
}

// PromptError represents an invalid command (distinct from invalid CEL expr)
type PromptError struct {
	msg   string
	usage *string
}

func (p PromptError) Error() string {
	r := p.msg
	if p.usage != nil {
		r = fmt.Sprintf("%s\nUsage: \n%s", r, *p.usage)
	}
	return r
}

// Parse a cli command into a canonical command and arguments.
func Parse(line string) (cmd string, args []string, expr string, err error) {
	// TODO(issue/538): Switch to a parsing library over regex as this gets more complicated.
	line = strings.TrimSpace(line)

	if strings.IndexRune(line, '%') == 0 {
		span := ws.FindStringIndex(line)
		if span != nil {
			cmd = line[1:span[0]]
			line = line[span[1]:]
		} else {
			cmd = line[1:]
			line = ""
		}
	} else if line == "" {
		cmd = "null"
		return
	} else {
		// default is to just evaluate cel
		cmd = "eval"
	}

	switch cmd {
	case "let":
		m := letRe.FindStringSubmatch(line)
		if m == nil {
			err = PromptError{msg: "invalid let statement", usage: &letUsage}
			return
		}
		if a := subexpByName(letRe, m, "ident"); a != nil {
			args = append(args, *a)
		}
		if t := subexpByName(letRe, m, "type_hint"); t != nil && *t != "" {
			args = append(args, *t)
		}
		if e := subexpByName(letRe, m, "expr"); e != nil {
			expr = *e
		}
	case "delete":
		m := delRe.FindStringSubmatch(line)
		if m == nil {
			err = PromptError{msg: "invalid delete statement", usage: &deleteUsage}
			return
		}
		if i := subexpByName(delRe, m, "ident"); i != nil {
			args = append(args, *i)
		}
	case "declare":
		m := declRe.FindStringSubmatch(line)
		if m == nil {
			err = PromptError{msg: "invalid declare statement", usage: &declareUsage}
			return
		}
		if a := subexpByName(declRe, m, "ident"); a != nil {
			args = append(args, *a)
		}
		if t := subexpByName(declRe, m, "type_hint"); t != nil {
			args = append(args, *t)
		}
	case "eval":
		expr = line
	case "exit":
	case "null":
	default:
		err = PromptError{msg: "Undefined command: " + cmd}
	}
	return
}
