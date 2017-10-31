package types

type MessageType struct {
	name string
}

var _ Type = &MessageType{}

func (m *MessageType) Kind() TypeKind {
	return KindMessage
}

func (m *MessageType) Equals(t Type) bool {
	if mt, ok := t.(*MessageType); ok {
		return t.(*MessageType).name == mt.name
	}

	return false
}

func (m *MessageType) String() string {
	return m.name
}

func (m *MessageType) Name() string {
	return m.name
}

func NewMessage(name string) *MessageType {
	return &MessageType{
		name: name,
	}
}
