package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/model"
	"context"
	"fmt"
)

type UserDao struct {
	common.BaseDao
}

const UserTable = "user"

func NewUserDao() *UserDao {
	return &UserDao{}
}

// CreateUser 创建用户
func (d *UserDao) CreateUser(ctx context.Context, param *model.UserInfo) (*model.UserInfo, error) {
	if param == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	db := d.GetDB().Table(UserTable)

	err := db.Create(param).Error
	if err != nil {
		logger.Error(ctx, "CreateUser", logger.LogArgs{"err": err, "msg": "创建用户失败"})
		return nil, err
	}

	logger.Info(ctx, "CreateUser", logger.LogArgs{"msg": "创建新用户", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
	return param, nil
}

// GetUserInfoByParam 根据参数获得用户信息
func (d *UserDao) GetUserInfoByParam(ctx context.Context, param *model.UserInfo) (*model.UserInfo, error) {
	if param == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	result := &model.UserInfo{}
	db := d.GetDB().Table(UserTable)
	if param.ID != 0 {
		db.Where("id = ?", param.ID)
	}

	if param.UID != "" {
		db.Where("uid = ?", param.UID)
	}

	if param.LoginType != 0 {
		db.Where("login_type = ?", param.LoginType)
	}

	if param.Phone != "" {
		db.Where("phone = ?", param.Phone)
	}

	if param.NickName != "" {
		db.Where("nick_name is like ?", param.NickName)
	}

	db = db.Take(&result)
	if err := db.Error; err != nil {
		logger.Error(ctx, "GetUserInfoByParam", logger.LogArgs{"err": err.Error, "msg": "查询用户信息失败", "id": param.ID, "uid": param.UID, "phone": param.Phone, "loginType": param.LoginType})
		return nil, err
	}

	return result, nil
}

// UpdateUserInfoByUID 根据uid更新用户信息
func (d *UserDao) UpdateUserInfoByUID(ctx context.Context, userInfo *model.UserInfo) (*model.UserInfo, error) {
	if userInfo == nil {
		err := fmt.Errorf("参数不能为空")
		logger.Error(ctx, "UpdateUserInfoByUID", logger.LogArgs{"err": err, "msg": "参数不能为空"})
		return nil, err
	}

	if userInfo.UID == "" {
		err := fmt.Errorf("UID不能为空")
		logger.Error(ctx, "UpdateUserInfoByUID", logger.LogArgs{"err": err, "msg": "UID不能为空"})
		return nil, err
	}

	db := d.GetDB().Table(UserTable)
	err := db.Where("uid = ?", userInfo.UID).Updates(userInfo).Error
	if err != nil {
		logger.Error(ctx, "UpdateUserInfoByUID", logger.LogArgs{"err": err, "msg": "更新用户信息失败"})
		return nil, err
	}

	return userInfo, nil
}
