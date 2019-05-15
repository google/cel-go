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
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// InterpretableDecorator is a functional interface for decorating or replacing
// Interpretable expression nodes at construction time.
type InterpretableDecorator func(Interpretable) (Interpretable, error)

// evalObserver is a functional interface that accepts an expression id and an observed value.
type evalObserver func(int64, ref.Val)

// decObserveEval records evaluation state into an EvalState object.
func decObserveEval(observer evalObserver) InterpretableDecorator {
	return func(i Interpretable) (Interpretable, error) {
		return &evalWatch{
			inst:     i,
			observer: observer,
		}, nil
	}
}

// decDisableShortcircuits ensures that all branches of an expression will be evaluated, no short-circuiting.
func decDisableShortcircuits() InterpretableDecorator {
	return func(i Interpretable) (Interpretable, error) {
		switch i.(type) {
		case *evalOr:
			or := i.(*evalOr)
			return &evalExhaustiveOr{
				id:  or.id,
				lhs: or.lhs,
				rhs: or.rhs,
			}, nil
		case *evalAnd:
			and := i.(*evalAnd)
			return &evalExhaustiveAnd{
				id:  and.id,
				lhs: and.lhs,
				rhs: and.rhs,
			}, nil
		case *evalConditional:
			cond := i.(*evalConditional)
			return &evalExhaustiveConditional{
				id:     cond.id,
				expr:   cond.expr,
				truthy: cond.truthy,
				falsy:  cond.falsy,
			}, nil
		case *evalFold:
			fold := i.(*evalFold)
			return &evalExhaustiveFold{
				id:        fold.id,
				accu:      fold.accu,
				accuVar:   fold.accuVar,
				iterRange: fold.iterRange,
				iterVar:   fold.iterVar,
				cond:      fold.cond,
				step:      fold.step,
				result:    fold.result,
			}, nil
		}
		return i, nil
	}
}

// decOptimize optimizes the program plan by looking for common evaluation patterns and
// conditionally precomputating the result.
// - build list and map values with constant elements.
// - convert 'in' operations to set membership tests if possible.
func decOptimize() InterpretableDecorator {
	return func(i Interpretable) (Interpretable, error) {
		switch i.(type) {
		case *evalList:
			return maybeBuildListLiteral(i, i.(*evalList))
		case *evalMap:
			return maybeBuildMapLiteral(i, i.(*evalMap))
		case *evalBinary:
			call := i.(*evalBinary)
			if call.overload == overloads.InList {
				return maybeOptimizeSetMembership(i, call)
			}
		}
		return i, nil
	}
}

func maybeBuildListLiteral(i Interpretable, l *evalList) (Interpretable, error) {
	for _, elem := range l.elems {
		_, isConst := elem.(*evalConst)
		if !isConst {
			return i, nil
		}
	}
	val := l.Eval(EmptyActivation())
	return &evalConst{
		id:  l.id,
		val: val,
	}, nil
}

func maybeBuildMapLiteral(i Interpretable, mp *evalMap) (Interpretable, error) {
	for idx, key := range mp.keys {
		_, isConst := key.(*evalConst)
		if !isConst {
			return i, nil
		}
		_, isConst = mp.vals[idx].(*evalConst)
		if !isConst {
			return i, nil
		}
	}
	val := mp.Eval(EmptyActivation())
	return &evalConst{
		id:  mp.id,
		val: val,
	}, nil
}

// maybeOptimizeSetMembership may convert an 'in' operation against a list to map key membership
// test if the following conditions are true:
// - the list is a constant with homogeneous element types.
// - the elements are all of primitive type.
func maybeOptimizeSetMembership(i Interpretable, inlist *evalBinary) (Interpretable, error) {
	l, isConst := inlist.rhs.(*evalConst)
	if !isConst {
		return i, nil
	}
	// When the incoming binary call is flagged with as the InList overload, the value will
	// always be convertible to a `traits.Lister` type.
	list := l.val.(traits.Lister)
	if list.Size() == types.IntZero {
		return &evalConst{
			id:  inlist.id,
			val: types.False,
		}, nil
	}
	it := list.Iterator()
	var typ ref.Type
	valueSet := make(map[ref.Val]ref.Val)
	for it.HasNext() == types.True {
		elem := it.Next()
		if !types.IsPrimitiveType(elem) {
			// Note, non-primitive type are not yet supported.
			return i, nil
		}
		if typ == nil {
			typ = elem.Type()
		} else if typ.TypeName() != elem.Type().TypeName() {
			return i, nil
		}
		valueSet[elem] = types.True
	}
	return &evalSetMembership{
		inst:        inlist,
		arg:         inlist.lhs,
		argTypeName: typ.TypeName(),
		valueSet:    valueSet,
	}, nil
}
