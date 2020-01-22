// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interpreter

import (
	"fmt"

	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

// AttributePattern represents a top-level variable with an optional set of qualifier patterns.
//
// The variable name must always be a string, and may be a qualified name according to the CEL
// namespacing conventions, e.g. 'ns.app.a'.
//
// The qualifier patterns for attribute matching must be one of the following:
//
//   - valid map key type: string, int, uint, bool
//   - wildcard (*)
//
// Examples:
//
//   1. ns.myvar["complex-value"]
//   2. ns.myvar["complex-value"][0]
//   3. ns.myvar["complex-value"].*.name
//
// The first example is simple: match an attribute where the variable is 'ns.myvar' with a
// field access on 'complex-value'. The second example expands the match to indicate that only
// a specific index `0` should match. And lastly, the third example matches any indexed access
// that later selects the 'name' field.
type AttributePattern struct {
	variable          string
	qualifierPatterns []*AttributeQualifierPattern
}

// NewAttributePattern produces a new mutable AttributePattern based on a variable name.
func NewAttributePattern(variable string) *AttributePattern {
	return &AttributePattern{
		variable:          variable,
		qualifierPatterns: []*AttributeQualifierPattern{},
	}
}

// Field adds a string qualifier pattern to the AttributePattern. The string may be a valid
// identifier, or string map key including empty string.
func (apat *AttributePattern) Field(pattern string) *AttributePattern {
	apat.qualifierPatterns = append(apat.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return apat
}

// Index adds an int qualifier pattern to the AttributePattern. The index may be either a map or
// list index.
func (apat *AttributePattern) Index(pattern int64) *AttributePattern {
	apat.qualifierPatterns = append(apat.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return apat
}

// IndexUint adds an uint qualifier pattern for a map index operation to the AttributePattern.
func (apat *AttributePattern) IndexUint(pattern uint64) *AttributePattern {
	apat.qualifierPatterns = append(apat.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return apat
}

// IndexBool adds a bool qualifier pattern for a map index operation to the AttributePattern.
func (apat *AttributePattern) IndexBool(pattern bool) *AttributePattern {
	apat.qualifierPatterns = append(apat.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return apat
}

// Wildcard adds a special sentinel qualifier pattern that indicates any value will yeild a
// qualifier match.
func (apat *AttributePattern) Wildcard() *AttributePattern {
	apat.qualifierPatterns = append(apat.qualifierPatterns,
		&AttributeQualifierPattern{wildcard: true})
	return apat
}

// Matches returns true if the variable matches the AttributePattern variable.
func (apat *AttributePattern) Matches(variable string) bool {
	return apat.variable == variable
}

// QualifierPatterns returns the set of AttributeQualifierPattern values on the AttributePattern.
func (apat *AttributePattern) QualifierPatterns() []*AttributeQualifierPattern {
	return apat.qualifierPatterns
}

// AttributeQualifierPattern holds a wilcard or valued qualifier pattern.
type AttributeQualifierPattern struct {
	wildcard bool
	value    interface{}
}

// Matches returns true if the qualifier pattern is a wildcard, or the Qualifier implements the
// qualifierValueEquator interface and its IsValueEqualTo returns true for the qualifier pattern.
func (qpat *AttributeQualifierPattern) Matches(q Qualifier) bool {
	if qpat.wildcard {
		return true
	}
	qve, ok := q.(qualifierValueEquator)
	return ok && qve.QualifierValueEquals(qpat.value)
}

// qualifierValueEquator defines an interface for determining if an input value, of valid map key
// type, is equal to the value held in the Qualifier. This interface is used by the
// AttributeQualifierPattern to determine pattern matches for non-wildcard qualifier patterns.
//
// Note: Attribute values are also Qualifier values; however, Attriutes are resolved before
// qualification happens. This is an implementation detail, but one relevant to why the Attribute
// types do no surface in the list of implementations.
type qualifierValueEquator interface {
	// QualifierValueEquals returns true if the input value is equal to the value held in the
	// Qualifier.
	QualifierValueEquals(value interface{}) bool
}

// QualifierValueEquals implementation for boolean qualifiers.
func (q *boolQualifier) QualifierValueEquals(value interface{}) bool {
	bval, ok := value.(bool)
	return ok && q.value == bval
}

// QualifierValueEquals implementation for field qualifiers.
func (q *fieldQualifier) QualifierValueEquals(value interface{}) bool {
	sval, ok := value.(string)
	return ok && q.Name == sval
}

// QualifierValueEquals implementation for string qualifiers.
func (q *stringQualifier) QualifierValueEquals(value interface{}) bool {
	sval, ok := value.(string)
	return ok && q.value == sval
}

// QualifierValueEquals implementation for int qualifiers.
func (q *intQualifier) QualifierValueEquals(value interface{}) bool {
	ival, ok := value.(int64)
	return ok && q.value == ival
}

// QualifierValueEquals implementation for uint qualifiers.
func (q *uintQualifier) QualifierValueEquals(value interface{}) bool {
	uval, ok := value.(uint64)
	return ok && q.value == uval
}

// NewPartialAttributeFactory returns an AttributeFactory implementation capable of performing
// AttributePattern matches with PartialActivation inputs.
func NewPartialAttributeFactory(pkg packages.Packager,
	adapter ref.TypeAdapter,
	provider ref.TypeProvider) AttributeFactory {
	fac := NewAttributeFactory(pkg, adapter, provider)
	return &partialAttributeFactory{
		AttributeFactory: fac,
		pkg:              pkg,
		adapter:          adapter,
		provider:         provider,
	}
}

type partialAttributeFactory struct {
	AttributeFactory
	pkg      packages.Packager
	adapter  ref.TypeAdapter
	provider ref.TypeProvider
}

// AbsoluteAttribute implementation of the AttributeFactory interface which wraps the
// NamespacedAttribute resolution in an internal attributeMatcher object to dynamically match
// unknown patterns from PartialActivation inputs if given.
func (fac *partialAttributeFactory) AbsoluteAttribute(id int64, names ...string) NamespacedAttribute {
	attr := fac.AttributeFactory.AbsoluteAttribute(id, names...)
	return &attributeMatcher{fac: fac, NamespacedAttribute: attr}
}

// MaybeAttribute implementation of the AttributeFactory interface which ensure that the set of
// 'maybe' NamespacedAttribute values are produced using the PartialAttributeFactory rather than
// the base AttributeFactory implementation.
func (fac *partialAttributeFactory) MaybeAttribute(id int64, name string) Attribute {
	return &maybeAttribute{
		id: id,
		attrs: []NamespacedAttribute{
			fac.AbsoluteAttribute(id, fac.pkg.ResolveCandidateNames(name)...),
		},
		adapter:  fac.adapter,
		provider: fac.provider,
		fac:      fac,
	}
}

// matchesUnknownPatterns returns true if the variable names and qualifiers for a given
// Attribute value match any of the ActivationPattern objects in the set of unknown activation
// patterns on the given PartialActivation.
func (fac *partialAttributeFactory) matchesUnknownPatterns(
	vars PartialActivation,
	attrID int64,
	variableNames []string,
	qualifiers []Qualifier) (types.Unknown, error) {
	patterns := vars.UnknownAttributePatterns()
	candIndices := map[int]struct{}{}
	for _, variable := range variableNames {
		for i, pat := range patterns {
			if pat.Matches(variable) {
				candIndices[i] = struct{}{}
			}
		}
	}
	// Determine whether to return early if there are no candidate unknown patterns.
	if len(candIndices) == 0 {
		return nil, nil
	}
	// Determine whether to return early if there are no qualifiers.
	if len(qualifiers) == 0 {
		return types.Unknown{attrID}, nil
	}
	// Resolve the attribute qualifiers into a static set. This prevents more dynamic
	// Attribute resolutions than necessary when there are multiple unknown patterns
	// that traverse the same Attribute-based qualifier field.
	newQuals := make([]Qualifier, len(qualifiers), len(qualifiers))
	for i, qual := range qualifiers {
		attr, isAttr := qual.(Attribute)
		if isAttr {
			val, err := attr.Resolve(vars)
			if err != nil {
				return nil, err
			}
			unk, isUnk := val.(types.Unknown)
			if isUnk {
				return unk, nil
			}
			qual, err = fac.NewQualifier(nil, qual.ID(), val)
			if err != nil {
				return nil, err
			}
		}
		newQuals[i] = qual
	}
	// Determine whether any of the unknown patterns match.
	for patIdx := range candIndices {
		pat := patterns[patIdx]
		isUnk := true
		matchExprID := attrID
		qualPats := pat.QualifierPatterns()
		for i, qual := range newQuals {
			if i >= len(qualPats) {
				break
			}
			matchExprID = qual.ID()
			qualPat := qualPats[i]
			if !qualPat.Matches(qual) {
				isUnk = false
				break
			}
		}
		if isUnk {
			return types.Unknown{matchExprID}, nil
		}
	}
	return nil, nil
}

type attributeMatcher struct {
	NamespacedAttribute
	qualifiers []Qualifier
	fac        *partialAttributeFactory
}

// AddQualifier implements the Attribute interface method which adds the Qualifier onto the
// underlying NamespacedAttribute as well as tracks the Qualifier in internal storage. The
// double accounting for the Qualifier values is to assist with AttributePattern matching in
// the TryResolve call.
func (m *attributeMatcher) AddQualifier(qual Qualifier) (Attribute, error) {
	_, err := m.NamespacedAttribute.AddQualifier(qual)
	if err != nil {
		return nil, err
	}
	m.qualifiers = append(m.qualifiers, qual)
	return m, nil
}

// Resolve is an implementation of the Attribute interface method which uses the
// attributeMatcher TryResolve implementation rather than the embedded NamespacedAttribute
// Resolve implementation.
func (m *attributeMatcher) Resolve(vars Activation) (interface{}, error) {
	obj, found, err := m.TryResolve(vars)
	if err != nil {
		return nil, err
	}
	if found {
		return obj, nil
	}
	return nil, fmt.Errorf("no such attribute: %v", m.NamespacedAttribute)
}

// TryResolve is an implementation of the NamespacedAttribute interface method which tests
// for matching unknown attribute patterns and returns types.Unknown if present. Otherwise,
// the standard Resolve logic applies.
func (m *attributeMatcher) TryResolve(vars Activation) (interface{}, bool, error) {
	id := m.NamespacedAttribute.ID()
	partial, isPartial := vars.(PartialActivation)
	if isPartial {
		unk, err := m.fac.matchesUnknownPatterns(
			partial,
			id,
			m.CandidateVariableNames(),
			m.qualifiers)
		if err != nil {
			return nil, true, err
		}
		if unk != nil {
			return unk, true, nil
		}
	}
	return m.NamespacedAttribute.TryResolve(vars)
}

// Qualify is an implementation of the Qualifier interface method.
func (m *attributeMatcher) Qualify(vars Activation, obj interface{}) (interface{}, error) {
	val, err := m.Resolve(vars)
	if err != nil {
		return nil, err
	}
	unk, isUnk := val.(types.Unknown)
	if isUnk {
		return unk, nil
	}
	qual, err := m.fac.NewQualifier(nil, m.ID(), val)
	if err != nil {
		return nil, err
	}
	return qual.Qualify(vars, obj)
}
