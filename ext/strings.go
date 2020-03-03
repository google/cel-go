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
// After
//
// Given the target string and a search string, return the substring that occurs after the search
// string, or empty string if no match is found:
//
//     <string>.after(<string>) -> <string>
//
// Examples:
//
//     'dragon'.after('taco')  // returns ''
//     'tacocat'.after('taco') // returns 'cat'
//
// Before
//
// Given the target string and a search string, return the substring that occurs before the search
// string, or the target string if no match is found:
//
//     <string>.before(<string>) -> <string>
//
// Examples:
//
//     'dragon'.before('taco') // returns 'dragon'
//     'tacocat'.before('cat') // returns 'taco'
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
// Returns the integer index of the search string if found, or -1 if not found. Accepts an optional
// index from which to start the search:
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
// Lower
//
// Return a new lowercase version of the target string:
//
//     <string>.lower() -> <string>
//
// Examples:
//
//     'User-Agent'.lower() // returns 'user-agent'
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
//     'hello hello'.replace('he', 'we')    // returns 'wello wello'
//     'hello hello'.replace('he', 'we', 1) // returns 'wello hello'
//     'hello hello'.replace('he', 'we', 0) // returns 'hello wello'
//     'hello hello'.replace('he', 'we', -1) // returns 'wello wello'
//
// Split
//
// Produces a list of strings which were split from the input by the given seperator. Accepts an
// optional argument specifying a limit on the number of substrings produced by the split:
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
// a string:
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
// Returns a new string which trims the whitespace before and after the target string:
//
//      <string>.trim() -> <string>
//
// Examples:
//
//     '    trim    '.trim() // returns 'trim'
//
// Upper
//
// Returns a new uppercase version of the target string:
//
//     <string>.upper() -> <string>
//
// Examples:
//
//     'Constant'.upper() // returns 'CONSTANT'
func Strings() cel.EnvOption {
	return cel.Lib(stringLib{})
}

type stringLib struct{}

func (stringLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Declarations(
			decls.NewFunction("after",
				decls.NewInstanceOverload("string_after_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.String)),
			decls.NewFunction("before",
				decls.NewInstanceOverload("string_before_string",
					[]*exprpb.Type{decls.String, decls.String},
					decls.String)),
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
			decls.NewFunction("lower",
				decls.NewInstanceOverload("string_lower",
					[]*exprpb.Type{decls.String},
					decls.String)),
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
			decls.NewFunction("upper",
				decls.NewInstanceOverload("string_upper",
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
				Operator: "after",
				Binary:   callInStrStrOutStr(after),
			},
			&functions.Overload{
				Operator: "string_after_stirng",
				Binary:   callInStrStrOutStr(after),
			},
			&functions.Overload{
				Operator: "before",
				Binary:   callInStrStrOutStr(before),
			},
			&functions.Overload{
				Operator: "string_before_string",
				Binary:   callInStrStrOutStr(before),
			},
			&functions.Overload{
				Operator: "charAt",
				Binary:   callInStrIntOutStr(charAt),
			},
			&functions.Overload{
				Operator: "string_chat_at_int",
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
				Operator: "lower",
				Unary:    callInStrOutStr(strings.ToLower),
			},
			&functions.Overload{
				Operator: "string_lower",
				Unary:    callInStrOutStr(strings.ToLower),
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
			&functions.Overload{
				Operator: "upper",
				Unary:    callInStrOutStr(strings.ToUpper),
			},
			&functions.Overload{
				Operator: "string_upper",
				Unary:    callInStrOutStr(strings.ToUpper),
			},
		),
	}
}

func after(str, sub string) (string, error) {
	ind := strings.Index(str, sub)
	if ind < 0 {
		return "", nil
	}
	return str[ind+len(sub):], nil
}

func before(str, sub string) (string, error) {
	ind := strings.Index(str, sub)
	if ind < 0 {
		return str, nil
	}
	return str[:ind], nil
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

func callInStrStrOutStr(fn func(string, string) (string, error)) functions.BinaryOp {
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
		return types.String(out)
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
