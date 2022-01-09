package common

import (
	"MyServer/src/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
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

	echoJson, err := json.Marshal(result)
	if err != nil {
		log.Errorf("EchoErrorStruct Marshal Json Failed: %v", err.Error())
		echoJson = make([]byte, 0)
	}
	ctx.JSON(errorStruct.ErrNo, echoJson)
}

func (c *BaseController) EchoSuccess(ctx *gin.Context, data interface{}) {
	result := EchoResult{
		ErrNo:  http.StatusOK,
		ErrMsg: "",
		Data:   data,
	}
	echoJson, err := json.Marshal(result)
	if err != nil {
		log.Errorf("EchoErrorStruct Marshal Json Failed: %v", err.Error())
		echoJson = make([]byte, 0)
	}
	ctx.JSON(http.StatusOK, echoJson)
}
