// Copyright 2018 Google LLC
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
	"fmt"
	"strings"

	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

// newPlanner creates an interpretablePlanner which references a Dispatcher, TypeProvider,
// TypeAdapter, Container, and CheckedExpr value. These pieces of data are used to resolve
// functions, types, and namespaced identifiers at plan time rather than at runtime since
// it only needs to be done once and may be semi-expensive to compute.
func newPlanner(disp Dispatcher,
	provider types.Provider,
	adapter types.Adapter,
	attrFactory AttributeFactory,
	cont *containers.Container,
	exprAST *ast.AST) *planner {
	return &planner{
		disp:        disp,
		provider:    provider,
		adapter:     adapter,
		attrFactory: attrFactory,
		container:   cont,
		refMap:      exprAST.ReferenceMap(),
		typeMap:     exprAST.TypeMap(),
		decorators:  make([]InterpretableDecorator, 0),
		observers:   make([]StatefulObserver, 0),
	}
}

// planner is an implementation of the interpretablePlanner interface.
type planner struct {
	disp        Dispatcher
	provider    types.Provider
	adapter     types.Adapter
	attrFactory AttributeFactory
	container   *containers.Container
	refMap      map[int64]*ast.ReferenceInfo
	typeMap     map[int64]*types.Type
	decorators  []InterpretableDecorator
	observers   []StatefulObserver
}

// Plan implements the interpretablePlanner interface. This implementation of the Plan method also
// applies decorators to each Interpretable generated as part of the overall plan. Decorators are
// useful for layering functionality into the evaluation that is not natively understood by CEL,
// such as state-tracking, expression re-write, and possibly efficient thread-safe memoization of
// repeated expressions.
func (p *planner) Plan(expr ast.Expr) (Interpretable, error) {
	i, err := p.plan(expr)
	if err != nil {
		return nil, err
	}
	if len(p.observers) == 0 {
		return i, nil
	}
	return &ObservableInterpretable{Interpretable: i, observers: p.observers}, nil
}

func (p *planner) plan(expr ast.Expr) (Interpretable, error) {
	switch expr.Kind() {
	case ast.CallKind:
		return p.decorate(p.planCall(expr))
	case ast.IdentKind:
		return p.decorate(p.planIdent(expr))
	case ast.LiteralKind:
		return p.decorate(p.planConst(expr))
	case ast.SelectKind:
		return p.decorate(p.planSelect(expr))
	case ast.ListKind:
		return p.decorate(p.planCreateList(expr))
	case ast.MapKind:
		return p.decorate(p.planCreateMap(expr))
	case ast.StructKind:
		return p.decorate(p.planCreateStruct(expr))
	case ast.ComprehensionKind:
		return p.decorate(p.planComprehension(expr))
	}
	return nil, fmt.Errorf("unsupported expr: %v", expr)
}

// decorate applies the InterpretableDecorator functions to the given Interpretable.
// Both the Interpretable and error generated by a Plan step are accepted as arguments
// for convenience.
func (p *planner) decorate(i Interpretable, err error) (Interpretable, error) {
	if err != nil {
		return nil, err
	}
	for _, dec := range p.decorators {
		i, err = dec(i)
		if err != nil {
			return nil, err
		}
	}
	return i, nil
}

// planIdent creates an Interpretable that resolves an identifier from an Activation.
func (p *planner) planIdent(expr ast.Expr) (Interpretable, error) {
	// Establish whether the identifier is in the reference map.
	if identRef, found := p.refMap[expr.ID()]; found {
		return p.planCheckedIdent(expr.ID(), identRef)
	}
	// Create the possible attribute list for the unresolved reference.
	ident := expr.AsIdent()
	return &evalAttr{
		adapter: p.adapter,
		attr:    p.attrFactory.MaybeAttribute(expr.ID(), ident),
	}, nil
}

func (p *planner) planCheckedIdent(id int64, identRef *ast.ReferenceInfo) (Interpretable, error) {
	// Plan a constant reference if this is the case for this simple identifier.
	if identRef.Value != nil {
		return NewConstValue(id, identRef.Value), nil
	}

	// Check to see whether the type map indicates this is a type name. All types should be
	// registered with the provider.
	cType := p.typeMap[id]
	if cType.Kind() == types.TypeKind {
		cVal, found := p.provider.FindIdent(identRef.Name)
		if !found {
			return nil, fmt.Errorf("reference to undefined type: %s", identRef.Name)
		}
		return NewConstValue(id, cVal), nil
	}

	// Otherwise, return the attribute for the resolved identifier name.
	return &evalAttr{
		adapter: p.adapter,
		attr:    p.attrFactory.AbsoluteAttribute(id, identRef.Name),
	}, nil
}

// planSelect creates an Interpretable with either:
//
//	a) selects a field from a map or proto.
//	b) creates a field presence test for a select within a has() macro.
//	c) resolves the select expression to a namespaced identifier.
func (p *planner) planSelect(expr ast.Expr) (Interpretable, error) {
	// If the Select id appears in the reference map from the CheckedExpr proto then it is either
	// a namespaced identifier or enum value.
	if identRef, found := p.refMap[expr.ID()]; found {
		return p.planCheckedIdent(expr.ID(), identRef)
	}

	sel := expr.AsSelect()
	// Plan the operand evaluation.
	op, err := p.plan(sel.Operand())
	if err != nil {
		return nil, err
	}
	opType := p.typeMap[sel.Operand().ID()]

	// If the Select was marked TestOnly, this is a presence test.
	//
	// Note: presence tests are defined for structured (e.g. proto) and dynamic values (map, json)
	// as follows:
	//  - True if the object field has a non-default value, e.g. obj.str != ""
	//  - True if the dynamic value has the field defined, e.g. key in map
	//
	// However, presence tests are not defined for qualified identifier names with primitive types.
	// If a string named 'a.b.c' is declared in the environment and referenced within `has(a.b.c)`,
	// it is not clear whether has should error or follow the convention defined for structured
	// values.

	// Establish the attribute reference.
	attr, isAttr := op.(InterpretableAttribute)
	if !isAttr {
		attr, err = p.relativeAttr(op.ID(), op, false)
		if err != nil {
			return nil, err
		}
	}

	// Build a qualifier for the attribute.
	qual, err := p.attrFactory.NewQualifier(opType, expr.ID(), sel.FieldName(), false)
	if err != nil {
		return nil, err
	}
	// Modify the attribute to be test-only.
	if sel.IsTestOnly() {
		attr = &evalTestOnly{
			id:                     expr.ID(),
			InterpretableAttribute: attr,
		}
	}
	// Append the qualifier on the attribute.
	_, err = attr.AddQualifier(qual)
	return attr, err
}

// planCall creates a callable Interpretable while specializing for common functions and invocation
// patterns. Specifically, conditional operators &&, ||, ?:, and (in)equality functions result in
// optimized Interpretable values.
func (p *planner) planCall(expr ast.Expr) (Interpretable, error) {
	call := expr.AsCall()
	target, fnName, oName := p.resolveFunction(expr)
	argCount := len(call.Args())
	var offset int
	if target != nil {
		argCount++
		offset++
	}

	args := make([]Interpretable, argCount)
	if target != nil {
		arg, err := p.plan(target)
		if err != nil {
			return nil, err
		}
		args[0] = arg
	}
	for i, argExpr := range call.Args() {
		arg, err := p.plan(argExpr)
		if err != nil {
			return nil, err
		}
		args[i+offset] = arg
	}

	// Generate specialized Interpretable operators by function name if possible.
	switch fnName {
	case operators.LogicalAnd:
		return p.planCallLogicalAnd(expr, args)
	case operators.LogicalOr:
		return p.planCallLogicalOr(expr, args)
	case operators.Conditional:
		return p.planCallConditional(expr, args)
	case operators.Equals:
		return p.planCallEqual(expr, args)
	case operators.NotEquals:
		return p.planCallNotEqual(expr, args)
	case operators.Index:
		return p.planCallIndex(expr, args, false)
	case operators.OptSelect, operators.OptIndex:
		return p.planCallIndex(expr, args, true)
	}

	// Otherwise, generate Interpretable calls specialized by argument count.
	// Try to find the specific function by overload id.
	var fnDef *functions.Overload
	if oName != "" {
		fnDef, _ = p.disp.FindOverload(oName)
	}
	// If the overload id couldn't resolve the function, try the simple function name.
	if fnDef == nil {
		fnDef, _ = p.disp.FindOverload(fnName)
	}
	switch argCount {
	case 0:
		return p.planCallZero(expr, fnName, oName, fnDef)
	case 1:
		// If the FunctionOp has been used, then use it as it may exist for the purposes
		// of dynamic dispatch within a singleton function implementation.
		if fnDef != nil && fnDef.Unary == nil && fnDef.Function != nil {
			return p.planCallVarArgs(expr, fnName, oName, fnDef, args)
		}
		return p.planCallUnary(expr, fnName, oName, fnDef, args)
	case 2:
		// If the FunctionOp has been used, then use it as it may exist for the purposes
		// of dynamic dispatch within a singleton function implementation.
		if fnDef != nil && fnDef.Binary == nil && fnDef.Function != nil {
			return p.planCallVarArgs(expr, fnName, oName, fnDef, args)
		}
		return p.planCallBinary(expr, fnName, oName, fnDef, args)
	default:
		return p.planCallVarArgs(expr, fnName, oName, fnDef, args)
	}
}

// planCallZero generates a zero-arity callable Interpretable.
func (p *planner) planCallZero(expr ast.Expr,
	function string,
	overload string,
	impl *functions.Overload) (Interpretable, error) {
	if impl == nil || impl.Function == nil {
		return nil, fmt.Errorf("no such overload: %s()", function)
	}
	return &evalZeroArity{
		id:       expr.ID(),
		function: function,
		overload: overload,
		impl:     impl.Function,
	}, nil
}

// planCallUnary generates a unary callable Interpretable.
func (p *planner) planCallUnary(expr ast.Expr,
	function string,
	overload string,
	impl *functions.Overload,
	args []Interpretable) (Interpretable, error) {
	var fn functions.UnaryOp
	var trait int
	var nonStrict bool
	if impl != nil {
		if impl.Unary == nil {
			return nil, fmt.Errorf("no such overload: %s(arg)", function)
		}
		fn = impl.Unary
		trait = impl.OperandTrait
		nonStrict = impl.NonStrict
	}
	return &evalUnary{
		id:        expr.ID(),
		function:  function,
		overload:  overload,
		arg:       args[0],
		trait:     trait,
		impl:      fn,
		nonStrict: nonStrict,
	}, nil
}

// planCallBinary generates a binary callable Interpretable.
func (p *planner) planCallBinary(expr ast.Expr,
	function string,
	overload string,
	impl *functions.Overload,
	args []Interpretable) (Interpretable, error) {
	var fn functions.BinaryOp
	var trait int
	var nonStrict bool
	if impl != nil {
		if impl.Binary == nil {
			return nil, fmt.Errorf("no such overload: %s(lhs, rhs)", function)
		}
		fn = impl.Binary
		trait = impl.OperandTrait
		nonStrict = impl.NonStrict
	}
	return &evalBinary{
		id:        expr.ID(),
		function:  function,
		overload:  overload,
		lhs:       args[0],
		rhs:       args[1],
		trait:     trait,
		impl:      fn,
		nonStrict: nonStrict,
	}, nil
}

// planCallVarArgs generates a variable argument callable Interpretable.
func (p *planner) planCallVarArgs(expr ast.Expr,
	function string,
	overload string,
	impl *functions.Overload,
	args []Interpretable) (Interpretable, error) {
	var fn functions.FunctionOp
	var trait int
	var nonStrict bool
	if impl != nil {
		if impl.Function == nil {
			return nil, fmt.Errorf("no such overload: %s(...)", function)
		}
		fn = impl.Function
		trait = impl.OperandTrait
		nonStrict = impl.NonStrict
	}
	return &evalVarArgs{
		id:        expr.ID(),
		function:  function,
		overload:  overload,
		args:      args,
		trait:     trait,
		impl:      fn,
		nonStrict: nonStrict,
	}, nil
}

// planCallEqual generates an equals (==) Interpretable.
func (p *planner) planCallEqual(expr ast.Expr, args []Interpretable) (Interpretable, error) {
	return &evalEq{
		id:  expr.ID(),
		lhs: args[0],
		rhs: args[1],
	}, nil
}

// planCallNotEqual generates a not equals (!=) Interpretable.
func (p *planner) planCallNotEqual(expr ast.Expr, args []Interpretable) (Interpretable, error) {
	return &evalNe{
		id:  expr.ID(),
		lhs: args[0],
		rhs: args[1],
	}, nil
}

// planCallLogicalAnd generates a logical and (&&) Interpretable.
func (p *planner) planCallLogicalAnd(expr ast.Expr, args []Interpretable) (Interpretable, error) {
	return &evalAnd{
		id:    expr.ID(),
		terms: args,
	}, nil
}

// planCallLogicalOr generates a logical or (||) Interpretable.
func (p *planner) planCallLogicalOr(expr ast.Expr, args []Interpretable) (Interpretable, error) {
	return &evalOr{
		id:    expr.ID(),
		terms: args,
	}, nil
}

// planCallConditional generates a conditional / ternary (c ? t : f) Interpretable.
func (p *planner) planCallConditional(expr ast.Expr, args []Interpretable) (Interpretable, error) {
	cond := args[0]
	t := args[1]
	var tAttr Attribute
	truthyAttr, isTruthyAttr := t.(InterpretableAttribute)
	if isTruthyAttr {
		tAttr = truthyAttr.Attr()
	} else {
		tAttr = p.attrFactory.RelativeAttribute(t.ID(), t)
	}

	f := args[2]
	var fAttr Attribute
	falsyAttr, isFalsyAttr := f.(InterpretableAttribute)
	if isFalsyAttr {
		fAttr = falsyAttr.Attr()
	} else {
		fAttr = p.attrFactory.RelativeAttribute(f.ID(), f)
	}

	return &evalAttr{
		adapter: p.adapter,
		attr:    p.attrFactory.ConditionalAttribute(expr.ID(), cond, tAttr, fAttr),
	}, nil
}

// planCallIndex either extends an attribute with the argument to the index operation, or creates
// a relative attribute based on the return of a function call or operation.
func (p *planner) planCallIndex(expr ast.Expr, args []Interpretable, optional bool) (Interpretable, error) {
	op := args[0]
	ind := args[1]
	opType := p.typeMap[op.ID()]

	// Establish the attribute reference.
	var err error
	attr, isAttr := op.(InterpretableAttribute)
	if !isAttr {
		attr, err = p.relativeAttr(op.ID(), op, false)
		if err != nil {
			return nil, err
		}
	}

	// Construct the qualifier type.
	var qual Qualifier
	switch ind := ind.(type) {
	case InterpretableConst:
		qual, err = p.attrFactory.NewQualifier(opType, expr.ID(), ind.Value(), optional)
	case InterpretableAttribute:
		qual, err = p.attrFactory.NewQualifier(opType, expr.ID(), ind, optional)
	default:
		qual, err = p.relativeAttr(expr.ID(), ind, optional)
	}
	if err != nil {
		return nil, err
	}

	// Add the qualifier to the attribute
	_, err = attr.AddQualifier(qual)
	return attr, err
}

// planCreateList generates a list construction Interpretable.
func (p *planner) planCreateList(expr ast.Expr) (Interpretable, error) {
	list := expr.AsList()
	optionalIndices := list.OptionalIndices()
	elements := list.Elements()
	optionals := make([]bool, len(elements))
	for _, index := range optionalIndices {
		if index < 0 || index >= int32(len(elements)) {
			return nil, fmt.Errorf("optional index %d out of element bounds [0, %d]", index, len(elements))
		}
		optionals[index] = true
	}
	elems := make([]Interpretable, len(elements))
	for i, elem := range elements {
		elemVal, err := p.plan(elem)
		if err != nil {
			return nil, err
		}
		elems[i] = elemVal
	}
	return &evalList{
		id:           expr.ID(),
		elems:        elems,
		optionals:    optionals,
		hasOptionals: len(optionalIndices) != 0,
		adapter:      p.adapter,
	}, nil
}

// planCreateStruct generates a map or object construction Interpretable.
func (p *planner) planCreateMap(expr ast.Expr) (Interpretable, error) {
	m := expr.AsMap()
	entries := m.Entries()
	optionals := make([]bool, len(entries))
	keys := make([]Interpretable, len(entries))
	vals := make([]Interpretable, len(entries))
	hasOptionals := false
	for i, e := range entries {
		entry := e.AsMapEntry()
		keyVal, err := p.plan(entry.Key())
		if err != nil {
			return nil, err
		}
		keys[i] = keyVal

		valVal, err := p.plan(entry.Value())
		if err != nil {
			return nil, err
		}
		vals[i] = valVal
		optionals[i] = entry.IsOptional()
		hasOptionals = hasOptionals || entry.IsOptional()
	}
	return &evalMap{
		id:           expr.ID(),
		keys:         keys,
		vals:         vals,
		optionals:    optionals,
		hasOptionals: hasOptionals,
		adapter:      p.adapter,
	}, nil
}

// planCreateObj generates an object construction Interpretable.
func (p *planner) planCreateStruct(expr ast.Expr) (Interpretable, error) {
	obj := expr.AsStruct()
	typeName, defined := p.resolveTypeName(obj.TypeName())
	if !defined {
		return nil, fmt.Errorf("unknown type: %s", obj.TypeName())
	}
	objFields := obj.Fields()
	optionals := make([]bool, len(objFields))
	fields := make([]string, len(objFields))
	vals := make([]Interpretable, len(objFields))
	hasOptionals := false
	for i, f := range objFields {
		field := f.AsStructField()
		fields[i] = field.Name()
		val, err := p.plan(field.Value())
		if err != nil {
			return nil, err
		}
		vals[i] = val
		optionals[i] = field.IsOptional()
		hasOptionals = hasOptionals || field.IsOptional()
	}
	return &evalObj{
		id:           expr.ID(),
		typeName:     typeName,
		fields:       fields,
		vals:         vals,
		optionals:    optionals,
		hasOptionals: hasOptionals,
		provider:     p.provider,
	}, nil
}

// planComprehension generates an Interpretable fold operation.
func (p *planner) planComprehension(expr ast.Expr) (Interpretable, error) {
	fold := expr.AsComprehension()
	accu, err := p.plan(fold.AccuInit())
	if err != nil {
		return nil, err
	}
	iterRange, err := p.plan(fold.IterRange())
	if err != nil {
		return nil, err
	}
	cond, err := p.plan(fold.LoopCondition())
	if err != nil {
		return nil, err
	}
	step, err := p.plan(fold.LoopStep())
	if err != nil {
		return nil, err
	}
	result, err := p.plan(fold.Result())
	if err != nil {
		return nil, err
	}
	return &evalFold{
		id:        expr.ID(),
		accuVar:   fold.AccuVar(),
		accu:      accu,
		iterVar:   fold.IterVar(),
		iterVar2:  fold.IterVar2(),
		iterRange: iterRange,
		cond:      cond,
		step:      step,
		result:    result,
		adapter:   p.adapter,
	}, nil
}

// planConst generates a constant valued Interpretable.
func (p *planner) planConst(expr ast.Expr) (Interpretable, error) {
	return NewConstValue(expr.ID(), expr.AsLiteral()), nil
}

// resolveTypeName takes a qualified string constructed at parse time, applies the proto
// namespace resolution rules to it in a scan over possible matching types in the TypeProvider.
func (p *planner) resolveTypeName(typeName string) (string, bool) {
	for _, qualifiedTypeName := range p.container.ResolveCandidateNames(typeName) {
		if _, found := p.provider.FindStructType(qualifiedTypeName); found {
			return qualifiedTypeName, true
		}
	}
	return "", false
}

// resolveFunction determines the call target, function name, and overload name from a given Expr
// value.
//
// The resolveFunction resolves ambiguities where a function may either be a receiver-style
// invocation or a qualified global function name.
// - The target expression may only consist of ident and select expressions.
// - The function is declared in the environment using its fully-qualified name.
// - The fully-qualified function name matches the string serialized target value.
func (p *planner) resolveFunction(expr ast.Expr) (ast.Expr, string, string) {
	// Note: similar logic exists within the `checker/checker.go`. If making changes here
	// please consider the impact on checker.go and consolidate implementations or mirror code
	// as appropriate.
	call := expr.AsCall()
	var target ast.Expr = nil
	if call.IsMemberFunction() {
		target = call.Target()
	}
	fnName := call.FunctionName()

	// Checked expressions always have a reference map entry, and _should_ have the fully qualified
	// function name as the fnName value.
	oRef, hasOverload := p.refMap[expr.ID()]
	if hasOverload {
		if len(oRef.OverloadIDs) == 1 {
			return target, fnName, oRef.OverloadIDs[0]
		}
		// Note, this namespaced function name will not appear as a fully qualified name in ASTs
		// built and stored before cel-go v0.5.0; however, this functionality did not work at all
		// before the v0.5.0 release.
		return target, fnName, ""
	}

	// Parse-only expressions need to handle the same logic as is normally performed at check time,
	// but with potentially much less information. The only reliable source of information about
	// which functions are configured is the dispatcher.
	if target == nil {
		// If the user has a parse-only expression, then it should have been configured as such in
		// the interpreter dispatcher as it may have been omitted from the checker environment.
		for _, qualifiedName := range p.container.ResolveCandidateNames(fnName) {
			_, found := p.disp.FindOverload(qualifiedName)
			if found {
				return nil, qualifiedName, ""
			}
		}
		// It's possible that the overload was not found, but this situation is accounted for in
		// the planCall phase; however, the leading dot used for denoting fully-qualified
		// namespaced identifiers must be stripped, as all declarations already use fully-qualified
		// names. This stripping behavior is handled automatically by the ResolveCandidateNames
		// call.
		return target, stripLeadingDot(fnName), ""
	}

	// Handle the situation where the function target actually indicates a qualified function name.
	qualifiedPrefix, maybeQualified := p.toQualifiedName(target)
	if maybeQualified {
		maybeQualifiedName := qualifiedPrefix + "." + fnName
		for _, qualifiedName := range p.container.ResolveCandidateNames(maybeQualifiedName) {
			_, found := p.disp.FindOverload(qualifiedName)
			if found {
				// Clear the target to ensure the proper arity is used for finding the
				// implementation.
				return nil, qualifiedName, ""
			}
		}
	}
	// In the default case, the function is exactly as it was advertised: a receiver call on with
	// an expression-based target with the given simple function name.
	return target, fnName, ""
}

// relativeAttr indicates that the attribute in this case acts as a qualifier and as such needs to
// be observed to ensure that it's evaluation value is properly recorded for state tracking.
func (p *planner) relativeAttr(id int64, eval Interpretable, opt bool) (InterpretableAttribute, error) {
	eAttr, ok := eval.(InterpretableAttribute)
	if !ok {
		eAttr = &evalAttr{
			adapter:  p.adapter,
			attr:     p.attrFactory.RelativeAttribute(id, eval),
			optional: opt,
		}
	}
	// This looks like it should either decorate the new evalAttr node, or early return the InterpretableAttribute
	decAttr, err := p.decorate(eAttr, nil)
	if err != nil {
		return nil, err
	}
	eAttr, ok = decAttr.(InterpretableAttribute)
	if !ok {
		return nil, fmt.Errorf("invalid attribute decoration: %v(%T)", decAttr, decAttr)
	}
	return eAttr, nil
}

// toQualifiedName converts an expression AST into a qualified name if possible, with a boolean
// 'found' value that indicates if the conversion is successful.
func (p *planner) toQualifiedName(operand ast.Expr) (string, bool) {
	// If the checker identified the expression as an attribute by the type-checker, then it can't
	// possibly be part of qualified name in a namespace.
	_, isAttr := p.refMap[operand.ID()]
	if isAttr {
		return "", false
	}
	// Since functions cannot be both namespaced and receiver functions, if the operand is not an
	// qualified variable name, return the (possibly) qualified name given the expressions.
	switch operand.Kind() {
	case ast.IdentKind:
		id := operand.AsIdent()
		return id, true
	case ast.SelectKind:
		sel := operand.AsSelect()
		// Test only expressions are not valid as qualified names.
		if sel.IsTestOnly() {
			return "", false
		}
		if qual, found := p.toQualifiedName(sel.Operand()); found {
			return qual + "." + sel.FieldName(), true
		}
	}
	return "", false
}

func stripLeadingDot(name string) string {
	if strings.HasPrefix(name, ".") {
		return name[1:]
	}
	return name
}
