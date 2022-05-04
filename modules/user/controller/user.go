package controller

import (
	"MyServer/common"
	commonConsts "MyServer/consts"
	"MyServer/middleware/logger"
	"MyServer/modules/user/dto"
	"MyServer/modules/user/model"
	"MyServer/modules/user/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	*common.BaseController
}

// UserLogin 用户登录（未注册用户自动注册）
func (u *UserController) UserLogin(ctx *gin.Context) {
	var userData dto.UserInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, "UserLogin", logger.LogArgs{"msg": "UserLogin Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, "UserLogin", logger.LogArgs{"name": userData.LoginType, "nickName": userData.Password, "city": userData.City, "phone": userData.Phone})

	result := &dto.UserInfo{}
	needRegister := false
	userService := service.NewUserService()
	phoneService := service.NewUserPhoneService()

	var userInfo *model.UserInfo
	uid, err := phoneService.GetUIDByPhone(ctx, userData.Phone)
	if err != nil || uid == 0 {
		// 获得用户信息失败则认为无该用户，需要注册
		needRegister = true
	} else {
		userInfo, err = userService.GetUserInfoByUID(ctx, uid)
		if err != nil || userInfo == nil {
			// 获得用户信息失败则认为无该用户，需要注册
			needRegister = true
		}
	}

	if needRegister {
		userModelData, err := userService.ConvertUserModelData(ctx, &userData)
		if err != nil {
			logger.Error(ctx, "UserLogin", logger.LogArgs{"err": err, "msg": "用户数据转换失败"})
			u.EchoErrorStruct(ctx, common.ErrUserRegisterFailed)
			return
		}

		userInfo, err = userService.CreateUser(ctx, userModelData)
		if err != nil {
			logger.Error(ctx, "UserLogin", logger.LogArgs{"err": err, "msg": "创建新用户失败", "id": userData.ID, "uid": userData.UID, "loginType": userData.LoginType, "phone": userData.Phone})
			u.EchoErrorStruct(ctx, common.ErrUserRegisterFailed)
			return
		}
	}

	result.ID = userInfo.ID
	result.UID = userInfo.UID
	result.Phone = userInfo.Phone
	result.LoginType = int(userInfo.LoginType)
	result.NickName = userInfo.NickName
	result.Sex = int(userInfo.Sex)
	result.City = int(userInfo.City)
	result.Birthday = userInfo.Birthday.Format(commonConsts.TimeFormatData)
	result.ProfilePhoto = userInfo.ProfilePhoto
	result.Signature = userInfo.Signature
	result.RegisterTime = userInfo.RegisterTime.Format(commonConsts.TimeFormat)

	u.EchoSuccess(ctx, result)
	return
}

// GetUserInfoByUID 获得用户信息
func (u *UserController) GetUserInfoByUID(ctx *gin.Context) {
	var userData dto.UserInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, "GetUserInfoByUID", logger.LogArgs{"msg": "Get UserInfo Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, "GetUserInfoByUID", logger.LogArgs{"uid": userData.UID, "nickName": userData.NickName, "city": userData.City, "birthDay": userData.Birthday, "sex": userData.Sex, "signature": userData.Signature, "photo": userData.ProfilePhoto})

	userService := service.NewUserService()
	result, err := userService.GetUserInfoByUID(ctx, userData.UID)
	if err != nil {
		logger.Error(ctx, "GetUserInfoByUID", logger.LogArgs{"err": err, "msg": "查询用户信息失败", "uid": userData.UID})
		u.EchoErrorStruct(ctx, common.ErrGetUserInfoFailed)
		return
	}

	u.EchoSuccess(ctx, result)
	return
}

// UpdateUserInfo 更新用户信息
func (u *UserController) UpdateUserInfo(ctx *gin.Context) {
	var userData dto.UserInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, "GetUserInfo", logger.LogArgs{"msg": "Update UserInfo Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, "GetUserInfo", logger.LogArgs{"uid": userData.UID, "nickName": userData.NickName, "city": userData.City, "birtyDay": userData.Birthday, "sex": userData.Sex, "signature": userData.Signature, "photo": userData.ProfilePhoto})

	userService := service.NewUserService()
	result, err := userService.UpdateUserInfoByUID(ctx, &userData)
	if err != nil {
		logger.Error(ctx, "GetUserInfo", logger.LogArgs{"err": err, "msg": "更新用户信息失败", "uid": userData.UID})
		u.EchoErrorStruct(ctx, common.ErrUpdateUserInfoFailed)
		return
	}

	u.EchoSuccess(ctx, result)
	return
}
