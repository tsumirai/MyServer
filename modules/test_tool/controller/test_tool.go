package controller

import (
	"MyServer/common"
	"github.com/gin-gonic/gin"
)

type TestToolController struct {
	common.BaseController
}

func NewTestToolController() *TestToolController {
	return &TestToolController{}
}

func (c *TestToolController) Ping(ctx *gin.Context) {
	c.EchoSuccess(ctx, "Pong")
}
