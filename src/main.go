package main

import (
	"MyServer/src/config"
	"MyServer/src/database"
	"MyServer/src/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
	middleware.InitLogger()
	database.InitMysql()
	database.InitRedis()
}

func main() {
	R := gin.Default()

	// 调用中间件
	R.Use(middleware.LoggerToFile())

	InitRouter(R)

	log.Info("Server Start!!")
	R.Run(":5455")
}
