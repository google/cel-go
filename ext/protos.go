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

package ext

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Protos returns a cel.EnvOption to configure extended macros and functions for
// proto manipulation.
//
// # Protos.GetExt
//
// Retrieves an extension field from the input proto2 syntax message. If the field
// is not set, the default value for the extension field is returned.
//
//	proto.getExt(<msg>, <fully.qualified.extension.name>) -> <field-type>
//
// Examples:
//
//	proto.hasExt(msg, google.expr.proto2.test.int32_ext) // returns int value
//
// # Protos.HasExt
//
// Determines whether an extension field is set on a proto2 syntax message.
//
//	proto.hasExt(<msg>, <fully.qualified.extension.name>) -> <bool>
//
// Examples:
//
//	proto.hasExt(msg, google.expr.proto2.test.int32_ext) // returns true || false
func Protos() cel.EnvOption {
	return cel.Lib(protoLib{})
}

type protoLib struct{}

func (protoLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Macros(
			// proto.getExt(msg, select_expression)
			cel.NewReceiverMacro("getExt", 2, getProtoExt),
			// proto.hasExt(msg, select_expression)
			cel.NewReceiverMacro("hasExt", 2, hasProtoExt),
		),
	}
}

func (protoLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}

func hasProtoExt(meh cel.MacroExprHelper, target *exprpb.Expr, args []*exprpb.Expr) (*exprpb.Expr, *common.Error) {
	if call, isExt := isExtCall(meh, "hasExt", target, args); !isExt {
		return call, nil
	}
	extensionField, err := getExtFieldName(meh, args[1])
	if err != nil {
		return nil, err
	}
	return meh.PresenceTest(args[0], extensionField), nil
}

func getProtoExt(meh cel.MacroExprHelper, target *exprpb.Expr, args []*exprpb.Expr) (*exprpb.Expr, *common.Error) {
	if call, isExt := isExtCall(meh, "getExt", target, args); !isExt {
		return call, nil
	}
	extFieldName, err := getExtFieldName(meh, args[1])
	if err != nil {
		return nil, err
	}
	return meh.Select(args[0], extFieldName), nil
}

func isExtCall(meh cel.MacroExprHelper, fnName string, target *exprpb.Expr, args []*exprpb.Expr) (*exprpb.Expr, bool) {
	switch target.GetExprKind().(type) {
	case *exprpb.Expr_IdentExpr:
		if target.GetIdentExpr().GetName() != "proto" {
			return meh.ReceiverCall(fnName, target, args...), false
		}
	default:
		return meh.ReceiverCall(fnName, target, args...), false
	}
	return nil, true
}

func getExtFieldName(meh cel.MacroExprHelper, expr *exprpb.Expr) (string, *common.Error) {
	isValid := false
	extensionField := ""
	switch expr.GetExprKind().(type) {
	case *exprpb.Expr_SelectExpr:
		extensionField, isValid = validateIdentifier(expr)
	}
	if !isValid {
		return "", &common.Error{
			Message:  "invalid extension field",
			Location: meh.OffsetLocation(expr.GetId()),
		}
	}
	return extensionField, nil
}

func validateIdentifier(expr *exprpb.Expr) (string, bool) {
	switch expr.GetExprKind().(type) {
	case *exprpb.Expr_IdentExpr:
		return expr.GetIdentExpr().GetName(), true
	case *exprpb.Expr_SelectExpr:
		sel := expr.GetSelectExpr()
		if sel.GetTestOnly() {
			return "", false
		}
		opStr, isIdent := validateIdentifier(sel.GetOperand())
		if !isIdent {
			return "", false
		}
		return opStr + "." + sel.GetField(), true
	default:
		return "", false
	}
}
