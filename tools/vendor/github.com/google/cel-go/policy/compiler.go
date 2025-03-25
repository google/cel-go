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
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// CompiledRule represents the variables and match blocks associated with a rule block.
type CompiledRule struct {
	exprID    int64
	id        *ValueString
	variables []*CompiledVariable
	matches   []*CompiledMatch
}

// SourceID returns the source metadata identifier associated with the compiled rule.
func (r *CompiledRule) SourceID() int64 {
	return r.exprID
}

// ID returns the expression id associated with the rule.
func (r *CompiledRule) ID() *ValueString {
	return r.id
}

// Variables rturns the list of CompiledVariable values associated with the rule.
func (r *CompiledRule) Variables() []*CompiledVariable {
	return r.variables[:]
}

// Matches returns the list of matches associated with the rule.
func (r *CompiledRule) Matches() []*CompiledMatch {
	return r.matches[:]
}

// OutputType returns the output type of the first match clause as all match clauses
// are validated for agreement prior to construction fo the CompiledRule.
func (r *CompiledRule) OutputType() *cel.Type {
	// It's a compilation error if the output types of the matches don't agree
	for _, m := range r.Matches() {
		return m.OutputType()
	}
	return cel.DynType
}

// HasOptionalOutput returns whether the rule returns a concrete or optional value.
// The rule may return an optional value if all match expressions under the rule are conditional.
func (r *CompiledRule) HasOptionalOutput() bool {
	optionalOutput := false
	for _, m := range r.Matches() {
		if m.NestedRule() != nil && m.NestedRule().HasOptionalOutput() {
			return true
		}
		if m.ConditionIsLiteral(types.True) {
			return false
		}
		optionalOutput = true
	}
	return optionalOutput
}

// CompiledVariable represents the variable name, expression, and associated type-check declaration.
type CompiledVariable struct {
	exprID  int64
	name    string
	expr    *cel.Ast
	varDecl *decls.VariableDecl
}

// SourceID returns the source metadata identifier associated with the variable.
func (v *CompiledVariable) SourceID() int64 {
	return v.exprID
}

// Name returns the variable name.
func (v *CompiledVariable) Name() string {
	return v.name
}

// Expr returns the compiled expression associated with the variable name.
func (v *CompiledVariable) Expr() *cel.Ast {
	return v.expr
}

// Declaration returns the type-check declaration associated with the variable.
func (v *CompiledVariable) Declaration() *decls.VariableDecl {
	return v.varDecl
}

// CompiledMatch represents a match block which has an optional condition (true, by default) as well
// as an output or a nested rule (one or the other, but not both).
type CompiledMatch struct {
	exprID     int64
	cond       *cel.Ast
	output     *OutputValue
	nestedRule *CompiledRule
}

// SourceID returns the source identifier associated with the compiled match.
func (m *CompiledMatch) SourceID() int64 {
	return m.exprID
}

// Condition returns the compiled predicate expression which must evaluate to true before the output
// or subrule is entered.
func (m *CompiledMatch) Condition() *cel.Ast {
	return m.cond
}

// ConditionIsLiteral indicates whether the condition for the match is a literal with a given value.
func (m *CompiledMatch) ConditionIsLiteral(val ref.Val) bool {
	c := m.cond.NativeRep().Expr()
	return c.Kind() == ast.LiteralKind && c.AsLiteral().Equal(val) == types.True
}

// Output returns the compiled output expression associated with the match block, if set.
func (m *CompiledMatch) Output() *OutputValue {
	return m.output
}

// NestedRule returns the nested rule, if set.
func (m *CompiledMatch) NestedRule() *CompiledRule {
	return m.nestedRule
}

// OutputType returns the cel.Type associated with output expression.
func (m *CompiledMatch) OutputType() *cel.Type {
	if m.output != nil {
		return m.output.Expr().OutputType()
	}
	if m.nestedRule != nil {
		return m.nestedRule.OutputType()
	}
	return cel.DynType
}

// OutputValue represents the output expression associated with a match block.
type OutputValue struct {
	exprID int64
	expr   *cel.Ast
}

// SourceID returns the expression id associated with the output expression.
func (o *OutputValue) SourceID() int64 {
	return o.exprID
}

// Expr returns the compiled expression associated with the output.
func (o *OutputValue) Expr() *cel.Ast {
	return o.expr
}

// CompilerOption specifies a functional option to be applied to new RuleComposer instances.
type CompilerOption func(*compiler) error

// MaxNestedExpressions limits the number of variable and nested rule expressions during compilation.
//
// Defaults to 100 if not set.
func MaxNestedExpressions(limit int) CompilerOption {
	return func(c *compiler) error {
		if limit <= 0 {
			return fmt.Errorf("nested expression limit must be non-negative, non-zero value: %d", limit)
		}
		c.maxNestedExpressions = limit
		return nil
	}
}

// Compile combines the policy compilation and composition steps into a single call.
//
// This generates a single CEL AST from a collection of policy expressions associated with a
// CEL environment.
func Compile(env *cel.Env, p *Policy, opts ...CompilerOption) (*cel.Ast, *cel.Issues) {
	rule, iss := CompileRule(env, p, opts...)
	if iss.Err() != nil {
		return nil, iss
	}
	// An error cannot happen when composing without supplying options
	composer, _ := NewRuleComposer(env)
	return composer.Compose(rule)
}

// CompileRule creates a compiled rules from the policy which contains a set of compiled variables and
// match statements. The compiled rule defines an expression graph, which can be composed into a single
// expression via the RuleComposer.Compose method.
func CompileRule(env *cel.Env, p *Policy, opts ...CompilerOption) (*CompiledRule, *cel.Issues) {
	c := &compiler{
		env:                  env,
		info:                 p.SourceInfo(),
		src:                  p.Source(),
		maxNestedExpressions: defaultMaxNestedExpressions,
	}
	var err error
	errs := common.NewErrors(c.src)
	iss := cel.NewIssuesWithSourceInfo(errs, c.info)
	for _, o := range opts {
		if err = o(c); err != nil {
			iss.ReportErrorAtID(p.Name().ID, "error configuring compiler option: %s", err)
			return nil, iss
		}
	}
	c.env, err = c.env.Extend(cel.EagerlyValidateDeclarations(true))
	if err != nil {
		iss.ReportErrorAtID(p.Name().ID, "error configuring environment: %s", err)
		return nil, iss
	}

	importCount := len(p.Imports())
	if importCount > 0 {
		importNames := make([]string, 0, importCount)
		for _, imp := range p.Imports() {
			typeName := imp.Name().Value
			_, err := containers.NewContainer(containers.Abbrevs(typeName))
			if err != nil {
				iss.ReportErrorAtID(imp.Name().ID, "error configuring import: %s", err)
			} else {
				importNames = append(importNames, typeName)
			}
		}
		env, err := c.env.Extend(cel.Abbrevs(importNames...))
		if err != nil {
			// validation happens earlier in the sequence, so this should be unreachable.
			iss.ReportErrorAtID(p.Imports()[0].SourceID(), "error configuring imports: %s", err)
		} else {
			c.env = env
		}
	}
	return c.compileRule(p.Rule(), c.env, iss)
}

type compiler struct {
	env  *cel.Env
	info *ast.SourceInfo
	src  *Source

	maxNestedExpressions int
	nestedCount          int
}

func (c *compiler) compileRule(r *Rule, ruleEnv *cel.Env, iss *cel.Issues) (*CompiledRule, *cel.Issues) {
	compiledVars := make([]*CompiledVariable, len(r.Variables()))
	for i, v := range r.Variables() {
		exprSrc := c.relSource(v.Expression())
		varAST, exprIss := ruleEnv.CompileSource(exprSrc)
		varName := v.Name().Value

		// Determine the variable type. If the expression is an error then record the error and
		// mark the variable type as dyn to permit compilation to continue.
		varType := types.DynType
		if exprIss.Err() != nil {
			iss = iss.Append(exprIss)
		} else {
			// Otherwise, the expression compiled successfully and we use its output type.
			varType = varAST.OutputType()
		}

		// Introduce the variable into the environment. By extending the environment, the variables
		// are effectively scoped such that they must be declared before use.
		varDecl := decls.NewVariable(fmt.Sprintf("%s.%s", variablePrefix, varName), varType)
		varEnv, err := ruleEnv.Extend(cel.Variable(varDecl.Name(), varDecl.Type()))
		if err != nil {
			iss.ReportErrorAtID(v.exprID, "invalid variable declaration: %s", err.Error())
		} else {
			ruleEnv = varEnv
		}
		compiledVar := &CompiledVariable{
			exprID:  v.name.ID,
			name:    v.name.Value,
			expr:    varAST,
			varDecl: varDecl,
		}
		compiledVars[i] = compiledVar

		// Increment the nesting count post-compile.
		c.nestedCount++
		if c.nestedCount == c.maxNestedExpressions+1 {
			iss.ReportErrorAtID(compiledVar.SourceID(), "variable exceeds nested expression limit")
		}
	}

	// Compile the set of match conditions under the rule.
	compiledMatches := []*CompiledMatch{}
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
			compiledMatches = append(compiledMatches, &CompiledMatch{
				exprID: m.exprID,
				cond:   condAST,
				output: &OutputValue{
					exprID: m.Output().ID,
					expr:   outAST,
				},
			})
			continue
		}
		if m.HasRule() {
			nestedRule, ruleIss := c.compileRule(m.Rule(), ruleEnv, iss)
			iss = iss.Append(ruleIss)
			compiledMatches = append(compiledMatches, &CompiledMatch{
				exprID:     m.exprID,
				cond:       condAST,
				nestedRule: nestedRule,
			})

			// Increment the nesting count post-compile.
			c.nestedCount++
			if c.nestedCount == c.maxNestedExpressions+1 {
				iss.ReportErrorAtID(nestedRule.SourceID(), "rule exceeds nested expression limit")
			}
		}
	}

	// Validate that all branches in the rule are reachable
	rule := &CompiledRule{
		exprID:    r.exprID,
		id:        r.id,
		variables: compiledVars,
		matches:   compiledMatches,
	}

	// Note: Consider supporting configurable policy validators that take the policy, rule, and issues
	// Validate type agreement between the different match outputs
	c.checkMatchOutputTypesAgree(rule, iss)
	// Validate that all branches in the policy are reachable
	c.checkUnreachableCode(rule, iss)

	return rule, iss
}

func (c *compiler) checkMatchOutputTypesAgree(rule *CompiledRule, iss *cel.Issues) {
	var outputType *cel.Type
	for _, m := range rule.Matches() {
		if outputType == nil {
			outputType = m.OutputType()
			if outputType.TypeName() == "error" {
				outputType = nil
				continue
			}
		}
		matchOutputType := m.OutputType()
		if matchOutputType.TypeName() == "error" {
			continue
		}
		// Handle assignability as the output type is assignable to the match output or vice versa.
		// During composition, this is roughly how the type-checker will handle the type agreement check.
		if !(outputType.IsAssignableType(matchOutputType) || matchOutputType.IsAssignableType(outputType)) {
			sourceID := m.SourceID()
			if m.Output() != nil {
				sourceID = m.Output().SourceID()
			} else if m.NestedRule() != nil {
				sourceID = m.NestedRule().SourceID()
			}
			iss.ReportErrorAtID(sourceID, "incompatible output types: block has output type %s, but previous outputs have type %s", matchOutputType, outputType)
			return
		}
	}
}

func (c *compiler) checkUnreachableCode(rule *CompiledRule, iss *cel.Issues) {
	ruleHasOptional := rule.HasOptionalOutput()
	compiledMatches := rule.Matches()
	matchCount := len(compiledMatches)
	for i := matchCount - 1; i >= 0; i-- {
		m := compiledMatches[i]
		triviallyTrue := m.ConditionIsLiteral(types.True)

		if triviallyTrue && !ruleHasOptional && i != matchCount-1 {
			if m.Output() != nil {
				iss.ReportErrorAtID(m.SourceID(), "match creates unreachable outputs")
			}
			if m.NestedRule() != nil {
				iss.ReportErrorAtID(m.NestedRule().SourceID(), "rule creates unreachable outputs")
			}
			break
		}
	}
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

const (
	// Consider making the variables namespace configurable.
	variablePrefix = "variables"

	defaultMaxNestedExpressions = 100
)
