// Copyright 2019 Google LLC
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

package cel

import (
	"errors"

	"github.com/google/cel-go/checker"
	"github.com/google/cel-go/checker/decls"
	"github.com/google/cel-go/common"
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/interpreter"
	"github.com/google/cel-go/interpreter/functions"
	"github.com/google/cel-go/parser"

	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Source interface representing a user-provided expression.
type Source interface {
	common.Source
}

// Ast representing the checked or unchecked expression, its source, and related metadata such as
// source position information.
type Ast struct {
	expr    *exprpb.Expr
	info    *exprpb.SourceInfo
	source  Source
	refMap  map[int64]*exprpb.Reference
	typeMap map[int64]*exprpb.Type
}

// Expr returns the proto serializable instance of the parsed/checked expression.
func (ast *Ast) Expr() *exprpb.Expr {
	return ast.expr
}

// IsChecked returns whether the Ast value has been successfully type-checked.
func (ast *Ast) IsChecked() bool {
	return ast.refMap != nil && ast.typeMap != nil
}

// SourceInfo returns character offset and newling position information about expression elements.
func (ast *Ast) SourceInfo() *exprpb.SourceInfo {
	return ast.info
}

// ResultType returns the output type of the expression if the Ast has been type-checked, else
// returns decls.Dyn as the parse step cannot infer the type.
func (ast *Ast) ResultType() *exprpb.Type {
	if !ast.IsChecked() {
		return decls.Dyn
	}
	return ast.typeMap[ast.expr.Id]
}

// Source returns a view of the input used to create the Ast. This source may be complete or
// constructed from the SourceInfo.
func (ast *Ast) Source() Source {
	return ast.source
}

// Env encapsulates the context necessary to perform parsing, type checking, or generation of
// evaluable programs for different expressions.
type Env struct {
	declarations []*exprpb.Decl
	macros       []parser.Macro
	pkg          packages.Packager
	provider     ref.TypeProvider
	adapter      ref.TypeAdapter
	chk          *checker.Env
	// environment options, true by default.
	enableBuiltins                 bool
	enableDynamicAggregateLiterals bool
}

// NewEnv creates an Env instance suitable for parsing and checking expressions against a set of
// user-defined constants, variables, and functions. Macros and the standard built-ins are enabled
// by default.
//
// See the EnvOptions for the options that can be used to configure the environment.
func NewEnv(opts ...EnvOption) (*Env, error) {
	registry := types.NewRegistry()
	return (&Env{
		declarations:                   checker.StandardDeclarations(),
		macros:                         parser.AllMacros,
		pkg:                            packages.DefaultPackage,
		provider:                       registry,
		adapter:                        registry,
		enableBuiltins:                 true,
		enableDynamicAggregateLiterals: true,
	}).configure(opts...)
}

// Extend the current environment with additional options to produce a new Env.
func (e *Env) Extend(opts ...EnvOption) (*Env, error) {
	ext := &Env{}
	*ext = *e
	return ext.configure(opts...)
}

// Check performs type-checking on the input Ast and yields a checked Ast and/or set of Issues.
//
// Checking has failed if the returned Issues value and its Issues.Err() value are non-nil.
// Issues should be inspected if they are non-nil, but may not represent a fatal error.
//
// It is possible to have both non-nil Ast and Issues values returned from this call: however,
// the mere presence of an Ast does not imply that it is valid for use.
func (e *Env) Check(ast *Ast) (*Ast, *Issues) {
	// Note, errors aren't currently possible on the Ast to ParsedExpr conversion.
	pe, _ := AstToParsedExpr(ast)
	res, errs := checker.Check(pe, ast.Source(), e.chk)
	if len(errs.GetErrors()) > 0 {
		return nil, &Issues{errs: errs}
	}
	// Manually create the Ast to ensure that the Ast source information (which may be more
	// detailed than the information provided by Check), is returned to the caller.
	return &Ast{
		source:  ast.Source(),
		expr:    res.GetExpr(),
		info:    res.GetSourceInfo(),
		refMap:  res.GetReferenceMap(),
		typeMap: res.GetTypeMap()}, nil
}

// Parse parses the input expression value `txt` to a Ast and/or a set of Issues.
//
// This form of Parse creates a common.Source value for the input `txt` and forwards to the
// ParseSource method.
func (e *Env) Parse(txt string) (*Ast, *Issues) {
	src := common.NewTextSource(txt)
	return e.ParseSource(src)
}

// ParseSource parses the input source to an Ast and/or set of Issues.
//
// Parsing has failed if the returned Issues value and its Issues.Err() value is non-nil.
// Issues should be inspected if they are non-nil, but may not represent a fatal error.
//
// It is possible to have both non-nil Ast and Issues values returned from this call; however,
// the mere presence of an Ast does not imply that it is valid for use.
func (e *Env) ParseSource(src common.Source) (*Ast, *Issues) {
	res, errs := parser.ParseWithMacros(src, e.macros)
	if len(errs.GetErrors()) > 0 {
		return nil, &Issues{errs: errs}
	}
	// Manually create the Ast to ensure that the text source information is propagated on
	// subsequent calls to Check.
	return &Ast{
		source: Source(src),
		expr:   res.GetExpr(),
		info:   res.GetSourceInfo()}, nil
}

// Program generates an evaluable instance of the Ast within the environment (Env).
func (e *Env) Program(ast *Ast, opts ...ProgramOption) (Program, error) {
	if e.enableBuiltins {
		opts = append(
			[]ProgramOption{Functions(functions.StandardOverloads()...)},
			opts...)
	}
	return newProgram(e, ast, opts...)
}

// TypeAdapter returns the `ref.TypeAdapter` configured for the environment.
func (e *Env) TypeAdapter() ref.TypeAdapter {
	return e.adapter
}

// TypeProvider returns the `ref.TypeProvider` configured for the environment.
func (e *Env) TypeProvider() ref.TypeProvider {
	return e.provider
}

// UnknownActivation returns an interpreter.PartialActivation which marks all variables
// declared in the Env as unknown AttributePattern values.
func (e *Env) UnknownActivation() interpreter.PartialActivation {
	var unknownPatterns []*interpreter.AttributePattern
	for _, d := range e.declarations {
		switch d.GetDeclKind().(type) {
		case *exprpb.Decl_Ident:
			unknownPatterns = append(unknownPatterns,
				interpreter.NewAttributePattern(d.GetName()))
		}
	}
	part, _ := interpreter.NewPartialActivation(
		interpreter.EmptyActivation(),
		unknownPatterns...)
	return part
}

// configure applies a series of EnvOptions to the current environment.
func (e *Env) configure(opts ...EnvOption) (*Env, error) {
	// Customized the environment using the provided EnvOption values. If an error is
	// generated at any step this, will be returned as a nil Env with a non-nil error.
	var err error
	for _, opt := range opts {
		e, err = opt(e)
		if err != nil {
			return nil, err
		}
	}

	// Construct the internal checker env, erroring if there is an issue adding the declarations.
	ce := checker.NewEnv(e.pkg, e.provider)
	ce.EnableDynamicAggregateLiterals(e.enableDynamicAggregateLiterals)
	err = ce.Add(e.declarations...)
	if err != nil {
		return nil, err
	}
	e.chk = ce
	return e, nil
}

// Issues defines methods for inspecting the error details of parse and check calls.
//
// Note: in the future, non-fatal warnings and notices may be inspectable via the Issues struct.
type Issues struct {
	errs *common.Errors
}

// NewIssues returns an Issues struct from a common.Errors object.
func NewIssues(errs *common.Errors) *Issues {
	return &Issues{
		errs: errs,
	}
}

// Err returns an error value if the issues list contains one or more errors.
func (i *Issues) Err() error {
	if len(i.errs.GetErrors()) > 0 {
		return errors.New(i.errs.ToDisplayString())
	}
	return nil
}

// Errors returns the collection of errors encountered in more granular detail.
func (i *Issues) Errors() []common.Error {
	return i.errs.GetErrors()
}

// Append collects the issues from another Issues struct into the current object.
func (i *Issues) Append(other *Issues) {
	i.errs.Append(other.errs.GetErrors())
}

// String converts the issues to a suitable display string.
func (i *Issues) String() string {
	return i.errs.ToDisplayString()
}
