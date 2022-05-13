package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/comment/consts"
	"MyServer/modules/comment/model"
	"context"
	"strconv"
)

type commentDao struct {
	common.BaseDao
}

func NewCommentDao() *commentDao {
	return &commentDao{}
}

func (d *commentDao) getTableName(contentID int64) string {
	return consts.CommentTable + "_" + strconv.FormatInt(contentID%consts.CommentTableNum, 10)
}

// CreateComment 创建评论
func (d *commentDao) CreateComment(ctx context.Context, param *model.Comment) error {
	err := d.GetDB().Table(d.getTableName(param.ContentID)).
		Create(param).Error
	if err != nil {
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "创建评论失败"})
		return err
	}
	return nil
}

// GetCommentsByContentID 根据contentID获得评论
func (d *commentDao) GetCommentsByContentID(ctx context.Context, contentID int64, pageNum, pageSize int) ([]*model.Comment, error) {
	result := make([]*model.Comment, 0)
	err := d.GetDB().Table(d.getTableName(contentID)).
		Where("content_id = ?", contentID).
		Offset(pageNum * pageSize).
		Limit(pageSize).
		Order("create_time desc").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "获取评论失败"})
		return result, err
	}
	return result, nil
}

// GetCommentIDByContentID 根据contentID获得评论id
func (d *commentDao) GetCommentIDByContentID(ctx context.Context, contentID int64, pageNum, pageSize int) ([]int64, error) {
	result := make([]int64, 0)
	err := d.GetDB().Table(d.getTableName(contentID)).
		Select("id").
		Where("content_id = ?", contentID).
		Offset(pageNum * pageSize).
		Limit(pageSize).
		Order("create_time desc").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "获取评论失败"})
		return result, err
	}
	return result, nil
}

// GetCommentsByCommentIDs 根据commentID获得评论数据
func (d *commentDao) GetCommentsByCommentIDs(ctx context.Context, contentID int64, commentIDs []int64) ([]*model.Comment, error) {
	result := make([]*model.Comment, 0)
	err := d.GetDB().Table(d.getTableName(contentID)).
		Where("id in ?", commentIDs).
		Order("create_time desc").
		Find(&result).Error
	if err != nil {
		logger.Error(ctx, "GetCommentsByCommentIDs", logger.LogArgs{"err": err, "msg": "获取评论失败"})
		return result, err
	}
	return result, nil
}

// GetCommentCountByContentID 获得内容下的评论数量
func (d *commentDao) GetCommentCountByContentID(ctx context.Context, contentID int64) (int64, error) {
	result := int64(0)
	err := d.GetDB().Table(d.getTableName(contentID)).
		Where("id = ?", contentID).
		Count(&result).Error
	if err != nil {
		logger.Error(ctx, "GetCommentCountByContentID", logger.LogArgs{"err": err, "msg": "获取评论数量失败"})
		return result, err
	}
	return result, nil
}
