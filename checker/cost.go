// Copyright 2022 Google LLC
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
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"math"

	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/parser"
)

// CostEstimator estimates the sizes of variable length input data and the costs of functions.
type CostEstimator interface {
	// EstimateSize returns a SizeEstimate for the given AstNode, or nil if
	// the estimator has no estimate to provide. The size is equivalent to the result of the CEL `size()` function:
	// length of strings and bytes, number of map entries or number of list items.
	// EstimateSize is only called for AstNodes where
	// CEL does not know the size; EstimateSize is not called for values defined inline in CEL where the size
	// is already obvious to CEL.
	EstimateSize(element AstNode) *SizeEstimate
	// EstimateCallCost returns the estimated cost of an invocation, or nil if
	// the estimator has no estimate to provide.
	EstimateCallCost(overloadId string, target *AstNode, args []AstNode) *CostEstimate
}

// AstNode represents an AST node for the purpose of cost estimations.
type AstNode interface {
	// Path returns a field path through the provided type declarations to the type of the AstNode, or nil if the AstNode does not
	// represent type directly reachable from the provided type declarations.
	// The first path element is a variable. All subsequent path elements are one of: field name, '@items', '@keys', '@values'.
	Path() []string
	// Type returns the deduced type of the AstNode.
	Type() *exprpb.Type
	// LiteralSize returns the size of the AstNode if it is a literal defined inline in CEL, or
	// nil if the AstNodes size is not statically known.
	LiteralSize() *uint64
}

type astNode struct {
	path []string
	t    *exprpb.Type
	expr *exprpb.Expr
}

func (e astNode) Path() []string {
	return e.path
}

func (e astNode) Type() *exprpb.Type {
	return e.t
}

func (e astNode) LiteralSize() *uint64 {
	var v uint64
	switch ek := e.expr.ExprKind.(type) {
	case *exprpb.Expr_ConstExpr:
		switch ck := ek.ConstExpr.ConstantKind.(type) {
		case *exprpb.Constant_StringValue:
			v = uint64(len(ck.StringValue))
		case *exprpb.Constant_BytesValue:
			v = uint64(len(ck.BytesValue))
		default:
			return nil
		}
	case *exprpb.Expr_ListExpr:
		v = uint64(len(ek.ListExpr.Elements))
	case *exprpb.Expr_StructExpr:
		if ek.StructExpr.MessageName == "" {
			v = uint64(len(ek.StructExpr.Entries))
		}
	default:
		return nil
	}

	return &v
}

// SizeEstimate represents an estimated size of a variable length string, bytes, map or list.
type SizeEstimate struct {
	Min, Max uint64
}

// Multiply multiplies by another SizeEstimate and returns the sum.
// If multiply would result in an uint64 overflow, the result is Maxuint64.
func (se SizeEstimate) Multiply(sizeEstimate SizeEstimate) SizeEstimate {
	return SizeEstimate{
		multiplyUint64NoOverflow(se.Min, sizeEstimate.Min),
		multiplyUint64NoOverflow(se.Max, sizeEstimate.Max),
	}
}

// MultiplyByCostFactor multiplies a SizeEstimate by a cost factor and returns the CostEstimate with the
// nearest integer of the result, rounded up.
func (se SizeEstimate) MultiplyByCostFactor(costPerUnit float64) CostEstimate {
	return CostEstimate{
		multiplyByCostFactor(se.Min, costPerUnit),
		multiplyByCostFactor(se.Max, costPerUnit),
	}
}

// MultiplyByCost multiplies by the cost and returns the sum.
// If multiply would result in an uint64 overflow, the result is Maxuint64.
func (se SizeEstimate) MultiplyByCost(cost CostEstimate) CostEstimate {
	return CostEstimate{
		multiplyUint64NoOverflow(se.Min, cost.Min),
		multiplyUint64NoOverflow(se.Max, cost.Max),
	}
}

// CostEstimate represents an estimated cost range and provides add and multiply operations
// that do not overflow.
type CostEstimate struct {
	Min, Max uint64
}

// Add adds the costs and returns the sum.
// If add would result in an uint64 overflow for the min or max, the value is set to Maxuint64.
func (ce CostEstimate) Add(cost CostEstimate) CostEstimate {
	return CostEstimate{
		addUint64NoOverflow(ce.Min, cost.Min),
		addUint64NoOverflow(ce.Max, cost.Max),
	}
}

// Multiply multiplies by the cost and returns the sum.
// If multiply would result in an uint64 overflow, the result is Maxuint64.
func (ce CostEstimate) Multiply(cost CostEstimate) CostEstimate {
	return CostEstimate{
		multiplyUint64NoOverflow(ce.Min, cost.Min),
		multiplyUint64NoOverflow(ce.Max, cost.Max),
	}
}

// MultiplyByCostFactor multiplies a CostEstimate by a cost factor and returns the CostEstimate with the
// nearest integer of the result, rounded up.
func (ce CostEstimate) MultiplyByCostFactor(costPerUnit float64) CostEstimate {
	return CostEstimate{
		multiplyByCostFactor(ce.Min, costPerUnit),
		multiplyByCostFactor(ce.Max, costPerUnit),
	}
}

// addUint64NoOverflow adds non-negative ints. If the result is exceeds math.MaxUint64, math.MaxUint64
// is returned.
func addUint64NoOverflow(x, y uint64) uint64 {
	if y > 0 && x > math.MaxUint64-y {
		return math.MaxUint64
	}
	return x + y
}

// multiplyUint64NoOverflow multiplies non-negative ints. If the result is exceeds math.MaxUint64, math.MaxUint64
// is returned.
func multiplyUint64NoOverflow(x, y uint64) uint64 {
	if x > 0 && y > 0 && x > math.MaxUint64/y {
		return math.MaxUint64
	}
	return x * y
}

// multiplyByFactor multiplies an integer by a cost factor float and returns the nearest integer value, rounded up.
func multiplyByCostFactor(x uint64, y float64) uint64 {
	xFloat := float64(x)
	if xFloat > 0 && y > 0 && xFloat > math.MaxUint64/y {
		return math.MaxUint64
	}
	return uint64(math.Ceil(xFloat * y))
}

var (
	identCost  = CostEstimate{Min: 1, Max: 1}
	selectCost = CostEstimate{Min: 1, Max: 1}
	constCost  = CostEstimate{Min: 0, Max: 0}

	createListBaseCost    = CostEstimate{Min: 10, Max: 10}
	createMapBaseCost     = CostEstimate{Min: 30, Max: 30}
	createMessageBaseCost = CostEstimate{Min: 40, Max: 40}
)

type coster struct {
	// exprPath maps from Expr Id to field path.
	exprPath map[int64][]string
	// iterRanges tracks the iterRange of each iterVar.
	iterRanges  iterRangeScopes
	checkedExpr *exprpb.CheckedExpr
	estimator   CostEstimator
}

// Use a stack of iterVar -> iterRange Expr Ids to handle shadowed variable names.
type iterRangeScopes map[string][]int64

func (vs iterRangeScopes) push(varName string, expr *exprpb.Expr) {
	vs[varName] = append(vs[varName], expr.GetId())
}

func (vs iterRangeScopes) pop(varName string) {
	varStack := vs[varName]
	vs[varName] = varStack[:len(varStack)-1]
}

func (vs iterRangeScopes) peek(varName string) (int64, bool) {
	varStack := vs[varName]
	if len(varStack) > 0 {
		return varStack[len(varStack)-1], true
	}
	return 0, false
}

// Cost estimates the cost of the parsed and type checked CEL expression.
func Cost(checker *exprpb.CheckedExpr, estimator CostEstimator) CostEstimate {
	c := coster{
		checkedExpr: checker,
		estimator:   estimator,
		exprPath:    map[int64][]string{},
		iterRanges:  map[string][]int64{},
	}
	return c.cost(checker.GetExpr())
}

func (c *coster) cost(e *exprpb.Expr) CostEstimate {
	if e == nil {
		return CostEstimate{}
	}
	var cost CostEstimate
	switch e.ExprKind.(type) {
	case *exprpb.Expr_ConstExpr:
		cost = constCost
	case *exprpb.Expr_IdentExpr:
		cost = c.costIdent(e)
	case *exprpb.Expr_SelectExpr:
		cost = c.costSelect(e)
	case *exprpb.Expr_CallExpr:
		cost = c.costCall(e)
	case *exprpb.Expr_ListExpr:
		cost = c.costCreateList(e)
	case *exprpb.Expr_StructExpr:
		cost = c.costCreateStruct(e)
	case *exprpb.Expr_ComprehensionExpr:
		cost = c.costComprehension(e)
	default:
		return CostEstimate{}
	}
	return cost
}

func (c *coster) costIdent(e *exprpb.Expr) CostEstimate {
	identExpr := e.GetIdentExpr()

	// build and track the field path
	if iterRange, ok := c.iterRanges.peek(identExpr.GetName()); ok {
		switch c.checkedExpr.TypeMap[iterRange].TypeKind.(type) {
		case *exprpb.Type_ListType_:
			c.addPath(e, append(c.exprPath[iterRange], "@items"))
		case *exprpb.Type_MapType_:
			c.addPath(e, append(c.exprPath[iterRange], "@keys"))
		}
	} else {
		c.addPath(e, []string{identExpr.GetName()})
	}

	return identCost
}

func (c *coster) costSelect(e *exprpb.Expr) CostEstimate {
	sel := e.GetSelectExpr()
	var sum CostEstimate
	if sel.GetTestOnly() {
		return sum
	}
	sum = sum.Add(c.cost(sel.GetOperand()))
	targetType := c.getType(sel.GetOperand())
	switch kindOf(targetType) {
	case kindMap, kindObject, kindTypeParam:
		sum = sum.Add(selectCost)
	}

	// build and track the field path
	c.addPath(e, append(c.getPath(sel.GetOperand()), sel.Field))

	return sum
}

func (c *coster) costCall(e *exprpb.Expr) CostEstimate {
	call := e.GetCallExpr()
	target := call.GetTarget()
	args := call.GetArgs()

	var sum CostEstimate

	argTypes := make([]AstNode, len(args))
	for i, arg := range args {
		// TODO: && || operators short circuit, so min cost should only include 1st arg eval
		// unless exhaustive evaluation is enabled
		// TODO: ternary operator also short circuits, Min cost should be cond + min(a, b) within <cond> ? a : b
		sum = sum.Add(c.cost(arg))
		argTypes[i] = c.newAstNode(arg)
	}

	ref := c.checkedExpr.ReferenceMap[e.GetId()]
	if ref == nil || len(ref.GetOverloadId()) == 0 {
		return CostEstimate{}
	}
	var targetType AstNode
	if target != nil {
		if call.Target != nil {
			sum = sum.Add(c.cost(call.GetTarget()))
			targetType = c.newAstNode(call.GetTarget())
		}
	}
	// Pick a cost estimate range that covers all the overload cost estimation ranges
	fnCost := CostEstimate{Min: uint64(math.MaxUint64), Max: 0}
	for _, overload := range ref.GetOverloadId() {
		overloadCost := c.functionCost(overload, &targetType, argTypes)
		if overloadCost.Max > fnCost.Max {
			fnCost.Max = overloadCost.Max
		}
		if overloadCost.Min < fnCost.Min {
			fnCost.Min = overloadCost.Min
		}
		// build and track the field path for index operations
		switch overload {
		case overloads.IndexList:
			if len(args) > 0 {
				c.addPath(e, append(c.getPath(args[0]), "@items"))
			}
		case overloads.IndexMap:
			if len(args) > 0 {
				c.addPath(e, append(c.getPath(args[0]), "@values"))
			}
		}
	}
	return sum.Add(fnCost)
}

func (c *coster) costCreateList(e *exprpb.Expr) CostEstimate {
	create := e.GetListExpr()
	var sum CostEstimate
	for _, e := range create.GetElements() {
		sum = sum.Add(c.cost(e))
	}
	return sum.Add(createListBaseCost)
}

func (c *coster) costCreateStruct(e *exprpb.Expr) CostEstimate {
	str := e.GetStructExpr()
	if str.MessageName != "" {
		return c.costCreateMessage(e)
	} else {
		return c.costCreateMap(e)
	}
}

func (c *coster) costCreateMap(e *exprpb.Expr) CostEstimate {
	mapVal := e.GetStructExpr()
	var sum CostEstimate
	for _, ent := range mapVal.GetEntries() {
		key := ent.GetMapKey()
		sum = sum.Add(c.cost(key))

		sum = sum.Add(c.cost(ent.GetValue()))
	}
	return sum.Add(createMapBaseCost)
}

func (c *coster) costCreateMessage(e *exprpb.Expr) CostEstimate {
	msgVal := e.GetStructExpr()
	var sum CostEstimate
	for _, ent := range msgVal.GetEntries() {
		sum = sum.Add(c.cost(ent.GetValue()))
	}
	return sum.Add(createMessageBaseCost)
}

func (c *coster) costComprehension(e *exprpb.Expr) CostEstimate {
	comp := e.GetComprehensionExpr()
	var sum CostEstimate
	sum = sum.Add(c.cost(comp.GetIterRange()))
	sum = sum.Add(c.cost(comp.GetAccuInit()))

	// Track the iterRange of each IterVar for field path construction
	c.iterRanges.push(comp.GetIterVar(), comp.GetIterRange())
	loopCost := c.cost(comp.GetLoopCondition())
	stepCost := c.cost(comp.GetLoopStep())
	c.iterRanges.pop(comp.GetIterVar())
	sum = sum.Add(c.cost(comp.Result))
	// TODO: comprehensions short circuit, so even if the min list size > 0, the minimum number of elements evaluated
	// will be 1.
	rangeCnt := c.sizeEstimate(c.newAstNode(comp.GetIterRange()))
	rangeCost := rangeCnt.MultiplyByCost(stepCost.Add(loopCost))
	sum = sum.Add(rangeCost)

	return sum
}

func (c *coster) sizeEstimate(t AstNode) SizeEstimate {
	if l := t.LiteralSize(); l != nil {
		return SizeEstimate{Min: *l, Max: *l}
	}
	if l := c.estimator.EstimateSize(t); l != nil {
		return *l
	}
	return SizeEstimate{Min: 0, Max: math.MaxUint64}
}

func (c *coster) functionCost(overloadId string, target *AstNode, args []AstNode) CostEstimate {
	if est := c.estimator.EstimateCallCost(overloadId, target, args); est != nil {
		return *est
	}
	switch overloadId {
	// O(n) functions
	case overloads.StartsWithString, overloads.EndsWithString, overloads.StringToBytes, overloads.BytesToString:
		if len(args) == 1 {
			return c.sizeEstimate(args[0]).MultiplyByCostFactor(0.1)
		}
	case overloads.InList:
		// If a list is composed entirely of constant values this is O(1), but we don't account for that here.
		// We just assume all list containment checks are O(n).
		if len(args) == 2 {
			return c.sizeEstimate(args[1]).MultiplyByCostFactor(1)
		}
	// O(nm) functions
	case overloads.MatchesString:
		// https://swtch.com/~rsc/regexp/regexp1.html applies to RE2 implementation supported by CEL
		if target != nil && len(args) == 1 {
			strCost := c.sizeEstimate(*target).MultiplyByCostFactor(0.1)
			// We don't know how many expressions are in the regex, just the string length (a huge
			// improvement here would be to somehow get a count the number of expressions in the regex or
			// how many states are in the regex state machine and use that to measure regex cost).
			// For now, we're making a guess that each expression in a regex is typically at least 4 chars
			// in length.
			regexCost := c.sizeEstimate(args[0]).MultiplyByCostFactor(0.25)
			return strCost.Multiply(regexCost)
		}
	case overloads.ContainsString:
		if target != nil && len(args) == 1 {
			strCost := c.sizeEstimate(*target).MultiplyByCostFactor(0.1)
			substrCost := c.sizeEstimate(args[0]).MultiplyByCostFactor(0.1)
			return strCost.Multiply(substrCost)
		}
	}
	// O(1) functions
	// Benchmarks suggest that most of the other operations take +/- 50% of a base cost unit
	// which on an Intel xeon 2.20GHz CPU is 50ns.
	return CostEstimate{Min: 1, Max: 1}
}

func (c *coster) getType(e *exprpb.Expr) *exprpb.Type {
	return c.checkedExpr.TypeMap[e.GetId()]
}

func (c *coster) getPath(e *exprpb.Expr) []string {
	return c.exprPath[e.GetId()]
}

func (c *coster) addPath(e *exprpb.Expr, path []string) {
	c.exprPath[e.GetId()] = path
}

func (c *coster) newAstNode(e *exprpb.Expr) *astNode {
	path := c.getPath(e)
	if len(path) > 0 && path[0] == parser.AccumulatorName {
		// only provide paths to root vars; omit accumulator vars
		path = nil
	}
	return &astNode{path: path, t: c.getType(e), expr: e}
}
