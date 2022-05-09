package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/content/dto"
	"MyServer/modules/content/service"

	"github.com/gin-gonic/gin"
)

type contentController struct {
	common.BaseController
}

func NewContentController() *contentController {
	return &contentController{}
}

// CreateContent 创建内容
func (c *contentController) CreateContent(ctx *gin.Context) {
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

// DeleteContentByID 用户根据ID删除内容
func (c *contentController) DeleteContentByID(ctx *gin.Context) {
	var param *dto.ContentReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	contentSvr := service.NewContentService()
	contentSvr.DeleteContentByID(ctx, param.ID, param.AuthorUID)

}

// Feed 获得内容的feed流
func (c *contentController) Feed(ctx *gin.Context) {

}

// GetContentByID 根据ID获得内容的具体信息
func (c *contentController) GetContentByID(ctx *gin.Context) {
	var param *dto.ContentReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "GetContentByID", logger.LogArgs{"err": err})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	contentSvr := service.NewContentService()
	content, err := contentSvr.GetContentByID(ctx, param.ID, param.AuthorUID)
	if err != nil {
		logger.Error(ctx, "GetContentByID", logger.LogArgs{"err": err, "msg": "获得内容数据失败"})
		c.EchoErrorStruct(ctx, common.ErrGetContentFailed)
		return
	}

	c.EchoSuccess(ctx, content)
}

// GetContentList 获得内容列表（获得收藏的内容，获得点赞的内容）
func (c *contentController) GetContentList(ctx *gin.Context) {
	var param *dto.ContentReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "GetContentList", logger.LogArgs{"err": err})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	contentSvr := service.NewContentService()
	content, err := contentSvr.GetContentList(ctx, param.AuthorUID, param.PageNum, param.PageSize, param.ContentListType)
	if err != nil {
		logger.Error(ctx, "GetContentList", logger.LogArgs{"err": err, "msg": "获得内容数据失败"})
		c.EchoErrorStruct(ctx, common.ErrGetContentFailed)
		return
	}

	c.EchoSuccess(ctx, content)
}

// ShareContent 分享内容
func (c *contentController) ShareContent(ctx *gin.Context) {

	c.EchoSuccess(ctx, nil)
}

// SetContentPermission 设置内容权限
func (c *contentController) SetContentPermission(ctx *gin.Context) {
	var param *dto.ContentPermissionReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "msg": "参数解析失败"})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	contentSvr := service.NewContentService()
	err = contentSvr.SetContentPermission(ctx, param.ID, param.AuthorUID, param.Permission)
	if err != nil {
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "msg": "权限设置失败"})
		c.EchoErrorStruct(ctx, common.ErrSetContentPermissionFailed)
		return
	}

	c.EchoSuccess(ctx, nil)
}

// SetContentSpace 设置内容的空间 0：普通空间  1：隐私空间
func (c *contentController) SetContentSpace(ctx *gin.Context) {
	c.EchoSuccess(ctx, nil)
}
