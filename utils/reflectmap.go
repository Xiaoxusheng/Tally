package utils

import (
	"reflect"
)

/*
运用反射把结构体反射成map
*/

// ReflectMap 反射
func ReflectMap(s any) map[any]any {
	val := reflect.ValueOf(s).Elem()
	types := reflect.TypeOf(s).Elem()
	//不是结构体
	if val.Kind() != reflect.Struct {
		return nil
	}
	var key any
	var value any
	m := make(map[any]any)
	for i := 0; i < val.NumField(); i++ {
		//这个是小写
		if !types.Field(i).IsExported() {
			continue
		}
		//
		if types.Field(i).Tag.Get("json") == "key" {
			key = val.Field(i).Interface()
		} else if types.Field(i).Tag.Get("json") == "value" {
			value = val.Field(i).Interface()
		}
	}
	m[key] = value
	return m
}
