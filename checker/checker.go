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

	expr "github.com/google/cel-spec/proto/v1/syntax"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/semantics"
	"github.com/google/cel-go/semantics/types"
)

type checker struct {
	env                *Env
	container          string
	mappings           *Mapping
	freeTypeVarCounter int
	sourceInfo 		   *expr.SourceInfo

	types      map[int64]types.Type
	references map[int64]semantics.Reference
}

func Check(env *Env, container string, parsedExpr *expr.ParsedExpr) *semantics.Semantics {
	c := checker{
		env:                env,
		container:          container,
		mappings:           newMapping(),
		freeTypeVarCounter: 0,
		sourceInfo:         parsedExpr.GetSourceInfo(),

		types:      make(map[int64]types.Type),
		references: make(map[int64]semantics.Reference),
	}
	c.check(parsedExpr.GetExpr())

	// Walk over the final type map substituting any type parameters either by their bound value or
	// by DYN.
	m := make(map[int64]types.Type)
	for k, v := range c.types {
		m[k] = substitute(c.mappings, v, true)
	}

	return semantics.New(m, c.references)
}

func (c *checker) check(e *expr.Expr) {
	if e == nil {
		return
	}

	switch e.ExprKind.(type) {
	case *expr.Expr_ConstExpr:
		constExpr := e.GetConstExpr()
		switch constExpr.ConstantKind.(type) {
		case *expr.Constant_BoolValue:
			c.checkBoolConstant(e)
		case *expr.Constant_BytesValue:
			c.checkBytesConstant(e)
		case *expr.Constant_DoubleValue:
			c.checkDoubleConstant(e)
		case *expr.Constant_Int64Value:
			c.checkInt64Constant(e)
		case *expr.Constant_NullValue:
			c.checkNullConstant(e)
		case *expr.Constant_StringValue:
			c.checkStringConstant(e)
		case *expr.Constant_Uint64Value:
			c.checkUint64Constant(e)
		}
	case *expr.Expr_IdentExpr:
		c.checkIdent(e)
	case *expr.Expr_SelectExpr:
		c.checkSelect(e)
	case *expr.Expr_CallExpr:
		c.checkCall(e)
	case *expr.Expr_ListExpr:
		c.checkCreateList(e)
	case *expr.Expr_StructExpr:
		c.checkCreateStruct(e)
	case *expr.Expr_ComprehensionExpr:
		c.checkComprehension(e)
	default:
		panic(fmt.Sprintf("Unrecognized ast type: %v", reflect.TypeOf(e)))
	}
}

func (c *checker) checkInt64Constant(e *expr.Expr) {
	c.setType(e, types.Int64)
}

func (c *checker) checkUint64Constant(e *expr.Expr) {
	c.setType(e, types.Uint64)
}

func (c *checker) checkStringConstant(e *expr.Expr) {
	c.setType(e, types.String)
}

func (c *checker) checkBytesConstant(e *expr.Expr) {
	c.setType(e, types.Bytes)
}

func (c *checker) checkDoubleConstant(e *expr.Expr) {
	c.setType(e, types.Double)
}

func (c *checker) checkBoolConstant(e *expr.Expr) {
	c.setType(e, types.Bool)
}

func (c *checker) checkNullConstant(e *expr.Expr) {
	c.setType(e, types.Null)
}

func (c *checker) checkIdent(e *expr.Expr) {
	identExpr := e.GetIdentExpr()
	if ident := c.env.LookupIdent(c.container, identExpr.Name); ident != nil {
		c.setType(e, ident.Type())
		c.setReference(e, semantics.NewIdentReference(ident.Name(), ident.Value()))
		return
	}

	c.setType(e, types.Error)
	c.env.errors.undeclaredReference(c.location(e), c.container, identExpr.Name)
}

func (c *checker) checkSelect(e *expr.Expr) {
	sel := e.GetSelectExpr()
	// Before traversing down the tree, try to interpret as qualified name.
	qname, found := asQualifiedName(e)
	if found {
		ident := c.env.LookupIdent(c.container, qname)
		if ident != nil {
			if sel.TestOnly {
				c.env.errors.expressionDoesNotSelectField(c.location(e))
				c.setType(e, types.Bool)
			} else {
				c.setType(e, ident.Type())
				c.setReference(e, semantics.NewIdentReference(ident.Name(), ident.Value()))
			}
			return
		}
	}

	// Interpret as field selection, first traversing down the operand.
	c.check(sel.Operand)
	targetType := c.getType(sel.Operand)
	resultType := types.Error

	switch targetType.Kind() {
	case types.KindError, types.KindDynamic:
		resultType = types.Dynamic

	case types.KindMessage:
		messageType := targetType.(*types.MessageType)
		if fieldType, found := c.lookupFieldType(c.location(e), messageType, sel.Field); found {
			resultType = fieldType.Type
			if sel.TestOnly && !fieldType.SupportsPresence {
				c.env.errors.fieldDoesNotSupportPresenceCheck(c.location(e), sel.Field)
			}
		}

	case types.KindMap:
		mapType := targetType.(*types.MapType)
		resultType = mapType.ValueType

	default:
		c.env.errors.typeDoesNotSupportFieldSelection(c.location(e), targetType)
	}

	if sel.TestOnly {
		resultType = types.Bool
	}

	c.setType(e, resultType)
}

func (c *checker) checkCall(e *expr.Expr) {
	call := e.GetCallExpr()
	// Traverse arguments.
	for _, arg := range call.Args {
		c.check(arg)
	}

	var resolution *overloadResolution

	if call.Target == nil {
		// Regular static call with simple name.
		if fn := c.env.LookupFunction(c.container, call.Function); fn != nil {
			resolution = c.resolveOverload(c.location(e), fn, nil, call.Args)
		} else {
			c.env.errors.undeclaredReference(c.location(e), c.container, call.Function)
		}
	} else {
		// Check whether the target is actually a qualified name for a static function.
		if qname, found := asQualifiedName(call.Target); found {
			fn := c.env.LookupFunction(c.container, qname+"."+call.Function)
			if fn != nil {
				resolution = c.resolveOverload(c.location(e), fn, nil, call.Args)
			}
		}

		if resolution == nil {
			// Regular instance call.
			c.check(call.Target)

			if fn := c.env.LookupFunction(c.container, call.Function); fn != nil {
				resolution = c.resolveOverload(c.location(e), fn, call.Target, call.Args)
			} else {
				c.env.errors.undeclaredReference(c.location(e), c.container, call.Function)
			}
		}
	}

	if resolution != nil {
		c.setType(e, resolution.Type)
		c.setReference(e, resolution.Reference)
	} else {
		c.setType(e, types.Error)
	}
}

func (c *checker) resolveOverload(
	loc common.Location,
	fn *semantics.Function, target *expr.Expr, args []*expr.Expr) *overloadResolution {

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

func (c *checker) checkCreateList(e *expr.Expr) {
	create := e.GetListExpr()
	var elemType types.Type = nil
	for _, e := range create.Elements {
		c.check(e)
		elemType = c.joinTypes(c.location(e), elemType, c.getType(e))
	}
	if elemType == nil {
		// If the list is empty, assign free type var to elem type.
		elemType = c.newTypeVar()
	}
	c.setType(e, types.NewList(elemType))
}

func (c *checker) checkCreateStruct(e *expr.Expr) {
	str := e.GetStructExpr()
	if str.MessageName != "" {
		c.checkCreateMessage(e)
	} else {
		c.checkCreateMap(e)
	}
}

func (c *checker) checkCreateMap(e *expr.Expr) {
	mapVal := e.GetStructExpr()
	var keyType types.Type = nil
	var valueType types.Type = nil
	for _, ent := range mapVal.GetEntries() {
		key := ent.GetMapKey()
		c.check(key)
		keyType = c.joinTypes(c.location(key), keyType, c.getType(key))

		c.check(ent.Value)
		valueType = c.joinTypes(c.location(ent.Value), valueType, c.getType(ent.Value))
	}
	if keyType == nil {
		// If the map is empty, assign free type variables to typeKey and value type.
		keyType = c.newTypeVar()
		valueType = c.newTypeVar()
	}
	c.setType(e, types.NewMap(keyType, valueType))
}

func (c *checker) checkCreateMessage(e *expr.Expr) {
	msgVal := e.GetStructExpr()
	// Determine the type of the message.
	messageType := types.Error
	decl := c.env.LookupIdent(c.container, msgVal.MessageName)
	if decl == nil {
		c.env.errors.undeclaredReference(c.location(e), c.container, msgVal.MessageName)
		return
	}
	c.setReference(e, semantics.NewIdentReference(decl.Name(), nil))
	if decl.Type().Kind() != types.KindError {
		if decl.Type().Kind() != types.KindType {
			c.env.errors.notAType(c.location(e), decl.Type())
		} else {
			messageType = decl.Type().(*types.TypeType).Target()
			if messageType.Kind() != types.KindMessage {
				c.env.errors.notAMessageType(c.location(e), messageType)
				messageType = types.Error
			}
		}
	}
	c.setType(e, messageType)

	// Check the field initializers.
	for _, ent := range msgVal.GetEntries() {
		field := ent.GetFieldKey()
		value := ent.Value
		c.check(value)

		fieldType := types.Error
		if mt, ok := messageType.(*types.MessageType); ok {
			if t, found := c.lookupFieldType(c.locationById(ent.Id), mt, field); found {
				fieldType = t.Type
			}
		}
		if !c.isAssignable(fieldType, c.getType(value)) {
			c.env.errors.fieldTypeMismatch(c.locationById(ent.Id), field, fieldType, c.getType(value))
		}
	}
}

func (c *checker) checkComprehension(e *expr.Expr) {
	comp := e.GetComprehensionExpr()
	c.check(comp.IterRange)
	c.check(comp.AccuInit)
	accuType := c.getType(comp.AccuInit)
	rangeType := c.getType(comp.IterRange)
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
		c.env.errors.notAComprehensionRange(c.location(comp.IterRange), rangeType)
	}

	c.env.enterScope()
	c.env.AddIdent(semantics.NewIdent(comp.AccuVar, accuType, nil))
	// Declare iteration variable on inner scope.
	c.env.enterScope()
	c.env.AddIdent(semantics.NewIdent(comp.IterVar, varType, nil))
	c.check(comp.LoopCondition)
	c.assertType(comp.LoopCondition, types.Bool)
	c.check(comp.LoopStep)
	c.assertType(comp.LoopStep, accuType)
	// Forget iteration variable, as result expression must only depend on accu.
	c.env.exitScope()
	c.check(comp.Result)
	c.env.exitScope()
	c.setType(e, c.getType(comp.Result))
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

func (c *checker) setType(e *expr.Expr, t types.Type) {
	if old, found := c.types[e.Id]; found && !old.Equals(t) {
		panic(fmt.Sprintf("(Incompatible) Type already exists for expression: %v(%d) old:%v, new:%v", e, e.Id, old, t))
	}

	c.types[e.Id] = t
}

func (c *checker) getType(e *expr.Expr) types.Type {
	return c.types[e.Id]
}

func (c *checker) setReference(e *expr.Expr, r semantics.Reference) {
	if old, found := c.references[e.Id]; found && !old.Equals(r) {
		panic(fmt.Sprintf("Reference already exists for expression: %v(%d) old:%v, new:%v", e, e.Id, old, r))
	}

	c.references[e.Id] = r
}

func (c *checker) assertType(e *expr.Expr, t types.Type) {
	if !c.isAssignable(t, c.getType(e)) {
		c.env.errors.typeMismatch(c.location(e), t, c.getType(e))
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

func (c *checker) location(e *expr.Expr) common.Location {
	return c.locationById(e.Id)
}

func (c *checker) locationById(id int64) common.Location {
	positions := c.sourceInfo.GetPositions()
	var line = 1
	var col = 0
	if offset, found := positions[id]; found {
		col = int(offset)
		for _, lineOffset := range c.sourceInfo.LineOffsets {
			if lineOffset < offset {
				line += 1
				col = int(offset - lineOffset)
			} else {
				break;
			}
		}
		return common.NewLocation(line, col)
	}
	return common.NoLocation
}
