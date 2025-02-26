// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package policy

import (
	"cmp"
	"fmt"
	"slices"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

// ComposerOption is a functional option used to configure a RuleComposer
type ComposerOption func(*RuleComposer) (*RuleComposer, error)

// ExpressionUnnestHeight determines the height at which nested expressions are split into local
// variables within the cel.@block declaration.
func ExpressionUnnestHeight(height int) ComposerOption {
	return func(c *RuleComposer) (*RuleComposer, error) {
		if height <= 0 {
			return nil, fmt.Errorf("invalid unnest height: value must be positive: %d", height)
		}
		c.exprUnnestHeight = height
		return c, nil
	}
}

// NewRuleComposer creates a rule composer which stitches together rules within a policy into
// a single CEL expression.
func NewRuleComposer(env *cel.Env, opts ...ComposerOption) (*RuleComposer, error) {
	composer := &RuleComposer{
		env: env,
		// set the default nesting height to something reasonable.
		exprUnnestHeight: 25,
	}
	var err error
	for _, opt := range opts {
		composer, err = opt(composer)
		if err != nil {
			return nil, err
		}
	}
	return composer, nil
}

// RuleComposer optimizes a set of expressions into a single expression.
type RuleComposer struct {
	env *cel.Env

	// exprUnnestHeight determines the height at which nested matches are split into
	// index variables within a cel.@block index declaration when composing matches under
	// the first-match semantic.
	exprUnnestHeight int
}

// Compose stitches together a set of expressions within a CompiledRule into a single CEL ast.
func (c *RuleComposer) Compose(r *CompiledRule) (*cel.Ast, *cel.Issues) {
	ruleRoot, _ := c.env.Compile("true")
	composer := &ruleComposerImpl{
		rule:       r,
		varIndices: []varIndex{},
	}
	opt := cel.NewStaticOptimizer(composer)
	ast, iss := opt.Optimize(c.env, ruleRoot)
	if iss.Err() != nil {
		return nil, iss
	}
	unnester := &ruleUnnesterImpl{
		varIndices:       []varIndex{},
		exprUnnestHeight: c.exprUnnestHeight,
	}
	opt = cel.NewStaticOptimizer(unnester)
	return opt.Optimize(c.env, ast)
}

type varIndex struct {
	index    int
	indexVar string
	localVar string
	expr     ast.Expr
	celType  *types.Type
}

type ruleComposerImpl struct {
	rule         *CompiledRule
	nextVarIndex int
	varIndices   []varIndex
}

// Optimize implements an AST optimizer for CEL which composes an expression graph into a single
// expression value.
func (opt *ruleComposerImpl) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	// The input to optimize is a dummy expression which is completely replaced according
	// to the configuration of the rule composition graph.
	ruleExpr := opt.optimizeRule(ctx, opt.rule)

	// If there were no variables, return the expression.
	if len(opt.varIndices) == 0 {
		return ctx.NewAST(ruleExpr)
	}

	// Otherwise populate the cel.@block with the variable declarations and wrap the expression
	// in the block.
	varExprs := make([]ast.Expr, len(opt.varIndices))
	for i, vi := range opt.varIndices {
		varExprs[i] = vi.expr
		err := ctx.ExtendEnv(cel.Variable(vi.indexVar, vi.celType))
		if err != nil {
			ctx.ReportErrorAtID(ruleExpr.ID(), "%s", err.Error())
		}
	}
	blockExpr := ctx.NewCall("cel.@block", ctx.NewList(varExprs, []int32{}), ruleExpr)
	return ctx.NewAST(blockExpr)
}

func (opt *ruleComposerImpl) optimizeRule(ctx *cel.OptimizerContext, r *CompiledRule) ast.Expr {
	// Visitor to rewrite variables-prefixed identifiers with index names.
	vars := r.Variables()
	for _, v := range vars {
		opt.registerVariable(ctx, v)
	}

	matches := r.Matches()
	matchCount := len(matches)
	var output compositionStep = nil
	// If the rule has an optional output, the last result in the ternary should return
	// `optional.none`. This output is implicit and created here to reflect the desired
	// last possible output of this type of rule.
	if r.HasOptionalOutput() {
		output = newOptionalCompositionStep(ctx, ctx.NewLiteral(types.True), ctx.NewCall("optional.none"))
	}
	// Build the rule subgraph.
	for i := matchCount - 1; i >= 0; i-- {
		m := matches[i]
		cond := ctx.CopyASTAndMetadata(m.Condition().NativeRep())

		// If the output is non-nil, then it is considered a non-optional output since
		// it is explictly stated. If the rule itself is optional, then the base case value
		// of output being optional.none() will convert the non-optional value to an optional
		// one.
		if m.Output() != nil {
			out := ctx.CopyASTAndMetadata(m.Output().Expr().NativeRep())
			step := newNonOptionalCompositionStep(ctx, cond, out)
			output = step.combine(output)
			continue
		}

		// If the match has a nested rule, then compute the rule and whether it has
		// an optional return value.
		child := m.NestedRule()
		nestedRule := opt.optimizeRule(ctx, child)
		nestedHasOptional := child.HasOptionalOutput()
		if nestedHasOptional {
			step := newOptionalCompositionStep(ctx, cond, nestedRule)
			output = step.combine(output)
			continue
		}
		step := newNonOptionalCompositionStep(ctx, cond, nestedRule)
		output = step.combine(output)
	}

	matchExpr := output.expr()
	identVisitor := opt.rewriteVariableName(ctx)
	ast.PostOrderVisit(matchExpr, identVisitor)

	return matchExpr
}

func (opt *ruleComposerImpl) rewriteVariableName(ctx *cel.OptimizerContext) ast.Visitor {
	return ast.NewExprVisitor(func(expr ast.Expr) {
		if expr.Kind() != ast.IdentKind || !strings.HasPrefix(expr.AsIdent(), "variables.") {
			return
		}
		varName := expr.AsIdent()
		for i := len(opt.varIndices) - 1; i >= 0; i-- {
			v := opt.varIndices[i]
			if v.localVar == varName {
				ctx.UpdateExpr(expr, ctx.NewIdent(v.indexVar))
				return
			}
		}
	})
}

// registerVariable creates an entry for a variable name within the cel.@block used to enumerate
// variables within composed policy expression.
func (opt *ruleComposerImpl) registerVariable(ctx *cel.OptimizerContext, v *CompiledVariable) {
	varName := fmt.Sprintf("variables.%s", v.Name())
	indexVar := fmt.Sprintf("@index%d", opt.nextVarIndex)
	varExpr := ctx.CopyASTAndMetadata(v.Expr().NativeRep())
	ast.PostOrderVisit(varExpr, opt.rewriteVariableName(ctx))
	vi := varIndex{
		index:    opt.nextVarIndex,
		indexVar: indexVar,
		localVar: varName,
		expr:     varExpr,
		celType:  v.Declaration().Type()}
	opt.varIndices = append(opt.varIndices, vi)
	opt.nextVarIndex++
}

type ruleUnnesterImpl struct {
	nextVarIndex     int
	varIndices       []varIndex
	exprUnnestHeight int
}

func (opt *ruleUnnesterImpl) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	// Since the optimizer is based on the original environment provided to the composer,
	// a second pass on the `cel.@block` will require a rebuilding of the cel environment
	ruleExpr := ast.NavigateAST(a)
	var varExprs []ast.Expr
	var varDecls []cel.EnvOption
	if ruleExpr.Kind() == ast.CallKind && ruleExpr.AsCall().FunctionName() == "cel.@block" {
		// Extract the expr from the cel.@block, args[1], as a navigable expr value.
		// Also extract the variable declarations and all associated types from the cel.@block as
		// varIndex values, but without doing any rewrites as the types are all correct already.
		block := ruleExpr.AsCall()
		ruleExpr = block.Args()[1].(ast.NavigableExpr)

		// Collect the list of variables associated with the block
		blockList := block.Args()[0].(ast.NavigableExpr)
		vars := blockList.AsList()
		varExprs = make([]ast.Expr, vars.Size())
		varDecls = make([]cel.EnvOption, vars.Size())
		copy(varExprs, vars.Elements())
		for i, v := range varExprs {
			// Track the variable he varDecls set.
			indexVar := fmt.Sprintf("@index%d", i)
			celType := a.GetType(v.ID())
			varDecls[i] = cel.Variable(indexVar, celType)
			opt.nextVarIndex++
		}
	}
	if len(varDecls) != 0 {
		err := ctx.ExtendEnv(varDecls...)
		if err != nil {
			ctx.ReportErrorAtID(ruleExpr.ID(), "%s", err.Error())
		}
	}

	// Attempt to unnest the rule.
	ruleExpr = opt.maybeUnnestRule(ctx, ruleExpr)
	// If there were no variables, return the expression.
	if len(opt.varIndices) == 0 {
		return a
	}

	// Otherwise populate the cel.@block with the variable declarations and wrap the expression
	// in the block.
	for i := 0; i < len(opt.varIndices); i++ {
		vi := opt.varIndices[i]
		varExprs = append(varExprs, vi.expr)
		err := ctx.ExtendEnv(cel.Variable(vi.indexVar, vi.celType))
		if err != nil {
			ctx.ReportErrorAtID(ruleExpr.ID(), "%s", err.Error())
		}
	}
	blockExpr := ctx.NewCall("cel.@block", ctx.NewList(varExprs, []int32{}), ruleExpr)
	return ctx.NewAST(blockExpr)
}

func (opt *ruleUnnesterImpl) maybeUnnestRule(ctx *cel.OptimizerContext, ruleExpr ast.NavigableExpr) ast.NavigableExpr {
	// Unnest expressions are ordered from leaf to root via the ast.MatchDescendants call.
	heights := ast.Heights(ast.NewAST(ruleExpr, nil))
	unnestMap := map[int64]bool{}
	unnestExprs := []ast.NavigableExpr{}
	ast.MatchDescendants(ruleExpr, func(e ast.NavigableExpr) bool {
		// If the expression is a comprehension, then all unnest candidates captured previously that relate
		// to the comprehension body should be removed from the list of candidate branches for unnesting.
		if e.Kind() == ast.ComprehensionKind {
			// This only removes branches from the map, but not from the list of branches.
			removeIneligibleSubExprs(e, unnestMap)
			return false
		}
		// Otherwise, if the expression is not a call, don't include it.
		if e.Kind() != ast.CallKind {
			return false
		}
		height := heights[e.ID()]
		if height < opt.exprUnnestHeight {
			return false
		}
		unnestMap[e.ID()] = true
		unnestExprs = append(unnestExprs, e)
		return true
	})

	slices.SortStableFunc(unnestExprs, func(a, b ast.NavigableExpr) int {
		heightA := heights[a.ID()]
		heightB := heights[b.ID()]
		return cmp.Compare(heightA, heightB)
	})

	// Prune the expression set to unnest down to only those not included in comprehensions.
	for idx := 0; idx < len(unnestExprs)-1; idx++ {
		e := unnestExprs[idx]
		if present, found := unnestMap[e.ID()]; !found || !present {
			continue
		}
		height := heights[e.ID()]
		if height < opt.exprUnnestHeight {
			continue
		}
		reduceHeight(heights, e, opt.exprUnnestHeight)
		opt.registerUnnestVariable(ctx, e)
	}
	return ruleExpr
}

// registerUnnestVariable creates an entry for a variable name within the cel.@block used to unnest
// a deeply nested logical branch or logical operator.
func (opt *ruleUnnesterImpl) registerUnnestVariable(ctx *cel.OptimizerContext, varExpr ast.NavigableExpr) {
	indexVar := fmt.Sprintf("@index%d", opt.nextVarIndex)
	varExprCopy := ctx.CopyASTAndMetadata(ctx.NewAST(varExpr))
	vi := varIndex{
		index:    opt.nextVarIndex,
		indexVar: indexVar,
		expr:     varExprCopy,
		celType:  varExpr.Type(),
	}
	ctx.UpdateExpr(varExpr, ctx.NewIdent(vi.indexVar))
	opt.varIndices = append(opt.varIndices, vi)
	opt.nextVarIndex++
}

// compositionStep interface represents an intermediate stage of rule and match expression composition
//
// The CompiledRule and CompiledMatch types are meant to represent standalone tuples of condition
// and output expressions, and have no notion of how the order of combination would impact composition
// since composition rules may vary based on the policy execution semantic, e.g. first-match versus
// logical-or, logical-and, or accumulation.
type compositionStep interface {
	// isOptional indicates whether the output step has an optional result.
	//
	// Individual conditional attributes are not optional; however, rules and subrules can have optional output.
	isOptional() bool

	// condition returns the condition associated with the output.
	condition() ast.Expr

	// isConditional returns true if the condition expression is not trivially true.
	isConditional() bool

	// expr returns the output expression for the step.
	expr() ast.Expr

	// combine assembles two output expressions into a single output step.
	combine(other compositionStep) compositionStep
}

// baseCompositionStep encapsulates the common features of an compositionStep implementation.
type baseCompositionStep struct {
	ctx  *cel.OptimizerContext
	cond ast.Expr
	out  ast.Expr
}

func (b baseCompositionStep) condition() ast.Expr {
	return b.cond
}

func (b baseCompositionStep) isConditional() bool {
	c := b.cond
	return c.Kind() != ast.LiteralKind || c.AsLiteral() != types.True
}

func (b baseCompositionStep) expr() ast.Expr {
	return b.out
}

// newNonOptionalCompositionStep returns an output step whose output is not optional.
func newNonOptionalCompositionStep(ctx *cel.OptimizerContext, cond, out ast.Expr) nonOptionalCompositionStep {
	return nonOptionalCompositionStep{
		baseCompositionStep: &baseCompositionStep{
			ctx:  ctx,
			cond: cond,
			out:  out,
		},
	}
}

type nonOptionalCompositionStep struct {
	*baseCompositionStep
}

// isOptional returns false
func (nonOptionalCompositionStep) isOptional() bool {
	return false
}

// combine assembles a new compositionStep from the target output step an an input output step.
//
// non-optional.combine(non-optional) // non-optional
// (non-optional && conditional).combine(optional) // optional
// (non-optional && unconditional).combine(optional) // non-optional
//
// The last combination case is unusual, but effectively it means that the non-optional value prunes away
// the potential optional output.
func (s nonOptionalCompositionStep) combine(step compositionStep) compositionStep {
	if step == nil {
		// The input `step` may be nil if this is the first compositionStep
		return s
	}
	ctx := s.ctx
	trueCondition := ctx.NewLiteral(types.True)
	if step.isOptional() {
		// If the step is optional, convert the non-optional value to an optional one and return a ternary
		if s.isConditional() {
			return newOptionalCompositionStep(ctx,
				trueCondition,
				ctx.NewCall(operators.Conditional,
					s.condition(),
					ctx.NewCall("optional.of", s.expr()),
					step.expr()),
			)
		}
		// The `step` is pruned away by a unconditional non-optional step `s`.
		// Likely a candidate for dead-code warnings.
		return s
	}
	return newNonOptionalCompositionStep(ctx,
		trueCondition,
		ctx.NewCall(operators.Conditional,
			s.condition(),
			s.expr(),
			step.expr()))
}

// newOptionalCompositionStep returns an output step with an optional policy output.
func newOptionalCompositionStep(ctx *cel.OptimizerContext, cond, out ast.Expr) optionalCompositionStep {
	return optionalCompositionStep{
		baseCompositionStep: &baseCompositionStep{
			ctx:  ctx,
			cond: cond,
			out:  out,
		},
	}
}

type optionalCompositionStep struct {
	*baseCompositionStep
}

// isOptional returns true.
func (optionalCompositionStep) isOptional() bool {
	return true
}

// combine assembles a new compositionStep from the target output step an an input output step.
//
// optional.combine(optional) // optional
// (optional && conditional).combine(non-optional) // optional
// (optional && unconditional).combine(non-optional) // non-optional
//
// The last combination case indicates that an optional value in one case should be resolved
// to a non-optional value as
func (s optionalCompositionStep) combine(step compositionStep) compositionStep {
	if step == nil {
		// This is likely unreachable for an optional step, but worth adding as a safeguard
		return s
	}
	ctx := s.ctx
	trueCondition := ctx.NewLiteral(types.True)
	if step.isOptional() {
		// Introduce a ternary to capture the conditional return when combining a
		// conditional optional with another optional.
		if s.isConditional() {
			return newOptionalCompositionStep(ctx,
				trueCondition,
				ctx.NewCall(operators.Conditional,
					s.condition(),
					s.expr(),
					step.expr()),
			)
		}
		// When an optional is unconditionally combined with another optional, rely
		// on the optional 'or' to fall-through from one optional to another.
		if !isOptionalNone(step.expr()) {
			return newOptionalCompositionStep(ctx,
				trueCondition,
				ctx.NewMemberCall("or", s.expr(), step.expr()))
		}
		// Otherwise, the current step 's' is unconditional and effectively prunes away
		// the other input 'step'.
		return s
	}
	if s.isConditional() {
		// Introduce a ternary to capture the conditional return while wrapping the
		// non-optional result from a lower step into an optional value.
		return newOptionalCompositionStep(ctx,
			trueCondition,
			ctx.NewCall(operators.Conditional,
				s.condition(),
				s.expr(),
				ctx.NewCall("optional.of", step.expr())))
	}
	// If the current step is unconditional and the step is non-optional, attempt
	// to convert to the optional step 's' to a non-optional value using `orValue`
	// with the 'step' expression value.
	return newNonOptionalCompositionStep(ctx,
		trueCondition,
		ctx.NewMemberCall("orValue", s.expr(), step.expr()),
	)
}

func isOptionalNone(e ast.Expr) bool {
	return e.Kind() == ast.CallKind &&
		e.AsCall().FunctionName() == "optional.none" &&
		len(e.AsCall().Args()) == 0
}

func removeIneligibleSubExprs(e ast.NavigableExpr, unnestMap map[int64]bool) {
	for _, id := range comprehensionSubExprIDs(e) {
		if _, found := unnestMap[id]; found {
			delete(unnestMap, id)
		}
	}
}

func comprehensionSubExprIDs(e ast.NavigableExpr) []int64 {
	compre := e.AsComprehension()
	// Almost the same as e.Children(), but skips the iteration range
	return enumerateExprIDs(
		compre.AccuInit().(ast.NavigableExpr),
		compre.LoopCondition().(ast.NavigableExpr),
		compre.LoopStep().(ast.NavigableExpr),
		compre.Result().(ast.NavigableExpr),
	)
}

func enumerateExprIDs(exprs ...ast.NavigableExpr) []int64 {
	ids := make([]int64, 0, len(exprs))
	for _, e := range exprs {
		ids = append(ids, e.ID())
		ids = append(ids, enumerateExprIDs(e.Children()...)...)
	}
	return ids
}

func reduceHeight(heights map[int64]int, e ast.NavigableExpr, amount int) {
	height := heights[e.ID()]
	if height < amount {
		return
	}
	heights[e.ID()] = height - amount
	if parent, hasParent := e.Parent(); hasParent {
		reduceHeight(heights, parent, amount)
	}
}
