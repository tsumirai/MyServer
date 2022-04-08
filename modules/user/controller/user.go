package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/model"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*common.BaseController
}

func (u *UserController) Ping(ctx *gin.Context) {
	u.EchoSuccess(ctx, "PONG")
}

func (u *UserController) RegisterUser(ctx *gin.Context) {
	var userData model.User
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"msg": "RegisterUser Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, logger.LogArgs{"name": userData.Name, "nickName": userData.NickName, "city": userData.City, "sex": userData.Sex, "birthDay": userData.Birthday})

	u.EchoSuccess(ctx, "")
}
