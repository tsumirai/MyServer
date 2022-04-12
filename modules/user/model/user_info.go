package model

import "time"

type UserInfo struct {
	ID           int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                             //  自增id
	UID          string    `gorm:"column:uid" db:"uid" json:"uid" form:"uid"`                                         //  用户uid
	NickName     string    `gorm:"column:nick_name" db:"nick_name" json:"nick_name" form:"nick_name"`                 //  用户昵称
	Birthday     time.Time `gorm:"column:birthday" db:"birthday" json:"birthday" form:"birthday"`                     //  用户生日
	ProfilePhoto string    `gorm:"column:profile_photo" db:"profile_photo" json:"profile_photo" form:"profile_photo"` //  用户头像
	Sex          int64     `gorm:"column:sex" db:"sex" json:"sex" form:"sex"`                                         //  用户性别  0：男性  1：女性
	City         int       `gorm:"column:city" db:"city" json:"city" form:"city"`                                     // 城市
	Signature    string    `gorm:"column:signature" db:"signature" json:"signature" form:"signature"`                 //  个性签名
	UpdateTime   time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`         //  更新时间
}
