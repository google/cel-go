package ext

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

var stringTests = []struct {
	expr string
	err  string
}{
	{
		expr: `'hello'.after('l') == 'lo'`,
	},
	{
		expr: `'hello'.after('l', 3) == 'o'`,
	},
	{
		expr: `'hello'.after('none') == ''`,
	},
	{
		expr: `'hello'.after('l', 30) == 'o'`,
		err:  "index out of range: 30",
	},
	{
		expr: `'tacocat'.before('c') == 'ta'`,
	},
	{
		expr: `'tacocat'.before('c', 3) == 'taco'`,
	},
	{
		expr: `'tacocat'.before('none') == 'tacocat'`,
	},
	{
		expr: `'tacocat'.before('c', 30) == 'taco'`,
		err:  "index out of range: 30",
	},
	{
		expr: `'tacocat'.charAt(3) == 'o'`,
	},
	{
		expr: `'tacocat'.charAt(30) == ''`,
		err:  "index out of range: 30",
	},
	{
		expr: `'tacocat'.indexOf('a') == 1`,
	},
	{
		expr: `'tacocat'.indexOf('a', 3) == 5`,
	},
	{
		expr: `'tacocat'.indexOf('none') == -1`,
	},
	{
		expr: `'tacocat'.indexOf('a', 30) == -1`,
		err:  "index out of range: 30",
	},
	{
		expr: `'MixedCase'.lower().after('mixed') == 'case'`,
	},
	{
		expr: `'MixedCase'.upper().after('MIXED') == 'CASE'`,
	},
	{
		expr: `"12 days 12 hours".replace("{0}", "2") == "12 days 12 hours"`,
	},
	{
		expr: `"{0} days {0} hours".replace("{0}", "2") == "2 days 2 hours"`,
	},
	{
		expr: `"{0} days {0} hours".replace("{0}", "2", 1).replace("{0}", "23") == "2 days 23 hours"`,
	},
	{
		expr: `"hello world".split(" ") == ["hello", "world"]`,
	},
	{
		expr: `"hello world events!".split(" ", 1) == ["hello world events!"]`,
	},
	{
		expr: `"hello world events!".split(" ", 2) == ["hello", "world events!"]`,
	},
}

func TestStrings_Library(t *testing.T) {
	env, err := cel.NewEnv(Strings())
	if err != nil {
		t.Fatal(err)
	}
	for i, tst := range stringTests {
		tc := tst
		t.Run(fmt.Sprintf("[%d]", i), func(tt *testing.T) {
			ast, iss := env.Compile(tc.expr)
			if iss.Err() != nil {
				tt.Fatal(iss.Err())
			}
			exe, err := env.Program(ast)
			if err != nil {
				tt.Fatal(err)
			}
			out, _, err := exe.Eval(cel.NoVars())
			if tc.err != "" {
				if err == nil {
					tt.Fatalf("got value %v, wanted error %s for expr: %s", out.Value(), tc.err, tc.expr)
				}
				if tc.err != err.Error() {
					tt.Errorf("got error %v, wanted error %s for expr: %s", err, tc.err, tc.expr)
				}
			} else if err != nil {
				tt.Fatal(err)
			} else if out.Value() != true {
				tt.Errorf("got %v, wanted true for expr: %s", out.Value(), tc.expr)
			}
		})
	}
}
