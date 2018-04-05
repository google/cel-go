// Copyright 2018 Google LLC
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

package checker

import (
	"github.com/google/cel-go/debug"
	"github.com/google/cel-go/semantics"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type semanticAdorner struct {
	s *semantics.Semantics
}

var _ debug.DebugAdorner = &semanticAdorner{}

func (a *semanticAdorner) GetMetadata(elem interface{}) string {
	result := ""
	e, isExpr := elem.(*expr.Expr)
	if !isExpr {
		return result
	}
	t := a.s.GetType(e)
	if t != nil {
		result += "~"
		result += t.String()
	}

	var ref semantics.Reference = nil
	switch e.ExprKind.(type) {
	case *expr.Expr_IdentExpr,
	     *expr.Expr_CallExpr,
	     *expr.Expr_StructExpr,
	     *expr.Expr_SelectExpr:
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

func print(e *expr.Expr, s *semantics.Semantics) string {
	a := &semanticAdorner{s: s}
	return debug.ToAdornedDebugString(e, a)
}
