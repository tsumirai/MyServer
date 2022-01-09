package common

type BaseError struct {
	ErrNo  int    `json:"err_no"`
	ErrMsg string `json:"err_msg"`
}

var ErrJSONUnmarshallFailed = &BaseError{-10001, "Unmarshall JSON failed"}
