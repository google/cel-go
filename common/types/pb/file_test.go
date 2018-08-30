package pb

import (
	"testing"

	testpb "github.com/google/cel-go/test"
)

func TestFileDescription_GetTypes(t *testing.T) {
	fd, err := DescribeFile(&testpb.TestAllTypes{})
	if err != nil {
		t.Error(err)
	}
	expected := []string{
		"google.api.tools.expr.test.TestAllTypes",
		"google.api.tools.expr.test.TestAllTypes.NestedMessage",
		"google.api.tools.expr.test.TestAllTypes.MapStringStringEntry",
		"google.api.tools.expr.test.TestAllTypes.MapInt64NestedTypeEntry",
		"google.api.tools.expr.test.NestedTestAllTypes"}
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
		td, err := fd.GetTypeDescription(typeName)
		if err != nil {
			t.Error(err)
		}
		if td.Name() != typeName {
			t.Error("Indexed type name not equal to descriptor type name", td, typeName)
		}
		if td.file != fd {
			t.Error("Indexed type does not refer to current file", td)
		}
	}
}

func TestFileDescription_GetEnumNames(t *testing.T) {
	fd, err := DescribeFile(&testpb.TestAllTypes{})
	if err != nil {
		t.Error(err)
	}
	expected := map[string]int32{
		"google.api.tools.expr.test.TestAllTypes.NestedEnum.FOO": 0,
		"google.api.tools.expr.test.TestAllTypes.NestedEnum.BAR": 1,
		"google.api.tools.expr.test.TestAllTypes.NestedEnum.BAZ": 2,
		"google.api.tools.expr.test.GlobalEnum.GOO":              0,
		"google.api.tools.expr.test.GlobalEnum.GAR":              1,
		"google.api.tools.expr.test.GlobalEnum.GAZ":              2}
	if len(expected) != len(fd.GetEnumNames()) {
		t.Error("Count of enum names does not match expected'",
			fd.GetEnumNames(), expected)
	}
	for _, enumName := range fd.GetEnumNames() {
		if enumVal, found := expected[enumName]; found {
			ed, err := fd.GetEnumDescription(enumName)
			if err != nil {
				t.Error(err)
			} else if ed.Value() != enumVal {
				t.Errorf("Enum did not have expected value. %s got '%v', wanted '%v'",
					enumName, ed.Value(), enumVal)
			}
		} else {
			t.Errorf("Enum value not found for: %s", enumName)
		}
	}
}
