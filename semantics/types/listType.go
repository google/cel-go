package types

import "fmt"

type ListType struct {
	ElementType Type
}

var _ Type = &ListType{}

func (l *ListType) Kind() TypeKind {
	return KindList
}

func (l *ListType) Equals(t Type) bool {
	if other, ok := t.(*ListType); ok {
		return l.ElementType.Equals(other.ElementType)
	}

	return false
}

func (l *ListType) String() string {
	return fmt.Sprintf("list(%s)", l.ElementType.String())
}

func NewList(elementType Type) *ListType {
	return &ListType{
		ElementType: elementType,
	}
}
