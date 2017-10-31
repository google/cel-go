package types

var Null Type = &nullType{}

type nullType struct {
}

var _ Type = &nullType{}

func (n *nullType) Kind() TypeKind {
	return KindNull
}

func (n *nullType) Equals(t Type) bool {
	_, ok := t.(*nullType)
	return ok
}

func (n *nullType) String() string {
	return "null"
}
