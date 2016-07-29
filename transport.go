package gorpc

import (
	"fmt"
)

type TransportOptions struct {
	
}


type Transport struct {
	connections map[string]IConnection
	router IRouter
	protocol IProtocol
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
	p.write(id, string(message))
	return nil
}

func (p *Transport) Close() {
	
}

func (p *Transport) addConnection(id string) IConnection {
	if (p.connections == nil) {
		p.connections = make(map[string]IConnection)
	}
	c := GetFactory().MakeConnection(p, id)
	p.connections[id] = c
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

func (p *Transport) receive(id, message string) {
	c, err := p.getConnection(id)
	if err != nil {
		// We should never get here...
		return
	}
	rw, err := p.protocol.Decode([]byte(message))
	if err != nil {
		// We cannot return a cParseError, since we could not decode the request we have no id 
		return
	}
	go func () {
		response := p.router.Route(c, rw)
		if response != nil {
			p.Send(id, response)
		}
	}()
}

func (p *Transport) write(id, message string) {
	
}

