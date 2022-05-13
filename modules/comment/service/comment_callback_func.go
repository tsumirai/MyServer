package service

import (
	"MyServer/middleware/logger"
	"MyServer/modules/comment/dao"
	"MyServer/util"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// GetCommentDataByIDCallback 根据contentID和commentID获得评论数据的回调函数
func (s *commentService) GetCommentDataByIDCallback(ctx context.Context, key string, subKeys ...string) (map[string][]byte, error) {
	commentDao := dao.NewCommentDao()
	result := make(map[string][]byte)

	keys := strings.Split(key, ":")
	if len(keys) != 3 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "GetCommentDataByIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
		return nil, err
	}

	contentID, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetCommentDataByIDCallback", logger.LogArgs{"err": err, "msg": "parse contentID failed", "contentID": keys[2]})
		return nil, err
	}

	if len(subKeys) == 0 {
		return result, nil
	}

	commentIDs, err := util.ConvertStringSliceToInt64Slice(ctx, subKeys)
	if err != nil {
		logger.Error(ctx, "GetCommentDataByIDCallback", logger.LogArgs{"err": err, "msg": "parse commentID failed"})
		return nil, err
	}

	commentData, err := commentDao.GetCommentsByCommentIDs(ctx, contentID, commentIDs)
	if err != nil {
		logger.Error(ctx, "GetCommentDataByIDCallback", logger.LogArgs{"err": err, "msg": "获得评论失败", "contentID": contentID, "commentIDs": commentIDs})
		return nil, err
	}

	for _, v := range commentData {
		jsonData, err := json.Marshal(v)
		if err != nil {
			logger.Error(ctx, "GetCommentDataByIDCallback", logger.LogArgs{"err": err, "msg": "json序列化失败", "contentID": contentID, "commentID": v})
			return nil, err
		}
		result[strconv.FormatInt(v.ID, 64)] = jsonData
	}

	return result, nil
}

// GetCommentIDByContentIDCallback 根据内容id获得评论id的回调函数
func (s *commentService) GetCommentIDByContentIDCallback(ctx context.Context, key string, subKey ...string) ([]byte, error) {
	commentDao := dao.NewCommentDao()

	keys := strings.Split(key, ":")
	if len(keys) != 7 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "GetCommentIDByContentIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
		return nil, err
	}

	contentID, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetCommentIDByContentIDCallback", logger.LogArgs{"err": err, "msg": "parse contentID failed", "contentID": keys[2]})
		return nil, err
	}

	pageNum, err := strconv.ParseInt(keys[4], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetCommentIDByContentIDCallback", logger.LogArgs{"err": err, "msg": "parse pageNum failed", "pageNum": keys[4]})
		return nil, err
	}

	pageSize, err := strconv.ParseInt(keys[6], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetCommentIDByContentIDCallback", logger.LogArgs{"err": err, "msg": "parse pageSize failed", "pageSize": keys[6]})
		return nil, err
	}

	commentIDs, err := commentDao.GetCommentIDByContentID(ctx, contentID, int(pageNum), int(pageSize))
	if err != nil {
		logger.Error(ctx, "GetCommentIDByContentIDCallback", logger.LogArgs{"err": err, "msg": "获得评论失败", "contentID": contentID, "pageNum": pageNum, "pageSize": pageSize})
		return nil, err
	}

	result, err := json.Marshal(commentIDs)
	if err != nil {
		logger.Error(ctx, "GetCommentIDByContentIDCallback", logger.LogArgs{"err": err, "msg": "json序列化失败", "contentID": contentID, "pageNum": pageNum, "pageSize": pageSize})
		return nil, err
	}

	return result, nil
}
