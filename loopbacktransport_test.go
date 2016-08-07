package gorpc

import (
    "testing"
	log "github.com/Sirupsen/logrus"
)

type TestTransport struct {
	LoopbackTransport
	received chan string
}

func (this *TestTransport) init() {
	this.received = make(chan string)
	this.Protocol = GetFactory().MakeProtocol()
	this.LoopbackTransport.SetFactory(GetFactory())
	this.LoopbackTransport.Init(this.Receive, nil)
}

func (this *TestTransport) Receive(id, message string) {
	this.received <- (id + ":" + message)
}

func TestLoopbackTransport(t *testing.T) {
	beginTest("TestLoopbackTransport")

	trans := &TestTransport{}
	trans.init()
	trans.Start()

	trans.write("1", "Hello")
	
	result := <-trans.received
	if result != "1:Hello" {
		t.Errorf("Result is wrong: (%T)%v", result,result)
	}

	trans.write("2", "World")
	
	result = <-trans.received
	if result != "2:World" {
		t.Errorf("Result is wrong: (%T)%v", result,result)
	}

	log.Debug("Sending quit to LoopbackTransport")
	trans.quit <- true
	
	endTest()
}

