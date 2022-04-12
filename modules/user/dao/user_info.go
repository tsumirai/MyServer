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

const UserInfoTable = "user_info"

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (d *UserDao) CreateUserInfo(ctx context.Context, userInfo *model.UserInfo) (*model.UserInfo, error) {
	if userInfo == nil {
		err := fmt.Errorf("参数不能为空！")
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "参数不能为空"})
		return nil, err
	}

	db := d.GetDB().Table(UserInfoTable)
	err := db.Create(userInfo).Error
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "插入用户信息失败"})
		return nil, err
	}

	return userInfo, nil
}

// UpdateUserInfoByUID 根据uid更新用户信息
func (d *UserDao) UpdateUserInfoByUID(ctx context.Context, userInfo *model.UserInfo) (*model.UserInfo, error) {
	if userInfo == nil {
		err := fmt.Errorf("参数不能为空")
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "参数不能为空"})
		return nil, err
	}

	if userInfo.UID == "" {
		err := fmt.Errorf("UID不能为空")
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "UID不能为空"})
		return nil, err
	}

	db := d.GetDB().Table(UserInfoTable)
	err := db.Where("uid = ?", userInfo.UID).Update(userInfo).Error
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "更新用户信息失败"})
		return nil, err
	}

	return userInfo, nil
}

// GetUserInfoByUID 根据uid查询用户信息
func (d *UserDao) GetUserInfoByUID(ctx context.Context, uid string) (*model.UserInfo, error) {
	if uid == "" {
		err := fmt.Errorf("uid为空！！")
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "uid为空"})
		return nil, err
	}

	db := d.GetDB().Table(UserInfoTable)
	db.Where("uid = ?", uid)
	result := &model.UserInfo{}
	db = db.Find(&result)
	if err := db.Error; err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "查询用户信息失败", "uid": uid})
		return nil, err
	}

	return result, nil
}
