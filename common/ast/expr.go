// Copyright 2023 Google LLC
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

package ast

import (
	"github.com/google/cel-go/common/types/ref"
)

// ExprKind represents the expression node kind.
type ExprKind int

const (
	// UnspecifiedExprKind represents an unset expression with no specified properties.
	UnspecifiedExprKind ExprKind = iota

	// CallKind represents a function call.
	CallKind

	// ComprehensionKind represents a comprehension expression generated by a macro.
	ComprehensionKind

	// IdentKind represents a simple variable, constant, or type identifier.
	IdentKind

	// ListKind represents a list literal expression.
	ListKind

	// LiteralKind represents a primitive scalar literal.
	LiteralKind

	// MapKind represents a map literal expression.
	MapKind

	// SelectKind represents a field selection expression.
	SelectKind

	// StructKind represents a struct literal expression.
	StructKind
)

// Expr represents the base expression node in a CEL abstract syntax tree.
//
// Depending on the `Kind()` value, the Expr may be converted to a concrete expression types
// as indicated by the `As<Kind>` methods.
type Expr interface {
	// ID of the expression as it appears in the AST
	ID() int64

	// Kind of the expression node. See ExprKind for the valid enum values.
	Kind() ExprKind

	// AsCall adapts the expr into a CallExpr
	//
	// The Kind() must be equal to a CallKind for the conversion to be well-defined.
	AsCall() CallExpr

	// AsComprehension adapts the expr into a ComprehensionExpr.
	//
	// The Kind() must be equal to a ComprehensionKind for the conversion to be well-defined.
	AsComprehension() ComprehensionExpr

	// AsIdent adapts the expr into an identifier string.
	//
	// The Kind() must be equal to an IdentKind for the conversion to be well-defined.
	AsIdent() string

	// AsLiteral adapts the expr into a constant ref.Val.
	//
	// The Kind() must be equal to a LiteralKind for the conversion to be well-defined.
	AsLiteral() ref.Val

	// AsList adapts the expr into a ListExpr.
	//
	// The Kind() must be equal to a ListKind for the conversion to be well-defined.
	AsList() ListExpr

	// AsMap adapts the expr into a MapExpr.
	//
	// The Kind() must be equal to a MapKind for the conversion to be well-defined.
	AsMap() MapExpr

	// AsSelect adapts the expr into a SelectExpr.
	//
	// The Kind() must be equal to a SelectKind for the conversion to be well-defined.
	AsSelect() SelectExpr

	// AsStruct adapts the expr into a StructExpr.
	//
	// The Kind() must be equal to a StructKind for the conversion to be well-defined.
	AsStruct() StructExpr

	// RenumberIDs performs an in-place update of the expression and all of its descendents numeric ids.
	RenumberIDs(IDGenerator)

	// SetKindCase replaces the contents of the current expression with the contents of the other.
	//
	// The SetKindCase takes ownership of any expression instances references within the input Expr.
	// A shallow copy is made of the Expr value itself, but not a deep one.
	//
	// This method should only be used during AST rewrites using temporary Expr values.
	SetKindCase(Expr)

	// isExpr is a marker interface.
	isExpr()
}

// EntryExprKind represents the possible EntryExpr kinds.
type EntryExprKind int

const (
	// UnspecifiedEntryExprKind indicates that the entry expr is not set.
	UnspecifiedEntryExprKind EntryExprKind = iota

	// MapEntryKind indicates that the entry is a MapEntry type with key and value expressions.
	MapEntryKind

	// StructFieldKind indicates that the entry is a StructField with a field name and initializer
	// expression.
	StructFieldKind
)

// EntryExpr represents the base entry expression in a CEL map or struct literal.
type EntryExpr interface {
	// ID of the entry as it appears in the AST.
	ID() int64

	// Kind of the entry expression node. See EntryExprKind for valid enum values.
	Kind() EntryExprKind

	// AsMapEntry casts the EntryExpr to a MapEntry.
	//
	// The Kind() must be equal to MapEntryKind for the conversion to be well-defined.
	AsMapEntry() MapEntry

	// AsStructField casts the EntryExpr to a StructField
	//
	// The Kind() must be equal to StructFieldKind for the conversion to be well-defined.
	AsStructField() StructField

	// RenumberIDs performs an in-place update of the expression and all of its descendents numeric ids.
	RenumberIDs(IDGenerator)

	isEntryExpr()
}

// IDGenerator produces monotonically increasing ids suitable for tagging expression nodes.
type IDGenerator func() int64

// CallExpr defines an interface for inspecting a function call and its arugments.
type CallExpr interface {
	// FunctionName returns the name of the function.
	FunctionName() string

	// IsMemberFunction returns whether the call has a non-nil target indicating it is a member function
	IsMemberFunction() bool

	// Target returns the target of the expression if one is present.
	Target() Expr

	// Args returns the list of call arguments, excluding the target.
	Args() []Expr

	// marker interface method
	isExpr()
}

// ListExpr defines an interface for inspecting a list literal expression.
type ListExpr interface {
	// Elements returns the list elements as navigable expressions.
	Elements() []Expr

	// OptionalIndicies returns the list of optional indices in the list literal.
	OptionalIndices() []int32

	// Size returns the number of elements in the list.
	Size() int

	// marker interface method
	isExpr()
}

// SelectExpr defines an interface for inspecting a select expression.
type SelectExpr interface {
	// Operand returns the selection operand expression.
	Operand() Expr

	// FieldName returns the field name being selected from the operand.
	FieldName() string

	// IsTestOnly indicates whether the select expression is a presence test generated by a macro.
	IsTestOnly() bool

	// marker interface method
	isExpr()
}

// MapExpr defines an interface for inspecting a map expression.
type MapExpr interface {
	// Entries returns the map key value pairs as EntryExpr values.
	Entries() []EntryExpr

	// Size returns the number of entries in the map.
	Size() int

	// marker interface method
	isExpr()
}

// MapEntry defines an interface for inspecting a map entry.
type MapEntry interface {
	// Key returns the map entry key expression.
	Key() Expr

	// Value returns the map entry value expression.
	Value() Expr

	// IsOptional returns whether the entry is optional.
	IsOptional() bool

	// marker interface method
	isEntryExpr()
}

// StructExpr defines an interfaces for inspecting a struct and its field initializers.
type StructExpr interface {
	// TypeName returns the struct type name.
	TypeName() string

	// Fields returns the set of field initializers in the struct expression as EntryExpr values.
	Fields() []EntryExpr

	// marker interface method
	isExpr()
}

// StructField defines an interface for inspecting a struct field initialization.
type StructField interface {
	// Name returns the name of the field.
	Name() string

	// Value returns the field initialization expression.
	Value() Expr

	// IsOptional returns whether the field is optional.
	IsOptional() bool

	// marker interface method
	isEntryExpr()
}

// ComprehensionExpr defines an interface for inspecting a comprehension expression.
type ComprehensionExpr interface {
	// IterRange returns the iteration range expression.
	IterRange() Expr

	// IterVar returns the iteration variable name.
	IterVar() string

	// AccuVar returns the accumulation variable name.
	AccuVar() string

	// AccuInit returns the accumulation variable initialization expression.
	AccuInit() Expr

	// LoopCondition returns the loop condition expression.
	LoopCondition() Expr

	// LoopStep returns the loop step expression.
	LoopStep() Expr

	// Result returns the comprehension result expression.
	Result() Expr

	// marker interface method
	isExpr()
}

var _ Expr = &expr{}

type expr struct {
	id int64
	exprKindCase
}

type exprKindCase interface {
	Kind() ExprKind

	renumberIDs(IDGenerator)

	isExpr()
}

func (e *expr) ID() int64 {
	if e == nil {
		return 0
	}
	return e.id
}

func (e *expr) Kind() ExprKind {
	if e == nil || e.exprKindCase == nil {
		return UnspecifiedExprKind
	}
	return e.exprKindCase.Kind()
}

func (e *expr) AsCall() CallExpr {
	if e.Kind() != CallKind {
		return nilCall
	}
	return e.exprKindCase.(CallExpr)
}

func (e *expr) AsComprehension() ComprehensionExpr {
	if e.Kind() != ComprehensionKind {
		return nilCompre
	}
	return e.exprKindCase.(ComprehensionExpr)
}

func (e *expr) AsIdent() string {
	if e.Kind() != IdentKind {
		return ""
	}
	return string(e.exprKindCase.(baseIdentExpr))
}

func (e *expr) AsLiteral() ref.Val {
	if e.Kind() != LiteralKind {
		return nil
	}
	return e.exprKindCase.(*baseLiteral).Val
}

func (e *expr) AsList() ListExpr {
	if e.Kind() != ListKind {
		return nilList
	}
	return e.exprKindCase.(ListExpr)
}

func (e *expr) AsMap() MapExpr {
	if e.Kind() != MapKind {
		return nilMap
	}
	return e.exprKindCase.(MapExpr)
}

func (e *expr) AsSelect() SelectExpr {
	if e.Kind() != SelectKind {
		return nilSel
	}
	return e.exprKindCase.(SelectExpr)
}

func (e *expr) AsStruct() StructExpr {
	if e.Kind() != StructKind {
		return nilStruct
	}
	return e.exprKindCase.(StructExpr)
}

func (e *expr) SetKindCase(other Expr) {
	if e == nil {
		return
	}
	if other == nil {
		e.exprKindCase = nil
		return
	}
	switch other.Kind() {
	case CallKind:
		c := other.AsCall()
		e.exprKindCase = &baseCallExpr{
			function: c.FunctionName(),
			target:   c.Target(),
			args:     c.Args(),
		}
	case ComprehensionKind:
		c := other.AsComprehension()
		e.exprKindCase = &baseComprehensionExpr{
			iterRange: c.IterRange(),
			iterVar:   c.IterVar(),
			accuVar:   c.AccuVar(),
			accuInit:  c.AccuInit(),
			loopCond:  c.LoopCondition(),
			loopStep:  c.LoopStep(),
			result:    c.Result(),
		}
	case IdentKind:
		e.exprKindCase = baseIdentExpr(other.AsIdent())
	case ListKind:
		l := other.AsList()
		e.exprKindCase = &baseListExpr{
			elements:   l.Elements(),
			optIndices: l.OptionalIndices(),
		}
	case LiteralKind:
		e.exprKindCase = &baseLiteral{Val: other.AsLiteral()}
	case MapKind:
		e.exprKindCase = &baseMapExpr{
			entries: other.AsMap().Entries(),
		}
	case SelectKind:
		s := other.AsSelect()
		e.exprKindCase = &baseSelectExpr{
			operand:  s.Operand(),
			field:    s.FieldName(),
			testOnly: s.IsTestOnly(),
		}
	case StructKind:
		s := other.AsStruct()
		e.exprKindCase = &baseStructExpr{
			typeName: s.TypeName(),
			fields:   s.Fields(),
		}
	case UnspecifiedExprKind:
		e.exprKindCase = nil
	}
}

func (e *expr) RenumberIDs(idGen IDGenerator) {
	if e.Kind() == UnspecifiedExprKind {
		return
	}
	e.id = idGen()
	e.exprKindCase.renumberIDs(idGen)
}

type baseCallExpr struct {
	function string
	target   Expr
	args     []Expr
	isMember bool
}

func (*baseCallExpr) Kind() ExprKind {
	return CallKind
}

func (e *baseCallExpr) FunctionName() string {
	if e == nil {
		return ""
	}
	return e.function
}

func (e *baseCallExpr) IsMemberFunction() bool {
	if e == nil {
		return false
	}
	return e.isMember
}

func (e *baseCallExpr) Target() Expr {
	if e == nil || !e.IsMemberFunction() {
		return nilExpr
	}
	return e.target
}

func (e *baseCallExpr) Args() []Expr {
	if e == nil {
		return []Expr{}
	}
	return e.args
}

func (e *baseCallExpr) renumberIDs(idGen IDGenerator) {
	if e.IsMemberFunction() {
		e.Target().RenumberIDs(idGen)
	}
	for _, arg := range e.Args() {
		arg.RenumberIDs(idGen)
	}
}

func (*baseCallExpr) isExpr() {}

var _ ComprehensionExpr = &baseComprehensionExpr{}

type baseComprehensionExpr struct {
	iterRange Expr
	iterVar   string
	accuVar   string
	accuInit  Expr
	loopCond  Expr
	loopStep  Expr
	result    Expr
}

func (*baseComprehensionExpr) Kind() ExprKind {
	return ComprehensionKind
}

func (e *baseComprehensionExpr) IterRange() Expr {
	if e == nil {
		return nilExpr
	}
	return e.iterRange
}

func (e *baseComprehensionExpr) IterVar() string {
	return e.iterVar
}

func (e *baseComprehensionExpr) AccuVar() string {
	return e.accuVar
}

func (e *baseComprehensionExpr) AccuInit() Expr {
	if e == nil {
		return nilExpr
	}
	return e.accuInit
}

func (e *baseComprehensionExpr) LoopCondition() Expr {
	if e == nil {
		return nilExpr
	}
	return e.loopCond
}

func (e *baseComprehensionExpr) LoopStep() Expr {
	if e == nil {
		return nilExpr
	}
	return e.loopStep
}

func (e *baseComprehensionExpr) Result() Expr {
	if e == nil {
		return nilExpr
	}
	return e.result
}

func (e *baseComprehensionExpr) renumberIDs(idGen IDGenerator) {
	e.IterRange().RenumberIDs(idGen)
	e.AccuInit().RenumberIDs(idGen)
	e.LoopCondition().RenumberIDs(idGen)
	e.LoopStep().RenumberIDs(idGen)
	e.Result().RenumberIDs(idGen)
}

func (*baseComprehensionExpr) isExpr() {}

var _ exprKindCase = baseIdentExpr("")

type baseIdentExpr string

func (baseIdentExpr) Kind() ExprKind {
	return IdentKind
}

func (e baseIdentExpr) renumberIDs(IDGenerator) {}

func (baseIdentExpr) isExpr() {}

var _ exprKindCase = &baseLiteral{}
var _ ref.Val = &baseLiteral{}

type baseLiteral struct {
	ref.Val
}

func (*baseLiteral) Kind() ExprKind {
	return LiteralKind
}

func (l *baseLiteral) renumberIDs(IDGenerator) {}

func (*baseLiteral) isExpr() {}

var _ ListExpr = &baseListExpr{}

type baseListExpr struct {
	elements   []Expr
	optIndices []int32
}

func (*baseListExpr) Kind() ExprKind {
	return ListKind
}

func (e *baseListExpr) Elements() []Expr {
	if e == nil {
		return []Expr{}
	}
	return e.elements
}

func (e *baseListExpr) OptionalIndices() []int32 {
	if e == nil {
		return []int32{}
	}
	return e.optIndices
}

func (e *baseListExpr) Size() int {
	return len(e.Elements())
}

func (e *baseListExpr) renumberIDs(idGen IDGenerator) {
	for _, elem := range e.Elements() {
		elem.RenumberIDs(idGen)
	}
}

func (*baseListExpr) isExpr() {}

type baseMapExpr struct {
	entries []EntryExpr
}

func (*baseMapExpr) Kind() ExprKind {
	return MapKind
}

func (e *baseMapExpr) Entries() []EntryExpr {
	if e == nil {
		return []EntryExpr{}
	}
	return e.entries
}

func (e *baseMapExpr) Size() int {
	return len(e.Entries())
}

func (e *baseMapExpr) renumberIDs(idGen IDGenerator) {
	for _, entry := range e.Entries() {
		entry.RenumberIDs(idGen)
	}
}

func (*baseMapExpr) isExpr() {}

type baseSelectExpr struct {
	operand  Expr
	field    string
	testOnly bool
}

func (*baseSelectExpr) Kind() ExprKind {
	return SelectKind
}

func (e *baseSelectExpr) Operand() Expr {
	if e == nil || e.operand == nil {
		return nilExpr
	}
	return e.operand
}

func (e *baseSelectExpr) FieldName() string {
	if e == nil {
		return ""
	}
	return e.field
}

func (e *baseSelectExpr) IsTestOnly() bool {
	if e == nil {
		return false
	}
	return e.testOnly
}

func (e *baseSelectExpr) renumberIDs(idGen IDGenerator) {
	e.Operand().RenumberIDs(idGen)
}

func (*baseSelectExpr) isExpr() {}

type baseStructExpr struct {
	typeName string
	fields   []EntryExpr
}

func (*baseStructExpr) Kind() ExprKind {
	return StructKind
}

func (e *baseStructExpr) TypeName() string {
	if e == nil {
		return ""
	}
	return e.typeName
}

func (e *baseStructExpr) Fields() []EntryExpr {
	if e == nil {
		return []EntryExpr{}
	}
	return e.fields
}

func (e *baseStructExpr) renumberIDs(idGen IDGenerator) {
	for _, f := range e.Fields() {
		f.RenumberIDs(idGen)
	}
}

func (*baseStructExpr) isExpr() {}

type entryExprKindCase interface {
	Kind() EntryExprKind

	renumberIDs(IDGenerator)

	isEntryExpr()
}

var _ EntryExpr = &entryExpr{}

type entryExpr struct {
	id int64
	entryExprKindCase
}

func (e *entryExpr) ID() int64 {
	return e.id
}

func (e *entryExpr) AsMapEntry() MapEntry {
	if e.Kind() != MapEntryKind {
		return nilMapEntry
	}
	return e.entryExprKindCase.(MapEntry)
}

func (e *entryExpr) AsStructField() StructField {
	if e.Kind() != StructFieldKind {
		return nilStructField
	}
	return e.entryExprKindCase.(StructField)
}

func (e *entryExpr) RenumberIDs(idGen IDGenerator) {
	e.id = idGen()
	e.entryExprKindCase.renumberIDs(idGen)
}

type baseMapEntry struct {
	key        Expr
	value      Expr
	isOptional bool
}

func (e *baseMapEntry) Kind() EntryExprKind {
	return MapEntryKind
}

func (e *baseMapEntry) Key() Expr {
	if e == nil {
		return nilExpr
	}
	return e.key
}

func (e *baseMapEntry) Value() Expr {
	if e == nil {
		return nilExpr
	}
	return e.value
}

func (e *baseMapEntry) IsOptional() bool {
	if e == nil {
		return false
	}
	return e.isOptional
}

func (e *baseMapEntry) renumberIDs(idGen IDGenerator) {
	e.Key().RenumberIDs(idGen)
	e.Value().RenumberIDs(idGen)
}

func (*baseMapEntry) isEntryExpr() {}

type baseStructField struct {
	field      string
	value      Expr
	isOptional bool
}

func (f *baseStructField) Kind() EntryExprKind {
	return StructFieldKind
}

func (f *baseStructField) Name() string {
	if f == nil {
		return ""
	}
	return f.field
}

func (f *baseStructField) Value() Expr {
	if f == nil {
		return nilExpr
	}
	return f.value
}

func (f *baseStructField) IsOptional() bool {
	if f == nil {
		return false
	}
	return f.isOptional
}

func (f *baseStructField) renumberIDs(idGen IDGenerator) {
	f.Value().RenumberIDs(idGen)
}

func (*baseStructField) isEntryExpr() {}

var (
	nilExpr        *expr                  = nil
	nilCall        *baseCallExpr          = nil
	nilCompre      *baseComprehensionExpr = nil
	nilList        *baseListExpr          = nil
	nilMap         *baseMapExpr           = nil
	nilMapEntry    *baseMapEntry          = nil
	nilSel         *baseSelectExpr        = nil
	nilStruct      *baseStructExpr        = nil
	nilStructField *baseStructField       = nil
)
