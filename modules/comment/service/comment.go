package service

import (
	"MyServer/cache"
	"MyServer/common"
	commonConsts "MyServer/consts"
	"MyServer/middleware/logger"
	"MyServer/modules/comment/consts"
	"MyServer/modules/comment/dao"
	"MyServer/modules/comment/dto"
	"MyServer/modules/comment/model"
	"MyServer/util"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type commentService struct {
	common.BaseService
}

func NewCommentService() *commentService {
	return &commentService{}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(ctx context.Context, param *dto.CreateCommentReq) error {
	if param.ContentID == 0 || param.AuthorUID == 0 {
		err := fmt.Errorf("参数非法")
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "参数非法"})
		return err
	}

	// 评论格式化
	param.Comment = strings.TrimSpace(param.Comment)
	if param.Comment == "" {
		err := fmt.Errorf("内容不能为空")
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "内容不能为空"})
		return err
	}

	// 检查评论长度
	if len([]rune(param.Comment)) > consts.CommentLengthLimit {
		err := fmt.Errorf("评论太长")
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "评论太长"})
		return err
	}

	// todo 检查用户是否被封禁

	// todo 审核，另开goroutine异步进行

	commentDao := dao.NewCommentDao()
	commentData := s.convertCommentToModel(ctx, param)
	if commentData == nil {
		return nil
	}

	commentData.CommentStatus = consts.CommentStatusAuditing
	err := commentDao.CreateComment(ctx, commentData)
	if err != nil {
		logger.Error(ctx, "CreateComment", logger.LogArgs{"err": err, "msg": "创建评论失败"})
		return err
	}

	return nil
}

// GetCommentsByContentID 根据内容ID获得评论
func (s *commentService) GetCommentsByContentID(ctx context.Context, contentID int64, pageNum, pageSize int) ([]*dto.CommentRes, error) {
	result := make([]*dto.CommentRes, 0)

	if contentID == 0 {
		err := fmt.Errorf("参数错误")
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "contentID": contentID})
		return result, err
	}

	// 默认每页拉取10条评论
	if pageSize == 0 {
		pageSize = consts.DefaultCommentPageSize
	}

	// 获得内容下的评论id列表
	cacheSvr := cache.NewCache()
	cacheSvr.RegisterCallbackFunc(s.GetCommentIDByContentIDCallback)

	commentIDBytes, err := cacheSvr.GetValueFromCache(ctx, cache.GetCommentIDByContentID(contentID, pageNum, pageSize), commonConsts.FiveMinute)
	if err != nil {
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "获得评论失败", "contentID": contentID, "pageNum": pageNum, "pageSize": pageSize})
		return result, err
	}

	commentID := make([]int64, 0)
	err = json.Unmarshal(commentIDBytes, &commentID)
	if err != nil {
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "反序列化失败"})
		return result, err
	}

	commentIDStr := util.ConvertInt64SliceToStringSlice(commentID)

	// 根据评论id列表获得评论内容
	cacheSvr.RegisterMultiCallbackFunc(s.GetCommentDataByIDCallback)
	commentByte, err := cacheSvr.GetValuesFromHashCache(ctx, cache.GetCommentDataByIDsRedisKey(contentID), commonConsts.FiveMinute, commentIDStr...)
	if err != nil {
		logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "获得评论失败"})
		return result, err
	}

	for _, v := range commentByte {
		var comment *model.Comment
		err = json.Unmarshal(v, &comment)
		if err != nil {
			logger.Error(ctx, "GetCommentsByContentID", logger.LogArgs{"err": err, "msg": "反序列化失败"})
			return result, err
		}

		result = append(result, s.convertCommentToDto(ctx, comment))
	}

	return result, nil
}

// GetCommentCountByContentID 获得内容下所有评论的数量
func (s *commentService) GetCommentCountByContentID(ctx context.Context, contentID int64) (int64, error) {
	var err error
	if contentID <= 0 {
		err = fmt.Errorf("参数错误")
		logger.Error(ctx, "GetCommentCountByContentID", logger.LogArgs{"err": err, "contentID": contentID})
		return 0, err
	}

	commentDao := dao.NewCommentDao()
	result, err := commentDao.GetCommentCountByContentID(ctx, contentID)
	if err != nil {
		logger.Error(ctx, "GetCommentCountByContentID", logger.LogArgs{"err": err, "cotentID": contentID})
		return 0, err
	}

	return result, nil
}

/*==================================================================================*/
// convertCommentToModel 将请求的结构体转为数据库结构
func (s *commentService) convertCommentToModel(ctx context.Context, param *dto.CreateCommentReq) *model.Comment {
	if param == nil {
		return nil
	}

	result := &model.Comment{
		ContentID:       param.ContentID,
		ParentCommentID: param.ParentCommentID,
		AuthorUID:       param.AuthorUID,
		Comment:         param.Comment,
	}

	return result
}

// convertCommentToDto 将评论的数据库结构转为前端结构
func (s *commentService) convertCommentToDto(ctx context.Context, param *model.Comment) *dto.CommentRes {
	if param == nil {
		return nil
	}

	result := &dto.CommentRes{
		ID:              param.ID,
		ContentID:       param.ContentID,
		ParentCommentID: param.ParentCommentID,
		AuthorUID:       param.AuthorUID,
		Comment:         param.Comment,
		CommentStatus:   param.CommentStatus,
		AuditReason:     param.AuditReason,
		CreateTime:      param.CreateTime,
		UpdateTime:      param.UpdateTime,
	}
	return result
}
