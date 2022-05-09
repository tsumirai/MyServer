package dao

import (
	"MyServer/cache"
	"MyServer/middleware/logger"
	"MyServer/modules/content/model"
	"context"
	"encoding/json"
	"strconv"
)

// UpdateContentCache 更新内容缓存
func (d *contentDao) UpdateContentCache(ctx context.Context, content *model.Content) error {
	cacheSvr := cache.NewCache()
	exit, err := cacheSvr.Exists(cache.GetContentDataByIDsRedisKey(content.AuthorUID))
	if err != nil {
		logger.Error(ctx, "UpdateContentCache", logger.LogArgs{"msg": "获取键失败", "err": err})
		return err
	}

	if !exit {
		logger.Warn(ctx, "UpdateContentCache", logger.LogArgs{"msg": "键不存在", "err": err})
		return nil
	}

	value, err := json.Marshal(content)
	if err != nil {
		logger.Error(ctx, "UpdateContentCache", logger.LogArgs{"msg": "序列化内容数据失败", "err": err})
		return err
	}

	err = cacheSvr.HSet(cache.GetContentDataByIDsRedisKey(content.AuthorUID), strconv.FormatInt(content.ID, 64), value)
	if err != nil {
		logger.Error(ctx, "UpdateContentCache", logger.LogArgs{"msg": "更新缓存失败", "err": err})
		return err
	}

	return nil
}
