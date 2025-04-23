package interpreter

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
)

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
				t.Error("unexpected nil reference returned by ResolveOverloads")
				t.FailNow()
			}

			expectedIds := testCase.expected.OverloadIds()
			actualIds := actual.OverloadIds()

			if len(expectedIds) != len(actualIds) {
				t.Errorf("number of overloads mismatch (got: %d, want: %d)", len(actualIds), len(expectedIds))
			}

			for _, ovlId := range expectedIds {

				expectedOverload, found := testCase.expected.FindOverload(ovlId)
				if !found {
					t.Errorf("unexpected: overload (id: %s) declared but not found", ovlId)
					t.FailNow()
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
			t.Errorf("unexpected type in test case (got: %T, want: %T)", target, &lateBindEvalZeroArity{})
			t.FailNow()
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
			t.Errorf("unexpected type in test case (got: %T, want: %T)", target, &lateBindEvalUnary{})
			t.FailNow()
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
			t.Errorf("unexpected type in test case (got: %T, want: %T)", target, &lateBindEvalBinary{})
			t.FailNow()
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
		t.Errorf("args size mismatch (got: %d, want: %d)", len(aArgs), len(eArgs))
		t.FailNow()
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
