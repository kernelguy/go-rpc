package jsonrpc2

import (
	"github.com/kernelguy/gorpc"
//	log "github.com/Sirupsen/logrus"
)

type Request struct {
	gorpc.Request
}


func (this *Request) SetRequest(id, method, params interface{}) {
	this.Request.SetRequest(id, method, params)
	this.Set("jsonrpc", "2.0")
}

func (this *Request) SetResponse(id, result, error interface{}) {
	this.Request.SetResponse(id, result, error)
	this.Set("jsonrpc", "2.0")
}

func (this *Request) JsonRPC() string {
	if v, ok := this.Get("jsonrpc").(string); ok {
		return v
	}
	return ""
}
