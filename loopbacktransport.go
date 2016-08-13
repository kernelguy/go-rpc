package gorpc

import (
	log "github.com/Sirupsen/logrus"
)

type Message struct {
	id   string
	data string
}

type LoopbackTransport struct {
	Transport
	wire chan Message
	quit chan bool
	LastReceivedMessage Message
}

func (this *LoopbackTransport) Init(onReceive, onWrite func(id, message string)) {
	this.Transport.Init(onReceive, this.write)
}

func (this *LoopbackTransport) CreateTestConnections() IConnection {
	this.addConnection(this.Factory().MakeAddress("1", "2", nil))
	this.addConnection(this.Factory().MakeAddress("2", "1", nil))

	conn, _ := this.getConnection("2")
	return conn
}


func (this *LoopbackTransport) Start() {
	this.wire = make(chan Message)
	this.quit = make(chan bool)

	go func() {
		for {
			select {
			case <-this.quit:
				log.Debug("LoopbackTransport quitting...")
				return
				
			case in := <-this.wire:
				this.LastReceivedMessage = in
				this.Receive(in.id, in.data)
			}
		}
	}()
}

func (this *LoopbackTransport) Stop() {
	this.quit <- true
}

func (this *LoopbackTransport) write(id, message string) {
	log.Debugf("LoopbackTransport.Write(%s, %s)", id, message)
	this.wire <- Message{id: id, data: message}
}
