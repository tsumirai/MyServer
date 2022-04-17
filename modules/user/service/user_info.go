package service

import (
	"MyServer/common"
	commonConsts "MyServer/consts"
	"MyServer/middleware/logger"
	"MyServer/modules/user/consts"
	"MyServer/modules/user/dao"
	"MyServer/modules/user/dto"
	"MyServer/modules/user/model"
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"time"
)

type UserService struct {
	common.BaseService
}

func NewUserService() *UserService {
	return &UserService{}
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(ctx context.Context, param *model.UserInfo) (*model.UserInfo, error) {
	userDao := dao.NewUserDao()
	if param.Password == "" {
		err := fmt.Errorf("用户密码不能为空")
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "用户密码不能为空", "passWord": param.Password})
		return nil, err
	}

	if param.LoginType == consts.LoginTypePhone && param.Phone == "" {
		err := fmt.Errorf("用户手机号不能为空")
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "用户手机号不能为空", "phone": param.Phone})
		return nil, err
	}

	// 生成uid
	node, err := snowflake.NewNode(1)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "创建用户失败"})
		return nil, err
	}

	uid := node.Generate().String()
	param.UID = uid

	param.RegisterTime = time.Now()
	param.UpdateTime = time.Now()

	userInfo, err := userDao.CreateUser(ctx, param)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "创建用户失败", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
		return nil, err
	}

	return userInfo, nil
}

// GetUserInfoByParam 获得用户信息
func (s *UserService) GetUserInfoByParam(ctx context.Context, param *model.UserInfo) (*model.UserInfo, error) {
	userDao := dao.NewUserDao()
	userInfo, err := userDao.GetUserInfoByParam(ctx, param)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "获得用户信息失败", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
		return nil, err
	}

	return userInfo, nil
}

// GetUserInfoByUID 根据UID查询用户信息
func (s *UserService) GetUserInfoByUID(ctx context.Context, uid string) (*model.UserInfo, error) {
	userDao := dao.NewUserDao()
	userInfo, err := userDao.GetUserInfoByParam(ctx, &model.UserInfo{UID: uid})
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
		Phone:        userInfo.Phone,
		Password:     userInfo.Password,
		LoginType:    int64(userInfo.LoginType),
		NickName:     userInfo.NickName,
		Birthday:     birthDay,
		ProfilePhoto: userInfo.ProfilePhoto,
		Sex:          int64(userInfo.Sex),
		City:         int64(userInfo.City),
		Signature:    userInfo.Signature,
	}

	result, err := s.GetUserInfoByUID(ctx, userInfo.UID)
	if err != nil || result == nil {
		// 查询失败或者userInfo为空，认为不存在，需要插入
		result, err := userDao.CreateUser(ctx, newUserInfo)
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

// ConvertUserModelData 转换用户数据
func (s *UserService) ConvertUserModelData(ctx context.Context, userDtoData *dto.UserInfo) (*model.UserInfo, error) {
	birthDay, err := time.Parse(commonConsts.TimeFormatData, userDtoData.Birthday)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"msg": "用户数据转换失败", "err": err.Error()})
		return nil, err
	}
	result := &model.UserInfo{
		ID:           userDtoData.ID,
		UID:          userDtoData.UID,
		Phone:        userDtoData.Phone,
		Password:     userDtoData.Password,
		LoginType:    int64(userDtoData.LoginType),
		NickName:     userDtoData.NickName,
		Sex:          int64(userDtoData.Sex),
		City:         int64(userDtoData.City),
		Birthday:     birthDay,
		ProfilePhoto: userDtoData.ProfilePhoto,
		Signature:    userDtoData.Signature,
		RegisterTime: userDtoData.RegisterTime,
		UpdateTime:   userDtoData.UpdateTime,
	}
	return result, nil
}
