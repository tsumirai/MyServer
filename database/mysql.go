package database

import (
	"MyServer/base"
	"MyServer/middleware/logger"
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitMysql() {
	var err error
	mysqlConfig := base.Config.GetString("mysql.user") + ":" + base.Config.GetString("mysql.password") +
		"@tcp(" + base.Config.GetString("mysql.ip") + ":" + base.Config.GetString("mysql.port") + ")" +
		"/" + base.Config.GetString("mysql.database") + "?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", mysqlConfig)
	if err != nil {
		logger.Error(context.TODO(), logger.LogArgs{"msg": "InitMysql failed", "err": err.Error()})
		panic(err)
	}
	DB.SingularTable(true)
	logger.Info(context.TODO(), logger.LogArgs{"msg": "InitMysql Success!"})
}
