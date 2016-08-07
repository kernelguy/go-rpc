// gorpc is a framework for (bi)directional remote procedure calls over any transport.
// It can be adapted in multiple ways to use other transports and protocols.
package gorpc

import (
)

// Interfaces used in the gorpc package.
// OOP in Go sucks: Interfaces can be thought of as pure abstract classes,
// struct's can have methods, but there is no constructors or method overriding.
// All in all, it's a pain to make generic packages that can be used as base for
// specific implementations.
//
// The idea with this package was to make a small framework to support various
// RPC protocols over all sorts of transports.

type IFactory interface {
	MakeAddress(src, dest string, options interface{}) IConnectionAddress
	MakeConnection(transport ITransport, addr IConnectionAddress) IConnection
	MakeController() IController
	MakeProtocol() IProtocol
	MakeRequest(id, method, params interface{}) IRequest
	MakeRequestWrapper() IRequestWrapper
	MakeResponse(id, result, error interface{}) IRequest
	MakeRouter(validator func(request IRequest)) IRouter
	MakeRpcError(code int, previous error) IRpcError
	MakeTransport(ITransportOptions) ITransport
	MakeTransportOptions() ITransportOptions
}

type IFactoryGetter interface {
	Factory() IFactory
	SetFactory(IFactory)
}

type ITransportOptions interface {
	
}

type ITransport interface {
	IFactoryGetter
	Close()
	Init(onReceive, onWrite func(id, message string))
	Receive(id, message string)
	Send(id string, rw IRequestWrapper) error
	Start()
}

type IDefaultController interface {
	RPC_Echo(value string) string
}

type IConnectionAddress interface {
	Destination() string
	Options() interface{}
	SetAddress(src, dest string, options interface{})
	Source() string
}

type IConnection interface {
	IFactoryGetter
	IConnectionAddress
	Close()
	Call(method string, params interface{}) (interface{}, error)
	Notify(method string, params interface{})
	Response(id interface{}, result interface{})
	RootController() IController
	SetRootController(obj IController)
}

type IController interface {
	Connection() IConnection
	SetConnection(IConnection)
}

type IProtocol interface {
	IFactoryGetter
	Decode([]byte) (IRequestWrapper, error)
	Encode(IRequestWrapper) ([]byte, error)
	Parse(connection IConnection, request IRequestWrapper) IRequestWrapper
	SetValidate(func(IRequest))
}


type IRouter interface {
	IFactoryGetter
	SetValidator(validator func(IRequest))
	Route(connection IConnection, request IRequestWrapper) IRequestWrapper
}


type IRpcError interface {
	error
	GetCode() int
	GetData() error
	GetMessage() string
}


type IRequestWrapper interface {
	AddRequest(IRequest) IRequestWrapper
	Clear() IRequestWrapper
	IsBatchRequest() bool
	IsEmpty() bool
	GetBatchRequests() []IRequest
	GetRequest() IRequest
	SetBatchRequest(bool) IRequestWrapper
}


type IRequest interface {
	Error() interface{}
	Id() interface{}
	IsError() bool
	IsRequest() bool
	IsResponse() bool
	Method() string
	Params() interface{}
	Populate(map[string]interface{})
	Result() interface{}
	SetRequest(id, method, params interface{})
	SetResponse(id, result, error interface{})
	String() string
}

type IJsonRPC2Request interface {
	IRequest
	JsonRPC() string
}


