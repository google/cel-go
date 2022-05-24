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
// Implements a CLI REPL session for CEL.
//
// introduces commands for manipulating evaluation context to simulate a realistic host
// environment for evaluating expresions.
//
// example session:
//
// ```
// $ go run .
// CEL REPL
// %exit or EOF to quit.

// cel-repl> %let x = 42
// cel-repl> %let y = {'a': x, 'b': y}
// Adding let failed:
// Error updating y = {'a': x, 'b': y}
// ERROR: <input>:1:15: undeclared reference to 'y' (in container '')
//  | {'a': x, 'b': y}
//  | ..............^
// cel-repl> %let z = 41
// cel-repl> %let y = {'a': x, 'b': z}
// cel-repl> y.map(key, y[key]).filter(x, x > 41)
// [42] (list_type:{elem_type:{primitive:INT64}})
// ```
package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func main() {
	var c readline.Config
	c.Prompt = "cel-repl> "

	err := c.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Init readline failed: %v\n", err)
		os.Exit(1)
	}

	rl, err := readline.NewEx(&c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewEx readline failed: %v\n", err)
		os.Exit(1)
	}

	defer rl.Close()

	fmt.Println("CEL REPL")
	fmt.Printf("%%exit or EOF to quit.\n\n")

	eval, err := NewEvaluator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewEvaluator failed: %v\n", err)
		os.Exit(1)
	}

PromptLoop:
	for {
		line, err := rl.Readline()
		if err != nil {
			// Likely eof or interrupt so exit.
			break
		}

		cmd, args, expr, err := Parse(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid command: %v\n", err)
			continue
		}
		switch cmd {
		case "eval":
			val, resultT, err := eval.Evaluate(expr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Expr failed:\n%v\n", err)
			}
			if val != nil {
				fmt.Printf("%v : %s\n", val.Value(), UnparseType(resultT))
			}
		case "let":
			var typeHint *exprpb.Type
			var err error
			if len(args) == 2 {
				typeHint, err = ParseType(args[1])
				if err != nil {
					fmt.Printf("Adding let failed:\n%v\n", err)
					break
				}
			}
			err = eval.AddLetVar(args[0], expr, typeHint)
			if err != nil {
				fmt.Printf("Adding let failed:\n%v\n", err)
			}
		case "declare":
			typePB, err := ParseType(args[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Adding decl failed:\n%b\n", err)
				break
			}
			err = eval.AddDeclVar(args[0], typePB)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Adding decl failed:\n%b\n", err)
			}
		case "delete":
			err = eval.DelLetVar(args[0])
			if err != nil {
				fmt.Printf("Deleting declaration failed:\n%v\n", err)
			}
		case "exit":
			break PromptLoop
		default:
			// nothing to do so just refresh the prompt
			continue
		}

	}
}
