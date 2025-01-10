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
	"math"

	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/parser"
)

// WARNING: Any changes to cost calculations in this file require a corresponding change in interpreter/runtimecost.go

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
	EstimateCallCost(function, overloadID string, target *AstNode, args []AstNode) *CallEstimate
}

// CallEstimate includes a CostEstimate for the call, and an optional estimate of the result object size.
// The ResultSize should only be provided if the call results in a map, list, string or bytes.
type CallEstimate struct {
	CostEstimate
	ResultSize *SizeEstimate
}

// AstNode represents an AST node for the purpose of cost estimations.
type AstNode interface {
	// Path returns a field path through the provided type declarations to the type of the AstNode, or nil if the AstNode does not
	// represent type directly reachable from the provided type declarations.
	// The first path element is a variable. All subsequent path elements are one of: field name, '@items', '@keys', '@values'.
	Path() []string
	// Type returns the deduced type of the AstNode.
	Type() *types.Type
	// Expr returns the expression of the AstNode.
	Expr() ast.Expr
	// ComputedSize returns a size estimate of the AstNode derived from information available in the CEL expression.
	// For constants and inline list and map declarations, the exact size is returned. For concatenated list, strings
	// and bytes, the size is derived from the size estimates of the operands. nil is returned if there is no
	// computed size available.
	ComputedSize() *SizeEstimate
}

type astNode struct {
	path        []string
	t           *types.Type
	expr        ast.Expr
	derivedSize *SizeEstimate
}

func (e astNode) Path() []string {
	return e.path
}

func (e astNode) Type() *types.Type {
	return e.t
}

func (e astNode) Expr() ast.Expr {
	return e.expr
}

func (e astNode) ComputedSize() *SizeEstimate {
	return e.derivedSize
}

// SizeEstimate represents an estimated size of a variable length string, bytes, map or list.
type SizeEstimate struct {
	Min, Max uint64
}

// UnknownSizeEstimate returns a size between 0 and max uint
func UnknownSizeEstimate() SizeEstimate {
	return unknownSizeEstimate
}

// FixedSizeEstimate returns a size estimate with a fixed min and max range.
func FixedSizeEstimate(size uint64) SizeEstimate {
	return SizeEstimate{Min: size, Max: size}
}

// Add adds to another SizeEstimate and returns the sum.
// If add would result in an uint64 overflow, the result is math.MaxUint64.
func (se SizeEstimate) Add(sizeEstimate SizeEstimate) SizeEstimate {
	return SizeEstimate{
		addUint64NoOverflow(se.Min, sizeEstimate.Min),
		addUint64NoOverflow(se.Max, sizeEstimate.Max),
	}
}

// Multiply multiplies by another SizeEstimate and returns the product.
// If multiply would result in an uint64 overflow, the result is math.MaxUint64.
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

// MultiplyByCost multiplies by the cost and returns the product.
// If multiply would result in an uint64 overflow, the result is math.MaxUint64.
func (se SizeEstimate) MultiplyByCost(cost CostEstimate) CostEstimate {
	return CostEstimate{
		multiplyUint64NoOverflow(se.Min, cost.Min),
		multiplyUint64NoOverflow(se.Max, cost.Max),
	}
}

// Union returns a SizeEstimate that encompasses both input the SizeEstimate.
func (se SizeEstimate) Union(size SizeEstimate) SizeEstimate {
	result := se
	if size.Min < result.Min {
		result.Min = size.Min
	}
	if size.Max > result.Max {
		result.Max = size.Max
	}
	return result
}

// CostEstimate represents an estimated cost range and provides add and multiply operations
// that do not overflow.
type CostEstimate struct {
	Min, Max uint64
}

// UnknownCostEstimate returns a cost with an unknown impact.
func UnknownCostEstimate() CostEstimate {
	return unknownCostEstimate
}

// FixedCostEstimate returns a cost with a fixed min and max range.
func FixedCostEstimate(cost uint64) CostEstimate {
	return CostEstimate{Min: cost, Max: cost}
}

// Add adds the costs and returns the sum.
// If add would result in an uint64 overflow for the min or max, the value is set to math.MaxUint64.
func (ce CostEstimate) Add(cost CostEstimate) CostEstimate {
	return CostEstimate{
		Min: addUint64NoOverflow(ce.Min, cost.Min),
		Max: addUint64NoOverflow(ce.Max, cost.Max),
	}
}

// Multiply multiplies by the cost and returns the product.
// If multiply would result in an uint64 overflow, the result is math.MaxUint64.
func (ce CostEstimate) Multiply(cost CostEstimate) CostEstimate {
	return CostEstimate{
		Min: multiplyUint64NoOverflow(ce.Min, cost.Min),
		Max: multiplyUint64NoOverflow(ce.Max, cost.Max),
	}
}

// MultiplyByCostFactor multiplies a CostEstimate by a cost factor and returns the CostEstimate with the
// nearest integer of the result, rounded up.
func (ce CostEstimate) MultiplyByCostFactor(costPerUnit float64) CostEstimate {
	return CostEstimate{
		Min: multiplyByCostFactor(ce.Min, costPerUnit),
		Max: multiplyByCostFactor(ce.Max, costPerUnit),
	}
}

// Union returns a CostEstimate that encompasses both input the CostEstimates.
func (ce CostEstimate) Union(size CostEstimate) CostEstimate {
	result := ce
	if size.Min < result.Min {
		result.Min = size.Min
	}
	if size.Max > result.Max {
		result.Max = size.Max
	}
	return result
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
	if y != 0 && x > math.MaxUint64/y {
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
	ceil := math.Ceil(xFloat * y)
	if ceil >= doubleTwoTo64 {
		return math.MaxUint64
	}
	return uint64(ceil)
}

// CostOption configures flags which affect cost computations.
type CostOption func(*coster) error

// PresenceTestHasCost determines whether presence testing has a cost of one or zero.
//
// Defaults to presence test has a cost of one.
func PresenceTestHasCost(hasCost bool) CostOption {
	return func(c *coster) error {
		if hasCost {
			c.presenceTestCost = selectAndIdentCost
			return nil
		}
		c.presenceTestCost = FixedCostEstimate(0)
		return nil
	}
}

// FunctionEstimator provides a CallEstimate given the target and arguments for a specific function, overload pair.
type FunctionEstimator func(estimator CostEstimator, target *AstNode, args []AstNode) *CallEstimate

// OverloadCostEstimate binds a FunctionCoster to a specific function overload ID.
//
// When a OverloadCostEstimate is provided, it will override the cost calculation of the CostEstimator provided to
// the Cost() call.
func OverloadCostEstimate(overloadID string, functionCoster FunctionEstimator) CostOption {
	return func(c *coster) error {
		c.overloadEstimators[overloadID] = functionCoster
		return nil
	}
}

// Cost estimates the cost of the parsed and type checked CEL expression.
func Cost(checked *ast.AST, estimator CostEstimator, opts ...CostOption) (CostEstimate, error) {
	c := &coster{
		checkedAST:         checked,
		estimator:          estimator,
		overloadEstimators: map[string]FunctionEstimator{},
		exprPaths:          map[int64][]string{},
		localVars:          make(scopes),
		iterVars:           make(iterScopes),
		computedSizes:      map[int64]SizeEstimate{},
		computedEntrySizes: map[int64]entrySizeEstimate{},
		presenceTestCost:   FixedCostEstimate(1),
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return CostEstimate{}, err
		}
	}
	return c.cost(checked.Expr()), nil
}

type coster struct {
	// exprPaths maps from Expr Id to field path.
	exprPaths map[int64][]string
	// iterVars tracks the iterRange of each iterVar.
	iterVars iterScopes
	// localVars tracks the local variables possibly assigned during cel.bind()-like calls
	localVars scopes
	// computedSizes tracks the computed sizes of call results.
	computedSizes map[int64]SizeEstimate
	// computedEntrySizes tracks the size of list and map entries
	computedEntrySizes map[int64]entrySizeEstimate

	checkedAST         *ast.AST
	estimator          CostEstimator
	overloadEstimators map[string]FunctionEstimator
	// presenceTestCost will either be a zero or one based on whether has() macros count against cost computations.
	presenceTestCost CostEstimate
}

// entrySizeEstimate captures the container kind and associated key/index and value SizeEstimate values.
//
// An entrySizeEstimate only exists if both the key/index and the value have SizeEstimate values, otherwise
// a nil entrySizeEstimate should be used.
type entrySizeEstimate struct {
	containerKind types.Kind
	key           SizeEstimate
	val           SizeEstimate
}

// container returns the container kind (list or map) of the entry.
func (s *entrySizeEstimate) container() types.Kind {
	if s == nil {
		return types.UnknownKind
	}
	return s.containerKind
}

// keySize returns the SizeEstimate for the key if one exists.
func (s *entrySizeEstimate) keySize() *SizeEstimate {
	if s == nil {
		return nil
	}
	return &s.key
}

// valSize returns the SizeEstimate for the value if one exists.
func (s *entrySizeEstimate) valSize() *SizeEstimate {
	if s == nil {
		return nil
	}
	return &s.val
}

// localVar captures the local variable size and entrySize estimates if they exist for variables
type localVar struct {
	id        int64
	path      []string
	size      *SizeEstimate
	entrySize *entrySizeEstimate
}

// scopes is a stack of variable name to integer id stack to handle scopes created by cel.bind() like macros
type scopes map[string][]*localVar

func (s scopes) push(varName string, expr ast.Expr, size *SizeEstimate, entrySize *entrySizeEstimate) {
	s[varName] = append(s[varName], &localVar{id: expr.ID(), size: size, entrySize: entrySize})
}

func (s scopes) pop(varName string) {
	varStack := s[varName]
	s[varName] = varStack[:len(varStack)-1]
}

func (s scopes) peek(varName string) (*localVar, bool) {
	varStack := s[varName]
	if len(varStack) > 0 {
		return varStack[len(varStack)-1], true
	}
	return nil, false
}

// iterVarKind enumerates the different kinds of iteration variables
type iterVarKind int

const (
	iterVarKindSingle = iota + 1
	iterVarKindKey
	iterVarKindValue
)

// iterVar indicates the iteration variable kind and the associated id
type iterVar struct {
	rangeID int64
	kind    iterVarKind
	size    *SizeEstimate
}

type iterScopes map[string][]*iterVar

func (is iterScopes) pushSingle(varName string, rangeExpr ast.Expr, size *SizeEstimate) {
	is[varName] = append(is[varName], &iterVar{
		rangeID: rangeExpr.ID(),
		kind:    iterVarKindSingle,
		size:    size,
	})
}

func (is iterScopes) pushKey(varName string, rangeExpr ast.Expr, size *SizeEstimate) {
	is[varName] = append(is[varName], &iterVar{
		rangeID: rangeExpr.ID(),
		kind:    iterVarKindKey,
		size:    size,
	})
}

func (is iterScopes) pushValue(varName string, rangeExpr ast.Expr, size *SizeEstimate) {
	is[varName] = append(is[varName], &iterVar{
		rangeID: rangeExpr.ID(),
		kind:    iterVarKindValue,
		size:    size,
	})
}

func (is iterScopes) pop(varName string) {
	varStack := is[varName]
	is[varName] = varStack[:len(varStack)-1]
}

func (is iterScopes) peek(varName string) (*iterVar, bool) {
	varStack := is[varName]
	if len(varStack) > 0 {
		return varStack[len(varStack)-1], true
	}
	return nil, false
}

func (c *coster) pushIterKey(varName string, rangeExpr ast.Expr, size *SizeEstimate) {
	c.iterVars.pushKey(varName, rangeExpr, size)
}

func (c *coster) pushIterValue(varName string, rangeExpr ast.Expr, size *SizeEstimate) {
	c.iterVars.pushValue(varName, rangeExpr, size)
}

func (c *coster) pushIterSingle(varName string, rangeExpr ast.Expr, size *SizeEstimate) {
	c.iterVars.pushSingle(varName, rangeExpr, size)
}

func (c *coster) peekIterVar(varName string) (*iterVar, bool) {
	return c.iterVars.peek(varName)
}

func (c *coster) popIterVar(varName string) {
	c.iterVars.pop(varName)
}

func (c *coster) pushLocalVar(varName string, e ast.Expr) {
	entrySize := c.computeEntrySize(e)
	c.localVars.push(varName, e, c.computeSize(e), entrySize)
}

func (c *coster) peekLocalVar(varName string) (*localVar, bool) {
	return c.localVars.peek(varName)
}

func (c *coster) popLocalVar(varName string) {
	c.localVars.pop(varName)
}

func (c *coster) cost(e ast.Expr) CostEstimate {
	if e == nil {
		return CostEstimate{}
	}
	var cost CostEstimate
	switch e.Kind() {
	case ast.LiteralKind:
		cost = constCost
	case ast.IdentKind:
		cost = c.costIdent(e)
	case ast.SelectKind:
		cost = c.costSelect(e)
	case ast.CallKind:
		cost = c.costCall(e)
	case ast.ListKind:
		cost = c.costCreateList(e)
	case ast.MapKind:
		cost = c.costCreateMap(e)
	case ast.StructKind:
		cost = c.costCreateStruct(e)
	case ast.ComprehensionKind:
		if c.isBind(e) {
			cost = c.costBind(e)
		} else {
			cost = c.costComprehension(e)
		}
	default:
		return CostEstimate{}
	}
	return cost
}

func (c *coster) costIdent(e ast.Expr) CostEstimate {
	identName := e.AsIdent()
	// build and track the field path
	if iterVar, ok := c.peekIterVar(identName); ok {
		rangePath := c.exprPaths[iterVar.rangeID]
		switch c.getTypeByID(iterVar.rangeID).Kind() {
		case types.ListKind:
			path := "@indices"
			if iterVar.kind == iterVarKindValue || iterVar.kind == iterVarKindSingle {
				path = "@items"
			}
			c.addPath(e, append(rangePath, path))
		case types.MapKind:
			path := "@keys"
			if iterVar.kind == iterVarKindValue {
				path = "@values"
			}
			c.addPath(e, append(rangePath, path))
		}
	} else {
		c.addPath(e, []string{identName})
	}

	return selectAndIdentCost
}

func (c *coster) costSelect(e ast.Expr) CostEstimate {
	sel := e.AsSelect()
	var sum CostEstimate
	if sel.IsTestOnly() {
		// recurse, but do not add any cost
		// this is equivalent to how evalTestOnly increments the runtime cost counter
		// but does not add any additional cost for the qualifier, except here we do
		// the reverse (ident adds cost)
		sum = sum.Add(c.presenceTestCost)
		sum = sum.Add(c.cost(sel.Operand()))
		return sum
	}
	sum = sum.Add(c.cost(sel.Operand()))
	targetType := c.getType(sel.Operand())
	switch targetType.Kind() {
	case types.MapKind, types.StructKind, types.TypeParamKind:
		sum = sum.Add(selectAndIdentCost)
	}

	// build and track the field path
	c.addPath(e, append(c.getPath(sel.Operand()), sel.FieldName()))
	return sum
}

func (c *coster) costCall(e ast.Expr) CostEstimate {
	// Dyn is just a way to disable type-checking, so return the cost of 1 with the cost of the argument
	if dynEstimate := c.maybeUnwrapDynCall(e); dynEstimate != nil {
		return *dynEstimate
	}

	// Continue estimating the cost of all other calls.
	call := e.AsCall()
	args := call.Args()
	var sum CostEstimate

	argTypes := make([]AstNode, len(args))
	argCosts := make([]CostEstimate, len(args))
	for i, arg := range args {
		argCosts[i] = c.cost(arg)
		argTypes[i] = c.newAstNode(arg)
	}

	overloadIDs := c.checkedAST.GetOverloadIDs(e.ID())
	if len(overloadIDs) == 0 {
		return CostEstimate{}
	}
	var targetType AstNode
	if call.IsMemberFunction() {
		sum = sum.Add(c.cost(call.Target()))
		targetType = c.newAstNode(call.Target())
	}
	// Pick a cost estimate range that covers all the overload cost estimation ranges
	fnCost := CostEstimate{Min: uint64(math.MaxUint64), Max: 0}
	var resultSize *SizeEstimate
	for _, overload := range overloadIDs {
		overloadCost := c.functionCost(call.FunctionName(), overload, &targetType, argTypes, argCosts)
		fnCost = fnCost.Union(overloadCost.CostEstimate)
		if overloadCost.ResultSize != nil {
			if resultSize == nil {
				resultSize = overloadCost.ResultSize
			} else {
				size := resultSize.Union(*overloadCost.ResultSize)
				resultSize = &size
			}
		}
		// build and track the field path for index operations
		switch overload {
		case overloads.IndexList:
			if len(args) > 0 {
				resultSize = c.computeEntrySize(args[0]).valSize()
				c.addPath(e, append(c.getPath(args[0]), "@items"))
			}
		case overloads.IndexMap:
			if len(args) > 0 {
				resultSize = c.computeEntrySize(args[0]).valSize()
				c.addPath(e, append(c.getPath(args[0]), "@values"))
			}
		}
		if resultSize == nil {
			resultSize = c.computeSize(e)
		}
	}
	c.setSize(e, resultSize)
	return sum.Add(fnCost)
}

func (c *coster) maybeUnwrapDynCall(e ast.Expr) *CostEstimate {
	call := e.AsCall()
	if call.FunctionName() != "dyn" {
		return nil
	}
	arg := call.Args()[0]
	argCost := c.cost(arg)
	c.setSize(e, c.newAstNode(arg).ComputedSize())
	if entrySize := c.computeEntrySize(arg); entrySize != nil {
		c.setEntrySize(e, entrySize)
	}
	callCost := FixedCostEstimate(1).Add(argCost)
	return &callCost
}

func (c *coster) costCreateList(e ast.Expr) CostEstimate {
	create := e.AsList()
	var sum CostEstimate
	itemSize := SizeEstimate{Min: math.MaxUint64, Max: 0}
	if create.Size() == 0 {
		itemSize.Min = 0
	}
	for _, e := range create.Elements() {
		sum = sum.Add(c.cost(e))
		itemNode := c.newAstNode(e)
		is := c.sizeEstimate(itemNode)
		itemSize = itemSize.Union(is)
	}
	c.setEntrySize(e, &entrySizeEstimate{containerKind: types.ListKind, key: FixedSizeEstimate(1), val: itemSize})
	return sum.Add(createListBaseCost)
}

func (c *coster) costCreateMap(e ast.Expr) CostEstimate {
	mapVal := e.AsMap()
	var sum CostEstimate
	keySize := SizeEstimate{Min: math.MaxUint64, Max: 0}
	valSize := SizeEstimate{Min: math.MaxUint64, Max: 0}
	if mapVal.Size() == 0 {
		valSize.Min = 0
		keySize.Min = 0
	}
	for _, ent := range mapVal.Entries() {
		entry := ent.AsMapEntry()
		sum = sum.Add(c.cost(entry.Key()))
		sum = sum.Add(c.cost(entry.Value()))

		// Compute the key size range
		keyNode := c.newAstNode(entry.Key())
		ks := c.sizeEstimate(keyNode)
		keySize = keySize.Union(ks)

		// Compute the value size range
		valueNode := c.newAstNode(entry.Value())
		vs := c.sizeEstimate(valueNode)
		valSize = valSize.Union(vs)
	}
	c.setEntrySize(e, &entrySizeEstimate{containerKind: types.MapKind, key: keySize, val: valSize})
	return sum.Add(createMapBaseCost)
}

func (c *coster) costCreateStruct(e ast.Expr) CostEstimate {
	msgVal := e.AsStruct()
	var sum CostEstimate
	for _, ent := range msgVal.Fields() {
		field := ent.AsStructField()
		sum = sum.Add(c.cost(field.Value()))
	}
	return sum.Add(createMessageBaseCost)
}

func (c *coster) costComprehension(e ast.Expr) CostEstimate {
	comp := e.AsComprehension()
	var sum CostEstimate
	sum = sum.Add(c.cost(comp.IterRange()))
	sum = sum.Add(c.cost(comp.AccuInit()))
	initSize := c.computeSize(comp.AccuInit())
	initEntrySize := c.computeEntrySize(comp.AccuInit())
	c.localVars.push(comp.AccuVar(), comp.AccuInit(), initSize, initEntrySize)
	entrySize := c.computeEntrySize(comp.IterRange())

	// Track the iterRange of each IterVar and AccuVar for field path construction
	if comp.HasIterVar2() {
		c.pushIterKey(comp.IterVar(), comp.IterRange(), entrySize.keySize())
		c.pushIterValue(comp.IterVar2(), comp.IterRange(), entrySize.valSize())
	} else {
		var sz *SizeEstimate
		switch entrySize.container() {
		case types.ListKind:
			sz = entrySize.valSize()
		case types.MapKind:
			sz = entrySize.keySize()
		}
		c.pushIterSingle(comp.IterVar(), comp.IterRange(), sz)
	}

	// Determine the cost for each element in the loop
	loopCost := c.cost(comp.LoopCondition())
	stepCost := c.cost(comp.LoopStep())
	stepSize := c.computeSize(comp.LoopStep())

	// Clear the intermediate variable tracking.
	c.popIterVar(comp.IterVar())
	if comp.HasIterVar2() {
		c.popIterVar(comp.IterVar2())
	}

	// Determine the result cost.
	sum = sum.Add(c.cost(comp.Result()))
	c.localVars.pop(comp.AccuVar())

	// Estimate the cost of the loop.
	rangeCnt := c.sizeEstimate(c.newAstNode(comp.IterRange()))
	rangeCost := rangeCnt.MultiplyByCost(stepCost.Add(loopCost))
	sum = sum.Add(rangeCost)

	switch k := comp.AccuInit().Kind(); k {
	case ast.LiteralKind:
		c.setSize(e, c.computeSize(comp.AccuInit()))
	case ast.ListKind, ast.MapKind:
		kind := types.ListKind
		if k == ast.MapKind {
			kind = types.MapKind
		}
		c.setSize(e, &rangeCnt)
		if entrySize != nil && stepSize != nil {
			c.setEntrySize(e, &entrySizeEstimate{
				containerKind: kind,
				key:           *entrySize.keySize(),
				val:           *stepSize,
			})
		}
	}
	return sum
}

func (c *coster) isBind(e ast.Expr) bool {
	comp := e.AsComprehension()
	iterRange := comp.IterRange()
	loopCond := comp.LoopCondition()
	return iterRange.Kind() == ast.ListKind && iterRange.AsList().Size() == 0 &&
		loopCond.Kind() == ast.LiteralKind && loopCond.AsLiteral() == types.False &&
		comp.AccuVar() != parser.AccumulatorName
}

func (c *coster) costBind(e ast.Expr) CostEstimate {
	comp := e.AsComprehension()
	var sum CostEstimate
	// Binds are lazily initialized, so we retain the cost of an empty iteration range.
	sum = sum.Add(c.cost(comp.IterRange()))
	sum = sum.Add(c.cost(comp.AccuInit()))

	c.pushLocalVar(comp.AccuVar(), comp.AccuInit())
	sum = sum.Add(c.cost(comp.Result()))
	c.popLocalVar(comp.AccuVar())

	// Associate the bind output size with the result size.
	if resultSize := c.computeSize(comp.Result()); resultSize != nil {
		c.setSize(e, resultSize)
	}
	if entrySize := c.computeEntrySize(comp.Result()); entrySize != nil {
		c.setEntrySize(e, entrySize)
	}
	return sum
}

func (c *coster) functionCost(function, overloadID string, target *AstNode, args []AstNode, argCosts []CostEstimate) CallEstimate {
	argCostSum := func() CostEstimate {
		var sum CostEstimate
		for _, a := range argCosts {
			sum = sum.Add(a)
		}
		return sum
	}
	if len(c.overloadEstimators) != 0 {
		if estimator, found := c.overloadEstimators[overloadID]; found {
			if est := estimator(c.estimator, target, args); est != nil {
				callEst := *est
				return CallEstimate{CostEstimate: callEst.Add(argCostSum()), ResultSize: est.ResultSize}
			}
		}
	}
	if est := c.estimator.EstimateCallCost(function, overloadID, target, args); est != nil {
		callEst := *est
		return CallEstimate{CostEstimate: callEst.Add(argCostSum()), ResultSize: est.ResultSize}
	}
	switch overloadID {
	// O(n) functions
	case overloads.ExtFormatString:
		if target != nil {
			// ResultSize not calculated because we can't bound the max size.
			return CallEstimate{
				CostEstimate: c.sizeEstimate(*target).MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum())}
		}
	case overloads.StringToBytes:
		if len(args) == 1 {
			sz := c.sizeEstimate(args[0])
			// ResultSize max is when each char converts to 4 bytes.
			return CallEstimate{
				CostEstimate: sz.MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum()),
				ResultSize:   &SizeEstimate{Min: sz.Min, Max: sz.Max * 4}}
		}
	case overloads.BytesToString:
		if len(args) == 1 {
			sz := c.sizeEstimate(args[0])
			// ResultSize min is when 4 bytes convert to 1 char.
			return CallEstimate{
				CostEstimate: sz.MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum()),
				ResultSize:   &SizeEstimate{Min: sz.Min / 4, Max: sz.Max}}
		}
	case overloads.ExtQuoteString:
		if len(args) == 1 {
			sz := c.sizeEstimate(args[0])
			// ResultSize max is when each char is escaped. 2 quote chars always added.
			return CallEstimate{
				CostEstimate: sz.MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum()),
				ResultSize:   &SizeEstimate{Min: sz.Min + 2, Max: sz.Max*2 + 2}}
		}
	case overloads.StartsWithString, overloads.EndsWithString:
		if len(args) == 1 {
			return CallEstimate{CostEstimate: c.sizeEstimate(args[0]).MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum())}
		}
	case overloads.InList:
		// If a list is composed entirely of constant values this is O(1), but we don't account for that here.
		// We just assume all list containment checks are O(n).
		if len(args) == 2 {
			return CallEstimate{CostEstimate: c.sizeEstimate(args[1]).MultiplyByCostFactor(1).Add(argCostSum())}
		}
	// O(nm) functions
	case overloads.MatchesString:
		// https://swtch.com/~rsc/regexp/regexp1.html applies to RE2 implementation supported by CEL
		if target != nil && len(args) == 1 {
			// Add one to string length for purposes of cost calculation to prevent product of string and regex to be 0
			// in case where string is empty but regex is still expensive.
			strCost := c.sizeEstimate(*target).Add(SizeEstimate{Min: 1, Max: 1}).MultiplyByCostFactor(common.StringTraversalCostFactor)
			// We don't know how many expressions are in the regex, just the string length (a huge
			// improvement here would be to somehow get a count the number of expressions in the regex or
			// how many states are in the regex state machine and use that to measure regex cost).
			// For now, we're making a guess that each expression in a regex is typically at least 4 chars
			// in length.
			regexCost := c.sizeEstimate(args[0]).MultiplyByCostFactor(common.RegexStringLengthCostFactor)
			return CallEstimate{CostEstimate: strCost.Multiply(regexCost).Add(argCostSum())}
		}
	case overloads.ContainsString:
		if target != nil && len(args) == 1 {
			strCost := c.sizeEstimate(*target).MultiplyByCostFactor(common.StringTraversalCostFactor)
			substrCost := c.sizeEstimate(args[0]).MultiplyByCostFactor(common.StringTraversalCostFactor)
			return CallEstimate{CostEstimate: strCost.Multiply(substrCost).Add(argCostSum())}
		}
	case overloads.LogicalOr, overloads.LogicalAnd:
		lhs := argCosts[0]
		rhs := argCosts[1]
		// min cost is min of LHS for short circuited && or ||
		argCost := CostEstimate{Min: lhs.Min, Max: lhs.Add(rhs).Max}
		return CallEstimate{CostEstimate: argCost}
	case overloads.Conditional:
		size := c.sizeEstimate(args[1]).Union(c.sizeEstimate(args[2]))
		conditionalCost := argCosts[0]
		ifTrueCost := argCosts[1]
		ifFalseCost := argCosts[2]
		argCost := conditionalCost.Add(ifTrueCost.Union(ifFalseCost))
		return CallEstimate{CostEstimate: argCost, ResultSize: &size}
	case overloads.AddString, overloads.AddBytes, overloads.AddList:
		if len(args) == 2 {
			lhsSize := c.sizeEstimate(args[0])
			rhsSize := c.sizeEstimate(args[1])
			resultSize := lhsSize.Add(rhsSize)
			switch overloadID {
			case overloads.AddList:
				// list concatenation is O(1), but we handle it here to track size
				return CallEstimate{CostEstimate: FixedCostEstimate(1).Add(argCostSum()), ResultSize: &resultSize}
			default:
				return CallEstimate{CostEstimate: resultSize.MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum()), ResultSize: &resultSize}
			}
		}
	case overloads.LessString, overloads.GreaterString, overloads.LessEqualsString, overloads.GreaterEqualsString,
		overloads.LessBytes, overloads.GreaterBytes, overloads.LessEqualsBytes, overloads.GreaterEqualsBytes,
		overloads.Equals, overloads.NotEquals:
		lhsCost := c.sizeEstimate(args[0])
		rhsCost := c.sizeEstimate(args[1])
		min := uint64(0)
		smallestMax := lhsCost.Max
		if rhsCost.Max < smallestMax {
			smallestMax = rhsCost.Max
		}
		if smallestMax > 0 {
			min = 1
		}
		// equality of 2 scalar values results in a cost of 1
		return CallEstimate{
			CostEstimate: CostEstimate{Min: min, Max: smallestMax}.MultiplyByCostFactor(common.StringTraversalCostFactor).Add(argCostSum()),
		}
	}
	// O(1) functions
	// See CostTracker.costCall for more details about O(1) cost calculations

	// Benchmarks suggest that most of the other operations take +/- 50% of a base cost unit
	// which on an Intel xeon 2.20GHz CPU is 50ns.
	return CallEstimate{CostEstimate: FixedCostEstimate(1).Add(argCostSum())}
}

func (c *coster) getTypeByID(id int64) *types.Type {
	return c.checkedAST.GetType(id)
}

func (c *coster) getType(e ast.Expr) *types.Type {
	return c.getTypeByID(e.ID())
}

func (c *coster) getPathByID(id int64) []string {
	return c.exprPaths[id]
}

func (c *coster) getPath(e ast.Expr) []string {
	return c.getPathByID(e.ID())
}

func (c *coster) addPath(e ast.Expr, path []string) {
	if len(path) == 1 {
		name := path[0]
		if iterVar, found := c.peekIterVar(name); found {
			path = c.getPathByID(iterVar.rangeID)
		}
		if localVar, found := c.peekLocalVar(name); found {
			path = c.getPathByID(localVar.id)
		}
	}
	c.exprPaths[e.ID()] = path
}

func (c *coster) newAstNode(e ast.Expr) *astNode {
	path := c.getPath(e)
	if len(path) > 0 && path[0] == parser.AccumulatorName {
		// only provide paths to root vars; omit accumulator vars
		path = nil
	}
	return &astNode{
		path:        path,
		t:           c.getType(e),
		expr:        e,
		derivedSize: c.computeSize(e)}
}

func (c *coster) setSize(e ast.Expr, size *SizeEstimate) {
	if size == nil {
		return
	}
	// Store the computed size with the expression
	c.computedSizes[e.ID()] = *size
	if e.Kind() != ast.IdentKind {
		return
	}
	// If the expression is an identifier, then determine whether it's a local variable
	// of some kind and store the size with the local variable for future accesses.
	varName := e.AsIdent()
	if localVar, ok := c.peekLocalVar(varName); ok {
		localVar.size = size
	}
}

func (c *coster) setEntrySize(e ast.Expr, size *entrySizeEstimate) {
	if size == nil {
		return
	}
	c.computedEntrySizes[e.ID()] = *size
}

// sizeEstimate computes the estimated size of the given AstNode using a combination of the node's
// computed size or the estimator's assessment of the size. If no size can be established the
// result is an `UnknownSizeEstimate`.
func (c *coster) sizeEstimate(node AstNode) SizeEstimate {
	if l := node.ComputedSize(); l != nil {
		return *l
	}
	if size := c.estimator.EstimateSize(node); size != nil {
		return *size
	}
	return UnknownSizeEstimate()
}

func (c *coster) computeSize(e ast.Expr) *SizeEstimate {
	if size, ok := c.computedSizes[e.ID()]; ok {
		return &size
	}
	if size := computeExprSize(e); size != nil {
		return size
	}
	if size := computeTypeSize(c.getType(e)); size != nil {
		return size
	}
	if e.Kind() == ast.IdentKind {
		varName := e.AsIdent()
		if iterVar, ok := c.peekIterVar(varName); ok {
			return iterVar.size
		}
		if localVar, ok := c.peekLocalVar(varName); ok {
			return localVar.size
		}
	}
	return nil
}

func (c *coster) computeEntrySize(e ast.Expr) *entrySizeEstimate {
	if sz, found := c.computedEntrySizes[e.ID()]; found {
		return &sz
	}
	if e.Kind() == ast.IdentKind {
		varName := e.AsIdent()
		if localVar, ok := c.peekLocalVar(varName); ok {
			return localVar.entrySize
		}
	}
	path, found := c.exprPaths[e.ID()]
	if !found {
		return nil
	}
	pathLen := len(path)
	elemPath := make([]string, pathLen+1)
	copy(elemPath, path)
	entryNode := astNode{path: elemPath}
	switch kind := c.getType(e).Kind(); kind {
	case types.ListKind:
		elemPath[pathLen] = "@items"
		itemSize := c.estimator.EstimateSize(entryNode)
		if itemSize == nil {
			return nil
		}
		return &entrySizeEstimate{
			containerKind: kind,
			key:           FixedSizeEstimate(1),
			val:           *itemSize,
		}
	case types.MapKind:
		elemPath[pathLen] = "@keys"
		keySize := c.estimator.EstimateSize(entryNode)
		if keySize == nil {
			return nil
		}
		elemPath[pathLen] = "@values"
		valSize := c.estimator.EstimateSize(entryNode)
		if valSize == nil {
			return nil
		}
		return &entrySizeEstimate{
			containerKind: kind,
			key:           *keySize,
			val:           *valSize,
		}
	}
	return nil
}

func computeExprSize(expr ast.Expr) *SizeEstimate {
	var v uint64
	switch expr.Kind() {
	case ast.LiteralKind:
		switch ck := expr.AsLiteral().(type) {
		case types.String:
			// converting to runes here is an O(n) operation, but
			// this is consistent with how size is computed at runtime,
			// and how the language definition defines string size
			v = uint64(len([]rune(ck)))
		case types.Bytes:
			v = uint64(len(ck))
		case types.Bool, types.Double, types.Duration,
			types.Int, types.Timestamp, types.Uint,
			types.Null:
			v = uint64(1)
		default:
			return nil
		}
	case ast.ListKind:
		v = uint64(expr.AsList().Size())
	case ast.MapKind:
		v = uint64(expr.AsMap().Size())
	default:
		return nil
	}
	cost := FixedSizeEstimate(v)
	return &cost
}

func computeTypeSize(t *types.Type) *SizeEstimate {
	if isScalar(t) {
		cost := FixedSizeEstimate(1)
		return &cost
	}
	return nil
}

// isScalar returns true if the given type is known to be of a constant size at
// compile time. isScalar will return false for strings (they are variable-width)
// in addition to protobuf.Any and protobuf.Value (their size is not knowable at compile time).
func isScalar(t *types.Type) bool {
	switch t.Kind() {
	case types.BoolKind, types.DoubleKind, types.DurationKind, types.IntKind,
		types.NullTypeKind, types.TimestampKind, types.TypeKind, types.UintKind:
		return true
	case types.OpaqueKind:
		if t.TypeName() == "optional_type" {
			return isScalar(t.Parameters()[0])
		}
	}
	return false
}

var (
	doubleTwoTo64 = math.Ldexp(1.0, 64)

	unknownSizeEstimate = SizeEstimate{Min: 0, Max: math.MaxUint64}
	unknownCostEstimate = unknownSizeEstimate.MultiplyByCostFactor(1)

	selectAndIdentCost = FixedCostEstimate(common.SelectAndIdentCost)
	constCost          = FixedCostEstimate(common.ConstCost)

	createListBaseCost    = FixedCostEstimate(common.ListCreateBaseCost)
	createMapBaseCost     = FixedCostEstimate(common.MapCreateBaseCost)
	createMessageBaseCost = FixedCostEstimate(common.StructCreateBaseCost)
)
