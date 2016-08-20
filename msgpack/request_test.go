package msgpack

import (
	"fmt"
	"github.com/kernelguy/gorpc"
    "testing"
)

func TestUninitialized(t *testing.T) {
	beginTest("TestUninitialized")
	
	r := gorpc.GetFactory().MakeRequest(nil, nil, nil).(gorpc.IRequest)

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

	r.SetResponse(1, nil, gorpc.GetFactory().MakeRpcError(gorpc.ErrInvalidRequest, fmt.Errorf("Embedded Error")))
	s = r.String()
	if s != `Request{data:{error:(*gorpc.RpcError)Code: -32600, Message: Invalid Request, Data: Embedded Error, id:1}}` {
 		t.Errorf("Request does not contain an error: %s", s)
	}
	
	m := make(map[string]interface{})
	m["id"] = 666
	m["method"] = "diablo"
	r = gorpc.GetFactory().MakeRequest(m, nil, nil)
	if r.Id() != 666 {
		t.Error("Id should be 666")
	}
}
