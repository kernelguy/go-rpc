package gorpc

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
	"strings"
)

type Router struct {
	FactoryGetter
	validate func(request IRequest)
	cache map[string]reflect.Value
}

func (this *Router) SetValidator(validator func(IRequest)) {
	this.validate = validator
}

func (this *Router) Route(connection IConnection, request IRequestWrapper) IRequestWrapper {
	if this.cache == nil {
		this.cache = make(map[string]reflect.Value, 10)
	}
	
	result := this.invokeBatch(connection, request.GetBatchRequests())

	if result.IsEmpty() == false {
		result.SetBatchRequest(request.IsBatchRequest())
		return result
	}

	return nil
}

func (this *Router) invokeBatch(connection IConnection, rm []IRequest) IRequestWrapper {
	result := this.Factory().MakeRequestWrapper()

	for i := 0; i < len(rm); i++ {
		if rm[i].IsRequest() {
			r, err := this.invoke(connection, rm[i])
			id := rm[i].Id()
			if id != nil {
				result.AddRequest(this.Factory().MakeResponse(id, r, err))
			}
			log.Debugf("Router.invokeBatch returning: %v", result)
		} else if rm[i].IsResponse() {
			log.Debugf("Router.Response: (%T)%v", rm[i], rm[i])
			connection.Response(rm[i].Id(), rm[i].Result())
		} else { // error...
			log.Debugf("Router.Error: (%T)%v", rm[i], rm[i])
			connection.Response(rm[i].Id(), rm[i].Error())
		}
	}
	return result
}

func (this *Router) invoke(connection IConnection, request IRequest) (response interface{}, err IRpcError) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("Router.invoke Recovered from panic: (%T)%v", r, r)
			switch x := r.(type) {
			case RpcError:
				err = &x
			case *RpcError:
				err = x
			case string:
				err = NewRpcError(ErrInternalError, errors.New(x))
			case error:
				err = NewRpcError(ErrInternalError, x)
			default:
				err = NewRpcError(ErrInternalError, fmt.Errorf("Unknown: (%T)%v", x,x))
			}
		}
	}()

	if this.validate != nil {
		this.validate(request)
	}
	log.Debug("Router.invoke step 1")

	var method reflect.Value
	method, ok := this.cache[request.Method()]
	if !ok {
		method = this.getRoute(connection.RootController(), request.Method())
		this.cache[request.Method()] = method
	}		
	log.Debug("Router.invoke step 2")

	p := this.checkParams(method, request.Params())
	log.Debug("Router.invoke step 3")

	r := method.Call(p)

	if len(r) > 0 {
		response = r[0].Interface()
	}
	log.Debugf("Router.invoke Returning: (%T)%v, %v", response, response, err)
	return
}

func (this *Router) getRoute(obj interface{}, method string) reflect.Value {

	path := strings.Split(method, ".")
	v := reflect.ValueOf(obj)
	for len(path) > 1 {
		f := v.Elem().FieldByName(path[0])
		if !f.CanAddr() {
			panic(NewRpcError(ErrMethodNotFound, nil))
		}
		v = f.Addr()
		path = path[1:]
	}
	m := v.MethodByName("RPC_" + path[0])
	if m.IsValid() == false {
		panic(NewRpcError(ErrMethodNotFound, nil))
	}

	return m
}

func (this *Router) checkParams(m reflect.Value, params interface{}) (result []reflect.Value) {

	n := m.Type().NumIn()
	if n == 0 {
		return
	} else if n > 1 {
		panic(fmt.Errorf("RPC methods should have zero or one parameter: Rpc_MyMethod(struct{a int, b string})"))
	}
	n = m.Type().In(0).NumField()

	if params != nil {
		rt := m.Type().In(0)
		rv := reflect.New(rt)
		InitializeStruct(rv.Interface())

		v := reflect.TypeOf(params)
		switch v.Kind() {
		case reflect.Map:
			p := params.(map[string]interface{})
			if len(p) != n {
				panic(NewRpcError(ErrInvalidParams, nil))
			}
			FillStructFromMap(rv.Interface(), p)
			result = append(result, rv.Elem())

		case reflect.Slice:
			p := params.([]interface{})
			if len(p) != n {
				panic(NewRpcError(ErrInvalidParams, nil))
			}
			FillStructFromArray(rv.Interface(), p)
			result = append(result, rv.Elem())

		default:
			log.Debugf("Router.checkParams Unknown parameter kind: %(T)%v", v, v)
			panic(NewRpcError(ErrInvalidParams, nil))
		}
	} else if n > 0 {
		panic(NewRpcError(ErrInvalidParams, nil))
	}
	log.Debugf("Router.checkParams returning %v", result)
	return
}
