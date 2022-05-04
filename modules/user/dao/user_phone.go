package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/consts"
	"MyServer/modules/user/model"
	"context"
	"strconv"
)

type UserPhone struct {
	common.BaseDao
}

func NewUserPhoneDao() *UserPhone {
	return &UserPhone{}
}

func (d *UserPhone) getTableName(ctx context.Context, phone int64) string {
	return consts.UserPhoneTable + "_" + strconv.FormatInt(phone%consts.UserPhoneTableNum, 10)
}

// CreateUserPhone 创建用户手机号数据
func (d *UserPhone) CreateUserPhone(ctx context.Context, data *model.UserPhone) (*model.UserPhone, error) {
	err := d.GetDB().Table(d.getTableName(ctx, data.UID)).Create(data).Error
	if err != nil {
		logger.Error(ctx, "CreateUserPhone", logger.LogArgs{"err": err})
		return nil, err
	}
	return data, nil
}

// GetUIDByPhone 根据手机号获得用户uid
func (d *UserPhone) GetUIDByPhone(ctx context.Context, phone string) (int64, error) {
	phoneNum, err := strconv.ParseInt(phone, 10, 64)
	if err != nil {
		logger.Error(ctx, "GetUIDByPhone", logger.LogArgs{"err": err})
		return 0, err
	}

	var uid int64

	err = d.GetDB().Table(d.getTableName(ctx, phoneNum)).Select("uid").Where("phone = ?", phone).Find(&uid).Error
	if err != nil {
		logger.Error(ctx, "GetUIDByPhone", logger.LogArgs{"err": err})
		return 0, err
	}

	return uid, nil
}
