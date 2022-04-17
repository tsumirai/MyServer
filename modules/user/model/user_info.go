package model

import "time"

type UserInfo struct {
	ID           int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                                                       //  自增id
	UID          string    `gorm:"column:uid" db:"uid" json:"uid" form:"uid"`                                                                   //  用户uid
	Phone        string    `gorm:"column:phone" db:"phone" json:"phone" form:"phone"`                                                           //  用户电话号码
	Password     string    `gorm:"column:password" db:"password" json:"password" form:"password"`                                               //  用户密码
	LoginType    int64     `gorm:"column:login_type" db:"login_type" json:"login_type" form:"login_type"`                                       //  登录方式
	NickName     string    `gorm:"column:nick_name" db:"nick_name" json:"nick_name" form:"nick_name"`                                           //  用户昵称
	Sex          int64     `gorm:"column:sex" db:"sex" json:"sex" form:"sex"`                                                                   //  用户性别  0：男性  1：女性
	City         int64     `gorm:"column:city" db:"city" json:"city" form:"city"`                                                               //  城市
	Birthday     time.Time `gorm:"column:birthday" db:"birthday" json:"birthday" form:"birthday"`                                               //  用户生日
	ProfilePhoto string    `gorm:"column:profile_photo" db:"profile_photo" json:"profile_photo" form:"profile_photo"`                           //  用户头像
	Signature    string    `gorm:"column:signature" db:"signature" json:"signature" form:"signature"`                                           //  个性签名
	RegisterTime time.Time `gorm:"column:register_time default:CURRENT_TIMESTAMP" db:"register_time" json:"register_time" form:"register_time"` //  注册时间
}
