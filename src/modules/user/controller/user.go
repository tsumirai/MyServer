package controller

import (
	"MyServer/src/common"
	"MyServer/src/middleware/logutil"
	"MyServer/src/modules/user/model"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	*common.BaseController
}

func (u *UserController) RegisterUser(ctx *gin.Context) {
	var userData model.User
	err := ctx.BindJSON(&userData)
	if err != nil {
		logutil.Errorf("RegisterUser Failed: %v", err.Error())
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logutil.Info(userData.Name, userData.NickName, userData.City, userData.Sex, userData.Birthday)
	u.EchoSuccess(ctx, "")
}
