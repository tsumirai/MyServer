package service

import (
	"MyServer/middleware/logger"
	"MyServer/modules/message/dao"
	"MyServer/util"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// GetMessageIDListByUIDCallback 根据接收者UID获得消息id列表的回调函数
func (s *messageService) GetMessageIDListByUIDCallback(ctx context.Context, key string, subKeys ...string) ([]byte, error) {
	messageDao := dao.NewMessageDao()
	result := make([]byte, 0)

	keys := strings.Split(key, ":")
	if len(keys) != 7 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
		return result, err
	}

	receiverUID, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "parse receiverUID failed", "receiverUID": keys[2]})
		return result, err
	}

	pageNum, err := strconv.ParseInt(keys[4], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "parse pageNum failed", "pageNum": keys[4]})
		return result, err
	}

	pageSize, err := strconv.ParseInt(keys[6], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "parse pageSize failed", "pageNum": keys[6]})
		return result, err
	}

	messages, err := messageDao.GetMessageIDsByReceiverUID(ctx, receiverUID, int(pageNum), int(pageSize))
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "receiverUID": receiverUID, "pageNum": pageNum, "pageSize": pageSize})
		return result, err
	}

	result, err = json.Marshal(messages)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "receiverUID": receiverUID, "pageNum": pageNum, "pageSize": pageSize})
		return result, err
	}

	return result, nil
}

// GetMessageListByUIDCallback 根据接收者UID获得消息列表的回调函数
func (s *messageService) GetMessageListByUIDCallback(ctx context.Context, key string, subKeys ...string) (map[string][]byte, error) {
	messageDao := dao.NewMessageDao()
	result := make(map[string][]byte, 0)

	keys := strings.Split(key, ":")
	if len(keys) != 3 {
		err := fmt.Errorf("redisKey长度有误")
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": err.Error()})
		return result, err
	}

	receiverUID, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "parse receiverUID failed", "receiverUID": keys[2]})
		return result, err
	}

	if len(subKeys) == 0 {
		return result, nil
	}

	messageIDs, err := util.ConvertStringSliceToInt64Slice(ctx, subKeys)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "parse messageIDs failed", "subKeys": subKeys})
		return result, err
	}

	messages, err := messageDao.GetMessageDataByIDs(ctx, receiverUID, messageIDs)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "获得消息数据失败", "receiverUID": receiverUID, "messageIDs": messageIDs})
		return result, err
	}

	for _, v := range messages {
		temp, err := json.Marshal(v)
		if err != nil {
			logger.Error(ctx, "GetMessageListByUIDCallback", logger.LogArgs{"err": err, "msg": "序列化消息数据失败", "receiverUID": receiverUID, "messageID": v.ID})
			continue
		}
		result[strconv.FormatInt(v.ID, 10)] = temp
	}

	return result, nil
}
