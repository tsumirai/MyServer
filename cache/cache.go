package cache

import (
	"MyServer/database"
	"MyServer/middleware/logger"
	"context"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type cache struct {
	redisPool         *redis.Pool
	callbackFunc      CallbackMysqlFunc
	multiCallbackFunc MultiCallbackMysqlFunc
}

func NewCache() *cache {
	return &cache{
		redisPool: database.RDB,
	}
}

type CallbackMysqlFunc func(ctx context.Context, key string, subKey ...string) ([]byte, error)

type MultiCallbackMysqlFunc func(ctx context.Context, key string, subKey ...string) (map[string][]byte, error)

func (c *cache) RegisterCallbackFunc(callbackFunc CallbackMysqlFunc) {
	c.callbackFunc = callbackFunc
}

func (c *cache) RegisterMultiCallbackFunc(multiCallbackFunc MultiCallbackMysqlFunc) {
	c.multiCallbackFunc = multiCallbackFunc
}

// SetDataToRedis 把数据存到redis中
func (c *cache) SetDataToRedis(ctx context.Context, key, subKey string, value []byte, expireTime int) ([]byte, error) {
	var err error
	fmt.Println("===========================开始设置缓存", key, subKey, value, expireTime)
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
func (c *cache) GetValueFromCache(ctx context.Context, key string, expireTime int, subKeys ...string) ([]byte, error) {
	result := make([]byte, 0)
	subKey := ""
	if subKeys != nil && len(subKeys) > 0 {
		subKey = subKeys[0]
	}
	exit, err := c.Exists(key)
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "查询键失败", "err": err.Error()})
		result, err = c.callbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
	}

	if !exit {
		result, err = c.callbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
	}

	var resultStr string
	if subKey != "" {
		resultStr, err = c.HGet(key, subKey)
	} else {
		resultStr, err = c.Get(key)
	}
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
		result, err = c.callbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return []byte{}, err
		}
		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
	}

	return []byte(resultStr), nil
}

// SetHashDataToRedis 批量把数据存到redis的hash中
func (c *cache) SetHashDataToRedis(ctx context.Context, key string, keyValues map[string][]byte, expireTime int) (map[string][]byte, error) {
	var err error
	fmt.Println("===========================开始设置缓存", key, keyValues, expireTime)
	if len(keyValues) != 0 {
		err = c.HMSet(key, keyValues)
		if err != nil {
			logger.Error(ctx, "SetHashDataToRedis", logger.LogArgs{"msg": "设置缓存失败", "err": err.Error()})
			return map[string][]byte{}, err
		}
		if expireTime != 0 {
			err = c.Expire(key, expireTime)
			if err != nil {
				logger.Error(ctx, "SetHashDataToRedis", logger.LogArgs{"msg": "设置缓存失败", "err": err.Error()})
				return map[string][]byte{}, err
			}
		}
		return keyValues, nil
	} else {
		return map[string][]byte{}, nil
	}
}

// GetValuesFromHashCache 从redis中批量获取hash数据，获取失败则从mysql中获取
func (c *cache) GetValuesFromHashCache(ctx context.Context, key string, expireTime int, subKeys ...string) (map[string][]byte, error) {
	result := make(map[string][]byte, 0)
	if len(subKeys) <= 0 {
		return result, nil
	}
	exit, err := c.Exists(key)
	if err != nil {
		logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "查询键失败", "err": err.Error()})
		result, err = c.multiCallbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return result, err
		}
		return c.SetHashDataToRedis(ctx, key, result, expireTime)
	}

	if !exit {
		result, err = c.multiCallbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return result, err
		}
		return c.SetHashDataToRedis(ctx, key, result, expireTime)
	}

	results, err := c.HMGet(key, subKeys...)
	if err != nil {
		logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
		result, err = c.multiCallbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return result, err
		}
		return c.SetHashDataToRedis(ctx, key, result, expireTime)
	}

	if len(subKeys) != len(results) {
		err = fmt.Errorf("从redis中获取数据失败")
		logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
		result, err = c.multiCallbackFunc(ctx, key, subKeys...)
		if err != nil {
			logger.Error(ctx, "GetValuesFromHashCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
			return result, err
		}
		return c.SetHashDataToRedis(ctx, key, result, expireTime)
	}

	for i := 0; i < len(subKeys); i++ {
		result[subKeys[i]] = []byte(results[i])
	}

	return result, nil
}

// GetValueFromHashCache 从hash缓存中取数据
//func (c *Cache) GetValueFromHashCache(ctx context.Context, key, subKey string, expireTime int) ([]byte, error) {
//	result := make([]byte, 0)
//	exit, err := c.Exists(key)
//	if err != nil {
//		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "查询键失败", "err": err.Error()})
//		result, err = c.callbackFunc(ctx, key, subKey)
//		if err != nil {
//			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
//			return []byte{}, err
//		}
//		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
//	}
//
//	if !exit {
//		result, err = c.callbackFunc(ctx, key, subKey)
//		if err != nil {
//			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
//			return []byte{}, err
//		}
//		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
//	}
//
//	resultStr, err := c.HGet(key, subKey)
//	if err != nil {
//		logger.Error(ctx, "GetValueFromHashCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
//		result, err = c.callbackFunc(ctx, key)
//		if err != nil {
//			logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从mysql中获取数据失败", "err": err.Error()})
//			return []byte{}, err
//		}
//		return c.SetDataToRedis(ctx, key, subKey, result, expireTime)
//	}
//
//	return []byte(resultStr), nil
//}

func (c *cache) Set(key string, value interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("set", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Get(key string) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("get", key))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *cache) SetEx(key string, value interface{}, expire int) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("setex", key, expire, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) SetNx(key string, value interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("setnx", key, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Exists(key string) (bool, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	exist, err := redis.Bool(conn.Do("exists", key))
	if err != nil {
		return false, err
	}
	return exist, nil
}

func (c *cache) HSet(key, subKey string, value interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("hset", key, subKey, value)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) HGet(key, subKey string) (string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	result, err := redis.String(conn.Do("hget", key, subKey))
	if err != nil {
		return "", err
	}
	return result, nil
}

func (c *cache) HMGet(key string, subKey ...string) ([]string, error) {
	conn := c.redisPool.Get()
	defer conn.Close()
	result, err := redis.Strings(conn.Do("hmget", key, subKey))
	if err != nil {
		return []string{}, err
	}

	return result, nil
}

func (c *cache) HMSet(key string, keyValue interface{}) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := redis.String(conn.Do("hmset", key, keyValue))
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Expire(key string, expire int) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("expire", key, expire)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) Del(key string) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("del", key)
	if err != nil {
		return err
	}
	return nil
}

func (c *cache) HDel(key string, subKeys ...string) error {
	conn := c.redisPool.Get()
	defer conn.Close()
	_, err := conn.Do("hdel", key, subKeys)
	if err != nil {
		return err
	}
	return nil
}
