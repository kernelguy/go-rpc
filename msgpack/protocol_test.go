package msgpack

import (
	"github.com/kernelguy/gorpc"
    "testing"
)

func TestProtocolEncode(t *testing.T) {
	beginTest("TestProtocolEncode")

//	params := struct{Value string `msgpack:"value"`}{"Hello World"}
	params := []interface{}{"Hello World"}
//	params := make(map[string]interface{})
//	params["value"] = "Hello World"
	
	f := gorpc.GetFactory()
	
	protocol := f.MakeProtocol()
	rw := f.MakeRequestWrapper()
	
	rw.AddRequest(f.MakeRequest(1, "Protocol.Encode", params))
	
	result, err := protocol.Encode(rw)
	
	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	
	if len(result) != 48 {
		t.Errorf("Protocol.Encode result does not match: %d (%T)%x", len(result), result, result)
	}

	rw.Clear().AddRequest(f.MakeResponse(1, "Hello World", nil))
	
	result, err = protocol.Encode(rw)
	
	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}
	
	if len(result) != 24 {
		t.Errorf("Protocol.Encode result does not match: %d (%T)%x", len(result), result, result)
	}

	rw.Clear().AddRequest(f.MakeResponse(1, nil, f.MakeRpcError(gorpc.ErrNotAllowed, nil)))
	
	result, err = protocol.Encode(rw)

	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}

	if len(result) != 46 {
		t.Errorf("Protocol.Encode result does not match: %d (%T)%x", len(result), result, result)
	}

	rw.Clear()
	rw.AddRequest(f.MakeResponse(1, "Hello World", nil))
	rw.AddRequest(f.MakeResponse(2, nil, f.MakeRpcError(gorpc.ErrNotAllowed, nil)))

	result, err = protocol.Encode(rw)

	if err != nil {
		t.Errorf("Protocol.Encode failed: %v", err)
	}

	if len(result) != 71 {
		t.Errorf("Protocol.Encode result does not match: %d (%T)%x", len(result), result, result)
	}

	rw.Clear()
	result, err = protocol.Encode(rw)
	if err == nil {
		t.Errorf("Protocol.Encode should fail: (%T)%v", rw,rw)
	}

	endTest()
}

func TestProtocolDecode(t *testing.T) {
	beginTest("TestProtocolDecode")

	f := gorpc.GetFactory()
	protocol := f.MakeProtocol()

	rw := f.MakeRequestWrapper()
	rw.AddRequest(f.MakeRequest(1, "Protocol.Decode", []interface{}{"Hello World"}))
	s, err := protocol.Encode(rw)

	
	rw, err = protocol.Decode(s)
	
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
	if req.Id().(uint64) != 1 {
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

	rw.Clear()
	rw.AddRequest(f.MakeResponse(1, "Hello World", nil))
	rw.AddRequest(f.MakeResponse(2, nil, f.MakeRpcError(gorpc.ErrNotAllowed, nil)))
	s, err = protocol.Encode(rw)
	
	rw, err = protocol.Decode(s)
	if err != nil {
		t.Errorf("Decode failed parsing batch request: (%T)%v", err, err)
	}
	if len(rw.GetBatchRequests()) != 2 {
		t.Errorf("Decode result is not correct: (%T)%v", rw,rw)
	}

	_, err = protocol.Decode([]byte(`["Hello","World"]`))
	if err == nil {
		t.Error("Decode should have failed parsing illegal batch request: (%T)%v", rw, rw)
	}

	rw.Clear()
	rw.AddRequest(f.MakeRequest(1, "Protocol.Decode", []interface{}{"Hello World"}))
	s, err = protocol.Encode(rw)
	rw, err = protocol.Decode(s)
	iw := protocol.Parse(f.MakeConnection(f.MakeTransport(f.MakeTransportOptions()), f.MakeAddress("1", "2", nil)), rw)
	if iw.GetRequest().IsError() == false {
		t.Error("Request should be invalid: (%T)%v", iw,iw)
	}
	
	_, err = protocol.Decode([]byte("WhaaHaHa"))
	if err == nil {
		t.Error("Message should be illegal")
	}
	
	r := &gorpc.Request{}
	r.SetRequest(1, "Echo", nil)
	rw = f.MakeRequestWrapper()
	rw.AddRequest(r)
	iw = protocol.Parse(f.MakeConnection(f.MakeTransport(f.MakeTransportOptions()), f.MakeAddress("1", "2", nil)), rw)
	if iw.GetRequest().IsError() == false {
		t.Error("Response should be an error: (%T)%v", iw,iw)
	}
	
	rw, err = protocol.Decode([]uint8(""))
	if err == nil {
		t.Errorf("Decode should fail: (%T)%v", rw, rw)
	}
	
	endTest()
}

