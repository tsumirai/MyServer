package database

import (
	config "MyServer/conf"
	"MyServer/middleware/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitMysql() {
	var err error
	mysqlConfig := config.Config.GetString("mysql.user") + ":" + config.Config.GetString("mysql.password") +
		"@tcp(" + config.Config.GetString("mysql.ip") + ":" + config.Config.GetString("mysql.port") + ")" +
		"/" + config.Config.GetString("mysql.dbbase") + "?charset=utf8&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", mysqlConfig)
	if err != nil {
		logger.Error(logger.LogArgs{"msg": "InitMysql failed", "err": err.Error()})
		panic(err)
	}
	DB.SingularTable(true)
	logger.Info(logger.LogArgs{"msg": "InitMysql Success!"})
}