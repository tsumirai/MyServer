package cache

import (
	"MyServer/database"
	"MyServer/middleware/logger"
	"context"
	"github.com/gomodule/redigo/redis"
)

type Cache struct {
	redisPool    *redis.Pool
	callbackFunc CallbackMysqlFunc
}

func NewCache() *Cache {
	return &Cache{
		redisPool: database.RDB,
	}
}

type CallbackMysqlFunc func(ctx context.Context, key string, subKey ...string) ([]byte, error)

func (c *Cache) RegisterCallbackFunc(callbackFunc CallbackMysqlFunc) {
	c.callbackFunc = callbackFunc
}

// SetDataToRedis 把数据存到redis中
func (c *Cache) SetDataToRedis(ctx context.Context, key, subKey string, value []byte, expireTime int) ([]byte, error) {
	var err error
	if value != nil {
		if subKey == "" {
			if expireTime != 0 {
				err = c.SetEx(key, value, expireTime)
				if err != nil {
					logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "设置缓存失败", "err": err.Error()})
					return []byte{}, err
				}
			} else {
				err = c.Set(key, value)
				if err != nil {
					logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "设置缓存失败", "err": err.Error()})
					return []byte{}, err
				}
			}
		} else {
			err = c.HSet(key, subKey, value)
			if err != nil {
				logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "设置缓存失败", "err": err.Error()})
				return []byte{}, err
			}

			err = c.Expire(key, expireTime)
			if err != nil {
				logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "设置缓存失败", "err": err.Error()})
				return []byte{}, err
			}
		}
		return value, nil
	} else {
		return []byte{}, nil
	}
}

// GetValueFromCache 从redis中获取数据，获取失败则从mysql中获取
func (c *Cache) GetValueFromCache(ctx context.Context, key string, expireTime int) ([]byte, error) {
	result := make([]byte, 0)
	exit, err := c.Exists(key)
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "查询键失败", "err": err.Error()})
		result, err = c.callbackFunc(ctx, key)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, "", result, expireTime)
	}

	if !exit {
		result, err = c.callbackFunc(ctx, key)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, "", result, expireTime)
	}

	resultStr, err := c.Get(key)
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
		result, err = c.callbackFunc(ctx, key)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, "", result, expireTime)
	}

	return []byte(resultStr), nil
}

// GetValueFromHashCache 从hash缓存中取数据
func (c *Cache) GetValueFromHashCache(ctx context.Context, key, subKey string, expireTime int) ([]byte, error) {
	result := make([]byte, 0)
	exit, err := c.Exists(key)
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "查询键失败", "err": err.Error()})
		result, err = c.callbackFunc(ctx, key, subKey)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
	}

	if !exit {
		result, err = c.callbackFunc(ctx, key, subKey)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
	}

	resultStr, err := c.HGet(key, subKey)
	if err != nil {
		logger.Error(ctx, "GetValueFromHashCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
		result, err = c.callbackFunc(ctx, key)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
	}

	return []byte(resultStr), nil
}

func (c *Cache) Set(key string, value interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("set", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Get(key string) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Cache) SetEx(key string, value interface{}, expire int) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("setex", key, expire, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) SetNx(key string, value interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("setnx", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Exists(key string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	exist, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (c *Cache) HSet(key, subKey string, value interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("hset", key, subKey, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) HGet(key, subKey string) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("hget", key, subKey))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *Cache) Expire(key string, expire int) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("expire", key, expire)
	if err != nil {
		return err
	}
	return nil
}
