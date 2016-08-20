package gorpc

import (
	"fmt"
    "testing"
)

func TestRequest(t *testing.T) {
	beginTest("TestRequest")

	f := GetFactory()
	req := f.MakeRequest(1, "Echo", "Hello World")
	
	if req.IsError() {
		t.Error("Request should not be error")
	}

	if req.IsResponse() {
		t.Error("Request should not be response")
	}

	if !req.IsRequest() {
		t.Error("Request should be request")
	}
	endTest()
}

func TestResponse(t *testing.T) {
	beginTest("TestResponse")

	f := GetFactory()
	req := f.MakeResponse(1, "Hello World", nil)
	
	if req.IsError() {
		t.Error("Request should not be error")
	}

	if !req.IsResponse() {
		t.Error("Request should be response")
	}

	if req.IsRequest() {
		t.Error("Request should not be request")
	}
	endTest()
}

func TestError(t *testing.T) {
	beginTest("TestError")

	f := GetFactory()

	req := f.MakeResponse(1, nil, f.MakeRpcError(ErrInvalidParams, nil))

	if !req.IsError() {
		t.Error("Request should be error")
	}

	if req.IsResponse() {
		t.Error("Request should not be response")
	}

	if req.IsRequest() {
		t.Error("Request should not be request")
	}
	endTest()
}

func TestUninitialized(t *testing.T) {
	beginTest("TestUninitialized")
	
	r := &Request{}

	if r.GetData() != nil {
		t.Error("Request.data should be nil")
	}
	
	if r.IsError() {
		t.Error("Request should not be an error")
	}

	if r.IsResponse() {
		t.Error("Request should not be a response")
	}

	if r.IsRequest() {
		t.Error("Request should not be a request")
	}

	if r.Id() != nil {
		t.Error("Id should be nil")
	}

	if r.Method() != "" {
		t.Error("Method should be nil")
	}

	if r.Params() != nil {
		t.Error("Params should be nil")
	}

	if r.Result() != nil {
		t.Error("Result should be nil")
	}

	if r.Error() != nil {
		t.Error("Error should be nil")
	}

	r.SetResponse(1, nil, nil)
	s := r.String()
	if s != `Request{data:{id:1, result:nil}}` {
		t.Errorf("Request is not filled with empty: %s", s)
	}

	r = &Request{}
	r.SetResponse(1, nil, GetFactory().MakeRpcError(ErrInvalidRequest, fmt.Errorf("Embedded Error")))
	s = r.String()
	if s != `Request{data:{error:(*gorpc.RpcError)Code: -32600, Message: Invalid Request, Data: Embedded Error, id:1}}` {
 		t.Errorf("Request does not contain an error: %s", s)
	}
	
	m := make(map[string]interface{})
	m["id"] = 666
	m["method"] = "diablo"
	ir := GetFactory().MakeRequest(m, nil, nil)
	if ir.Id() != 666 {
		t.Error("Id should be 666")
	}
}
