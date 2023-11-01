// Copyright 2023 Google LLC
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

import (
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// StaticOptimizer contains a sequence of ASTOptimizer instances which will be applied in order.
//
// The static optimizer normalizes expression ids and type-checking run between optimization
// passes to ensure that the final optimized output is a valid expression with metadata consistent
// with what would have been generated from a parsed and checked expression.
//
// Note: source position information is best-effort and likely wrong, but optimized expressions
// should be suitable for calls to parser.Unparse.
type StaticOptimizer struct {
	optimizers []ASTOptimizer
}

// NewStaticOptimizer creates a StaticOptimizer with a sequence of ASTOptimizer's to be applied
// to a checked expression.
func NewStaticOptimizer(optimizers ...ASTOptimizer) *StaticOptimizer {
	return &StaticOptimizer{
		optimizers: optimizers,
	}
}

// Optimize applies a sequence of optimizations to an Ast within a given environment.
//
// If issues are encountered, the Issues.Err() return value will be non-nil.
func (opt *StaticOptimizer) Optimize(env *Env, a *Ast) (*Ast, *Issues) {
	// Make a copy of the AST to be optimized.
	optimized := ast.Copy(a.impl)
	ids := newIDGenerator(ast.MaxID(a.impl))

	// Create the optimizer context, could be pooled in the future.
	issues := NewIssues(common.NewErrors(a.Source()))
	baseFac := ast.NewExprFactory()
	exprFac := &optimizerExprFactory{
		idGenerator: ids,
		fac:         baseFac,
		sourceInfo:  optimized.SourceInfo(),
	}
	ctx := &OptimizerContext{
		optimizerExprFactory: exprFac,
		Env:                  env,
		Issues:               issues,
	}

	// Apply the optimizations sequentially.
	for _, o := range opt.optimizers {
		optimized = o.Optimize(ctx, optimized)
		if issues.Err() != nil {
			return nil, issues
		}
		// Normalize expression id metadata including coordination with macro call metadata.
		stableIDGen := newIDGenerator(0)
		info := optimized.SourceInfo()
		expr := optimized.Expr()
		normalizeIDs(stableIDGen.renumberStable, expr, info)

		// Sanitize the macro call references once the optimized expression has been computed
		// and the ids normalized between the expression and the macros.
		sanitized := baseFac.CopyExpr(expr)
		sanitizedExprMap := make(map[int64]ast.Expr)
		ast.PostOrderVisit(sanitized, ast.NewExprVisitor(func(e ast.Expr) {
			if _, found := info.GetMacroCall(e.ID()); found {
				e.SetKindCase(nil)
			}
			sanitizedExprMap[e.ID()] = baseFac.CopyExpr(e)
		}))
		// Update the macro call id references to ensure that macro pointers are
		// updated consistently across macros.
		for id, call := range info.MacroCalls() {
			resetMacroCall(call, sanitizedExprMap)
			info.SetMacroCall(id, call)
		}

		// Recheck the updated expression for any possible type-agreement or validation errors.
		parsed := &Ast{
			source: a.Source(),
			impl:   ast.NewAST(expr, info)}
		checked, iss := ctx.Check(parsed)
		if iss.Err() != nil {
			return nil, iss
		}
		optimized = checked.impl
	}

	// Return the optimized result.
	return &Ast{
		source: a.Source(),
		impl:   optimized,
	}, nil
}

// normalizeIDs ensures that the metadata present with an AST is reset in a manner such
// that the ids within the expression correspond to the ids within macros.
func normalizeIDs(idGen ast.IDGenerator, optimized ast.Expr, info *ast.SourceInfo) {
	optimized.RenumberIDs(idGen)

	// First, update the macro call ids themselves.
	for id, call := range info.MacroCalls() {
		info.ClearMacroCall(id)
		callID := idGen(id)
		info.SetMacroCall(callID, call)
	}
	// Then update the macro call definitions which refer to these ids
	for id, call := range info.MacroCalls() {
		call.RenumberIDs(idGen)
		info.SetMacroCall(id, call)
	}
}

func resetMacroCall(call ast.Expr, sanitizedExprMap map[int64]ast.Expr) {
	// Identify the set of expressions in the core expression which were updated,
	// excluding nodes which correspond to macros.
	ast.PostOrderVisit(call, ast.NewExprVisitor(func(e ast.Expr) {
		if update, found := sanitizedExprMap[e.ID()]; found {
			e.SetKindCase(update)
		}
	}))
}

// newIDGenerator ensures that new ids are only created the first time they are encountered.
func newIDGenerator(seed int64) *idGenerator {
	return &idGenerator{
		idMap: make(map[int64]int64),
		seed:  seed,
	}
}

type idGenerator struct {
	idMap map[int64]int64
	seed  int64
}

func (gen *idGenerator) nextID() int64 {
	gen.seed++
	return gen.seed
}

func (gen *idGenerator) renumberMonotonic(id int64) int64 {
	if id == 0 {
		return 0
	}
	return gen.nextID()
}

func (gen *idGenerator) renumberStable(id int64) int64 {
	if id == 0 {
		return 0
	}
	if newID, found := gen.idMap[id]; found {
		return newID
	}
	nextID := gen.nextID()
	gen.idMap[id] = nextID
	return nextID
}

// OptimizerContext embeds Env and Issues instances to make it easy to type-check and evaluate
// subexpressions and report any errors encountered along the way. The context also embeds the
// optimizerExprFactory which can be used to generate new sub-expressions with expression ids
// consistent with the expectations of a parsed expression.
type OptimizerContext struct {
	*Env
	*optimizerExprFactory
	*Issues
}

// ASTOptimizer applies an optimization over an AST and returns the optimized result.
type ASTOptimizer interface {
	// Optimize optimizes a type-checked AST within an Environment and accumulates any issues.
	Optimize(*OptimizerContext, *ast.AST) *ast.AST
}

type optimizerExprFactory struct {
	*idGenerator
	fac        ast.ExprFactory
	sourceInfo *ast.SourceInfo
}

// CopyAST creates a renumbered copy of `Expr` and `SourceInfo` values of the input AST, where the
// renumbering uses the same scheme as the core optimizer logic ensuring there are no collisions
// between copies.
//
// Use this method before attempting to merge the expression from AST into another.
func (opt *optimizerExprFactory) CopyAST(a *ast.AST) (ast.Expr, *ast.SourceInfo) {
	idGen := newIDGenerator(opt.nextID())
	defer func() { opt.seed = idGen.nextID() }()
	copyExpr := opt.fac.CopyExpr(a.Expr())
	copyInfo := ast.CopySourceInfo(a.SourceInfo())
	normalizeIDs(idGen.renumberStable, copyExpr, copyInfo)
	return copyExpr, copyInfo
}

// CopyExpr copies the structure of the input ast.Expr and renumbers the identifiers in a manner
// consistent with the CEL parser / checker.
func (opt *optimizerExprFactory) CopyExpr(e ast.Expr) ast.Expr {
	copy := opt.fac.CopyExpr(e)
	copy.RenumberIDs(opt.renumberMonotonic)
	return copy
}

// NewBindMacro creates a cel.bind() call with a variable name, initialization expression, and remaining expression.
//
// Note: the macroID indicates the insertion point, the call id that matched the macro signature, which will be used
// for coordinating macro metadata with the bind call. This piece of data is what makes it possible to unparse
// optimized expressions which use the bind() call.
//
// Example:
//
// cel.bind(myVar, a && b || c, !myVar || (myVar && d))
// - varName: myVar
// - varInit: a && b || c
// - remaining: !myVar || (myVar && d)
func (opt *optimizerExprFactory) NewBindMacro(macroID int64, varName string, varInit, remaining ast.Expr) ast.Expr {
	bindID := opt.nextID()
	varID := opt.nextID()

	var bindVarInit ast.Expr
	varInit, bindVarInit = opt.sanitizeMacroExpr(varInit)

	var bindRemaining ast.Expr
	remaining, bindRemaining = opt.sanitizeMacroExpr(remaining)

	// Place the expanded macro form in the macro calls list so that the inlined
	// call can be unparsed.
	opt.sourceInfo.SetMacroCall(macroID,
		opt.fac.NewMemberCall(0, "bind",
			opt.fac.NewIdent(opt.nextID(), "cel"),
			opt.fac.NewIdent(varID, varName),
			bindVarInit,
			bindRemaining))

	// Replace the parent node with the intercepted inlining using cel.bind()-like
	// generated comprehension AST.
	return opt.fac.NewComprehension(bindID,
		opt.fac.NewList(opt.nextID(), []ast.Expr{}, []int32{}),
		"#unused",
		varName,
		opt.fac.CopyExpr(varInit),
		opt.fac.NewLiteral(opt.nextID(), types.False),
		opt.fac.NewIdent(varID, varName),
		opt.fac.CopyExpr(remaining))
}

// NewCall creates a global function call invocation expression.
//
// Example:
//
// countByField(list, fieldName)
// - function: countByField
// - args: [list, fieldName]
func (opt *optimizerExprFactory) NewCall(function string, args ...ast.Expr) ast.Expr {
	return opt.fac.NewCall(opt.nextID(), function, args...)
}

// NewMemberCall creates a member function call invocation expression where 'target' is the receiver of the call.
//
// Example:
//
// list.countByField(fieldName)
// - function: countByField
// - target: list
// - args: [fieldName]
func (opt *optimizerExprFactory) NewMemberCall(function string, target ast.Expr, args ...ast.Expr) ast.Expr {
	return opt.fac.NewMemberCall(opt.nextID(), function, target, args...)
}

// NewIdent creates a new identifier expression.
//
// Examples:
//
// - simple_var_name
// - qualified.subpackage.var_name
func (opt *optimizerExprFactory) NewIdent(name string) ast.Expr {
	return opt.fac.NewIdent(opt.nextID(), name)
}

// NewLiteral creates a new literal expression value.
//
// The range of valid values for a literal generated during optimization is different than for expressions
// generated via parsing / type-checking, as the ref.Val may be _any_ CEL value so long as the value can
// be converted back to a literal-like form.
func (opt *optimizerExprFactory) NewLiteral(value ref.Val) ast.Expr {
	return opt.fac.NewLiteral(opt.nextID(), value)
}

// NewList creates a list expression with a set of optional indices.
//
// Examples:
//
// [a, b]
// - elems: [a, b]
// - optIndices: []
//
// [a, ?b, ?c]
// - elems: [a, b, c]
// - optIndices: [1, 2]
func (opt *optimizerExprFactory) NewList(elems []ast.Expr, optIndices []int32) ast.Expr {
	return opt.fac.NewList(opt.nextID(), elems, optIndices)
}

// NewMap creates a map from a set of entry expressions which contain a key and value expression.
func (opt *optimizerExprFactory) NewMap(entries []ast.EntryExpr) ast.Expr {
	return opt.fac.NewMap(opt.nextID(), entries)
}

// NewMapEntry creates a map entry with a key and value expression and a flag to indicate whether the
// entry is optional.
//
// Examples:
//
// {a: b}
// - key: a
// - value: b
// - optional: false
//
// {?a: ?b}
// - key: a
// - value: b
// - optional: true
func (opt *optimizerExprFactory) NewMapEntry(key, value ast.Expr, isOptional bool) ast.EntryExpr {
	return opt.fac.NewMapEntry(opt.nextID(), key, value, isOptional)
}

// NewPresenceTest creates a new presence test macro call.
//
// Example:
//
// has(msg.field_name)
// - operand: msg
// - field: field_name
func (opt *optimizerExprFactory) NewPresenceTest(macroID int64, operand ast.Expr, field string) ast.Expr {
	// Copy the input operand and renumber it.
	var hasOperand ast.Expr
	operand, hasOperand = opt.sanitizeMacroExpr(operand)

	// Place the expanded macro form in the macro calls list so that the inlined call can be unparsed.
	opt.sourceInfo.SetMacroCall(macroID,
		opt.fac.NewCall(0, "has",
			opt.fac.NewSelect(opt.nextID(), hasOperand, field)))

	// Generate a new presence test macro.
	return opt.fac.NewPresenceTest(opt.nextID(), opt.CopyExpr(operand), field)
}

// NewSelect creates a select expression where a field value is selected from an operand.
//
// Example:
//
// msg.field_name
// - operand: msg
// - field: field_name
func (opt *optimizerExprFactory) NewSelect(operand ast.Expr, field string) ast.Expr {
	return opt.fac.NewSelect(opt.nextID(), operand, field)
}

// NewStruct creates a new typed struct value with an set of field initializations.
//
// Example:
//
// pkg.TypeName{field: value}
// - typeName: pkg.TypeName
// - fields: [{field: value}]
func (opt *optimizerExprFactory) NewStruct(typeName string, fields []ast.EntryExpr) ast.Expr {
	return opt.fac.NewStruct(opt.nextID(), typeName, fields)
}

// NewStructField creates a struct field initialization.
//
// Examples:
//
// {count: 3u}
// - field: count
// - value: 3u
// - optional: false
//
// {?count: x}
// - field: count
// - value: x
// - optional: true
func (opt *optimizerExprFactory) NewStructField(field string, value ast.Expr, isOptional bool) ast.EntryExpr {
	return opt.fac.NewStructField(opt.nextID(), field, value, isOptional)
}

// sanitizeMacroExpr copies the input expression, renumbers it, and also generates a sanitized version
// suitable for use within macro bodies where the body must not contain the content of another macro,
// but rather a 'pointer' consisting of an empty expression node with an id.
func (opt *optimizerExprFactory) sanitizeMacroExpr(baseExpr ast.Expr) (copyExpr, macroExpr ast.Expr) {
	// Something is off in this logic
	// - renumber the base expression using a stable identifier
	// - base expression changes need to be reflected in macros which might have referenced
	//   the old ids
	// - old and new call entries don't appear to have been updated???
	idGen := newIDGenerator(opt.nextID())
	defer func() { opt.seed = idGen.nextID() }()
	copyExpr = opt.fac.CopyExpr(baseExpr)
	copyExpr.RenumberIDs(idGen.renumberStable)

	// Traverse the base expression and determine whether a macro id was updated by using
	// the stable id generator to verify the id move.
	oldToNewMacroIDs := make(map[int64]int64)
	newToOldMacroIDs := make(map[int64]int64)
	ast.PreOrderVisit(baseExpr, ast.NewExprVisitor(func(e ast.Expr) {
		if call, isMacroRef := opt.sourceInfo.GetMacroCall(e.ID()); isMacroRef {
			newID := idGen.renumberStable(e.ID())
			newToOldMacroIDs[newID] = e.ID()
			oldToNewMacroIDs[e.ID()] = newID
			opt.sourceInfo.SetMacroCall(newID, call)
			opt.sourceInfo.ClearMacroCall(e.ID())
		}
	}))

	// Clear the expression nodes which correspond to other macros from the macro-sanitized expression.
	macroExpr = opt.fac.CopyExpr(copyExpr)
	ast.PreOrderVisit(macroExpr, ast.NewExprVisitor(func(e ast.Expr) {
		if _, isMacroRef := newToOldMacroIDs[e.ID()]; isMacroRef {
			e.SetKindCase(nil)
		}
	}))

	// Macro ids were renumbered during the copy, but not within the macro calls themselves, so this
	// step ensures that nested macro references are updated as well.
	updatedIDVisitor := ast.NewExprVisitor(func(e ast.Expr) {
		if newID, found := oldToNewMacroIDs[e.ID()]; found {
			e.RenumberIDs(func(int64) int64 {
				return newID
			})
		}
	})
	for _, call := range opt.sourceInfo.MacroCalls() {
		ast.PostOrderVisit(call, updatedIDVisitor)
	}
	return
}
