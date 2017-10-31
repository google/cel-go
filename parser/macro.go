package parser

import (
	"fmt"

	"celgo/ast"
	"celgo/common"
	"celgo/operators"
)

type Macros []Macro

type Macro struct {
	name          string
	instanceStyle bool
	args          int
	expander      func(*parser, common.Location, ast.Expression, []ast.Expression) ast.Expression
}

var AllMacros = []Macro{

	// The macro "has(m.f)" which tests the presence of a field, avoiding the need to specify
	// the field as a string.
	{
		name:          operators.Has,
		instanceStyle: false,
		args:          1,
		expander:      makeHas,
	},

	// The macro "range.all(var, predicate)", which is true if for all elements in range the  predicate holds.
	{
		name:          operators.All,
		instanceStyle: true,
		args:          2,
		expander:      makeAll,
	},

	// The macro "range.exists(var, predicate)", which is true if for at least one element in
	// range the predicate holds.
	{
		name:          operators.Exists,
		instanceStyle: true,
		args:          2,
		expander:      makeExists,
	},

	// The macro "range.exists_one(var, predicate)", which is true if for exactly one element
	// in range the predicate holds.
	{
		name:          operators.ExistsOne,
		instanceStyle: true,
		args:          2,
		expander:      makeExistsOne,
	},

	// The macro "range.map(var, function)", applies the function to the vars in the range.
	{
		name:          operators.Map,
		instanceStyle: true,
		args:          2,
		expander:      makeMap,
	},

	// The macro "range.map(var, predicate, function)", applies the function to the vars in
	// the range for which the predicate holds true. The other variables are filtered out.
	{
		name:          operators.Map,
		instanceStyle: true,
		args:          3,
		expander:      makeMap,
	},

	// The macro "range.filter(var, predicate)", filters out the variables for which the
	// predicate is false.
	{
		name:          operators.Filter,
		instanceStyle: true,
		args:          2,
		expander:      makeFilter,
	},
}

var NoMacros = []Macro{}

// Field Presence
// ==============

func makeHas(p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	switch args[0].(type) {
	case *ast.SelectExpression:
		s := args[0].(*ast.SelectExpression)
		return ast.NewSelect(p.id(), loc, s.Target, s.Field, true)
	}

	p.errors.invalidHasArgument(loc)
	return &ast.ErrorExpression{}
}

// Logical Quantifiers
// ===================

const accumulatorName = "__result__"

type quantifierKind int

const (
	quantifierAll quantifierKind = iota
	quantifierExists
	quantifierExistsOne
)

func makeAll(p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	return makeQuantifier(quantifierAll, p, loc, target, args)
}

func makeExists(p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	return makeQuantifier(quantifierExists, p, loc, target, args)
}

func makeExistsOne(p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	return makeQuantifier(quantifierExistsOne, p, loc, target, args)
}

func makeQuantifier(kind quantifierKind, p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	v, found := extractIdent(args[0])
	if !found {
		p.errors.argumentIsNotIdent(args[0].Location())
		return &ast.ErrorExpression{}
	}

	accuIdent := func() ast.Expression {
		return ast.NewIdent(p.id(), loc, accumulatorName)
	}

	var init ast.Expression
	var condition ast.Expression
	var step ast.Expression
	var result ast.Expression
	switch kind {
	case quantifierAll:
		init = ast.NewBoolConstant(p.id(), loc, true)
		condition = accuIdent()
		step = ast.NewCallFunction(p.id(), loc, operators.LogicalAnd, accuIdent(), args[1])
		result = accuIdent()
	case quantifierExists:
		init = ast.NewBoolConstant(p.id(), loc, false)
		condition = ast.NewCallFunction(p.id(), loc, operators.LogicalNot, accuIdent())
		step = ast.NewCallFunction(p.id(), loc, operators.LogicalOr, accuIdent(), args[1])
		result = accuIdent()
	case quantifierExistsOne:
		zeroExpr := ast.NewInt64Constant(p.id(), loc, 0)
		oneExpr := ast.NewInt64Constant(p.id(), loc, 1)
		init = zeroExpr
		condition = ast.NewCallFunction(p.id(), loc, operators.LessEquals, accuIdent(), oneExpr)
		step = ast.NewCallFunction(p.id(), loc, operators.Conditional, args[1],
			ast.NewCallFunction(p.id(), loc, operators.Add, accuIdent(), oneExpr), accuIdent())
		result = ast.NewCallFunction(p.id(), loc, operators.Equals, accuIdent(), oneExpr)
	default:
		panic("unrecognized quantifier")
	}

	return ast.NewComprehension(p.id(), loc, v, target, accumulatorName, init, condition, step, result)
}

// Map
// ===

func makeMap(p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	v, found := extractIdent(args[0])
	if !found {
		p.errors.argumentIsNotIdent(args[0].Location())
		return &ast.ErrorExpression{}
	}

	var fn ast.Expression
	var filter ast.Expression

	if len(args) == 3 {
		filter = args[1]
		fn = args[2]
	} else {
		filter = nil
		fn = args[1]
	}

	accuExpr := ast.NewIdent(p.id(), loc, accumulatorName)
	init := ast.NewCreateList(p.id(), loc)
	condition := ast.NewBoolConstant(p.id(), loc, true)
	// TODO: use compiler internal method for faster, stateful add.
	step := ast.NewCallFunction(p.id(), loc, operators.Add, accuExpr, ast.NewCreateList(p.id(), loc, fn))

	if filter != nil {
		step = ast.NewCallFunction(p.id(), loc, operators.Conditional, filter, step, accuExpr)
	}
	return ast.NewComprehension(p.id(), loc, v, target, accumulatorName, init, condition, step, accuExpr)
}

// Filter
// ======

func makeFilter(p *parser, loc common.Location, target ast.Expression, args []ast.Expression) ast.Expression {
	v, found := extractIdent(args[0])
	if !found {
		p.errors.argumentIsNotIdent(args[0].Location())
		return &ast.ErrorExpression{}
	}

	filter := args[1]
	accuExpr := ast.NewIdent(p.id(), loc, accumulatorName)
	init := ast.NewCreateList(p.id(), loc)
	condition := ast.NewBoolConstant(p.id(), loc, true)
	// TODO: use compiler internal method for faster, stateful add.
	step := ast.NewCallFunction(p.id(), loc, operators.Add, accuExpr, ast.NewCreateList(p.id(), loc, args[0]))
	step = ast.NewCallFunction(p.id(), loc, operators.Conditional, filter, step, accuExpr)
	return ast.NewComprehension(p.id(), loc, v, target, accumulatorName, init, condition, step, accuExpr)
}

func extractIdent(e ast.Expression) (string, bool) {
	switch e.(type) {
	case *ast.IdentExpression:
		return e.(*ast.IdentExpression).Name, true
	}

	return "", false
}

func makeMacroKey(name string, args int, instanceStyle bool) string {
	return fmt.Sprintf("%s:%d:%v", name, args, instanceStyle)
}
