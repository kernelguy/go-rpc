package gorpc

import (
	"strings"
	"reflect"
	"fmt"
	"errors"
	log "github.com/Sirupsen/logrus"
)


type Router struct {
	validate func(request IRequest)
}


func (this *Router) Init(validator func(IRequest)) {
	this.validate = validator
}


func (this *Router) Route(connection IConnection, request IRequestWrapper) IRequestWrapper {
	result := this.invokeBatch(connection, request.GetBatchRequests())
	
	if result.IsEmpty() == false {
		result.SetBatchRequest(request.IsBatchRequest())
		return result
	}
	
	return nil
}


func (p *Router) invokeBatch(connection IConnection, rm []IRequest) IRequestWrapper {
	result := GetFactory().MakeRequestWrapper()
	 
	for i:=0; i < len(rm); i++ {
		if rm[i].IsRequest() {
			r, err := p.invoke(connection, rm[i].(*Request))
			id := rm[i].(*Request).Id()
			if id != nil {
				req := GetFactory().MakeRequest()
				req.CreateResponse(id, r, err)
				result.AddRequest(req)
			}
		} else if rm[i].IsResponse() {
			p.onResponse(connection, rm[i])
		} else { // error...
			p.onError(connection, rm[i])
		}
	}
	return result
}

func (this *Router) invoke(connection IConnection, request *Request) (response reflect.Value, err IRpcError) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("Router.invoke Recovered from panic: %v", r) 
			switch x := r.(type) {
		        case RpcError:
			        err = &x    
		        case string:
		            err = NewRpcError(ErrInternalError, errors.New(x))
		        case error:
		            err = NewRpcError(ErrInternalError, x)
		        default:
		            err = NewRpcError(ErrInternalError, errors.New("Unknown panic"))
			}
		}
	}()

	this.validate(request)
	log.Debug("Step 1")	

	_, method := this.getRoute(connection.RootController(), request.Method())
	log.Debug("Step 2")	
	
	p := this.checkParams(method, request.Params())
	log.Debug("Step 3")	

	r := method.Call(p)
	log.Debug("Step 4")	
	
	if len(r) > 0 {
		response = r[0]
	}
	return
}

func (this *Router) getRoute(obj interface{}, method string) (interface{}, reflect.Value) {

	path := strings.Split(method, ".")
	for len(path) > 1 {
		obj = reflect.ValueOf(obj).FieldByName(path[0])
		if obj == nil {
			panic(NewRpcError(ErrMethodNotFound, nil))
		}
		path = path[1:]
	}
	m := reflect.ValueOf(obj).MethodByName("RPC_" + path[0])
	if m.IsValid() == false {
		panic(NewRpcError(ErrMethodNotFound, nil))
	}

	return obj, m;
}

func (this *Router) checkParams(m reflect.Value, params interface{}) (result []reflect.Value) {

	requiredParams := make(map[int]string, 10)

	n := m.Type().NumIn()
	if n == 0 {
		return
	} else if n > 1 {
		panic(fmt.Errorf("RPC methods should have zero or one parameter: Rpc_MyMethod(struct{a int, b string})"))
	}
	
	n = m.Type().In(0).NumField()
	for i:=0 ; i < n ; i++ {
		requiredParams[i] = m.Type().In(0).Field(i).Name 
	}

	if params != nil {
		t := reflect.New(m.Type().In(0))
		log.Debug("Step 2.5.1 ", t.Type())
		v := reflect.TypeOf(params)
		switch v.Kind() {
			case reflect.Map:
				p := params.(map[string]interface{})
				if len(p) != n {
					panic(NewRpcError(ErrInvalidParams, nil))
				}
				for i:=0; i < n; i++ {
					f := reflect.Indirect(reflect.ValueOf(t)).FieldByName(requiredParams[i])
					log.Debugf("Step 2.5.2 %v, %v, %v", requiredParams[i], f, p[requiredParams[i]])
					reflect.ValueOf(t).FieldByName(requiredParams[i]).Set(reflect.ValueOf(p[requiredParams[i]]))
					//result = append(result, reflect.ValueOf(p[requiredParams[i]]))
				}
				result = append(result, reflect.ValueOf(t))

			case reflect.Struct:
				log.Debug("Step 2.5.1 ", v)
				if v.NumField() != n {
					panic(NewRpcError(ErrInvalidParams, nil))
				}
				log.Debug("Step 2.5.2")
				for i:=0; i < n; i++ {
					result = append(result, reflect.ValueOf(params).FieldByName(requiredParams[i]))
				}
	
			case reflect.Slice:
				log.Debug("Step 2.5.3")
				if len(params.([]reflect.Value)) != n {
					panic(NewRpcError(ErrInvalidParams, nil))
				}
				log.Debug("Step 2.5.4")
				result = params.([]reflect.Value)
	
			default:
				log.Debug("Step 2.5.5")
				panic(NewRpcError(ErrInvalidParams, nil))
		}
	}
	return
}


func (p *Router) onResponse(connection IConnection, request IRequest) {

}

func (p *Router) onError(connection IConnection, request IRequest) {

}
