package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/dto"
	"MyServer/modules/user/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*common.BaseController
}

// GetUserInfo 获得用户信息
func (u *UserController) GetUserInfo(ctx *gin.Context) {
	var userData dto.UserInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"msg": "Get UserInfo Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, logger.LogArgs{"uid": userData.UID, "nickName": userData.NickName, "city": userData.City, "birtyDay": userData.Birthday, "sex": userData.Sex, "signature": userData.Signature, "photo": userData.ProfilePhoto})

	userService := service.NewUserService()
	result, err := userService.GetUserInfoByUID(ctx, userData.UID)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "查询用户信息失败", "uid": userData.UID})
		u.EchoErrorStruct(ctx, common.ErrGetUserInfoFailed)
		return
	}

	u.EchoSuccess(ctx, result)
}

// UpdateUserInfo 更新用户信息
func (u *UserController) UpdateUserInfo(ctx *gin.Context) {
	var userData dto.UserInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"msg": "Update UserInfo Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, logger.LogArgs{"uid": userData.UID, "nickName": userData.NickName, "city": userData.City, "birtyDay": userData.Birthday, "sex": userData.Sex, "signature": userData.Signature, "photo": userData.ProfilePhoto})

	userService := service.NewUserService()
	result, err := userService.UpdateUserInfoByUID(ctx, &userData)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"err": err, "msg": "更新用户信息失败", "uid": userData.UID})
		u.EchoErrorStruct(ctx, common.ErrUpdateUserInfoFailed)
		return
	}

	u.EchoSuccess(ctx, result)
}
