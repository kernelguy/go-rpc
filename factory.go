package gorpc

import (
)


type Factory struct {
}

var instance *Factory

func init() {
	instance = &Factory{}
}

func SetFactory(factory *Factory) {
	instance = factory
}

func GetFactory() IFactory {
    return instance
}


func (this *Factory) MakeTransportOptions() ITransportOptions {
	return &TransportOptions{}
}

func (this *Factory) MakeTransport() ITransport {
	return &Transport{router: this.MakeRouter(), protocol: this.MakeProtocol()}
}

func (this *Factory) MakeConnection(transport ITransport, id string) IConnection {
	return &Connection{name: id, transport: transport.(*Transport)}
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

func (this *Factory) MakeRequest() IRequest {
	return &Request{}
}


