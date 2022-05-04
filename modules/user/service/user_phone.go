package service

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/dao"
	"context"
)

type UserPhoneService struct {
	common.BaseService
}

func NewUserPhoneService() *UserPhoneService {
	return &UserPhoneService{}
}

// GetUIDByPhone 根据手机号获得UID
func (s *UserPhoneService) GetUIDByPhone(ctx context.Context, phone string) (int64, error) {
	phoneDao := dao.NewUserPhone()
	uid, err := phoneDao.GetUIDByPhone(ctx, phone)
	if err != nil {
		logger.Error(ctx, "GetUIDByPhone", logger.LogArgs{"err": err})
		return uid, err
	}
	return uid, nil
}
