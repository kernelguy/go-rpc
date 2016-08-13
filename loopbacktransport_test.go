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

func emptyReceive(id, message string) {
	
}

func BenchmarkTransport(b *testing.B) {
	var count int = 0
	f := GetFactory()
	transport := f.MakeTransport(nil).(*LoopbackTransport)
	transport.Init(func(id, message string) {
		count--
	}, nil)
	transport.Start()

	b.ResetTimer()
	
	for i:=0; i < b.N; i++ {
		count++
		transport.write("Id", "Message");
	}

	transport.quit <- true

	if count != 0 {
		b.Errorf("Count should be zero: %v", count)
	}
}

