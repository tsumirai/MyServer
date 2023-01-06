package service

import (
	"MyServer/middleware/logger"
	"MyServer/modules/content/dao"
	"MyServer/util"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// getContentDataByIDsAndAuthorUIDCallback 根据authorUID和contentID获得内容数据
func (s *contentService) getContentDataByIDsAndAuthorUIDCallback(ctx context.Context, key string, subKey ...string) (map[string][]byte, error) {
	contentDao := dao.NewContentDao()
	result := make(map[string][]byte)

	keys := strings.Split(key, ":")
	if len(keys) != 3 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "getContentDataByIDsAndAuthorUIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
		return result, err
	}

	authorUID, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		logger.Error(ctx, "getContentDataByIDsAndAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "parse uid failed", "uid": keys[2]})
		return result, err
	}

	if len(subKey) == 0 {
		return result, nil
	}

	contentIDs, err := util.ConvertStringSliceToInt64Slice(ctx, subKey)
	if err != nil {
		logger.Error(ctx, "getContentDataByIDsAndAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "解析contentID失败"})
		return result, err
	}

	contentData, err := contentDao.GetContentsByIDsAndAuthorUID(ctx, contentIDs, authorUID)
	if err != nil {
		logger.Error(ctx, "getContentDataByIDsAndAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "查询用户信息失败", "ids": contentIDs, "uid": authorUID})
		return result, err
	}

	for _, v := range contentData {
		jsonData, err := json.Marshal(v)
		if err != nil {
			logger.Error(ctx, "getContentDataByIDsAndAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "json序列化失败", "ids": contentIDs, "uid": authorUID})
			return result, err
		}
		result[strconv.FormatInt(v.ID, 64)] = jsonData
	}

	return result, nil
}

// getContentDataByIDCallback 根据内容ID获得内容数据的回调函数
//func (s *contentService) getContentIDsByIDCallback(ctx context.Context, key string, subKey ...string) ([]byte, error) {
//	contentDao := dao.NewContentDao()
//
//	keys := strings.Split(key, ":")
//	if len(keys) != 5 {
//		err := fmt.Errorf("redisKey长度有误")
//		logger.Error(ctx, "GetContentDataByIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
//		return nil, err
//	}
//
//	ID, err := strconv.ParseInt(keys[2], 10, 64)
//	if err != nil {
//		logger.Error(ctx, "GetContentDataByIDCallback", logger.LogArgs{"err": err, "msg": "parse id failed", "id": keys[2]})
//		return nil, err
//	}
//
//	authorUID, err := strconv.ParseInt(keys[4], 10, 64)
//	if err != nil {
//		logger.Error(ctx, "GetContentDataByIDCallback", logger.LogArgs{"err": err, "msg": "parse authorUID failed", "authorUID": keys[4]})
//		return nil, err
//	}
//
//	contentData, err := contentDao.GetContentByID(ctx, ID, authorUID)
//	if err != nil {
//		logger.Error(ctx, "GetContentDataByIDCallback", logger.LogArgs{"err": err, "msg": "查询用户信息失败", "id": ID, "uid": authorUID})
//		return nil, err
//	}
//
//	result, err := json.Marshal(contentData)
//	if err != nil {
//		logger.Error(ctx, "GetContentDataByIDCallback", logger.LogArgs{"err": err, "msg": "json序列化失败", "id": ID, "uid": authorUID})
//		return nil, err
//	}
//
//	return result, nil
//}

// getContentListByAuthorUIDCallback 根据作者的UID获得内容列表的回调函数
func (s *contentService) getContentIDsByAuthorUIDCallback(ctx context.Context, key string, subKey ...string) ([]byte, error) {
	contentDao := dao.NewContentDao()

	keys := strings.Split(key, ":")
	if len(keys) != 7 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "getContentIDsByAuthorUIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
		return nil, err
	}

	authorUID, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "parse authorUID failed", "authorUID": keys[2]})
		return nil, err
	}

	pageNum, err := strconv.ParseInt(keys[4], 10, 64)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "parse pageNum failed", "pageNum": keys[4]})
		return nil, err
	}

	pageSize, err := strconv.ParseInt(keys[6], 10, 64)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "parse pageSize failed", "pageSize": keys[6]})
		return nil, err
	}

	contentData, err := contentDao.GetContentIDsByAuthorUID(ctx, authorUID, int(pageNum), int(pageSize))
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "查询用户信息失败", "uid": authorUID, "pageNum": pageNum, "pageSize": pageSize})
		return nil, err
	}

	result, err := json.Marshal(contentData)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUIDCallback", logger.LogArgs{"err": err, "msg": "json序列化失败", "uid": authorUID, "pageNum": pageNum, "pageSize": pageSize})
		return nil, err
	}

	return result, nil
}
