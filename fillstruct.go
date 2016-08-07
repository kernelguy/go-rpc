package gorpc

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"reflect"
)

func InitializeStruct(s interface{}) {
	log.Debugf("Init Struct: (%T)%v", s,s)

	v := reflect.ValueOf(s).Elem()
	t := v.Type()

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
			InitializeStruct(f.Addr().Interface())
		case reflect.Ptr:
			fv := reflect.New(ft.Type.Elem())
			InitializeStruct(fv.Interface())
			f.Set(fv)
		default:
		}
	}
}

func mapCopy(dst, src interface{}) {
	//log.Debugf("mapCopy: (%T)%v => (%T)%v", src,src, dst,dst)
    dv, sv := reflect.ValueOf(dst), reflect.ValueOf(src)

    for _, k := range sv.MapKeys() {
        dv.SetMapIndex(k, sv.MapIndex(k))
    }
}

func setField(obj interface{}, n interface{}, value interface{}) error {
	//log.Debugf("obj: (%T)%v", obj, obj)
	structValue := reflect.ValueOf(obj).Elem()
	var structFieldValue reflect.Value
	if name, ok := n.(string); ok {
		structFieldValue = structValue.FieldByName(name)
	} else {
		structFieldValue = structValue.Field(n.(int))
	}

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %v in obj", n)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %v field value", n)
	}

	structFieldType := structFieldValue.Type()
	//log.Debugf("setField: (%T)%v to (%T)%v", structFieldValue.Interface(), structFieldValue, value, value)

	if structFieldType.Kind() == reflect.Map {
		mapCopy(structFieldValue.Interface(), value)
		return nil
	}

	val := reflect.ValueOf(value)
	if value != nil {
		if val.Type().Kind() == reflect.Map {
			
			if structFieldType.Kind() == reflect.Ptr {
				structFieldValue = structFieldValue.Elem()
			}
			err := FillStructFromMap(structFieldValue.Addr().Interface(), val.Interface().(map[string]interface{}))
			return err
		}
		//log.Debugf("Setting Value (%T)%v to (%T)%v", structFieldValue.Interface(), structFieldValue, val.Interface(), val)
		structFieldValue.Set(val.Convert(structFieldType))
	} else {
		structFieldValue.Set(reflect.Zero(structFieldType))
	}

	return nil
}

func FillStructFromMap(s interface{}, m map[string]interface{}) error {
	log.Debugf("Traverse into (%T)%v filling with (%T)%v", s,s, m,m)
	for k, v := range m {
		err := setField(s, k, v)
		if err != nil {
			log.Debugf("FillStructFromMap.Error((%T)%v)", err,err)
			return err
		}
	}
	//log.Debugf("Traverse result (%T)%v", s,s)
	return nil
}

func FillStructFromArray(s interface{}, m []interface{}) error {
	log.Debugf("Traverse into (%T)%v filling with (%T)%v", s,s, m, m)
	for i:=0; i < len(m); i++ {
		err := setField(s, i, m[i])
		if err != nil {
			log.Debugf("FillStructFromArray.Error((%T)%v)", err, err)
			return err
		}
	}
	//log.Debugf("Traverse result (%T)%v", s,s)
	return nil
}


