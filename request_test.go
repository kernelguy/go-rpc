package gorpc

import (
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


