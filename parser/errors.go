package parser

import (
	"fmt"

	"celgo/common"
)

// ParseErrors is a specialization of Errors.
type ParseErrors struct {
	*common.Errors
}

func (e *ParseErrors) syntaxError(l common.Location, message string) {
	e.ReportError(l, fmt.Sprintf("Syntax error: %s", message))
}

func (e *ParseErrors) invalidHasArgument(l common.Location) {
	e.ReportError(l, "The argument to the function 'has' must be a field selection")
}

func (e *ParseErrors) argumentIsNotIdent(l common.Location) {
	e.ReportError(l, "The argument must be a simple name")
}

func (e *ParseErrors) notAQualifiedName(l common.Location) {
	e.ReportError(l, "expected a qualified name")
}
