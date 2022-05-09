package dao

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/consts"
	"MyServer/modules/user/model"
	"context"
	"fmt"
	"strconv"
)

type userDao struct {
	common.BaseDao
}

func NewUserDao() *userDao {
	return &userDao{}
}

func (d *userDao) getTableName(ctx context.Context, uid int64) string {
	return consts.UserTable + "_" + strconv.FormatInt(uid%consts.UserTableNum, 10)
}

// CreateUser 创建用户
func (d *userDao) CreateUser(ctx context.Context, param *model.UserInfo) (*model.UserInfo, error) {
	if param == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	if param.UID == 0 {
		return nil, fmt.Errorf("UID数不能为0")
	}

	tx := d.GetDB().Begin()
	err := tx.Table(d.getTableName(ctx, param.UID)).
		Create(param).Error
	if err != nil {
		logger.Error(ctx, "CreateUser", logger.LogArgs{"err": err, "msg": "创建用户失败"})
		tx.Rollback()
		return nil, err
	}

	phoneDao := NewUserPhoneDao()
	err = tx.Table(phoneDao.getTableName(ctx, param.UID)).
		Create(&model.UserPhone{
			Phone: param.Phone,
			UID:   param.UID,
		}).Error
	if err != nil {
		logger.Error(ctx, "CreateUser", logger.LogArgs{"err": err, "msg": "创建用户失败"})
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	logger.Info(ctx, "CreateUser", logger.LogArgs{"msg": "创建新用户", "id": param.ID, "uid": param.UID, "loginType": param.LoginType, "phone": param.Phone})
	return param, nil
}

// GetUserInfoByParam 根据参数获得用户信息
func (d *userDao) GetUserInfoByParam(ctx context.Context, param *model.UserInfo) (*model.UserInfo, error) {
	if param == nil {
		return nil, fmt.Errorf("参数不能为空")
	}

	if param.UID == 0 {
		return nil, fmt.Errorf("UID数不能为0")
	}

	result := &model.UserInfo{}
	db := d.GetDB().Table(d.getTableName(ctx, param.UID))
	if param.ID != 0 {
		db.Where("id = ?", param.ID)
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
func (d *userDao) UpdateUserInfoByUID(ctx context.Context, userInfo *model.UserInfo) (*model.UserInfo, error) {
	if userInfo == nil {
		err := fmt.Errorf("参数不能为空")
		logger.Error(ctx, "UpdateUserInfoByUID", logger.LogArgs{"err": err, "msg": "参数不能为空"})
		return nil, err
	}

	if userInfo.UID == 0 {
		err := fmt.Errorf("UID不能为空")
		logger.Error(ctx, "UpdateUserInfoByUID", logger.LogArgs{"err": err, "msg": "UID不能为空"})
		return nil, err
	}

	db := d.GetDB().Table(d.getTableName(ctx, userInfo.UID))
	err := db.Where("uid = ?", userInfo.UID).Updates(userInfo).Error
	if err != nil {
		logger.Error(ctx, "UpdateUserInfoByUID", logger.LogArgs{"err": err, "msg": "更新用户信息失败"})
		return nil, err
	}

	return userInfo, nil
}
