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

package functions

import (
	operatorspb "github.com/google/cel-go/common/operators"
	overloadspb "github.com/google/cel-go/common/overloads"
	typespb "github.com/google/cel-go/common/types"
	refpb "github.com/google/cel-go/common/types/ref"
	traitspb "github.com/google/cel-go/common/types/traits"
)

// StandardOverloads returns the definitions of the built-in overloads.
func StandardOverloads() []*Overload {
	return []*Overload{
		// Logical not (!a)
		{
			Operator:     operatorspb.LogicalNot,
			OperandTrait: traitspb.NegatorType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.(traitspb.Negater).Negate()
			}},
		// Logical and (a && b)
		{
			Operator: operatorspb.LogicalAnd,
			Binary:   logicalAnd},
		// Logical or (a || b)
		{
			Operator: operatorspb.LogicalOr,
			Binary:   logicalOr},
		// Conditional operator (a ? b : c)
		{
			Operator: operatorspb.Conditional,
			Function: conditional},

		// Equality overloads
		{Operator: operatorspb.Equals,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.Equal(rhs)
			}},

		{Operator: operatorspb.NotEquals,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				eq := lhs.Equal(rhs)
				if typespb.IsBool(eq) {
					return !eq.(typespb.Bool)
				}
				return eq
			}},

		// Less than operator
		{Operator: operatorspb.Less,
			OperandTrait: traitspb.ComparerType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				cmp := lhs.(traitspb.Comparer).Compare(rhs)
				if cmp == typespb.IntNegOne {
					return typespb.True
				}
				if cmp == typespb.IntOne || cmp == typespb.IntZero {
					return typespb.False
				}
				return cmp
			}},

		// Less than or equal operator
		{Operator: operatorspb.LessEquals,
			OperandTrait: traitspb.ComparerType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				cmp := lhs.(traitspb.Comparer).Compare(rhs)
				if cmp == typespb.IntNegOne || cmp == typespb.IntZero {
					return typespb.True
				}
				if cmp == typespb.IntOne {
					return typespb.False
				}
				return cmp
			}},

		// Greater than operator
		{Operator: operatorspb.Greater,
			OperandTrait: traitspb.ComparerType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				cmp := lhs.(traitspb.Comparer).Compare(rhs)
				if cmp == typespb.IntOne {
					return typespb.True
				}
				if cmp == typespb.IntNegOne || cmp == typespb.IntZero {
					return typespb.False
				}
				return cmp
			}},

		// Greater than equal operators
		{Operator: operatorspb.GreaterEquals,
			OperandTrait: traitspb.ComparerType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				cmp := lhs.(traitspb.Comparer).Compare(rhs)
				if cmp == typespb.IntOne || cmp == typespb.IntZero {
					return typespb.True
				}
				if cmp == typespb.IntNegOne {
					return typespb.False
				}
				return cmp
			}},

		// TODO: Verify overflow, NaN, underflow cases for numeric values.

		// Add operator
		{Operator: operatorspb.Add,
			OperandTrait: traitspb.AdderType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Adder).Add(rhs)
			}},

		// Subtract operators
		{Operator: operatorspb.Subtract,
			OperandTrait: traitspb.SubtractorType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Subtractor).Subtract(rhs)
			}},

		// Multiply operator
		{Operator: operatorspb.Multiply,
			OperandTrait: traitspb.MultiplierType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Multiplier).Multiply(rhs)
			}},

		// Divide operator
		{Operator: operatorspb.Divide,
			OperandTrait: traitspb.DividerType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Divider).Divide(rhs)
			}},

		// Modulo operator
		{Operator: operatorspb.Modulo,
			OperandTrait: traitspb.ModderType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Modder).Modulo(rhs)
			}},

		// Negate operator
		{Operator: operatorspb.Negate,
			OperandTrait: traitspb.NegatorType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.(traitspb.Negater).Negate()
			}},

		// Index operator
		{Operator: operatorspb.Index,
			OperandTrait: traitspb.IndexerType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Indexer).Get(rhs)
			}},

		// Size function
		{Operator: overloadspb.Size,
			OperandTrait: traitspb.SizerType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.(traitspb.Sizer).Size()
			}},

		// In operator
		{Operator: operatorspb.In,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				if rhs.Type().HasTrait(traitspb.ContainerType) {
					return rhs.(traitspb.Container).Contains(lhs)
				}
				return typespb.NewErr("no such overload")
			}},

		// Matches function
		{Operator: overloadspb.MatchString,
			OperandTrait: traitspb.MatcherType,
			Binary: func(lhs refpb.Value, rhs refpb.Value) refpb.Value {
				return lhs.(traitspb.Matcher).Match(rhs)
			}},

		// Type conversion functions
		// TODO: verify type conversion safety of numeric values.

		// Int conversions.
		{Operator: overloadspb.TypeConvertInt,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.IntType)
			}},

		// Uint conversions.
		{Operator: overloadspb.TypeConvertUint,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.UintType)
			}},

		// Double conversions.
		{Operator: overloadspb.TypeConvertDouble,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.DoubleType)
			}},

		// Bool conversions.
		{Operator: overloadspb.TypeConvertBool,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.BoolType)
			}},

		// Bytes conversions.
		{Operator: overloadspb.TypeConvertBytes,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.BytesType)
			}},

		// String conversions.
		{Operator: overloadspb.TypeConvertString,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.StringType)
			}},

		// Timestamp conversions.
		{Operator: overloadspb.TypeConvertTimestamp,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.TimestampType)
			}},

		// Duration conversions.
		{Operator: overloadspb.TypeConvertDuration,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.DurationType)
			}},

		// Type operations.
		{Operator: overloadspb.TypeConvertType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.ConvertToType(typespb.TypeType)
			}},

		{Operator: overloadspb.Iterator,
			OperandTrait: traitspb.IterableType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.(traitspb.Iterable).Iterator()
			}},

		{Operator: overloadspb.HasNext,
			OperandTrait: traitspb.IteratorType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.(traitspb.Iterator).HasNext()
			}},

		{Operator: overloadspb.Next,
			OperandTrait: traitspb.IteratorType,
			Unary: func(value refpb.Value) refpb.Value {
				return value.(traitspb.Iterator).Next()
			}},
	}

}

func logicalAnd(lhs refpb.Value, rhs refpb.Value) refpb.Value {
	lhsIsBool := typespb.Bool(typespb.IsBool(lhs))
	rhsIsBool := typespb.Bool(typespb.IsBool(rhs))
	// both are boolean use natural logic.
	if lhsIsBool && rhsIsBool {
		return lhs.(typespb.Bool) && rhs.(typespb.Bool)
	}
	// one or the other is boolean and false, return false.
	if lhsIsBool && !lhs.(typespb.Bool) ||
		rhsIsBool && !rhs.(typespb.Bool) {
		return typespb.False
	}

	if typespb.IsUnknown(lhs) {
		return lhs
	}

	if typespb.IsUnknown(rhs) {
		return rhs
	}

	// if the left-hand side is non-boolean return it as the error.
	if !lhsIsBool {
		return typespb.NewErr("Got '%v', expected argument of type 'bool'", lhs)
	}
	return typespb.NewErr("Got '%v', expected argument of type 'bool'", rhs)
}

func logicalOr(lhs refpb.Value, rhs refpb.Value) refpb.Value {
	lhsIsBool := typespb.Bool(typespb.IsBool(lhs.Type()))
	rhsIsBool := typespb.Bool(typespb.IsBool(rhs.Type()))
	// both are boolean, use natural logic.
	if lhsIsBool && rhsIsBool {
		return lhs.(typespb.Bool) || rhs.(typespb.Bool)
	}
	// one or the other is boolean and true, return true
	if lhsIsBool && lhs.(typespb.Bool) ||
		rhsIsBool && rhs.(typespb.Bool) {
		return typespb.True
	}

	if typespb.IsUnknown(lhs) {
		return lhs
	}

	if typespb.IsUnknown(rhs) {
		return rhs
	}

	// if the left-hand side is non-boolean return it as the error.
	if !lhsIsBool {
		return typespb.NewErr("Got '%v', expected argument of type 'bool'", lhs)
	}
	return typespb.NewErr("Got '%v', expected argument of type 'bool'", rhs)
}

func conditional(values ...refpb.Value) refpb.Value {
	if len(values) != 3 {
		return typespb.NewErr("no such overload")
	}
	cond := values[0]
	condType := cond.Type()
	if typespb.IsBool(condType) {
		if cond == typespb.True {
			return values[1]
		}
		return values[2]
	} else if typespb.IsError(condType) || typespb.IsUnknown(condType) {
		return cond
	} else {
		return typespb.NewErr("no such overload")
	}
}
