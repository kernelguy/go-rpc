package gorpc

import (
)


type Factory struct {
}

type FactoryGetter struct {
	factory IFactory
}

var instance IFactory

func init() {
	instance = &Factory{}
}

func SetFactory(factory IFactory) {
	instance = factory
}

func GetFactory() IFactory {
    return instance
}



func (this *Factory) MakeAddress(src, dest string, options interface{}) IConnectionAddress {
	return &ConnectionAddress{src: src, dest: dest, options: options}
}

func (this *Factory) MakeConnection(transport ITransport, addr IConnectionAddress) IConnection {
	c := &Connection{}
	c.Init(transport, addr)
	return c
}

func (this *Factory) MakeController() IController {
	return &Controller{}
}

func (this *Factory) MakeProtocol() IProtocol {
	p := &Protocol{}
	p.SetFactory(this)
	return p
}

func (this *Factory) MakeRequest(id, method, params interface{}) IRequest {
	r := &Request{}
	if v, ok := id.(map[string]interface{}); ok {
		r.Populate(v)
	} else {
		r.SetRequest(id, method, params)
	}
	return r
}

func (this *Factory) MakeRequestWrapper() IRequestWrapper {
	return &RequestWrapper{}
}

func (this *Factory) MakeResponse(id, result, error interface{}) IRequest {
	r := &Request{}
	r.SetResponse(id, result, error)
	return r
}

/*
	Make a router takes a protocol validation function as an argument. 
	The validator can be nill if not needed.
*/
func (this *Factory) MakeRouter() IRouter {
	r := &Router{}
	r.SetFactory(this)
	return r
}

func (this *Factory) MakeRpcError(code int, previous error) IRpcError {
	return NewRpcError(code, previous)
}

func (this *Factory) MakeTransport(options ITransportOptions) ITransport {
	t := &Transport{Options: options, Protocol: this.MakeProtocol()}
	t.SetFactory(this)
	t.Init(nil, nil)
	return t
}

func (this *Factory) MakeTransportOptions() ITransportOptions {
	return &TransportOptions{}
}


func (this *FactoryGetter) SetFactory(factory IFactory) {
	this.factory = factory
}

func (this *FactoryGetter) Factory() IFactory {
	return this.factory
}
