package recover

import (
	"MyServer/middleware/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recover 自定义recover
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				http.Error(c.Writer, "500 Server internal error", http.StatusInternalServerError)
				logger.Error(c, "Recover", logger.LogArgs{"msg": "程序panic", "err": r, "stack": string(debug.Stack())})
			}
		}()
		c.Next()
	}
}
