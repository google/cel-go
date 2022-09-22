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

package parser

import (
	"fmt"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type errListener struct {
	antlr.DefaultErrorListener

	errs []error
}

func (l *errListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
	l.errs = append(l.errs, fmt.Errorf("(%d:%d) %s", line, column, msg))
}

func tryParse(t testing.TB, s string) error {
	t.Helper()

	l := &errListener{}
	is := antlr.NewInputStream(s)
	lexer := NewCommandsLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(l)
	p := NewCommandsParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	p.RemoveErrorListeners()
	p.AddErrorListener(l)
	p.StartCommand()

	if len(l.errs) == 0 {
		return nil
	}

	err := l.errs[0].Error()
	for _, el := range l.errs[1:] {
		err = err + "\n" + el.Error()
	}
	return fmt.Errorf("parse errors: %s", err)
}

func TestAccept(t *testing.T) {
	var testCases = []string{
		"%exit",
		"%let id = 2",
		"%arbitrary",
		"%arbitrary --flag -alt_flag 'string arg'",
		" ",
		"%let y : int = [1, 2, 3]",
		"%let fn (y : int) : int -> y + 10",
		"%let fn () : int -> 10",
		"%let fn (x:int, y : int) : int -> x + y",
		"%let fn (x:int, y : int) : int -> x + y",
		"%let com.google.fn (x:int, y : int) : int -> x + y",
		"%let int.plus (x: int) : int -> this + x",
		"%delete id",
		"%delete com.google.id",
		"%delete int.fn(x:int): int",
		"%declare x : int",
		"%declare fn (x : int) : int",
		"%declare x", // accepted by grammar, but business logic will error
		"x + 2",      // also an expr
	}
	for _, tc := range testCases {
		err := tryParse(t, tc)
		if err != nil {
			t.Errorf("parse %s:\ngot %v\nexpected nil", tc, err)
		}

	}
}

func TestReject(t *testing.T) {
	var testCases = []string{
		"%declare 1",
		"%let 1 = 2",
		"%1badid",
		"%let fn x : int : int -> x + 2", // parens required
		"%let fn (x : int) -> x + 2",     // return type required
		"x{{",                            // won't parse as CEL expr
		"%declare fn (x: int)",           // return type required
		"%declare fn (x) : int",          // arg type required
	}
	for _, tc := range testCases {
		err := tryParse(t, tc)
		if err == nil {
			t.Errorf("parse %s got nil, expected error", tc)
		}

	}
}
