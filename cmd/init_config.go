package main

import (
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/rest"
	"log"
	"net_disk/middleware"
	"xorm.io/xorm"
)

type Config struct {
	rest.RestConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr string
	}
}

type ServiceContext struct {
	Config      Config
	Engine      *xorm.Engine
	RedisEngine *redis.Client
	Auth        rest.Middleware
}

func NewServiceContext(c Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		Engine:      initMysql(c.Mysql.DataSource),
		RedisEngine: initRedis(c.Redis.Addr),
		Auth:        middleware.NewAuthMiddleware().Handle,
	}
}

func initMysql(datasource string) *xorm.Engine {
	engine, err := xorm.NewEngine("mysql", datasource)
	if err != nil {
		panic(err)
	}
	if err := engine.Ping(); err != nil {
		log.Println("xorm connect mysql fail")
		return nil
	}
	return engine
}

func initRedis(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
