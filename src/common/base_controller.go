package common

import (
	"MyServer/src/database"
	"github.com/jinzhu/gorm"
)

type BaseController struct {
}

func (c *BaseController) GetDB() *gorm.DB {
	return database.DB
}
