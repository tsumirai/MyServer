package main

import (
	config "MyServer/conf"
	"MyServer/database"
	"MyServer/middleware/logger"
	"MyServer/middleware/recover"
	"os"

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
	R.Use(logger.NewLogModel().LoggerToFile(), recover.Recover())

	InitRouter(R)

	logger.Info(logger.LogArgs{"msg": "Server Start!!"})
	err := R.Run(":5455")
	if err != nil {
		logger.Fatal(logger.LogArgs{"msg": "启动服务失败", "err": err.Error()})
		os.Exit(1)
	}
}
