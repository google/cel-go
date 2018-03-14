package interpreter

import (
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/interpreter/types"
	"fmt"
	"reflect"
)

type Dispatcher interface {
	Add(overloads ...*functions.Overload) error
	Dispatch(ctx *CallContext) (interface{}, error)
}

type CallContext struct {
	call       *CallExpr
	args       []interface{}
	activation Activation
	metadata   Metadata
}

func (ctx *CallContext) Function() string {
	return ctx.call.Function
}

func (ctx *CallContext) Overload() (string, bool) {
	return ctx.call.Overload, ctx.call.Overload != ""
}

func (ctx *CallContext) Args() []interface{} {
	return ctx.args
}

func NewDispatcher() *defaultDispatcher {
	return &defaultDispatcher{
		make(map[string]*overload),
		make(map[string]map[int][]*overload)}
}

type OverloadMap map[string]*overload
type OverloadMapByFunctionAndArgCount map[string]map[int][]*overload

type defaultDispatcher struct {
	overloads OverloadMap
	functions OverloadMapByFunctionAndArgCount
}

var _ Dispatcher = &defaultDispatcher{}

type overload struct {
	function   string
	overloadId string
	argCount   int
	argTypes   []types.Type
	impl       functions.OverloadImpl
}

func (o *overload) handlesArgs(args []interface{}) bool {
	for i, arg := range args {
		argType := o.argTypes[i]
		if !argType.IsDyn() {
			if t, found := types.TypeOf(arg); !found {
				return false
			} else if !argType.Equal(t) {
				return false
			}
		}
	}
	return true
}

func (d *defaultDispatcher) Add(overloads ...*functions.Overload) error {
	for _, o := range overloads {
		iface := reflect.TypeOf(o.Signature)
		argCount := iface.NumIn()
		argTypes := make([]types.Type, argCount)
		for i := 0; i < argCount; i++ {
			refType := iface.In(i)
			var argType types.Type = nil
			if refType.Kind() == reflect.Interface {
				argType = types.DynType
			} else {
				refVal := reflect.New(refType).Elem()
				if argTypeVal, found := types.TypeOf(refVal.Interface()); found {
					argType = argTypeVal
				}
			}
			if argType != nil {
				argTypes[i] = argType
			} else {
				return fmt.Errorf("unrecognized type '%T'"+
					" in function signature for overload '%s'",
					refType, o.Name)
			}
		}
		overloadRef := &overload{o.Function,
			o.Name,
			argCount,
			argTypes,
			o.Impl}
		if _, found := d.overloads[o.Name]; found {
			return fmt.Errorf("overload already exists '%s'", o.Name)
		}
		d.overloads[o.Name] = overloadRef
		if byFunction, found := d.functions[o.Function]; !found {
			byFunction = make(map[int][]*overload)
			byFunction[argCount] = []*overload{overloadRef}
			d.functions[o.Function] = byFunction
		} else if byArgCount, found := byFunction[argCount]; !found {
			byFunction[argCount] = []*overload{overloadRef}
		} else {
			byFunction[argCount] = append(byArgCount, overloadRef)
		}
	}
	return nil
}

func (d *defaultDispatcher) Dispatch(ctx *CallContext) (interface{}, error) {
	if overload, err := d.findOverload(ctx); err == nil {
		return overload.impl(ctx.args...)
	} else {
		return nil, err
	}
}

func (d *defaultDispatcher) findOverload(ctx *CallContext) (*overload, error) {
	// TODO: Add location metadata to overloads.
	if overloadId, found := ctx.Overload(); found {
		if match, found := d.overloads[overloadId]; found {
			return match, nil
		} else {
			return nil, fmt.Errorf(
				"unknown overload id '%s' for function '%s'",
				ctx.call.Overload, ctx.Function())
		}
	}
	if byArgCount, found := d.functions[ctx.Function()]; found {
		args := ctx.Args()
		candidates := byArgCount[len(args)]
		var matches []*overload
		for _, candidate := range candidates {
			if candidate.handlesArgs(args) {
				matches = append(matches, candidate)
			}
		}
		if len(matches) == 1 {
			return matches[0], nil
		} else {
			return nil, fmt.Errorf("ambiguous overloads for function '%s'. "+
				"candidates: ['%s']",
				ctx.Function(), matches)
		}
	} else {
		return nil, fmt.Errorf(
			"no matching overload for function '%s'",
			ctx.Function())
	}
}
