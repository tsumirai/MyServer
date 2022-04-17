package common

import (
	"MyServer/database"

	"gorm.io/gorm"
)

type BaseService struct{}

func (c *BaseService) GetDB() *gorm.DB {
	return database.DB
}
