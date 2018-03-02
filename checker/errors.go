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
	"celgo/common"
	"celgo/semantics/types"
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

func (errors *TypeErrors) typeDoesNotSupportFieldSelection(l common.Location, t types.Type) {
	errors.ReportError(l, "type '%s' does not support field selection", t)
}

func (errors *TypeErrors) undefinedField(l common.Location, field string) {
	errors.ReportError(l, "undefined field '%s'", field)
}

func (errors *TypeErrors) fieldDoesNotSupportPresenceCheck(l common.Location, field string) {
	errors.ReportError(l, "field '%s' does not support presence check", field)
}

func (errors *TypeErrors) noMatchingOverload(l common.Location, name string, args []types.Type, isInstance bool) {
	signature := types.FormatFunction(nil, args, isInstance)
	errors.ReportError(l, "found no matching overload for '%s' applied to '%s'", name, signature)
}

func (errors *TypeErrors) aggregateTypeMismatch(l common.Location, aggregate types.Type, member types.Type) {
	errors.ReportError(
		l,
		"type '%s' does not match previous type '%s' in aggregate. Use 'dyn(x)' to make the aggregate dynamic.",
		member,
		aggregate)
}

func (errors *TypeErrors) notAType(l common.Location, t types.Type) {
	errors.ReportError(l, "'%s(%v)' is not a type", t.String(), t.Kind())
}

func (errors *TypeErrors) notAMessageType(l common.Location, t types.Type) {
	errors.ReportError(l, "'%s' is not a message type", t.String())
}

func (errors *TypeErrors) fieldTypeMismatch(l common.Location, name string, field types.Type, value types.Type) {
	errors.ReportError(l, "expected type of field '%s' is '%s' but provided type is '%s'",
		name, field.String(), value.String())
}

func (errors *TypeErrors) unexpectedFailedResolution(l common.Location, typeName string) {
	errors.ReportError(l, "[internal] unexpected failed resolution of '%s'", typeName)
}

func (errors *TypeErrors) notAComprehensionRange(l common.Location, t types.Type) {
	errors.ReportError(l, "expression of type '%s' cannot be range of a comprehension (must be list, map, or dynamic)",
		t.String())
}

func (errors *TypeErrors) typeMismatch(l common.Location, expected types.Type, actual types.Type) {
	errors.ReportError(l, "expected type '%s' but found '%s'", expected, actual)
}
