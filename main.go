package main

import (
	"TM_chat/models"
	"TM_chat/router"
	"TM_chat/utils"
	"flag"
)

var configpath = flag.String("config", "config/app.yaml", "")

func main() {
	flag.Parse()
	//sql.Createtable()s
	utils.InitConfig(*configpath)
	utils.InitMysql()
	utils.InitRedis()
	// utils.DB.AutoMigrate(models.UserBasic{})
	utils.DB.AutoMigrate(models.Contact{})
	utils.DB.AutoMigrate(models.GroupBasic{})
	utils.DB.AutoMigrate(models.Message{})
	utils.DB.AutoMigrate(models.Message{})
	r := router.Router()
	r.Run(":8083")
}
