package gorpc

import (
	"testing"
)

func TestGetConnections(t *testing.T) {
	beginTest("TestGetConnections")

	transport := &Transport{}

	c, err := transport.getConnection("1")
	if err == nil {
		t.Error("We should have got an error")
	} else if err.Error() != "No connections found (1)" {
		t.Errorf("Error message does not match: %v", err)
	}

	transport.addConnection(GetFactory().MakeAddress("1", "2", nil))

	c, err = transport.getConnection("1")
	if err != nil {
		t.Error("We should not have got an error!")
	} else if c.Source() != "1" {
		t.Errorf("Connection.Source does not match: %s", c.Source())
	}

	_, err = transport.getConnection("2")
	if err == nil {
		t.Error("We should have got an error")
	} else if err.Error() != "Connection 2 not found" {
		t.Errorf("Error message does not match: %v", err)
	}

	endTest()
}

func TestAddConnections(t *testing.T) {
	beginTest("TestAddConnections")

	f := GetFactory()
	transport := &Transport{}
	transport.SetFactory(f)
	transport.Init(nil, nil)

	addr := f.MakeAddress("", "", nil)
	addr.SetAddress("1", "2", "yes")
	if addr.Options().(string) != "yes" {
		t.Errorf("Options are wrong: %v", addr.Options())
	}
	transport.addConnection(addr)
	transport.addConnection(f.MakeAddress("2", "1", nil))

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

func TestReceive(t *testing.T) {
	beginTest("TestReceive")

	f := GetFactory()
	transport := &Transport{Protocol: f.MakeProtocol()}
	transport.SetFactory(f)
	transport.Init(nil, nil)

	transport.addConnection(f.MakeAddress("1", "2", nil))
	transport.addConnection(f.MakeAddress("2", "1", nil))

	iw := f.MakeRequestWrapper()
	err := transport.Send("1", iw)
	if err == nil {
		t.Error("We should have got an error.")
	}
	
	// Calls with no check, simply go get code coverage	
	transport.Receive("1","Hello")
	transport.Receive("3", "Hello")
	transport.write("1","Hello")
	
	endTest()
}

