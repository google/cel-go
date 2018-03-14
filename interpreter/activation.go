package interpreter

type Activation interface {
	ResolveReference(exprId int64) (interface{}, bool)
	ResolveName(name string) (interface{}, bool)
	Parent() Activation
}

func NewActivation(bindings map[string]interface{}) *MapActivation {
	return &MapActivation{bindings: bindings}
}

type MapActivation struct {
	references map[int64]interface{}
	bindings   map[string]interface{}
}

var _ Activation = &MapActivation{}

func (a *MapActivation) Parent() Activation {
	return nil
}

func (a *MapActivation) ResolveReference(exprId int64) (interface{}, bool) {
	if object, found := a.references[exprId]; found {
		return object, true
	}
	return nil, false
}

func (a *MapActivation) ResolveName(name string) (interface{}, bool) {
	// TODO: Look at how name resolution logic works for enums
	if object, found := a.bindings[name]; found {
		switch object.(type) {
		case func() interface{}:
			return object.(func() interface{})(), true
		default:
			return object, true
		}
	}
	return nil, false
}

type HierarchicalActivation struct {
	parent Activation
	child  Activation
}

var _ Activation = &HierarchicalActivation{}

func (a *HierarchicalActivation) Parent() Activation {
	return a.parent
}

func (a *HierarchicalActivation) ResolveReference(exprId int64) (interface{}, bool) {
	if object, found := a.child.ResolveReference(exprId); found {
		return object, found
	}
	return a.parent.ResolveReference(exprId)
}

func (a *HierarchicalActivation) ResolveName(name string) (interface{}, bool) {
	if object, found := a.child.ResolveName(name); found {
		return object, found
	}
	return a.parent.ResolveName(name)
}

func ExtendActivation(parent Activation, child Activation) Activation {
	return &HierarchicalActivation{parent, child}
}
