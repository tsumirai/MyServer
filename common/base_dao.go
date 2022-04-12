package common

import (
	"MyServer/database"

	"github.com/jinzhu/gorm"
)

type BaseDao struct{}

func (c *BaseDao) GetDB() *gorm.DB {
	return database.DB
}
