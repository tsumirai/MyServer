package common

import (
	"MyServer/src/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type BaseController struct {
}

func (c *BaseController) GetDB() *gorm.DB {
	return database.DB
}

func (c *BaseController) EchoErrorStruct(ctx *gin.Context, errorStruct *BaseError) {
	ctx.String(errorStruct.ErrNo, errorStruct.ErrMsg)
}

func (c *BaseController) EchoSuccess(ctx *gin.Context) {
	ctx.String(http.StatusOK, "")
}
