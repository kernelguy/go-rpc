package gorpc

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
)

type TransportOptions struct {
	
}


type Transport struct {
	FactoryGetter
	Options ITransportOptions
	connections map[string]IConnection
	Protocol IProtocol
	onReceive func(id, message string)
	onWrite func(id, message string)
}

func (this *Transport) Init(onReceive, onWrite func(id, message string)) {
	if onReceive != nil {
		this.onReceive = onReceive
	} else {
		this.onReceive = this.defaultReceive
	} 
	
	if onWrite != nil {
		this.onWrite = onWrite
	} else {
		this.onWrite = this.defaultWrite
	}
}


func (p *Transport) Start() {
	
}

func (p *Transport) Send(id string, rw IRequestWrapper) error {
	log.Debugf("Transport.Send(%s, %v)", id, rw)
	message, err := p.Protocol.Encode(rw)
	if err != nil {
		return err
	}
	p.write(id, string(message))
	return nil
}

func (this *Transport) Close() {
	
}

func (p *Transport) addConnection(addr IConnectionAddress) IConnection {
	if (p.connections == nil) {
		p.connections = make(map[string]IConnection)
	}
	c := GetFactory().MakeConnection(p, addr)
	p.connections[addr.Source()] = c
	return c
}

func (p *Transport) removeConnection(id string) {
	_, err := p.getConnection(id)
	if err == nil {
		delete(p.connections, id)
	}
}

func (p *Transport) getConnection(id string) (IConnection, error) {
	if (p.connections == nil) {
		return nil, fmt.Errorf("No connections found (%s)", id)
	}
	c := p.connections[id]
	if c == nil {
		return nil, fmt.Errorf("Connection %s not found", id)
	}
	return c, nil
}

func (this *Transport) defaultReceive(id, message string) {
	log.Debugf("Transport.defaultReceive(%s, %s)", id, message)
	c, err := this.getConnection(id)
	//log.Debug("step 1")
	if err != nil {
		log.Debug("No connection for received message. Aborting...")
		return
	}
	//log.Debug("step 2")
	rw, err := this.Protocol.Decode([]byte(message))
	//log.Debug("step 3")
	if err != nil {
		log.Debug("Incoming message could not decode. Aborting...")
		// We cannot return an ErrParseError, since we could not decode the request we have no id 
		return
	}
	//log.Debug("step 4")
	go func () {
		response := this.Protocol.Parse(c, rw)
		if response != nil {
			this.Send(c.Destination(), response)
		}
	}()
}

func (this *Transport) defaultWrite(id, message string) {
	log.Debugf("Transport.defaultWrite(%s, %s)", id, message)
	// Doing nothing here...
}


func (this *Transport) Receive(id, message string) {
	this.onReceive(id, message)
}

func (this *Transport) write(id, message string) {
	this.onWrite(id, message)
}
