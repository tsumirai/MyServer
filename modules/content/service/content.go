package service

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/content/consts"
	"MyServer/modules/content/dao"
	"MyServer/modules/content/dto"
	"MyServer/modules/content/model"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContentService struct {
	common.BaseService
}

func NewContentService() *ContentService {
	return &ContentService{}
}

// CreateContent 创建内容
func (s *ContentService) CreateContent(ctx *gin.Context, param *dto.CreateContentReq) error {
	// 检查图片数量上限
	if len(param.ImageUrls) > consts.ContentImageNumLimit {
		err := fmt.Errorf("图片不能超过%d张", consts.ContentImageNumLimit)
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	if param.AuthorUID == 0 {
		err := fmt.Errorf("未找到用户")
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	// todo 检查用户是否被封禁

	contentDao := dao.NewContentDao()
	createContent := s.ConvertToContentModel(ctx, param)
	err := contentDao.CreateContent(ctx, createContent)
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "创建内容失败"})
		return err
	}

	return nil
}

// ConvertToContentModel 转换content的请求参数为mysql格式
func (s *ContentService) ConvertToContentModel(ctx *gin.Context, param *dto.CreateContentReq) *model.Content {
	result := &model.Content{
		Title:            param.Title,
		Content:          param.Content,
		AuthorUID:        param.AuthorUID,
		ImageUrls:        strings.Join(param.ImageUrls, ";"),
		VideoUrl:         param.VideoUrl,
		ContentType:      param.ContentType,
		ContentSpaceType: param.ContentSpaceType,
		LocCityID:        param.LocCityID,
	}

	return result
}
