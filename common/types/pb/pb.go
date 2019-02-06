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

// Package pb reflects over protocol buffer descriptors to generate objects
// that simplify type, enum, and field lookup.
package pb

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	descpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	anypb "github.com/golang/protobuf/ptypes/any"
	durpb "github.com/golang/protobuf/ptypes/duration"
	structpb "github.com/golang/protobuf/ptypes/struct"
	tspb "github.com/golang/protobuf/ptypes/timestamp"
	wrapperspb "github.com/golang/protobuf/ptypes/wrappers"
)

// map from file / message / enum name to file description.
type PbDb struct {
	revFileDescriptorMap map[string]*FileDescription
}

var (
	DefaultPbDb = &PbDb{
		revFileDescriptorMap:	make(map[string]*FileDescription),
	}
)

func NewPbDb() *PbDb {
	pbdb := &PbDb{
		revFileDescriptorMap:	make(map[string]*FileDescription),
	}
	for k,v := range DefaultPbDb.revFileDescriptorMap {
		pbdb.revFileDescriptorMap[k] = v
	}
	return pbdb
}

// DescribeEnum takes a qualified enum name and returns an EnumDescription.
func (pbdb *PbDb) DescribeEnum(enumName string) (*EnumDescription, error) {
	enumName = sanitizeProtoName(enumName)
	if fd, found := pbdb.revFileDescriptorMap[enumName]; found {
		return fd.GetEnumDescription(enumName)
	}
	return nil, fmt.Errorf("unrecognized enum '%s'", enumName)
}

// DescribeFile takes a protocol buffer message and indexes all of the message
// types and enum values contained within the message's file descriptor.
func (pbdb *PbDb) DescribeFile(message proto.Message) (*FileDescription, error) {
	if fd, found := pbdb.revFileDescriptorMap[proto.MessageName(message)]; found {
		return fd, nil
	}
	fileDesc, _ := descriptor.ForMessage(message.(descriptor.Message))
	fd, err := pbdb.describeFileInternal(fileDesc)
	if err != nil {
		return nil, err
	}
	pkg := fd.Package()
	fd.indexTypes(pkg, fileDesc.MessageType)
	fd.indexEnums(pkg, fileDesc.EnumType)
	return fd, nil
}

// DescribeType provides a TypeDescription given a qualified type name.
func (pbdb *PbDb) DescribeType(typeName string) (*TypeDescription, error) {
	typeName = sanitizeProtoName(typeName)
	if fd, found := pbdb.revFileDescriptorMap[typeName]; found {
		return fd.GetTypeDescription(typeName)
	}
	return nil, fmt.Errorf("unrecognized type '%s'", typeName)
}

// DescribeValue takes an instance of a protocol buffer message and returns
// the associated TypeDescription.
func (pbdb *PbDb) DescribeValue(value proto.Message) (*TypeDescription, error) {
	fd, err := pbdb.DescribeFile(value)
	if err != nil {
		return nil, err
	}
	typeName := proto.MessageName(value)
	return fd.GetTypeDescription(typeName)
}

func (pbdb *PbDb) describeFileInternal(fileDesc *descpb.FileDescriptorProto) (*FileDescription, error) {
	fd := &FileDescription{
		pbdb:  pbdb,
		desc:  fileDesc,
		types: make(map[string]*TypeDescription),
		enums: make(map[string]*EnumDescription)}
	return fd, nil
}

func fileDescriptor(protoFileName string) (*descpb.FileDescriptorProto, error) {
	gzipped := proto.FileDescriptor(protoFileName)
	r, err := gzip.NewReader(bytes.NewReader(gzipped))
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	unzipped, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	fd := &descpb.FileDescriptorProto{}
	if err := proto.Unmarshal(unzipped, fd); err != nil {
		return nil, fmt.Errorf("bad gzipped descriptor: %v", err)
	}
	return fd, nil
}

func init() {
	// Describe well-known types to ensure they can always be resolved by the check and interpret
	// execution phases.
	//
	// The following subset of message types is enough to ensure that all well-known types can
	// resolved in the runtime, since describing the value results in describing the whole file
	// where the message is declared.
	DefaultPbDb.DescribeValue(&anypb.Any{})
	DefaultPbDb.DescribeValue(&durpb.Duration{})
	DefaultPbDb.DescribeValue(&tspb.Timestamp{})
	DefaultPbDb.DescribeValue(&structpb.Value{})
	DefaultPbDb.DescribeValue(&wrapperspb.BoolValue{})
}
