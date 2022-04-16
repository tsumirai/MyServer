package router

import (
	"MyServer/modules/user/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(R *gin.Engine) {
	R.GET("/", new(controller.UserController).Ping)
	R.POST("/userLogin", new(controller.UserController).UserLogin)           // 用户登录，新用户自动注册
	R.POST("/getUserInfo", new(controller.UserController).GetUserInfo)       // 获取用户信息
	R.POST("/updateUserInfo", new(controller.UserController).UpdateUserInfo) // 更新用户信息
}
