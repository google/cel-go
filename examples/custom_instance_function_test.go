package examples

import (
	"fmt"
	"log"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func ExampleCustomInstanceFunction() {
	d := cel.Declarations(decls.NewIdent("i", decls.String, nil),
		decls.NewIdent("you", decls.String, nil),
		decls.NewFunction("greet",
			decls.NewInstanceOverload("string_greet_string",
				[]*exprpb.Type{decls.String, decls.String},
				decls.String)))
	env, err := cel.NewEnv(d)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}
	// Check iss for error in both Parse and Check.
	p, iss := env.Parse(`i.greet(you)`)
	if iss != nil && iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	c, iss := env.Check(p)
	if iss != nil && iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	greetFunc := &functions.Overload{
		Operator: "string_greet_string",
		Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
			return types.String(
				fmt.Sprintf("Hello %s! Nice to meet you, I'm %s.\n", rhs, lhs))
		}}
	prg, err := env.Program(c, cel.Functions(greetFunc))
	if err != nil {
		log.Fatalf("Program creation error: %v\n", err)
	}

	out, _, err := prg.Eval(map[string]interface{}{
		"i":   "CEL",
		"you": "world",
	})
	if err != nil {
		log.Fatalf("Evaluation error: %v\n", err)
	}

	fmt.Println(out)
	// Output:Hello world! Nice to meet you, I'm CEL.
}
