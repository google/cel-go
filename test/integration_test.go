package test

import (
	"testing"

	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/parser"
)

func Test_stringOverfloats(t *testing.T) {
	for _, expression := range []string{
		"req_params.Nick.matches('k.*')",
		"nick.matches('k.*')",
		"nick in ['klmno', 'kpacha']",
	} {
		source := common.NewTextSource(expression)
		p, errors := parser.Parse(source)
		if len(errors.GetErrors()) != 0 {
			t.Error(errors.ToDisplayString())
			continue
		}

		typeProvider := types.NewProvider()
		env := checker.NewStandardEnv(packages.DefaultPackage, typeProvider)
		env.Add(decls.NewIdent("req_params", decls.NewMapType(decls.String, decls.String), nil))
		env.Add(decls.NewIdent("nick", decls.String, nil))

		checkedExp, errors := checker.Check(p, source, env)
		if len(errors.GetErrors()) != 0 {
			t.Error(errors.ToDisplayString())
			continue
		}

		interpreterInstance := interpreter.NewStandardInterpreter(packages.DefaultPackage, typeProvider)
		eval := interpreterInstance.NewInterpretable(interpreter.NewCheckedProgram(checkedExp))

		for _, tc := range []struct {
			nick string
			res  bool
		}{
			{nick: "kpacha", res: true},
			{nick: "foo", res: false},
			{nick: "bar", res: false},
			{nick: "klmno", res: true},
		} {
			result, _ := eval.Eval(interpreter.NewActivation(map[string]interface{}{
				"req_params": map[string]string{"Nick": tc.nick},
				"nick":       tc.nick,
			}))
			if v, ok := result.Value().(bool); !ok || v != tc.res {
				t.Errorf("%s [%s] -> unexpected result: %+v", expression, tc.nick, result.Value())
			}
		}
	}
}
