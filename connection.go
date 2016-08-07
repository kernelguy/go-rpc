package gorpc

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
	"sync/atomic"
)

type Connection struct {
	FactoryGetter
	ConnectionAddress
	transport ITransport
	pendingRequests map[uint32]chan interface{}
	rootController IController
}

// Global id generator
var _rpc_id uint32 = 0


func (this *Connection) Init(transport ITransport, addr IConnectionAddress) {
	this.transport = transport
	this.SetFactory(transport.Factory())
	this.ConnectionAddress = *addr.(*ConnectionAddress)
}

func (this *Connection) Close() {
	this.transport.(*Transport).removeConnection(this.Source())
}

func (this *Connection) Call(method string, params interface{}) (interface{},error) {
	if this.pendingRequests == nil {
		this.pendingRequests = make(map[uint32]chan interface{})
	}
	id := atomic.AddUint32(&_rpc_id, 1)
	this.pendingRequests[id] = make(chan interface{})
	f := this.Factory()
	r := f.MakeRequest(id, method, params)
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
	f := this.Factory()
	r := f.MakeRequest(nil, method, params)
	rw := f.MakeRequestWrapper()
	rw.AddRequest(r)
	this.transport.Send(this.Destination(), rw)
}

func (this *Connection) Response(id interface{}, result interface{}) {
	var i uint32
	v := reflect.ValueOf(id).Convert(reflect.TypeOf(i))
	i = v.Interface().(uint32)
	ch, ok := this.pendingRequests[i]
	if ok  {
		ch<- result
	}
}

func (this *Connection) RootController() IController {
	if this.rootController == nil {
		this.SetRootController( this.Factory().MakeController() )
	}
	return this.rootController
}

func (this *Connection) SetRootController(obj IController) {
	obj.SetConnection(this)
	this.rootController = obj
}

