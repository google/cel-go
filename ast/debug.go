package ast

import (
	"bytes"
	"fmt"
	"strings"
)

// DebugAdorner returns debug metadata that will be tacked on to the string representation of an expression.
type DebugAdorner interface {
	GetMetadata(e Expression) string
}

type emptyDebugAdorner struct {
}

var EmptyAdorner DebugAdorner = &emptyDebugAdorner{}

func (a *emptyDebugAdorner) GetMetadata(e Expression) string {
	return ""
}

func ToDebugString(e Expression) string {
	return ToAdornedDebugString(e, EmptyAdorner)
}

func ToAdornedDebugString(e Expression, adorner DebugAdorner) string {
	w := newDebugWriter(adorner)
	e.writeDebugString(w)
	return w.String()
}

// debugWriter is used to print out pretty-printed debug strings.
type debugWriter struct {
	adorner   DebugAdorner
	buffer    bytes.Buffer
	indent    int
	lineStart bool
}

func newDebugWriter(a DebugAdorner) *debugWriter {
	return &debugWriter{
		adorner:   a,
		indent:    0,
		lineStart: true,
	}
}

func (w *debugWriter) append(s string) {
	w.doIndent()
	w.buffer.WriteString(s)
}

func (w *debugWriter) appendFormat(f string, args ...interface{}) {
	w.append(fmt.Sprintf(f, args...))
}

func (w *debugWriter) doIndent() {
	if w.lineStart {
		w.lineStart = false
		w.buffer.WriteString(strings.Repeat("  ", w.indent))
	}
}

func (w *debugWriter) adorn(e Expression) {
	w.append(w.adorner.GetMetadata(e))
}

func (w *debugWriter) appendExpression(e Expression) {
	e.writeDebugString(w)
}

func (w *debugWriter) appendLine() {
	w.buffer.WriteString("\n")
	w.lineStart = true
}

func (w *debugWriter) addIndent() {
	w.indent++
}

func (w *debugWriter) removeIndent() {
	w.indent--
	if w.indent < 0 {
		panic("negative indent")
	}
}

func (w *debugWriter) String() string {
	return w.buffer.String()
}
