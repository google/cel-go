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
	"math"
	"sort"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

var comparableTypes = []*cel.Type{
	cel.IntType,
	cel.UintType,
	cel.DoubleType,
	cel.BoolType,
	cel.DurationType,
	cel.TimestampType,
	cel.StringType,
	cel.BytesType,
}

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
// Flattens a list recursively.
// If an optional depth is provided, the list is flattened to a the specificied level.
// A negative depth value will result in an error.
//
//	<list>.flatten(<list>) -> <list>
//	<list>.flatten(<list>, <int>) -> <list>
//
// Examples:
//
// [1,[2,3],[4]].flatten() // return [1, 2, 3, 4]
// [1,[2,[3,4]]].flatten() // return [1, 2, [3, 4]]
// [1,2,[],[],[3,4]].flatten() // return [1, 2, 3, 4]
// [1,[2,[3,[4]]]].flatten(2) // return [1, 2, 3, [4]]
// [1,[2,[3,[4]]]].flatten(-1) // error
//
// # Sort
//
// Introduced in version: 2
//
// Sorts a list with comparable elements. If the element type is not comparable
// or the element types are not the same, the function will produce an error.
//
//	<list(T)>.sort() -> <list(T)>
//	T in {int, uint, double, bool, duration, timestamp, string, bytes}
//
// Examples:
//
//	[3, 2, 1].sort() // return [1, 2, 3]
//	["b", "c", "a"].sort() // return ["a", "b", "c"]
//	[1, "b"].sort() // error
//	[[1, 2, 3]].sort() // error

func Lists(options ...ListsOption) cel.EnvOption {
	l := &listsLib{
		version: math.MaxUint32,
	}
	for _, o := range options {
		l = o(l)
	}

	return cel.Lib(l)
}

type listsLib struct {
	version uint32
}

// LibraryName implements the SingletonLibrary interface method.
func (listsLib) LibraryName() string {
	return "cel.lib.ext.lists"
}

// ListsOption is a functional interface for configuring the strings library.
type ListsOption func(*listsLib) *listsLib

// ListsVersion configures the version of the string library.
//
// The version limits which functions are available. Only functions introduced
// below or equal to the given version included in the library. If this option
// is not set, all functions are available.
//
// See the library documentation to determine which version a function was introduced.
// If the documentation does not state which version a function was introduced, it can
// be assumed to be introduced at version 0, when the library was first created.
func ListsVersion(version uint32) ListsOption {
	return func(lib *listsLib) *listsLib {
		lib.version = version
		return lib
	}
}

// CompileOptions implements the Library interface method.
func (lib listsLib) CompileOptions() []cel.EnvOption {
	listType := cel.ListType(cel.TypeParamType("T"))
	listListType := cel.ListType(listType)
	listDyn := cel.ListType(cel.DynType)
	opts := []cel.EnvOption{
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
	}
	if lib.version >= 1 {
		opts = append(opts,
			cel.Function("flatten",
				cel.MemberOverload("list_flatten",
					[]*cel.Type{listListType}, listType,
					cel.UnaryBinding(func(arg ref.Val) ref.Val {
						list, ok := arg.(traits.Lister)
						if !ok {
							return types.MaybeNoSuchOverloadErr(arg)
						}
						flatList, err := flatten(list, 1)
						if err != nil {
							return types.WrapErr(err)
						}

						return types.DefaultTypeAdapter.NativeToValue(flatList)
					}),
				),
				cel.MemberOverload("list_flatten_int",
					[]*cel.Type{listDyn, types.IntType}, listDyn,
					cel.BinaryBinding(func(arg1, arg2 ref.Val) ref.Val {
						list, ok := arg1.(traits.Lister)
						if !ok {
							return types.MaybeNoSuchOverloadErr(arg1)
						}
						depth, ok := arg2.(types.Int)
						if !ok {
							return types.MaybeNoSuchOverloadErr(arg2)
						}
						flatList, err := flatten(list, int64(depth))
						if err != nil {
							return types.WrapErr(err)
						}

						return types.DefaultTypeAdapter.NativeToValue(flatList)
					}),
				),
				// To handle the case where a variable of just `list(T)` is provided at runtime
				// with a graceful failure more, disable the type guards since the implementation
				// can handle lists which are already flat.
				decls.DisableTypeGuards(true),
			),
		)
	}
	if lib.version >= 2 {
		sortDecl := cel.Function("sort",
			append(
				templatedOverloads(comparableTypes, func(t *cel.Type) cel.FunctionOpt {
					return cel.MemberOverload(
						fmt.Sprintf("list_%s_sort", t.TypeName()),
						[]*cel.Type{cel.ListType(t)}, cel.ListType(t),
					)
				}),
				cel.SingletonUnaryBinding(
					func(arg ref.Val) ref.Val {
						list, ok := arg.(traits.Lister)
						if !ok {
							return types.MaybeNoSuchOverloadErr(arg)
						}
						sorted, err := sortList(list)
						if err != nil {
							return types.WrapErr(err)
						}

						return sorted
					},
					// List traits
					traits.ListerType,
				),
			)...,
		)
		opts = append(opts, sortDecl)
	}

	return opts
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

func flatten(list traits.Lister, depth int64) ([]ref.Val, error) {
	if depth < 0 {
		return nil, fmt.Errorf("level must be non-negative")
	}

	var newList []ref.Val
	iter := list.Iterator()

	for iter.HasNext() == types.True {
		val := iter.Next()
		nestedList, isList := val.(traits.Lister)

		if !isList || depth == 0 {
			newList = append(newList, val)
			continue
		} else {
			flattenedList, err := flatten(nestedList, depth-1)
			if err != nil {
				return nil, err
			}

			newList = append(newList, flattenedList...)
		}
	}

	return newList, nil
}

func sortList(list traits.Lister) (ref.Val, error) {
	listLength := list.Size().(types.Int)
	if listLength == 0 {
		return list, nil
	}
	elem := list.Get(types.IntZero)
	if _, ok := elem.(traits.Comparer); !ok {
		return nil, fmt.Errorf("list elements must be comparable")
	}

	sorted := make([]ref.Val, 0, listLength)
	for i := types.IntZero; i < listLength; i++ {
		val := list.Get(i)
		if val.Type() != elem.Type() {
			return nil, fmt.Errorf("list elements must have the same type")
		}
		sorted = append(sorted, val)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].(traits.Comparer).Compare(sorted[j]) == types.IntNegOne
	})

	return types.DefaultTypeAdapter.NativeToValue(sorted), nil
}

func templatedOverloads(types []*cel.Type, template func(t *cel.Type) cel.FunctionOpt) []cel.FunctionOpt {
	overloads := make([]cel.FunctionOpt, len(types))
	for i, t := range types {
		overloads[i] = template(t)
	}
	return overloads
}
