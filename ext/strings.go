// Copyright 2020 Google LLC
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

// Package ext contains CEL extension libraries where each library defines a related set of
// constants, functions, macros, or other configuration settings which may not be covered by
// the core CEL spec.
package ext

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter/functions"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Strings returns a cel.EnvOption to configure extended functions for string manipulation.
//
// CharAt
//
// Returns the character at the given position. If the position is less than 0 or greater than the
// length of the string, the function will produce an error. CharAt is an instance method with the
// following signature:
//
//     <string>.charAt(<int>) -> <string>
//
// Examples:
//
//     'hello'.charAt(4)  // return 'o'
//     'hello'.charAt(-1) // error
//     'hello'.charAt(5)  // error
//
// IndexOf
//
// Returns the integer index of the first occurrence of the search string. If the search string is
// not found the function returns -1. The function also accepts an optional index from which to
// begin the substring search:
//
//     <string>.indexOf(<string>) -> <int>
//     <string>.indexOf(<string>, <int>) -> <int>
//
// Examples:
//
//     'hello mellow'.indexOf('ello')     // returns 1
//     'hello mellow'.indexOf('jello')    // returns -1
//     'hello mellow'.indexOf('ello', 2)  // returns 7
//     'hello mellow'.indexOf('ello', 20) // error
//
// LastIndexOf
//
// Returns the integer index of the last occurrence of the search string. If the search string is
// not found the function returns -1. The function also accepts an optional index which represents
// where to end the search.
//
//     <string>.lastIndexOf(<string>) -> <int>
//     <string>.lastIndexOf(<string>, <int>) -> <int>
//
// Examples:
//
//     'hello mellow'.lastIndexOf('ello')     // returns 7
//     'hello mellow'.lastIndexOf('jello')    // returns -1
//     'hello mellow'.lastIndexOf('ello', 8)  // returns 2
//     'hello mellow'.lastIndexOf('ello', -1) // error
//
// Replace
//
// Produces a new string based on the target, which replaces the occurrences of a search string
// with a replacement string if present. Accepts an optional argument specifying a limit on the
// number of substring replacements to be made:
//
//     <string>.replace(<string>, <string>) -> <string>
//     <string>.replace(<string>, <string>, <int>) -> <string>
//
// Examples:
//
//     'hello hello'.replace('he', 'we')     // returns 'wello wello'
//     'hello hello'.replace('he', 'we', -1) // returns 'wello wello'
//     'hello hello'.replace('he', 'we', 1)  // returns 'wello hello'
//     'hello hello'.replace('he', 'we', 0)  // returns 'hello hello'
//
// Split
//
// Produces a list of strings which were split from the input by the given seperator. The function
// accepts an optional argument specifying a limit on the number of substrings produced by the
// split.
//
// When the split limit is 0, the result is an empty list. When the limit is 1, the result is the
// target string to split. When the limit is -1, the function behaves the same as split all.
//
//     <string>.split(<string>) -> <list<string>>
//     <string>.split(<string>, <int>) -> <list<string>>
//
// Examples:
//
//     'hello hello hello'.split(' ')     // returns ['hello', 'hello', 'hello']
//     'hello hello hello'.split(' ', 0)  // returns []
//     'hello hello hello'.split(' ', 1)  // returns ['hello hello hello']
//     'hello hello hello'.split(' ', 2)  // returns ['hello', 'hello hello']
//     'hello hello hello'.split(' ', -1) // returns ['hello', 'hello', 'hello']
//
// Substring
//
// Returns the substring given a numeric range corresponding to character positions. Optionally
// may omit the trailing range for a substring from a given character position until the end of
// a string.
//
// Character offsets are 0-based with an inclusive start range and exclusive end range. It is an
// error to specify an end range that is lower than the start range, or for either the start or end
// index to be negative or exceed the string length.
//
//     <string>.substring(<int>) -> <string>
//     <string>.substring(<int>, <int>) -> <string>
//
// Examples:
//
//     'tacocat'.substring(4)    // returns 'cat'
//     'tacocat'.substring(0, 4) // returns 'taco'
//     'tacocat'.substring(-1)   // error
//     'tacocat'.substring(2, 1) // error
//
// Trim
//
// Returns a new string which removes the leading and trailing whitespace in the target string.
// The trim function uses the Unicode definition of whitespace which does not include the
// zero-width spaces. See: https://en.wikipedia.org/wiki/Whitespace_character#Unicode
//
//      <string>.trim() -> <string>
//
// Examples:
//
//     '  \ttrim\n    '.trim() // returns 'trim'
func Strings() cel.EnvOption {
	return cel.Lib(stringLib{})
}

type stringLib struct{}

func (stringLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Declarations(
			decls.NewFunction("charAt",
				decls.NewInstanceOverload("string_char_at_int",
					[]*exprpb.Type{decls.String, decls.Int},
					decls.String)),
			decls.NewFunction("indexOf",
				decls.NewInstanceOverload("string_index_of_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Int),
				decls.NewInstanceOverload("string_index_of_string_int",
					[]*exprpb.Type{decls.String, decls.String, decls.Int},
					decls.Int)),
			decls.NewFunction("lastIndexOf",
				decls.NewInstanceOverload("string_last_index_of_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.Int),
				decls.NewInstanceOverload("string_last_index_of_string_int",
					[]*exprpb.Type{decls.String, decls.String, decls.Int},
					decls.Int)),
			decls.NewFunction("replace",
				decls.NewInstanceOverload("string_replace_string_string",
					[]*exprpb.Type{decls.String, decls.String, decls.String},
					decls.String),
				decls.NewInstanceOverload("string_replace_string_string_int",
					[]*exprpb.Type{decls.String, decls.String, decls.String, decls.Int},
					decls.String)),
			decls.NewFunction("split",
				decls.NewInstanceOverload("string_split_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.NewListType(decls.String)),
				decls.NewInstanceOverload("string_split_string_int",
					[]*exprpb.Type{decls.String, decls.String, decls.Int},
					decls.NewListType(decls.String))),
			decls.NewFunction("substring",
				decls.NewInstanceOverload("string_substring_int",
					[]*exprpb.Type{decls.String, decls.Int},
					decls.String),
				decls.NewInstanceOverload("string_substring_int_int",
					[]*exprpb.Type{decls.String, decls.Int, decls.Int},
					decls.String)),
			decls.NewFunction("trim",
				decls.NewInstanceOverload("string_trim",
					[]*exprpb.Type{decls.String},
					decls.String)),
		),
	}
}

func (stringLib) ProgramOptions() []cel.ProgramOption {
	wrappedReplace := callInStrStrStrOutStr(replace)
	wrappedReplaceN := callInStrStrStrIntOutStr(replaceN)
	return []cel.ProgramOption{
		cel.Functions(
			&functions.Overload{
				Operator: "charAt",
				Binary:   callInStrIntOutStr(charAt),
			},
			&functions.Overload{
				Operator: "string_char_at_int",
				Binary:   callInStrIntOutStr(charAt),
			},
			&functions.Overload{
				Operator: "indexOf",
				Binary:   callInStrStrOutInt(indexOf),
				Function: callInStrStrIntOutInt(indexOfOffset),
			},
			&functions.Overload{
				Operator: "string_index_of_string",
				Binary:   callInStrStrOutInt(indexOf),
			},
			&functions.Overload{
				Operator: "string_index_of_string_int",
				Function: callInStrStrIntOutInt(indexOfOffset),
			},
			&functions.Overload{
				Operator: "lastIndexOf",
				Binary:   callInStrStrOutInt(lastIndexOf),
				Function: callInStrStrIntOutInt(lastIndexOfOffset),
			},
			&functions.Overload{
				Operator: "string_last_index_of_string",
				Binary:   callInStrStrOutInt(lastIndexOf),
			},
			&functions.Overload{
				Operator: "string_last_index_of_string_int",
				Function: callInStrStrIntOutInt(lastIndexOfOffset),
			},
			&functions.Overload{
				Operator: "replace",
				Function: func(values ...ref.Val) ref.Val {
					if len(values) == 3 {
						return wrappedReplace(values...)
					}
					if len(values) == 4 {
						return wrappedReplaceN(values...)
					}
					return types.NewErr("no such overload")
				},
			},
			&functions.Overload{
				Operator: "string_replace_string_string",
				Function: wrappedReplace,
			},
			&functions.Overload{
				Operator: "string_replace_string_string_int",
				Function: wrappedReplaceN,
			},
			&functions.Overload{
				Operator: "split",
				Binary:   callInStrStrOutListStr(split),
				Function: callInStrStrIntOutListStr(splitN),
			},
			&functions.Overload{
				Operator: "string_split_string",
				Binary:   callInStrStrOutListStr(split),
			},
			&functions.Overload{
				Operator: "string_split_string_int",
				Function: callInStrStrIntOutListStr(splitN),
			},
			&functions.Overload{
				Operator: "substring",
				Binary:   callInStrIntOutStr(substr),
				Function: callInStrIntIntOutStr(substrRange),
			},
			&functions.Overload{
				Operator: "string_substring_int",
				Binary:   callInStrIntOutStr(substr),
			},
			&functions.Overload{
				Operator: "string_substring_int_int",
				Function: callInStrIntIntOutStr(substrRange),
			},
			&functions.Overload{
				Operator: "trim",
				Unary:    callInStrOutStr(strings.TrimSpace),
			},
			&functions.Overload{
				Operator: "string_trim",
				Unary:    callInStrOutStr(strings.TrimSpace),
			},
		),
	}
}

func charAt(str string, ind int64) (string, error) {
	i := int(ind)
	if i < 0 || i >= len(str) {
		return "", fmt.Errorf("index out of range: %d", ind)
	}
	return str[i : i+1], nil
}

func indexOf(str, substr string) (int64, error) {
	return int64(strings.Index(str, substr)), nil
}

func indexOfOffset(str, substr string, offset int64) (int64, error) {
	off := int(offset)
	if off < 0 || off >= len(str) {
		return -1, fmt.Errorf("index out of range: %d", off)
	}
	return offset + int64(strings.Index(str[offset:], substr)), nil
}

func lastIndexOf(str, substr string) (int64, error) {
	return int64(strings.LastIndex(str, substr)), nil
}

func lastIndexOfOffset(str, substr string, offset int64) (int64, error) {
	off := int(offset)
	if off < 0 || off >= len(str) {
		return -1, fmt.Errorf("index out of range: %d", off)
	}
	return int64(strings.Index(str[:off], substr)), nil
}

func replace(str, old, new string) (string, error) {
	return strings.ReplaceAll(str, old, new), nil
}

func replaceN(str, old, new string, n int64) (string, error) {
	return strings.Replace(str, old, new, int(n)), nil
}

func split(str, sep string) ([]string, error) {
	return strings.Split(str, sep), nil
}

func splitN(str, sep string, n int64) ([]string, error) {
	return strings.SplitN(str, sep, int(n)), nil
}

func substr(str string, start int64) (string, error) {
	if int(start) < 0 || int(start) >= len(str) {
		return "", fmt.Errorf("index out of range: %d", start)
	}
	return str[int(start):], nil
}

func substrRange(str string, start, end int64) (string, error) {
	l := len(str)
	if start > end {
		return "", fmt.Errorf("invalid substring range. start: %d, end: %d", start, end)
	}
	if int(start) < 0 || int(start) >= l {
		return "", fmt.Errorf("index out of range: %d", start)
	}
	if int(end) < 0 || int(end) >= l {
		return "", fmt.Errorf("index out of range: %d", end)
	}
	return str[int(start):int(end)], nil
}

func callInStrOutStr(fn func(string) string) functions.UnaryOp {
	return func(val ref.Val) ref.Val {
		vVal, ok := val.(types.String)
		if !ok {
			return types.ValOrErr(val, "no such overload")
		}
		return types.String(fn(string(vVal)))
	}
}

func callInStrIntOutStr(fn func(string, int64) (string, error)) functions.BinaryOp {
	return func(val, arg ref.Val) ref.Val {
		vVal, ok := val.(types.String)
		if !ok {
			return types.ValOrErr(val, "no such overload")
		}
		argVal, ok := arg.(types.Int)
		if !ok {
			return types.ValOrErr(arg, "no such overload")
		}
		out, err := fn(string(vVal), int64(argVal))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(out)
	}
}

func callInStrStrOutInt(fn func(string, string) (int64, error)) functions.BinaryOp {
	return func(val, arg ref.Val) ref.Val {
		vVal, ok := val.(types.String)
		if !ok {
			return types.ValOrErr(val, "no such overload")
		}
		argVal, ok := arg.(types.String)
		if !ok {
			return types.ValOrErr(arg, "no such overload")
		}
		out, err := fn(string(vVal), string(argVal))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.Int(out)
	}
}

func callInStrStrOutListStr(fn func(string, string) ([]string, error)) functions.BinaryOp {
	return func(val, arg ref.Val) ref.Val {
		vVal, ok := val.(types.String)
		if !ok {
			return types.ValOrErr(val, "no such overload")
		}
		argVal, ok := arg.(types.String)
		if !ok {
			return types.ValOrErr(arg, "no such overload")
		}
		out, err := fn(string(vVal), string(argVal))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.DefaultTypeAdapter.NativeToValue(out)
	}
}

func callInStrIntIntOutStr(fn func(string, int64, int64) (string, error)) functions.FunctionOp {
	return func(args ...ref.Val) ref.Val {
		if len(args) != 3 {
			return types.NewErr("no such overload")
		}
		vVal, ok := args[0].(types.String)
		if !ok {
			return types.ValOrErr(args[0], "no such overload")
		}
		arg1Val, ok := args[1].(types.Int)
		if !ok {
			return types.ValOrErr(args[1], "no such overload")
		}
		arg2Val, ok := args[2].(types.Int)
		if !ok {
			return types.ValOrErr(args[2], "no such overload")
		}
		out, err := fn(string(vVal), int64(arg1Val), int64(arg2Val))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(out)
	}
}

func callInStrStrStrOutStr(fn func(string, string, string) (string, error)) functions.FunctionOp {
	return func(args ...ref.Val) ref.Val {
		if len(args) != 3 {
			return types.NewErr("no such overload")
		}
		vVal, ok := args[0].(types.String)
		if !ok {
			return types.ValOrErr(args[0], "no such overload")
		}
		arg1Val, ok := args[1].(types.String)
		if !ok {
			return types.ValOrErr(args[1], "no such overload")
		}
		arg2Val, ok := args[2].(types.String)
		if !ok {
			return types.ValOrErr(args[2], "no such overload")
		}
		out, err := fn(string(vVal), string(arg1Val), string(arg2Val))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(out)
	}
}

func callInStrStrIntOutInt(fn func(string, string, int64) (int64, error)) functions.FunctionOp {
	return func(args ...ref.Val) ref.Val {
		if len(args) != 3 {
			return types.NewErr("no such overload")
		}
		vVal, ok := args[0].(types.String)
		if !ok {
			return types.ValOrErr(args[0], "no such overload")
		}
		arg1Val, ok := args[1].(types.String)
		if !ok {
			return types.ValOrErr(args[1], "no such overload")
		}
		arg2Val, ok := args[2].(types.Int)
		if !ok {
			return types.ValOrErr(args[2], "no such overload")
		}
		out, err := fn(string(vVal), string(arg1Val), int64(arg2Val))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.Int(out)
	}
}

func callInStrStrIntOutListStr(fn func(string, string, int64) ([]string, error)) functions.FunctionOp {
	return func(args ...ref.Val) ref.Val {
		if len(args) != 3 {
			return types.NewErr("no such overload")
		}
		vVal, ok := args[0].(types.String)
		if !ok {
			return types.ValOrErr(args[0], "no such overload")
		}
		arg1Val, ok := args[1].(types.String)
		if !ok {
			return types.ValOrErr(args[1], "no such overload")
		}
		arg2Val, ok := args[2].(types.Int)
		if !ok {
			return types.ValOrErr(args[2], "no such overload")
		}
		out, err := fn(string(vVal), string(arg1Val), int64(arg2Val))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.DefaultTypeAdapter.NativeToValue(out)
	}
}

func callInStrStrStrIntOutStr(fn func(string, string, string, int64) (string, error)) functions.FunctionOp {
	return func(args ...ref.Val) ref.Val {
		if len(args) != 4 {
			return types.NewErr("no such overload")
		}
		vVal, ok := args[0].(types.String)
		if !ok {
			return types.ValOrErr(args[0], "no such overload")
		}
		arg1Val, ok := args[1].(types.String)
		if !ok {
			return types.ValOrErr(args[1], "no such overload")
		}
		arg2Val, ok := args[2].(types.String)
		if !ok {
			return types.ValOrErr(args[2], "no such overload")
		}
		arg3Val, ok := args[3].(types.Int)
		if !ok {
			return types.ValOrErr(args[3], "no such overload")
		}
		out, err := fn(string(vVal), string(arg1Val), string(arg2Val), int64(arg3Val))
		if err != nil {
			return types.NewErr(err.Error())
		}
		return types.String(out)
	}
}
