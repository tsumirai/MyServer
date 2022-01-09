package main

import (
	"MyServer/src/config"
	"MyServer/src/database"
	"MyServer/src/middleware/logutil"
	"github.com/gin-gonic/gin"
)

func init() {
	config.InitConfig()
	logutil.InitLogger()
	database.InitMysql()
	database.InitRedis()
}

func main() {
	R := gin.Default()

	// 调用中间件
	R.Use(logutil.LoggerToFile())

	InitRouter(R)

	logutil.Info("Server Start!!")
	R.Run(":5455")
}
