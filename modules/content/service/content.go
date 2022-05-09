package service

import (
	"MyServer/cache"
	"MyServer/common"
	commonConsts "MyServer/consts"
	"MyServer/middleware/logger"
	"MyServer/modules/content/consts"
	"MyServer/modules/content/dao"
	"MyServer/modules/content/dto"
	"MyServer/modules/content/model"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type contentService struct {
	common.BaseService
}

func NewContentService() *contentService {
	return &contentService{}
}

// CreateContent 创建内容
func (s *contentService) CreateContent(ctx context.Context, param *dto.CreateContentReq) error {
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

	param.Title = strings.TrimSpace(param.Title)
	param.Content = strings.TrimSpace(param.Content)
	if param.Title == "" || param.Content == "" {
		err := fmt.Errorf("标题或内容不能为空")
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	if len([]rune(param.Title)) > consts.ContentTitleLengthLimit {
		err := fmt.Errorf("标题太长")
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	if len([]rune(param.Content)) > consts.ContentLengthLimit {
		err := fmt.Errorf("内容太长")
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err})
		return err
	}

	// todo 检查用户是否被封禁

	contentDao := dao.NewContentDao()
	createContent := s.convertToContentModel(ctx, param)
	if createContent == nil {
		return nil
	}
	createContent.ContentStatus = consts.ContentStatusAuditing
	err := contentDao.CreateContent(ctx, createContent)
	if err != nil {
		logger.Error(ctx, "CreateContent", logger.LogArgs{"err": err, "msg": "创建内容失败"})
		return err
	}

	// todo 内容审核 单独开goroutine异步审核

	return nil
}

// GetContentByID 根据内容id获得内容数据
func (s *contentService) GetContentByID(ctx context.Context, ID, authorUID int64) (*dto.ContentRes, error) {
	cacheSvr := cache.NewCache()
	cacheSvr.RegisterMultiCallbackFunc(s.getContentDataByIDsCallback)

	var contentData *model.Content
	contentByteData, err := cacheSvr.GetValuesFromHashCache(ctx, cache.GetContentDataByIDsRedisKey(authorUID), commonConsts.FiveMinute, strconv.FormatInt(ID, 64))
	if err != nil {
		logger.Error(ctx, "GetContentByID", logger.LogArgs{"err": err, "msg": "获得内容数据失败"})
		return nil, err
	}

	if contentByte, exit := contentByteData[strconv.FormatInt(ID, 64)]; exit {
		if err = json.Unmarshal(contentByte, &contentData); err != nil {
			logger.Error(ctx, "GetContentByID", logger.LogArgs{"err": err, "msg": "反序列化内容数据失败"})
			return nil, err
		}
	}

	content := s.convertToContentDto(ctx, contentData)
	return content, nil
}

// DeleteContentByID 根据内容ID删除内容
func (s *contentService) DeleteContentByID(ctx context.Context, contentID, authorUID int64) error {
	if contentID == 0 || authorUID == 0 {
		err := fmt.Errorf("参数错误")
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err, "contentID": contentID, "authorUID": authorUID})
		return err
	}

	// 检查是否是作者本人
	contentData, err := s.GetContentByID(ctx, contentID, authorUID)
	if err != nil {
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err, "msg": "查询内容失败"})
		return err
	}

	if contentData == nil {
		err := fmt.Errorf("查询内容失败")
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err, "msg": "查询内容失败"})
		return err
	}

	if contentData.AuthorUID != authorUID {
		err := fmt.Errorf("只能删除本人发布的内容")
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err, "msg": "只能删除本人发布的内容"})
		return err
	}

	// 设置内容为用户删除
	err = s.setContentStatus(ctx, contentID, authorUID, consts.ContentStatusUserDelete)
	if err != nil {
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err, "msg": "删除内容失败"})
		return err
	}

	// 删除缓存
	cacheSvr := cache.NewCache()
	err = cacheSvr.HDel(cache.GetContentDataByIDsRedisKey(authorUID), strconv.FormatInt(contentID, 64))
	if err != nil {
		logger.Error(ctx, "DeleteContentByID", logger.LogArgs{"err": err, "msg": "删除内容缓存失败"})
	}

	return nil
}

// GetContentList 获得内容
func (s *contentService) GetContentList(ctx context.Context, UID int64, pageNo, pageSize int, contentListType int) ([]*dto.ContentRes, error) {
	switch contentListType {
	case consts.ContentListTypeLike:
		// 获得喜欢的内容列表
		return s.getLikeContentList(ctx, UID, pageNo, pageSize)
	case consts.ContentListTypeCollect:
		// 获得收藏的内容列表
		return s.getCollectContentList(ctx, UID, pageNo, pageSize)
	default:
		// 默认为拉取该UID发布的内容列表
		return s.getContentIDsByAuthorUID(ctx, UID, pageNo, pageSize)
	}
}

// SetContentPermission 设置内容的权限
func (s *contentService) SetContentPermission(ctx context.Context, contentID, authorUID int64, permission int) error {
	if contentID == 0 || authorUID == 0 {
		err := fmt.Errorf("参数非法")
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "contentID": contentID, "authorUID": authorUID})
		return err
	}

	if permission < consts.ContentPermissionNone || permission > consts.ContentPermissionSelf {
		err := fmt.Errorf("参数非法")
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "permission": permission})
		return err
	}

	contentDao := dao.NewContentDao()
	err := contentDao.SetContentPermission(ctx, contentID, authorUID, permission)
	if err != nil {
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "contentID": contentID, "authorUID": authorUID, "permission": permission})
		return err
	}

	// 修改缓存
	newContent, err := contentDao.GetContentByID(ctx, contentID, authorUID)
	if err != nil {
		logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "contentID": contentID, "authorUID": authorUID, "permission": permission})
	} else {
		err = contentDao.UpdateContentCache(ctx, newContent)
		if err != nil {
			logger.Error(ctx, "SetContentPermission", logger.LogArgs{"err": err, "contentID": contentID, "authorUID": authorUID, "permission": permission})
		}
	}

	return nil
}

/*=========================================================================================*/
// setContentStatus 设置内容的状态
func (s *contentService) setContentStatus(ctx context.Context, contentID, authorUID int64, contentStatus int) error {
	contentDao := dao.NewContentDao()
	return contentDao.SetContentStatus(ctx, contentID, authorUID, contentStatus)
}

// getContentIDsByAuthorUID 根据作者UID获得内容列表
func (s *contentService) getContentIDsByAuthorUID(ctx context.Context, authorUID int64, pageNum, pageSize int) ([]*dto.ContentRes, error) {
	result := make([]*dto.ContentRes, 0)
	cacheSvr := cache.NewCache()
	cacheSvr.RegisterCallbackFunc(s.getContentIDsByAuthorUIDCallback)

	// 默认每页拉取10条内容
	if pageSize == 0 {
		pageSize = consts.DefaultContentSize
	}

	// 先获得contentID的列表
	contentIDBytes, err := cacheSvr.GetValueFromCache(ctx, cache.GetContentIDsByAuthorUID(authorUID, pageNum, pageSize), commonConsts.FiveMinute)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUID", logger.LogArgs{"err": err, "msg": "获得内容列表失败"})
		return result, err
	}

	tempResult := make([]int64, 0)
	err = json.Unmarshal(contentIDBytes, &tempResult)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUID", logger.LogArgs{"err": err, "msg": "反序列化内容列表失败"})
		return result, err
	}

	// 获得内容的具体数据
	contentID := make([]string, 0, len(tempResult))
	for _, v := range tempResult {
		contentID = append(contentID, strconv.FormatInt(v, 64))
	}
	cacheSvr.RegisterMultiCallbackFunc(s.getContentDataByIDsCallback)
	contentByteData, err := cacheSvr.GetValuesFromHashCache(ctx, cache.GetContentDataByIDsRedisKey(authorUID), commonConsts.FiveMinute, contentID...)
	if err != nil {
		logger.Error(ctx, "getContentIDsByAuthorUID", logger.LogArgs{"err": err, "msg": "获取内容数据失败"})
		return result, err
	}

	for _, v := range contentByteData {
		var tempContent *model.Content
		err = json.Unmarshal(v, &tempContent)
		if err != nil {
			logger.Error(ctx, "getContentIDsByAuthorUID", logger.LogArgs{"err": err, "msg": "反序列化内容列表失败"})
			return result, err
		}
		result = append(result, s.convertToContentDto(ctx, tempContent))
	}

	return result, nil
}

// todo getLikeContentList 获得喜欢的内容列表
func (s *contentService) getLikeContentList(ctx context.Context, UID int64, pageNo, pageSize int) ([]*dto.ContentRes, error) {
	result := make([]*dto.ContentRes, 0)

	return result, nil
}

// todo getCollectContentList 获得收藏的内容列表
func (s *contentService) getCollectContentList(ctx context.Context, UID int64, pageNo, pageSize int) ([]*dto.ContentRes, error) {
	result := make([]*dto.ContentRes, 0)

	return result, nil
}

// convertToContentModel 转换content的请求参数为mysql格式
func (s *contentService) convertToContentModel(ctx context.Context, param *dto.CreateContentReq) *model.Content {
	if param == nil {
		return nil
	}
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

// convertToContentDto 将content从mysql的格式转为dto的格式
func (s *contentService) convertToContentDto(ctx context.Context, param *model.Content) *dto.ContentRes {
	if param == nil {
		return nil
	}
	result := &dto.ContentRes{
		ID:               param.ID,
		Title:            param.Title,
		ImageUrls:        make([]string, 0),
		Content:          param.Content,
		AuthorUID:        param.AuthorUID,
		VideoUrl:         param.VideoUrl,
		ContentType:      param.ContentType,
		ContentSpaceType: param.ContentSpaceType,
		LocCityID:        param.LocCityID,
		ContentStatus:    param.ContentStatus,
		AuditReason:      param.AuditReason,
		CreateTime:       param.CreateTime,
		UpdateTime:       param.UpdateTime,
	}
	imageUrls := strings.Split(param.ImageUrls, ";")
	result.ImageUrls = append(result.ImageUrls, imageUrls...)
	return result
}
