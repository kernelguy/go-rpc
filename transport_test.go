package gorpc

import (
    "testing"
)

type Message struct {
	id string
	data string
}

type MyTransport struct {
	Transport
	wire chan Message
	terminated bool
}


func TestGetConnections(t *testing.T) {
	transport := &Transport{}

	_, err := transport.getConnection("1")
	if err == nil {
		t.Error("We should have got an error")
	} else if err.Error() != "No connections found (1)" {
		t.Error("Error message does not match: " + err.Error())
	}

	transport.addConnection("2")

	_, err2 := transport.getConnection("1")
	if err2 == nil {
		t.Error("We should have got an error!")
	} else if err2.Error() != "Connection 1 not found" {
		t.Error("Error message does not match: " + err2.Error())
	}
}


func TestAddConnections(t *testing.T) {
	transport := &Transport{}
	
	transport.addConnection("1")
	transport.addConnection("2")
	
	connection, err := transport.getConnection("1")
	if err != nil {
		t.Error(err)
	}
	if connection.(*Connection).name != "1" {
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
	if connection2.(*Connection).name != "2" {
		t.Fatal("Connection name does not match!")
	}

	transport.removeConnection("2")
		
	if len(transport.connections) != 0 {
		t.Fatal("Connections count should be zero!")
	}
}

func (p *MyTransport) run() {
	p.wire = make(chan Message, 10)
	p.terminated = false
	
	go func() {
		var in Message
		for !p.terminated {
			in = <-p.wire
			p.receive(in.id, in.data)
		}
	}()
}

func (p *MyTransport) write(id, message string) {
	go func() {
		out := Message{id: id, data: message}
		p.wire <- out
	}()
}

func TestRpcEcho( t * testing.T) {
	transport := &MyTransport{}
	transport.run()

	transport.addConnection("1")
	transport.addConnection("2")

	conn, _ := transport.getConnection("2");
	
	result, _ := conn.Call("Echo", "Hello World")
	
	transport.terminated = true
	
	if result.(string) != "Hello World" {
		t.Errorf("Result Mismatch: " + result.(string))
	}
}

