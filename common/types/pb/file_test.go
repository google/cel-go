package pb

import (
	"testing"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	proto3pb "github.com/google/cel-go/test/proto3pb"
	descpb "google.golang.org/protobuf/types/descriptorpb"
)

func TestFileDescriptionGetTypes(t *testing.T) {
	pbdb := NewDb()
	fd, err := pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Error(err)
	}
	expected := []string{
		"google.expr.proto3.test.TestAllTypes",
		"google.expr.proto3.test.TestAllTypes.NestedMessage",
		"google.expr.proto3.test.TestAllTypes.MapStringStringEntry",
		"google.expr.proto3.test.TestAllTypes.MapInt64NestedTypeEntry",
		"google.expr.proto3.test.NestedTestAllTypes"}
	if len(fd.GetTypeNames()) != len(expected) {
		t.Errorf("got '%v', wanted '%v'", fd.GetTypeNames(), expected)
	}
	for _, tn := range fd.GetTypeNames() {
		var found = false
		for _, elem := range expected {
			if elem == tn {
				found = true
				break
			}
		}
		if !found {
			t.Error("Unexpected type name", tn)
		}
	}
	for _, typeName := range fd.GetTypeNames() {
		td, found := fd.GetTypeDescription(typeName)
		if !found {
			t.Fatalf("fd.GetTypeDescription(%v) returned not found", typeName)
		}
		if td.Name() != typeName {
			t.Error("Indexed type name not equal to descriptor type name", td, typeName)
		}
	}
}

func TestFileDescriptionGetEnumNames(t *testing.T) {
	pbdb := NewDb()
	fd, err := pbdb.RegisterMessage(&proto3pb.TestAllTypes{})
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]int32{
		"google.expr.proto3.test.TestAllTypes.NestedEnum.FOO": 0,
		"google.expr.proto3.test.TestAllTypes.NestedEnum.BAR": 1,
		"google.expr.proto3.test.TestAllTypes.NestedEnum.BAZ": 2,
		"google.expr.proto3.test.GlobalEnum.GOO":              0,
		"google.expr.proto3.test.GlobalEnum.GAR":              1,
		"google.expr.proto3.test.GlobalEnum.GAZ":              2}
	if len(expected) != len(fd.GetEnumNames()) {
		t.Error("Count of enum names does not match expected'",
			fd.GetEnumNames(), expected)
	}
	for _, enumName := range fd.GetEnumNames() {
		if enumVal, found := expected[enumName]; found {
			ed, found := fd.GetEnumDescription(enumName)
			if !found {
				t.Fatalf("fd.GetEnumDescription(%v) returned not found", enumName)
			}
			if ed.Value() != enumVal {
				t.Errorf("Enum did not have expected value. %s got '%v', wanted '%v'",
					enumName, ed.Value(), enumVal)
			}
		} else {
			t.Errorf("Enum value not found for: %s", enumName)
		}
	}
}

func TestFileDescriptionGetImportedEnumNames(t *testing.T) {
	pbdb := NewDb()
	fdMap := CollectFileDescriptorSet(&proto3pb.TestAllTypes{})
	fileSet := make([]*descpb.FileDescriptorProto, 0, len(fdMap))
	for _, fd := range fdMap {
		fileSet = append(fileSet, protodesc.ToFileDescriptorProto(fd))
	}
	fds := &descpb.FileDescriptorSet{File: fileSet}
	files, err := protodesc.NewFiles(fds)
	if err != nil {
		t.Fatalf("protodesc.NewFiles() failed: %v", err)
	}
	files.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		t.Logf("registering file: %v", fd.Path())
		_, err = pbdb.RegisterDescriptor(fd)
		if err != nil {
			t.Fatalf("pbdb.RegisterDescriptor(%v) failed: %v", fd, err)
			return false
		}
		return true
	})
	imported := map[string]int32{
		"google.expr.proto3.test.ImportedGlobalEnum.IMPORT_FOO": 0,
		"google.expr.proto3.test.ImportedGlobalEnum.IMPORT_BAR": 1,
		"google.expr.proto3.test.ImportedGlobalEnum.IMPORT_BAZ": 2,
	}
	for enumName, value := range imported {
		ed, found := pbdb.DescribeEnum(enumName)
		if !found {
			t.Fatalf("pbdb.DescribeEnum(%q) returned not found", enumName)
		}
		if ed.Value() != value {
			t.Errorf("Got %v, wanted %v for enum %s", ed, value, enumName)
		}
	}
}
