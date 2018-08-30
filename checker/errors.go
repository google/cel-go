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
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	commonpb "github.com/google/cel-go/common"
)

// typeErrors is a specialization of Errors.
type typeErrors struct {
	*commonpb.Errors
}

func (e *typeErrors) undeclaredReference(l commonpb.Location, container string, name string) {
	e.ReportError(l, "undeclared reference to '%s' (in container '%s')", name, container)
}

func (e *typeErrors) expressionDoesNotSelectField(l commonpb.Location) {
	e.ReportError(l, "expression does not select a field")
}

func (e *typeErrors) typeDoesNotSupportFieldSelection(l commonpb.Location, t *checkedpb.Type) {
	e.ReportError(l, "type '%s' does not support field selection", t)
}

func (e *typeErrors) undefinedField(l commonpb.Location, field string) {
	e.ReportError(l, "undefined field '%s'", field)
}

func (e *typeErrors) fieldDoesNotSupportPresenceCheck(l commonpb.Location, field string) {
	e.ReportError(l, "field '%s' does not support presence check", field)
}

func (e *typeErrors) overlappingOverload(l commonpb.Location, name string, overloadId1 string, f1 *checkedpb.Type,
	overloadId2 string, f2 *checkedpb.Type) {
	e.ReportError(l, "overlapping overload for name '%s' (type '%s' with overloadId: '%s' cannot be distinguished from '%s' with "+
		"overloadId: '%s')", name, FormatCheckedType(f1), overloadId1, FormatCheckedType(f2), overloadId2)
}

func (e *typeErrors) overlappingMacro(l commonpb.Location, name string, args int) {
	e.ReportError(l, "overload for name '%s' with %d argument(s) overlaps with predefined macro",
		name, args)
}

func (e *typeErrors) noMatchingOverload(l commonpb.Location, name string, args []*checkedpb.Type, isInstance bool) {
	signature := formatFunction(nil, args, isInstance)
	e.ReportError(l, "found no matching overload for '%s' applied to '%s'", name, signature)
}

func (e *typeErrors) aggregateTypeMismatch(l commonpb.Location, aggregate *checkedpb.Type, member *checkedpb.Type) {
	e.ReportError(
		l,
		"type '%s' does not match previous type '%s' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.",
		FormatCheckedType(member),
		FormatCheckedType(aggregate))
}

func (e *typeErrors) notAType(l commonpb.Location, t *checkedpb.Type) {
	e.ReportError(l, "'%s(%v)' is not a type", FormatCheckedType(t), t)
}

func (e *typeErrors) notAMessageType(l commonpb.Location, t *checkedpb.Type) {
	e.ReportError(l, "'%s' is not a message type", FormatCheckedType(t))
}

func (e *typeErrors) fieldTypeMismatch(l commonpb.Location, name string, field *checkedpb.Type, value *checkedpb.Type) {
	e.ReportError(l, "expected type of field '%s' is '%s' but provided type is '%s'",
		name, FormatCheckedType(field), FormatCheckedType(value))
}

func (e *typeErrors) unexpectedFailedResolution(l commonpb.Location, typeName string) {
	e.ReportError(l, "[internal] unexpected failed resolution of '%s'", typeName)
}

func (e *typeErrors) notAComprehensionRange(l commonpb.Location, t *checkedpb.Type) {
	e.ReportError(l, "expression of type '%s' cannot be range of a comprehension (must be list, map, or dynamic)",
		FormatCheckedType(t))
}

func (e *typeErrors) typeMismatch(l commonpb.Location, expected *checkedpb.Type, actual *checkedpb.Type) {
	e.ReportError(l, "expected type '%s' but found '%s'",
		FormatCheckedType(expected), FormatCheckedType(actual))
}

func formatFunction(resultType *checkedpb.Type, argTypes []*checkedpb.Type, isInstance bool) string {
	result := ""
	if isInstance {
		target := argTypes[0]
		argTypes = argTypes[1:]

		result += FormatCheckedType(target)
		result += "."
	}

	result += "("
	for i, arg := range argTypes {
		if i > 0 {
			result += ", "
		}
		result += FormatCheckedType(arg)
	}
	result += ")"
	if resultType != nil {
		result += " -> "
		result += FormatCheckedType(resultType)
	}

	return result
}
