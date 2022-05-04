package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/content/consts"
	"MyServer/modules/content/model"
	"context"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ContentDao struct {
	common.BaseDao
}

func NewContentDao() *ContentDao {
	return &ContentDao{}
}

func (d *ContentDao) getTableName(ctx context.Context, uid int64) string {
	return consts.ContentTable + "_" + strconv.FormatInt(uid%consts.ContentTableNum, 10)
}

// CreateContent 创建内容
func (d *ContentDao) CreateContent(ctx *gin.Context, param *model.Content) error {
	if param.AuthorUID == 0 {
		err := fmt.Errorf("UID不能为0")
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	err := d.GetDB().Table(d.getTableName(ctx, param.AuthorUID)).Create(param).Error
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "创建内容失败"})
		return err
	}
	return nil
}
