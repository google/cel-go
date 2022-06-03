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
// Provides Evaluator type with basic operations for setting up the evaluation
// environment then evaluating small CEL expressions.
package main

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/interpreter/functions"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"google.golang.org/protobuf/proto"
)

// letVariable let variable representation
type letVariable struct {
	identifier string
	src        string
	typeHint   *exprpb.Type

	// memoized results from building the expression AST and program
	resultType *exprpb.Type
	env        *cel.Env
	ast        *cel.Ast
	prog       cel.Program
}

type letFunctionParam struct {
	identifier string
	typeHint   *exprpb.Type
}

// letFunction coordinates let function data (type definition and CEL function implementation).
type letFunction struct {
	identifier string
	src        string
	resultType *exprpb.Type
	params     []letFunctionParam
	receiver   *exprpb.Type // if not nil indicates an instance function

	// memoized results from building the expression AST and program
	env   *cel.Env // the context env for repl evaluation
	fnEnv *cel.Env // the env for implementing the extension fn
	prog  cel.Program
	impl  functions.FunctionOp
}

func typeAssignable(rtType ref.Type, declType *exprpb.Type) bool {
	// TODO(issue/535): add better type agreement support
	return UnparseType(declType) == rtType.TypeName()
}

func checkArgsMatch(params []letFunctionParam, args []ref.Val) error {
	if len(params) != len(args) {
		return fmt.Errorf("got %d args, expected %d", len(args), len(params))
	}
	for i, arg := range args {
		if !typeAssignable(arg.Type(), params[i].typeHint) {
			return fmt.Errorf("got %s, expected %s for argument %d", arg.Type().TypeName(), UnparseType(params[i].typeHint), i)
		}
	}
	return nil
}

func (l *letFunction) updateImpl(env *cel.Env, deps []*functions.Overload) error {
	var paramVars []*exprpb.Decl

	if l.receiver != nil {
		paramVars = append(paramVars, decls.NewVar("this", l.receiver))
	}

	for _, p := range l.params {
		paramVars = append(paramVars, decls.NewVar(p.identifier, p.typeHint))
	}

	var err error
	l.fnEnv, err = env.Extend(cel.Declarations(paramVars...))
	if err != nil {
		return err
	}

	ast, iss := l.fnEnv.Compile(l.src)

	if iss != nil {
		return iss.Err()
	}

	if !proto.Equal(ast.ResultType(), l.resultType) {
		return fmt.Errorf("got result type %s for %s", UnparseType(ast.ResultType()), l)
	}

	l.prog, err = l.fnEnv.Program(ast, cel.Functions(deps...))

	if err != nil {
		return err
	}

	l.impl = func(args ...ref.Val) ref.Val {
		var err error
		var instance ref.Val
		if l.receiver != nil {
			instance = args[0]
			args = args[1:]
		}
		err = checkArgsMatch(l.params, args)
		if err != nil {
			return types.NewErr("error evaluating %s: %v", l, err)
		}

		activation := make(map[string]interface{})
		for i, param := range l.params {
			activation[param.identifier] = args[i]
		}

		if instance != nil {
			if !typeAssignable(instance.Type(), l.receiver) {
				return types.NewErr("error evaluating %s: got receiver type: %s wanted %s", l, instance.Type().TypeName(), UnparseType(l.receiver))
			}
			activation["this"] = instance
		}

		val, _, err := l.prog.Eval(activation)

		if err != nil {
			return types.NewErr("error evaluating %s: %v", l, err)
		}

		return val
	}
	return nil
}

func (l *letFunction) update(env *cel.Env, deps []*functions.Overload) error {
	var err error
	if l.src != "" {
		err = l.updateImpl(env, deps)
		if err != nil {
			return err
		}
	}

	paramTypes := make([]*exprpb.Type, len(l.params))
	for i, p := range l.params {
		paramTypes[i] = p.typeHint
	}

	var opt cel.EnvOption
	if l.receiver != nil {
		paramTypes = append([]*exprpb.Type{l.receiver}, paramTypes...)
		opt = cel.Declarations(
			decls.NewFunction(l.identifier,
				decls.NewInstanceOverload(
					l.identifier,
					paramTypes,
					l.resultType,
				)))
	} else {
		opt = cel.Declarations(
			decls.NewFunction(
				l.identifier,
				decls.NewOverload(l.identifier,
					paramTypes,
					l.resultType)))
	}

	l.env, err = env.Extend(opt)
	if err != nil {
		return err
	}

	return nil
}

func (l letVariable) String() string {
	return fmt.Sprintf("%s = %s", l.identifier, l.src)
}

func formatParams(params []letFunctionParam) string {
	fmtParams := make([]string, len(params))

	for i, p := range params {
		fmtParams[i] = fmt.Sprintf("%s: %s", p.identifier, UnparseType(p.typeHint))
	}

	return fmt.Sprintf("(%s)", strings.Join(fmtParams, ", "))
}

func (l letFunction) String() string {
	receiverFmt := ""
	if l.receiver != nil {
		receiverFmt = fmt.Sprintf("%s.", UnparseType(l.receiver))
	}

	return fmt.Sprintf("%s%s%s : %s -> %s", receiverFmt, l.identifier, formatParams(l.params), UnparseType(l.resultType), l.src)
}

func (l *letFunction) generateFunction() *functions.Overload {
	argLen := len(l.params)
	if l.receiver != nil {
		argLen += 1
	}
	switch argLen {
	case 1:
		return &functions.Overload{
			Operator: l.identifier,
			Unary:    func(v ref.Val) ref.Val { return l.impl(v) },
		}
	case 2:
		return &functions.Overload{
			Operator: l.identifier,
			Binary:   func(lhs ref.Val, rhs ref.Val) ref.Val { return l.impl(lhs, rhs) },
		}
	default:
		return &functions.Overload{
			Operator: l.identifier,
			Function: l.impl,
		}
	}
}

// Reset plan if we need to recompile based on a dependency change.
func (l *letVariable) clearPlan() {
	l.resultType = nil
	l.env = nil
	l.ast = nil
	l.prog = nil
}

// EvaluationContext context for the repl.
// Handles maintaining state for multiple let expressions.
type EvaluationContext struct {
	letVars []letVariable
	letFns  []letFunction
}

func (ctx *EvaluationContext) indexLetVar(name string) int {
	for idx, el := range ctx.letVars {
		if el.identifier == name {
			return idx
		}
	}
	return -1
}

func (ctx *EvaluationContext) indexLetFn(name string) int {
	for idx, el := range ctx.letFns {
		if el.identifier == name {
			return idx
		}
	}
	return -1
}

func (ctx *EvaluationContext) copy() *EvaluationContext {
	var cpy EvaluationContext
	cpy.letVars = make([]letVariable, len(ctx.letVars))
	copy(cpy.letVars, ctx.letVars)
	cpy.letFns = make([]letFunction, len(ctx.letFns))
	copy(cpy.letFns, ctx.letFns)
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

func (ctx *EvaluationContext) delLetFn(name string) {
	idx := ctx.indexLetFn(name)
	if idx < 0 {
		// no-op if deleting something that's not defined
		return
	}

	ctx.letFns = append(ctx.letFns[:idx], ctx.letFns[idx+1:]...)

	for i := range ctx.letVars {
		ctx.letVars[i].clearPlan()
	}
}

// Add or update an existing let then invalidate any computed plans.
func (ctx *EvaluationContext) addLetVar(name string, expr string, typeHint *exprpb.Type) {
	idx := ctx.indexLetVar(name)
	newVar := letVariable{identifier: name, src: expr, typeHint: typeHint}
	if idx < 0 {
		ctx.letVars = append(ctx.letVars, newVar)
	} else {
		ctx.letVars[idx] = newVar
		for i := idx + 1; i < len(ctx.letVars); i++ {
			// invalidate dependant let exprs
			ctx.letVars[i].clearPlan()
		}
	}
}

// Try to normalize a defined function name as either a namespaced function or a receiver call.
func (ctx *EvaluationContext) resolveFn(name string) (string, *exprpb.Type) {
	leadingDot := ""
	id := name
	if strings.HasPrefix(name, ".") {
		id = strings.TrimLeft(name, ".")
		leadingDot = "."
	}
	qualifiers := strings.Split(id, ".")
	if len(qualifiers) == 1 {
		return qualifiers[0], nil
	}

	namespace := strings.Join(qualifiers[:len(qualifiers)-1], ".")
	id = qualifiers[len(qualifiers)-1]

	maybeType, err := ParseType(leadingDot + namespace)
	if err != nil {
		return name, nil
	}

	switch maybeType.TypeKind.(type) {
	// unsupported type assume it's just namespaced
	case *exprpb.Type_AbstractType_:
	case *exprpb.Type_MessageType:
	case *exprpb.Type_Error:
	case *exprpb.Type_Function:
	default:
		return id, maybeType
	}

	return name, nil
}

// Add or update an existing let then invalidate any computed plans.
func (ctx *EvaluationContext) addLetFn(name string, params []letFunctionParam, resultType *exprpb.Type, expr string) {
	name, receiver := ctx.resolveFn(name)
	idx := ctx.indexLetFn(name)
	newFn := letFunction{identifier: name, params: params, receiver: receiver, resultType: resultType, src: expr}
	if idx < 0 {
		ctx.letFns = append(ctx.letFns, newFn)
	} else {
		ctx.letFns[idx] = newFn
	}

	for i := 0; i < len(ctx.letVars); i++ {
		// invalidate dependant let exprs
		ctx.letVars[i].clearPlan()
	}
}

// programOptions generates the program options for planning.
// Assumes context has been planned.
func (ctx *EvaluationContext) programOptions() cel.ProgramOption {
	var fns = make([]*functions.Overload, len(ctx.letFns))
	for i, fn := range ctx.letFns {
		fns[i] = fn.generateFunction()
	}
	return cel.Functions(fns...)
}

// Evaluator provides basic environment for evaluating an expression with
// applied context.
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

// Attempt to update context in place after an update.
// This is done eagerly to help avoid introducing an invalid 'let' expression.
// The planned expressions are evaluated as needed when evaluating a (non-let) CEL expression.
// Return an error if any of the updates fail.
func updateContextPlans(ctx *EvaluationContext, env *cel.Env) error {
	overloads := make([]*functions.Overload, 0)
	for i := range ctx.letFns {
		letFn := &ctx.letFns[i]
		err := letFn.update(env, overloads)
		if err != nil {
			return fmt.Errorf("error updating %s: %w", letFn, err)
		}
		env = letFn.env
		// if no src, this is declared but not defined.
		if letFn.src != "" {
			overloads = append(overloads, letFn.generateFunction())
		}

	}
	for i := range ctx.letVars {
		el := &ctx.letVars[i]
		// Check if the let variable has a definition and needs to be re-planned
		if el.prog == nil && el.src != "" {
			ast, iss := env.Compile(el.src)
			if iss != nil {
				return fmt.Errorf("error updating %v\n%w", el, iss.Err())
			}

			if el.typeHint != nil && !proto.Equal(ast.ResultType(), el.typeHint) {
				return fmt.Errorf("error updating %v\ntype mismatch got %v expected %v",
					el,
					UnparseType(ast.ResultType()),
					UnparseType(el.typeHint))
			}

			el.ast = ast
			el.resultType = ast.ResultType()

			plan, err := env.Program(ast, ctx.programOptions())
			if err != nil {
				return err
			}
			el.prog = plan
		} else if el.src == "" {
			// Variable is declared but not defined, just update the type checking environment
			el.resultType = el.typeHint
		}
		if el.env == nil {
			env, err := env.Extend(cel.Declarations(decls.NewVar(el.identifier, el.resultType)))
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
func (e *Evaluator) AddLetVar(name string, expr string, typeHint *exprpb.Type) error {
	// copy the current context and attempt to update dependant expressions.
	// if successful, swap the current context with the updated copy.
	ctx := e.ctx.copy()
	ctx.addLetVar(name, expr, typeHint)
	err := updateContextPlans(ctx, e.env)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

// AddLetFn adds a let function to the evaluation context.
func (e *Evaluator) AddLetFn(name string, params []letFunctionParam, resultType *exprpb.Type, expr string) error {
	// copy the current context and attempt to update dependant expressions.
	// if successful, swap the current context with the updated copy.
	cpy := e.ctx.copy()
	cpy.addLetFn(name, params, resultType, expr)
	err := updateContextPlans(cpy, e.env)
	if err != nil {
		return err
	}
	e.ctx = *cpy
	return nil
}

// AddDeclVar declares a variable in the environment but doesn't register an expr with it.
// This allows planning to succeed, but with no value for the variable at runtime.
func (e *Evaluator) AddDeclVar(name string, typeHint *exprpb.Type) error {
	ctx := e.ctx.copy()
	ctx.addLetVar(name, "", typeHint)
	err := updateContextPlans(ctx, e.env)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

// AddDeclFn declares a function in the environment but doesn't register an expr with it.
// This allows planning to succeed, but with no value for the function at runtime.
func (e *Evaluator) AddDeclFn(name string, params []letFunctionParam, typeHint *exprpb.Type) error {
	ctx := e.ctx.copy()
	ctx.addLetFn(name, params, typeHint, "")
	err := updateContextPlans(ctx, e.env)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

// DelLetVar removes a variable from the evaluation context.
// If deleting the variable breaks a later expression, this function will return an error without modifying the context.
func (e *Evaluator) DelLetVar(name string) error {
	ctx := e.ctx.copy()
	ctx.delLetVar(name)
	err := updateContextPlans(ctx, e.env)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

// DelLetVar removes a function from the evaluation context.
// If deleting the function breaks a later expression, this function will return an error without modifying the context.
func (e *Evaluator) DelLetFn(name string) error {
	ctx := e.ctx.copy()
	ctx.delLetFn(name)
	err := updateContextPlans(ctx, e.env)
	if err != nil {
		return err
	}
	e.ctx = *ctx
	return nil
}

// applyContext evaluates the let expressions in the context to build an activation for the given expression.
// returns the environment for compiling and planning the top level CEL expression and an activation with the
// values of the let expressions.
func (e *Evaluator) applyContext() (*cel.Env, interpreter.Activation, error) {
	var vars = make(map[string]interface{})

	for _, el := range e.ctx.letVars {
		if el.prog == nil {
			// Declared but not defined variable so nothing to evaluate
			continue
		}

		val, _, err := el.prog.Eval(vars)
		if val != nil {
			vars[el.identifier] = val
		} else if err != nil {
			return nil, nil, err
		}
	}

	act, err := interpreter.NewActivation(vars)
	if err != nil {
		return nil, nil, err
	}

	env := e.env

	if len(e.ctx.letVars) > 0 {
		env = e.ctx.letVars[len(e.ctx.letVars)-1].env
	} else if len(e.ctx.letFns) > 0 {
		env = e.ctx.letFns[len(e.ctx.letFns)-1].env
	}

	return env, act, nil
}

// Evaluate sets up a CEL evaluation using the current evaluation context.
func (e *Evaluator) Evaluate(expr string) (ref.Val, *exprpb.Type, error) {
	env, act, err := e.applyContext()
	if err != nil {
		return nil, nil, err
	}

	ast, iss := env.Compile(expr)
	if iss != nil {
		return nil, nil, iss.Err()
	}

	p, err := env.Program(ast, e.ctx.programOptions())
	if err != nil {
		return nil, nil, err
	}

	val, _, err := p.Eval(act)
	// expression can be well-formed and result in an error
	return val, ast.ResultType(), err
}
