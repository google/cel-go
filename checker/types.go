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
	"github.com/google/cel-go/semantics/types"
)

func isAssignable(m *Mapping, t1 types.Type, t2 types.Type) *Mapping {
	mCopy := m.Copy()
	if internalIsAssignable(mCopy, t1, t2) {
		return mCopy
	}
	return nil
}

func isAssignableList(m *Mapping, l1 []types.Type, l2 []types.Type) *Mapping {
	mCopy := m.Copy()
	if internalIsAssignableList(mCopy, l1, l2) {
		return mCopy
	}
	return nil
}

func internalIsAssignableList(m *Mapping, l1 []types.Type, l2 []types.Type) bool {
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

func internalIsAssignable(m *Mapping, t1 types.Type, t2 types.Type) bool {
	// Process type parameters.
	if t2.Kind() == types.KindTypeParameter {
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

	if t1.Kind() == types.KindTypeParameter {
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

	if t1.Kind() == types.KindDynamic || t1.Kind() == types.KindError {
		return true
	}

	if t2.Kind() == types.KindDynamic || t2.Kind() == types.KindError {
		return true
	}

	if t1.Kind() == types.KindNull && isNullable(t2) {
		return true
	}
	// Unwrap box types
	if t1.Kind() == types.KindWrapper {
		return internalIsAssignable(m, t1.(*types.WrapperType).Primitive(), t2)
	}
	// Finally check equality and type args recursively.
	if t1.Kind() != t2.Kind() {
		return false
	}

	switch t1.Kind() {
	case types.KindPrimitive:
		return t1.Equals(t2)
	case types.KindWellKnown:
		return t1.Equals(t2)
	case types.KindMessage:
		return t1.Equals(t2)
	case types.KindType:
		return internalIsAssignable(m, t1.(*types.TypeType).Target(), t2.(*types.TypeType).Target())
	case types.KindList:
		return internalIsAssignable(m, t1.(*types.ListType).ElementType, t2.(*types.ListType).ElementType)
	case types.KindMap:
		m1 := t1.(*types.MapType)
		m2 := t2.(*types.MapType)
		return internalIsAssignableList(m, []types.Type{m1.KeyType, m1.ValueType}, []types.Type{m2.KeyType, m2.ValueType})
	case types.KindFunction:
		fn1 := t1.(*types.FunctionType)
		fn2 := t2.(*types.FunctionType)
		return internalIsAssignableList(m, append(fn1.ArgTypes(), fn1.ResultType()), append(fn2.ArgTypes(), fn2.ResultType()))
	default:
		return false
	}
}

/**
 * Check whether one type is equal or less specific than the other one. A type is less specific if
 * it matches the other type using the DYN type.
 */
func isEqualOrLessSpecific(t1 types.Type, t2 types.Type) bool {
	if t1.Kind() == types.KindDynamic || t1.Kind() == types.KindTypeParameter {
		return true
	}

	if t2.Kind() == types.KindDynamic || t2.Kind() == types.KindTypeParameter {
		return false
	}

	if t1.Kind() != t2.Kind() {
		return false
	}

	switch t1.Kind() {
	case types.KindMessage, types.KindPrimitive, types.KindWellKnown, types.KindWrapper:
		return t1.Equals(t2)
	case types.KindType:
		return isEqualOrLessSpecific(t1.(*types.TypeType).Target(), t2.(*types.TypeType).Target())
	case types.KindList:
		return isEqualOrLessSpecific(t1.(*types.ListType).ElementType, t2.(*types.ListType).ElementType)
	case types.KindMap:
		m1 := t1.(*types.MapType)
		m2 := t2.(*types.MapType)
		return isEqualOrLessSpecific(m1.KeyType, m2.KeyType) && isEqualOrLessSpecific(m1.KeyType, m2.KeyType)
	case types.KindFunction:
		fn1 := t1.(*types.FunctionType)
		fn2 := t2.(*types.FunctionType)
		if len(fn1.ArgTypes()) != len(fn2.ArgTypes()) {
			return false
		}
		if !isEqualOrLessSpecific(fn1.ResultType(), fn2.ResultType()) {
			return false
		}
		for i, a1 := range fn1.ArgTypes() {
			if !isEqualOrLessSpecific(a1, fn2.ArgTypes()[i]) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

/** Returns the more general of two types which are known to unify. */
func mostGeneral(t1 types.Type, t2 types.Type) types.Type {
	if isEqualOrLessSpecific(t1, t2) {
		return t1
	}
	return t2
}

/**
 * Apply substitution to given type, replacing all direct and indirect occurrences of bound type
 * parameters. Unbound type parameters are replaced by DYN if typeParamToDyn is true.
 */
func substitute(m *Mapping, t types.Type, typeParamToDyn bool) types.Type {
	if tSub, found := m.Find(t); found {
		return substitute(m, tSub, typeParamToDyn)
	}
	if typeParamToDyn && t.Kind() == types.KindTypeParameter {
		return types.Dynamic
	}
	switch t.Kind() {
	case types.KindType:
		return types.NewTypeType(substitute(m, t.(*types.TypeType).Target(), typeParamToDyn))
	case types.KindList:
		return types.NewList(substitute(m, t.(*types.ListType).ElementType, typeParamToDyn))
	case types.KindMap:
		mt := t.(*types.MapType)
		return types.NewMap(substitute(m, mt.KeyType, typeParamToDyn), substitute(m, mt.ValueType, typeParamToDyn))
	case types.KindFunction:
		fn := t.(*types.FunctionType)
		rt := substitute(m, fn.ResultType(), typeParamToDyn)
		args := make([]types.Type, len(fn.ArgTypes()))
		for i, a := range fn.ArgTypes() {
			args[i] = substitute(m, a, typeParamToDyn)
		}
		return types.NewFunctionType(rt, args)
	default:
		return t
	}
}

func isNullable(t types.Type) bool {
	switch t.Kind() {
	case types.KindMessage, types.KindWellKnown, types.KindWrapper:
		return true
	default:
		return false
	}
}

func notReferencedIn(t types.Type, withinType types.Type) bool {
	if t.Equals(withinType) {
		return false
	}
	switch withinType.Kind() {
	case types.KindWrapper:
		return notReferencedIn(t, withinType.(*types.WrapperType).Primitive())
	case types.KindType:
		return notReferencedIn(t, withinType.(*types.TypeType).Target())
	case types.KindList:
		return notReferencedIn(t, withinType.(*types.ListType).ElementType)
	case types.KindMap:
		m := withinType.(*types.MapType)
		return notReferencedIn(t, m.KeyType) && notReferencedIn(t, m.ValueType)
	case types.KindFunction:
		fn := withinType.(*types.FunctionType)
		types := append(fn.ArgTypes(), fn.ResultType())
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

func typeKey(t types.Type) string {
	return fmt.Sprintf("%v:%v", t.Kind(), t.String())
}
