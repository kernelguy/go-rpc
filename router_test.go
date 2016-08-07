package gorpc

import (
    "testing"
	"time"
)

func TestNonExisting(t *testing.T) {
	beginTest("TestNonExisting")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	conn, _ := transport.getConnection("2")

	_, err := conn.Call("NonExisting", []interface{}{"Hello World"})
	if err == nil {
		t.Errorf("Call should have failed.")
	} else if err.(IRpcError).GetCode() != ErrMethodNotFound {
		t.Errorf("Error is wrong: (%T)%v", err, err)
	} 

	_, err = conn.Call("NonExisting.Echo", []interface{}{"Hello World"})
	if err == nil {
		t.Errorf("Call should have failed.")
	} else if err.(IRpcError).GetCode() != ErrMethodNotFound {
		t.Errorf("Error is wrong: (%T)%v", err, err)
	} 

	transport.quit <- true
	
	endTest()
}

func TestIllegalParams(t *testing.T) {
	beginTest("TestIllegalParams")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	conn := transport.CreateTestConnections()

	_, err := conn.Call("IllegalParams", []interface{}{"Hello World"})
	if err == nil {
		t.Errorf("Call should have failed.")
	} else if err.(IRpcError).GetCode() != ErrInternalError {
		t.Errorf("Error is wrong: (%T)%v", err, err)
	} 

	_, err = conn.Call("NoParams", nil)
	if err != nil {
		t.Errorf("Call to NoParams failed: (%T)%v", err, err)
	}

	_, err = conn.Call("Echo", nil)
	if err == nil {
		t.Errorf("Call should have failed.")
	} else if err.(IRpcError).GetCode() == ErrInternalError {
		t.Errorf("Error is wrong: (%T)%v", err, err)
	} 

	r, _ := conn.Call("TwoParams", []interface{}{1, "Hello"})
	if r.(float64) != 42 {
		t.Errorf("TwoParams result is wrong: (%T)%v", r, r)
	}

	_, err = conn.Call("TwoParams", []interface{}{1})
	if err == nil {
		t.Errorf("Call should have failed.")
	} else if err.(IRpcError).GetCode() == ErrInternalError {
		t.Errorf("Error is wrong: (%T)%v", err, err)
	} 

	_, err = conn.Call("TwoParams", struct{V1 int}{1})
	if err == nil {
		t.Errorf("Call should have failed.")
	} else if err.(IRpcError).GetCode() == ErrInternalError {
		t.Errorf("Error is wrong: (%T)%v", err, err)
	} 

	transport.Stop()

	endTest()
}

func TestErrorResponse(t *testing.T) {
	beginTest("TestErrorResponse")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	conn, _ := transport.getConnection("1")

	r, err := conn.Call("FailWithError", nil)
	if err == nil {
		t.Errorf("Response should be an error: (%T)%v", r, r) 
	} else if err.Error() != "Code: -32603, Message: Internal Error, Data: Chained Error" {
		t.Errorf("Response was not the correct error: (%T)%v", err, err) 
	}

	r, err = conn.Call("FailWithRPCError", nil)
	if err == nil {
		t.Errorf("Response should be an error: (%T)%v", r, r) 
	} else if err.Error() != "Code: -32603, Message: Internal Error, Data: nil" {
		t.Errorf("Response was not the correct error: (%T)%v", err, err) 
	}

	r, err = conn.Call("FailWithString", nil)
	if err == nil {
		t.Errorf("Response should be an error: (%T)%v", r, r) 
	} else if err.Error() != "Code: -32603, Message: Internal Error, Data: Chained String" {
		t.Errorf("Response was not the correct error: (%T)%v", err, err) 
	}

	r, err = conn.Call("FailWithInt", nil)
	if err == nil {
		t.Errorf("Response should be an error: (%T)%v", r, r) 
	} else if err.Error() != "Code: -32603, Message: Internal Error, Data: Unknown: (int)0" {
		t.Errorf("Response was not the correct error: (%T)%v", err, err) 
	}

	transport.quit <- true

	endTest()
}

func TestRpcEcho(t *testing.T) {
	beginTest("TestRpcEcho")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	conn, _ := transport.getConnection("2")

	result, err := conn.RootController().(*testController).Echo("Hello World")

	if err != nil {
		t.Errorf("We should not get an error here: %v", err)
	} else if result != "Hello World" {
		t.Errorf("Result Mismatch: %v", result)
	}

	r, err := conn.Call("Echo", []interface{}{"Hello World"})
	if err != nil {
		r = err.Error()
	}

	if r.(string) != "Hello World" {
		t.Errorf("Result Mismatch: %v", r)
	}

	r, err = conn.Call("Echo", "Hello World")
	if err == nil {
		t.Error("Call should have returned an error.")
	}
	if err.Error() != "Code: -32602, Message: Invalid Params, Data: nil" {
		t.Errorf("Call should have returned an RpcError(Invalid Params): \"%s\"", err.Error())
	}

	transport.quit <- true

	endTest()
}

func TestNotify(t *testing.T) {
	beginTest("TestNotify")

	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Start()

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	conn, _ := transport.getConnection("1")

	conn.Notify("Echo", []interface{}{"Hello World"})

	if len(conn.(*Connection).pendingRequests) > 0 {
		t.Error("There is more than zero pending requests.")
	}

	time.Sleep(time.Millisecond * 10)

	if (transport.LastReceivedMessage.id != "2") || (transport.LastReceivedMessage.data != `{"method":"Echo","params":["Hello World"]}`) {
		t.Errorf("Wrong message received: (%T)%s, %s", transport.LastReceivedMessage, transport.LastReceivedMessage.id, transport.LastReceivedMessage.data)
	}

	transport.quit <- true

	endTest()
}
