package interpreter

import (
	"github.com/google/cel-go/interpreter/types"
	"github.com/google/cel-go/interpreter/types/adapters"
	"github.com/google/cel-go/interpreter/types/objects"
	"fmt"
)

type Interpreter interface {
	NewInterpretable(program Program) Interpretable
}

type Interpretable interface {
	Eval(activation Activation) (interface{}, EvalState)
}

type exprInterpreter struct {
	dispatcher   Dispatcher
	typeProvider types.TypeProvider
}

var _ Interpreter = &exprInterpreter{}

func NewInterpreter(dispatcher Dispatcher, typeProvider types.TypeProvider) *exprInterpreter {
	return &exprInterpreter{dispatcher, typeProvider}
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
		case *CreateTypeExpr:
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
			currActivation = ExtendActivation(currActivation, NewActivation(childActivaton))
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
	} else if indexer, ok := operand.(objects.Indexer); ok {
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
		fmt.Println(err)
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
	adaptingList := adapters.NewListAdapter(elements)
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
	adaptingMap := adapters.NewMapAdapter(entries)
	i.setValue(step.GetId(), adaptingMap)
}

func (i *exprInterpretable) evalCreateType(step Instruction) {
	typeExpr := step.(*CreateTypeExpr)
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
	} else {
		return nil
	}
}

func (i *exprInterpretable) setValue(id int64, value interface{}) {
	i.state.SetValue(id, value)
}

func (i *exprInterpretable) newValue(typeName string,
	fields map[string]interface{}) (adapters.MsgAdapter, error) {
	return i.interpreter.typeProvider.NewValue(typeName, fields)
}
