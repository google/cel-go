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

// TypeErrors is a specialization of Errors.
type TypeErrors struct {
	*common.Errors
}

func (errors *TypeErrors) undeclaredReference(l common.Location, container string, name string) {
	errors.ReportError(l, "undeclared reference to '%s' (in container '%s')", name, container)
}

func (errors *TypeErrors) expressionDoesNotSelectField(l common.Location) {
	errors.ReportError(l, "expression does not select a field")
}

func (errors *TypeErrors) typeDoesNotSupportFieldSelection(l common.Location, t *checked.Type) {
	errors.ReportError(l, "type '%s' does not support field selection", t)
}

func (errors *TypeErrors) undefinedField(l common.Location, field string) {
	errors.ReportError(l, "undefined field '%s'", field)
}

func (errors *TypeErrors) fieldDoesNotSupportPresenceCheck(l common.Location, field string) {
	errors.ReportError(l, "field '%s' does not support presence check", field)
}

func (errors *TypeErrors) noMatchingOverload(l common.Location, name string, args []*checked.Type, isInstance bool) {
	signature := formatFunction(nil, args, isInstance)
	errors.ReportError(l, "found no matching overload for '%s' applied to '%s'", name, signature)
}

func (errors *TypeErrors) aggregateTypeMismatch(l common.Location, aggregate *checked.Type, member *checked.Type) {
	errors.ReportError(
		l,
		"type '%s' does not match previous type '%s' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.",
		types.FormatType(member),
		types.FormatType(aggregate))
}

func (errors *TypeErrors) notAType(l common.Location, t *checked.Type) {
	errors.ReportError(l, "'%s(%v)' is not a type", types.FormatType(t), t)
}

func (errors *TypeErrors) notAMessageType(l common.Location, t *checked.Type) {
	errors.ReportError(l, "'%s' is not a message type", types.FormatType(t))
}

func (errors *TypeErrors) fieldTypeMismatch(l common.Location, name string, field *checked.Type, value *checked.Type) {
	errors.ReportError(l, "expected type of field '%s' is '%s' but provided type is '%s'",
		name, types.FormatType(field), types.FormatType(value))
}

func (errors *TypeErrors) unexpectedFailedResolution(l common.Location, typeName string) {
	errors.ReportError(l, "[internal] unexpected failed resolution of '%s'", typeName)
}

func (errors *TypeErrors) notAComprehensionRange(l common.Location, t *checked.Type) {
	errors.ReportError(l, "expression of type '%s' cannot be range of a comprehension (must be list, map, or dynamic)",
		types.FormatType(t))
}

func (errors *TypeErrors) typeMismatch(l common.Location, expected *checked.Type, actual *checked.Type) {
	errors.ReportError(l, "expected type '%s' but found '%s'",
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
