package interpreter

import (
	"errors"
	"fmt"

	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types/ref"
)

const (
	errorInvalidSignature = "function overload (id: %s) is not matched (got: %s, want: %s)"
	errorMismatch         = "function overload (id: %s) has different attributes (name: %s, got: %v, want: %v)"
	errorNilActivation    = "cannot create a late bind activation with a nil activation"
	errorOverloadNotFound = "unexpected: overload (id: %s) not found."

	unarySignature    = "unary{ func(ref.Val) ref.Val }"
	binarySignature   = "binary{ func(ref.Val, ref.Val) ref.Val }"
	functionSignature = "varargs{ func(...ref.Val) ref.Val }"
)

// NewLateBindActivation creates an activation that wraps the given activation and
// exposes the given function overloads to the evaluation. If the list of overloads
// has duplicates or the given activation is nil, it will return an error.
func NewLateBindActivation(activation Activation, overloads ...*functions.Overload) (LateBindActivation, error) {

	dispatcher := NewDispatcher()
	err := dispatcher.Add(overloads...)
	if err != nil {
		return nil, err
	}

	if activation == nil {
		return nil, errors.New(errorNilActivation)
	}

	return &lateBindActivation{
		vars:       activation,
		dispatcher: dispatcher,
	}, nil
}

// LateBindActivation provides an interface that defines
// the contract for exposing function overloads during
// the evaluation.
//
// This interface enables the integration of external
// implementations of the late bind behaviour, without
// limiting the design to a given concrete type.
type LateBindActivation interface {
	Activation
	// ResolveOverload resolves the function overload that is
	// mapped to overloadId. Implementations of this function
	// are expected to recursively navigate the activation tree
	// by respecting the parent-child relationships to find the
	// first overload definition that is mapped to overloadId.
	ResolveOverload(overloadId string) *functions.Overload
	// ResolveOverloads returns a Dispatcher implementation that maintains all
	// the overload functions that are defined starting from the instance of the
	// concrete type implementing this method. The list is guaranteed to be
	// unique (i.e. with no duplicates). Should duplicates be found, only the
	// first occurrence of the overload is added to the list, thus ensuring
	// that the correct behaviour is being implemented.
	ResolveOverloads() Dispatcher
}

// lateBindActivation is an Activation implementation
// that carries a dispatcher which can be used to
// supply overrides for function overloads during
// evaluation.
type lateBindActivation struct {
	vars       Activation
	dispatcher Dispatcher
}

// ResolveName implemments Activation.ResolveName(string). The
// method defers the name resolution to the activation instance
// that is wrapped.
func (activation *lateBindActivation) ResolveName(name string) (any, bool) {
	return activation.vars.ResolveName(name)
}

// Parent implements Activation.Parent() and returns the
// activation that is wrapped by this struct.
func (activation *lateBindActivation) Parent() Activation {
	return activation.vars
}

// ResolveOverload resolves function overload that is mapped by
// the given overloadId. The implementation first checks if the
// dispatcher configured with the current activation defines an
// overload for overloadId, and if found it returns such overload.
// If the dispatcher does not define such overloads the function
// recursively checks the activation to find any LateBindActivation
// that might declare such overload.
func (activation *lateBindActivation) ResolveOverload(overloadId string) *functions.Overload {

	if activation.dispatcher != nil {
		ovl, found := activation.dispatcher.FindOverload(overloadId)
		if found {
			return ovl
		}
	}

	return resolveOverload(overloadId, activation.vars)
}

// ResolveOverloads returns a Dispatcher implementation that aggregates
// all function overloads definition that are accessible from the current
// activation reference. The preference is given to the overloads of the
// defined dispatcher, and then the hierarchy of activations originating
// from the configured parent activation. If there are any duplicates
func (activation *lateBindActivation) ResolveOverloads() Dispatcher {

	dispatcher := NewDispatcher()
	for _, ovlId := range activation.dispatcher.OverloadIds() {
		ovl, _ := activation.dispatcher.FindOverload(ovlId)
		dispatcher.Add(ovl)
	}

	resolveAllOverloads(dispatcher, activation.vars)

	return dispatcher
}

// decLateBinding returns an InterpretableDecorator
// that transforms the Interpretable to wrap all the
// calls to function to late bindg evaluation structures.
func decLateBinding() InterpretableDecorator {

	return lateBind
}

// lateBind matches the signature of InterpretableDecorator
// and wraps any occurrence of a call to a function with an
// InterpretableCall implementation that inspect the activat
//
// The implementation is recursive and cater for all instances
// of Interpretable that carry expressions. The implemented
// logic operates as follows:
//
//   - evalZeroArity, evalUnary, evalBinary, and evalVarArgs are
//     directly replaced with the corresponding lateBindXXX
//     implementation.
//
//   - evalAnd, evalOr, evalEq, evalNe, evalExhaustiveOr, and
//     evalExhaustiveAnd are mutated by applying lateBind to
//     their term expressions.
//
//   - evalList, evalMap, evalObj are mutated by applying lateBind
//     to their elements, keys and values, or field values.
//
//   - evalFold is mutated by applying lateBind to the condition
//     the iteration range expressions, and the step expression.
//
//   - evalSetMembership is mutated by applying lateBind to both
//     the argument and the set definition.
//
//   - evalWatch is mutated by applying lateBind to wrapped
//     Interpretable implementation.
//
//   - evalWatchConstructor is mutated by applying lateBind to the
//     watcheed InterepretableConstructor implementation.
//
// All other evalXXX entities are left untouched.
//
// If there is any error in applying the transformation the function
// returns a nil Intepretable and such error.
func lateBind(i Interpretable) (Interpretable, error) {

	if i == nil {
		return nil, nil
	}

	switch expr := i.(type) {

	// Group 1: function calls
	// -----------------------
	// evalZeroArity, evalUnary, evalBinary, and evalVarArgs
	// are explicit calls to functions, these are directly
	// wrapped with the corresponding lateBindEvalXXX struct.
	// In addition, we need to apply recursively late binding
	// to the arguments because they are expressions.

	case *evalZeroArity:
		return &lateBindEvalZeroArity{
			target: expr,
		}, nil

	case *evalUnary:

		arg, err := lateBind(expr.arg)
		if err != nil {
			return nil, err
		}
		expr.arg = arg

		return &lateBindEvalUnary{
			target: expr,
		}, nil

	case *evalBinary:

		lhs, rhs, err := lateBindPair(expr.lhs, expr.rhs)
		if err != nil {
			return nil, err
		}
		expr.lhs = lhs
		expr.rhs = rhs

		return &lateBindEvalBinary{
			target: expr,
		}, nil

	case *evalVarArgs:

		args, err := lateBindSlice(expr.args)
		if err != nil {
			return nil, err
		}
		expr.args = args

		return &lateBindEvalVarArgs{
			target: expr,
		}, nil

	// Group 2: logical operators
	// --------------------------
	// These have expressions as arguments (or terms). We need
	// to apply late binding to all the terms a of the operator.

	case *evalEq:
		lhs, rhs, err := lateBindPair(expr.lhs, expr.rhs)
		if err != nil {
			return nil, err
		}
		expr.lhs = lhs
		expr.rhs = rhs

		return expr, nil

	case *evalNe:

		lhs, rhs, err := lateBindPair(expr.lhs, expr.rhs)
		if err != nil {
			return nil, err
		}
		expr.lhs = lhs
		expr.rhs = rhs

		return expr, nil

	case *evalOr:
		mapped, err := lateBindSlice(expr.terms)
		if err != nil {
			return nil, err
		}
		expr.terms = mapped
		return expr, nil

	case *evalAnd:
		mapped, err := lateBindSlice(expr.terms)
		if err != nil {
			return nil, err
		}
		expr.terms = mapped
		return expr, nil

	// exhaustive cases need to be handled too
	// to ensure that when we apply the decorator
	// for exhaustive evaluation we don't loose
	// calls in the modified versions of OR and AND.
	case *evalExhaustiveOr:
		mapped, err := lateBindSlice(expr.terms)
		if err != nil {
			return nil, err
		}
		expr.terms = mapped
		return expr, nil

	case *evalExhaustiveAnd:
		mapped, err := lateBindSlice(expr.terms)
		if err != nil {
			return nil, err
		}
		expr.terms = mapped
		return expr, nil

	// Group 3: complex structures
	// ---------------------------
	// List, maps, and objects in general can have expressions
	// as values for their elements, keys and values, and fields.
	// We need to apply late binding transformations to all of
	// these.

	case *evalList:

		mapped, err := lateBindSlice(expr.elems)
		if err != nil {
			return nil, err
		}

		expr.elems = mapped
		return expr, nil

	case *evalMap:

		keys, err := lateBindSlice(expr.keys)
		if err != nil {
			return nil, err
		}
		values, err := lateBindSlice(expr.vals)
		if err != nil {
			return nil, err
		}

		expr.keys = keys
		expr.vals = values

		return expr, nil

	case *evalObj:
		values, err := lateBindSlice(expr.vals)
		if err != nil {
			return nil, err
		}
		expr.vals = values
		return expr, nil

	// Group 5: Macro
	// --------------
	// Macros can have expressions in it, different types of macros
	// have different parameters. In principle we should only operate
	// on the arguments representing the predicate (and the function
	// for map macros). The other interpretables in the definition of
	// the struct are internally generate and we don't want to touch
	// them.

	case *evalFold:

		iterRange, err := lateBind(expr.iterRange)
		if err != nil {
			return nil, err
		}

		cond, err := lateBind(expr.cond)
		if err != nil {
			return nil, err
		}

		// this is needed for map macros?
		step, err := lateBind(expr.step)
		if err != nil {
			return nil, err
		}

		expr.iterRange = iterRange
		expr.cond = cond
		expr.step = step

		return expr, nil

	// Group 6: Set Membership
	// -----------------------
	// the 'in' operator can have calls to function functions on both
	// sides of the operator, we need to apply late binding transforms
	// to both.
	case *evalSetMembership:

		inst, arg, err := lateBindPair(expr.inst, expr.arg)
		if err != nil {
			return nil, err
		}

		expr.inst = inst
		expr.arg = arg

		return expr, nil

	// evalWatch is a pass-through we need to recursively
	// apply the late binding to the expression that is
	// being watched which may be anything.

	case *evalWatch:

		interpretable, err := lateBind(expr.Interpretable)
		if err != nil {
			return nil, err
		}
		expr.Interpretable = interpretable

		return expr, nil

	case *evalWatchConstructor:

		interpretable, err := lateBind(expr.constructor)
		if err != nil {
			return nil, err
		}
		constructor, ok := interpretable.(InterpretableConstructor)
		if !ok {
			return nil, fmt.Errorf("late bind decorator failed to map ('%T')", expr)
		}
		expr.constructor = constructor

		return expr, nil
	}

	return i, nil
}

// lateBindSlice is a convenience function that iterates lateBind over
// each of the elements of the array of Interpretable passed as argument.
// If there is any error in the execution of lateBind the function stops
// the execution and returns a nil Intepretable and such error. The
// elements rather than being mutated in place are returned in a new
// slice of the same size of the original by preserving the order.
func lateBindSlice(interpretables []Interpretable) ([]Interpretable, error) {

	mapped := make([]Interpretable, len(interpretables))
	for index, interpretable := range interpretables {
		m, err := lateBind(interpretable)
		if err != nil {
			return nil, err
		}
		mapped[index] = m
	}
	return mapped, nil
}

// lateBindPair is a convenience function that executes lateBind on the two
// arguments. This function executes lateBind on the first argument and then
// then on the second argument. If there is any error during the process the
// execution stops and a nil Intepretable pair with the error is returned.
func lateBindPair(lhs Interpretable, rhs Interpretable) (Interpretable, Interpretable, error) {
	mappedLhs, err := lateBind(lhs)
	if err != nil {
		return nil, nil, err
	}
	mappedRhs, err := lateBind(rhs)
	if err != nil {
		return nil, nil, err
	}
	return mappedLhs, mappedRhs, err
}

// LateBindCalls returns a PlannerOption that allows for mutating
// the Intepretable with injections for replacing at evaluation
// time the bindings to the function calls.
func LateBindCalls() PlannerOption {
	return CustomDecorator(decLateBinding())
}

// lateBindEvalZeroArity is the late bind counterpart of
// evalZeroArity and wraps a reference to evalZeroArity.
type lateBindEvalZeroArity struct {
	target *evalZeroArity
}

// ID implements the Interpretable.ID() interface method.
// The unique identifier returned is the one associated
// to the wrapped evalZeroArity reference.
func (zero *lateBindEvalZeroArity) ID() int64 {
	return zero.target.ID()
}

// Function implements the InterpretableCall.Function() interface method.
// The name of the function returned is the one associated to the wrapped
// evalZeroArity reference.
func (zero *lateBindEvalZeroArity) Function() string {
	return zero.target.Function()
}

// OverloadID implements the IntepretableCall.OverloadID() interface method.
// The overload identifier returned is the one associated to the wrapped
// evalZeroArity reference.
func (zero *lateBindEvalZeroArity) OverloadID() string {
	return zero.target.OverloadID()
}

// Args implements the InterpretableCall.Args() interface method.
// The arguments returned are those associated to the wrapped
// evalZeroArity reference.
func (zero *lateBindEvalZeroArity) Args() []Interpretable {
	return zero.target.Args()
}

// Eval implements the Intepretable.Eval(Activation) interface method.
// The implementation first resolves the overload of the function being
// invoked from the activation context, if there is any override and then
// creates a new instance of evalZeroArity with the replaced function
// implementation for the overload. It then invokes the Eval on the newly
// created struct.
//
// NOTE: the reason why we create a fresh new instance of evalZeroArity is
//
//	to make sure that the substitution of the overload only affects
//	the current call to Eval and it is not permanently stored in the
//	original evalZeroArity reference. This enables to reuse cached
//	programs multiple times with different types of activations and
//	maintains a consistent result all the times:
//
//	- if the activation context has an overload for this call that one
//	  is used.
//	- if the activation context does not have an overload for this call
//	  the one originally bound during the planning phase is used.
func (zero *lateBindEvalZeroArity) Eval(ctx Activation) ref.Val {

	overloadId := zero.target.OverloadID()
	subject := zero.target
	overload := resolveOverload(overloadId, ctx)
	if overload != nil {
		subject = &evalZeroArity{
			id:       zero.target.ID(),
			function: zero.target.Function(),
			overload: overloadId,
			impl:     overload.Function,
		}
	}
	return subject.Eval(ctx)

}

// lateBindEvalUnary is the late bind counterpart of
// evalUnary and wraps a reference to evalUnary.
type lateBindEvalUnary struct {
	target *evalUnary
}

// ID implements the Interpretable.ID() interface method.
// The unique identifier returned is the one associated
// to the wrapped evalUnary reference.
func (un *lateBindEvalUnary) ID() int64 {
	return un.target.ID()
}

// Function implements the InterpretableCall.Function() interface method.
// The name of the function returned is the one associated to the wrapped
// evalUnary reference.
func (un *lateBindEvalUnary) Function() string {
	return un.target.Function()
}

// OverloadID implements the IntepretableCall.OverloadID() interface method.
// The overload identifier returned is the one associated to the wrapped
// evalUnary reference.
func (un *lateBindEvalUnary) OverloadID() string {
	return un.target.OverloadID()
}

// Args implements the InterpretableCall.Args() interface method.
// The arguments returned are those associated to the wrapped
// evalUnary reference.
func (un *lateBindEvalUnary) Args() []Interpretable {
	return un.target.Args()
}

// Eval implements the Intepretable.Eval(Activation) interface method.
// The implementation first resolves the overload of the function being
// invoked from the activation context, if there is any override and then
// creates a new instance of evalUnary with the replaced function
// implementation for the overload. It then invokes the Eval on the newly
// created struct.
//
// NOTE: the reason why we create a fresh new instance of evalUnary is to
//
//	make sure that the substitution of the overload only affects the
//	current call to Eval and it is not permanently stored in the
//	original evalUnary reference. This enables to reuse cached program
//	multiple times with different types of activations and maintains a
//	consistent result all the times:
//
//	- if the activation context has an overload for this call that one
//	  is used.
//	- if the activation context does not have an overload for this call
//	  the one originally bound during the planning phase is used.
func (un *lateBindEvalUnary) Eval(ctx Activation) ref.Val {

	overloadId := un.target.OverloadID()
	subject := un.target
	overload := resolveOverload(overloadId, ctx)
	if overload != nil {
		subject = &evalUnary{
			id:        un.target.ID(),
			function:  un.target.Function(),
			overload:  overloadId,
			impl:      overload.Unary,
			trait:     overload.OperandTrait,
			nonStrict: overload.NonStrict,
		}
	}
	return subject.Eval(ctx)

}

// lateBindEvalBinary is the late bind counterpart of
// evalBinary and wraps a reference to evalBinary.
type lateBindEvalBinary struct {
	target *evalBinary
}

// ID implements the Interpretable.ID() interface method.
// The unique identifier returned is the one associated
// to the wrapped evalBinary reference.
func (bin *lateBindEvalBinary) ID() int64 {
	return bin.target.ID()
}

// Function implements the InterpretableCall.Function() interface method.
// The name of the function returned is the one associated to the wrapped
// evalBinary reference.
func (bin *lateBindEvalBinary) Function() string {
	return bin.target.Function()
}

// OverloadID implements the IntepretableCall.OverloadID() interface method.
// The overload identifier returned is the one associated to the wrapped
// evalBinary reference.
func (bin *lateBindEvalBinary) OverloadID() string {
	return bin.target.OverloadID()
}

// Args implements the InterpretableCall.Args() interface method.
// The arguments returned are those associated to the wrapped
// evalBinary reference.
func (bin *lateBindEvalBinary) Args() []Interpretable {
	return bin.target.Args()
}

// Eval implements the Intepretable.Eval(Activation) interface method.
// The implementation first resolves the overload of the function being
// invoked from the activation context, if there is any override and then
// creates a new instance of evalBinary with the replaced function
// implementation for the overload. It then invokes the Eval on the newly
// created struct.
//
// NOTE: the reason why we create a fresh new instance of evalBinary is to
//
//	make sure that the substitution of the overload only affects the
//	current call to Eval and it is not permanently stored in the
//	original evalBinary reference. This enables to reuse cached program
//	multiple times with different types of activations and maintains a
//	consistent result all the times:
//
//	- if the activation context has an overload for this call that one
//	  is used.
//	- if the activation context does not have an overload for this call
//	  the one originally bound during the planning phase is used.
func (bin *lateBindEvalBinary) Eval(ctx Activation) ref.Val {

	overloadId := bin.target.OverloadID()
	subject := bin.target
	overload := resolveOverload(overloadId, ctx)
	if overload != nil {
		args := bin.target.Args()
		subject = &evalBinary{
			id:        bin.target.ID(),
			function:  bin.target.Function(),
			overload:  overloadId,
			lhs:       args[0],
			rhs:       args[1],
			impl:      overload.Binary,
			trait:     overload.OperandTrait,
			nonStrict: overload.NonStrict,
		}
	}
	return subject.Eval(ctx)

}

// lateBindEvalVarArgs is the late bind counterpart of
// evalVarArgs and wraps a reference to evalVarArgs.
type lateBindEvalVarArgs struct {
	target *evalVarArgs
}

// ID implements the Interpretable.ID() interface method.
// The unique identifier returned is the one associated
// to the wrapped evalVarArgs reference.
func (fn *lateBindEvalVarArgs) ID() int64 {
	return fn.target.ID()
}

// Function implements the InterpretableCall.Function() interface method.
// The name of the function returned is the one associated to the wrapped
// evalVarArgs reference.
func (fn *lateBindEvalVarArgs) Function() string {
	return fn.target.Function()
}

// OverloadID implements the IntepretableCall.OverloadID() interface method.
// The overload identifier returned is the one associated to the wrapped
// evalVarArgs reference.
func (fn *lateBindEvalVarArgs) OverloadID() string {
	return fn.target.OverloadID()
}

// Args implements the InterpretableCall.Args() interface method.
// The arguments returned are those associated to the wrapped
// evalVarArgs reference.
func (fn *lateBindEvalVarArgs) Args() []Interpretable {
	return fn.target.Args()
}

// Eval implements the Intepretable.Eval(Activation) interface method.
// The implementation first resolves the overload of the function being
// invoked from the activation context, if there is any override and then
// creates a new instance of evalVarArgs with the replaced function
// implementation for the overload. It then invokes the Eval on the newly
// created struct.
//
// NOTE: the reason why we create a fresh new instance of evalVarArgs is to
//
//	make sure that the substitution of the overload only affects the
//	current call to Eval and it is not permanently stored in the
//	original evalBinary reference. This enables to reuse cached program
//	multiple times with different types of activations and maintains a
//	consistent result all the times:
//
//	- if the activation context has an overload for this call that one
//	  is used.
//	- if the activation context does not have an overload for this call
//	  the one originally bound during the planning phase is used.
func (fn *lateBindEvalVarArgs) Eval(ctx Activation) ref.Val {

	overloadId := fn.target.OverloadID()
	subject := fn.target
	overload := resolveOverload(overloadId, ctx)
	if overload != nil {
		subject = &evalVarArgs{
			id:        fn.target.ID(),
			function:  fn.target.Function(),
			overload:  overloadId,
			args:      fn.Args(),
			impl:      overload.Function,
			trait:     overload.OperandTrait,
			nonStrict: overload.NonStrict,
		}
	}
	return subject.Eval(ctx)

}

// resolveOverload travels the hierarchy of activations originating from the given
// Activation implementation to find the overload associatd to overloadId. Since the
// Activation APIs allow for different types of activations and compositions we need
// to ensure that if there is any valid overload that is mapped to overloadId we can
// find it.
func resolveOverload(overloadId string, activation Activation) *functions.Overload {

	if activation == nil {
		return nil
	}

	switch act := activation.(type) {
	case *mapActivation:
		return nil
	case *emptyActivation:
		return nil
	case *partActivation:
		return resolveOverload(overloadId, act.Activation)
	case *hierarchicalActivation:
		ovl := resolveOverload(overloadId, act.child)
		if ovl == nil {
			return resolveOverload(overloadId, act.parent)
		}
		return ovl
	case LateBindActivation:

		return act.ResolveOverload(overloadId)
	}

	// this is to ensure that if there are
	// more types of activations added we
	// can default to nil.
	return nil

}

// ValidateOverloads ensures that if the activation contains an overload function
// its signature matches the one associated to the same overload identifier in
// the dispatcher otherwise throws an error. If the activation defines more function
// overloads, those won't be considered in the validation.
func ValidateOverloads(original Dispatcher, activation Activation) error {

	// we create a
	aggregate := NewDispatcher()
	resolveAllOverloads(aggregate, activation)

	overloads := original.OverloadIds()
	for _, overloadId := range overloads {

		refOvl, found := original.FindOverload(overloadId)
		if !found {
			return fmt.Errorf(errorOverloadNotFound, overloadId)
		}

		ovl, found := aggregate.FindOverload(overloadId)
		if found {
			// we need to make sure that the overloads are
			// matching.

			result := matchSignature(overloadId, refOvl, ovl)
			if result != nil {
				return result
			}
		}
	}

	return nil
}

// resolveAllOverloads aggregates all function overloads defined in the
// activation into a single dispatcher so that they can be easily checked
// at once when we validate the overloads.
func resolveAllOverloads(aggregate Dispatcher, activation Activation) {

	if activation == nil {
		return
	}
	switch act := activation.(type) {
	case *mapActivation:
		return
	case *emptyActivation:
		return
	case *partActivation:
		resolveAllOverloads(aggregate, act.Activation)
	case *hierarchicalActivation:
		resolveAllOverloads(aggregate, act.child)
		resolveAllOverloads(aggregate, act.parent)
	case LateBindActivation:

		// the implementation of Overloads() is expected to be
		// recursive, therefore we don't need to look any further.
		dispatcher := act.ResolveOverloads()
		for _, ovlId := range dispatcher.OverloadIds() {
			ovl, found := dispatcher.FindOverload(ovlId)
			if found {
				// note we don't need to check an error because if there
				// is an error the overload is already defined. This may
				// happen because we nest multiple activation with late
				// binding capabilities and one may shadow another as it
				// happens for variable names. Since the activations are
				// visitedin the correct order this is expected behaviour
				aggregate.Add(ovl)
			}
		}
	}

}

// matchSignature compares the two overload definitions and returns an error
// if the overload function does not have a matching signature with the
// reference overload. The only check we can implement is over the number of
// parameters and the attributes of the overload.
//
// The impmlementation verifies the following:
//
// - if refOvl.Unary is not nil, the expectation is that ovl.Unary is not nil.
// - if refOvl.Binry is not nil, the expectation is that ovl.Binary is not nil.
// - if refOvl.Function not nil, the expectation is that ovl.Fnuction is not nil.
// - refOvl.NotStrict and ovl.NonStrict must be the same.
// - refOvl.OperandTrait and ovl.OperandTrait must be the same.
// - refOvl.Operator and ovl.Operator must be the same.
func matchSignature(overloadId string, refOvl *functions.Overload, ovl *functions.Overload) error {

	got := "<nil>"

	if refOvl.Unary != nil {

		if ovl.Unary == nil {

			if ovl.Binary != nil {
				got = binarySignature

			} else if ovl.Function != nil {
				got = functionSignature
			}
			return fmt.Errorf(errorInvalidSignature, overloadId, got, unarySignature)
		}
	} else if refOvl.Binary != nil {

		if ovl.Binary == nil {

			if ovl.Unary != nil {
				got = unarySignature
			} else if ovl.Function != nil {
				got = functionSignature
			}

			return fmt.Errorf(errorInvalidSignature, overloadId, got, binarySignature)
		}

	} else if refOvl.Function != nil {

		if ovl.Function == nil {

			if ovl.Unary != nil {
				got = unarySignature
			} else if ovl.Binary != nil {
				got = binarySignature
			}

			return fmt.Errorf(errorInvalidSignature, overloadId, got, functionSignature)

		}
	}

	if refOvl.NonStrict != ovl.NonStrict {

		return fmt.Errorf(errorMismatch, overloadId, "NonStrict", ovl.NonStrict, refOvl.NonStrict)
	}
	if refOvl.OperandTrait != ovl.OperandTrait {

		return fmt.Errorf(errorMismatch, overloadId, "OperandTrait", ovl.OperandTrait, refOvl.OperandTrait)
	}
	if refOvl.Operator != ovl.Operator {

		return fmt.Errorf(errorMismatch, overloadId, "Operator", refOvl.Operator, ovl.Operator)
	}

	return nil
}
