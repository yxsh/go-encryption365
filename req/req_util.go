package req

import "reflect"

func StructToValue(r interface{}) Value{
	v := reflect.ValueOf(r)
	if v.Type().Kind() == reflect.Ptr {
		ev := v.Elem()
		if ev.Kind() != reflect.Struct{
			return nil
		}
		val := make(Value)
		for i := 0; i < ev.NumField(); i++ {
			field := ev.Field(i)
			val[ev.Type().Field(i).Tag.Get("json")] = field.Interface().(string)
		}
		return val
	}
	return nil
}
