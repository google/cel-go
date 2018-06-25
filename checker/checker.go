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

// Package checker defines functions to type-checked a parsed expression
// against a set of identifier and function declarations.
package checker

import (
	"fmt"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-spec/proto/checked/v1/checked"
	expr "github.com/google/cel-spec/proto/v1/syntax"
)

type checker struct {
	env                *Env
	mappings           *mapping
	freeTypeVarCounter int
	sourceInfo         *expr.SourceInfo

	types      map[int64]*checked.Type
	references map[int64]*checked.Reference
}

func Check(parsedExpr *expr.ParsedExpr, env *Env) *checked.CheckedExpr {
	c := checker{
		env:                env,
		mappings:           newMapping(),
		freeTypeVarCounter: 0,
		sourceInfo:         parsedExpr.GetSourceInfo(),

		types:      make(map[int64]*checked.Type),
		references: make(map[int64]*checked.Reference),
	}
	c.check(parsedExpr.GetExpr())

	// Walk over the final type map substituting any type parameters either by their bound value or
	// by DYN.
	m := make(map[int64]*checked.Type)
	for k, v := range c.types {
		m[k] = substitute(c.mappings, v, true)
	}

	return &checked.CheckedExpr{
		Expr:         parsedExpr.GetExpr(),
		SourceInfo:   parsedExpr.GetSourceInfo(),
		TypeMap:      m,
		ReferenceMap: c.references,
	}
}

func (c *checker) check(e *expr.Expr) {
	if e == nil {
		return
	}

	switch e.ExprKind.(type) {
	case *expr.Expr_LiteralExpr:
		literal := e.GetLiteralExpr()
		switch literal.LiteralKind.(type) {
		case *expr.Literal_BoolValue:
			c.checkBoolLiteral(e)
		case *expr.Literal_BytesValue:
			c.checkBytesLiteral(e)
		case *expr.Literal_DoubleValue:
			c.checkDoubleLiteral(e)
		case *expr.Literal_Int64Value:
			c.checkInt64Literal(e)
		case *expr.Literal_NullValue:
			c.checkNullLiteral(e)
		case *expr.Literal_StringValue:
			c.checkStringLiteral(e)
		case *expr.Literal_Uint64Value:
			c.checkUint64Literal(e)
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

func (c *checker) checkInt64Literal(e *expr.Expr) {
	c.setType(e, decls.Int)
}

func (c *checker) checkUint64Literal(e *expr.Expr) {
	c.setType(e, decls.Uint)
}

func (c *checker) checkStringLiteral(e *expr.Expr) {
	c.setType(e, decls.String)
}

func (c *checker) checkBytesLiteral(e *expr.Expr) {
	c.setType(e, decls.Bytes)
}

func (c *checker) checkDoubleLiteral(e *expr.Expr) {
	c.setType(e, decls.Double)
}

func (c *checker) checkBoolLiteral(e *expr.Expr) {
	c.setType(e, decls.Bool)
}

func (c *checker) checkNullLiteral(e *expr.Expr) {
	c.setType(e, decls.Null)
}

func (c *checker) checkIdent(e *expr.Expr) {
	identExpr := e.GetIdentExpr()
	if ident := c.env.LookupIdent(identExpr.Name); ident != nil {
		c.setType(e, ident.GetIdent().Type)
		c.setReference(e, newIdentReference(ident.Name, ident.GetIdent().Value))
		return
	}

	c.setType(e, decls.Error)
	c.env.errors.undeclaredReference(
		c.location(e), c.env.packager.Package(), identExpr.Name)
}

func (c *checker) checkSelect(e *expr.Expr) {
	sel := e.GetSelectExpr()
	// Before traversing down the tree, try to interpret as qualified name.
	qname, found := toQualifiedName(e)
	if found {
		ident := c.env.LookupIdent(qname)
		if ident != nil {
			if sel.TestOnly {
				c.env.errors.expressionDoesNotSelectField(c.location(e))
				c.setType(e, decls.Bool)
			} else {
				c.setType(e, ident.GetIdent().Type)
				c.setReference(e,
					newIdentReference(ident.Name, ident.GetIdent().Value))
			}
			return
		}
	}

	// Interpret as field selection, first traversing down the operand.
	c.check(sel.Operand)
	targetType := c.getType(sel.Operand)
	resultType := decls.Error

	switch kindOf(targetType) {
	case kindError, kindDyn:
		resultType = decls.Dyn

	case kindObject:
		messageType := targetType
		if fieldType, found := c.lookupFieldType(c.location(e), messageType, sel.Field); found {
			resultType = fieldType.Type
			if sel.TestOnly && !fieldType.SupportsPresence {
				c.env.errors.fieldDoesNotSupportPresenceCheck(c.location(e), sel.Field)
			}
		}

	case kindMap:
		mapType := targetType.GetMapType()
		resultType = mapType.ValueType

	default:
		c.env.errors.typeDoesNotSupportFieldSelection(c.location(e), targetType)
	}

	if sel.TestOnly {
		resultType = decls.Bool
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
		if fn := c.env.LookupFunction(call.Function); fn != nil {
			resolution = c.resolveOverload(c.location(e), fn, nil, call.Args)
		} else {
			c.env.errors.undeclaredReference(
				c.location(e), c.env.packager.Package(), call.Function)
		}
	} else {
		// Check whether the target is actually a qualified name for a static function.
		if qname, found := toQualifiedName(call.Target); found {
			fn := c.env.LookupFunction(qname + "." + call.Function)
			if fn != nil {
				resolution = c.resolveOverload(c.location(e), fn, nil, call.Args)
			}
		}

		if resolution == nil {
			// Regular instance call.
			c.check(call.Target)

			if fn := c.env.LookupFunction(call.Function); fn != nil {
				resolution = c.resolveOverload(c.location(e), fn, call.Target, call.Args)
			} else {
				c.env.errors.undeclaredReference(
					c.location(e), c.env.packager.Package(), call.Function)
			}
		}
	}

	if resolution != nil {
		c.setType(e, resolution.Type)
		c.setReference(e, resolution.Reference)
	} else {
		c.setType(e, decls.Error)
	}
}

func (c *checker) resolveOverload(
	loc common.Location,
	fn *checked.Decl, target *expr.Expr, args []*expr.Expr) *overloadResolution {

	var argTypes []*checked.Type
	if target != nil {
		argTypes = append(argTypes, c.getType(target))
	}
	for _, arg := range args {
		argTypes = append(argTypes, c.getType(arg))
	}

	var resultType *checked.Type = nil
	var checkedRef *checked.Reference = nil
	for _, overload := range fn.GetFunction().Overloads {
		if (target == nil && overload.IsInstanceFunction) ||
			(target != nil && !overload.IsInstanceFunction) {
			// not a compatible call style.
			continue
		}

		overloadType := decls.NewFunctionType(overload.ResultType, overload.Params...)
		if len(overload.TypeParams) > 0 {
			// Instantiate overload's type with fresh type variables.
			substitutions := newMapping()
			for _, typePar := range overload.TypeParams {
				substitutions.add(decls.NewTypeParamType(typePar), c.newTypeVar())
			}

			overloadType = substitute(substitutions, overloadType, false)
		}

		candidateArgTypes := overloadType.GetFunction().ArgTypes
		if c.isAssignableList(argTypes, candidateArgTypes) {
			if checkedRef == nil {
				checkedRef = newFunctionReference(overload.OverloadId)
			} else {
				checkedRef.OverloadId = append(checkedRef.OverloadId, overload.OverloadId)
			}

			if resultType == nil {
				// First matching overload, determines result type.
				resultType = substitute(c.mappings,
					overloadType.GetFunction().ResultType,
					false)
			} else {
				// More than one matching overload, narrow result type to DYN.
				resultType = decls.Dyn
			}

		}
	}

	if resultType == nil {
		c.env.errors.noMatchingOverload(loc, fn.Name, argTypes, target != nil)
		resultType = decls.Error
		return nil
	}

	return newResolution(checkedRef, resultType)
}

func (c *checker) checkCreateList(e *expr.Expr) {
	create := e.GetListExpr()
	var elemType *checked.Type = nil
	for _, e := range create.Elements {
		c.check(e)
		elemType = c.joinTypes(c.location(e), elemType, c.getType(e))
	}
	if elemType == nil {
		// If the list is empty, assign free type var to elem type.
		elemType = c.newTypeVar()
	}
	c.setType(e, decls.NewListType(elemType))
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
	var keyType *checked.Type = nil
	var valueType *checked.Type = nil
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
	c.setType(e, decls.NewMapType(keyType, valueType))
}

func (c *checker) checkCreateMessage(e *expr.Expr) {
	msgVal := e.GetStructExpr()
	// Determine the type of the message.
	messageType := decls.Error
	decl := c.env.LookupIdent(msgVal.MessageName)
	if decl == nil {
		c.env.errors.undeclaredReference(
			c.location(e), c.env.packager.Package(), msgVal.MessageName)
		return
	}

	c.setReference(e, newIdentReference(decl.Name, nil))
	ident := decl.GetIdent()
	identKind := kindOf(ident.Type)
	if identKind != kindError {
		if identKind != kindType {
			c.env.errors.notAType(c.location(e), ident.Type)
		} else {
			messageType = ident.Type.GetType()
			if kindOf(messageType) != kindObject {
				c.env.errors.notAMessageType(c.location(e), messageType)
				messageType = decls.Error
			}
		}
	}
	c.setType(e, messageType)

	// Check the field initializers.
	for _, ent := range msgVal.GetEntries() {
		field := ent.GetFieldKey()
		value := ent.Value
		c.check(value)

		fieldType := decls.Error
		if t, found := c.lookupFieldType(c.locationById(ent.Id), messageType, field); found {
			fieldType = t.Type
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
	var varType *checked.Type

	switch kindOf(rangeType) {
	case kindList:
		varType = rangeType.GetListType().ElemType
	case kindMap:
		// Ranges over the keys.
		varType = rangeType.GetMapType().KeyType
	case kindDyn, kindError:
		varType = decls.Dyn
	default:
		c.env.errors.notAComprehensionRange(c.location(comp.IterRange), rangeType)
	}

	c.env.enterScope()
	c.env.Add(decls.NewIdent(comp.AccuVar, accuType, nil))
	// Declare iteration variable on inner scope.
	c.env.enterScope()
	c.env.Add(decls.NewIdent(comp.IterVar, varType, nil))
	c.check(comp.LoopCondition)
	c.assertType(comp.LoopCondition, decls.Bool)
	c.check(comp.LoopStep)
	c.assertType(comp.LoopStep, accuType)
	// Forget iteration variable, as result expression must only depend on accu.
	c.env.exitScope()
	c.check(comp.Result)
	c.env.exitScope()
	c.setType(e, c.getType(comp.Result))
}

// Checks compatibility of joined types, and returns the most general common type.
func (c *checker) joinTypes(loc common.Location, previous *checked.Type, current *checked.Type) *checked.Type {
	if previous == nil {
		return current
	}
	if !c.isAssignable(previous, current) {
		c.env.errors.aggregateTypeMismatch(loc, previous, current)
		return previous
	}
	return mostGeneral(previous, current)
}

func (c *checker) newTypeVar() *checked.Type {
	id := c.freeTypeVarCounter
	c.freeTypeVarCounter++
	return decls.NewTypeParamType(fmt.Sprintf("_var%d", id))
}

func (c *checker) isAssignable(t1 *checked.Type, t2 *checked.Type) bool {
	subs := isAssignable(c.mappings, t1, t2)
	if subs != nil {
		c.mappings = subs
		return true
	}

	return false
}

func (c *checker) isAssignableList(l1 []*checked.Type, l2 []*checked.Type) bool {
	subs := isAssignableList(c.mappings, l1, l2)
	if subs != nil {
		c.mappings = subs
		return true
	}

	return false
}

func (c *checker) lookupFieldType(l common.Location, messageType *checked.Type, fieldName string) (*ref.FieldType, bool) {
	if _, found := c.env.typeProvider.FindType(messageType.GetMessageType()); !found {
		// This should not happen, anyway, report an error.
		c.env.errors.unexpectedFailedResolution(l, messageType.GetMessageType())
		return nil, false
	}

	if ft, found := c.env.typeProvider.FindFieldType(messageType, fieldName); found {
		return ft, found
	}

	c.env.errors.undefinedField(l, fieldName)
	return nil, false
}

func (c *checker) setType(e *expr.Expr, t *checked.Type) {
	if old, found := c.types[e.Id]; found && !proto.Equal(old, t) {
		panic(fmt.Sprintf("(Incompatible) Type already exists for expression: %v(%d) old:%v, new:%v", e, e.Id, old, t))
	}
	c.types[e.Id] = t
}

func (c *checker) getType(e *expr.Expr) *checked.Type {
	return c.types[e.Id]
}

func (c *checker) setReference(e *expr.Expr, r *checked.Reference) {
	if old, found := c.references[e.Id]; found && !proto.Equal(old, r) {
		panic(fmt.Sprintf("Reference already exists for expression: %v(%d) old:%v, new:%v", e, e.Id, old, r))
	}
	c.references[e.Id] = r
}

func (c *checker) assertType(e *expr.Expr, t *checked.Type) {
	if !c.isAssignable(t, c.getType(e)) {
		c.env.errors.typeMismatch(c.location(e), t, c.getType(e))
	}
}

type overloadResolution struct {
	Reference *checked.Reference
	Type      *checked.Type
}

func newResolution(checkedRef *checked.Reference, t *checked.Type) *overloadResolution {
	return &overloadResolution{
		Reference: checkedRef,
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
				break
			}
		}
		return common.NewLocation(line, col)
	}
	return common.NoLocation
}

func newIdentReference(name string, value *expr.Literal) *checked.Reference {
	return &checked.Reference{Name: name, Value: value}
}

func newFunctionReference(overloads ...string) *checked.Reference {
	return &checked.Reference{OverloadId: overloads}
}

// Attempt to interpret an expression as a qualified name. This traverses select and getIdent
// expression and returns the name they constitute, or null if the expression cannot be
// interpreted like this.
func toQualifiedName(e *expr.Expr) (string, bool) {
	switch e.ExprKind.(type) {
	case *expr.Expr_IdentExpr:
		i := e.GetIdentExpr()
		return i.Name, true
	case *expr.Expr_SelectExpr:
		s := e.GetSelectExpr()
		if qname, found := toQualifiedName(s.Operand); found {
			return qname + "." + s.Field, true
		}
	}
	return "", false
}
