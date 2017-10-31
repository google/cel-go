package types

import "fmt"

type TypeType struct {
	target Type
}

var _ Type = &TypeType{}

func (tt *TypeType) Target() Type {
	return tt.target
}

func (tt *TypeType) Kind() TypeKind {
	return KindType
}

func (tt *TypeType) Equals(t Type) bool {
	other, ok := t.(*TypeType)
	if !ok {
		return false
	}

	return tt.target.Equals(other.target)
}

func (t *TypeType) String() string {
	return fmt.Sprintf("type(%s)", t.target)
}

func NewTypeType(target Type) *TypeType {
	return &TypeType{
		target: target,
	}
}
