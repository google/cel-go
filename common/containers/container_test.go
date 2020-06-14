// Copyright 2020 Google LLC
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

package containers

import (
	"reflect"
	"testing"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func TestContainer_ResolveCandidateNames(t *testing.T) {
	c, err := NewContainer(Name("a.b.c.M.N"))
	if err != nil {
		t.Fatal(err)
	}
	names := c.ResolveCandidateNames("R.s")
	want := []string{
		"a.b.c.M.N.R.s",
		"a.b.c.M.R.s",
		"a.b.c.R.s",
		"a.b.R.s",
		"a.R.s",
		"R.s",
	}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("got %v, wanted %v", names, want)
	}
}

func TestContainer_ResolveCandidateNames_FullyQualifiedName(t *testing.T) {
	c, err := NewContainer(Name("a.b.c.M.N"))
	if err != nil {
		t.Fatal(err)
	}
	// The leading '.' indicates the name is already fully-qualified.
	names := c.ResolveCandidateNames(".R.s")
	want := []string{"R.s"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("got %v, wanted %v", names, want)
	}
}

func TestContainer_ResolveCandidateNames_EmptyContainer(t *testing.T) {
	names := DefaultContainer.ResolveCandidateNames("R.s")
	want := []string{"R.s"}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("got %v, wanted %v", names, want)
	}
}

func TestContainer_Aliases(t *testing.T) {
	c, err := NewContainer(Name("a.b.c"), Aliases("my.alias.R"))
	if err != nil {
		t.Fatal(err)
	}
	names := c.ResolveCandidateNames("R")
	want := []string{
		"a.b.c.R",
		"a.b.R",
		"a.R",
		"R",
		"my.alias.R",
	}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("got %v, wanted %v", names, want)
	}
	names = c.ResolveCandidateNames("S")
	want = []string{
		"a.b.c.S",
		"a.b.S",
		"a.S",
		"S",
	}
	if !reflect.DeepEqual(names, want) {
		t.Errorf("got %v, wanted %v", names, want)
	}
}

func TestContainer_Aliases_Errors(t *testing.T) {
	_, err := NewContainer(Name("a.b.c.M.N"), Aliases("my.alias.R", "yer.other.R"))
	wantErr := "alias collides with existing reference: " +
		"name=yer.other.R, alias=R, existing=my.alias.R"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got error %v, expected %s.", err, wantErr)
	}

	_, err = NewContainer(Name("a.b.c.M.N"), Aliases("my.alias.a", "yer.other.b"))
	wantErr = "alias collides with container name: name=my.alias.a, alias=a, container=a.b.c.M.N"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got error %v, expected %s.", err, wantErr)
	}

	_, err = NewContainer(Name("a.b.c.M.N"), Aliases(".bad"))
	wantErr = "invalid qualified name: .bad, wanted name of the form 'qualified.name'"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got error %v, expected %s.", err, wantErr)
	}

	_, err = NewContainer(Name("a.b.c.M.N"), Aliases("bad.alias."))
	wantErr = "invalid qualified name: bad.alias., wanted name of the form 'qualified.name'"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got error %v, expected %s.", err, wantErr)
	}

	_, err = NewContainer(Name("a.b.c.M.N"), AliasAs("a", "b"))
	wantErr = "aliases must refer to qualified names: a"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got error %v, expected %s.", err, wantErr)
	}

	_, err = NewContainer(Name("a.b.c.M.N"), AliasAs("my.alias", "b.c"))
	wantErr = "alias names must be non-empty and simple (not qualified): alias=b.c"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got error %v, expected %s.", err, wantErr)
	}
}

func TestContainers_Extend_Alias(t *testing.T) {
	c, err := DefaultContainer.Extend(AliasAs("test.alias", "alias"))
	if err != nil {
		t.Fatal(err)
	}
	if c.aliasSet()["alias"] != "test.alias" {
		t.Errorf("got alias %v wanted 'test.alias'", c.aliasSet())
	}
	c, err = c.Extend(Name("with.container"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Name() != "with.container" {
		t.Errorf("got container name %s, wanted 'goodbye.container'", c.Name())
	}
	if c.aliasSet()["alias"] != "test.alias" {
		t.Errorf("got alias %v wanted 'test.alias'", c.aliasSet())
	}
}

func TestContainers_Extend_Name(t *testing.T) {
	c, err := DefaultContainer.Extend(Name(""))
	if err != nil {
		t.Fatal(err)
	}
	if c != nil {
		t.Error("got non-nil container, wanted nil")
	}
	c, err = DefaultContainer.Extend(Name("hello.container"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Name() != "hello.container" {
		t.Errorf("got container name %s, wanted 'hello.container'", c.Name())
	}
	c, err = c.Extend(Name("goodbye.container"))
	if err != nil {
		t.Fatal(err)
	}
	if c.Name() != "goodbye.container" {
		t.Errorf("got container name %s, wanted 'goodbye.container'", c.Name())
	}
}

func TestContainers_ToQualifiedName(t *testing.T) {
	ident := &exprpb.Expr{
		ExprKind: &exprpb.Expr_IdentExpr{
			IdentExpr: &exprpb.Expr_Ident{
				Name: "var",
			},
		},
	}
	idName, found := ToQualifiedName(ident)
	if !found {
		t.Errorf("got not found from %v expr, wanted found", ident)
	}
	if idName != "var" {
		t.Errorf("got %v, wanted 'var'", idName)
	}
	sel := &exprpb.Expr{
		ExprKind: &exprpb.Expr_SelectExpr{
			SelectExpr: &exprpb.Expr_Select{
				Operand: ident,
				Field:   "qualifier",
			},
		},
	}
	qualName, found := ToQualifiedName(sel)
	if !found {
		t.Errorf("got not found from %v expr, wanted found", sel)
	}
	if qualName != "var.qualifier" {
		t.Errorf("got %v, wanted 'var.qualifier'", qualName)
	}
	unary := &exprpb.Expr{
		ExprKind: &exprpb.Expr_CallExpr{
			CallExpr: &exprpb.Expr_Call{
				Function: "!_",
				Args:     []*exprpb.Expr{ident},
			},
		},
	}
	sel.GetSelectExpr().Operand = unary
	_, found = ToQualifiedName(sel)
	if found {
		t.Errorf("got found, wanted not found for %v", sel)
	}
}
