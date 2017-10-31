package semantics

import (
	"celgo/semantics/types"
	"fmt"
)

type Function struct {
	name string

	overloads []*Overload
}

var _ Declaration = &Function{}

func (f *Function) Name() string {
	return f.name
}

func (f *Function) Overloads() []*Overload {
	return f.overloads[:]
}

func (f *Function) Merge(other *Function) *Function {
	if f.name != other.name {
		return nil
	}

	// TODO: Check for conflicts.
	overloads := append(f.overloads, other.overloads...)

	return NewFunction(f.name, overloads...)
}

func (f *Function) String() string {
	result := f.name + "("

	for i, o := range f.overloads {
		if i > 0 {
			result += "|"
		}
		result += o.id
	}

	result += ")"
	return result
}

type Overload struct {
	id         string
	isInstance bool
	argTypes   []types.Type
	resultType types.Type
	typeParams []string
}

func (o *Overload) Id() string {
	return o.id
}

func (o *Overload) IsInstance() bool {
	return o.isInstance
}

func (o *Overload) ArgTypes() []types.Type {
	return o.argTypes[:]
}

func (o *Overload) ResultType() types.Type {
	return o.resultType
}

func (o *Overload) TypeParams() []string {
	return o.typeParams[:]
}

func (o *Overload) String() string {
	result := ""

	argTypes := o.argTypes[:]
	if o.isInstance {
		result += fmt.Sprintf("%v.%s(", argTypes[0], o.id)
		argTypes = argTypes[1:]
	} else {
		result += fmt.Sprintf("%s(", o.id)
	}

	for i, a := range argTypes {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%v", a)
	}
	result += ")"

	return result
}

func NewFunction(name string, overloads ...*Overload) *Function {
	return &Function{
		name:      name,
		overloads: overloads,
	}
}

func NewOverload(id string, isInstance bool, resultType types.Type, argTypes ...types.Type) *Overload {
	return &Overload{
		id:         id,
		isInstance: isInstance,
		resultType: resultType,
		argTypes:   argTypes,
		typeParams: []string{},
	}
}

func NewParameterizedOverload(id string, isInstance bool, typeParams []string, resultType types.Type, argTypes ...types.Type) *Overload {
	return &Overload{
		id:         id,
		isInstance: isInstance,
		resultType: resultType,
		argTypes:   argTypes,
		typeParams: typeParams,
	}
}
