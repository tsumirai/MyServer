package controller

import (
	"MyServer/src/common"
	"MyServer/src/modules/user/model"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type UserController struct {
	*common.BaseController
}

func (u *UserController) RegisterUser(ctx *gin.Context) {
	var userData model.User
	err := ctx.BindJSON(&userData)
	if err != nil {
		log.Errorf("RegisterUser Failed: %v", err.Error())
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	fmt.Println(userData.Name, userData.NickName, userData.City, userData.Sex, userData.Birthday)
	u.EchoSuccess(ctx, "")
}
