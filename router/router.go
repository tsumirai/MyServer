package router

import (
	contentController "MyServer/modules/content/controller"
	userController "MyServer/modules/user/controller"

	"github.com/gin-gonic/gin"
)

func InitRouter(R *gin.Engine) {
	R.GET("/ping", new(userController.UserController).Ping)

	/*------------------- user ---------------------*/
	R.POST("/userLogin", new(userController.UserController).UserLogin)               // 用户登录，新用户自动注册
	R.POST("/getUserInfoByUID", new(userController.UserController).GetUserInfoByUID) // 根据UID获取用户信息
	R.POST("/updateUserInfo", new(userController.UserController).UpdateUserInfo)     // 更新用户信息

	/*-------------------- content -----------------*/
	R.POST("/createContent", new(contentController.ContentController).CreateContent) // 创建内容
}
