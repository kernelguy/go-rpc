package gorpc

import (
    "testing"
)

func TestRequest(t *testing.T) {

	m := make(map[string]interface{})
	m["jsonrpc"] = "2.0"
	m["method"] = "Echo"
	m["params"] = "Hello World"
	m["id"] = 1

	factory := GetFactory()
	req := factory.MakeRequest()
	req.Populate(m)
	
	if req.IsError() {
		t.Error("Request should not be error")
	}

	if req.IsResponse() {
		t.Error("Request should not be response")
	}

	if !req.IsRequest() {
		t.Error("Request should be request")
	}
}

func TestResponse(t *testing.T) {

	m := make(map[string]interface{})
	m["jsonrpc"] = "2.0"
	m["result"] = "Hello World"
	m["id"] = 1

	factory := GetFactory()
	req := factory.MakeRequest()
	req.Populate(m)
	
	if req.IsError() {
		t.Error("Request should not be error")
	}

	if !req.IsResponse() {
		t.Error("Request should be response")
	}

	if req.IsRequest() {
		t.Error("Request should not be request")
	}
}

func TestError(t *testing.T) {

	factory := GetFactory()

	m := make(map[string]interface{})
	m["jsonrpc"] = "2.0"
	m["error"] = factory.MakeRpcError(ErrInvalidParams, nil)
	m["id"] = 1

	req := factory.MakeRequest()
	req.Populate(m)

	if !req.IsError() {
		t.Error("Request should be error")
	}

	if req.IsResponse() {
		t.Error("Request should not be response")
	}

	if req.IsRequest() {
		t.Error("Request should not be request")
	}
}


