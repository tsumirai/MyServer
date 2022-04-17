package common

import (
	"MyServer/database"
	"gorm.io/gorm"
)

type BaseDao struct{}

func (c *BaseDao) GetDB() *gorm.DB {
	return database.DB
}
