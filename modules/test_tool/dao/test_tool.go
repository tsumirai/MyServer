package dao

import (
	"MyServer/common"
)

type TestToolDao struct {
	common.BaseDao
}

func NewTestToolDao() *TestToolDao {
	return &TestToolDao{}
}
