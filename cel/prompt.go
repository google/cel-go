// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cel

import (
	"fmt"
	"strings"

	"github.com/google/cel-go/common/decls"
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/types"
)

// Prompt represents the core components of an LLM prompt based on a CEL environment.
//
// All fields of the prompt may be overwritten / modified with support for rendering the
// prompt to a human-readable string.
type Prompt struct {
	// Persona indicates something about the kind of user making the request
	Persona string

	// FormatRules indicate how the LLM should generate its output
	FormatRules string

	// Variables, macros, and functions represent additional context helpful in constraining generation
	Variables []string
	Macros    []string
	Functions []string

	// GeneralUsage specifies additional context on how CEL should be used.
	GeneralUsage string
}

// Render renders the user prompt with the associated context from the prompt template
// for use with LLM generators.
func (p *Prompt) Render(userPrompt string) string {
	return fmt.Sprintf(`%s
%s

Only use the following variables in expressions:
  
%s
Only use the following macros in expressions:

%s
Only use function signatures from the list of supported functions:

%s
General Advice for authoring CEL:
%s
Prompt: %s`,
		p.Persona,
		p.FormatRules,
		"  "+strings.Join(p.Variables, "\n  "),
		strings.Join(p.Macros, "\n"),
		strings.Join(p.Functions, "\n"),
		p.GeneralUsage,
		userPrompt)
}

// NewPrompt creates a prompt template from a CEL environment.
func NewPrompt(env *Env) *Prompt {
	return &Prompt{
		Persona:      defaultPersona,
		FormatRules:  defaultFormatRules,
		GeneralUsage: defaultGeneralUsage,
		Variables:    describeVariables(env.Variables()),
		Macros:       describeMacros(env.Macros()),
		Functions:    describeFunctions(env.Functions()),
	}
}

func describeFunctions(funcs map[string]*decls.FunctionDecl) []string {
	calls := make([]string, 0, len(funcs))
	selectors := []string{
		"object.FIELD -> <FIELD_TYPE>, select a FIELD whose name must be a valid CEL identifier from a structured object.",
		"map.string -> <V>, select a map key whose name is a valid CEL identifier.",
	}
	opsByArity := map[int][]string{}
	for _, f := range funcs {
		if opName, found := operators.FindReverseBinaryOperator(f.Name()); found {
			ops, found := opsByArity[2]
			if !found {
				ops = []string{}
			}
			ops = append(ops, formatBinaryOperator(opName, f.OverloadDecls()))
			opsByArity[2] = ops
			continue
		}
		switch f.Name() {
		case operators.Negate, operators.LogicalNot:
			opName := "!"
			signatures := []string{"!bool -> bool"}
			if f.Name() == operators.Negate {
				opName = "-"
				signatures = []string{"-int -> int", "-double -> double"}
			}
			ops, found := opsByArity[1]
			if !found {
				ops = []string{}
			}
			ops = append(ops, fmt.Sprintf("%q signatures:\n  %s", opName, strings.Join(signatures, ", ")))
		case operators.Index:
			selectors = append(selectors,
				"list[int] -> <V>, select a value by integer index from a list, the index must be non-negative",
				"map[<K>] -> <V>, select a parameterized value <V> from a map by parameterized key type <K>, only string, int, uint, and bool types are valid types for <K>",
			)
		case operators.OptIndex:
			selectors = append(selectors,
				"list[?int] -> optional_type(<V>), select a value by integer index from a list, the index must be non-negative. If the index is present return an optional_type(<V>) value, otherwise optional.none()",
				"map[?<K>] -> optional_type(<V>), select a parameterized value <V> from a map by parameterized key type <K>, only string, int, uint, and bool types are valid types for <K>. If the key <K> is present return an optional_type(<V>) value, otherwise optional.none()",
			)
		case operators.OptSelect:
			selectors = append(selectors,
				"object.?FIELD -> optional_type(<FIELD_TYPE>), select a FIELD whose name must be a valid CEL identifier from a structured object. If the field is set to a non-empty value return an optional_type(<FIELD_TYPE>) for the value, otherwise optional.none().",
				"map.?string -> optional_type(<V>), select a parameterized value <V> from a map by a string-typed key. If the string key is present return an optional_type(<V>) value, otherwise optional.none()",
			)
		case operators.Conditional:
			opsByArity[3] = []string{
				"conditional signatures:\n  bool ? <T(TRUE)> : <T(FALSE)> -> <T>, when the bool condition evaluates to true, return the truthy value, otherwise return the falsey value. Both the truthy and falsy values must have the same type <T>",
			}
		default:
			calls = append(calls, formatCall(f.Name(), f.OverloadDecls()))
		}
	}
	operators := []string{}
	if unary, found := opsByArity[1]; found {
		operators = append(operators, fmt.Sprintf("unary operators:\n%s\n", strings.Join(unary, "\n")))
	}
	if binary, found := opsByArity[2]; found {
		operators = append(operators, fmt.Sprintf("binary operators:\n%s\n", strings.Join(binary, "\n")))
	}
	if ternary, found := opsByArity[3]; found {
		operators = append(operators, fmt.Sprintf("ternary operators:\n%s\n", strings.Join(ternary, "\n")))
	}

	content := []string{}
	if len(operators) > 0 {
		content = append(content, fmt.Sprintf(
			"The following operators exist for equality, comparison, key/value containment, and logic tests:\n\n%s\n",
			strings.Join(operators, "\n"),
		))
	}
	if len(calls) > 0 {
		content = append(content, fmt.Sprintf(
			"The following global and receiver-style functions are supported:\n\n%s\n",
			strings.Join(calls, "\n"),
		))
	}
	if len(selectors) > 0 {
		content = append(content, fmt.Sprintf(
			"List and map types support indexing, and field selection is supported for structured objects and maps:\n\n  %s\n",
			strings.Join(selectors, "\n  "),
		))
	}
	return content
}

func formatBinaryOperator(opName string, overloads []*decls.OverloadDecl) string {
	signatures := make([]string, len(overloads))
	for i, o := range overloads {
		args := o.ArgTypes()
		arg0 := describeCELType(args[0])
		arg1 := describeCELType(args[1])
		ret := describeCELType(o.ResultType())
		signatures[i] = fmt.Sprintf("%s %s %s -> %s", arg0, opName, arg1, ret)
	}
	return fmt.Sprintf("%q signatures:\n  %s\n", opName, strings.Join(signatures, "\n  "))
}

func formatCall(funcName string, overloads []*decls.OverloadDecl) string {
	signatures := make([]string, len(overloads))
	for i, o := range overloads {
		args := make([]string, len(o.ArgTypes()))
		ret := describeCELType(o.ResultType())
		for j, a := range o.ArgTypes() {
			args[j] = describeCELType(a)
		}
		if o.IsMemberFunction() {
			target := args[0]
			args = args[1:]
			signatures[i] = fmt.Sprintf("%s.%s(%s) -> %s", target, funcName, strings.Join(args, ", "), ret)
		} else {
			signatures[i] = fmt.Sprintf("%s(%s) -> %s", funcName, strings.Join(args, ", "), ret)
		}
	}
	return fmt.Sprintf("%q signatures:\n  %s\n", funcName, strings.Join(signatures, "\n  "))
}

func describeMacros(macros []Macro) []string {
	macroMap := map[string][]Macro{}
	for _, m := range macros {
		macroFuncs, found := macroMap[m.Function()]
		if !found {
			macroFuncs = []Macro{}
		}
		macroFuncs = append(macroFuncs, m)
		macroMap[m.Function()] = macroFuncs
	}
	macDescs := []string{}
	for funcName, macroSet := range macroMap {
		macroSignatures := make([]string, len(macroSet))
		for i, m := range macroSet {
			signature := ""
			args := []string{}
			if m.ArgCount() == 0 && strings.Contains(m.MacroKey(), "*") {
				args = []string{"ARG0", "ARG1", "... ARGN"}
			}
			for i := 0; i < m.ArgCount(); i++ {
				args = append(args, fmt.Sprintf("ARG%d", i))
			}
			if m.IsReceiverStyle() {
				signature = fmt.Sprintf("TARGET.%s(%s)", funcName, strings.Join(args, ", "))
			} else {
				signature = fmt.Sprintf("%s(%s)", funcName, strings.Join(args, ", "))
			}
			macroSignatures[i] = signature
		}
		macDescs = append(macDescs,
			fmt.Sprintf("%q signatures:\n  %s", funcName, strings.Join(macroSignatures, "\n  ")))
	}
	return macDescs
}

func describeVariables(vars []*decls.VariableDecl) []string {
	varDescs := make([]string, len(vars))
	for i, v := range vars {
		varDescs[i] = fmt.Sprintf("%s // type: %s", v.Name(), describeCELType(v.Type()))
	}
	return varDescs
}

func describeCELType(t *Type) string {
	if t.Kind() == types.TypeParamKind {
		return fmt.Sprintf("<%s>", t.TypeName())
	}
	return t.String()
}

var (
	wrapperTypes = map[types.Kind]string{
		types.BoolKind:   "google.protobuf.BoolValue",
		types.BytesKind:  "google.protobuf.BytesValue",
		types.DoubleKind: "google.protobuf.DoubleValue",
		types.IntKind:    "google.protobuf.Int64Value",
		types.StringKind: "google.protobuf.StringValue",
		types.UintKind:   "google.protobuf.UInt64Value",
	}
)

const (
	defaultPersona string = `
You are a software engineer with expertise in networking and application security
authoring boolean Common Expression Language (CEL) expressions to ensure firewall,
networking, authentication, and data access is only permitted when all conditions
are satisified.`

	defaultFormatRules string = `
Output your response as a CEL expression.
Write the expression in a professional format familiar to software engineers.
Provide a small comment above the CEL expression documenting the expression intent.`

	defaultGeneralUsage string = `
CEL supports Protocol Buffer and JSON types, as well as simple types and aggregate types.

Simple types include: bool, bytes, double, int, string, uint.
  * double literal must always include a decimal point: "1.0", "3.5", "-2.2"
  * uint literal values must be positive values suffixed with a 'u': "42u"
  * byte literal values are strings prefixed with a 'b': "b'1235'"

Aggregate types include: list and map.
  * list literal example: "['a', 'b', 'c']"
  * map literal example: "{'key1': 1, 'key2': 2}"
  * Only int, uint, string, and bool literals are valid map keys.

If the user asks to check whether one value contains another, or asks about presence, use the following rules to disambiguate what they mean:
  * For list values, use the "in" operator to test for a value in a list: "1 in [0, 1, 2]"
  * For maps, use the "in" operator to test for a key in a map: "'user-agent' in request.headers"
  * For objects, use the "has" macro to test for field presence: "has(msg.field)"
  * For strings, use the "string.contains(string)" method to test if the target string contains the argument substring
  * You MUST NOT use "has()" as a member function, e.g. request.headers.has("key")
  * Maps containing HTTP headers must always use lower-cased string keys.

Comments start with two-forward slashes followed by text and a newline.`
)
