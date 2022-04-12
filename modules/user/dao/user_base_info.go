package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/model"
	"context"
	"fmt"
)

type UserBaseDao struct {
	common.BaseDao
}

const UserBaseInfoTable = "user_base_info"

func NewUserBaseDao() *UserBaseDao {
	return &UserBaseDao{}
}

// CreateUser 创建用户
func (d *UserBaseDao) CreateUser(ctx context.Context, param *model.UserBaseInfo) (*model.UserBaseInfo, error) {
	if param == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	db := d.GetDB().Table(UserBaseInfoTable)

	err := db.Create(param).Error
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "创建用户失败"})
		return nil, err
	}

	logger.Info(ctx, logger.LogArgs{"msg": "创建新用户", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
	return param, nil
}

// GetUserBaseInfoByParam 根据参数获得用户信息
func (d *UserBaseDao) GetUserBaseInfoByParam(ctx context.Context, param *model.UserBaseInfo) (*model.UserBaseInfo, error) {
	if param == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	result := &model.UserBaseInfo{}
	db := d.GetDB().Table(UserBaseInfoTable)
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

	db = db.Find(&result)
	if err := db.Error; err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err.Error, "msg": "查询用户信息失败", "id": param.ID, "uid": param.UID, "phone": param.Phone, "loginType": param.LoginType})
		return nil, err
	}
	return result, nil
}
