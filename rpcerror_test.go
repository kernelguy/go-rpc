package gorpc

import (
	"fmt"
    "testing"
)

func TestRpcError(t *testing.T) {
	beginTest("TestRpcError")

	e := &RpcError{Code: 42, Message: "MyError"}
	
	if e.GetMessage() != "MyError" {
		t.Errorf("Error message is wrong: %v", e.GetMessage())
	}
	
	if e.GetData() != nil {
		t.Error("Error should be nil")
	}
	
	e.Data = fmt.Errorf("Chained Error")
	
	if e.GetData() == nil {
		t.Error("Error should not be nil")
	} else if e.GetData().Error() != "Chained Error" {
		t.Errorf("Error data is wrong: (%T)%v", e.Data, e.Data)
	}

	endTest()
}

