package database

import (
	"MyServer/base"
	"MyServer/middleware/logger"
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
)

var RDB *redis.Pool

func InitRedis() {
	RDB = redisPool()
	logger.Info(context.TODO(), "InitRedis", logger.LogArgs{"msg": "InitRedis Success!"})
}

func redisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     1,
		MaxActive:   base.Config.GetInt("redis.pool_size"),
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", base.Config.GetString("redis.ip")+":"+base.Config.GetString("redis.port"))
			if err != nil {
				logger.Error(context.TODO(), "redisPool", logger.LogArgs{"msg": "redisPool failed", "err": err.Error()})
				panic(err)
				return nil, err
			}

			if _, err := conn.Do("AUTH", base.Config.GetString("redis.password")); err != nil {
				conn.Close()
				logger.Error(context.TODO(), "redisPool", logger.LogArgs{"msg": "redisPool failed", "err": err.Error()})
				panic(err)
				return nil, err
			}
			return conn, err
		},
	}
}
