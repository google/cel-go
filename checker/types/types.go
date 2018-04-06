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

// Package types declares type inspection and type assignability for use with
// type checking and runtime type resolution.
package types

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/struct"
	"github.com/google/cel-spec/proto/checked/v1/checked"
	"strings"
)

const (
	KindUnknown = iota + 1
	KindError
	KindFunction
	KindDyn
	KindPrimitive
	KindWellKnown
	KindWrapper
	KindNull
	KindAbstract // TODO: Update the checked.proto to include abstract
	KindType
	KindList
	KindMap
	KindObject
	KindTypeParam
)

var (
	// Commonly used types.
	Error = &checked.Type{
		TypeKind: &checked.Type_Error{
			Error: &empty.Empty{}}}
	Dyn = &checked.Type{
		TypeKind: &checked.Type_Dyn{
			Dyn: &empty.Empty{}}}
	Null = &checked.Type{
		TypeKind: &checked.Type_Null{
			Null: structpb.NullValue_NULL_VALUE}}
	Bool   = NewPrimitive(checked.Type_BOOL)
	Bytes  = NewPrimitive(checked.Type_BYTES)
	String = NewPrimitive(checked.Type_STRING)
	Double = NewPrimitive(checked.Type_DOUBLE)
	Int64  = NewPrimitive(checked.Type_INT64)
	Uint64 = NewPrimitive(checked.Type_UINT64)

	// Well-known types.
	// TODO: Replace with an abstract type registry.
	Any       = NewWellKnown(checked.Type_ANY)
	Duration  = NewWellKnown(checked.Type_DURATION)
	Timestamp = NewWellKnown(checked.Type_TIMESTAMP)
)

func NewFunction(resultType *checked.Type, argTypes ...*checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Function{
			Function: &checked.Type_FunctionType{
				ResultType: resultType,
				ArgTypes:   argTypes}}}
}

func NewTypeParam(name string) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_TypeParam{
			TypeParam: name}}
}

func NewType(nested *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Type{
			Type: nested}}
}

func NewPrimitive(primitive checked.Type_PrimitiveType) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_Primitive{
			Primitive: primitive}}
}

func NewWellKnown(wellKnown checked.Type_WellKnownType) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_WellKnown{
			WellKnown: wellKnown}}
}

func NewWrapper(wrapped *checked.Type) *checked.Type {
	primitive := wrapped.GetPrimitive()
	if primitive == checked.Type_PRIMITIVE_TYPE_UNSPECIFIED {
		// TODO: return an error
		panic("Wrapped type must be a primitive")
	}
	return &checked.Type{
		TypeKind: &checked.Type_Wrapper{
			Wrapper: primitive}}
}

func NewList(elem *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_ListType_{
			ListType: &checked.Type_ListType{
				ElemType: elem}}}
}

func NewMap(key *checked.Type, value *checked.Type) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_MapType_{
			MapType: &checked.Type_MapType{
				KeyType:   key,
				ValueType: value}}}
}

func NewObject(typeName string) *checked.Type {
	return &checked.Type{
		TypeKind: &checked.Type_MessageType{
			MessageType: typeName}}
}

func KindOf(t *checked.Type) int {
	switch t.TypeKind.(type) {
	case *checked.Type_Error:
		return KindError
	case *checked.Type_Function:
		return KindFunction
	case *checked.Type_Dyn:
		return KindDyn
	case *checked.Type_Primitive:
		return KindPrimitive
	case *checked.Type_WellKnown:
		return KindWellKnown
	case *checked.Type_Wrapper:
		return KindWrapper
	case *checked.Type_Null:
		return KindNull
	case *checked.Type_Type:
		return KindType
	case *checked.Type_ListType_:
		return KindList
	case *checked.Type_MapType_:
		return KindMap
	case *checked.Type_MessageType:
		return KindObject
	case *checked.Type_TypeParam:
		return KindTypeParam
	}
	return KindUnknown
}

/** Returns the more general of two types which are known to unify. */
func MostGeneral(t1 *checked.Type, t2 *checked.Type) *checked.Type {
	if isEqualOrLessSpecific(t1, t2) {
		return t1
	}
	return t2
}

func IsAssignable(m *Mapping, t1 *checked.Type, t2 *checked.Type) *Mapping {
	mCopy := m.Copy()
	if internalIsAssignable(mCopy, t1, t2) {
		return mCopy
	}
	return nil
}

func IsAssignableList(m *Mapping, l1 []*checked.Type, l2 []*checked.Type) *Mapping {
	mCopy := m.Copy()
	if internalIsAssignableList(mCopy, l1, l2) {
		return mCopy
	}
	return nil
}

/**
 * Apply substitution to given type, replacing all direct and indirect occurrences of bound type
 * parameters. Unbound type parameters are replaced by DYN if typeParamToDyn is true.
 */
func Substitute(m *Mapping, t *checked.Type, typeParamToDyn bool) *checked.Type {
	if tSub, found := m.Find(t); found {
		return Substitute(m, tSub, typeParamToDyn)
	}
	kind := KindOf(t)
	if typeParamToDyn && kind == KindTypeParam {
		return Dyn
	}
	switch kind {
	case KindType:
		return NewType(Substitute(m, t.GetType(), typeParamToDyn))
	case KindList:
		return NewList(Substitute(m, t.GetListType().ElemType, typeParamToDyn))
	case KindMap:
		mt := t.GetMapType()
		return NewMap(Substitute(m, mt.KeyType, typeParamToDyn),
			Substitute(m, mt.ValueType, typeParamToDyn))
	case KindFunction:
		fn := t.GetFunction()
		rt := Substitute(m, fn.ResultType, typeParamToDyn)
		args := make([]*checked.Type, len(fn.ArgTypes))
		for i, a := range fn.ArgTypes {
			args[i] = Substitute(m, a, typeParamToDyn)
		}
		return NewFunction(rt, args...)
	default:
		return t
	}
}

func internalIsAssignableList(m *Mapping, l1 []*checked.Type, l2 []*checked.Type) bool {
	if len(l1) != len(l2) {
		return false
	}
	for i, t1 := range l1 {
		if !internalIsAssignable(m, t1, l2[i]) {
			return false
		}
	}
	return true
}

func internalIsAssignable(m *Mapping, t1 *checked.Type, t2 *checked.Type) bool {
	// Process type parameters.
	kind1, kind2 := KindOf(t1), KindOf(t2)
	if kind2 == KindTypeParam {
		if t2Sub, found := m.Find(t2); found {
			// Adjust the existing substitution to a more common type if possible. This is sound
			// because any previous substitution will be compatible with the common type. This
			// deals with the case the we have e.g. A -> int assigned, but now encounter a test
			// against DYN, and want to widen A to DYN.
			if isEqualOrLessSpecific(t1, t2Sub) && notReferencedIn(t2, t1) {
				m.Add(t2, t1)
				return true
			} else {
				// Continue regular process with the assignment for type2.
				return internalIsAssignable(m, t1, t2Sub)
			}
		} else if notReferencedIn(t2, t1) {
			m.Add(t2, t1)
			return true
		}
	}

	if kind1 == KindTypeParam {
		// For the lower type bound, we currently do not perform adjustment. The restricted
		// way we use type parameters in lower type bounds, it is not necessary, but may
		// become if we generalize type unification.
		if t1Sub, found := m.Find(t1); found {
			return internalIsAssignable(m, t1Sub, t2)
		} else if notReferencedIn(t1, t2) {
			m.Add(t1, t2)
			return true
		}
	}

	if kind1 == KindDyn || kind1 == KindError {
		return true
	}
	if kind2 == KindDyn || kind2 == KindError {
		return true
	}
	if kind1 == KindNull && isNullable(kind2) {
		return true
	}
	// Unwrap box types
	if kind1 == KindWrapper {
		return internalIsAssignable(m, NewPrimitive(t1.GetWrapper()), t2)
	}
	// Finally check equality and type args recursively.
	if kind1 != kind2 {
		return false
	}

	switch kind1 {
	case KindPrimitive, KindWellKnown, KindObject:
		return proto.Equal(t1, t2)
	case KindType:
		return internalIsAssignable(m, t1.GetType(), t2.GetType())
	case KindList:
		return internalIsAssignable(m, t1.GetListType().ElemType, t2.GetListType().ElemType)
	case KindMap:
		m1 := t1.GetMapType()
		m2 := t2.GetMapType()
		return internalIsAssignableList(m,
			[]*checked.Type{m1.KeyType, m1.ValueType},
			[]*checked.Type{m2.KeyType, m2.ValueType})
	case KindFunction:
		fn1 := t1.GetFunction()
		fn2 := t2.GetFunction()
		return internalIsAssignableList(m,
			append(fn1.ArgTypes, fn1.ResultType),
			append(fn2.ArgTypes, fn2.ResultType))
	default:
		return false
	}
}

/**
 * Check whether one type is equal or less specific than the other one. A type is less specific if
 * it matches the other type using the DYN type.
 */
func isEqualOrLessSpecific(t1 *checked.Type, t2 *checked.Type) bool {
	kind1, kind2 := KindOf(t1), KindOf(t2)
	if kind1 == KindDyn || kind1 == KindTypeParam {
		return true
	}
	if kind2 == KindDyn || kind2 == KindTypeParam {
		return false
	}
	if kind1 != kind2 {
		return false
	}

	switch kind1 {
	case KindObject, KindPrimitive, KindWellKnown, KindWrapper:
		return proto.Equal(t1, t2)
	case KindType:
		return isEqualOrLessSpecific(t1.GetType(), t2.GetType())
	case KindList:
		return isEqualOrLessSpecific(t1.GetListType().ElemType, t2.GetListType().ElemType)
	case KindMap:
		m1 := t1.GetMapType()
		m2 := t2.GetMapType()
		return isEqualOrLessSpecific(m1.KeyType, m2.KeyType) &&
			isEqualOrLessSpecific(m1.KeyType, m2.KeyType)
	case KindFunction:
		fn1 := t1.GetFunction()
		fn2 := t2.GetFunction()
		if len(fn1.ArgTypes) != len(fn2.ArgTypes) {
			return false
		}
		if !isEqualOrLessSpecific(fn1.ResultType, fn2.ResultType) {
			return false
		}
		for i, a1 := range fn1.ArgTypes {
			if !isEqualOrLessSpecific(a1, fn2.ArgTypes[i]) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func isNullable(kind int) bool {
	switch kind {
	case KindObject, KindWrapper, KindWellKnown:
		return true
	default:
		return false
	}
}

func notReferencedIn(t *checked.Type, withinType *checked.Type) bool {
	if proto.Equal(t, withinType) {
		return false
	}
	withinKind := KindOf(withinType)
	switch withinKind {
	case KindWrapper:
		return notReferencedIn(t, NewPrimitive(withinType.GetWrapper()))
	case KindType:
		return notReferencedIn(t, withinType.GetType())
	case KindList:
		return notReferencedIn(t, withinType.GetListType().ElemType)
	case KindMap:
		m := withinType.GetMapType()
		return notReferencedIn(t, m.KeyType) && notReferencedIn(t, m.ValueType)
	case KindFunction:
		fn := withinType.GetFunction()
		types := append(fn.ArgTypes, fn.ResultType)
		for _, a := range types {
			if !notReferencedIn(t, a) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func typeKey(t *checked.Type) string {
	return fmt.Sprintf("%v:%v", KindOf(t), t.String())
}

func FormatType(t *checked.Type) string {
	switch KindOf(t) {
	case KindPrimitive:
		switch t.GetPrimitive() {
		case checked.Type_UINT64:
			return "uint"
		case checked.Type_INT64:
			return "int"
		}
		return strings.Trim(strings.ToLower(t.GetPrimitive().String()), " ")
	case KindWrapper:
		return fmt.Sprintf("wrapper(%s)", FormatType(NewPrimitive(t.GetWrapper())))
	case KindObject:
		return t.GetMessageType()
	case KindList:
		return fmt.Sprintf("list(%s)", FormatType(t.GetListType().ElemType))
	case KindMap:
		return fmt.Sprintf("map(%s, %s)",
			FormatType(t.GetMapType().KeyType),
			FormatType(t.GetMapType().ValueType))
	case KindNull:
		return "null"
	case KindDyn:
		return "dyn"
	case KindType:
		return fmt.Sprintf("type(%s)", FormatType(t.GetType()))
	case KindError:
		return "!error!"
	}
	return t.String()
}
