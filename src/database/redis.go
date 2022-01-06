package database

import (
	"MyServer/src/config"
	"fmt"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:  config.Config.GetString("redis.ip")+":"+config.Config.GetString("redis.port"),
		Password: config.Config.GetString("redis.password"),
		DB: config.Config.GetInt("redis.db"),
		PoolSize: config.Config.GetInt("redis.poolsize"),
	})

	_,err := rdb.Ping().Result()
	if err != nil{
		fmt.Println("InitRedis failed: ",err.Error())
		panic(err)
	}

	log.Info("InitRedis Success!")
}
