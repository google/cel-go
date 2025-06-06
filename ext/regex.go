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
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func Regex(options ...RegexOptions) cel.EnvOption {
	s := &regexLib{
		version: math.MaxUint32,
	}
	for _, o := range options {
		s = o(s)
	}
	return cel.Lib(s)
}

const (
	regexReplace    = "regex.replace"
	capture         = "regex.capture"
	captureAll      = "regex.captureAll"
	captureAllNamed = "regex.captureAllNamed"
)

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
		cel.Function(capture, cel.Overload("regex_capture_string_string", []*cel.Type{cel.StringType, cel.StringType}, optionalString,
			cel.BinaryBinding(func(target, regexStr ref.Val) ref.Val {
				t := target.(types.String)
				r := regexStr.(types.String)
				return captureFirstMatch(string(t), string(r))
			}))),
		cel.Function(captureAll, cel.Overload("regex_captureAll_string_string", []*cel.Type{cel.StringType, cel.StringType}, cel.ListType(cel.StringType),
			cel.BinaryBinding(func(target, regexStr ref.Val) ref.Val {
				t := target.(types.String)
				r := regexStr.(types.String)
				return captureAllMatches(string(t), string(r))
			}))),
		cel.Function(captureAllNamed, cel.Overload("regex_captureAllNamed_string_string", []*cel.Type{cel.StringType, cel.StringType}, cel.MapType(cel.StringType, cel.StringType),
			cel.BinaryBinding(func(target, regexStr ref.Val) ref.Val {
				t := target.(types.String)
				r := regexStr.(types.String)
				return captureAllNamedGroups(string(t), string(r))
			}))),
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

// ProgramOptions implements cel.Library.
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

func captureFirstMatch(target, regexStr string) ref.Val {
	re, err := compileRegex(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	matches := re.FindStringSubmatch(target)
	if len(matches) == 0 {
		return types.OptionalNone
	}

	// If there is a capturing group, return the first group; otherwise, return the whole match.
	if len(matches) > 1 {
		return types.OptionalOf(types.String(matches[1]))
	}
	return types.OptionalOf(types.String(matches[0]))
}

func captureAllMatches(target, regexStr string) ref.Val {
	re, err := compileRegex(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	matches := re.FindAllStringSubmatch(target, -1)
	var result []string
	if len(matches) == 0 {
		return types.NewStringList(types.DefaultTypeAdapter, result)
	}

	hasCaptureGroups := len(matches[0]) > 1
	for _, match := range matches {
		if hasCaptureGroups {
			for i := 1; i < len(match); i++ {
				if match[i] != "" {
					result = append(result, match[i])
				}
			}
		} else {
			result = append(result, match[0])
		}
	}
	return types.NewStringList(types.DefaultTypeAdapter, result)
}

func captureAllNamedGroups(target, regexStr string) ref.Val {
	re, err := compileRegex(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	result := make(map[string]string)
	matches := re.FindAllStringSubmatch(target, -1)
	if len(matches) == 0 {
		return types.NewStringStringMap(types.DefaultTypeAdapter, result)
	}

	groupNames := re.SubexpNames()
	for _, match := range matches {
		for i, name := range groupNames {
			if i < len(match) && name != "" && match[i] != "" {
				result[name] = match[i]
			}
		}
	}
	return types.NewStringStringMap(types.DefaultTypeAdapter, result)
}

func replaceAll(target, regexStr, replaceStr string) ref.Val {
	re, err := compileRegex(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	return types.String(re.ReplaceAllString(target, replaceStr))
}

func replaceCount(target, regexStr, replaceStr string, replaceCount int64) ref.Val {
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return types.WrapErr(err)
	}

	if replaceCount == -1 {
		return types.String(re.ReplaceAllString(target, replaceStr))
	}

	matches := re.FindAllStringSubmatchIndex(target, int(replaceCount))
	if len(matches) == 0 || replaceCount == 0 {
		return types.String(target)
	}

	var builder strings.Builder
	lastIndex := 0
	for _, match := range matches {
		builder.WriteString(target[lastIndex:match[0]])

		var expanded []byte
		expanded = re.ExpandString(expanded, replaceStr, target, match)
		builder.Write(expanded)

		lastIndex = match[1]
	}

	builder.WriteString(target[lastIndex:])
	return types.String(builder.String())
}
