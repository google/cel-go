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

package checker

import (
	"fmt"
	"reflect"

	"celgo/ast"
	"celgo/common"
	"celgo/semantics"
	"celgo/semantics/types"
)

type checker struct {
	env                *Env
	container          string
	mappings           *Mapping
	freeTypeVarCounter int

	types      map[int64]types.Type
	references map[int64]semantics.Reference
}

func Check(env *Env, container string, expression ast.Expression) *semantics.Semantics {
	c := checker{
		env:                env,
		container:          container,
		mappings:           newMapping(),
		freeTypeVarCounter: 0,

		types:      make(map[int64]types.Type),
		references: make(map[int64]semantics.Reference),
	}

	c.check(expression)

	// Walk over the final type map substituting any type parameters either by their bound value or
	// by DYN.
	m := make(map[int64]types.Type)
	for k, v := range c.types {
		m[k] = substitute(c.mappings, v, true)
	}

	return semantics.New(m, c.references)
}

func (c *checker) check(e ast.Expression) {
	if e == nil {
		return
	}

	switch e.(type) {
	case *ast.StringConstant:
		c.checkStringConstant(e.(*ast.StringConstant))
	case *ast.Int64Constant:
		c.checkInt64Constant(e.(*ast.Int64Constant))
	case *ast.Uint64Constant:
		c.checkUint64Constant(e.(*ast.Uint64Constant))
	case *ast.BytesConstant:
		c.checkBytesConstant(e.(*ast.BytesConstant))
	case *ast.DoubleConstant:
		c.checkDoubleConstant(e.(*ast.DoubleConstant))
	case *ast.BoolConstant:
		c.checkBoolConstant(e.(*ast.BoolConstant))
	case *ast.NullConstant:
		c.checkNullConstant(e.(*ast.NullConstant))
	case *ast.IdentExpression:
		c.checkIdent(e.(*ast.IdentExpression))
	case *ast.SelectExpression:
		c.checkSelect(e.(*ast.SelectExpression))
	case *ast.CallExpression:
		c.checkCall(e.(*ast.CallExpression))
	case *ast.CreateListExpression:
		c.checkCreateList(e.(*ast.CreateListExpression))
	case *ast.CreateStructExpression:
		c.checkCreateStruct(e.(*ast.CreateStructExpression))
	case *ast.CreateMessageExpression:
		c.checkCreateMessage(e.(*ast.CreateMessageExpression))
	case *ast.ComprehensionExpression:
		c.checkComprehension(e.(*ast.ComprehensionExpression))
	default:
		panic(fmt.Sprintf("Unrecognized ast type: %v", reflect.TypeOf(e)))
	}
}

func (c *checker) checkInt64Constant(e *ast.Int64Constant) {
	c.setType(e, types.Int64)
}

func (c *checker) checkUint64Constant(e *ast.Uint64Constant) {
	c.setType(e, types.Uint64)
}

func (c *checker) checkStringConstant(e *ast.StringConstant) {
	c.setType(e, types.String)
}

func (c *checker) checkBytesConstant(e *ast.BytesConstant) {
	c.setType(e, types.Bytes)
}

func (c *checker) checkDoubleConstant(e *ast.DoubleConstant) {
	c.setType(e, types.Double)
}

func (c *checker) checkBoolConstant(e *ast.BoolConstant) {
	c.setType(e, types.Bool)
}

func (c *checker) checkNullConstant(e *ast.NullConstant) {
	c.setType(e, types.Null)
}

func (c *checker) checkIdent(e *ast.IdentExpression) {
	if ident := c.env.LookupIdent(c.container, e.Name); ident != nil {
		c.setType(e, ident.Type())
		c.setReference(e, semantics.NewIdentReference(ident.Name(), ident.Value()))
		return
	}

	c.setType(e, types.Error)
	c.env.errors.undeclaredReference(e.Location(), c.container, e.Name)
}

func (c *checker) checkSelect(e *ast.SelectExpression) {
	// Before traversing down the tree, try to interpret as qualified name.
	qname, found := asQualifiedName(e)
	if found {
		ident := c.env.LookupIdent(c.container, qname)
		if ident != nil {
			if e.TestOnly {
				c.env.errors.expressionDoesNotSelectField(e.Location())
				c.setType(e, types.Bool)
			} else {
				c.setType(e, ident.Type())
				c.setReference(e, semantics.NewIdentReference(ident.Name(), ident.Value()))
			}
			return
		}
	}

	// Interpret as field selection, first traversing down the operand.
	c.check(e.Target)
	targetType := c.getType(e.Target)
	resultType := types.Error

	switch targetType.Kind() {
	case types.KindError, types.KindDynamic:
		resultType = types.Dynamic

	case types.KindMessage:
		messageType := targetType.(*types.MessageType)
		if fieldType, found := c.lookupFieldType(e.Location(), messageType, e.Field); found {
			resultType = fieldType.Type
			if e.TestOnly && !fieldType.SupportsPresence {
				c.env.errors.fieldDoesNotSupportPresenceCheck(e.Location(), e.Field)
			}
		}

	case types.KindMap:
		mapType := targetType.(*types.MapType)
		resultType = mapType.ValueType

	default:
		c.env.errors.typeDoesNotSupportFieldSelection(e.Location(), targetType)
	}

	if e.TestOnly {
		resultType = types.Bool
	}

	c.setType(e, resultType)
}

func (c *checker) checkCall(call *ast.CallExpression) {
	// Traverse arguments.
	for _, arg := range call.Args {
		c.check(arg)
	}

	var resolution *overloadResolution

	if call.Target == nil {
		// Regular static call with simple name.
		if fn := c.env.LookupFunction(c.container, call.Function); fn != nil {
			resolution = c.resolveOverload(call.Location(), fn, nil, call.Args)
		} else {
			c.env.errors.undeclaredReference(call.Location(), c.container, call.Function)
		}
	} else {
		// Check whether the target is actually a qualified name for a static function.
		if qname, found := asQualifiedName(call.Target); found {
			fn := c.env.LookupFunction(c.container, qname+"."+call.Function)
			if fn != nil {
				resolution = c.resolveOverload(call.Location(), fn, nil, call.Args)
			}
		}

		if resolution == nil {
			// Regular instance call.
			c.check(call.Target)

			if fn := c.env.LookupFunction(c.container, call.Function); fn != nil {
				resolution = c.resolveOverload(call.Location(), fn, call.Target, call.Args)
			} else {
				c.env.errors.undeclaredReference(call.Location(), c.container, call.Function)
			}
		}
	}

	if resolution != nil {
		c.setType(call, resolution.Type)
		c.setReference(call, resolution.Reference)
	} else {
		c.setType(call, types.Error)
	}
}

func (c *checker) resolveOverload(
	loc common.Location,
	fn *semantics.Function, target ast.Expression, args []ast.Expression) *overloadResolution {

	argTypes := []types.Type{}
	if target != nil {
		argTypes = append(argTypes, c.getType(target))
	}
	for _, arg := range args {
		argTypes = append(argTypes, c.getType(arg))
	}

	var resultType types.Type = nil
	var ref *semantics.FunctionReference = nil
	for _, overload := range fn.Overloads() {
		if (target == nil && overload.IsInstance()) || (target != nil && !overload.IsInstance()) {
			// not a compatible call style.
			continue
		}

		overloadType := types.NewFunctionType(overload.ResultType(), overload.ArgTypes())
		if len(overload.TypeParams()) > 0 {
			// Instantiate overload's type with fresh type variables.
			substitutions := newMapping()
			for _, typePar := range overload.TypeParams() {
				substitutions.Add(types.NewTypeParam(typePar), c.newTypeVar())
			}

			substitution := substitute(substitutions, overloadType, false)
			overloadType = substitution.(*types.FunctionType)
		}

		candidateArgTypes := overloadType.ArgTypes()
		if c.isAssignableList(argTypes, candidateArgTypes) {
			if ref == nil {
				ref = semantics.NewFunctionReference(overload.Id())
			} else {
				ref = ref.AddOverloadReference(overload.Id())
			}

			if resultType == nil {
				// First matching overload, determines result type.
				resultType = substitute(c.mappings, overloadType.ResultType(), false)
			} else {
				// More than one matching overload, narrow result type to DYN.
				resultType = types.Dynamic
			}

		}
	}

	if resultType == nil {
		c.env.errors.noMatchingOverload(loc, fn.Name(), argTypes, target != nil)
		resultType = types.Error
		return nil
	}

	return newResolution(ref, resultType)
}

func (c *checker) checkCreateList(create *ast.CreateListExpression) {
	var elemType types.Type = nil
	for _, e := range create.Entries {
		c.check(e)
		elemType = c.joinTypes(e.Location(), elemType, c.getType(e))
	}
	if elemType == nil {
		// If the list is empty, assign free type var to elem type.
		elemType = c.newTypeVar()
	}
	c.setType(create, types.NewList(elemType))
}

func (c *checker) checkCreateStruct(str *ast.CreateStructExpression) {

	var keyType types.Type = nil
	var valueType types.Type = nil
	for _, ent := range str.Entries {
		c.check(ent.Key)
		keyType = c.joinTypes(ent.Key.Location(), keyType, c.getType(ent.Key))

		c.check(ent.Value)
		valueType = c.joinTypes(ent.Value.Location(), valueType, c.getType(ent.Value))
	}
	if keyType == nil {
		// If the map is empty, assign free type variables to typeKey and value type.
		keyType = c.newTypeVar()
		valueType = c.newTypeVar()
	}
	c.setType(str, types.NewMap(keyType, valueType))
}

func (c *checker) checkCreateMessage(msg *ast.CreateMessageExpression) {
	// Determine the type of the message.
	messageType := types.Error
	decl := c.env.LookupIdent(c.container, msg.MessageName)
	if decl == nil {
		c.env.errors.undeclaredReference(msg.Location(), c.container, msg.MessageName)
		return
	}
	c.setReference(msg, semantics.NewIdentReference(decl.Name(), nil))
	if decl.Type().Kind() != types.KindError {
		if decl.Type().Kind() != types.KindType {
			c.env.errors.notAType(msg.Location(), decl.Type())
		} else {
			messageType = decl.Type().(*types.TypeType).Target()
			if messageType.Kind() != types.KindMessage {
				c.env.errors.notAMessageType(msg.Location(), messageType)
				messageType = types.Error
			}
		}
	}
	c.setType(msg, messageType)

	// Check the field initializers.
	for _, f := range msg.Fields {
		value := f.Initializer
		c.check(value)

		fieldType := types.Error
		if mt, ok := messageType.(*types.MessageType); ok {
			if t, found := c.lookupFieldType(f.Location(), mt, f.Name); found {
				fieldType = t.Type
			}
		}
		if !c.isAssignable(fieldType, c.getType(value)) {
			c.env.errors.fieldTypeMismatch(f.Location(), f.Name, fieldType, c.getType(value))
		}
	}
}

func (c *checker) checkComprehension(comp *ast.ComprehensionExpression) {
	c.check(comp.Target)
	c.check(comp.Init)
	accuType := c.getType(comp.Init)
	rangeType := c.getType(comp.Target)
	var varType types.Type

	switch rangeType.Kind() {
	case types.KindList:
		varType = rangeType.(*types.ListType).ElementType
	case types.KindMap:
		// Ranges over the keys.
		varType = rangeType.(*types.MapType).KeyType
	case types.KindDynamic, types.KindError:
		varType = types.Dynamic
	default:
		c.env.errors.notAComprehensionRange(comp.Target.Location(), rangeType)
	}

	c.env.enterScope()
	c.env.AddIdent(semantics.NewIdent(comp.Accumulator, accuType, nil))
	// Declare iteration variable on inner scope.
	c.env.enterScope()
	c.env.AddIdent(semantics.NewIdent(comp.Variable, varType, nil))
	c.check(comp.LoopCondition)
	c.assertType(comp.LoopCondition, types.Bool)
	c.check(comp.LoopStep)
	c.assertType(comp.LoopStep, accuType)
	// Forget iteration variable, as result expression must only depend on accu.
	c.env.exitScope()
	c.check(comp.Result)
	c.env.exitScope()
	c.setType(comp, c.getType(comp.Result))
}

// Checks compatibility of joined types, and returns the most general common type.
func (c *checker) joinTypes(loc common.Location, previous types.Type, current types.Type) types.Type {
	if previous == nil {
		return current
	}

	if !c.isAssignable(previous, current) {
		c.env.errors.aggregateTypeMismatch(loc, previous, current)
		return previous
	}

	return mostGeneral(previous, current)
}

func (c *checker) newTypeVar() types.Type {
	id := c.freeTypeVarCounter
	c.freeTypeVarCounter++
	return types.NewTypeParam(fmt.Sprintf("_var%d", id))
}

func (c *checker) isAssignable(t1 types.Type, t2 types.Type) bool {
	subs := isAssignable(c.mappings, t1, t2)
	if subs != nil {
		c.mappings = subs
		return true
	}

	return false
}

func (c *checker) isAssignableList(l1 []types.Type, l2 []types.Type) bool {
	subs := isAssignableList(c.mappings, l1, l2)
	if subs != nil {
		c.mappings = subs
		return true
	}

	return false
}

func (c *checker) lookupFieldType(l common.Location, messageType *types.MessageType, fieldName string) (*FieldType, bool) {
	if c.env.typeProvider.LookupType(messageType.Name()) == nil {
		// This should not happen, anyway, report an error.
		c.env.errors.unexpectedFailedResolution(l, messageType.Name())
		return nil, false
	}

	if ft, found := c.env.typeProvider.LookupFieldType(messageType, fieldName); found {
		return ft, found
	}

	c.env.errors.undefinedField(l, fieldName)
	return nil, false
}

func (c *checker) setType(e ast.Expression, t types.Type) {
	if old, found := c.types[e.Id()]; found && !old.Equals(t) {
		panic(fmt.Sprintf("(Incompatible) Type already exists for expression: %v(%d) old:%v, new:%v", e, e.Id(), old, t))
	}

	c.types[e.Id()] = t
}

func (c *checker) getType(e ast.Expression) types.Type {
	return c.types[e.Id()]
}

func (c *checker) setReference(e ast.Expression, r semantics.Reference) {
	if old, found := c.references[e.Id()]; found && !old.Equals(r) {
		panic(fmt.Sprintf("Reference already exists for expression: %v(%d) old:%v, new:%v", e, e.Id, old, r))
	}

	c.references[e.Id()] = r
}

func (c *checker) assertType(e ast.Expression, t types.Type) {
	if !c.isAssignable(t, c.getType(e)) {
		c.env.errors.typeMismatch(e.Location(), t, c.getType(e))
	}
}

type overloadResolution struct {
	Reference *semantics.FunctionReference
	Type      types.Type
}

func newResolution(ref *semantics.FunctionReference, t types.Type) *overloadResolution {
	return &overloadResolution{
		Reference: ref,
		Type:      t,
	}
}
