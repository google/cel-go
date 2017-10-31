package types

import "fmt"

type WrapperType struct {
	primitive *PrimitiveType
}

var _ Type = &WrapperType{}

func (w *WrapperType) Kind() TypeKind {
	return KindWrapper
}

func (w *WrapperType) Equals(t Type) bool {
	other, ok := t.(*WrapperType)
	if !ok {
		return false
	}

	return w.primitive.Equals(other.primitive)
}

func (w *WrapperType) String() string {
	return fmt.Sprintf("wrapper(%s)", w.primitive.String())
}

func (w *WrapperType) Primitive() *PrimitiveType {
	return w.primitive
}

func NewWrapper(p *PrimitiveType) *WrapperType {
	return &WrapperType{
		primitive: p,
	}
}
