package config

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Addr string            `yaml:"addr"`
	Cron map[string]string `yaml:"cron"`
	Log  LogConfig         `yaml:"log"`
}

type LogConfig struct {
	SingleCapacity int    `yaml:"single_capacity" default:"0"`
	RuntimePath    string `yaml:"runtime_path" default:"./"`
	Test           bool   `yaml:"test"`
}

var GlobalConfig = new(Config)

func init() {
	yamlFile, err := ioutil.ReadFile("/home/spider1998/goweb/pkg/config/conf.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, GlobalConfig)
	if err != nil {
		panic(err)
	}
	structByReflect(GlobalConfig)
	return
}

func structByReflect(inStructPtr interface{}) {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		rType = rType.Elem()
		rVal = rVal.Elem()
	}
	// 遍历结构体
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		switch f.Kind() {
		case reflect.Struct:
			structByReflect(f.Addr().Interface())
		default:
			if !isBlank(f) {
				continue
			}
			key := t.Tag.Get("default")
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
