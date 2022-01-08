package model

type User struct {
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Birthday string `json:"birthday"`
	Sex      int    `json:"sex"`
	City     int    `json:"city"`
}
