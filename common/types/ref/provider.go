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

package ref

import (
	"github.com/golang/protobuf/proto"

	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// TypeProvider specifies functions for creating new object instances and for
// resolving enum values by name.
type TypeProvider interface {
	// EnumValue returns the numeric value of the given enum value name.
	EnumValue(enumName string) Val

	// FindIdent takes a qualified identifier name and returns a Value if one
	// exists.
	FindIdent(identName string) (Val, bool)

	// FindType looks up the Type given a qualified typeName. Returns false
	// if not found.
	//
	// Used during type-checking only.
	FindType(typeName string) (*exprpb.Type, bool)

	// FieldFieldType returns the field type for a checked type value. Returns
	// false if the field could not be found.
	//
	// Used during type-checking only.
	FindFieldType(messageType string, fieldName string) (*FieldType, bool)

	// NewValue creates a new type value from a qualified name and a map of
	// field initializers.
	NewValue(typeName string, fields map[string]Val) Val
}

// TypeAdapter converts native Go values of varying type and complexity to equivalent CEL values.
type TypeAdapter interface {
	// NativeToValue converts the input `value` to a CEL `ref.Val`.
	NativeToValue(value interface{}) Val
}

// TypeRegistry allows third-parties to registry custom types. Not all `TypeProvider`
// implementations support type-customization, so these features are optional. However, a
// `TypeRegistry` should be a `TypeProvider` and a `TypeAdapter` to ensure that types
// which are registered can be converted to CEL representations.
type TypeRegistry interface {
	TypeAdapter
	TypeProvider

	// RegisterDescriptor registers the contents of a protocol buffer `FileDescriptor`.
	RegisterDescriptor(fileDesc *descpb.FileDescriptorProto) error

	// RegisterMessage registers a protocol buffer message and its dependencies.
	RegisterMessage(message proto.Message) error

	// RegisterType registers a type value with the provider which ensures the
	// provider is aware of how to map the type to an identifier.
	//
	// If a type is provided more than once with an alternative definition, the
	// call will result in an error.
	RegisterType(types ...Type) error
}

// FieldType represents a field's type value and whether that field supports
// presence detection.
type FieldType struct {
	// SupportsPresence indicates if the field having been set can be detected.
	SupportsPresence bool

	// Type of the field.
	Type *exprpb.Type
}
