package gorpc

import (
	"fmt"
    "testing"
)

func TestProtocolEncode(t *testing.T) {
	beginTest("TestProtocolEncode")

//	params := struct{Value string `json:"value"`}{"Hello World"}
	params := []interface{}{"Hello World"}
//	params := make(map[string]interface{})
//	params["value"] = "Hello World"
	
	f := GetFactory()
	
	protocol := f.MakeProtocol()
	rw := f.MakeRequestWrapper()

	result, err := protocol.Encode(rw)
	if err.Error() != "json: error calling MarshalJSON for type *gorpc.RequestWrapper: No requests in wrapper!!" {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	
	rw.AddRequest(f.MakeRequest(1, "Protocol.Encode", params))
	
	result, err = protocol.Encode(rw)
	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	if string(result) != `{"id":1,"method":"Protocol.Encode","params":["Hello World"]}` {
		t.Errorf("Protocol.Encode result does not match: %s", result)
	}

	rw.Clear().AddRequest(f.MakeResponse(1, "Hello World", nil))
	
	result, err = protocol.Encode(rw)
	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	if string(result) != `{"id":1,"result":"Hello World"}` {
		t.Errorf("Protocol.Encode result does not match: %s", result)
	}

	rw.Clear().AddRequest(f.MakeResponse(1, nil, f.MakeRpcError(ErrNotAllowed, nil)))
	
	result, err = protocol.Encode(rw)
	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	if string(result) != `{"error":{"code":-32000,"message":"Not Allowed","data":null},"id":1}` {
		t.Errorf("Protocol.Encode result does not match: %s", result)
	}

	rw.Clear()
	rw.AddRequest(f.MakeResponse(1, "Hello World", nil))
	rw.AddRequest(f.MakeResponse(2, nil, f.MakeRpcError(ErrNotAllowed, nil)))

	result, err = protocol.Encode(rw)
	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	if string(result) != `[{"id":1,"result":"Hello World"},{"error":{"code":-32000,"message":"Not Allowed","data":null},"id":2}]` {
		t.Errorf("Protocol.Encode result does not match: %s", result)
	}
	
	r := f.MakeResponse("1", "Whaahahah", nil)
	r.(*Request).Set("result", float64(42)) 
	s := r.String()
	if s != `Request{data:{id:"1", result:42.000000}}` {
		t.Errorf("Request if formatted wrong: %s", s)
	}
	
	e := make(map[string]interface{})
	e["code"] = float64(ErrInvalidRequest)
	e["message"] = "Invalid Request"
	e["data"] = fmt.Errorf("Test Error")
	
	rm := make(map[string]interface{})
	rm["id"] = 1
	rm["error"] = e
	r = f.MakeRequest(rm, nil, nil)
	s = r.String()
	if s != `Request{data:{error:(*gorpc.RpcError)Code: -32600, Message: Invalid Request, Data: Test Error, id:1}}` {
		t.Errorf("Request if formatted wrong: %s", s)
	}

	
	endTest()
}

func validate(r IRequest) {
	
}


func TestProtocolDecode(t *testing.T) {
	beginTest("TestProtocolDecode")

	s := []byte(`{"id":1,"jsonrpc":"2.0","method":"Protocol.Decode","params":["Hello World"]}`)

	f := GetFactory()
	protocol := f.MakeProtocol()
	protocol.SetValidate(validate)
	
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
	
	rw, err = protocol.Decode([]byte(`[{"id":1,"jsonrpc":"2.0","result":"Hello World"},{"error":{"code":-32000,"message":"Not Allowed","data":null},"id":2,"jsonrpc":"2.0"}]`))
	if err != nil {
		t.Errorf("Decode failed parsing batch request: (%T)%v", err, err)
	}
	if len(rw.GetBatchRequests()) != 2 {
		t.Errorf("Decode result is not correct: (%T)%v", rw,rw)
	}

	rw, err = protocol.Decode([]byte(`["Hello","World"]`))
	if err == nil {
		t.Error("Decode should have failed parsing illegal batch request: (%T)%v", rw, rw)
	}

	rw, err = protocol.Decode([]byte(`{"id":1,"method":"Protocol.Decode","params":["Hello World"]}`))
	iw := protocol.Parse(f.MakeConnection(f.MakeTransport(nil), f.MakeAddress("1", "2", nil)), rw)
	if iw.GetRequest().IsError() == false {
		t.Error("Request should be invalid: (%T)%v", iw,iw)
	}
	
	_, err = protocol.Decode([]byte("WhaaHaHa"))
	if err == nil {
		t.Error("Json should be illegal")
	}
	endTest()
}

