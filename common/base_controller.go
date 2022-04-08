package common

import (
	"MyServer/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type BaseController struct {
}

type EchoResult struct {
	ErrNo  int         `json:"err_no"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func (c *BaseController) GetDB() *gorm.DB {
	return database.DB
}

func (c *BaseController) EchoErrorStruct(ctx *gin.Context, errorStruct *BaseError) {
	result := EchoResult{
		ErrNo:  errorStruct.ErrNo,
		ErrMsg: errorStruct.ErrMsg,
		Data:   "",
	}

	ctx.JSON(errorStruct.ErrNo, result)
}

func (c *BaseController) EchoSuccess(ctx *gin.Context, data interface{}) {
	result := EchoResult{
		ErrNo:  0,
		ErrMsg: "",
		Data:   data,
	}

	ctx.JSON(http.StatusOK, result)
}
