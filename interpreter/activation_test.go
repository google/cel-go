package interpreter

import (
	"testing"
)

func TestNewActivation(t *testing.T) {
	activation := NewActivation(map[string]interface{}{"a": true})
	if val, found := activation.ResolveName("a"); !found || !val.(bool) {
		t.Error("Activation failed to resolve 'a'")
	}
}

func TestHierarchicalActivation(t *testing.T) {
	// compose a parent with more properties than the child
	parent := NewActivation(map[string]interface{}{"a": "world", "b": -42})
	// compose the child such that it shadows the parent
	child := NewActivation(map[string]interface{}{"a": true, "c": "universe"})
	combined := ExtendActivation(parent, child)

	// Resolve the shadowed child value.
	if val, found := combined.ResolveName("a"); !found || !val.(bool) {
		t.Error("Activation failed to resolve shadow value of 'a'")
	}

	// Resolve the parent only value.
	if val, found := combined.ResolveName("b"); !found || val.(int) != -42 {
		t.Error("Activation failed to resolve parent value of 'b'")
	}

	// Resolve the child only value.
	if val, found := combined.ResolveName("c"); !found || val.(string) != "universe" {
		t.Error("Activation failed to resolve child value of 'c'")
	}
}
