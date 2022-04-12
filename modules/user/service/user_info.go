package service

import (
	"MyServer/common"
	commonConsts "MyServer/consts"
	"MyServer/middleware/logger"
	"MyServer/modules/user/dao"
	"MyServer/modules/user/dto"
	"MyServer/modules/user/model"
	"context"
	"time"
)

type UserService struct {
	common.BaseService
}

func NewUserService() *UserService {
	return &UserService{}
}

// GetUserInfoByUID 根据UID查询用户信息
func (s *UserService) GetUserInfoByUID(ctx context.Context, uid string) (*model.UserInfo, error) {
	userDao := dao.NewUserDao()
	userInfo, err := userDao.GetUserInfoByUID(ctx, uid)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "查询用户信息失败", "uid": uid})
		return nil, err
	}

	return userInfo, nil
}

// UpdateUserInfoByUID 根据UID更新用户信息
func (s *UserService) UpdateUserInfoByUID(ctx context.Context, userInfo *dto.UserInfo) (*model.UserInfo, error) {
	userDao := dao.NewUserDao()

	birthDay, err := time.Parse(commonConsts.TimeFormatData, userInfo.Birthday)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "解析生日失败", "birthDay": userInfo.Birthday})
		return nil, err
	}
	newUserInfo := &model.UserInfo{
		UID:          userInfo.UID,
		NickName:     userInfo.NickName,
		Birthday:     birthDay,
		ProfilePhoto: userInfo.ProfilePhoto,
		Sex:          userInfo.Sex,
		City:         userInfo.City,
		Signature:    userInfo.Signature,
		UpdateTime:   time.Now(),
	}

	result, err := s.GetUserInfoByUID(ctx, userInfo.UID)
	if err != nil || result == nil {
		// 查询失败或者userInfo为空，认为不存在，需要插入
		result, err := userDao.CreateUserInfo(ctx, newUserInfo)
		if err != nil {
			logger.Error(ctx, logger.LogArgs{"err": err, "msg": "插入用户信息失败"})
			return nil, err
		}
		return result, nil
	}

	result, err = userDao.UpdateUserInfoByUID(ctx, newUserInfo)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "更新用户信息失败"})
		return nil, err
	}

	return result, nil
}
