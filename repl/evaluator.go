// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Provides Evaluator type with basic operations for setting up the evaluation environment then
// evaluating small CEL expressions.
package main

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// LetVariable let variable representation
type LetVariable struct {
	identifier string
	src        string

	// memoized results from building the expression tree
	typeHint *expr.Type
	env      *cel.Env
	ast      *cel.Ast
	prog     *cel.Program
}

func (l LetVariable) String() string {
	return fmt.Sprintf("%s = %s", l.identifier, l.src)
}

// Reset plan if we need to recompile based on a dependency change.
func (l *LetVariable) clearPlan() {
	l.typeHint = nil
	l.env = nil
	l.ast = nil
	l.prog = nil
}

// EvaluationContext context for the repl.
// Handles maintaining state for multiple let expressions.
type EvaluationContext struct {
	letVars []LetVariable
}

func (ctx *EvaluationContext) indexLetVar(name string) int {
	for idx, el := range ctx.letVars {
		if el.identifier == name {
			return idx
		}
	}
	return -1
}

func (ctx *EvaluationContext) copy() *EvaluationContext {
	var cpy EvaluationContext
	cpy.letVars = make([]LetVariable, len(ctx.letVars))
	copy(cpy.letVars, ctx.letVars)
	return &cpy
}

func (ctx *EvaluationContext) delLetVar(name string) {
	idx := ctx.indexLetVar(name)
	if idx < 0 {
		// no-op if deleting something that's not defined
		return
	}

	ctx.letVars = append(ctx.letVars[:idx], ctx.letVars[idx+1:]...)

	for i := idx; i < len(ctx.letVars); i++ {
		ctx.letVars[i].clearPlan()
	}

}

// Add or update an existing let then invalidate any computed plans.
func (ctx *EvaluationContext) addLetVar(name string, expr string) {
	idx := ctx.indexLetVar(name)
	if idx < 0 {
		ctx.letVars = append(ctx.letVars, LetVariable{identifier: name, src: expr})
	} else {
		ctx.letVars[idx] = LetVariable{identifier: name, src: expr}
		for i := idx + 1; i < len(ctx.letVars); i++ {
			// invalidate dependant let exprs
			ctx.letVars[i].clearPlan()
		}
	}
}

// Evaluator provides basic environment for evaluating an expression with applied context.
type Evaluator struct {
	env *cel.Env
	ctx EvaluationContext
}

// NewEvaluator returns an inialized evaluator
func NewEvaluator() (*Evaluator, error) {
	env, err := cel.NewEnv()
	if err != nil {
		return nil, err
	}

	return &Evaluator{env: env}, nil
}

// Attempt to update context in place after an update. return an error if this fails.
func (e *Evaluator) updateContextPlans(ctx *EvaluationContext) error {
	env := e.env
	for i := range ctx.letVars {
		el := &ctx.letVars[i]
		if el.prog == nil {
			ast, iss := env.Compile(el.src)
			if iss != nil {
				return fmt.Errorf("Error updating %v\n%w", el, iss.Err())
			}
			el.ast = ast
			el.typeHint = ast.ResultType()

			plan, err := env.Program(ast)
			if err != nil {
				return err
			}
			el.prog = &plan
		}
		if el.env == nil {
			env, err := env.Extend(cel.Declarations(decls.NewVar(el.identifier, el.typeHint)))
			if err != nil {
				return err
			}
			el.env = env
		}
		env = el.env
	}
	return nil

}

// AddLetVar adds a let variable to the evaluation context.
// The expression is planned but evaluated lazily.
func (e *Evaluator) AddLetVar(name string, expr string) error {
	ctx := e.ctx.copy()
	ctx.addLetVar(name, expr)
	err := e.updateContextPlans(ctx)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

func (e *Evaluator) DelLetVar(name string) error {
	ctx := e.ctx.copy()
	ctx.delLetVar(name)
	err := e.updateContextPlans(ctx)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

func (e *Evaluator) applyContext() (*cel.Env, interpreter.Activation, error) {
	env := e.env
	var vars = make(map[string]interface{})
	for _, el := range e.ctx.letVars {
		val, _, err := (*el.prog).Eval(vars)
		if val != nil {
			vars[el.identifier] = val
		} else if err != nil {
			return nil, nil, err
		}

		env = el.env
	}

	act, err := interpreter.NewActivation(vars)
	if err != nil {
		return nil, nil, err
	}

	return env, act, nil
}

// Evaluate sets up a CEL evaluation using the provided evaluation context
func (e *Evaluator) Evaluate(expr string) (ref.Val, *expr.Type, error) {
	env, act, err := e.applyContext()
	if err != nil {
		return nil, nil, err
	}

	ast, iss := env.Compile(expr)
	if iss != nil {
		return nil, nil, iss.Err()
	}

	p, err := e.env.Program(ast)
	if err != nil {
		return nil, nil, err
	}

	val, _, err := p.Eval(act)
	// expression can be well-formed and result in an error
	return val, ast.ResultType(), err
}
