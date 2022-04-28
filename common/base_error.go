package common

type BaseError struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}

var (
	ErrJSONUnmarshallFailed = &BaseError{-10001, "Unmarshall JSON failed"}

	/*--------------- User -------------------*/
	ErrUserLoginFailed      = &BaseError{-10002, "User Login failed"}    // 用户登录失败
	ErrUserRegisterFailed   = &BaseError{-10003, "User Register failed"} // 用户注册失败
	ErrGetUserInfoFailed    = &BaseError{-10004, "获得用户信息失败"}             // 获得用户信息失败
	ErrUpdateUserInfoFailed = &BaseError{-10005, "更新用户信息失败"}             // 更新用户信息失败

	/*--------------- Content ----------------*/
	ErrCreateContentFailed = &BaseError{-101001, "创建内容失败"}
)
