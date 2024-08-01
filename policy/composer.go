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
	opt := cel.NewStaticOptimizer(&ruleComposerImpl{rule: r})
	return opt.Optimize(c.env, ruleRoot)
}

type ruleComposerImpl struct {
	rule                     *CompiledRule
	maxNestedExpressionLimit int
}

// Optimize implements an AST optimizer for CEL which composes an expression graph into a single
// expression value.
func (opt *ruleComposerImpl) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	// The input to optimize is a dummy expression which is completely replaced according
	// to the configuration of the rule composition graph.
	ruleExpr := opt.optimizeRule(ctx, opt.rule)
	return ctx.NewAST(ruleExpr)
}

func (opt *ruleComposerImpl) optimizeRule(ctx *cel.OptimizerContext, r *CompiledRule) ast.Expr {
	matchExpr := ctx.NewCall("optional.none")
	matches := r.Matches()
	matchCount := len(matches)
	vars := r.Variables()

	optionalResult := true
	// Build the rule subgraph.
	for i := matchCount - 1; i >= 0; i-- {
		m := matches[i]
		cond := ctx.CopyASTAndMetadata(m.Condition().NativeRep())
		// If the condition is trivially true, not of the matches in the rule causes the result
		// to become optional, and the rule is not the last match, then this will introduce
		// unreachable outputs or rules.
		triviallyTrue := m.ConditionIsLiteral(types.True)

		// If the output is non-nil, then determine whether the output should be wrapped
		// into an optional value, a conditional, or both.
		if m.Output() != nil {
			out := ctx.CopyASTAndMetadata(m.Output().Expr().NativeRep())
			if triviallyTrue {
				matchExpr = out
				optionalResult = false
				continue
			}
			if optionalResult {
				out = ctx.NewCall("optional.of", out)
			}
			matchExpr = ctx.NewCall(
				operators.Conditional,
				cond,
				out,
				matchExpr)
			continue
		}

		// If the match has a nested rule, then compute the rule and whether it has
		// an optional return value.
		child := m.NestedRule()
		nestedRule := opt.optimizeRule(ctx, child)
		nestedHasOptional := child.HasOptionalOutput()
		if optionalResult && !nestedHasOptional {
			nestedRule = ctx.NewCall("optional.of", nestedRule)
		}
		if !optionalResult && nestedHasOptional {
			matchExpr = ctx.NewCall("optional.of", matchExpr)
			optionalResult = true
		}
		// If either the nested rule or current condition output are optional then
		// use optional.or() to specify the combination of the first and second results
		// Note, the argument order is reversed due to the traversal of matches in
		// reverse order.
		if optionalResult && triviallyTrue {
			matchExpr = ctx.NewMemberCall("or", nestedRule, matchExpr)
			continue
		}
		matchExpr = ctx.NewCall(
			operators.Conditional,
			cond,
			nestedRule,
			matchExpr,
		)
	}

	// Bind variables in reverse order to declaration on top of rule-subgraph.
	for i := len(vars) - 1; i >= 0; i-- {
		v := vars[i]
		varAST := ctx.CopyASTAndMetadata(v.Expr().NativeRep())
		// Build up the bindings in reverse order, starting from root, all the way up to the outermost
		// binding:
		//    currExpr = cel.bind(outerVar, outerExpr, currExpr)
		varName := v.Declaration().Name()
		inlined, bindMacro := ctx.NewBindMacro(matchExpr.ID(), varName, varAST, matchExpr)
		ctx.UpdateExpr(matchExpr, inlined)
		ctx.SetMacroCall(matchExpr.ID(), bindMacro)
	}
	return matchExpr
}
