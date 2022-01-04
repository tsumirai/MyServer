package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(R *gin.Engine) {
	R.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world!")
	})
}
