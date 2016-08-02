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
	MakeTransportOptions() ITransportOptions
	MakeTransport() ITransport
	MakeAddress(src, dest string, options interface{}) IConnectionAddress
	MakeConnection(transport ITransport, addr IConnectionAddress) IConnection
	MakeController() IController
	MakeProtocol() IProtocol
	MakeRouter() IRouter
	MakeRpcError(code int, previous error) IRpcError
	MakeRequestWrapper() IRequestWrapper
	MakeRequest(id, method, params interface{}) IRequest
	MakeResponse(id, result, error interface{}) IRequest
}


type ITransportOptions interface {
	
}

type ITransport interface {
	Serve(ITransportOptions)
	Connect(ITransportOptions)
	Send(id string, rw IRequestWrapper) error
	Close()
	Init(onReceive, onWrite func(id, message string)) // Because Go does not support method overriding
}

type IDefaultController interface {
	RPC_Echo(value string) string
}

type IConnectionAddress interface {
	SetAddress(src, dest string, options interface{})
	Source() string
	Destination() string
	Options() interface{}
}

type IConnection interface {
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
	Encode(IRequestWrapper) ([]byte, error)
	Decode([]byte) (IRequestWrapper, error)
	Parse(connection IConnection, request IRequestWrapper) IRequestWrapper
}


type IRouter interface {
	Route(connection IConnection, request IRequestWrapper) IRequestWrapper
	Init(validator func(IRequest))
}


type IRpcError interface {
	error
	GetCode() int
	GetMessage() string
	GetData() error
}


type IRequestWrapper interface {
	AddRequest(IRequest) IRequestWrapper
	Clear() IRequestWrapper
	SetBatchRequest(bool) IRequestWrapper
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
	SetRequest(id, method, params interface{})
	SetResponse(id, result, error interface{})

	Id() interface{}
	Method() string
	Params() interface{}
	Result() interface{}
	Error() interface{}
}

type IJsonRPC2Request interface {
	IRequest
	JsonRPC() string
}


