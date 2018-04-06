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

// Package debug provides tools to print a parsed expression graph and
// adorn each expression element with additional metadata.
package debug

import (
	"bytes"
	"fmt"
	expr "github.com/google/cel-spec/proto/v1/syntax"
	"strconv"
	"strings"
)

// DebugAdorner returns debug metadata that will be tacked on to the string
// representation of an expression.
type DebugAdorner interface {
	// GetMetadata for the input context.
	GetMetadata(ctx interface{}) string
}

// DebugWriter manages writing expressions to an internal string.
type DebugWriter interface {
	fmt.Stringer

	// Buffer pushes an expression into an internal queue of expressions to
	// write to a string.
	Buffer(e *expr.Expr)
}

type emptyDebugAdorner struct {
}

var emptyAdorner DebugAdorner = &emptyDebugAdorner{}

func (a *emptyDebugAdorner) GetMetadata(e interface{}) string {
	return ""
}

func ToDebugString(e *expr.Expr) string {
	return ToAdornedDebugString(e, emptyAdorner)
}

func ToAdornedDebugString(e *expr.Expr, adorner DebugAdorner) string {
	w := newDebugWriter(adorner)
	w.Buffer(e)
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

func (w *debugWriter) Buffer(e *expr.Expr) {
	if e == nil {
		return
	}
	switch e.ExprKind.(type) {
	case *expr.Expr_ConstExpr:
		w.append(formatConstant(e.GetConstExpr()))
	case *expr.Expr_IdentExpr:
		w.append(e.GetIdentExpr().Name)
	case *expr.Expr_SelectExpr:
		w.appendSelect(e.GetSelectExpr())
	case *expr.Expr_CallExpr:
		w.appendCall(e.GetCallExpr())
	case *expr.Expr_ListExpr:
		w.appendList(e.GetListExpr())
	case *expr.Expr_StructExpr:
		w.appendStruct(e.GetStructExpr())
	case *expr.Expr_ComprehensionExpr:
		w.appendComprehension(e.GetComprehensionExpr())
	}
	w.adorn(e)
}

func (w *debugWriter) appendSelect(sel *expr.Expr_Select) {
	w.Buffer(sel.Operand)
	w.append(".")
	w.append(sel.Field)
	if sel.TestOnly {
		w.append("~test-only~")
	}
}

func (w *debugWriter) appendCall(call *expr.Expr_Call) {
	if call.Target != nil {
		w.Buffer(call.Target)
		w.append(".")
	}
	w.append(call.Function)
	w.append("(")
	if len(call.GetArgs()) > 0 {
		w.addIndent()
		w.appendLine()
		for i, arg := range call.Args {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.Buffer(arg)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append(")")
}

func (w *debugWriter) appendList(list *expr.Expr_CreateList) {
	w.append("[")
	if len(list.GetElements()) > 0 {
		w.appendLine()
		w.addIndent()
		for i, elem := range list.Elements {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.Buffer(elem)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append("]")
}

func (w *debugWriter) appendStruct(obj *expr.Expr_CreateStruct) {
	if obj.MessageName != "" {
		w.appendObject(obj)
	} else {
		w.appendMap(obj)
	}
}

func (w *debugWriter) appendObject(obj *expr.Expr_CreateStruct) {
	w.append(obj.MessageName)
	w.append("{")
	if len(obj.Entries) > 0 {
		w.appendLine()
		w.addIndent()
		for i, entry := range obj.Entries {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.append(entry.GetFieldKey())
			w.append(":")
			w.Buffer(entry.Value)
			w.adorn(entry)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append("}")
}

func (w *debugWriter) appendMap(obj *expr.Expr_CreateStruct) {
	w.append("{")
	if len(obj.Entries) > 0 {
		w.appendLine()
		w.addIndent()
		for i, entry := range obj.Entries {
			if i > 0 {
				w.append(",")
				w.appendLine()
			}
			w.Buffer(entry.GetMapKey())
			w.append(":")
			w.Buffer(entry.Value)
			w.adorn(entry)
		}
		w.removeIndent()
		w.appendLine()
	}
	w.append("}")
}

func (w *debugWriter) appendComprehension(comprehension *expr.Expr_Comprehension) {
	w.append("__comprehension__(")
	w.addIndent()
	w.appendLine()
	w.append("// Variable")
	w.appendLine()
	w.append(comprehension.IterVar)
	w.append(",")
	w.appendLine()
	w.append("// Target")
	w.appendLine()
	w.Buffer(comprehension.IterRange)
	w.append(",")
	w.appendLine()
	w.append("// Accumulator")
	w.appendLine()
	w.append(comprehension.AccuVar)
	w.append(",")
	w.appendLine()
	w.append("// Init")
	w.appendLine()
	w.Buffer(comprehension.AccuInit)
	w.append(",")
	w.appendLine()
	w.append("// LoopCondition")
	w.appendLine()
	w.Buffer(comprehension.LoopCondition)
	w.append(",")
	w.appendLine()
	w.append("// LoopStep")
	w.appendLine()
	w.Buffer(comprehension.LoopStep)
	w.append(",")
	w.appendLine()
	w.append("// Result")
	w.appendLine()
	w.Buffer(comprehension.Result)
	w.append(")")
	w.removeIndent()
}

func formatConstant(c *expr.Constant) string {
	switch c.ConstantKind.(type) {
	case *expr.Constant_BoolValue:
		return fmt.Sprintf("%t", c.GetBoolValue())
	case *expr.Constant_BytesValue:
		return fmt.Sprintf("b\"%s\"", string(c.GetBytesValue()))
	case *expr.Constant_DoubleValue:
		return fmt.Sprintf("%v", c.GetDoubleValue())
	case *expr.Constant_Int64Value:
		return fmt.Sprintf("%d", c.GetInt64Value())
	case *expr.Constant_StringValue:
		return strconv.Quote(c.GetStringValue())
	case *expr.Constant_Uint64Value:
		return fmt.Sprintf("%du", c.GetUint64Value())
	case *expr.Constant_NullValue:
		return "null"
	default:
		panic("Unknown constant type")
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

func (w *debugWriter) adorn(e interface{}) {
	w.append(w.adorner.GetMetadata(e))
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
