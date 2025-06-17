// Copyright 2025 Google LLC
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
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

const (
	regexReplace    = "regex.replace"
	regexExtract    = "regex.extract"
	regexExtractAll = "regex.extractAll"
)

// Regex introduces support for regular expressions in CEL.
//
// This library provides functions for capturing groups, replacing strings using regex patterns,
// Regex configures namespaced regex helper functions.
// Note, all functions use the 'regex' namespace. If you are
// currently using a variable named 'regex', the macro will likely work just as
// intended; however, there is some chance for collision.
//
// # Replace
//
// The `regex.replace` function replaces all occurrences of a regex pattern in a
// string with a replacement string. Optionally, you can limit the number of
// replacements by providing a count argument. Both numeric ($N) and named
// (${name}) capture group references are supported in the replacement string, with
// validation for correctness. An error will be thrown for invalid regex or replace
// string.
//
//  regex.replace(target: string, pattern: string, replacement: string) -> string
//  regex.replace(target: string, pattern: string, replacement: string, count: int) -> string
//
// Examples:
//
//  regex.replace('banana', 'a', 'x', 0) == 'banana'
//  regex.replace('banana', 'a', 'x', 1) == 'bxnana'
//  regex.replace('banana', 'a', 'x', 2) == 'bxnxna'
//  regex.replace('foo bar', '(fo)o (ba)r', '$2 $1') == 'ba fo'
//
//  regex.replace('test', '(.)', '$2') // Runtime Error invalid replace string
//  regex.replace('foo bar', '(', '$2 $1') // Runtime Error invalid regex string
//  regex.replace('id=123', 'id=(?P<value>\\\\d+)', 'value: ${values}') // Runtime Error invalid replace string
//
// # Extract
//
// The `regex.extract` function returns the first match of a regex pattern in a
// string. If no match is found, it returns an optional none value. An error will
// be thrown for invalid regex or for multiple capture groups.
//
//  regex.extract(target: string, pattern: string) -> optional<string>
//
// Examples:
//
//  regex.extract('hello world', 'hello(.*)') == optional.of(' world')
//  regex.extract('item-A, item-B', 'item-(\\w+)') == optional.of('A')
//  regex.extract('HELLO', 'hello') == optional.empty()
//
//  regex.extract('testuser@testdomain', '(.*)@([^.]*)') // Runtime Error multiple capture group
//
// # Extract All
//
// The `regex.extractAll` function returns a list of all matches of a regex
// pattern in a target string. If no matches are found, it returns an empty list. An error will
// be thrown for invalid regex or for multiple capture groups.
//
//  regex.extractAll(target: string, pattern: string) -> list<string>
//
// Examples:
//
//  regex.extractAll('id:123, id:456', 'id:\\d+') == ['id:123', 'id:456']
//  regex.extractAll('id:123, id:456', 'assa') == []
//
//  regex.extractAll('testuser@testdomain', '(.*)@([^.]*)') // Runtime Error multiple capture group
//

func Regex(options ...RegexOptions) cel.EnvOption {
	s := &regexLib{
		version: math.MaxUint32,
	}
	for _, o := range options {
		s = o(s)
	}
	return cel.Lib(s)
}

type RegexOptions func(*regexLib) *regexLib

func RegexVersion(version uint32) RegexOptions {
	return func(lib *regexLib) *regexLib {
		lib.version = version
		return lib
	}
}

type regexLib struct {
	version uint32
}

func (r *regexLib) LibraryName() string {
	return "cel.lib.ext.regex"
}

// CompileOptions implements cel.Library.
func (r *regexLib) CompileOptions() []cel.EnvOption {
	optionalString := cel.OptionalType(cel.StringType)
	opts := []cel.EnvOption{
		cel.Function(regexExtract, cel.Overload("regex_extract_string_string", []*cel.Type{cel.StringType, cel.StringType}, optionalString,
			cel.BinaryBinding(extract))),
		cel.Function(regexExtractAll, cel.Overload("regex_extractAll_string_string", []*cel.Type{cel.StringType, cel.StringType}, cel.ListType(cel.StringType),
			cel.BinaryBinding(extractAll))),
		cel.Function(regexReplace,
			cel.Overload("regex_replace_string_string_string", []*cel.Type{cel.StringType, cel.StringType, cel.StringType}, cel.StringType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					target := args[0].(types.String)
					regexStr := args[1].(types.String)
					replaceStr := args[2].(types.String)
					return replaceAll(string(target), string(regexStr), string(replaceStr))
				})),
			cel.Overload("regex_replace_string_string_string_int", []*cel.Type{cel.StringType, cel.StringType, cel.StringType, cel.IntType}, cel.StringType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					target := args[0].(types.String)
					regexStr := args[1].(types.String)
					replaceStr := args[2].(types.String)
					count := args[3].(types.Int)
					return replaceCount(string(target), string(regexStr), string(replaceStr), int64(count))
				}))),
	}
	return opts
}

// ProgramOptions implements the cel.Library interface method
func (r *regexLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func compileRegex(regexStr string) (*regexp.Regexp, error) {
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, fmt.Errorf("given regex is invalid: %w", err)
	}
	return re, nil
}

// Initializing regular expression patterns for validating replacement strings.
var (
	reGroupNum     = regexp.MustCompile(`\$(\d+)`)
	reGroupName    = regexp.MustCompile(`\$\{([a-zA-Z_][a-zA-Z0-9_]*)\}`)
	reGroupInvalid = regexp.MustCompile(`\$[a-zA-Z_][a-zA-Z0-9_]*`)
)

func validateReplacementString(re *regexp.Regexp, replaceStr string) error {
	// If there are no $N or ${name} patterns, skip validation.
	if !strings.Contains(replaceStr, "$") {
		return nil
	}

	groupNames := re.SubexpNames()
	groupCount := len(groupNames) - 1 // Exclude group 0 (whole match)
	// Find all $N patterns in the replacement string and validate them.
	matches := reGroupNum.FindAllStringSubmatch(replaceStr, -1)
	for _, m := range matches {
		if len(m) > 1 {
			idx, _ := strconv.Atoi(m[1])
			if idx < 0 || idx > groupCount {
				return fmt.Errorf("replacement string references group $%d, but regex has only %d group(s)", idx, groupCount)
			}
		}
	}

	if strings.Contains(replaceStr, "${") {
		validNames := make(map[string]struct{})
		for _, name := range groupNames {
			if name != "" {
				validNames[name] = struct{}{}
			}
		}
		// If there are named groups, validate them against the defined group names.
		nameMatches := reGroupName.FindAllStringSubmatch(replaceStr, -1)
		for _, m := range nameMatches {
			if len(m) > 1 {
				if _, ok := validNames[m[1]]; !ok {
					return fmt.Errorf("invalid capture group name in replacement string: %s", m[1])
				}
			}
		}
	}

	// Check for invalid $word references (e.g., $a)
	invalids := reGroupInvalid.FindAllString(replaceStr, -1)
	for _, m := range invalids {
		// If not matched by $N, it's invalid
		if !reGroupNum.MatchString(m) {
			return fmt.Errorf("invalid group reference: %s", m)
		}
	}

	return nil
}

func extract(target, regexStr ref.Val) ref.Val {
	t := string(target.(types.String))
	r := string(regexStr.(types.String))
	re, err := compileRegex(r)
	if err != nil {
		return types.WrapErr(err)
	}

	if len(re.SubexpNames())-1 > 1 {
		return types.WrapErr(fmt.Errorf("regular expression has more than one capturing group: %q", r))
	}

	matches := re.FindStringSubmatch(t)
	if len(matches) == 0 {
		return types.OptionalNone
	}

	// If there is a capturing group, return the first match; otherwise, return the whole match.
	if len(matches) > 1 {
		capturedGroup := matches[1]
		// If optional group is empty, return OptionalNone.
		if capturedGroup == "" {
			return types.OptionalNone
		}
		return types.OptionalOf(types.String(capturedGroup))
	}
	return types.OptionalOf(types.String(matches[0]))
}

func extractAll(target, regexStr ref.Val) ref.Val {
	t := string(target.(types.String))
	r := string(regexStr.(types.String))
	re, err := compileRegex(r)
	if err != nil {
		return types.WrapErr(err)
	}

	groupCount := len(re.SubexpNames()) - 1
	if groupCount > 1 {
		return types.WrapErr(fmt.Errorf("regular expression has more than one capturing group: %q", r))
	}

	matches := re.FindAllStringSubmatch(t, -1)
	result := make([]string, 0, len(matches))
	if len(matches) == 0 {
		return types.NewStringList(types.DefaultTypeAdapter, result)
	}

	if groupCount == 1 {
		for _, match := range matches {
			if match[1] != "" {
				result = append(result, match[1])
			}
		}
	} else {
		for _, match := range matches {
			result = append(result, match[0])
		}
	}
	return types.NewStringList(types.DefaultTypeAdapter, result)
}

func replaceAll(target, regexStr, replaceStr string) ref.Val {
	re, err := compileRegex(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	if err := validateReplacementString(re, replaceStr); err != nil {
		return types.WrapErr(err)
	}
	a := types.String(re.ReplaceAllString(target, replaceStr))
	return a
}

func replaceCount(target, regexStr, replaceStr string, replaceCount int64) ref.Val {
	re, err := compileRegex(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	if err := validateReplacementString(re, replaceStr); err != nil {
		return types.WrapErr(err)
	}

	if replaceCount == -1 {
		return types.String(re.ReplaceAllString(target, replaceStr))
	}

	if replaceCount <= 0 {
		return types.String(target)
	}

	if replaceCount > math.MaxInt32 {
		return errIntOverflow
	}

	matches := re.FindAllStringSubmatchIndex(target, int(replaceCount))
	if len(matches) == 0 {
		return types.String(target)
	}

	var builder strings.Builder
	builder.Grow(len(target))
	lastIndex := 0
	// Reuse this buffer to reduce allocations
	var expanded []byte
	for _, match := range matches {
		builder.WriteString(target[lastIndex:match[0]])
		// Reset slice length but keep capacity
		expanded = expanded[:0]
		expanded = re.ExpandString(expanded, replaceStr, target, match)
		builder.Write(expanded)
		lastIndex = match[1]
	}
	builder.WriteString(target[lastIndex:])
	return types.String(builder.String())
}
