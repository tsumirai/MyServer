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

type CallbackMysqlFunc func(ctx context.Context, key string) (string, error)

func (c *Cache) RegisterCallbackFunc(callbackFunc CallbackMysqlFunc) {
	c.callbackFunc = callbackFunc
}

// GetValueFromCache 从redis中获取数据，获取失败则从mysql中获取
func (c *Cache) GetValueFromCache(ctx context.Context, key string) (string, error) {
	exit, err := c.Exists(key)
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "查询键失败", "err": err.Error()})
		return c.callbackFunc(ctx, key)
	}

	if !exit {
		return c.callbackFunc(ctx, key)
	}

	result, err := c.Get(key)
	if err != nil {
		logger.Error(ctx, "GetValueFromCache", logger.LogArgs{"msg": "从redis中获取数据失败", "err": err.Error()})
		return c.callbackFunc(ctx, key)
	}

	return result, nil
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
