package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/message/dto"
	"MyServer/modules/message/service"
	"github.com/gin-gonic/gin"
)

type messageController struct {
	common.BaseController
}

func NewMessageController() *messageController {
	return &messageController{}
}

// GetMessageListByUID 获得消息列表
func (c *messageController) GetMessageListByUID(ctx *gin.Context) {
	var param *dto.GetMessageListByUID
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "msg": "参数解析失败"})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	messageSvr := service.NewMessageService()
	messages, err := messageSvr.GetMessageListByUID(ctx, param.UID, param.PageNum, param.PageSize)
	if err != nil {
		logger.Error(ctx, "GetMessageListByUID", logger.LogArgs{"err": err, "msg": "获得消息数据失败"})
		c.EchoErrorStruct(ctx, common.ErrGetMessageFailed)
		return
	}

	c.EchoSuccess(ctx, messages)
}

// UpdateMessageStatus 修改消息状态
func (c *messageController) UpdateMessageStatus(ctx *gin.Context) {

}
