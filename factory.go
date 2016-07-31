package gorpc

import (
)


type Factory struct {
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


func (this *Factory) MakeTransportOptions() ITransportOptions {
	return &TransportOptions{}
}

func (this *Factory) MakeTransport() ITransport {
	t := &Transport{}
	t.Init(nil, nil)
	return t
}

func (this *Factory) MakeAddress(src, dest string, options interface{}) IConnectionAddress {
	return &ConnectionAddress{src: src, dest: dest, options: options}
}

func (this *Factory) MakeConnection(transport ITransport, addr IConnectionAddress) IConnection {
	result := &Connection{transport: transport.(*Transport)}
	result.SetAddress(addr.Source(), addr.Destination(), addr.Options())
	return result
}

func (this *Factory) MakeProtocol() IProtocol {
	return &Protocol{}
}

func (this *Factory) MakeRouter() IRouter {
	return &Router{}
}

func (this *Factory) MakeRpcError(code int, previous error) IRpcError {
	return NewRpcError(code, previous)
}

func (this *Factory) MakeRequestWrapper() IRequestWrapper {
	return &RequestWrapper{}
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

func (this *Factory) MakeResponse(id, result, error interface{}) IRequest {
	r := &Request{}
	r.SetResponse(id, result, error)
	return r
}

