package gorpc

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
	"strings"
)

type Router struct {
	FactoryGetter
	routeCache map[string]reflect.Value
}

func (this *Router) GetRoute(obj interface{}, method string) reflect.Value {
	if this.routeCache == nil {
		this.routeCache = make(map[string]reflect.Value, 10)
	}
	
	m, ok := this.routeCache[method]
	if !ok {
		m = this.findRoute(obj, method)
		this.routeCache[method] = m
	}
	return m
}		

func (this *Router) findRoute(obj interface{}, method string) reflect.Value {

	path := strings.Split(method, ".")
	v := reflect.ValueOf(obj)
	for len(path) > 1 {
		f := v.Elem().FieldByName(path[0])
		if !f.CanAddr() {
			panic(this.Factory().MakeRpcError(ErrMethodNotFound, nil))
		}
		v = f.Addr()
		path = path[1:]
	}
	m := v.MethodByName("RPC_" + path[0])
	if m.IsValid() == false {
		panic(this.Factory().MakeRpcError(ErrMethodNotFound, nil))
	}

	return m
}

func (this *Router) CheckParams(m reflect.Value, params interface{}) (result []reflect.Value) {

	n := m.Type().NumIn()
	if n == 0 {
		return // No parameters
	} else if n > 1 {
		panic(fmt.Errorf("Design Error: RPC methods should have zero or one parameter: Rpc_MyMethod(struct{a int, b string})"))
	}
	rt := m.Type().In(0)
	n = rt.NumField()

	if params != nil {
		rv := reflect.New(rt)
		InitializeStruct(rv.Interface())
		var err error = nil
		 
		v := reflect.TypeOf(params)
		switch v.Kind() {
		case reflect.Map:
			p := params.(map[string]interface{})
			if len(p) != n {
				err = fmt.Errorf("Wrong parameter count.")
			} else {
				err = FillStructFromMap(rv.Interface(), p)
			}

		case reflect.Slice:
			p := params.([]interface{})
			if len(p) != n {
				err = fmt.Errorf("Wrong parameter count.")
			} else {
				err = FillStructFromArray(rv.Interface(), p)
			}

		default:
			err = fmt.Errorf("Unknown parameter type: (%T)%v", v, v)
		}
		if err != nil {
			panic(this.Factory().MakeRpcError(ErrInvalidParams, err))
		}
		result = append(result, rv.Elem())
	} else if n > 0 {
		panic(this.Factory().MakeRpcError(ErrInvalidParams, nil))
	}
	log.Debugf("Router.checkParams returning %v", result)
	return
}
