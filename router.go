package gorpc

import (
	"strings"
	"reflect"
	"fmt"
	"errors"
)


type Router struct {
	
}


func (p *Router) Route(connection IConnection, request IRequestWrapper) IRequestWrapper {

	result := p.invokeBatch(connection, request.GetBatchRequests())
	
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
			if rm[i].(*Request).id != nil {
				req := Request{id: rm[i].(*Request).id}
				if err != nil {
					req.error = err
				} else {
					req.result = r
				}
				result.AddRequest(&req)
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
			switch x := r.(type) {
		        case RpcError:
			        err = &x    
		        case string:
		            err = NewRpcError(cInternalError, errors.New(x))
		        case error:
		            err = NewRpcError(cInternalError, x)
		        default:
		            err = NewRpcError(cInternalError, errors.New("Unknown panic"))
			}
		}
	}()

	_, method := this.getRoute(connection, request.method)
	
	p := this.checkParams(method, request.params)

	r := reflect.ValueOf(method).Call(p)
	
	if len(r) > 0 {
		response = r[0]
	}
	return
}

func (this *Router) getRoute(connection IConnection, method string) (interface{}, reflect.Value) {

	var obj interface{} = connection
	path := strings.Split(method, ".")
	for len(path) > 1 {
		obj = reflect.ValueOf(&obj).FieldByName(path[0])
		if obj == nil {
			panic(NewRpcError(cMethodNotFound, nil))
		}
		path = path[1:]
	}
	m := reflect.ValueOf(&obj).MethodByName("Rpc_" + path[0])
	if m.IsValid() == false {
		panic(NewRpcError(cMethodNotFound, nil))
	}

	return obj, m;
}

func (this *Router) checkParams(m reflect.Value, params interface{}) (result []reflect.Value) {

	requiredParams := make(map[int]string, 10)

	n := m.NumField()
	if n == 0 {
		return
	} else if n > 1 {
		panic(fmt.Errorf("RPC methods should have zero or one parameter: Rpc_MyMethod(struct{a int, b string})"))
	}
	
	n = m.Field(0).Type().NumField()
	for i:=0 ; i < n ; i++ {
		requiredParams[i] = m.Field(0).Type().Field(i).Name 
	}

	if params != nil {
		v := reflect.TypeOf(params)
		switch v.Kind() {
			case reflect.Struct:
				if v.NumField() != n {
					panic(NewRpcError(cInvalidParams, nil))
				}
				for i:=0; i < n; i++ {
					result = append(result, reflect.ValueOf(params).FieldByName(requiredParams[i]))
				}
	
			case reflect.Slice:
				if len(params.([]reflect.Value)) != n {
					panic(NewRpcError(cInvalidParams, nil))
				}
				result = params.([]reflect.Value)
	
			default:
				panic(NewRpcError(cInvalidParams, nil))
		}
	}
	return
}


func (p *Router) onResponse(connection IConnection, request IRequest) {

}

func (p *Router) onError(connection IConnection, request IRequest) {

}
