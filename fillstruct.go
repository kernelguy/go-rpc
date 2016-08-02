package gorpc

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
)

func InitializeStruct(t reflect.Type, v reflect.Value) {
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
			InitializeStruct(ft.Type, f)
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			InitializeStruct(ft.Type.Elem(), fv.Elem())
			f.Set(fv)
		default:
		}
	}
}

func setField(obj interface{}, name string, value interface{}) error {
	log.Debugf("obj: (%T)%v", obj, obj)
	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if val.Type().Kind() == reflect.Map {
		err := FillStructFromMap(structFieldValue.Addr().Interface(), val.Interface().(map[string]interface{}))
		if err != nil {
			return err
		}
		return nil
	}

	log.Debugf("Setting Value (%T)%v to (%T)%v", structFieldValue.Interface(), structFieldValue, val.Interface(), val)
	structFieldValue.Set(val.Convert(structFieldType))
	return nil
}

func FillStructFromMap(s interface{}, m map[string]interface{}) error {
	log.Debugf("Traverse into (%T)%v filling with (%T)%v", s,s, m, m)
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			log.Debugf("FillStruct.Error((%T)%v)", err, err)
			return err
		}
	}
	log.Debugf("Traverse result (%T)%v", s,s)
	return nil
}

func FillStructFromArray(s interface{}, m []interface{}, t map[int]reflect.StructField) error {
	log.Debugf("Traverse into (%T)%v filling with (%T)%v", s,s, m, m)
	for i:=0; i < len(m); i++ {
		err := setField(s, t[i].Name, m[i])
		if err != nil {
			log.Debugf("FillStruct.Error((%T)%v)", err, err)
			return err
		}
	}
	log.Debugf("Traverse result (%T)%v", s,s)
	return nil
}


