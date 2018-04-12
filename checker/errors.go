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
	"github.com/google/cel-go/checker/types"
	"github.com/google/cel-go/common"
	"github.com/google/cel-spec/proto/checked/v1/checked"
)

// typeErrors is a specialization of Errors.
type typeErrors struct {
	*common.Errors
}

func (e *typeErrors) undeclaredReference(l common.Location, container string, name string) {
	e.ReportError(l, "undeclared reference to '%s' (in container '%s')", name, container)
}

func (e *typeErrors) expressionDoesNotSelectField(l common.Location) {
	e.ReportError(l, "expression does not select a field")
}

func (e *typeErrors) typeDoesNotSupportFieldSelection(l common.Location, t *checked.Type) {
	e.ReportError(l, "type '%s' does not support field selection", t)
}

func (e *typeErrors) undefinedField(l common.Location, field string) {
	e.ReportError(l, "undefined field '%s'", field)
}

func (e *typeErrors) fieldDoesNotSupportPresenceCheck(l common.Location, field string) {
	e.ReportError(l, "field '%s' does not support presence check", field)
}

func (e *typeErrors) overlappingOverload(l common.Location, name string, overloadId1 string, f1 *checked.Type,
	overloadId2 string, f2 *checked.Type) {
	e.ReportError(l, "overlapping overload for name '%s' (type '%s' with overloadId: '%s' cannot be distinguished from '%s' with "+
		"overloadId: '%s')", name, types.FormatType(f1), overloadId1, types.FormatType(f2), overloadId2)
}

func (e *typeErrors) overlappingMacro(l common.Location, name string, args int) {
	e.ReportError(l, "overload for name '%s' with %d argument(s) overlaps with predefined macro",
		name, args)
}

func (e *typeErrors) noMatchingOverload(l common.Location, name string, args []*checked.Type, isInstance bool) {
	signature := formatFunction(nil, args, isInstance)
	e.ReportError(l, "found no matching overload for '%s' applied to '%s'", name, signature)
}

func (e *typeErrors) aggregateTypeMismatch(l common.Location, aggregate *checked.Type, member *checked.Type) {
	e.ReportError(
		l,
		"type '%s' does not match previous type '%s' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.",
		types.FormatType(member),
		types.FormatType(aggregate))
}

func (e *typeErrors) notAType(l common.Location, t *checked.Type) {
	e.ReportError(l, "'%s(%v)' is not a type", types.FormatType(t), t)
}

func (e *typeErrors) notAMessageType(l common.Location, t *checked.Type) {
	e.ReportError(l, "'%s' is not a message type", types.FormatType(t))
}

func (e *typeErrors) fieldTypeMismatch(l common.Location, name string, field *checked.Type, value *checked.Type) {
	e.ReportError(l, "expected type of field '%s' is '%s' but provided type is '%s'",
		name, types.FormatType(field), types.FormatType(value))
}

func (e *typeErrors) unexpectedFailedResolution(l common.Location, typeName string) {
	e.ReportError(l, "[internal] unexpected failed resolution of '%s'", typeName)
}

func (e *typeErrors) notAComprehensionRange(l common.Location, t *checked.Type) {
	e.ReportError(l, "expression of type '%s' cannot be range of a comprehension (must be list, map, or dynamic)",
		types.FormatType(t))
}

func (e *typeErrors) typeMismatch(l common.Location, expected *checked.Type, actual *checked.Type) {
	e.ReportError(l, "expected type '%s' but found '%s'",
		types.FormatType(expected), types.FormatType(actual))
}

func formatFunction(resultType *checked.Type, argTypes []*checked.Type, isInstance bool) string {
	result := ""
	if isInstance {
		target := argTypes[0]
		argTypes = argTypes[1:]

		result += types.FormatType(target)
		result += "."
	}

	result += "("
	for i, arg := range argTypes {
		if i > 0 {
			result += ", "
		}
		result += types.FormatType(arg)
	}
	result += ")"
	if resultType != nil {
		result += " -> "
		result += types.FormatType(resultType)
	}

	return result
}
