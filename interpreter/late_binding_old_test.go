package interpreter

import (
	"testing"

	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// TestLateBindEvalZeroArityID verifies the implemented behaviour of lateBindEvalZeroArity.ID().
// The expectation is for the function to forward the call to the wrapped evalZeroArity reference.
func TestLateBindEvalZeroArityID(t *testing.T) {

	expected := &evalZeroArity{
		id:       12,
		function: "f1",
		overload: "f1_zero",
		impl: func(args ...ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallID(t, &lateBindEvalZeroArity{target: expected}, expected)
}

// TestLateBindEvalZeroArityFunction verifies the implemented behaviour of latebindEvalZeroArity.Function().
// The expectation is for the function to forward the call to the wrapped evalZeroArity reference.
func TestLateBindEvalZeroArityFunction(t *testing.T) {

	expected := &evalZeroArity{
		id:       12,
		function: "f1",
		overload: "f1_zero",
		impl: func(args ...ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallFunction(t, &lateBindEvalZeroArity{target: expected}, expected)
}

// TestLateBindEvalZeroArityOverloadID verifies the implemented behaviour of latebindEvalZeroArity.OverloadID().
// The expectation is for the function to forward the call to the wrapped evalZeroArity reference.
func TestLateBindEvalZeroArityOverloadID(t *testing.T) {

	expected := &evalZeroArity{
		id:       12,
		function: "f1",
		overload: "f1_zero",
		impl: func(args ...ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallOverloadID(t, &lateBindEvalZeroArity{target: expected}, expected)

}

// TestLateBindEvalZeroArityArgs verifies the implemented behaviour of latebindEvalZeroArity.Args().
// The expectation is for the function to forward the call to the wrapped evalZeroArity reference.
func TestLateBindEvalZeroArityArgs(t *testing.T) {

	expected := &evalZeroArity{
		id:       12,
		function: "f1",
		overload: "f1_zero",
		impl: func(args ...ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallArgs(t, &lateBindEvalZeroArity{target: expected}, expected)
}

// TestLateBindEvalZeroArityEval verifies the implemented behaviour of lateBindZeroArity.Eval(Activation).
// The expectation is that the method inspects the activation tree and if it finds a function overload
// declaration for the specified overload identifier, it creates a new instance of evalZeroArity and it
// configures it with the new overload, before calling Eval(Activation) on the new evalZeroArity struct
// reference.
func TestLateBindEvalZeroArityEval(t *testing.T) {

	unwrap := func(target Interpretable) *evalZeroArity {

		lba, ok := target.(*lateBindEvalZeroArity)
		if !ok {
			t.Fatalf("unexpected type in test case (got: %T, want: %T)", target, &lateBindEvalZeroArity{})
		}
		actual := lba.target
		if actual == nil {
			t.Errorf("unexpected nil reference in lateBindEvalBinary")
		}

		return actual
	}

	// zeroArity is a function returns a function that upon call produces always a new instance
	// of evalZeroArity configured in the same way. This behaviour is required to ensure that the
	// invocation to Eval(Activation) does not mutate the original evalZeroArity wrapped in the
	// lateBindEvalZeroArity. The reason why we need to introduce this complexity is because the
	// Interpretable are pointers. Therefore, we can't rely on the reference stored in the test
	// case to determine whether it has remained unchanged or not, but we need a fresh instance
	// that is identical to the original configured with the test case.
	zeroArity := func(id int64, function string, overload string, impl func(...ref.Val) ref.Val) func() *evalZeroArity {

		return func() *evalZeroArity {

			return &evalZeroArity{
				id:       id,
				function: function,
				overload: overload,
				impl:     impl,
			}
		}
	}

	nestedActivation, _ := prepareNestedActivation()

	t1 := zeroArity(50, "f1", "f1_int", func(values ...ref.Val) ref.Val { return types.Int(30) })
	t2 := zeroArity(51, "f1", "f1_int", func(values ...ref.Val) ref.Val { return types.Int(50) })
	t3 := zeroArity(52, "f3", "f3_string", func(values ...ref.Val) ref.Val { return types.String("f3_original_string") })

	testInterpretableEval(t, []interpretableTestCase{
		{
			name:       "OK_Simple_Case_No_Overload",
			activation: &emptyActivation{},
			candidate: &lateBindEvalZeroArity{
				target: t1(),
			},
			expect: chain(
				expectEqual(t1(), unwrap, nil, &emptyActivation{}),
				expectValue(types.Int(30)),
			),
		}, {
			name:       "OK_Complex_Case_No_Overload",
			activation: nestedActivation(),
			candidate: &lateBindEvalZeroArity{
				target: t2(),
			},
			expect: chain(
				expectEqual(t2(), unwrap, nil, nestedActivation()),
				expectValue(types.Int(50)),
			),
		}, {
			name:       "OK_Complex_Case_With_Overload",
			activation: nestedActivation(),
			candidate: &lateBindEvalZeroArity{
				target: t3(),
			},
			expect: chain(
				expectEqual(t3(), unwrap, nil, nestedActivation()),
				expectValue(types.String("f3_string")),
			),
		},
	})
}

// TestLateBindEvalUnaryID verifies the implemented behaviour of lateBindEvalUnary.ID().
// The expectation is for the function to forward the call to the wrapped evalUnary and
// therefore return the same value returned by an invocation of the same method on the
// wrapped reference.
func TestLateBindEvalUnaryID(t *testing.T) {

	expected := &evalUnary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		arg:      NewConstValue(53, types.Int(3)),
		impl: func(arg ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallID(t, &lateBindEvalUnary{target: expected}, expected)
}

// TestLateBindEvalUnaryFunction verifies the implemented behaviour of lateBindEvalUnary.Function().
// The expectation is for the function to forward the call to the wrapped evalUnary and therefore
// return the same value returned by an invocation of the same method on the wrapped reference.
func TestLateBindEvalUnaryFunction(t *testing.T) {

	expected := &evalUnary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		arg:      NewConstValue(53, types.Int(3)),
		impl: func(arg ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallFunction(t, &lateBindEvalUnary{target: expected}, expected)

}

// TestLateBindEvalUnaryOverloadID verifies the implemented behaviour of lateBindEvalUnary.OverloadID().
// The expectation is for the function to forward the call to the wrapped evalUnary and therefore return
// the same value returned by an invocation of the same method on the wrapped reference.
func TestLateBindEvalUnaryOverloadID(t *testing.T) {

	expected := &evalUnary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		arg:      NewConstValue(53, types.Int(3)),
		impl: func(arg ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallOverloadID(t, &lateBindEvalUnary{target: expected}, expected)
}

// TestLateBindEvalUnaryArgs verifies the implemented behaviour of lateBindEvalUnary.Args().
// The expectation is for the function to forward the call to the wrapped evalUnary and
// therefore return the same slice it would be returned by an invocation of the same method
// on the wrapped reference.
func TestLateBindEvalUnaryArgs(t *testing.T) {

	expected := &evalUnary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		arg:      NewConstValue(53, types.Int(3)),
		impl: func(arg ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallArgs(t, &lateBindEvalUnary{target: expected}, expected)

}

// TestLateBindEvalUnaryEval verifies the implemented behaviour of lateBindEvalUnary.Eval(Activation).
// The expectation is that the method inspects the activation tree and if it finds a function overload
// declaration for the specified overload identifier, it creates a new instance of evalUnary and it
// configures it with the new overload, before calling Eval(Activation) on the new evalUnary struct.
func TestLateBindEvalUnaryEval(t *testing.T) {

	unwrap := func(target Interpretable) *evalUnary {

		lba, ok := target.(*lateBindEvalUnary)
		if !ok {
			t.Fatalf("unexpected type in test case (got: %T, want: %T)", target, &lateBindEvalUnary{})
		}
		actual := lba.target
		if actual == nil {
			t.Errorf("unexpected nil reference in lateBindEvalBinary")
		}

		return actual
	}

	extra := func(t *testing.T, actual *evalUnary, expected *evalUnary) {

		if expected.trait != actual.trait {
			t.Errorf("trait mismatch (got: %d, want: %d)", actual.trait, expected.trait)
		}

		if expected.nonStrict != actual.nonStrict {
			t.Errorf("nonStrict mismatch (got: %v, want: %v)", actual.nonStrict, expected.nonStrict)
		}
	}

	// unaryInterpretable generates a function that always returns a reference to evalUnary configured
	// with the parameters passed as argument and defaulting trait to 0 and nonStrict to false. This
	// function is used to produce a neew struct every time with always the same values, so that it is
	// possible to have a constant term of comparison which is not affected by the execution.
	unaryInterpretable := func(id int64, function string, overload string, arg Interpretable, impl functions.UnaryOp) func() *evalUnary {

		return func() *evalUnary {
			return &evalUnary{
				id:        id,
				function:  function,
				overload:  overload,
				arg:       arg,
				impl:      impl,
				trait:     0,
				nonStrict: false,
			}
		}
	}

	t1 := unaryInterpretable(45, "f1", "f1_int_int", NewConstValue(46, types.Int(5)), func(arg ref.Val) ref.Val {

		number, _ := arg.(types.Int)
		return number.Multiply(types.Int(2))
	})

	t2 := unaryInterpretable(47, "f1", "f1_int_int", NewConstValue(48, types.Int(7)), func(arg ref.Val) ref.Val {

		number, _ := arg.(types.Int)
		return number.Subtract(types.Int(2))
	})

	t3 := unaryInterpretable(49, "f1", "f1_string_string", NewConstValue(50, types.String("hola")), func(arg ref.Val) ref.Val {

		text, _ := arg.(types.String)
		return text.Add(types.String("_amigo"))
	})

	nestedActivation, _ := prepareNestedActivation()

	testInterpretableEval(t, []interpretableTestCase{
		{

			name:       "OK_Simple_Case_No_Overload",
			activation: &emptyActivation{},
			candidate: &lateBindEvalUnary{
				target: t1(),
			},
			expect: chain(
				expectEqual(t1(), unwrap, extra, &emptyActivation{}),
				expectValue(types.Int(10)),
			),
		}, {
			name:       "OK_Complex_Case_No_Overload",
			activation: nestedActivation(),
			candidate: &lateBindEvalUnary{
				target: t2(),
			},
			expect: chain(
				expectEqual(t2(), unwrap, extra, nestedActivation()),
				expectValue(types.Int(5)),
			),
		}, {
			name:       "OK_Complex_Case_With_Overload",
			activation: nestedActivation(),
			candidate: &lateBindEvalUnary{
				target: t3(),
			},
			expect: chain(
				expectEqual(t3(), unwrap, extra, nestedActivation()),
				expectValue(types.String("HOLA")),
			),
		},
	})
}

// TestLateBindEvalBinaryID verifies the implemented behaviour of lateBindEvalBinary.ID().
// The expectation is for the function to forward the call to the wrapped evalBinary and
// therefore return the same value returned by an invocation of the same method on the
// wrapped reference.
func TestLateBindEvalBinaryID(t *testing.T) {

	expected := &evalBinary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		lhs:      NewConstValue(53, types.Int(3)),
		rhs:      NewConstValue(54, types.Int(6)),
		impl: func(lhs ref.Val, rhs ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallID(t, &lateBindEvalBinary{target: expected}, expected)
}

// TestLateBindEvalBinaryFunction verifies the implemented behaviour of lateBindEvalBinary.Function().
// The expectation is for the function to forward the call to the wrapped evalBinary and therefore return
// the same value returned by an invocation of the same method on the wrapped reference.
func TestLateBindEvalBinaryFunction(t *testing.T) {

	expected := &evalBinary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		lhs:      NewConstValue(53, types.Int(3)),
		rhs:      NewConstValue(54, types.Int(6)),
		impl: func(lhs ref.Val, rhs ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallFunction(t, &lateBindEvalBinary{target: expected}, expected)
}

// TestLateBindEvalBinaryOverloadID verifies the implemented behaviour of lateBindEvalBinary.OverloadID().
// The expectation is for the function to forward the call to the wrapped evalBinary and therefore return
// the same value returned by an invocation of the same method on the wrapped reference.
func TestLateBindEvalBinaryOverloadID(t *testing.T) {

	expected := &evalBinary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		lhs:      NewConstValue(53, types.Int(3)),
		rhs:      NewConstValue(54, types.Int(6)),
		impl: func(lhs ref.Val, rhs ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallOverloadID(t, &lateBindEvalBinary{target: expected}, expected)
}

// TestLateBindEvalBinaryArgs verifies the implemented behaviour of lateBindEvalBinary.Args().
// The expectation is for the function to forward the call to the wrapped evalBinary and
// therefore return the same slice it would be returned by an invocation of the same method
// on the wrapped reference.
func TestLateBindEvalBinaryArgs(t *testing.T) {

	expected := &evalBinary{
		id:       45,
		function: "f3",
		overload: "f3_int",
		lhs:      NewConstValue(53, types.Int(3)),
		rhs:      NewConstValue(54, types.Int(6)),
		impl: func(lhs ref.Val, rhs ref.Val) ref.Val {
			return types.Int(0)
		},
	}

	testInterpretableCallArgs(t, &lateBindEvalBinary{target: expected}, expected)
}

// TestLateBindBinaryEval verifies the implemented behaviour of lateBindEvalBinary.Eval(Activation).
// The expectation is that the method inspects the activation tree and if it finds a function overload
// declaration for the specified overload identifier, it creates a new instance of evalBinary and it
// configures it with the new overload, before calling Eval(Activation) on the new evalBinary struct
// reference.
func TestLateBindBinaryEval(t *testing.T) {

	unwrap := func(target Interpretable) *evalBinary {

		lba, ok := target.(*lateBindEvalBinary)
		if !ok {
			t.Fatalf("unexpected type in test case (got: %T, want: %T)", target, &lateBindEvalBinary{})
		}
		actual := lba.target
		if actual == nil {
			t.Errorf("unexpected nil reference in lateBindEvalBinary")
		}

		return actual
	}

	extra := func(t *testing.T, actual *evalBinary, expected *evalBinary) {

		if expected.trait != actual.trait {
			t.Errorf("trait mismatch (got: %d, want: %d)", actual.trait, expected.trait)
		}

		if expected.nonStrict != actual.nonStrict {
			t.Errorf("nonStrict mismatch (got: %v, want: %v)", actual.nonStrict, expected.nonStrict)
		}
	}

	// binaryInterpretable generates a function that always returns a reference to evalBinary configured
	// with the parameters passed as argument and defaulting trait to 0 and nonStrict to false. This
	// function is used to produce a neew struct every time with always the same values, so that it is
	// possible to have a constant term of comparison which is not affected by the execution.
	binaryInterpretable := func(id int64, function string, overload string, lhs Interpretable, rhs Interpretable, impl functions.BinaryOp) func() *evalBinary {

		return func() *evalBinary {
			return &evalBinary{
				id:        id,
				function:  function,
				overload:  overload,
				lhs:       lhs,
				rhs:       rhs,
				impl:      impl,
				trait:     0,
				nonStrict: false,
			}
		}
	}

	t1 := binaryInterpretable(45, "f1", "f1_int_int_int", NewConstValue(46, types.Int(5)), NewConstValue(47, types.Int(8)), func(lhs ref.Val, rhs ref.Val) ref.Val {

		a, _ := lhs.(types.Int)
		b, _ := rhs.(types.Int)
		return a.Add(b)
	})

	t2 := binaryInterpretable(48, "f1", "f1_int_int_int", NewConstValue(49, types.Int(7)), NewConstValue(50, types.Int(7)), func(lhs ref.Val, rhs ref.Val) ref.Val {

		a, _ := lhs.(types.Int)
		b, _ := rhs.(types.Int)
		return a.Subtract(b)
	})

	t3 := binaryInterpretable(51, "f1", "f1_string_string_string", NewConstValue(52, types.String("hola")), NewConstValue(53, types.String("amigo")), func(lhs ref.Val, rhs ref.Val) ref.Val {

		first, _ := lhs.(types.String)
		second, _ := rhs.(types.String)
		return first.Add(second)
	})

	nestedActivation, _ := prepareNestedActivation()

	testInterpretableEval(t, []interpretableTestCase{
		{

			name:       "OK_Simple_Case_No_Overload",
			activation: &emptyActivation{},
			candidate: &lateBindEvalBinary{
				target: t1(),
			},
			expect: chain(
				expectEqual(t1(), unwrap, extra, &emptyActivation{}),
				expectValue(types.Int(13)),
			),
		}, {
			name:       "OK_Complex_Case_No_Overload",
			activation: nestedActivation(),
			candidate: &lateBindEvalBinary{
				target: t2(),
			},
			expect: chain(
				expectEqual(t2(), unwrap, extra, nestedActivation()),
				expectValue(types.Int(0)),
			),
		}, {
			name:       "OK_Complex_Case_With_Overload",
			activation: nestedActivation(),
			candidate: &lateBindEvalBinary{
				target: t3(),
			},
			expect: chain(
				expectEqual(t3(), unwrap, extra, nestedActivation()),
				expectValue(types.String("amigo hola")),
			),
		},
	})
}

// TestLateBindEvalVarArgsID verifies the implemented behaviour of lateBindEvalVarArgs.ID().
// The expectation is for the function to forward the call to the wrapped evalVarArgs and
// therefore return the same value returned by an invocation of the same method on the
// wrapped reference.
func TestLateBindEvalVarArgsID(t *testing.T) {

	expected := &evalVarArgs{
		id:       65,
		function: "f4",
		overload: "f4_bool",
		args: []Interpretable{
			NewConstValue(66, types.True),
			NewConstValue(67, types.False),
			NewConstValue(67, types.False),
		},
		trait:     0,
		nonStrict: false,
		impl: func(_ ...ref.Val) ref.Val {
			return types.False
		},
	}

	testInterpretableCallID(t, &lateBindEvalVarArgs{target: expected}, expected)
}

// TestLateBindEvalVarArgsFunction verifies the implemented behaviour of lateBindEvalVarArgs.Function().
// The expectation is for the function to forward the call to the wrapped evalVarArgs and therefore
// return the same value returned by an invocation of the same method on the wrapped reference.
func TestLateBindEvalVarArgsFunction(t *testing.T) {

	expected := &evalVarArgs{
		id:       65,
		function: "f4",
		overload: "f4_bool",
		args: []Interpretable{
			NewConstValue(66, types.True),
			NewConstValue(67, types.False),
			NewConstValue(67, types.False),
		},
		trait:     0,
		nonStrict: false,
		impl: func(_ ...ref.Val) ref.Val {
			return types.False
		},
	}

	testInterpretableCallFunction(t, &lateBindEvalVarArgs{target: expected}, expected)
}

// TestLateBindEvalVarArgsOverloadID verifies the implemented behaviour of lateBindEvalVarArgs.OverloadID().
// The expectation is for the function to forward the call to the wrapped evalVarArgs and therefore return
// the same value returned by an invocation of the same method on the wrapped reference.
func TestLateBindEvalVarArgsOverloadID(t *testing.T) {

	expected := &evalVarArgs{
		id:       65,
		function: "f4",
		overload: "f4_bool",
		args: []Interpretable{
			NewConstValue(66, types.True),
			NewConstValue(67, types.False),
			NewConstValue(67, types.False),
		},
		trait:     0,
		nonStrict: false,
		impl: func(_ ...ref.Val) ref.Val {
			return types.False
		},
	}

	testInterpretableCallOverloadID(t, &lateBindEvalVarArgs{target: expected}, expected)
}

// TestLateBindEvalVarArgsArgs verifies the implemented behaviour of lateBindEvalVarArgs.Args().
// The expectation is for the function to forward the call to the wrapped evalVarArgs and
// therefore return the same slice it would be returned by an invocation of the same method on
// the wrapped reference.
func TestLateBindEvalVarArgsArgs(t *testing.T) {

	expected := &evalVarArgs{
		id:       65,
		function: "f4",
		overload: "f4_bool",
		args: []Interpretable{
			NewConstValue(66, types.True),
			NewConstValue(67, types.False),
			NewConstValue(67, types.False),
		},
		trait:     0,
		nonStrict: false,
		impl: func(_ ...ref.Val) ref.Val {
			return types.False
		},
	}

	testInterpretableCallArgs(t, &lateBindEvalVarArgs{target: expected}, expected)
}

// TestLateBindEvalVargArgsEval verifies the implemented behaviour of lateBindEvalVarArgs.Eval(Activation).
// The expectation is that the method inspects the activation tree and if it finds a function overload
// declaration for the specified overload identifier, it creates a new instance of evalVarArgs and it
// configures it with the new overload, before calling Eval(Activation) on the new evalVarArgs struct
// reference.
func TestLateBindEvalVargArgsEval(t *testing.T) {

	nestedActivation, _ := prepareNestedActivation()

	varArgsIntepretable := func(id int64, function string, overload string, impl functions.FunctionOp, args ...Interpretable) *evalVarArgs {

		return &evalVarArgs{
			id:        id,
			function:  function,
			overload:  overload,
			args:      args,
			impl:      impl,
			nonStrict: false,
			trait:     0,
		}
	}

	// original concatenate in reverse order the strings
	// that are passed as argument.
	original := func(args ...ref.Val) ref.Val {

		result := types.String("")
		for i := len(args); i > 0; i-- {
			text := args[i-1].(types.String)
			result = result.Add(text).(types.String)
		}
		return result
	}

	t1 := varArgsIntepretable(67, "f9", "f9_varargs_string", original,
		NewConstValue(68, types.String("one")),
		NewConstValue(69, types.String(" ")),
		NewConstValue(70, types.String("two")),
		NewConstValue(71, types.String(" ")),
		NewConstValue(72, types.String("three")),
	)
	t2 := varArgsIntepretable(73, "f10", "f10_varargs_string", original,
		NewConstValue(74, types.String("one")),
		NewConstValue(75, types.String(",")),
		NewConstValue(76, types.String("two")),
	)
	t3 := varArgsIntepretable(77, "f6", "f6_string", original,
		NewConstValue(78, types.String("one")),
		NewConstValue(79, types.String(":")),
		NewConstValue(80, types.String("two")),
	)

	testInterpretableEval(t, []interpretableTestCase{
		{
			name: "OK_Simple_Case_No_Overload",
			candidate: &lateBindEvalVarArgs{
				target: t1,
			},
			activation: &mapActivation{
				bindings: map[string]any{
					"a": 5,
					"b": true,
				},
			},
			expect: expectValue(types.String("three two one")),
		}, {
			name: "OK_Complex_Case_No_Overload",
			candidate: &lateBindEvalVarArgs{
				target: t2,
			},
			activation: nestedActivation(),
			expect:     expectValue(types.String("two,one")),
		}, {
			name: "OK_Complex_Case_With_Overload",
			candidate: &lateBindEvalVarArgs{
				target: t3,
			},
			activation: nestedActivation(),
			expect:     expectValue(types.String("one:two")),
		},
	})
}

// TestLateBindingOldDecorator verifies the implemented behaviour of decLateBindOld().
// The expectation is for the function to return an InterepretableDecorator that applies
// the late bind behaviour (old version).
func TestLateBindingOldDecorator(t *testing.T) {

	/*

		// declaractions generates an array of function declaration that are used
		// as a baseline of function overloads statically bound to the expression
		// in the AST. The array contains a single entry (f1) with the following
		// overloads:
		// - f1_int: zero arguments, returns 0
		// - fi_string_int: one string argument, returns the size of the string
		// - f1_int_int_int: two int arguments, returns the sum of the the two
		// - f1_bool_bool_bool_int: three bool arguments, returns the number of
		//   true values in the arguments.
		declarations := func(t *testing.T) map[string]*decls.FunctionDecl {

			return []*decls.FunctionDecl{
				funcDecl(t, "f1",
					decls.Overload(
						"f1_int",
						[]*types.Type{},
						types.IntType,
						decls.FunctionBinding(func(args ...ref.Val) ref.Val {
							return types.Int(0)
						}),
					),
					decls.Overload(
						"f1_string_int",
						[]*types.Type{types.StringType},
						types.IntType,
						decls.UnaryBinding(func(arg ref.Val) ref.Val {
							text := arg.(types.String)

							return text.Size()
						}),
					),
					decls.Overload(
						"f1_int_int_int",
						[]*types.Type{types.IntType, types.IntType},
						types.IntType,
						decls.BinaryBinding(func(lhs ref.Val, rhs ref.Val) ref.Val {

							l := lhs.(types.Int)
							r := rhs.(types.Int)

							return l.Add(r)
						}),
					),
					decls.Overload(
						"f1_bool_bool_bool_int",
						[]*types.Type{types.BoolType, types.BoolType, types.BoolType},
						types.IntType,
						decls.FunctionBinding(func(args ...ref.Val) ref.Val {

							count := 0
							for _, arg := range args {
								b := arg.(types.Bool)
								if b == types.True {
									count = count + 1
								}
							}

							return types.Int(count)
						}),
					),
				),
			}
		}
	*/

	f1_int := func(args ...ref.Val) ref.Val {
		return types.Int(0)
	}

	f1_string_int := func(arg ref.Val) ref.Val {
		text := arg.(types.String)
		return text.Size()
	}

	f1_int_int_int := func(lhs ref.Val, rhs ref.Val) ref.Val {

		l := lhs.(types.Int)
		r := rhs.(types.Int)

		return l.Add(r)
	}

	f1_bool_bool_bool_int := func(args ...ref.Val) ref.Val {

		count := 0
		for _, arg := range args {
			b := arg.(types.Bool)
			if b == types.True {
				count = count + 1
			}
		}

		return types.Int(count)
	}

	// activation returns a lateBindActivation that is configured
	// with the given variables and function overloads.
	activation := func(vars map[string]any, ovls ...*functions.Overload) *lateBindActivation {

		d := &defaultDispatcher{
			parent:    nil,
			overloads: overloadMap{},
		}
		for _, ovl := range ovls {
			d.overloads[ovl.Operator] = ovl
		}
		return &lateBindActivation{
			vars:       &mapActivation{bindings: vars},
			dispatcher: d,
		}
	}

	// this is used to simplify the creation of attributes for the purpose of
	// testing.
	attrFactory := NewAttributeFactory(containers.DefaultContainer, types.DefaultTypeAdapter, types.NewEmptyRegistry())

	// this is used to verify that the late-bind of the object creation
	// expression does what is expected to do. By no means this is a full
	// implementation of the Provider interface.
	exampleProvider := &testProvider{
		adapter:  types.DefaultTypeAdapter,
		typeName: "example",
		fields: map[string]*types.Type{
			"quantity":  types.IntType,
			"frequency": types.IntType,
		},
	}

	testCases := []struct {
		name          string
		interpretable Interpretable
		activation    Activation
		options       []PlannerOption
		expect        func(t *testing.T, expr Interpretable, result ref.Val)
	}{

		// Test Group 01.A - No Overloads, expressions with single function call.
		{
			name: "01.A.01__No_Overloads_Zero_Arity",
			interpretable: &evalZeroArity{
				id:       1,
				function: "f1",
				overload: "f1_int",
				impl:     f1_int,
			},
			activation: EmptyActivation(),
			expect: chain(
				expectType(&lateBindEvalZeroArity{}),
				expectValue(types.Int(0)),
			),
		},
		{
			name: "01.A.02__No_Overloads_Unary",
			interpretable: &evalUnary{
				id:        1,
				function:  "f1",
				overload:  "f1_string_int",
				trait:     0,
				nonStrict: false,
				arg:       NewConstValue(2, types.String("banana")),
				impl:      f1_string_int,
			},
			activation: EmptyActivation(),
			expect: chain(
				expectType(&lateBindEvalUnary{}),
				expectValue(types.Int(6)),
			),
		},
		{
			name: "01.A.03__No_Overloads_Binary",
			interpretable: &evalBinary{
				id:        1,
				function:  "f1",
				overload:  "f1_int_int_int",
				trait:     0,
				nonStrict: false,
				lhs:       NewConstValue(2, types.Int(2)),
				rhs:       NewConstValue(3, types.Int(3)),
				impl:      f1_int_int_int,
			},
			activation: NewHierarchicalActivation(
				&mapActivation{bindings: map[string]any{"a": 2, "b": 5}},
				&mapActivation{bindings: map[string]any{"c": 3, "d": 6}},
			),
			expect: chain(
				expectType(&lateBindEvalBinary{}),
				expectValue(types.Int(5)),
			),
		},
		{
			name: "01.A.04__No_Overloads_VarArgs",
			interpretable: &evalVarArgs{
				id:        1,
				function:  "f1",
				overload:  "f1_bool_bool_bool_int",
				trait:     0,
				nonStrict: false,
				args: []Interpretable{
					NewConstValue(2, types.True),
					NewConstValue(2, types.False),
					NewConstValue(2, types.True),
				},
				impl: f1_bool_bool_bool_int,
			},
			activation: NewHierarchicalActivation(
				&mapActivation{bindings: map[string]any{"a": 6, "b": 3}},
				activation(
					map[string]any{"d": 3, "f": 2},
					function("f2", 0, false, func(args ...ref.Val) ref.Val { return types.String("f2") }),
					function("f3", 0, false, func(args ...ref.Val) ref.Val { return types.String("f3") }),
				),
			),
			expect: chain(
				expectType(&lateBindEvalVarArgs{}),
				expectValue(types.Int(2)),
			),
		},

		// Test Group 01.B - Simple Function Calls with Activation Overloads
		{
			name: "01.B.01__With_Overloads_Zero_Arity",
			// without override evaluates to 0
			interpretable: &evalZeroArity{
				id:       1,
				function: "f1",
				overload: "f1_int",
				impl:     f1_int,
			},
			activation: &lateBindActivation{
				vars: EmptyActivation(),
				dispatcher: &defaultDispatcher{
					overloads: overloadMap{
						"f1_int": function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(-1) }),
					},
				},
			},
			expect: chain(
				expectType(&lateBindEvalZeroArity{}),
				expectValue(types.Int(-1)),
			),
		},
		{
			name: "01.B.02__With_Overloads_Unary",
			// without override evaluates to 6
			interpretable: &evalUnary{
				id:        1,
				function:  "f1",
				overload:  "f1_string_int",
				trait:     0,
				nonStrict: false,
				arg:       NewConstValue(2, types.String("banana")),
				impl:      f1_string_int,
			},
			activation: &hierarchicalActivation{
				parent: &mapActivation{
					bindings: map[string]any{"a": 3, "b": 4},
				},
				child: activation(
					map[string]any{},
					unary("f1_string_int", 0, false, func(arg ref.Val) ref.Val {
						text := arg.(types.String)
						size := text.Size().(types.Int)
						return size.Add(size)
					}),
				),
			},
			expect: chain(
				expectType(&lateBindEvalUnary{}),
				expectValue(types.Int(12)),
			),
		},
		{
			name: "01.B.03__With_Overloads_Binary",
			// without override evaluates to 5
			interpretable: &evalBinary{
				id:        1,
				function:  "f1",
				overload:  "f1_int_int_int",
				trait:     0,
				nonStrict: false,
				lhs:       NewConstValue(2, types.Int(2)),
				rhs:       NewConstValue(3, types.Int(3)),
				impl:      f1_int_int_int,
			},
			activation: NewHierarchicalActivation(
				activation(
					map[string]any{},
					binary("f1_int_int_int", 0, false, func(lhs ref.Val, rhs ref.Val) ref.Val {

						l := lhs.(types.Int)
						r := rhs.(types.Int)
						t := l.Add(r).(types.Int)
						return t.Multiply(types.Int(2))
					}),
				),
				&mapActivation{bindings: map[string]any{"c": 3, "d": 6}},
			),
			expect: chain(
				expectType(&lateBindEvalBinary{}),
				expectValue(types.Int(10)),
			),
		},
		{
			name: "01.B.04__With_Overloads_VarArgs",
			// without override evaluates to 2
			interpretable: &evalVarArgs{
				id:        1,
				function:  "f1",
				overload:  "f1_bool_bool_bool_int",
				trait:     0,
				nonStrict: false,
				args: []Interpretable{
					NewConstValue(2, types.True),
					NewConstValue(2, types.False),
					NewConstValue(2, types.True),
				},
				impl: f1_bool_bool_bool_int,
			},
			activation: NewHierarchicalActivation(
				&mapActivation{bindings: map[string]any{"a": 6, "b": 3}},
				activation(
					map[string]any{"d": 3, "f": 2},
					function("f2", 0, false, func(args ...ref.Val) ref.Val { return types.String("f2") }),
					function("f3", 0, false, func(args ...ref.Val) ref.Val { return types.String("f3") }),
					function("f1_bool_bool_bool_int", 0, false, func(args ...ref.Val) ref.Val {
						counter := 0
						for _, arg := range args {

							b := arg.(ref.Val)
							if b == types.False {
								counter++
							}
						}
						return types.Int(counter)
					}),
				),
			),
			expect: chain(
				expectType(&lateBindEvalVarArgs{}),
				expectValue(types.Int(1)),
			),
		},

		// Test Group 02 - Equality Operators with Function Calls
		{
			name: "02.01__Equal_Operator_With_Function_Calls",
			// witohut overrides evaluates to true
			interpretable: &evalEq{
				id: 50,
				lhs: &evalZeroArity{
					id:       51,
					function: "f1",
					overload: "f1_int",
					// returns 0
					impl: f1_int,
				},
				rhs: NewConstValue(52, types.Int(0)),
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(args ...ref.Val) ref.Val { return types.Int(10) }),
			),
			expect: chain(
				expectType(&evalEq{}),
				expectValue(types.False),
			),
		},
		{
			name: "02.02__Not_Equal_Operator_With_Function_Calls",
			// without override evaluates to false
			interpretable: &evalNe{
				id: 50,
				lhs: &evalZeroArity{
					id:       51,
					function: "f1",
					overload: "f1_int",
					// returns 0
					impl: f1_int,
				},
				rhs: NewConstValue(52, types.Int(0)),
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(args ...ref.Val) ref.Val { return types.Int(10) }),
			),
			expect: chain(
				expectType(&evalNe{}),
				expectValue(types.True),
			),
		},
		// Test Group 03 - Logical Operators with Function Calls
		{
			name: "03.01__And_Operator_With_Function_Calls",
			// 2 == f1() && f1(example) == 14
			interpretable: &evalAnd{
				id: 50,
				terms: []Interpretable{
					// without override evaluates to false
					&evalEq{
						id:  51,
						lhs: NewConstValue(52, types.Int(2)),
						rhs: &evalZeroArity{
							id:       53,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
					},
					// without override evaluates to false
					&evalEq{
						id: 54,
						lhs: &evalUnary{
							id:       55,
							function: "f1",
							overload: "f1_string_int",
							// returns len(example) = 7
							impl: f1_string_int,
							arg:  NewConstValue(56, types.String("example")),
						},
						rhs: NewConstValue(57, types.Int(14)),
					},
				},
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(args ...ref.Val) ref.Val { return types.Int(2) }),
				unary("f1_string_int", 0, false, func(arg ref.Val) ref.Val {
					text := arg.(types.String)
					return types.Int(2 * len(text.Value().(string)))
				}),
			),
			expect: chain(
				expectType(&evalAnd{}),
				expectValue(types.True),
			),
		},
		{
			name: "03.02__Or_Operator_With_Function_Calls",
			// f1() == 0 || f1(true, true, true) == 3
			interpretable: &evalOr{
				id: 51,
				terms: []Interpretable{
					// without override evaluates to true
					&evalEq{
						id: 51,
						lhs: &evalZeroArity{
							id:       52,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
						rhs: NewConstValue(53, types.Int(0)),
					},
					// without override evaluates to true
					&evalEq{
						id: 54,
						lhs: &evalVarArgs{
							id: 55,
							args: []Interpretable{
								NewConstValue(56, types.True),
								NewConstValue(57, types.True),
								NewConstValue(58, types.True),
							},
							function: "f1",
							overload: "f1_bool_bool_bool_int",
							// returns 3 (all three parameters are true)
							impl:      f1_bool_bool_bool_int,
							trait:     0,
							nonStrict: false,
						},
						rhs: NewConstValue(59, types.Int(3)),
					},
				},
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(args ...ref.Val) ref.Val { return types.Int(2) }),
				function("f1_bool_bool_bool_int", 0, false, func(args ...ref.Val) ref.Val { return types.Int(6) }),
			),
			expect: chain(
				expectType(&evalOr{}),
				expectValue(types.False),
			),
		},
		{
			name: "03.03__Exhaustive_Or_With_Function_Calls",
			// f1() == 0 || f1(2,3) == 5 || f1("hello") != 5
			interpretable: &evalExhaustiveOr{
				id: 51,
				terms: []Interpretable{
					// without override evaluates to true
					&evalEq{
						id: 52,
						lhs: &evalZeroArity{
							id:       53,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
						rhs: NewConstValue(54, types.Int(0)),
					},
					// without override evaluates to true
					&evalEq{
						id: 55,
						lhs: &evalBinary{
							id:       56,
							lhs:      NewConstValue(56, types.Int(2)),
							rhs:      NewConstValue(57, types.Int(3)),
							function: "f1",
							overload: "f1_int_int_int",
							// returns 2 + 3 = 5
							impl:      f1_int_int_int,
							trait:     0,
							nonStrict: false,
						},
						rhs: NewConstValue(58, types.Int(5)),
					},
					// without override evaluates to true
					&evalNe{
						id: 59,
						lhs: &evalUnary{
							id:       60,
							arg:      NewConstValue(61, types.String("hello")),
							function: "f1",
							overload: "f1_string_int",
							// returns len(hello) = 5
							impl:      f1_string_int,
							trait:     0,
							nonStrict: false,
						},
						rhs: NewConstValue(62, types.Int(7)),
					},
				},
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(10) }),
				binary("f1_int_int_int", 0, false, func(lhs ref.Val, rhs ref.Val) ref.Val { return types.Int(12) }),
				unary("f1_string_int", 0, false, func(arg ref.Val) ref.Val { return types.Int(7) }),
			),
			expect: chain(
				expectType(&evalExhaustiveOr{}),
				expectValue(types.False),
			),
		},
		{
			name: "03.04__Exhaustive_And_With_Function_Calls",
			// f1() == f1(0,0) && f1() != f1("")
			interpretable: &evalExhaustiveAnd{
				id: 51,
				terms: []Interpretable{
					// evaluates to true (there are no overrides)
					&evalEq{
						id: 52,
						lhs: &evalZeroArity{
							id:       53,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
						rhs: &evalBinary{
							id:       54,
							lhs:      NewConstValue(55, types.IntZero),
							rhs:      NewConstValue(56, types.IntZero),
							function: "f1",
							overload: "f1_int_int_int",
							// returns 0 + 0 = 0
							impl:      f1_int_int_int,
							trait:     0,
							nonStrict: false,
						},
					},
					// without overrides evaluates to false
					&evalNe{
						id: 57,
						lhs: &evalZeroArity{
							id:       58,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
						rhs: &evalUnary{
							id:       59,
							arg:      NewConstValue(60, types.String("")),
							function: "f1",
							overload: "f1_string_int",
							// returns len("") = 0
							impl:      f1_string_int,
							trait:     0,
							nonStrict: false,
						},
					},
				},
			},
			activation: activation(
				map[string]any{},
				unary("f1_string_int", 0, false, func(arg ref.Val) ref.Val { return types.Int(100) }),
			),
			expect: chain(
				expectType(&evalExhaustiveAnd{}),
				expectValue(types.True),
			),
		},

		// Test Group 04 - Conditional Operators with Function Calls
		{
			name: "04.01__Conditional_With_Function_Calls",
			// without overrides evaluates to a = 10
			interpretable: &evalAttr{
				adapter:  types.DefaultTypeAdapter,
				optional: false,
				// f1() == 0 ? a : b
				attr: &conditionalAttribute{
					id: 51,
					expr: &evalEq{
						id: 52,
						lhs: &evalZeroArity{
							id:       53,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
						rhs: NewConstValue(54, types.Int(0)),
					},
					truthy: attrFactory.AbsoluteAttribute(55, "a"),
					falsy:  attrFactory.AbsoluteAttribute(56, "b"),
				},
			},
			activation: activation(
				map[string]any{
					"a": 10,
					"b": 20,
				},
				function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(1) }),
			),
			expect: chain(
				expectType(&evalAttr{}),
				expectValue(types.Int(20)),
			),
		},
		{
			name: "04.02__Exhaustive_Conditional_With_Function_Calls",
			// without overrides evaluates to b = 20
			interpretable: &evalExhaustiveConditional{
				id:      51,
				adapter: types.DefaultTypeAdapter,
				// f1() != 0 ? a : b
				attr: &conditionalAttribute{
					id: 52,
					expr: &evalNe{
						id: 53,
						lhs: &evalZeroArity{
							id:       54,
							function: "f1",
							overload: "f1_int",
							// returns 0
							impl: f1_int,
						},
						rhs: NewConstValue(55, types.Int(0)),
					},
					truthy: attrFactory.AbsoluteAttribute(56, "a"),
					falsy:  attrFactory.AbsoluteAttribute(57, "b"),
				},
			},
			activation: activation(
				map[string]any{
					"a": 10,
					"b": 20,
				},
				function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.IntOne }),
			),
			expect: chain(
				expectType(&evalExhaustiveConditional{}),
				expectValue(types.Int(10)),
			),
		},

		// Test Group 05 - Complex Data Structures with Function Calls
		{
			name: "05.01__Map_With_Function_Calls",
			//
			// {
			//    f1():       f1(1,2),
			//    f1("hi"):   10,
			// }
			// without overrides evaluates to:
			// {
			//    0: 3,
			//    2: 10,
			// }
			interpretable: &evalMap{
				id:           51,
				adapter:      types.DefaultTypeAdapter,
				hasOptionals: false,
				optionals:    []bool{},
				keys: []Interpretable{
					&evalZeroArity{
						id:       52,
						function: "f1",
						overload: "f1_int",
						// returns 0
						impl: f1_int,
					},
					&evalUnary{
						id:       53,
						function: "f1",
						overload: "f1_string_int",
						arg:      NewConstValue(54, types.String("hi")),
						// returns len(hi) = 2
						impl:      f1_string_int,
						trait:     0,
						nonStrict: false,
					},
				},
				vals: []Interpretable{
					&evalBinary{
						id:       55,
						lhs:      NewConstValue(56, types.Int(1)),
						rhs:      NewConstValue(57, types.Int(2)),
						function: "f1",
						overload: "f1_int_int_int",
						// returns 1 + 2 = 3
						impl:      f1_int_int_int,
						trait:     0,
						nonStrict: false,
					},
					NewConstValue(58, types.Int(10)),
				},
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(1) }),
				unary("f1_string_int", 0, false, func(_ ref.Val) ref.Val { return types.Int(2) }),
				binary("f1_int_int_int", 0, false, func(_ ref.Val, _ ref.Val) ref.Val { return types.Int(0) }),
			),
			expect: chain(
				expectType(&evalMap{}),
				expectValue(types.DefaultTypeAdapter.NativeToValue(map[ref.Val]ref.Val{
					types.Int(1): types.Int(0),
					types.Int(2): types.Int(10),
				})),
			),
		},
		{
			name: "05.02__List_With_Function_Calls",
			//
			// [ f1(), f1("hi"), f1(1,3), f1(true,true,true) ]
			//
			// without overrides evaluates to [ 0, 2, 4, 3 ]
			interpretable: &evalList{
				id: 51,
				elems: []Interpretable{
					&evalZeroArity{
						id:       52,
						function: "f1",
						overload: "f1_int",
						// returns 0
						impl: f1_int,
					},
					&evalUnary{
						id:       53,
						function: "f1",
						overload: "f1_string_int",
						arg:      NewConstValue(54, types.String("hi")),
						// returns len(hi) = 2
						impl:      f1_string_int,
						trait:     0,
						nonStrict: false,
					},
					&evalBinary{
						id:       55,
						function: "f1",
						overload: "f1_int_int_int",
						lhs:      NewConstValue(56, types.Int(1)),
						rhs:      NewConstValue(57, types.Int(3)),
						// returns 1 + 3 = 4
						impl:      f1_int_int_int,
						trait:     0,
						nonStrict: false,
					},
					&evalVarArgs{
						id:       58,
						function: "f1",
						overload: "f1_bool_bool_bool_int",
						args: []Interpretable{
							NewConstValue(59, types.True),
							NewConstValue(60, types.True),
							NewConstValue(61, types.True),
						},
						// returns 3 (three true values in arguments)
						impl:      f1_bool_bool_bool_int,
						trait:     0,
						nonStrict: false,
					},
				},
				optionals:    []bool{},
				hasOptionals: false,
				adapter:      types.DefaultTypeAdapter,
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(1) }),
				unary("f1_string_int", 0, false, func(_ ref.Val) ref.Val { return types.Int(2) }),
				binary("f1_int_int_int", 0, false, func(_ ref.Val, _ ref.Val) ref.Val { return types.Int(3) }),
				function("f1_bool_bool_bool_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(4) }),
			),
			expect: chain(
				expectType(&evalList{}),
				expectValue(types.DefaultTypeAdapter.NativeToValue([]ref.Val{
					types.Int(1),
					types.Int(2),
					types.Int(3),
					types.Int(4),
				})),
			),
		},
		{
			name: "05.03__Object_With_Function_Calls",
			interpretable: &evalObj{
				id:       51,
				typeName: "example",
				fields: []string{
					"quantity",
					"frequency",
				},
				vals: []Interpretable{
					&evalZeroArity{
						id:       52,
						function: "f1",
						overload: "f1_int",
						// returns 0
						impl: f1_int,
					},
					&evalUnary{
						id:       53,
						function: "f1",
						overload: "f1_string_int",
						// returns len("weekly") = 6
						impl:      f1_string_int,
						trait:     0,
						nonStrict: false,
						arg:       NewConstValue(54, types.String("weekly")),
					},
				},
				optionals:    []bool{},
				hasOptionals: false,
				provider:     exampleProvider,
			},
			activation: activation(
				map[string]any{},
				function("f1_int", 0, false, func(_ ...ref.Val) ref.Val { return types.Int(10) }),
				unary("f1_string_int", 0, false, func(_ ref.Val) ref.Val { return types.Int(3) }),
			),
			expect: chain(
				expectType(&evalObj{}),
				expectValue(exampleProvider.NewValue("example", map[string]ref.Val{
					"weekly":   types.Int(3),
					"quantity": types.Int(10),
				})),
			),
		},
		// Test Group 06 - Macros
		// Test Group 07 - Set Membership
		// Test Group 08 - EvalObserver Alterations
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {

			candidate := decLateBindingOld()
			if candidate == nil {
				t.Fatal("pre-condition failed: intepretable decorator is nil")
			}

			interpretable, err := candidate(testCase.interpretable)
			if err != nil {
				t.Fatalf("unexpected error while applying late binding decorator (cause: %v)", err)
			}
			result := interpretable.Eval(testCase.activation)
			testCase.expect(t, interpretable, result)
		})
	}

}

// testInterpretableCallID is a convenience function that verifies that the value
// returned by actual.ID() is the same as the value return by expected.ID().
func testInterpretableCallID[A, E InterpretableCall](t *testing.T, actual A, expected E) {

	aID := actual.ID()
	eID := expected.ID()

	if aID != eID {
		t.Errorf("identifier mismatch (got: %d, want: %d)", aID, eID)
	}
}

// testInterpretableCallFunction is a convenience function that verifies that the value
// returned by actual.Function() is the same as the value return by expected.Function().
func testInterpretableCallFunction[A, E InterpretableCall](t *testing.T, actual A, expected E) {

	af := actual.Function()
	ef := expected.Function()

	if af != ef {
		t.Errorf("function mismatch (got %s, want: %s)", af, ef)
	}
}

// testInterpretableCallOverloadID is a convenience function that verifies that the value
// returned by actual.OverloadID() is the same as the value return by expected.OverloadID().
func testInterpretableCallOverloadID[A, E InterpretableCall](t *testing.T, actual A, expected E) {

	aOvl := actual.OverloadID()
	eOvl := expected.OverloadID()

	if aOvl != eOvl {
		t.Errorf("overload identifier mismatch (got %s, want: %s)", aOvl, eOvl)
	}

}

// testInterpretableCallArgs is a convenience function that verifies that the slice
// returned by actual.Args() has the same content (and order) of the slice returned
// by expected.Args().
func testInterpretableCallArgs[A, E InterpretableCall](t *testing.T, actual A, expected E) {

	aArgs := actual.Args()
	eArgs := expected.Args()

	if len(eArgs) != len(aArgs) {
		t.Fatalf("args size mismatch (got: %d, want: %d)", len(aArgs), len(eArgs))
	}

	for index, eArg := range eArgs {

		aArg := aArgs[index]
		if eArg != aArg {
			t.Errorf("arg (index: %d) mismatch (got: %v, want: %v)", index, aArg, eArg)
		}
	}
}

// interpretableTestCase defines the structure used to model test cases
// for all Interpretable.Eval(Activation) that are tested in this file.
type interpretableTestCase struct {
	name       string
	candidate  Interpretable
	activation Activation
	expect     func(t *testing.T, target Interpretable, actual ref.Val)
}

// testInterpretableEval tests the implementation of Interpretable that is
// configured in each of the test cases by invoking the Eval(Activation)
// method and then validating the outcome via the expectation function
// defined in the test case.
func testInterpretableEval(t *testing.T, testCases []interpretableTestCase) {

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {

			actual := testCase.candidate.Eval(testCase.activation)
			testCase.expect(t, testCase.candidate, actual)
		})
	}
}

// expectEqual produces an expectation function that verifies that the given Interpretable contains
// a reference that is equal to the supplied value passed to expectEqual. The function first invokes
// the unwrap function to resolve the actual entity to compare, then it tests all the fields that are
// exposed by the InterpretableCall interface to ensure that they are the same as those exposed by the
// expected IntepretableCall. If the extra function is not nil, it is then invoked to allow for the
// customisation of further checks, and finally it invokes the Eval method on both to compare the
// results.
func expectEqual[I InterpretableCall](
	expected I,
	unwrap func(wrapper Interpretable) I,
	extra func(t *testing.T, actual I, expected I),
	ctx Activation,
) func(t *testing.T, target Interpretable, _ ref.Val) {

	return func(t *testing.T, target Interpretable, _ ref.Val) {

		t.Helper()

		actual := unwrap(target)
		testInterpretableCallID(t, actual, expected)
		testInterpretableCallFunction(t, actual, expected)
		testInterpretableCallOverloadID(t, actual, expected)
		testInterpretableCallArgs(t, actual, expected)
		if extra != nil {
			extra(t, actual, expected)
		}

		// we cannot compare functions, therefore we need
		// invoke their execution.
		aVal := actual.Eval(ctx)
		eVal := actual.Eval(ctx)

		if eVal.Equal(aVal) == types.False {
			t.Errorf("function overload has been mutated on original target (invoke, got: %v, want: %v)", aVal, eVal)
		}
	}
}

// testProvider is a stripped down implementation used for the purpose
// of testing the late bind decorator. It can be configured with a type
// name and a list of fields, which is then implemented as a map.
type testProvider struct {
	adapter  types.Adapter
	typeName string
	fields   map[string]*types.Type
}

// EnumValue is a dummy implementation of Provider.EnumValue(string) and
// returns a types.Err since the provider does not support any enumeration
// definition.
func (et *testProvider) EnumValue(enumName string) ref.Val {
	return types.NewErr("unknown enum name '%s'", enumName)
}

// FindIdent is a dummy implementatio of Provider.FindIdent(string) and
// returns a pair (nil, false) for all values of identName.
func (et *testProvider) FindIdent(identName string) (ref.Val, bool) {

	return nil, false
}

// FindStructType checks whether the structType matches the pre-configured
// type and if so, it does return the corresponding map type definition
// for the type.
func (et *testProvider) FindStructType(structType string) (*types.Type, bool) {

	if et.typeName != structType {

		return nil, false
	}

	return types.NewMapType(types.StringType, types.DynType), true
}

// FindStructFieldNames checks whether the structType matches the pre-configured
// type and if so, it does return the list of fields configured with the provider.
func (et *testProvider) FindStructFieldNames(structType string) ([]string, bool) {

	if et.typeName != structType {

		return nil, false
	}

	names := make([]string, len(et.fields))
	index := 0
	for k, _ := range et.fields {
		names[index] = k
		index++
	}

	return names, true
}

// FindStructFieldType checks whether the structType matches the pre-configured type and if so, it
// returns a FieldType reference that provides information about the type, but without the IsSet and
// GetFrom accessor functions. If either structType or fieldName do not match the pre-configured type
// and fields the pair (nil, false) is returned.
func (er *testProvider) FindStructFieldType(structType, fieldName string) (*types.FieldType, bool) {

	if er.typeName != structType {
		return nil, false
	}

	fType, found := er.fields[fieldName]
	if !found {
		return nil, false
	}

	return &types.FieldType{
		Type:    fType,
		IsSet:   nil,
		GetFrom: nil,
	}, true
}

// NewValue checks whether the structType matches the pre-configured type and if so, it
// creates a dynamic map that is configured with the fields that are passed as arguments.
// If any of these fields is not contained in the map of pre-configured fields the function
// returns an error.
func (et *testProvider) NewValue(structType string, fields map[string]ref.Val) ref.Val {

	if structType != et.typeName {
		return types.NewErr("unknown struct '%s'", structType)
	}

	for name, _ := range fields {
		_, found := et.fields[name]
		if !found {
			return types.NewErr("field '%s' does not exist in struct '%s'", name, structType)
		}
	}

	return types.NewDynamicMap(et.adapter, fields)

}
