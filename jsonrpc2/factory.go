package jsonrpc2

import (
	"github.com/kernelguy/gorpc"
)


type Factory struct {
	gorpc.Factory
}

func (this *Factory) MakeProtocol() gorpc.IProtocol {
	p := &Protocol{}
	p.SetFactory(this)
	return p
}

func (this *Factory) MakeRequest(id, method, params interface{}) gorpc.IRequest {
	r := &Request{}
	if v, ok := id.(map[string]interface{}); ok {
		r.Populate(v)
	} else {
		r.SetRequest(id, method, params)
	}
	return r
}

func (this *Factory) MakeResponse(id, result, error interface{}) gorpc.IRequest {
	r := &Request{}
	r.SetResponse(id, result, error)
	return r
}
