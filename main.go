package main

import (
	"TM_chat/models"
	"TM_chat/router"
	"TM_chat/utils"
	"flag"
	"time"

	"github.com/spf13/viper"
)

var configpath = flag.String("config", "config/app.yaml", "")

func main() {
	flag.Parse()
	//sql.Createtable()s
	utils.InitConfig(*configpath)
	utils.InitMysql()
	utils.InitRedis()
	// utils.DB.AutoMigrate(models.UserBasic{})
	// utils.DB.AutoMigrate(models.Contact{})
	// utils.DB.AutoMigrate(models.GroupBasic{})
	// utils.DB.AutoMigrate(models.Message{})
	// utils.DB.AutoMigrate(models.Community{})
	// InitTimer()
	r := router.Router()
	r.Run(viper.GetString("port.server"))
}
func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second,
		time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}
