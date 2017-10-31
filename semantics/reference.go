package semantics

import (
	"celgo/ast"
)

type Reference interface {
	Equals(Reference) bool
	String() string
}

type FunctionReference struct {
	overloads []string
}

var _ Reference = &FunctionReference{}

func (f *FunctionReference) Equals(r Reference) bool {
	if other, ok := r.(*FunctionReference); ok {
		if len(other.overloads) == len(f.overloads) {
			for i, o := range f.overloads {
				if o != other.overloads[i] {
					return false
				}
			}
			return true
		}
	}

	return false
}

func (f *FunctionReference) String() string {
	result := ""
	for i, o := range f.overloads {
		if i > 0 {
			result += "|"
		}
		result += o
	}
	return result
}

type IdentReference struct {
	name  string
	value ast.Constant
}

var _ Reference = &IdentReference{}

func (f *IdentReference) Equals(r Reference) bool {
	if other, ok := r.(*IdentReference); ok {
		return other.name == f.name
	}

	return false
}

func (i *IdentReference) String() string {
	return i.name
}

func NewIdentReference(name string, constant ast.Constant) *IdentReference {
	return &IdentReference{
		name:  name,
		value: constant,
	}
}

func NewFunctionReference(overload string) *FunctionReference {
	return &FunctionReference{
		overloads: []string{overload},
	}
}

func (f *FunctionReference) AddOverloadReference(overload string) *FunctionReference {
	return &FunctionReference{
		overloads: append(f.overloads, overload),
	}
}

func (f *FunctionReference) Overloads() []string {
	return f.overloads[:]
}

func (r *IdentReference) Name() string {
	return r.name
}

func (r *IdentReference) Value() ast.Constant {
	return r.value
}
