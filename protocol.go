package gorpc

import (
	"encoding/json"
	"reflect"
	log "github.com/Sirupsen/logrus"
)

type Protocol struct {
	FactoryGetter
	router IRouter
	validator func(r IRequest)
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
	f := this.Factory()
	result := f.MakeRequestWrapper()
	log.Debugf("Protocol.Decode Result: %v", v)
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		result.SetBatchRequest(true)
		vr, _ := v.([]interface{})
		for i:=0 ; i < len(vr); i++ {
			m, ok := vr[i].(map[string]interface{})
			if !ok {
				log.Errorf("Batch element is not correct type: (%T)%v", vr[i],vr[i])
				return nil, this.Factory().MakeRpcError(ErrParseError, nil)
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
		this.router = this.Factory().MakeRouter(this.validator)
	}

	return this.router.Route(connection, request)	
}

func (this *Protocol) SetValidate(f func(r IRequest)) {
	this.validator = f
}
