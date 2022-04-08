package main

import (
	config "MyServer/conf"
	"MyServer/database"
	"MyServer/middleware/logger"
	"MyServer/middleware/recover"
	"MyServer/router"
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

	router.InitRouter(R)

	logger.Info(logger.LogArgs{"msg": "Server Start!!"})

	port := config.Config.GetString("server.port")
	if port == "" {
		port = "8991"
	}
	err := R.Run(":" + port)
	if err != nil {
		logger.Fatal(logger.LogArgs{"msg": "启动服务失败", "err": err.Error()})
		os.Exit(1)
	}
}
