package storage

import (
	"fmt"
	"reflect"
	"strings"
)

func GenericBuildClient[T any](objDesc any) T {

	// after recovery, zero value of enclosing function
	// will be returned
	defer RecoverFromPanic()

	var obj T
	if map_rep, ok := objDesc.(map[string]any); ok {

		jsonKeyToStructField := GetJsonKeyToStructField(obj)
		for key, val := range map_rep {
			SetProperty(&obj, jsonKeyToStructField[key], val)
		}

	}
	return obj
}

func GetJsonKeyToStructField(obj any) map[string]string {
	var res = map[string]string{}
	typ := reflect.TypeOf(obj)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonKey := strings.Split(field.Tag.Get("json"), ",")[0]
		res[jsonKey] = field.Name
	}

	return res
}

func SetProperty(obj any, propName string, propValue any) {
	field := reflect.ValueOf(obj).Elem().FieldByName(propName)

	defer RecoverFromPanic()

	numVal, _ := getFloat64Equivalent(propValue)
	switch field.Kind() {
	case reflect.Int:
		field.Set(reflect.ValueOf(int(numVal)))
	case reflect.Int8:
		field.Set(reflect.ValueOf(int8(numVal)))
	case reflect.Int16:
		field.Set(reflect.ValueOf(int16(numVal)))
	case reflect.Int32:
		field.Set(reflect.ValueOf(int32(numVal)))
	case reflect.Int64:
		field.Set(reflect.ValueOf(int64(numVal)))
	case reflect.Float32:
		field.Set(reflect.ValueOf(float32(numVal)))
	case reflect.Float64:
		field.Set(reflect.ValueOf(float64(numVal)))
	default:
		field.Set(reflect.ValueOf(propValue))
	}
}

func RecoverFromPanic() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
