package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/content/consts"
	"MyServer/modules/content/model"
	"context"
	"fmt"
	"strconv"
)

type contentDao struct {
	common.BaseDao
}

func NewContentDao() *contentDao {
	return &contentDao{}
}

func (d *contentDao) getTableName(uid int64) string {
	return consts.ContentTable + "_" + strconv.FormatInt(uid%consts.ContentTableNum, 10)
}

// CreateContent 创建内容
func (d *contentDao) CreateContent(ctx context.Context, param *model.Content) error {
	if param.AuthorUID == 0 {
		err := fmt.Errorf("UID不能为0")
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	err := d.GetDB().Table(d.getTableName(param.AuthorUID)).
		Create(param).Error
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "创建内容失败"})
		return err
	}
	return nil
}

// GetContentsByIDs 根据内容ID批量获得内容数据
func (d *contentDao) GetContentsByIDs(ctx context.Context, IDs []int64, authorUID int64) ([]*model.Content, error) {
	result := make([]*model.Content, 0)
	err := d.GetDB().Table(d.getTableName(authorUID)).
		Where("id in (?)", IDs).
		Order("create_time desc").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetContentsByIDs", logger.LogArgs{"err": err, "msg": "获得内容数据失败"})
		return nil, err
	}

	return result, nil
}

// GetContentByID 根据内容ID获得内容数据
func (d *contentDao) GetContentByID(ctx context.Context, ID, authorUID int64) (*model.Content, error) {
	result := &model.Content{}
	err := d.GetDB().Table(d.getTableName(authorUID)).
		Where("id = ?", ID).
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetContentByID", logger.LogArgs{"err": err, "msg": "获得内容数据失败"})
		return nil, err
	}

	return result, nil
}

// SetContentStatus 设置内容的状态
func (d *contentDao) SetContentStatus(ctx context.Context, contentID, authorUID int64, contentStatus int) error {
	err := d.GetDB().Table(d.getTableName(authorUID)).
		Where("id = ?", contentID).
		Update("content_status", contentStatus).Error
	if err != nil {
		logger.Error(ctx, "SetContentStatus", logger.LogArgs{"err": err, "msg": "设置内容状态失败"})
		return err
	}
	return nil
}

// GetContentIDsByAuthorUID 根据作者UID获得内容id列表
func (d *contentDao) GetContentIDsByAuthorUID(ctx context.Context, authorUID int64, pageNum, pageSize int) ([]int64, error) {
	result := make([]int64, 0)
	err := d.GetDB().Table(d.getTableName(authorUID)).
		Select("id").
		Where("author_uid = ?", authorUID).
		Offset(pageNum * pageSize).
		Limit(pageSize).
		Order("create_time desc").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetContentIDsByAuthorUID", logger.LogArgs{"err": err, "msg": "获得内容id列表失败"})
		return result, err
	}

	return result, nil
}

// SetContentPermission 设置内容的权限
func (d *contentDao) SetContentPermission(ctx context.Context, contentID, authorUID int64, permission int) error {
	err := d.GetDB().Table(d.getTableName(authorUID)).
		Where("id = ?", contentID).
		Update("content_permission", permission).Error
	if err != nil {
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "msg": "设置内容权限失败"})
		return err
	}
	return nil
}
