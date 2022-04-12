package model

import "time"

type UserBaseInfo struct {
	ID           int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                             //  用户id
	UID          string    `gorm:"column:uid" db:"uid" json:"uid" form:"uid"`                                         //  用户uid
	Password     string    `gorm:"column:password" db:"password" json:"password" form:"password"`                     //  用户密码
	Phone        string    `gorm:"column:phone" db:"phone" json:"phone" form:"phone"`                                 //  用户手机号码
	LoginType    int64     `gorm:"column:login_type" db:"login_type" json:"login_type" form:"login_type"`             //  登录方式
	RegisterTime time.Time `gorm:"column:register_time" db:"register_time" json:"register_time" form:"register_time"` //  注册时间
	UpdateTime   time.Time `gorm:"column:update_time" db:"update_time" json:"update_time" form:"update_time"`         //  更新时间
}
