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

package ext

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// Lists returns a cel.EnvOption to configure extended functions for list manipulation.
// As a general note, all indices are zero-based.
// # Slice
//
// Returns a new sub-list using the indexes provided.
//
//	<list>.slice(<int>, <int>) -> <list>
//
// Examples:
//
//	[1,2,3,4].slice(1, 3) // return [2, 3]
//	[1,2,3,4].slice(2, 4) // return [3 ,4]
//
// # Flatten
//
// Flattens the given list to a single level.
//
//	<list>.flatten(<list>) -> <list>
//
// Examples:
//
// [1,[2,3],[4]].flatten() // return [1, 2, 3, 4]
// [1,[2,[3,4]]].flatten() // return [1, 2, [3, 4]]
// [1,2,[],[],[3,4]].flatten() // return [1, 2, 3, 4]
//
// # FlattenDeep
//
// Flattens the given list to the deepest level.
//
//	<list>.flattenDeep(<list>) -> <list>
//
// Examples:
//
// [1,[2,3],[4]].flattenDeep() // return [1, 2, 3, 4]
// [1,[2,[3,4]]].flattenDeep() // return [1, 2, 3, 4]
// [1,[2,[3,[4]]]].flattenDeep() // return [1, 2, 3, 4]
func Lists() cel.EnvOption {
	return cel.Lib(listsLib{})
}

type listsLib struct{}

// LibraryName implements the SingletonLibrary interface method.
func (listsLib) LibraryName() string {
	return "cel.lib.ext.lists"
}

// CompileOptions implements the Library interface method.
func (listsLib) CompileOptions() []cel.EnvOption {
	listType := cel.ListType(cel.TypeParamType("T"))
	listDyn := cel.ListType(cel.DynType)
	return []cel.EnvOption{
		cel.Function("slice",
			cel.MemberOverload("list_slice",
				[]*cel.Type{listType, cel.IntType, cel.IntType}, listType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					list := args[0].(traits.Lister)
					start := args[1].(types.Int)
					end := args[2].(types.Int)
					result, err := slice(list, start, end)
					if err != nil {
						return types.WrapErr(err)
					}
					return result
				}),
			),
		),
		cel.Function("flatten",
			cel.MemberOverload("list_flatten",
				[]*cel.Type{listDyn}, listDyn,
				cel.UnaryBinding(func(arg ref.Val) ref.Val {
					list := arg.(traits.Lister)
					flatList := flattenHelper(list, false)
					return types.DefaultTypeAdapter.NativeToValue(flatList)
				}),
			),
		),
		cel.Function("flattenDeep",
			cel.MemberOverload("list_flatten_deep",
				[]*cel.Type{listDyn}, listDyn,
				cel.UnaryBinding(func(arg ref.Val) ref.Val {
					list := arg.(traits.Lister)
					flatList := flattenHelper(list, true)
					return types.DefaultTypeAdapter.NativeToValue(flatList)
				}),
			),
		),
	}
}

// ProgramOptions implements the Library interface method.
func (listsLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func slice(list traits.Lister, start, end types.Int) (ref.Val, error) {
	listLength := list.Size().(types.Int)
	if start < 0 || end < 0 {
		return nil, fmt.Errorf("cannot slice(%d, %d), negative indexes not supported", start, end)
	}
	if start > end {
		return nil, fmt.Errorf("cannot slice(%d, %d), start index must be less than or equal to end index", start, end)
	}
	if listLength < end {
		return nil, fmt.Errorf("cannot slice(%d, %d), list is length %d", start, end, listLength)
	}

	var newList []ref.Val
	for i := types.Int(start); i < end; i++ {
		val := list.Get(i)
		newList = append(newList, val)
	}
	return types.DefaultTypeAdapter.NativeToValue(newList), nil
}

func flattenHelper(list traits.Lister, deep bool) []ref.Val {
	iter := list.Iterator()
	var newList []ref.Val

	for iter.HasNext() == types.True {
		val := iter.Next()

		if nestedList, ok := val.(traits.Lister); ok {
			if deep {
				newList = append(newList, flattenHelper(nestedList, true)...)
			} else {
				nestedIter := nestedList.Iterator()
				for nestedIter.HasNext() == types.True {
					newList = append(newList, nestedIter.Next())
				}
			}
		} else {
			newList = append(newList, val)
		}
	}

	return newList
}
