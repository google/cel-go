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

// The interpreter package provides functions to evaluate CEL programs against
// a series of inputs and functions supplied at runtime.
package interpreter

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/interpreter/types"
	"github.com/google/cel-go/interpreter/types/aspects"
	"github.com/google/cel-go/interpreter/types/providers"
	"github.com/google/cel-go/interpreter/functions"
)

// Interpreter generates a new Interpretable from a Program.
type Interpreter interface {
	// NewInterpretable returns an Interpretable from a Program.
	NewInterpretable(program Program) Interpretable
}

// Interpretable can accept a given Activation and produce a result along with
// an accompanying EvalState which can be used to inspect whether additional
// data might be necessary to complete the evaluation.
type Interpretable interface {
	// Eval an Activation to produce an output and EvalState.
	Eval(activation Activation) (interface{}, EvalState)
}

type exprInterpreter struct {
	dispatcher   Dispatcher
	typeProvider providers.TypeProvider
}

var _ Interpreter = &exprInterpreter{}

// NewInterpreter builds an Interpreter from a Dispatcher and TypeProvider
// which will be used throughout the Eval of all Interpretable instances
// gerenated from it.
func NewInterpreter(dispatcher Dispatcher, typeProvider providers.TypeProvider) *exprInterpreter {
	return &exprInterpreter{dispatcher, typeProvider}
}

// StandardInterpreter builds a Dispatcher and TypeProvider with support
// for all of the CEL builtins defined in the language definition.
func StandardIntepreter(types ...proto.Message) *exprInterpreter {
	dispatcher := NewDispatcher()
	dispatcher.Add(functions.StandardBuiltins()...)
	typeProvider := providers.NewTypeProvider(types...)
	return NewInterpreter(dispatcher, typeProvider)
}

func (i *exprInterpreter) NewInterpretable(program Program) Interpretable {
	// program needs to be pruned with the TypeProvider
	return &exprInterpretable{
		interpreter: i,
		program:     program,
		state:       NewEvalState()}
}

type exprInterpretable struct {
	interpreter *exprInterpreter
	program     Program
	state       MutableEvalState
}

var _ Interpretable = &exprInterpretable{}

func (i *exprInterpretable) Eval(activation Activation) (interface{}, EvalState) {
	// register machine-like evaluation of the program with the given activation.
	currActivation := activation
	stepper := i.program.Instructions()
	var resultId int64
	for step, hasNext := stepper.Next(); hasNext; step, hasNext = stepper.Next() {
		resultId = step.GetId()
		switch step.(type) {
		case *ConstExpr:
			i.evalConst(step)
		case *IdentExpr:
			i.evalIdent(step, currActivation)
		case *SelectExpr:
			i.evalSelect(step, currActivation)
		case *CallExpr:
			i.evalCall(step, currActivation)
		case *CreateListExpr:
			i.evalCreateList(step)
		case *CreateMapExpr:
			i.evalCreateMap(step)
		case *CreateObjectExpr:
			i.evalCreateType(step)
		case *MovInst:
			i.evalMov(step)
			// Special instruction for modifying the program cursor
		case *JumpInst:
			jmpExpr := step.(*JumpInst)
			// TODO: Add test for whether the jump should be made based
			// on nil as the jump value, equality as the value, or unknown/err
			if jmpExpr.OnValue == nil || i.value(resultId) == jmpExpr.OnValue {
				if !stepper.JumpCount(jmpExpr.Count) {
					// TODO: Error, the jump count should never exceed the
					// program length.
					panic("jumped too far")
				}
			}
			// Special instructions for modifying the activation stack
		case *PushScopeInst:
			pushScope := step.(*PushScopeInst)
			scopeDecls := pushScope.Declarations
			childActivaton := make(map[string]interface{})
			for key, declId := range scopeDecls {
				childActivaton[key] = func() interface{} {
					return i.value(declId)
				}
			}
			currActivation = NewHierarchicalActivation(currActivation, NewActivation(childActivaton))
		case *PopScopeInst:
			currActivation = currActivation.Parent()
		}
	}
	return i.value(resultId), i.state
}

func (i *exprInterpretable) evalConst(step Instruction) {
	i.setValue(step.GetId(), step.(*ConstExpr).Value)
}

func (i *exprInterpretable) evalIdent(step Instruction, currActivation Activation) {
	idExpr := step.(*IdentExpr)
	// TODO: Refactor this code for sharing.
	if result, found := currActivation.ResolveName(idExpr.Name); found {
		i.setValue(step.GetId(), result)
	} else if enum, found := i.interpreter.typeProvider.EnumValue(idExpr.Name); found {
		i.setValue(step.GetId(), enum)
	} else {
		i.setValue(step.GetId(), &UnknownValue{[]Instruction{step}})
	}
}

func (i *exprInterpretable) evalSelect(step Instruction, currActivation Activation) {
	selExpr := step.(*SelectExpr)
	operand := i.value(selExpr.Operand)
	if unknown, ok := operand.(*UnknownValue); ok {
		i.resolveUnknown(unknown, selExpr, currActivation)
	} else if indexer, ok := operand.(aspects.Indexer); ok {
		if fieldValue, err := indexer.Get(selExpr.Field); err == nil {
			i.setValue(step.GetId(), fieldValue)
		} else {
			i.setValue(step.GetId(), err)
		}
	} else {
		// determine whether the operand was unknown or just the wrong type
		i.setValue(step.GetId(),
			ErrorValue{[]error{fmt.Errorf("invalid operand in select")}})
	}
}

// resolveUnknown attempts to resolve a qualified name from a select expression
// which may have generated unknown values during the course of execution if
// the expression was not type-checked and the select, in fact, refers to a
// qualified identifier name instead of a series of field selections.
func (i *exprInterpretable) resolveUnknown(unknown *UnknownValue,
	selExpr *SelectExpr,
	currActivation Activation) {
	if object, found := currActivation.ResolveReference(selExpr.Id); found {
		i.setValue(selExpr.Id, object)
	} else {
		var validIdent = true
		var identifier = selExpr.Field
		for _, arg := range unknown.Args {
			switch arg.(type) {
			case *IdentExpr:
				identifier = arg.(*IdentExpr).Name + "." + identifier
			case *SelectExpr:
				identifier = arg.(*SelectExpr).Field + "." + identifier
			default:
				argVal := i.value(arg.GetId())
				if argStr, ok := argVal.(string); ok {
					identifier = argStr + "." + identifier
				} else {
					validIdent = false
					break
				}
			}
		}
		if validIdent {
			if i.program.Container() != "" {
				identifier = i.program.Container() + "." + identifier
			}
			if object, found := currActivation.ResolveName(identifier); found {
				i.setValue(selExpr.Id, object)
			} else if enum, found := i.interpreter.typeProvider.EnumValue(identifier); found {
				i.setValue(selExpr.Id, enum)
			} else {
				i.setValue(selExpr.Id,
					&UnknownValue{append([]Instruction{selExpr}, unknown.Args...)})
			}
		}
	}
}

func (i *exprInterpretable) evalCall(step Instruction, currActivation Activation) {
	callExpr := step.(*CallExpr)
	argVals := make([]interface{}, len(callExpr.Args), len(callExpr.Args))
	for idx, argId := range callExpr.Args {
		argVals[idx] = i.value(argId)
	}
	ctx := &CallContext{
		call:       step.(*CallExpr),
		activation: currActivation,
		args:       argVals,
		metadata:   i.program.Metadata()}
	if result, err := i.interpreter.dispatcher.Dispatch(ctx); err == nil {
		i.setValue(step.GetId(), result)
	} else {
		i.setValue(step.GetId(), err)
	}
}

func (i *exprInterpretable) evalCreateList(step Instruction) {
	listExpr := step.(*CreateListExpr)
	elements := make([]interface{}, len(listExpr.Elements))
	for idx, elementId := range listExpr.Elements {
		elements[idx] = i.value(elementId)
	}
	// TODO: Add an error state for the list if any element is an error
	adaptingList := types.NewListValue(elements)
	i.setValue(step.GetId(), adaptingList)
}

func (i *exprInterpretable) evalCreateMap(step Instruction) {
	mapExpr := step.(*CreateMapExpr)
	entries := make(map[interface{}]interface{})
	for keyId, valueId := range mapExpr.KeyValues {
		entries[i.value(keyId)] = i.value(valueId)
	}
	// TODO: Add an error state if any key is repeated and any element in the
	// map (key or value) is an error.
	adaptingMap := types.NewMapValue(entries)
	i.setValue(step.GetId(), adaptingMap)
}

func (i *exprInterpretable) evalCreateType(step Instruction) {
	typeExpr := step.(*CreateObjectExpr)
	fields := make(map[string]interface{})
	for field, valueId := range typeExpr.FieldValues {
		fields[field] = i.value(valueId)
	}
	if value, err := i.newValue(typeExpr.Name, fields); err == nil {
		i.setValue(step.GetId(), value)
	} else {
		i.setValue(step.GetId(), err)
	}
}

func (i *exprInterpretable) evalMov(step Instruction) {
	movExpr := step.(*MovInst)
	i.setValue(movExpr.ToExprId, i.value(step.GetId()))
}

func (i *exprInterpretable) value(id int64) interface{} {
	if object, found := i.state.Value(id); found {
		return object
	}
	return nil
}

func (i *exprInterpretable) setValue(id int64, value interface{}) {
	i.state.SetValue(id, value)
}

func (i *exprInterpretable) newValue(typeName string,
	fields map[string]interface{}) (types.ObjectValue, error) {
	return i.interpreter.typeProvider.NewValue(typeName, fields)
}
