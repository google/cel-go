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

package ext

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/ast"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
)

// Bindings returns a cel.EnvOption to configure support for local variable
// bindings in expressions.
//
// # Cel.Bind
//
// Binds a simple identifier to an initialization expression which may be used
// in a subsequenct result expression. Bindings may also be nested within each
// other.
//
//	cel.bind(<varName>, <initExpr>, <resultExpr>)
//
// Examples:
//
//	cel.bind(a, 'hello',
//	cel.bind(b, 'world', a + b + b + a)) // "helloworldworldhello"
//
//	// Avoid a list allocation within the exists comprehension.
//	cel.bind(valid_values, [a, b, c],
//	[d, e, f].exists(elem, elem in valid_values))
//
// Local bindings are not guaranteed to be evaluated before use.
func Bindings(options ...BindingsOption) cel.EnvOption {
	b := &celBindings{version: math.MaxUint32}
	for _, o := range options {
		b = o(b)
	}
	return cel.Lib(b)
}

const (
	celNamespace  = "cel"
	bindMacro     = "bind"
	blockFunc     = "@block"
	unusedIterVar = "#unused"
)

// BindingsOption declares a functional operator for configuring the Bindings library behavior.
type BindingsOption func(*celBindings) *celBindings

// BindingsVersion sets the version of the bindings library to an explicit version.
func BindingsVersion(version uint32) BindingsOption {
	return func(lib *celBindings) *celBindings {
		lib.version = version
		return lib
	}
}

type celBindings struct {
	version uint32
}

func (*celBindings) LibraryName() string {
	return "cel.lib.ext.cel.bindings"
}

func (lib *celBindings) CompileOptions() []cel.EnvOption {
	opts := []cel.EnvOption{
		cel.Macros(
			// cel.bind(var, <init>, <expr>)
			cel.ReceiverMacro(bindMacro, 3, celBind),
		),
	}
	if lib.version >= 1 {
		// The cel.@block signature takes a list of subexpressions and a typed expression which is
		// used as the output type.
		paramType := cel.TypeParamType("T")
		opts = append(opts,
			cel.Function("cel.@block",
				cel.Overload("cel_block_list",
					[]*cel.Type{cel.ListType(cel.DynType), paramType}, paramType)),
		)
	}
	return opts
}

func (lib *celBindings) ProgramOptions() []cel.ProgramOption {
	if lib.version >= 1 {
		celBlockPlan := func(i interpreter.Interpretable) (interpreter.Interpretable, error) {
			call, ok := i.(interpreter.InterpretableCall)
			if !ok {
				return i, nil
			}
			switch call.Function() {
			case "cel.@block":
				args := call.Args()
				if len(args) != 2 {
					return nil, fmt.Errorf("cel.@block expects two arguments, but got %d", len(args))
				}
				block, ok := args[0].(interpreter.InterpretableConstructor)
				if !ok {
					return nil, errors.New("cel.@block expects a list constructor as the first argument")
				}
				slotExprs := block.InitVals()
				expr := args[1]
				return newBlockScope(slotExprs, expr), nil
			default:
				return i, nil
			}
		}
		return []cel.ProgramOption{cel.CustomDecorator(celBlockPlan)}
	}
	return []cel.ProgramOption{}
}

func celBind(mef cel.MacroExprFactory, target ast.Expr, args []ast.Expr) (ast.Expr, *cel.Error) {
	if !macroTargetMatchesNamespace(celNamespace, target) {
		return nil, nil
	}
	varIdent := args[0]
	varName := ""
	switch varIdent.Kind() {
	case ast.IdentKind:
		varName = varIdent.AsIdent()
	default:
		return nil, mef.NewError(varIdent.ID(), "cel.bind() variable names must be simple identifiers")
	}
	varInit := args[1]
	resultExpr := args[2]
	return mef.NewComprehension(
		mef.NewList(),
		unusedIterVar,
		varName,
		varInit,
		mef.NewLiteral(types.False),
		mef.NewIdent(varName),
		resultExpr,
	), nil
}

func newBlockScope(slotExprs []interpreter.Interpretable, expr interpreter.Interpretable) *blockScope {
	bs := &blockScope{
		slotExprs: slotExprs,
		expr:      expr,
	}
	bs.slotActivationPool = &sync.Pool{
		New: func() any {
			sa := &slotActivation{
				slotExprs: slotExprs,
				slotVals:  make([]ref.Val, len(slotExprs)),
			}
			return sa
		},
	}
	return bs
}

type blockScope struct {
	slotExprs          []interpreter.Interpretable
	expr               interpreter.Interpretable
	slotActivationPool *sync.Pool
}

func (bs *blockScope) ID() int64 {
	return bs.expr.ID()
}

func (bs *blockScope) Eval(activation interpreter.Activation) ref.Val {
	sa := bs.slotActivationPool.Get().(*slotActivation)
	sa.Activation = activation
	defer bs.clearSlots(sa)
	return bs.expr.Eval(sa)
}

func (bs *blockScope) clearSlots(sa *slotActivation) {
	sa.reset()
	bs.slotActivationPool.Put(sa)
}

type slotActivation struct {
	interpreter.Activation
	slotExprs []interpreter.Interpretable
	slotVals  []ref.Val
}

func (sa *slotActivation) ResolveName(name string) (any, bool) {
	if idx, found := strings.CutPrefix(name, indexPrefix); found {
		idx, err := strconv.Atoi(idx)
		if err != nil {
			return nil, false
		}
		v := sa.slotVals[idx]
		if v != nil {
			return v, true
		}
		v = sa.slotExprs[idx].Eval(sa)
		sa.slotVals[idx] = v
		return v, true
	}
	return sa.Activation.ResolveName(name)
}

func (sa *slotActivation) reset() {
	sa.Activation = nil
	clear(sa.slotVals)
}

var (
	indexPrefix = "@index"
)
