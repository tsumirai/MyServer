package main

import (
	"MyServer/src/modules/user/controller"
	"github.com/gin-gonic/gin"
)

func InitRouter(R *gin.Engine) {
	R.GET("/", new(controller.UserController).RegisterUser)
}
