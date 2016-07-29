package gorpc

import (
)

type Connection struct {
	transport *Transport
	name string
	pendingRequests map[int]chan interface{}
}

// Global id generator
var _rpc_id int = 0


func (p *Connection) Close() {
	p.transport.removeConnection(p.name)
}

func (p *Connection) Call(method string, params interface{}) (interface{},error) {
	if p.pendingRequests == nil {
		p.pendingRequests = make(map[int]chan interface{})
	}
	_rpc_id++
	id := _rpc_id
	p.pendingRequests[id] = make(chan interface{})
	r := Request{id: string(id), method: method}
	if params != nil {
		r.params = params
	}
	rw := GetFactory().MakeRequestWrapper()
	rw.AddRequest(&r)
	p.transport.Send(p.name, rw)
	
	var result interface{} = <-p.pendingRequests[id]
	var err error
	
	switch x := result.(type) {
        case error:
            err = x
            result = nil
            
        default:
            err = nil
	}
	return result, err
}

func (p *Connection) Notify(method string, params interface{}) {
	r := Request{method: method}
	if params != nil {
		r.params = params
	}
	rw := GetFactory().MakeRequestWrapper()
	rw.AddRequest(&r)
	p.transport.Send(p.name, rw)
}

func (this *Connection) Response(id int, result interface{}) {
	ch, ok := this.pendingRequests[id]
	if ok  {
		ch<- result
	}
}

func (this *Connection) Rpc_Echo(value string) string {
	return value
}
