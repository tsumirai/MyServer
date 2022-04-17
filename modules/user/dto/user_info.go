package dto

import "time"

type UserInfo struct {
	ID           int64     `json:"id"`            // 自增id
	UID          string    `json:"uid"`           // 用户uid
	Phone        string    `json:"phone"`         // 用户电话
	LoginType    int       `json:"login_type"`    // 用户登陆方式
	Password     string    `json:"password"`      // 用户密码
	NickName     string    `json:"nick_name"`     // 用户昵称
	Sex          int       `json:"sex"`           // 用户性别  0：男性  1：女性
	City         int       `json:"city"`          // 城市
	Birthday     string    `json:"birthday"`      // 用户生日
	ProfilePhoto string    `json:"profile_photo"` // 用户头像
	Signature    string    `json:"signature"`     // 个性签名
	RegisterTime time.Time `json:"register_time"` // 注册时间
	UpdateTime   time.Time `json:"update_time"`   // 更新时间
}
