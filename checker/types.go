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
	"strings"

	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	declspb "github.com/google/cel-go/checker/decls"
	protopb "github.com/golang/protobuf/proto"
)

const (
	kindUnknown = iota + 1
	kindError
	kindFunction
	kindDyn
	kindPrimitive
	kindWellKnown
	kindWrapper
	kindNull
	kindAbstract // TODO: Update the checkedpb.proto to include abstract
	kindType
	kindList
	kindMap
	kindObject
	kindTypeParam
)

// FormatCheckedType converts a type message into a string representation.
func FormatCheckedType(t *checkedpb.Type) string {
	switch kindOf(t) {
	case kindPrimitive:
		switch t.GetPrimitive() {
		case checkedpb.Type_UINT64:
			return "uint"
		case checkedpb.Type_INT64:
			return "int"
		}
		return strings.Trim(strings.ToLower(t.GetPrimitive().String()), " ")
	case kindFunction:
		return formatFunction(t.GetFunction().GetResultType(),
			t.GetFunction().GetArgTypes(),
			false)
	case kindWrapper:
		return fmt.Sprintf("wrapper(%s)",
			FormatCheckedType(declspb.NewPrimitiveType(t.GetWrapper())))
	case kindObject:
		return t.GetMessageType()
	case kindList:
		return fmt.Sprintf("list(%s)", FormatCheckedType(t.GetListType().ElemType))
	case kindMap:
		return fmt.Sprintf("map(%s, %s)",
			FormatCheckedType(t.GetMapType().KeyType),
			FormatCheckedType(t.GetMapType().ValueType))
	case kindNull:
		return "null"
	case kindDyn:
		return "dyn"
	case kindType:
		if t.GetType() == nil {
			return "type"
		}
		return fmt.Sprintf("type(%s)", FormatCheckedType(t.GetType()))
	case kindError:
		return "!error!"
	}
	return t.String()
}

/**
 * Check whether one type is equal or less specific than the other one. A type is less specific if
 * it matches the other type using the DYN type.
 */
func isEqualOrLessSpecific(t1 *checkedpb.Type, t2 *checkedpb.Type) bool {
	kind1, kind2 := kindOf(t1), kindOf(t2)
	if kind1 == kindDyn || kind1 == kindTypeParam {
		return true
	}
	if kind2 == kindDyn || kind2 == kindTypeParam {
		return false
	}
	if kind1 != kind2 {
		return false
	}

	switch kind1 {
	case kindObject, kindPrimitive, kindWellKnown, kindWrapper:
		return protopb.Equal(t1, t2)
	case kindList:
		return isEqualOrLessSpecific(t1.GetListType().ElemType, t2.GetListType().ElemType)
	case kindMap:
		m1 := t1.GetMapType()
		m2 := t2.GetMapType()
		return isEqualOrLessSpecific(m1.KeyType, m2.KeyType) &&
			isEqualOrLessSpecific(m1.KeyType, m2.KeyType)
	case kindFunction:
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

func internalIsAssignableList(m *mapping, l1 []*checkedpb.Type, l2 []*checkedpb.Type) bool {
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

func internalIsAssignable(m *mapping, t1 *checkedpb.Type, t2 *checkedpb.Type) bool {
	// Process type parameters.
	kind1, kind2 := kindOf(t1), kindOf(t2)
	if kind2 == kindTypeParam {
		if t2Sub, found := m.find(t2); found {
			// Adjust the existing substitution to a more common type if possible. This is sound
			// because any previous substitution will be compatible with the common type. This
			// deals with the case the we have e.g. A -> int assigned, but now encounter a test
			// against DYN, and want to widen A to DYN.
			if isEqualOrLessSpecific(t1, t2Sub) && notReferencedIn(t2, t1) {
				m.add(t2, t1)
				return true
			}
			// Continue regular process with the assignment for type2.
			return internalIsAssignable(m, t1, t2Sub)
		}
		if notReferencedIn(t2, t1) {
			m.add(t2, t1)
			return true
		}
	}

	if kind1 == kindTypeParam {
		// For the lower type bound, we currently do not perform adjustment. The restricted
		// way we use type parameters in lower type bounds, it is not necessary, but may
		// become if we generalize type unification.
		if t1Sub, found := m.find(t1); found {
			return internalIsAssignable(m, t1Sub, t2)
		}
		if notReferencedIn(t1, t2) {
			m.add(t1, t2)
			return true
		}
	}

	if kind1 == kindDyn || kind1 == kindError {
		return true
	}
	if kind2 == kindDyn || kind2 == kindError {
		return true
	}
	if kind1 == kindNull && isNullable(kind2) {
		return true
	}
	// Unwrap box types
	if kind1 == kindWrapper {
		return internalIsAssignable(m, declspb.NewPrimitiveType(t1.GetWrapper()), t2)
	}
	// Finally check equality and type args recursively.
	if kind1 != kind2 {
		return false
	}

	switch kind1 {
	case kindPrimitive, kindWellKnown, kindObject:
		return protopb.Equal(t1, t2)
	case kindList:
		return internalIsAssignable(m, t1.GetListType().ElemType, t2.GetListType().ElemType)
	case kindMap:
		m1 := t1.GetMapType()
		m2 := t2.GetMapType()
		return internalIsAssignableList(m,
			[]*checkedpb.Type{m1.KeyType, m1.ValueType},
			[]*checkedpb.Type{m2.KeyType, m2.ValueType})
	case kindFunction:
		fn1 := t1.GetFunction()
		fn2 := t2.GetFunction()
		return internalIsAssignableList(m,
			append(fn1.ArgTypes, fn1.ResultType),
			append(fn2.ArgTypes, fn2.ResultType))
	case kindType:
		// A type is a type is a type, any additional parameterization of the
		// type cannot affect method resolution or assignability.
		return true
	default:
		return false
	}
}

func isAssignable(m *mapping, t1 *checkedpb.Type, t2 *checkedpb.Type) *mapping {
	mCopy := m.copy()
	if internalIsAssignable(mCopy, t1, t2) {
		return mCopy
	}
	return nil
}

func isAssignableList(m *mapping, l1 []*checkedpb.Type, l2 []*checkedpb.Type) *mapping {
	mCopy := m.copy()
	if internalIsAssignableList(mCopy, l1, l2) {
		return mCopy
	}
	return nil
}

func isNullable(kind int) bool {
	switch kind {
	case kindObject, kindWrapper, kindWellKnown:
		return true
	default:
		return false
	}
}

func kindOf(t *checkedpb.Type) int {
	if t == nil || t.TypeKind == nil {
		return kindUnknown
	}
	switch t.TypeKind.(type) {
	case *checkedpb.Type_Error:
		return kindError
	case *checkedpb.Type_Function:
		return kindFunction
	case *checkedpb.Type_Dyn:
		return kindDyn
	case *checkedpb.Type_Primitive:
		return kindPrimitive
	case *checkedpb.Type_WellKnown:
		return kindWellKnown
	case *checkedpb.Type_Wrapper:
		return kindWrapper
	case *checkedpb.Type_Null:
		return kindNull
	case *checkedpb.Type_Type:
		return kindType
	case *checkedpb.Type_ListType_:
		return kindList
	case *checkedpb.Type_MapType_:
		return kindMap
	case *checkedpb.Type_MessageType:
		return kindObject
	case *checkedpb.Type_TypeParam:
		return kindTypeParam
	}
	return kindUnknown
}

/** Returns the more general of two types which are known to unify. */
func mostGeneral(t1 *checkedpb.Type, t2 *checkedpb.Type) *checkedpb.Type {
	if isEqualOrLessSpecific(t1, t2) {
		return t1
	}
	return t2
}

/**
 * Apply substitution to given type, replacing all direct and indirect occurrences of bound type
 * parameters. Unbound type parameters are replaced by DYN if typeParamToDyn is true.
 */
func substitute(m *mapping, t *checkedpb.Type, typeParamToDyn bool) *checkedpb.Type {
	if tSub, found := m.find(t); found {
		return substitute(m, tSub, typeParamToDyn)
	}
	kind := kindOf(t)
	if typeParamToDyn && kind == kindTypeParam {
		return declspb.Dyn
	}
	switch kind {
	case kindType:
		if t.GetType() != nil {
			return declspb.NewTypeType(substitute(m, t.GetType(), typeParamToDyn))
		}
		return t
	case kindList:
		return declspb.NewListType(substitute(m, t.GetListType().ElemType, typeParamToDyn))
	case kindMap:
		mt := t.GetMapType()
		return declspb.NewMapType(substitute(m, mt.KeyType, typeParamToDyn),
			substitute(m, mt.ValueType, typeParamToDyn))
	case kindFunction:
		fn := t.GetFunction()
		rt := substitute(m, fn.ResultType, typeParamToDyn)
		args := make([]*checkedpb.Type, len(fn.ArgTypes))
		for i, a := range fn.ArgTypes {
			args[i] = substitute(m, a, typeParamToDyn)
		}
		return declspb.NewFunctionType(rt, args...)
	default:
		return t
	}
}

func notReferencedIn(t *checkedpb.Type, withinType *checkedpb.Type) bool {
	if protopb.Equal(t, withinType) {
		return false
	}
	withinKind := kindOf(withinType)
	switch withinKind {
	case kindWrapper:
		return notReferencedIn(t, declspb.NewPrimitiveType(withinType.GetWrapper()))
	case kindList:
		return notReferencedIn(t, withinType.GetListType().ElemType)
	case kindMap:
		m := withinType.GetMapType()
		return notReferencedIn(t, m.KeyType) && notReferencedIn(t, m.ValueType)
	case kindFunction:
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

func typeKey(t *checkedpb.Type) string {
	return fmt.Sprintf("%v:%v", kindOf(t), t.String())
}
