package gorpc

import (
	"testing"
	"time"
)

func TestGetConnections(t *testing.T) {
	beginTest("TestGetConnections")

	transport := &Transport{}

	_, err := transport.getConnection("1")
	if err == nil {
		t.Error("We should have got an error")
	} else if err.Error() != "No connections found (1)" {
		t.Errorf("Error message does not match: %v", err)
	}

	transport.addConnection(GetFactory().MakeAddress("1", "2", nil))

	c, err2 := transport.getConnection("1")
	if err2 != nil {
		t.Error("We should not have got an error!")
	} else if c.Source() != "1" {
		t.Errorf("Connection.Source does not match: %s", c.Source())
	}
	endTest()
}

func TestAddConnections(t *testing.T) {
	beginTest("TestAddConnections")

	transport := &Transport{}

	transport.addConnection(GetFactory().MakeAddress("1", "2", nil))
	transport.addConnection(GetFactory().MakeAddress("2", "1", nil))

	connection, err := transport.getConnection("1")
	if err != nil {
		t.Error(err)
	}
	if connection.(*Connection).Source() != "1" {
		t.Fatal("Connection name does not match!")
	}

	if len(transport.connections) != 2 {
		t.Fatal("Connections count does not match!")
	}

	transport.removeConnection("1")

	if len(transport.connections) != 1 {
		t.Fatal("Connections count should be zero!")
	}

	connection2, err := transport.getConnection("2")
	if err != nil {
		t.Error(err)
	}
	if connection2.(*Connection).Source() != "2" {
		t.Fatal("Connection name does not match!")
	}

	transport.removeConnection("2")

	if len(transport.connections) != 0 {
		t.Fatal("Connections count should be zero!")
	}
	endTest()
}

func TestRpcEcho(t *testing.T) {
	beginTest("TestRpcEcho")

	transport := &LoopbackTransport{}
	transport.init()
	transport.run()

	transport.addConnection(GetFactory().MakeAddress("1", "2", nil))
	transport.addConnection(GetFactory().MakeAddress("2", "1", nil))

	conn, _ := transport.getConnection("2")

	result, err := conn.RootController().(*Controller).Echo("Hello World")

	if result != "Hello World" {
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

	transport := &LoopbackTransport{}
	transport.init()
	transport.run()

	transport.addConnection(GetFactory().MakeAddress("1", "2", nil))
	transport.addConnection(GetFactory().MakeAddress("2", "1", nil))

	conn, _ := transport.getConnection("1")

	conn.Notify("Echo", []interface{}{"Hello World"})

	if len(conn.(*Connection).pendingRequests) > 0 {
		t.Error("There is more that zero pending requests.")
	}

	time.Sleep(time.Millisecond * 10)

	if (transport.LastReceivedMessage.id != "2") || (transport.LastReceivedMessage.data != `{"jsonrpc":"2.0","method":"Echo","params":["Hello World"]}`) {
		t.Errorf("Wrong message received: (%T)%s, %s", transport.LastReceivedMessage, transport.LastReceivedMessage.id, transport.LastReceivedMessage.data)
	}

	f := GetFactory()
	r := f.MakeRequest(nil, "Notify", nil)
	r.(*Request).data["jsonrpc"] = "1.0"
	rw := f.MakeRequestWrapper()
	rw.AddRequest(r)
	transport.Send(conn.Destination(), rw)

	time.Sleep(time.Millisecond * 10)


	transport.quit <- true

	endTest()
}
