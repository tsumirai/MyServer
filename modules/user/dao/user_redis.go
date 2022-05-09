package dao

import (
	"MyServer/cache"
	"MyServer/middleware/logger"
	"context"
)

// DelUserInfoRedisByUID 根据uid删除用户缓存
func (d *userDao) DelUserInfoRedisByUID(ctx context.Context, uid int64) error {
	cacheSvr := cache.NewCache()
	exit, err := cacheSvr.Exists(cache.GetUserInfoRedisKey(uid))
	if err != nil {
		logger.Error(ctx, "DelUserInfoRedis", logger.LogArgs{"err": err.Error(), "msg": "获取redis数据失败", "uid": uid})
		return err
	}

	if exit {
		err := cacheSvr.Del(cache.GetUserInfoRedisKey(uid))
		if err != nil {
			logger.Error(ctx, "DelUserInfoRedis", logger.LogArgs{"err": err.Error(), "msg": "删除用户缓存失败", "uid": uid})
			return err
		}
	}
	return nil
}
