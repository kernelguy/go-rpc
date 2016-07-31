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
	this.LoopbackTransport.init()
	this.Transport.Init(this.Receive, nil)
}

func (this *TestTransport) Receive(id, message string) {
	this.received <- (id + ":" + message)
}

func TestLoopbackTransport(t *testing.T) {
	beginTest("TestLoopbackTransport")

	trans := &TestTransport{}
	trans.init()
	trans.run()

	trans.Write("1", "Hello World")
	
	result := <-trans.received

	log.Debug("Sending quit to LoopbackTransport")
	trans.quit <- true
	
	if result != "1:Hello World" {
		t.Errorf("Result was incorrect: (%s)", result)
	}
	endTest()
}

