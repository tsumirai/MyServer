package database

import (
	"MyServer/src/config"
	"fmt"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

func InitMysql(){
	var err error
	mysqlConfig := config.Config.GetString("mysql.user")+":"+config.Config.GetString("mysql.password") +
		"@tcp("+config.Config.GetString("mysql.ip")+":"+config.Config.GetString("mysql.port")+")"+
		"/"+config.Config.GetString("mysql.dbbase")+"?charset=utf8&parseTime=True&loc=Local"
	db,err = gorm.Open("mysql",mysqlConfig)
	if err != nil {
		fmt.Println("InitMysql failed: ",err.Error())
		panic(err)
	}
	db.SingularTable(true)
	log.Info("InitMysql Success!")
}