package sql

import (
	"TM_chat/models"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Createtable() {
	mysqllogger := logger.Default.LogMode(logger.Warn)
	db, err := gorm.Open(mysql.Open("root:!Zhang123456@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{Logger: mysqllogger})

	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(models.UserBasic{})
	user := &models.UserBasic{}
	// user.Name = "zjx"
	// db.Create(&user)
	db.First(&user)
	fmt.Println(user)
	db.Model(user).Update("Password", "123456")
}
