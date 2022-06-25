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
	"github.com/google/cel-go/common/operators"
	"github.com/google/cel-go/common/overloads"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// StandardOverloads returns the definitions of the built-in overloads.
func StandardOverloads() []Overloader {
	return []Overloader{
		// Logical not (!a)
		&Overload{
			Operator:     operators.LogicalNot,
			OperandTrait: traits.NegatorType,
			Unary: func(value ref.Val) ref.Val {
				if !types.IsBool(value) {
					return types.ValOrErr(value, "no such overload")
				}
				return value.(traits.Negater).Negate()
			}},
		// Not strictly false: IsBool(a) ? a : true
		&Overload{
			Operator: operators.NotStrictlyFalse,
			Unary:    notStrictlyFalse},
		// Deprecated: not strictly false, may be overridden in the environment.
		&Overload{
			Operator: operators.OldNotStrictlyFalse,
			Unary:    notStrictlyFalse},

		// Less than operator
		&Overload{Operator: operators.Less,
			OperandTrait: traits.ComparerType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntNegOne {
					return types.True
				}
				if cmp == types.IntOne || cmp == types.IntZero {
					return types.False
				}
				return cmp
			}},

		// Less than or equal operator
		&Overload{Operator: operators.LessEquals,
			OperandTrait: traits.ComparerType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntNegOne || cmp == types.IntZero {
					return types.True
				}
				if cmp == types.IntOne {
					return types.False
				}
				return cmp
			}},

		// Greater than operator
		&Overload{Operator: operators.Greater,
			OperandTrait: traits.ComparerType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntOne {
					return types.True
				}
				if cmp == types.IntNegOne || cmp == types.IntZero {
					return types.False
				}
				return cmp
			}},

		// Greater than equal operators
		&Overload{Operator: operators.GreaterEquals,
			OperandTrait: traits.ComparerType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				cmp := lhs.(traits.Comparer).Compare(rhs)
				if cmp == types.IntOne || cmp == types.IntZero {
					return types.True
				}
				if cmp == types.IntNegOne {
					return types.False
				}
				return cmp
			}},

		// Add operator
		&Overload{Operator: operators.Add,
			OperandTrait: traits.AdderType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Adder).Add(rhs)
			}},

		// Subtract operators
		&Overload{Operator: operators.Subtract,
			OperandTrait: traits.SubtractorType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Subtractor).Subtract(rhs)
			}},

		// Multiply operator
		&Overload{Operator: operators.Multiply,
			OperandTrait: traits.MultiplierType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Multiplier).Multiply(rhs)
			}},

		// Divide operator
		&Overload{Operator: operators.Divide,
			OperandTrait: traits.DividerType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Divider).Divide(rhs)
			}},

		// Modulo operator
		&Overload{Operator: operators.Modulo,
			OperandTrait: traits.ModderType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Modder).Modulo(rhs)
			}},

		// Negate operator
		&Overload{Operator: operators.Negate,
			OperandTrait: traits.NegatorType,
			Unary: func(value ref.Val) ref.Val {
				if types.IsBool(value) {
					return types.ValOrErr(value, "no such overload")
				}
				return value.(traits.Negater).Negate()
			}},

		// Index operator
		&Overload{Operator: operators.Index,
			OperandTrait: traits.IndexerType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Indexer).Get(rhs)
			}},

		// Size function
		&Overload{Operator: overloads.Size,
			OperandTrait: traits.SizerType,
			Unary: func(value ref.Val) ref.Val {
				return value.(traits.Sizer).Size()
			}},

		// In operator
		&Overload{Operator: operators.In, Binary: inAggregate},
		// Deprecated: in operator, may be overridden in the environment.
		&Overload{Operator: operators.OldIn, Binary: inAggregate},

		// Matches function
		&Overload{Operator: overloads.Matches,
			OperandTrait: traits.MatcherType,
			Binary: func(lhs ref.Val, rhs ref.Val) ref.Val {
				return lhs.(traits.Matcher).Match(rhs)
			}},

		// Type conversion functions
		// TODO: verify type conversion safety of numeric values.

		// Int conversions.
		&Overload{Operator: overloads.TypeConvertInt,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.IntType)
			}},

		// Uint conversions.
		&Overload{Operator: overloads.TypeConvertUint,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.UintType)
			}},

		// Double conversions.
		&Overload{Operator: overloads.TypeConvertDouble,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.DoubleType)
			}},

		// Bool conversions.
		&Overload{Operator: overloads.TypeConvertBool,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.BoolType)
			}},

		// Bytes conversions.
		&Overload{Operator: overloads.TypeConvertBytes,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.BytesType)
			}},

		// String conversions.
		&Overload{Operator: overloads.TypeConvertString,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.StringType)
			}},

		// Timestamp conversions.
		&Overload{Operator: overloads.TypeConvertTimestamp,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.TimestampType)
			}},

		// Duration conversions.
		&Overload{Operator: overloads.TypeConvertDuration,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.DurationType)
			}},

		// Type operations.
		&Overload{Operator: overloads.TypeConvertType,
			Unary: func(value ref.Val) ref.Val {
				return value.ConvertToType(types.TypeType)
			}},

		// Dyn conversion (identity function).
		&Overload{Operator: overloads.TypeConvertDyn,
			Unary: func(value ref.Val) ref.Val {
				return value
			}},

		&Overload{Operator: overloads.Iterator,
			OperandTrait: traits.IterableType,
			Unary: func(value ref.Val) ref.Val {
				return value.(traits.Iterable).Iterator()
			}},

		&Overload{Operator: overloads.HasNext,
			OperandTrait: traits.IteratorType,
			Unary: func(value ref.Val) ref.Val {
				return value.(traits.Iterator).HasNext()
			}},

		&Overload{Operator: overloads.Next,
			OperandTrait: traits.IteratorType,
			Unary: func(value ref.Val) ref.Val {
				return value.(traits.Iterator).Next()
			}},
	}

}

func notStrictlyFalse(value ref.Val) ref.Val {
	if types.IsBool(value) {
		return value
	}
	return types.True
}

func inAggregate(lhs ref.Val, rhs ref.Val) ref.Val {
	if rhs.Type().HasTrait(traits.ContainerType) {
		return rhs.(traits.Container).Contains(lhs)
	}
	return types.ValOrErr(rhs, "no such overload")
}
