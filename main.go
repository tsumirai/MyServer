package main

import (
	config "MyServer/conf"
	"MyServer/database"
	"MyServer/middleware/logger"

	"github.com/gin-gonic/gin"
)

func init() {
	config.InitConfig()
	logger.NewLogModel().InitLogger()
	database.InitMysql()
	database.InitRedis()
}

func main() {
	R := gin.Default()

	// 调用中间件
	R.Use(logger.NewLogModel().LoggerToFile())

	InitRouter(R)

	logger.Info(logger.LogArgs{"msg": "Server Start!!"})
	R.Run(":5455")
}
