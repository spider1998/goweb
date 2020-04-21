package util

import (
	"fmt"
	"reflect"
	"strconv"
)

func ParseTagReflect(inStructPtr interface{}, tag string) {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		rType = rType.Elem()
		rVal = rVal.Elem()
	}
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			ParseTagReflect(f.Addr().Interface(), tag)
		default:
			if !isBlank(f) {
				continue
			}
			key := t.Tag.Get(tag)
			fmt.Println(f.Type())
			structType := f.Type()
			switch structType.Kind() {
			case reflect.Int:
				v, _ := strconv.Atoi(key)
				f.SetInt(int64(v))
			case reflect.Bool:
				if key == "true" {
					f.SetBool(true)
				} else {
					f.SetBool(false)
				}
			case reflect.String:
				f.Set(reflect.ValueOf(key).Convert(structType))
			}

		}
	}
}

func isBlank(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String:
		return value.Len() == 0
	case reflect.Bool:
		return !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return value.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return value.IsNil()
	}
	return reflect.DeepEqual(value.Interface(), reflect.Zero(value.Type()).Interface())
}
