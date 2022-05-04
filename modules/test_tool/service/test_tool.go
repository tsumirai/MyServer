package service

import (
	"MyServer/common"
)

type TestToolService struct {
	common.BaseService
}

func NewTestToolService() *TestToolService {
	return &TestToolService{}
}
