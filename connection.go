package gorpc

import (
	log "github.com/Sirupsen/logrus"
	"strconv"
)

type Connection struct {
	ConnectionAddress
	transport ITransport
	pendingRequests map[int]chan interface{}
	rootController interface{}
}

// Global id generator
var _rpc_id int = 0


func (this *Connection) Close() {
	this.transport.(*Transport).removeConnection(this.Source())
}

func (p *Connection) Call(method string, params interface{}) (interface{},error) {
	if p.pendingRequests == nil {
		p.pendingRequests = make(map[int]chan interface{})
	}
	_rpc_id++
	id := _rpc_id
	p.pendingRequests[id] = make(chan interface{})
	r := GetFactory().MakeRequest()
	r.CreateRequest(strconv.Itoa(id), method, params)
	log.Debugf("Connection(%s).Call(%v)", p.Source(), r)
	rw := GetFactory().MakeRequestWrapper()
	rw.AddRequest(r)
	p.transport.Send(p.Destination(), rw)

	log.Debugf("Connection(%s).Call() waiting for result. Id:%d", p.Source(), id)
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
	r := GetFactory().MakeRequest()
	r.CreateRequest(nil, method, params)
	rw := GetFactory().MakeRequestWrapper()
	rw.AddRequest(r)
	p.transport.Send(p.Destination(), rw)
}

func (this *Connection) Response(id int, result interface{}) {
	ch, ok := this.pendingRequests[id]
	if ok  {
		ch<- result
	}
}

func (this *Connection) RootController() interface{} {
	if this.rootController == nil {
		this.rootController = &Controller{connection: this}
	}
	return this.rootController
}

func (this *Connection) SetRootController(obj interface{}) {
	this.rootController = obj
}

