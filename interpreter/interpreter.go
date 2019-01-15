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

// Package interpreter provides functions to evaluate parsed expressions with
// the option to augment the evaluation with inputs and functions supplied at
// evaluation time.
package interpreter

import (
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
	"github.com/google/cel-go/common/types/traits"
	"github.com/google/cel-go/interpreter/functions"
)

// Interpreter generates a new Interpretable from a Program.
type Interpreter interface {
	// NewInterpretable returns an Interpretable from a Program.
	NewInterpretable(program *Program) Interpretable
}

// Interpretable can accept a given Activation and produce a value along with
// an accompanying EvalState which can be used to inspect whether additional
// data might be necessary to complete the evaluation.
type Interpretable interface {
	// Eval an Activation to produce an output and EvalState.
	Eval(activation Activation) (ref.Value, EvalState)
}

type exprInterpreter struct {
	dispatcher   Dispatcher
	packager     packages.Packager
	typeProvider ref.TypeProvider
}

// NewInterpreter builds an Interpreter from a Dispatcher and TypeProvider
// which will be used throughout the Eval of all Interpretable instances
// gerenated from it.
func NewInterpreter(dispatcher Dispatcher,
	packager packages.Packager,
	typeProvider ref.TypeProvider) Interpreter {
	return &exprInterpreter{
		dispatcher:   dispatcher,
		packager:     packager,
		typeProvider: typeProvider}
}

// NewStandardInterpreter builds a Dispatcher and TypeProvider with support
// for all of the CEL builtins defined in the language definition.
func NewStandardInterpreter(packager packages.Packager,
	typeProvider ref.TypeProvider) Interpreter {
	dispatcher := NewDispatcher()
	dispatcher.Add(functions.StandardOverloads()...)
	return NewInterpreter(dispatcher, packager, typeProvider)
}

func (i *exprInterpreter) NewInterpretable(program *Program) Interpretable {
	// program needs to be pruned with the TypeProvider
	evalState := NewEvalState(program.MaxInstructionID() + 1)
	intr := &exprInterpreter{
		dispatcher:   i.dispatcher.Init(evalState),
		packager:     i.packager,
		typeProvider: i.typeProvider}
	program.Init(evalState)
	return &exprInterpretable{
		interpreter: intr,
		program:     program,
		state:       evalState}
}

type exprInterpretable struct {
	interpreter *exprInterpreter
	program     *Program
	state       MutableEvalState
}

func (i *exprInterpretable) Eval(activation Activation) (ref.Value, EvalState) {
	// register machine-like evaluation of the program with the given activation.
	currActivation := activation
	// program counter
	pc := 0
	end := len(i.program.Instructions)
	steps := i.program.Instructions
	var resultID int64
	for pc < end {
		step := steps[pc]
		resultID = step.GetID()
		switch step.(type) {
		case *IdentExpr:
			i.evalIdent(step.(*IdentExpr), currActivation)
		case *SelectExpr:
			i.evalSelect(step.(*SelectExpr), currActivation)
		case *CallExpr:
			i.evalCall(step.(*CallExpr), currActivation)
		case *CreateListExpr:
			i.evalCreateList(step.(*CreateListExpr))
		case *CreateMapExpr:
			i.evalCreateMap(step.(*CreateMapExpr))
		case *CreateObjectExpr:
			i.evalCreateType(step.(*CreateObjectExpr))
		case *MovInst:
			i.evalMov(step.(*MovInst))
			// Special instruction for modifying the program cursor
		case *JumpInst:
			jmpExpr := step.(*JumpInst)
			if jmpExpr.OnCondition(i.state) {
				jmpPc := pc + jmpExpr.Count
				if jmpPc < 0 || jmpPc >= end {
					// TODO: Error, the jump count should never exceed the
					// program length.
					panic("jumped too far")
				}
				pc = jmpPc
				continue
			}
			// Special instructions for modifying the activation stack
		case *PushScopeInst:
			pushScope := step.(*PushScopeInst)
			scopeDecls := pushScope.Declarations
			childActivaton := make(map[string]interface{})
			for key, declID := range scopeDecls {
				childActivaton[key] = func() interface{} {
					return i.value(declID)
				}
			}
			currActivation =
				NewHierarchicalActivation(currActivation, NewActivation(childActivaton))
		case *PopScopeInst:
			currActivation = currActivation.Parent()
		}
		pc++
	}
	result := i.value(resultID)
	if result == nil {
		result, _ = i.state.OnlyValue()
	}
	return result, i.state
}

func (i *exprInterpretable) evalConst(constExpr *ConstExpr) {
	i.setValue(constExpr.GetID(), constExpr.Value)
}

func (i *exprInterpretable) evalIdent(idExpr *IdentExpr, currActivation Activation) {
	// TODO: Refactor this code for sharing.
	if result, found := currActivation.ResolveName(idExpr.Name); found {
		i.setValue(idExpr.GetID(), result)
	} else if idVal, found := i.interpreter.typeProvider.FindIdent(idExpr.Name); found {
		i.setValue(idExpr.GetID(), idVal)
	} else {
		i.setValue(idExpr.GetID(), types.Unknown{idExpr.ID})
	}
}

func (i *exprInterpretable) evalSelect(selExpr *SelectExpr, currActivation Activation) {
	operand := i.value(selExpr.Operand)
	if !operand.Type().HasTrait(traits.IndexerType) {
		// If the operand is unknown, this could be an identifer.
		if types.IsUnknown(operand) {
			resVal := i.resolveUnknown(operand.(types.Unknown), selExpr, currActivation)
			i.setValue(selExpr.GetID(), resVal)
			return
		}
		// If the operand is an error, early return.
		if types.IsError(operand) {
			i.setValue(selExpr.GetID(), operand)
			return
		}
		// Otherwise, create an error.
		i.setValue(selExpr.GetID(), types.NewErr("invalid operand in select"))
		return
	}
	field := types.String(selExpr.Field)
	if selExpr.TestOnly {
		if operand.Type() == types.MapType {
			i.setValue(selExpr.GetID(), operand.(traits.Container).Contains(field))
			return
		}
		if operand.Type().HasTrait(traits.FieldTesterType) {
			i.setValue(selExpr.GetID(), operand.(traits.FieldTester).IsSet(field))
			return
		}
		i.setValue(selExpr.GetID(), types.NewErr("invalid operand in select"))
		return
	}
	fieldValue := operand.(traits.Indexer).Get(field)
	i.setValue(selExpr.GetID(), fieldValue)
}

// resolveUnknown attempts to resolve a qualified name from a select expression
// which may have generated unknown values during the course of execution if
// the expression was not type-checked and the select, in fact, refers to a
// qualified identifier name instead of a series of field selections.
//
// Returns one of the following:
// - The resolved identifier value from the activation
// - An unknown value if the expression is a valid identifier, but was not found.
// - Otherwise, an error.
func (i *exprInterpretable) resolveUnknown(unknown types.Unknown,
	selExpr *SelectExpr,
	currActivation Activation) ref.Value {
	if object, found := currActivation.ResolveReference(selExpr.ID); found {
		return object
	}
	validIdent := true
	identifier := selExpr.Field
	for _, arg := range unknown {
		inst := i.program.GetInstruction(arg)
		switch inst.(type) {
		case *IdentExpr:
			identifier = inst.(*IdentExpr).Name + "." + identifier
		case *SelectExpr:
			identifier = inst.(*SelectExpr).Field + "." + identifier
		default:
			argVal := i.value(arg)
			if argVal.Type() == types.StringType {
				identifier = string(argVal.(types.String)) + "." + identifier
			} else {
				validIdent = false
				break
			}
		}
	}
	if !validIdent {
		return types.NewErr("invalid identifier encountered: %v", selExpr)
	}
	pkg := i.interpreter.packager
	tp := i.interpreter.typeProvider
	for _, id := range pkg.ResolveCandidateNames(identifier) {
		if object, found := currActivation.ResolveName(id); found {
			return object
		}
		if identVal, found := tp.FindIdent(id); found {
			return identVal
		}
	}
	return append(types.Unknown{selExpr.ID}, unknown...)
}

func (i *exprInterpretable) evalCall(callExpr *CallExpr, currActivation Activation) {
	if callExpr.Strict {
		for _, argID := range callExpr.Args {
			argVal := i.value(argID)
			argType := argVal.Type()
			if types.IsUnknownOrError(argType) {
				i.setValue(callExpr.GetID(), argVal)
				return
			}
		}
	}
	i.interpreter.dispatcher.Dispatch(callExpr)
}

func (i *exprInterpretable) evalCreateList(listExpr *CreateListExpr) {
	elements := make([]ref.Value, len(listExpr.Elements))
	for idx, elementID := range listExpr.Elements {
		elem := i.value(elementID)
		if types.IsUnknownOrError(elem) {
			i.setValue(listExpr.GetID(), elem)
			return
		}
		elements[idx] = i.value(elementID)
	}
	adaptingList := types.NewDynamicList(elements)
	i.setValue(listExpr.GetID(), adaptingList)
}

func (i *exprInterpretable) evalCreateMap(mapExpr *CreateMapExpr) {
	entries := make(map[ref.Value]ref.Value)
	for keyID, valueID := range mapExpr.KeyValues {
		key := i.value(keyID)
		if types.IsUnknownOrError(key) {
			i.setValue(mapExpr.GetID(), key)
			return
		}
		val := i.value(valueID)
		if types.IsUnknownOrError(val) {
			i.setValue(mapExpr.GetID(), val)
			return
		}
		entries[key] = val
	}
	adaptingMap := types.NewDynamicMap(entries)
	i.setValue(mapExpr.GetID(), adaptingMap)
}

func (i *exprInterpretable) evalCreateType(objExpr *CreateObjectExpr) {
	fields := make(map[string]ref.Value)
	for field, valueID := range objExpr.FieldValues {
		val := i.value(valueID)
		if types.IsUnknownOrError(val) {
			i.setValue(objExpr.GetID(), val)
			return
		}
		fields[field] = val
	}
	i.setValue(objExpr.GetID(), i.newValue(objExpr.Name, fields))
}

func (i *exprInterpretable) evalMov(movExpr *MovInst) {
	i.setValue(movExpr.ToExprID, i.value(movExpr.GetID()))
}

func (i *exprInterpretable) value(id int64) ref.Value {
	if object, found := i.state.Value(id); found {
		return object
	}
	return types.Unknown{id}
}

func (i *exprInterpretable) setValue(id int64, value ref.Value) {
	i.state.SetValue(id, value)
}

func (i *exprInterpretable) newValue(typeName string,
	fields map[string]ref.Value) ref.Value {
	pkg := i.interpreter.packager
	tp := i.interpreter.typeProvider
	for _, qualifiedTypeName := range pkg.ResolveCandidateNames(typeName) {
		if _, found := tp.FindType(qualifiedTypeName); found {
			typeName = qualifiedTypeName
			break
		}
	}
	return i.interpreter.typeProvider.NewValue(typeName, fields)
}
