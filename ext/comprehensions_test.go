package ext

import (
	"fmt"
	"testing"

	"github.com/google/cel-go/cel"
)

func TestTwoVarComprehensions(t *testing.T) {
	compreTests := []struct {
		expr string
	}{
		// list.all()
		{expr: "[1, 2, 3, 4].all(i, v, i < 5 && v > 0)"},
		{expr: "[1, 2, 3, 4].all(i, v, i < v)"},
		{expr: "[1, 2, 3, 4].all(i, v, i > v) == false"},
		{expr: `
		cel.bind(listA, [1, 2, 3, 4],
		cel.bind(listB, [1, 2, 3, 4, 5],
		   listA.all(i, v, listB[?i].hasValue() && listB[i] == v)
		))
		`},
		{expr: `
		cel.bind(listA, [1, 2, 3, 4, 5, 6],
		cel.bind(listB, [1, 2, 3, 4, 5],
		   listA.all(i, v, listB[?i].hasValue() && listB[i] == v)
		)) == false
		`},
		// list.exists()
		{expr: `
		cel.bind(l, ['hello', 'world', 'hello!', 'worlds'],
		  l.exists(i, v,
		    v.startsWith('hello') && l[?(i+1)].optMap(next, next.endsWith('world')).orValue(false)
		  )
		)
		`},
		// list.existsOne()
		{expr: `
		cel.bind(l, ['hello', 'world', 'hello!', 'worlds'],
		  l.existsOne(i, v,
		    v.startsWith('hello') && l[?(i+1)].optMap(next, next.endsWith('world')).orValue(false)
		  )
		)
		`},
		{expr: `
		cel.bind(l, ['hello', 'goodbye', 'hello!', 'goodbye'],
		  l.exists_one(i, v,
		    v.startsWith('hello') && l[?(i+1)].optMap(next, next == "goodbye").orValue(false)
		  )
		) == false
		`},
		// list.transformList()
		{expr: `
		cel.bind(l, ['Hello', 'world'],
		  l.transformList(i, v, "[%d]%s".format([i, v.lowerAscii()]))
		) == ["[0]hello", "[1]world"]
		`},
		{expr: `
		cel.bind(l, ['hello', 'world'],
		  l.transformList(i, v, v.startsWith('greeting'), "[%d]%s".format([i, v]))
		) == []
		`},
		// map.all()
		{expr: `
		cel.bind(m, {'hello': 'world', 'hello!': 'world'},
		  m.all(k, v, k.startsWith('hello') && v == 'world')
		)
		`},
		{expr: `
		cel.bind(m, {'hello': 'world', 'hello!': 'worlds'},
		  m.all(k, v, k.startsWith('hello') && v.endsWith('world'))
		) == false
		`},
		// map.exists()
		{expr: `
		cel.bind(m, {'hello': 'world', 'hello!': 'worlds'},
		  m.exists(k, v, k.startsWith('hello') && v.endsWith('world'))
		)
		`},
		// map.existsOne()
		{expr: `
		cel.bind(m, {'hello': 'world', 'hello!': 'worlds'},
		  m.existsOne(k, v, k.startsWith('hello') && v.endsWith('world'))
		)
		`},
		// map.exists_one()
		{expr: `
		cel.bind(m, {'hello': 'world', 'hello!': 'worlds'},
		  m.exists_one(k, v, k.startsWith('hello') && v.endsWith('world'))
		)
		`},
		{expr: `
		cel.bind(m, {'hello': 'world', 'hello!': 'wow, world'},
		  m.exists_one(k, v, k.startsWith('hello') && v.endsWith('world'))
		) == false
		`},
		// map.transformList()
		{expr: `
		cel.bind(m, {'Hello': 'world'},
		  m.transformList(k, v, "%s=%s".format([k.lowerAscii(), v]))
		) == ["hello=world"]
		`},
		{expr: `
		cel.bind(m, {'hello': 'world'},
		  m.transformList(k, v, k.startsWith('greeting'), "%s=%s".format([k, v]))
		) == []
		`},
		// map.transformMap()
		{expr: `
		cel.bind(m, {'hello': 'world', 'goodbye': 'cruel world'},
			m.transformMap(k, v, "%s, %s!".format([k, v]))
		) == {'hello': 'hello, world!', 'goodbye': 'goodbye, cruel world!'}
		`},
		{expr: `
		cel.bind(m, {'hello': 'world', 'goodbye': 'cruel world'},
			m.transformMap(k, v, v.startsWith('world'), "%s, %s!".format([k, v]))
		) == {'hello': 'hello, world!'}
		`},
	}

	env := testCompreEnv(t)
	for i, tst := range compreTests {
		tc := tst
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			var asts []*cel.Ast
			pAst, iss := env.Parse(tc.expr)
			if iss.Err() != nil {
				t.Fatalf("env.Parse(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, pAst)
			cAst, iss := env.Check(pAst)
			if iss.Err() != nil {
				t.Fatalf("env.Check(%v) failed: %v", tc.expr, iss.Err())
			}
			asts = append(asts, cAst)

			for _, ast := range asts {
				prg, err := env.Program(ast)
				if err != nil {
					t.Fatalf("env.Program() failed: %v", err)
				}
				out, _, err := prg.Eval(cel.NoVars())
				if err != nil {
					t.Fatalf("prg.Eval() failed: %v", err)
				}
				if out.Value() != true {
					t.Errorf("prg.Eval() got %v, wanted true for expr: %s", out.Value(), tc.expr)
				}
			}
		})
	}
}

func testCompreEnv(t *testing.T, opts ...cel.EnvOption) *cel.Env {
	t.Helper()
	baseOpts := []cel.EnvOption{
		TwoVarComprehensions(),
		Bindings(),
		Strings(),
		cel.OptionalTypes(),
		cel.EnableMacroCallTracking()}
	env, err := cel.NewEnv(append(baseOpts, opts...)...)
	if err != nil {
		t.Fatalf("cel.NewEnv(TwoVarComprehensions()) failed: %v", err)
	}
	return env
}
