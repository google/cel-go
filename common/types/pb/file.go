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
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
	dynamicpb "google.golang.org/protobuf/types/dynamicpb"
)

// newFileDescription returns a FileDescription instance with a complete listing of all the message
// types, enum values, and extension fields (for proto2) declared within any scope in the file.
func newFileDescription(fileDesc protoreflect.FileDescriptor, pbdb *Db) *FileDescription {
	metadata := collectFileMetadata(fileDesc)
	enums := make(map[string]*EnumValueDescription)
	for name, enumVal := range metadata.enumValues {
		enums[name] = newEnumValueDescription(name, enumVal)
	}
	types := make(map[string]*TypeDescription)
	for name, msgType := range metadata.msgTypes {
		types[name] = newTypeDescription(name, msgType)
	}
	fileExtensionMap := map[string]map[string]*FieldDescription{}
	for typeName, extensions := range metadata.msgExtensionMap {
		messageExtMap, found := fileExtensionMap[typeName]
		if !found {
			messageExtMap = make(map[string]*FieldDescription)
		}
		for _, ext := range extensions {
			extDesc := dynamicpb.NewExtensionType(ext).TypeDescriptor()
			messageExtMap[string(ext.FullName())] = newFieldDescription(extDesc)
		}
		fileExtensionMap[typeName] = messageExtMap
	}
	return &FileDescription{
		types:      types,
		enums:      enums,
		extensions: fileExtensionMap,
	}
}

// FileDescription holds a map of all types and enum values declared within a proto file.
type FileDescription struct {
	types      map[string]*TypeDescription
	enums      map[string]*EnumValueDescription
	extensions map[string]map[string]*FieldDescription
}

// GetEnumDescription returns an EnumDescription for a qualified enum value
// name declared within the .proto file.
func (fd *FileDescription) GetEnumDescription(enumName string) (*EnumValueDescription, bool) {
	ed, found := fd.enums[sanitizeProtoName(enumName)]
	return ed, found
}

// GetEnumNames returns the string names of all enum values in the file.
func (fd *FileDescription) GetEnumNames() []string {
	enumNames := make([]string, len(fd.enums))
	i := 0
	for _, e := range fd.enums {
		enumNames[i] = e.Name()
		i++
	}
	return enumNames
}

// GetTypeDescription returns a TypeDescription for a qualified protobuf message type name
// declared within the .proto file.
func (fd *FileDescription) GetTypeDescription(typeName string) (*TypeDescription, bool) {
	td, found := fd.types[sanitizeProtoName(typeName)]
	return td, found
}

// GetTypeNames returns the list of all type names contained within the file.
func (fd *FileDescription) GetTypeNames() []string {
	typeNames := make([]string, len(fd.types))
	i := 0
	for _, t := range fd.types {
		typeNames[i] = t.Name()
		i++
	}
	return typeNames
}

// GetExtension returns a protobuf extension field for the given message type and field name.
func (fd *FileDescription) GetExtension(messageType, fieldName string) (*FieldDescription, bool) {
	msgExtensions, found := fd.extensions[messageType]
	if !found {
		return nil, false
	}
	extField, found := msgExtensions[fieldName]
	return extField, found
}

// sanitizeProtoName strips the leading '.' from the proto message name.
func sanitizeProtoName(name string) string {
	if name != "" && name[0] == '.' {
		return name[1:]
	}
	return name
}

// fileMetadata is a flattened view of message types and enum values within a file descriptor.
type fileMetadata struct {
	// msgTypes maps from fully-qualified message name to descriptor.
	msgTypes map[string]protoreflect.MessageDescriptor
	// enumValues maps from fully-qualified enum value to enum value descriptor.
	enumValues map[string]protoreflect.EnumValueDescriptor
	// msgExtensionMap maps from the protobuf message name being extended to a set of extensions
	// for the type.
	msgExtensionMap map[string][]protoreflect.ExtensionDescriptor

	// TODO: support enum type definitions for use in future type-check enhancements.
}

// collectFileMetadata traverses the proto file object graph to collect message types and enum
// values and index them by their fully qualified names.
func collectFileMetadata(fileDesc protoreflect.FileDescriptor) *fileMetadata {
	msgTypes := make(map[string]protoreflect.MessageDescriptor)
	enumValues := make(map[string]protoreflect.EnumValueDescriptor)
	msgExtensionMap := make(map[string][]protoreflect.ExtensionDescriptor)
	collectMsgTypes(fileDesc.Messages(), msgTypes, enumValues, msgExtensionMap)
	collectEnumValues(fileDesc.Enums(), enumValues)
	collectExtensions(fileDesc.Extensions(), msgExtensionMap)
	return &fileMetadata{
		msgTypes:        msgTypes,
		enumValues:      enumValues,
		msgExtensionMap: msgExtensionMap,
	}
}

// collectMsgTypes recursively collects messages, nested messages, and nested enums into a map of
// fully qualified protobuf names to descriptors.
func collectMsgTypes(msgTypes protoreflect.MessageDescriptors,
	msgTypeMap map[string]protoreflect.MessageDescriptor,
	enumValueMap map[string]protoreflect.EnumValueDescriptor,
	msgExtensionMap map[string][]protoreflect.ExtensionDescriptor) {
	for i := 0; i < msgTypes.Len(); i++ {
		msgType := msgTypes.Get(i)
		msgTypeMap[string(msgType.FullName())] = msgType
		nestedMsgTypes := msgType.Messages()
		if nestedMsgTypes.Len() != 0 {
			collectMsgTypes(nestedMsgTypes, msgTypeMap, enumValueMap, msgExtensionMap)
		}
		nestedEnumTypes := msgType.Enums()
		if nestedEnumTypes.Len() != 0 {
			collectEnumValues(nestedEnumTypes, enumValueMap)
		}
		nestedExtensions := msgType.Extensions()
		if nestedExtensions.Len() != 0 {
			collectExtensions(nestedExtensions, msgExtensionMap)
		}
	}
}

// collectEnumValues accumulates the enum values within an enum declaration.
func collectEnumValues(enumTypes protoreflect.EnumDescriptors, enumValueMap map[string]protoreflect.EnumValueDescriptor) {
	for i := 0; i < enumTypes.Len(); i++ {
		enumType := enumTypes.Get(i)
		enumTypeValues := enumType.Values()
		for j := 0; j < enumTypeValues.Len(); j++ {
			enumValue := enumTypeValues.Get(j)
			enumValueName := fmt.Sprintf("%s.%s", string(enumType.FullName()), string(enumValue.Name()))
			enumValueMap[enumValueName] = enumValue
		}
	}
}

func collectExtensions(extensions protoreflect.ExtensionDescriptors, msgExtensionMap map[string][]protoreflect.ExtensionDescriptor) {
	for i := 0; i < extensions.Len(); i++ {
		ext := extensions.Get(i)
		extendsMsg := string(ext.ContainingMessage().FullName())
		msgExts, found := msgExtensionMap[extendsMsg]
		if !found {
			msgExts = []protoreflect.ExtensionDescriptor{}
		}
		msgExts = append(msgExts, ext)
		msgExtensionMap[extendsMsg] = msgExts
	}
}
