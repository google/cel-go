package types

type TypeParameter struct {
	name string
}

var _ Type = &TypeParameter{}

func (p *TypeParameter) Kind() TypeKind {
	return KindTypeParameter
}

func (p *TypeParameter) Equals(t Type) bool {
	other, ok := t.(*TypeParameter)
	if !ok {
		return false
	}
	return other.name == p.name
}

func (p *TypeParameter) String() string {
	return p.name
}

func NewTypeParam(name string) *TypeParameter {
	return &TypeParameter{
		name: name,
	}
}
