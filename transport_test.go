package gorpc

import (
	"testing"
)

func TestGetConnections(t *testing.T) {
	beginTest("TestGetConnections")

	transport := &Transport{}

	_, err := transport.getConnection("1")
	if err == nil {
		t.Error("We should have got an error")
	} else if err.Error() != "No connections found (1)" {
		t.Errorf("Error message does not match: %v",  err)
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

	result := conn.RootController().(*Controller).Echo("Hello World")

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

	transport.quit <- true
	
	endTest()
}
