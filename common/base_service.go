package common

import (
	"MyServer/database"

	"github.com/jinzhu/gorm"
)

type BaseService struct{}

func (c *BaseService) GetDB() *gorm.DB {
	return database.DB
}
