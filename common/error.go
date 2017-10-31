package common

import (
	"fmt"
	"strings"
)

// Error represents an error message.
type Error struct {
	// Location within the source.
	Location Location

	// Message is the actual error message.
	Message string
}

// ToDisplayString returns the error in a user-friendly format.
func (e *Error) ToDisplayString() string {

	result := fmt.Sprintf("ERROR: %s:%d:%d: %s", e.Location.Source().Name(), e.Location.Line(), e.Location.Column(), e.Message)

	if snippet, found := e.Location.Source().Snippet(e.Location.Line()); found {
		result += "\n | "
		result += snippet
		result += "\n | "
		result += strings.Repeat(".", int(e.Location.Column()-1))
		result += "^"

	}
	return result
}
