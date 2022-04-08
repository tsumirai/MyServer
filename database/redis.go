package database

import (
	config "MyServer/conf"
	"MyServer/middleware/logger"
	"context"

	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     config.Config.GetString("redis.ip") + ":" + config.Config.GetString("redis.port"),
		Password: config.Config.GetString("redis.password"),
		DB:       config.Config.GetInt("redis.db"),
		PoolSize: config.Config.GetInt("redis.pool_size"),
	})

	_, err := RDB.Ping().Result()
	if err != nil {
		logger.Error(context.TODO(), logger.LogArgs{"msg": "InitRedis failed", "err": err.Error()})
		panic(err)
	}

	logger.Info(context.TODO(), logger.LogArgs{"msg": "InitRedis Success!"})
}
