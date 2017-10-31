package types

var Error Type = &errorType{}

type errorType struct {
}

var _ Type = &errorType{}

func (e *errorType) Kind() TypeKind {
	return KindError
}

func (e *errorType) Equals(t Type) bool {
	_, ok := t.(*errorType)
	return ok
}

func (e *errorType) String() string {
	return "!error!"
}
