package gorpc

import (

)


type Request struct {
	id interface{}
	method string
	params interface{}
	result interface{}
	error interface{}
}


func (p *Request) IsError() bool {
	return (p.error != nil)
}

func (p *Request) IsRequest() bool {
	return (p.error == nil) && (p.result == nil)
}

func (p *Request) IsResponse() bool {
	return (p.result != nil)
}

func (p *Request) Populate(vr map[string]interface{}) {
	p.method, _ = vr["method"].(string)
	p.id, _ 	= vr["id"].(string)
	p.params, _ = vr["params"]
	p.error, _  = vr["error"]
	p.result, _ = vr["result"]
}