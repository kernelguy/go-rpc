package gorpc

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
	"strings"
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

	for i := 0; i < len(rm); i++ {
		if rm[i].IsRequest() {
			r, err := p.invoke(connection, rm[i].(*Request))
			id := rm[i].Id()
			if id != nil {
				result.AddRequest(GetFactory().MakeResponse(id, r, err))
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

func (this *Router) invoke(connection IConnection, request *Request) (response interface{}, err IRpcError) {
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
				err = NewRpcError(ErrInternalError, errors.New("Unknown panic"))
			}
		}
	}()

	this.validate(request)
	log.Debug("Router.invoke step 1")

	_, method := this.getRoute(connection.RootController(), request.Method())
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

func (this *Router) getRoute(obj interface{}, method string) (interface{}, reflect.Value) {

	path := strings.Split(method, ".")
	v := reflect.ValueOf(obj)
	for len(path) > 1 {
		v = v.Elem().FieldByName(path[0]).Addr()
		if !v.IsValid() {
			panic(NewRpcError(ErrMethodNotFound, nil))
		}
		path = path[1:]
	}
	m := v.MethodByName("RPC_" + path[0])
	if m.IsValid() == false {
		panic(NewRpcError(ErrMethodNotFound, nil))
	}

	return obj, m
}

func (this *Router) checkParams(m reflect.Value, params interface{}) (result []reflect.Value) {

	requiredParams := make(map[int]reflect.StructField, 10)

	n := m.Type().NumIn()
	if n == 0 {
		return
	} else if n > 1 {
		panic(fmt.Errorf("RPC methods should have zero or one parameter: Rpc_MyMethod(struct{a int, b string})"))
	}

	n = m.Type().In(0).NumField()
	for i := 0; i < n; i++ {
		requiredParams[i] = m.Type().In(0).Field(i)
	}

	if params != nil {
		rt := m.Type().In(0)
		rv := reflect.New(rt)
		InitializeStruct(rt, rv.Elem())

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
			FillStructFromArray(rv.Interface(), p, requiredParams)
			result = append(result, rv.Elem())

		default:
			log.Debugf("Router.checkParams Unknown parameter kind: %(T)%v", v, v)
			panic(NewRpcError(ErrInvalidParams, nil))
		}
	}
	log.Debugf("Router.checkParams returning %v", result)
	return
}
