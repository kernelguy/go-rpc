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
			connection.Response(rm[i].Id(), rm[i].Result())
		} else { // error...
			connection.Response(rm[i].Id(), rm[i].Error())
		}
	}
	return result
}

func (this *Router) invoke(connection IConnection, request *Request) (response interface{}, err IRpcError) {
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

	return obj, m
}

func (this *Router) initializeStruct(t reflect.Type, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		switch ft.Type.Kind() {
		case reflect.Map:
			f.Set(reflect.MakeMap(ft.Type))
		case reflect.Slice:
			f.Set(reflect.MakeSlice(ft.Type, 0, 0))
		case reflect.Chan:
			f.Set(reflect.MakeChan(ft.Type, 0))
		case reflect.Struct:
			this.initializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			this.initializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
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
	for i := 0; i < n; i++ {
		requiredParams[i] = m.Type().In(0).Field(i).Name
	}

	if params != nil {
		//t := reflect.TypeOf(EchoParams{})
		rt := m.Type().In(0)
		rv := reflect.New(rt)
		this.initializeStruct(rt, rv.Elem())

		v := reflect.TypeOf(params)
		switch v.Kind() {
		case reflect.Map:
			p := params.(map[string]interface{})
			if len(p) != n {
				panic(NewRpcError(ErrInvalidParams, nil))
			}
			for i := 0; i < n; i++ {
				rv.Elem().FieldByName(requiredParams[i]).Set(reflect.ValueOf(p[requiredParams[i]]))
			}
			result = append(result, rv.Elem())

		case reflect.Slice:
			if len(params.([]interface{})) != n {
				panic(NewRpcError(ErrInvalidParams, nil))
			}
			p := params.([]interface{})
			for i := 0; i < n; i++ {
				rv.Elem().FieldByName(requiredParams[i]).Set(reflect.ValueOf(p[i]))
			}
			result = append(result, rv.Elem())

		default:
			log.Debugf("Router.checkParams Unknown parameter kind: %(T)%v", v, v)
			panic(NewRpcError(ErrInvalidParams, nil))
		}
	}
	log.Debugf("Router.checkParams returning %v", result)
	return
}
