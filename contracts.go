package gorpc

import (

)

type IFactory interface {
	MakeTransportOptions() ITransportOptions
	MakeTransport() ITransport
	MakeConnection(transport ITransport, id string) IConnection
	MakeProtocol() IProtocol
	MakeRouter() IRouter
	MakeRpcError(code int, previous error) IRpcError
	MakeRequestWrapper() IRequestWrapper
	MakeRequest() IRequest
}


type ITransportOptions interface {
	
}

type ITransport interface {
	Serve(ITransportOptions)
	Connect(ITransportOptions)
	Send(id string, rw IRequestWrapper) error
	Close()
}


type IConnection interface {
	Close()
	Call(method string, params interface{}) (interface{}, error)
	Notify(method string, params interface{})
	Response(id int, result interface{})
}


type IProtocol interface {
	Encode(IRequestWrapper) ([]byte, error)
	Decode([]byte) (*RequestWrapper, error)
}


type IRouter interface {
	Route(connection IConnection, request IRequestWrapper) IRequestWrapper
}


type IRpcError interface {
	error
	GetCode() int
	GetMessage() string
	GetData() error
}


type IRequestWrapper interface {
	AddRequest(IRequest)
	SetBatchRequest(bool)
	IsEmpty() bool
	IsBatchRequest() bool
	GetRequest() IRequest
	GetBatchRequests() []IRequest
}


type IRequest interface {
	IsError() bool
	IsRequest() bool
	IsResponse() bool
	Populate(map[string]interface{})
}

