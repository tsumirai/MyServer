package controller

import (
	"MyServer/common"
	"github.com/gin-gonic/gin"
)

type interactController struct {
	common.BaseController
}

func NewInteractController() *interactController {
	return &interactController{}
}

// Like 点赞、取消点赞
func (c *interactController) Like(ctx *gin.Context) {}

// Report 举报用户、举报内容、举报评论
func (c *interactController) Report(ctx *gin.Context) {}

// FollowUser 关注用户
func (c *interactController) FollowUser(ctx *gin.Context) {}
