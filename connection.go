package gorpc

import (
	log "github.com/Sirupsen/logrus"
	"fmt"
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

func (this *Connection) Call(method string, params interface{}) (interface{},error) {
	if this.pendingRequests == nil {
		this.pendingRequests = make(map[int]chan interface{})
	}
	_rpc_id++
	id := _rpc_id
	this.pendingRequests[id] = make(chan interface{})
	f := GetFactory()
	r := f.MakeRequest(strconv.Itoa(id), method, params)
	log.Debugf("Connection(%s).Call(%v)", this.Source(), r)
	rw := f.MakeRequestWrapper()
	rw.AddRequest(r)
	this.transport.Send(this.Destination(), rw)

	log.Debugf("Connection(%s).Call() waiting for result. Id:%d", this.Source(), id)
	var result interface{} = <-this.pendingRequests[id]
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

func (this *Connection) Notify(method string, params interface{}) {
	f := GetFactory()
	r := f.MakeRequest(nil, method, params)
	rw := f.MakeRequestWrapper()
	rw.AddRequest(r)
	this.transport.Send(this.Destination(), rw)
}

func (this *Connection) Response(id interface{}, result interface{}) {
	var i int
	switch id.(type) {
		case string:
			i, _ = strconv.Atoi(id.(string))
		case int, int64:
			i = id.(int)
		default:
			panic(fmt.Errorf("RPC response to this package should always be integers, not: (%T)%v", id, id))
	}
	ch, ok := this.pendingRequests[i]
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

