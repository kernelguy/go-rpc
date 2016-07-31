package gorpc

import (
	"encoding/json"
	"reflect"
	"fmt"
	log "github.com/Sirupsen/logrus"
)

type Protocol struct {
	router IRouter
}


func (this *Protocol) Encode(request IRequestWrapper) ([]byte, error) {
	b, err := json.Marshal(request)
	return b, err
}

func (this *Protocol) Decode(data []byte) (IRequestWrapper, error) {
	log.Debugf("Protocol.Decode(%s)", string(data))
	var v interface{}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	result := RequestWrapper{}
	log.Debugf("Protocol.Decode Result: %v", v)
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		result.SetBatchRequest(true)
		vr, _ := v.([]map[string]interface{})
		for i:=0 ; i < len(vr); i++ {
			req := GetFactory().MakeRequest()
			req.Populate(vr[i])
			result.AddRequest(req)
		} 
	} else {
		req := GetFactory().MakeRequest()
		vr, _ := v.(map[string]interface{})
		req.Populate(vr)
		result.AddRequest(req)
	}
	return &result, nil
}

func (this *Protocol) Parse(connection IConnection, request IRequestWrapper) IRequestWrapper {
	
	if this.router == nil {
		this.router = GetFactory().MakeRouter()
		this.router.Init(this.validate)
	}

	return this.router.Route(connection, request)	
}


func (this *Protocol) validate(r IRequest) {
	json, ok := r.(IJsonRPC2Request)
	if !ok {
		panic(GetFactory().MakeRpcError(ErrInvalidRequest, fmt.Errorf("Jsonrpc property does not exist")))
	}
	if json.JsonRPC() != "2.0" {
		panic(GetFactory().MakeRpcError(ErrInvalidRequest, fmt.Errorf("Jsonrpc version mismatch")))
	}
}
