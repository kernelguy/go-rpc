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
}

func (this *LoopbackTransport) init() {
	this.Transport.Init(nil, this.Write)
}

func (this *LoopbackTransport) run() {
	this.wire = make(chan Message)
	this.quit = make(chan bool)

	go func() {
		for {
			select {
			case <-this.quit:
				log.Debug("LoopbackTransport quitting...")
				return
				
			case in := <-this.wire:
				this.Receive(in.id, in.data)
			}
		}
	}()

}

func (this *LoopbackTransport) Write(id, message string) {
	go func() {
		log.Debugf("LoopbackTransport.Write(%s, %s)", id, message)
		out := Message{id: id, data: message}
		this.wire <- out
	}()
}
