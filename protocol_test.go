package gorpc

import (
    "testing"
)

func TestProtocolEncode(t *testing.T) {

//	params := struct{Value string `json:"value"`}{"Hello World"}
	params := []interface{}{"Hello World"}
//	params := make(map[string]interface{})
//	params["value"] = "Hello World"
	
	m := make(map[string]interface{})
	m["jsonrpc"] = "2.0"
	m["method"] = "Protocol.Encode"
	m["params"] = params
	m["id"] = 1

	factory := GetFactory()
	
	protocol := factory.MakeProtocol()
	req := factory.MakeRequest()
	rw := factory.MakeRequestWrapper()
	
	req.Populate(m)
	rw.AddRequest(req)
	
	result, err := protocol.Encode(rw)
	
	if err != nil {
		t.Errorf("Encode failed: %v", err)
	}
	
	if string(result) != `{"id":1,"jsonrpc":"2.0","method":"Protocol.Encode","params":["Hello World"]}` {
		t.Errorf("Encode result does not match: %s", result)
	}
}

func TestProtocolDecode(t *testing.T) {
	s := []byte(`{"id":1,"jsonrpc":"2.0","method":"Protocol.Decode","params":["Hello World"]}`)

	factory := GetFactory()
	protocol := factory.MakeProtocol()
	
	rw, err := protocol.Decode(s)
	
	if err != nil {
		t.Errorf("Decode failed: %+v", err)
	}

	if rw.IsEmpty() {
		t.Error("Decode result should not be empty")
	}
	if rw.IsBatchRequest() {
		t.Error("Decode result should not be a batch request")
	}

	req := rw.GetRequest()
	if req.Id().(float64) != 1 {
		t.Errorf("Decode result request id does not match")
	}
	if req.Method() != "Protocol.Decode" {
		t.Error("Decode result request method does not match")
	}
	p := req.Params()
	params, ok := p.([]interface{})
	if !ok {
		t.Error("Decode result request params is not an array")
	}
	if params[0] != "Hello World" {
		t.Error("Decode result request params does not match")
	}
	
}

