package service

import (
	"MyServer/common"
	"MyServer/modules/test_tool/dto"

	"github.com/gin-gonic/gin"
)

type TestToolService struct {
	common.BaseService
}

func NewTestToolService() *TestToolService {
	return &TestToolService{}
}

func (s *TestToolService) DivideTable(ctx *gin.Context, param *dto.DivideTableReq) {

}
