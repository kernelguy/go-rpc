package jsonrpc2

import (
	"fmt"
	"github.com/kernelguy/gorpc"
)

type Protocol struct {
	gorpc.Protocol
}

func (this *Protocol) Parse(connection gorpc.IConnection, request gorpc.IRequestWrapper) gorpc.IRequestWrapper {
	this.SetValidate(this.validate)
	return this.Protocol.Parse(connection, request)	
}

func (this *Protocol) validate(r gorpc.IRequest) {
	json, ok := r.(gorpc.IJsonRPC2Request)
	if !ok {
		panic(gorpc.GetFactory().MakeRpcError(gorpc.ErrInvalidRequest, fmt.Errorf("Jsonrpc property does not exist")))
	}
	if json.JsonRPC() != "2.0" {
		panic(gorpc.GetFactory().MakeRpcError(gorpc.ErrInvalidRequest, fmt.Errorf("Jsonrpc version mismatch")))
	}
}


