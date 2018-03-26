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
	"strings"

	"github.com/google/cel-go/ast"
)

//
// Returns the candidates for name resolution of a name within a container(e.g. package, message,
// enum, service elements) context following PB (== C++) conventions. Iterates those names which
// shadow other names first; recognizes and removes a leading '.' for overriding shadowing. Given
// a container name a.b.c.M.N and a type name R.s, this will deliver in order
// a.b.c.M.N.R.s, a.b.c.M.R.s, a.b.c.R.s, a.b.R.s, a.R.s, R.s.
//
func qualifiedTypeNameCandidates(container string, typeName string) []string {
	if strings.HasPrefix(typeName, ".") {
		return []string{typeName[1:]}
	}

	if container == "" {
		return []string{typeName}
	}

	i := strings.LastIndex(container, ".")
	first := []string{container + "." + typeName}
	var rest []string
	if i >= 0 {
		rest = qualifiedTypeNameCandidates(container[:i], typeName)
	} else {
		rest = qualifiedTypeNameCandidates("", typeName)
	}

	return append(first, rest...)
}

// Attempt to interpret an expression as a qualified name. This traverses select and getIdent
// expression and returns the name they constitute, or null if the expression cannot be
// interpreted like this.
func asQualifiedName(e ast.Expression) (string, bool) {
	switch e.(type) {
	case *ast.IdentExpression:
		i := e.(*ast.IdentExpression)
		return i.Name, true
	case *ast.SelectExpression:
		s := e.(*ast.SelectExpression)
		if qname, found := asQualifiedName(s.Target); found {
			return qname + "." + s.Field, true
		}
	}
	return "", false
}
