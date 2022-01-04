package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var Config *viper.Viper

func InitConfig() {
	Config = viper.New()
	Config.SetConfigName("config")
	Config.SetConfigType("toml")
	Config.AddConfigPath(".\\src\\config\\")
	Config.SetDefault("redis.port", 6381)
	err := Config.ReadInConfig()
	if err != nil {
		log.Fatalln("Init Failed: read config failed: ", err.Error())
		panic(err)
	}
}
