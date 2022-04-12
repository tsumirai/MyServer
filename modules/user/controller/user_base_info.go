package controller

import (
	"MyServer/common"
	"MyServer/middleware/logger"
	"MyServer/modules/user/dto"
	"MyServer/modules/user/model"
	"MyServer/modules/user/service"

	"github.com/gin-gonic/gin"
)

type UserBaseController struct {
	*common.BaseController
}

func (u *UserBaseController) Ping(ctx *gin.Context) {
	u.EchoSuccess(ctx, "PONG")
}

// UserLogin 用户登录（未注册用户自动注册）
func (u *UserBaseController) UserLogin(ctx *gin.Context) {
	var userData model.UserBaseInfo
	err := ctx.BindJSON(&userData)
	if err != nil {
		logger.Error(ctx, logger.LogArgs{"msg": "UserLogin Failed", "err": err.Error()})
		u.EchoErrorStruct(ctx, common.ErrJSONUnmarshallFailed)
		return
	}

	logger.Info(ctx, logger.LogArgs{"name": userData.LoginType, "nickName": userData.Password, "city": userData.Phone})

	result := &dto.UserBaseInfo{}
	needRegister := false
	userBaseService := service.NewUserBaseService()

	userBaseInfo, err := userBaseService.GetUserBaseInfoByParam(ctx, &userData)
	if err != nil {
		// 获得用户信息失败则认为无该用户，需要注册
		needRegister = true
	} else {
		result.ID = userBaseInfo.ID
		result.UID = userBaseInfo.UID
		result.Phone = userBaseInfo.Phone
		result.RegisterTime = userBaseInfo.RegisterTime
	}

	if needRegister {
		userBaseInfo, err := userBaseService.CreateUser(ctx, &userData)
		if err != nil {
			logger.Error(ctx, logger.LogArgs{"err": err, "msg": "创建新用户失败", "id": userData.ID, "uid": userData.UID, "loginType": userData.LoginType, "phone": userData.Phone})
			u.EchoErrorStruct(ctx, common.ErrUserRegisterFailed)
			return
		}

		result.ID = userBaseInfo.ID
		result.UID = userBaseInfo.UID
		result.Phone = userBaseInfo.Phone
		result.RegisterTime = userBaseInfo.RegisterTime
	}

	u.EchoSuccess(ctx, result)
}
