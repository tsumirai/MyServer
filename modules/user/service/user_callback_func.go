package service

import (
	"MyServer/middleware/logger"
	"MyServer/modules/user/dao"
	"MyServer/modules/user/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// GetUserInfoByUIDCallback 获得用户信息的回调函数
func (s *UserService) GetUserInfoByUIDCallback(ctx context.Context, key string, subKey ...string) ([]byte, error) {
	userDao := dao.NewUserDao()

	keys := strings.Split(key, ":")
	if len(keys) != 2 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "GetUserInfoByUID", logger.LogArgs{"err": err, "msg": err.Error()})
		return nil, err
	}

	uid, err := strconv.ParseInt(keys[1], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetUserInfoByUID", logger.LogArgs{"err": err, "msg": "parse uid failed"})
		return nil, err
	}

	userInfo, err := userDao.GetUserInfoByParam(ctx, &model.UserInfo{UID: uid})
	if err != nil {
		logger.Error(ctx, "GetUserInfoByUID", logger.LogArgs{"err": err, "msg": "查询用户信息失败", "uid": uid})
		return nil, err
	}

	result, err := json.Marshal(userInfo)
	if err != nil {
		logger.Error(ctx, "GetUserInfoByUID", logger.LogArgs{"err": err, "msg": "json序列化失败", "uid": uid})
		return nil, err
	}

	return result, nil
}
