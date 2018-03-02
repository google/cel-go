// Copyright 2018 Google LLC
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
