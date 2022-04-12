package dto

import "time"

type UserBaseInfo struct {
	ID           int64     `gorm:"column:id" db:"id" json:"id" form:"id"`                                             //  用户id
	UID          string    `gorm:"column:uid" db:"uid" json:"uid" form:"uid"`                                         //  用户uid
	Phone        string    `gorm:"column:phone" db:"phone" json:"phone" form:"phone"`                                 //  用户手机号码
	RegisterTime time.Time `gorm:"column:register_time" db:"register_time" json:"register_time" form:"register_time"` //  注册时间
}
