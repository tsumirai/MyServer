package controller

import (
	"MyServer/src/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct  {
	*common.BaseController
}

func (u *UserController) RegisterUser(ctx *gin.Context){
	ctx.String(http.StatusOK,"Hello World!")
}