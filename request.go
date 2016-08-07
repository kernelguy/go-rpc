package gorpc

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"sort"
)


type Request struct {
	data map[string]interface{}
}

type __notSet__ struct {
}

func (this *Request) IsError() bool {
	if this.data == nil {
		return false
	}
	if _, ok := this.data["error"]; ok {
		return true
	}
	return false
}

func (this *Request) IsRequest() bool {
	if this.data == nil {
		return false
	}
	if _, ok := this.data["method"]; ok {
		return true
	}
	return false
}

func (this *Request) IsResponse() bool {
	if this.data == nil {
		return false
	}
	if _, ok := this.data["result"]; ok {
		return true
	}
	return false
}

func (this *Request) Populate(vr map[string]interface{}) {
	this.data = vr
	err, _ := vr["error"].(map[string]interface{})
	if err != nil {
		e := &RpcError{Code: int(err["code"].(float64)), Message: err["message"].(string)}
		if err["data"] != nil {
			switch x := err["data"].(type) {
			case string:
				e.Data = errors.New(x)
			case error:
				e.Data = x
			}
		}
		this.data["error"] = e
	}
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
	log.Debugf("Request.SetRequest result: %v", this.data)
}

func (this *Request) SetResponse(id, result, error interface{}) {
	this.data = make(map[string]interface{}, 3)
	this.data["id"] = id
	if error != nil {
		this.data["error"] = error
	} else {
		this.data["result"] = result
	}
	log.Debugf("Request.SetResponse result: %v", this.data)
}

func (this *Request) Set(name string, value interface{}) {
	if this.data != nil {
		this.data[name] = value
	}
}

func (this *Request) Get(name string) interface{} {
	if this.data == nil {
		return nil
	}
	return this.data[name]
}

func (this *Request) Id() interface{} {
	return this.Get("id")
}

func (this *Request) Method() string {
	if v, ok := this.Get("method").(string); ok {
		return v
	}
	return ""
}

func (this *Request) Params() interface{} {
	return this.Get("params")
}

func (this *Request) Result() interface{} {
	return this.Get("result")
}

func (this *Request) Error() interface{} {
	return this.Get("error")
}

func (this *Request) String() string {
	s := "Request{data:"
	sep := ""
	if this.data != nil {
		idx := make([]string, len(this.data))
		i := 0;
		for k, _ := range this.data {
			idx[i] = k
			i++
		}
		sort.Strings(idx)

		s = s + "{"
		for _, k := range idx {
			v := this.data[k]
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
