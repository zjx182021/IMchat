package models

import (
	"TM_chat/utils"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name       string    `gorm:"size:10"`
	Password   string    `gorm:"size:40"`
	Phone      string    `gorm:"size:15" valid:"matches(^1[3-9]{1}[0-9]{11}$)"`
	Email      *string   `gorm:"size:128" valid:"^[1-9a-zA-Z]+@[a-zA-Z0-9]{1,13}\.([a-zA-Z0-9]{1,3}\.)[,5][a-zA-Z]{,5}$"`
	Identity   string    `gorm:"size:128"`
	ClientIp   string    `gorm:"size:15"`
	ClientPort int       `gorm:"size:10"`
	Salt       string    `gorm:"size:32"`
	LoginTime  time.Time `gorm:"size:64;column:login_time;default:0000-00-00"`
	Heartbeat  time.Time `gorm:"size:64;default:0000-00-00"`
	LogoutTime time.Time `gorm:"size:64;column:logout_time;default:0000-00-00"`
	IsLogout   bool      `gorm:"size:1"`
	DeviceInfo string    `gorm:"size:64"`
}

func (table *UserBasic) TableName() string {
	return "user_basics"
}
func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func CreateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Create(user)
}

func DeleteUser(user *UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user *UserBasic) *gorm.DB {
	return utils.DB.Model(&UserBasic{}).Where("id = ?", user.ID).Updates(&UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Phone:    user.Phone,
		Email:    user.Email,
	})
}

func FindUserByname(user *UserBasic) *gorm.DB {
	return utils.DB.Where("name =?", user.Name).First(&user)
}

func FindUserByphone(user *UserBasic) *gorm.DB {
	return utils.DB.Where("phone =?", user.Phone).First(&user)
}

func FindUserByemail(user *UserBasic) *gorm.DB {
	return utils.DB.Where("email =?", user.Email).First(&user)
}

func FindUserByid(user *UserBasic) *gorm.DB {
	return utils.DB.Where("id =?", user.ID).First(&user)
}

func Updateidentity(user *UserBasic) *gorm.DB {
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)

	return utils.DB.Model(&UserBasic{}).Where("name =? and password=?", user.Name, user.Password).Update("identity", temp)

}
