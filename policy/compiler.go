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

// Package policy provides an extensible parser and compiler for composing
// a graph of CEL expressions into a single evaluable expression.
package policy

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

type compiler struct {
	env  *cel.Env
	info *ast.SourceInfo
	src  *Source
}

type compiledRule struct {
	variables []*compiledVariable
	matches   []*compiledMatch
}

type compiledVariable struct {
	name string
	expr *cel.Ast
}

type compiledMatch struct {
	cond       *cel.Ast
	output     *cel.Ast
	nestedRule *compiledRule
}

// Compile generates a single CEL AST from a collection of policy expressions associated with a CEL environment.
func Compile(env *cel.Env, p *Policy) (*cel.Ast, *cel.Issues) {
	c := &compiler{
		env:  env,
		info: p.SourceInfo(),
		src:  p.Source(),
	}
	errs := common.NewErrors(c.src)
	iss := cel.NewIssuesWithSourceInfo(errs, c.info)
	rule, ruleIss := c.compileRule(p.Rule(), c.env, iss)
	iss = iss.Append(ruleIss)
	if iss.Err() != nil {
		return nil, iss
	}
	ruleRoot, _ := env.Compile("true")
	opt := cel.NewStaticOptimizer(&ruleComposer{rule: rule})
	ruleExprAST, optIss := opt.Optimize(env, ruleRoot)
	return ruleExprAST, iss.Append(optIss)
}

func (c *compiler) compileRule(r *Rule, ruleEnv *cel.Env, iss *cel.Issues) (*compiledRule, *cel.Issues) {
	var err error
	compiledVars := make([]*compiledVariable, len(r.Variables()))
	for i, v := range r.Variables() {
		exprSrc := c.relSource(v.Expression())
		varAST, exprIss := ruleEnv.CompileSource(exprSrc)
		if exprIss.Err() == nil {
			ruleEnv, err = ruleEnv.Extend(cel.Variable(fmt.Sprintf("%s.%s", variablePrefix, v.Name().Value), varAST.OutputType()))
			if err != nil {
				iss.ReportErrorAtID(v.Expression().ID, "invalid variable declaration")
			}
			compiledVars[i] = &compiledVariable{
				name: v.name.Value,
				expr: varAST,
			}
		}
		iss = iss.Append(exprIss)
	}
	compiledMatches := []*compiledMatch{}
	for _, m := range r.Matches() {
		condSrc := c.relSource(m.Condition())
		condAST, condIss := ruleEnv.CompileSource(condSrc)
		iss = iss.Append(condIss)
		// This case cannot happen when the Policy object is parsed from yaml, but could happen
		// with a non-YAML generation of the Policy object.
		// TODO: Test this case once there's an alternative method of constructing Policy objects
		if m.HasOutput() && m.HasRule() {
			iss.ReportErrorAtID(m.Condition().ID, "either output or rule may be set but not both")
			continue
		}
		if m.HasOutput() {
			outSrc := c.relSource(m.Output())
			outAST, outIss := ruleEnv.CompileSource(outSrc)
			iss = iss.Append(outIss)
			compiledMatches = append(compiledMatches, &compiledMatch{
				cond:   condAST,
				output: outAST,
			})
			continue
		}
		if m.HasRule() {
			nestedRule, ruleIss := c.compileRule(m.Rule(), ruleEnv, iss)
			iss = iss.Append(ruleIss)
			compiledMatches = append(compiledMatches, &compiledMatch{
				cond:       condAST,
				nestedRule: nestedRule,
			})
		}
	}
	return &compiledRule{
		variables: compiledVars,
		matches:   compiledMatches,
	}, iss
}

func (c *compiler) relSource(pstr ValueString) *RelativeSource {
	line := 0
	col := 1
	if offset, found := c.info.GetOffsetRange(pstr.ID); found {
		if loc, found := c.src.OffsetLocation(offset.Start); found {
			line = loc.Line()
			col = loc.Column()
		}
	}
	return c.src.Relative(pstr.Value, line, col)
}

type ruleComposer struct {
	rule *compiledRule
}

// Optimize implements an AST optimizer for CEL which composes an expression graph into a single
// expression value.
func (opt *ruleComposer) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	// The input to optimize is a dummy expression which is completely replaced according
	// to the configuration of the rule composition graph.
	ruleExpr, _ := optimizeRule(ctx, opt.rule)
	return ctx.NewAST(ruleExpr)
}

func optimizeRule(ctx *cel.OptimizerContext, r *compiledRule) (ast.Expr, bool) {
	matchExpr := ctx.NewCall("optional.none")
	matches := r.matches
	optionalResult := true
	for i := len(matches) - 1; i >= 0; i-- {
		m := matches[i]
		cond := ctx.CopyASTAndMetadata(m.cond.NativeRep())
		triviallyTrue := cond.Kind() == ast.LiteralKind && cond.AsLiteral() == types.True
		if m.output != nil {
			out := ctx.CopyASTAndMetadata(m.output.NativeRep())
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
		nestedRule, nestedOptional := optimizeRule(ctx, m.nestedRule)
		if optionalResult && !nestedOptional {
			nestedRule = ctx.NewCall("optional.of", nestedRule)
		}
		if !optionalResult && nestedOptional {
			matchExpr = ctx.NewCall("optional.of", matchExpr)
			optionalResult = true
		}
		if !optionalResult && !nestedOptional {
			ctx.ReportErrorAtID(nestedRule.ID(), "subrule early terminates policy")
			continue
		}
		matchExpr = ctx.NewMemberCall("or", nestedRule, matchExpr)
	}

	vars := r.variables
	for i := len(vars) - 1; i >= 0; i-- {
		v := vars[i]
		varAST := ctx.CopyASTAndMetadata(v.expr.NativeRep())
		// Build up the bindings in reverse order, starting from root, all the way up to the outermost
		// binding:
		//    currExpr = cel.bind(outerVar, outerExpr, currExpr)
		varName := fmt.Sprintf("%s.%s", variablePrefix, v.name)
		inlined, bindMacro := ctx.NewBindMacro(matchExpr.ID(), varName, varAST, matchExpr)
		ctx.UpdateExpr(matchExpr, inlined)
		ctx.SetMacroCall(matchExpr.ID(), bindMacro)
	}
	return matchExpr, optionalResult
}

const (
	// Consider making the variables namespace configurable.
	variablePrefix = "variables"
)
