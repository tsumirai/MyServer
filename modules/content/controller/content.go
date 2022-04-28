package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/content/dto"
	"MyServer/modules/content/service"

	"github.com/gin-gonic/gin"
)

type ContentController struct {
	common.BaseController
}

// CreateContent 创建内容
func (c *ContentController) CreateContent(ctx *gin.Context) {
	var param *dto.CreateContentReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "解析参数失败"})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, "CreateContent", logger.LogArgs{"title": param.Title, "content": param.Content, "image_urls": param.ImageUrls, "video_url": param.VideoUrl, "author_uid": param.AuthorUID, "content_type": param.ContentType, "content_space_type": param.ContentSpaceType})

	contentSvr := service.NewContentService()
	err = contentSvr.CreateContent(ctx, param)
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "创建内容失败"})
		c.EchoErrorStruct(ctx, common.ErrCreateContentFailed)
		return
	}

	c.EchoSuccess(ctx, nil)
}
