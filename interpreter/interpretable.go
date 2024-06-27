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

package interpreter

import (
	"context"
	"fmt"

	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Interpretable can accept a given Activation and produce a value along with
// an accompanying EvalState which can be used to inspect whether additional
// data might be necessary to complete the evaluation.
type Interpretable interface {
	// ID value corresponding to the expression node.
	ID() int64

	// Eval an Activation to produce an output.
	Eval(activation Activation) ref.Val
}

// InterpretableContext is an adaptation of Interpretable for context.Context.
type InterpretableContext interface {
	Interpretable
	// EvalContext an Activation to produce an output.
	EvalContext(ctx context.Context, activation Activation) ref.Val
}

// InterpretableConst interface for tracking whether the Interpretable is a constant value.
type InterpretableConst interface {
	Interpretable

	// Value returns the constant value of the instruction.
	Value() ref.Val
}

type InterpretableConstContext interface {
	InterpretableContext
	// Value returns the constant value of the instruction.
	Value() ref.Val
}

// InterpretableAttribute interface for tracking whether the Interpretable is an attribute.
type InterpretableAttribute interface {
	Interpretable

	// Attr returns the Attribute value.
	Attr() Attribute

	// Adapter returns the type adapter to be used for adapting resolved Attribute values.
	Adapter() types.Adapter

	// AddQualifier proxies the Attribute.AddQualifier method.
	//
	// Note, this method may mutate the current attribute state. If the desire is to clone the
	// Attribute, the Attribute should first be copied before adding the qualifier. Attributes
	// are not copyable by default, so this is a capable that would need to be added to the
	// AttributeFactory or specifically to the underlying Attribute implementation.
	AddQualifier(Qualifier) (Attribute, error)

	// Qualify replicates the Attribute.Qualify method to permit extension and interception
	// of object qualification.
	Qualify(vars Activation, obj any) (any, error)

	// QualifyIfPresent qualifies the object if the qualifier is declared or defined on the object.
	// The 'presenceOnly' flag indicates that the value is not necessary, just a boolean status as
	// to whether the qualifier is present.
	QualifyIfPresent(vars Activation, obj any, presenceOnly bool) (any, bool, error)

	// IsOptional indicates whether the resulting value is an optional type.
	IsOptional() bool

	// Resolve returns the value of the Attribute given the current Activation.
	Resolve(Activation) (any, error)
}

type InterpretableAttributeContext interface {
	InterpretableAttribute
	InterpretableContext

	// Attr returns the Attribute value.
	AttrContext(context.Context) AttributeContext

	// AddQualifierContext proxies the Attribute.AddQualifier method.
	//
	// Note, this method may mutate the current attribute state. If the desire is to clone the
	// Attribute, the Attribute should first be copied before adding the qualifier. Attributes
	// are not copyable by default, so this is a capable that would need to be added to the
	// AttributeFactory or specifically to the underlying Attribute implementation.
	AddQualifierContext(ctx context.Context, q Qualifier) (AttributeContext, error)

	// QualifyContext replicates the Attribute.Qualify method to permit extension and interception
	// of object qualification.
	QualifyContext(ctx context.Context, vars Activation, obj any) (any, error)

	// QualifyIfPresentContext qualifies the object if the qualifier is declared or defined on the object.
	// The 'presenceOnly' flag indicates that the value is not necessary, just a boolean status as
	// to whether the qualifier is present.
	QualifyIfPresentContext(ctx context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error)

	// ResolveContext returns the value of the Attribute given the current Activation.
	ResolveContext(ctx context.Context, vars Activation) (any, error)
}

// InterpretableCall interface for inspecting Interpretable instructions related to function calls.
type InterpretableCall interface {
	Interpretable

	// Function returns the function name as it appears in text or mangled operator name as it
	// appears in the operators.go file.
	Function() string

	// OverloadID returns the overload id associated with the function specialization.
	// Overload ids are stable across language boundaries and can be treated as synonymous with a
	// unique function signature.
	OverloadID() string

	// Args returns the normalized arguments to the function overload.
	// For receiver-style functions, the receiver target is arg 0.
	Args() []Interpretable
}

type InterpretableCallContext interface {
	InterpretableCall
	InterpretableContext
}

// InterpretableConstructor interface for inspecting  Interpretable instructions that initialize a list, map
// or struct.
type InterpretableConstructor interface {
	Interpretable

	// InitVals returns all the list elements, map key and values or struct field values.
	InitVals() []Interpretable

	// Type returns the type constructed.
	Type() ref.Type
}

type InterpretableConstructorContext interface {
	InterpretableContext
	InterpretableConstructor
}

// Core Interpretable implementations used during the program planning phase.

type evalTestOnly struct {
	id int64
	InterpretableAttributeContext
}

// ID implements the Interpretable interface method.
func (test *evalTestOnly) ID() int64 {
	return test.id
}

// Eval implements the Interpretable interface method.
func (test *evalTestOnly) Eval(vars Activation) ref.Val {
	return test.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (test *evalTestOnly) EvalContext(ctx context.Context, vars Activation) ref.Val {
	val, err := test.ResolveContext(ctx, vars)
	// Return an error if the resolve step fails
	if err != nil {
		return types.LabelErrNode(test.id, types.WrapErr(err))
	}
	if optVal, isOpt := val.(*types.Optional); isOpt {
		return types.Bool(optVal.HasValue())
	}
	return test.Adapter().NativeToValue(val)
}

// AddQualifier appends a qualifier that will always and only perform a presence test.
func (test *evalTestOnly) AddQualifier(q Qualifier) (Attribute, error) {
	return test.AddQualifierContext(context.Background(), q)
}

// AddQualifierContext appends a qualifier that will always and only perform a presence test.
func (test *evalTestOnly) AddQualifierContext(ctx context.Context, q Qualifier) (AttributeContext, error) {
	cq, ok := q.(ConstantQualifierContext)
	if !ok {
		return nil, fmt.Errorf("test only expressions must have constant qualifiers: %v", q)
	}
	return test.InterpretableAttributeContext.AddQualifierContext(ctx, &testOnlyQualifier{ConstantQualifierContext: cq})
}

type testOnlyQualifier struct {
	ConstantQualifierContext
}

// Qualify determines whether the test-only qualifier is present on the input object.
func (q *testOnlyQualifier) Qualify(vars Activation, obj any) (any, error) {
	return q.QualifyContext(context.Background(), vars, obj)
}

// QualifyContext determines whether the test-only qualifier is present on the input object.
func (q *testOnlyQualifier) QualifyContext(ctx context.Context, vars Activation, obj any) (any, error) {
	out, present, err := q.ConstantQualifierContext.QualifyIfPresentContext(ctx, vars, obj, true)
	if err != nil {
		return nil, err
	}
	if unk, isUnk := out.(types.Unknown); isUnk {
		return unk, nil
	}
	if opt, isOpt := out.(types.Optional); isOpt {
		return opt.HasValue(), nil
	}
	return present, nil
}

// QualifyIfPresent returns whether the target field in the test-only expression is present.
func (q *testOnlyQualifier) QualifyIfPresent(vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return q.QualifyIfPresentContext(context.Background(), vars, obj, presenceOnly)
}

// QualifyIfPresentContext returns whether the target field in the test-only expression is present.
func (q *testOnlyQualifier) QualifyIfPresentContext(ctx context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	// Only ever test for presence.
	return q.ConstantQualifierContext.QualifyIfPresentContext(ctx, vars, obj, true)
}

// QualifierValueEquals determines whether the test-only constant qualifier equals the input value.
func (q *testOnlyQualifier) QualifierValueEquals(value any) bool {
	// The input qualifier will always be of type string
	return q.ConstantQualifierContext.Value().Value() == value
}

// NewConstValue creates a new constant valued Interpretable.
func NewConstValue(id int64, val ref.Val) InterpretableConst {
	return &evalConst{
		id:  id,
		val: val,
	}
}

// NewConstValueContext creates a new constant valued InterpretableContext.
func NewConstValueContext(id int64, val ref.Val) InterpretableConstContext {
	return &evalConst{
		id:  id,
		val: val,
	}
}

type evalConst struct {
	id  int64
	val ref.Val
}

// ID implements the Interpretable interface method.
func (cons *evalConst) ID() int64 {
	return cons.id
}

// Eval implements the Interpretable interface method.
func (cons *evalConst) Eval(vars Activation) ref.Val {
	return cons.EvalContext(context.Background(), vars)
}

// EvalContext implements the Interpretable interface method.
func (cons *evalConst) EvalContext(ctx context.Context, vars Activation) ref.Val {
	return cons.val
}

// Value implements the InterpretableConst interface method.
func (cons *evalConst) Value() ref.Val {
	return cons.val
}

type evalOr struct {
	id    int64
	terms []InterpretableContext
}

// ID implements the Interpretable interface method.
func (or *evalOr) ID() int64 {
	return or.id
}

// Eval implements the Interpretable interface method.
func (or *evalOr) Eval(vars Activation) ref.Val {
	return or.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (or *evalOr) EvalContext(ctx context.Context, vars Activation) ref.Val {
	var err ref.Val = nil
	var unk *types.Unknown
	for _, term := range or.terms {
		val := term.EvalContext(ctx, vars)
		boolVal, ok := val.(types.Bool)
		// short-circuit on true.
		if ok && boolVal == types.True {
			return types.True
		}
		if !ok {
			isUnk := false
			unk, isUnk = types.MaybeMergeUnknowns(val, unk)
			if !isUnk && err == nil {
				if types.IsError(val) {
					err = val
				} else {
					err = types.MaybeNoSuchOverloadErr(val)
				}
				err = types.LabelErrNode(or.id, err)
			}
		}
	}
	if unk != nil {
		return unk
	}
	if err != nil {
		return err
	}
	return types.False
}

type evalAnd struct {
	id    int64
	terms []InterpretableContext
}

// ID implements the Interpretable interface method.
func (and *evalAnd) ID() int64 {
	return and.id
}

func (and *evalAnd) Eval(vars Activation) ref.Val {
	return and.EvalContext(context.Background(), vars)
}

// EvalContext implements the Interpretable interface method.
func (and *evalAnd) EvalContext(ctx context.Context, vars Activation) ref.Val {
	var err ref.Val = nil
	var unk *types.Unknown
	for _, term := range and.terms {
		val := term.EvalContext(ctx, vars)
		boolVal, ok := val.(types.Bool)
		// short-circuit on false.
		if ok && boolVal == types.False {
			return types.False
		}
		if !ok {
			isUnk := false
			unk, isUnk = types.MaybeMergeUnknowns(val, unk)
			if !isUnk && err == nil {
				if types.IsError(val) {
					err = val
				} else {
					err = types.MaybeNoSuchOverloadErr(val)
				}
				err = types.LabelErrNode(and.id, err)
			}
		}
	}
	if unk != nil {
		return unk
	}
	if err != nil {
		return err
	}
	return types.True
}

type evalEq struct {
	id  int64
	lhs InterpretableContext
	rhs InterpretableContext
}

// ID implements the Interpretable interface method.
func (eq *evalEq) ID() int64 {
	return eq.id
}

// Eval implements the Interpretable interface method.
func (eq *evalEq) Eval(vars Activation) ref.Val {
	return eq.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (eq *evalEq) EvalContext(ctx context.Context, vars Activation) ref.Val {
	lVal := eq.lhs.EvalContext(ctx, vars)
	rVal := eq.rhs.EvalContext(ctx, vars)
	if types.IsUnknownOrError(lVal) {
		return lVal
	}
	if types.IsUnknownOrError(rVal) {
		return rVal
	}
	return types.Equal(lVal, rVal)
}

// Function implements the InterpretableCall interface method.
func (*evalEq) Function() string {
	return operators.Equals
}

// OverloadID implements the InterpretableCall interface method.
func (*evalEq) OverloadID() string {
	return overloads.Equals
}

// Args implements the InterpretableCall interface method.
func (eq *evalEq) Args() []Interpretable {
	return []Interpretable{eq.lhs, eq.rhs}
}

type evalNe struct {
	id  int64
	lhs InterpretableContext
	rhs InterpretableContext
}

// ID implements the Interpretable interface method.
func (ne *evalNe) ID() int64 {
	return ne.id
}

// Eval implements the Interpretable interface method.
func (ne *evalNe) Eval(vars Activation) ref.Val {
	return ne.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (ne *evalNe) EvalContext(ctx context.Context, vars Activation) ref.Val {
	lVal := ne.lhs.EvalContext(ctx, vars)
	rVal := ne.rhs.EvalContext(ctx, vars)
	if types.IsUnknownOrError(lVal) {
		return lVal
	}
	if types.IsUnknownOrError(rVal) {
		return rVal
	}
	return types.Bool(types.Equal(lVal, rVal) != types.True)
}

// Function implements the InterpretableCall interface method.
func (*evalNe) Function() string {
	return operators.NotEquals
}

// OverloadID implements the InterpretableCall interface method.
func (*evalNe) OverloadID() string {
	return overloads.NotEquals
}

// Args implements the InterpretableCall interface method.
func (ne *evalNe) Args() []Interpretable {
	return []Interpretable{ne.lhs, ne.rhs}
}

type evalZeroArity struct {
	id       int64
	function string
	overload string
	impl     functions.FunctionOpContext
}

// ID implements the Interpretable interface method.
func (zero *evalZeroArity) ID() int64 {
	return zero.id
}

// Eval implements the Interpretable interface method.
func (zero *evalZeroArity) Eval(vars Activation) ref.Val {
	return zero.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (zero *evalZeroArity) EvalContext(ctx context.Context, vars Activation) ref.Val {
	return types.LabelErrNode(zero.id, zero.impl(ctx))
}

// Function implements the InterpretableCall interface method.
func (zero *evalZeroArity) Function() string {
	return zero.function
}

// OverloadID implements the InterpretableCall interface method.
func (zero *evalZeroArity) OverloadID() string {
	return zero.overload
}

// Args returns the argument to the unary function.
func (zero *evalZeroArity) Args() []Interpretable {
	return []Interpretable{}
}

type evalUnary struct {
	id        int64
	function  string
	overload  string
	arg       InterpretableContext
	trait     int
	impl      functions.UnaryOpContext
	nonStrict bool
}

// ID implements the Interpretable interface method.
func (un *evalUnary) ID() int64 {
	return un.id
}

// Eval implements the Interpretable interface method.
func (un *evalUnary) Eval(vars Activation) ref.Val {
	return un.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (un *evalUnary) EvalContext(ctx context.Context, vars Activation) ref.Val {
	argVal := un.arg.EvalContext(ctx, vars)
	// Early return if the argument to the function is unknown or error.
	strict := !un.nonStrict
	if strict && types.IsUnknownOrError(argVal) {
		return argVal
	}
	// If the implementation is bound and the argument value has the right traits required to
	// invoke it, then call the implementation.
	if un.impl != nil && (un.trait == 0 || (!strict && types.IsUnknownOrError(argVal)) || argVal.Type().HasTrait(un.trait)) {
		return types.LabelErrNode(un.id, un.impl(ctx, argVal))
	}
	// Otherwise, if the argument is a ReceiverType attempt to invoke the receiver method on the
	// operand (arg0).
	if argVal.Type().HasTrait(traits.ReceiverType) {
		return types.LabelErrNode(un.id, argVal.(traits.Receiver).Receive(un.function, un.overload, []ref.Val{}))
	}
	return types.NewErrWithNodeID(un.id, "no such overload: %s", un.function)
}

// Function implements the InterpretableCall interface method.
func (un *evalUnary) Function() string {
	return un.function
}

// OverloadID implements the InterpretableCall interface method.
func (un *evalUnary) OverloadID() string {
	return un.overload
}

// Args returns the argument to the unary function.
func (un *evalUnary) Args() []Interpretable {
	return []Interpretable{un.arg}
}

type evalBinary struct {
	id        int64
	function  string
	overload  string
	lhs       InterpretableContext
	rhs       InterpretableContext
	trait     int
	impl      functions.BinaryOpContext
	nonStrict bool
}

// ID implements the Interpretable interface method.
func (bin *evalBinary) ID() int64 {
	return bin.id
}

// Eval implements the Interpretable interface method.
func (bin *evalBinary) Eval(vars Activation) ref.Val {
	return bin.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (bin *evalBinary) EvalContext(ctx context.Context, vars Activation) ref.Val {
	lVal := bin.lhs.EvalContext(ctx, vars)
	rVal := bin.rhs.EvalContext(ctx, vars)
	// Early return if any argument to the function is unknown or error.
	strict := !bin.nonStrict
	if strict {
		if types.IsUnknownOrError(lVal) {
			return lVal
		}
		if types.IsUnknownOrError(rVal) {
			return rVal
		}
	}
	// If the implementation is bound and the argument value has the right traits required to
	// invoke it, then call the implementation.
	if bin.impl != nil && (bin.trait == 0 || (!strict && types.IsUnknownOrError(lVal)) || lVal.Type().HasTrait(bin.trait)) {
		return types.LabelErrNode(bin.id, bin.impl(ctx, lVal, rVal))
	}
	// Otherwise, if the argument is a ReceiverType attempt to invoke the receiver method on the
	// operand (arg0).
	if lVal.Type().HasTrait(traits.ReceiverType) {
		return types.LabelErrNode(bin.id, lVal.(traits.Receiver).Receive(bin.function, bin.overload, []ref.Val{rVal}))
	}
	return types.NewErrWithNodeID(bin.id, "no such overload: %s", bin.function)
}

// Function implements the InterpretableCall interface method.
func (bin *evalBinary) Function() string {
	return bin.function
}

// OverloadID implements the InterpretableCall interface method.
func (bin *evalBinary) OverloadID() string {
	return bin.overload
}

// Args returns the argument to the unary function.
func (bin *evalBinary) Args() []Interpretable {
	return []Interpretable{bin.lhs, bin.rhs}
}

type evalVarArgs struct {
	id        int64
	function  string
	overload  string
	args      []InterpretableContext
	trait     int
	impl      functions.FunctionOpContext
	nonStrict bool
}

// NewCall creates a new call Interpretable.
func NewCall(id int64, function, overload string, args []Interpretable, impl functions.FunctionOp) InterpretableCall {
	newArgs := make([]InterpretableContext, len(args))
	for idx, arg := range args {
		newArgs[idx] = ToInterpretableContext(arg)
	}
	return NewCallContext(id, function, overload, newArgs, func(_ context.Context, values ...ref.Val) ref.Val {
		return impl(values...)
	})
}

// NewCallContext creates a new call InterpretableContext.
func NewCallContext(id int64, function, overload string, args []InterpretableContext, impl functions.FunctionOpContext) InterpretableCall {
	return &evalVarArgs{
		id:       id,
		function: function,
		overload: overload,
		args:     args,
		impl:     impl,
	}
}

// ID implements the Interpretable interface method.
func (fn *evalVarArgs) ID() int64 {
	return fn.id
}

// Eval implements the Interpretable interface method.
func (fn *evalVarArgs) Eval(vars Activation) ref.Val {
	return fn.EvalContext(context.Background(), vars)
}

// EvalContext implements the Interpretable interface method.
func (fn *evalVarArgs) EvalContext(ctx context.Context, vars Activation) ref.Val {
	argVals := make([]ref.Val, len(fn.args))
	// Early return if any argument to the function is unknown or error.
	strict := !fn.nonStrict
	for i, arg := range fn.args {
		argVals[i] = arg.EvalContext(ctx, vars)
		if strict && types.IsUnknownOrError(argVals[i]) {
			return argVals[i]
		}
	}
	// If the implementation is bound and the argument value has the right traits required to
	// invoke it, then call the implementation.
	arg0 := argVals[0]
	if fn.impl != nil && (fn.trait == 0 || (!strict && types.IsUnknownOrError(arg0)) || arg0.Type().HasTrait(fn.trait)) {
		return types.LabelErrNode(fn.id, fn.impl(ctx, argVals...))
	}
	// Otherwise, if the argument is a ReceiverType attempt to invoke the receiver method on the
	// operand (arg0).
	if arg0.Type().HasTrait(traits.ReceiverType) {
		return types.LabelErrNode(fn.id, arg0.(traits.Receiver).Receive(fn.function, fn.overload, argVals[1:]))
	}
	return types.NewErrWithNodeID(fn.id, "no such overload: %s %d", fn.function, fn.id)
}

// Function implements the InterpretableCall interface method.
func (fn *evalVarArgs) Function() string {
	return fn.function
}

// OverloadID implements the InterpretableCall interface method.
func (fn *evalVarArgs) OverloadID() string {
	return fn.overload
}

// Args returns the argument to the unary function.
func (fn *evalVarArgs) Args() []Interpretable {
	ret := make([]Interpretable, len(fn.args))
	for idx, arg := range fn.args {
		ret[idx] = arg.(Interpretable)
	}
	return ret
}

type evalList struct {
	id           int64
	elems        []InterpretableContext
	optionals    []bool
	hasOptionals bool
	adapter      types.Adapter
}

// ID implements the Interpretable interface method.
func (l *evalList) ID() int64 {
	return l.id
}

// Eval implements the Interpretable interface method.
func (l *evalList) Eval(vars Activation) ref.Val {
	return l.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (l *evalList) EvalContext(ctx context.Context, vars Activation) ref.Val {
	elemVals := make([]ref.Val, 0, len(l.elems))
	// If any argument is unknown or error early terminate.
	for i, elem := range l.elems {
		elemVal := elem.EvalContext(ctx, vars)
		if types.IsUnknownOrError(elemVal) {
			return elemVal
		}
		if l.hasOptionals && l.optionals[i] {
			optVal, ok := elemVal.(*types.Optional)
			if !ok {
				return types.LabelErrNode(l.id, invalidOptionalElementInit(elemVal))
			}
			if !optVal.HasValue() {
				continue
			}
			elemVal = optVal.GetValue()
		}
		elemVals = append(elemVals, elemVal)
	}
	return l.adapter.NativeToValue(elemVals)
}

func (l *evalList) InitVals() []Interpretable {
	ret := make([]Interpretable, len(l.elems))
	for idx, elem := range l.elems {
		ret[idx] = elem.(Interpretable)
	}
	return ret
}

func (l *evalList) Type() ref.Type {
	return types.ListType
}

type evalMap struct {
	id           int64
	keys         []InterpretableContext
	vals         []InterpretableContext
	optionals    []bool
	hasOptionals bool
	adapter      types.Adapter
}

// ID implements the Interpretable interface method.
func (m *evalMap) ID() int64 {
	return m.id
}

// Eval implements the Interpretable interface method.
func (m *evalMap) Eval(vars Activation) ref.Val {
	return m.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (m *evalMap) EvalContext(ctx context.Context, vars Activation) ref.Val {
	entries := make(map[ref.Val]ref.Val)
	// If any argument is unknown or error early terminate.
	for i, key := range m.keys {
		keyVal := key.EvalContext(ctx, vars)
		if types.IsUnknownOrError(keyVal) {
			return keyVal
		}
		valVal := m.vals[i].EvalContext(ctx, vars)
		if types.IsUnknownOrError(valVal) {
			return valVal
		}
		if m.hasOptionals && m.optionals[i] {
			optVal, ok := valVal.(*types.Optional)
			if !ok {
				return types.LabelErrNode(m.id, invalidOptionalEntryInit(keyVal, valVal))
			}
			if !optVal.HasValue() {
				delete(entries, keyVal)
				continue
			}
			valVal = optVal.GetValue()
		}
		entries[keyVal] = valVal
	}
	return m.adapter.NativeToValue(entries)
}

func (m *evalMap) InitVals() []Interpretable {
	if len(m.keys) != len(m.vals) {
		return nil
	}
	result := make([]Interpretable, len(m.keys)+len(m.vals))
	idx := 0
	for i, k := range m.keys {
		v := m.vals[i].(Interpretable)
		result[idx] = k.(Interpretable)
		idx++
		result[idx] = v
		idx++
	}
	return result
}

func (m *evalMap) Type() ref.Type {
	return types.MapType
}

type evalObj struct {
	id           int64
	typeName     string
	fields       []string
	vals         []InterpretableContext
	optionals    []bool
	hasOptionals bool
	provider     types.Provider
}

// ID implements the Interpretable interface method.
func (o *evalObj) ID() int64 {
	return o.id
}

// Eval implements the Interpretable interface method.
func (o *evalObj) Eval(vars Activation) ref.Val {
	return o.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (o *evalObj) EvalContext(ctx context.Context, vars Activation) ref.Val {
	fieldVals := make(map[string]ref.Val)
	// If any argument is unknown or error early terminate.
	for i, field := range o.fields {
		val := o.vals[i].EvalContext(ctx, vars)
		if types.IsUnknownOrError(val) {
			return val
		}
		if o.hasOptionals && o.optionals[i] {
			optVal, ok := val.(*types.Optional)
			if !ok {
				return types.LabelErrNode(o.id, invalidOptionalEntryInit(field, val))
			}
			if !optVal.HasValue() {
				delete(fieldVals, field)
				continue
			}
			val = optVal.GetValue()
		}
		fieldVals[field] = val
	}
	return types.LabelErrNode(o.id, o.provider.NewValue(o.typeName, fieldVals))
}

func (o *evalObj) InitVals() []Interpretable {
	ret := make([]Interpretable, len(o.vals))
	for idx, val := range o.vals {
		ret[idx] = val.(Interpretable)
	}
	return ret
}

func (o *evalObj) Type() ref.Type {
	return types.NewObjectTypeValue(o.typeName)
}

type evalFold struct {
	id            int64
	accuVar       string
	iterVar       string
	iterRange     InterpretableContext
	accu          InterpretableContext
	cond          InterpretableContext
	step          InterpretableContext
	result        InterpretableContext
	adapter       types.Adapter
	exhaustive    bool
	interruptable bool
}

// ID implements the Interpretable interface method.
func (fold *evalFold) ID() int64 {
	return fold.id
}

// Eval implements the Interpretable interface method.
func (fold *evalFold) Eval(vars Activation) ref.Val {
	return fold.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (fold *evalFold) EvalContext(ctx context.Context, vars Activation) ref.Val {
	foldRange := fold.iterRange.EvalContext(ctx, vars)
	if !foldRange.Type().HasTrait(traits.IterableType) {
		return types.ValOrErr(foldRange, "got '%T', expected iterable type", foldRange)
	}
	// Configure the fold activation with the accumulator initial value.
	accuCtx := varActivationPool.Get().(*varActivation)
	accuCtx.parent = vars
	accuCtx.name = fold.accuVar
	accuCtx.val = fold.accu.EvalContext(ctx, vars)
	// If the accumulator starts as an empty list, then the comprehension will build a list
	// so create a mutable list to optimize the cost of the inner loop.
	l, ok := accuCtx.val.(traits.Lister)
	buildingList := false
	if !fold.exhaustive && ok && l.Size() == types.IntZero {
		buildingList = true
		accuCtx.val = types.NewMutableList(fold.adapter)
	}
	iterCtx := varActivationPool.Get().(*varActivation)
	iterCtx.parent = accuCtx
	iterCtx.name = fold.iterVar

	interrupted := false
	it := foldRange.(traits.Iterable).Iterator()
	for it.HasNext() == types.True {
		// Modify the iter var in the fold activation.
		iterCtx.val = it.Next()

		// Evaluate the condition, terminate the loop if false.
		cond := fold.cond.EvalContext(ctx, iterCtx)
		condBool, ok := cond.(types.Bool)
		if !fold.exhaustive && ok && condBool != types.True {
			break
		}
		// Evaluate the evaluation step into accu var.
		accuCtx.val = fold.step.EvalContext(ctx, iterCtx)
		if fold.interruptable {
			if stop, found := vars.ResolveName("#interrupted"); found && stop == true {
				interrupted = true
				break
			}
		}
	}
	varActivationPool.Put(iterCtx)
	if interrupted {
		varActivationPool.Put(accuCtx)
		return types.NewErr("operation interrupted")
	}

	// Compute the result.
	res := fold.result.EvalContext(ctx, accuCtx)
	varActivationPool.Put(accuCtx)
	// Convert a mutable list to an immutable one, if the comprehension has generated a list as a result.
	if !types.IsUnknownOrError(res) && buildingList {
		if _, ok := res.(traits.MutableLister); ok {
			res = res.(traits.MutableLister).ToImmutableList()
		}
	}
	return res
}

// Optional Interpretable implementations that specialize, subsume, or extend the core evaluation
// plan via decorators.

// evalSetMembership is an Interpretable implementation which tests whether an input value
// exists within the set of map keys used to model a set.
type evalSetMembership struct {
	inst     InterpretableContext
	arg      InterpretableContext
	valueSet map[ref.Val]ref.Val
}

// ID implements the Interpretable interface method.
func (e *evalSetMembership) ID() int64 {
	return e.inst.ID()
}

// Eval implements the Interpretable interface method.
func (e *evalSetMembership) Eval(vars Activation) ref.Val {
	return e.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (e *evalSetMembership) EvalContext(ctx context.Context, vars Activation) ref.Val {
	val := e.arg.EvalContext(ctx, vars)
	if types.IsUnknownOrError(val) {
		return val
	}
	if ret, found := e.valueSet[val]; found {
		return ret
	}
	return types.False
}

// evalWatch is an Interpretable implementation that wraps the execution of a given
// expression so that it may observe the computed value and send it to an observer.
type evalWatch struct {
	InterpretableContext
	observer EvalObserver
}

// Eval implements the Interpretable interface method.
func (e *evalWatch) Eval(vars Activation) ref.Val {
	return e.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (e *evalWatch) EvalContext(ctx context.Context, vars Activation) ref.Val {
	val := e.InterpretableContext.EvalContext(ctx, vars)
	e.observer(e.ID(), e.InterpretableContext, val)
	return val
}

// evalWatchAttr describes a watcher of an InterpretableAttribute Interpretable.
//
// Since the watcher may be selected against at a later stage in program planning, the watcher
// must implement the InterpretableAttribute interface by proxy.
type evalWatchAttr struct {
	InterpretableAttributeContext
	observer EvalObserver
}

// AddQualifier creates a wrapper over the incoming qualifier which observes the qualification
// result.
func (e *evalWatchAttr) AddQualifier(q Qualifier) (Attribute, error) {
	return e.AddQualifierContext(context.Background(), q)
}

// AddQualifierContext creates a wrapper over the incoming qualifier which observes the qualification
// result.
func (e *evalWatchAttr) AddQualifierContext(ctx context.Context, q Qualifier) (AttributeContext, error) {
	switch qual := q.(type) {
	// By default, the qualifier is either a constant or an attribute
	// There may be some custom cases where the attribute is neither.
	case ConstantQualifierContext:
		// Expose a method to test whether the qualifier matches the input pattern.
		q = &evalWatchConstQual{
			ConstantQualifierContext: qual,
			observer:                 e.observer,
			adapter:                  e.Adapter(),
		}
	case *evalWatchAttr:
		// Unwrap the evalWatchAttr since the observation will be applied during Qualify or
		// QualifyIfPresent rather than Eval.
		q = &evalWatchAttrQual{
			Attribute: qual.InterpretableAttributeContext,
			observer:  e.observer,
			adapter:   e.Adapter(),
		}
	case AttributeContext:
		// Expose methods which intercept the qualification prior to being applied as a qualifier.
		// Using this interface ensures that the qualifier is converted to a constant value one
		// time during attribute pattern matching as the method embeds the Attribute interface
		// needed to trip the conversion to a constant.
		q = &evalWatchAttrQual{
			Attribute: qual,
			observer:  e.observer,
			adapter:   e.Adapter(),
		}
	case Attribute:
		// Expose methods which intercept the qualification prior to being applied as a qualifier.
		// Using this interface ensures that the qualifier is converted to a constant value one
		// time during attribute pattern matching as the method embeds the Attribute interface
		// needed to trip the conversion to a constant.
		q = &evalWatchAttrQual{
			Attribute: qual,
			observer:  e.observer,
			adapter:   e.Adapter(),
		}
	default:
		// This is likely a custom qualifier type.
		q = &evalWatchQual{
			Qualifier: qual,
			observer:  e.observer,
			adapter:   e.Adapter(),
		}
	}
	_, err := e.InterpretableAttributeContext.AddQualifierContext(ctx, q)
	return e, err
}

// Eval implements the Interpretable interface method.
func (e *evalWatchAttr) Eval(vars Activation) ref.Val {
	return e.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (e *evalWatchAttr) EvalContext(ctx context.Context, vars Activation) ref.Val {
	val := e.InterpretableAttributeContext.EvalContext(ctx, vars)
	e.observer(e.ID(), e.InterpretableAttributeContext, val)
	return val
}

// evalWatchConstQual observes the qualification of an object using a constant boolean, int,
// string, or uint.
type evalWatchConstQual struct {
	ConstantQualifierContext
	observer EvalObserver
	adapter  types.Adapter
}

// Qualify observes the qualification of a object via a constant boolean, int, string, or uint.
func (e *evalWatchConstQual) Qualify(vars Activation, obj any) (any, error) {
	return e.QualifyContext(context.Background(), vars, obj)
}

// QualifyContext observes the qualification of a object via a constant boolean, int, string, or uint.
func (e *evalWatchConstQual) QualifyContext(ctx context.Context, vars Activation, obj any) (any, error) {
	out, err := e.ConstantQualifierContext.QualifyContext(ctx, vars, obj)
	var val ref.Val
	if err != nil {
		val = types.LabelErrNode(e.ID(), types.WrapErr(err))
	} else {
		val = e.adapter.NativeToValue(out)
	}
	e.observer(e.ID(), e.ConstantQualifierContext, val)
	return out, err
}

// QualifyIfPresent conditionally qualifies the variable and only records a value if one is present.
func (e *evalWatchConstQual) QualifyIfPresent(vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return e.QualifyIfPresentContext(context.Background(), vars, obj, presenceOnly)
}

// QualifyIfPresentContext conditionally qualifies the variable and only records a value if one is present.
func (e *evalWatchConstQual) QualifyIfPresentContext(ctx context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	out, present, err := e.ConstantQualifierContext.QualifyIfPresentContext(ctx, vars, obj, presenceOnly)
	var val ref.Val
	if err != nil {
		val = types.LabelErrNode(e.ID(), types.WrapErr(err))
	} else if out != nil {
		val = e.adapter.NativeToValue(out)
	} else if presenceOnly {
		val = types.Bool(present)
	}
	if present || presenceOnly {
		e.observer(e.ID(), e.ConstantQualifierContext, val)
	}
	return out, present, err
}

// QualifierValueEquals tests whether the incoming value is equal to the qualifying constant.
func (e *evalWatchConstQual) QualifierValueEquals(value any) bool {
	qve, ok := e.ConstantQualifierContext.(qualifierValueEquator)
	return ok && qve.QualifierValueEquals(value)
}

// evalWatchAttrQual observes the qualification of an object by a value computed at runtime.
type evalWatchAttrQual struct {
	Attribute
	observer EvalObserver
	adapter  ref.TypeAdapter
}

// Qualify observes the qualification of a object via a value computed at runtime.
func (e *evalWatchAttrQual) Qualify(vars Activation, obj any) (any, error) {
	out, err := e.Attribute.Qualify(vars, obj)
	var val ref.Val
	if err != nil {
		val = types.LabelErrNode(e.ID(), types.WrapErr(err))
	} else {
		val = e.adapter.NativeToValue(out)
	}
	e.observer(e.ID(), e.Attribute, val)
	return out, err
}

// QualifyIfPresent conditionally qualifies the variable and only records a value if one is present.
func (e *evalWatchAttrQual) QualifyIfPresent(vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	out, present, err := e.Attribute.QualifyIfPresent(vars, obj, presenceOnly)
	var val ref.Val
	if err != nil {
		val = types.LabelErrNode(e.ID(), types.WrapErr(err))
	} else if out != nil {
		val = e.adapter.NativeToValue(out)
	} else if presenceOnly {
		val = types.Bool(present)
	}
	if present || presenceOnly {
		e.observer(e.ID(), e.Attribute, val)
	}
	return out, present, err
}

// evalWatchQual observes the qualification of an object by a value computed at runtime.
type evalWatchQual struct {
	Qualifier
	observer EvalObserver
	adapter  types.Adapter
}

// Qualify observes the qualification of a object via a value computed at runtime.
func (e *evalWatchQual) Qualify(vars Activation, obj any) (any, error) {
	out, err := e.Qualifier.Qualify(vars, obj)
	var val ref.Val
	if err != nil {
		val = types.LabelErrNode(e.ID(), types.WrapErr(err))
	} else {
		val = e.adapter.NativeToValue(out)
	}
	e.observer(e.ID(), e.Qualifier, val)
	return out, err
}

// QualifyIfPresent conditionally qualifies the variable and only records a value if one is present.
func (e *evalWatchQual) QualifyIfPresent(vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	out, present, err := e.Qualifier.QualifyIfPresent(vars, obj, presenceOnly)
	var val ref.Val
	if err != nil {
		val = types.LabelErrNode(e.ID(), types.WrapErr(err))
	} else if out != nil {
		val = e.adapter.NativeToValue(out)
	} else if presenceOnly {
		val = types.Bool(present)
	}
	if present || presenceOnly {
		e.observer(e.ID(), e.Qualifier, val)
	}
	return out, present, err
}

// evalWatchConst describes a watcher of an instConst Interpretable.
type evalWatchConst struct {
	InterpretableConst
	observer EvalObserver
}

// Eval implements the Interpretable interface method.
func (e *evalWatchConst) Eval(vars Activation) ref.Val {
	return e.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (e *evalWatchConst) EvalContext(ctx context.Context, vars Activation) ref.Val {
	val := e.Value()
	e.observer(e.ID(), e.InterpretableConst, val)
	return val
}

// evalExhaustiveOr is just like evalOr, but does not short-circuit argument evaluation.
type evalExhaustiveOr struct {
	id    int64
	terms []InterpretableContext
}

// ID implements the Interpretable interface method.
func (or *evalExhaustiveOr) ID() int64 {
	return or.id
}

// Eval implements the Interpretable interface method.
func (or *evalExhaustiveOr) Eval(vars Activation) ref.Val {
	return or.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (or *evalExhaustiveOr) EvalContext(ctx context.Context, vars Activation) ref.Val {
	var err ref.Val = nil
	var unk *types.Unknown
	isTrue := false
	for _, term := range or.terms {
		val := term.EvalContext(ctx, vars)
		boolVal, ok := val.(types.Bool)
		// flag the result as true
		if ok && boolVal == types.True {
			isTrue = true
		}
		if !ok && !isTrue {
			isUnk := false
			unk, isUnk = types.MaybeMergeUnknowns(val, unk)
			if !isUnk && err == nil {
				if types.IsError(val) {
					err = val
				} else {
					err = types.MaybeNoSuchOverloadErr(val)
				}
			}
		}
	}
	if isTrue {
		return types.True
	}
	if unk != nil {
		return unk
	}
	if err != nil {
		return err
	}
	return types.False
}

// evalExhaustiveAnd is just like evalAnd, but does not short-circuit argument evaluation.
type evalExhaustiveAnd struct {
	id    int64
	terms []InterpretableContext
}

// ID implements the Interpretable interface method.
func (and *evalExhaustiveAnd) ID() int64 {
	return and.id
}

// Eval implements the Interpretable interface method.
func (and *evalExhaustiveAnd) Eval(vars Activation) ref.Val {
	return and.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (and *evalExhaustiveAnd) EvalContext(ctx context.Context, vars Activation) ref.Val {
	var err ref.Val = nil
	var unk *types.Unknown
	isFalse := false
	for _, term := range and.terms {
		val := term.EvalContext(ctx, vars)
		boolVal, ok := val.(types.Bool)
		// short-circuit on false.
		if ok && boolVal == types.False {
			isFalse = true
		}
		if !ok && !isFalse {
			isUnk := false
			unk, isUnk = types.MaybeMergeUnknowns(val, unk)
			if !isUnk && err == nil {
				if types.IsError(val) {
					err = val
				} else {
					err = types.MaybeNoSuchOverloadErr(val)
				}
			}
		}
	}
	if isFalse {
		return types.False
	}
	if unk != nil {
		return unk
	}
	if err != nil {
		return err
	}
	return types.True
}

// evalExhaustiveConditional is like evalConditional, but does not short-circuit argument
// evaluation.
type evalExhaustiveConditional struct {
	id      int64
	adapter types.Adapter
	attr    *conditionalAttribute
}

// ID implements the Interpretable interface method.
func (cond *evalExhaustiveConditional) ID() int64 {
	return cond.id
}

// Eval implements the Interpretable interface method.
func (cond *evalExhaustiveConditional) Eval(vars Activation) ref.Val {
	return cond.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (cond *evalExhaustiveConditional) EvalContext(ctx context.Context, vars Activation) ref.Val {
	cVal := cond.attr.expr.EvalContext(ctx, vars)
	tVal, tErr := cond.attr.truthy.ResolveContext(ctx, vars)
	fVal, fErr := cond.attr.falsy.ResolveContext(ctx, vars)
	cBool, ok := cVal.(types.Bool)
	if !ok {
		return types.ValOrErr(cVal, "no such overload")
	}
	if cBool {
		if tErr != nil {
			return types.LabelErrNode(cond.id, types.WrapErr(tErr))
		}
		return cond.adapter.NativeToValue(tVal)
	}
	if fErr != nil {
		return types.LabelErrNode(cond.id, types.WrapErr(fErr))
	}
	return cond.adapter.NativeToValue(fVal)
}

// evalAttr evaluates an Attribute value.
type evalAttr struct {
	adapter  types.Adapter
	attr     AttributeContext
	optional bool
}

var _ InterpretableAttribute = &evalAttr{}

// ID of the attribute instruction.
func (a *evalAttr) ID() int64 {
	return a.attr.ID()
}

// AddQualifier implements the InterpretableAttribute interface method.
func (a *evalAttr) AddQualifier(qual Qualifier) (Attribute, error) {
	attr, err := a.AddQualifierContext(context.Background(), qual)
	return attr.(Attribute), err
}

// AddQualifierContext implements the InterpretableAttributeContext interface method.
func (a *evalAttr) AddQualifierContext(ctx context.Context, qual Qualifier) (AttributeContext, error) {
	attr, err := a.attr.AddQualifierContext(ctx, qual)
	a.attr = attr
	return attr, err
}

// Attr implements the InterpretableAttribute interface method.
func (a *evalAttr) Attr() Attribute {
	return unwrapContextImpl(a.attr).(Attribute)
}

// AttrContext implements the InterpretableAttributeContext interface method.
func (a *evalAttr) AttrContext(ctx context.Context) AttributeContext {
	return a.attr
}

// Adapter implements the InterpretableAttribute interface method.
func (a *evalAttr) Adapter() types.Adapter {
	return a.adapter
}

// Eval implements the Interpretable interface method.
func (a *evalAttr) Eval(vars Activation) ref.Val {
	return a.EvalContext(context.Background(), vars)
}

// EvalContext implements the InterpretableContext interface method.
func (a *evalAttr) EvalContext(ctx context.Context, vars Activation) ref.Val {
	v, err := a.attr.Resolve(vars)
	if err != nil {
		return types.LabelErrNode(a.ID(), types.WrapErr(err))
	}
	return a.adapter.NativeToValue(v)
}

// Qualify proxies to the Attribute's Qualify method.
func (a *evalAttr) Qualify(ctx Activation, obj any) (any, error) {
	return a.attr.Qualify(ctx, obj)
}

// QualifyContext proxies to the AttributeContext's Qualify method.
func (a *evalAttr) QualifyContext(ctx context.Context, vars Activation, obj any) (any, error) {
	return a.attr.Qualify(vars, obj)
}

// QualifyIfPresent proxies to the Attribute's QualifyIfPresent method.
func (a *evalAttr) QualifyIfPresent(ctx Activation, obj any, presenceOnly bool) (any, bool, error) {
	return a.attr.QualifyIfPresent(ctx, obj, presenceOnly)
}

// QualifyIfPresentContext proxies to the AttributeContext's QualifyIfPresent method.
func (a *evalAttr) QualifyIfPresentContext(ctx context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return a.attr.QualifyIfPresentContext(ctx, vars, obj, presenceOnly)
}

func (a *evalAttr) IsOptional() bool {
	return a.optional
}

// Resolve proxies to the Attribute's Resolve method.
func (a *evalAttr) Resolve(ctx Activation) (any, error) {
	return a.attr.Resolve(ctx)
}

// ResolveContext proxies to the AttributeContext's Resolve method.
func (a *evalAttr) ResolveContext(ctx context.Context, vars Activation) (any, error) {
	return a.attr.ResolveContext(ctx, vars)
}

type evalWatchConstructor struct {
	constructor InterpretableConstructorContext
	observer    EvalObserver
}

// InitVals implements the InterpretableConstructor InitVals function.
func (c *evalWatchConstructor) InitVals() []Interpretable {
	return c.constructor.InitVals()
}

// Type implements the InterpretableConstructor Type function.
func (c *evalWatchConstructor) Type() ref.Type {
	return c.constructor.Type()
}

// ID implements the Interpretable ID function.
func (c *evalWatchConstructor) ID() int64 {
	return c.constructor.ID()
}

// Eval implements the Interpretable Eval function.
func (c *evalWatchConstructor) Eval(vars Activation) ref.Val {
	return c.EvalContext(context.Background(), vars)
}

// Eval implements the Interpretable Eval function.
func (c *evalWatchConstructor) EvalContext(ctx context.Context, vars Activation) ref.Val {
	val := c.constructor.EvalContext(ctx, vars)
	c.observer(c.ID(), c.constructor, val)
	return val
}

func invalidOptionalEntryInit(field any, value ref.Val) ref.Val {
	return types.NewErr("cannot initialize optional entry '%v' from non-optional value %v", field, value)
}

func invalidOptionalElementInit(value ref.Val) ref.Val {
	return types.NewErr("cannot initialize optional list element from non-optional value %v", value)
}
