// Copyright 2023 Google LLC
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

// Package decls contains function and variable declaration structs and helper methods.
package decls

import (
	"fmt"
	"strings"

	chkdecls "github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// NewFunction creates a new function declaration with a set of function options to configure overloads
// and function definitions (implementations).
//
// Functions are checked for name collisions and singleton redefinition.
func NewFunction(name string, opts ...FunctionOpt) (*FunctionDecl, error) {
	fn := &FunctionDecl{
		Name:             name,
		Overloads:        map[string]*OverloadDecl{},
		overloadOrdinals: []string{},
	}
	var err error
	for _, opt := range opts {
		fn, err = opt(fn)
		if err != nil {
			return nil, err
		}
	}
	if len(fn.Overloads) == 0 {
		return nil, fmt.Errorf("function %s must have at least one overload", name)
	}
	return fn, nil
}

// FunctionDecl defines a function name, overload set, and optionally a singleton definition for all
// overload instances.
type FunctionDecl struct {
	// Name of the function in human-readable terms, e.g. 'contains' of 'math.least'
	Name string

	// Overloads associated with the function name.
	Overloads map[string]*OverloadDecl

	// Singleton implementation of the function for all overloads.
	//
	// If this option is set, an error will occur if any overloads specify a per-overload implementation
	// or if another function with the same name attempts to redefine the singleton.
	Singleton *functions.Overload

	// disableTypeGuards is a performance optimization to disable detailed runtime type checks which could
	// add overhead on common operations. Setting this option true leaves error checks and argument checks
	// intact.
	disableTypeGuards bool

	// state indicates that the binding should be provided as a declaration, as a runtime binding, or both.
	state declarationState

	// overloadOrdinals indicates the order in which the overload was declared.
	overloadOrdinals []string
}

type declarationState int

const (
	declarationStateUnset declarationState = iota
	declarationDisabled
	declarationEnabled
)

// IsDeclarationDisabled indicates that the function implementation should be added to the dispatcher, but the
// declaration should not be exposed for use in expressions.
func (f *FunctionDecl) IsDeclarationDisabled() bool {
	return f.state == declarationDisabled
}

// Merge combines an existing function declaration with another.
//
// If a function is extended, by say adding new overloads to an existing function, then it is merged with the
// prior definition of the function at which point its overloads must not collide with pre-existing overloads
// and its bindings (singleton, or per-overload) must not conflict with previous definitions either.
func (f *FunctionDecl) Merge(other *FunctionDecl) (*FunctionDecl, error) {
	if f == other {
		return f, nil
	}
	if f.Name != other.Name {
		return nil, fmt.Errorf("cannot merge unrelated functions. %s and %s", f.Name, other.Name)
	}
	merged := &FunctionDecl{
		Name:             f.Name,
		Overloads:        make(map[string]*OverloadDecl, len(f.Overloads)),
		Singleton:        f.Singleton,
		overloadOrdinals: make([]string, len(f.Overloads)),
		// if one function is expecting type-guards and the other is not, then they
		// must not be disabled.
		disableTypeGuards: f.disableTypeGuards && other.disableTypeGuards,
		// default to the current functions declaration state.
		state: f.state,
	}
	// If the other state indicates that the declaration should be explicitly enabled or
	// disabled, then update the merged state with the most recent value.
	if other.state != declarationStateUnset {
		merged.state = other.state
	}
	// baseline copy of the overloads and their ordinals
	copy(merged.overloadOrdinals, f.overloadOrdinals)
	for oID, o := range f.Overloads {
		merged.Overloads[oID] = o
	}
	// overloads and their ordinals are added from the left
	for _, oID := range other.overloadOrdinals {
		o := other.Overloads[oID]
		err := merged.AddOverload(o)
		if err != nil {
			return nil, fmt.Errorf("function declaration merge failed: %v", err)
		}
	}
	if other.Singleton != nil {
		if merged.Singleton != nil {
			return nil, fmt.Errorf("function already has singleton binding: %s", f.Name)
		}
		merged.Singleton = other.Singleton
	}
	return merged, nil
}

// AddOverload ensures that the new overload does not collide with an existing overload signature;
// however, if the function signatures are identical, the implementation may be rewritten as its
// difficult to compare functions by object identity.
func (f *FunctionDecl) AddOverload(overload *OverloadDecl) error {
	for oID, o := range f.Overloads {
		if oID != overload.ID && o.SignatureOverlaps(overload) {
			return fmt.Errorf("overload signature collision in function %s: %s collides with %s", f.Name, oID, overload.ID)
		}
		if oID == overload.ID {
			if o.SignatureEquals(overload) && o.NonStrict == overload.NonStrict {
				// Allow redefinition of an overload implementation so long as the signatures match.
				f.Overloads[oID] = overload
				return nil
			}
			return fmt.Errorf("overload redefinition in function. %s: %s has multiple definitions", f.Name, oID)
		}
	}
	f.overloadOrdinals = append(f.overloadOrdinals, overload.ID)
	f.Overloads[overload.ID] = overload
	return nil
}

// Bindings produces a set of function bindings, if any are defined.
func (f *FunctionDecl) Bindings() ([]*functions.Overload, error) {
	overloads := []*functions.Overload{}
	nonStrict := false
	for _, oID := range f.overloadOrdinals {
		o := f.Overloads[oID]
		if o.hasBinding() {
			overload := &functions.Overload{
				Operator:     o.ID,
				Unary:        o.guardedUnaryOp(f.Name, f.disableTypeGuards),
				Binary:       o.guardedBinaryOp(f.Name, f.disableTypeGuards),
				Function:     o.guardedFunctionOp(f.Name, f.disableTypeGuards),
				OperandTrait: o.OperandTrait,
				NonStrict:    o.NonStrict,
			}
			overloads = append(overloads, overload)
			nonStrict = nonStrict || o.NonStrict
		}
	}
	if f.Singleton != nil {
		if len(overloads) != 0 {
			return nil, fmt.Errorf("singleton function incompatible with specialized overloads: %s", f.Name)
		}
		overloads = []*functions.Overload{
			{
				Operator:     f.Name,
				Unary:        f.Singleton.Unary,
				Binary:       f.Singleton.Binary,
				Function:     f.Singleton.Function,
				OperandTrait: f.Singleton.OperandTrait,
			},
		}
		// fall-through to return single overload case.
	}
	if len(overloads) == 0 {
		return overloads, nil
	}
	// Single overload. Replicate an entry for it using the function name as well.
	if len(overloads) == 1 {
		if overloads[0].Operator == f.Name {
			return overloads, nil
		}
		return append(overloads, &functions.Overload{
			Operator:     f.Name,
			Unary:        overloads[0].Unary,
			Binary:       overloads[0].Binary,
			Function:     overloads[0].Function,
			NonStrict:    overloads[0].NonStrict,
			OperandTrait: overloads[0].OperandTrait,
		}), nil
	}
	// All of the defined overloads are wrapped into a top-level function which
	// performs dynamic dispatch to the proper overload based on the argument types.
	bindings := append([]*functions.Overload{}, overloads...)
	funcDispatch := func(args ...ref.Val) ref.Val {
		for _, o := range f.Overloads {
			// During dynamic dispatch over multiple functions, signature agreement checks
			// are preserved in order to assist with the function resolution step.
			switch len(args) {
			case 1:
				if o.UnaryOp != nil && o.matchesRuntimeSignature( /* disableTypeGuards=*/ false, args...) {
					return o.UnaryOp(args[0])
				}
			case 2:
				if o.BinaryOp != nil && o.matchesRuntimeSignature( /* disableTypeGuards=*/ false, args...) {
					return o.BinaryOp(args[0], args[1])
				}
			}
			if o.FunctionOp != nil && o.matchesRuntimeSignature( /* disableTypeGuards=*/ false, args...) {
				return o.FunctionOp(args...)
			}
			// eventually this will fall through to the noSuchOverload below.
		}
		return MaybeNoSuchOverload(f.Name, args...)
	}
	function := &functions.Overload{
		Operator:  f.Name,
		Function:  funcDispatch,
		NonStrict: nonStrict,
	}
	return append(bindings, function), nil
}

// MaybeNoSuchOverload determines whether to propagate an error if one is provided as an argument, or
// to return an unknown set, or to produce a new error for a missing function signature.
func MaybeNoSuchOverload(funcName string, args ...ref.Val) ref.Val {
	argTypes := make([]string, len(args))
	var unk types.Unknown
	for i, arg := range args {
		if types.IsError(arg) {
			return arg
		}
		if types.IsUnknown(arg) {
			unk = append(unk, arg.(types.Unknown)...)
		}
		argTypes[i] = arg.Type().TypeName()
	}
	if len(unk) != 0 {
		return unk
	}
	signature := strings.Join(argTypes, ", ")
	return types.NewErr("no such overload: %s(%s)", funcName, signature)
}

// FunctionOpt defines a functional option for mutating a function declaration.
type FunctionOpt func(*FunctionDecl) (*FunctionDecl, error)

// DisableTypeGuards disables automatically generated function invocation guards on direct overload calls.
// Type guards remain on during dynamic dispatch for parsed-only expressions.
func DisableTypeGuards(value bool) FunctionOpt {
	return func(fn *FunctionDecl) (*FunctionDecl, error) {
		fn.disableTypeGuards = value
		return fn, nil
	}
}

// DisableDeclaration indicates that the function declaration should be disabled, but the runtime function
// binding should be provided. Marking a function as runtime-only is a safe way to manage deprecations
// of function declarations while still preserving the runtime behavior for previously compiled expressions.
func DisableDeclaration(value bool) FunctionOpt {
	return func(fn *FunctionDecl) (*FunctionDecl, error) {
		if value {
			fn.state = declarationDisabled
		} else {
			fn.state = declarationEnabled
		}
		return fn, nil
	}
}

// SingletonUnaryBinding creates a singleton function definition to be used for all function overloads.
//
// Note, this approach works well if operand is expected to have a specific trait which it implements,
// e.g. traits.ContainerType. Otherwise, prefer per-overload function bindings.
func SingletonUnaryBinding(fn functions.UnaryOp, traits ...int) FunctionOpt {
	trait := 0
	for _, t := range traits {
		trait = trait | t
	}
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		if f.Singleton != nil {
			return nil, fmt.Errorf("function already has a singleton binding: %s", f.Name)
		}
		f.Singleton = &functions.Overload{
			Operator:     f.Name,
			Unary:        fn,
			OperandTrait: trait,
		}
		return f, nil
	}
}

// SingletonBinaryBinding creates a singleton function definition to be used with all function overloads.
//
// Note, this approach works well if operand is expected to have a specific trait which it implements,
// e.g. traits.ContainerType. Otherwise, prefer per-overload function bindings.
func SingletonBinaryBinding(fn functions.BinaryOp, traits ...int) FunctionOpt {
	trait := 0
	for _, t := range traits {
		trait = trait | t
	}
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		if f.Singleton != nil {
			return nil, fmt.Errorf("function already has a singleton binding: %s", f.Name)
		}
		f.Singleton = &functions.Overload{
			Operator:     f.Name,
			Binary:       fn,
			OperandTrait: trait,
		}
		return f, nil
	}
}

// SingletonFunctionBinding creates a singleton function definition to be used with all function overloads.
//
// Note, this approach works well if operand is expected to have a specific trait which it implements,
// e.g. traits.ContainerType. Otherwise, prefer per-overload function bindings.
func SingletonFunctionBinding(fn functions.FunctionOp, traits ...int) FunctionOpt {
	trait := 0
	for _, t := range traits {
		trait = trait | t
	}
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		if f.Singleton != nil {
			return nil, fmt.Errorf("function already has a singleton binding: %s", f.Name)
		}
		f.Singleton = &functions.Overload{
			Operator:     f.Name,
			Function:     fn,
			OperandTrait: trait,
		}
		return f, nil
	}
}

// Overload defines a new global overload with an overload id, argument types, and result type. Through the
// use of OverloadOpt options, the overload may also be configured with a binding, an operand trait, and to
// be non-strict.
//
// Note: function bindings should be commonly configured with Overload instances whereas operand traits and
// strict-ness should be rare occurrences.
func Overload(overloadID string, args []*Type, resultType *Type, opts ...OverloadOpt) FunctionOpt {
	return newOverload(overloadID, false, args, resultType, opts...)
}

// MemberOverload defines a new receiver-style overload (or member function) with an overload id, argument types,
// and result type. Through the use of OverloadOpt options, the overload may also be configured with a binding,
// an operand trait, and to be non-strict.
//
// Note: function bindings should be commonly configured with Overload instances whereas operand traits and
// strict-ness should be rare occurrences.
func MemberOverload(overloadID string, args []*Type, resultType *Type, opts ...OverloadOpt) FunctionOpt {
	return newOverload(overloadID, true, args, resultType, opts...)
}

func newOverload(overloadID string, memberFunction bool, args []*Type, resultType *Type, opts ...OverloadOpt) FunctionOpt {
	return func(f *FunctionDecl) (*FunctionDecl, error) {
		overload, err := newOverloadInternal(overloadID, memberFunction, args, resultType, opts...)
		if err != nil {
			return nil, err
		}
		err = f.AddOverload(overload)
		if err != nil {
			return nil, err
		}
		return f, nil
	}
}

func newOverloadInternal(overloadID string, memberFunction bool, args []*Type, resultType *Type, opts ...OverloadOpt) (*OverloadDecl, error) {
	overload := &OverloadDecl{
		ID:               overloadID,
		ArgTypes:         args,
		ResultType:       resultType,
		IsMemberFunction: memberFunction,
	}
	var err error
	for _, opt := range opts {
		overload, err = opt(overload)
		if err != nil {
			return nil, err
		}
	}
	return overload, nil
}

// OverloadDecl contains the definition of a single overload id with a specific signature, and an optional
// implementation.
type OverloadDecl struct {
	// ID mirrors the overload signature and provides a unique id which may be referenced within the type-checker
	// and interpreter to optimize performance.
	//
	// The ID format is usually one of two styles:
	// global: <functionName>_<argType>_<argTypeN>
	// member: <memberType>_<functionName>_<argType>_<argTypeN>
	ID string

	// ArgTypes contains the set of argument types expected by the overload.
	//
	// For member functions ArgTypes[0] represents the member operand type.
	ArgTypes []*Type

	// ResultType indicates the output type from calling the function.
	ResultType *Type

	// IsMemberFunction indicates whether the overload is a member function
	IsMemberFunction bool

	// Function implementation options. Optional, but encouraged.
	// UnaryOp is a function binding that takes a single argument.
	UnaryOp functions.UnaryOp
	// BinaryOp is a function binding that takes two arguments.
	BinaryOp functions.BinaryOp
	// FunctionOp is a catch-all for zero-arity and three-plus arity functions.
	FunctionOp functions.FunctionOp

	// NonStrict indicates that the function will accept error and unknown arguments as inputs.
	NonStrict bool

	// OperandTrait indicates whether the member argument should have a specific type-trait.
	//
	// This is useful for creating overloads which operate on a type-interface rather than a concrete type.
	OperandTrait int
}

// SignatureEquals determines whether the incoming overload declaration signature is equal to the current signature.
//
// Result type, operand trait, and strict-ness are not considered as part of signature equality.
func (o *OverloadDecl) SignatureEquals(other *OverloadDecl) bool {
	if o == other {
		return true
	}
	if o.ID != other.ID || o.IsMemberFunction != other.IsMemberFunction || len(o.ArgTypes) != len(other.ArgTypes) {
		return false
	}
	for i, at := range o.ArgTypes {
		oat := other.ArgTypes[i]
		if !at.IsType(oat) {
			return false
		}
	}
	return o.ResultType.IsType(other.ResultType)
}

// SignatureOverlaps indicates whether two functions have non-equal, but overloapping function signatures.
//
// For example, list(dyn) collides with list(string) since the 'dyn' type can contain a 'string' type.
func (o *OverloadDecl) SignatureOverlaps(other *OverloadDecl) bool {
	if o.IsMemberFunction != other.IsMemberFunction || len(o.ArgTypes) != len(other.ArgTypes) {
		return false
	}
	argsOverlap := true
	for i, argType := range o.ArgTypes {
		otherArgType := other.ArgTypes[i]
		argsOverlap = argsOverlap &&
			(argType.IsAssignableType(otherArgType) ||
				otherArgType.IsAssignableType(argType))
	}
	return argsOverlap
}

// hasBinding indicates whether the overload already has a definition.
func (o *OverloadDecl) hasBinding() bool {
	return o.UnaryOp != nil || o.BinaryOp != nil || o.FunctionOp != nil
}

// guardedUnaryOp creates an invocation guard around the provided unary operator, if one is defined.
func (o *OverloadDecl) guardedUnaryOp(funcName string, disableTypeGuards bool) functions.UnaryOp {
	if o.UnaryOp == nil {
		return nil
	}
	return func(arg ref.Val) ref.Val {
		if !o.matchesRuntimeUnarySignature(disableTypeGuards, arg) {
			return MaybeNoSuchOverload(funcName, arg)
		}
		return o.UnaryOp(arg)
	}
}

// guardedBinaryOp creates an invocation guard around the provided binary operator, if one is defined.
func (o *OverloadDecl) guardedBinaryOp(funcName string, disableTypeGuards bool) functions.BinaryOp {
	if o.BinaryOp == nil {
		return nil
	}
	return func(arg1, arg2 ref.Val) ref.Val {
		if !o.matchesRuntimeBinarySignature(disableTypeGuards, arg1, arg2) {
			return MaybeNoSuchOverload(funcName, arg1, arg2)
		}
		return o.BinaryOp(arg1, arg2)
	}
}

// guardedFunctionOp creates an invocation guard around the provided variadic function binding, if one is provided.
func (o *OverloadDecl) guardedFunctionOp(funcName string, disableTypeGuards bool) functions.FunctionOp {
	if o.FunctionOp == nil {
		return nil
	}
	return func(args ...ref.Val) ref.Val {
		if !o.matchesRuntimeSignature(disableTypeGuards, args...) {
			return MaybeNoSuchOverload(funcName, args...)
		}
		return o.FunctionOp(args...)
	}
}

// matchesRuntimeUnarySignature indicates whether the argument type is runtime assiganble to the overload's expected argument.
func (o *OverloadDecl) matchesRuntimeUnarySignature(disableTypeGuards bool, arg ref.Val) bool {
	return matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[0], arg) &&
		matchOperandTrait(o.OperandTrait, arg)
}

// matchesRuntimeBinarySignature indicates whether the argument types are runtime assiganble to the overload's expected arguments.
func (o *OverloadDecl) matchesRuntimeBinarySignature(disableTypeGuards bool, arg1, arg2 ref.Val) bool {
	return matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[0], arg1) &&
		matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[1], arg2) &&
		matchOperandTrait(o.OperandTrait, arg1)
}

// matchesRuntimeSignature indicates whether the argument types are runtime assiganble to the overload's expected arguments.
func (o *OverloadDecl) matchesRuntimeSignature(disableTypeGuards bool, args ...ref.Val) bool {
	if len(args) != len(o.ArgTypes) {
		return false
	}
	if len(args) == 0 {
		return true
	}
	for i, arg := range args {
		if !matchRuntimeArgType(o.NonStrict, disableTypeGuards, o.ArgTypes[i], arg) {
			return false
		}
	}
	return matchOperandTrait(o.OperandTrait, args[0])
}

func matchRuntimeArgType(nonStrict, disableTypeGuards bool, argType *Type, arg ref.Val) bool {
	if nonStrict && (disableTypeGuards || types.IsUnknownOrError(arg)) {
		return true
	}
	if types.IsUnknownOrError(arg) {
		return false
	}
	return disableTypeGuards || argType.IsAssignableRuntimeType(arg)
}

func matchOperandTrait(trait int, arg ref.Val) bool {
	return trait == 0 || arg.Type().HasTrait(trait) || types.IsUnknownOrError(arg)
}

// OverloadOpt is a functional option for configuring a function overload.
type OverloadOpt func(*OverloadDecl) (*OverloadDecl, error)

// UnaryBinding provides the implementation of a unary overload. The provided function is protected by a runtime
// type-guard which ensures runtime type agreement between the overload signature and runtime argument types.
func UnaryBinding(binding functions.UnaryOp) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		if o.hasBinding() {
			return nil, fmt.Errorf("overload already has a binding: %s", o.ID)
		}
		if len(o.ArgTypes) != 1 {
			return nil, fmt.Errorf("unary function bound to non-unary overload: %s", o.ID)
		}
		o.UnaryOp = binding
		return o, nil
	}
}

// BinaryBinding provides the implementation of a binary overload. The provided function is protected by a runtime
// type-guard which ensures runtime type agreement between the overload signature and runtime argument types.
func BinaryBinding(binding functions.BinaryOp) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		if o.hasBinding() {
			return nil, fmt.Errorf("overload already has a binding: %s", o.ID)
		}
		if len(o.ArgTypes) != 2 {
			return nil, fmt.Errorf("binary function bound to non-binary overload: %s", o.ID)
		}
		o.BinaryOp = binding
		return o, nil
	}
}

// FunctionBinding provides the implementation of a variadic overload. The provided function is protected by a runtime
// type-guard which ensures runtime type agreement between the overload signature and runtime argument types.
func FunctionBinding(binding functions.FunctionOp) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		if o.hasBinding() {
			return nil, fmt.Errorf("overload already has a binding: %s", o.ID)
		}
		o.FunctionOp = binding
		return o, nil
	}
}

// OverloadIsNonStrict enables the function to be called with error and unknown argument values.
//
// Note: do not use this option unless absoluately necessary as it should be an uncommon feature.
func OverloadIsNonStrict() OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		o.NonStrict = true
		return o, nil
	}
}

// OverloadOperandTrait configures a set of traits which the first argument to the overload must implement in order to be
// successfully invoked.
func OverloadOperandTrait(trait int) OverloadOpt {
	return func(o *OverloadDecl) (*OverloadDecl, error) {
		o.OperandTrait = trait
		return o, nil
	}
}

// NewVariable creates a new variable declaration.
func NewVariable(name string, t *Type) *VariableDecl {
	return &VariableDecl{Name: name, Type: t}
}

// VariableDecl defines a variable declaration which may optionally have a constant value.
type VariableDecl struct {
	Name string
	Type *Type
}

// DeclarationEquals returns true if one variable declaration has the same name and same type as the input.
func (v *VariableDecl) DeclarationEquals(other *VariableDecl) bool {
	if v == other {
		return true
	}
	return v.Name == other.Name && v.Type.IsType(other.Type)
}

// VariableDeclToExprDecl converts a go-native variable declaration into a protobuf-type variable declaration.
func VariableDeclToExprDecl(v *VariableDecl) (*exprpb.Decl, error) {
	varType, err := TypeToExprType(v.Type)
	if err != nil {
		return nil, err
	}
	return chkdecls.NewVar(v.Name, varType), nil
}

// FunctionDeclToExprDecl converts a go-native function declaration into a protobuf-typed function declaration.
func FunctionDeclToExprDecl(f *FunctionDecl) (*exprpb.Decl, error) {
	overloads := make([]*exprpb.Decl_FunctionDecl_Overload, len(f.Overloads))
	for i, oID := range f.overloadOrdinals {
		o := f.Overloads[oID]
		paramNames := map[string]struct{}{}
		argTypes := make([]*exprpb.Type, len(o.ArgTypes))
		for j, a := range o.ArgTypes {
			collectParamNames(paramNames, a)
			at, err := TypeToExprType(a)
			if err != nil {
				return nil, err
			}
			argTypes[j] = at
		}
		collectParamNames(paramNames, o.ResultType)
		resultType, err := TypeToExprType(o.ResultType)
		if err != nil {
			return nil, err
		}
		if len(paramNames) == 0 {
			if o.IsMemberFunction {
				overloads[i] = chkdecls.NewInstanceOverload(oID, argTypes, resultType)
			} else {
				overloads[i] = chkdecls.NewOverload(oID, argTypes, resultType)
			}
		} else {
			params := []string{}
			for pn := range paramNames {
				params = append(params, pn)
			}
			if o.IsMemberFunction {
				overloads[i] = chkdecls.NewParameterizedInstanceOverload(oID, argTypes, resultType, params)
			} else {
				overloads[i] = chkdecls.NewParameterizedOverload(oID, argTypes, resultType, params)
			}
		}
	}
	return chkdecls.NewFunction(f.Name, overloads...), nil
}

func collectParamNames(paramNames map[string]struct{}, arg *Type) {
	if arg.Kind == TypeParamKind {
		paramNames[arg.RuntimeTypeName()] = struct{}{}
	}
	for _, param := range arg.Parameters {
		collectParamNames(paramNames, param)
	}
}
