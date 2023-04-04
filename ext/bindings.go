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
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func Bindings() cel.EnvOption {
	return cel.Lib(celBindings{})
}

const (
	celNamespace  = "cel"
	bindMacro     = "bind"
	unusedIterVar = "#unused"
)

type celBindings struct{}

func (celBindings) LibraryName() string {
	return "cel.lib.ext.cel.bindings"
}

func (celBindings) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Macros(
			// cel.bind(var, <init>, <expr>)
			cel.NewReceiverMacro(bindMacro, 3, celBind),
		),
	}
}

func (celBindings) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func celBind(meh cel.MacroExprHelper, target *exprpb.Expr, args []*exprpb.Expr) (*exprpb.Expr, *common.Error) {
	if !macroTargetMatchesNamespace(celNamespace, target) {
		return nil, nil
	}
	varIdent := args[0]
	varName := ""
	switch varIdent.GetExprKind().(type) {
	case *exprpb.Expr_IdentExpr:
		varName = varIdent.GetIdentExpr().GetName()
	default:
		return nil, &common.Error{
			Message:  fmt.Sprintf("cel.bind() variable names must be simple identifers: %v", varIdent),
			Location: meh.OffsetLocation(target.GetId()),
		}
	}
	varInit := args[1]
	resultExpr := args[2]
	return meh.Fold(
		unusedIterVar,
		meh.NewList(),
		varName,
		varInit,
		meh.LiteralBool(false),
		meh.Ident(varName),
		resultExpr,
	), nil
}
