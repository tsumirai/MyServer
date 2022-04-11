package context

import (
	"MyServer/util"
	"github.com/gin-gonic/gin"
)

// InitContext 初始化context
func InitContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("trace_id", util.GenerateTraceID())
	}
}
