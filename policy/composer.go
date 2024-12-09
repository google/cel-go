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
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

// NewRuleComposer creates a rule composer which stitches together rules within a policy into
// a single CEL expression.
func NewRuleComposer(env *cel.Env, p *Policy) *RuleComposer {
	return &RuleComposer{
		env: env,
		p:   p,
	}
}

// RuleComposer optimizes a set of expressions into a single expression.
type RuleComposer struct {
	env *cel.Env
	p   *Policy
}

// Compose stitches together a set of expressions within a CompiledRule into a single CEL ast.
func (c *RuleComposer) Compose(r *CompiledRule) (*cel.Ast, *cel.Issues) {
	ruleRoot, _ := c.env.Compile("true")
	opt := cel.NewStaticOptimizer(&ruleComposerImpl{rule: r, varIndices: []varIndex{}})
	return opt.Optimize(c.env, ruleRoot)
}

type varIndex struct {
	index    int
	indexVar string
	localVar string
	expr     ast.Expr
	cv       *CompiledVariable
}

type ruleComposerImpl struct {
	rule         *CompiledRule
	nextVarIndex int
	varIndices   []varIndex

	maxNestedExpressionLimit int
}

// Optimize implements an AST optimizer for CEL which composes an expression graph into a single
// expression value.
func (opt *ruleComposerImpl) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	// The input to optimize is a dummy expression which is completely replaced according
	// to the configuration of the rule composition graph.
	ruleExpr := opt.optimizeRule(ctx, opt.rule)
	allVars := opt.sortedVariables()
	// If there were no variables, return the expression.
	if len(allVars) == 0 {
		return ctx.NewAST(ruleExpr)
	}

	// Otherwise populate the block.
	varExprs := make([]ast.Expr, len(allVars))
	for i, vi := range allVars {
		varExprs[i] = vi.expr
		err := ctx.ExtendEnv(cel.Variable(vi.indexVar, vi.cv.Declaration().Type()))
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
	var output outputStep = nil
	if r.HasOptionalOutput() {
		output = newOptionalOutputStep(ctx, ctx.NewLiteral(types.True), ctx.NewCall("optional.none"))
	}
	// Build the rule subgraph.
	for i := matchCount - 1; i >= 0; i-- {
		m := matches[i]
		cond := ctx.CopyASTAndMetadata(m.Condition().NativeRep())

		// If the output is non-nil, then determine whether the output should be wrapped
		// into an optional value, a conditional, or both.
		if m.Output() != nil {
			out := ctx.CopyASTAndMetadata(m.Output().Expr().NativeRep())
			step := newNonOptionalOutputStep(ctx, cond, out)
			output = step.combine(output)
			continue
		}

		// If the match has a nested rule, then compute the rule and whether it has
		// an optional return value.
		child := m.NestedRule()
		nestedRule := opt.optimizeRule(ctx, child)
		nestedHasOptional := child.HasOptionalOutput()
		if nestedHasOptional {
			step := newOptionalOutputStep(ctx, cond, nestedRule)
			output = step.combine(output)
			continue
		}
		step := newNonOptionalOutputStep(ctx, cond, nestedRule)
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
		cv:       v}
	opt.varIndices = append(opt.varIndices, vi)
	opt.nextVarIndex++
}

// sortedVariables returns the variables ordered by their declaration index.
func (opt *ruleComposerImpl) sortedVariables() []varIndex {
	return opt.varIndices
}

// outputStep interface represents a policy output expression.
type outputStep interface {
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
	combine(other outputStep) outputStep
}

// baseOutputStep encapsulates the common features of an outputStep implementation.
type baseOutputStep struct {
	ctx  *cel.OptimizerContext
	cond ast.Expr
	out  ast.Expr
}

func (b baseOutputStep) condition() ast.Expr {
	return b.cond
}

func (b baseOutputStep) isConditional() bool {
	c := b.cond
	return c.Kind() != ast.LiteralKind || c.AsLiteral() != types.True
}

func (b baseOutputStep) expr() ast.Expr {
	return b.out
}

// newNonOptionalOutputStep returns an output step whose output is not optional.
func newNonOptionalOutputStep(ctx *cel.OptimizerContext, cond, out ast.Expr) nonOptionalOutputStep {
	return nonOptionalOutputStep{
		baseOutputStep: &baseOutputStep{
			ctx:  ctx,
			cond: cond,
			out:  out,
		},
	}
}

type nonOptionalOutputStep struct {
	*baseOutputStep
}

// isOptional returns false
func (nonOptionalOutputStep) isOptional() bool {
	return false
}

// combine assembles a new outputStep from the target output step an an input output step.
//
// non-optional.combine(non-optional) // non-optional
// (non-optional && conditional).combine(optional) // optional
// (non-optional && unconditional).combine(optional) // non-optional
//
// The last combination case is unusual, but effectively it means that the non-optional value prunes away
// the potential optional output.
func (s nonOptionalOutputStep) combine(step outputStep) outputStep {
	if step == nil {
		// The input `step`` may be nil if this is the first outputStep
		return s
	}
	ctx := s.ctx
	trueCondition := ctx.NewLiteral(types.True)
	if step.isOptional() {
		// If the step is optional, convert the non-optional value to an optional one and return a ternary
		if s.isConditional() {
			return newOptionalOutputStep(ctx,
				trueCondition,
				ctx.NewCall(operators.Conditional,
					s.condition(),
					ctx.NewCall("optional.of", s.expr()),
					step.expr()),
			)
		}
		// The `step` is pruned away by a unconditional non-optional step `s`.
		return s
	}
	return newNonOptionalOutputStep(ctx,
		trueCondition,
		ctx.NewCall(operators.Conditional,
			s.condition(),
			s.expr(),
			step.expr()))
}

// newOptionalOutputStep returns an output step with an optional policy output.
func newOptionalOutputStep(ctx *cel.OptimizerContext, cond, out ast.Expr) optionalOutputStep {
	return optionalOutputStep{
		baseOutputStep: &baseOutputStep{
			ctx:  ctx,
			cond: cond,
			out:  out,
		},
	}
}

type optionalOutputStep struct {
	*baseOutputStep
}

// isOptional returns true.
func (optionalOutputStep) isOptional() bool {
	return true
}

// combine assembles a new outputStep from the target output step an an input output step.
//
// optional.combine(optional) // optional
// (optional && conditional).combine(non-optional) // optional
// (optional && unconditional).combine(non-optional) // non-optional
//
// The last combination case indicates that an optional value in one case should be resolved
// to a non-optional value as
func (s optionalOutputStep) combine(step outputStep) outputStep {
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
			return newOptionalOutputStep(ctx,
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
			return newOptionalOutputStep(ctx,
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
		return newOptionalOutputStep(ctx,
			trueCondition,
			ctx.NewCall(operators.Conditional,
				s.condition(),
				s.expr(),
				ctx.NewCall("optional.of", step.expr())))
	}
	// If the current step is unconditional and the step is non-optional, attempt
	// to convert to the optional step 's' to a non-optional value using `orValue`
	// with the 'step' expression value.
	return newNonOptionalOutputStep(ctx,
		trueCondition,
		ctx.NewMemberCall("orValue", s.expr(), step.expr()),
	)
}

func isOptionalNone(e ast.Expr) bool {
	return e.Kind() == ast.CallKind &&
		e.AsCall().FunctionName() == "optional.none" &&
		len(e.AsCall().Args()) == 0
}
