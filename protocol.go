package gorpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	log "github.com/Sirupsen/logrus"
)

type Protocol struct {
	FactoryGetter
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
				log.Errorf("Protocol: Batch element is not correct type: (%T)%v", vr[i],vr[i])
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

	result := this.invokeBatch(connection, request.GetBatchRequests())

	if result.IsEmpty() == false {
		result.SetBatchRequest(request.IsBatchRequest())
		return result
	}

	return nil
}

func (this *Protocol) SetValidate(f func(r IRequest)) {
	this.validator = f
}


func (this *Protocol) invokeBatch(connection IConnection, rm []IRequest) IRequestWrapper {
	result := this.Factory().MakeRequestWrapper()

	for i := 0; i < len(rm); i++ {
		if rm[i].IsRequest() {
			r, err := this.invoke(connection, rm[i])
			id := rm[i].Id()
			if id != nil {
				result.AddRequest(this.Factory().MakeResponse(id, r, err))
			}
			log.Debugf("Protocol.invokeBatch returning: %v", result)
		} else if rm[i].IsResponse() {
			log.Debugf("Protocol.Response: (%T)%v", rm[i], rm[i])
			connection.Response(rm[i].Id(), rm[i].Result())
		} else { // error...
			log.Debugf("Protocol.Error: (%T)%v", rm[i], rm[i])
			connection.Response(rm[i].Id(), rm[i].Error())
		}
	}
	return result
}

func (this *Protocol) invoke(connection IConnection, request IRequest) (response interface{}, err IRpcError) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("Protocol.invoke Recovered from panic: (%T)%v", r, r)
			switch x := r.(type) {
			case RpcError:
				err = &x
			case *RpcError:
				err = x
			case string:
				err = this.Factory().MakeRpcError(ErrInternalError, errors.New(x))
			case error:
				err = this.Factory().MakeRpcError(ErrInternalError, x)
			default:
				err = this.Factory().MakeRpcError(ErrInternalError, fmt.Errorf("Unknown: (%T)%v", x,x))
			}
		}
	}()

	if this.validator != nil {
		this.validator(request)
	}

	response = connection.Invoke(request)
	
	log.Debugf("Protocol.invoke Returning: (%T)%v, %v", response, response, err)
	return
}
