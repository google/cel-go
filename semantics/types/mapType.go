package types

import "fmt"

type MapType struct {
	KeyType   Type
	ValueType Type
}

var _ Type = &MapType{}

func (m *MapType) Kind() TypeKind {
	return KindMap
}

func (m *MapType) Equals(t Type) bool {
	if other, ok := t.(*MapType); ok {
		return m.KeyType.Equals(other.KeyType) && m.ValueType.Equals(other.ValueType)
	}

	return false
}

func (m *MapType) String() string {
	return fmt.Sprintf("map(%s, %s)", m.KeyType.String(), m.ValueType.String())
}

func NewMap(keyType Type, valueType Type) *MapType {
	return &MapType{
		KeyType:   keyType,
		ValueType: valueType,
	}
}
