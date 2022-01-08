package controller

import (
	"MyServer/src/common"
	"MyServer/src/modules/user/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct  {
	*common.BaseController
}

func (u *UserController) RegisterUser(ctx *gin.Context){
	var userData model.User
	ctx.BindJSON(&userData)
	ctx.String(http.StatusOK,"Hello World! %v %v %v %v %v",userData.Name,userData.NickName,userData.Sex,userData.Birthday,userData.City)
}