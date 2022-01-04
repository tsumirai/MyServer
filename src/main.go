package main

import (
	"MyServer/src/config"
	"MyServer/src/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.InitConfig()
	middleware.InitLogger()
}

func main() {
	R := gin.Default()

	// 调用中间件
	R.Use(middleware.LoggerToFile())

	InitRouter(R)

	log.Info("Server Start!!")
	R.Run(":5455")
}
