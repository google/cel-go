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
func Compile(env *cel.Env, p *policy) (*cel.Ast, *cel.Issues) {
	c := &compiler{
		env:  env,
		info: p.info,
		src:  p.source,
	}
	errs := common.NewErrors(c.src)
	iss := cel.NewIssuesWithSourceInfo(errs, c.info)
	rule, ruleIss := c.compileRule(p.rule, c.env, iss)
	iss = iss.Append(ruleIss)
	if iss.Err() != nil {
		return nil, iss
	}
	ruleRoot, _ := env.Compile("true")
	opt := cel.NewStaticOptimizer(&ruleComposer{rule: rule})
	ruleExprAST, iss := opt.Optimize(env, ruleRoot)
	return ruleExprAST, iss.Append(iss)
}

func (c *compiler) compileRule(r *rule, ruleEnv *cel.Env, iss *cel.Issues) (*compiledRule, *cel.Issues) {
	var err error
	compiledVars := make([]*compiledVariable, len(r.variables))
	for i, v := range r.variables {
		exprSrc := c.relSource(v.expression)
		varAST, exprIss := ruleEnv.CompileSource(exprSrc)
		if exprIss.Err() == nil {
			ruleEnv, err = ruleEnv.Extend(cel.Variable(fmt.Sprintf("variables.%s", v.name.value), varAST.OutputType()))
			if err != nil {
				iss.ReportErrorAtID(v.expression.id, "invalid variable declaration")
			}
			compiledVars[i] = &compiledVariable{
				name: v.name.value,
				expr: varAST,
			}
		}
		iss = iss.Append(exprIss)
	}
	compiledMatches := []*compiledMatch{}
	for _, m := range r.matches {
		condSrc := c.relSource(m.condition)
		condAST, condIss := ruleEnv.CompileSource(condSrc)
		iss = iss.Append(condIss)
		if m.output != nil && m.rule != nil {
			iss.ReportErrorAtID(m.condition.id, "either output or rule may be set but not both")
			continue
		}
		if m.output != nil {
			outSrc := c.relSource(*m.output)
			outAST, outIss := ruleEnv.CompileSource(outSrc)
			iss = iss.Append(outIss)
			compiledMatches = append(compiledMatches, &compiledMatch{
				cond:   condAST,
				output: outAST,
			})
			continue
		}
		if m.rule != nil {
			nestedRule, ruleIss := c.compileRule(m.rule, ruleEnv, iss)
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

func (c *compiler) relSource(pstr policyString) *RelativeSource {
	line := 0
	col := 1
	if offset, found := c.info.GetOffsetRange(pstr.id); found {
		if loc, found := c.src.OffsetLocation(offset.Start); found {
			line = loc.Line()
			col = loc.Column()
		}
	}
	return c.src.Relative(pstr.value, line, col)
}

type ruleComposer struct {
	rule *compiledRule
}

func (opt *ruleComposer) Optimize(ctx *cel.OptimizerContext, a *ast.AST) *ast.AST {
	ruleExpr, _ := optimizeRule(ctx, opt.rule)
	ctx.UpdateExpr(a.Expr(), ruleExpr)
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
		optionalResult = optionalResult || nestedOptional
		if triviallyTrue {
			matchExpr = nestedRule
			continue
		}
		matchExpr = ctx.NewCall(
			operators.Conditional,
			cond,
			nestedRule,
			matchExpr)
	}

	vars := r.variables
	for i := len(vars) - 1; i >= 0; i-- {
		v := vars[i]
		varAST := ctx.CopyASTAndMetadata(v.expr.NativeRep())
		// Build up the bindings in reverse order, starting from root, all the way up to the outermost
		// binding:
		//    currExpr = cel.bind(outerVar, outerExpr, currExpr)
		inlined, bindMacro := ctx.NewBindMacro(matchExpr.ID(), fmt.Sprintf("variables.%s", v.name), varAST, matchExpr)
		ctx.SetMacroCall(inlined.ID(), bindMacro)
		matchExpr = inlined
	}
	return matchExpr, optionalResult
}