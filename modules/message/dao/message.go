package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/message/consts"
	"MyServer/modules/message/model"
	"context"
	"fmt"
	"strconv"
)

type messageDao struct {
	common.BaseDao
}

func NewMessageDao() *messageDao {
	return &messageDao{}
}

func (d *messageDao) getTableName(receiverUID int64) string {
	return fmt.Sprintf(consts.MessageTable + "_" + strconv.FormatInt(receiverUID%consts.MessageTableNum, 10))
}

// GetMessageIDsByReceiverUID 根据接收者的uid获得消息id
func (d *messageDao) GetMessageIDsByReceiverUID(ctx context.Context, receiverUID int64, pageNum, pageSize int) ([]*model.Message, error) {
	result := make([]*model.Message, 0)
	err := d.GetDB().Table(d.getTableName(receiverUID)).
		Where("receiver_uid = ?", receiverUID).
		Offset(pageNum * pageSize).
		Limit(pageSize).
		Order("create_time").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetMessageIDsByReceiverUID", logger.LogArgs{"err": err})
		return result, err
	}

	return result, nil
}

// GetMessageDataByIDs 根据messageID批量获得消息数据
func (d *messageDao) GetMessageDataByIDs(ctx context.Context, receiverUID int64, messageIDs []int64) ([]*model.Message, error) {
	result := make([]*model.Message, 0)
	err := d.GetDB().Table(d.getTableName(receiverUID)).
		Where("id in (?)", messageIDs).
		Where("receiver_uid = ?", receiverUID).
		Order("create_time").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetMessageDataByIDs", logger.LogArgs{"err": err, "receiverUID": receiverUID, "messageIDs": messageIDs})
		return result, err
	}
	return result, nil
}
