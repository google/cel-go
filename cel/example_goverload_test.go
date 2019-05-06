package cel_test

import (
	"fmt"
	"log"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"
)

// ExampleGlobalOverload demonstrates how to define global overload function.
func Example_globalOverload() {
	// Create the CEL environment with declarations for the input attributes and
	// the desired extension functions. In many cases the desired functionality will
	// be present in a built-in function.
	decls := cel.Declarations(
		// Identifiers used within this expression.
		decls.NewIdent("i", decls.String, nil),
		decls.NewIdent("you", decls.String, nil),
		// Function to generate shake_hands between two people.
		//    shake_hands(i,you)
		decls.NewFunction("shake_hands",
			decls.NewOverload("greet_string_string",
				[]*exprpb.Type{decls.String, decls.String},
				decls.String)))
	e, err := cel.NewEnv(decls)
	if err != nil {
		log.Fatalf("environment creation error: %s\n", err)
	}

	// Parse and check the expression.
	p, iss := e.Parse(`shake_hands(i,you)`)
	if iss != nil && iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	c, iss := e.Check(p)
	if iss != nil && iss.Err() != nil {
		log.Fatalln(iss.Err())
	}

	// Create the program.
	funcs := cel.Functions(
		&functions.Overload{
			Operator: "shake_hands",
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				s1, ok := lhs.(types.String)
				if !ok {
					return types.NewErr("unexpected type '%v' passed to shake_hands", lhs.Type())
				}
				// If expression was type checked and dynamic data is not used (JSON, maps),
				// type checking can be skipped.
				s2 := rhs.(types.String)
				return types.String(
					fmt.Sprintf("%s and %s are shaking hands.\n", s1, s2))
			}})
	prg, err := e.Program(c, funcs)
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
