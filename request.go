package gorpc

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)


type Request struct {
	data map[string]interface{}
}


func (this *Request) IsError() bool {
	return (this.Error() != nil)
}

func (this *Request) IsRequest() bool {
	return (this.Error() == nil) && (this.Result() == nil)
}

func (this *Request) IsResponse() bool {
	return (this.Result() != nil)
}

func (this *Request) Populate(vr map[string]interface{}) {
	this.data = vr
}

func (this *Request) CreateRequest(id, method, params interface{}) {
	this.data = make(map[string]interface{}, 4)
	if id != nil {
		this.data["id"] = id
	}
	this.data["method"] = method
	if params != nil {
		this.data["params"] = params
	}
	this.data["jsonrpc"] = "2.0"
	log.Debugf("Request created: %v", this.data)
}

func (this *Request) CreateResponse(id, result, error interface{}) {
	this.data = make(map[string]interface{}, 3)
	this.data["id"] = id
	if error != nil {
		this.data["error"] = error
	} else {
		this.data["result"] = result
	}
	this.data["jsonrpc"] = "2.0"
}

func (this *Request) Id() interface{} {
	if this.data == nil {
		return nil
	}
	return this.data["id"]
}

func (this *Request) Method() string {
	if this.data == nil {
		return ""
	}
	return this.data["method"].(string)
}

func (this *Request) Params() interface{} {
	if this.data == nil {
		return nil
	}
	return this.data["params"]
}

func (this *Request) Result() interface{} {
	if this.data == nil {
		return nil
	}
	return this.data["result"]
}

func (this *Request) Error() interface{} {
	if this.data == nil {
		return nil
	}
	return this.data["error"]
}

func (this *Request) JsonRPC() string {
	if this.data == nil {
		return ""
	}
	return this.data["jsonrpc"].(string)
}

func (this *Request) MarshalJSON() (result []byte, err error) {
	result, err = json.Marshal(this.data)
	log.Debugf("Request Encoded: %s, %v", string(result), err)
	return result, err
}
