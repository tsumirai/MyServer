package database

import (
	"MyServer/src/config"
	"MyServer/src/middleware/logutil"
	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.GetString("redis.ip") + ":" + config.Config.GetString("redis.port"),
		Password: config.Config.GetString("redis.password"),
		DB:       config.Config.GetInt("redis.db"),
		PoolSize: config.Config.GetInt("redis.poolsize"),
	})

	_, err := RDB.Ping().Result()
	if err != nil {
		logutil.Error("InitRedis failed: ", err.Error())
		panic(err)
	}

	logutil.Info("InitRedis Success!")
}
