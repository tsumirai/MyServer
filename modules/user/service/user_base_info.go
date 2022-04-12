package service

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/consts"
	"MyServer/modules/user/dao"
	"MyServer/modules/user/model"
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/snowflake"
)

type UserBaseService struct {
	common.BaseService
}

func NewUserBaseService() *UserBaseService {
	return &UserBaseService{}
}

// CreateUser 创建新用户
func (s *UserBaseService) CreateUser(ctx context.Context, param *model.UserBaseInfo) (*model.UserBaseInfo, error) {
	userDao := dao.NewUserBaseDao()
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

	userBaseInfo, err := userDao.CreateUser(ctx, param)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "创建用户失败", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
		return nil, err
	}

	return userBaseInfo, nil
}

// GetUserBaseInfoByParam 获得用户信息
func (s *UserBaseService) GetUserBaseInfoByParam(ctx context.Context, param *model.UserBaseInfo) (*model.UserBaseInfo, error) {
	userDao := dao.NewUserBaseDao()
	userBaseInfo, err := userDao.GetUserBaseInfoByParam(ctx, param)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "获得用户信息失败", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
		return nil, err
	}

	return userBaseInfo, nil
}
