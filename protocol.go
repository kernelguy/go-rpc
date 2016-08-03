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
	f := GetFactory()
	result := f.MakeRequestWrapper()
	log.Debugf("Protocol.Decode Result: %v", v)
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		result.SetBatchRequest(true)
		vr, _ := v.([]interface{})
		for i:=0 ; i < len(vr); i++ {
			m, ok := vr[i].(map[string]interface{})
			if !ok {
				log.Errorf("Batch element is not correct type: (%T)%v", vr[i],vr[i])
				return nil, GetFactory().MakeRpcError(ErrParseError, nil)
			}
			result.AddRequest(f.MakeRequest(m, nil, nil))
		} 
	} else {
		result.AddRequest(f.MakeRequest(v, nil, nil))
	}
	return result, nil
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
