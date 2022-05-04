package model

// UserPhone 用户手机表
type UserPhone struct {
	ID    int64  `json:"id"`
	Phone string `json:"phone""`
	UID   int64  `json:"uid"`
}
