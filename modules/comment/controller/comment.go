package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/comment/dto"
	"MyServer/modules/comment/service"
	"github.com/gin-gonic/gin"
)

type commentController struct {
	common.BaseController
}

func NewCommentController() *commentController {
	return &commentController{}
}

// CreateComment 创建评论
func (c *commentController) CreateComment(ctx *gin.Context) {
	var param *dto.CreateCommentReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "解析参数失败"})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	commentSvr := service.NewCommentService()
	err = commentSvr.CreateComment(ctx, param)
	if err != nil {
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "创建评论失败"})
		c.EchoErrorStruct(ctx, common.ErrCreateCommentFailed)
		return
	}

	c.EchoSuccess(ctx, nil)
}

// GetCommentsByContentID 获得内容下的评论
func (c *commentController) GetCommentsByContentID(ctx *gin.Context) {
	var param *dto.GetCommentByContentIDReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "GetCommentByContentID", logger.LogArgs{"err": err, "msg": "参数解析失败"})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	commentSvr := service.NewCommentService()
	comments, err := commentSvr.GetCommentsByContentID(ctx, param.ContentID, param.PageNum, param.PageSize)
	if err != nil {
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "获得评论失败"})
		c.EchoErrorStruct(ctx, common.ErrGetCommentFailed)
		return
	}

	c.EchoSuccess(ctx, comments)
}

// GetCommentCountByContentID 获得内容下所有的评论数量
func (c *commentController) GetCommentCountByContentID(ctx *gin.Context) {
	var param *dto.GetCommentContByContentIDReq
	err := ctx.BindJSON(&param)
	if err != nil {
		logger.Error(ctx, "GetCommentCountByContentID", logger.LogArgs{"err": err, "msg": "参数解析失败"})
		c.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	commentSvr := service.NewCommentService()
	commentCount, err := commentSvr.GetCommentCountByContentID(ctx, param.ContentID)
	if err != nil {
		logger.Error(ctx, "GetCommentCountByContentID", logger.LogArgs{"err": err, "msg": "获得评论数量失败"})
		c.EchoErrorStruct(ctx, common.ErrGetCommentCountFailed)
		return
	}

	c.EchoSuccess(ctx, commentCount)
}
