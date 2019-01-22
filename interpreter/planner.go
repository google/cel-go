package interpreter

import (
	"fmt"

	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter/functions"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type planner struct {
	disp         Dispatcher
	types        ref.TypeProvider
	pkg          packages.Packager
	identMap     map[string]inst
	refMap       map[int64]*exprpb.Reference
	typeMap      map[int64]*exprpb.Type
	shortcircuit bool
}

func (p *planner) plan(expr *exprpb.Expr) (inst, error) {
	switch expr.ExprKind.(type) {
	case *exprpb.Expr_CallExpr:
		return p.planCall(expr)
	case *exprpb.Expr_IdentExpr:
		return p.planIdent(expr)
	case *exprpb.Expr_SelectExpr:
		return p.planSelect(expr)
	case *exprpb.Expr_ListExpr:
		return p.planCreateList(expr)
	case *exprpb.Expr_StructExpr:
		return p.planCreateStruct(expr)
	case *exprpb.Expr_ComprehensionExpr:
		return p.planComprehension(expr)
	case *exprpb.Expr_ConstExpr:
		return &evalConst{
			id:  expr.Id,
			val: p.constValue(expr.GetConstExpr()),
		}, nil
	}
	return nil, fmt.Errorf("unsupported expr: %v", expr)
}

func (p *planner) planIdent(expr *exprpb.Expr) (inst, error) {
	ident := expr.GetIdentExpr()
	idName := ident.Name
	inst, found := p.identMap[idName]
	if found {
		return inst, nil
	}
	inst = &evalIdent{
		id:   expr.Id,
		name: idName,
	}
	p.identMap[idName] = inst
	return inst, nil
}

func (p *planner) planSelect(expr *exprpb.Expr) (inst, error) {
	sel := expr.GetSelectExpr()
	if sel.TestOnly {
		return p.planTestOnly(expr)
	}

	ref, found := p.refMap[expr.Id]
	if found {
		idName := ref.Name
		if ref.Value != nil {
			return &evalConst{
				id:  expr.Id,
				val: p.constValue(ref.Value),
			}, nil
		}
		inst, found := p.identMap[idName]
		if found {
			return inst, nil
		}
		inst = &evalIdent{
			id:   expr.Id,
			name: idName,
		}
		p.identMap[idName] = inst
		return inst, nil
	}
	op, err := p.plan(sel.GetOperand())
	if err != nil {
		return nil, err
	}
	return &evalSelect{
		id:        expr.Id,
		field:     types.String(sel.Field),
		op:        op,
		resolveID: p.idResolver(sel),
	}, nil
}

func (p *planner) planTestOnly(expr *exprpb.Expr) (inst, error) {
	sel := expr.GetSelectExpr()
	op, err := p.plan(sel.GetOperand())
	if err != nil {
		return nil, err
	}
	return &evalTestOnly{
		id:    expr.Id,
		field: types.String(sel.Field),
		op:    op,
	}, nil
}

func (p *planner) planCall(expr *exprpb.Expr) (inst, error) {
	call := expr.GetCallExpr()
	fnName := call.Function
	fnDef, _ := p.disp.FindOverload(fnName)
	argCount := len(call.GetArgs())
	var offset int
	if call.Target != nil {
		argCount++
		offset++
	}
	args := make([]inst, argCount, argCount)
	if call.Target != nil {
		arg, err := p.plan(call.Target)
		if err != nil {
			return nil, err
		}
		args[0] = arg
	}
	for i, argExpr := range call.GetArgs() {
		arg, err := p.plan(argExpr)
		if err != nil {
			return nil, err
		}
		args[i+offset] = arg
	}
	var oName string
	if oRef, found := p.refMap[expr.Id]; found &&
		len(oRef.GetOverloadId()) == 1 {
		oName = oRef.GetOverloadId()[0]
	}

	switch fnName {
	case operators.LogicalAnd:
		return p.planCallLogicalAnd(expr, args)
	case operators.LogicalOr:
		return p.planCallLogicalOr(expr, args)
	case operators.Conditional:
		return p.planCallConditional(expr, args)
	case operators.Equals:
		return p.planCallEqual(expr, args)
	case operators.NotEquals:
		return p.planCallNotEqual(expr, args)
	}

	switch argCount {
	case 0:
		return p.planCallZero(expr, fnName, oName, fnDef)
	case 1:
		return p.planCallUnary(expr, fnName, oName, fnDef, args)
	case 2:
		return p.planCallBinary(expr, fnName, oName, fnDef, args)
	default:
		return p.planCallVarArgs(expr, fnName, oName, fnDef, args)
	}
}

func (p *planner) planCallZero(expr *exprpb.Expr,
	function string,
	overload string,
	impl *functions.Overload) (inst, error) {
	if impl == nil {
		return nil, fmt.Errorf("no such overload: %s()", function)
	}
	if impl.Function == nil {
		return nil, fmt.Errorf("no such overload: %s()", function)
	}
	return &evalZeroArity{
		impl: impl.Function,
	}, nil
}

func (p *planner) planCallUnary(expr *exprpb.Expr,
	function string,
	overload string,
	impl *functions.Overload,
	args []inst) (inst, error) {
	var fn functions.UnaryOp
	var trait int
	if impl != nil {
		if impl.Unary == nil {
			return nil, fmt.Errorf("no such overload: %s(arg)", function)
		}
		fn = impl.Unary
		trait = impl.OperandTrait
	}
	return &evalUnary{
		id:       expr.Id,
		function: function,
		overload: overload,
		arg:      args[0],
		trait:    trait,
		impl:     fn,
	}, nil
}

func (p *planner) planCallBinary(expr *exprpb.Expr,
	function string,
	overload string,
	impl *functions.Overload,
	args []inst) (inst, error) {
	var fn functions.BinaryOp
	var trait int
	if impl != nil {
		if impl.Binary == nil {
			return nil, fmt.Errorf("no such overload: %s(lhs, rhs)", function)
		}
		fn = impl.Binary
		trait = impl.OperandTrait
	}
	return &evalBinary{
		id:       expr.Id,
		function: function,
		overload: overload,
		lhs:      args[0],
		rhs:      args[1],
		trait:    trait,
		impl:     fn,
	}, nil
}

func (p *planner) planCallVarArgs(expr *exprpb.Expr,
	function string,
	overload string,
	impl *functions.Overload,
	args []inst) (inst, error) {
	var fn functions.FunctionOp
	var trait int
	if impl != nil {
		if impl.Function == nil {
			return nil, fmt.Errorf("no such overload: %s(...)", function)
		}
		fn = impl.Function
		trait = impl.OperandTrait
	}
	return &evalVarArgs{
		id:       expr.Id,
		function: function,
		overload: overload,
		args:     args,
		trait:    trait,
		impl:     fn,
	}, nil
}

func (p *planner) planCallEqual(expr *exprpb.Expr,
	args []inst) (inst, error) {
	return &evalEq{
		id:  expr.Id,
		lhs: args[0],
		rhs: args[1],
	}, nil
}

func (p *planner) planCallNotEqual(expr *exprpb.Expr,
	args []inst) (inst, error) {
	return &evalNe{
		id:  expr.Id,
		lhs: args[0],
		rhs: args[1],
	}, nil
}

func (p *planner) planCallLogicalAnd(expr *exprpb.Expr,
	args []inst) (inst, error) {
	return &evalAnd{
		id:  expr.Id,
		lhs: args[0],
		rhs: args[1],
	}, nil
}

func (p *planner) planCallLogicalOr(expr *exprpb.Expr,
	args []inst) (inst, error) {
	return &evalOr{
		id:  expr.Id,
		lhs: args[0],
		rhs: args[1],
	}, nil
}

func (p *planner) planCallConditional(expr *exprpb.Expr,
	args []inst) (inst, error) {
	return &evalConditional{
		id:     expr.Id,
		expr:   args[0],
		truthy: args[1],
		falsy:  args[2],
	}, nil
}

func (p *planner) planCreateList(expr *exprpb.Expr) (inst, error) {
	list := expr.GetListExpr()
	elems := make([]inst, len(list.GetElements()), len(list.GetElements()))
	for i, elem := range list.GetElements() {
		elemVal, err := p.plan(elem)
		if err != nil {
			return nil, err
		}
		elems[i] = elemVal
	}
	return &evalList{
		id:    expr.Id,
		elems: elems,
	}, nil
}

func (p *planner) planCreateStruct(expr *exprpb.Expr) (inst, error) {
	str := expr.GetStructExpr()
	if len(str.MessageName) != 0 {
		return p.planCreateObj(expr)
	}
	entries := str.GetEntries()
	keys := make([]inst, len(entries))
	vals := make([]inst, len(entries))
	for i, entry := range entries {
		keyVal, err := p.plan(entry.GetMapKey())
		if err != nil {
			return nil, err
		}
		keys[i] = keyVal

		valVal, err := p.plan(entry.GetValue())
		if err != nil {
			return nil, err
		}
		vals[i] = valVal
	}
	return &evalMap{
		id:   expr.Id,
		keys: keys,
		vals: vals,
	}, nil
}

func (p *planner) planCreateObj(expr *exprpb.Expr) (inst, error) {
	obj := expr.GetStructExpr()
	typeName := obj.MessageName
	var defined bool
	for _, qualifiedTypeName := range p.pkg.ResolveCandidateNames(typeName) {
		if _, found := p.types.FindType(qualifiedTypeName); found {
			typeName = qualifiedTypeName
			defined = true
			break
		}
	}
	if !defined {
		panic("unknown type name")
	}
	entries := obj.GetEntries()
	fields := make([]string, len(entries))
	vals := make([]inst, len(entries))
	for i, entry := range entries {
		fields[i] = entry.GetFieldKey()
		val, err := p.plan(entry.GetValue())
		if err != nil {
			return nil, err
		}
		vals[i] = val
	}
	return &evalObj{
		id:       expr.Id,
		typeName: typeName,
		fields:   fields,
		vals:     vals,
	}, nil
}

func (p *planner) planComprehension(expr *exprpb.Expr) (inst, error) {
	fold := expr.GetComprehensionExpr()
	accu, err := p.plan(fold.GetAccuInit())
	if err != nil {
		return nil, err
	}
	iterRange, err := p.plan(fold.GetIterRange())
	if err != nil {
		return nil, err
	}
	cond, err := p.plan(fold.GetLoopCondition())
	if err != nil {
		return nil, err
	}
	step, err := p.plan(fold.GetLoopStep())
	if err != nil {
		return nil, err
	}
	result, err := p.plan(fold.GetResult())
	if err != nil {
		return nil, err
	}
	return &evalFold{
		id:        expr.Id,
		accuVar:   fold.AccuVar,
		accu:      accu,
		iterVar:   fold.IterVar,
		iterRange: iterRange,
		cond:      cond,
		step:      step,
		result:    result,
	}, nil
}

func (p *planner) constValue(c *exprpb.Constant) ref.Value {
	switch c.ConstantKind.(type) {
	case *exprpb.Constant_BoolValue:
		return types.Bool(c.GetBoolValue())
	case *exprpb.Constant_BytesValue:
		return types.Bytes(c.GetBytesValue())
	case *exprpb.Constant_DoubleValue:
		return types.Double(c.GetDoubleValue())
	case *exprpb.Constant_Int64Value:
		return types.Int(c.GetInt64Value())
	case *exprpb.Constant_NullValue:
		return types.Null(c.GetNullValue())
	case *exprpb.Constant_StringValue:
		return types.String(c.GetStringValue())
	case *exprpb.Constant_Uint64Value:
		return types.Uint(c.GetUint64Value())
	}
	return nil
}

func (p *planner) idResolver(sel *exprpb.Expr_Select) func(Activation) (ref.Value, bool) {
	validIdent := true
	ident := sel.Field
	op := sel.Operand
	for validIdent {
		switch op.ExprKind.(type) {
		case *exprpb.Expr_IdentExpr:
			ident = op.GetIdentExpr().Name + "." + ident
			break
		case *exprpb.Expr_SelectExpr:
			nested := op.GetSelectExpr()
			ident = nested.GetField() + "." + ident
			op = nested.Operand
		default:
			validIdent = false
		}
	}
	return func(ctx Activation) (ref.Value, bool) {
		for _, id := range p.pkg.ResolveCandidateNames(ident) {
			if object, found := ctx.ResolveName(id); found {
				return object, found
			}
			if typeIdent, found := p.types.FindIdent(id); found {
				return typeIdent, found
			}
		}
		return nil, false
	}
}

type inst interface {
	eval(ctx Activation) ref.Value
}

type evalIdent struct {
	id   int64
	name string
}

func (id *evalIdent) eval(ctx Activation) ref.Value {
	val, found := ctx.ResolveName(id.name)
	if !found {
		return types.Unknown{id.id}
	}
	return val
}

type evalSelect struct {
	id        int64
	op        inst
	field     types.String
	resolveID func(Activation) (ref.Value, bool)
}

func (sel *evalSelect) eval(ctx Activation) ref.Value {
	obj := sel.op.eval(ctx)
	indexer, ok := obj.(traits.Indexer)
	if !ok {
		resolve, ok := sel.resolveID(ctx)
		if !ok {
			return types.ValOrErr(resolve, "invalid type for field selection.")
		}
		indexer, ok = resolve.(traits.Indexer)
		if !ok {
			return types.ValOrErr(resolve, "invalid type for field selection.")
		}
	}
	return indexer.Get(sel.field)
}

type evalTestOnly struct {
	id    int64
	op    inst
	field types.String
}

func (test *evalTestOnly) eval(ctx Activation) ref.Value {
	obj := test.op.eval(ctx)
	tester, ok := obj.(traits.FieldTester)
	if !ok {
		return types.ValOrErr(obj, "invalid type for field selection.")
	}
	return tester.IsSet(test.field)
}

type evalConst struct {
	id  int64
	val ref.Value
}

func (cons *evalConst) eval(ctx Activation) ref.Value {
	return cons.val
}

type evalOr struct {
	id  int64
	lhs inst
	rhs inst
}

func (or *evalOr) eval(ctx Activation) ref.Value {
	// short-circuit lhs.
	lVal := or.lhs.eval(ctx)
	lBool, lok := lVal.(types.Bool)
	if lok && lBool == types.True {
		return types.True
	}
	// short-circuit on rhs.
	rVal := or.rhs.eval(ctx)
	rBool, rok := rVal.(types.Bool)
	if rok && rBool == types.True {
		return types.True
	}
	// return if both sides are bool false.
	if lok && rok {
		return types.False
	}
	// prefer left unknown to right unknown.
	if types.IsUnknown(lVal) {
		return lVal
	}
	if types.IsUnknown(rVal) {
		return rVal
	}
	// if the left-hand side is non-boolean return it as the error.
	return types.ValOrErr(lVal, "Got '%v', expected argument of type 'bool'", lVal)
}

type evalAnd struct {
	id  int64
	lhs inst
	rhs inst
}

func (or *evalAnd) eval(ctx Activation) ref.Value {
	// short-circuit lhs.
	lVal := or.lhs.eval(ctx)
	lBool, lok := lVal.(types.Bool)
	if lok && lBool == types.False {
		return types.False
	}
	// short-circuit on rhs.
	rVal := or.rhs.eval(ctx)
	rBool, rok := rVal.(types.Bool)
	if rok && rBool == types.False {
		return types.False
	}
	// return if both sides are bool true.
	if lok && rok {
		return types.True
	}
	// prefer left unknown to right unknown.
	if types.IsUnknown(lVal) {
		return lVal
	}
	if types.IsUnknown(rVal) {
		return rVal
	}
	// if the left-hand side is non-boolean return it as the error.
	return types.ValOrErr(lVal, "Got '%v', expected argument of type 'bool'", lVal)
}

type evalConditional struct {
	id     int64
	expr   inst
	truthy inst
	falsy  inst
}

func (cond *evalConditional) eval(ctx Activation) ref.Value {
	condVal := cond.expr.eval(ctx)
	condBool, ok := condVal.(types.Bool)
	if !ok {
		return types.ValOrErr(condVal, "no such overload")
	}
	if condBool {
		return cond.truthy.eval(ctx)
	}
	return cond.falsy.eval(ctx)
}

type evalZeroArity struct {
	impl functions.FunctionOp
}

func (zero *evalZeroArity) eval(ctx Activation) ref.Value {
	return zero.impl()
}

type evalUnary struct {
	id       int64
	function string
	overload string
	arg      inst
	trait    int
	impl     functions.UnaryOp
}

func (un *evalUnary) eval(ctx Activation) ref.Value {
	argVal := un.arg.eval(ctx)
	if un.impl != nil && (un.trait == 0 || argVal.Type().HasTrait(un.trait)) {
		return un.impl(argVal)
	}
	if argVal.Type().HasTrait(traits.ReceiverType) {
		argVal.(traits.Receiver).Receive(un.function, un.overload, []ref.Value{})
	}
	return types.NewErr("no such overload: %s", un.function)
}

type evalBinary struct {
	id       int64
	function string
	overload string
	lhs      inst
	rhs      inst
	trait    int
	impl     functions.BinaryOp
}

func (bin *evalBinary) eval(ctx Activation) ref.Value {
	lVal := bin.lhs.eval(ctx)
	rVal := bin.rhs.eval(ctx)
	if bin.impl != nil && (bin.trait == 0 || lVal.Type().HasTrait(bin.trait)) {
		return bin.impl(lVal, rVal)
	}
	if lVal.Type().HasTrait(traits.ReceiverType) {
		lVal.(traits.Receiver).Receive(bin.function, bin.overload, []ref.Value{lVal, rVal})
	}
	return types.NewErr("no such overload: %s", bin.function)
}

type evalEq struct {
	id  int64
	lhs inst
	rhs inst
}

func (eq *evalEq) eval(ctx Activation) ref.Value {
	lVal := eq.lhs.eval(ctx)
	rVal := eq.rhs.eval(ctx)
	return lVal.Equal(rVal)
}

type evalNe struct {
	id  int64
	lhs inst
	rhs inst
}

func (eq *evalNe) eval(ctx Activation) ref.Value {
	lVal := eq.lhs.eval(ctx)
	rVal := eq.rhs.eval(ctx)
	eqVal := lVal.Equal(rVal)
	eqBool, ok := eqVal.(types.Bool)
	if !ok {
		return types.ValOrErr(eqVal, "no such overload.")
	}
	return !eqBool
}

type evalVarArgs struct {
	id       int64
	function string
	overload string
	args     []inst
	trait    int
	impl     functions.FunctionOp
}

func (fn *evalVarArgs) eval(ctx Activation) ref.Value {
	argVals := make([]ref.Value, len(fn.args), len(fn.args))
	for i, arg := range fn.args {
		argVals[i] = arg.eval(ctx)
	}
	arg0 := argVals[0]
	if fn.impl != nil && (fn.trait == 0 || arg0.Type().HasTrait(fn.trait)) {
		return fn.impl(argVals...)
	}
	if arg0.Type().HasTrait(traits.ReceiverType) {
		return arg0.(traits.Receiver).Receive(fn.function, fn.overload, argVals[1:])
	}
	return types.NewErr("no such overload: %s", fn.function)
}

type evalList struct {
	id    int64
	elems []inst
}

func (l *evalList) eval(ctx Activation) ref.Value {
	elemVals := make([]ref.Value, len(l.elems), len(l.elems))
	for i, elem := range l.elems {
		elemVal := elem.eval(ctx)
		if types.IsUnknownOrError(elemVal) {
			return elemVal
		}
		elemVals[i] = elemVal
	}
	return types.NewDynamicList(elemVals)
}

type evalMap struct {
	id   int64
	keys []inst
	vals []inst
}

func (m *evalMap) eval(ctx Activation) ref.Value {
	entries := make(map[ref.Value]ref.Value)
	for i, key := range m.keys {
		keyVal := key.eval(ctx)
		if types.IsUnknownOrError(keyVal) {
			return keyVal
		}
		valVal := m.vals[i].eval(ctx)
		if types.IsUnknownOrError(valVal) {
			return valVal
		}
		entries[keyVal] = valVal
	}
	return types.NewDynamicMap(entries)
}

type evalObj struct {
	id       int64
	typeName string
	fields   []string
	vals     []inst
	types    ref.TypeProvider
}

func (o *evalObj) eval(ctx Activation) ref.Value {
	fieldVals := make(map[string]ref.Value)
	for i, field := range o.fields {
		val := o.vals[i].eval(ctx)
		if types.IsUnknownOrError(val) {
			return val
		}
		fieldVals[field] = val
	}
	return o.types.NewValue(o.typeName, fieldVals)
}

type evalFold struct {
	id        int64
	accuVar   string
	iterVar   string
	iterRange inst
	accu      inst
	cond      inst
	step      inst
	result    inst
}

func (fold *evalFold) eval(ctx Activation) ref.Value {
	foldRange := fold.iterRange.eval(ctx)
	if !foldRange.Type().HasTrait(traits.IterableType) {
		return types.ValOrErr(foldRange, "got '%T', expected iterable type", foldRange)
	}
	// Configure the fold activation with the accumulator initial value.
	accuCtx := newVarActivation(ctx, fold.accuVar)
	accuCtx.val = fold.accu.eval(ctx)
	iterCtx := newVarActivation(ctx, fold.iterVar)
	it := foldRange.(traits.Iterable).Iterator()
	for it.HasNext() == types.True {
		// Modify the iter var in the fold activation.
		iterCtx.val = it.Next()

		// Evaluate the condition, terminate the loop if false.
		cond := fold.cond.eval(iterCtx)
		condBool, ok := cond.(types.Bool)
		if !ok && !types.IsUnknown(cond) && condBool != types.True {
			break
		}

		// Evalute the evaluation step into accu var.
		accuCtx.val = fold.step.eval(iterCtx)
	}
	// Compute the result.
	return fold.result.eval(accuCtx)
}
