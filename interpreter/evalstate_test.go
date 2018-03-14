package interpreter

import (
	"testing"
)

func TestGetterSetter(t *testing.T) {
	var evalState EvalState = NewEvalState()
	if val, found := evalState.Value(1); found || val != nil {
		t.Error("Unexpected value found", val)
	}
	var mutableState = evalState.(MutableEvalState)
	mutableState.SetValue(1, "hello")
	if greeting, found := evalState.Value(1); !found || greeting != "hello" {
		t.Error("Unexpected value found", greeting)
	}
}
