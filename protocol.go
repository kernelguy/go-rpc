package gorpc

import (
	"encoding/json"
	"reflect"
	"fmt"
)

type Protocol struct {
}


func (this *Protocol) Encode(request IRequestWrapper) ([]byte, error) {
	b, err := json.Marshal(request)
	return b, err
}

func (this *Protocol) Decode(data []byte) (*RequestWrapper, error) {
	var v interface{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return nil, err
	}
	result := RequestWrapper{}
	fmt.Printf("Protocol decode type: %s\n", reflect.TypeOf(v).String())
	if reflect.TypeOf(v).Kind() == reflect.Slice {
		result.SetBatchRequest(true)
		vr, _ := v.([]map[string]interface{})
		for i:=0 ; i < len(vr); i++ {
			if err := this.validate(vr[i]); err != nil {
				return nil, err
			}
			req := GetFactory().MakeRequest()
			req.Populate(vr[i])
			result.AddRequest(req)
		} 
	} else {
		if err := this.validate(v.(map[string]interface{})); err != nil {
			return nil, err
		}
		req := GetFactory().MakeRequest()
		req.Populate(v.(map[string]interface{}))
		result.AddRequest(req)
	}
	return &result, nil
}

func (this *Protocol) validate(r map[string]interface{}) error {
	json, ok := r["json"]
	if !ok {
		return GetFactory().MakeRpcError(cInvalidRequest, fmt.Errorf("Json property does not exist"))
	}
	if json != "2.0" {
		return GetFactory().MakeRpcError(cInvalidRequest, fmt.Errorf("Json version mismatch"))
	}
	return nil
}
