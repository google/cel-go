package interpreter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/cel-go/common/containers"
	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

// TestUncheckedAstError verifies the implemented behaviour of UncheckedAstError. The expectation
// is for the function to return a non-nil error implementation that contains the message defined
// by the constant errrorUncheckedAst.
func TestUncheckedAstError(t *testing.T) {

	candidate := UncheckedAstError()
	if candidate == nil {
		t.Fatalf("UncheckedAstError returned nil, a non-nil error implementation is expected")
	}
	if candidate.Error() != errorUncheckedAst {
		t.Errorf("UncheckedAstError has unexpected message (got: %s, want: %s)", candidate.Error(), errorUncheckedAst)
	}
}

// TestNewLateBindingActivation verifies the implementation of NewLateBindingActivation. The
// expectation is for the constructor function to produce a LateBindActivation implementation
// (i.e. lateBindActivation) that is configured with the given parent activation and with a
// dispatcher declaring containing the specified function overloads. The function should return
// a nil implementation and an error in case of duplicate overload function definitions or a nil
// activation.
func TestNewLateBindingActivation(t *testing.T) {

	// expectActivation generates an expectation function that verifies that
	// the outcome of NewLateBindingActivation has not generated any error and
	// contains the given activation as well as the specified function overloads.
	expectActivation := func(expected Activation, overloads ...*functions.Overload) func(t *testing.T, actual LateBindActivation, err error) {

		return func(t *testing.T, actual LateBindActivation, err error) {

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if actual == nil {
				t.Errorf("expected non nil activation")
			} else {

				lateBind, ok := actual.(*lateBindActivation)
				if !ok {
					t.Errorf("unexpected activation type (got: %T, want: %T)", actual, &lateBindActivation{})
				} else {
					if lateBind.vars != expected {
						t.Errorf("unexpected wrapped activation (got: %v, want: %v)", lateBind.vars, expected)
					}
					if lateBind.dispatcher == nil {
						t.Errorf("expected non-nil dispatcher")
					}

					actualIds := lateBind.dispatcher.OverloadIds()
					if len(actualIds) != len(overloads) {
						t.Errorf("number of overloads do not match (got: %d, want: %d)", len(actualIds), len(overloads))
					} else {

						for _, expOvl := range overloads {

							actOvl, found := lateBind.dispatcher.FindOverload(expOvl.Operator)
							if !found {
								t.Errorf("expected overload (id: %s)", expOvl.Operator)
							} else {

								if expOvl != actOvl {
									t.Errorf("overload (id: %s) mismatch", expOvl.Operator)
								}
							}
						}
					}
				}
			}
		}
	}

	// expectError generates an expectation function that checks that the
	// outcome of NewLateBindActivation produces a nil activation and an
	// error that contains the specified message.
	expectError := func(msg string) func(t *testing.T, actual LateBindActivation, err error) {

		return func(t *testing.T, actual LateBindActivation, err error) {

			if actual != nil {
				t.Errorf("expected nil activation")
			}
			if err == nil {
				t.Errorf("expected non-nil error")
			} else {
				if !strings.Contains(err.Error(), msg) {
					t.Errorf("error message (value: %s) does not contain '%s'", err.Error(), msg)
				}
			}
		}
	}

	f1_string_string := unary("f1_string_string", 0, false, func(value ref.Val) ref.Val {
		return types.String("f1_string_string")
	})

	f1_string_string_string := binary("f1_string_string_string", 0, false, func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.String("f1_string_string_string")
	})

	f1_varargs_string := function("f1_varargs_string", 0, false, func(args ...ref.Val) ref.Val {
		return types.String("f1_varargs_string")
	})

	f2_string := function("f2_string", 0, false, func(args ...ref.Val) ref.Val {
		return types.String("f2_string")
	})

	f2_string_string := &functions.Overload{
		Operator:  "f2_string_string",
		NonStrict: true,
		Unary: func(arg ref.Val) ref.Val {
			return types.String("f2_string_string")
		},
	}

	actHierarchical := &hierarchicalActivation{
		parent: &emptyActivation{},
		child: &mapActivation{
			bindings: map[string]any{},
		},
	}

	actEmpty := &emptyActivation{}

	testCases := []struct {
		name       string
		activation Activation
		overloads  []*functions.Overload
		expect     func(t *testing.T, activation LateBindActivation, err error)
	}{
		{
			name:       "OK_No_Overloads",
			activation: actEmpty,
			overloads:  nil,
			expect:     expectActivation(actEmpty),
		},
		{
			name:       "OK_Happy_Path",
			activation: actHierarchical,
			overloads: []*functions.Overload{
				f1_string_string,
				f1_string_string_string,
				f1_varargs_string,
			},
			expect: expectActivation(actHierarchical, f1_string_string, f1_string_string_string, f1_varargs_string),
		},
		{
			name:       "ERROR_Activation_Nil",
			activation: nil,
			overloads: []*functions.Overload{
				{
					Operator: "f2",
					Function: func(values ...ref.Val) ref.Val {
						return types.String("f2")
					},
					NonStrict: false,
				},
			},
			expect: expectError(errorNilActivation),
		},
		{
			name:       "ERROR_Duplicate_Overloads",
			activation: &mapActivation{},
			overloads: []*functions.Overload{
				f2_string_string,
				f2_string,
				f2_string_string,
			},
			expect: expectError(fmt.Sprintf("overload already exists '%s'", "f2_string_string")),
		},
	}

	for _, testCase := range testCases {

		t.Run(testCase.name, func(t *testing.T) {

			actual, err := NewLateBindActivation(testCase.activation, testCase.overloads...)
			testCase.expect(t, actual, err)
		})
	}
}

// TestLateBindActivationParent verifies the implementation of lateBindActivation.Parent(). The
// expectation is for the function to return the Activation implementation configured with the
// vars field of the activation.
func TestLateBindActivatioParent(t *testing.T) {

	actHierarchical := &hierarchicalActivation{
		parent: &mapActivation{
			bindings: map[string]any{
				"a": 5,
				"b": 10,
			},
		},
		child: &mapActivation{
			bindings: map[string]any{
				"a": 4,
				"c": 22,
			},
		},
	}

	testCases := []struct {
		name       string
		activation func() *lateBindActivation
		expect     func(t *testing.T, actual Activation)
	}{
		// NOTE: this test is implemented for completeness but unless
		//       we have access to the private type there is no way to
		//       produce a nil parent.
		{
			name: "OK_Nil_Parent",
			activation: func() *lateBindActivation {
				return &lateBindActivation{
					vars: nil,
					dispatcher: &defaultDispatcher{
						parent:    nil,
						overloads: overloadMap{},
					},
				}
			},
			expect: func(t *testing.T, actual Activation) {

				if actual != nil {
					t.Error("expected nil parent.")
				}
			},
		},
		{
			name: "OK_Non_Nil_Parent",
			activation: func() *lateBindActivation {
				return &lateBindActivation{
					vars: actHierarchical,
					dispatcher: &defaultDispatcher{
						parent:    nil,
						overloads: overloadMap{},
					},
				}
			},
			expect: func(t *testing.T, actual Activation) {
				if actual != actHierarchical {
					t.Errorf("unexpected parent (got: %v, want: %v)", actual, actHierarchical)
				}
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			candidate := testCase.activation()
			actual := candidate.Parent()
			testCase.expect(t, actual)
		})
	}
}

// TestLateBindActivationResolveName verifies the implemented behaviour of
// lateBindActivation.ResolveName(string). The expectation is for the function
// to defer all the name resolution to the configured activation.
func TestLateBindActivationResolveName(t *testing.T) {

	activation := func() Activation {
		return &hierarchicalActivation{
			parent: &mapActivation{
				bindings: map[string]any{
					"a": 5,
					"b": 10,
				},
			},
			child: &mapActivation{
				bindings: map[string]any{
					"a": 4,
					"c": 22,
				},
			},
		}
	}

	testCases := []struct {
		name     string
		vars     Activation
		varName  string
		found    bool
		expected any
	}{
		{
			name:     "TRUE_Single_Name_Occurrence",
			vars:     activation(),
			varName:  "c",
			found:    true,
			expected: 22,
		},
		{
			name:     "TRUE_Multiple_Name_Occurrences",
			vars:     activation(),
			varName:  "a",
			found:    true,
			expected: 4,
		},
		{
			name:     "FALSE_Missing_Name",
			vars:     activation(),
			varName:  "d",
			found:    false,
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			candidate := &lateBindActivation{
				vars:       testCase.vars,
				dispatcher: NewDispatcher(),
			}

			actual, found := candidate.ResolveName(testCase.varName)
			if testCase.found != found {
				t.Errorf("found mistmatch for var (name: '%s', got: %v, want: %v", testCase.varName, found, testCase.found)
			}

			if testCase.expected != actual {
				t.Errorf("value mismatch for var (name: '%s', got: %v, want: %v)", testCase.varName, actual, testCase.expected)
			}
		})
	}
}

// TestLateBindActivationResolveOverload verifies the implemented behaviour of
// lateBindActivation.ResolveOverload(string). The expectation is for the function
// to resolve the overload that is mapped to the given overload identifier if this
// is present. The resolution rules are as follows:
//
//   - the overload is first searched in the dispatcher associated to the instance.
//   - if a non nil function overload is found, it is returned.
//   - if a nil overload is found, the search continues by inspecting the activation
//     bound to the instance.
//   - if the activation bound to the instance is an empty activation the search is
//     complete and nil is returned.
//   - if the activation bound to the instance is a mapActivation the search is complete
//     and nil is returned.
//   - if the activation bound to the instance is a hierarchical activation, first the
//     child is searched to determine whether there is a LateBindActivation implementation
//     in the tree that originates from the parent.
//   - if the child search returns a nil overload, the parent is searched to determine
//     whether there is a LateBindActivation implementation in the tree that originates
//     from the child.
//   - if a LateBindActivation implementation is found, the ResolveOverload(string) name
//     is invoked to repeat the search detailed in the previous step.
//
// If the activation tree is exhausted and no overload is found matching the given
// identifier, nil is returned.
func TestLateBindActivationResolveOverload(t *testing.T) {

	nestedActivation, overloads := prepareNestedActivation()

	testCases := []struct {
		name       string
		candidate  func() *lateBindActivation
		overloadId string
		expected   *functions.Overload
	}{
		{
			name: "TRUE_Simple_Case",
			candidate: func() *lateBindActivation {

				return &lateBindActivation{
					vars: &emptyActivation{},
					dispatcher: &defaultDispatcher{
						parent: nil,
						overloads: overloadMap{
							"f1_string":        overloads["f1_string"],
							"f1_string_string": overloads["f1_string_string"],
						},
					},
				}
			},
			overloadId: "f1_string",
			expected:   overloads["f1_string"],
		}, {
			name: "FALSE_Simple_Case",
			candidate: func() *lateBindActivation {

				return &lateBindActivation{
					vars: &emptyActivation{},
					dispatcher: &defaultDispatcher{
						parent: nil,
						overloads: overloadMap{
							"f1_string":        overloads["f1_string"],
							"f1_string_string": overloads["f1_string_string"],
						},
					},
				}
			},
			overloadId: "f1_string_string_string",
			expected:   nil,
		}, {
			name:       "TRUE_Complex_Case_With_Nesting_Top_Level",
			candidate:  nestedActivation,
			overloadId: "f1_string",
			expected:   overloads["f1_string"],
		}, {
			name:       "TRUE_Complex_Case_With_Nesting_Top_Level_Parent",
			candidate:  nestedActivation,
			overloadId: "f2_string",
			expected:   overloads["f2_string_parent"],
		}, {
			name:       "TRUE_Complex_Case_With_Nesting_Top_Level_Shadows_Vars_Parent",
			candidate:  nestedActivation,
			overloadId: "f3_string",
			expected:   overloads["f3_string"],
		}, {
			name:       "TRUE_Complex_Case_With_Nesting_Top_Level_Shadows_Vars_Child",
			candidate:  nestedActivation,
			overloadId: "f4_string",
			expected:   overloads["f4_string"],
		}, {
			name:       "TRUE_Complex_Case_With_Nesting_Vars_Child_Shadows_Parent",
			candidate:  nestedActivation,
			overloadId: "f5_string",
			expected:   overloads["f5_string_nested_child"],
		}, {
			name:       "TRUE_Complex_Case_With_Nexting_Vars_Child_Only_Find",
			candidate:  nestedActivation,
			overloadId: "f6_string",
			expected:   overloads["f6_string_nested_child"],
		}, {
			name:       "TRUE_Complex_Case_With_Nesting_Vars_Parent_Only_Find",
			candidate:  nestedActivation,
			overloadId: "f7_string",
			expected:   overloads["f7_string_nested_parent"],
		}, {
			name:       "FALSE_Complex_Case_With_Nesting_Missing",
			candidate:  nestedActivation,
			overloadId: "f8_string",
			expected:   nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			activation := testCase.candidate()
			actual := activation.ResolveOverload(testCase.overloadId)
			if actual != testCase.expected {
				t.Errorf("mismatch function for overload (id: %s, nil: %v)", testCase.overloadId, actual == nil)
			}
		})
	}

}

// TestLateBindActivationResolverOverloads verifies the implemented behaviour of
// latebindActivation.ResolveOverloads(). The expectation is for the function to
// generate a dispatcher that aggregates all the function overloads definition by
// following the precedence rules implemented for ResolveOverload(string) when
// duplicates are encountered.
func TestLateBindActivationResolveOverloads(t *testing.T) {

	nestedActivation, overloads := prepareNestedActivation()

	testCases := []struct {
		name      string
		candidate func() *lateBindActivation
		expected  Dispatcher
	}{
		{
			name: "OK_Simple_Activation_Empty",
			candidate: func() *lateBindActivation {
				return &lateBindActivation{
					vars: &mapActivation{},
					dispatcher: &defaultDispatcher{
						parent:    nil,
						overloads: overloadMap{},
					},
				}
			},
			expected: &defaultDispatcher{
				parent:    nil,
				overloads: overloadMap{},
			},
		},
		{
			name: "OK_Simple_Activation_Not_Empty",
			candidate: func() *lateBindActivation {
				return &lateBindActivation{
					vars: &mapActivation{},
					dispatcher: &defaultDispatcher{
						parent: &defaultDispatcher{
							parent: nil,
							overloads: overloadMap{
								"f2_string": overloads["f2_string_parent"],
							},
						},
						overloads: overloadMap{
							"f1_string": overloads["f1_string"],
							"f3_string": overloads["f3_string"],
						},
					},
				}
			},
			expected: &defaultDispatcher{
				parent: nil,
				overloads: overloadMap{
					"f1_string": overloads["f1_string"],
					"f2_string": overloads["f2_string_parent"],
					"f3_string": overloads["f3_string"],
				},
			},
		},
		{
			name:      "OK_Nested_Activation",
			candidate: nestedActivation,
			expected: &defaultDispatcher{
				parent: nil,
				overloads: overloadMap{
					"f1_string":               overloads["f1_string"],
					"f1_string_string":        overloads["f1_string_string"],
					"f1_string_string_string": overloads["f1_string_string_string"],
					"f2_string":               overloads["f2_string_parent"],
					"f3_string":               overloads["f3_string"],
					"f4_string":               overloads["f4_string"],
					"f5_string":               overloads["f5_string_nested_child"],
					"f6_string":               overloads["f6_string_nested_child"],
					"f7_string":               overloads["f7_string_nested_parent"],
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			activation := testCase.candidate()
			actual := activation.ResolveOverloads()

			if actual == nil {
				t.Fatal("unexpected nil reference returned by ResolveOverloads")
			}

			expectedIds := testCase.expected.OverloadIds()
			actualIds := actual.OverloadIds()

			if len(expectedIds) != len(actualIds) {
				t.Errorf("number of overloads mismatch (got: %d, want: %d)", len(actualIds), len(expectedIds))
			}

			for _, ovlId := range expectedIds {

				expectedOverload, found := testCase.expected.FindOverload(ovlId)
				if !found {
					t.Fatalf("unexpected: overload (id: %s) declared but not found", ovlId)
				}
				actualOverload, found := actual.FindOverload(ovlId)
				if !found {
					t.Errorf("overload (id: %s) not found in result", ovlId)
				}
				if actualOverload == nil {
					t.Errorf("overload (id: %s) is found, but nil", ovlId)
				}
				if expectedOverload != actualOverload {
					t.Errorf("overload (id: %s) mismatch (got: %v, want: %v)", ovlId, actualOverload, expectedOverload)
				}
			}
		})
	}

}

// TestValidateOverloads verifies the implemented behaviour of ValidateOverloads.
// the expectation is for the function to resolve all function overload definitions
// that are exposed by the given activation and compare them with the reference
// dispatcher maintaining the set of functions statically bound during expression
// parsing. For any function overload that is redefined in the activation the
// expectation is for the two function signatures to match and the associated
// overload attributes to be the same.
func TestValidateOverloads(t *testing.T) {

	f2_string_string := unary("f2_string_string", 0, true, func(arg ref.Val) ref.Val {
		return types.String("f2_string_string")
	})
	f2_string_string_string := binary("f2_string_string_string", 0, true, func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.String("f2_string_string_string")
	})
	f2_varargs_string := function("f2_varargs_string", 0, true, func(args ...ref.Val) ref.Val {
		return types.String("f2_varargs_string")
	})

	f1_string_string := unary("f1_string_string", 0, true, func(arg ref.Val) ref.Val {
		return types.String("f1_string_string")
	})
	f1_string_string_string := binary("f1_string_string_string", 0, true, func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.String("f1_string")
	})
	f1_varargs_string := function("f1_varargs_string", 0, true, func(args ...ref.Val) ref.Val {
		return types.String("f1_varargs_string")
	})

	// matching overloads

	f1_string_string_overload := unary("f1_string_string", 0, true, func(arg ref.Val) ref.Val {
		return types.String("f1_string_string_overload")
	})
	f1_string_string_string_overload := binary("f1_string_string_string", 0, true, func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.String("f1_string_overload")
	})
	f1_varargs_string_overload := function("f1_varargs_string", 0, true, func(args ...ref.Val) ref.Val {
		return types.String("f1_varargs_string_overload")
	})

	f3_string_string := unary("f3_string_string", 0, true, func(arg ref.Val) ref.Val {
		return types.String("f3_string_string")
	})

	// mismatched overloads

	f1_string_string_binary := binary("f1_string_string", 0, true, func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.String("f1_string_string_binary")
	})

	f1_string_string_varargs := function("f1_string_string", 0, true, func(args ...ref.Val) ref.Val {
		return types.String("f1_string_string_binary")
	})

	f1_string_string_string_unary := unary("f1_string_string_string", 0, true, func(arg ref.Val) ref.Val {
		return types.String("f1_string_string_string_unary")
	})
	f1_string_string_string_varargs := function("f1_string_string_string", 0, true, func(args ...ref.Val) ref.Val {
		return types.String("f1_string_string_string_varargs")
	})

	f1_varargs_string_unary := unary("f1_varargs_string", 0, true, func(arg ref.Val) ref.Val {
		return types.String("f1_varargs_string_unary")
	})
	f1_varargs_string_binary := binary("f1_varargs_string", 0, true, func(lhs ref.Val, rhs ref.Val) ref.Val {
		return types.String("f1_varargs_string_unary")
	})

	f1_string_operand_trait := unary("f1_string_string", 2, true, func(arg ref.Val) ref.Val {
		return types.String("f1_string_operand_trait")
	})

	f1_string_non_strict := unary("f1_string_string", 0, false, func(arg ref.Val) ref.Val {
		return types.String("f1_string_non_strict")
	})

	// referenceDispatcher generates a Dispatcher implementation that contains
	// a baseline configuration for all the function overloads that are defined.
	referenceDispatcher := func() Dispatcher {

		return &defaultDispatcher{
			parent: nil,
			overloads: overloadMap{
				"f1_string_string":        f1_string_string,
				"f1_string_string_string": f1_string_string_string,
				"f1_varargs_string":       f1_varargs_string,
				"f2_string_string":        f2_string_string,
				"f2_string_string_string": f2_string_string_string,
				"f2_varargs_string":       f2_varargs_string,
			},
		}
	}

	// candidateActivation generates a function that produces a repeatable configuration
	// for the Activation implementation based on lateBindActivation. When supplied with
	// a non-empty overload identifier, it also maps the given overload to such identifier.
	//
	// NOTE: these parameters are used to produce alterations over a baseline and produce
	//       different test cases, to ensure that all branches are explored.
	candidateActivation := func(overloadId string, overload *functions.Overload) func() Activation {

		return func() Activation {

			dispatcher := &defaultDispatcher{
				parent: nil,
				overloads: overloadMap{
					"f1_string_string":        f1_string_string_overload,
					"f1_string_string_string": f1_string_string_string_overload,
					"f1_varargs_string":       f1_varargs_string_overload,
					"f3_string_string":        f3_string_string,
				},
			}
			if len(overloadId) > 0 {
				dispatcher.overloads[overloadId] = overload
			}
			return &lateBindActivation{
				vars:       &mapActivation{},
				dispatcher: dispatcher,
			}
		}
	}

	testCases := []struct {
		name      string
		reference Dispatcher
		candidate func() Activation
		err       error
	}{
		{
			name:      "OK_Matching_Overloads",
			reference: referenceDispatcher(),
			candidate: candidateActivation("", nil),
			err:       nil,
		},

		// binary function is unmatched in signature
		{
			name:      "ERROR_Unary_Mismatch_Binary",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string", f1_string_string_binary),
			err:       fmt.Errorf(errorInvalidSignature, "f1_string_string", binarySignature, unarySignature),
		}, {
			name:      "ERROR_Unary_Mismatch_VarArgs",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string", f1_string_string_varargs),
			err:       fmt.Errorf(errorInvalidSignature, "f1_string_string", functionSignature, unarySignature),
		}, {
			name:      "ERROR_Unary_Mismatch_Nil",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string", &functions.Overload{
				Operator:     "f1_string_string",
				OperandTrait: 0,
				NonStrict:    true,
			}),
			err: fmt.Errorf(errorInvalidSignature, "f1_string_string", "<nil>", unarySignature),
		},

		// binary function is unmatched in signature
		{
			name:      "ERROR_Binary_Mismatch_Unary",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string_string", f1_string_string_string_unary),
			err:       fmt.Errorf(errorInvalidSignature, "f1_string_string_string", unarySignature, binarySignature),
		}, {
			name:      "ERROR_Binary_Mismatch_VarArgs",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string_string", f1_string_string_string_varargs),
			err:       fmt.Errorf(errorInvalidSignature, "f1_string_string_string", functionSignature, binarySignature),
		}, {
			name:      "ERROR_Binary_Mismatch_Nil",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string_string", &functions.Overload{
				Operator:     "f1_string_string_string",
				OperandTrait: 0,
				NonStrict:    true,
			}),
			err: fmt.Errorf(errorInvalidSignature, "f1_string_string_string", "<nil>", binarySignature),
		},

		// varargs function is unmatched in signature
		{
			name:      "ERROR_VarArgs_Mismatch_Unary",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_varargs_string", f1_varargs_string_unary),
			err:       fmt.Errorf(errorInvalidSignature, "f1_varargs_string", unarySignature, functionSignature),
		}, {
			name:      "ERROR_VarArgs_Mismatch_Binary",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_varargs_string", f1_varargs_string_binary),
			err:       fmt.Errorf(errorInvalidSignature, "f1_varargs_string", binarySignature, functionSignature),
		}, {
			name:      "ERROR_VarArgs_Mismatch_Nil",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_varargs_string", &functions.Overload{
				Operator:     "f1_varargs_string",
				OperandTrait: 0,
				NonStrict:    true,
			}),
			err: fmt.Errorf(errorInvalidSignature, "f1_varargs_string", "<nil>", functionSignature),
		},

		// unmatched attributes
		{
			name:      "ERROR_Mismatch_OperandTrait",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string", f1_string_operand_trait),
			err:       fmt.Errorf(errorMismatch, "f1_string_string", "OperandTrait", 2, 0),
		}, {
			name:      "ERROR_Mismatch_NonStrict",
			reference: referenceDispatcher(),
			candidate: candidateActivation("f1_string_string", f1_string_non_strict),
			err:       fmt.Errorf(errorMismatch, "f1_string_string", "NonStrict", false, true),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			activation := testCase.candidate()
			err := ValidateOverloads(testCase.reference, activation)

			// if we expect an error
			if testCase.err != nil {
				if err == nil {
					// it should not be nil
					t.Errorf("outcome mismatch (got: <nil>,  want: %s)", testCase.err.Error())
				} else if testCase.err.Error() != err.Error() {

					// it should match
					t.Errorf("outcome mismatch (got: %s want: %s)", err.Error(), testCase.err.Error())
				}
			} else if err != nil {

				t.Errorf("outcome mismatch (got: %s, want: <nil>)", err.Error())
			}
		})
	}
}

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

// TestLateBindingDecorator verifies the implemented behaviour of decLateBind().
// The expectation is for the function to return an InterepretableDecorator that
// applies the late bind behaviour.
func TestLateBindingDecorator(t *testing.T) {

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

			candidate := decLateBinding()
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

// expectValue is a convenience function that produces an expectation function that checks
// that the value returned by the execution of Interpretable.Eval(Activation) matches the
// given expected value otherwise it fails the test.
func expectValue(expected ref.Val) func(t *testing.T, target Interpretable, actual ref.Val) {

	return func(t *testing.T, _ Interpretable, actual ref.Val) {

		t.Helper()

		if expected.Equal(actual) == types.False {
			t.Errorf("unexpected value (got: %v, want: %v)", actual, expected)
		}
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

// expectType generates an expectation function that is used to validate that the name
// of Interpretable expression is the same of the name of the given type.
func expectType(expected Interpretable) func(t *testing.T, target Interpretable, actual ref.Val) {

	return func(t *testing.T, target Interpretable, _ ref.Val) {

		t.Helper()

		expectedType := reflect.TypeOf(expected).Name()
		actualType := reflect.TypeOf(target).Name()

		if expectedType != actualType {
			t.Errorf("unexpected type: (got: %s, want: %s)", actualType, expectedType)
		}
	}
}

// chain is a convenience function that can be used to run in sequence multiple expectation
// functions that are passed as argument. The function returns an expectation function that
// executes all the functions in the checks array, with the actual arguments passed to the
// function.
func chain(checks ...func(t *testing.T, target Interpretable, actual ref.Val)) func(t *testing.T, target Interpretable, actual ref.Val) {

	return func(t *testing.T, target Interpretable, actual ref.Val) {

		for _, check := range checks {
			check(t, target, actual)
		}
	}
}

// prepareNestedActivation generates a map of overloads and a function that produces a
// lateBindActivation reference which holds a tree of activations with implementations
// of LateBindActivation in the tree. The resulting activation is as structured as shown
// below:
//
// lateBindActivation:
//
//	├─ vars ---> hierarchicalActivation:
//	│  ├─ parent ---> emptyActivation
//	│  └─ child  ---> hierarchicalActivation:
//	│                 ├─ parent ---> lateBindActivation:
//	│                 │               ├─ vars: mapActivation,
//	│                 │               └─ dispatcher: defaultDispatcher
//	│                 │                               ├─ parent: nil
//	│                 │                               └─ overloads:
//	│                 │                                   ├─ "f3_string" --> f3_string_nested_parent
//	│                 │                                   ├─ "f5_string" --> f5_string_nested_parent
//	│                 │                                   └─ "f7_string" --> f7_string_nested_parent
//	│                 └─ child ---> lateBindActivation:
//	│                                ├─ vars: mapActivation,
//	│                                └─  dispatcher: defaultDispatcher
//	│                                                 ├─ parent: nil
//	│                                                 └─ overloads
//	│                                                     ├─ "f3_string" --> f3_string_nested_child
//	│                                                     ├─ "f4_string" --> f5_string_nested_child
//	│                                                     └─ "f6_string" --> f7_string_nested_child
//	└─ dispatcher: defaultDispatcher:
//	                ├─ parent: defaultDispatcher:
//	                │           ├─ parent: nil
//	                │           └─ overloads:
//	                │               ├─ "f1_string": f1_string_parent
//	                │               └─ "f2_string": f2_string_parent
//	                └─ overloads:
//	                    ├─ "f1_string": f1_string
//	                    ├─ "f1_string_string": f1_string_string
//	                    ├─ "f1_string_string_string": f1_string_stirng_string
//	                    ├─ "f3_string": f3_string
//	                    └─ "f4_string": f4_string
func prepareNestedActivation() (func() *lateBindActivation, map[string]*functions.Overload) {

	f1_string := function("f1_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f1_string") })
	f3_string := function("f3_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f3_string") })
	f4_string := function("f4_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f4_string") })

	// this function creates an upper case version of the original string passed as
	// argument (see: TestLateBindEvalUnaryEval).
	f1_string_string := unary("f1_string_string", 0, false, func(arg ref.Val) ref.Val {
		text, _ := arg.(types.String)
		return types.String(strings.ToUpper(string(text)))
	})

	// this function composes the two strings passed as arguments in inverse order and with
	// a space in the middle (see TestLateBindEvalBinaryEval).
	f1_string_string_string := binary("f1_string_string_string", 0, false, func(lhs ref.Val, rhs ref.Val) ref.Val {

		a, _ := lhs.(types.String)
		b, _ := rhs.(types.String)

		return b.Add(types.String(" ")).(types.String).Add(a)
	})
	f1_string_parent := function("f1_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f1_string_parent") })
	f2_string_parent := function("f2_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f2_string_parent") })

	f3_string_nested_parent := function("f3_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f3_string_nested_parent") })
	f5_string_nested_parent := function("f5_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f5_string_nested_parent") })
	f7_string_nested_parent := function("f7_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f7_string_nested_parent") })

	f4_string_nested_child := function("f4_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f4_string_nested_child") })
	f5_string_nested_child := function("f5_string", 0, false, func(args ...ref.Val) ref.Val { return types.String("f5_string_nested_child") })
	f6_string_nested_child := function("f6_string", 0, false, func(args ...ref.Val) ref.Val {

		var result traits.Adder = types.String("")

		for _, arg := range args {
			text, _ := arg.(types.String)
			result = result.Add(text).(traits.Adder)
		}
		return result.(ref.Val)
	})

	overloads := map[string]*functions.Overload{
		"f1_string":               f1_string,
		"f1_string_string":        f1_string_string,
		"f1_string_string_string": f1_string_string_string,
		"f1_string_parent":        f1_string_parent,
		"f2_string_parent":        f2_string_parent,
		"f3_string":               f3_string,
		"f3_string_nested_parent": f3_string_nested_parent,
		"f4_string":               f4_string,
		"f4_string_nested_child":  f4_string_nested_child,
		"f5_string_nested_child":  f5_string_nested_child,
		"f5_string_nested_parent": f5_string_nested_parent,
		"f6_string_nested_child":  f6_string_nested_child,
		"f7_string_nested_parent": f7_string_nested_parent,
	}

	nestedActivation := func() *lateBindActivation {

		return &lateBindActivation{
			vars: &hierarchicalActivation{
				parent: &emptyActivation{},
				child: &hierarchicalActivation{
					parent: &lateBindActivation{
						vars: &mapActivation{},
						dispatcher: &defaultDispatcher{
							parent: nil,
							overloads: overloadMap{
								"f3_string": f3_string_nested_parent,
								"f5_string": f5_string_nested_parent,
								"f7_string": f7_string_nested_parent,
							},
						},
					},
					child: &lateBindActivation{
						vars: &mapActivation{},
						dispatcher: &defaultDispatcher{
							parent: nil,
							overloads: overloadMap{
								"f4_string": f4_string_nested_child,
								"f5_string": f5_string_nested_child,
								"f6_string": f6_string_nested_child,
							},
						},
					},
				},
			},
			dispatcher: &defaultDispatcher{
				parent: &defaultDispatcher{
					parent: nil,
					overloads: overloadMap{
						"f1_string": f1_string_parent,
						"f2_string": f2_string_parent,
					},
				},
				overloads: overloadMap{
					"f1_string":               f1_string,
					"f1_string_string":        f1_string_string,
					"f1_string_string_string": f1_string_string_string,
					"f3_string":               f3_string,
					"f4_string":               f4_string,
				},
			},
		}
	}

	return nestedActivation, overloads

}

// unary is a convenience function to produce a reference to functions.Overload configured with
// the given parameters. This function sets the Overload.Unary to the given function.
func unary(operator string, operandTrait int, nonStrict bool, function functions.UnaryOp) *functions.Overload {

	return &functions.Overload{
		Operator:     operator,
		OperandTrait: operandTrait,
		NonStrict:    nonStrict,
		Unary:        function,
	}
}

// binary is a convenience function to produce a reference to functions.Overload configured with
// the given parameters. This function sets the Overload.Binary to the given function.
func binary(operator string, operandTrait int, nonStrict bool, function functions.BinaryOp) *functions.Overload {

	return &functions.Overload{
		Operator:     operator,
		OperandTrait: operandTrait,
		NonStrict:    nonStrict,
		Binary:       function,
	}
}

// function is a convenience function to produce a reference to functions.Overload configured with
// the given parameters. This function sets the Overload.Function to the given function.
func function(operator string, operandTrait int, nonStrict bool, function functions.FunctionOp) *functions.Overload {

	return &functions.Overload{
		Operator:     operator,
		OperandTrait: operandTrait,
		NonStrict:    nonStrict,
		Function:     function,
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
