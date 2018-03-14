package interpreter

type EvalState interface {
	Value(int64) (interface{}, bool)
}

type MutableEvalState interface {
	EvalState
	SetValue(int64, interface{})
}

var _ EvalState = &defaultEvalState{}
var _ MutableEvalState = &defaultEvalState{}

func NewEvalState() *defaultEvalState {
	return &defaultEvalState{make(map[int64]interface{})}
}

type defaultEvalState struct {
	exprValues map[int64]interface{}
}

func (s *defaultEvalState) Value(exprId int64) (interface{}, bool) {
	object, found := s.exprValues[exprId]
	return object, found
}

func (s *defaultEvalState) SetValue(exprId int64, value interface{}) {
	s.exprValues[exprId] = value
}

type UnknownValue struct {
	Args []Instruction
}

func (o *UnknownValue) Equal(other interface{}) bool {
	return false
}

type ErrorValue struct {
	ErrorSet []error
}

func (o *ErrorValue) Equal(other interface{}) bool {
	return false
}
