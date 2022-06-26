// Copyright 2018 Google LLC
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

// Package functions defines the standard builtin functions supported by the
// interpreter and as declared within the checker#StandardDeclarations.
package functions

import (
	"context"

	"github.com/google/cel-go/common/types/ref"
)

type Overloader interface {
	GetOperator() string
	GetOperandTrait() int
	GetUnary() ContextUnaryOp
	GetBinary() ContextBinaryOp
	GetFunction() ContextFunctionOp
	IsNonStrict() bool
}

// Overload defines a named overload of a function, indicating an operand trait
// which must be present on the first argument to the overload as well as one
// of either a unary, binary, or function implementation.
//
// The majority of  operators within the expression language are unary or binary
// and the specializations simplify the call contract for implementers of
// types with operator overloads. Any added complexity is assumed to be handled
// by the generic FunctionOp.
type Overload struct {
	// Operator name as written in an expression or defined within
	// operators.go.
	Operator string

	// Operand trait used to dispatch the call. The zero-value indicates a
	// global function overload or that one of the Unary / Binary / Function
	// definitions should be used to execute the call.
	OperandTrait int

	// Unary defines the overload with a UnaryOp implementation. May be nil.
	Unary UnaryOp

	// Binary defines the overload with a BinaryOp implementation. May be nil.
	Binary BinaryOp

	// Function defines the overload with a FunctionOp implementation. May be
	// nil.
	Function FunctionOp

	// NonStrict specifies whether the Overload will tolerate arguments that
	// are types.Err or types.Unknown.
	NonStrict bool
}

func (o *Overload) GetOperator() string  { return o.Operator }
func (o *Overload) GetOperandTrait() int { return o.OperandTrait }
func (o *Overload) GetUnary() ContextUnaryOp {
	if o.Unary != nil {
		return func(ctx context.Context, value ref.Val) ref.Val { return o.Unary(value) }
	}
	return nil
}
func (o *Overload) GetBinary() ContextBinaryOp {
	if o.Binary != nil {
		return func(ctx context.Context, lhs, rhs ref.Val) ref.Val { return o.Binary(lhs, rhs) }
	}
	return nil
}
func (o *Overload) GetFunction() ContextFunctionOp {
	if o.Function != nil {
		return func(ctx context.Context, values ...ref.Val) ref.Val { return o.Function(values...) }
	}
	return nil
}
func (o *Overload) IsNonStrict() bool { return o.NonStrict }

type ContextOverload struct {
	// Operator name as written in an expression or defined within
	// operators.go.
	Operator string

	// Operand trait used to dispatch the call. The zero-value indicates a
	// global function overload or that one of the Unary / Binary / Function
	// definitions should be used to execute the call.
	OperandTrait int

	// Unary defines the overload with a UnaryOp implementation. May be nil.
	Unary ContextUnaryOp

	// Binary defines the overload with a BinaryOp implementation. May be nil.
	Binary ContextBinaryOp

	// Function defines the overload with a FunctionOp implementation. May be
	// nil.
	Function ContextFunctionOp

	// NonStrict specifies whether the Overload will tolerate arguments that
	// are types.Err or types.Unknown.
	NonStrict bool
}

func (o *ContextOverload) GetOperator() string            { return o.Operator }
func (o *ContextOverload) GetOperandTrait() int           { return o.OperandTrait }
func (o *ContextOverload) GetUnary() ContextUnaryOp       { return o.Unary }
func (o *ContextOverload) GetBinary() ContextBinaryOp     { return o.Binary }
func (o *ContextOverload) GetFunction() ContextFunctionOp { return o.Function }
func (o *ContextOverload) IsNonStrict() bool              { return o.NonStrict }

// UnaryOp is a function that takes a single value and produces an output.
type UnaryOp func(value ref.Val) ref.Val
type ContextUnaryOp func(ctx context.Context, value ref.Val) ref.Val

// BinaryOp is a function that takes two values and produces an output.
type BinaryOp func(lhs ref.Val, rhs ref.Val) ref.Val
type ContextBinaryOp func(ctx context.Context, lhs ref.Val, rhs ref.Val) ref.Val

// FunctionOp is a function with accepts zero or more arguments and produces
// an value (as interface{}) or error as a result.
type FunctionOp func(values ...ref.Val) ref.Val
type ContextFunctionOp func(ctx context.Context, values ...ref.Val) ref.Val
