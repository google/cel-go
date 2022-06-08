// Copyright 2019 Google LLC
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

package cel_test

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func Example() {
	// Create the CEL environment with declarations for the input attributes and the extension functions.
	// In many cases the desired functionality will be present in a built-in function.
	e, err := cel.NewEnv(
		// Variable identifiers used within this expression.
		cel.Variable("i", cel.StringType),
		cel.Variable("you", cel.StringType),
		// Function to generate a greeting from one person to another: i.greet(you)
		cel.Function("greet",
			cel.MemberOverload("string_greet_string", []*cel.Type{cel.StringType, cel.StringType}, cel.StringType,
				cel.BinaryImpl(func(lhs, rhs ref.Val) ref.Val {
					return types.String(fmt.Sprintf("Hello %s! Nice to meet you, I'm %s.\n", rhs, lhs))
				}),
			),
		),
	)
	if err != nil {
		log.Fatalf("environment creation error: %s\n", err)
	}

	// Compile the expression.
	ast, iss := e.Compile("i.greet(you)")
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}

	// Create the program.
	prg, err := e.Program(ast)
	if err != nil {
		log.Fatalf("program creation error: %s\n", err)
	}

	// Evaluate the program against some inputs. Note: the details return is not used.
	out, _, err := prg.Eval(map[string]interface{}{
		// Native values are converted to CEL values under the covers.
		"i": "CEL",
		// Values may also be lazily supplied.
		"you": func() ref.Val { return types.String("world") },
	})
	if err != nil {
		log.Fatalf("runtime error: %s\n", err)
	}

	fmt.Println(out)
	// Output:Hello world! Nice to meet you, I'm CEL.
}

func Example_globalOverload() {
	// The GlobalOverload example demonstrates how to define global overload function.
	// Create the CEL environment with declarations for the input attributes and
	// the desired extension functions. In many cases the desired functionality will
	// be present in a built-in function.
	e, err := cel.NewEnv(
		// Identifiers used within this expression.
		cel.Variable("i", cel.StringType),
		cel.Variable("you", cel.StringType),
		// Function to generate shake_hands between two people.
		//    shake_hands(i,you)
		cel.Function("shake_hands",
			cel.Overload("shake_hands_string_string", []*cel.Type{cel.StringType, cel.StringType}, cel.StringType,
				cel.BinaryImpl(func(arg1, arg2 ref.Val) ref.Val {
					return types.String(fmt.Sprintf("%v and %v are shaking hands.\n", arg1, arg2))
				}),
			),
		),
	)
	if err != nil {
		log.Fatalf("environment creation error: %s\n", err)
	}

	// Compile the expression.
	ast, iss := e.Compile(`shake_hands(i,you)`)
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}

	// Create the program.
	prg, err := e.Program(ast)
	if err != nil {
		log.Fatalf("program creation error: %s\n", err)
	}

	// Evaluate the program against some inputs. Note: the details return is not used.
	out, _, err := prg.Eval(map[string]interface{}{
		"i":   "CEL",
		"you": func() ref.Val { return types.String("world") },
	})
	if err != nil {
		log.Fatalf("runtime error: %s\n", err)
	}

	fmt.Println(out)
	// Output:CEL and world are shaking hands.
}
