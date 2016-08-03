package gorpc

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
)

type Connection struct {
	ConnectionAddress
	transport ITransport
	pendingRequests map[int]chan interface{}
	rootController IController
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
	f := GetFactory()
	r := f.MakeRequest(nil, method, params)
	rw := f.MakeRequestWrapper()
	rw.AddRequest(r)
	this.transport.Send(this.Destination(), rw)
}

func (this *Connection) Response(id interface{}, result interface{}) {
	var i int
	v := reflect.ValueOf(id).Convert(reflect.TypeOf(i))
	i = v.Interface().(int)
	ch, ok := this.pendingRequests[i]
	if ok  {
		ch<- result
	}
}

func (this *Connection) RootController() IController {
	if this.rootController == nil {
		this.SetRootController( GetFactory().MakeController() )
	}
	return this.rootController
}

func (this *Connection) SetRootController(obj IController) {
	obj.SetConnection(this)
	this.rootController = obj
}

