package service

import (
	"MyServer/cache"
	"MyServer/common"
	commonConsts "MyServer/consts"
	"MyServer/middleware/logger"
	"MyServer/modules/message/dto"
	"MyServer/modules/message/model"
	"MyServer/util"
	"context"
	"encoding/json"
	"fmt"
)

type messageService struct {
	common.BaseService
}

func NewMessageService() *messageService {
	return &messageService{}
}

// GetMessageListByUID 根据接收者的uid获得消息
func (s *messageService) GetMessageListByUID(ctx context.Context, receiverUID int64, pageNum, pageSize int) ([]*dto.Message, error) {
	var err error
	result := make([]*dto.Message, 0)
	if receiverUID <= 0 {
		err = fmt.Errorf("参数非法")
		logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "receiverUID": receiverUID})
		return result, err
	}

	// 获得消息的id
	cacheSvr := cache.NewCache()
	cacheSvr.RegisterCallbackFunc(s.GetMessageIDListByUIDCallback)

	messageIDByte, err := cacheSvr.GetValueFromCache(ctx, cache.GetMessageIDsByUID(receiverUID, pageNum, pageSize), commonConsts.FiveMinute)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "receiverUID": receiverUID})
		return result, err
	}

	messageIDs := make([]int64, 0)
	err = json.Unmarshal(messageIDByte, &messageIDs)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "receiverUID": receiverUID})
		return result, err
	}

	messageIDStr := util.ConvertInt64SliceToStringSlice(messageIDs)

	// 获得消息的具体数据
	cacheSvr.RegisterMultiCallbackFunc(s.GetMessageListByUIDCallback)
	messageByte, err := cacheSvr.GetValuesFromHashCache(ctx, cache.GetMessageDataByIDsRedisKey(receiverUID), commonConsts.FiveMinute, messageIDStr...)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "receiverUID": receiverUID, "messageIDStr": messageIDStr})
		return result, err
	}

	for _, v := range messageByte {
		var tempMessage *model.Message
		err = json.Unmarshal(v, &tempMessage)
		if err != nil {
			logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "receiverUID": receiverUID, "messageIDStr": messageIDStr})
			continue
		}
		result = append(result, s.convertMessageToDto(tempMessage))
	}

	return result, nil
}

/* ========================================================================== */
func (s *messageService) convertMessageToDto(param *model.Message) *dto.Message {
	result := &dto.Message{
		ID:            param.ID,
		SenderUID:     param.SenderUID,
		ReceiverUID:   param.ReceiverUID,
		MessageType:   param.MessageType,
		RelateID:      param.RelateID,
		MessageStatus: param.MessageStatus,
		ExtraInfo:     param.ExtraInfo,
		CreateTime:    param.CreateTime,
		UpdateTime:    param.UpdateTime,
	}
	return result
}
