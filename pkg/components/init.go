package components

import (
	"github.com/go-xorm/xorm"
	"goweb/pkg/resource"
)

var (
	Conf  *Config
	Log   Logger
	DB    *xorm.Engine
	Redis *RedisClient
	MQ    *NsqMq
)

func init() {
	var err error
	{
		Conf = NewConfig()
	}
	{
		Log = NewLogger(Conf)
	}

	{
		resource.Load()
	}
	{
		DB, err = NewDB(Conf.Mysql, resource.MigrationBox, Log)
		if err != nil {
			panic(err)
		}

		Log.Infof("db ready.")
	}
	{
		Redis, err = NewRedis(Conf.Redis, 10, Log)
		if err != nil {
			panic(err)
		}

		Log.Infof("redis ready.")
	}
	{
		MQ = NewMq(Conf.MQ,Log)
	}

}
