package dto

import "time"

type UserInfo struct {
	ID           int64     `json:"id"`            //  自增id
	UID          string    `json:"uid"`           //  用户uid
	NickName     string    `json:"nick_name"`     //  用户昵称
	Birthday     string    `json:"birthday"`      //  用户生日
	ProfilePhoto string    `json:"profile_photo"` //  用户头像
	Sex          int64     `json:"sex"`           //  用户性别  0：男性  1：女性
	City         int       `json:"city"`          // 城市
	Signature    string    `json:"signature"`     //  个性签名
	UpdateTime   time.Time `json:"update_time"`   //  更新时间
}
