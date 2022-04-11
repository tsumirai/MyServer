package database

import (
	"MyServer/base"
	"MyServer/middleware/logger"
	"context"

	"github.com/go-redis/redis"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     base.Config.GetString("redis.ip") + ":" + base.Config.GetString("redis.port"),
		Password: base.Config.GetString("redis.password"),
		DB:       base.Config.GetInt("redis.db"),
		PoolSize: base.Config.GetInt("redis.pool_size"),
	})

	_, err := RDB.Ping().Result()
	if err != nil {
		logger.Error(context.TODO(), logger.LogArgs{"msg": "InitRedis failed", "err": err.Error()})
		panic(err)
	}

	logger.Info(context.TODO(), logger.LogArgs{"msg": "InitRedis Success!"})
}
