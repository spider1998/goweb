package components

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"goweb/pkg/util"
)

type Config struct {
	Addr  string            `yaml:"addr"`
	Cron  map[string]string `yaml:"cron"`
	Log   LogConfig         `yaml:"log"`
	Mysql string            `yaml:"mysql" default:"root:123456@tcp(localhost:3306)/test?charset=utf8mb4"`
	Redis string            `yaml:"redis" default:"1localhost:6379"`
	MQ string `yaml:"mq"`
}

type LogConfig struct {
	SingleCapacity int    `yaml:"single_capacity" default:"0"`
	RuntimePath    string `yaml:"runtime_path" default:"./"`
	LimitSize      int    `yaml:"limit_size" default:"1024"`
	LimitDay       int    `yaml:"limit_dat" default:"30"`
}

func NewConfig() *Config {
	var GlobalConfig = new(Config)
	yamlFile, err := ioutil.ReadFile("/home/spider1998/goweb/pkg/config/conf.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, GlobalConfig)
	if err != nil {
		panic(err)
	}
	util.ParseTagReflect(GlobalConfig, "default")
	return GlobalConfig
}
