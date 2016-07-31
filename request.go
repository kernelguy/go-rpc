package gorpc

import (
	"encoding/json"
	"fmt"
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

func (this *Request) SetRequest(id, method, params interface{}) {
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

func (this *Request) SetResponse(id, result, error interface{}) {
	this.data = make(map[string]interface{}, 3)
	this.data["id"] = id
	if error != nil {
		this.data["error"] = error
	} else {
		this.data["result"] = result
	}
	this.data["jsonrpc"] = "2.0"
	log.Debugf("Request.CreateResponse result: %v", this.data)
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

func (this *Request) String() string {
	s := "Request{data:"
	sep := ""
	if this.data != nil {
		s = s + "{"
		for k, v := range this.data {
			switch v.(type) {
				case nil:
					s = fmt.Sprintf("%s%s%s:nil", s, sep, k)
				case string:
					s = fmt.Sprintf("%s%s%s:\"%s\"", s, sep, k, v)
				case int, int64:
					s = fmt.Sprintf("%s%s%s:%d", s, sep, k, v)
				case float32, float64:
					s = fmt.Sprintf("%s%s%s:%f", s, sep, k, v)
				default:
					s = fmt.Sprintf("%s%s%s:(%T)%v", s, sep, k, v, v)
			}
			sep = ", "
		}
		s = s + "}"
	}
	s = s + "}"
	return s
}

func (this *Request) MarshalJSON() (result []byte, err error) {
	result, err = json.Marshal(this.data)
	log.Debugf("Request Encoded: %s, %v", string(result), err)
	return result, err
}
