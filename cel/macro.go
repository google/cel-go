// Copyright 2022 Google LLC
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

package cel

import "github.com/google/cel-go/parser"

type Macro = parser.Macro

// MacroExpander converts a call and its associated arguments into a new CEL abstract syntax tree, or an error
// if the input arguments are not suitable for the expansion requirements for the macro in question.
//
// The MacroExpander accepts as arguments a MacroExprHelper as well as the arguments used in the function call
// and produces as output an Expr ast node.
//
// Note: when the Macro.IsReceiverStyle() method returns true, the target argument will be nil.
type MacroExpander = parser.MacroExpander

// MacroExprHelper exposes helper methods for creating new expressions within a CEL abstract syntax tree.
type MacroExprHelper = parser.ExprHelper

var (
	// Factory functions for creating new cel.Macro instances that will match against the CEL abstract syntax tree
	// at parse time and then be expanded into an alternative AST defined by the MacroExpander

	NewGlobalMacro         = parser.NewGlobalMacro
	NewReceiverMacro       = parser.NewReceiverMacro
	NewGlobalVarArgMacro   = parser.NewGlobalVarArgMacro
	NewReceiverVarArgMacro = parser.NewReceiverVarArgMacro

	// Aliases to the functions used to create the CEL standard macros.

	MakeHas       = parser.MakeHas
	MakeExists    = parser.MakeExists
	MakeExistsOne = parser.MakeExistsOne
	MakeFilter    = parser.MakeFilter
	MakeMap       = parser.MakeMap

	// Aliases to each macro in the CEL standard environment.

	// HasMacro expands "has(m.f)" which tests the presence of a field, avoiding the need to
	// specify the field as a string.
	HasMacro = parser.HasMacro

	// AllMacro expands "range.all(var, predicate)" into a comprehension which ensures that all
	// elements in the range satisfy the predicate.
	AllMacro = parser.AllMacro

	// ExistsMacro expands "range.exists(var, predicate)" into a comprehension which ensures that
	// some element in the range satisfies the predicate.
	ExistsMacro = parser.ExistsMacro

	// ExistsOneMacro expands "range.exists_one(var, predicate)", which is true if for exactly one
	// element in range the predicate holds.
	ExistsOneMacro = parser.ExistsOneMacro

	// MapMacro expands "range.map(var, function)" into a comprehension which applies the function
	// to each element in the range to produce a new list.
	MapMacro = parser.MapMacro

	// MapFilterMacro expands "range.map(var, predicate, function)" into a comprehension which
	// first filters the elements in the range by the predicate, then applies the transform function
	// to produce a new list.
	MapFilterMacro = parser.MapFilterMacro

	// FilterMacro expands "range.filter(var, predicate)" into a comprehension which filters
	// elements in the range, producing a new list from the elements that satisfy the predicate.
	FilterMacro = parser.FilterMacro

	// StandardMacros provides an alias to all the CEL macros defined in the standard environment.
	StandardMacros = []Macro{
		HasMacro, AllMacro, ExistsMacro, ExistsOneMacro, MapMacro, MapFilterMacro, FilterMacro,
	}

	// NoMacros provides an alias to an empty list of macros
	NoMacros = []Macro{}
)
