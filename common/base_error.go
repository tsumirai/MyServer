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

	ErrCreateContentFailed        = &BaseError{-101001, "创建内容失败"}
	ErrGetContentFailed           = &BaseError{-101002, "获得内容数据失败"}
	ErrSetContentPermissionFailed = &BaseError{-101003, "设置内容权限失败"}

	/*-------------- Comment -----------------*/

	ErrCreateCommentFailed = &BaseError{-102001, "创建评论失败"}
	ErrGetCommentFailed    = &BaseError{-102002, "获得评论失败"}
)
