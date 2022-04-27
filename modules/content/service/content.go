package service

import (
	"MyServer/common"
	"MyServer/modules/content/dto"
	"github.com/gin-gonic/gin"
)

type ContentService struct {
	common.BaseService
}

func NewContentService() *ContentService {
	return &ContentService{}
}

// CreateContent 创建内容
func (s *ContentService) CreateContent(ctx *gin.Context, param *dto.CreateContentReq) {

}
