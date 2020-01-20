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
	"github.com/google/cel-go/common/packages"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

type partialAttributeFactory struct {
	AttributeFactory
	pkg      packages.Packager
	adapter  ref.TypeAdapter
	provider ref.TypeProvider
}

type attributeMatcher struct {
	NamespacedAttribute
	qualifiers []Qualifier
	fac        *partialAttributeFactory
}

func (m *attributeMatcher) AddQualifier(qual Qualifier) (Attribute, error) {
	_, err := m.NamespacedAttribute.AddQualifier(qual)
	if err != nil {
		return nil, err
	}
	m.qualifiers = append(m.qualifiers, qual)
	return m, nil
}

func (m *attributeMatcher) TryResolve(vars Activation) (interface{}, bool, error) {
	id := m.NamespacedAttribute.ID()
	partial, isPartial := vars.(PartialActivation)
	if isPartial {
		isUnk, err := m.fac.matchesUnknownPatterns(
			partial,
			m.CandidateVariableNames(),
			m.qualifiers)
		if err != nil {
			return nil, true, err
		}
		if isUnk {
			return types.Unknown{id}, true, nil
		}
	}
	return m.NamespacedAttribute.TryResolve(vars)
}

func (fac *partialAttributeFactory) AbsoluteAttribute(id int64, names ...string) NamespacedAttribute {
	attr := fac.AttributeFactory.AbsoluteAttribute(id, names...)
	return &attributeMatcher{fac: fac, NamespacedAttribute: attr}
}

func (fac *partialAttributeFactory) MaybeAttribute(id int64, name string) Attribute {
	return &maybeAttribute{
		id: id,
		attrs: []NamespacedAttribute{
			fac.AbsoluteAttribute(id, fac.pkg.ResolveCandidateNames(name)...),
		},
		adapter:  fac.adapter,
		provider: fac.provider,
		res:      fac,
	}
}

func (fac *partialAttributeFactory) matchesUnknownPatterns(
	vars PartialActivation,
	variableNames []string,
	qualifiers []Qualifier) (bool, error) {
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
		return false, nil
	}

	// Determine whether to return early if there are no qualifiers.
	if len(qualifiers) == 0 {
		return true, nil
	}
	// Resolve the attribute qualifiers into a static set.
	newQuals := make([]Qualifier, len(qualifiers), len(qualifiers))
	for i, qual := range qualifiers {
		attr, isAttr := qual.(Attribute)
		if isAttr {
			val, err := attr.Resolve(vars)
			if err != nil {
				return false, err
			}
			qual, err = fac.NewQualifier(nil, qual.ID(), val)
			if err != nil {
				return false, err
			}
		}
		newQuals[i] = qual
	}

	// Determine whether any of the unknown patterns match.
	for patIdx := range candIndices {
		pat := patterns[patIdx]
		isUnk := true
		// lastIdx := 0
		qualPats := pat.QualifierPatterns()
		for i, qual := range newQuals {
			//lastIdx = i
			if i >= len(qualPats) {
				break
			}
			qualPat := qualPats[i]
			if !qualPat.Matches(qual) {
				isUnk = false
				break
			}
		}
		if isUnk {
			return true, nil
			// return match type, attribute trail, and unknown pattern
		}
	}
	return false, nil
}

type AttributePattern struct {
	variable          string
	qualifierPatterns []*AttributeQualifierPattern
}

func NewAttributePattern(variable string) *AttributePattern {
	return &AttributePattern{
		variable:          variable,
		qualifierPatterns: []*AttributeQualifierPattern{},
	}
}

func (ap *AttributePattern) Any() *AttributePattern {
	ap.qualifierPatterns = append(ap.qualifierPatterns,
		&AttributeQualifierPattern{any: true})
	return ap
}

func (ap *AttributePattern) Name(pattern string) *AttributePattern {
	ap.qualifierPatterns = append(ap.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return ap
}

func (ap *AttributePattern) Index(pattern int64) *AttributePattern {
	ap.qualifierPatterns = append(ap.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return ap
}

func (ap *AttributePattern) IndexUint(pattern uint64) *AttributePattern {
	ap.qualifierPatterns = append(ap.qualifierPatterns,
		&AttributeQualifierPattern{value: pattern})
	return ap
}

func (ap *AttributePattern) True() *AttributePattern {
	ap.qualifierPatterns = append(ap.qualifierPatterns,
		&AttributeQualifierPattern{value: true})
	return ap
}

func (ap *AttributePattern) False() *AttributePattern {
	ap.qualifierPatterns = append(ap.qualifierPatterns,
		&AttributeQualifierPattern{value: false})
	return ap
}

func (ap *AttributePattern) Matches(variable string) bool {
	return ap.variable == variable
}

func (ap *AttributePattern) QualifierPatterns() []*AttributeQualifierPattern {
	return ap.qualifierPatterns
}

type AttributeQualifierPattern struct {
	any   bool
	value interface{}
}

func (qp *AttributeQualifierPattern) Matches(q Qualifier) bool {
	if qp.any {
		return true
	}
	cq, ok := q.(QualifierEquivalence)
	if !ok {
		panic("only constant valued qualifiers should appear here.")
		return false
	}
	return cq.IsEquivalentTo(qp.value)
}

type QualifierEquivalence interface {
	IsEquivalentTo(value interface{}) bool
}

func (q *boolQualifier) IsEquivalentTo(value interface{}) bool {
	bval, ok := value.(bool)
	return ok && q.value == bval
}

func (q *fieldQualifier) IsEquivalentTo(value interface{}) bool {
	sval, ok := value.(string)
	return ok && q.Name == sval
}

func (q *stringQualifier) IsEquivalentTo(value interface{}) bool {
	sval, ok := value.(string)
	return ok && q.value == sval
}

func (q *intQualifier) IsEquivalentTo(value interface{}) bool {
	ival, ok := value.(int64)
	return ok && q.value == ival
}

func (q *uintQualifier) IsEquivalentTo(value interface{}) bool {
	uval, ok := value.(uint64)
	return ok && q.value == uval
}
