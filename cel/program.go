// Copyright 2019 Google LLC
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
	"fmt"

	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Program is an evaluable view of an Ast.
type Program interface {
	// Eval returns the result of an evaluation of the Ast and environment against the input vars.
	//
	// The vars value may either be an `interpreter.Activation` or a `map[string]interface{}`.
	//
	// If the `OptTrackState` or `OptExhaustiveEval` flags are used, the `details` response will
	// be non-nil. Given this caveat on `details`, the return state from evaluation will be:
	//
	// *  `val`, `details`, `nil` - Successful evaluation of a non-error result.
	// *  `val`, `details`, `err` - Successful evaluation to an error result.
	// *  `nil`, `details`, `err` - Unsuccessful evaluation.
	//
	// An unsuccessful evaluation is typically the result of a series of incompatible `EnvOption`
	// or `ProgramOption` values used in the creation of the evaluation environment or executable
	// program.
	Eval(vars interface{}) (ref.Val, EvalDetails, error)
}

// EvalDetails holds additional information observed during the Eval() call.
type EvalDetails interface {
	// State of the evaluation, non-nil if the OptTrackState or OptExhaustiveEval is specified
	// within EvalOptions.
	State() interpreter.EvalState
}

// NoKnownVars returns an activation that indicates all attributes are unknown.
func NoKnownVars() interpreter.Activation {
	return interpreter.UnknownActivation()
}

// NoVars returns an empty Activation.
func NoVars() interpreter.Activation {
	return interpreter.EmptyActivation()
}

// Vars returns an Activation containing the program input and an optional list of attributes which
// are known at evaluation time to be unknown.
func Vars(input interface{}, unknowns ...interpreter.Attribute) (interpreter.Activation, error) {
	vars, err := interpreter.NewActivation(input)
	if err != nil {
		return nil, err
	}
	if len(unknowns) == 0 {
		return vars, nil
	}
	return interpreter.NewPartialActivation(vars, unknowns)
}

// Unknown specifies an Attribute whose value is not known at runtime.
//
// Not all inputs to an expression have the same cost. Some are cheap to provide and others
// require I/O or heavy compute to resolve. Marking expensive attributes as unknown can be a
// useful tool for improving the operational efficiency of expression evaluation, as it allows
// the cheap data to be evaluated quickly and often conclusively. When the expression result
// cannot be computed, the result will be a types.Unknown value containing the ids of the
// sub-expressions that resulted in the Unknown value.
//
// Unknowns are treated like recoverable errors within CEL and take precedence over types.Err
// values.
func Unknown(varName string, pathElems ...interface{}) (interpreter.Attribute, error) {
	return interpreter.NewUnknownAttribute(varName, pathElems...)
}

// evalDetails is the internal implementation of the EvalDetails interface.
type evalDetails struct {
	state interpreter.EvalState
}

// State implements the Result interface method.
func (ed *evalDetails) State() interpreter.EvalState {
	return ed.state
}

// newProgram creates a program instance with an environment, an ast, and an optional list of
// ProgramOption values.
//
// If the program cannot be configured the prog will be nil, with a non-nil error response.
func newProgram(e *env, ast Ast, opts ...ProgramOption) (Program, error) {
	// Build the dispatcher and default program value.
	disp := interpreter.NewDispatcher()
	p := &prog{env: e, dispatcher: disp}

	// Configure the program via the ProgramOption values.
	var err error
	for _, opt := range opts {
		if opt == nil {
			return nil, fmt.Errorf("program options should be non-nil")
		}
		p, err = opt(p)
		if err != nil {
			return nil, err
		}
	}

	// Configure the resolver to support unknown attributes if the eval option is set.
	resolver := interpreter.NewResolver(e.adapter)
	if p.evalOpts&OptUnknownAttributes == OptUnknownAttributes {
		resolver = interpreter.NewUnknownResolver(resolver)
	}
	p.interpreter = interpreter.NewInterpreter(
		disp, e.pkg, e.provider, e.adapter, resolver)

	// Translate the EvalOption flags into InterpretableDecorator instances.
	decorators := []interpreter.InterpretableDecorator{}
	// Enable constant folding first.
	if p.evalOpts&OptOptimize == OptOptimize {
		decorators = append(decorators, interpreter.Optimize())
	}
	// Enable exhaustive eval over state tracking since it offers a superset of features.
	if p.evalOpts&OptExhaustiveEval == OptExhaustiveEval {
		// State tracking requires that each Eval() call operate on an isolated EvalState
		// object; hence, the presence of the factory.
		factory := func(state interpreter.EvalState) (Program, error) {
			decs := append(decorators, interpreter.ExhaustiveEval(state))
			return p.reset(ast, decs)
		}
		return initProgGen(factory)
	}
	// Enable state tracking last since it too requires the factory approach but is less
	// featured than the ExhaustiveEval decorator.
	if p.evalOpts&OptTrackState == OptTrackState {
		factory := func(state interpreter.EvalState) (Program, error) {
			decs := append(decorators, interpreter.TrackState(state))
			return p.reset(ast, decs)
		}
		return initProgGen(factory)
	}
	return p.init(ast, decorators)
}

// prog is the internal implementation of the Program interface.
type prog struct {
	*env
	evalOpts      EvalOption
	defaultVars   interpreter.Activation
	dispatcher    interpreter.Dispatcher
	interpreter   interpreter.Interpreter
	interpretable interpreter.Interpretable
}

// Eval implements the Program interface method.
func (p *prog) Eval(input interface{}) (v ref.Val, det EvalDetails, err error) {
	// Configure error recovery for unexpected panics during evaluation. Note, the use of named
	// return values makes it possible to modify the error response during the recovery
	// function.
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("internal error: %v", r)
		}
	}()
	// Build a hierarchical activation if there are default vars set.
	vars, err := interpreter.NewActivation(input)
	if err != nil {
		return
	}
	if p.defaultVars != nil {
		vars = interpreter.NewHierarchicalActivation(p.defaultVars, vars)
	}
	v = p.interpretable.Eval(vars)
	// The output of an internal Eval may have a value (`v`) that is a types.Err. This step
	// translates the CEL value to a Go error response. This interface does not quite match the
	// RPC signature which allows for multiple errors to be returned, but should be sufficient.
	if types.IsError(v) {
		err = v.Value().(error)
	}
	return
}

// init builds a checked or unchecked interpretable plan based on whether the Ast is type-checked.
func (p *prog) init(ast Ast, decorators []interpreter.InterpretableDecorator) (Program, error) {
	var err error
	// Unchecked programs do not contain type and reference information and may be
	// slower to execute than their checked counterparts.
	if !ast.IsChecked() {
		p.interpretable, err =
			p.interpreter.NewUncheckedInterpretable(ast.Expr(), decorators...)
		if err != nil {
			return nil, err
		}
		return p, nil
	}
	// When the AST has been checked it contains metadata that can be used to speed up program
	// execution.
	var checked *exprpb.CheckedExpr
	checked, err = AstToCheckedExpr(ast)
	if err != nil {
		return nil, err
	}
	p.interpretable, err = p.interpreter.NewInterpretable(checked, decorators...)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// reset clones the current program and regenerates the interpretable plan using the input Ast
// and decorators.
func (p *prog) reset(ast Ast, decorators []interpreter.InterpretableDecorator) (Program, error) {
	return (&prog{
		evalOpts:    p.evalOpts,
		defaultVars: p.defaultVars,
		env:         p.env,
		dispatcher:  p.dispatcher,
		interpreter: p.interpreter,
	}).init(ast, decorators)
}

// Eval implements the Program interface method.
func (gen *progGen) Eval(input interface{}) (ref.Val, EvalDetails, error) {
	// The factory based Eval() differs from the standard evaluation model in that it generates a
	// new EvalState instance for each call to ensure that unique evaluations yield unique stateful
	// results.
	state := interpreter.NewEvalState()
	det := &evalDetails{state: state}

	// Generate a new instance of the interpretable using the factory configured during the call to
	// newProgram(). It is incredibly unlikely that the factory call will generate an error given
	// the factory test performed within the Program() call.
	p, err := gen.factory(state)
	if err != nil {
		return nil, det, err
	}

	// Evaluate the input, returning the result and the 'state' within EvalDetails.
	v, _, err := p.Eval(input)
	if err != nil {
		return v, det, err
	}
	return v, det, nil
}

// progFactory is a helper alias for marking a program creation factory function.
type progFactory func(interpreter.EvalState) (Program, error)

// progGen holds a reference to a progFactory instance and implements the Program interface.
type progGen struct {
	factory progFactory
}

// initProgGen tests the factory object by calling it once and returns a factory-based Program if
// the test is successful.
func initProgGen(factory progFactory) (Program, error) {
	// Test the factory to make sure that configuration errors are spotted at config
	_, err := factory(interpreter.NewEvalState())
	if err != nil {
		return nil, err
	}
	return &progGen{factory: factory}, nil
}
