package types

type FunctionType struct {
	resultType Type
	argTypes   []Type
}

var _ Type = &FunctionType{}

func (f *FunctionType) Kind() TypeKind {
	return KindFunction
}

func (f *FunctionType) Equals(t Type) bool {
	other, ok := t.(*FunctionType)
	if !ok {
		return false
	}

	if !f.resultType.Equals(other.resultType) {
		return false
	}

	if len(f.argTypes) != len(other.argTypes) {
		return false
	}

	for i := 0; i < len(f.argTypes); i++ {
		if !f.argTypes[i].Equals(other.argTypes[i]) {
			return false
		}
	}

	return true
}

func (f *FunctionType) String() string {
	return FormatFunction(f.resultType, f.argTypes, false)
}

func (f *FunctionType) ResultType() Type {
	return f.resultType
}

func (f *FunctionType) ArgTypes() []Type {
	return f.argTypes[:]
}

func NewFunctionType(resultType Type, argTypes []Type) *FunctionType {
	return &FunctionType{
		resultType: resultType,
		argTypes:   argTypes,
	}
}

func FormatFunction(resultType Type, argTypes []Type, isInstance bool) string {
	result := ""
	if isInstance {
		target := argTypes[0]
		argTypes = argTypes[1:]

		result += target.String()
		result += "."
	}

	result += "("
	for i, arg := range argTypes {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	if resultType != nil {
		result += " -> "
		result += resultType.String()
	}

	return result
}
