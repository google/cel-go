// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repl

import (
	"fmt"
	"strings"

	antlr "github.com/antlr4-go/antlr/v4"

	"github.com/google/cel-go/repl/parser"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func unparseParamList(ts []*exprpb.Type, sb *strings.Builder) {
	for i, v := range ts {
		unparseHelper(v, sb)
		if i < len(ts)-1 {
			sb.WriteString(", ")
		}
	}
}

func unparseHelper(t *exprpb.Type, sb *strings.Builder) {
	switch k := t.GetTypeKind().(type) {
	case *exprpb.Type_AbstractType_:
		sb.WriteString(k.AbstractType.Name)
		sb.WriteRune('<')
		unparseParamList(k.AbstractType.GetParameterTypes(), sb)
		sb.WriteRune('>')
	case *exprpb.Type_Function:
		sb.WriteString("function<")
		unparseHelper(k.Function.GetResultType(), sb)
		sb.WriteRune('(')
		unparseParamList(k.Function.GetArgTypes(), sb)
		sb.WriteRune(')')
		sb.WriteString(">")
	case *exprpb.Type_ListType_:
		sb.WriteString("list<")
		unparseHelper(k.ListType.GetElemType(), sb)
		sb.WriteString(">")
	case *exprpb.Type_MapType_:
		sb.WriteString("map<")
		unparseHelper(k.MapType.GetKeyType(), sb)
		sb.WriteString(", ")
		unparseHelper(k.MapType.GetValueType(), sb)
		sb.WriteString(">")
	case *exprpb.Type_Type:
		if k.Type.TypeKind == nil {
			sb.WriteString("type")
			break
		}
		sb.WriteString("type<")
		unparseHelper(k.Type, sb)
		sb.WriteString(">")
	case *exprpb.Type_WellKnown:
		switch k.WellKnown {
		case exprpb.Type_ANY:
			sb.WriteString("any")
		case exprpb.Type_DURATION:
			sb.WriteString("duration")
		case exprpb.Type_TIMESTAMP:
			sb.WriteString("timestamp")
		default:
			sb.WriteString("*unspecified wellknown*")
		}
	case *exprpb.Type_Primitive:
		switch k.Primitive {
		case exprpb.Type_BOOL:
			sb.WriteString("bool")
		case exprpb.Type_BYTES:
			sb.WriteString("bytes")
		case exprpb.Type_DOUBLE:
			sb.WriteString("double")
		case exprpb.Type_INT64:
			sb.WriteString("int")
		case exprpb.Type_STRING:
			sb.WriteString("string")
		case exprpb.Type_UINT64:
			sb.WriteString("uint")
		default:
			sb.WriteString("*unspecified primitive*")
		}
	case *exprpb.Type_Wrapper:
		switch k.Wrapper {
		case exprpb.Type_BOOL:
			sb.WriteString("wrapper_bool")
		case exprpb.Type_BYTES:
			sb.WriteString("wrapper_bytes")
		case exprpb.Type_DOUBLE:
			sb.WriteString("wrapper_double")
		case exprpb.Type_INT64:
			sb.WriteString("wrapper_int")
		case exprpb.Type_STRING:
			sb.WriteString("wrapper_string")
		case exprpb.Type_UINT64:
			sb.WriteString("wrapper_uint")
		default:
			sb.WriteString("*unspecified wrapper*")
		}
	case *exprpb.Type_TypeParam:
		sb.WriteRune('~')
		sb.WriteString(k.TypeParam)
	case *exprpb.Type_Dyn:
		sb.WriteString("dyn")
	case *exprpb.Type_Error:
		sb.WriteString("*error*")
	case *exprpb.Type_MessageType:
		sb.WriteString(k.MessageType)
	case *exprpb.Type_Null:
		sb.WriteString("null")
	default:
		sb.WriteString("*unknown type kind*")
	}
}

// UnparseType pretty-prints a type for the REPL.
func UnparseType(t *exprpb.Type) string {
	var sb strings.Builder
	unparseHelper(t, &sb)
	return sb.String()
}

type errorListener struct {
	*antlr.DefaultErrorListener
	errs []error
}

func (l *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
	l.errs = append(l.errs, fmt.Errorf("type parse error: %s", msg))
	l.DefaultErrorListener.SyntaxError(recognizer, offendingSymbol, line, column, msg, e)
}

type typesVisitor struct {
	parser.BaseCommandsVisitor

	errs []error
}

var _ parser.CommandsVisitor = &typesVisitor{}

type typeParams []*exprpb.Type

func (t *typesVisitor) Visit(tree antlr.ParseTree) any {
	switch ctx := tree.(type) {
	case *parser.StartTypeContext:
		return t.VisitStartType(ctx)
	case *parser.TypeContext:
		return t.VisitType(ctx)
	case *parser.TypeIdContext:
		return t.VisitTypeID(ctx)
	case *parser.TypeParamListContext:
		return t.VisitTypeParamList(ctx)
	default:
		t.errs = append(t.errs, fmt.Errorf("unhandled parse node kind"))
		return nil
	}

}

func (t *typesVisitor) VisitStartType(ctx *parser.StartTypeContext) any {
	return t.Visit(ctx.GetT())
}

func (t *typesVisitor) expectUnparameterized(p typeParams, id string) {
	if p != nil {
		t.errs = append(t.errs, fmt.Errorf("unexpected type params for %s", id))
	}
}

func makeAbstract(name string, params typeParams) *exprpb.Type {
	return &exprpb.Type{TypeKind: &exprpb.Type_AbstractType_{
		AbstractType: &exprpb.Type_AbstractType{
			Name:           name,
			ParameterTypes: params}}}
}

func checkWellKnown(name string) *exprpb.Type {
	switch name {
	case "google.protobuf.Timestamp", ".google.protobuf.Timestamp", "timestamp":
		return &exprpb.Type{TypeKind: &exprpb.Type_WellKnown{WellKnown: exprpb.Type_TIMESTAMP}}
	case "google.protobuf.Duration", ".google.protobuf.Duration", "duration":
		return &exprpb.Type{TypeKind: &exprpb.Type_WellKnown{WellKnown: exprpb.Type_DURATION}}
	case "google.protobuf.Any", ".google.protobuf.Any", "any":
		return &exprpb.Type{TypeKind: &exprpb.Type_WellKnown{WellKnown: exprpb.Type_ANY}}
	}
	return nil
}

func (t *typesVisitor) VisitTypeID(ctx *parser.TypeIdContext) any {
	id := ""
	if ctx.GetLeadingDot() != nil {
		id += "."
	}
	tl := ctx.GetId()
	if tl == nil {
		return nil
	}
	id += tl.GetText()
	for _, tok := range ctx.GetQualifiers() {
		id += "." + tok.GetText()
	}
	return id
}

func (t *typesVisitor) VisitTypeParamList(ctx *parser.TypeParamListContext) any {
	var params typeParams
	for _, ty := range ctx.GetTypes() {
		p := t.Visit(ty)
		params = append(params, p.(*exprpb.Type))
	}
	return params
}

func (t *typesVisitor) VisitType(ctx *parser.TypeContext) any {
	emptyType := &exprpb.Type{}

	if ctx.ParamId() != nil {
		param := ctx.ParamId().GetText()
		if !strings.HasPrefix(param, "~") {
			return emptyType
		}
		return &exprpb.Type{TypeKind: &exprpb.Type_TypeParam{TypeParam: strings.TrimPrefix(
			param, "~")}}
	}
	if ctx.TypeId() == nil {
		return emptyType
	}
	r := t.Visit(ctx.TypeId())
	if r == nil {
		return emptyType
	}

	typeID := r.(string)

	paramsCtx := ctx.GetParams()

	var params typeParams
	if paramsCtx != nil {
		r = t.Visit(paramsCtx)
		if r == nil {
			return emptyType
		}
		params = r.(typeParams)
	}

	switch typeID {
	case "int":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: exprpb.Type_INT64}}
	case "uint":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: exprpb.Type_UINT64}}
	case "double":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: exprpb.Type_DOUBLE}}
	case "bytes":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: exprpb.Type_BYTES}}
	case "string":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: exprpb.Type_STRING}}
	case "bool":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Primitive{Primitive: exprpb.Type_BOOL}}
	case "dyn":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Dyn{}}
	case "null", "null_type":
		t.expectUnparameterized(params, typeID)
		return &exprpb.Type{TypeKind: &exprpb.Type_Null{}}
	case "wrapper_int", "google.protobuf.Int64Value", "google.protobuf.Int32Value":
		return &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{Wrapper: exprpb.Type_INT64}}
	case "wrapper_uint", "google.protobuf.Uint64Value", "google.protobuf.Uint32Value":
		return &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{Wrapper: exprpb.Type_UINT64}}
	case "wrapper_double", "google.protobuf.DoubleValue", "google.protobuf.FloatValue":
		return &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{Wrapper: exprpb.Type_DOUBLE}}
	case "wrapper_bool", "google.protobuf.BoolValue":
		return &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{Wrapper: exprpb.Type_BOOL}}
	case "wrapper_string", "google.protobuf.StringValue":
		return &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{Wrapper: exprpb.Type_STRING}}
	case "wrapper_bytes", "google.protobuf.BytesValue":
		return &exprpb.Type{TypeKind: &exprpb.Type_Wrapper{Wrapper: exprpb.Type_BYTES}}
	case "list":
		if params == nil || len(params) != 1 {
			t.errs = append(t.errs, fmt.Errorf("expected exactly one parameter for list"))
			return emptyType
		}
		p := params[0]
		return &exprpb.Type{TypeKind: &exprpb.Type_ListType_{ListType: &exprpb.Type_ListType{ElemType: p}}}
	case "map":
		if params == nil || len(params) != 2 {
			t.errs = append(t.errs, fmt.Errorf("expected exactly two parameters for map"))
			return emptyType
		}
		k, v := params[0], params[1]
		return &exprpb.Type{TypeKind: &exprpb.Type_MapType_{
			MapType: &exprpb.Type_MapType{
				KeyType:   k,
				ValueType: v,
			}}}
	case "type":
		if params == nil || len(params) != 1 {
			t.errs = append(t.errs, fmt.Errorf("expected exactly one parameter for type"))
			return emptyType
		}
		p := params[0]
		return &exprpb.Type{TypeKind: &exprpb.Type_Type{Type: p}}
	case "optional_type":
		if len(params) == 0 {
			params = []*exprpb.Type{{TypeKind: &exprpb.Type_Dyn{}}}
		}
		if len(params) != 1 {
			t.errs = append(t.errs, fmt.Errorf("expected exactly one parameter for optional_type"))
			return emptyType
		}
		return makeAbstract("optional_type", params)
	default:
		wkt := checkWellKnown(typeID)
		if wkt != nil {
			return wkt
		}

		if params != nil {
			return makeAbstract(typeID, params)
		}
		return &exprpb.Type{TypeKind: &exprpb.Type_MessageType{MessageType: typeID}}
	}
}

// ParseType parses a human readable type string into the protobuf representation.
func ParseType(t string) (*exprpb.Type, error) {
	var errListener errorListener
	visitor := &typesVisitor{}
	is := antlr.NewInputStream(t)
	lexer := parser.NewCommandsLexer(is)
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&errListener)
	p := parser.NewCommandsParser(antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel))
	p.RemoveErrorListeners()
	p.AddErrorListener(&errListener)

	var result *exprpb.Type
	s := visitor.Visit(p.StartType())
	if s != nil {
		result = s.(*exprpb.Type)
	}

	errs := append(errListener.errs, visitor.errs...)
	var err error = nil

	if len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i, e := range errs {
			msgs[i] = e.Error()
		}
		err = fmt.Errorf("errors parsing type:\n%s", strings.Join(msgs, "\n"))
	}

	return result, err
}
