package model

type User struct {
	Name     string `json:"name"`
	NickName string `json:"nick_name"`
	Phone    string `json:"phone"`
	PassWord string `json:"pass_word"`
	City     int    `json:"city"`
}
