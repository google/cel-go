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
	env, err := cel.NewEnv(cel.Lib(customLib{}))
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}
	// Check iss for error in both Parse and Check.
	ast, iss := env.Compile(`i.greet(you)`)
	if iss != nil && iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
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

type customLib struct{}

func (customLib) EnvOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Declarations(
			decls.NewIdent("i", decls.String, nil),
			decls.NewIdent("you", decls.String, nil),
			decls.NewFunction("greet",
				decls.NewInstanceOverload("string_greet_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.String))),
	}
}

func (customLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "string_greet_string",
				Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
					return types.String(
						fmt.Sprintf("Hello %s! Nice to meet you, I'm %s.\n", rhs, lhs))
				},
			},
		),
	}
}
