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

// Package containers defines types and functions for resolving qualified names within a namespace
// or type provided to CEL.
package containers

import (
	"fmt"
	"strings"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

var (
	// DefaultContainer has an empty container name.
	DefaultContainer *Container = nil

	// Empty map to search for aliases when needed.
	noAliases = make(map[string]string)
)

// NewContainer creates a new Container with the fully-qualified name.
func NewContainer(opts ...ContainerOption) (*Container, error) {
	var c *Container
	var err error
	for _, opt := range opts {
		c, err = opt(c)
		if err != nil {
			return nil, err
		}
	}
	return c, nil
}

// Container holds a reference to an optional qualified container name and set of aliases.
//
// The program container can be used to simplify variable, function, and type specification within
// CEL programs and behaves more or less like a C++ namespace. See ResolveCandidateNames for more
// details.
type Container struct {
	name    string
	aliases map[string]string
}

// Extend creates a new Container with the existing settings and applies a series of
// ContainerOptions to further configure the new container.
func (c *Container) Extend(opts ...ContainerOption) (*Container, error) {
	if c == nil {
		return NewContainer(opts...)
	}
	// Copy the name and aliases of the existing container.
	ext := &Container{name: c.Name()}
	if len(c.aliasSet()) > 0 {
		aliasSet := make(map[string]string, len(c.aliasSet()))
		for k, v := range c.aliasSet() {
			aliasSet[k] = v
		}
		ext.aliases = aliasSet
	}
	// Apply the new options to the container.
	var err error
	for _, opt := range opts {
		ext, err = opt(ext)
		if err != nil {
			return nil, err
		}
	}
	return ext, nil
}

// Name returns the fully-qualified name of the container.
//
// The name may conceptually be a namespace, package, or type.
func (c *Container) Name() string {
	if c == nil {
		return ""
	}
	return c.name
}

// ResolveCandidateNames returns the candidates name of namespaced identifiers in C++ resolution
// order.
//
// Names which shadow other names are returned first. If a name includes a leading dot ('.'),
// the name is treated as an absolute identifier which cannot be shadowed.
//
// Given a container name a.b.c.M.N and a type name R.s, this will deliver in order:
//
//     a.b.c.M.N.R.s
//     a.b.c.M.R.s
//     a.b.c.R.s
//     a.b.R.s
//     a.R.s
//     R.s
//
// If aliases are configured for the container, then alias names appear after the potential set
// of container-based names.
func (c *Container) ResolveCandidateNames(name string) []string {
	if strings.HasPrefix(name, ".") {
		qn := name[1:]
		return c.candidatesWithAlias([]string{qn}, qn)
	}
	if c.Name() == "" {
		return c.candidatesWithAlias([]string{name}, name)
	}

	nextCont := c.name
	candidates := []string{nextCont + "." + name}
	for i := strings.LastIndex(nextCont, "."); i >= 0; i = strings.LastIndex(nextCont, ".") {
		nextCont = nextCont[:i]
		candidates = append(candidates, nextCont+"."+name)
	}
	candidates = append(candidates, name)
	return c.candidatesWithAlias(candidates, name)
}

// aliasSet returns the alias -> fully-qualified name mapping stored in the container.
func (c *Container) aliasSet() map[string]string {
	if c == nil || c.aliases == nil {
		return noAliases
	}
	return c.aliases
}

// candidatesWithAliases returns the resolved candidates from the container with any applicable
// aliases appended on the end of the list.
func (c *Container) candidatesWithAlias(candidates []string, name string) []string {
	if len(c.aliasSet()) == 0 {
		return candidates
	}
	// If an alias exists for the name, ensure it is searched last.
	alias, found := c.aliasSet()[name]
	if found {
		return append(candidates, alias)
	}
	return candidates
}

// ContainerOption specifies a functional configuration option for a Container.
type ContainerOption func(*Container) (*Container, error)

// Aliases configures a set of simple names as aliases for fully-qualified names.
//
// An alias can be useful when working with variables, functions, and especially types from
// multiple namespaces:
//
//    // CEL object construction
//    qual.pkg.version.ObjTypeName{
//       field: alt.container.ver.FieldTypeName{value: ...}
//    }
//
// Only one qualified path may be used as CEL container, so at least one of these references is
// a long qualified name within an otherwise short CEL program. Using the following aliases, the
// program becomes much simpler:
//
//    // CEL Go option
//    Aliases("qual.pkg.version.ObjTypeName", "alt.container.ver.FieldTypeName")
//    // Simplified Object construction
//    ObjTypeName{field: FieldTypeName{value: ...}}
//
// There are a few rules for the qualified names and the simple name aliases generated from them:
// - Qualified names must be dot-delimited, e.g. `package.subpkg.name`.
// - The last element in the qualified name is the simple name used as an alias.
// - Alias names must not collide with each other.
// - The alias name must not collide with the root-level container name.
//
// Aliases are distinct from container-based references in the following important ways:
// - Containers follow C++ namespace resolution rules with searches from the most qualified name
//   to the least qualified name.
// - Container references within the CEL program may be relative, and are resolved to fully
//   qualified names at either type-check time or program plan time, whichever comes first.
// - Alias simple names must resolve to a fully-qualified name.
// - Resolved aliases do not participate in namespace resolution.
// - Resolved aliases are searched after container names, including container names in the global
//   scope.
//
// If there is ever a case where an identifier could be in both the container and in the alias,
// the container wins as the container will continue to evolve over time and the program must be
// forward compatible with changes in the container.
func Aliases(qualifiedNames ...string) ContainerOption {
	return func(c *Container) (*Container, error) {
		for _, qn := range qualifiedNames {
			alias := qn
			ind := strings.LastIndex(qn, ".")
			if ind <= 0 || ind >= len(qn)-1 {
				return nil, fmt.Errorf(
					"invalid qualified name: %s, wanted name of the form 'qualified.name'", qn)
			}
			alias = qn[ind+1:]
			var err error
			c, err = AliasAs(qn, alias)(c)
			if err != nil {
				return nil, err
			}
		}
		return c, nil
	}
}

// AliasAs associates a fully-qualified name with a user-defined alias.
//
// In general, Aliases is preferred to AliasAs since the names generated from the Aliases option
// are more easily traced back to source code. The AliasAs option is useful for propagating alias
// configuration from one Container instance to another, and may also be useful for remapping
// poorly chosen protobuf message / package names.
//
// Note: all of the rules that apply to Aliases also apply to AliasAs.
func AliasAs(qualifiedName, alias string) ContainerOption {
	return func(c *Container) (*Container, error) {
		if len(alias) <= 0 || strings.Contains(alias, ".") {
			return nil, fmt.Errorf(
				"alias names must be non-empty and simple (not qualified): alias=%s", alias)
		}
		ind := strings.LastIndex(qualifiedName, ".")
		if ind <= 0 || ind == len(qualifiedName)-1 {
			return nil, fmt.Errorf("aliases must refer to qualified names: %s",
				qualifiedName)
		}
		aliasRef, found := c.aliasSet()[alias]
		if found {
			return nil, fmt.Errorf(
				"alias collides with existing reference: name=%s, alias=%s, existing=%s",
				qualifiedName, alias, aliasRef)
		}
		if strings.HasPrefix(c.Name(), alias+".") || c.Name() == alias {
			return nil, fmt.Errorf(
				"alias collides with container name: name=%s, alias=%s, container=%s",
				qualifiedName, alias, c.Name())
		}
		if c == nil {
			c = &Container{}
		}
		if c.aliases == nil {
			c.aliases = make(map[string]string)
		}
		c.aliases[alias] = qualifiedName
		return c, nil
	}
}

// Name sets the fully-qualified name of the Container.
func Name(name string) ContainerOption {
	return func(c *Container) (*Container, error) {
		if c.Name() == name {
			return c, nil
		}
		if c == nil {
			return &Container{name: name}, nil
		}
		c.name = name
		return c, nil
	}
}

// ToQualifiedName converts an expression AST into a qualified name if possible, with a boolean
// 'found' value that indicates if the conversion is successful.
func ToQualifiedName(e *exprpb.Expr) (string, bool) {
	switch e.ExprKind.(type) {
	case *exprpb.Expr_IdentExpr:
		id := e.GetIdentExpr()
		return id.Name, true
	case *exprpb.Expr_SelectExpr:
		sel := e.GetSelectExpr()
		if qual, found := ToQualifiedName(sel.Operand); found {
			return qual + "." + sel.Field, true
		}
	}
	return "", false
}
