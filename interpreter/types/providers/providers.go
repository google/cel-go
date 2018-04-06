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

// Package providers declares utilities for creating / resolving external types
// as expression values.
package providers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/google/cel-go/interpreter/types"
	"io/ioutil"
	"reflect"
)

// TypeProvider specifies functions for creating new object instances and for
// resolving enum values by name.
type TypeProvider interface {
	// Create a new instance of the qualified type name and initialize the
	// fields with the values provided.
	//
	// Handles conversion of expression types to proto types.
	NewValue(typeName string, fields map[string]interface{}) (types.ObjectValue, error)

	// Resolve a qualified enum name to its integer value, return false if the
	// name could not be resolved.
	EnumValue(enumName string) (int64, bool)
}

type defaultTypeProvider struct {
	typesByName map[string]reflect.Type
	enumsByName map[string]int64
}

// NewTypeProvider accepts a list of proto message instances and returns a type
// provider which can create new instances of the provided message or any
// message that proto depends upon in its FileDescriptor.
func NewTypeProvider(types ...proto.Message) TypeProvider {
	protoTypes := make(map[string]reflect.Type)
	enumValues := make(map[string]int64)
	descriptorSet := make(map[string]*protobuf.FileDescriptorProto)
	for _, protoType := range types {
		fileDesc, _ := descriptor.ForMessage(protoType.(descriptor.Message))
		descriptorSet[fileDesc.GetName()] = fileDesc
		buildDescriptorSet(fileDesc, descriptorSet)
	}
	for _, fileDesc := range descriptorSet {
		buildTypeInfo(fileDesc.GetPackage(), fileDesc.MessageType, protoTypes, enumValues)
		buildEnumInfo(fileDesc.GetPackage(), fileDesc.EnumType, enumValues)
	}
	return &defaultTypeProvider{protoTypes, enumValues}
}

func (tp *defaultTypeProvider) EnumValue(enumName string) (int64, bool) {
	enumVal, found := tp.enumsByName[enumName]
	return enumVal, found
}

func (tp *defaultTypeProvider) NewValue(typeName string,
	fields map[string]interface{}) (types.ObjectValue, error) {
	if refType, found := tp.typesByName[typeName]; found {
		// create the new type instance.
		value := reflect.New(refType.Elem())
		pbValue := value.Elem()

		// for all of the field names referenced, set the provided value.
		for fieldName, fieldValue := range fields {
			refField := pbValue.FieldByName(fieldName)
			if !refField.IsValid() {
				// TODO: fix the error message
				return nil, fmt.Errorf("no such field")
			}
			if refFieldValue, err :=
				types.ExprToProto(refField.Type(), fieldValue); err == nil {
				refField.Set(reflect.ValueOf(refFieldValue))
			} else {
				return nil, err
			}
		}
		return types.NewProtoValue(value.Interface()), nil
	} else {
		return nil, fmt.Errorf("unknown type '%s'", typeName)
	}
}

func buildTypeInfo(packageName string, protoMsgTypes []*protobuf.DescriptorProto,
	msgTypes map[string]reflect.Type,
	enumValues map[string]int64) {
	for _, msgType := range protoMsgTypes {
		msgName := fmt.Sprintf("%s.%s", packageName, msgType.GetName())
		msgTypes[msgName] = proto.MessageType(msgName)
		buildTypeInfo(packageName, msgType.NestedType, msgTypes, enumValues)
		buildEnumInfo(packageName, msgType.EnumType, enumValues)
	}
}

func buildEnumInfo(packageName string, protoEnumTypes []*protobuf.EnumDescriptorProto,
	enumValues map[string]int64) {
	for _, enumType := range protoEnumTypes {
		for _, enumValue := range enumType.GetValue() {
			// Embeds the fully qualified name into the enum values map
			enumName := fmt.Sprintf("%s.%s", packageName, enumValue.String())
			enumValues[enumName] = int64(enumValue.GetNumber())
		}
	}
}

func buildDescriptorSet(fileDesc *protobuf.FileDescriptorProto,
	descriptorSet map[string]*protobuf.FileDescriptorProto) {
	descriptorSet[fileDesc.GetName()] = fileDesc
	for _, protoFileName := range fileDesc.Dependency {
		if _, found := descriptorSet[protoFileName]; !found {
			if fd, err := fileDescriptor(protoFileName); err != nil {
				panic(err)
			} else {
				buildDescriptorSet(fd, descriptorSet)
			}
		}
	}
}

func fileDescriptor(protoFileName string) (*protobuf.FileDescriptorProto, error) {
	gzipped := proto.FileDescriptor(protoFileName)
	r, err := gzip.NewReader(bytes.NewReader(gzipped))
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	unzipped, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	fd := &protobuf.FileDescriptorProto{}
	if err := proto.Unmarshal(unzipped, fd); err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	return fd, nil
}
