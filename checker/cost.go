package checker

import (
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"math"

	"google.golang.org/protobuf/proto"

	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common/overloads"
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
	// Type returns the deduced type of the AstNode.
	Type() *exprpb.Type
	// LiteralSize returns the size of the AstNode if it is a literal defined inline in CEL, or
	// nil if the AstNodes size is not statically known.
	LiteralSize() *uint64
}

type expr struct {
	t    *exprpb.Type
	expr *exprpb.Expr
}

func (e expr) Type() *exprpb.Type {
	return e.t
}

func (e expr) LiteralSize() *uint64 {
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
	env       *Env
	checker   *exprpb.CheckedExpr
	estimator CostEstimator
}

// Cost estimates the cost of the parsed CEL expression.
// The parsedExpr is expected to have passed type checking. If it has not been
// type checked, the results are undefined.
func Cost(checker *exprpb.CheckedExpr,
	env *Env, estimator CostEstimator) CostEstimate {
	c := coster{
		env:       env,
		checker:   checker,
		estimator: estimator,
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
	// set types to the exact type references provided
	// other than this set, types are unmodified. This set is guaranteed not to change
	// the type if it is already set.
	// TODO: Isolate where the type is being set such that this is needed.
	if ident := c.env.LookupIdent(identExpr.GetName()); ident != nil {
		c.setType(e, ident.GetIdent().Type)
		return identCost
	}
	return CostEstimate{}
}

func (c *coster) costSelect(e *exprpb.Expr) CostEstimate {
	sel := e.GetSelectExpr()
	var sum CostEstimate

	if sel.TestOnly {
		return sum
	}
	sum = sum.Add(c.cost(sel.Operand))
	targetType := c.getType(sel.Operand)
	switch kindOf(targetType) {
	case kindMap, kindObject, kindTypeParam:
		sum = sum.Add(selectCost)
	}
	return sum
}

func (c *coster) costCall(e *exprpb.Expr) CostEstimate {
	call := e.GetCallExpr()
	target := call.GetTarget()
	args := call.GetArgs()

	var sum CostEstimate

	argTypes := make([]AstNode, len(args))
	for i, arg := range args {
		// TODO: && || short circuit, so min cost should only include 1st arg eval
		// unless exhaustive evaluation is enabled
		sum = sum.Add(c.cost(arg))
		argTypes[i] = expr{t: c.getType(arg), expr: arg}
	}

	ref := c.checker.ReferenceMap[e.Id]
	if ref == nil || len(ref.GetOverloadId()) == 0 {
		return CostEstimate{}
	}
	var targetType AstNode
	if target != nil {
		if call.Target != nil {
			sum = sum.Add(c.cost(call.Target))
			targetType = expr{t: c.getType(call.Target), expr: call.Target}
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
	}
	return sum.Add(fnCost)
}

func (c *coster) costCreateList(e *exprpb.Expr) CostEstimate {
	create := e.GetListExpr()
	var sum CostEstimate
	for _, e := range create.Elements {
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

		sum = sum.Add(c.cost(ent.Value))
	}
	return sum.Add(createMapBaseCost)
}

func (c *coster) costCreateMessage(e *exprpb.Expr) CostEstimate {
	msgVal := e.GetStructExpr()
	var sum CostEstimate
	for _, ent := range msgVal.GetEntries() {
		sum = sum.Add(c.cost(ent.Value))
	}
	return sum.Add(createMessageBaseCost)
}

func (c *coster) costComprehension(e *exprpb.Expr) CostEstimate {
	comp := e.GetComprehensionExpr()
	var sum CostEstimate
	sum = sum.Add(c.cost(comp.IterRange))
	sum = sum.Add(c.cost(comp.AccuInit))
	accuType := c.getType(comp.AccuInit)
	rangeType := c.getType(comp.IterRange)
	var varType *exprpb.Type

	switch kindOf(rangeType) {
	case kindList:
		varType = rangeType.GetListType().ElemType
	case kindMap:
		// Ranges over the keys.
		varType = rangeType.GetMapType().KeyType
	case kindDyn, kindError, kindTypeParam:
		// Set the range iteration variable to type DYN as well.
		varType = decls.Dyn
	default:
		// not a valid range
		varType = decls.Error
	}

	// Create a scope for the comprehension since it has a local accumulation variable.
	// This scope will contain the accumulation variable used to compute the result.
	c.env = c.env.enterScope()
	c.env.Add(decls.NewVar(comp.AccuVar, accuType))
	// Create a block scope for the loop.
	c.env = c.env.enterScope()
	c.env.Add(decls.NewVar(comp.IterVar, varType))
	// Check the variable references in the condition and step.
	loopCost := c.cost(comp.LoopCondition)
	stepCost := c.cost(comp.LoopStep)
	// Exit the loop's block scope before checking the result.
	c.env = c.env.exitScope()
	sum = sum.Add(c.cost(comp.Result))
	// Exit the comprehension scope.
	c.env = c.env.exitScope()

	rangeCnt := c.sizeEstimate(expr{t: c.getType(comp.IterRange), expr: comp.IterRange})
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
	// O(n) functions
	// Benchmarks suggest that most of the other operations take +/- 50% of a base cost unit
	// which on an Intel xeon 2.20GHz CPU is 50ns.
	return CostEstimate{Min: 1, Max: 1}
}

func (c *coster) setType(e *exprpb.Expr, t *exprpb.Type) {
	if old, found := c.checker.TypeMap[e.Id]; found && !proto.Equal(old, t) {
		return
	}
	c.checker.TypeMap[e.Id] = t
}

func (c *coster) getType(e *exprpb.Expr) *exprpb.Type {
	return c.checker.TypeMap[e.Id]
}
