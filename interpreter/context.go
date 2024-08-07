package interpreter

import (
	"context"

	"github.com/google/cel-go/common/functions"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func ToInterpretableDecoratorContext(dec InterpretableDecorator) InterpretableDecoratorContext {
	return func(i InterpretableContext) (InterpretableContext, error) {
		ret, err := dec(i)
		if err != nil {
			return nil, err
		}
		ic, ok := ret.(InterpretableContext)
		if ok {
			return ic, nil
		}
		return ToInterpretableContext(ret), nil
	}
}

func ToInterpretableContext(i Interpretable) InterpretableContext {
	if i == nil {
		return nil
	}
	return &interpretableContextImpl{
		Interpretable: i,
	}
}

func ToInterpretableAttributeContext(i InterpretableAttribute) InterpretableAttributeContext {
	if i == nil {
		return nil
	}
	return &interpretableAttributeContextImpl{
		InterpretableAttribute: i,
	}
}

func ToInterpretableConstructorContext(i InterpretableConstructor) InterpretableConstructorContext {
	if i == nil {
		return nil
	}
	return &interpretableConstructorContextImpl{
		InterpretableConstructor: i,
	}
}

func ToDispatcherContext(d Dispatcher) DispatcherContext {
	if d == nil {
		return nil
	}
	return &dispatcherContextImpl{
		Dispatcher: d,
	}
}

func ToQualifierContext(q Qualifier) QualifierContext {
	if q == nil {
		return nil
	}
	return &qualifierContextImpl{
		Qualifier: q,
	}
}

func ToNamespacedAttributeContext(attr NamespacedAttribute) NamespacedAttributeContext {
	if attr == nil {
		return nil
	}
	return &namespacedAttributeContextImpl{
		NamespacedAttribute: attr,
	}
}

func ToAttributeFactoryContext(fac AttributeFactory) AttributeFactoryContext {
	return &attributeFactoryContextImpl{
		AttributeFactory: fac,
	}
}

func ToAttributeContext(attr Attribute) AttributeContext {
	return &attributeContextImpl{
		Attribute: attr,
	}
}

type interpretableContextImpl struct {
	Interpretable
}

func (impl *interpretableContextImpl) EvalContext(_ context.Context, activation Activation) ref.Val {
	return impl.Interpretable.Eval(activation)
}

type interpretableAttributeContextImpl struct {
	InterpretableAttribute
}

func (impl *interpretableAttributeContextImpl) AttrContext(_ context.Context) AttributeContext {
	return ToAttributeContext(impl.InterpretableAttribute.Attr())
}

func (impl *interpretableAttributeContextImpl) AddQualifierContext(_ context.Context, q Qualifier) (AttributeContext, error) {
	ret, err := impl.InterpretableAttribute.AddQualifier(q)
	if err != nil {
		return nil, err
	}
	return ToAttributeContext(ret), nil
}

func (impl *interpretableAttributeContextImpl) QualifyContext(_ context.Context, vars Activation, obj any) (any, error) {
	return impl.InterpretableAttribute.Qualify(vars, obj)
}

func (impl *interpretableAttributeContextImpl) QualifyIfPresentContext(_ context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return impl.InterpretableAttribute.QualifyIfPresent(vars, obj, presenceOnly)
}

func (impl *interpretableAttributeContextImpl) ResolveContext(_ context.Context, vars Activation) (any, error) {
	return impl.InterpretableAttribute.Resolve(vars)
}

func (impl *interpretableAttributeContextImpl) EvalContext(_ context.Context, vars Activation) ref.Val {
	return impl.InterpretableAttribute.Eval(vars)
}

type interpretableConstructorContextImpl struct {
	InterpretableConstructor
}

func (impl *interpretableConstructorContextImpl) EvalContext(_ context.Context, vars Activation) ref.Val {
	return impl.InterpretableConstructor.Eval(vars)
}

type dispatcherContextImpl struct {
	Dispatcher
}

func (impl *dispatcherContextImpl) AddContext(overloads ...*functions.OverloadContext) error {
	converted := make([]*functions.Overload, len(overloads))
	for idx, overload := range overloads {
		converted[idx] = overload.ToOverload()
	}
	return impl.Dispatcher.Add(converted...)
}

func (impl *dispatcherContextImpl) FindOverloadContext(overload string) (*functions.OverloadContext, bool) {
	o, found := impl.Dispatcher.FindOverload(overload)
	if found {
		return o.ToOverloadContext(), found
	}
	return nil, false
}

type qualifierContextImpl struct {
	Qualifier
}

func (impl *qualifierContextImpl) QualifyContext(_ context.Context, vars Activation, obj any) (any, error) {
	return impl.Qualifier.Qualify(vars, obj)
}

func (impl *qualifierContextImpl) QualifyIfPresentContext(_ context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return impl.Qualifier.QualifyIfPresent(vars, obj, presenceOnly)
}

type namespacedAttributeContextImpl struct {
	NamespacedAttribute
}

func (impl *namespacedAttributeContextImpl) QualifyContext(_ context.Context, vars Activation, obj any) (any, error) {
	return impl.NamespacedAttribute.Qualify(vars, obj)
}

func (impl *namespacedAttributeContextImpl) QualifyIfPresentContext(_ context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return impl.NamespacedAttribute.QualifyIfPresent(vars, obj, presenceOnly)
}

func (impl *namespacedAttributeContextImpl) AddQualifierContext(_ context.Context, qual Qualifier) (AttributeContext, error) {
	ret, err := impl.NamespacedAttribute.AddQualifier(qual)
	if err != nil {
		return nil, err
	}
	return ToAttributeContext(ret), nil
}

func (impl *namespacedAttributeContextImpl) ResolveContext(_ context.Context, vars Activation) (any, error) {
	return impl.NamespacedAttribute.Resolve(vars)
}

type attributeFactoryContextImpl struct {
	AttributeFactory
}

func (impl *attributeFactoryContextImpl) AbsoluteAttributeContext(_ context.Context, id int64, names ...string) NamespacedAttributeContext {
	return ToNamespacedAttributeContext(impl.AttributeFactory.AbsoluteAttribute(id, names...))
}

func (impl *attributeFactoryContextImpl) ConditionalAttributeContext(_ context.Context, id int64, expr InterpretableContext, t, f AttributeContext) AttributeContext {
	return ToAttributeContext(impl.AttributeFactory.ConditionalAttribute(id, expr.(Interpretable), t, f.(Attribute)))
}

func (impl *attributeFactoryContextImpl) MaybeAttributeContext(_ context.Context, id int64, name string) AttributeContext {
	return ToAttributeContext(impl.AttributeFactory.MaybeAttribute(id, name))
}

func (impl *attributeFactoryContextImpl) RelativeAttributeContext(_ context.Context, id int64, operand InterpretableContext) AttributeContext {
	return ToAttributeContext(impl.AttributeFactory.RelativeAttribute(id, operand.(Interpretable)))
}

func (impl *attributeFactoryContextImpl) NewQualifierContext(_ context.Context, objType *types.Type, qualID int64, val any, opt bool) (QualifierContext, error) {
	ret, err := impl.AttributeFactory.NewQualifier(objType, qualID, val, opt)
	if err != nil {
		return nil, err
	}
	return ToQualifierContext(ret), nil
}

type attributeContextImpl struct {
	Attribute
}

func (impl *attributeContextImpl) QualifyContext(_ context.Context, vars Activation, obj any) (any, error) {
	return impl.Attribute.Qualify(vars, obj)
}

func (impl *attributeContextImpl) QualifyIfPresentContext(_ context.Context, vars Activation, obj any, presenceOnly bool) (any, bool, error) {
	return impl.Attribute.QualifyIfPresent(vars, obj, presenceOnly)
}

func (impl *attributeContextImpl) AddQualifierContext(ctx context.Context, qual Qualifier) (AttributeContext, error) {
	ret, err := impl.Attribute.AddQualifier(qual)
	if err != nil {
		return nil, err
	}
	return ToAttributeContext(ret), nil
}

func (impl *attributeContextImpl) ResolveContext(_ context.Context, vars Activation) (any, error) {
	return impl.Attribute.Resolve(vars)
}

func unwrapContextImpl(v any) any {
	switch impl := v.(type) {
	case *interpretableContextImpl:
		return impl.Interpretable
	case *interpretableAttributeContextImpl:
		return impl.InterpretableAttribute
	case *interpretableConstructorContextImpl:
		return impl.InterpretableConstructor
	case *attributeContextImpl:
		return impl.Attribute
	case *namespacedAttributeContextImpl:
		return impl.NamespacedAttribute
	case *qualifierContextImpl:
		return impl.Qualifier
	case *dispatcherContextImpl:
		return impl.Dispatcher
	}
	return v
}
