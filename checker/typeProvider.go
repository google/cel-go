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

package checker

import (
	"github.com/google/cel-go/semantics/types"
)

// TypeProvider defines methods to lookup types and enums, and resolve field types.
type TypeProvider interface {

	// LookupType looks up the Type given a qualified typeName. Returns nil if not found.
	LookupType(typeName string) types.Type

	LookupFieldType(t *types.MessageType, fieldName string) (*FieldType, bool)

	// LookupEnumValue looks up the integer enum value given an enumName. Returns nil if not found.
	LookupEnumValue(enumName string) (int64, bool)
}

type FieldType struct {
	Type             types.Type
	SupportsPresence bool
}

type InMemoryTypeProvider struct {
	// TODO: This is pretty broken. We should fix it at some point
	types  map[string]types.Type
	fields map[string]map[string]*FieldType
	enums  map[string]int64
}

var _ TypeProvider = &InMemoryTypeProvider{}

func NewInMemoryTypeProvider() *InMemoryTypeProvider {
	return &InMemoryTypeProvider{
		types:  make(map[string]types.Type),
		fields: make(map[string]map[string]*FieldType),
		enums:  make(map[string]int64),
	}
}

// LookupType looks up the Type given a qualified typeName. Returns nil if not found.
func (p *InMemoryTypeProvider) LookupType(typeName string) types.Type {
	if t, found := p.types[typeName]; found {
		return types.NewTypeType(t)
	}
	return nil
}

func (p *InMemoryTypeProvider) LookupFieldType(t *types.MessageType, fieldName string) (*FieldType, bool) {
	if m, found1 := p.fields[t.Name()]; found1 {
		if f, found2 := m[fieldName]; found2 {
			return f, true
		}
	}

	return nil, false
}

// LookupEnumValue looks up the integer enum value given an enumName.
func (p *InMemoryTypeProvider) LookupEnumValue(enumName string) (int64, bool) {
	if i, ok := p.enums[enumName]; ok {
		return i, ok
	}

	return 0, false
}

func (p *InMemoryTypeProvider) AddType(name string, t types.Type) {
	p.types[name] = t
}

func (p *InMemoryTypeProvider) AddFieldType(t *types.MessageType, fieldName string, fieldType types.Type, supportsPresence bool) {

	var typeFields map[string]*FieldType
	var found bool
	if typeFields, found = p.fields[t.Name()]; !found {
		typeFields = make(map[string]*FieldType)
		p.fields[t.Name()] = typeFields
	}

	typeFields[fieldName] = &FieldType{Type: fieldType, SupportsPresence: supportsPresence}
}

func (p *InMemoryTypeProvider) AddEnum(name string, value int64) {
	p.enums[name] = value
}

///** Map of well-known type names to {@code Type} instances. */
//var wellKnownTypeMap = map[string]types.Type {
//	"google.protobuf.DoubleValue": types.NewWrapper(types.Double),
//	"google.protobuf.FloatValue": types.NewWrapper(types.Double),
//	"google.protobuf.Int64Value": types.NewWrapper(types.Int64),
//	"google.protobuf.Int32Value": types.NewWrapper(types.Int64),
//	"google.protobuf.UInt64Value": types.NewWrapper(types.Uint64),
//	"google.protobuf.UInt32Value": types.NewWrapper(types.Uint64),
//	"google.protobuf.BoolValue": types.NewWrapper(types.Bool),
//	"google.protobuf.StringValue": types.NewWrapper(types.String),
//	"google.protobuf.BytesValue": types.NewWrapper(types.Bytes),
//	"google.protobuf.Timestamp": types.Timestamp,
//	"google.protobuf.Duration": types.Duration,
//	"google.protobuf.Struct": types.Dynamic,
//	"google.protobuf.Value": types.Dynamic,
//	"google.protobuf.ListValue": types.Dynamic,
//	"google.protobuf.Any": types.Any,
//}
//
