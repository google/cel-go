package pb

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/cel-go/test"
	checkedpb "github.com/google/cel-spec/proto/checked/v1/checked"
	"testing"
)

func TestTypeDescription_FieldCount(t *testing.T) {
	td, err := DescribeValue(&test.NestedTestAllTypes{})
	if err != nil {
		t.Error(err)
	}
	if td.FieldCount() != 2 {
		t.Errorf("Unexpected field count. got '%d', wanted '%d'",
			td.FieldCount(), 2)
	}
}

func TestTypeDescription_Field(t *testing.T) {
	td, err := DescribeValue(&test.NestedTestAllTypes{})
	if err != nil {
		t.Error(err)
	}
	fd, found := td.FieldByName("payload")
	if !found {
		t.Error("Field 'payload' not found")
	}
	if fd.OrigName() != "payload" {
		t.Error("Unexpected proto name for field 'payload'", fd.OrigName())
	}
	if fd.Name() != "Payload" {
		t.Error("Unexpected struct name for field 'payload'", fd.Name())
	}
	if fd.GetterName() != "GetPayload" {
		t.Error("Unexpected accessor name for field 'payload'", fd.GetterName())
	}
	if fd.IsOneof() {
		t.Error("Field payload is listed as a oneof and it is not.")
	}
	if fd.IsMap() {
		t.Error("Field 'payload' is listed as a map and it is not.")
	}
	if !fd.IsMessage() {
		t.Error("Field 'payload' is not marked as a message.")
	}
	if fd.IsEnum() {
		t.Error("Field 'payload' is marked as an enum.")
	}
	if fd.IsRepeated() {
		t.Error("Field 'payload' is marked as repeated.")
	}
	if fd.Index() != 1 {
		t.Error("Field 'payload' was fetched at index 1, but not listed there.")
	}
	if !proto.Equal(fd.CheckedType(), &checkedpb.Type{
		TypeKind: &checkedpb.Type_MessageType{
			MessageType: "google.api.tools.expr.test.TestAllTypes"}}) {
		t.Error("Field 'payload' had an unexpected checked type.")
	}
}
