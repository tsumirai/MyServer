package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/test_tool/dto"
	"MyServer/modules/test_tool/service"

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

// DivideTable 建立分表
func (c *TestToolController) DivideTable(ctx *gin.Context) {
	var param *dto.DivideTableReq
	if err := ctx.BindJSON(&param); err != nil {
		logger.Error(ctx, "DivideTable", logger.LogArgs{"err": err})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	divideSvr := service.NewTestToolService()
	divideSvr.DivideTable(ctx, param)
}
