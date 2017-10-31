package ast

import (
	"celgo/common"
)

// Expression is the common interface for all AST expressions.
type Expression interface {
	// Id is the id of an expression, unique within a parse tree.
	Id() int64

	// Location is the source-text location of the expression.
	Location() common.Location

	// String returns a string representation of the expression.
	String() string

	// writeDebugString writes the detailed string representation of an the expression to the supplied debugWriter.
	writeDebugString(w *debugWriter)
}
