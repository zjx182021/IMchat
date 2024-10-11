package models

import (
	"TM_chat/utils"
	"fmt"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
	Desc     string
}

func (table *Contact) ContactTableName() string {
	return "Contact_basics"
}

func SearchFriend(userId uint) []*UserBasic {
	contacts := make([]*Contact, 0)
	objId := make([]uint, 0)
	utils.DB.Where("ower_id = ? and type = ?", userId, 1).Find(&contacts)
	for _, data := range contacts {
		fmt.Println(data)
		objId = append(objId, data.TargetId)
	}
	users := make([]*UserBasic, 0)
	utils.DB.Where("id IN (?)", objId).Find(&users)
	return users
}
