package types

var Dynamic Type = &dynamicType{}

type dynamicType struct {
}

var _ Type = &dynamicType{}

func (d *dynamicType) Kind() TypeKind {
	return KindDynamic
}

func (d *dynamicType) Equals(t Type) bool {
	_, ok := t.(*dynamicType)
	return ok
}

func (d *dynamicType) String() string {
	return "dyn"
}
