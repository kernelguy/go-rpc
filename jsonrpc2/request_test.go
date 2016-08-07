package jsonrpc2

import (
	"fmt"
	"github.com/kernelguy/gorpc"
    "testing"
)

func TestUninitialized(t *testing.T) {
	beginTest("TestUninitialized")
	
	r := gorpc.GetFactory().MakeRequest(nil, nil, nil).(gorpc.IJsonRPC2Request)

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

	if r.JsonRPC() != "2.0" {
		t.Error("JsonRPC should be 2.0")
	}

	r.SetResponse(1, nil, nil)
	s := r.String()
	if s != `Request{data:{id:1, jsonrpc:"2.0", result:nil}}` {
		t.Errorf("Request is not filled with empty: %s", s)
	}

	r = &Request{}

	if r.JsonRPC() != "" {
		t.Error("JsonRPC should be nil")
	}

	r.SetResponse(1, nil, gorpc.GetFactory().MakeRpcError(gorpc.ErrInvalidRequest, fmt.Errorf("Embedded Error")))
	s = r.String()
	if s != `Request{data:{error:(*gorpc.RpcError)Code: -32600, Message: Invalid Request, Data: Embedded Error, id:1, jsonrpc:"2.0"}}` {
 		t.Errorf("Request does not contain an error: %s", s)
	}
}
