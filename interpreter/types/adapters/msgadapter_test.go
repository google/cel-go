package adapters

import (
	testExpr "github.com/google/cel-go/interpreter/testing"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"testing"
)

func TestMsgAdapter_Get(t *testing.T) {
	existsMsg := NewMsgAdapter(testExpr.Exists.Expr)
	compreExpr, err := existsMsg.Get("ComprehensionExpr")
	if err != nil {
		t.Error(err)
	}
	iterVar, err := compreExpr.(MsgAdapter).Get("IterVar")
	if err != nil {
		t.Error(err)
	}
	if iterVar != "x" {
		t.Error("Could not retrieve iter var from comprehension")
	}
	// This field is not set, but should return the default instance.
	selectExpr, err := existsMsg.Get("SelectExpr")
	if err != nil {
		t.Error(err)
	}
	field, err := selectExpr.(MsgAdapter).Get("Field")
	if field != "" {
		t.Error("Selected field on unset message, but field was non-empty")
	}
	field2, err := selectExpr.(MsgAdapter).Get("Field")
	if field2 != field {
		t.Error("Selected cached field not equal to original field")
	}
}

func TestMsgAdapter_Any(t *testing.T) {
	anyVal, err := ptypes.MarshalAny(testExpr.Exists.Expr)
	if err != nil {
		t.Error(err)
	}
	existsMsg := NewMsgAdapter(anyVal)
	compre, err := existsMsg.Get("ComprehensionExpr")
	if err != nil {
		t.Error(err)
	}
	if compre == nil {
		t.Error("Comprehension was null")
	}
}

func TestMsgAdapter_Iterator(t *testing.T) {
	existsMsg := NewMsgAdapter(testExpr.Exists.Expr)
	it := existsMsg.Iterator()
	for it.HasNext() {
		fmt.Println(it.Next())
	}
}
