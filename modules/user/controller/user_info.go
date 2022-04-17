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

func (u *UserController) Ping(ctx *gin.Context) {
	u.EchoSuccess(ctx, "Pong")
}

// UserLogin 用户登录（未注册用户自动注册）
func (u *UserController) UserLogin(ctx *gin.Context) {
	var userData dto.UserInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"msg": "UserLogin Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, logger.LogArgs{"name": userData.LoginType, "nickName": userData.Password, "city": userData.Phone})

	result := &dto.UserInfo{}
	needRegister := false
	userService := service.NewUserService()

	userInfo, err := userService.GetUserInfoByParam(ctx, &model.UserInfo{
		Phone:     userData.Phone,
		LoginType: int64(userData.LoginType),
	})
	if err != nil {
		// 获得用户信息失败则认为无该用户，需要注册
		needRegister = true
	}

	if needRegister {
		userModelData, err := userService.ConvertUserModelData(ctx, &userData)
		if err != nil {
			logger.Error(ctx, logger.LogArgs{"err": err, "msg": "用户数据转换失败"})
			u.EchoErrorStruct(ctx, common.ErrUserRegisterFailed)
			return
		}

		userInfo, err = userService.CreateUser(ctx, userModelData)
		if err != nil {
			logger.Error(ctx, logger.LogArgs{"err": err, "msg": "创建新用户失败", "id": userData.ID, "uid": userData.UID, "loginType": userData.LoginType, "phone": userData.Phone})
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

	u.EchoSuccess(ctx, result)
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