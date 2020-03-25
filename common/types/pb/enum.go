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

package pb

import (
	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// NewEnumValueDescription produces a new named enum value description.
func NewEnumValueDescription(name string,
	desc *descpb.EnumValueDescriptorProto) *EnumValueDescription {
	return &EnumValueDescription{
		enumValueName: name,
		desc:          desc,
	}
}

// EnumValueDescription maps a qualified enum name to its numeric value.
type EnumValueDescription struct {
	enumValueName string
	desc          *descpb.EnumValueDescriptorProto
}

// Name returns the human-readable identifier name for the enum value.
func (ed *EnumValueDescription) Name() string {
	return ed.enumValueName
}

// Value returns the (numeric) value of the enum.
func (ed *EnumValueDescription) Value() int32 {
	return ed.desc.GetNumber()
}
