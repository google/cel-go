package checker

import (
	"celgo/ast"
	"celgo/semantics"
)

type semanticAdorner struct {
	s *semantics.Semantics
}

var _ ast.DebugAdorner = &semanticAdorner{}

func (a *semanticAdorner) GetMetadata(e ast.Expression) string {
	result := ""

	t := a.s.GetType(e)
	if t != nil {
		result += "~"
		result += t.String()
	}

	var ref semantics.Reference = nil
	switch e.(type) {
	case *ast.IdentExpression, *ast.CallExpression, *ast.CreateMessageExpression, *ast.SelectExpression:
		ref = a.s.GetReference(e)
		if ref != nil {
			if iref, ok := ref.(*semantics.IdentReference); ok {
				result += "^" + iref.Name()
			} else if fref, ok := ref.(*semantics.FunctionReference); ok {
				for i, overload := range fref.Overloads() {
					if i == 0 {
						result += "^"
					} else {
						result += "|"
					}
					result += overload
				}
			}
		}
	}

	return result
}

func print(e ast.Expression, s *semantics.Semantics) string {
	a := semanticAdorner{
		s: s,
	}

	return ast.ToAdornedDebugString(e, &a)
}
