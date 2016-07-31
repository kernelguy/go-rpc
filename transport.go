package gorpc

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
)

type TransportOptions struct {
	
}


type Transport struct {
	connections map[string]IConnection
	protocol IProtocol
	onReceive func(id, message string)
	onWrite func(id, message string)
}


func (this *Transport) Init(onReceive, onWrite func(id, message string)) {
	this.protocol = GetFactory().MakeProtocol()

	if onReceive != nil {
		this.onReceive = onReceive
	}
	if onWrite != nil {
		this.onWrite = onWrite
	}
}


func (p *Transport) Serve(options ITransportOptions) {
	
}

func (p *Transport) Connect(options ITransportOptions) {
	
}

func (p *Transport) Send(id string, rw IRequestWrapper) error {
	message, err := p.protocol.Encode(rw)
	if err != nil {
		return err
	}
	p.Write(id, string(message))
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

func (this *Transport) Receive(id, message string) {
	if this.onReceive == nil {
		this.onReceive = func(id, message string) {
			log.Debugf("Transport.Receive(%s, %s)", id, message)
			c, err := this.getConnection(id)
			if err != nil {
				// We should never get here...
				return
			}
			rw, err := this.protocol.Decode([]byte(message))
			if err != nil {
				// We cannot return an ErrParseError, since we could not decode the request we have no id 
				return
			}
			go func () {
				response := this.protocol.Parse(c, rw)
				if response != nil {
					this.Send(id, response)
				}
			}()
		}
	}
	this.onReceive(id, message)
}

func (this *Transport) Write(id, message string) {
	if this.onWrite == nil {
		this.onWrite = func(id, message string) {
			log.Debugf("Transport.Write(%s, %s)", id, message)
			// Doing nothing here...
		}
	}
	this.onWrite(id, message)
}

