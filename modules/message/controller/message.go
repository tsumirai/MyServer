package controller

import (
	"MyServer/common"
	"github.com/gin-gonic/gin"
)

type messageController struct {
	common.BaseController
}

func NewMessageController() *messageController {
	return &messageController{}
}

// GetMessageList 获得消息列表
func (c *messageController) GetMessageList(ctx *gin.Context) {}
