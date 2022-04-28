package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/content/consts"
	"MyServer/modules/content/model"

	"github.com/gin-gonic/gin"
)

type ContentDao struct {
	common.BaseDao
}

func NewContentDao() *ContentDao {
	return &ContentDao{}
}

// CreateContent 创建内容
func (d *ContentDao) CreateContent(ctx *gin.Context, param *model.Content) error {
	err := d.GetDB().Table(consts.ContentTable).Create(param).Error
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "创建内容失败"})
		return err
	}
	return nil
}
