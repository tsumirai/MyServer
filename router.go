package main

import (
	"MyServer/modules/user/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(R *gin.Engine) {
	R.GET("/", new(controller.UserController).Ping)
	R.POST("/registerUser", new(controller.UserController).RegisterUser)
}
