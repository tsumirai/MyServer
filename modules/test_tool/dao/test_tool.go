package dao

import (
	"MyServer/common"

	"github.com/gin-gonic/gin"
)

type TestToolDao struct {
	common.BaseDao
}

func NewTestToolDao() *TestToolDao {
	return &TestToolDao{}
}

func (d *TestToolDao) DivideTable(ctx *gin.Context) {
	d.GetDB()
}
