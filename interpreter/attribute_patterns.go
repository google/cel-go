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

type partialAttributeFactory struct {
	AttributeFactory
}

func (fac *partialAttributeFactory) matchesUnknownAttributePatterns(
	vars PartialActivation,
	variable string,
	qualifiers []Qualifier) (bool, error) {
	patterns := vars.UnknownAttributePatterns()
	candIndices := []int{}
	for i, pat := range patterns {
		if pat.Matches(variable) {
			candIndices = append(candIndices, i)
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
	for _, patIdx := range candIndices {
		pat := patterns[patIdx]
		isUnk := true
		// lastIdx := 0
		qualPats := pat.Qualifiers()
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

type AttributePattern interface {
	Matches(string) bool
	Qualifiers() []AttributeQualifierPattern
}

type AttributeQualifierPattern interface {
	Matches(Qualifier) bool
}
